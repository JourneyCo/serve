package testutils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http/httptest"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // File for migrations use
	"github.com/gorilla/mux"
	"github.com/integralist/go-findroot/find"
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
	err := os.Setenv("DEV_MODE", "true")
	if err != nil {
		log.Fatal("Failed to load dev_mode env var:", err)
	}

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

	connStr := "host=0.0.0.0 port=5432 user=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Check if the database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		// Create the database (must use string interpolation, not parameters)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database created:", dbName)
	} else {
		fmt.Println("Database already exists:", dbName)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to get driver:", err)
	}

	root, err := find.Repo()
	if err != nil {
		// handle error
	}

	// Run database migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+root.Path+"/server/migrations",
		"postgres", driver,
	)
	if err != nil || m == nil {
		log.Fatal("failed to create migration instance:", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("error migrating up: ", err) // only fatal if there is a migration up where there is a change
	}

	return db
}
