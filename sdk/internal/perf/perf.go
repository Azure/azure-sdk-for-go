// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func init() {
	flag.IntVar(&duration, "d", 10, "Duration of test in seconds.")
	flag.IntVar(&duration, "duration", 10, "Duration of test in seconds.")

	flag.StringVar(&testProxyURLs, "test-proxies", "", "test proxy URLs to target. This can be a semi-colon separated list for multiple proxies.")
	flag.StringVar(&testProxyURLs, "x", "", "test proxy URLs to target. This can be a semi-colon separated list for multiple proxies.")

	flag.IntVar(&warmUpDuration, "warmup", 5, "Duration of warmup in seconds.")
	flag.IntVar(&warmUpDuration, "w", 5, "Duration of warmup in seconds.")

	flag.IntVar(&parallelInstances, "parallel", 1, "Degree of parallelism to run with.")
	flag.IntVar(&parallelInstances, "p", 1, "Degree of parallelism to run with.")

	flag.IntVar(&numProcesses, "maxprocs", runtime.NumCPU(), "Number of CPUs to use.")

	flag.BoolVar(&debug, "debug", false, "Print debugging information")
}

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

// PerfTestOptions are the options for a performance test. Name and Transporter can be
// used by an individual performance test.
type PerfTestOptions struct {
	// Transporter is the azcore.Transporter instance used for sending requests.
	Transporter HTTPClient

	// Name is the name of an individual test
	Name string

	// parallelIndex is the index of the goroutine
	parallelIndex int

	// number of warmup operations completed
	warmupCount   int64
	warmupStart   *time.Time
	warmupElapsed time.Duration

	// number of operations runCount
	runCount   int64
	runStart   *time.Time
	runElapsed time.Duration

	finished bool
}

func newPerfTestOptions(name string) PerfTestOptions {
	return PerfTestOptions{
		Name:        name,
		warmupStart: &time.Time{},
		runStart:    &time.Time{},
	}
}

// increment does an atomic increment of the warmup or non-warmup performance test
func (p *PerfTestOptions) increment(warmup bool) {
	if warmup {
		atomic.AddInt64(&p.warmupCount, 1)
	} else {
		atomic.AddInt64(&p.runCount, 1)
	}
}

// NewPerfTest returns an instance of PerfTest and embeds the given `options` in the struct
type NewPerfTest func(context.Context, PerfTestOptions) (GlobalPerfTest, error)

// PerfMethods is the struct given to the Run method.
type PerfMethods struct {
	Register func()
	New      func(context.Context, PerfTestOptions) (GlobalPerfTest, error)
}

// Run runs an individual test, registers, and parses command line flags
func Run(tests map[string]PerfMethods) {
	if len(os.Args) < 2 {
		// Error out and show available perf tests
		fmt.Println("Available performance tests:")
		for name := range tests {
			fmt.Printf("\t%s\n", name)
		}
		flag.PrintDefaults()
		return
	}

	testNameToRun := os.Args[1]
	var perfTestToRun PerfMethods
	var ok bool
	if perfTestToRun, ok = tests[testNameToRun]; !ok {
		// Error out and show available perf tests
		fmt.Println("Available performance tests:")
		for name := range tests {
			fmt.Printf("\t%s\n", name)
		}
		flag.PrintDefaults()
		return
	}

	if perfTestToRun.Register != nil {
		perfTestToRun.Register()
	}

	// We strip off the first argument because that is used in determining the test
	os.Args = os.Args[1:]

	flag.Parse()

	if numProcesses > 0 {
		val := runtime.GOMAXPROCS(numProcesses)
		if debug {
			fmt.Printf("Changed GOMAXPROCS from %d to %d\n", val, numProcesses)
		}
	}

	fmt.Printf("\tRunning %s\n", testNameToRun)

	runner := newPerfRunner(perfTestToRun, testNameToRun)
	err := runner.Run()
	if err != nil {
		panic(err)
	}
}
