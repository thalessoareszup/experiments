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
  2. Use 'wk say' to communicate progress:

    wk say "Analyzing the authentication module"

  3. Once you have completed the step, advance the workflow:

    wk next

- IMPORTANT: Some steps are marked with 'requires-confirmation'. When you reach
  such a step, you MUST stop and ask the user to review your work before
  proceeding. Explain what you've done and what should be reviewed. Wait for
  explicit approval before calling 'wk next'.

- At any time, use 'wk status' to see current progress and step details.`,
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
			fmt.Fprintln(cmd.OutOrStdout(), "    2) Communicate progress:")
			fmt.Fprintln(cmd.OutOrStdout(), "       wk say \"<message>\"")
			fmt.Fprintln(cmd.OutOrStdout(), "    3) Advance to the next step:")
			fmt.Fprintln(cmd.OutOrStdout(), "       wk next")
			fmt.Fprintln(cmd.OutOrStdout(), "- IMPORTANT: If a step has 'requires-confirmation', stop and ask the user")
			fmt.Fprintln(cmd.OutOrStdout(), "  to review your work before proceeding. Explain what should be reviewed.")
			fmt.Fprintln(cmd.OutOrStdout(), "- Use 'wk status' anytime to see current progress and step details.")
			return nil
		},
	}

	return cmd
}
