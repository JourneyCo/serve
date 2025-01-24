package database

import (
	"context"
	"log"
	"serve/models"
)

// GetLocationByAddress will search for a location in the database by street number and street name
func GetLocationByAddress(ctx context.Context, number int, street string) (models.Location, error) {
	var lm models.Location

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return lm, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `SELECT * FROM locations WHERE number=$1 AND street=$2`
	if err = tx.QueryRow(sqlStatement, number, street).Scan(&lm); err != nil {
		return models.Location{}, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return models.Location{}, err
	}

	log.Printf("Inserted location: %v", lm.FormattedAddress)
	return lm, nil
}

func PostLocation(ctx context.Context, l models.Location) (models.Location, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Location{}, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	var id int64
	sqlStatement := `
INSERT INTO locations (latitude, longitude, info, street, number, city, state, postal_code, formatted_address, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	if err = tx.QueryRow(sqlStatement, l.Latitude, l.Longitude, l.Info, l.Street, l.Number, l.City, l.State, l.PostalCode, l.FormattedAddress, l.CreatedAt, l.UpdatedAt).Scan(&id); err != nil {
		return models.Location{}, err
	}
	// Get the ID of the order item just created.

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return l, err
	}

	l.ID = id
	log.Printf("Inserted location: %d, %v", l.ID, l.FormattedAddress)
	return l, nil
}
