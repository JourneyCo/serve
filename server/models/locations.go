package models

import (
	"time"
)

type Location struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Address   string
	Latitude  float64
	Longitude float64
	Info      string
}
