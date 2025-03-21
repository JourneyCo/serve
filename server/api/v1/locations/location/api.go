package location

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	db "serve/app/database"
)

// toCtx will place the location into the context.
func toCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			vars := mux.Vars(r)
			id, err := strconv.ParseInt(vars["id"], 10, 64)
			if err != nil {
				log.Println("error in parsing location id: ", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			l, err := db.GetLocation(ctx, int(id))
			if err != nil {
				log.Println("error getting location from db: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, "location", l)
			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
