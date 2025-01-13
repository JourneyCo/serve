package projects

import (
	"github.com/gorilla/mux"
	"serve/app/middleware"
)

func Route(r *mux.Router) {

	r.Path("").
		Methods("GET").
		Handler(ProtectedApiHandler(Show()))

	r.Use(middleware.HandleCacheControl)
}
