package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	v1 "serve/api/v1"
	"serve/web"
)

func main() {
	if os.Getenv("profile") != "prod" {
		godotenv.Load() // load env's from .env file
	}
	// Getting environment variables
	fmt.Println("Getting Env variables...")
	//dbHost := cmp.Or(os.Getenv("DB_HOST"), "localhost")
	//dbPort := cmp.Or(os.Getenv("DB_PORT"), "5432")
	//dbUser := cmp.Or(os.Getenv("DB_USER"), "journey")
	//dbPassword := cmp.Or(os.Getenv("DB_PASSWORD"), "pass")
	//dbName := cmp.Or(os.Getenv("DB_NAME"), "journey")

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
		log.Fatal("error migrating up: ", err)
	}

	r := mux.NewRouter()
	v := r.PathPrefix("/api/v1").Subrouter()
	v1.Route(v)

	// CORS is enabled only in prod profile
	cors := os.Getenv("profile") == "prod"
	app := web.NewApp(cors)

	if err = app.Serve(); err != nil {
		log.Println("error serving application: ", err)
	}
}

func dataSource() string {
	//TODO: remove hardocding before prod
	host := "localhost"
	dbUser := "journey"
	dbPass := "pass"
	if os.Getenv("profile") == "prod" {
		host = "db"
		dbPass = os.Getenv("db_pass")
	}
	return "postgresql://" + host + ":5432/journey" +
		"?user=" + dbUser + "&sslmode=disable&password=" + dbPass
}
