package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	workflowFile string
	dbFile       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wk",
	Short: "wk helps agents follow a step-by-step workflow",
	Long: `wk is a CLI tool for agents to follow a workflow.
Use it to start and progress through the steps of your workflow.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		workflowFile = homeDir + "/.local/wk/workflow.yaml"
		dbFile = homeDir + "/.local/wk/wk.db"
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&workflowFile, "file", "f", "", "Path to workflow YAML file (default: $HOME/.local/wk/workflow.yaml)")

	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newStartCmd())
	rootCmd.AddCommand(newNextCmd())
	rootCmd.AddCommand(newOnboardCmd())
	rootCmd.AddCommand(newWebCmd())
	rootCmd.AddCommand(newReportCmd())
}
