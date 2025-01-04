package models

import (
	"time"
)

type Locations struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt *time.Time
	Address   string
	Lattitude float64
	Longitude float64
	Info      string
}
