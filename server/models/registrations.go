package models

import "time"

type Registration struct {
	AccountID   int64
	ProjectID   int64
	QtyEnrolled int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
