// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/azperf/internal/recording"
)

var (
	Duration       int
	TimeoutSeconds int
	TestProxy      string
	WarmUp         int
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
	fmt.Println("Running `GlobalSetup`")

	fmt.Printf("Deadline of %d\n", TimeoutSeconds)
	_, cancel := getLimitedContext(time.Duration(TimeoutSeconds) * time.Second)
	defer cancel()

	recording.SetRecordingMode("live")
	err := p.GlobalSetup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func runSetup(p PerfTest) error {
	fmt.Println("Running `Setup`")

	_, cancel := getLimitedContext(time.Duration(TimeoutSeconds) * time.Second)
	defer cancel()

	err := p.Setup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// runTest takes care of the semantics of running a single iteration. It returns the number of times the test ran as an int, the exact number
// of seconds the test ran as a float64, and any errors.
func runTest(p PerfTest) (int, float64, error) {
	fmt.Println("Beginning `Run`")
	_, cancel := getLimitedContext(time.Duration(TimeoutSeconds) * time.Second)
	defer cancel()

	// If we are using the test proxy need to set up the in-memory recording.
	usingProxy := TestProxy == "http" || TestProxy == "https"
	if usingProxy {
		// First request goes through in Live mode
		fmt.Println("Running in live")
		recording.SetRecordingMode("live")
		err := p.Run(context.TODO())
		if err != nil {
			return 0, 0.0, err
		}

		// 2nd request goes through in Record mode
		fmt.Println("Running in record")
		recording.SetRecordingMode("record")
		recording.Start(p.GetMetadata(), nil)
		err = p.Run(context.TODO())
		if err != nil {
			return 0, 0.0, err
		}
		recording.Stop(p.GetMetadata(), nil)

		fmt.Println("Running the rest in playback")
		// All ensuing requests go through in Playback mode
		recording.SetRecordingMode("playback")
		recording.Start(p.GetMetadata(), nil)
	}

	start := time.Now()
	count := 0
	for time.Since(start).Seconds() < float64(Duration) {
		err := p.Run(context.Background())
		if err != nil {
			return 0, 0.0, err
		}
		count += 1
	}

	elapsed := time.Since(start).Seconds()
	// Stop the proxy now
	recording.Stop(p.GetMetadata(), nil)
	recording.SetRecordingMode("live")
	return count, elapsed, nil
}

// runTearDown takes care of the semantics for tearing down a single iteration of a performance test.
func runTearDown(p PerfTest) error {
	fmt.Println("Running `Teardown`")

	_, cancel := getLimitedContext(time.Duration(TimeoutSeconds) * time.Second)
	defer cancel()

	err := p.TearDown(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func runGlobalTearDown(p PerfTest) error {
	fmt.Println("Running `GlobalTeardown`")
	_, cancel := getLimitedContext(time.Duration(TimeoutSeconds) * time.Second)
	defer cancel()

	err := p.GlobalTearDown(context.Background())
	if err != nil {
		return err
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

	runSetup(p)

	count, end, err := runTest(p)
	if err != nil {
		return err
	}

	printIteration(end, count)

	totalCount += count
	totalTime += end

	err = runTearDown(p)
	if err != nil {
		return err

	}

	err = runGlobalTearDown(p)
	printFinal(totalCount, totalTime)

	return err
}

func printIteration(t float64, count int) {
	perSecond := float64(count) / t
	fmt.Printf("Completed %d operations in %.4f seconds. Averaged %.4f operations per second.\n", count, t, perSecond)
}

func printFinal(totalCount int, totalElapsed float64) {
	perSecond := float64(totalCount) / totalElapsed
	fmt.Printf("Summary: Completed %d operations in %.4f seconds. Averaged %.4f ops/sec.\n", totalCount, totalElapsed, perSecond)
}
