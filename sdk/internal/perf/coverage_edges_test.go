// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/metrics"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPercentileBoundaries(t *testing.T) {
	cases := []struct {
		name string
		vals []time.Duration
		p    int
		want time.Duration
	}{
		{name: "empty", p: 50},
		{name: "single", vals: []time.Duration{5 * time.Millisecond}, p: 50, want: 5 * time.Millisecond},
		{name: "minimum", vals: []time.Duration{time.Millisecond, 2 * time.Millisecond, 3 * time.Millisecond}, p: 0, want: time.Millisecond},
		{name: "maximum", vals: []time.Duration{time.Millisecond, 2 * time.Millisecond, 3 * time.Millisecond}, p: 100, want: 3 * time.Millisecond},
		{name: "interpolated", vals: []time.Duration{time.Millisecond, 100 * time.Millisecond}, p: 50, want: 50*time.Millisecond + 500*time.Microsecond},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, percentile(test.vals, test.p))
		})
	}
}

func TestCollectorEdgeCases(t *testing.T) {
	latency := newLatencyCollector(-1)
	latency.Add(-time.Millisecond)
	latency.Add(0)
	require.Equal(t, int64(1), latency.seen)
	require.Empty(t, latency.durations)
	require.Equal(t, "Latency: no data", latency.Summary())

	source := newLatencyCollector(2)
	source.Add(time.Millisecond)
	destination := &latencyCollector{max: 2}
	destination.MergeFrom(source)
	require.NotNil(t, destination.rng)
	require.Equal(t, []time.Duration{time.Millisecond}, destination.durations)

	calls := newBoundedCallTypeLatencyCollector(-1)
	calls.Add("ignored", -time.Millisecond)
	calls.Add("", time.Millisecond)
	require.Equal(t, int64(1), calls.seen["operation"])
	require.Empty(t, calls.values["operation"])
	require.Equal(t, "Latency by call type: no data", newCallTypeLatencyCollector().Summary())

	results := newOperationResultsCollector(-1)
	results.Add("ignored", -time.Millisecond, 0)
	results.Add("", time.Millisecond, 1)
	require.Equal(t, int64(1), results.seen)
	require.Empty(t, results.results)

	unbounded := newOperationResultsCollector(0)
	unbounded.Add("source", time.Millisecond, 1)
	results.MergeFrom(unbounded)
	require.Equal(t, int64(2), results.seen)
	require.Len(t, results.results, 1)
}

func TestResourceTelemetry(t *testing.T) {
	require.Equal(t, 1.5, bytesToMiB(1572864))
	before := resourceTelemetrySnapshot{allocMiB: 1, totalAllocMiB: 2, sysMiB: 3, numGC: 4, goroutines: 5}
	after := resourceTelemetrySnapshot{allocMiB: 2.5, totalAllocMiB: 5, sysMiB: 7, numGC: 6, goroutines: 8}
	require.Equal(t, "Resource telemetry: alloc(MiB)=1.50 totalAlloc(MiB)=3.00 sys(MiB)=4.00 gc=2 goroutines=3", before.DiffSummary(after))
	snapshot := captureResourceTelemetry()
	require.Positive(t, snapshot.sysMiB)
	require.Positive(t, snapshot.goroutines)
}

func TestProcessStatsHelpers(t *testing.T) {
	sampler := newProcessStatsSampler(time.Hour)
	sampler.Stop()
	require.Equal(t, -1.0, sampler.AverageCPU())
	require.Equal(t, int64(-1), sampler.AverageMemory())
	cpu, memory := sampler.LastSample()
	require.Equal(t, -1.0, cpu)
	require.Zero(t, memory)
	require.Contains(t, sampler.Summary(), "samples=0")

	sampler.cpuSamples = []float64{10, 30}
	sampler.memorySamples = []uint64{1024, 3072}
	require.Equal(t, 20.0, sampler.AverageCPU())
	require.Equal(t, int64(2048), sampler.AverageMemory())
	cpu, memory = sampler.LastSample()
	require.Equal(t, 30.0, cpu)
	require.Equal(t, uint64(3072), memory)

	invalid := metrics.Sample{}
	require.False(t, func() bool { _, ok := readCPUSeconds(invalid); return ok }())
	require.False(t, func() bool { _, ok := readMemoryBytes(invalid); return ok }())
	require.False(t, func() bool { _, ok := readBusyCPUSeconds(invalid, invalid); return ok }())

	samples := []metrics.Sample{{Name: cpuTotalMetric}, {Name: cpuIdleMetric}, {Name: memoryTotalMetric}}
	metrics.Read(samples)
	total, ok := readCPUSeconds(samples[0])
	require.True(t, ok)
	require.GreaterOrEqual(t, total, 0.0)
	busy, ok := readBusyCPUSeconds(samples[0], invalid)
	require.True(t, ok)
	require.Equal(t, total, busy)
	_, ok = readMemoryBytes(samples[2])
	require.True(t, ok)
}

func TestParameterValueToStringAllTypes(t *testing.T) {
	cases := []struct {
		value any
		want  string
	}{
		{nil, ""}, {"text", "text"}, {true, "true"}, {json.Number("12.5"), "12.5"},
		{float64(12), "12"}, {float64(12.5), "12.5"}, {[]int{1, 2}, "[1 2]"},
	}
	for _, test := range cases {
		require.Equal(t, test.want, parameterValueToString(test.value))
	}
}

func TestRunSummaryRenderingAndArtifacts(t *testing.T) {
	defer snapshotFlags(t)()
	duration, warmUpDuration, parallelInstances, iterations = 10, 2, 4, 3
	workloadConfigPath, workloadName = "config.json", "upload"
	summary := newRunSummary("test", 50, 25, "latency", "calls", "resources")
	summary.AverageCPUPercent = 12.5
	summary.AverageMemoryBytes = 4096
	summary.ProcessStatsSummary = "process"
	require.Equal(t, 0.04, summary.SecondsPerOp)
	require.Equal(t, 2.0, summary.TotalElapsedSec)

	text := renderText(summary)
	markdown := renderMarkdown(summary)
	for _, expected := range []string{"Iterations: 3", "latency", "calls", "resources", "WorkloadConfig: config.json", "SelectedWorkload: upload", "process"} {
		require.Contains(t, text, expected)
	}
	for _, expected := range []string{"| Iterations | 3 |", "## Latency", "## By Call Type", "## Resource Telemetry", "## Process Stats"} {
		require.Contains(t, markdown, expected)
	}

	prefix := filepath.Join(t.TempDir(), "run")
	require.NoError(t, writeRunArtifacts(prefix, summary))
	for _, extension := range []string{".json", ".csv", ".txt", ".md"} {
		data, err := os.ReadFile(prefix + extension)
		require.NoError(t, err)
		require.NotEmpty(t, data)
	}
	require.NoError(t, writeRunArtifacts("", summary))

	workloadConfigPath, workloadName = "", ""
	empty := newRunSummary("empty", 0, 0, "", "", "")
	require.Zero(t, empty.SecondsPerOp)
	require.Zero(t, empty.TotalElapsedSec)
	require.NotContains(t, renderText(empty), "WorkloadConfig:")
	require.NotContains(t, renderMarkdown(empty), "## Latency")
}

func TestArtifactWriteErrors(t *testing.T) {
	badPrefix := filepath.Join(t.TempDir(), "missing", "run")
	require.Error(t, writeRunArtifacts(badPrefix, runSummary{}))
	badPath := filepath.Join(t.TempDir(), "missing", "artifact")
	require.Error(t, writeJSON(badPath, runSummary{}))
	require.Error(t, writeCSV(badPath, runSummary{}))
	require.Error(t, writeText(badPath, runSummary{}))
	require.Error(t, writeMarkdown(badPath, runSummary{}))
}

func TestReservoirMergeEdgeCases(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	values, seen := mergeReservoirs([]int{1}, 1, []int{2, 3}, 2, 0, rng)
	require.Equal(t, []int{1, 2, 3}, values)
	require.Equal(t, int64(3), seen)
	values, seen = mergeReservoirs(values, seen, nil, 5, 2, rng)
	require.Len(t, values, 2)
	require.Equal(t, int64(8), seen)
	require.True(t, strings.Contains(renderText(runSummary{}), "Test:"))
}
