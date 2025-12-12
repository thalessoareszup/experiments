package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"example.com/wk/internal/storage"
	"example.com/wk/internal/workflow"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a workflow",
		Long: `Load the workflow definition and create a new run.

By default, the workflow file is $HOME/.local/wk/workflow.yaml.
You can override this with -f/--file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load workflow definition
			wf, err := workflow.Load(workflowFile)
			if err != nil {
				return err
			}

			// Open state DB
			db, err := storage.Open(workflowFile)
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			run, err := db.StartRun(ctx, wf)
			if err != nil {
				return err
			}

			step := wf.Steps[run.CurrentStepIndex]
			fmt.Fprintf(cmd.OutOrStdout(), "Started workflow '%s' (run #%d)\n", wf.Name, run.ID)
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
