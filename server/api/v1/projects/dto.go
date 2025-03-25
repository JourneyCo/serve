package projects

import "time"

type Request struct {
	// ID is the serial number primary key generated from postgres
	ID int64 `json:"id"`

	// Name is the project name
	Name string `json:"name"`

	// Required is the total number of people needed for the project
	Required int `json:"required"`

	// Registered is the number of people registered for the project
	Registered int `json:"registered"`

	Enabled          bool       `json:"enabled"`
	Status           string     `json:"status"`
	LeaderID         string     `json:"leader_id"`
	StartTime        time.Time  `json:"start_time"`
	EndTime          time.Time  `json:"end_time"`
	Category         string     `json:"category"`
	AgesID           *int64     `json:"ages_id,omitempty"`
	Wheelchair       bool       `json:"wheelchair"`
	ShortDescription string     `json:"short_description"`
	LongDescription  *string    `json:"long_description,omitempty"`
	Street           string     `json:"street,omitempty"`
	StreetNumber     int        `json:"street_number"`
	City             string     `json:"city,omitempty"`
	State            string     `json:"state,omitempty"`
	PostalCode       string     `json:"postal_code"`
	LocationID       *int64     `json:"location_id,omitempty"`
	Date             *time.Time `json:"date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}
