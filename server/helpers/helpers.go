package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func GetEnvVar(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("the environment variable '%s' doesn't exist or is not set", key)
	}
	return os.Getenv(key)
}

func WriteJSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		log.Println("error marshaling json: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(js)
	if err != nil {
		log.Println("error writing json: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
