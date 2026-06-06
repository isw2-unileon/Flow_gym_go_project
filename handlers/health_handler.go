package handlers

import (
	"net/http"
)

// HealthHandler returns a simple OK response to confirm that the server is running.
func HealthHandler(w http.ResponseWriter, r *http.Request) {

	// Set the HTTP status code to 200 OK.
	w.WriteHeader(http.StatusOK)

	// Write a plain text response body.
	w.Write([]byte("OK"))
}