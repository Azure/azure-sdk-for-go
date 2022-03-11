// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"golang.org/x/text/message"
)

var messagePrinter *message.Printer = message.NewPrinter(message.MatchLanguage("en"))

// statusRunner is the struct responsible for handling messages
type statusRunner struct {
	// results is a slice of all results from the goroutines
	results []runResult

	// start is the time the statusRunner was started
	start time.Time

	// perRoutineResults map the parallel index to a slice of runResults
	perRoutineResults map[int][]runResult

	// lastPrint holds when the last information was printed to stdout
	// the initial value is the same as start. When this value exceeds
	// is more than 1 second after time.Now(), a new update is printed.
	lastPrint time.Time

	// total is a running count of the count of performance tests run
	total int

	// prevTotal is the total of the last output
	prevTotal int

	// hasFinished indicates if the final results have been printed out
	totalRunTime int

	// routinesFinished indicates how many routines have sent a message
	// indicating they have completed execution
	routinesFinished int

	// isWarmup indicates whether the messages are from warmup
	isWarmup bool
}

func newStatusRunner(t time.Time, runTime int) *statusRunner {
	return &statusRunner{
		results:           make([]runResult, 0),
		start:             t,
		perRoutineResults: map[int][]runResult{},
		lastPrint:         t,
		totalRunTime:      runTime,
	}
}

func (s *statusRunner) handleMessage(msg runResult, w *tabwriter.Writer) {
	s.results = append(s.results, msg)

	if msg.completed {
		s.routinesFinished += 1
	}

	s.total += msg.count

	s.perRoutineResults[msg.parallelIndex] = append(s.perRoutineResults[msg.parallelIndex], msg)
}

func (s *statusRunner) printUpdates() {
	w := tabwriter.NewWriter(os.Stdout, 16, 8, 1, ' ', tabwriter.AlignRight)
	firstPrint := false
	for s.routinesFinished != parallelInstances {
		// Poll and print
		if time.Since(s.lastPrint).Seconds() > 1.0 {

			if !firstPrint {
				if s.isWarmup {
					fmt.Println("\n=== Warm Up ===")
				} else {
					fmt.Println("\n=== Test ===")
				}
				fmt.Fprintln(w, "Current\tTotal\tAverage\t")
				w.Flush()
				firstPrint = true
			}

			avg := float64(s.total) / time.Since(s.start).Seconds()
			_, err := fmt.Fprintf(
				w,
				"%s\t%s\t%s\t\n",
				messagePrinter.Sprintf("%d", s.total-s.prevTotal),
				messagePrinter.Sprintf("%d", s.total),
				messagePrinter.Sprintf("%.2f", avg),
			)
			if err != nil {
				panic(err)
			}

			w.Flush()

			s.lastPrint = time.Now()
			s.prevTotal = s.total
		}
	}
}

func (s *statusRunner) printFinalUpdate() {
	opsPerRoutine := make([]int, parallelInstances)
	secondsPerRoutine := make([]float64, parallelInstances)

	for pIdx, msgs := range s.perRoutineResults {
		secondsPerRoutine[pIdx] = msgs[len(msgs)-1].timeInSeconds
		for _, msg := range msgs {
			opsPerRoutine[pIdx] += msg.count
		}
	}

	opsPerSecond := 0.0
	for i := 0; i < parallelInstances; i++ {
		opsPerSecond += float64(opsPerRoutine[i]) / secondsPerRoutine[i]
	}

	fmt.Println("\n=== Results ===")
	secondsPerOp := 1.0 / opsPerSecond
	weightedAvgSec := float64(s.total) / opsPerSecond
	fmt.Printf(
		"Completed %s operations in a weighted-average of %ss (%s ops/s, %s s/op)\n",
		messagePrinter.Sprintf("%d", s.total),
		messagePrinter.Sprintf("%.2f", weightedAvgSec),
		messagePrinter.Sprintf("%.2f", opsPerSecond),
		messagePrinter.Sprintf("%.3f", secondsPerOp),
	)
}
