package errors

import (
	"log"
	"net/http"

	"serve/helpers"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

const (
	notFoundErrorMessage       = "Not Found"
	internalServerErrorMessage = "Internal Server Error"
)

func ServerError(rw http.ResponseWriter, err error) {
	errorMessage := ErrorMessage{Message: internalServerErrorMessage}
	helpers.WriteJSON(rw, errorMessage)
	log.Print("Internal error server: ", err.Error())
}

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	errorMessage := ErrorMessage{Message: notFoundErrorMessage}
	helpers.WriteJSON(rw, errorMessage)
}
