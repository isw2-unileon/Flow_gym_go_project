package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Flow_gym_go_project/services"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func RecommendationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exerciseName := r.URL.Query().Get("exercise")
		if exerciseName == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Message: "missing exercise parameter",
			})
			return
		}

		recommendation, err := services.GetRecommendation(db, exerciseName)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recommendation)
	}
}