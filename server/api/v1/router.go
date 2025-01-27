package v1

import (
	"net/http"
	"serve/api/v1/locations"
	"serve/api/v1/system"

	"github.com/gorilla/mux"
	"serve/api/v1/projects"
	"serve/app/auth0"
	"serve/app/errors"
	"serve/app/middleware"
)

func Route(domain string, audience string, r *mux.Router) {

	r.HandleFunc("/", errors.NotFoundHandler)

	r.Path("/messages/admin").
		Methods(http.MethodGet).
		Handler(auth0.ValidateJWT(audience, domain, http.HandlerFunc(projects.AdminApiHandler)))

	s := r.Path("/projects").Subrouter()
	projects.Route(s)

	l := r.Path("/locations").Subrouter()
	locations.Route(l)

	sys := r.Path("/system").Subrouter()
	system.Route(sys)

	r.Use(middleware.HandleCacheControl)
}
