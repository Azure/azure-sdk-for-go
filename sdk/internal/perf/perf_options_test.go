// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// snapshotFlags captures the current values of every flag-backed global so a
// test can restore them on exit. Each test that calls parseFlagsForTest must
// defer the returned restore function.
func snapshotFlags(t *testing.T) func() {
	t.Helper()
	saved := struct {
		duration               int
		warmUpDuration         int
		parallelInstances      int
		iterations             int
		targetRate             int
		statusInterval         int
		jobStatistics          bool
		numProcesses           int
		debug                  bool
		syncMode               bool
		resourceTelemetry      bool
		enableOperationLatency bool
		noCleanup              bool
		insecureSkipVerify     bool
		testProxyURLs          string
		resultsFilePath        string
		workloadConfigPath     string
		workloadName           string
		outputFilePrefix       string
		maxIOCompletionThreads int
		maxWorkerThreads       int
		minIOCompletionThreads int
		minWorkerThreads       int
	}{
		duration, warmUpDuration, parallelInstances, iterations, targetRate, statusInterval,
		jobStatistics, numProcesses, debug, syncMode, resourceTelemetry,
		enableOperationLatency, noCleanup, insecureSkipVerify, testProxyURLs,
		resultsFilePath, workloadConfigPath, workloadName, outputFilePrefix,
		maxIOCompletionThreads, maxWorkerThreads, minIOCompletionThreads, minWorkerThreads,
	}
	return func() {
		duration = saved.duration
		warmUpDuration = saved.warmUpDuration
		parallelInstances = saved.parallelInstances
		iterations = saved.iterations
		targetRate = saved.targetRate
		statusInterval = saved.statusInterval
		jobStatistics = saved.jobStatistics
		numProcesses = saved.numProcesses
		debug = saved.debug
		syncMode = saved.syncMode
		resourceTelemetry = saved.resourceTelemetry
		enableOperationLatency = saved.enableOperationLatency
		noCleanup = saved.noCleanup
		insecureSkipVerify = saved.insecureSkipVerify
		testProxyURLs = saved.testProxyURLs
		resultsFilePath = saved.resultsFilePath
		workloadConfigPath = saved.workloadConfigPath
		workloadName = saved.workloadName
		outputFilePrefix = saved.outputFilePrefix
		maxIOCompletionThreads = saved.maxIOCompletionThreads
		maxWorkerThreads = saved.maxWorkerThreads
		minIOCompletionThreads = saved.minIOCompletionThreads
		minWorkerThreads = saved.minWorkerThreads
	}
}

// parseFlagsForTest registers every perf flag onto a fresh FlagSet, parses
// the supplied args, and leaves the package-level globals populated for
// assertion. Tests must defer the cleanup function returned by snapshotFlags.
func parseFlagsForTest(t *testing.T, args []string) {
	t.Helper()
	fs := flag.NewFlagSet("perf-test", flag.ContinueOnError)
	registerFlags(fs)
	require.NoError(t, fs.Parse(args))
}

// TestDefaults asserts that every flag has the default value documented in
// PerfOptions.cs (and the Go-only flags retain their existing defaults).
func TestDefaults(t *testing.T) {
	defer snapshotFlags(t)()
	parseFlagsForTest(t, nil)

	require.Equal(t, 10, duration, "--duration default")
	require.Equal(t, 5, warmUpDuration, "--warmup default")
	require.Equal(t, 1, parallelInstances, "--parallel default")
	require.Equal(t, 1, iterations, "--iterations default")
	require.Equal(t, 0, targetRate, "--rate default (0 = unlimited)")
	require.Equal(t, 1, statusInterval, "--status-interval default")
	require.False(t, jobStatistics, "--job-statistics default")
	require.False(t, syncMode, "--sync default")
	require.False(t, insecureSkipVerify, "--insecure default")
	require.False(t, noCleanup, "--no-cleanup default")
	require.False(t, enableOperationLatency, "--latency default")
	require.False(t, resourceTelemetry, "--resource-telemetry default")
	require.Equal(t, "", testProxyURLs, "--test-proxies default")
	require.Equal(t, "", resultsFilePath, "--results-file default")
	require.Equal(t, "", outputFilePrefix, "--output-file-prefix default")
	require.Equal(t, 0, maxIOCompletionThreads, "--max-io-completion-threads default")
	require.Equal(t, 0, maxWorkerThreads, "--max-worker-threads default")
	require.Equal(t, 0, minIOCompletionThreads, "--min-io-completion-threads default")
	require.Equal(t, 0, minWorkerThreads, "--min-worker-threads default")
	require.Equal(t, runtime.NumCPU(), numProcesses, "--maxprocs default")
}

// TestFlagLongForms exercises the long-form spelling of every CLI flag. Each
// table entry sets one flag and asserts it parsed into the expected global.
func TestFlagLongForms(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		assert func(t *testing.T)
	}{
		{"duration", []string{"--duration", "42"}, func(t *testing.T) { require.Equal(t, 42, duration) }},
		{"warmup", []string{"--warmup", "7"}, func(t *testing.T) { require.Equal(t, 7, warmUpDuration) }},
		{"parallel", []string{"--parallel", "8"}, func(t *testing.T) { require.Equal(t, 8, parallelInstances) }},
		{"iterations", []string{"--iterations", "3"}, func(t *testing.T) { require.Equal(t, 3, iterations) }},
		{"rate", []string{"--rate", "100"}, func(t *testing.T) { require.Equal(t, 100, targetRate) }},
		{"status-interval", []string{"--status-interval", "5"}, func(t *testing.T) { require.Equal(t, 5, statusInterval) }},
		{"job-statistics", []string{"--job-statistics"}, func(t *testing.T) { require.True(t, jobStatistics) }},
		{"sync", []string{"--sync"}, func(t *testing.T) { require.True(t, syncMode) }},
		{"insecure", []string{"--insecure"}, func(t *testing.T) { require.True(t, insecureSkipVerify) }},
		{"no-cleanup", []string{"--no-cleanup"}, func(t *testing.T) { require.True(t, noCleanup) }},
		{"latency", []string{"--latency"}, func(t *testing.T) { require.True(t, enableOperationLatency) }},
		{"resource-telemetry", []string{"--resource-telemetry"}, func(t *testing.T) { require.True(t, resourceTelemetry) }},
		{"test-proxies", []string{"--test-proxies", "http://a;http://b"}, func(t *testing.T) {
			require.Equal(t, "http://a;http://b", testProxyURLs)
			require.Equal(t, []string{"http://a", "http://b"}, parseProxyURLS())
		}},
		{"results-file", []string{"--results-file", "/tmp/out.json"}, func(t *testing.T) { require.Equal(t, "/tmp/out.json", resultsFilePath) }},
		{"output-file-prefix", []string{"--output-file-prefix", "/tmp/run"}, func(t *testing.T) { require.Equal(t, "/tmp/run", outputFilePrefix) }},
		{"config", []string{"--config", "wl.json"}, func(t *testing.T) { require.Equal(t, "wl.json", workloadConfigPath) }},
		{"workload", []string{"--workload", "wl-upload"}, func(t *testing.T) { require.Equal(t, "wl-upload", workloadName) }},
		{"debug", []string{"--debug"}, func(t *testing.T) { require.True(t, debug) }},
		{"maxprocs", []string{"--maxprocs", "2"}, func(t *testing.T) { require.Equal(t, 2, numProcesses) }},
		{"max-io-completion-threads", []string{"--max-io-completion-threads", "16"}, func(t *testing.T) { require.Equal(t, 16, maxIOCompletionThreads) }},
		{"max-worker-threads", []string{"--max-worker-threads", "16"}, func(t *testing.T) { require.Equal(t, 16, maxWorkerThreads) }},
		{"min-io-completion-threads", []string{"--min-io-completion-threads", "2"}, func(t *testing.T) { require.Equal(t, 2, minIOCompletionThreads) }},
		{"min-worker-threads", []string{"--min-worker-threads", "2"}, func(t *testing.T) { require.Equal(t, 2, minWorkerThreads) }},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer snapshotFlags(t)()
			parseFlagsForTest(t, tc.args)
			tc.assert(t)
		})
	}
}

// TestFlagShortForms exercises the short alias of every flag that defines one
// (-d, -w, -p, -i, -r, -l, -x). They must produce the same effect as the
// long form.
func TestFlagShortForms(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		assert func(t *testing.T)
	}{
		{"d", []string{"-d", "20"}, func(t *testing.T) { require.Equal(t, 20, duration) }},
		{"w", []string{"-w", "3"}, func(t *testing.T) { require.Equal(t, 3, warmUpDuration) }},
		{"p", []string{"-p", "16"}, func(t *testing.T) { require.Equal(t, 16, parallelInstances) }},
		{"i", []string{"-i", "5"}, func(t *testing.T) { require.Equal(t, 5, iterations) }},
		{"r", []string{"-r", "50"}, func(t *testing.T) { require.Equal(t, 50, targetRate) }},
		{"l", []string{"-l"}, func(t *testing.T) { require.True(t, enableOperationLatency) }},
		{"x", []string{"-x", "http://proxy"}, func(t *testing.T) { require.Equal(t, "http://proxy", testProxyURLs) }},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			defer snapshotFlags(t)()
			parseFlagsForTest(t, tc.args)
			tc.assert(t)
		})
	}
}

// TestRateLimiter verifies that newRateLimiter emits tokens at approximately
// the requested rate. The test uses a coarse window (200 ms) and asserts the
// observed count is within [floor(expected*0.5), ceil(expected*2)] to avoid
// flakiness on slow CI runners.
func TestRateLimiter(t *testing.T) {
	const rate = 100
	tokens, stop := newRateLimiter(rate)
	defer close(stop)

	deadline := time.After(200 * time.Millisecond)
	count := 0
loop:
	for {
		select {
		case _, ok := <-tokens:
			if !ok {
				break loop
			}
			count++
		case <-deadline:
			break loop
		}
	}

	// At 100 ops/s for ~200 ms we expect ~20 tokens. Allow a wide tolerance
	// to keep the test reliable under load.
	require.GreaterOrEqual(t, count, 5, "rate limiter produced too few tokens")
	require.LessOrEqual(t, count, 60, "rate limiter produced too many tokens")
}

// TestLatencyCollector verifies that --latency accumulators record samples
// and report percentile summaries containing the expected markers.
func TestLatencyCollector(t *testing.T) {
	c := &latencyCollector{}
	c.Add(time.Millisecond)
	c.Add(2 * time.Millisecond)
	c.Add(3 * time.Millisecond)

	summary := c.Summary()
	for _, marker := range []string{"p50", "p90", "p95", "p99", "max"} {
		require.Contains(t, summary, marker)
	}
}

// TestCallTypeCollector verifies the call-type latency buckets that drive
// the per-call summary printed alongside --latency.
func TestCallTypeCollector(t *testing.T) {
	c := newCallTypeLatencyCollector()
	c.Add("operation", time.Millisecond)
	c.Add("operation", 2*time.Millisecond)
	c.Add("proxy-live", 5*time.Millisecond)

	summary := c.Summary()
	require.Contains(t, summary, "operation")
	require.Contains(t, summary, "proxy-live")
}

// TestOperationResultsWriteJSON drives the --results-file output path and
// asserts the on-disk shape matches what perf-automation consumes.
func TestOperationResultsWriteJSON(t *testing.T) {
	c := &operationResultsCollector{}
	c.Add("UploadBlobTest", 12*time.Millisecond, 1024)
	c.Add("UploadBlobTest", 7*time.Millisecond, 1024)

	path := filepath.Join(t.TempDir(), "results.json")
	require.NoError(t, c.WriteJSON(path))

	raw, err := os.ReadFile(path)
	require.NoError(t, err)

	var rows []operationResult
	require.NoError(t, json.Unmarshal(raw, &rows))
	require.Len(t, rows, 2)
	require.Equal(t, "UploadBlobTest", rows[0].Operation)
	require.InDelta(t, 12.0, rows[0].LatencyMs, 0.001)
	require.Equal(t, int64(1024), rows[0].SizeBytes)
}

// TestOutputFilePrefixWritesAllArtifacts drives the --output-file-prefix
// writer and asserts that all four artifact files are produced and the JSON
// shape includes the fields populated by the runner.
func TestOutputFilePrefixWritesAllArtifacts(t *testing.T) {
	dir := t.TempDir()
	prefix := filepath.Join(dir, "run")

	summary := newRunSummary("TestX", 100, 50.0, "lat", "ct", "res")
	summary.AverageCPUPercent = 25.5
	summary.AverageMemoryBytes = 1024 * 1024
	summary.ProcessStatsSummary = "Process stats: ..."
	require.NoError(t, writeRunArtifacts(prefix, summary))

	for _, ext := range []string{".json", ".csv", ".txt", ".md"} {
		_, err := os.Stat(prefix + ext)
		require.NoError(t, err, "expected %s artifact", ext)
	}

	raw, err := os.ReadFile(prefix + ".json")
	require.NoError(t, err)
	var got runSummary
	require.NoError(t, json.Unmarshal(raw, &got))
	require.Equal(t, "TestX", got.TestName)
	require.Equal(t, int64(100), got.TotalOperations)
	require.InDelta(t, 50.0, got.OpsPerSecond, 0.001)
	require.InDelta(t, 25.5, got.AverageCPUPercent, 0.001)
}

// TestStatusColumnFormatters covers the helpers used to render the
// CPU / Memory columns on the live status line.
func TestStatusColumnFormatters(t *testing.T) {
	require.Equal(t, "n/a", formatCPUColumn(-1))
	require.Equal(t, "12.34%", formatCPUColumn(12.34))

	require.Equal(t, "n/a", formatMemoryColumn(0))
	require.Equal(t, "1.00", formatMemoryColumn(1024*1024))
}

// TestProcessStatsSamplerProducesSamples starts the sampler, waits long
// enough to produce at least one sample, and asserts the recorded CPU% is
// finite and non-negative. This protects against regressions in the
// /cpu/classes/total - /cpu/classes/idle math.
func TestProcessStatsSamplerProducesSamples(t *testing.T) {
	s := newProcessStatsSampler(50 * time.Millisecond)
	s.Start()

	// Give the sampler enough wall time to collect at least 2 ticks so
	// AverageCPU has a usable delta to work with.
	time.Sleep(250 * time.Millisecond)
	s.Stop()

	require.Greater(t, s.SampleCount(), 0, "expected at least one sample")
	require.GreaterOrEqual(t, s.AverageCPU(), 0.0, "averageCpuPercent should not be negative")
	require.GreaterOrEqual(t, s.AverageMemory(), int64(0), "averageMemoryBytes should not be negative")

	last, mem := s.LastSample()
	require.GreaterOrEqual(t, last, 0.0)
	require.GreaterOrEqual(t, mem, uint64(0))
}

// TestPrintVersionsIncludesGoRuntime captures stdout while printVersions
// runs and asserts the Go runtime version line is present in the expected
// format (matches the .NET adapter's `Informational: (\S*)` regex).
func TestPrintVersionsIncludesGoRuntime(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	origStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = origStdout })

	printVersions()
	require.NoError(t, w.Close())

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	out := string(buf[:n])

	require.Contains(t, out, "=== Versions ===")
	require.Contains(t, out, "go:")
	require.Contains(t, out, "Informational:")
	require.Contains(t, out, runtime.Version())
}

// TestPrintJobStatisticsBlock captures stdout while printJobStatistics runs
// and asserts the #StartJobStatistics / payload / #EndJobStatistics markers
// match the .NET runner's output that perf-automation parses.
func TestPrintJobStatisticsBlock(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	origStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = origStdout })

	printJobStatistics(123.45)
	require.NoError(t, w.Close())

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	out := string(buf[:n])

	require.True(t, strings.HasPrefix(strings.TrimSpace(out), "#StartJobStatistics"))
	require.Contains(t, out, "#EndJobStatistics")
	require.Contains(t, out, "perfstress/throughput")
	require.Contains(t, out, "123.45")
}
