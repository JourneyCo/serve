
package testutils

import (
	"database/sql"
	"log"
	"net/http/httptest"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"serve/config"
	"serve/handlers"
	"serve/services"
)

// TestServer represents a test server instance
type TestServer struct {
	Server *httptest.Server
	DB     *sql.DB
	Router *mux.Router
}

// NewTestServer creates a new test server with a test database
func NewTestServer() *TestServer {
	// Set test environment
	os.Setenv("DEV_MODE", "true")
	
	// Initialize test config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup test database connection
	db := setupTestDB()

	// Initialize router and handlers
	router := mux.NewRouter()
	
	// Initialize services
	emailService := services.NewEmailService(cfg)
	
	// Register routes
	api := router.PathPrefix("/api").Subrouter()
	projectRouter := api.PathPrefix("/projects").Subrouter()
	handlers.RegisterProjectRoutes(projectRouter, db, emailService)

	// Create test server
	ts := httptest.NewServer(router)

	return &TestServer{
		Server: ts,
		DB:     db,
		Router: router,
	}
}

// Close closes the test server and database connections
func (ts *TestServer) Close() {
	ts.Server.Close()
	ts.DB.Close()
}

// setupTestDB initializes a test database connection
func setupTestDB() *sql.DB {
	// Use environment variables or default test database
	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "serve_test"
	}

	connStr := "host=0.0.0.0 port=5432 user=postgres password=postgres dbname=" + dbName + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	return db
}
