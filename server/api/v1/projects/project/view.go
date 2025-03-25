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
				ID:               &project.ID,
				Enabled:          &project.Enabled,
				Name:             &project.Name,
				Required:         &project.Required,
				Status:           &project.Status,
				StartTime:        &project.StartTime,
				EndTime:          &project.EndTime,
				Category:         &project.Category,
				AgesID:           project.AgesID,
				Wheelchair:       &project.Wheelchair,
				ShortDescription: &project.ShortDescription,
				LongDescription:  project.LongDescription,
				Registered:       &project.Registered,
				LeaderID:         &project.LeaderID,
				LocationID:       &project.LocationID,
				CreatedAt:        project.CreatedAt,
				UpdatedAt:        project.UpdatedAt,
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
