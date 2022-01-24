// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

var (
	debug     bool
	Duration  int
	TestProxy string
	WarmUp    int
	Parallel  int
	wg        sync.WaitGroup
)

type PerfTest interface {
	// GetMetadata returns the name of the test
	// TODO: Add local flags to the GetMetadata
	GetMetadata() string

	GlobalSetup(context.Context) error

	Setup(context.Context) error

	Run(context.Context) error

	Cleanup(context.Context) error

	GlobalCleanup(context.Context) error
}

func getLimitedContext(t time.Duration) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, t)
}

func runGlobalSetup(p PerfTest) error {
	fmt.Println("Running `GlobalSetup`")

	setRecordingMode("live")
	err := p.GlobalSetup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func runSetup(p PerfTest) error {
	fmt.Println("Running `Setup`")

	err := p.Setup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// runTest takes care of the semantics of running a single iteration. It returns the number of times the test ran as an int, the exact number
// of seconds the test ran as a float64, and any errors.
func runTest(p PerfTest, c chan runResult) {
	defer wg.Done()
	fmt.Println("Beginning `Run`")

	// If we are using the test proxy need to set up the in-memory recording.
	if TestProxy != "" {
		// First request goes through in Live mode
		setRecordingMode("live")
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{count: 0, timeInSeconds: 0.0, err: err}
		}

		// 2nd request goes through in Record mode
		setRecordingMode("record")
		start(p.GetMetadata(), nil)
		err = p.Run(context.Background())
		if err != nil {
			c <- runResult{count: 0, timeInSeconds: 0.0, err: err}
		}
		stop(p.GetMetadata(), nil)

		// All ensuing requests go through in Playback mode
		setRecordingMode("playback")
		start(p.GetMetadata(), nil)
	}

	start := time.Now()
	count := 0
	for time.Since(start).Seconds() < float64(Duration) {
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{count: 0, timeInSeconds: 0.0, err: err}
		}
		count += 1
	}

	elapsed := time.Since(start).Seconds()

	// Stop the proxy now
	stop(p.GetMetadata(), nil)
	setRecordingMode("live")
	c <- runResult{count: count, timeInSeconds: elapsed, err: nil}
}

// runCleanup takes care of the semantics for tearing down a single iteration of a performance test.
func runCleanup(p PerfTest) error {
	fmt.Println("Running `Cleanup`")

	err := p.Cleanup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func runGlobalCleanup(p PerfTest) error {
	fmt.Println("Running `GlobalCleanup`")

	err := p.GlobalCleanup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

type runResult struct {
	count         int
	timeInSeconds float64
	err           error
}

func leftPad(i int) string {
	if i >= 100 {
		return fmt.Sprintf("%d", i)
	} else if i >= 10 {
		return fmt.Sprintf("0%d", i)
	} else if i > 0 {
		return fmt.Sprintf("00%d", i)
	}
	return "000"
}

func commaIze(i int) string {
	if i < 1000 {
		return fmt.Sprintf("%d", i)
	}

	copy := i
	ret := ""
	for copy >= 1000 {
		temp := copy % 1000
		tempS := leftPad(temp)
		ret = fmt.Sprintf("%s,%s", ret, tempS)

		copy /= 1000
	}

	if copy == 0 {
		return ret
	}

	return fmt.Sprintf("%d%s", copy, ret)
}

func runPerfTest(p PerfTest) error {
	err := runGlobalSetup(p)
	if err != nil {
		return err
	}

	runSetup(p)

	var channels []chan runResult

	for idx := 0; idx < Parallel; idx++ {
		wg.Add(1)
		c := make(chan runResult, 1) // Create a buffered channel
		channels = append(channels, c)
		go runTest(p, c)
	}

	wg.Wait()

	opsPerSecond := 0.0
	totalOperations := 0
	// Get the results from the channels
	for idx, channel := range channels {
		result := <-channel
		if err != nil {
			return err
		}
		printIteration(idx, result.timeInSeconds, result.count)

		opsPerSecond += float64(result.count) / result.timeInSeconds
		totalOperations += result.count
	}

	err = runCleanup(p)
	if err != nil {
		return err
	}

	err = runGlobalCleanup(p)

	secondsPerOp := 1.0 / opsPerSecond
	weightedAvgSec := float64(totalOperations) / opsPerSecond
	fmt.Printf("Completed %d operations in a weighted-average of %.2fs (%d ops/s, %.3f s/op)", totalOperations, weightedAvgSec, int(opsPerSecond), secondsPerOp)

	return err
}

func printIteration(idx int, t float64, count int) {
	perSecond := float64(count) / t
	fmt.Printf("Thread #%d: Completed %d operations in %.4f seconds. Averaged %.4f operations per second.\n", idx+1, count, t, perSecond)
}

// testsToRun trims the slice of PerfTest to only those that are flagged as running.
func testsToRun(registered []PerfTest) []PerfTest {
	var ret []PerfTest

	args := os.Args[1:]
	for _, r := range registered {
		for _, arg := range args {
			if r.GetMetadata() == arg {
				ret = append(ret, r)
			}
		}
	}

	if len(ret) > 1 {
		fmt.Println("Performance only supports running one test per process. Run the performance multiple times per performance for each test you want to run.")
		os.Exit(1)
	}

	return ret
}

// Run runs all individual tests
func Run(perfTests []PerfTest) {
	// Start with adding all of our arguments
	pflag.IntVarP(&Duration, "duration", "d", 10, "Duration of the test in seconds. Default is 10.")
	pflag.StringVarP(&TestProxy, "test-proxies", "x", "", "whether to target http or https proxy (default is neither)")
	pflag.IntVarP(&WarmUp, "warmup", "w", 5, "Duration of warmup in seconds. Default is 5.")
	pflag.IntVarP(&Parallel, "parallel", "p", 1, "Degree of parallelism to run with. Default is 1.")

	pflag.BoolVarP(&debug, "debug", "g", false, "Print debugging information")
	pflag.CommandLine.MarkHidden("debug")

	// Need to add individual performance tests local flags

	pflag.Parse()

	perfTestsToRun := testsToRun(perfTests)

	if len(perfTestsToRun) == 0 {
		fmt.Println("Available performance tests:")
		for _, p := range perfTests {
			fmt.Printf("\t%s\n", p.GetMetadata())
		}
		return
	}

	fmt.Println("======= Pre Run Summary =======")
	for _, p := range perfTestsToRun {
		fmt.Printf("\tRunning %s\n", p.GetMetadata())
	}
	fmt.Printf("===============================\n\n")

	for _, p := range perfTestsToRun {
		err := runPerfTest(p)
		if err != nil {
			panic(err)
		}
	}
}
