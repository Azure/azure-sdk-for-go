package cmd

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "snippet-generator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute()
		},
	}

	rootCmd.PersistentFlags().String("base-directory", ".", "The base directory.")
	rootCmd.PersistentFlags().Bool("strict-mode", true, "Are we running in strict mode?")

	return rootCmd
}

func execute() error {
	return nil
}
