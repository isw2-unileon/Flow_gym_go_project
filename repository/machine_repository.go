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

func (r *MachineRepository) GetAvailableByExerciseID(exerciseID int) (*models.Machine, error) {
	query := `
		SELECT m.id, m.name, m.is_available
		FROM machines m
		JOIN exercise_machines em ON m.id = em.machine_id
		WHERE em.exercise_id = $1
		  AND m.is_available = true
		ORDER BY m.id
		LIMIT 1
	`

	var machine models.Machine
	err := r.DB.QueryRow(query, exerciseID).Scan(
		&machine.ID,
		&machine.Name,
		&machine.IsAvailable,
	)
	if err != nil {
		return nil, err
	}

	return &machine, nil
}

func (r *MachineRepository) UpdateAvailability(machineID int, isAvailable bool) error {
	query := `
		UPDATE machines
		SET is_available = $1
		WHERE id = $2
	`

	_, err := r.DB.Exec(query, isAvailable, machineID)
	return err
}

func (r *MachineRepository) GetByID(id int) (*models.Machine, error) {
	query := `
		SELECT id, name, is_available
		FROM machines
		WHERE id = $1
	`

	var machine models.Machine
	err := r.DB.QueryRow(query, id).Scan(
		&machine.ID,
		&machine.Name,
		&machine.IsAvailable,
	)
	if err != nil {
		return nil, err
	}

	return &machine, nil
}