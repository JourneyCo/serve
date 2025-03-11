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

	a := app.New()

	if err := a.Serve(); err != nil {
		log.Fatal("error serving application: ", err)
	}
}
