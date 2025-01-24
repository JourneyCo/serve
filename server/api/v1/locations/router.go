package locations

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Route(r *mux.Router) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(idxToCtx(index()))
}
