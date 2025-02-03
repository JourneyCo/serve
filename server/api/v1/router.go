package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"serve/api/v1/locations"
	"serve/api/v1/projects"
	"serve/api/v1/system"
	"serve/app/auth0"
	"serve/app/middleware"
)

func Route(domain string, audience string, r *mux.Router) {

	r.Path("/messages/admin").
		Methods(http.MethodGet).
		Handler(auth0.ValidateJWT(audience, domain, http.HandlerFunc(projects.AdminAPIHandler)))

	s := r.PathPrefix("/projects").Subrouter()
	projects.Route(s)

	l := r.PathPrefix("/locations").Subrouter()
	locations.Route(l)

	sys := r.PathPrefix("/system").Subrouter()
	system.Route(sys)

	r.Use(middleware.HandleCacheControl)
}
