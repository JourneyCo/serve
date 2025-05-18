package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// Project represents a project in the system
type Project struct {
	ID                   int                `json:"id"`
	GoogleID             *int               `json:"google_id"`
	Title                string             `json:"title"`
	ShortDescription     string             `json:"short_description"`
	Description          string             `json:"description"`
	Time                 string             `json:"time"`
	ProjectDate          time.Time          `json:"project_date"`
	MaxCapacity          int                `json:"max_capacity"`
	CurrentReg           int                `json:"current_registrations"`
	LocationName         string             `json:"location_name"`
	LocationAddress      string             `json:"location_address"`
	Latitude             float64            `json:"latitude"`
	Longitude            float64            `json:"longitude"`
	WheelchairAccessible bool               `json:"wheelchair_accessible"`
	LeadUserID           string             `json:"lead_user_id"`
	LeadUser             *User              `json:"lead_user,omitempty"`
	Tools                []ProjectAccessory `json:"tools,omitempty"`
	Supplies             []ProjectAccessory `json:"supplies,omitempty"`
	Categories           []ProjectAccessory `json:"categories,omitempty"`
	Ages                 []ProjectAccessory `json:"ages,omitempty"`
	Skills               []ProjectAccessory `json:"skills,omitempty"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

const (
	AccTools      = "tools"
	AccCategories = "categories"
	AccSupplies   = "supplies"
	AccSkills     = "skills"
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
                SELECT p.id, p.google_id, p.title, p.short_description, p.description, p.time, p.project_date, 
                p.max_capacity, p.location_name, p.location_address, p.latitude, p.longitude,
                p.wheelchair_accessible, p.created_at, p.updated_at, 
                COALESCE(SUM(CASE WHEN r.status = 'registered' THEN r.guest_count + 1 ELSE 0 END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id 
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
		if err = rows.Scan(
			&p.ID, &p.GoogleID, &p.Title, &p.ShortDescription, &p.Description, &p.Time, &p.ProjectDate,
			&p.MaxCapacity, &p.LocationName, &p.LocationAddress, &p.Latitude, &p.Longitude,
			&p.WheelchairAccessible, &p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
		); err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

// GetProjectByID retrieves a project by its ID
func GetProjectByID(ctx context.Context, db *sql.DB, id int) (*Project, error) {
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
	err := db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.ShortDescription, &p.Time, &p.ProjectDate,
		&p.MaxCapacity, &p.LocationName, &p.LocationAddress, &p.Latitude, &p.Longitude, &p.LeadUserID,
		&p.WheelchairAccessible, &p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Project not found
		}
		return nil, err
	}

	// Get tools for this project
	toolsQuery := `
		SELECT t.id, t.name FROM tools t
		JOIN project_tools pt ON t.id = pt.tool_id
		WHERE pt.project_id = $1
	`
	toolRows, err := db.QueryContext(ctx, toolsQuery, id)
	if err != nil {
		return nil, err
	}
	defer toolRows.Close()

	for toolRows.Next() {
		var tool ProjectAccessory
		if err := toolRows.Scan(&tool.ID, &tool.Name); err != nil {
			return nil, err
		}
		p.Tools = append(p.Tools, tool)
	}

	// Get supplies for this project
	suppliesQuery := `
		SELECT s.id, s.name FROM supplies s
		JOIN project_supplies ps ON s.id = ps.supply_id
		WHERE ps.project_id = $1
	`
	supplyRows, err := db.QueryContext(ctx, suppliesQuery, id)
	if err != nil {
		return nil, err
	}
	defer supplyRows.Close()

	for supplyRows.Next() {
		var supply ProjectAccessory
		if err := supplyRows.Scan(&supply.ID, &supply.Name); err != nil {
			return nil, err
		}
		p.Supplies = append(p.Supplies, supply)
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

	// Get skills for this project
	skillsQuery := `SELECT s.id, s.name FROM skills s JOIN project_skills ps ON s.id = ps.skill_id WHERE ps.project_id = $1`
	skillsRows, err := db.QueryContext(ctx, skillsQuery, id)
	if err != nil {
		return nil, err
	}
	defer skillsRows.Close()

	for skillsRows.Next() {
		var skill ProjectAccessory
		if err = skillsRows.Scan(&skill.ID, &skill.Name); err != nil {
			return nil, err
		}
		p.Skills = append(p.Skills, skill)
	}

	return &p, nil
}

// CreateProject creates a new project in the database
func CreateProject(ctx context.Context, db *sql.DB, project *Project) error {

	query := `
                INSERT INTO projects (google_id, title, short_description, description, time, project_date, max_capacity, 
                                    location_name, location_address, latitude, longitude, wheelchair_accessible, lead_user_id)
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
		project.Time,
		project.ProjectDate,
		project.MaxCapacity,
		project.LocationName,
		project.LocationAddress,
		project.Latitude,
		project.Longitude,
		project.WheelchairAccessible,
		project.LeadUserID,
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
                SET google_id=$13, title = $1, short_description = $2, description = $3, time = $4, project_date = $5, 
                max_capacity = $6, location_name = $7, location_address = $8, latitude = $9, longitude = $10,
                wheelchair_accessible = $11, updated_at = CURRENT_TIMESTAMP
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
		project.Time,
		project.ProjectDate,
		project.MaxCapacity,
		project.LocationName,
		project.LocationAddress,
		project.Latitude,
		project.Longitude,
		project.WheelchairAccessible,
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
	if len(p.Tools) > 0 {
		accs = append(accs, "tools")
	}
	if len(p.Supplies) > 0 {
		accs = append(accs, "supplies")
	}
	if len(p.Categories) > 0 {
		accs = append(accs, "categories")
	}
	if len(p.Skills) > 0 {
		accs = append(accs, "skills")
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
	case AccTools:
		for i, tool := range p.Tools {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, tool.ID)
		}
		tbl = "project_tools"
		id = "tool_id"
	case AccAges:
		for i, age := range p.Ages {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, age.ID)
		}
		tbl = "project_ages"
		id = "ages_id"
	case AccSkills:
		for i, skill := range p.Skills {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, skill.ID)
		}
		tbl = "project_skills"
		id = "skill_id"
	case AccCategories:
		for i, cat := range p.Categories {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, cat.ID)
		}
		tbl = "project_categories"
		id = "category_id"
	case AccSupplies:
		for i, sup := range p.Supplies {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, p.ID, sup.ID)
		}
		tbl = "project_supplies"
		id = "supply_id"
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
		"project_tools",
		"project_ages",
		"project_categories",
		"project_supplies",
		"project_skills",
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
