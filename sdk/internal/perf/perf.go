// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
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
	GetMetadata() PerfTestOptions

	GlobalSetup(context.Context) error

	Setup(context.Context) error

	Run(context.Context) error

	Cleanup(context.Context) error

	GlobalCleanup(context.Context) error
}

type PerfTestOptions struct {
	ParallelIndex int
	ProxyInstance *RecordingHTTPClient
	Name          string
}

type NewPerfTest func(options *PerfTestOptions) PerfTest

func runGlobalSetup(p PerfTest) error {
	setRecordingMode("live")
	err := p.GlobalSetup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func runSetup(p PerfTest) error {
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
		start(p.GetMetadata().Name, nil)
		err = p.Run(context.Background())
		if err != nil {
			c <- runResult{count: 0, timeInSeconds: 0.0, err: err}
		}
		stop(p.GetMetadata().Name, nil)

		// All ensuing requests go through in Playback mode
		setRecordingMode("playback")
		start(p.GetMetadata().Name, nil)
	}

	if WarmUp > 0 {
		warmUpStart := time.Now()
		for time.Since(warmUpStart).Seconds() < float64(WarmUp) {
			err := p.Run(context.Background())
			if err != nil {
				c <- runResult{count: 0, timeInSeconds: 0.0, err: err}
			}
		}
	}

	start := time.Now()
	totalCount := 0
	lastPrint := 1.0
	perSecondCount := make([]int, 0)
	w := tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for time.Since(start).Seconds() < float64(Duration) {
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{count: 0, timeInSeconds: 0.0, err: err}
		}
		totalCount += 1

		// Every second (roughly) we print out an update
		if time.Since(start).Seconds() > float64(lastPrint) {
			thisCount := totalCount - sumInts(perSecondCount)
			perSecondCount = append(perSecondCount, thisCount)
			_, err = fmt.Fprintf(
				w,
				"%s\t%s\t%.2f\t\n",
				commaIze(thisCount),
				commaIze(totalCount),
				float64(sumInts(perSecondCount))/time.Since(start).Seconds(),
			)
			if err != nil {
				c <- runResult{count: 0, timeInSeconds: 0.0, err: err}
			}
			lastPrint = time.Since(start).Seconds() + 1.0
			w.Flush()
		}
	}

	elapsed := time.Since(start).Seconds()

	// Stop the proxy now
	stop(p.GetMetadata().Name, nil)
	setRecordingMode("live")
	c <- runResult{count: totalCount, timeInSeconds: elapsed, err: nil}
}

// runCleanup takes care of the semantics for tearing down a single iteration of a performance test.
func runCleanup(p PerfTest) error {
	err := p.Cleanup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func runGlobalCleanup(p PerfTest) error {
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

func runPerfTest(p NewPerfTest) error {
	err := runGlobalSetup(p(nil))
	if err != nil {
		return err
	}

	var channels []chan runResult
	var perfTests []PerfTest

	fmt.Println("=== Test ===")

	w := tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Current\tTotal\tAverage\t")
	err = w.Flush()
	if err != nil {
		return err
	}

	for idx := 0; idx < Parallel; idx++ {
		var transporter *RecordingHTTPClient
		var err error
		if TestProxy != "" {
			transporter, err = NewProxyTransport(nil)
			if err != nil {
				return err
			}
		}

		perfTest := p(&PerfTestOptions{
			ParallelIndex: idx,
			ProxyInstance: transporter,
		})
		perfTests = append(perfTests, perfTest)

		// Run the setup for a single instance
		runSetup(perfTest)

		// Create a thread for running a single test
		wg.Add(1)
		c := make(chan runResult, 1) // Create a buffered channel
		channels = append(channels, c)
		go runTest(perfTest, c)
	}

	wg.Wait()

	opsPerSecond := 0.0
	totalOperations := 0
	// Get the results from the channels
	for _, channel := range channels {
		result := <-channel
		if err != nil {
			return err
		}

		opsPerSecond += float64(result.count) / result.timeInSeconds
		totalOperations += result.count
	}

	// Run Cleanup on each parallel instance
	for _, pTest := range perfTests {
		err = runCleanup(pTest)
		if err != nil {
			return err
		}
	}

	err = runGlobalCleanup(p(nil))

	fmt.Println("\n=== Results ===")
	secondsPerOp := 1.0 / opsPerSecond
	weightedAvgSec := float64(totalOperations) / opsPerSecond
	fmt.Printf(
		"Completed %s operations in a weighted-average of %.2fs (%s ops/s, %.3f s/op)",
		commaIze(totalOperations),
		weightedAvgSec,
		commaIze(int(opsPerSecond)),
		secondsPerOp,
	)

	return err
}

// testsToRun trims the slice of PerfTest to only those that are flagged as running.
func testsToRun(registered []NewPerfTest) NewPerfTest {
	var ret []NewPerfTest

	args := os.Args[1:]
	for _, r := range registered {
		p := r(nil)
		for _, arg := range args {
			if p.GetMetadata().Name == arg {
				ret = append(ret, r)
			}
		}
	}

	if len(ret) > 1 {
		fmt.Println("Performance only supports running one test per process. Run the performance multiple times per performance for each test you want to run.")
		os.Exit(1)
	} else if len(ret) == 0 {
		return nil
	}

	return ret[0]
}

// Run runs all individual tests
func Run(perfTests []NewPerfTest) {
	// Start with adding all of our arguments
	pflag.IntVarP(&Duration, "duration", "d", 10, "Duration of the test in seconds. Default is 10.")
	pflag.StringVarP(&TestProxy, "test-proxies", "x", "", "whether to target http or https proxy (default is neither)")
	pflag.IntVarP(&WarmUp, "warmup", "w", 5, "Duration of warmup in seconds. Default is 5.")
	pflag.IntVarP(&Parallel, "parallel", "p", 1, "Degree of parallelism to run with. Default is 1.")

	pflag.BoolVarP(&debug, "debug", "g", false, "Print debugging information")
	pflag.CommandLine.MarkHidden("debug")

	// TODO: add individual performance tests local flags

	pflag.Parse()

	perfTestToRun := testsToRun(perfTests)

	if perfTestToRun == nil {
		fmt.Println("Available performance tests:")
		for _, p := range perfTests {
			c := p(nil)
			fmt.Printf("\t%s\n", c.GetMetadata().Name)
		}
		return
	}

	fmt.Printf("\tRunning %s\n", perfTestToRun(nil).GetMetadata().Name)

	err := runPerfTest(perfTestToRun)
	if err != nil {
		panic(err)
	}
}
