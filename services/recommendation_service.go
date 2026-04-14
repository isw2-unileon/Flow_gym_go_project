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

	alternativeExercise, err := exerciseRepo.GetAlternativeByMuscleGroup(
		requestedExercise.MuscleGroupID,
		requestedExercise.ID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no alternative recommendation found")
		}
		return nil, err
	}

	machine, err := machineRepo.GetAvailableByExerciseID(alternativeExercise.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no available machine found for alternative exercise")
		}
		return nil, err
	}

	recommendation := &models.Recommendation{
		RequestedExercise:   requestedExercise.Name,
		RecommendedExercise: alternativeExercise.Name,
		MuscleGroup:         requestedExercise.MuscleGroupName,
		Machine:             machine.Name,
	}

	return recommendation, nil
}