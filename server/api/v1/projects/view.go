package projects

import (
	"net/http"
	"serve/helpers"
)

func Show() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, APIResponse{})
	})
}
