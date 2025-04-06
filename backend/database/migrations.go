package database

import (
	"database/sql"
	"log"
)

// RunMigrations runs the database migrations
func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Create users table
	if err := createUsersTable(db); err != nil {
		return err
	}

	// Create projects table
	if err := createProjectsTable(db); err != nil {
		return err
	}

	// Create registrations table
	if err := createRegistrationsTable(db); err != nil {
		return err
	}

	// Create tools table
	if err := createToolsTable(db); err != nil {
		return err
	}

	// Create project_tools table
	if err := createProjectToolsTable(db); err != nil {
		return err
	}

	// Add example project (only runs if no projects exist)
	if err := addExampleProject(db); err != nil {
		log.Printf("Warning: Failed to add example project: %v", err)
		// Continue even if this fails, it's not critical
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createUsersTable creates the users table if it doesn't exist
func createUsersTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
                id TEXT PRIMARY KEY,
                email TEXT NOT NULL DEFAULT '' UNIQUE,
                first_name TEXT NOT NULL DEFAULT '',
                last_name TEXT NOT NULL DEFAULT '',
                name TEXT NOT NULL DEFAULT '',
                picture TEXT DEFAULT '',
                phone TEXT DEFAULT '',
                contact_email TEXT DEFAULT '',
                is_admin BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        )
        `
	_, err := db.Exec(query)
	return err
}

// createProjectsTable creates the projects table if it doesn't exist
func createProjectsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS projects (
                id SERIAL PRIMARY KEY,
                title TEXT NOT NULL,
                short_description TEXT NOT NULL,
                description TEXT NOT NULL,
                start_time TIME NOT NULL,
                end_time TIME NOT NULL,
                project_date DATE NOT NULL,
                max_capacity INTEGER NOT NULL,
                location_name TEXT,
                latitude DOUBLE PRECISION,
                longitude DOUBLE PRECISION,
                location_address TEXT,
                wheelchair_accessible BOOLEAN NOT NULL DEFAULT FALSE,
                lead_user_id TEXT REFERENCES users(id),
                created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        )
        `
	_, err := db.Exec(query)
	return err
}

// createRegistrationsTable creates the registrations table if it doesn't exist
func createRegistrationsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS registrations (
                id SERIAL PRIMARY KEY,
                user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                status TEXT NOT NULL DEFAULT 'registered',
                guest_count INTEGER NOT NULL DEFAULT 0,
                is_project_lead BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                UNIQUE(user_id, project_id, status)
        )
        `
	_, err := db.Exec(query)
	return err
}

// addExampleProject adds an example project if no projects exist
func addExampleProject(db *sql.DB) error {
	// First check if any projects exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		return err
	}

	// If projects already exist, don't add the example one
	if count > 0 {
		log.Println("Projects already exist, skipping example project creation")
		return nil
	}

	// Add example user first
	userQuery := `
		INSERT INTO users (id, email, name, first_name, last_name, picture, phone, contact_email, is_admin)
		VALUES (
			'example-user-123',
			'project.lead@example.com',
			'Example Project Lead',
			'Example',
			'Project Lead',
			'https://example.com/avatar.jpg',
			'555-0123',
			'project.lead@example.com',
			true
		)
		RETURNING id`

	var userID string
	err = db.QueryRow(userQuery).Scan(&userID)
	if err != nil {
		return err
	}

	// Add example project for July 12, 2025
	query := `INSERT INTO projects (
      title,
      short_description,
      description, 
      start_time, 
      end_time, 
      project_date, 
      max_capacity,
      location_name,
      latitude,
      longitude,
      lead_user_id,
      wheelchair_accessible,
      location_address
    ) VALUES (
      'Community Park Cleanup',
      'cleanup project',
      'Join us for a community park cleanup event! We will be cleaning up trash, planting flowers, and making general improvements to our local park. All supplies will be provided. Please wear comfortable clothes and bring water.',
      '09:00:00',
      '12:00:00',
      '2025-07-12',
      25,
      'Central Community Park',
      40.7128,
      -74.0060,
      $1,
      true,
      '123 Main Street, New York, NY 10001'
    )`

	_, err = db.Exec(query, userID) // Use userID here
	if err != nil {
		return err
	}

	log.Println("Successfully added example project")
	return nil
}

// createToolsTable creates the tools table if it doesn't exist
func createToolsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS tools (
                id SERIAL PRIMARY KEY,
                name TEXT NOT NULL UNIQUE,
                created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        )
        `
	_, err := db.Exec(query)
	return err
}

// createProjectToolsTable creates the project_tools junction table if it doesn't exist
func createProjectToolsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS project_tools (
                project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
                tool_id INTEGER REFERENCES tools(id) ON DELETE CASCADE,
                PRIMARY KEY (project_id, tool_id)
        )
        `
	_, err := db.Exec(query)
	return err
}
