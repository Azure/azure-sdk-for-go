// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"math"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type globalSleepPerfTest struct {
	perf.PerfTestOptions
	count int
}

func NewSleepTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	return &globalSleepPerfTest{PerfTestOptions: options}, nil
}

func (g *globalSleepPerfTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

type sleepPerfTest struct {
	sleepInterval int
}

func (g *globalSleepPerfTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	s := &sleepPerfTest{sleepInterval: int(math.Pow(2.0, float64(g.count)))}
	g.count += 1
	return s, nil
}

func (s *sleepPerfTest) Run(ctx context.Context) error {
	time.Sleep(time.Duration(s.sleepInterval) * time.Second)
	return nil
}

func (s *sleepPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
