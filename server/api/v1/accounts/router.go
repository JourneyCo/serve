package accounts

import (
	"net/http"

	"serve/api/v1/accounts/account"
	"serve/app/router"
)

func Route(r router.ServeRouter) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(idxToCtx(index()))

	// single account
	l := r.SubPath("/{id}")
	account.Route(l)
}
