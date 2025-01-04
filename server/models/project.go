package models

import (
	"time"
)

type Project struct {
	ID        int64
	Name      string
	Required  int
	Needed    int
	AdminID   int64
	Date      *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}
