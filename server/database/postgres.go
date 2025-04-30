package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // File for migrations use
	_ "github.com/lib/pq"                                // PostgresSQL driver
	"serve/config"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*sql.DB, error) {

	if !isDockerRunning() {
		log.Fatal("Docker is not running; unable to initiate app")
	}
	// Open database connection
	db, err := sql.Open("postgres", cfg.GetDBConnString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to get driver:", err)
	}

	// Run database migrations
	m, err := migrate.NewWithDatabaseInstance(
		"file:./migrations",
		"postgres", driver,
	)
	if err != nil || m == nil {
		log.Fatal("failed to create migration instance:", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("error migrating up: ", err) // only fatal if there is a migration up where there is a change
	}

	return db, nil
}

func isDockerRunning() bool {
	_, err := exec.Command("docker", "info").Output()
	if err != nil {
		return false
	}
	return true
}
