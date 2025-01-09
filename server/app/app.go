package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"serve/app/auth0"
	"serve/app/router"
)

type App struct {
	handlers    map[string]http.HandlerFunc
	Auth0Config auth0.Config
	Router      *mux.Router
}

func New() App {
	app := App{
		handlers:    make(map[string]http.HandlerFunc),
		Auth0Config: auth0.New(),
	}

	app.Router = router.Route(app.Auth0Config.Domain, app.Auth0Config.Audience)

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
