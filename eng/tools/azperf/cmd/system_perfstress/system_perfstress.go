// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package systemperf

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/eng/tools/azperf/internal/perf"
	"github.com/spf13/cobra"
)

var NoOpTestCmd = &cobra.Command{
	Use:   "NoOpTest",
	Short: "No op test for verifying performance framework",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		return perf.RunPerfTest(&templatePerfTest{})
	},
}

type templatePerfTest struct {}

func (m *templatePerfTest) GlobalSetup(ctx context.Context) error {
	return nil
}

func (m *templatePerfTest) GlobalTearDown(ctx context.Context) error {
	return nil
}

func (m *templatePerfTest) Setup(ctx context.Context) error {
	return nil
}

func (m *templatePerfTest) Run(ctx context.Context) error {
	return nil
}

func (m *templatePerfTest) TearDown(ctx context.Context) error {
	return nil
}

func (m *templatePerfTest) GetMetadata() string {
	return "NoOpTest"
}
