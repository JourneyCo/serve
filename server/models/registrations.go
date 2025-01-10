package models

import "time"

type Registrations struct {
	AccountID   int64
	ProjectID   int64
	QtyEnrolled int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
