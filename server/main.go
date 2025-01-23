package main

import (
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"serve/app"
	"serve/app/database"
	"serve/app/google"
)

func main() {
	godotenv.Load() // load env's from .env file

	// Getting environment variables
	fmt.Println("Getting Env variables...")

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
