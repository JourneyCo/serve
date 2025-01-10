package models

import (
	"time"
)

type Projects struct {
	// ID is the serial number primary key generated from postgres
	ID int64

	// GoogleID is a string representation of the ID created within Google Sheets
	GoogleID string
	Name     string

	// Required is the total number of people needed for the project
	Required int

	// Needed is the number of people still needed for the project
	Needed     int
	AdminID    int64
	LocationID int64
	Date       *time.Time
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
