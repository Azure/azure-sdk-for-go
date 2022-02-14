// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/spf13/pflag"
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

// PerfTestOptions are the options for a performance test. Name and Transporter can be
// used by an individual performance test.
type PerfTestOptions struct {
	// Transporter is the azcore.Transporter instance used for sending requests.
	Transporter HTTPClient

	// Name is the name of an individual test
	Name string

	// parallelIndex is the index of the goroutine
	parallelIndex int
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

	testNameToRun := os.Args[1]
	var perfTestToRun PerfMethods
	var ok bool
	if perfTestToRun, ok = tests[testNameToRun]; !ok {
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

	fmt.Printf("\tRunning %s\n", testNameToRun)

	err = runPerfTest(testNameToRun, perfTestToRun.New)
	if err != nil {
		panic(err)
	}
}
