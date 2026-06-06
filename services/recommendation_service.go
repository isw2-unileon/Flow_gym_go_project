package services

import (
	"database/sql"
	"fmt"

	"Flow_gym_go_project/models"
	"Flow_gym_go_project/repository"
)

// GetRecommendation generates an alternative exercise recommendation
// based on the requested exercise and current machine availability.
func GetRecommendation(db *sql.DB, exerciseName string) (*models.Recommendation, error) {

	// Create repository instances.
	exerciseRepo := repository.NewExerciseRepository(db)
	machineRepo := repository.NewMachineRepository(db)

	// Retrieve the exercise requested by the user.
	requestedExercise, err := exerciseRepo.GetByName(exerciseName)
	if err != nil {

		// Return a user-friendly error if the exercise does not exist.
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exercise not found")
		}

		return nil, err
	}

	// Retrieve alternative exercises that belong to the same muscle group.
	alternativeExercises, err := exerciseRepo.GetAlternativesByMuscleGroup(
		requestedExercise.MuscleGroupID,
		requestedExercise.ID,
	)
	if err != nil {
		return nil, err
	}

	// Iterate through the alternative exercises.
	for _, alternativeExercise := range alternativeExercises {

		// Look for an available machine that supports the alternative exercise.
		machine, err := machineRepo.GetAvailableByExerciseID(alternativeExercise.ID)

		// If an available machine is found, build the recommendation.
		if err == nil {
			recommendation := &models.Recommendation{
				RequestedExercise:   requestedExercise.Name,
				RecommendedExercise: alternativeExercise.Name,
				MuscleGroup:         requestedExercise.MuscleGroupName,
				Machine:             machine.Name,
			}

			return recommendation, nil
		}

		// Ignore "no rows" errors and continue searching for alternatives.
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// No suitable alternative exercise with an available machine was found.
	return nil, fmt.Errorf("no available recommendation found for this muscle group")
}