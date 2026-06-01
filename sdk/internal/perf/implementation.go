// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"strings"
	"sync"
)

var (
	// debug is true if --debug is specified
	debug bool
	// syncMode is accepted for compatibility with other language perf runners.
	syncMode bool
	// duration is the -d/--duration flag
	duration int
	// testProxyURLs is the -x/--test-proxy flag, a semi-colon separated list
	testProxyURLs string
	// warmUpDuration is the -w/--warmup flag
	warmUpDuration int
	// parallelInstances is the -p/--parallel flag
	parallelInstances int

	// wg is used to keep track of the number of goroutines created
	wg sync.WaitGroup

	// number of processes to use, the --maxprocs flag
	numProcesses int

	// resourceTelemetry prints a runtime.MemStats / goroutine-count diff
	// at the end of the run when --resource-telemetry is set.
	resourceTelemetry bool

	// enableOperationLatency is the -l/--latency flag; mirrors PerfOptions.cs
	// Latency. When set, per-operation latency is tracked and a percentile
	// summary is printed at the end of the run (and per-operation entries are
	// written to --results-file).
	enableOperationLatency bool

	// resultsFilePath points to the file where per-operation results should be written
	resultsFilePath string

	// workloadConfigPath points to a workload configuration JSON file
	workloadConfigPath string

	// workloadName is the name of a workload entry in a workload configuration JSON file
	workloadName string

	// outputFilePrefix writes result artifacts to <prefix>.json/.csv/.txt/.md when set
	outputFilePrefix string

	// noCleanup skips per-test Cleanup() and GlobalCleanup() at the end of a run.
	// Mirrors --no-cleanup in other language perf runners; used by perf-automation
	// to leave server-side state in place between iterations.
	noCleanup bool

	// insecureSkipVerify is accepted for compatibility with other language perf
	// runners (--insecure). The Go perf default transport already skips TLS
	// verification to support the local test proxy, so this flag is a no-op.
	insecureSkipVerify bool

	// iterations is the -i/--iterations flag; the number of times the measurement
	// phase is repeated (warmup runs once before the first iteration). Matches the
	// .NET PerfOptions.Iterations option.
	iterations int

	// jobStatistics is the --job-statistics flag. When true the runner emits a
	// `#StartJobStatistics` / `#EndJobStatistics` JSON block consumed by
	// perf-automation, in the same format as the .NET runner.
	jobStatistics bool

	// targetRate is the -r/--rate flag, the target throughput in ops/sec aggregated
	// across all parallel workers. Zero (the default) means "unlimited".
	targetRate int

	// statusInterval is the --status-interval flag, the number of seconds between
	// status lines printed to the console. Matches PerfOptions.StatusInterval.
	statusInterval int

	// The following four flags exist purely for CLI compatibility with the .NET
	// perf runner's ThreadPool tuning options. The Go runtime has no equivalent
	// knobs (goroutines are not OS threads), so these flags are parsed and ignored.
	maxIOCompletionThreads int
	maxWorkerThreads       int
	minIOCompletionThreads int
	minWorkerThreads       int
)

// parseProxyURLs splits the --test-proxy input with the delimiter ';'
func parseProxyURLS() []string {
	if testProxyURLs == "" {
		return nil
	}

	testProxyURLs = strings.TrimSuffix(testProxyURLs, ";")

	return strings.Split(testProxyURLs, ";")
}
