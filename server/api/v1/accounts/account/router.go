package account

import (
	"net/http"

	"serve/app/router"
)

func Route(r router.ServeRouter) {
	r.Use(toCtx)

	r.Path("").
		Methods(http.MethodGet).
		Handler(show())
}
