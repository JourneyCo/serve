package app

import (
	"database/sql"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"os"
	"serve/app/middleware"
	"time"

	"github.com/gorilla/mux"
	"serve/api/v1"
	"serve/app/auth0"
	db "serve/app/database"
)

type App struct {
	Auth0Config auth0.Config
	Router      *mux.Router
	Database    *sql.DB
}

func New() App {
	app := App{
		Auth0Config: auth0.New(),
		Database:    db.DB,
	}

	r := mux.NewRouter()
	secureMiddleware := secure.New()
	r.Use(secureMiddleware.Handler)
	cors := os.Getenv("profile") == "prod"
	if !cors {
		r.Use(middleware.DisableCORS)
	}

	s := r.PathPrefix("/api/v1").Subrouter()
	v1.Route(app.Auth0Config.Domain, app.Auth0Config.Audience, s)

	app.Router = r
	return app
}

func (a *App) Serve() error {
	port := ":8080"
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("web server is available on port %s\n", port)
	return srv.ListenAndServe()
}
