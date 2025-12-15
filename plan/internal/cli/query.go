package cli

import (
	"encoding/json"
	"fmt"

	"plan/internal/db"
	"plan/internal/models"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query plans and steps",
	Long:  `Query plans and steps with optional filters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		planID, _ := cmd.Flags().GetString("plan")
		statusStr, _ := cmd.Flags().GetString("status")
		children, _ := cmd.Flags().GetBool("children")
		limit, _ := cmd.Flags().GetInt("limit")

		// If specific plan ID is provided, return that plan
		if planID != "" {
			var plan *models.Plan
			var err error

			if children {
				plan, err = db.GetPlanWithChildren(planID)
			} else {
				plan, err = db.GetPlanWithSteps(planID)
			}

			if err != nil {
				return fmt.Errorf("failed to get plan: %w", err)
			}

			output, _ := json.MarshalIndent(plan, "", "  ")
			fmt.Println(string(output))
			return nil
		}

		// Otherwise, list plans
		var status *models.Status
		if statusStr != "" {
			s := models.Status(statusStr)
			status = &s
		}

		plans, err := db.ListPlans(status, limit)
		if err != nil {
			return fmt.Errorf("failed to list plans: %w", err)
		}

		// Optionally fetch steps for each plan
		for i := range plans {
			steps, _ := db.GetStepsByPlan(plans[i].ID)
			plans[i].Steps = steps
		}

		type Result struct {
			Plans []models.Plan `json:"plans"`
		}

		result := Result{Plans: plans}
		if result.Plans == nil {
			result.Plans = []models.Plan{}
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))

		return nil
	},
}

func init() {
	queryCmd.Flags().StringP("plan", "p", "", "Get specific plan with its steps")
	queryCmd.Flags().StringP("status", "s", "", "Filter by status: pending|in_progress|completed|failed")
	queryCmd.Flags().Bool("children", false, "Include child plans in output")
	queryCmd.Flags().IntP("limit", "l", 50, "Limit results")
}
