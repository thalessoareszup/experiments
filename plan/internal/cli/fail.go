package cli

import (
	"encoding/json"
	"fmt"

	"plan/internal/db"
	"plan/internal/models"

	"github.com/spf13/cobra"
)

var failCmd = &cobra.Command{
	Use:   "fail",
	Short: "Mark a step or plan as failed",
	Long:  `Mark a step or entire plan as failed, optionally with a reason.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stepID, _ := cmd.Flags().GetString("step")
		planID, _ := cmd.Flags().GetString("plan")
		reason, _ := cmd.Flags().GetString("reason")

		if stepID == "" && planID == "" {
			return fmt.Errorf("either --step or --plan is required")
		}

		type Result struct {
			Type      string  `json:"type"`
			ID        string  `json:"id"`
			Status    string  `json:"status"`
			Reason    *string `json:"reason,omitempty"`
			UpdatedAt string  `json:"updated_at"`
		}

		var reasonPtr *string
		if reason != "" {
			reasonPtr = &reason
		}

		if stepID != "" {
			if err := db.UpdateStepStatus(stepID, models.StatusFailed); err != nil {
				return fmt.Errorf("failed to mark step as failed: %w", err)
			}
			step, _ := db.GetStep(stepID)
			result := Result{
				Type:      "step",
				ID:        step.ID,
				Status:    string(step.Status),
				Reason:    reasonPtr,
				UpdatedAt: step.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			output, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(output))
		}

		if planID != "" {
			if err := db.UpdatePlanStatus(planID, models.StatusFailed); err != nil {
				return fmt.Errorf("failed to mark plan as failed: %w", err)
			}
			plan, _ := db.GetPlan(planID)
			result := Result{
				Type:      "plan",
				ID:        plan.ID,
				Status:    string(plan.Status),
				Reason:    reasonPtr,
				UpdatedAt: plan.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
			output, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(output))
		}

		return nil
	},
}

func init() {
	failCmd.Flags().StringP("step", "s", "", "Step ID to mark as failed")
	failCmd.Flags().StringP("plan", "p", "", "Plan ID to mark as failed")
	failCmd.Flags().StringP("reason", "r", "", "Failure reason")
}
