package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system
type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	TextPermission bool      `json:"text_permission"`
	LeadInterest   bool      `json:"lead_interest"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *sql.DB, id string) (*User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, text_permission, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.TextPermission,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, text_permission, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user User
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.TextPermission,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, user *User) error {
	query := `
		INSERT INTO users (id, email, first_name, last_name, phone, text_permission)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	return db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.TextPermission,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// UpdateUser updates an existing user
func UpdateUser(db *sql.DB, user *User) error {
	query := `
		UPDATE users
		SET email = $1, first_name = $2, last_name = $3, phone = $4, text_permission = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at
	`

	return db.QueryRow(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.TextPermission,
		user.ID,
	).Scan(&user.UpdatedAt)
}

// CreateOrUpdateUser creates a new user or updates an existing one
func CreateOrUpdateUser(db *sql.DB, user *User) error {
	existingUser, err := GetUserByID(db, user.ID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return CreateUser(db, user)
	}

	return UpdateUser(db, user)
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, text_permission, created_at, updated_at
		FROM users
		ORDER BY last_name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone, &u.TextPermission,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
