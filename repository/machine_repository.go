package repository

import (
	"database/sql"
	"errors"

	"Flow_gym_go_project/models"
)

// ErrUserAlreadyOccupiesMachine is returned when a normal user tries to occupy
// a second machine while already occupying another one.
var ErrUserAlreadyOccupiesMachine = errors.New("user already occupies another machine")

// MachineRepository provides database access methods for machines.
type MachineRepository struct {
	DB *sql.DB
}

// NewMachineRepository creates a new MachineRepository instance.
func NewMachineRepository(db *sql.DB) *MachineRepository {
	return &MachineRepository{DB: db}
}

// GetAll retrieves all machines from the database.
// Before returning the machines, it releases any machine whose occupation time has expired.
func (r *MachineRepository) GetAll() ([]models.Machine, error) {

	// Release expired machine occupations before reading the current state.
	err := r.ReleaseExpiredMachines()
	if err != nil {
		return nil, err
	}

	// Retrieve all machine data, including occupation and cooldown metadata.
	query := `
		SELECT id, name, is_available, occupied_by_user_id, last_used_by_user_id, last_released_at, occupied_until
		FROM machines
		ORDER BY id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []models.Machine

	// Convert each database row into a Machine model.
	for rows.Next() {
		var machine models.Machine
		err := rows.Scan(
			&machine.ID,
			&machine.Name,
			&machine.IsAvailable,
			&machine.OccupiedByUserID,
			&machine.LastUsedByUserID,
			&machine.LastReleasedAt,
			&machine.OccupiedUntil,
		)
		if err != nil {
			return nil, err
		}

		machines = append(machines, machine)
	}

	return machines, nil
}

// GetAvailable retrieves only machines currently marked as available.
func (r *MachineRepository) GetAvailable() ([]models.Machine, error) {

	// Query machines that are available.
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

	// Convert each row into a Machine model.
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

// GetAvailableByExerciseID returns the first available machine that can be used
// for a specific exercise.
func (r *MachineRepository) GetAvailableByExerciseID(exerciseID int) (*models.Machine, error) {

	// Join machines with exercise_machines to find machines linked to the exercise.
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

	// Retrieve a single available machine.
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

// UpdateAvailability updates a machine availability state directly.
// This method does not apply user ownership, cooldown, or role rules.
func (r *MachineRepository) UpdateAvailability(machineID int, isAvailable bool) error {

	// Simple availability update used by the legacy endpoint.
	query := `
		UPDATE machines
		SET is_available = $1
		WHERE id = $2
	`

	_, err := r.DB.Exec(query, isAvailable, machineID)
	return err
}

// GetByID retrieves a single machine by its ID.
func (r *MachineRepository) GetByID(id int) (*models.Machine, error) {

	// Retrieve machine data including occupation and cooldown metadata.
	query := `
		SELECT id, name, is_available, occupied_by_user_id, last_used_by_user_id, last_released_at, occupied_until
		FROM machines
		WHERE id = $1
	`

	var machine models.Machine

	// Map the database row into a Machine model.
	err := r.DB.QueryRow(query, id).Scan(
		&machine.ID,
		&machine.Name,
		&machine.IsAvailable,
		&machine.OccupiedByUserID,
		&machine.LastUsedByUserID,
		&machine.LastReleasedAt,
		&machine.OccupiedUntil,
	)
	if err != nil {
		return nil, err
	}

	return &machine, nil
}

// ReleaseExpiredMachines automatically releases machines whose occupation time has expired.
func (r *MachineRepository) ReleaseExpiredMachines() error {

	// Machines with an expired occupied_until timestamp become available again.
	// The last user and release time are stored for cooldown purposes.
	query := `
		UPDATE machines
		SET
			is_available = true,
			last_used_by_user_id = occupied_by_user_id,
			last_released_at = timezone('utc', now()),
			occupied_by_user_id = NULL,
			occupied_until = NULL
		WHERE occupied_until IS NOT NULL
		  AND occupied_until < timezone('utc', now())
	`

	_, err := r.DB.Exec(query)
	return err
}

// UpdateAvailabilityWithUser updates machine availability while applying user permissions,
// admin privileges, machine ownership, reservation time, and cooldown rules.
func (r *MachineRepository) UpdateAvailabilityWithUser(machineID int, userID int, isAvailable bool, userRole string) error {

	// Always release expired machines before applying a new update.
	err := r.ReleaseExpiredMachines()
	if err != nil {
		return err
	}

	// If isAvailable is true, the user is trying to release the machine.
	if isAvailable {

		// A normal user can only release a machine occupied by themselves.
		// An admin can release any occupied machine.
		query := `
			UPDATE machines
			SET
				is_available = true,
				last_used_by_user_id = $2,
				last_released_at = timezone('utc', now()),
				occupied_by_user_id = NULL,
				occupied_until = NULL
			WHERE id = $1
			  AND (
			  	occupied_by_user_id = $2
			  	OR $3 = 'admin'
			  )
		`

		result, err := r.DB.Exec(query, machineID, userID, userRole)
		if err != nil {
			return err
		}

		// Check whether the update actually affected a row.
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		// If no row was affected, the user was not allowed to release the machine.
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}

		return nil
	}

	// Normal users cannot occupy more than one machine at the same time.
	if userRole != "admin" {
		var occupiedMachineID int

		// Look for another machine already occupied by the same user.
		err := r.DB.QueryRow(`
			SELECT id
			FROM machines
			WHERE occupied_by_user_id = $1
			  AND is_available = false
			  AND id != $2
			LIMIT 1
		`, userID, machineID).Scan(&occupiedMachineID)

		// If a row is found, the user is already occupying another machine.
		if err == nil {
			return ErrUserAlreadyOccupiesMachine
		}

		// Any error other than no rows is treated as a database error.
		if err != sql.ErrNoRows {
			return err
		}
	}

	// Occupy the machine for 15 minutes.
	// Normal users must respect the cooldown period before reusing the same machine.
	// Admin users bypass the cooldown condition.
	query := `
		UPDATE machines
		SET
			is_available = false,
			occupied_by_user_id = $2,
			occupied_until = timezone('utc', now()) + INTERVAL '15 minutes'
		WHERE id = $1
		AND is_available = true
		AND (
			$3 = 'admin'
			OR last_used_by_user_id IS NULL
			OR last_used_by_user_id != $2
			OR last_released_at IS NULL
			OR last_released_at <= timezone('utc', now()) - INTERVAL '10 seconds'
		)
	`

	result, err := r.DB.Exec(query, machineID, userID, userRole)
	if err != nil {
		return err
	}

	// Check whether the machine was actually updated.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no row was affected, the machine was unavailable, blocked by cooldown,
	// or could not be updated under the current rules.
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}