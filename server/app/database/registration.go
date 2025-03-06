package database

import (
	"context"
	"log"

	"serve/models"
)

func PutRegistration(ctx context.Context, r models.Registration) (models.Registration, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning tx")
		return r, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `
INSERT INTO registrations (project_id, account_id, quantity, lead, updated_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (project_id, account_id) DO UPDATE SET quantity=excluded.quantity, lead=excluded.lead, updated_at=excluded.updated_at`
	_, err = tx.ExecContext(
		ctx, sqlStatement, r.ProjectID, r.AccountID, r.QtyEnrolled, r.Lead, r.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error executing update")
		return r, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing tx")
		return r, err
	}

	return r, nil
}

func GetRegistrations(ctx context.Context) ([]models.Registration, error) {
	registrations := []models.Registration{}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning tx")
		return registrations, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `
SELECT * FROM registrations`
	rows, err := tx.QueryContext(ctx, sqlStatement)
	if err != nil {
		log.Printf("Error getting registrations: %v", err)
		return registrations, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Registration
		if err = rows.Scan(&r.AccountID, &r.ProjectID, &r.QtyEnrolled, &r.UpdatedAt); err != nil {
			log.Printf("Error scanning")
			return registrations, err
		}
		registrations = append(registrations, r)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		log.Printf("Error row err")
		return registrations, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("Error commiting tx")
		return registrations, err
	}

	return registrations, nil
}
