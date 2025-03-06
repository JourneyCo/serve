package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"serve/helpers"
)

var DB *sql.DB

func StartDB() {
	// Connect to Database
	fmt.Println("Connecting to database...")
	var err error
	DB, err = sql.Open("postgres", dataSource())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to get driver:", err)
	}

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
}

func dataSource() string {
	// TODO: remove hardcoding before prod
	host := helpers.GetEnvVar("DB_HOST")
	// dbUser := helpers.GetEnvVar("DB_USER")
	// dbPass := helpers.GetEnvVar("DB_PASS")
	if os.Getenv("profile") == "prod" {
		host = "database"
	}
	return "postgresql://" + host + ":5432/serve" +
		"?user=" + "postgres" + "&sslmode=disable&password=" + "pass"
}
