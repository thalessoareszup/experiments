package cmd

import (
	"encoding/json"
	"fmt"

	"plan/internal/db"
	"plan/internal/models"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a step or plan as completed",
	Long:  `Mark a step or entire plan as completed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stepID, _ := cmd.Flags().GetString("step")
		planID, _ := cmd.Flags().GetString("plan")

		if stepID == "" && planID == "" {
			return fmt.Errorf("either --step or --plan is required")
		}

		type Result struct {
			Type      string `json:"type"`
			ID        string `json:"id"`
			Status    string `json:"status"`
			UpdatedAt string `json:"updated_at"`
		}

		if stepID != "" {
			if err := db.UpdateStepStatus(stepID, models.StatusCompleted); err != nil {
				return fmt.Errorf("failed to complete step: %w", err)
			}
			step, _ := db.GetStep(stepID)
			result := Result{
				Type:      "step",
				ID:        step.ID,
				Status:    string(step.Status),
				UpdatedAt: step.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			output, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(output))
		}

		if planID != "" {
			if err := db.UpdatePlanStatus(planID, models.StatusCompleted); err != nil {
				return fmt.Errorf("failed to complete plan: %w", err)
			}
			plan, _ := db.GetPlan(planID)
			result := Result{
				Type:      "plan",
				ID:        plan.ID,
				Status:    string(plan.Status),
				UpdatedAt: plan.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			output, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(output))
		}

		return nil
	},
}

func init() {
	completeCmd.Flags().StringP("step", "s", "", "Step ID to complete")
	completeCmd.Flags().StringP("plan", "p", "", "Plan ID to complete")
}
