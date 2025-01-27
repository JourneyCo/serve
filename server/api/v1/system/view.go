package system

import (
	"log"
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

		if err := helpers.WriteJSON(w, http.StatusOK, dto); err != nil {
			log.Println(err)
		}
	})
}
