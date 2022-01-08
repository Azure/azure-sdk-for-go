// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/automation"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/refresh"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/release"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/template"
	automation_v2 "github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/automation"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/readme"
	refresh_v2 "github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/refresh"
	release_v2 "github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/release"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/spf13/cobra"
)

// Command returns the command for the generator
func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "generator",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := os.Setenv("NODE_OPTIONS", "--max-old-space-size=8192"); err != nil {
				return fmt.Errorf("failed to set environment variable: %v", err)
			}
			log.SetFlags(0) // remove the time stamp prefix
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please specify a subcommand")
		},
		Hidden: true,
	}

	//bind global flags
	common.BindGlobalFlags(rootCmd.PersistentFlags())

	rootCmd.AddCommand(
		automation.Command(),
		issue.Command(),
		template.Command(),
		refresh.Command(),
		release.Command(),
		automation_v2.Command(),
		release_v2.Command(),
		refresh_v2.Command(),
		readme.Command(),
	)

	return rootCmd
}
