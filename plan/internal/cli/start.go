package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"plan/internal/db"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new plan",
	Long:  `Start a new plan session. Returns the plan ID for subsequent commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		parent, _ := cmd.Flags().GetString("parent")

		var descPtr *string
		if description != "" {
			descPtr = &description
		}

		var parentPtr *string
		if parent != "" {
			parentPtr = &parent
		}

		plan, err := db.CreatePlan(title, descPtr, parentPtr)
		if err != nil {
			return fmt.Errorf("failed to create plan: %w", err)
		}

		output, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(output))

		// Print env hint to stderr so it doesn't interfere with JSON output
		fmt.Fprintf(os.Stderr, "\nexport PLAN_SESSION_ID=%s\n", plan.ID)

		return nil
	},
}

func init() {
	startCmd.Flags().StringP("title", "t", "", "Plan title (required)")
	startCmd.Flags().StringP("description", "d", "", "Plan description")
	startCmd.Flags().StringP("parent", "p", "", "Parent plan ID for sub-agent coordination")
	startCmd.MarkFlagRequired("title")
}
