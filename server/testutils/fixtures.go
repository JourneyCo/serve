
package testutils

import (
	"database/sql"
	"time"
)

// CreateTestProject creates a test project in the database
func CreateTestProject(db *sql.DB) (int, error) {
	var projectID int
	err := db.QueryRow(`
		INSERT INTO projects (
			title, description, short_description, project_date, 
			time, max_capacity, wheelchair_accessible, created_at, updated_at
		) VALUES (
			'Test Project', 'Test Description', 'Short Desc', 
			$1, '09:00 AM', 10, true, NOW(), NOW()
		) RETURNING id`,
		time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
	).Scan(&projectID)

	return projectID, err
}

// CleanTestData removes all test data from the database
func CleanTestData(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM registrations;
		DELETE FROM projects;
	`)
	return err
}
