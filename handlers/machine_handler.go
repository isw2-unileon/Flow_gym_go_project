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