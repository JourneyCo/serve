package project

import (
	"net/http"

	"serve/app/auth0"
	"serve/app/middleware"
	"serve/app/router"
)

func Route(r router.ServeRouter) {
	r.Use(toCtx)

	r.Path("").
		Methods(http.MethodGet).
		Handler(show())

	// administer edit a project
	r.RBAC(auth0.Admin).
		Path("").
		Methods(http.MethodPut).
		Handler(show())

	r.RBAC(auth0.Admin).
		Path("").
		Methods(http.MethodDelete).
		Handler(show())

	u := r.Path("/register").Subrouter()
	u.Use(auth0.ValidateToken)

	// user register under a project (limited edit)
	u.Path("").
		Methods(http.MethodPut).
		Handler(middleware.JSONToCtx(Request{}, register(show())))

	r.Use(middleware.HandleCacheControl)
}
