package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// ==========================================
// TESTS GetExercisesHandler (GetAll)
// ==========================================

// Tests that GetExercisesHandler returns all exercises successfully.
func TestGetExercisesHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Mock rows returned by the database.
	rows := sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
		AddRow(1, "Bench Press", 1).
		AddRow(2, "Squat", 2)

	// Expect the SELECT query used to retrieve exercises.
	mock.ExpectQuery("SELECT id, name, muscle_group_id FROM exercises").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/exercises", nil)
	rr := httptest.NewRecorder()

	// Execute the handler.
	handler := GetExercisesHandler(db)
	handler.ServeHTTP(rr, req)

	// The request should succeed.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}
}

// Tests that GetExercisesHandler returns 500 when the database fails.
func TestGetExercisesHandler_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// Simulate a database error.
	mock.ExpectQuery("SELECT id, name, muscle_group_id FROM exercises").WillReturnError(errors.New("database down"))

	req := httptest.NewRequest("GET", "/exercises", nil)
	rr := httptest.NewRecorder()

	// Execute the handler.
	handler := GetExercisesHandler(db)
	handler.ServeHTTP(rr, req)

	// The handler should return Internal Server Error.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error, got %v", status)
	}
}

// ==========================================
// TESTS GetExerciseByNameHandler
// ==========================================

// Tests that GetExerciseByNameHandler returns Bad Request when name is missing.
func TestGetExerciseByNameHandler_MissingName(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	req := httptest.NewRequest("GET", "/exercise", nil)
	rr := httptest.NewRecorder()

	// Execute the handler without a name query parameter.
	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	// Missing name should return Bad Request.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %v", status)
	}
}

// Tests that GetExerciseByNameHandler returns Not Found when the exercise does not exist.
func TestGetExerciseByNameHandler_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// Simulate no exercise found.
	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name").
		WithArgs("Unknown").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("GET", "/exercise?name=Unknown", nil)
	rr := httptest.NewRecorder()

	// Execute the handler.
	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	// Unknown exercise should return Not Found.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %v", status)
	}
}

// Tests that GetExerciseByNameHandler returns 500 on unexpected database errors.
func TestGetExerciseByNameHandler_InternalError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// Simulate an unexpected SQL error.
	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name").
		WithArgs("Squat").
		WillReturnError(errors.New("random sql error"))

	req := httptest.NewRequest("GET", "/exercise?name=Squat", nil)
	rr := httptest.NewRecorder()

	// Execute the handler.
	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	// Unexpected database errors should return Internal Server Error.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error, got %v", status)
	}
}

// Tests that GetExerciseByNameHandler returns the requested exercise successfully.
func TestGetExerciseByNameHandler_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	// Mock exercise row with its muscle group name.
	rows := sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "name"}).
		AddRow(1, "Squat", 2, "Legs")

	// Expect query by exercise name.
	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name").
		WithArgs("Squat").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/exercise?name=Squat", nil)
	rr := httptest.NewRecorder()

	// Execute the handler.
	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	// The request should succeed.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}
}