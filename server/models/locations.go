package models

import (
	"github.com/kelvins/geocoder"
	"time"
)

type Location struct {
	// google's lat and long
	geocoder.Location

	ID               int64
	Info             string
	Street           string
	Number           int
	City             string
	State            string
	PostalCode       string
	FormattedAddress string
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}
