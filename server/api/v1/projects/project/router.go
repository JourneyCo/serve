package project

import (
	"github.com/gorilla/mux"
	"net/http"
	"serve/app/middleware"
)

func Route(r *mux.Router) {

	u := r.Path("/register").Subrouter()

	r.Path("").
		Methods(http.MethodGet).
		Handler(show())

	// user register
	u.Path("").
		Methods(http.MethodPut).
		Handler(show())

	// administer project
	r.Path("").
		Methods(http.MethodPut).
		Handler(show())

	r.Path("").
		Methods(http.MethodDelete).
		Handler(show())

	r.Use(middleware.HandleCacheControl)
}
