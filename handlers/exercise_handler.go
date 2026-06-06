package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Flow_gym_go_project/repository"
)

// GetExercisesHandler returns all exercises stored in the database.
func GetExercisesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create the exercise repository using the shared database connection.
		exerciseRepo := repository.NewExerciseRepository(db)

		// Retrieve all exercises from the database.
		exercises, err := exerciseRepo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the exercises as a JSON response.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
	}
}

// GetExerciseByNameHandler returns a specific exercise by its name.
func GetExerciseByNameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the exercise name from the query parameter.
		name := r.URL.Query().Get("name")

		// The name parameter is required.
		if name == "" {
			http.Error(w, "missing exercise name", http.StatusBadRequest)
			return
		}

		// Create the exercise repository using the shared database connection.
		exerciseRepo := repository.NewExerciseRepository(db)

		// Retrieve the exercise by name.
		exercise, err := exerciseRepo.GetByName(name)
		if err != nil {

			// Return 404 if the exercise does not exist.
			if err == sql.ErrNoRows {
				http.Error(w, "exercise not found", http.StatusNotFound)
				return
			}

			// Return 500 for unexpected database errors.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the exercise as a JSON response.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercise)
	}
}