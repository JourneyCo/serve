package auth0

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/golang-jwt/jwt/v5"
	db "serve/app/database"
	"serve/helpers"
	"serve/models"
)

const (
	publicKeyFile = "public.key"
	Admin         = 200
)

var publicKeyPEM []byte

type Config struct {
	port     string
	Audience string
	Domain   string
}

type Claims struct {
	jwt.RegisteredClaims
}

type Session struct {
	UserID    string
	First     *string
	Last      *string
	Email     *string
	CellPhone *string
	TextPerm  *bool
	CreatedAt time.Time
	UpdatedAt *time.Time
}

var config Config

func New() Config {
	config = Config{
		Audience: helpers.GetEnvVar("AUTH0_AUDIENCE"),
		Domain:   helpers.GetEnvVar("AUTH0_DOMAIN"),
		port:     helpers.GetEnvVar("AUTH0_PORT"),
	}
	return config
}

type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasPermissions(expectedClaims []string) bool {
	if len(expectedClaims) == 0 {
		return false
	}
	for _, scope := range expectedClaims {
		if !helpers.Contains(c.Permissions, scope) {
			return false
		}
	}
	return true
}

func SetPublicKeyPEM() {
	var err error
	publicKeyPEM, err = os.ReadFile(publicKeyFile)
	if err != nil || len(publicKeyPEM) == 0 {
		log.Fatalf("public key file is not present") // called from main on startup, so fatal here
	}
}

func ValidatePermissions(expectedClaims []string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*CustomClaims)
			if !claims.HasPermissions(expectedClaims) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			issuerURL, err := url.Parse("https://" + config.Domain + "/")
			if err != nil {
				log.Printf("Failed to parse the issuer url: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

			jwtValidator, err := validator.New(
				provider.KeyFunc,
				validator.RS256,
				issuerURL.String(),
				[]string{config.Audience},
				validator.WithCustomClaims(
					func() validator.CustomClaims {
						return new(CustomClaims)
					},
				),
			)
			if err != nil {
				log.Print("Failed to set up the jwt validator")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if authHeaderParts := strings.Fields(r.Header.Get("Authorization")); len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
				log.Printf("invalid JWT - jwt not found")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
				log.Printf("Encountered error while validating JWT: %v", err)
				if errors.Is(err, jwtmiddleware.ErrJWTMissing) {
					log.Printf("missing JWT")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
					log.Print("invalid JWT")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
			}

			middleware := jwtmiddleware.New(
				jwtValidator.ValidateToken,
				jwtmiddleware.WithErrorHandler(errorHandler),
			)

			middleware.CheckJWT(next).ServeHTTP(w, r)
		},
	)
}

func CreateSession(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeaderParts := strings.Fields(r.Header.Get("Authorization"))

			if len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
				log.Printf("invalid JWT - jwt not found")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
			if err != nil {
				log.Printf("error parsing public key: %s", publicKeyPEM)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var claims jwt.RegisteredClaims

			token, err := jwt.ParseWithClaims(
				authHeaderParts[1], &claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return publicKey, nil
				},
			)

			if err != nil {
				log.Printf("error parsing claims: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, ok := token.Claims.(*jwt.RegisteredClaims); !ok {
				log.Print("unknown claims type, cannot proceed")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var s Session
			account, err := db.GetAccount(ctx, claims.Subject)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Println("error getting account from db: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// if the account doesn't already exist, we will add it
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				a := models.Account{
					ID: claims.Subject,
				}
				account, err = db.PostAccount(ctx, a)
				if err != nil {
					log.Println("error posting account to db: ", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				// url := "https://" + config.Domain + "/userinfo"
				// req, err := http.NewRequest("GET", url, nil)
				// if err != nil {
				// 	log.Println("error composing user info request: ", err)
				// 	w.WriteHeader(http.StatusInternalServerError)
				// 	return
				// }
				//
				// req.Header.Set("Authorization", r.Header.Get("Authorization"))
				//
				// client := &http.Client{}
				// resp, err := client.Do(req)
				// if err != nil {
				// 	log.Println("error getting userinfo from domain: ", err)
				// 	w.WriteHeader(http.StatusInternalServerError)
				// 	return
				// }
				// defer resp.Body.Close()
				//
				// body, err := io.ReadAll(resp.Body)
				// if err != nil {
				// 	log.Print("error reading response body from userinfo endpoint: ", err)
				// 	w.WriteHeader(http.StatusInternalServerError)
				// 	return
				// }
				//
				// var dto userInfo
				//
				// if err = json.Unmarshal(body, &dto); err != nil {
				// 	fmt.Printf("error unmarshalling body: %v", err)
				// 	return
				// }
				//
				// account.Email = dto.Email
				// account.FirstName = dto.GivenName
				// account.LastName = dto.FamilyName
			}

			s = Session{
				UserID:    account.ID,
				First:     account.FirstName,
				Last:      account.LastName,
				Email:     account.Email,
				CellPhone: account.CellPhone,
				TextPerm:  account.TextPermission,
				CreatedAt: account.CreatedAt,
				UpdatedAt: account.UpdatedAt,
			}
			ctx = context.WithValue(ctx, "session", s)

			h.ServeHTTP(w, r.WithContext(ctx))

		},
	)
}
