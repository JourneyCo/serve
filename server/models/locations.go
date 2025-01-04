package models

import (
	"time"
)

type Locations struct {
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Address   string
	Lattitude float64
	Longitude float64
	Info      string
}
