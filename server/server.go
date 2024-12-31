package main

import (
	"database/sql"
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
	d, err := sql.Open("postgres", dataSource())
	if err != nil {
		log.Fatal("error starting db: ", err)
	}
	defer d.Close()

	_, err = gorm.Open(postgres.New(postgres.Config{
		Conn: d,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("gorm could not connect to db: ", err)
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
