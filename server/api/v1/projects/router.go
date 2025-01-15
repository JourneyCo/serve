package projects

import (
	"github.com/gorilla/mux"
	"net/http"
	"serve/api/v1/projects/project"
	"serve/app/middleware"
)

func Route(r *mux.Router) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(LogAMessage(Show()))

	r.Path("").
		Methods(http.MethodPost).
		Handler(middleware.JSONToCtx(Request{}, create(Show())))

	// single project
	p := r.Path("/{id:[0-9]+}").Subrouter()
	project.Route(p)

	r.Use(middleware.HandleCacheControl)
}
