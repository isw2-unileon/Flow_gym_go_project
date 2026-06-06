package repository

import (
	"database/sql"

	"Flow_gym_go_project/models"
)

// UserRepository provides database access methods related to users.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create inserts a new user into the database and returns
// the newly created user record.
func (r *UserRepository) Create(name string, email string, passwordHash string, role string) (*models.User, error) {

	// Insert a new user and immediately return the created record.
	query := `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role
	`

	var user models.User

	// Execute the query and map the returned row into the User model.
	err := r.DB.QueryRow(query, name, email, passwordHash, role).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user using their email address.
// This method is primarily used during authentication.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {

	// Query the user by email.
	query := `
		SELECT id, name, email, password_hash, role
		FROM users
		WHERE email = $1
	`

	var user models.User

	// Map the database result into the User model.
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user using their unique database ID.
// This method is mainly used for session validation and
// retrieving the currently authenticated user.
func (r *UserRepository) GetByID(id int) (*models.User, error) {

	// Query the user by primary key.
	query := `
		SELECT id, name, email, password_hash, role
		FROM users
		WHERE id = $1
	`

	var user models.User

	// Map the database result into the User model.
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}