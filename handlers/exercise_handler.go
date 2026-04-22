package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Flow_gym_go_project/repository"
)

func GetExercisesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exerciseRepo := repository.NewExerciseRepository(db)

		exercises, err := exerciseRepo.GetAll()
		if err != nil {
			http.Error(w, "could not fetch exercises", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercises)
	}
}

func GetExerciseByNameHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing exercise name", http.StatusBadRequest)
			return
		}

		exerciseRepo := repository.NewExerciseRepository(db)
		exercise, err := exerciseRepo.GetByName(name)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "exercise not found", http.StatusNotFound)
				return
			}
			http.Error(w, "could not fetch exercise", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercise)
	}
}