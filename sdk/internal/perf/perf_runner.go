package perf

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"golang.org/x/text/message"
)

type perfRunner struct {
	// ticker is the runner for giving updates every second
	ticker *time.Ticker
	done   chan bool

	// name of the performance test
	name string

	// the perf test, options, and transports being tested/used
	perfToRun       PerfMethods
	allOptions      []*PerfTestOptions
	proxyTransports map[string]*RecordingHTTPClient

	// All created tests
	tests []PerfTest

	globalInstance GlobalPerfTest

	// this is the previous prints total
	warmupOperationStatusTracker int64
	operationStatusTracker       int64

	// writer
	w              *tabwriter.Writer
	messagePrinter *message.Printer

	// tracker for whether the warmup has finished
	warmupFinished int32
	warmupPrinted  bool
}

func newPerfRunner(p PerfMethods, name string) *perfRunner {
	return &perfRunner{
		ticker:                       time.NewTicker(time.Second),
		done:                         make(chan bool),
		name:                         name,
		perfToRun:                    p,
		operationStatusTracker:       -1,
		warmupOperationStatusTracker: -1,
		w:                            tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight),
		messagePrinter:               message.NewPrinter(message.MatchLanguage("en")),
		warmupFinished:               0,
		warmupPrinted:                false,
	}
}

func (r *perfRunner) Run() error {
	// Poller for printing
	go func() {
		for {
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
	}()

	err := r.globalSetup()
	if err != nil {
		return err
	}

	err = r.createPerfTests()
	if err != nil {
		return err
	}
	r.done <- true

	r.printFinalUpdate(false)
	err = r.cleanup()
	if err != nil {
		panic(err)
	}
	err = r.globalInstance.GlobalCleanup(context.Background())
	if err != nil {
		panic(err)
	}
	return nil
}

// global setup by instantiating a single global instance
func (r *perfRunner) globalSetup() error {
	globalInst, err := r.perfToRun.New(context.TODO(), PerfTestOptions{Name: r.name})
	if err != nil {
		return err
	}
	r.globalInstance = globalInst
	return nil
}

func (r *perfRunner) createPerfTests() error {
	var IDs []string
	proxyURLS := parseProxyURLS()

	for idx := 0; idx < parallelInstances; idx++ {
		options := &PerfTestOptions{}
		ID := fmt.Sprintf("%s-%d", r.name, idx)
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

		perfTest, err := r.globalInstance.NewPerfTest(context.TODO(), options)
		if err != nil {
			return err
		}
		r.tests = append(r.tests, perfTest)
		r.allOptions = append(r.allOptions, options)
	}

	for idx, test := range r.tests {
		wg.Add(1)
		go r.runTest(test, idx, IDs[idx])
	}

	wg.Wait()

	return nil
}

func (r *perfRunner) cleanup() error {
	for _, t := range r.tests {
		err := t.Cleanup(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *perfRunner) printStatus() error {
	if !r.warmupPrinted {
		finishedWarmup := r.printWarmupStatus()
		if !finishedWarmup {
			return nil
		}
	}

	if r.operationStatusTracker == -1 {
		r.printFinalUpdate(true)
		r.operationStatusTracker = 0
		fmt.Fprintln(r.w, "Current\tTotal\tAverage\t")
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
		fmt.Fprintln(r.w, "Current\tTotal\tAverage\t")
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

func (r *perfRunner) totalOperations(warmup bool) int64 {
	var ret int64

	for _, opt := range r.allOptions {
		if warmup {
			ret += atomic.LoadInt64(&opt.warmupCount)
		} else {
			ret += atomic.LoadInt64(&opt.runCount)
		}
	}

	return ret
}

func (r *perfRunner) opsPerSecond(warmup bool) float64 {
	var ret float64
	for _, opt := range r.allOptions {
		if warmup {
			ret += float64(atomic.LoadInt64(&opt.warmupCount)) / opt.warmupElapsed.GetFloat()
		} else {
			ret += float64(atomic.LoadInt64(&opt.runCount)) / opt.runElapsed.GetFloat()
		}
	}
	return ret
}

func (r *perfRunner) printFinalUpdate(warmup bool) error {
	totalOperations := r.totalOperations(warmup)
	opsPerSecond := r.opsPerSecond(warmup)
	if opsPerSecond == 0.0 {
		return fmt.Errorf("completed without generating operation statistics")
	}

	secondsPerOp := 1.0 / opsPerSecond
	weightedAvg := float64(totalOperations) / opsPerSecond

	fmt.Println("\n=== Results ===")
	fmt.Printf(
		"Completed %s operations in a weighted-average of %ss (%s ops/s, %s s/op)\n",
		r.messagePrinter.Sprintf("%d", totalOperations),
		r.messagePrinter.Sprintf("%.2f", weightedAvg),
		r.messagePrinter.Sprintf("%.3f", opsPerSecond),
		r.messagePrinter.Sprintf("%.3f", secondsPerOp),
	)
	return nil
}

// runTest takes care of the semantics of running a single iteration. It returns the number of times the test ran as an int, the exact number
// of seconds the test ran as a float64, and any errors.
func (r *perfRunner) runTest(p PerfTest, index int, ID string) {
	defer wg.Done()
	if debug {
		log.Printf("number of proxies %d", len(r.proxyTransports))
	}

	opts := r.allOptions[index]

	// If we are using the test proxy need to set up the in-memory recording.
	if testProxyURLs != "" {
		// First request goes through in Live mode
		r.proxyTransports[ID].SetMode("live")
		err := p.Run(context.Background())
		if err != nil {
			panic(err)
		}

		// 2nd request goes through in Record mode
		r.proxyTransports[ID].SetMode("record")
		err = r.proxyTransports[ID].start()
		if err != nil {
			panic(err)
		}
		err = p.Run(context.Background())
		if err != nil {
			panic(err)
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

	if warmUpDuration > 0 {
		opts.warmupStart = time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(warmUpDuration))
		defer cancel()

		lastSavedTime := time.Now()
		for time.Since(opts.warmupStart).Seconds() < float64(warmUpDuration) {
			err := p.Run(ctx)
			if err != nil {
				panic(err)
			}
			opts.incrememt(true)

			if time.Since(lastSavedTime).Seconds() > 0.3 {
				opts.warmupElapsed.SetFloat(time.Since(opts.warmupStart).Seconds())
				lastSavedTime = time.Now()
			}
		}

		opts.warmupElapsed.SetFloat(time.Since(opts.warmupStart).Seconds())
	}
	_ = atomic.AddInt32(&r.warmupFinished, 1)

	opts.runStart = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(duration))
	defer cancel()

	lastSavedTime := time.Now()
	for time.Since(opts.runStart).Seconds() < float64(duration) {
		err := p.Run(ctx)
		if err != nil {
			panic(err)
		}
		opts.incrememt(false)

		if time.Since(lastSavedTime).Seconds() > 0.1 {
			opts.runElapsed.SetFloat(time.Since(opts.runStart).Seconds())
			lastSavedTime = time.Now()
		}
	}

	opts.runElapsed.SetFloat(time.Since(opts.runStart).Seconds())

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
