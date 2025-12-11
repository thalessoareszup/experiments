package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"example.com/wk/internal/storage"
)

func newStatusCmd() *cobra.Command {
	var waitConfirmation bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show workflow run status and current step (from SQLite state)",
		Long: `Show the latest run status and current step using only the SQLite state DB.

This does not read the YAML file; it only inspects the state stored by 'wk start'.

With --wait-confirmation, blocks until the current step is confirmed (if it requires confirmation).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := storage.Open(workflowFile)
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			run, step, total, err := db.LatestRunWithCurrentStep(ctx)
			if err != nil {
				if errors.Is(err, storage.ErrNoRuns) {
					fmt.Fprintln(cmd.OutOrStdout(), "No runs found. Use 'wk start' to create a new run.")
					return nil
				}
				return err
			}

			// Display current status
			displayRunStatus(cmd, run, step, total)

			// If --wait-confirmation is set, poll until confirmed
			if waitConfirmation {
				if err := waitForConfirmation(ctx, db, run, step); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&waitConfirmation, "wait-confirmation", false, "Block until current step is confirmed (if required)")

	return cmd
}

func displayRunStatus(cmd *cobra.Command, run *storage.Run, step *storage.StepSnapshot, total int) {
	// Header: run summary
	fmt.Fprintf(cmd.OutOrStdout(), "Run #%d   Status: %s   Step: %d/%d\n", run.ID, run.Status, step.Index+1, total)
	fmt.Fprintln(cmd.OutOrStdout(), "")

	// Current step details
	fmt.Fprintf(cmd.OutOrStdout(), "Current step [%s]: %s\n", step.ID, step.Name)

	if step.RequiresConfirmation {
		if step.ConfirmedAt == nil {
			fmt.Fprintln(cmd.OutOrStdout(), "⚠️  This step requires confirmation via web UI")
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "✓ Confirmed at: %s\n", step.ConfirmedAt.Format(time.RFC3339))
		}
	}

	if step.Description != "" {
		fmt.Fprintln(cmd.OutOrStdout(), "---")
		fmt.Fprintln(cmd.OutOrStdout(), step.Description)
	}
}

func waitForConfirmation(ctx context.Context, db *storage.DB, run *storage.Run, step *storage.StepSnapshot) error {
	// If step doesn't require confirmation, return immediately
	if !step.RequiresConfirmation {
		return nil
	}

	// If already confirmed, return immediately
	if step.ConfirmedAt != nil {
		return nil
	}

	// Setup signal handling for Ctrl+C
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Fprintln(os.Stderr, "Waiting for confirmation via web UI... (Press Ctrl+C to cancel)")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("interrupted")
		case <-ticker.C:
			// Check if step has been confirmed
			awaiting, err := db.IsAwaitingConfirmation(ctx, run.ID, step.Index)
			if err != nil {
				return err
			}

			if !awaiting {
				fmt.Fprintln(os.Stderr, "✓ Step confirmed!")
				return nil
			}
		}
	}
}
