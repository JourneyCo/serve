package registrations

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
			registrations := []models.Registration{}

			registrations, err := db.GetRegistrations(ctx)
			if err != nil {
				fmt.Printf("error retrieving registrations: %v", err)
				return
			}

			ctx = context.WithValue(ctx, "registrations", registrations)
			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
