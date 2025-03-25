package account

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	db "serve/app/database"
)

// toCtx will place the account into the context.
func toCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			vars := mux.Vars(r)
			id := vars["id"]

			l, err := db.GetAccount(ctx, id)
			if err != nil {
				log.Println("error getting account from db: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, "account", l)
			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
