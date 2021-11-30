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
	"time"
)

var duration = flag.Int64("duration", 10, "Number of seconds to run as many operations as possible. Default is 10.")
var iterations = flag.Int64("iterations", 3, "Number of iterations to run. Default is 1.")
var parallel = flag.Int64("parallel", 1, "Number of tests to run in parallel. Default is 1.")
var warmUp = flag.Int64("warm-up", 5, "Number of seconds to spend warming up the connection before measuring begins. Default is 5.")
var noCleanup = flag.Bool("no-cleanup", false, "Whether to clean up newly created resources after test run. Default is false (resources will be deleted).")
var insecure = flag.Bool("insecure", false, "Whether to run without SSL validation. Default is false.")
var testProxies = flag.String("test-proxies", "", "Whether to run tests against the proxy sever. Specify the URLs for the proxy endpoints in a semi-colon separated list")

func init() {
	flag.Parse()

	fmt.Println("Running perf tests")
	fmt.Printf("Running for %d seconds\n", *duration)
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

// Helper function for pretty printing times in log messages
func printFormattedTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999999")
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

// PerfTest declares the interface for a performance test
type PerfTest interface {
	// Setup is a function run exactly once for any configurations of a performance test instance
	Setup()
	// Run is the function to be measured
	Run()
	// TearDown is a function run after completing the performance test, used for cleaning up resources.
	TearDown()
}

// Use RunPerfTest with a type that implements the PerfTest type to prepare a performance test
func RunPerfTest(p PerfTest) {
	fmt.Printf("Running setup...\n\n")
	p.Setup()

	if *warmUp > 0 {
		fmt.Printf("Warming up for %d seconds...\n", *warmUp)
		startTime := time.Now()
		warmUpTime := *warmUp * time.Second.Nanoseconds()
		for time.Since(startTime) < time.Duration(warmUpTime) {
			p.Run()
		}
	}

	for i := int64(0); i < *iterations; i++ {
		fmt.Printf("Starting iteration #%d\n", i+1)
		startTime := time.Now()
		count := 0
		for time.Since(startTime) < 10*time.Second {
			p.Run()
			count++
		}
		runTime := time.Since(startTime)
		opsPerSecond := float64(count) / runTime.Seconds()
		fmt.Printf("Round #%d: Completed %d operations in %.2f seconds\tAveraged %.4f operations / second\n", *iterations, count, runTime.Seconds(), opsPerSecond)
	}
	p.TearDown()

}
