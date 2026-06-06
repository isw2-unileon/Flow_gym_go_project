package database

import (
	"os"
	"testing"
)

// Tests that ConnectDB returns an error when DATABASE_URL is empty.
func TestConnectDB_EmptyURL(t *testing.T) {

	// Save the original environment variable value.
	oldURL := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", oldURL)

	// Simulate a missing DATABASE_URL.
	os.Setenv("DATABASE_URL", "")

	db, err := ConnectDB()

	// No database connection should be returned.
	if db != nil {
		t.Errorf("Expected db to be nil, got %v", db)
	}

	// An error should be returned.
	if err == nil {
		t.Error("Expected an error because DATABASE_URL is empty, got nil")
	}
}

// Tests that ConnectDB fails when the database server is unreachable.
func TestConnectDB_PingFailure(t *testing.T) {

	// Save the original environment variable value.
	oldURL := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", oldURL)

	// Set an invalid PostgreSQL connection string.
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:54321/fake_db?sslmode=disable")

	db, err := ConnectDB()

	// No database connection should be returned.
	if db != nil {
		t.Errorf("Expected db to be nil due to ping failure, got %v", db)
	}

	// An error should be returned because the connection cannot be established.
	if err == nil {
		t.Error("Expected a ping connection error, got nil")
	}
}

// Tests that ConnectDB successfully connects to a valid database.
func TestConnectDB_Success(t *testing.T) {

	// Retrieve the configured database URL.
	url := os.Getenv("DATABASE_URL")

	// Skip the test if no valid database configuration exists.
	if url == "" || url == "postgres://user:pass@localhost:54321/fake_db?sslmode=disable" {
		t.Skip("Skipping success test: No valid DATABASE_URL found for testing connection")
	}

	// Attempt to connect to the database.
	db, err := ConnectDB()

	// The connection should succeed.
	if err != nil {
		t.Fatalf("Expected successful database connection, got error: %v", err)
	}

	// A valid database instance should be returned.
	if db == nil {
		t.Fatal("Expected db instance, got nil")
	}

	// Close the connection after the test.
	db.Close()
}