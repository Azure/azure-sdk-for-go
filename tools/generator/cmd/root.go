// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/automation"
	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/issue"
	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/template"
	"github.com/spf13/cobra"
)

// Command returns the command for the generator
func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "generator",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand")
		},
		Hidden: true,
	}

	rootCmd.AddCommand(
		automation.Command(),
		issue.Command(),
		template.Command(),
	)

	return rootCmd
}
