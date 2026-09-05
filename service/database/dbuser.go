package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ematri2108281/wasatext/service/components/schema"
	"github.com/google/uuid"
)

// CreateUser inserts a new user into the database
func (db *appdbimpl) CreateUser(u *schema.User) (string, error) {
	if u == nil {
		return "", fmt.Errorf("user cannot be nil")
	}

	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=?)`, u.Username).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return "", schema.ErrUsernameTaken
	}
	u.ID = uuid.New().String()
	_, err = db.c.Exec(`INSERT INTO users(id, username, photo) VALUES (?, ?, ?)`, u.ID, u.Username, u.Photo)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return u.ID, nil
}

// GetAllUsers returns all users
func (db *appdbimpl) GetAllUsers() ([]schema.User, error) {
	rows, err := db.c.Query(`SELECT id, username, photo FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []schema.User
	for rows.Next() {
		var u schema.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Photo); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

// GetUserById returns a user by their ID
func (db *appdbimpl) GetUserById(userID string) (*schema.User, error) {
	var u schema.User
	err := db.c.QueryRow(`SELECT id, username, photo FROM users WHERE id=?`, userID).Scan(&u.ID, &u.Username, &u.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schema.ErrUserDoesNotExist
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &u, nil
}

// GetUserByName returns a user by their username
func (db *appdbimpl) GetUserByName(username string) (*schema.User, error) {
	var u schema.User
	err := db.c.QueryRow(`SELECT id, username, photo FROM users WHERE username=?`, username).Scan(&u.ID, &u.Username, &u.Photo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, schema.ErrUserDoesNotExist
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &u, nil
}

// SearchUserByUsername returns users whose username contains the search string
func (db *appdbimpl) SearchUserByUsername(username string) ([]schema.User, error) {
	rows, err := db.c.Query(`SELECT id, username, photo FROM users WHERE username LIKE ?`, "%"+username+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []schema.User
	for rows.Next() {
		var u schema.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Photo); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

// UpdateUsername changes a user's username
func (db *appdbimpl) UpdateUsername(userID, newUsername string) error {
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=? AND id<>?)`, newUsername, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check username availability: %w", err)
	}
	if exists {
		return schema.ErrUsernameTaken
	}

	res, err := db.c.Exec(`UPDATE users SET username=? WHERE id=?`, newUsername, userID)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return schema.ErrUserDoesNotExist
	}
	return nil
}

// UpdateUserPhoto updates the profile photo of a user
func (db *appdbimpl) UpdateUserPhoto(userID string, photo []byte) error {
	var exists bool
	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id=?)`, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return schema.ErrUserDoesNotExist
	}

	_, err = db.c.Exec(`UPDATE users SET photo=? WHERE id=?`, photo, userID)
	if err != nil {
		return fmt.Errorf("failed to update photo: %w", err)
	}
	return nil
}
