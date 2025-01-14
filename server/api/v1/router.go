package v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"serve/api/v1/projects"
	"serve/app/auth0"
	"serve/app/errors"
	"serve/app/middleware"
)

func Route(domain string, audience string, r *mux.Router) {

	r.HandleFunc("/", errors.NotFoundHandler)

	r.Path("/messages/protected").
		Methods(http.MethodGet).
		Handler(projects.LogAMessage(projects.Show()))

	r.Path("/messages/admin").
		Methods(http.MethodGet).
		Handler(auth0.ValidateJWT(audience, domain, http.HandlerFunc(projects.AdminApiHandler)))

	s := r.Path("/projects").Subrouter()
	projects.Route(s)

	r.Use(middleware.HandleCacheControl)
}
