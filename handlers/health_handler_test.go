package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthHandler verifies that the health endpoint returns status 200 and body "OK".
func TestHealthHandler(t *testing.T) {

	// Create a test HTTP GET request for the health endpoint.
	req := httptest.NewRequest("GET", "/health", nil)

	// Create a response recorder to capture the handler response.
	rr := httptest.NewRecorder()

	// Execute the health handler.
	HealthHandler(rr, req)

	// The handler should return HTTP 200 OK.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// The response body should be exactly "OK".
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}