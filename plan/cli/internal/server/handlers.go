package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"plan/internal/db"
	"plan/internal/events"
	"plan/internal/models"
)

func handlePlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		status := r.URL.Query().Get("status")
		var statusPtr *models.Status
		if status != "" {
			s := models.Status(status)
			statusPtr = &s
		}

		plans, err := db.ListPlans(statusPtr, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get steps for each plan
		for i := range plans {
			steps, _ := db.GetStepsByPlan(plans[i].ID)
			plans[i].Steps = steps
		}

		if plans == nil {
			plans = []models.Plan{}
		}

		json.NewEncoder(w).Encode(plans)

	case "POST":
		var req struct {
			Title       string  `json:"title"`
			Description *string `json:"description"`
			ParentID    *string `json:"parent_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		plan, err := db.CreatePlan(req.Title, req.Description, req.ParentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		events.Emit(events.PlanCreated, plan)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(plan)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePlanByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract plan ID from path: /api/plans/{id} or /api/plans/{id}/steps
	path := strings.TrimPrefix(r.URL.Path, "/api/plans/")
	parts := strings.Split(path, "/")
	planID := parts[0]

	// Handle /api/plans/{id}/steps
	if len(parts) > 1 && parts[1] == "steps" {
		handlePlanSteps(w, r, planID)
		return
	}

	// Handle /api/plans/{id}/tree
	if len(parts) > 1 && parts[1] == "tree" {
		handlePlanTree(w, r, planID)
		return
	}

	switch r.Method {
	case "GET":
		plan, err := db.GetPlanWithSteps(planID)
		if err != nil {
			http.Error(w, "Plan not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(plan)

	case "PATCH":
		var req struct {
			Status *models.Status `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Status != nil {
			if err := db.UpdatePlanStatus(planID, *req.Status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		plan, _ := db.GetPlanWithSteps(planID)
		events.Emit(events.PlanUpdated, plan)
		json.NewEncoder(w).Encode(plan)

	case "DELETE":
		if err := db.DeletePlan(planID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events.Emit(events.PlanDeleted, map[string]string{"id": planID})
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePlanSteps(w http.ResponseWriter, r *http.Request, planID string) {
	switch r.Method {
	case "GET":
		steps, err := db.GetStepsByPlan(planID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if steps == nil {
			steps = []models.Step{}
		}
		json.NewEncoder(w).Encode(steps)

	case "POST":
		var req struct {
			Title       string  `json:"title"`
			Description *string `json:"description"`
			Order       int     `json:"order"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		step, err := db.CreateStep(planID, req.Title, req.Description, req.Order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		events.Emit(events.StepCreated, step)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(step)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePlanTree(w http.ResponseWriter, r *http.Request, planID string) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	plan, err := db.GetPlanWithChildren(planID)
	if err != nil {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(plan)
}

func handleStepByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract step ID from path: /api/steps/{id}
	stepID := strings.TrimPrefix(r.URL.Path, "/api/steps/")

	switch r.Method {
	case "GET":
		step, err := db.GetStep(stepID)
		if err != nil {
			http.Error(w, "Step not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(step)

	case "PATCH":
		var req struct {
			Status   *models.Status `json:"status"`
			Progress *int           `json:"progress"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		step, err := db.UpdateStep(stepID, req.Status, req.Progress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		events.Emit(events.StepUpdated, step)
		json.NewEncoder(w).Encode(step)

	case "DELETE":
		if err := db.DeleteStep(stepID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events.Emit(events.StepDeleted, map[string]string{"id": stepID})
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
