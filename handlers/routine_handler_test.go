package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func makeJSONReq(method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	return httptest.NewRequest(method, url, &buf)
}

// ==========================================
// TESTS GetRoutinesByUserIDHandler
// ==========================================
func TestGetRoutinesByUserIDHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := GetRoutinesByUserIDHandler(db)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/routines", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing userId, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/routines?userId=abc", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid userId, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, user_id, name FROM routines").WithArgs(1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/routines?userId=1", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for DB error, got %d", rr.Code)
	}

	routineRows := sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(1, 1, "Leg Day")
	mock.ExpectQuery("SELECT id, user_id, name FROM routines").WithArgs(1).WillReturnRows(routineRows)
	
	exerciseRows := sqlmock.NewRows([]string{"id", "routine_id", "exercise_id", "exercise_order", "id", "name", "muscle_group_id"}).
		AddRow(1, 1, 5, 1, 5, "Squat", 2)
	mock.ExpectQuery("SELECT re.id, re.routine_id, re.exercise_id").WithArgs(1).WillReturnRows(exerciseRows)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/routines?userId=1", nil))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS CreateRoutineHandler
// ==========================================
func TestCreateRoutineHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := CreateRoutineHandler(db)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("GET", "/routines/create", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("POST", "/routines/create", make(chan int)))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad JSON, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("POST", "/routines/create", RoutineCreateInput{UserID: 1, Name: "Test"}))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing fields, got %d", rr.Code)
	}

	mock.ExpectQuery("INSERT INTO routines").WithArgs(1, "Test").WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("POST", "/routines/create", RoutineCreateInput{UserID: 1, Name: "Test", ExerciseIDs: []int{1}}))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for DB error on insert, got %d", rr.Code)
	}

	mock.ExpectQuery("INSERT INTO routines").WithArgs(1, "Test").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO routine_exercises").WithArgs(1, 2, 1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("POST", "/routines/create", RoutineCreateInput{UserID: 1, Name: "Test", ExerciseIDs: []int{2}}))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for DB error on exercises, got %d", rr.Code)
	}

	mock.ExpectQuery("INSERT INTO routines").WithArgs(1, "Test").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO routine_exercises").WithArgs(1, 2, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("POST", "/routines/create", RoutineCreateInput{UserID: 1, Name: "Test", ExerciseIDs: []int{2}}))
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", rr.Code)
	}
}

// ==========================================
// TESTS DeleteRoutineHandler
// ==========================================
func TestDeleteRoutineHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := DeleteRoutineHandler(db)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("POST", "/routines/delete", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("DELETE", "/routines/delete", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing params, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("DELETE", "/routines/delete?id=abc&userId=1", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad id, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("DELETE", "/routines/delete?id=1&userId=abc", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad userId, got %d", rr.Code)
	}

	mock.ExpectExec("DELETE FROM routines").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 0))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("DELETE", "/routines/delete?id=1&userId=1", nil))
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", rr.Code)
	}

	mock.ExpectExec("DELETE FROM routines").WithArgs(1, 1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("DELETE", "/routines/delete?id=1&userId=1", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 DB error, got %d", rr.Code)
	}

	mock.ExpectExec("DELETE FROM routines").WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("DELETE", "/routines/delete?id=1&userId=1", nil))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS UpdateRoutineHandler
// ==========================================
func TestUpdateRoutineHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := UpdateRoutineHandler(db)

	validInput := RoutineUpdateInput{RoutineID: 1, UserID: 1, Name: "Updated", ExerciseIDs: []int{2}}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("POST", "/routines/update", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("PUT", "/routines/update", make(chan int)))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("PUT", "/routines/update", RoutineUpdateInput{}))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("Updated", 1, 1).WillReturnResult(sqlmock.NewResult(0, 0))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("PUT", "/routines/update", validInput))
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("Updated", 1, 1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("PUT", "/routines/update", validInput))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("Updated", 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM routine_exercises").WithArgs(1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("PUT", "/routines/update", validInput))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("Updated", 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM routine_exercises").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO routine_exercises").WithArgs(1, 2, 1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("PUT", "/routines/update", validInput))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE routines SET name").WithArgs("Updated", 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM routine_exercises").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO routine_exercises").WithArgs(1, 2, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeJSONReq("PUT", "/routines/update", validInput))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}