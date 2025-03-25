package accounts

import (
	"context"
	"fmt"
	"net/http"

	db "serve/app/database"
	"serve/models"
)

func idxToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accounts := []models.Account{}

			accounts, err := db.GetAccounts(ctx)
			if err != nil {
				fmt.Printf("error retrieving accounts: %v", err)
				return
			}

			ctx = context.WithValue(ctx, "accounts", accounts)

			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
