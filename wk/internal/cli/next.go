package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"example.com/wk/internal/storage"
)

func newNextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Advance to the next workflow step",
		Long: `Advance the current run to the next step.

If the run is already at the last step, it is marked as completed.
If the current step requires confirmation but has not been confirmed, the command fails.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.Open(workflowFile)
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()

			// Check if current step is awaiting confirmation
			run, step, _, err := db.LatestRunWithCurrentStep(ctx)
			if err != nil {
				if errors.Is(err, storage.ErrNoRuns) {
					fmt.Fprintln(cmd.OutOrStdout(), "No runs found. Use 'wk start' to create a new run.")
					return nil
				}
				return err
			}

			// Block if awaiting confirmation
			if step.RequiresConfirmation && step.ConfirmedAt == nil {
				return fmt.Errorf("current step requires confirmation before advancing (use web UI to confirm)")
			}

			// Proceed with advance logic
			run, step, total, err := db.AdvanceLatestRun(ctx)
			if err != nil {
				if errors.Is(err, storage.ErrNoRuns) {
					fmt.Fprintln(cmd.OutOrStdout(), "No runs found. Use 'wk start' to create a new run.")
					return nil
				}
				if errors.Is(err, storage.ErrAlreadyAtLastStep) {
					fmt.Fprintf(cmd.OutOrStdout(), "Run #%d is already at the last step and marked as completed.\n", run.ID)
					return nil
				}
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Run #%d advanced to step %d/%d\n", run.ID, step.Index+1, total)
			fmt.Fprintf(cmd.OutOrStdout(), "Current step [%s]: %s\n", step.ID, step.Name)
			if step.Description != "" {
				fmt.Fprintln(cmd.OutOrStdout(), "---")
				fmt.Fprintln(cmd.OutOrStdout(), step.Description)
			}

			return nil
		},
	}

	return cmd
}
