package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"serve/helpers"
	"serve/models"
)

var DB *sql.DB

func StartDB() {
	// Connect to Database
	fmt.Println("Connecting to database...")
	var err error
	DB, err = sql.Open("postgres", dataSource())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		log.Fatal("failed to get driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:./migrations",
		"postgres", driver)
	if err != nil || m == nil {
		log.Fatal("failed to create migration instance:", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("error migrating up: ", err) // only fatal if there is a migration up where there is a change
	}
}

func dataSource() string {
	//TODO: remove hardcoding before prod
	host := helpers.GetEnvVar("DB_HOST")
	dbUser := helpers.GetEnvVar("DB_USER")
	dbPass := helpers.GetEnvVar("DB_PASS")
	if os.Getenv("profile") == "prod" {
		host = "database"
	}
	return "postgresql://" + host + ":5432/journey" +
		"?user=" + dbUser + "&sslmode=disable&password=" + dbPass
}

func PostProject(ctx context.Context, p models.Project) (models.Project, error) {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Project{}, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `
INSERT INTO projects (name, required, needed, admin_id, location_id, created_at, updated_at, google_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	_, err = tx.Exec(sqlStatement, p.Name, p.Required, p.Needed, p.AdminID, p.LocationID, p.CreatedAt, p.UpdatedAt, p.GoogleID)
	if err != nil {
		log.Printf("Error inserting project: %v", err)
		return p, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return p, err
	}

	//p.ID = int64(id)
	log.Printf("Inserted project: %v", p.Name)
	return p, nil
}

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
