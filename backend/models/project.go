package models

import (
        "database/sql"
        "time"
)

// Project represents a project in the system
type Project struct {
        ID          int       `json:"id"`
        Title       string    `json:"title"`
        Description string    `json:"description"`
        StartTime   string    `json:"startTime"`
        EndTime     string    `json:"endTime"`
        ProjectDate time.Time `json:"projectDate"`
        MaxCapacity int       `json:"maxCapacity"`
        CurrentReg  int       `json:"currentRegistrations"`
        LocationName string    `json:"locationName"`
        Latitude    float64   `json:"latitude"`
        Longitude   float64   `json:"longitude"`
        CreatedAt   time.Time `json:"createdAt"`
        UpdatedAt   time.Time `json:"updatedAt"`
}

// GetAllProjects retrieves all projects from the database
func GetAllProjects(db *sql.DB) ([]Project, error) {
        query := `
                SELECT p.id, p.title, p.description, p.start_time, p.end_time, p.project_date, 
                p.max_capacity, p.location_name, p.latitude, p.longitude,
                p.created_at, p.updated_at, 
                COALESCE(COUNT(CASE WHEN r.status = 'registered' THEN r.id END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                GROUP BY p.id
                ORDER BY p.project_date ASC, p.start_time ASC
        `

        rows, err := db.Query(query)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var projects []Project
        for rows.Next() {
                var p Project
                if err := rows.Scan(
                        &p.ID, &p.Title, &p.Description, &p.StartTime, &p.EndTime, &p.ProjectDate,
                        &p.MaxCapacity, &p.LocationName, &p.Latitude, &p.Longitude,
                        &p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
                ); err != nil {
                        return nil, err
                }
                projects = append(projects, p)
        }

        if err := rows.Err(); err != nil {
                return nil, err
        }

        return projects, nil
}

// GetProjectByID retrieves a project by its ID
func GetProjectByID(db *sql.DB, id int) (*Project, error) {
        query := `
                SELECT p.id, p.title, p.description, p.start_time, p.end_time, p.project_date, 
                p.max_capacity, p.location_name, p.latitude, p.longitude,
                p.created_at, p.updated_at, 
                COALESCE(COUNT(CASE WHEN r.status = 'registered' THEN r.id END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                WHERE p.id = $1
                GROUP BY p.id
        `

        var p Project
        err := db.QueryRow(query, id).Scan(
                &p.ID, &p.Title, &p.Description, &p.StartTime, &p.EndTime, &p.ProjectDate,
                &p.MaxCapacity, &p.LocationName, &p.Latitude, &p.Longitude,
                &p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
        )

        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, nil // Project not found
                }
                return nil, err
        }

        return &p, nil
}

// CreateProject creates a new project in the database
func CreateProject(db *sql.DB, project *Project) error {
        query := `
                INSERT INTO projects (title, description, start_time, end_time, project_date, max_capacity, 
                                    location_name, latitude, longitude)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                RETURNING id, created_at, updated_at
        `

        return db.QueryRow(
                query,
                project.Title,
                project.Description,
                project.StartTime,
                project.EndTime,
                project.ProjectDate,
                project.MaxCapacity,
                project.LocationName,
                project.Latitude,
                project.Longitude,
        ).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

// UpdateProject updates an existing project
func UpdateProject(db *sql.DB, project *Project) error {
        query := `
                UPDATE projects
                SET title = $1, description = $2, start_time = $3, end_time = $4, project_date = $5, 
                max_capacity = $6, location_name = $7, latitude = $8, longitude = $9,
                updated_at = CURRENT_TIMESTAMP
                WHERE id = $10
                RETURNING updated_at
        `

        return db.QueryRow(
                query,
                project.Title,
                project.Description,
                project.StartTime,
                project.EndTime,
                project.ProjectDate,
                project.MaxCapacity,
                project.LocationName,
                project.Latitude,
                project.Longitude,
                project.ID,
        ).Scan(&project.UpdatedAt)
}

// DeleteProject deletes a project by its ID
func DeleteProject(db *sql.DB, id int) error {
        query := `DELETE FROM projects WHERE id = $1`
        _, err := db.Exec(query, id)
        return err
}

// GetUpcomingProjects retrieves projects that are starting within the given days
func GetUpcomingProjects(db *sql.DB, days int) ([]Project, error) {
        query := `
                SELECT p.id, p.title, p.description, p.start_time, p.end_time, p.project_date, 
                p.max_capacity, p.location_name, p.latitude, p.longitude,
                p.created_at, p.updated_at, 
                COALESCE(COUNT(CASE WHEN r.status = 'registered' THEN r.id END), 0) as current_registrations
                FROM projects p
                LEFT JOIN registrations r ON p.id = r.project_id
                WHERE p.project_date BETWEEN CURRENT_DATE AND CURRENT_DATE + $1::integer
                GROUP BY p.id
                ORDER BY p.project_date ASC, p.start_time ASC
        `

        rows, err := db.Query(query, days)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var projects []Project
        for rows.Next() {
                var p Project
                if err := rows.Scan(
                        &p.ID, &p.Title, &p.Description, &p.StartTime, &p.EndTime, &p.ProjectDate,
                        &p.MaxCapacity, &p.LocationName, &p.Latitude, &p.Longitude,
                        &p.CreatedAt, &p.UpdatedAt, &p.CurrentReg,
                ); err != nil {
                        return nil, err
                }
                projects = append(projects, p)
        }

        if err := rows.Err(); err != nil {
                return nil, err
        }

        return projects, nil
}
