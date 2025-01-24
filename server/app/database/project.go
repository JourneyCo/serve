package database

import (
	"context"
	"log"
	"serve/models"
)

func PostProject(ctx context.Context, p models.Project) (models.Project, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Project{}, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	var id int64
	sqlStatement := `
INSERT INTO projects (name, required, needed, admin_id, location_id, created_at, updated_at, google_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	if err = tx.QueryRow(sqlStatement, p.Name, p.Required, p.Needed, p.AdminID, p.LocationID, p.CreatedAt, p.UpdatedAt, p.GoogleID).Scan(&id); err != nil {
		log.Printf("Error inserting project: %v", err)
		return p, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return p, err
	}

	p.ID = id
	log.Printf("Inserted project: %v", p.Name)
	return p, nil
}

func GetProjects(ctx context.Context) ([]models.Project, error) {
	projects := []models.Project{}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning tx")
		return projects, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `
SELECT * FROM projects`
	rows, err := tx.Query(sqlStatement)
	if err != nil {
		log.Printf("Error getting project: %v", err)
		return projects, err
	}
	defer rows.Close()

	for rows.Next() {
		var project models.Project
		if err = rows.Scan(&project.ID, &project.GoogleID, &project.Name, &project.Required, &project.Needed, &project.AdminID, &project.LocationID, &project.CreatedAt, &project.UpdatedAt); err != nil {
			log.Printf("Error scanning")
			return projects, err
		}
		projects = append(projects, project)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		log.Printf("Error row err")
		return projects, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("Error commiting tx")
		return projects, err
	}

	return projects, nil
}
