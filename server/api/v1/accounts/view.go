package accounts

import (
	"net/http"
	"time"

	"serve/helpers"
	"serve/models"
)

func index() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accounts := ctx.Value("accounts").([]models.Account)

			dto := []any{}

			for _, account := range accounts {
				l := struct {
					ID string `json:"id"`

					UpdatedAt *time.Time `json:"updated_at"`
				}{
					ID: account.ID,
				}
				dto = append(dto, l)
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
