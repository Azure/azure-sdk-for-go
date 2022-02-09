// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
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

var numProcesses int

type PerfTest interface {
	// GetMetadata returns the name of the test
	GetMetadata() PerfTestOptions

	GlobalSetup(context.Context) error

	Setup(context.Context) error

	Run(context.Context) error

	Cleanup(context.Context) error

	GlobalCleanup(context.Context) error
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type PerfTestOptions struct {
	ParallelIndex int
	ProxyInstance HTTPClient
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
func runTest(p PerfTest, index int, c chan runResult) {
	defer wg.Done()

	// If we are using the test proxy need to set up the in-memory recording.
	if TestProxy != "" {
		// First request goes through in Live mode
		setRecordingMode("live")
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}

		// 2nd request goes through in Record mode
		setRecordingMode("record")
		err = start(p.GetMetadata().Name, nil)
		if err != nil {
			panic(err)
		}
		err = p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}
		err = stop(p.GetMetadata().Name, nil)
		if err != nil {
			panic(err)
		}

		// All ensuing requests go through in Playback mode
		setRecordingMode("playback")
		err = start(p.GetMetadata().Name, nil)
		if err != nil {
			panic(err)
		}
	}

	if WarmUp > 0 {
		warmUpStart := time.Now()
		warmUpLastPrint := 1.0
		warmUpPerSecondCount := make([]int, 0)
		warmupCount := 0
		for time.Since(warmUpStart).Seconds() < float64(WarmUp) {
			err := p.Run(context.Background())
			if err != nil {
				c <- runResult{err: err}
			}
			warmupCount += 1

			if time.Since(warmUpStart).Seconds() > warmUpLastPrint {
				thisSecondCount := warmupCount - sumInts(warmUpPerSecondCount)
				c <- runResult{
					warmup:        true,
					count:         thisSecondCount,
					parallelIndex: index,
					timeInSeconds: time.Since(warmUpStart).Seconds(),
				}

				warmUpPerSecondCount = append(warmUpPerSecondCount, thisSecondCount)
				warmUpLastPrint += 1.0
				if int(warmUpLastPrint) == WarmUp {
					// We can have odd scenarios where we send this
					// message, and the final message below right after
					warmUpLastPrint += 1.0
				}
			}
		}

		thisSecondCount := warmupCount - sumInts(warmUpPerSecondCount)
		c <- runResult{
			warmup:        true,
			completed:     true,
			count:         thisSecondCount,
			parallelIndex: index,
			timeInSeconds: time.Since(warmUpStart).Seconds(),
		}
	}

	timeStart := time.Now()
	totalCount := 0
	lastPrint := 1.0
	perSecondCount := make([]int, 0)
	for time.Since(timeStart).Seconds() < float64(Duration) {
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}
		totalCount += 1

		// Every second (roughly) we send an update through the channel
		if time.Since(timeStart).Seconds() > lastPrint {
			thisCount := totalCount - sumInts(perSecondCount)
			c <- runResult{
				count:         thisCount,
				parallelIndex: index,
				completed:     false,
				timeInSeconds: time.Since(timeStart).Seconds(),
			}
			lastPrint += 1.0
			perSecondCount = append(perSecondCount, thisCount)
			if int(lastPrint) == Duration {
				lastPrint += 1.0
			}
		}
	}

	elapsed := time.Since(timeStart).Seconds()

	if TestProxy != "" {
		// Stop the proxy now
		err := stop(p.GetMetadata().Name, nil)
		if err != nil {
			c <- runResult{err: err}
		}
		setRecordingMode("live")
	}
	c <- runResult{count: totalCount, timeInSeconds: elapsed, completed: true, parallelIndex: index}
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
	// number of iterations completed since the previous message
	count int

	// The time the update comes from
	timeInSeconds float64

	// if there is an error it will be here
	err error

	// true when this is the last result from a go routine
	completed bool

	// Index of the goroutine
	parallelIndex int

	// indicates if result is from a warmup start
	warmup bool
}

func (r runResult) String() string {
	return fmt.Sprintf("{count: %d, timeInSeconds: %.2f, err: %v, complete: %v, parallelIndex: %d, warmup: %v}", r.count, r.timeInSeconds, r.err, r.completed, r.parallelIndex, r.warmup)
}

func runPerfTest(p NewPerfTest) error {
	err := runGlobalSetup(p(nil))
	if err != nil {
		panic(err)
	}

	var perfTests []PerfTest

	w := tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight)
	if err != nil {
		panic(err)
	}

	messages := make(chan runResult)
	for idx := 0; idx < Parallel; idx++ {
		options := &PerfTestOptions{}
		if TestProxy != "" {
			transporter, err := NewProxyTransport(&TransportOptions{TestName: p(nil).GetMetadata().Name})
			if err != nil {
				panic(err)
			}
			options.ProxyInstance = transporter
		} else {
			options.ProxyInstance = defaultHTTPClient
		}
		options.ParallelIndex = idx

		perfTest := p(options)
		perfTests = append(perfTests, perfTest)

		// Run the setup for a single instance
		err := runSetup(perfTest)
		if err != nil {
			return err
		}
	}

	for idx, perfTest := range perfTests {
		// Create a thread for running a single test
		wg.Add(1)
		go runTest(perfTest, idx, messages)
	}

	// Add another goroutine to close the channel after completion
	go func() {
		wg.Wait()
		close(messages)
	}()

	// Read incoming messages and handle status updates
	for msg := range messages {
		if msg.err != nil {
			panic(err)
		}
		handleMessage(w, msg)
	}

	// Print before running the cleanup in case cleanup takes a while
	printFinalResults()

	// Run Cleanup on each parallel instance
	for _, pTest := range perfTests {
		err = runCleanup(pTest)
		if err != nil {
			panic(err)
		}
	}

	err = runGlobalCleanup(p(nil))

	return err
}

// testsToRun trims the slice of PerfTest to only those that are flagged as running.
func testsToRun(registered []NewPerfTest) NewPerfTest {
	args := os.Args[1:]
	for _, r := range registered {
		p := r(nil)
		for _, arg := range args {
			if p.GetMetadata().Name == arg {
				return r
			}
		}
	}

	return nil
}

var registerCalled bool

// RegisterArguments is used to
func RegisterArguments(f func()) {
	f()
	registerCalled = true
}

// Run runs all individual tests
func Run(perfTests []NewPerfTest) {
	// Start with adding all of our arguments
	pflag.IntVarP(&Duration, "duration", "d", 10, "Duration of the test in seconds. Default is 10.")
	pflag.StringVarP(&TestProxy, "test-proxies", "x", "", "whether to target http or https proxy (default is neither)")
	pflag.IntVarP(&WarmUp, "warmup", "w", 5, "Duration of warmup in seconds. Default is 5.")
	pflag.IntVarP(&Parallel, "parallel", "p", 1, "Degree of parallelism to run with. Default is 1.")
	pflag.IntVar(&numProcesses, "maxprocs", runtime.NumCPU(), "Number of CPUs to use.")

	pflag.BoolVarP(&debug, "debug", "g", false, "Print debugging information")
	err := pflag.CommandLine.MarkHidden("debug")
	if err != nil {
		panic(err)
	}

	pflag.Parse()

	if !registerCalled && debug {
		fmt.Println("There were no local flags added.")
	}

	if numProcesses > 0 {
		val := runtime.GOMAXPROCS(numProcesses)
		if debug {
			fmt.Printf("Changed GOMAXPROCS from %d to %d\n", val, numProcesses)
		}
	}

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

	err = runPerfTest(perfTestToRun)
	if err != nil {
		panic(err)
	}
}
