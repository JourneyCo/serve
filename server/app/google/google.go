package google

import (
	"github.com/kelvins/geocoder"
	"serve/helpers"
)

func SetKey() {
	key := helpers.GetEnvVar("GOOGLE_KEY")
	geocoder.ApiKey = key
}
