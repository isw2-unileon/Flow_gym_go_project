package models

import "time"

type Machine struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	IsAvailable        bool       `json:"is_available"`
	OccupiedByUserID   *int       `json:"occupied_by_user_id"`
	LastUsedByUserID   *int       `json:"last_used_by_user_id"`
	LastReleasedAt     *time.Time `json:"last_released_at"`
	OccupiedUntil      *time.Time `json:"occupied_until"`
}