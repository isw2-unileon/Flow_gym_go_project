package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Flow_gym_go_project/services"
)

// ErrorResponse represents a standard JSON error response.
type ErrorResponse struct {
	Message string `json:"message"`
}

// RecommendationHandler returns a recommended exercise and machine
// based on the exercise requested by the user.
func RecommendationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the exercise name from the query parameter.
		exerciseName := r.URL.Query().Get("exercise")

		// The exercise parameter is required.
		if exerciseName == "" {

			// Return a JSON error response.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(ErrorResponse{
				Message: "missing exercise parameter",
			})

			return
		}

		// Request a recommendation from the service layer.
		recommendation, err := services.GetRecommendation(db, exerciseName)
		if err != nil {

			// Return a JSON error if no recommendation can be generated.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(ErrorResponse{
				Message: err.Error(),
			})

			return
		}

		// Return the recommendation as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recommendation)
	}
}