package locations

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
			locations := []models.Location{}

			locations, err := db.GetLocations(ctx)
			if err != nil {
				fmt.Printf("error retrieving locations: %v", err)
				return
			}

			ctx = context.WithValue(ctx, "locations", locations)

			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
