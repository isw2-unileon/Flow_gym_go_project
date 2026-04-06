package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Flow_gym_go_project/services"
)

func RecommendationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exerciseName := r.URL.Query().Get("exercise")
		if exerciseName == "" {
			http.Error(w, "missing exercise parameter", http.StatusBadRequest)
			return
		}

		recommendation, err := services.GetRecommendation(db, exerciseName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recommendation)
	}
}