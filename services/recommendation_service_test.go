package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetRecommendation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Unknown").
		WillReturnError(sql.ErrNoRows)
	
	_, err = GetRecommendation(db, "Unknown")
	if err == nil || err.Error() != "exercise not found" {
		t.Errorf("Expected 'exercise not found' error, got %v", err)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Bench Press").
		WillReturnError(errors.New("db error"))
	
	_, err = GetRecommendation(db, "Bench Press")
	if err == nil || err.Error() != "db error" {
		t.Errorf("Expected 'db error', got %v", err)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Bench Press").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "mg.name"}).
			AddRow(1, "Bench Press", 1, "Chest"))

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 1).
		WillReturnError(errors.New("db error on alternatives"))

	_, err = GetRecommendation(db, "Bench Press")
	if err == nil || err.Error() != "db error on alternatives" {
		t.Errorf("Expected 'db error on alternatives', got %v", err)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Bench Press").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "mg.name"}).
			AddRow(1, "Bench Press", 1, "Chest"))

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(2, "Push Up", 1))

	mock.ExpectQuery("SELECT m.id, m.name, m.is_available FROM machines").
		WithArgs(2).
		WillReturnError(sql.ErrNoRows) 

	_, err = GetRecommendation(db, "Bench Press")
	if err == nil || err.Error() != "no available recommendation found for this muscle group" {
		t.Errorf("Expected 'no available recommendation found...', got %v", err)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Bench Press").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "mg.name"}).
			AddRow(1, "Bench Press", 1, "Chest"))

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(2, "Push Up", 1))

	mock.ExpectQuery("SELECT m.id, m.name, m.is_available FROM machines").
		WithArgs(2).
		WillReturnError(errors.New("machine db crash"))

	_, err = GetRecommendation(db, "Bench Press")
	if err == nil || err.Error() != "machine db crash" {
		t.Errorf("Expected 'machine db crash', got %v", err)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Bench Press").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "mg.name"}).
			AddRow(1, "Bench Press", 1, "Chest"))

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(2, "Push Up", 1))

	mock.ExpectQuery("SELECT m.id, m.name, m.is_available FROM machines").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available"}).
			AddRow(1, "Push Up Mat", true))

	rec, err := GetRecommendation(db, "Bench Press")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if rec.RequestedExercise != "Bench Press" || rec.RecommendedExercise != "Push Up" || rec.Machine != "Push Up Mat" {
		t.Errorf("Received unexpected recommendation values: %+v", rec)
	}
}