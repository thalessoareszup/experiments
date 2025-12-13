package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"example.com/wk/internal/storage"
	"github.com/spf13/cobra"
)

func newSayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "say <message>",
		Short: "Report current progress",
		Long: `Report what you are currently doing for the active step.

Messages are displayed in the web monitor for visibility.

Example:
  wk say "Analyzing the codebase structure"
  wk say "Implementing the authentication module"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			db, err := storage.Open(dbFile)
			if err != nil {
				return fmt.Errorf("open db: %w", err)
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

			message := strings.Join(args, " ")

			if err := db.AddMessage(ctx, run.ID, step.Index, message); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
