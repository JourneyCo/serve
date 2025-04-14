package v1

import (
	"serve/api/v1/accounts"
	"serve/api/v1/locations"
	"serve/api/v1/projects"
	"serve/api/v1/registrations"
	"serve/app/middleware"
	"serve/app/router"
)

func Route(r router.ServeRouter) {
	// r.Use(auth0.ValidateJWT)
	// r.Use(auth0.CreateSession)

	a := r.SubPath("/accounts")
	accounts.Route(a)

	p := r.SubPath("/projects")
	projects.Route(p)

	l := r.SubPath("/locations")
	locations.Route(l)

	reg := r.SubPath("/registrations")
	registrations.Route(reg)

	r.Use(middleware.HandleCacheControl)
}
