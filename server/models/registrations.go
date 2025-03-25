package models

import "time"

type Registration struct {
	AccountID   string
	ProjectID   int64
	QtyEnrolled int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
