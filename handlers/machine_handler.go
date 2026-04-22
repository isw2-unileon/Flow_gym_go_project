package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Flow_gym_go_project/repository"
)

func UpdateMachineAvailabilityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		machineIDParam := r.URL.Query().Get("id")
		availableParam := r.URL.Query().Get("available")

		if machineIDParam == "" || availableParam == "" {
			http.Error(w, "missing id or available parameter", http.StatusBadRequest)
			return
		}

		machineID, err := strconv.Atoi(machineIDParam)
		if err != nil {
			http.Error(w, "invalid machine id", http.StatusBadRequest)
			return
		}

		isAvailable, err := strconv.ParseBool(availableParam)
		if err != nil {
			http.Error(w, "invalid available value", http.StatusBadRequest)
			return
		}

		machineRepo := repository.NewMachineRepository(db)
		err = machineRepo.UpdateAvailability(machineID, isAvailable)
		if err != nil {
			http.Error(w, "could not update machine availability", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Machine availability updated successfully"))
	}
}

func GetMachinesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		machineRepo := repository.NewMachineRepository(db)

		machines, err := machineRepo.GetAll()
		if err != nil {
			http.Error(w, "could not fetch machines", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(machines)
	}
}

func GetAvailableMachinesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		machineRepo := repository.NewMachineRepository(db)

		machines, err := machineRepo.GetAvailable()
		if err != nil {
			http.Error(w, "could not fetch available machines", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(machines)
	}
}

type UpdateMachineAvailabilityRequest struct {
	ID        int  `json:"id"`
	Available bool `json:"available"`
}

func UpdateMachineAvailabilityPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req UpdateMachineAvailabilityRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		machineRepo := repository.NewMachineRepository(db)
		err = machineRepo.UpdateAvailability(req.ID, req.Available)
		if err != nil {
			http.Error(w, "could not update machine availability", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Machine availability updated successfully",
		})
	}
}

func GetMachineByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "missing machine id", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "invalid machine id", http.StatusBadRequest)
			return
		}

		machineRepo := repository.NewMachineRepository(db)
		machine, err := machineRepo.GetByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "machine not found", http.StatusNotFound)
				return
			}
			http.Error(w, "could not fetch machine", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(machine)
	}
}