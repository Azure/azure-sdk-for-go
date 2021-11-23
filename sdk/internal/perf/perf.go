//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package perf

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var runPerf = flag.Bool("perf", false, "Flag that must be specified to run perf tests.")
var duration = flag.Int64("duration", 10, "Number of seconds to run as many operations as possible. Default is 10.")
var iterations = flag.Int64("iterations", 1, "Number of iterations to run. Default is 1.")
var parallel = flag.Int64("parallel", 1, "Number of tests to run in parallel. Default is 1.")
var warmUp = flag.Int64("warm-up", 5, "Number of seconds to spend warming up the connection before measuring begins. Default is 5.")
var noCleanup = flag.Bool("no-cleanup", false, "Whether to clean up newly created resources after test run. Default is false (resources will be deleted).")
var insecure = flag.Bool("insecure", false, "Whether to run without SSL validation. Default is false.")
var testProxies = flag.String("test-proxies", "", "Whether to run tests against the proxy sever. Specify the URLs for the proxy endpoints in a semi-colon separated list")

func init() {
	if *runPerf {
		fmt.Println("Running perf tests")
		fmt.Printf("Running %d times\n", *duration)
		fmt.Printf("Running %d iterations\n", *iterations)
		fmt.Printf("Running %d parallel threads\n", *parallel)
		fmt.Printf("Warming up for %d seconds\n", *warmUp)
		if *noCleanup {
			fmt.Println("Skipping cleanup")
		}
		if *insecure {
			fmt.Println("Running without SSL validation")
		}
		if *testProxies != "" {
			fmt.Printf("Running against proxies hosted at %s\n", *testProxies)
		}
	}
}

// runPerfTests checks if the flag was set and the test is a perf test.
// The setup, and run methods will skip if these conditions are not true.
func runPerfTests(t *testing.T) bool {
	return strings.Contains(t.Name(), "Perf") && *runPerf
}

// Helper function for pretty printing times in log messages
func printFormattedTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999999")
}

var globalSetupRan = false

// The global set up is run only once. Use this for any setup that can be reused multiple times by all test instances
func GlobalSetup(t *testing.T, setupFunc func() error) {
	if !runPerfTests(t) {
		t.Skip()
	}
	if globalSetupRan {
		// Only want to run this once so return early.
		return
	}

	fmt.Printf("%v: %s\n", printFormattedTime(time.Now()), "Running global set up function")
	err := setupFunc()
	if err != nil {
		t.Logf("Global set up was not successful: %s", err.Error())
		panic(err)
	}
	globalSetupRan = true
}

// The global tear down is run after all perf tests have completed.
func GlobalTeardown(t *testing.T, teardownFunc func() error) {
	fmt.Printf("%v: %s\n", printFormattedTime(time.Now()), "Running global tear down function")
	err := teardownFunc()
	if err != nil {
		t.Logf("Global teardown was not successful: %s", err.Error())
		panic(err)
	}
}

func RunFunc(t *testing.T, perfFunc func()) {
	// Run WarmUp
	fmt.Printf("%v: %s\n", printFormattedTime(time.Now()), "Measuring performance")
	end := time.Now().Add(time.Second * time.Duration(*warmUp))
	for time.Since(end) < 0 {
		perfFunc()
	}

	totalCounter := 0
	var i int64
	totalDuration := 0.0
	for i = 0; i < *iterations; i++ {
		start := time.Now()
		counter := 0
		for time.Since(start) < time.Second*time.Duration(*duration) {
			perfFunc()
			counter += 1
		}
		now := time.Now()
		iterDuration := now.Sub(start).Seconds()

		t.Logf(
			"%s: Iteration #%d:\t%s: Finished running %d operations in %d seconds. %.2f ops/sec\n",
			printFormattedTime(time.Now()),
			i,
			t.Name(),
			counter,
			*duration,
			float64(counter)/iterDuration,
		)
		totalCounter += counter
		totalDuration += iterDuration
	}
	t.Logf(
		"%s: Finished running a total of %d operations in %d seconds. Average of %.2f ops/sec\n",
		printFormattedTime(time.Now()),
		totalCounter,
		*duration**iterations,
		float64(totalCounter)/totalDuration,
	)
}

type randomStream struct {
	size     int64
	position int64
	r        rand.Rand
}

// randomStream implements the io.ReadSeekCloser methods to use in place of reading files/text data.
// the stream will only store the current portion of data in memory
func NewRandomStream(size int64) randomStream {
	return randomStream{
		size:     size,
		position: 0,
		r:        *rand.New(rand.NewSource(0)),
	}
}

func (r *randomStream) Read(p []byte) (n int, err error) {
	if (r.size - r.position) < int64(len(p)) {
		n, err = r.r.Read(p[:(r.size - r.position)])
		if err != nil {
			return n, err
		}
		return int(r.size - r.position), io.EOF
	}

	r.position += int64(len(p))
	n, err = r.r.Read(p)
	if err != nil {
		return n, err
	}

	return len(p), nil
}

func (r *randomStream) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		r.position = offset
	} else if whence == io.SeekCurrent {
		r.position += offset
	} else if whence == io.SeekEnd {
		r.position = r.size + offset
	}
	return int64(r.position), nil
}

func (r *randomStream) Close() error {
	return nil
}
