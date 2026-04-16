package services

import (
	"database/sql"
	"fmt"

	"Flow_gym_go_project/models"
	"Flow_gym_go_project/repository"
)

func GetRecommendation(db *sql.DB, exerciseName string) (*models.Recommendation, error) {
	exerciseRepo := repository.NewExerciseRepository(db)
	machineRepo := repository.NewMachineRepository(db)

	requestedExercise, err := exerciseRepo.GetByName(exerciseName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exercise not found")
		}
		return nil, err
	}

	alternativeExercises, err := exerciseRepo.GetAlternativesByMuscleGroup(
		requestedExercise.MuscleGroupID,
		requestedExercise.ID,
	)
	if err != nil {
		return nil, err
	}

	for _, alternativeExercise := range alternativeExercises {
		machine, err := machineRepo.GetAvailableByExerciseID(alternativeExercise.ID)
		if err == nil {
			recommendation := &models.Recommendation{
				RequestedExercise:   requestedExercise.Name,
				RecommendedExercise: alternativeExercise.Name,
				MuscleGroup:         requestedExercise.MuscleGroupName,
				Machine:             machine.Name,
			}
			return recommendation, nil
		}

		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return nil, fmt.Errorf("no available recommendation found for this muscle group")
}