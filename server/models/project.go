package models

import (
	"database/sql"
	"time"
)

// Project represents a project in the system
type Project struct {
	ID                   int       `json:"id"`
	Title                string    `json:"title"`
	ShortDescription     string    `json:"short_description"`
	Description          string    `json:"description"`
	Time                 string    `json:"time"`
	ProjectDate          time.Time `json:"project_date"`
	MaxCapacity          int       `json:"max_capacity"`
	CurrentReg           int       `json:"current_registrations"`
	LocationName         string    `json:"location_name"`
	LocationAddress      string    `json:"location_address"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
	WheelchairAccessible bool      `json:"wheelchair_accessible"`
	LeadUserID           string    `json:"lead_user_id"`
	LeadUser             *User     `json:"lead_user,omitempty"`
	Tools                []Tool    `json:"tools,omitempty"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

// Tool represents a tool or skill in the system
type Tool struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetAllProjects retrieves all projects from the database
func GetAllProjects(db *sql.DB) ([]Project, error) {
	query := `
                SELECT p.id, p.title, p.short_description, p.description, p.time, p.project_date, 
                p.max_capacity, p.location_name, p.location_address, p.latitude, p.longitude,
                p.wheelchair_accessible, p.created_at, p.updated_at, 
                COALESCE(SUM(CASE WHEN r.status = 'registered' THEN r.guest_count + 1 ELSE 0 END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                GROUP BY p.id
                ORDER BY p.project_date ASC
        `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.Title, &p.ShortDescription, &p.Description, &p.Time, &p.ProjectDate,
			&p.MaxCapacity, &p.LocationName, &p.LocationAddress, &p.Latitude, &p.Longitude,
			&p.WheelchairAccessible, &p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
		); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProjectByID retrieves a project by its ID
func GetProjectByID(db *sql.DB, id int) (*Project, error) {
	query := `
                SELECT p.id, p.title, p.description, p.short_description, p.time, p.project_date, 
                p.max_capacity, p.location_name, p.location_address, p.latitude, p.longitude, p.lead_user_id,
                p.wheelchair_accessible, p.created_at, p.updated_at, 
                COALESCE(SUM(CASE WHEN r.status = 'registered' THEN r.guest_count + 1 ELSE 0 END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                WHERE p.id = $1
                GROUP BY p.id
        `

	var p Project
	err := db.QueryRow(query, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.ShortDescription, &p.Time, &p.ProjectDate,
		&p.MaxCapacity, &p.LocationName, &p.LocationAddress, &p.Latitude, &p.Longitude, &p.LeadUserID,
		&p.WheelchairAccessible, &p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Project not found
		}
		return nil, err
	}

	return &p, nil
}

// CreateProject creates a new project in the database
func CreateProject(db *sql.DB, project *Project) error {
	query := `
                INSERT INTO projects (title, short_description, description, time, project_date, max_capacity, 
                                    location_name, location_address, latitude, longitude, wheelchair_accessible)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
                RETURNING id, created_at, updated_at
        `

	return db.QueryRow(
		query,
		project.Title,
		project.ShortDescription,
		project.Description,
		project.Time,
		project.ProjectDate,
		project.MaxCapacity,
		project.LocationName,
		project.LocationAddress,
		project.Latitude,
		project.Longitude,
		project.WheelchairAccessible,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

// UpdateProject updates an existing project
func UpdateProject(db *sql.DB, project *Project) error {
	query := `
                UPDATE projects
                SET title = $1, short_description = $2, description = $3, time = $4, project_date = $5, 
                max_capacity = $6, location_name = $7, location_address = $8, latitude = $9, longitude = $10,
                wheelchair_accessible = $11, updated_at = CURRENT_TIMESTAMP
                WHERE id = $12
                RETURNING updated_at
        `

	return db.QueryRow(
		query,
		project.Title,
		project.ShortDescription,
		project.Description,
		project.Time,
		project.ProjectDate,
		project.MaxCapacity,
		project.LocationName,
		project.LocationAddress,
		project.Latitude,
		project.Longitude,
		project.WheelchairAccessible,
		project.ID,
	).Scan(&project.UpdatedAt)
}

// DeleteProject deletes a project by its ID
func DeleteProject(db *sql.DB, id int) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// GetUpcomingProjects retrieves projects that are starting within the given days
func GetUpcomingProjects(db *sql.DB, days int) ([]Project, error) {
	query := `
                SELECT p.id, p.title, p.description, p.time, p.project_date, 
                p.max_capacity, p.location_name, p.location_address, p.latitude, p.longitude,
                p.created_at, p.updated_at, 
                COALESCE(COUNT(CASE WHEN r.status = 'registered' THEN r.id END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                WHERE p.project_date BETWEEN CURRENT_DATE AND CURRENT_DATE + $1::integer
                GROUP BY p.id
        `

	rows, err := db.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.Time, &p.ProjectDate,
			&p.MaxCapacity, &p.LocationName, &p.LocationAddress, &p.Latitude, &p.Longitude,
			&p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
		); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
