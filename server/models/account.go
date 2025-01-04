package models

import (
	"time"
)

type Account struct {
	ID        int64      `json:"id"`
	FirstName string     `json:"first"`
	LastName  string     `json:"last"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
