package models

import (
	"time"
)

type Project struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Required  int        `json:"required"`
	Needed    int        `json:"needed"`
	AdminID   int64      `json:"admin_id"`
	Date      *time.Time `json:"date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
