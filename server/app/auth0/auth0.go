package auth0

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"serve/helpers"
)

type Config struct {
	port     string
	Audience string
	Domain   string
}

func New() Config {
	return Config{
		Audience: helpers.GetEnvVar("AUTH0_AUDIENCE"),
		Domain:   helpers.GetEnvVar("AUTH0_DOMAIN"),
		port:     helpers.GetEnvVar("AUTH0_PORT"),
	}
}

const (
	Admin = 200
)

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

func ValidateJWT(audience, domain string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			issuerURL, err := url.Parse("https://" + domain + "/")
			if err != nil {
				log.Fatalf("Failed to parse the issuer url: %v", err)
			}

			provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

			jwtValidator, err := validator.New(
				provider.KeyFunc,
				validator.RS256,
				issuerURL.String(),
				[]string{audience},
				validator.WithCustomClaims(
					func() validator.CustomClaims {
						return new(CustomClaims)
					},
				),
			)
			if err != nil {
				log.Fatalf("Failed to set up the jwt validator")
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
