package projects

import (
	"log"
	"net/http"
	"serve/helpers"
	"serve/models"
)

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projects := ctx.Value("projects").([]models.Project)

		dto := []Request{}

		for _, project := range projects {
			p := Request{
				ID:         project.ID,
				GoogleID:   project.GoogleID,
				Name:       project.Name,
				Required:   project.Required,
				Needed:     project.Needed,
				AdminID:    project.AdminID,
				LocationID: &project.LocationID,
				Date:       project.Date,
				CreatedAt:  project.CreatedAt,
				UpdatedAt:  project.UpdatedAt,
			}
			dto = append(dto, p)
		}

		if err := helpers.WriteJSON(w, http.StatusOK, dto); err != nil {
			log.Println(err)
		}
	})
}

func show() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		project := ctx.Value("project").(models.Project)

		dto := Request{
			ID:         project.ID,
			GoogleID:   project.GoogleID,
			Name:       project.Name,
			Required:   project.Required,
			Needed:     project.Needed,
			AdminID:    project.AdminID,
			LocationID: &project.LocationID,
			Date:       project.Date,
			CreatedAt:  project.CreatedAt,
			UpdatedAt:  project.UpdatedAt,
		}

		if err := helpers.WriteJSON(w, http.StatusOK, dto); err != nil {
			log.Println(err)
		}
	})
}
