package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration for the application
type Config struct {
	DevMode bool
	// Date
	ServeDay string

	// Server config
	ServerPort string

	// Database config
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBConnectionOptions string
	DBURL               string

	// Auth0 config
	Auth0Domain       string
	Auth0Audience     string
	Auth0ClientID     string
	Auth0ClientSecret string

	// Email config
	MailHost string
	MailKey  string
	MailFrom string
	MailUser string
	MailPass string
	MailPort string

	// Text config
	ClearStreamAPIKey string
	TextFrom          string

	// Google Maps API config
	GoogleMapsAPIKey string

	// Recaptcha config
	RecaptchaProject string
	RecaptchaKey     string
	RecaptchaAction  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		DevMode: getEnv("DEV_MODE", "true") == "true",
		// Serve Day Date
		ServeDay: getEnv("SERVE_DAY", "07-12-25"),

		// Server config with default
		ServerPort: getEnv("PORT", "8080"),

		// Database config
		DBHost:              getEnv("PGHOST", "localhost"),
		DBPort:              getEnv("PGPORT", "5432"),
		DBUser:              getEnv("PGUSER", "postgres"),
		DBPassword:          getEnv("PGPASSWORD", "postgres"),
		DBName:              getEnv("PGDATABASE", "serve"),
		DBConnectionOptions: getEnv("DBCONN_OPTS", ""),
		DBURL:               getEnv("DATABASE_URL", ""),

		// Auth0 config - in dev mode use placeholders
		Auth0Domain:       getEnv("AUTH0_DOMAIN", "dev-placeholder.auth0.com"),
		Auth0Audience:     getEnv("AUTH0_AUDIENCE", "https://api.projectregistration.com"),
		Auth0ClientID:     getEnv("AUTH0_CLIENT_ID", "dev-client-id"),
		Auth0ClientSecret: getEnv("AUTH0_CLIENT_SECRET", "dev-client-secret"),

		// Email config - in dev mode use placeholders
		MailHost: getEnv("MAIL_HOST", "smtp.example.com"),
		MailKey:  getEnv("MAIL_KEY", "apikey"),
		MailFrom: getEnv("MAIL_FROM", "from@example.com"),
		MailPort: getEnv("MAIL_PORT", "587"),
		MailUser: getEnv("MAIL_USER", "user"),
		MailPass: getEnv("MAIL_PASS", "password"),

		// Text config
		ClearStreamAPIKey: getEnv("CS_API_KEY", "apikey"),
		TextFrom:          getEnv("CS_TEXT_FROM", "9007"),

		// Google Maps API config
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),

		// Google Maps API config
		RecaptchaKey:     getEnv("RECAPTCHA_KEY", ""),
		RecaptchaProject: getEnv("RECAPTCHA_PROJECT", ""),
		RecaptchaAction:  getEnv("RECAPTCHA_ACTION", ""),
	}

	// In production mode, validate required configuration
	if !config.DevMode {
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
		if getEnv("MAIL_HOST", "") == "" {
			missingVars = append(missingVars, "MAIL_HOST")
		}
		if getEnv("MAIL_USER", "") == "" {
			missingVars = append(missingVars, "MAIL_USER")
		}
		if getEnv("MAIL_PASS", "") == "" {
			missingVars = append(missingVars, "MAIL_PASS")
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
		"host=%s port=%s user=%s password=%s dbname=%s %s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBConnectionOptions,
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
