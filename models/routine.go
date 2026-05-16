package models

// Routine represents a user's workout routine
type Routine struct {
	ID        int               `json:"id"`
	UserID    int               `json:"user_id"`
	Name      string            `json:"name"`
	Exercises []RoutineExercise `json:"exercises,omitempty"` // List of exercises in the routine
}

// RoutineExercise represents a specific exercise within a routine and its order
type RoutineExercise struct {
	ID         int      `json:"id"`
	RoutineID  int      `json:"routine_id"`
	ExerciseID int      `json:"exercise_id"`
	Order      int      `json:"order"`
	Exercise   Exercise `json:"exercise"` // We nest the Exercise model to include all the information
}