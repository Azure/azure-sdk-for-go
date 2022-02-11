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
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/pflag"
)

var (
	debug             bool
	duration          int
	testProxyURLs     string
	warmUpDuration    int
	parallelInstances int
	wg                sync.WaitGroup
	numProcesses      int
)

// GlobalPerfTest methods execute once per process
type GlobalPerfTest interface {
	// NewPerfTest creates an instance of a PerfTest for each goroutine.
	NewPerfTest(context.Context, *PerfTestOptions) (PerfTest, error)

	// GlobalCleanup is run one time per performance test, as the final method.
	GlobalCleanup(context.Context) error
}

// PerfTest methods once per goroutine
type PerfTest interface {
	// Run is the function that is being measured.
	Run(context.Context) error

	// Cleanup is run once for each parallel instance.
	Cleanup(context.Context) error
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
type NewPerfTest func(context.Context, PerfTestOptions) (GlobalPerfTest, error)

// runTest takes care of the semantics of running a single iteration. It returns the number of times the test ran as an int, the exact number
// of seconds the test ran as a float64, and any errors.
func runTest(p PerfTest, index int, c chan runResult, ID string) {
	defer wg.Done()
	if debug {
		log.Printf("number of proxies %d", len(proxyTransportsSuite))
	}

	// If we are using the test proxy need to set up the in-memory recording.
	if testProxyURLs != "" {
		// First request goes through in Live mode
		proxyTransportsSuite[ID].SetMode("live")
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}

		// 2nd request goes through in Record mode
		proxyTransportsSuite[ID].SetMode("record")
		err = proxyTransportsSuite[ID].start()
		if err != nil {
			panic(err)
		}
		err = p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}
		err = proxyTransportsSuite[ID].stop()
		if err != nil {
			panic(err)
		}

		// All ensuing requests go through in Playback mode
		proxyTransportsSuite[ID].SetMode("playback")
		err = proxyTransportsSuite[ID].start()
		if err != nil {
			panic(err)
		}
	}

	if warmUpDuration > 0 {
		warmUpStart := time.Now()
		warmUpLastPrint := 1.0
		warmUpPerSecondCount := make([]int, 0)
		warmupCount := 0
		for time.Since(warmUpStart).Seconds() < float64(warmUpDuration) {
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
				if int(warmUpLastPrint) == warmUpDuration {
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
	for time.Since(timeStart).Seconds() < float64(duration) {
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

			if int(lastPrint) == duration {
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

	if testProxyURLs != "" {
		// Stop the proxy now
		err := proxyTransportsSuite[ID].stop()
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

// parse the TestProxy input into a slice of strings
func parseProxyURLS() []string {
	var ret []string
	if testProxyURLs == "" {
		return ret
	}

	testProxyURLs = strings.TrimSuffix(testProxyURLs, ";")

	ret = strings.Split(testProxyURLs, ";")

	return ret
}

// Spins off each Parallel instance as a separate goroutine, reads messages and runs cleanup/setup methods.
func runPerfTest(name string, p NewPerfTest) error {
	globalInstance, err := p(context.TODO(), PerfTestOptions{Name: name})
	if err != nil {
		return err
	}

	var perfTests []PerfTest
	var IDs []string
	proxyURLS := parseProxyURLS()

	w := tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight)

	messages := make(chan runResult)
	for idx := 0; idx < parallelInstances; idx++ {
		options := &PerfTestOptions{}

		ID := fmt.Sprintf("%s-%d", name, idx)
		IDs = append(IDs, ID)

		if testProxyURLs != "" {
			proxyURL := proxyURLS[idx%len(proxyURLS)]
			transporter := NewProxyTransport(&TransportOptions{
				TestName: ID,
				proxyURL: proxyURL,
			})
			options.ProxyInstance = transporter
		} else {
			options.ProxyInstance = defaultHTTPClient
		}
		options.parallelIndex = idx

		perfTest, err := globalInstance.NewPerfTest(context.TODO(), options)
		if err != nil {
			panic(err)
		}
		perfTests = append(perfTests, perfTest)
	}

	for idx, perfTest := range perfTests {
		// Create a thread for running a single test
		wg.Add(1)
		go runTest(perfTest, idx, messages, IDs[idx])
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
			panic(msg.err)
		}
		handleMessage(w, msg)
	}

	// Print before running the cleanup in case cleanup takes a while
	printFinalResults(elapsedTimes, perSecondCount, false)

	// Run Cleanup on each parallel instance
	for _, pTest := range perfTests {
		err := runCleanup(pTest)
		if err != nil {
			panic(err)
		}
	}

	err = globalInstance.GlobalCleanup(context.TODO())
	if err != nil {
		return fmt.Errorf("there was an error with the GlobalCleanup method: %v", err.Error())
	}

	return nil
}

type MapInterface struct {
	Register func()
	New      func(context.Context, PerfTestOptions) (GlobalPerfTest, error)
}

// Run runs all individual tests
func Run(tests map[string]MapInterface) {
	// Start with adding all of our arguments
	pflag.IntVarP(&duration, "duration", "d", 10, "Duration of the test in seconds. Default is 10.")
	pflag.StringVarP(&testProxyURLs, "test-proxies", "x", "", "whether to target http or https proxy (default is neither)")
	pflag.IntVarP(&warmUpDuration, "warmup", "w", 5, "Duration of warmup in seconds. Default is 5.")
	pflag.IntVarP(&parallelInstances, "parallel", "p", 1, "Degree of parallelism to run with. Default is 1.")
	pflag.IntVar(&numProcesses, "maxprocs", runtime.NumCPU(), "Number of CPUs to use.")

	pflag.BoolVar(&debug, "debug", false, "Print debugging information")
	err := pflag.CommandLine.MarkHidden("debug")
	if err != nil {
		panic(err)
	}

	if numProcesses > 0 {
		val := runtime.GOMAXPROCS(numProcesses)
		if debug {
			fmt.Printf("Changed GOMAXPROCS from %d to %d\n", val, numProcesses)
		}
	}

	var perfTestToRun MapInterface
	var ok bool
	if perfTestToRun, ok = tests[os.Args[1]]; !ok {
		// Error out and show available perf tests
		fmt.Println("Available performance tests:")
		for name := range tests {
			fmt.Printf("\t%s\n", name)
		}
		return
	}

	if perfTestToRun.Register != nil {
		perfTestToRun.Register()
	}
	pflag.Parse()

	fmt.Printf("\tRunning %s\n", os.Args[1])

	err = runPerfTest(os.Args[1], perfTestToRun.New)
	if err != nil {
		panic(err)
	}
}
