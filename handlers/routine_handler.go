package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Flow_gym_go_project/repository"
)

// GetRoutinesByUserIDHandler handles GET requests to fetch routines for a specific user
func GetRoutinesByUserIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We retrieve the userId from the URL
		userIdStr := r.URL.Query().Get("userId")
		if userIdStr == "" {
			http.Error(w, "missing userId parameter", http.StatusBadRequest)
			return
		}

		// Convert the URL string to an integer
		userID, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, "invalid userId format", http.StatusBadRequest)
			return
		}

		// We instantiate the repository and look for the routines
		routineRepo := repository.NewRoutineRepository(db)
		routines, err := routineRepo.GetByUserID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the routines in JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routines)
	}
}