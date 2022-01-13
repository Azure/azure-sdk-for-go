// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package systemperf

import (
	"context"
	"math"
	"time"

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
		return perf.RunPerfTest(&noOpTestCmd{})
	},
}

type noOpTestCmd struct{}

func (m *noOpTestCmd) GlobalSetup(ctx context.Context) error {
	return nil
}

func (m *noOpTestCmd) GlobalTearDown(ctx context.Context) error {
	return nil
}

func (m *noOpTestCmd) Setup(ctx context.Context) error {
	return nil
}

func (m *noOpTestCmd) Run(ctx context.Context) error {
	return nil
}

func (m *noOpTestCmd) TearDown(ctx context.Context) error {
	return nil
}

func (m *noOpTestCmd) GetMetadata() string {
	return "NoOpTest"
}

var SleepTestCmd = &cobra.Command{
	Use:   "SleepTest",
	Short: "Sleep test for verifying performance framework",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		return perf.RunPerfTest(&sleepPerfTest{})
	},
}

type sleepPerfTest struct {
	// Count is the number of instances of sleepPerfTest running
	Count        int
	secondsPerOp float64
}

func (m *sleepPerfTest) GlobalSetup(ctx context.Context) error {
	m.secondsPerOp = math.Pow(2.0, float64(m.Count))
	return nil
}

func (m *sleepPerfTest) GlobalTearDown(ctx context.Context) error {
	return nil
}

func (m *sleepPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (m *sleepPerfTest) Run(ctx context.Context) error {
	time.Sleep(time.Duration(m.secondsPerOp) * time.Second)
	return nil
}

func (m *sleepPerfTest) TearDown(ctx context.Context) error {
	return nil
}

func (m *sleepPerfTest) GetMetadata() string {
	return "NoOpTest"
}
