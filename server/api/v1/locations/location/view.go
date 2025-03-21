package location

import (
	"net/http"
	"time"

	"serve/helpers"
	"serve/models"
)

func show() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			location := ctx.Value("location").(models.Location)

			dto := struct {
				Latitude         float64    `json:"latitude"`
				Longitude        float64    `json:"longitude"`
				ID               int64      `json:"id"`
				Info             string     `json:"info"`
				Street           string     `json:"street"`
				Number           int        `json:"number"`
				City             string     `json:"city"`
				State            string     `json:"state"`
				PostalCode       string     `json:"postal_code"`
				FormattedAddress string     `json:"formatted_address"`
				CreatedAt        time.Time  `json:"created_at"`
				UpdatedAt        *time.Time `json:"updated_at"`
			}{
				Latitude:         location.Latitude,
				Longitude:        location.Longitude,
				ID:               location.ID,
				Info:             location.Info,
				Street:           location.Street,
				Number:           location.Number,
				City:             location.City,
				State:            location.State,
				PostalCode:       location.PostalCode,
				FormattedAddress: location.FormattedAddress,
				CreatedAt:        location.CreatedAt,
				UpdatedAt:        location.UpdatedAt,
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
