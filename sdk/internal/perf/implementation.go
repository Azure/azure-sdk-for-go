// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

var (
	debug             bool
	duration          int
	testProxyURLs     string
	warmUpDuration    int
	parallelInstances int
	wg                sync.WaitGroup
	numProcesses      int
)

// runTest takes care of the semantics of running a single iteration. It returns the number of times the test ran as an int, the exact number
// of seconds the test ran as a float64, and any errors.
func runTest(p PerfTest, index int, c chan runResult, ID string) {
	defer wg.Done()
	if debug {
		log.Printf("number of proxies %d", len(proxyTransportsSuite))
	}

	// If we are using the test proxy need to set up the in-memory recording.
	if testProxyURLs != "" {
		// First request goes through in Live mode
		proxyTransportsSuite[ID].SetMode("live")
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}

		// 2nd request goes through in Record mode
		proxyTransportsSuite[ID].SetMode("record")
		err = proxyTransportsSuite[ID].start()
		if err != nil {
			panic(err)
		}
		err = p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}
		err = proxyTransportsSuite[ID].stop()
		if err != nil {
			panic(err)
		}

		// All ensuing requests go through in Playback mode
		proxyTransportsSuite[ID].SetMode("playback")
		err = proxyTransportsSuite[ID].start()
		if err != nil {
			panic(err)
		}
	}

	if warmUpDuration > 0 {
		warmUpStart := time.Now()
		warmUpLastPrint := 1.0
		warmUpPerSecondCount := make([]int, 0)
		warmupCount := 0
		for time.Since(warmUpStart).Seconds() < float64(warmUpDuration) {
			err := p.Run(context.Background())
			if err != nil {
				c <- runResult{err: err}
			}
			warmupCount += 1

			if time.Since(warmUpStart).Seconds() > warmUpLastPrint {
				thisSecondCount := warmupCount - sumInts(warmUpPerSecondCount)
				c <- runResult{
					warmup:        true,
					count:         thisSecondCount,
					parallelIndex: index,
					timeInSeconds: time.Since(warmUpStart).Seconds(),
				}

				warmUpPerSecondCount = append(warmUpPerSecondCount, thisSecondCount)
				warmUpLastPrint += 1.0
				if int(warmUpLastPrint) == warmUpDuration {
					// We can have odd scenarios where we send this
					// message, and the final message below right after
					warmUpLastPrint += 1.0
				}
			}
		}

		thisSecondCount := warmupCount - sumInts(warmUpPerSecondCount)
		c <- runResult{
			warmup:        true,
			completed:     true,
			count:         thisSecondCount,
			parallelIndex: index,
			timeInSeconds: time.Since(warmUpStart).Seconds(),
		}
	}

	timeStart := time.Now()
	totalCount := 0
	lastPrint := 1.0
	perSecondCount := make([]int, 0)
	for time.Since(timeStart).Seconds() < float64(duration) {
		err := p.Run(context.Background())
		if err != nil {
			c <- runResult{err: err}
		}
		totalCount += 1

		// Every second (roughly) we send an update through the channel
		if time.Since(timeStart).Seconds() > lastPrint {
			thisCount := totalCount - sumInts(perSecondCount)
			c <- runResult{
				count:         thisCount,
				parallelIndex: index,
				completed:     false,
				timeInSeconds: time.Since(timeStart).Seconds(),
			}
			lastPrint += 1.0
			perSecondCount = append(perSecondCount, thisCount)

			if int(lastPrint) == duration {
				// if we are w/in one second of the end time, we do not
				// want to send any more results, we'll just send a final result
				lastPrint += 10.0
			}
		}
	}

	elapsed := time.Since(timeStart).Seconds()
	lastSecondCount := totalCount - sumInts(perSecondCount)
	c <- runResult{
		completed:     true,
		count:         lastSecondCount,
		parallelIndex: index,
		timeInSeconds: elapsed,
	}

	if testProxyURLs != "" {
		// Stop the proxy now
		err := proxyTransportsSuite[ID].stop()
		if err != nil {
			c <- runResult{err: err}
		}
		proxyTransportsSuite[ID].SetMode("live")
	}
}

// runCleanup takes care of the semantics for tearing down a single iteration of a performance test.
func runCleanup(p PerfTest) error {
	return p.Cleanup(context.Background())
}

// runResult is the result sent back through the channel for updates and final results
type runResult struct {
	// number of iterations completed since the previous message
	count int

	// The time the update comes from
	timeInSeconds float64

	// if there is an error it will be here
	err error

	// true when this is the last result from a go routine
	completed bool

	// Index of the goroutine
	parallelIndex int

	// indicates if result is from a warmup start
	warmup bool
}

// parse the TestProxy input into a slice of strings
func parseProxyURLS() []string {
	var ret []string
	if testProxyURLs == "" {
		return ret
	}

	testProxyURLs = strings.TrimSuffix(testProxyURLs, ";")

	ret = strings.Split(testProxyURLs, ";")

	return ret
}

// Spins off each Parallel instance as a separate goroutine, reads messages and runs cleanup/setup methods.
func runPerfTest(name string, p NewPerfTest) error {
	globalInstance, err := p(context.TODO(), PerfTestOptions{Name: name})
	if err != nil {
		return err
	}

	var perfTests []PerfTest
	var IDs []string
	proxyURLS := parseProxyURLS()

	w := tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight)

	messages := make(chan runResult)
	for idx := 0; idx < parallelInstances; idx++ {
		options := &PerfTestOptions{}

		ID := fmt.Sprintf("%s-%d", name, idx)
		IDs = append(IDs, ID)

		if testProxyURLs != "" {
			proxyURL := proxyURLS[idx%len(proxyURLS)]
			transporter := NewProxyTransport(&TransportOptions{
				TestName: ID,
				proxyURL: proxyURL,
			})
			options.Transporter = transporter
		} else {
			options.Transporter = defaultHTTPClient
		}
		options.parallelIndex = idx

		perfTest, err := globalInstance.NewPerfTest(context.TODO(), options)
		if err != nil {
			panic(err)
		}
		perfTests = append(perfTests, perfTest)
	}

	for idx, perfTest := range perfTests {
		// Create a thread for running a single test
		wg.Add(1)
		go runTest(perfTest, idx, messages, IDs[idx])
	}

	// Add another goroutine to close the channel after completion
	go func() {
		wg.Wait()
		close(messages)
	}()

	// Read incoming messages and handle status updates
	for msg := range messages {
		if debug {
			log.Println("Handling message: ", msg)
		}
		if msg.err != nil {
			panic(msg.err)
		}
		handleMessage(w, msg)
	}

	// Print before running the cleanup in case cleanup takes a while
	printFinalResults(elapsedTimes, perSecondCount, false)

	// Run Cleanup on each parallel instance
	for _, pTest := range perfTests {
		err := runCleanup(pTest)
		if err != nil {
			panic(err)
		}
	}

	err = globalInstance.GlobalCleanup(context.TODO())
	if err != nil {
		return fmt.Errorf("there was an error with the GlobalCleanup method: %v", err.Error())
	}

	return nil
}
