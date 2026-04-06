package repository

import (
	"database/sql"

	"Flow_gym_go_project/models"
)

type MachineRepository struct {
	DB *sql.DB
}

func NewMachineRepository(db *sql.DB) *MachineRepository {
	return &MachineRepository{DB: db}
}

func (r *MachineRepository) GetAll() ([]models.Machine, error) {
	query := `
		SELECT id, name, is_available
		FROM machines
		ORDER BY id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []models.Machine

	for rows.Next() {
		var machine models.Machine
		err := rows.Scan(
			&machine.ID,
			&machine.Name,
			&machine.IsAvailable,
		)
		if err != nil {
			return nil, err
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

func (r *MachineRepository) GetAvailable() ([]models.Machine, error) {
	query := `
		SELECT id, name, is_available
		FROM machines
		WHERE is_available = true
		ORDER BY id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []models.Machine

	for rows.Next() {
		var machine models.Machine
		err := rows.Scan(
			&machine.ID,
			&machine.Name,
			&machine.IsAvailable,
		)
		if err != nil {
			return nil, err
		}

		machines = append(machines, machine)
	}

	return machines, nil
}