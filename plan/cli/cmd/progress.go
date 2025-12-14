package cmd

import (
	"encoding/json"
	"fmt"

	"plan/internal/db"
	"plan/internal/models"

	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Update progress on a step",
	Long:  `Update the progress percentage and/or status of a step.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stepID, _ := cmd.Flags().GetString("step")
		percent, _ := cmd.Flags().GetInt("percent")
		statusStr, _ := cmd.Flags().GetString("status")

		if stepID == "" {
			return fmt.Errorf("step ID required: use --step flag")
		}

		var status *models.Status
		if statusStr != "" {
			s := models.Status(statusStr)
			status = &s
		}

		var progress *int
		if cmd.Flags().Changed("percent") {
			progress = &percent
		}

		step, err := db.UpdateStep(stepID, status, progress)
		if err != nil {
			return fmt.Errorf("failed to update step: %w", err)
		}

		output, _ := json.MarshalIndent(step, "", "  ")
		fmt.Println(string(output))

		return nil
	},
}

func init() {
	progressCmd.Flags().StringP("step", "s", "", "Step ID (required)")
	progressCmd.Flags().IntP("percent", "p", 0, "Progress percentage (0-100)")
	progressCmd.Flags().String("status", "", "New status: pending|in_progress|completed|failed")
	progressCmd.MarkFlagRequired("step")
}
