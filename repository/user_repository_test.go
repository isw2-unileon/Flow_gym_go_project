package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("John", "john@test.com", "hash", "user").
		WillReturnError(errors.New("db error"))

	_, err := repo.Create("John", "john@test.com", "hash", "user")
	if err == nil || err.Error() != "db error" {
		t.Errorf("Expected db error, got %v", err)
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("John", "john@test.com", "hash", "user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).
			AddRow(1, "John", "john@test.com", "hash", "user"))

	user, err := repo.Create("John", "john@test.com", "hash", "user")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.ID != 1 || user.Name != "John" {
		t.Errorf("Expected user John with ID 1, got %+v", user)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").
		WithArgs("test@test.com").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByEmail("test@test.com")
	if err != sql.ErrNoRows {
		t.Errorf("Expected ErrNoRows, got %v", err)
	}

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").
		WithArgs("test@test.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).
			AddRow(1, "Test", "test@test.com", "hash", "user"))

	user, err := repo.GetByEmail("test@test.com")
	if err != nil || user.ID != 1 {
		t.Errorf("Expected user ID 1, got error %v", err)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(1)
	if err != sql.ErrNoRows {
		t.Errorf("Expected ErrNoRows, got %v", err)
	}

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).
			AddRow(1, "Test", "test@test.com", "hash", "user"))

	user, err := repo.GetByID(1)
	if err != nil || user.ID != 1 {
		t.Errorf("Expected user ID 1, got error %v", err)
	}
}