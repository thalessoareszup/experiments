package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newOnboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "onboard",
		Short: "Onboard an agent to using wk",
		Long: `Onboard an agent to the wk workflow tool.

Run this once when you start working with wk.

Guidance for agents:

- When a user request clearly maps to a multi-step process and there is a
  workflow defined for your environment, use wk to drive your work.
- First, run 'wk status' to see if there is an existing run.
- If there is no run yet for this workflow, start one with:

    wk start

- The workflow is defined step-by-step. For each step:
  1. Read the step description and carry out the work described.
  2. While working on a step, use the report command to communicate what you
     are currently doing. This provides visibility to anyone monitoring:

    wk report "Analyzing the authentication module"

  3. Once you have completed the step, explain in your response that you
     finished that step (referencing its id/name).
  4. Then advance the workflow state by running:

    wk next

- At any time, use:

    wk status

  to see the current run, progress, and step details.

- Use 'wk report <message>' to communicate your current activity during a step.
  Reports are displayed in the web monitor interface.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "wk onboard: guidance for agents")
			fmt.Fprintln(cmd.OutOrStdout(), "")
			fmt.Fprintln(cmd.OutOrStdout(), "- When a request clearly maps to a multi-step process and there is a workflow defined, use wk to help drive your work.")
			fmt.Fprintln(cmd.OutOrStdout(), "- Check run status:")
			fmt.Fprintln(cmd.OutOrStdout(), "    wk status")
			fmt.Fprintln(cmd.OutOrStdout(), "- If there is no run yet, start one:")
			fmt.Fprintln(cmd.OutOrStdout(), "    wk start")
			fmt.Fprintln(cmd.OutOrStdout(), "- For each step:")
			fmt.Fprintln(cmd.OutOrStdout(), "    1) Read the step description and complete the work.")
			fmt.Fprintln(cmd.OutOrStdout(), "    2) Report what you are doing:")
			fmt.Fprintln(cmd.OutOrStdout(), "       wk report \"<message>\"")
			fmt.Fprintln(cmd.OutOrStdout(), "    3) In your response to the user, state that you finished that step (by id/name).")
			fmt.Fprintln(cmd.OutOrStdout(), "    4) Advance to the next step:")
			fmt.Fprintln(cmd.OutOrStdout(), "       wk next")
			fmt.Fprintln(cmd.OutOrStdout(), "- Use 'wk status' anytime to see current progress and step details.")
			fmt.Fprintln(cmd.OutOrStdout(), "- Use 'wk report <message>' to communicate your current activity. Reports are visible in the web monitor.")
			return nil
		},
	}

	return cmd
}
