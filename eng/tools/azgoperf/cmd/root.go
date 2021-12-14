// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "azgoperf [OPTIONS] [PERFTEST] [LOCAL OPTIONS]",
	Short: "Generates a series of performance tests for different SDKs",
	Long:  `This tool creates a single executable for running performance tests for existing Track 2 Go SDKs`,
}

var (
	duration       int
	iterations     int
	TestProxy      string
	timeoutSeconds int
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 10, "The duration to run a single performance test for")
	rootCmd.PersistentFlags().IntVarP(&iterations, "iterations", "i", 3, "The number of iterations to run a single performance test for")
	rootCmd.PersistentFlags().StringVarP(&TestProxy, "testproxy", "x", "", "whether to target http or https proxy (default is neither)")
	rootCmd.PersistentFlags().IntVarP(&timeoutSeconds, "timeout", "t", 10, "How long to allow an operation to block before cancelling.")

	if !(TestProxy == "" || TestProxy == "http" || TestProxy == "https") {
		panic(fmt.Errorf("received invalid value for testproxy flag, received %s, expected 'http' or 'https'", TestProxy))
	}
}

// Execute executes the specified command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
