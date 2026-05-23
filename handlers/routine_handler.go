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

// RoutineCreateInput defines the JSON structure expected from the frontend
type RoutineCreateInput struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	ExerciseIDs []int  `json:"exercise_ids"`
}

// CreateRoutineHandler handles POST requests to save a new workout routine
func CreateRoutineHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input RoutineCreateInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.Name == "" || len(input.ExerciseIDs) == 0 || input.UserID == 0 {
			http.Error(w, "Missing required fields (user_id, name or exercises)", http.StatusBadRequest)
			return
		}

		routineRepo := repository.NewRoutineRepository(db)

		// 1. Insertar la rutina en la tabla 'routines' y obtener su ID
		routineID, err := routineRepo.CreateRoutine(input.UserID, input.Name)
		if err != nil {
			http.Error(w, "Error saving routine header: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Insertar los ejercicios vinculados en la tabla intermedia manteniendo el orden del array
		for index, exerciseID := range input.ExerciseIDs {
			order := index + 1 // El primer ejercicio será el orden 1, el segundo el 2, etc.
			err = routineRepo.AddExerciseToRoutine(routineID, exerciseID, order)
			if err != nil {
				http.Error(w, "Error saving routine exercise: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Routine successfully created!"})
	}
}