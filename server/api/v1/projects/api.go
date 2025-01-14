package projects

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"serve/app/database"
	"serve/app/errors"
	"serve/helpers"
	"serve/models"
)

type ApiResponse struct {
	Text string `json:"text"`
}

func toCtx(rw http.ResponseWriter, r *http.Request) {
}

func create(rw http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var project models.Project
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			http.Error(rw, "Invalid request body", http.StatusBadRequest)
			log.Printf("invalid request body: %v", err)
		}

		project, err = database.PostProject(ctx)

		if err != nil {
			http.Error(rw, "Failed to create project", http.StatusInternalServerError)
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
			errors.ServerError(rw, err)
		}
	} else {
		errors.NotFoundHandler(rw, r)
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
