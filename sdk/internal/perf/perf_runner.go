// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"errors"
	"fmt"
	"log"
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

type perfRunner struct {
	// ticker is the runner for giving updates every second
	ticker *time.Ticker
	done   chan bool

	// name of the performance test
	name string

	// the perf test, options, and transports being tested/used
	perfToRun       PerfMethods
	allOptions      optionsSlice
	proxyTransports map[string]*RecordingHTTPClient

	// All created tests
	tests []PerfTest

	// globalInstance is the single globalInstance for GlobalCleanup
	globalInstance GlobalPerfTest

	// this is the previous prints total
	warmupOperationStatusTracker int64
	operationStatusTracker       int64

	// writer and messagePrinter
	w              *tabwriter.Writer
	messagePrinter *message.Printer

	// tracker for whether the warmup has finished
	warmupFinished int32
	warmupPrinted  bool
}

func newPerfRunner(p PerfMethods, name string) *perfRunner {
	warmupFinished, warmupPrinted := 0, false
	if warmUpDuration == 0 {
		warmupFinished, warmupPrinted = parallelInstances, true
	}
	return &perfRunner{
		ticker:                       time.NewTicker(time.Second),
		done:                         make(chan bool),
		name:                         name,
		proxyTransports:              map[string]*RecordingHTTPClient{},
		perfToRun:                    p,
		operationStatusTracker:       -1,
		warmupOperationStatusTracker: -1,
		w:                            tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight),
		messagePrinter:               message.NewPrinter(message.MatchLanguage("en")),
		warmupFinished:               int32(warmupFinished),
		warmupPrinted:                warmupPrinted,
	}
}

func (r *perfRunner) Run() error {
	err := r.globalSetup()
	if err != nil {
		return err
	}
	defer func() {
		err = r.globalInstance.GlobalCleanup(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	err = r.createPerfTests()
	if err != nil {
		return err
	}

	r.ticker = time.NewTicker(time.Second)

	// Poller for printing
	go func() {
		for {
			if r.ticker != nil {
				select {
				case <-r.done:
					return
				case <-r.ticker.C:
					err := r.printStatus()
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}()
	wg.Wait()

	r.done <- true

	err = r.printFinalUpdate(false)
	if err != nil {
		return err
	}
	return r.cleanup()
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

// createPerfTests spins up `parallelInstances` (specified by --parallel flag) goroutines
func (r *perfRunner) createPerfTests() error {
	var IDs []string
	proxyURLS := parseProxyURLS()

	for idx := 0; idx < parallelInstances; idx++ {
		ID := fmt.Sprintf("%s-%d", r.name, idx)
		options := newPerfTestOptions(ID)
		IDs = append(IDs, ID)

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

	for idx, test := range r.tests {
		wg.Add(1)
		go r.runTest(test, idx, IDs[idx])
	}
	return nil
}

// cleanup runs the Cleanup on each of the r.tests
func (r *perfRunner) cleanup() error {
	for _, t := range r.tests {
		err := t.Cleanup(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

// print an update for the last second
func (r *perfRunner) printStatus() error {
	if !r.warmupPrinted {
		finishedWarmup := r.printWarmupStatus()
		if !finishedWarmup {
			return nil
		}
	}

	if r.operationStatusTracker == -1 {
		err := r.printFinalUpdate(true)
		if err != nil {
			return err
		}
		r.operationStatusTracker = 0
		fmt.Fprintln(r.w, "\nCurrent\tTotal\tAverage\t")
	}
	totalOperations := r.totalOperations(false)

	_, err := fmt.Fprintf(
		r.w,
		"%s\t%s\t%s\t\n",
		r.messagePrinter.Sprintf("%d", totalOperations-r.operationStatusTracker),
		r.messagePrinter.Sprintf("%d", totalOperations),
		r.messagePrinter.Sprintf("%.2f", r.opsPerSecond(false)),
	)
	if err != nil {
		return err
	}
	r.operationStatusTracker = totalOperations
	return r.w.Flush()
}

// return true if all warmup information has been printed
func (r *perfRunner) printWarmupStatus() bool {
	if r.warmupOperationStatusTracker == -1 {
		r.warmupOperationStatusTracker = 0
		fmt.Println("===== WARMUP =====")
		fmt.Fprintln(r.w, "\nCurrent\tTotal\tAverage\t")
	}
	totalOperations := r.totalOperations(true)

	if r.warmupOperationStatusTracker == totalOperations {
		return true
	}

	_, err := fmt.Fprintf(
		r.w,
		"%s\t%s\t%s\t\n",
		r.messagePrinter.Sprintf("%d", totalOperations-r.warmupOperationStatusTracker),
		r.messagePrinter.Sprintf("%d", totalOperations),
		r.messagePrinter.Sprintf("%.2f", r.opsPerSecond(true)),
	)
	if err != nil {
		panic(err)
	}
	r.warmupOperationStatusTracker = totalOperations
	err = r.w.Flush()
	if err != nil {
		panic(err)
	}
	return false
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
		return fmt.Errorf("completed without generating operation statistics")
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
	return nil
}

// runTest takes care of the semantics of running a single iteration.
// It changes configuration on the proxy, increments counters, and
// updates the running-time.
func (r *perfRunner) runTest(p PerfTest, index int, ID string) {
	defer wg.Done()
	if debug {
		log.Printf("number of proxies %d", len(r.proxyTransports))
	}

	r.allOptions.mu.Lock()
	opts := r.allOptions.opts[index]
	r.allOptions.mu.Unlock()

	// If we are using the test proxy need to set up the in-memory recording.
	if testProxyURLs != "" {
		// First request goes through in Live mode
		r.proxyTransports[ID].SetMode("live")
		err := p.Run(context.Background())
		if err != nil {
			if err != nil {
				panic(err)
			}
		}

		// 2nd request goes through in Record mode
		r.proxyTransports[ID].SetMode("record")
		err = r.proxyTransports[ID].start()
		if err != nil {
			panic(err)

		}

		err = p.Run(context.Background())
		if err != nil {
			if err != nil {
				panic(err)
			}
		}
		err = r.proxyTransports[ID].stop()
		if err != nil {
			panic(err)
		}

		// All ensuing requests go through in Playback mode
		r.proxyTransports[ID].SetMode("playback")
		err = r.proxyTransports[ID].start()
		if err != nil {
			panic(err)
		}
	}

	// true parameter indicates were running the warmup here
	err := r.runTestForDuration(p, opts, true)
	if err != nil {
		panic(err)
	}

	// increment the warmupFinished counter that one goroutine has finished warmup
	val := atomic.AddInt32(&r.warmupFinished, 1)
	if debug {
		fmt.Printf("finished %d warmups\n", val)
	}

	// run the actual test
	err = r.runTestForDuration(p, opts, false)
	if err != nil {
		panic(err)
	}

	if testProxyURLs != "" {
		// Stop the proxy now
		err := proxyTransportsSuite[ID].stop()
		if err != nil {
			panic(err)
		}
		proxyTransportsSuite[ID].SetMode("live")
	}
	opts.finished = true
}

func (r *perfRunner) runTestForDuration(p PerfTest, opts *PerfTestOptions, warmup bool) error {
	if warmup && warmUpDuration <= 0 {
		return nil
	}

	// startPtr is our base time for keeping track of how long a test has run
	var startPtr *time.Time
	if warmup {
		t := time.Now()
		opts.warmupStart = &t
		startPtr = opts.warmupStart
	} else {
		t := time.Now()
		opts.runStart = &t
		startPtr = opts.runStart
	}

	var runDuration int
	if warmup {
		runDuration = warmUpDuration
	} else {
		runDuration = duration
	}

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(runDuration))
	defer cancel()

	lastSavedTime := time.Now()
	for time.Since(*startPtr).Seconds() < float64(runDuration) {
		err := p.Run(ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			} else {
				return err
			}
		}
		opts.increment(warmup)

		if time.Since(lastSavedTime).Seconds() > 0.3 {
			duration := time.Since(*startPtr)
			if warmup {
				atomic.StoreInt64((*int64)(&opts.warmupElapsed), int64(duration))
			} else {
				atomic.StoreInt64((*int64)(&opts.runElapsed), int64(duration))
			}
			lastSavedTime = time.Now()
		}
	}

	duration := time.Since(*startPtr)
	if warmup {
		atomic.StoreInt64((*int64)(&opts.warmupElapsed), int64(duration))
	} else {
		atomic.StoreInt64((*int64)(&opts.runElapsed), int64(duration))
	}

	return nil
}
