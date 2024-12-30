package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"os"
	v1 "serve/api/v1"
	"serve/db"
	"serve/web"
)

func main() {
	d, err := sql.Open("postgres", dataSource())
	if err != nil {
		log.Fatal("error starting db: ", err)
	}
	defer d.Close()

	r := mux.NewRouter()
	v := r.PathPrefix("/api/v1").Subrouter()
	v1.Route(v)

	// CORS is enabled only in prod profile
	cors := os.Getenv("profile") == "prod"
	app := web.NewApp(db.NewDB(d), cors)

	if err = app.Serve(); err != nil {
		log.Println("error serving application: ", err)
	}
}

func dataSource() string {
	//TODO: remove hardocding before prod
	host := "localhost"
	dbUser := "goxygen"
	dbPass := "pass"
	if os.Getenv("profile") == "prod" {
		host = "db"
		dbPass = os.Getenv("db_pass")
	}
	return "postgresql://" + host + ":5432/goxygen" +
		"?user=" + dbUser + "&sslmode=disable&password=" + dbPass
}
