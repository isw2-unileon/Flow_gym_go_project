package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Flow_gym_go_project/services"
	"github.com/DATA-DOG/go-sqlmock"
)

func makeAuthReq(method, url string, body interface{}, cookieValue string) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, url, &buf)
	if cookieValue != "" {
		req.AddCookie(&http.Cookie{Name: "user_id", Value: cookieValue})
	}
	return req
}

// ==========================================
// TESTS RegisterHandler
// ==========================================
func TestRegisterHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := RegisterHandler(db)

	validReq := RegisterRequest{Name: "Test", Email: "test@test.com", Password: "pass"}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/register", nil, ""))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/register", make(chan int), ""))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad JSON, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/register", RegisterRequest{}, ""))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing fields, got %d", rr.Code)
	}

	tooLongPass := strings.Repeat("a", 73)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/register", RegisterRequest{Name: "T", Email: "e", Password: tooLongPass}, ""))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for hashing error, got %d", rr.Code)
	}

	mock.ExpectQuery("INSERT INTO users").WithArgs("Test", "test@test.com", sqlmock.AnyArg(), "user").WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/register", validReq, ""))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for DB error, got %d", rr.Code)
	}

	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).AddRow(1, "Test", "test@test.com", "hash", "user")
	mock.ExpectQuery("INSERT INTO users").WithArgs("Test", "test@test.com", sqlmock.AnyArg(), "user").WillReturnRows(mockRows)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/register", validReq, ""))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS LoginHandler
// ==========================================
func TestLoginHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := LoginHandler(db)

	validReq := LoginRequest{Email: "test@test.com", Password: "mysecretpassword"}
	validHash, _ := services.HashPassword("mysecretpassword")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/login", nil, ""))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/login", make(chan int), ""))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad JSON, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/login", LoginRequest{}, ""))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing fields, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs("test@test.com").WillReturnError(sql.ErrNoRows)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/login", validReq, ""))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for not found, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs("test@test.com").WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/login", validReq, ""))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for DB error, got %d", rr.Code)
	}

	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).AddRow(1, "Test", "test@test.com", validHash, "user")
	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs("test@test.com").WillReturnRows(mockRows)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/login", LoginRequest{Email: "test@test.com", Password: "wrongpassword"}, ""))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong password, got %d", rr.Code)
	}

	mockRows2 := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).AddRow(1, "Test", "test@test.com", validHash, "user")
	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs("test@test.com").WillReturnRows(mockRows2)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/login", validReq, ""))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS LogoutHandler
// ==========================================
func TestLogoutHandler(t *testing.T) {
	handler := LogoutHandler()

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/logout", nil, ""))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("POST", "/logout", nil, ""))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// ==========================================
// TESTS MeHandler
// ==========================================
func TestMeHandler(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	handler := MeHandler(db)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/me", nil, ""))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for missing cookie, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/me", nil, "abc"))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for invalid cookie, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs(1).WillReturnError(sql.ErrNoRows)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/me", nil, "1"))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for not found, got %d", rr.Code)
	}

	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs(1).WillReturnError(errors.New("db error"))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/me", nil, "1"))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 for DB error, got %d", rr.Code)
	}

	mockRows := sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "role"}).AddRow(1, "Test", "test@test.com", "hash", "user")
	mock.ExpectQuery("SELECT id, name, email, password_hash, role FROM users").WithArgs(1).WillReturnRows(mockRows)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, makeAuthReq("GET", "/me", nil, "1"))
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}