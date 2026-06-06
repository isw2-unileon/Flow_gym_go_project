package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Flow_gym_go_project/repository"
)

// GetRoutinesByUserIDHandler handles GET requests to retrieve all routines
// belonging to a specific user.
func GetRoutinesByUserIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the user ID from the query parameters.
		userIdStr := r.URL.Query().Get("userId")

		// The userId parameter is required.
		if userIdStr == "" {
			http.Error(w, "missing userId parameter", http.StatusBadRequest)
			return
		}

		// Convert the user ID from string to integer.
		userID, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, "invalid userId format", http.StatusBadRequest)
			return
		}

		// Create the routine repository.
		routineRepo := repository.NewRoutineRepository(db)

		// Retrieve all routines belonging to the specified user.
		routines, err := routineRepo.GetByUserID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the routines as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routines)
	}
}

// RoutineCreateInput defines the JSON payload required to create a new routine.
type RoutineCreateInput struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	ExerciseIDs []int  `json:"exercise_ids"`
}

// CreateRoutineHandler handles POST requests to create a new workout routine.
func CreateRoutineHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only POST requests are allowed.
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input RoutineCreateInput

		// Decode the JSON request body.
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields.
		if input.Name == "" || len(input.ExerciseIDs) == 0 || input.UserID == 0 {
			http.Error(w, "Missing required fields (user_id, name or exercises)", http.StatusBadRequest)
			return
		}

		// Create the routine repository.
		routineRepo := repository.NewRoutineRepository(db)

		// Create the routine header and retrieve the generated routine ID.
		routineID, err := routineRepo.CreateRoutine(input.UserID, input.Name)
		if err != nil {
			http.Error(w, "Error saving routine header: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Save all selected exercises while preserving their order.
		for index, exerciseID := range input.ExerciseIDs {

			// Exercise order starts at 1.
			order := index + 1

			err = routineRepo.AddExerciseToRoutine(routineID, exerciseID, order)
			if err != nil {
				http.Error(w, "Error saving routine exercise: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Return a success response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Routine successfully created!",
		})
	}
}

// DeleteRoutineHandler handles DELETE requests to remove a routine.
func DeleteRoutineHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only DELETE requests are allowed.
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read routine ID and user ID from query parameters.
		routineIdStr := r.URL.Query().Get("id")
		userIdStr := r.URL.Query().Get("userId")

		// Both parameters are required.
		if routineIdStr == "" || userIdStr == "" {
			http.Error(w, "Missing id or userId parameters", http.StatusBadRequest)
			return
		}

		// Convert routine ID to integer.
		routineID, err := strconv.Atoi(routineIdStr)
		if err != nil {
			http.Error(w, "Invalid routine ID format", http.StatusBadRequest)
			return
		}

		// Convert user ID to integer.
		userID, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}

		// Create the routine repository.
		routineRepo := repository.NewRoutineRepository(db)

		// Delete the routine if it belongs to the specified user.
		err = routineRepo.DeleteRoutine(routineID, userID)

		if err != nil {

			// Return 404 if the routine does not exist or does not belong to the user.
			if err == sql.ErrNoRows {
				http.Error(w, "Routine not found or unauthorized", http.StatusNotFound)
				return
			}

			http.Error(w, "Error deleting routine", http.StatusInternalServerError)
			return
		}

		// Return a success response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Routine successfully deleted!",
		})
	}
}

// RoutineUpdateInput defines the JSON payload required to update a routine.
type RoutineUpdateInput struct {
	RoutineID   int    `json:"routine_id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	ExerciseIDs []int  `json:"exercise_ids"`
}

// UpdateRoutineHandler handles PUT requests to completely replace a routine.
func UpdateRoutineHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only PUT requests are allowed.
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var input RoutineUpdateInput

		// Decode the JSON request body.
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate required fields.
		if input.RoutineID == 0 || input.Name == "" || len(input.ExerciseIDs) == 0 || input.UserID == 0 {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Create the routine repository.
		routineRepo := repository.NewRoutineRepository(db)

		// Update the routine name.
		err = routineRepo.UpdateRoutineName(input.RoutineID, input.UserID, input.Name)
		if err != nil {

			// Return 404 if the routine does not exist or belongs to another user.
			if err == sql.ErrNoRows {
				http.Error(w, "Routine not found or unauthorized", http.StatusNotFound)
				return
			}

			http.Error(w, "Error updating routine name: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Remove all previously linked exercises.
		err = routineRepo.ClearRoutineExercises(input.RoutineID)
		if err != nil {
			http.Error(w, "Error clearing old exercises: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Reinsert the updated list of exercises in the specified order.
		for index, exerciseID := range input.ExerciseIDs {

			order := index + 1

			err = routineRepo.AddExerciseToRoutine(input.RoutineID, exerciseID, order)
			if err != nil {
				http.Error(w, "Error saving updated exercise: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Return a success response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Routine successfully updated!",
		})
	}
}