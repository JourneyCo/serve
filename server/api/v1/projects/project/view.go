package project

import (
	"net/http"

	"serve/helpers"
	"serve/models"
)

func show() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			project := ctx.Value("project").(models.Project)

			dto := Request{
				ID:         &project.ID,
				Name:       &project.Name,
				Required:   &project.Required,
				Needed:     &project.Needed,
				LeaderID:   &project.LeaderID,
				LocationID: &project.LocationID,
				CreatedAt:  project.CreatedAt,
				UpdatedAt:  project.UpdatedAt,
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
