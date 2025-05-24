package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Project represents a project in the system
type Project struct {
	ID               int                `json:"id"`
	GoogleID         *int               `json:"google_id"`
	Title            string             `json:"title"`
	ShortDescription string             `json:"short_description"`
	Description      string             `json:"description"`
	Website          string             `json:"website"`
	Time             string             `json:"time"`
	ProjectDate      time.Time          `json:"project_date"`
	MaxCapacity      int                `json:"max_capacity"`
	CurrentReg       int                `json:"current_registrations"`
	Area             string             `json:"area"`
	LocationAddress  string             `json:"location_address"`
	Latitude         float64            `json:"latitude"`
	Longitude        float64            `json:"longitude"`
	ServeLeadID      string             `json:"serve_lead_id"`
	ServeLead        *User              `json:"serve_lead,omitempty"`
	Categories       []ProjectAccessory `json:"categories,omitempty"`
	Ages             []ProjectAccessory `json:"ages,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

const (
	AccCategories = "categories"
	AccAges       = "ages"
)

// ProjectAccessory represents an accessory to the project
type ProjectAccessory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetAllProjects retrieves all projects from the database
func GetAllProjects(ctx context.Context, db *sql.DB) ([]Project, error) {
	query := `
                SELECT p.id, p.google_id, p.title, p.short_description, p.description, p.website, p.time, 
                p.max_capacity, p.area, p.location_address, p.latitude, p.longitude,
                p.created_at, p.updated_at, 
                COALESCE(SUM(CASE WHEN r.status = 'registered' THEN r.guest_count + 1 ELSE 0 END), 0) as current_registrations,
                array_to_string(COALESCE(array_agg(DISTINCT pc.category_id) FILTER (WHERE pc.category_id IS NOT NULL), ARRAY[]::integer[]), ',') as category_ids
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id 
                LEFT JOIN project_categories pc ON p.id = pc.project_id
                GROUP BY p.id
        `

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var categoryIDsStr string
		if err = rows.Scan(
			&p.ID, &p.GoogleID, &p.Title, &p.ShortDescription, &p.Description, &p.Website, &p.Time,
			&p.MaxCapacity, &p.Area, &p.LocationAddress, &p.Latitude, &p.Longitude,
			&p.CreatedAt, &p.UpdatedAt, &p.CurrentReg, &categoryIDsStr,
		); err != nil {
			return nil, err
		}

		// Convert comma-separated string to Categories
		if categoryIDsStr != "" {
			for _, idStr := range strings.Split(categoryIDsStr, ",") {
				if id, err := strconv.Atoi(idStr); err == nil {
					p.Categories = append(p.Categories, ProjectAccessory{ID: id})
				}
			}
		}

		// TODO: Remove hardcoding
		p.ProjectDate = time.Date(2025, 7, 12, 0, 0, 0, 0, time.UTC)
		projects = append(projects, p)
	}

	return projects, nil
}

// GetProjectByID retrieves a project by its ID
func GetProjectByID(ctx context.Context, db *sql.DB, id int) (*Project, error) {
	query := `
                SELECT p.id, p.title, p.description, p.short_description, p.website, p.time, p.project_date, 
                p.max_capacity, p.area, p.location_address, p.latitude, p.longitude, p.serve_lead_id,
                p.created_at, p.updated_at, 
                COALESCE(SUM(CASE WHEN r.status = 'registered' THEN r.guest_count + 1 ELSE 0 END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                WHERE p.id = $1
                GROUP BY p.id
        `

	var p Project
	err := db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.ShortDescription, &p.Website, &p.Time, &p.ProjectDate,
		&p.MaxCapacity, &p.Area, &p.LocationAddress, &p.Latitude, &p.Longitude, &p.ServeLeadID,
		&p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Project not found
		}
		return nil, err
	}

	// Get categories for this project
	categoriesQuery := `
		SELECT c.id, c.category FROM categories c
		JOIN project_categories pc ON c.id = pc.category_id
		WHERE pc.project_id = $1
	`
	categoryRows, err := db.QueryContext(ctx, categoriesQuery, id)
	if err != nil {
		return nil, err
	}
	defer categoryRows.Close()

	for categoryRows.Next() {
		var category ProjectAccessory
		if err := categoryRows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		p.Categories = append(p.Categories, category)
	}

	// Get ages for this project
	agesQuery := `
		SELECT a.id, a.name FROM ages a
		JOIN project_ages pa ON a.id = pa.ages_id
		WHERE pa.project_id = $1
	`
	ageRows, err := db.QueryContext(ctx, agesQuery, id)
	if err != nil {
		return nil, err
	}
	defer ageRows.Close()

	for ageRows.Next() {
		var age ProjectAccessory
		if err := ageRows.Scan(&age.ID, &age.Name); err != nil {
			return nil, err
		}
		p.Ages = append(p.Ages, age)
	}

	return &p, nil
}

// CreateProject creates a new project in the database
func CreateProject(ctx context.Context, db *sql.DB, project *Project) error {

	query := `
                INSERT INTO projects (google_id, title, short_description, description, website, time, project_date, max_capacity, 
                                    area, location_address, latitude, longitude, serve_lead_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
                RETURNING id, created_at, updated_at
        `

	// TODO: Need rollbacks here
	err := db.QueryRowContext(
		ctx,
		query,
		project.GoogleID,
		project.Title,
		project.ShortDescription,
		project.Description,
		project.Website,
		project.Time,
		project.ProjectDate,
		project.MaxCapacity,
		project.Area,
		project.LocationAddress,
		project.Latitude,
		project.Longitude,
		project.ServeLeadID,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		log.Println("error creating project: ", err)
		return err
	}

	if err = insertAccessories(db, project); err != nil {
		log.Println("error creating project associations: ", err)
		return err
	}
	return nil
}

// UpdateProject updates an existing project
func UpdateProject(ctx context.Context, db *sql.DB, project *Project) error {
	query := `
                UPDATE projects
                SET google_id=$14, title = $1, short_description = $2, description = $3, website = $4, time = $5, project_date = $6, 
                max_capacity = $7, area = $8, location_address = $9, latitude = $10, longitude = $11,
                updated_at = CURRENT_TIMESTAMP
                WHERE id = $12
                RETURNING updated_at
        `

	// TODO: Need to apply accessories here as well
	// TODO: Need rollbacks here
	return db.QueryRowContext(
		ctx,
		query,
		project.Title,
		project.ShortDescription,
		project.Description,
		project.Website,
		project.Time,
		project.ProjectDate,
		project.MaxCapacity,
		project.Area,
		project.LocationAddress,
		project.Latitude,
		project.Longitude,
		project.ID,
		project.GoogleID,
	).Scan(&project.UpdatedAt)
}

// DeleteProject deletes a project by its ID
func DeleteProject(ctx context.Context, db *sql.DB, id int) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func insertAccessories(db *sql.DB, p *Project) error {
	accs := []string{}
	var stmt string
	var valueArgs []any
	if len(p.Categories) > 0 {
		accs = append(accs, "categories")
	}
	if len(p.Ages) > 0 {
		accs = append(accs, "ages")
	}

	if len(accs) == 0 {
		return nil // nothing to add
	}

	for _, a := range accs {
		stmt, valueArgs = createSQLStatement(p, a)
		_, err := db.Exec(stmt, valueArgs...)
		if err != nil {
			log.Println("error creating project accessories: ", err)
			return err
		}
	}

	return nil
}

func createSQLStatement(p *Project, a string) (string, []interface{}) {
	valueStrings := []string{}
	valueArgs := []interface{}{}
	var tbl string
	var id string
	switch a {
	case AccAges:
		for i, age := range p.Ages {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, age.ID)
		}
		tbl = "project_ages"
		id = "ages_id"
	case AccCategories:
		for i, cat := range p.Categories {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, cat.ID)
		}
		tbl = "project_categories"
		id = "category_id"
	}

	return fmt.Sprintf(
		"INSERT INTO %s (project_id, %s) VALUES %s ON CONFLICT DO NOTHING", tbl, id, strings.Join(valueStrings, ","),
	), valueArgs
}

// DeleteProjectAssociations removes all associated records for a project
func DeleteProjectAssociations(ctx context.Context, db *sql.DB, projectID int) error {
	// Start transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Define tables to clean
	tables := []string{
		"project_ages",
		"project_categories",
	}

	// Delete from each table
	for _, table := range tables {
		_, err = tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE project_id = $1", table), projectID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete from %s: %v", table, err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
