package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// ==========================================
// TESTS UpdateMachineAvailabilityHandler
// ==========================================
func TestUpdateMachineAvailabilityHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := UpdateMachineAvailabilityHandler(db)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines/update", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines/update?id=abc&available=true", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid ID, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines/update?id=1&available=abc", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid boolean, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE machines SET is_available = \\$1").WithArgs(true, 1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines/update?id=1&available=true", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for DB error, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE machines SET is_available = \\$1").WithArgs(true, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines/update?id=1&available=true", nil))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS GetMachinesHandler
// ==========================================
func TestGetMachinesHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := GetMachinesHandler(db)

	mock.ExpectExec("UPDATE machines SET is_available = true").WillReturnError(errors.New("db error"))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr.Code)
	}

	mock.ExpectExec("UPDATE machines SET is_available = true").WillReturnResult(sqlmock.NewResult(0, 0))
	rows := sqlmock.NewRows([]string{"id", "name", "is_available", "occupied_by_user_id", "last_used_by_user_id", "last_released_at", "occupied_until"}).
		AddRow(1, "Treadmill", true, nil, nil, nil, nil)
	mock.ExpectQuery("SELECT id, name, is_available, occupied_by_user_id").WillReturnRows(rows)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines", nil))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS GetAvailableMachinesHandler
// ==========================================
func TestGetAvailableMachinesHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := GetAvailableMachinesHandler(db)

	mock.ExpectQuery("SELECT id, name, is_available FROM machines").WillReturnError(errors.New("db error"))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines/available", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr.Code)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "is_available"}).AddRow(1, "Treadmill", true)
	mock.ExpectQuery("SELECT id, name, is_available FROM machines").WillReturnRows(rows)
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machines/available", nil))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS GetMachineByIDHandler
// ==========================================
func TestGetMachineByIDHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := GetMachineByIDHandler(db)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machine", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machine?id=abc", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, name, is_available").WithArgs(999).WillReturnError(sql.ErrNoRows)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machine?id=999", nil))
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, name, is_available").WithArgs(1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machine?id=1", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr.Code)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "is_available", "occupied_by_user_id", "last_used_by_user_id", "last_released_at", "occupied_until"}).
		AddRow(1, "Bench", true, nil, nil, nil, nil)
	mock.ExpectQuery("SELECT id, name, is_available").WithArgs(1).WillReturnRows(rows)
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/machine?id=1", nil))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}


// ==========================================
// TESTS UpdateMachineAvailabilityPostHandler
// ==========================================
func TestUpdateMachineAvailabilityPostHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := UpdateMachineAvailabilityPostHandler(db)

	makeReq := func(method string, body interface{}, cookieValue string) *http.Request {
		var buf bytes.Buffer
		if body != nil {
			json.NewEncoder(&buf).Encode(body)
		}
		req := httptest.NewRequest(method, "/machines/update-post", &buf)
		if cookieValue != "" {
			req.AddCookie(&http.Cookie{Name: "user_id", Value: cookieValue})
		}
		return req
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("GET", nil, ""))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", make(chan int), "")) // Body que falla al decodificar
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad JSON, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", UpdateMachineAvailabilityRequest{ID: 1, Available: true}, ""))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for missing cookie, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", UpdateMachineAvailabilityRequest{ID: 1, Available: true}, "abc"))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for invalid cookie, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs(1).WillReturnError(sql.ErrNoRows)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", UpdateMachineAvailabilityRequest{ID: 1, Available: true}, "1"))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for user DB error, got %d", rr.Code)
	}

	mockUser := func() {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).
			AddRow(1, "Test User", "test@test.com", "hash", "user")
		mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs(1).WillReturnRows(rows)
	}

	mockUser()
	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0)) // Release expired
	mock.ExpectQuery("SELECT id FROM machines WHERE occupied_by_user_id").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2)) // Devuelve otra máquina
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", UpdateMachineAvailabilityRequest{ID: 1, Available: false}, "1"))
	if rr.Code != http.StatusConflict {
		t.Errorf("Expected 409 Conflict, got %d", rr.Code)
	}

	mockUser()
	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0)) // Release expired
	mock.ExpectExec("UPDATE machines SET is_available = true").WithArgs(1, 1, "user").WillReturnResult(sqlmock.NewResult(0, 0)) // 0 afectadas
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", UpdateMachineAvailabilityRequest{ID: 1, Available: true}, "1"))
	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", rr.Code)
	}

	mockUser()
	mock.ExpectExec("UPDATE machines SET").WillReturnError(errors.New("db down")) // Release expired falla
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", UpdateMachineAvailabilityRequest{ID: 1, Available: true}, "1"))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Error, got %d", rr.Code)
	}

	mockUser()
	mock.ExpectExec("UPDATE machines SET").WillReturnResult(sqlmock.NewResult(0, 0)) // Release expired
	mock.ExpectExec("UPDATE machines SET is_available = true").WithArgs(1, 1, "user").WillReturnResult(sqlmock.NewResult(1, 1)) // 1 afectada
	
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeReq("POST", UpdateMachineAvailabilityRequest{ID: 1, Available: true}, "1"))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}