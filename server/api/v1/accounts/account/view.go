package account

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
			account := ctx.Value("account").(models.Account)

			dto := struct {
				ID             string     `json:"id"`
				FirstName      *string    `json:"first"`
				LastName       *string    `json:"last"`
				Email          *string    `json:"email"`
				CellPhone      *string    `json:"cellphone"`
				TextPermission *bool      `json:"text_permission"`
				Lead           *bool      `json:"lead"`
				CreatedAt      time.Time  `json:"created_at"`
				UpdatedAt      *time.Time `json:"updated_at"`
			}{
				ID:             account.ID,
				FirstName:      account.FirstName,
				LastName:       account.LastName,
				Email:          account.Email,
				CellPhone:      account.CellPhone,
				TextPermission: account.TextPermission,
				CreatedAt:      account.CreatedAt,
				UpdatedAt:      account.UpdatedAt,
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
