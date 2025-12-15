package models

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

type Plan struct {
	ID          string    `json:"id"`
	ParentID    *string   `json:"parent_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Steps       []Step    `json:"steps,omitempty"`
	Children    []Plan    `json:"children,omitempty"`
}

type Step struct {
	ID          string    `json:"id"`
	PlanID      string    `json:"plan_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Status      Status    `json:"status"`
	StepOrder   int       `json:"step_order"`
	Progress    int       `json:"progress"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
