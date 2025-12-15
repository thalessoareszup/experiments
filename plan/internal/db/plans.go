package db

import (
	"database/sql"
	"time"

	"plan/internal/models"

	"github.com/google/uuid"
)

func CreatePlan(title string, description *string, parentID *string) (*models.Plan, error) {
	plan := &models.Plan{
		ID:          uuid.New().String(),
		ParentID:    parentID,
		Title:       title,
		Description: description,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := DB.Exec(`
		INSERT INTO plans (id, parent_id, title, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, plan.ID, plan.ParentID, plan.Title, plan.Description, plan.Status, plan.CreatedAt, plan.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return plan, nil
}

func GetPlan(id string) (*models.Plan, error) {
	plan := &models.Plan{}
	var parentID, description sql.NullString

	err := DB.QueryRow(`
		SELECT id, parent_id, title, description, status, created_at, updated_at
		FROM plans WHERE id = ?
	`, id).Scan(&plan.ID, &parentID, &plan.Title, &description, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		plan.ParentID = &parentID.String
	}
	if description.Valid {
		plan.Description = &description.String
	}

	return plan, nil
}

func GetPlanWithSteps(id string) (*models.Plan, error) {
	plan, err := GetPlan(id)
	if err != nil {
		return nil, err
	}

	steps, err := GetStepsByPlan(id)
	if err != nil {
		return nil, err
	}
	plan.Steps = steps

	return plan, nil
}

func GetPlanWithChildren(id string) (*models.Plan, error) {
	plan, err := GetPlanWithSteps(id)
	if err != nil {
		return nil, err
	}

	children, err := GetChildPlans(id)
	if err != nil {
		return nil, err
	}
	plan.Children = children

	return plan, nil
}

func GetChildPlans(parentID string) ([]models.Plan, error) {
	rows, err := DB.Query(`
		SELECT id, parent_id, title, description, status, created_at, updated_at
		FROM plans WHERE parent_id = ?
		ORDER BY created_at ASC
	`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var plan models.Plan
		var parentID, description sql.NullString

		if err := rows.Scan(&plan.ID, &parentID, &plan.Title, &description, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, err
		}

		if parentID.Valid {
			plan.ParentID = &parentID.String
		}
		if description.Valid {
			plan.Description = &description.String
		}

		// Recursively get children
		children, err := GetChildPlans(plan.ID)
		if err != nil {
			return nil, err
		}
		plan.Children = children

		// Get steps
		steps, err := GetStepsByPlan(plan.ID)
		if err != nil {
			return nil, err
		}
		plan.Steps = steps

		plans = append(plans, plan)
	}

	return plans, nil
}

func ListPlans(status *models.Status, limit int) ([]models.Plan, error) {
	query := `
		SELECT id, parent_id, title, description, status, created_at, updated_at
		FROM plans
	`
	var args []interface{}

	if status != nil {
		query += " WHERE status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY updated_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var plan models.Plan
		var parentID, description sql.NullString

		if err := rows.Scan(&plan.ID, &parentID, &plan.Title, &description, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, err
		}

		if parentID.Valid {
			plan.ParentID = &parentID.String
		}
		if description.Valid {
			plan.Description = &description.String
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

func UpdatePlanStatus(id string, status models.Status) error {
	_, err := DB.Exec(`
		UPDATE plans SET status = ?, updated_at = ? WHERE id = ?
	`, status, time.Now(), id)
	return err
}

func DeletePlan(id string) error {
	_, err := DB.Exec(`DELETE FROM plans WHERE id = ?`, id)
	return err
}
