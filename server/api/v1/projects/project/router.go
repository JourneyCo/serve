package project

import (
	"net/http"

	"github.com/gorilla/mux"
	"serve/app/middleware"
)

func Route(r *mux.Router) {
	r.Use(toCtx)

	r.Path("").
		Methods(http.MethodGet).
		Handler(show())

	// administer edit a project
	r.Path("").
		Methods(http.MethodPut).
		Handler(show())

	r.Path("").
		Methods(http.MethodDelete).
		Handler(show())

	u := r.Path("/register").Subrouter()

	// user register under a project (limited edit)
	u.Path("").
		Methods(http.MethodPut).
		Handler(middleware.JSONToCtx(Request{}, register(show())))

	r.Use(middleware.HandleCacheControl)
}
