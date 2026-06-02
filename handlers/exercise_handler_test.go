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

func TestGetExercisesHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
		AddRow(1, "Bench Press", 1).
		AddRow(2, "Squat", 2)
	
	mock.ExpectQuery("SELECT id, name, muscle_group_id FROM exercises").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/exercises", nil)
	rr := httptest.NewRecorder()

	handler := GetExercisesHandler(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}
}

func TestGetExercisesHandler_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, muscle_group_id FROM exercises").WillReturnError(errors.New("database down"))

	req := httptest.NewRequest("GET", "/exercises", nil)
	rr := httptest.NewRecorder()

	handler := GetExercisesHandler(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error, got %v", status)
	}
}

// ==========================================
// TESTS GetExerciseByNameHandler
// ==========================================

func TestGetExerciseByNameHandler_MissingName(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	req := httptest.NewRequest("GET", "/exercise", nil)
	rr := httptest.NewRecorder()

	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request, got %v", status)
	}
}

func TestGetExerciseByNameHandler_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name").
		WithArgs("Unknown").
		WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest("GET", "/exercise?name=Unknown", nil)
	rr := httptest.NewRecorder()

	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %v", status)
	}
}

func TestGetExerciseByNameHandler_InternalError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name").
		WithArgs("Squat").
		WillReturnError(errors.New("random sql error"))

	req := httptest.NewRequest("GET", "/exercise?name=Squat", nil)
	rr := httptest.NewRecorder()

	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status 500 Internal Server Error, got %v", status)
	}
}

func TestGetExerciseByNameHandler_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "name"}).
		AddRow(1, "Squat", 2, "Legs")
	
	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name").
		WithArgs("Squat").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/exercise?name=Squat", nil)
	rr := httptest.NewRecorder()

	handler := GetExerciseByNameHandler(db)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}
}