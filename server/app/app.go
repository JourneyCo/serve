package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
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
