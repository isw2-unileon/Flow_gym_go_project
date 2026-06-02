package database

import (
	"os"
	"testing"
)

func TestConnectDB_EmptyURL(t *testing.T) {
	oldURL := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", oldURL)

	os.Setenv("DATABASE_URL", "")

	db, err := ConnectDB()
	if db != nil {
		t.Errorf("Expected db to be nil, got %v", db)
	}
	if err == nil {
		t.Error("Expected an error because DATABASE_URL is empty, got nil")
	}
}

func TestConnectDB_PingFailure(t *testing.T) {
	oldURL := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", oldURL)

	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:54321/fake_db?sslmode=disable")

	db, err := ConnectDB()
	if db != nil {
		t.Errorf("Expected db to be nil due to ping failure, got %v", db)
	}
	if err == nil {
		t.Error("Expected a ping connection error, got nil")
	}
}

func TestConnectDB_Success(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" || url == "postgres://user:pass@localhost:54321/fake_db?sslmode=disable" {
		t.Skip("Skipping success test: No valid DATABASE_URL found for testing connection")
	}

	db, err := ConnectDB()
	if err != nil {
		t.Fatalf("Expected successful database connection, got error: %v", err)
	}
	if db == nil {
		t.Fatal("Expected db instance, got nil")
	}
	db.Close()
}