package accounts

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	db "serve/app/database"
	"serve/models"
)

func idxToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accounts := []models.Account{}
			proj := r.URL.Query().Get("project")
			var err error
			if proj == "" {
				accounts, err = db.GetAccounts(ctx)
				if err != nil {
					fmt.Printf("error retrieving accounts: %v", err)
					return
				}
			} else {
				parseInt, err := strconv.Atoi(proj)
				if err != nil {
					fmt.Print("error parsing int")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				accounts, err = db.GetAccountsByProject(ctx, parseInt)
				if err != nil {
					fmt.Printf("error retrieving accounts by proj: %v", err)
					return
				}
			}

			ctx = context.WithValue(ctx, "accounts", accounts)
			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
