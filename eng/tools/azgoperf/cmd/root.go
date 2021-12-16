// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"os"

	azkeys "github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf/cmd/azkeys"
	aztables "github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf/cmd/aztables"
	template "github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf/cmd/template"
	"github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf/internal/perf"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "azgoperf [OPTIONS] [PERFTEST] [LOCAL OPTIONS]",
	Short: "Generates a series of performance tests for different SDKs",
	Long:  `This tool creates a single executable for running performance tests for existing Track 2 Go SDKs`,
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&perf.Duration, "duration", "d", 10, "The duration to run a single performance test for")
	rootCmd.PersistentFlags().IntVarP(&perf.Iterations, "iterations", "i", 3, "The number of iterations to run a single performance test for")
	rootCmd.PersistentFlags().StringVarP(&perf.TestProxy, "testproxy", "x", "", "whether to target http or https proxy (default is neither)")
	rootCmd.PersistentFlags().IntVarP(&perf.TimeoutSeconds, "timeout", "t", 10, "How long to allow an operation to block before cancelling.")
	rootCmd.PersistentFlags().IntVarP(&perf.WarmUp, "warmup", "w", 3, "How long to allow a connection to warm up.")

	if !(perf.TestProxy == "" || perf.TestProxy == "http" || perf.TestProxy == "https") {
		panic(fmt.Errorf("received invalid value for testproxy flag, received %s, expected 'http' or 'https'", perf.TestProxy))
	}
	rootCmd.AddCommand(template.TemplateCmd)
	rootCmd.AddCommand(aztables.AztablesCmd)
	rootCmd.AddCommand(azkeys.AzkeysCmd)
}

// Execute executes the specified command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
