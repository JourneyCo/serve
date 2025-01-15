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

	"serve/app/database"
	lerrors "serve/app/errors"
	"serve/helpers"
	"serve/models"
	"time"
)

type ApiResponse struct {
	Text string `json:"text"`
}

func toCtx(rw http.ResponseWriter, r *http.Request) {
}

func LogAMessage(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received:", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

// create will create a new project, along with a new location if it
// does not already exist in the system
func create(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		//dto, ok := ctx.Value("dto").(*Request)
		//if !ok {
		//	fmt.Println(reflect.TypeOf(dto))
		//}
		body := ctx.Value("body").([]byte)

		var dto Request
		if err := json.Unmarshal(body, &dto); err != nil {
			fmt.Printf("error unmarshalling body: %v", err)
			return
		}

		// check to see if this is a new location or not
		if dto.LocationID == nil {
			a, err := getExistingLocation(ctx)
			if err != nil {
				fmt.Printf("error getting existing location: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if a.Number == 0 {
				// do something
			}
		} else {

		}

		now := time.Now()
		project := models.Project{ //TODO: Remove hardcode once we can get dto into context
			GoogleID:   "gold",
			Name:       "something",
			Required:   123,
			Needed:     13,
			AdminID:    1,
			LocationID: 1,
			Date:       &now,
			CreatedAt:  now,
			UpdatedAt:  &now,
		}

		project, err := database.PostProject(ctx, project)
		if err != nil {
			log.Printf("failed to post project: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "project", project)

		r = r.WithContext(ctx)
		return
	})
}

// Wrapper function
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received:", r.Method, r.URL.Path)
		next(w, r)
		fmt.Println("Request completed")
	}
}

func sendMessage(rw http.ResponseWriter, r *http.Request, data ApiResponse) {
	if r.Method == http.MethodGet {
		err := helpers.WriteJSON(rw, http.StatusOK, data)
		if err != nil {
			lerrors.ServerError(rw, err)
		}
	} else {
		lerrors.NotFoundHandler(rw, r)
	}
}

func PublicApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, PublicMessage())
}

func AdminApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, AdminMessage())
}

func PublicMessage() ApiResponse {
	return ApiResponse{
		Text: "This is a public message.",
	}
}

func ProtectedMessage() ApiResponse {
	return ApiResponse{
		Text: "This is a protected message.",
	}
}

func AdminMessage() ApiResponse {
	return ApiResponse{
		Text: "This is an admin message.",
	}
}

// getExistingLocation will search the db to see if a location is already existing
func getExistingLocation(ctx context.Context) (models.Location, error) {
	dto, ok := ctx.Value("dto").(*Request)
	if !ok {
		fmt.Println(reflect.TypeOf(dto))
		fmt.Println(dto)
	}

	a, err := database.GetLocationByAddress(ctx, dto.StreetNumber, dto.Street)
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
