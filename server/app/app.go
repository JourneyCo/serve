package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/unrolled/secure"
	"serve/api/v1"
	"serve/app/auth0"
	db "serve/app/database"
	"serve/app/router"
)

type App struct {
	Auth0Config auth0.Config
	Router      http.Handler
	Database    *sql.DB
}

func New() App {
	app := App{
		Auth0Config: auth0.New(),
		Database:    db.DB,
	}

	r := router.ServeRouter{Router: mux.NewRouter()}
	// cors := os.Getenv("profile") == "prod"
	// if !cors {
	// Use handlers.CORS to configure CORS settings

	secureMiddleware := secure.New()
	r.Use(secureMiddleware.Handler)

	api := r.SubPath("/api/v1")
	v1.Route(api)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://yourdomain.com", "http://localhost:3000"}),   // Allowed origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),           // Allowed methods
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), // Allowed headers
		handlers.AllowCredentials(), // Allow credentials
	)(r)

	app.Router = corsHandler
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
