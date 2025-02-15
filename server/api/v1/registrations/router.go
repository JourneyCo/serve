package registrations

import (
	"net/http"

	"serve/app/middleware"
	"serve/app/router"
)

func Route(r router.ServeRouter) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(idxToCtx(index()))

	r.Use(middleware.HandleCacheControl)
}
