// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"math"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type sleepPerfTest struct {
	perf.PerfTestOptions
	secondsPerOp int
}

func (s *sleepPerfTest) GlobalSetup(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) Run(ctx context.Context) error {
	time.Sleep(time.Duration(s.secondsPerOp) * time.Second)
	return nil
}

func (s *sleepPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) RegisterArguments(ctx context.Context) error {
	return nil
}

func (s *sleepPerfTest) GetMetadata() perf.PerfTestOptions {
	return s.PerfTestOptions
}

func NewSleepTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "SleepTest"
	return &sleepPerfTest{
		PerfTestOptions: *options,
		secondsPerOp:    int(math.Pow(2.0, float64(options.ParallelIndex))),
	}
}
