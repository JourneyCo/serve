package project

import (
	"net/http"
	"serve/helpers"
)

func show() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{})
	})
}
