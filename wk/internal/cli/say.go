package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"example.com/wk/internal/storage"
	"github.com/spf13/cobra"
)

func newSayCmd() *cobra.Command {
	var askFlag bool

	cmd := &cobra.Command{
		Use:   "say <message>",
		Short: "Send a message or ask a question",
		Long: `Send a message to report progress or ask a question.

Without --ask, the message is recorded as a status update and the command returns immediately.
With --ask, the message is recorded as a question and the command waits for a reply from the web UI.

Examples:
  wk say "Working on the authentication module"
  wk say "Should I use JWT or sessions?" --ask`,
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

			msgID, err := db.AddMessage(ctx, run.ID, step.Index, message, askFlag)
			if err != nil {
				return err
			}

			if !askFlag {
				fmt.Fprintf(cmd.OutOrStdout(), "Message sent for step %d (%s)\n", step.Index, step.ID)
				return nil
			}

			// Wait for reply
			fmt.Fprintln(cmd.ErrOrStderr(), "Waiting for reply...")

			ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer cancel()

			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return fmt.Errorf("interrupted")
				case <-ticker.C:
					reply, err := db.GetMessageReply(ctx, msgID)
					if err != nil {
						return err
					}
					if reply != nil {
						fmt.Fprintln(cmd.OutOrStdout(), *reply)
						return nil
					}
				}
			}
		},
	}

	cmd.Flags().BoolVar(&askFlag, "ask", false, "Wait for a reply from the web UI")

	return cmd
}
