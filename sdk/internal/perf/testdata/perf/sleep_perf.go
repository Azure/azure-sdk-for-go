// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"flag"
	"math"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type sleepTestOptions struct {
	initialDelayMs        int
	instanceGrowthFactor  float64
	iterationGrowthFactor float64
}

var sleepTestOpts sleepTestOptions = sleepTestOptions{initialDelayMs: 1000, instanceGrowthFactor: 1, iterationGrowthFactor: 1}

// sleepTestRegister is called once per process
func sleepTestRegister() {
	flag.IntVar(&sleepTestOpts.initialDelayMs, "initial-delay-ms", 1000, "Initial delay (in milliseconds)")

	// Used for verifying the perf framework correctly computes average throughput across parallel tests of different speed.
	// Each instance of this test completes operations at a different rate, to allow for testing scenarios where
	// some instances are still waiting when time expires.  The first instance completes in 1 second per operation,
	// the second instance in 2 seconds, the third instance in 4 seconds, and so on.
	flag.Float64Var(&sleepTestOpts.instanceGrowthFactor, "instance-growth-factor", 1, "Instance growth factor.  The delay of instance N will be (InitialDelayMS * (InstanceGrowthFactor ^ InstanceCount)).")

	flag.Float64Var(&sleepTestOpts.iterationGrowthFactor, "iteration-growth-factor", 1, "Iteration growth factor.  The delay of iteration N will be (InitialDelayMS * (IterationGrowthFactor ^ IterationCount)).")
}

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
	sleepInterval time.Duration
}

func (g *globalSleepPerfTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	i := time.Duration(float64(sleepTestOpts.initialDelayMs)*math.Pow(sleepTestOpts.instanceGrowthFactor, float64(g.count))) * time.Millisecond
	s := &sleepPerfTest{sleepInterval: i}
	g.count += 1
	return s, nil
}

func (s *sleepPerfTest) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		time.Sleep(s.sleepInterval)
		s.sleepInterval = time.Duration(float64(s.sleepInterval.Nanoseconds()) * sleepTestOpts.iterationGrowthFactor)
	}
	return nil
}

func (s *sleepPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
