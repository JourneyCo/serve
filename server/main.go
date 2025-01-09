package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"serve/app"
	"serve/helpers"
)

func main() {
	godotenv.Load() // load env's from .env file

	// Getting environment variables
	fmt.Println("Getting Env variables...")

	// Connect to Database
	fmt.Println("Connecting to database...")
	db, err := sql.Open("postgres", dataSource())
	defer db.Close()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to get driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:./migrations",
		"postgres", driver)
	if err != nil || m == nil {
		log.Fatal("failed to create migration instance:", err)
	}

	if err = m.Up(); err != nil {
		log.Print("error migrating up: ", err)
	}

	// CORS is enabled only in prod profile
	//cors := helpers.GetEnvVar("profile") == "prod"
	a := app.New()

	if err = a.Serve(); err != nil {
		log.Fatal("error serving application: ", err)
	}
}

func dataSource() string {
	//TODO: remove hardcoding before prod
	host := helpers.GetEnvVar("DB_HOST")
	dbUser := helpers.GetEnvVar("DB_USER")
	dbPass := helpers.GetEnvVar("DB_PASS")
	if os.Getenv("profile") == "prod" {
		host = "db"
	}
	return "postgresql://" + host + ":5432/journey" +
		"?user=" + dbUser + "&sslmode=disable&password=" + dbPass
}
