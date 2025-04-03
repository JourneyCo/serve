package handlers

import (
        "encoding/json"
        "fmt"
        "net/http"
        "os"
        "strings"

        "github.com/gorilla/mux"

        "project-registration-system/config"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
        Config *config.Config
}

// RegisterAuthRoutes registers the routes for authentication handlers
func RegisterAuthRoutes(router *mux.Router, cfg *config.Config) {
        handler := &AuthHandler{
                Config: cfg,
        }

        router.HandleFunc("/config", handler.GetAuthConfig).Methods("GET")
}

// AuthConfig holds the Auth0 configuration for the frontend
type AuthConfig struct {
        Domain      string `json:"domain"`
        ClientID    string `json:"clientId"`
        Audience    string `json:"audience"`
        RedirectURI string `json:"redirectUri"`
}

// GetAuthConfig returns the Auth0 configuration for the frontend
func (h *AuthHandler) GetAuthConfig(w http.ResponseWriter, r *http.Request) {
        // Determine the redirect URI based on the request
        scheme := "https"  // Default to https for production
        proxyPort := "5000" // Port where our proxy server is running
        var redirectURI string
        
        // Get the host from the request
        host := r.Host
        
        // Check if we're in a Replit environment
        replitNix := os.Getenv("REPL_SLUG")
        replitId := os.Getenv("REPL_ID")
        isReplit := replitNix != "" || replitId != ""
        
        // Add debug logging for request headers
        fmt.Printf("Request Host: %s\n", host)
        fmt.Printf("X-Forwarded-Host: %s\n", r.Header.Get("X-Forwarded-Host"))
        fmt.Printf("X-Forwarded-Proto: %s\n", r.Header.Get("X-Forwarded-Proto"))
        
        if isReplit {
                // We're in Replit
                forwardedHost := r.Header.Get("X-Forwarded-Host")
                if forwardedHost != "" {
                        host = forwardedHost
                }
                
                // Normalize the host if it includes a port
                if strings.Contains(host, ":") {
                        parts := strings.Split(host, ":")
                        if len(parts) > 0 {
                                host = parts[0]
                        }
                }
                
                // In Replit, don't include the port in the redirect URI
                redirectURI = fmt.Sprintf("%s://%s/callback", scheme, host)
                
                // Add debug logging
                fmt.Printf("Auth0 Redirect URI (Replit): %s\n", redirectURI)
                
                // For localhost development
                if host == "localhost" || host == "127.0.0.1" {
                        scheme = "http"
                        redirectURI = fmt.Sprintf("%s://%s:%s/callback", scheme, host, proxyPort)
                        fmt.Printf("Auth0 Redirect URI (Local): %s\n", redirectURI)
                }
        } else {
                // Not in Replit, likely local development
                scheme = "http"
                
                // Normalize the host if it includes a port
                if strings.Contains(host, ":") {
                        parts := strings.Split(host, ":")
                        if len(parts) > 0 {
                                host = parts[0]
                        }
                }
                
                // Use the port where the proxy server is running
                redirectURI = fmt.Sprintf("%s://%s:%s/callback", scheme, host, proxyPort)
                fmt.Printf("Auth0 Redirect URI (Dev): %s\n", redirectURI)
        }

        config := AuthConfig{
                Domain:      h.Config.Auth0Domain,
                ClientID:    h.Config.Auth0ClientID,
                Audience:    h.Config.Auth0Audience,
                RedirectURI: redirectURI,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(config)
}


