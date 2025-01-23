package models

import (
	"time"
)

type Location struct {
	Latitude  float64
	Longitude float64

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
