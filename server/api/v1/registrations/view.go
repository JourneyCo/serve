package registrations

import (
	"net/http"
	"time"

	"serve/helpers"
	"serve/models"
)

type req struct {
	AccountID   string     `json:"account_id"`
	ProjectID   int64      `json:"project_id"`
	QtyEnrolled int        `json:"qty_enroll"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func index() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			registrations := ctx.Value("registrations").([]models.Registration)
			dto := []req{}

			for _, r := range registrations {
				registration := req{
					AccountID:   r.AccountID,
					ProjectID:   r.ProjectID,
					QtyEnrolled: r.QtyEnrolled,
					UpdatedAt:   r.UpdatedAt,
				}
				dto = append(dto, registration)
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
