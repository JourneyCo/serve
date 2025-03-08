package models

import (
	"time"
)

type Project struct {
	// ID is the serial number primary key generated from postgres
	ID int64

	Enabled bool

	// Name is the name of the project
	Name string

	// Required is the total number of people needed for the project
	Required int

	// Needed is the number of people still needed for the project
	Needed int

	// LeaderID is the account ID of the project lead
	LeaderID int64

	Status string

	StartTime        time.Time
	EndTime          time.Time
	Category         string
	AgesID           *int64
	Wheelchair       bool
	ShortDescription string
	LongDescription  *string
	LocationID       int64
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}
