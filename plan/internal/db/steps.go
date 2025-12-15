package db

import (
	"database/sql"
	"time"

	"plan/internal/models"

	"github.com/google/uuid"
)

func CreateStep(planID, title string, description *string, order int) (*models.Step, error) {
	// If order is 0, get the next order number
	if order == 0 {
		var maxOrder sql.NullInt64
		DB.QueryRow(`SELECT MAX(step_order) FROM steps WHERE plan_id = ?`, planID).Scan(&maxOrder)
		if maxOrder.Valid {
			order = int(maxOrder.Int64) + 1
		} else {
			order = 1
		}
	}

	step := &models.Step{
		ID:          uuid.New().String(),
		PlanID:      planID,
		Title:       title,
		Description: description,
		Status:      models.StatusPending,
		StepOrder:   order,
		Progress:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := DB.Exec(`
		INSERT INTO steps (id, plan_id, title, description, status, step_order, progress, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, step.ID, step.PlanID, step.Title, step.Description, step.Status, step.StepOrder, step.Progress, step.CreatedAt, step.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return step, nil
}

func GetStep(id string) (*models.Step, error) {
	step := &models.Step{}
	var description sql.NullString

	err := DB.QueryRow(`
		SELECT id, plan_id, title, description, status, step_order, progress, created_at, updated_at
		FROM steps WHERE id = ?
	`, id).Scan(&step.ID, &step.PlanID, &step.Title, &description, &step.Status, &step.StepOrder, &step.Progress, &step.CreatedAt, &step.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if description.Valid {
		step.Description = &description.String
	}

	return step, nil
}

func GetStepsByPlan(planID string) ([]models.Step, error) {
	rows, err := DB.Query(`
		SELECT id, plan_id, title, description, status, step_order, progress, created_at, updated_at
		FROM steps WHERE plan_id = ?
		ORDER BY step_order ASC
	`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []models.Step
	for rows.Next() {
		var step models.Step
		var description sql.NullString

		if err := rows.Scan(&step.ID, &step.PlanID, &step.Title, &description, &step.Status, &step.StepOrder, &step.Progress, &step.CreatedAt, &step.UpdatedAt); err != nil {
			return nil, err
		}

		if description.Valid {
			step.Description = &description.String
		}

		steps = append(steps, step)
	}

	return steps, nil
}

func UpdateStepStatus(id string, status models.Status) error {
	_, err := DB.Exec(`
		UPDATE steps SET status = ?, updated_at = ? WHERE id = ?
	`, status, time.Now(), id)
	return err
}

func UpdateStepProgress(id string, progress int) error {
	_, err := DB.Exec(`
		UPDATE steps SET progress = ?, updated_at = ? WHERE id = ?
	`, progress, time.Now(), id)
	return err
}

func UpdateStep(id string, status *models.Status, progress *int) (*models.Step, error) {
	if status != nil {
		if err := UpdateStepStatus(id, *status); err != nil {
			return nil, err
		}
	}

	if progress != nil {
		if err := UpdateStepProgress(id, *progress); err != nil {
			return nil, err
		}
	}

	return GetStep(id)
}

func DeleteStep(id string) error {
	_, err := DB.Exec(`DELETE FROM steps WHERE id = ?`, id)
	return err
}
