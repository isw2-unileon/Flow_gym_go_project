package models

type Machine struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
}