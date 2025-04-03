package database

import (
        "database/sql"
        "fmt"
        _ "github.com/lib/pq" // PostgreSQL driver
        "project-registration-system/config"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*sql.DB, error) {
        // Open database connection
        db, err := sql.Open("postgres", cfg.GetDBConnString())
        if err != nil {
                return nil, fmt.Errorf("failed to open database connection: %w", err)
        }

        // Test database connection
        if err = db.Ping(); err != nil {
                return nil, fmt.Errorf("failed to ping database: %w", err)
        }

        return db, nil
}
