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
	query := `
		SELECT id, user_id, name
		FROM routines
		WHERE user_id = $1
		ORDER BY id
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routines []models.Routine

	for rows.Next() {
		var routine models.Routine
		err := rows.Scan(
			&routine.ID,
			&routine.UserID,
			&routine.Name,
		)
		if err != nil {
			return nil, err
		}

		routines = append(routines, routine)
	}

	return routines, nil
}