// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type lifecycleGlobal struct {
	runErr        error
	cleanupErr    error
	globalErr     error
	newErrAt      int
	created       atomic.Int32
	runs          atomic.Int64
	cleanups      atomic.Int32
	globalCleanup atomic.Int32
}

func (g *lifecycleGlobal) NewPerfTest(context.Context, *PerfTestOptions) (PerfTest, error) {
	created := int(g.created.Add(1))
	if g.newErrAt > 0 && created == g.newErrAt {
		return nil, errCreatePerfTest
	}
	return &lifecycleTest{global: g}, nil
}

func (g *lifecycleGlobal) GlobalCleanup(context.Context) error {
	g.globalCleanup.Add(1)
	return g.globalErr
}

type lifecycleTest struct {
	global *lifecycleGlobal
}

type phaseErrorTest struct{}

func (phaseErrorTest) Run(context.Context) error     { return errRunPerfTest }
func (phaseErrorTest) Cleanup(context.Context) error { return nil }

type phaseWaitingTest struct {
	canceled atomic.Bool
}

func (t *phaseWaitingTest) Run(ctx context.Context) error {
	<-ctx.Done()
	t.canceled.Store(true)
	return ctx.Err()
}

func (*phaseWaitingTest) Cleanup(context.Context) error { return nil }

func (t *lifecycleTest) Run(context.Context) error {
	t.global.runs.Add(1)
	return t.global.runErr
}

func (t *lifecycleTest) Cleanup(context.Context) error {
	t.global.cleanups.Add(1)
	return t.global.cleanupErr
}

var (
	errCreatePerfTest = errors.New("create perf test failed")
	errRunPerfTest    = errors.New("run perf test failed")
	errCleanupTest    = errors.New("cleanup failed")
	errGlobalCleanup  = errors.New("global cleanup failed")
)

func configureRunnerTest(t *testing.T) {
	t.Helper()
	deferRestore := snapshotFlags(t)
	t.Cleanup(deferRestore)
	duration = 1
	warmUpDuration = 0
	parallelInstances = 1
	iterations = 1
	targetRate = 0
	statusInterval = 60
	noCleanup = false
	testProxyURLs = ""
	enableOperationLatency = false
	resultsFilePath = ""
	outputFilePrefix = ""
	resourceTelemetry = false
	jobStatistics = false
	maxResults = defaultMaxOperationResults
}

func newLifecycleRunner(global *lifecycleGlobal) *perfRunner {
	return newPerfRunner(PerfMethods{New: func(context.Context, PerfTestOptions) (GlobalPerfTest, error) {
		return global, nil
	}}, "LifecycleTest")
}

func TestRunnerAggregatesIterations(t *testing.T) {
	configureRunnerTest(t)
	iterations = 2
	global := &lifecycleGlobal{}
	runner := newLifecycleRunner(global)

	require.NoError(t, runner.Run())
	require.Greater(t, global.runs.Load(), int64(1))
	require.Equal(t, int32(1), global.cleanups.Load())
	require.Equal(t, int32(1), global.globalCleanup.Load())
	runner.allOptions.mu.Lock()
	elapsed := time.Duration(atomic.LoadInt64((*int64)(&runner.allOptions.opts[0].runElapsed)))
	runner.allOptions.mu.Unlock()
	require.GreaterOrEqual(t, elapsed, 1900*time.Millisecond)
}

func TestRunnerReturnsWorkerErrorAndCleansUp(t *testing.T) {
	configureRunnerTest(t)
	global := &lifecycleGlobal{runErr: errRunPerfTest}
	runner := newLifecycleRunner(global)

	err := runner.Run()

	require.ErrorIs(t, err, errRunPerfTest)
	require.Equal(t, int32(1), global.cleanups.Load())
	require.Equal(t, int32(1), global.globalCleanup.Load())
	select {
	case <-runner.done:
	default:
		t.Fatal("runner background goroutines were not stopped")
	}
}

func TestRunnerJoinsCleanupErrors(t *testing.T) {
	configureRunnerTest(t)
	global := &lifecycleGlobal{runErr: errRunPerfTest, cleanupErr: errCleanupTest, globalErr: errGlobalCleanup}

	err := newLifecycleRunner(global).Run()

	require.ErrorIs(t, err, errRunPerfTest)
	require.ErrorIs(t, err, errCleanupTest)
	require.ErrorIs(t, err, errGlobalCleanup)
}

func TestRunnerCleansEveryWorkerAfterErrors(t *testing.T) {
	configureRunnerTest(t)
	parallelInstances = 3
	global := &lifecycleGlobal{runErr: errRunPerfTest, cleanupErr: errCleanupTest}

	err := newLifecycleRunner(global).Run()

	require.ErrorIs(t, err, errRunPerfTest)
	require.ErrorIs(t, err, errCleanupTest)
	require.Equal(t, int32(3), global.cleanups.Load())
}

func TestRunnerCleansUpPartiallyCreatedTests(t *testing.T) {
	configureRunnerTest(t)
	parallelInstances = 3
	global := &lifecycleGlobal{newErrAt: 2}

	err := newLifecycleRunner(global).Run()

	require.ErrorIs(t, err, errCreatePerfTest)
	require.Equal(t, int32(1), global.cleanups.Load())
	require.Equal(t, int32(1), global.globalCleanup.Load())
}

func TestRunnerNoCleanup(t *testing.T) {
	configureRunnerTest(t)
	noCleanup = true
	global := &lifecycleGlobal{}

	require.NoError(t, newLifecycleRunner(global).Run())
	require.Zero(t, global.cleanups.Load())
	require.Zero(t, global.globalCleanup.Load())
}

func TestRunPhaseCancelsSiblingWorkers(t *testing.T) {
	configureRunnerTest(t)
	duration = 60
	waiting := &phaseWaitingTest{}
	runner := newPerfRunner(PerfMethods{}, "PhaseTest")
	runner.tests = []PerfTest{phaseErrorTest{}, waiting}
	first, second := newPerfTestOptions("first"), newPerfTestOptions("second")
	runner.allOptions.opts = []*PerfTestOptions{&first, &second}

	started := time.Now()
	err := runner.runPhase(false, nil)

	require.ErrorIs(t, err, errRunPerfTest)
	require.True(t, waiting.canceled.Load())
	require.Less(t, time.Since(started), time.Second)
}
