package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRoutineRepository_GetByUserID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectQuery("SELECT id, user_id, name FROM routines").WillReturnError(errors.New("db error"))
	_, err := repo.GetByUserID(1)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	mock.ExpectQuery("SELECT id, user_id, name FROM routines").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(1, 1, "Legs"))
	
	mock.ExpectQuery("SELECT re.id, re.routine_id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"re.id", "re.routine_id", "re.exercise_id", "re.exercise_order", "e.id", "e.name", "e.muscle_group_id"}).
			AddRow(1, 1, 2, 1, 2, "Squat", 2))

	routines, err := repo.GetByUserID(1)
	if err != nil || len(routines) != 1 || len(routines[0].Exercises) != 1 {
		t.Errorf("Expected 1 routine with 1 exercise, got %v", routines)
	}

	mock.ExpectQuery("SELECT id, user_id, name FROM routines").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(nil, 1, "Legs"))
	_, err = repo.GetByUserID(2)
	if err == nil {
		t.Error("Expected scan error, got nil")
	}
}

func TestRoutineRepository_CreateRoutine(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectQuery("INSERT INTO routines").WithArgs(1, "Test").WillReturnError(errors.New("db error"))
	_, err := repo.CreateRoutine(1, "Test")
	if err == nil {
		t.Error("Expected error")
	}

	mock.ExpectQuery("INSERT INTO routines").WithArgs(1, "Test").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
	id, err := repo.CreateRoutine(1, "Test")
	if err != nil || id != 5 {
		t.Errorf("Expected id 5, got %d", id)
	}
}

func TestRoutineRepository_AddExerciseToRoutine(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec("INSERT INTO routine_exercises").WithArgs(1, 2, 3).WillReturnError(errors.New("db error"))
	err := repo.AddExerciseToRoutine(1, 2, 3)
	if err == nil {
		t.Error("Expected error")
	}

	mock.ExpectExec("INSERT INTO routine_exercises").WithArgs(1, 2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.AddExerciseToRoutine(1, 2, 3)
	if err != nil {
		t.Error("Expected no error")
	}
}

func TestRoutineRepository_DeleteRoutine(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec("DELETE FROM routines").WithArgs(1, 1).WillReturnError(errors.New("db error"))
	err := repo.DeleteRoutine(1, 1)
	if err == nil {
		t.Error("Expected error")
	}

	mock.ExpectExec("DELETE FROM routines").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 0))
	err = repo.DeleteRoutine(1, 1)
	if err != sql.ErrNoRows {
		t.Error("Expected ErrNoRows")
	}

	mock.ExpectExec("DELETE FROM routines").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.DeleteRoutine(1, 1)
	if err != nil {
		t.Error("Expected no error")
	}

	mock.ExpectExec("DELETE FROM routines").WithArgs(2, 2).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
	err = repo.DeleteRoutine(2, 2)
	if err == nil || err.Error() != "rows affected error" {
		t.Error("Expected rows affected error")
	}
}

func TestRoutineRepository_UpdateRoutineName(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec("UPDATE routines SET name").WithArgs("New", 1, 1).WillReturnError(errors.New("db error"))
	err := repo.UpdateRoutineName(1, 1, "New")
	if err == nil {
		t.Error("Expected error")
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("New", 1, 1).WillReturnResult(sqlmock.NewResult(0, 0))
	err = repo.UpdateRoutineName(1, 1, "New")
	if err != sql.ErrNoRows {
		t.Error("Expected ErrNoRows")
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("New", 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateRoutineName(1, 1, "New")
	if err != nil {
		t.Error("Expected no error")
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("Error", 2, 2).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
	err = repo.UpdateRoutineName(2, 2, "Error")
	if err == nil || err.Error() != "rows affected error" {
		t.Error("Expected rows affected error")
	}
}

func TestRoutineRepository_ClearRoutineExercises(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec("DELETE FROM routine_exercises").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.ClearRoutineExercises(1)
	if err != nil {
		t.Error("Expected no error")
	}
}