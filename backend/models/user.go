package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Picture      string    `json:"picture"`
	Phone        string    `json:"phone"`
	ContactEmail string    `json:"contactEmail"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IsAdmin      bool      `json:"isAdmin"`
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *sql.DB, id string) (*User, error) {
	query := `
		SELECT id, email, name, picture, phone, contactEmail, created_at, updated_at, is_admin
		FROM users
		WHERE id = $1
	`

	var user User
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture, &user.Phone, &user.ContactEmail,
		&user.CreatedAt, &user.UpdatedAt, &user.IsAdmin,
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
		SELECT id, email, name, picture, phone, contactEmail, created_at, updated_at, is_admin
		FROM users
		WHERE email = $1
	`

	var user User
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Picture, &user.Phone, &user.ContactEmail,
		&user.CreatedAt, &user.UpdatedAt, &user.IsAdmin,
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
		INSERT INTO users (id, email, name, picture, phone, contactEmail, is_admin)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	return db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Name,
		user.Picture,
		user.Phone,
		user.ContactEmail,
		user.IsAdmin,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// UpdateUser updates an existing user
func UpdateUser(db *sql.DB, user *User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, picture = $3, phone = $4, contactEmail = $5, is_admin = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`

	return db.QueryRow(
		query,
		user.Email,
		user.Name,
		user.Picture,
		user.Phone,
		user.ContactEmail,
		user.IsAdmin,
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
		SELECT id, email, name, picture, phone, contactEmail, created_at, updated_at, is_admin
		FROM users
		ORDER BY name
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
			&u.ID, &u.Email, &u.Name, &u.Picture, &u.Phone, &u.ContactEmail,
			&u.CreatedAt, &u.UpdatedAt, &u.IsAdmin,
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

// SetUserAdmin sets the admin status for a user
func SetUserAdmin(db *sql.DB, userID string, isAdmin bool) error {
	query := `
		UPDATE users
		SET is_admin = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := db.Exec(query, isAdmin, userID)
	return err
}
