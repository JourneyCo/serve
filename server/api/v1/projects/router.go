package projects

import (
	"net/http"

	"github.com/gorilla/mux"
	"serve/api/v1/projects/project"
	"serve/app/middleware"
)

func Route(r *mux.Router) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(idxToCtx(index()))

	r.Path("").
		Methods(http.MethodPost).
		Handler(middleware.JSONToCtx(Request{}, create(show())))

	// single project
	p := r.Path("/{id:[0-9]+}").Subrouter()
	project.Route(p)

	r.Use(middleware.HandleCacheControl)
}
