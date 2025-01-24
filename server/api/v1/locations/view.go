package locations

import (
	"log"
	"net/http"
	"serve/helpers"
	"serve/models"
)

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		locations := ctx.Value("locations").([]models.Location)

		dto := []models.Location{}

		for _, location := range locations {
			l := models.Location{
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
			dto = append(dto, l)
		}

		if err := helpers.WriteJSON(w, http.StatusOK, dto); err != nil {
			log.Println(err)
		}
	})
}
