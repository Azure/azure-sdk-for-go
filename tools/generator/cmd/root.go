package cmd

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/automation"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "generator [command]",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := os.Setenv("NODE_OPTIONS", "--max-old-space-size=8192"); err != nil {
				return fmt.Errorf("failed to set environment variable: %v", err)
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(usage)
		},
	}

	// add other command groups
	rootCmd.AddCommand(
		automation.Command(),
	)

	return rootCmd
}

const usage = ``
