package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// ==========================================
// TESTS RecommendationHandler
// ==========================================
func TestRecommendationHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := RecommendationHandler(db)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/recommendation", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing exercise param, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Unknown").
		WillReturnError(sql.ErrNoRows)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/recommendation?exercise=Unknown", nil))
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for service error, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Bench Press").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "mg.name"}).
			AddRow(1, "Bench Press", 1, "Chest"))

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 1). // muscleGroupID = 1, excludedExerciseID = 1
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(2, "Push Up", 1))

	mock.ExpectQuery("SELECT m.id, m.name, m.is_available FROM machines").
		WithArgs(2). // exerciseID = 2
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available"}).
			AddRow(1, "Push Up Mat", true))

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/recommendation?exercise=Bench+Press", nil))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}