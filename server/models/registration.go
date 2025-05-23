package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Registration represents a user's registration for a project
type Registration struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	ProjectID    int       `json:"project_id"`
	Status       string    `json:"status"` // "registered", "cancelled", "completed"
	GuestCount   int       `json:"guest_count"`
	LeadInterest bool      `json:"lead_interest"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         *User     `json:"user,omitempty"`
	Project      *Project  `json:"project,omitempty"`
}

// RegisterForProject registers a user for a project
func RegisterForProject(db *sql.DB, userID string, projectID int, guestCount int, isLeadInterested bool) (
	*Registration, error,
) {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if project exists and has capacity
	var currentCount int
	var maxCapacity int
	err = tx.QueryRow(
		`
								SELECT 
												COALESCE(COUNT(r.id), 0) + COALESCE(SUM(r.guest_count), 0), 
												p.max_capacity 
								FROM projects p
								LEFT JOIN registrations r ON p.id = r.project_id AND r.status = 'registered'
								WHERE p.id = $1
								GROUP BY p.id
				`, projectID,
	).Scan(&currentCount, &maxCapacity)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	// Calculate total spots needed (user + guests)
	totalSpots := 1 + guestCount

	// Check if there's enough capacity
	if currentCount+totalSpots > maxCapacity {
		return nil, errors.New("Capacity not available for total # of volunteers requested")
	}

	// Check if user is already registered for this project
	var existingID int
	err = tx.QueryRow(
		`
								SELECT id FROM registrations
								WHERE user_id = $1 AND project_id = $2 AND status = 'registered'
				`, userID, projectID,
	).Scan(&existingID)

	if err == nil {
		return nil, errors.New("user is already registered for this project")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Create registration
	reg := &Registration{
		UserID:       userID,
		ProjectID:    projectID,
		Status:       "registered",
		GuestCount:   guestCount,
		LeadInterest: isLeadInterested,
	}

	err = tx.QueryRow(
		`
								INSERT INTO registrations (user_id, project_id, status, guest_count, lead_interest)
								VALUES ($1, $2, $3, $4, $5)
								RETURNING id, created_at, updated_at
				`, reg.UserID, reg.ProjectID, reg.Status, reg.GuestCount, reg.LeadInterest,
	).Scan(&reg.ID, &reg.CreatedAt, &reg.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return reg, nil
}

// CancelRegistration deletes a registration
func CancelRegistration(db *sql.DB, userID string, projectID int) error {
	query := `
								DELETE FROM registrations 
								WHERE user_id = $1 AND project_id = $2 AND status = 'registered'
				`

	result, err := db.Exec(query, userID, projectID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no active registration found for this project")
	}

	return nil
}

// GetUserRegistration gets the registration for a user
func GetUserRegistration(ctx context.Context, db *sql.DB, userID string) (Registration, error) {
	r := Registration{}
	query := `
									SELECT r.id, r.user_id, r.project_id, r.status, r.guest_count, r.lead_interest,
									r.created_at, r.updated_at,
									p.title, p.description, p.time, p.project_date, p.max_capacity,
									p.area, p.latitude, p.longitude
									FROM registrations r
									JOIN projects p ON r.project_id = p.id
									WHERE r.user_id = $1
									ORDER BY p.project_date
					`

	err := db.QueryRowContext(ctx, query, userID).Scan(
		&r.ID, &r.UserID, &r.ProjectID, &r.Status, &r.GuestCount, &r.LeadInterest,
		&r.CreatedAt, &r.UpdatedAt, &r.Project.Title, &r.Project.Description, &r.Project.Time, &r.Project.ProjectDate,
		&r.Project.MaxCapacity, &r.Project.Area, &r.Project.Latitude, &r.Project.Longitude,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r, nil // not found
		}
		return r, err
	}

	return r, nil
}

// GetProjectRegistrations gets all registrations for a project
func GetProjectRegistrations(ctx context.Context, db *sql.DB, projectID int) ([]Registration, error) {
	query := `
									SELECT r.id, r.user_id, r.project_id, r.status, r.guest_count, r.lead_interest,
									r.created_at, r.updated_at,
									u.email, u.first_name, u.last_name, u.phone, u.text_permission
									FROM registrations r
									JOIN users u ON r.user_id = u.id
									WHERE r.project_id = $1
									ORDER BY r.status, r.created_at
					`

	rows, err := db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []Registration
	for rows.Next() {
		var r Registration
		r.User = &User{}

		if err := rows.Scan(
			&r.ID, &r.UserID, &r.ProjectID, &r.Status, &r.GuestCount, &r.LeadInterest,
			&r.CreatedAt, &r.UpdatedAt,
			&r.User.Email, &r.User.FirstName, &r.User.LastName, &r.User.Phone, &r.User.TextPermission,
		); err != nil {
			return nil, err
		}

		r.User.ID = r.UserID
		registrations = append(registrations, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

// GetRegistrationsForReminders gets registrations for projects starting in specified days
func GetRegistrationsForReminders(db *sql.DB, days int) ([]Registration, error) {
	query := `
									SELECT r.id, r.user_id, r.project_id, r.status, r.guest_count, r.lead_interest,
									r.created_at, r.updated_at,
									u.email, u.first_name, u.last_name,
									p.title, p.description, p.time, p.project_date,
									p.area, p.latitude, p.longitude
									FROM registrations r
									JOIN users u ON r.user_id = u.id
									JOIN projects p ON r.project_id = p.id
									WHERE r.status = 'registered' 
									AND p.project_date = CURRENT_DATE + $1::integer
									ORDER BY p.project_date
					`

	rows, err := db.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var registrations []Registration
	for rows.Next() {
		var r Registration
		r.User = &User{}
		r.Project = &Project{}

		if err := rows.Scan(
			&r.ID, &r.UserID, &r.ProjectID, &r.Status, &r.GuestCount, &r.LeadInterest,
			&r.CreatedAt, &r.UpdatedAt,
			&r.User.Email, &r.User.FirstName, &r.User.LastName,
			&r.Project.Title, &r.Project.Description, &r.Project.Time, &r.Project.ProjectDate,
			&r.Project.Area, &r.Project.Latitude, &r.Project.Longitude,
		); err != nil {
			return nil, err
		}

		r.User.ID = r.UserID
		r.Project.ID = r.ProjectID
		registrations = append(registrations, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

// GetUserRegistrationByEmail gets the first active registration for a given email
func GetUserRegistrationByEmail(ctx context.Context, db *sql.DB, email string) (int, error) {
	query := `
		SELECT r.project_id
		FROM registrations r
		JOIN users u ON r.user_id = u.id
		WHERE u.email = $1
		AND r.status = 'registered'
	`

	var projectID int
	err := db.QueryRowContext(ctx, query, email).Scan(&projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return projectID, nil
}
