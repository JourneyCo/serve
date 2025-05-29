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
	ID              int                `json:"id"`
	GoogleID        *int               `json:"google_id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Website         string             `json:"website"`
	Time            string             `json:"time"`
	ProjectDate     time.Time          `json:"project_date"`
	MaxCapacity     int                `json:"max_capacity"`
	CurrentReg      int                `json:"current_registrations"`
	Area            string             `json:"area"`
	LocationAddress string             `json:"location_address"`
	Latitude        float64            `json:"latitude"`
	Longitude       float64            `json:"longitude"`
	ServeLeadID     string             `json:"serve_lead_id"`
	ServeLead       *User              `json:"serve_lead,omitempty"`
	Types           []ProjectAccessory `json:"types,omitempty"`
	Ages            string             `json:"ages,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

const (
	AccTypes = "types"
)

// ProjectAccessory represents an accessory to the project
type ProjectAccessory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetAllProjects retrieves all projects from the database
func GetAllProjects(ctx context.Context, db *sql.DB) ([]Project, error) {
	query := `
                SELECT p.id, p.google_id, p.title, p.description, p.website, p.time, 
                p.max_capacity, p.area, p.location_address, p.latitude, p.longitude,
                p.created_at, p.updated_at, p.ages,
                COALESCE(COUNT(CASE WHEN r.status = 'registered' THEN 1 END) + SUM(CASE WHEN r.status = 'registered' THEN r.guest_count ELSE 0 END), 0) as current_registrations,
                array_to_string(COALESCE(array_agg(DISTINCT pc.type_id) FILTER (WHERE pc.type_id IS NOT NULL), ARRAY[]::integer[]), ',') as type_ids
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id 
                LEFT JOIN project_types pc ON p.id = pc.project_id
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
		var typeIDsStr string
		if err = rows.Scan(
			&p.ID, &p.GoogleID, &p.Title, &p.Description, &p.Website, &p.Time,
			&p.MaxCapacity, &p.Area, &p.LocationAddress, &p.Latitude, &p.Longitude,
			&p.CreatedAt, &p.UpdatedAt, &p.Ages, &p.CurrentReg, &typeIDsStr,
		); err != nil {
			return nil, err
		}

		// Convert comma-separated string to Types
		if typeIDsStr != "" {
			for _, idStr := range strings.Split(typeIDsStr, ",") {
				if id, err := strconv.Atoi(idStr); err == nil {
					p.Types = append(p.Types, ProjectAccessory{ID: id})
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
                SELECT p.id, p.title, p.description, p.website, p.time, p.project_date, 
                p.max_capacity, p.area, p.location_address, p.latitude, p.longitude, p.serve_lead_id,
                p.created_at, p.updated_at, p.ages,
                COALESCE(COUNT(CASE WHEN r.status = 'registered' THEN 1 END) + SUM(CASE WHEN r.status = 'registered' THEN r.guest_count ELSE 0 END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                WHERE p.id = $1
                GROUP BY p.id
        `

	var p Project
	err := db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.Website, &p.Time, &p.ProjectDate,
		&p.MaxCapacity, &p.Area, &p.LocationAddress, &p.Latitude, &p.Longitude, &p.ServeLeadID,
		&p.CreatedAt, &p.UpdatedAt, &p.Ages, &p.CurrentReg,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Project not found
		}
		return nil, err
	}

	// Get types for this project
	typesQuery := `
		SELECT c.id, c.type FROM types c
		JOIN project_types pc ON c.id = pc.type_id
		WHERE pc.project_id = $1
	`
	typeRows, err := db.QueryContext(ctx, typesQuery, id)
	if err != nil {
		return nil, err
	}
	defer typeRows.Close()

	for typeRows.Next() {
		var typ ProjectAccessory
		if err := typeRows.Scan(&typ.ID, &typ.Name); err != nil {
			return nil, err
		}
		p.Types = append(p.Types, typ)
	}

	return &p, nil
}

// CreateProject creates a new project in the database
func CreateProject(ctx context.Context, db *sql.DB, project *Project) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("error creating tx")
		return err
	}

	query := `
                INSERT INTO projects (google_id, title, description, website, time, project_date, max_capacity, 
                                    area, location_address, latitude, longitude, serve_lead_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
                RETURNING id, created_at, updated_at
        `

	err = tx.QueryRowContext(
		ctx,
		query,
		project.GoogleID,
		project.Title,
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
		tx.Rollback()
		return err
	}

	if err = insertAccessories(ctx, tx, project); err != nil {
		log.Println("error creating project associations: ", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// UpdateProject updates an existing project
func UpdateProject(ctx context.Context, db *sql.DB, project *Project) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("error creating tx")
		return err
	}

	query := `
                UPDATE projects
                SET google_id=$13, title = $1, description = $2, website = $3, time = $4, project_date = $5, 
                max_capacity = $6, area = $7, location_address = $8, latitude = $9, longitude = $10,
                updated_at = CURRENT_TIMESTAMP, ages = $11,
                WHERE id = $12
                RETURNING updated_at
        `
	err = tx.QueryRowContext(
		ctx,
		query,
		project.Title,
		project.Description,
		project.Website,
		project.Time,
		project.ProjectDate,
		project.MaxCapacity,
		project.Area,
		project.LocationAddress,
		project.Latitude,
		project.Longitude,
		project.Ages,
		project.ID,
		project.GoogleID,
	).Scan(&project.UpdatedAt)
	if err != nil {
		tx.Rollback()
		log.Println("error updating project: ", err)
		return err
	}

	if err = DeleteProjectAssociations(ctx, tx, project.ID); err != nil {
		tx.Rollback()
		log.Println("error updating project types: ", err)
		return err
	}

	if err = insertAccessories(ctx, tx, project); err != nil {
		tx.Rollback()
		log.Println("error inserting accessories")
		return err
	}

	tx.Commit()
	return nil
}

// DeleteProject deletes a project by its ID
func DeleteProject(ctx context.Context, db *sql.DB, id int) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := db.ExecContext(ctx, query, id)
	return err
}

func insertAccessories(ctx context.Context, tx *sql.Tx, p *Project) error {
	accs := []string{}
	var stmt string
	var valueArgs []any
	if len(p.Types) > 0 {
		accs = append(accs, "types")
	}

	if len(accs) == 0 {
		return nil // nothing to add
	}

	for _, a := range accs {
		stmt, valueArgs = createSQLStatement(p, a)
		_, err := tx.ExecContext(ctx, stmt, valueArgs...)
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
	case AccTypes:
		for i, cat := range p.Types {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, cat.ID)
		}
		tbl = "project_types"
		id = "type_id"
	}

	return fmt.Sprintf(
		"INSERT INTO %s (project_id, %s) VALUES %s ON CONFLICT DO NOTHING", tbl, id, strings.Join(valueStrings, ","),
	), valueArgs
}

// GetAllTypes retrieves all types from the types table
func GetAllTypes(ctx context.Context, db *sql.DB) ([]ProjectAccessory, error) {
	query := `SELECT id, type FROM types ORDER BY id`
	
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []ProjectAccessory
	for rows.Next() {
		var typ ProjectAccessory
		if err := rows.Scan(&typ.ID, &typ.Name); err != nil {
			return nil, err
		}
		types = append(types, typ)
	}

	return types, nil
}

// DeleteProjectAssociations removes all associated records for a project
func DeleteProjectAssociations(ctx context.Context, tx *sql.Tx, projectID int) error {
	// Define tables to clean
	tables := []string{
		"project_types",
	}

	// Delete from each table
	for _, table := range tables {
		_, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE project_id = $1", table), projectID)
		if err != nil {
			return fmt.Errorf("failed to delete from %s: %v", table, err)
		}
	}

	return nil
}
