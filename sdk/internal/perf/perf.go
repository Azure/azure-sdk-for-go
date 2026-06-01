// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func init() {
	registerFlags(flag.CommandLine)
}

// registerFlags binds every perf CLI flag onto the supplied FlagSet. It is
// extracted from init() so tests can construct their own FlagSet, parse a
// slice of args, and assert on the resulting global state without disturbing
// the process-wide flag.CommandLine.
func registerFlags(fs *flag.FlagSet) {
	fs.IntVar(&duration, "d", 10, "Duration of test in seconds")
	fs.IntVar(&duration, "duration", 10, "Duration of test in seconds")

	fs.StringVar(&testProxyURLs, "test-proxies", "", "URIs of TestProxy Servers (separated by ';')")
	fs.StringVar(&testProxyURLs, "x", "", "URIs of TestProxy Servers (separated by ';')")

	fs.IntVar(&warmUpDuration, "warmup", 5, "Duration of warmup in seconds")
	fs.IntVar(&warmUpDuration, "w", 5, "Duration of warmup in seconds")

	fs.IntVar(&parallelInstances, "parallel", 1, "Number of operations to execute in parallel")
	fs.IntVar(&parallelInstances, "p", 1, "Number of operations to execute in parallel")

	fs.IntVar(&iterations, "iterations", 1, "Number of iterations of main test loop")
	fs.IntVar(&iterations, "i", 1, "Number of iterations of main test loop")

	fs.IntVar(&targetRate, "rate", 0, "Target throughput (ops/sec). 0 means unlimited.")
	fs.IntVar(&targetRate, "r", 0, "Target throughput (ops/sec). 0 means unlimited.")

	fs.IntVar(&statusInterval, "status-interval", 1, "Interval to write status to console in seconds")

	fs.BoolVar(&jobStatistics, "job-statistics", false, "Print job statistics (used by automation)")

	fs.IntVar(&numProcesses, "maxprocs", runtime.NumCPU(), "Number of CPUs to use.")

	fs.BoolVar(&debug, "debug", false, "Print debugging information")
	fs.BoolVar(&syncMode, "sync", false, "Runs sync version of test. Accepted for CLI compatibility; no-op for Go perf tests.")
	fs.BoolVar(&resourceTelemetry, "resource-telemetry", false, "Print process resource telemetry summary (memory, GC, and goroutines).")
	fs.BoolVar(&enableOperationLatency, "latency", false, "Track and print per-operation latency statistics")
	fs.BoolVar(&enableOperationLatency, "l", false, "Track and print per-operation latency statistics")
	fs.BoolVar(&noCleanup, "no-cleanup", false, "Disables test cleanup")
	fs.BoolVar(&insecureSkipVerify, "insecure", false, "Allow untrusted SSL certs")
	fs.StringVar(&resultsFilePath, "results-file", "", "File path location to store the results for the test run.")
	fs.StringVar(&workloadConfigPath, "config", "", "Path to workload config JSON file.")
	fs.StringVar(&workloadName, "workload", "", "Workload name from config JSON.")
	fs.StringVar(&outputFilePrefix, "output-file-prefix", "", "Write run artifacts to <prefix>.json/.csv/.txt/.md.")

	// .NET ThreadPool tuning flags — accepted for CLI compatibility, ignored at runtime.
	fs.IntVar(&maxIOCompletionThreads, "max-io-completion-threads", 0, "Compatibility flag; no-op in Go.")
	fs.IntVar(&maxWorkerThreads, "max-worker-threads", 0, "Compatibility flag; no-op in Go.")
	fs.IntVar(&minIOCompletionThreads, "min-io-completion-threads", 0, "Compatibility flag; no-op in Go.")
	fs.IntVar(&minWorkerThreads, "min-worker-threads", 0, "Compatibility flag; no-op in Go.")
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

// Parallel returns the value of the --parallel/-p flag (the number of
// goroutines running PerfTest instances concurrently). It is safe to call
// only after flag parsing has occurred (i.e. from inside a registered test's
// New/Register/Run callbacks). Before flag parsing it returns the default
// value of 1.
func Parallel() int {
	if parallelInstances <= 0 {
		return 1
	}
	return parallelInstances
}

// Run runs an individual test, registers, and parses command line flags
func Run(tests map[string]PerfMethods) {
	invocation, err := resolveRunInvocation(os.Args[1:])
	if err != nil {
		panic(err)
	}

	if invocation.TestName == "" {
		// Error out and show available perf tests
		fmt.Println("Available performance tests:")
		for name := range tests {
			fmt.Printf("\t%s\n", name)
		}
		flag.PrintDefaults()
		return
	}

	testNameToRun := invocation.TestName
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

	// Parse only config-provided and CLI-provided flags.
	// CLI arguments are appended last to allow explicit command-line overrides.
	os.Args = append([]string{os.Args[0]}, invocation.ConfigArgs...)
	os.Args = append(os.Args, invocation.CLIArgs...)

	flag.Parse()

	// --config / --workload are consumed by resolveRunInvocation before
	// flag.Parse runs, so the corresponding globals are not populated by
	// the standard flag machinery. Mirror the resolved invocation onto
	// the globals so downstream artifacts (run summary, =Options= dump)
	// reflect the workload that actually drove the run.
	if invocation.UsesConfig {
		workloadConfigPath = invocation.ConfigPath
		workloadName = invocation.Workload
	}

	// Normalize sentinel defaults for the .NET-compatible Rate/Iterations/
	// StatusInterval flags.
	if iterations < 1 {
		iterations = 1
	}
	if statusInterval < 1 {
		statusInterval = 1
	}
	if targetRate < 0 {
		targetRate = 0
	}

	if numProcesses > 0 {
		val := runtime.GOMAXPROCS(numProcesses)
		if debug {
			fmt.Printf("Changed GOMAXPROCS from %d to %d\n", val, numProcesses)
		}
	}

	if jobStatistics {
		// Matches the .NET runner's first-line marker used by perf-automation.
		fmt.Println("Application started.")
	}

	// Mirror the .NET runner's "=== Options ===" block so perf-automation
	// sees the same machine-parseable configuration dump regardless of language.
	printOptions(testNameToRun)

	if invocation.UsesConfig {
		fmt.Printf("\tRunning %s (workload: %s, config: %s)\n", testNameToRun, invocation.Workload, invocation.ConfigPath)
	} else {
		fmt.Printf("\tRunning %s\n", testNameToRun)
	}

	runner := newPerfRunner(perfTestToRun, testNameToRun)
	err = runner.Run()
	if err != nil {
		panic(err)
	}
}

// printOptions writes a "=== Options ===" header followed by a JSON object
// describing the parsed CLI options. The format mirrors the .NET runner's
// `JsonSerializer.Serialize(options, ...)` output so a shared parser can
// extract the run configuration from either language's stdout.
func printOptions(testName string) {
	type optionsDump struct {
		TestName               string `json:"testName"`
		Duration               int    `json:"duration"`
		Warmup                 int    `json:"warmup"`
		Parallel               int    `json:"parallel"`
		Iterations             int    `json:"iterations"`
		Rate                   int    `json:"rate"`
		StatusInterval         int    `json:"statusInterval"`
		Latency                bool   `json:"latency"`
		NoCleanup              bool   `json:"noCleanup"`
		JobStatistics          bool   `json:"jobStatistics"`
		Sync                   bool   `json:"sync"`
		Insecure               bool   `json:"insecure"`
		TestProxies            string `json:"testProxies,omitempty"`
		ResultsFile            string `json:"resultsFile,omitempty"`
		MaxIOCompletionThreads int    `json:"maxIOCompletionThreads"`
		MaxWorkerThreads       int    `json:"maxWorkerThreads"`
		MinIOCompletionThreads int    `json:"minIOCompletionThreads"`
		MinWorkerThreads       int    `json:"minWorkerThreads"`
	}
	opts := optionsDump{
		TestName:               testName,
		Duration:               duration,
		Warmup:                 warmUpDuration,
		Parallel:               parallelInstances,
		Iterations:             iterations,
		Rate:                   targetRate,
		StatusInterval:         statusInterval,
		Latency:                enableOperationLatency,
		NoCleanup:              noCleanup,
		JobStatistics:          jobStatistics,
		Sync:                   syncMode,
		Insecure:               insecureSkipVerify,
		TestProxies:            testProxyURLs,
		ResultsFile:            resultsFilePath,
		MaxIOCompletionThreads: maxIOCompletionThreads,
		MaxWorkerThreads:       maxWorkerThreads,
		MinIOCompletionThreads: minIOCompletionThreads,
		MinWorkerThreads:       minWorkerThreads,
	}
	b, err := json.MarshalIndent(opts, "", "  ")
	if err != nil {
		return
	}
	fmt.Println("=== Options ===")
	fmt.Println(string(b))
	fmt.Println()
}
