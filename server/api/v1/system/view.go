package system

import (
	"net/http"

	"serve/helpers"
)

func show() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		dto := struct {
			GoogleMapKey string `json:"google_map_key"`
		}{
			GoogleMapKey: helpers.GetEnvVar("GOOGLE_KEY"),
		}

		helpers.WriteJSON(w, dto)
	})
}
