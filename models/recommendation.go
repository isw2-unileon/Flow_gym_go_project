package models

type Recommendation struct {
	RequestedExercise   string `json:"requested_exercise"`
	RecommendedExercise string `json:"recommended_exercise"`
	MuscleGroup         string `json:"muscle_group"`
	Machine             string `json:"machine"`
}