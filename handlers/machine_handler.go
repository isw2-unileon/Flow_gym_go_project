package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Flow_gym_go_project/repository"
)

// UpdateMachineAvailabilityHandler updates machine availability using query parameters.
// This is the legacy GET-based version of the endpoint.
func UpdateMachineAvailabilityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the machine ID and availability value from query parameters.
		machineIDParam := r.URL.Query().Get("id")
		availableParam := r.URL.Query().Get("available")

		// Both parameters are required.
		if machineIDParam == "" || availableParam == "" {
			http.Error(w, "missing id or available parameter", http.StatusBadRequest)
			return
		}

		// Convert the machine ID from string to integer.
		machineID, err := strconv.Atoi(machineIDParam)
		if err != nil {
			http.Error(w, "invalid machine id", http.StatusBadRequest)
			return
		}

		// Convert the availability value from string to boolean.
		isAvailable, err := strconv.ParseBool(availableParam)
		if err != nil {
			http.Error(w, "invalid available value", http.StatusBadRequest)
			return
		}

		// Update the machine availability using the repository.
		machineRepo := repository.NewMachineRepository(db)
		err = machineRepo.UpdateAvailability(machineID, isAvailable)
		if err != nil {
			http.Error(w, "could not update machine availability", http.StatusInternalServerError)
			return
		}

		// Return a success response.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Machine availability updated successfully"))
	}
}

// GetMachinesHandler returns all machines from the database.
func GetMachinesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create the machine repository.
		machineRepo := repository.NewMachineRepository(db)

		// Retrieve all machines.
		machines, err := machineRepo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the machines as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(machines)
	}
}

// GetAvailableMachinesHandler returns only machines that are currently available.
func GetAvailableMachinesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Create the machine repository.
		machineRepo := repository.NewMachineRepository(db)

		// Retrieve available machines.
		machines, err := machineRepo.GetAvailable()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the available machines as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(machines)
	}
}

// UpdateMachineAvailabilityRequest represents the JSON body used to update machine availability.
type UpdateMachineAvailabilityRequest struct {
	ID        int  `json:"id"`
	Available bool `json:"available"`
}

// UpdateMachineAvailabilityPostHandler updates machine availability using an authenticated POST request.
// This is the main endpoint used by the frontend.
func UpdateMachineAvailabilityPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Only POST requests are allowed.
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req UpdateMachineAvailabilityRequest

		// Decode the JSON request body.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Read the authenticated user ID from the session cookie.
		cookie, err := r.Cookie("user_id")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Convert the cookie value to an integer user ID.
		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}

		// Load the current user from the database.
		userRepo := repository.NewUserRepository(db)
		user, err := userRepo.GetByID(userID)
		if err != nil {
			http.Error(w, "could not fetch current user", http.StatusInternalServerError)
			return
		}

		// Update the machine availability using user-aware business rules.
		machineRepo := repository.NewMachineRepository(db)
		err = machineRepo.UpdateAvailabilityWithUser(req.ID, user.ID, req.Available, user.Role)
		if err != nil {

			// A normal user cannot occupy more than one machine at the same time.
			if err == repository.ErrUserAlreadyOccupiesMachine {
				http.Error(w, "You already occupy another machine. Release it before selecting a new one.", http.StatusConflict)
				return
			}

			// This usually means the user does not have permission or cooldown rules blocked the update.
			if err == sql.ErrNoRows {
				http.Error(w, "machine cannot be updated by this user", http.StatusForbidden)
				return
			}

			// Any other error is treated as an internal server error.
			http.Error(w, "could not update machine availability", http.StatusInternalServerError)
			return
		}

		// Return a success message as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Machine availability updated successfully",
		})
	}
}

// GetMachineByIDHandler returns a single machine by its ID.
func GetMachineByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the machine ID from the query parameter.
		idParam := r.URL.Query().Get("id")

		// The ID parameter is required.
		if idParam == "" {
			http.Error(w, "missing machine id", http.StatusBadRequest)
			return
		}

		// Convert the ID from string to integer.
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "invalid machine id", http.StatusBadRequest)
			return
		}

		// Retrieve the machine from the database.
		machineRepo := repository.NewMachineRepository(db)
		machine, err := machineRepo.GetByID(id)
		if err != nil {

			// Return 404 if the machine does not exist.
			if err == sql.ErrNoRows {
				http.Error(w, "machine not found", http.StatusNotFound)
				return
			}

			// Return 500 for unexpected database errors.
			http.Error(w, "could not fetch machine", http.StatusInternalServerError)
			return
		}

		// Return the machine as JSON.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(machine)
	}
}