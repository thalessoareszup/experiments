package cmd

import (
	"fmt"
	"os"

	"plan/internal/db"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan tracking tool for AI agents",
	Long:  `A CLI tool for AI agents to track and coordinate multi-step workflows.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip db init for help commands
		if cmd.Name() == "help" || cmd.Name() == "version" {
			return nil
		}
		return db.Init()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.Close()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stepCmd)
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(failCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(serveCmd)
}
