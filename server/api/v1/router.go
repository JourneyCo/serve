package v1

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Route(r *mux.Router) {
	r.HandleFunc("", ProcessCrons).Methods("GET")
}

func ProcessCrons(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Crons route"))
	w.WriteHeader(http.StatusOK)
}
