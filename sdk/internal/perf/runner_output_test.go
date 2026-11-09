// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T, action func()) string {
	t.Helper()
	read, write, err := os.Pipe()
	require.NoError(t, err)
	saved := os.Stdout
	os.Stdout = write
	action()
	require.NoError(t, write.Close())
	os.Stdout = saved
	data, err := io.ReadAll(read)
	require.NoError(t, err)
	require.NoError(t, read.Close())
	return string(data)
}

func runnerWithOutput(t *testing.T) (*perfRunner, *bytes.Buffer, *PerfTestOptions) {
	t.Helper()
	runner := newPerfRunner(PerfMethods{}, "OutputTest")
	buffer := &bytes.Buffer{}
	runner.w = tabwriter.NewWriter(buffer, 16, 8, 1, ' ', tabwriter.AlignRight)
	options := newPerfTestOptions("worker")
	runner.allOptions.opts = []*PerfTestOptions{&options}
	return runner, buffer, &options
}

func TestRunnerWarmupAndStatusOutput(t *testing.T) {
	configureRunnerTest(t)
	warmUpDuration = 1
	parallelInstances = 1
	runner, output, options := runnerWithOutput(t)
	atomic.StoreInt64(&options.warmupCount, 2)
	atomic.StoreInt64((*int64)(&options.warmupElapsed), int64(time.Second))

	finished, err := runner.printWarmupStatus()
	require.NoError(t, err)
	require.False(t, finished)
	require.Contains(t, output.String(), "Current")

	atomic.StoreInt32(&runner.warmupFinished, 1)
	atomic.StoreInt64(&options.warmupCount, 3)
	finished, err = runner.printWarmupStatus()
	require.NoError(t, err)
	require.True(t, finished)

	runner.warmupPrinted = true
	atomic.StoreInt64(&options.runCount, 4)
	atomic.StoreInt64((*int64)(&options.runElapsed), int64(2*time.Second))
	require.NoError(t, runner.printStatus())
	require.NoError(t, runner.printStatus())
	require.Contains(t, output.String(), "Average")
	require.Equal(t, int64(4), runner.operationStatusTracker)
}

func TestRunnerFinalUpdateBranches(t *testing.T) {
	configureRunnerTest(t)
	runner, _, options := runnerWithOutput(t)

	require.NoError(t, runner.printFinalUpdate(true))
	require.True(t, runner.warmupPrinted)
	require.NoError(t, runner.printFinalUpdate(true))
	require.ErrorContains(t, runner.printFinalUpdate(false), "without generating operation statistics")

	runner.warmupPrinted = false
	atomic.StoreInt64(&options.runCount, 2)
	atomic.StoreInt64((*int64)(&options.runElapsed), int64(time.Second))
	enableOperationLatency = true
	resourceTelemetry = true
	runner.latencyCollector.Add(time.Millisecond)
	runner.callTypeCollector.Add("operation", time.Millisecond)
	runner.processStats = newProcessStatsSampler(time.Hour)
	runner.processStats.cpuSamples = []float64{10}
	runner.processStats.memorySamples = []uint64{1024}
	require.NoError(t, runner.printFinalUpdate(false))
	require.True(t, runner.resourceTelemetryDone)
}

func TestRunnerWritesArtifactsWithOptionalMetrics(t *testing.T) {
	configureRunnerTest(t)
	runner, _, options := runnerWithOutput(t)
	atomic.StoreInt64(&options.runCount, 2)
	atomic.StoreInt64((*int64)(&options.runElapsed), int64(time.Second))
	enableOperationLatency = true
	resourceTelemetry = true
	runner.latencyCollector.Add(time.Millisecond)
	runner.callTypeCollector.Add("operation", time.Millisecond)
	runner.processStats = newProcessStatsSampler(time.Hour)
	runner.processStats.cpuSamples = []float64{25}
	runner.processStats.memorySamples = []uint64{2048}
	outputFilePrefix = filepath.Join(t.TempDir(), "runner")

	require.NoError(t, runner.writeArtifacts())
	require.True(t, runner.resourceTelemetryDone)
	data, err := os.ReadFile(outputFilePrefix + ".json")
	require.NoError(t, err)
	require.Contains(t, string(data), `"averageCpuPercent": 25`)
	require.Contains(t, string(data), `"latencySummary"`)

	outputFilePrefix = filepath.Join(t.TempDir(), "missing", "runner")
	require.Error(t, runner.writeArtifacts())
	outputFilePrefix = ""
	require.NoError(t, runner.writeArtifacts())
}

func TestOptionsAndDispatchOutput(t *testing.T) {
	configureRunnerTest(t)
	defer snapshotFlags(t)()
	parallelInstances = 0
	require.Equal(t, 1, Parallel())
	parallelInstances = 3
	require.Equal(t, 3, Parallel())

	output := captureStdout(t, func() { printOptions("OutputTest") })
	require.Contains(t, output, "=== Options ===")
	require.Contains(t, output, `"testName": "OutputTest"`)

	savedArgs, savedFlags := os.Args, flag.CommandLine
	t.Cleanup(func() { os.Args, flag.CommandLine = savedArgs, savedFlags })
	flag.CommandLine = flag.NewFlagSet("dispatch", flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)
	os.Args = []string{"perf"}
	output = captureStdout(t, func() { Run(map[string]PerfMethods{}) })
	require.Contains(t, output, "Available performance tests:")

	os.Args = []string{"perf", "MissingTest"}
	output = captureStdout(t, func() { Run(map[string]PerfMethods{"KnownTest": {}}) })
	require.Contains(t, output, "KnownTest")
}

func TestStatusFormattingAndSamples(t *testing.T) {
	require.Equal(t, "n/a", formatCPUColumn(-1))
	require.Equal(t, "12.50%", formatCPUColumn(12.5))
	require.Equal(t, "n/a", formatMemoryColumn(0))
	require.Equal(t, "1.50", formatMemoryColumn(1572864))
	runner := newPerfRunner(PerfMethods{}, "samples")
	cpu, memory := runner.lastProcessSample()
	require.Equal(t, -1.0, cpu)
	require.Zero(t, memory)
	runner.processStats = newProcessStatsSampler(time.Hour)
	runner.processStats.cpuSamples = []float64{1}
	runner.processStats.memorySamples = []uint64{2}
	cpu, memory = runner.lastProcessSample()
	require.Equal(t, 1.0, cpu)
	require.Equal(t, uint64(2), memory)
}
