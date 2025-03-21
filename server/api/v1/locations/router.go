package locations

import (
	"net/http"

	"serve/api/v1/locations/location"
	"serve/app/router"
)

func Route(r router.ServeRouter) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(idxToCtx(index()))

	// single location
	l := r.SubPath("/{id:[0-9]+}")
	location.Route(l)
}
