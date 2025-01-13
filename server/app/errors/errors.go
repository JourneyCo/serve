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
	werr := helpers.WriteJSON(rw, http.StatusInternalServerError, errorMessage)
	if werr != nil {
		log.Println("Error writing error message: ", werr.Error())
	}
	log.Print("Internal error server: ", err.Error())
}

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	errorMessage := ErrorMessage{Message: notFoundErrorMessage}
	err := helpers.WriteJSON(rw, http.StatusNotFound, errorMessage)
	if err != nil {
		ServerError(rw, err)
	}
}
