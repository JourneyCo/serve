package v1

import (
	"serve/api/v1/locations"
	"serve/api/v1/projects"
	"serve/api/v1/registrations"
	"serve/api/v1/system"
	"serve/app/auth0"
	"serve/app/middleware"
	"serve/app/router"
)

func Route(r router.ServeRouter) {
	r.Use(auth0.ValidateJWT)
	r.Use(auth0.CreateSession)

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
