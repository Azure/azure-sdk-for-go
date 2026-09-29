// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"golang.org/x/text/message"
)

// optionsSlice is a way to access the options in a thread safe way
type optionsSlice struct {
	opts []*PerfTestOptions
	mu   sync.Mutex
}

// perfRunner is the per-process orchestrator that drives setup, warmup, one or
// more measurement iterations, and cleanup. The control-flow is modeled after
// the .NET runner (common/Perf/Azure.Test.Perf/PerfProgram.cs):
//
//	GlobalSetup -> Setup -> Warmup (if --warmup > 0)
//	            -> for i in 1..Iterations: RunTests
//	            -> Cleanup -> GlobalCleanup
type perfRunner struct {
	// ticker is the runner for giving updates every status-interval seconds
	ticker             *time.Ticker
	done               chan struct{}
	statusWG           sync.WaitGroup
	statusErr          chan error
	backgroundStopOnce sync.Once
	backgroundErr      error

	// name of the performance test
	name string

	// the perf test, options, and transports being tested/used
	perfToRun       PerfMethods
	allOptions      optionsSlice
	proxyTransports map[string]*RecordingHTTPClient
	activePlaybacks map[string]bool

	// All created tests
	tests []PerfTest

	// globalInstance is the single globalInstance for GlobalCleanup
	globalInstance GlobalPerfTest

	// this is the previous prints total
	warmupOperationStatusTracker int64
	operationStatusTracker       int64
	statusMu                     sync.Mutex

	// writer and messagePrinter
	w              *tabwriter.Writer
	messagePrinter *message.Printer

	// tracker for whether the warmup has finished
	warmupFinished int32
	warmupPrinted  bool

	latencyCollector  *latencyCollector
	callTypeCollector *callTypeLatencyCollector
	operationResults  *operationResultsCollector
	// mergeMu guards folding per-worker collectors into the shared collectors
	// above after each measurement phase completes.
	mergeMu               sync.Mutex
	resourceBeforeRun     resourceTelemetrySnapshot
	resourceAfterRun      resourceTelemetrySnapshot
	resourceTelemetryDone bool
	processStats          *processStatsSampler
}

func newPerfRunner(p PerfMethods, name string) *perfRunner {
	warmupFinished, warmupPrinted := 0, false
	if warmUpDuration == 0 {
		warmupFinished, warmupPrinted = parallelInstances, true
	}
	// The status ticker is created lazily inside Run() once the runner is
	// about to start emitting status lines. Allocating it here as well would
	// leak the original *time.Ticker (it is never Stop'd before being
	// overwritten by Run).
	return &perfRunner{
		done:                         make(chan struct{}),
		statusErr:                    make(chan error, 1),
		name:                         name,
		proxyTransports:              map[string]*RecordingHTTPClient{},
		activePlaybacks:              map[string]bool{},
		perfToRun:                    p,
		operationStatusTracker:       -1,
		warmupOperationStatusTracker: -1,
		w:                            tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight),
		messagePrinter:               message.NewPrinter(message.MatchLanguage("en")),
		warmupFinished:               int32(warmupFinished),
		warmupPrinted:                warmupPrinted,
		latencyCollector:             newLatencyCollector(defaultMaxLatencySamples),
		callTypeCollector:            newBoundedCallTypeLatencyCollector(defaultMaxLatencySamples),
		operationResults:             newOperationResultsCollector(maxResults),
	}
}

func (r *perfRunner) Run() (retErr error) {
	err := r.globalSetup()
	if err != nil {
		return err
	}
	defer func() {
		if noCleanup {
			return
		}
		retErr = errors.Join(retErr, r.globalInstance.GlobalCleanup(context.Background()))
	}()

	if err = r.createPerfTests(); err != nil {
		return errors.Join(err, r.cleanup())
	}
	defer func() {
		if !noCleanup {
			retErr = errors.Join(retErr, r.cleanup())
		}
	}()
	defer func() {
		retErr = errors.Join(retErr, r.stopProxyPlaybacks())
	}()

	// The process-stats sampler runs for the lifetime of every perf run.
	// Its output drives the CPU%/WorkingSet/PrivateMemory columns in the
	// status line and the final =Process stats= summary, mirroring the
	// .NET runner where the same data is always collected.
	r.processStats = newProcessStatsSampler(time.Second)
	r.processStats.Start()
	defer func() {
		retErr = errors.Join(retErr, r.stopBackground())
	}()

	statusTick := time.Duration(statusInterval) * time.Second
	if statusTick <= 0 {
		statusTick = time.Second
	}
	r.ticker = time.NewTicker(statusTick)
	// Ensure the ticker is released even if the run returns early (e.g. a
	// setup/bootstrap error short-circuits the function before r.done is
	// signaled). time.Ticker.Stop is safe to call multiple times.
	// Poller for printing
	r.statusWG.Add(1)
	go func() {
		defer r.statusWG.Done()
		for {
			select {
			case <-r.done:
				return
			case <-r.ticker.C:
				if err := r.printStatus(); err != nil {
					select {
					case r.statusErr <- err:
					default:
					}
					return
				}
			}
		}
	}()

	// Run any test-proxy bootstrap (live -> record -> playback) once per worker.
	if testProxyURLs != "" {
		if err = r.bootstrapProxies(); err != nil {
			return err
		}
	}

	if resourceTelemetry {
		r.resourceBeforeRun = captureResourceTelemetry()
	}

	// Warmup runs once before the first measurement iteration.
	if warmUpDuration > 0 {
		if err = r.runPhase(true, nil); err != nil {
			return err
		}
	}

	// Run the measurement phase `iterations` times. The Go --iterations /
	// -i flag mirrors the .NET PerfOptions.Iterations option.
	loopCount := iterations
	if loopCount < 1 {
		loopCount = 1
	}
	r.resetMeasurementState()
	for iter := 1; iter <= loopCount; iter++ {
		if loopCount > 1 {
			fmt.Printf("\n=== Test %d ===\n", iter)
		}
		var rateStop chan struct{}
		var rateCh <-chan struct{}
		if targetRate > 0 {
			rateCh, rateStop = newRateLimiter(targetRate)
		}
		err = r.runPhase(false, rateCh)
		if rateStop != nil {
			close(rateStop)
		}
		if err != nil {
			return err
		}
	}

	if err = r.stopBackground(); err != nil {
		return err
	}

	if err = r.printFinalUpdate(false); err != nil {
		return err
	}

	if resultsFilePath != "" && enableOperationLatency {
		if err = r.operationResults.WriteJSON(resultsFilePath); err != nil {
			return err
		}
	}

	if jobStatistics {
		printJobStatistics(r.opsPerSecond(false))
	}

	if err = r.writeArtifacts(); err != nil {
		return err
	}

	printVersions()

	if noCleanup {
		return nil
	}
	return nil
}

// runPhase spawns one goroutine per parallel worker that runs a single
// warmup or measurement pass, then waits for them all to complete.
func (r *perfRunner) runPhase(warmup bool, rateCh <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var phaseWG sync.WaitGroup
	errs := make(chan error, len(r.tests))
	for idx, test := range r.tests {
		phaseWG.Add(1)
		go func(p PerfTest, index int) {
			defer phaseWG.Done()
			if err := r.runWorkerPhase(ctx, p, index, warmup, rateCh); err != nil {
				errs <- err
				cancel()
			}
		}(test, idx)
	}
	phaseWG.Wait()
	close(errs)
	// Once a worker fails, cancel() tears down the shared context so in-flight
	// siblings return context.Canceled. Those are noise in front of the real
	// cause, so keep them only if no genuine error was recorded.
	var phaseErr error
	var canceledErr error
	for err := range errs {
		if errors.Is(err, context.Canceled) {
			canceledErr = err
			continue
		}
		phaseErr = errors.Join(phaseErr, err)
	}
	if phaseErr == nil {
		phaseErr = canceledErr
	}
	return phaseErr
}

// runWorkerPhase executes a single warmup or measurement pass for one worker.
func (r *perfRunner) runWorkerPhase(ctx context.Context, p PerfTest, index int, warmup bool, rateCh <-chan struct{}) error {
	r.allOptions.mu.Lock()
	opts := r.allOptions.opts[index]
	r.allOptions.mu.Unlock()

	// During the measurement phase each worker accumulates latency samples
	// into goroutine-local collectors so the hot loop never contends on a
	// shared mutex. The locals are merged into the shared collectors once,
	// after the loop completes. Warmup does not record latency samples.
	var localLatency *latencyCollector
	var localCallType *callTypeLatencyCollector
	var localResults *operationResultsCollector
	if !warmup {
		localLatency = newLatencyCollector(perWorkerSampleLimit(defaultMaxLatencySamples, len(r.tests), index))
		localCallType = newBoundedCallTypeLatencyCollector(perWorkerSampleLimit(defaultMaxLatencySamples, len(r.tests), index))
		localResults = newOperationResultsCollector(perWorkerSampleLimit(maxResults, len(r.tests), index))
	}

	if err := r.runTestForDuration(ctx, p, opts, warmup, rateCh, localLatency, localCallType, localResults); err != nil {
		return err
	}

	if !warmup {
		r.mergeMu.Lock()
		r.latencyCollector.MergeFrom(localLatency)
		r.callTypeCollector.MergeFrom(localCallType)
		r.operationResults.MergeFrom(localResults)
		r.mergeMu.Unlock()
	}

	if warmup {
		val := atomic.AddInt32(&r.warmupFinished, 1)
		if debug {
			fmt.Printf("finished %d warmups\n", val)
		}
	}
	return nil
}

// resetMeasurementState clears per-worker measurement counters and the
// shared latency/operation-result collectors so a fresh iteration starts
// with empty stats.
func (r *perfRunner) resetMeasurementState() {
	r.statusMu.Lock()
	defer r.statusMu.Unlock()

	r.allOptions.mu.Lock()
	for _, opt := range r.allOptions.opts {
		atomic.StoreInt64(&opt.runCount, 0)
		atomic.StoreInt64((*int64)(&opt.runElapsed), 0)
		t := time.Time{}
		opt.runStart = &t
		opt.finished = false
	}
	r.allOptions.mu.Unlock()
	r.operationStatusTracker = -1
	r.latencyCollector = newLatencyCollector(defaultMaxLatencySamples)
	r.callTypeCollector = newBoundedCallTypeLatencyCollector(defaultMaxLatencySamples)
	r.operationResults = newOperationResultsCollector(maxResults)
}

// bootstrapProxies performs the one-time live -> record -> playback dance
// against the test proxy for each worker so subsequent iterations replay
// from the in-memory recording.
func (r *perfRunner) bootstrapProxies() error {
	for idx := range r.tests {
		id := fmt.Sprintf("%s-%d", r.name, idx)
		p := r.tests[idx]

		r.proxyTransports[id].SetMode("live")
		start := time.Now()
		if err := p.Run(context.Background()); err != nil {
			return err
		}
		r.callTypeCollector.Add("proxy-live", time.Since(start))

		r.proxyTransports[id].SetMode("record")
		if err := r.proxyTransports[id].start(); err != nil {
			return err
		}
		start = time.Now()
		if err := p.Run(context.Background()); err != nil {
			return errors.Join(err, r.proxyTransports[id].stop())
		}
		r.callTypeCollector.Add("proxy-record", time.Since(start))
		if err := r.proxyTransports[id].stop(); err != nil {
			return err
		}

		r.proxyTransports[id].SetMode("playback")
		if err := r.proxyTransports[id].start(); err != nil {
			return err
		}
		r.activePlaybacks[id] = true
	}
	return nil
}

func (r *perfRunner) stopProxyPlaybacks() error {
	var stopErr error
	for id := range r.activePlaybacks {
		if err := r.proxyTransports[id].stop(); err != nil {
			stopErr = errors.Join(stopErr, err)
		}
		delete(r.activePlaybacks, id)
	}
	return stopErr
}

func (r *perfRunner) stopBackground() error {
	var err error
	r.backgroundStopOnce.Do(func() {
		if r.ticker != nil {
			r.ticker.Stop()
		}
		close(r.done)
		r.statusWG.Wait()
		if r.processStats != nil {
			r.processStats.Stop()
		}
		select {
		case r.backgroundErr = <-r.statusErr:
		default:
		}
		// Return the error only from the first call so it is not joined twice
		// (explicit call plus the deferred call in run()).
		err = r.backgroundErr
	})
	return err
}

// newRateLimiter starts a producer goroutine that emits one token every
// 1/rate seconds across a buffered channel shared by all workers. Workers
// receive one token before invoking each operation, throttling the
// aggregate throughput to approximately `rate` ops/sec. Closing the
// returned stop channel terminates the producer.
func newRateLimiter(rate int) (<-chan struct{}, chan struct{}) {
	// Clamp rate to a safe maximum so that a huge --rate value cannot:
	//   * compute an interval of 0ns (which would panic time.NewTicker), or
	//   * allocate a multi-GiB token channel buffer and OOM the process.
	// The ticker already gates production, so a small fixed buffer (1) is
	// sufficient regardless of the requested rate.
	const maxRate = 1_000_000_000 // 1e9 ops/sec is well beyond any realistic perf workload
	if rate > maxRate {
		rate = maxRate
	}
	interval := time.Second / time.Duration(rate)
	if interval < time.Nanosecond {
		interval = time.Nanosecond
	}
	tokens := make(chan struct{}, 1)
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				close(tokens)
				return
			case <-ticker.C:
				select {
				case tokens <- struct{}{}:
				default:
				}
			}
		}
	}()
	return tokens, stop
}

// global setup by instantiating a single global instance
func (r *perfRunner) globalSetup() error {
	globalInst, err := r.perfToRun.New(context.TODO(), newPerfTestOptions(r.name))
	if err != nil {
		return err
	}
	r.globalInstance = globalInst
	return nil
}

// createPerfTests spins up `parallelInstances` (specified by --parallel flag) PerfTest instances.
// The goroutines that drive them are spawned later, once per phase, by runPhase.
func (r *perfRunner) createPerfTests() error {
	proxyURLS := parseProxyURLS()

	for idx := 0; idx < parallelInstances; idx++ {
		ID := fmt.Sprintf("%s-%d", r.name, idx)
		options := newPerfTestOptions(ID)

		if testProxyURLs != "" {
			proxyURL := proxyURLS[idx%len(proxyURLS)]
			transporter := NewProxyTransport(&TransportOptions{
				TestName: ID,
				proxyURL: proxyURL,
			})
			options.Transporter = transporter
			r.proxyTransports[ID] = transporter
		} else {
			options.Transporter = defaultHTTPClient
		}
		options.parallelIndex = idx

		perfTest, err := r.globalInstance.NewPerfTest(context.TODO(), &options)
		if err != nil {
			return err
		}
		r.tests = append(r.tests, perfTest)
		r.allOptions.mu.Lock()
		r.allOptions.opts = append(r.allOptions.opts, &options)
		r.allOptions.mu.Unlock()
	}
	return nil
}

// cleanup runs the Cleanup on each of the r.tests
func (r *perfRunner) cleanup() error {
	var cleanupErr error
	for _, t := range r.tests {
		cleanupErr = errors.Join(cleanupErr, t.Cleanup(context.Background()))
	}
	return cleanupErr
}

// print an update for the last second
func (r *perfRunner) printStatus() error {
	r.statusMu.Lock()
	defer r.statusMu.Unlock()

	if !r.warmupPrinted {
		if finishedWarmup, err := r.printWarmupStatus(); err != nil {
			return err
		} else if !finishedWarmup {
			return nil
		}
	}

	if r.operationStatusTracker == -1 {
		err := r.printFinalUpdate(true)
		if err != nil {
			return err
		}
		r.operationStatusTracker = 0
		if _, err := fmt.Fprintln(r.w, "\nCurrent\tTotal\tAverage\tCPU\tMemory(MiB)\t"); err != nil {
			return err
		}
	}
	totalOperations := r.totalOperations(false)
	cpuPct, memBytes := r.lastProcessSample()

	_, err := fmt.Fprintf(
		r.w,
		"%s\t%s\t%s\t%s\t%s\t\n",
		r.messagePrinter.Sprintf("%d", totalOperations-r.operationStatusTracker),
		r.messagePrinter.Sprintf("%d", totalOperations),
		r.messagePrinter.Sprintf("%.2f", r.opsPerSecond(false)),
		formatCPUColumn(cpuPct),
		formatMemoryColumn(memBytes),
	)
	if err != nil {
		return err
	}
	r.operationStatusTracker = totalOperations
	return r.w.Flush()
}

// return true if all warmup information has been printed
func (r *perfRunner) printWarmupStatus() (bool, error) {
	if r.warmupOperationStatusTracker == -1 {
		r.warmupOperationStatusTracker = 0
		fmt.Println("===== WARMUP =====")
		if _, err := fmt.Fprintln(r.w, "\nCurrent\tTotal\tAverage\tCPU\tMemory(MiB)\t"); err != nil {
			return false, err
		}
	}
	totalOperations := r.totalOperations(true)
	cpuPct, memBytes := r.lastProcessSample()

	// Warmup is finished only when every goroutine has returned from its
	// warmup loop (runTest increments warmupFinished on completion). The
	// previous heuristic — "no new operations in the last status tick" —
	// fires spuriously for long-running iterations (e.g. multi-GiB blob
	// transfers) where a single op exceeds the 1s status interval and
	// would cause the warmup phase to terminate after only one second with
	// zero recorded operations.
	if atomic.LoadInt32(&r.warmupFinished) >= int32(parallelInstances) {
		// Drain one last status line so the user sees the final counters
		// before "=== Warm Up Results ===" is printed.
		if totalOperations != r.warmupOperationStatusTracker {
			_, err := fmt.Fprintf(
				r.w,
				"%s\t%s\t%s\t%s\t%s\t\n",
				r.messagePrinter.Sprintf("%d", totalOperations-r.warmupOperationStatusTracker),
				r.messagePrinter.Sprintf("%d", totalOperations),
				r.messagePrinter.Sprintf("%.2f", r.opsPerSecond(true)),
				formatCPUColumn(cpuPct),
				formatMemoryColumn(memBytes),
			)
			if err != nil {
				return false, err
			}
			r.warmupOperationStatusTracker = totalOperations
			if err := r.w.Flush(); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	_, err := fmt.Fprintf(
		r.w,
		"%s\t%s\t%s\t%s\t%s\t\n",
		r.messagePrinter.Sprintf("%d", totalOperations-r.warmupOperationStatusTracker),
		r.messagePrinter.Sprintf("%d", totalOperations),
		r.messagePrinter.Sprintf("%.2f", r.opsPerSecond(true)),
		formatCPUColumn(cpuPct),
		formatMemoryColumn(memBytes),
	)
	if err != nil {
		return false, err
	}
	r.warmupOperationStatusTracker = totalOperations
	err = r.w.Flush()
	if err != nil {
		return false, err
	}
	return false, nil
}

// totalOperations iterates over all options structs to get the number of operations completed
func (r *perfRunner) totalOperations(warmup bool) int64 {
	var ret int64

	r.allOptions.mu.Lock()
	defer r.allOptions.mu.Unlock()
	for _, opt := range r.allOptions.opts {
		if warmup {
			ret += atomic.LoadInt64(&opt.warmupCount)
		} else {
			ret += atomic.LoadInt64(&opt.runCount)
		}
	}

	return ret
}

// opsPerSecond calculates the average number of operations per second
func (r *perfRunner) opsPerSecond(warmup bool) float64 {
	var ret float64

	r.allOptions.mu.Lock()
	defer r.allOptions.mu.Unlock()
	for _, opt := range r.allOptions.opts {
		if warmup {
			e := float64(atomic.LoadInt64((*int64)(&opt.warmupElapsed))) / float64(time.Second)
			if e != 0 {
				ret += float64(atomic.LoadInt64(&opt.warmupCount)) / e
			}
		} else {
			e := float64(atomic.LoadInt64((*int64)(&opt.runElapsed))) / float64(time.Second)
			if e != 0 {
				ret += float64(atomic.LoadInt64(&opt.runCount)) / e
			}
		}
	}
	return ret
}

// printFinalUpdate prints the final update for the warmup/test run
func (r *perfRunner) printFinalUpdate(warmup bool) error {
	if r.warmupPrinted && warmup {
		return nil
	}
	totalOperations := r.totalOperations(warmup)
	opsPerSecond := r.opsPerSecond(warmup)
	if opsPerSecond == 0.0 {
		// Zero completed operations means the elapsed window was too short
		// to fit even one iteration (e.g. uploading a multi-GiB blob with a
		// 60 s --warmup). For the warmup phase this is a soft condition —
		// warn and continue so the measurement phase still has a chance to
		// run (with a longer --duration, or after retries succeed). For the
		// run phase it is a real failure: nothing useful can be reported.
		if warmup {
			fmt.Println("\n=== Warm Up Results ===")
			fmt.Printf(
				"warning: warmup completed without generating operation statistics. "+
					"Consider increasing --warmup beyond the per-iteration latency "+
					"(--warmup is currently %d s).\n",
				warmUpDuration,
			)
			r.warmupPrinted = true
			return nil
		}
		return fmt.Errorf(
			"completed without generating operation statistics: no iteration finished in --duration %d s. "+
				"Increase --duration, reduce --size, or use a chunked method (--upload-method buffer|stream, --download-method buffer)",
			duration,
		)
	}

	secondsPerOp := 1.0 / opsPerSecond
	weightedAvg := float64(totalOperations) / opsPerSecond

	if warmup {
		fmt.Println("\n=== Warm Up Results ===")
	} else {
		fmt.Println("\n=== Results ===")
	}
	fmt.Printf(
		"Completed %s operations in a weighted-average of %ss (%s ops/s, %s s/op)\n",
		r.messagePrinter.Sprintf("%d", totalOperations),
		r.messagePrinter.Sprintf("%.2f", weightedAvg),
		r.messagePrinter.Sprintf("%.3f", opsPerSecond),
		r.messagePrinter.Sprintf("%.3f", secondsPerOp),
	)

	if !warmup {
		if enableOperationLatency {
			fmt.Println(r.latencyCollector.Summary())
			fmt.Println(r.callTypeCollector.Summary())
		}
		if resourceTelemetry {
			if !r.resourceTelemetryDone {
				r.resourceAfterRun = captureResourceTelemetry()
				r.resourceTelemetryDone = true
			}
			fmt.Println(r.resourceBeforeRun.DiffSummary(r.resourceAfterRun))
		}
		if r.processStats != nil && r.processStats.SampleCount() > 0 {
			fmt.Println(r.processStats.Summary())
		}
	}

	return nil
}

// runTestForDuration runs a single warmup or measurement pass for one worker
// for the configured duration. When rateCh is non-nil each iteration waits to
// receive a token from it before invoking p.Run, throttling the aggregate
// throughput to the value of --rate. During the measurement phase latency
// samples are recorded into the goroutine-local collectors passed in by the
// caller (nil during warmup) so the hot loop never contends on shared state.
func (r *perfRunner) runTestForDuration(parent context.Context, p PerfTest, opts *PerfTestOptions, warmup bool, rateCh <-chan struct{}, localLatency *latencyCollector, localCallType *callTypeLatencyCollector, localResults *operationResultsCollector) error {
	if warmup && warmUpDuration <= 0 {
		return nil
	}

	// startPtr is our base time for keeping track of how long a test has run
	var startPtr *time.Time
	var previousElapsed time.Duration
	if warmup {
		t := time.Now()
		opts.warmupStart = &t
		startPtr = opts.warmupStart
	} else {
		t := time.Now()
		opts.runStart = &t
		startPtr = opts.runStart
		previousElapsed = time.Duration(atomic.LoadInt64((*int64)(&opts.runElapsed)))
	}

	var runDuration int
	if warmup {
		runDuration = warmUpDuration
	} else {
		runDuration = duration
	}

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(parent, time.Second*time.Duration(runDuration))
	defer cancel()

	lastSavedTime := time.Now()
	for time.Since(*startPtr).Seconds() < float64(runDuration) {
		if rateCh != nil {
			select {
			case _, ok := <-rateCh:
				if !ok {
					return nil
				}
			case <-ctx.Done():
				return nil
			}
		}
		operationStart := time.Now()
		err := p.Run(ctx)
		operationDuration := time.Since(operationStart)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			} else {
				return err
			}
		}
		opts.increment(warmup)

		if !warmup {
			localLatency.Add(operationDuration)
			localCallType.Add("operation", operationDuration)
			if resultsFilePath != "" && enableOperationLatency {
				localResults.Add(r.name, operationDuration, 0)
			}
		}

		if time.Since(lastSavedTime).Seconds() > 0.3 {
			elapsed := time.Since(*startPtr)
			if warmup {
				atomic.StoreInt64((*int64)(&opts.warmupElapsed), int64(elapsed))
			} else {
				atomic.StoreInt64((*int64)(&opts.runElapsed), int64(previousElapsed+elapsed))
			}
			lastSavedTime = time.Now()
		}
	}

	elapsed := time.Since(*startPtr)
	if warmup {
		atomic.StoreInt64((*int64)(&opts.warmupElapsed), int64(elapsed))
	} else {
		atomic.StoreInt64((*int64)(&opts.runElapsed), int64(previousElapsed+elapsed))
	}

	opts.finished = true
	return nil
}

func (r *perfRunner) writeArtifacts() error {
	if outputFilePrefix == "" {
		return nil
	}

	totalOperations := r.totalOperations(false)
	opsPerSecond := r.opsPerSecond(false)

	latencySummary := ""
	callTypeSummary := ""
	resourceSummary := ""

	if enableOperationLatency {
		latencySummary = r.latencyCollector.Summary()
		callTypeSummary = r.callTypeCollector.Summary()
	}
	if resourceTelemetry {
		if !r.resourceTelemetryDone {
			r.resourceAfterRun = captureResourceTelemetry()
			r.resourceTelemetryDone = true
		}
		resourceSummary = r.resourceBeforeRun.DiffSummary(r.resourceAfterRun)
	}

	avgCPU := -1.0
	avgMem := int64(-1)
	processStatsSummary := ""
	if r.processStats != nil && r.processStats.SampleCount() > 0 {
		avgCPU = r.processStats.AverageCPU()
		avgMem = r.processStats.AverageMemory()
		processStatsSummary = r.processStats.Summary()
	}

	summary := newRunSummary(r.name, totalOperations, opsPerSecond, latencySummary, callTypeSummary, resourceSummary)
	summary.AverageCPUPercent = avgCPU
	summary.AverageMemoryBytes = avgMem
	summary.ProcessStatsSummary = processStatsSummary
	if err := writeRunArtifacts(outputFilePrefix, summary); err != nil {
		return err
	}

	fmt.Printf("Wrote result artifacts: %s.json/.csv/.txt/.md\n", outputFilePrefix)
	return nil
}

// printJobStatistics writes the `#StartJobStatistics` / `#EndJobStatistics`
// block consumed by perf-automation. The JSON shape mirrors the .NET
// runner's BenchmarkOutput payload (Source, Name, ShortDescription,
// LongDescription, Format, Measurements[]) so a single parser can read
// throughput from either language's stdout.
func printJobStatistics(opsPerSecond float64) {
	type metadata struct {
		Source           string `json:"Source"`
		Name             string `json:"Name"`
		ShortDescription string `json:"ShortDescription"`
		LongDescription  string `json:"LongDescription"`
		Format           string `json:"Format"`
	}
	type measurement struct {
		Timestamp string  `json:"Timestamp"`
		Name      string  `json:"Name"`
		Value     float64 `json:"Value"`
	}
	type output struct {
		Metadata     []metadata    `json:"Metadata"`
		Measurements []measurement `json:"Measurements"`
	}
	out := output{
		Metadata: []metadata{{
			Source:           "PerfStress",
			Name:             "perfstress/throughput",
			ShortDescription: "Throughput (ops/sec)",
			LongDescription:  "Throughput (ops/sec)",
			Format:           "n2",
		}},
		Measurements: []measurement{{
			Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
			Name:      "perfstress/throughput",
			Value:     opsPerSecond,
		}},
	}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	fmt.Println("#StartJobStatistics")
	fmt.Println(string(b))
	fmt.Println("#EndJobStatistics")
}

// lastProcessSample returns the most recent CPU% and memory-in-bytes from the
// process-stats sampler, or sentinel values when the sampler is disabled or
// has not yet produced a sample.
func (r *perfRunner) lastProcessSample() (float64, uint64) {
	if r.processStats == nil {
		return -1, 0
	}
	return r.processStats.LastSample()
}

// formatCPUColumn renders the CPU% column for the status line. A negative
// value (no sample yet) is displayed as "n/a".
func formatCPUColumn(cpuPercent float64) string {
	if cpuPercent < 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.2f%%", cpuPercent)
}

// formatMemoryColumn renders the Memory(MiB) column for the status line.
// Zero bytes (no sample yet) is displayed as "n/a".
func formatMemoryColumn(memoryBytes uint64) string {
	if memoryBytes == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.2f", float64(memoryBytes)/(1024*1024))
}
