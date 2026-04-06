package services

import (
	"database/sql"
	"fmt"

	"Flow_gym_go_project/models"
)

func GetRecommendation(db *sql.DB, exerciseName string) (*models.Recommendation, error) {
	var requestedExercise string
	var muscleGroup string
	var requestedExerciseID int

	queryRequested := `
		SELECT e.id, e.name, mg.name
		FROM exercises e
		JOIN muscle_groups mg ON e.muscle_group_id = mg.id
		WHERE e.name = $1
	`

	err := db.QueryRow(queryRequested, exerciseName).Scan(&requestedExerciseID, &requestedExercise, &muscleGroup)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exercise not found")
		}
		return nil, err
	}

	queryRecommendation := `
		SELECT e.name, m.name
		FROM exercises e
		JOIN muscle_groups mg ON e.muscle_group_id = mg.id
		JOIN exercise_machines em ON e.id = em.exercise_id
		JOIN machines m ON em.machine_id = m.id
		WHERE mg.name = $1
		  AND e.id != $2
		  AND m.is_available = true
		LIMIT 1
	`

	var recommendedExercise string
	var machine string

	err = db.QueryRow(queryRecommendation, muscleGroup, requestedExerciseID).Scan(&recommendedExercise, &machine)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no alternative recommendation found")
		}
		return nil, err
	}

	recommendation := &models.Recommendation{
		RequestedExercise:   requestedExercise,
		RecommendedExercise: recommendedExercise,
		MuscleGroup:         muscleGroup,
		Machine:             machine,
	}

	return recommendation, nil
}