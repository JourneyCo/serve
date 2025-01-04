package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	v1 "serve/api/v1"
	"serve/web"
)

func main() {
	// Getting environment variables
	fmt.Println("Getting Env variables...")
	//dbHost := cmp.Or(os.Getenv("DB_HOST"), "localhost")
	//dbPort := cmp.Or(os.Getenv("DB_PORT"), "5432")
	//dbUser := cmp.Or(os.Getenv("DB_USER"), "journey")
	//dbPassword := cmp.Or(os.Getenv("DB_PASSWORD"), "pass")
	//dbName := cmp.Or(os.Getenv("DB_NAME"), "journey")

	// Connect to Database
	fmt.Println("Connecting to database...")
	_, err := gorm.Open(postgres.Open(dataSource()), &gorm.Config{}) // grab db

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
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
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=prefer TimeZone=UTC", host, dbUser, dbPass, "journey", 5432)
	//return "postgresql://" + host + ":5432/journey" +
	//	"?user=" + dbUser + "&sslmode=disable&password=" + dbPass
}
