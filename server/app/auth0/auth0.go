package auth0

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"serve/helpers"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	lerrors "serve/app/errors"
)

type Config struct {
	port     string
	Audience string
	Domain   string
}

func New() Config {
	//clientOriginUrl := helpers.GetEnvVar("CLIENT_ORIGIN_URL")
	return Config{
		Audience: helpers.GetEnvVar("AUTH0_AUDIENCE"),
		Domain:   helpers.GetEnvVar("AUTH0_DOMAIN"),
		port:     helpers.GetEnvVar("AUTH0_PORT"),
	}
}

const (
	missingJWTErrorMessage = "Requires authentication"
	invalidJWTErrorMessage = "Bad credentials"
)

func ValidateJWT(audience, domain string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		)
		if err != nil {
			log.Fatalf("Failed to set up the jwt validator")
		}

		if authHeaderParts := strings.Fields(r.Header.Get("Authorization")); len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
			errorMessage := lerrors.ErrorMessage{Message: invalidJWTErrorMessage}
			if err = helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Encountered error while validating JWT: %v", err)
			if errors.Is(err, jwtmiddleware.ErrJWTMissing) {
				errorMessage := lerrors.ErrorMessage{Message: missingJWTErrorMessage}
				if err := helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
					log.Printf("Failed to write error message: %v", err)
				}
				return
			}
			if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
				errorMessage := lerrors.ErrorMessage{Message: invalidJWTErrorMessage}
				if err := helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
					log.Printf("Failed to write error message: %v", err)
				}
				return
			}
			lerrors.ServerError(w, err)
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)

		middleware.CheckJWT(next).ServeHTTP(w, r)
	})
}
