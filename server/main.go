package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"serve/app"
	"serve/app/database"
	"serve/app/google"
)

func main() {
	// load env's from .env file
	godotenv.Load()

	database.StartDB()
	google.SetKey()

	if err := google.FetchProjects(); err != nil {
		log.Print(err) // we will just print the error - no need to fatal/panic
	}

	// CORS is enabled only in prod profile
	//cors := helpers.GetEnvVar("profile") == "prod"
	a := app.New()

	if err := a.Serve(); err != nil {
		log.Fatal("error serving application: ", err)
	}
}
