package models

type Exercise struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	MuscleGroupID int    `json:"muscle_group_id"`
}