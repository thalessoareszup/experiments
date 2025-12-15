package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"plan/internal/db"

	"github.com/spf13/cobra"
)

var stepCmd = &cobra.Command{
	Use:   "step",
	Short: "Add a step to a plan",
	Long:  `Add a new step to an existing plan.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		planID, _ := cmd.Flags().GetString("plan")
		order, _ := cmd.Flags().GetInt("order")

		// Use env var if plan ID not provided
		if planID == "" {
			planID = os.Getenv("PLAN_SESSION_ID")
		}

		if planID == "" {
			return fmt.Errorf("plan ID required: use --plan flag or set PLAN_SESSION_ID")
		}

		var descPtr *string
		if description != "" {
			descPtr = &description
		}

		step, err := db.CreateStep(planID, title, descPtr, order)
		if err != nil {
			return fmt.Errorf("failed to create step: %w", err)
		}

		output, _ := json.MarshalIndent(step, "", "  ")
		fmt.Println(string(output))

		return nil
	},
}

func init() {
	stepCmd.Flags().StringP("title", "t", "", "Step title (required)")
	stepCmd.Flags().StringP("description", "d", "", "Step description")
	stepCmd.Flags().StringP("plan", "p", "", "Plan ID (uses PLAN_SESSION_ID if not set)")
	stepCmd.Flags().IntP("order", "o", 0, "Step order (auto-increments if not set)")
	stepCmd.MarkFlagRequired("title")
}
