package system

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Route(r *mux.Router) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(show())
}
