package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExerciseRepository_GetByName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Squat").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByName("Squat")
	if err != sql.ErrNoRows {
		t.Errorf("Expected ErrNoRows, got %v", err)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id, mg.name FROM exercises").
		WithArgs("Squat").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id", "mg_name"}).
			AddRow(1, "Squat", 2, "Legs"))

	ex, err := repo.GetByName("Squat")
	if err != nil || ex.Name != "Squat" {
		t.Errorf("Expected Squat, got error %v", err)
	}
}

func TestExerciseRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	mock.ExpectQuery("SELECT id, name, muscle_group_id FROM exercises").
		WillReturnError(errors.New("db error"))

	_, err := repo.GetAll()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	mock.ExpectQuery("SELECT id, name, muscle_group_id FROM exercises").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(1, "Bench", 1).
			AddRow(2, "Squat", 2))

	exercises, err := repo.GetAll()
	if err != nil || len(exercises) != 2 {
		t.Errorf("Expected 2 exercises, got %v (error: %v)", len(exercises), err)
	}

	mock.ExpectQuery("SELECT id, name, muscle_group_id FROM exercises").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(nil, "Bench", 1)) // El nil hará que rows.Scan falle

	_, err = repo.GetAll()
	if err == nil {
		t.Error("Expected scan error, got nil")
	}
}

func TestExerciseRepository_GetAlternativesByMuscleGroup(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 2).
		WillReturnError(errors.New("db error"))

	_, err := repo.GetAlternativesByMuscleGroup(1, 2)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(1, "Bench", 1))

	exercises, err := repo.GetAlternativesByMuscleGroup(1, 2)
	if err != nil || len(exercises) != 1 {
		t.Errorf("Expected 1 exercise, got %v (error: %v)", len(exercises), err)
	}

	mock.ExpectQuery("SELECT e.id, e.name, e.muscle_group_id FROM exercises").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "muscle_group_id"}).
			AddRow(nil, "Bench", 1)) // El nil hará que rows.Scan falle

	_, err = repo.GetAlternativesByMuscleGroup(1, 2)
	if err == nil {
		t.Error("Expected scan error, got nil")
	}
}