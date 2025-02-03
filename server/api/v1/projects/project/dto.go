package project

import "time"

type Request struct {
	// ID is the serial number primary key generated from postgres
	ID *int64 `json:"id"`

	// GoogleID is a string representation of the ID created within Google Sheets
	GoogleID *string `json:"google_id"`
	Name     *string `json:"name"`

	// Required is the total number of people needed for the project
	Required *int `json:"required"`

	// Needed is the number of people still needed for the project
	Needed *int `json:"needed"`

	// Registering is the number of people that are registering for a project
	Registering *int `json:"registering,omitempty"`

	UserID *int64 `json:"user_id,omitempty"`

	Street       *string    `json:"street,omitempty"`
	StreetNumber *int       `json:"street_number,omitempty"`
	City         *string    `json:"city,omitempty"`
	State        *string    `json:"state,omitempty"`
	PostalCode   *string    `json:"postal_code,omitempty"`
	AdminID      *int64     `json:"admin_id"`
	LocationID   *int64     `json:"location_id,omitempty"`
	Date         *time.Time `json:"date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}
