package database

import (
	"database/sql"
	"fmt"
	"os"

	// PostgreSQL driver.
	_ "github.com/lib/pq"
)

// ConnectDB creates and validates a PostgreSQL database connection.
func ConnectDB() (*sql.DB, error) {

	// Read the database connection URL from the environment.
	databaseURL := os.Getenv("DATABASE_URL")

	// Return an error if the environment variable is missing.
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	// Open a PostgreSQL database connection.
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Verify that the database is reachable.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Return the validated database connection.
	return db, nil
}