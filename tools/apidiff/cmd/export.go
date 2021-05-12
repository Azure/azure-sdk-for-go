// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/internal/repo"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export <package searching dir>",
	Short: "Generates a report of all the export content for every package under the specified directory.",
	Long: `The export command generates a report of all the export content for all of the packages under the directory 
specified in <package searching dir>.

The export content will be categorized by consts, funcs, interfaces, structs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return exportCommand(args[0])
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func exportCommand(searchingDir string) error {
	wt, err := repo.Get(searchingDir)
	if err != nil {
		return fmt.Errorf("failed to get repo: %+v", err)
	}
	r, err := getRepoContent(&wt, searchingDir)
	if err != nil {
		return err
	}
	return r.Print(os.Stdout)
}
