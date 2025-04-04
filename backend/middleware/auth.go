package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"

	"project-registration-system/config"
)

// CustomClaims contains custom claims extended from the standard JWT claims
type CustomClaims struct {
	Email    string   `json:"email"`
	Nickname string   `json:"nickname"`
	Name     string   `json:"name"`
	Picture  string   `json:"picture"`
	Roles    []string `json:"https://projectapp/roles"`
}

// Validate does custom validation for the token
func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// AuthMiddleware returns a middleware function that validates JWT tokens
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	issuerURL := fmt.Sprintf("https://%s/", cfg.Auth0Domain)
	audience := cfg.Auth0Audience

	issuer, err := url.Parse(issuerURL)
	if err != nil {
		log.Fatalf("Failed to parse the issuer URL: %v", err)
	}

	provider := jwks.NewCachingProvider(issuer, 5*60)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL,
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(30),
	)
	if err != nil {
		log.Fatalf("Failed to set up JWT validator: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				encounteredError := true
				var handler http.Handler = http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						encounteredError = false
						next.ServeHTTP(w, r)
					},
				)

				middleware.CheckJWT(handler).ServeHTTP(w, r)

				if encounteredError {
					// No need to handle the error as the middleware already did
					return
				}
			},
		)
	}
}

// AdminMiddleware checks if the user has admin role
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(jwtmiddleware.ContextKey{})
			if token == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.(*validator.ValidatedClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			customClaims, ok := claims.CustomClaims.(*CustomClaims)
			if !ok {
				http.Error(w, "Invalid custom claims", http.StatusUnauthorized)
				return
			}

			// Check if user has admin role
			isAdmin := false
			for _, role := range customClaims.Roles {
				if role == "admin" {
					isAdmin = true
					break
				}
			}

			if !isAdmin {
				http.Error(w, "Forbidden: admin role required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

// GetUserIDFromRequest extracts the user ID from the JWT token
func GetUserIDFromRequest(r *http.Request) (string, error) {
	token := r.Context().Value(jwtmiddleware.ContextKey{})
	if token == nil {
		return "", errors.New("no token found in request context")
	}

	claims, ok := token.(*validator.ValidatedClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Get the subject (user ID) from the standard claims
	return claims.RegisteredClaims.Subject, nil
}

// GetUserFromRequest extracts user information from the JWT token
func GetUserFromRequest(r *http.Request) (*CustomClaims, error) {
	token := r.Context().Value(jwtmiddleware.ContextKey{})
	if token == nil {
		return nil, errors.New("no token found in request context")
	}

	claims, ok := token.(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid custom claims")
	}

	return customClaims, nil
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, status int, message string) {
	response := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	if status != http.StatusOK {
		log.Print(message)
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// ParseJSON parses a JSON request body into the provided struct
func ParseJSON(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}
