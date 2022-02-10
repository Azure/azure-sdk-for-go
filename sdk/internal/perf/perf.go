// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/pflag"
)

var (
	debug        bool
	Duration     int
	TestProxy    string
	WarmUp       int
	Parallel     int
	wg           sync.WaitGroup
	numProcesses int
)

type PerfTest interface {
	// GetMetadata returns the metadta for a test.
	GetMetadata() PerfTestOptions

	// GlobalSetup is run one time per performance test, as the first method.
	GlobalSetup(context.Context) error

	// Setup is run once for each parallel instance.
	Setup(context.Context) error

	// Run is the function that is being measured.
	Run(context.Context) error

	// Cleanup is run once for each parallel instance.
	Cleanup(context.Context) error

	// GlobalCleanup is run one time per performance test, as the final method.
	GlobalCleanup(context.Context) error
}

// HTTPClient is the same interface as azcore.Transporter
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// PerfTestOptions are the options for a performance test. Name and ProxyInstance can be
// used by an individual performance test.
type PerfTestOptions struct {
	parallelIndex int
	ProxyInstance HTTPClient
	Name          string
}

// NewPerfTest returns an instance of PerfTest and embeds the given `options` in the struct
type NewPerfTest func(options *PerfTestOptions) PerfTest

func runGlobalSetup(p PerfTest) error {
	if debug {
		log.Println("running GlobalSetup")
	}
	return p.GlobalSetup(context.Background())
}

func runSetup(p PerfTest) error {
	if debug {
		log.Println("running Setup")
	}
	return p.Setup(context.Background())
}

// runTest takes care of the semantics of running a single iteration. It returns the number of times the test ran as an int, the exact number
// of seconds the test ran as a float64, and any errors.
func runTest(p PerfTest, index int, c chan runResult) {
	defer wg.Done()
	if debug {
		log.Printf("number of proxies %d", len(proxyTransportsSuite))
	}

	ID := fmt.Sprintf("%s-%d", p.GetMetadata().Name, p.GetMetadata().parallelIndex)

	// If we are using the test proxy need to set up the in-memory recording.
	if TestProxy != "" {
		// First request goes through in Live mode
		proxyTransportsSuite[ID].SetMode("live")
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}

		// 2nd request goes through in Record mode
		proxyTransportsSuite[ID].SetMode("record")
		err = proxyTransportsSuite[ID].start(p.GetMetadata().Name)
		if err != nil {
			panic(err)
		}
		err = p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}
		err = proxyTransportsSuite[ID].stop(p.GetMetadata().Name)
		if err != nil {
			panic(err)
		}

		// All ensuing requests go through in Playback mode
		proxyTransportsSuite[ID].SetMode("playback")
		err = proxyTransportsSuite[ID].start(p.GetMetadata().Name)
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
				// if we are w/in one second of the end time, we do not
				// want to send any more results, we'll just send a final result
				lastPrint += 10.0
			}
		}
	}

	elapsed := time.Since(timeStart).Seconds()
	lastSecondCount := totalCount - sumInts(perSecondCount)
	c <- runResult{
		completed:     true,
		count:         lastSecondCount,
		parallelIndex: index,
		timeInSeconds: elapsed,
	}

	if TestProxy != "" {
		// Stop the proxy now
		err := proxyTransportsSuite[ID].stop(p.GetMetadata().Name)
		if err != nil {
			c <- runResult{err: err}
		}
		proxyTransportsSuite[ID].SetMode("live")
	}
}

// runCleanup takes care of the semantics for tearing down a single iteration of a performance test.
func runCleanup(p PerfTest) error {
	return p.Cleanup(context.Background())
}

func runGlobalCleanup(p PerfTest) error {
	return p.GlobalCleanup(context.Background())
}

// runResult is the result sent back through the channel for updates and final results
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

// Spins off each Parallel instance as a separate goroutine, reads messages and runs cleanup/setup methods.
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
			ID := fmt.Sprintf("%s-%d", p(nil).GetMetadata().Name, idx)
			transporter, err := NewProxyTransport(&TransportOptions{
				TestName: ID,
				proxyURL: TestProxy,
			})
			if err != nil {
				panic(err)
			}
			options.ProxyInstance = transporter
		} else {
			options.ProxyInstance = defaultHTTPClient
		}
		options.parallelIndex = idx

		perfTest := p(options)
		perfTests = append(perfTests, perfTest)

		// Run the setup for a single instance
		err := runSetup(perfTest)
		if err != nil {
			return fmt.Errorf("there was an error with the Setup method: %v", err.Error())
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
		if debug {
			log.Println("Handling message: ", msg)
		}
		if msg.err != nil {
			panic(err)
		}
		handleMessage(w, msg)
	}

	// Print before running the cleanup in case cleanup takes a while
	printFinalResults(elapsedTimes, perSecondCount, false)

	// Run Cleanup on each parallel instance
	for _, pTest := range perfTests {
		err = runCleanup(pTest)
		if err != nil {
			panic(err)
		}
	}

	err = runGlobalCleanup(p(nil))
	if err != nil {
		return fmt.Errorf("there was an error with the GlobalCleanup method: %v", err.Error())
	}

	return nil
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

// RegisterArguments is used to register local arguments. This is called before `pflag.Parse()`
func RegisterArguments(f func()) {
	f()
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
