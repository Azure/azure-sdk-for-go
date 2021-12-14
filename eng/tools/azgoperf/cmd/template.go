// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "template perf test",
	Long:  "template perf test longer description",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		return RunPerfTest(&templatePerfTest{})
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}

type templatePerfTest struct {
	counter int
}

func (m *templatePerfTest) GlobalSetup(ctx context.Context) error {
	return nil
}

func (m *templatePerfTest) GlobalTearDown(ctx context.Context) error {
	return nil
}

func (m *templatePerfTest) Setup(ctx context.Context) error {
	m.counter = 0
	return nil
}

func (m *templatePerfTest) Run(ctx context.Context) error {
	m.counter += 1
	return nil
}

func (m *templatePerfTest) TearDown(ctx context.Context) error {
	m.counter = 0
	return nil
}

func (m *templatePerfTest) GetMetadata() string {
	return "template"
}
