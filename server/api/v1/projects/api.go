package projects

import (
	"net/http"
	"serve/app/errors"
	"serve/helpers"
)

type ApiResponse struct {
	Text string `json:"text"`
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

func ProtectedApiHandler(rw http.ResponseWriter, r *http.Request) {
	sendMessage(rw, r, ProtectedMessage())
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
