package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"example.com/wk/internal/storage"
	"github.com/spf13/cobra"
)

func newReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report <message>",
		Short: "Report what the agent is currently doing",
		Long: `Report the current status or activity of the agent for the active step.

This command allows agents to communicate their progress during a workflow step.
Reports are stored and displayed in the web monitor interface, providing visibility
into what the agent is doing at any given moment.

Example:
  wk report "Analyzing the codebase structure"
  wk report "Running unit tests for the authentication module"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			db, err := storage.Open(workflowFile)
			if err != nil {
				return fmt.Errorf("open state db: %w", err)
			}
			defer db.Close()

			run, step, _, err := db.LatestRunWithCurrentStep(ctx)
			if err != nil {
				if errors.Is(err, storage.ErrNoRuns) {
					return fmt.Errorf("no active run found; use 'wk start' to begin a workflow")
				}
				return err
			}

			if run.Status == "completed" {
				return fmt.Errorf("workflow run #%d is already completed", run.ID)
			}

			// Join all arguments into a single message
			message := strings.Join(args, " ")

			if err := db.AddReport(ctx, run.ID, step.Index, message); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Report added for step %d (%s): %s\n", step.Index, step.ID, message)
			return nil
		},
	}

	return cmd
}
