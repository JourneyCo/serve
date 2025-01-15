package projects

import "time"

type Request struct {
	// ID is the serial number primary key generated from postgres
	ID int64 `json:"id"`

	// GoogleID is a string representation of the ID created within Google Sheets
	GoogleID string `json:"google_id"`
	Name     string `json:"name"`

	// Required is the total number of people needed for the project
	Required int `json:"required"`

	// Needed is the number of people still needed for the project
	Needed int `json:"needed"`

	Street       string     `json:"street"`
	StreetNumber int        `json:"street_number"`
	City         string     `json:"city"`
	State        string     `json:"state"`
	PostalCode   string     `json:"postal_code"`
	AdminID      int64      `json:"admin_id"`
	LocationID   *int64     `json:"location_id"`
	Date         *time.Time `json:"date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
