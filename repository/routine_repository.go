package repository

import (
	"database/sql"

	"Flow_gym_go_project/models"
)

type RoutineRepository struct {
	DB *sql.DB
}

func NewRoutineRepository(db *sql.DB) *RoutineRepository {
	return &RoutineRepository{DB: db}
}

// GetByUserID fetches all routines for a specific user
func (r *RoutineRepository) GetByUserID(userID int) ([]models.Routine, error) {
	// 1. We look for the user's routines
	queryRoutines := `
		SELECT id, user_id, name
		FROM routines
		WHERE user_id = $1
		ORDER BY id
	`

	rows, err := r.DB.Query(queryRoutines, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routines []models.Routine

	for rows.Next() {
		var routine models.Routine
		if err := rows.Scan(&routine.ID, &routine.UserID, &routine.Name); err != nil {
			return nil, err
		}

		// 2. We look for the exercises in THIS routine, keeping them in the same order
		queryExercises := `
			SELECT re.id, re.routine_id, re.exercise_id, re.exercise_order, 
			       e.id, e.name, e.muscle_group_id
			FROM routine_exercises re
			JOIN exercises e ON re.exercise_id = e.id
			WHERE re.routine_id = $1
			ORDER BY re.exercise_order
		`
		exRows, err := r.DB.Query(queryExercises, routine.ID)
		if err == nil {
			for exRows.Next() {
				var re models.RoutineExercise
				var e models.Exercise
				// We scan the data from the intermediate table and the exercise table
				exRows.Scan(
					&re.ID, &re.RoutineID, &re.ExerciseID, &re.Order, 
					&e.ID, &e.Name, &e.MuscleGroupID,
				)
				re.Exercise = e // We nest the exercise information
				routine.Exercises = append(routine.Exercises, re)
			}
			exRows.Close()
		}

		routines = append(routines, routine)
	}

	return routines, nil
}