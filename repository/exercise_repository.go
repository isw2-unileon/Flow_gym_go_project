package repository

import (
	"database/sql"

	"Flow_gym_go_project/models"
)

// ExerciseRepository provides database access methods for exercises.
type ExerciseRepository struct {
	DB *sql.DB
}

// NewExerciseRepository creates a new ExerciseRepository instance.
func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{DB: db}
}

// GetByName retrieves a single exercise by its name, including
// the associated muscle group name.
func (r *ExerciseRepository) GetByName(name string) (*models.Exercise, error) {

	// Query that joins exercises with muscle groups in order
	// to return additional muscle group information.
	query := `
		SELECT e.id, e.name, e.muscle_group_id, mg.name
		FROM exercises e
		JOIN muscle_groups mg ON e.muscle_group_id = mg.id
		WHERE e.name = $1
	`

	var exercise models.Exercise

	// Execute the query and map the result into the Exercise model.
	err := r.DB.QueryRow(query, name).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.MuscleGroupID,
		&exercise.MuscleGroupName,
	)
	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

// GetAll retrieves every exercise stored in the database.
func (r *ExerciseRepository) GetAll() ([]models.Exercise, error) {

	// Query all exercises ordered by ID.
	query := `
		SELECT id, name, muscle_group_id
		FROM exercises
		ORDER BY id
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise

	// Iterate through all returned rows.
	for rows.Next() {

		var exercise models.Exercise

		// Map each database row into an Exercise model.
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.MuscleGroupID,
		)
		if err != nil {
			return nil, err
		}

		// Add the exercise to the result slice.
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

// GetAlternativesByMuscleGroup returns alternative exercises
// that belong to the same muscle group while excluding
// the originally requested exercise.
func (r *ExerciseRepository) GetAlternativesByMuscleGroup(muscleGroupID int, excludedExerciseID int) ([]models.Exercise, error) {

	// Query all exercises from the same muscle group except
	// the one specified by excludedExerciseID.
	query := `
		SELECT e.id, e.name, e.muscle_group_id
		FROM exercises e
		WHERE e.muscle_group_id = $1
		  AND e.id != $2
		ORDER BY e.id
	`

	rows, err := r.DB.Query(query, muscleGroupID, excludedExerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise

	// Iterate through all matching exercises.
	for rows.Next() {

		var exercise models.Exercise

		// Map the database row into an Exercise model.
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.MuscleGroupID,
		)
		if err != nil {
			return nil, err
		}

		// Add the exercise to the alternatives list.
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}