package models

import (
	"time"
)

type Accounts struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
