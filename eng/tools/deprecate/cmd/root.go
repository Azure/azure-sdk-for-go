// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "deprecate",
	Short: "Used to deprecate package contents.",
	Long:  `This tool will add deprecation comments to the specified package based on the command.`,
}

var verboseFlag bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
}

// Execute executes the specified command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func vprintf(format string, a ...interface{}) {
	if verboseFlag {
		fmt.Printf(format, a...)
	}
}
