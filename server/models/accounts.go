package models

import (
	"time"
)

type Account struct {
	// ID is the user ID that auth0 provides. This ID is mapped to the auth0 token under the 'sub' claim.
	ID             string
	FirstName      *string
	LastName       *string
	Email          *string
	CellPhone      *string
	TextPermission *bool
	Lead           *bool
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}
