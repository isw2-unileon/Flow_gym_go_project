package repository

import (
	"database/sql"

	"Flow_gym_go_project/models"
)

type ExerciseRepository struct {
	DB *sql.DB
}

func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{DB: db}
}

func (r *ExerciseRepository) GetByName(name string) (*models.Exercise, error) {
	query := `
		SELECT e.id, e.name, e.muscle_group_id, mg.name
		FROM exercises e
		JOIN muscle_groups mg ON e.muscle_group_id = mg.id
		WHERE e.name = $1
	`

	var exercise models.Exercise
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

func (r *ExerciseRepository) GetAll() ([]models.Exercise, error) {
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

	for rows.Next() {
		var exercise models.Exercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.MuscleGroupID,
		)
		if err != nil {
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (r *ExerciseRepository) GetAlternativeByMuscleGroup(muscleGroupID int, excludedExerciseID int) (*models.Exercise, error) {
	query := `
		SELECT e.id, e.name, e.muscle_group_id
		FROM exercises e
		WHERE e.muscle_group_id = $1
		  AND e.id != $2
		ORDER BY e.id
		LIMIT 1
	`

	var exercise models.Exercise
	err := r.DB.QueryRow(query, muscleGroupID, excludedExerciseID).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.MuscleGroupID,
	)
	if err != nil {
		return nil, err
	}

	return &exercise, nil
}