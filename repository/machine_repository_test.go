package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMachineRepository_ReleaseExpiredMachines(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewMachineRepository(db)

	mock.ExpectExec("UPDATE machines SET is_available = true").WillReturnError(errors.New("db error"))
	err := repo.ReleaseExpiredMachines()
	if err == nil {
		t.Error("Expected error")
	}

	mock.ExpectExec("UPDATE machines SET is_available = true").WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.ReleaseExpiredMachines()
	if err != nil {
		t.Error("Expected no error")
	}
}

func TestMachineRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewMachineRepository(db)

	mock.ExpectExec("UPDATE machines SET").WillReturnError(errors.New("db error"))
	_, err := repo.GetAll()
	if err == nil {
		t.Error("Expected error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT id, name, is_available").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available", "occ", "last", "rel", "until"}).AddRow(1, "Bench", true, nil, nil, nil, nil))
	
	machines, err := repo.GetAll()
	if err != nil || len(machines) != 1 {
		t.Errorf("Expected 1 machine, got %v", machines)
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT id, name, is_available").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available", "occ", "last", "rel", "until"}).
			AddRow(nil, "Bench", true, nil, nil, nil, nil))
	_, err = repo.GetAll()
	if err == nil {
		t.Error("Expected scan error, got nil")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT id, name, is_available, occupied_by_user_id").WillReturnError(errors.New("query error"))
	_, err = repo.GetAll()
	if err == nil || err.Error() != "query error" {
		t.Error("Expected query error")
	}
}

func TestMachineRepository_GetAvailable(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewMachineRepository(db)

	mock.ExpectQuery("SELECT id, name, is_available FROM machines").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available"}).AddRow(1, "Bench", true))
	machines, err := repo.GetAvailable()
	if err != nil || len(machines) != 1 {
		t.Errorf("Expected 1 machine, got %v", machines)
	}

	mock.ExpectQuery("SELECT id, name, is_available FROM machines").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available"}).
			AddRow(nil, "Bench", true))
	_, err = repo.GetAvailable()
	if err == nil {
		t.Error("Expected scan error, got nil")
	}

	mock.ExpectQuery("SELECT id, name, is_available FROM machines").WillReturnError(errors.New("query error"))
	_, err = repo.GetAvailable()
	if err == nil || err.Error() != "query error" {
		t.Error("Expected query error")
	}
}

func TestMachineRepository_GetAvailableByExerciseID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewMachineRepository(db)

	mock.ExpectQuery("SELECT m.id, m.name, m.is_available").WithArgs(1).WillReturnError(sql.ErrNoRows)
	_, err := repo.GetAvailableByExerciseID(1)
	if err != sql.ErrNoRows {
		t.Error("Expected ErrNoRows")
	}

	mock.ExpectQuery("SELECT m.id, m.name, m.is_available").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available"}).AddRow(1, "Bench", true))
	machine, err := repo.GetAvailableByExerciseID(1)
	if err != nil || machine.Name != "Bench" {
		t.Error("Expected Bench")
	}
}

func TestMachineRepository_UpdateAvailability(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewMachineRepository(db)

	mock.ExpectExec("UPDATE machines SET is_available").WithArgs(true, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.UpdateAvailability(1, true)
	if err != nil {
		t.Error("Expected no error")
	}
}

func TestMachineRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewMachineRepository(db)

	mock.ExpectQuery("SELECT id, name, is_available").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "is_available", "occ", "last", "rel", "until"}).AddRow(1, "Bench", true, nil, nil, nil, nil))
	machine, err := repo.GetByID(1)
	if err != nil || machine.ID != 1 {
		t.Error("Expected machine ID 1")
	}

	mock.ExpectQuery("SELECT id, name, is_available").WithArgs(999).WillReturnError(errors.New("db error"))
	_, err = repo.GetByID(999)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestMachineRepository_UpdateAvailabilityWithUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewMachineRepository(db)

	mock.ExpectExec("UPDATE machines SET").WillReturnError(errors.New("db error"))
	err := repo.UpdateAvailabilityWithUser(1, 1, true, "user")
	if err == nil {
		t.Error("Expected error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE machines SET is_available = true").WithArgs(1, 1, "user").WillReturnResult(sqlmock.NewResult(0, 0))
	err = repo.UpdateAvailabilityWithUser(1, 1, true, "user")
	if err != sql.ErrNoRows {
		t.Error("Expected ErrNoRows")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE machines SET is_available = true").WithArgs(1, 1, "user").WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateAvailabilityWithUser(1, 1, true, "user")
	if err != nil {
		t.Error("Expected no error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT id FROM machines").WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
	err = repo.UpdateAvailabilityWithUser(2, 1, false, "user")
	if err != ErrUserAlreadyOccupiesMachine {
		t.Errorf("Expected ErrUserAlreadyOccupiesMachine, got %v", err)
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT id FROM machines").WithArgs(1, 2).WillReturnError(sql.ErrNoRows) // No tiene máquinas
	mock.ExpectExec("UPDATE machines SET is_available = false").WithArgs(2, 1, "user").WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.UpdateAvailabilityWithUser(2, 1, false, "user")
	if err != nil {
		t.Error("Expected no error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE machines SET is_available = true").WithArgs(1, 1, "user").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
	err = repo.UpdateAvailabilityWithUser(1, 1, true, "user")
	if err == nil || err.Error() != "rows affected error" {
		t.Error("Expected rows affected error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT id FROM machines").WithArgs(1, 2).WillReturnError(errors.New("db crash"))
	err = repo.UpdateAvailabilityWithUser(2, 1, false, "user")
	if err == nil || err.Error() != "db crash" {
		t.Error("Expected db crash")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT id FROM machines").WithArgs(1, 2).WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("UPDATE machines SET is_available = false").WithArgs(2, 1, "user").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
	err = repo.UpdateAvailabilityWithUser(2, 1, false, "user")
	if err == nil || err.Error() != "rows affected error" {
		t.Error("Expected rows affected error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0)) // Libera expiradas ok
	mock.ExpectExec("UPDATE machines SET is_available = true").WithArgs(3, 1, "user").WillReturnError(errors.New("exec error"))
	err = repo.UpdateAvailabilityWithUser(3, 1, true, "user")
	if err == nil || err.Error() != "exec error" {
		t.Error("Expected exec error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0)) // Libera expiradas ok
	mock.ExpectQuery("SELECT id FROM machines").WithArgs(1, 4).WillReturnError(sql.ErrNoRows) // No tiene máquina
	mock.ExpectExec("UPDATE machines SET is_available = false").WithArgs(4, 1, "user").WillReturnError(errors.New("exec error"))
	err = repo.UpdateAvailabilityWithUser(4, 1, false, "user")
	if err == nil || err.Error() != "exec error" {
		t.Error("Expected exec error")
	}

	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0)) // Libera expiradas ok
	mock.ExpectQuery("SELECT id FROM machines").WithArgs(1, 5).WillReturnError(sql.ErrNoRows) // No tiene máquina
	mock.ExpectExec("UPDATE machines SET is_available = false").WithArgs(5, 1, "user").WillReturnResult(sqlmock.NewResult(0, 0)) // 0 filas afectadas
	err = repo.UpdateAvailabilityWithUser(5, 1, false, "user")
	if err != sql.ErrNoRows {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}
}