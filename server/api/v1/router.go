package v1

import (
	"net/http"

	"serve/api/v1/locations"
	"serve/api/v1/projects"
	"serve/api/v1/registrations"
	"serve/api/v1/system"
	"serve/app/auth0"
	"serve/app/middleware"
	"serve/app/router"
)

func Route(domain string, audience string, r router.ServeRouter) {

	r.RBAC(auth0.Admin).
		Path("/messages/admin").
		Methods(http.MethodGet).
		Handler(auth0.ValidateJWT(audience, domain, projects.SendMessage()))

	p := r.SubPath("/projects")
	projects.Route(p)

	l := r.PathPrefix("/locations").Subrouter()
	locations.Route(l)

	reg := r.SubPath("/registrations")
	registrations.Route(reg)

	sys := r.PathPrefix("/system").Subrouter()
	system.Route(sys)

	r.Use(middleware.HandleCacheControl)
}
