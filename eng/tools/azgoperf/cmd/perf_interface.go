// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"context"
	"fmt"
	"time"
)

type PerfTest interface {
	// GetMetadata returns the name of the test
	GetMetadata() string

	GlobalSetup(context.Context) error

	Setup(context.Context) error

	Run(context.Context) error

	TearDown(context.Context) error

	GlobalTearDown(context.Context) error
}

func getLimitedContext(t time.Duration) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, t)
}

func runGlobalSetup(p PerfTest) error {
	fmt.Println("Running Global Setup")

	fmt.Println("Deadline of ", timeoutSeconds)
	ctx, cancel := getLimitedContext(time.Duration(timeoutSeconds) * time.Second)
	defer cancel()

	err := p.GlobalSetup(ctx)
	if err != nil {
		return fmt.Errorf("received an error while running global setup: %s", err.Error())
	}
	return nil
}

func runSetup(p PerfTest, i int) error {
	fmt.Println("\nRunning setup for iteration #", i)

	ctx, cancel := getLimitedContext(time.Duration(timeoutSeconds) * time.Second)
	defer cancel()

	err := p.Setup(ctx)
	if err != nil {
		return fmt.Errorf("received an error while running setup: %s", err.Error())
	}
	return nil
}

// runTest takes care of the semantics of running a single iteration. It returns the exact number
// of seconds the test ran for as a float64.
func runTest(p PerfTest) (int, float64, error) {
	ctx, cancel := getLimitedContext(time.Duration(timeoutSeconds) * time.Second)
	defer cancel()

	start := time.Now()
	count := 0
	for time.Since(start).Seconds() < float64(duration) {
		p.Run(ctx)
		count += 1
	}

	return count, time.Since(start).Seconds(), nil
}

// runTearDown takes care of the semantics for tearing down a single iteration of a performance test.
func runTearDown(p PerfTest) error {
	fmt.Println("\nRunning teardown")

	ctx, cancel := getLimitedContext(time.Duration(timeoutSeconds) * time.Second)
	defer cancel()

	err := p.TearDown(ctx)
	if err != nil {
		return fmt.Errorf("received an error while running teardown: %s", err.Error())
	}
	return nil
}

func runGlobalTearDown(p PerfTest) error {
	fmt.Println("Running global teardown")
	ctx, cancel := getLimitedContext(time.Duration(timeoutSeconds) * time.Second)
	defer cancel()

	err := p.GlobalTearDown(ctx)
	if err != nil {
		return fmt.Errorf("received an error while running global teardown: %s", err.Error())
	}
	return nil
}

func RunPerfTest(p PerfTest) error {
	err := runGlobalSetup(p)
	if err != nil {
		return err
	}

	totalCount := 0
	totalTime := 0.0
	for i := 0; i < iterations; i++ {
		runSetup(p, i)

		count, end, err := runTest(p)
		if err != nil {
			return err
		}

		printIteration(end, count, i)

		totalCount += count
		totalTime += end

		err = runTearDown(p)
		if err != nil {
			return err
		}
	}

	runGlobalTearDown(p)
	printFinal(totalCount, totalTime, iterations)

	return runGlobalSetup(p)
}

func printIteration(t float64, count int, itNum int) {
	perSecond := float64(count) / t
	fmt.Printf("Iteration #%d: Completed %d operations in %.4f seconds. Averaged %.4f operations per second.\n", itNum, count, t, perSecond)
}

func printFinal(totalCount int, totalElapsed float64, i int) {
	perSecond := float64(totalCount) / totalElapsed
	fmt.Printf("\nSummary: Completed %d operations in %.4f seconds. Averaged %.4f ops/sec.\n", totalCount, totalElapsed, perSecond)
}
