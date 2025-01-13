package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
	"os"
	"serve/helpers"
	"serve/models"
	"time"
)

var db *sql.DB

func StartDB() *sql.DB {
	// Connect to Database
	fmt.Println("Connecting to database...")
	db, err := sql.Open("postgres", dataSource())
	defer db.Close()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
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

	return db
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

func PostProject(ctx context.Context) (models.Project, error) {
	tx, err := db.BeginTx(ctx, nil)
	var p models.Project
	p, ok := ctx.Value("project").(models.Project)
	if !ok {
		return p, errors.New("project not found in context")
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	sqlStatement := `
INSERT INTO projects (name, required, needed, admin_id, location_id, created_at, updated_at, google_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	now := time.Now()
	result, err := tx.Exec(sqlStatement, p.Name, p.Required, p.Needed, p.AdminID, p.LocationID, now, now, p.GoogleID)
	if err != nil {
		log.Printf("Error inserting project: %v", err)
		return p, err
	}
	// Get the ID of the order item just created.
	id, err := result.LastInsertId()
	if err != nil {
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
