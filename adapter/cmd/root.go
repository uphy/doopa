package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doopa",
	Short: "doopa is super simple Docker PaaS.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute executes the command.
func Execute() error {
	return rootCmd.Execute()
}
