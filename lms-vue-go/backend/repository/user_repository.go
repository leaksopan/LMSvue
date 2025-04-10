package repository

import (
	"database/sql"
	"errors"
	"lms-vue-go/backend/config"
	"lms-vue-go/backend/models"
	"log"
)

// UserRepository handles database operations for users
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	// Check if DB is initialized
	if config.DB == nil {
		log.Println("WARNING: Database connection is nil in UserRepository")
	}
	return &UserRepository{
		DB: config.DB,
	}
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in FindByUsername")
		return nil, errors.New("database connection not initialized")
	}

	query := `
		SELECT id, username, password, email, role
		FROM users
		WHERE username = ?
	`

	var user models.User
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in FindByID")
		return nil, errors.New("database connection not initialized")
	}

	query := `
		SELECT id, username, password, email, role
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in Create")
		return errors.New("database connection not initialized")
	}

	query := `
		INSERT INTO users (username, password, email, role)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.DB.Exec(query,
		user.Username,
		user.Password,
		user.Email,
		user.Role,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in Update")
		return errors.New("database connection not initialized")
	}

	query := `
		UPDATE users
		SET username = ?, password = ?, email = ?, role = ?
		WHERE id = ?
	`

	_, err := r.DB.Exec(query,
		user.Username,
		user.Password,
		user.Email,
		user.Role,
		user.ID,
	)

	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	// Check if DB is nil
	if r.DB == nil {
		log.Println("ERROR: Database connection is nil in Delete")
		return errors.New("database connection not initialized")
	}

	query := `DELETE FROM users WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
