package dto

import (
	"time"

	"serve/models"
)

type Lead struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Active bool   `json:"active"`
}

// ProjectAccessory represents an accessory to the project
type ProjectAccessory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Project represents a project in the system
type Project struct {
	ID              int                       `json:"id"`
	GoogleID        *int                      `json:"google_id"`
	Title           string                    `json:"title"`
	Description     string                    `json:"description"`
	Website         string                    `json:"website"`
	Time            string                    `json:"time"`
	ProjectDate     time.Time                 `json:"project_date"`
	MaxCapacity     int                       `json:"max_capacity"`
	CurrentReg      int                       `json:"current_registrations"`
	Area            string                    `json:"area"`
	LocationAddress string                    `json:"location_address"`
	Latitude        float64                   `json:"latitude"`
	Longitude       float64                   `json:"longitude"`
	ServeLeadID     string                    `json:"serve_lead_id"`
	ServeLeadName   string                    `json:"serve_lead_name"`
	ServeLeadEmail  string                    `json:"serve_lead_email"`
	ServeLead       *models.User              `json:"serve_lead,omitempty"`
	Types           []models.ProjectAccessory `json:"types,omitempty"`
	Ages            string                    `json:"ages,omitempty"`
	Leads           []Lead                    `json:"leads,omitempty"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	Status          string                    `json:"status"`
}
