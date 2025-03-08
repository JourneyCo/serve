package models

import (
	"time"
)

type Account struct {
	ID             int64
	FirstName      string
	LastName       string
	Email          string
	CellPhone      string
	TextPermission bool
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}
