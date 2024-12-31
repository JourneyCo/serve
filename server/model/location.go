package model

import "time"

type Location struct {
	ID        int64      `json:"id"`
	Address   string     `json:"address"`
	Lattitude float64    `json:"lattitude"`
	Longitude float64    `json:"longitude"`
	Info      string     `json:"info"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
