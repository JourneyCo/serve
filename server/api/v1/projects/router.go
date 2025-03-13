package projects

import (
	"net/http"

	"serve/api/v1/projects/project"
	"serve/app/auth0"
	"serve/app/middleware"
	"serve/app/router"
)

func Route(r router.ServeRouter) {
	r.Path("").
		Methods(http.MethodGet).
		Handler(idxToCtx(index()))

	r.RBAC(auth0.Admin).
		Path("").
		Methods(http.MethodPost).
		Handler(middleware.JSONToCtx(Request{}, create(show())))

	// single project
	l := r.SubPath("/{id:[0-9]+}")
	project.Route(l)

	r.Use(middleware.HandleCacheControl)
}
