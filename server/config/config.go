package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	// Server config
	ServerPort string

	// Database config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBURL      string

	// Auth0 config
	Auth0Domain       string
	Auth0Audience     string
	Auth0ClientID     string
	Auth0ClientSecret string

	// Email config
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	EmailFrom    string

	// Text config
	ClearStreamAPIKey string
	TextFrom          string

	// Google Maps API config
	GoogleMapsAPIKey string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Check for development mode
	devMode := getEnv("DEV_MODE", "true") == "true"

	config := &Config{
		// Server config with default
		ServerPort: getEnv("PORT", "8080"),

		// Database config
		DBHost:     getEnv("PGHOST", "localhost"),
		DBPort:     getEnv("PGPORT", "5432"),
		DBUser:     getEnv("PGUSER", "postgres"),
		DBPassword: getEnv("PGPASSWORD", "postgres"),
		DBName:     getEnv("PGDATABASE", "serve"),
		DBURL:      getEnv("DATABASE_URL", ""),

		// Auth0 config - in dev mode use placeholders
		Auth0Domain:       getEnv("AUTH0_DOMAIN", "dev-placeholder.auth0.com"),
		Auth0Audience:     getEnv("AUTH0_AUDIENCE", "https://api.projectregistration.com"),
		Auth0ClientID:     getEnv("AUTH0_CLIENT_ID", "dev-client-id"),
		Auth0ClientSecret: getEnv("AUTH0_CLIENT_SECRET", "dev-client-secret"),

		// Email config - in dev mode use placeholders
		SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", "dev@example.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "dev-password"),
		EmailFrom:    getEnv("EMAIL_FROM", "noreply@projectregistration.com"),

		// Text config
		ClearStreamAPIKey: getEnv("CS_API_KEY", "apikey"),
		TextFrom:          getEnv("CS_TEXT_FROM", "9007"),

		// Google Maps API config
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),
	}

	// In production mode, validate required configuration
	if !devMode {
		var missingVars []string

		// For Auth0
		if getEnv("AUTH0_DOMAIN", "") == "" {
			missingVars = append(missingVars, "AUTH0_DOMAIN")
		}
		if getEnv("AUTH0_AUDIENCE", "") == "" {
			missingVars = append(missingVars, "AUTH0_AUDIENCE")
		}
		if getEnv("AUTH0_CLIENT_ID", "") == "" {
			missingVars = append(missingVars, "AUTH0_CLIENT_ID")
		}
		if getEnv("AUTH0_CLIENT_SECRET", "") == "" {
			missingVars = append(missingVars, "AUTH0_CLIENT_SECRET")
		}

		// For Email
		if getEnv("SMTP_HOST", "") == "" {
			missingVars = append(missingVars, "SMTP_HOST")
		}
		if getEnv("SMTP_USERNAME", "") == "" {
			missingVars = append(missingVars, "SMTP_USERNAME")
		}
		if getEnv("SMTP_PASSWORD", "") == "" {
			missingVars = append(missingVars, "SMTP_PASSWORD")
		}

		// For Google Maps API
		if getEnv("GOOGLE_MAPS_API_KEY", "") == "" {
			missingVars = append(missingVars, "GOOGLE_MAPS_API_KEY")
		}

		// If any required variables are missing, return an error
		if len(missingVars) > 0 {
			return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
		}
	}

	return config, nil
}

// GetDBConnString returns the database connection string
func (c *Config) GetDBConnString() string {
	if c.DBURL != "" {
		return c.DBURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
