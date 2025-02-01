package projects

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"googlemaps.github.io/maps"
	ldb "serve/app/database"
	lerrors "serve/app/errors"
	"serve/helpers"
	"serve/models"
	"time"
)

type APIResponse struct {
	Text string `json:"text"`
}

func idxToCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projects := []models.Project{}

		projects, err := ldb.GetProjects(ctx)
		if err != nil {
			fmt.Printf("error retrieving projects: %v", err)
			return
		}

		ctx = context.WithValue(ctx, "projects", projects)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// create will create a new project, along with a new location if it
// does not already exist in the system.
func create(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body := ctx.Value("body").([]byte)
		var dto Request

		if err := json.Unmarshal(body, &dto); err != nil {
			fmt.Printf("error unmarshalling body: %v", err)
			return
		}

		//TODO: Handle existing locations

		l, err, status := createLocation(ctx, dto)
		if err != nil {
			log.Printf("error creating location: %v", err)
			w.WriteHeader(status)
			return
		}

		now := time.Now()
		project := models.Project{ //TODO: Remove hardcode once we can get dto into context
			GoogleID:   dto.GoogleID,
			Name:       dto.Name,
			Required:   dto.Required,
			Needed:     dto.Needed,
			AdminID:    1,
			LocationID: l.ID,
			Date:       &now,
			CreatedAt:  now,
			UpdatedAt:  &now,
		}

		project, err = ldb.PostProject(ctx, project)
		if err != nil {
			log.Printf("failed to post project: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "project", project)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func sendMessage(rw http.ResponseWriter, r *http.Request, data APIResponse) {
	if r.Method == http.MethodGet {
		helpers.WriteJSON(rw, data)
	} else {
		lerrors.NotFoundHandler(rw, r)
	}
}

func AdminAPIHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, AdminMessage())
}

func AdminMessage() APIResponse {
	return APIResponse{
		Text: "This is an admin message.",
	}
}

// getExistingLocation will search the db to see if a location is already existing.
func getExistingLocation(ctx context.Context) (models.Location, error) {
	dto, ok := ctx.Value("dto").(*Request)
	if !ok {
		fmt.Println(reflect.TypeOf(dto))
		fmt.Println(dto)
	}

	a, err := ldb.GetLocationByAddress(ctx, dto.StreetNumber, dto.Street)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		fmt.Println("no rows found")
		return a, nil
	}
	if err != nil {
		fmt.Printf("error looking up address in db: %v\n", err)
		return a, err
	}
	return a, nil
}

// createLocation will create a location in the database.
func createLocation(ctx context.Context, dto Request) (models.Location, error, int) {

	c, err := maps.NewClient(maps.WithAPIKey(helpers.GetEnvVar("GOOGLE_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	q := fmt.Sprintf("%d %s %s %s %s USA", dto.StreetNumber, dto.Street, dto.State, dto.City, dto.PostalCode)
	r := &maps.TextSearchRequest{
		Query: q,
	}

	result, err := c.TextSearch(ctx, r)
	if err != nil {
		log.Printf("error obtaining lat long: %v", err)
		return models.Location{}, err, http.StatusInternalServerError
	}

	now := time.Time{}
	loc := models.Location{
		Latitude:         result.Results[0].Geometry.Location.Lat,
		Longitude:        result.Results[0].Geometry.Location.Lng,
		Info:             "",
		Street:           dto.Street,
		Number:           dto.StreetNumber,
		City:             dto.City,
		State:            dto.State,
		PostalCode:       dto.PostalCode,
		FormattedAddress: result.Results[0].FormattedAddress,
		CreatedAt:        now,
		UpdatedAt:        &now,
	}

	location, err := ldb.PostLocation(ctx, loc)
	if err != nil {
		return models.Location{}, err, http.StatusInternalServerError
	}

	return location, nil, http.StatusOK
}
