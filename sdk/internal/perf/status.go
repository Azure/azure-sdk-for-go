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
}

func newStatusRunner(t time.Time, runTime int) *statusRunner {
	return &statusRunner{
		results:           make([]runResult, 0),
		start:             t,
		perRoutineResults: map[int][]runResult{},
		// Printing happens half a second after each run second passes to
		// allow all messages to be handled.
		lastPrint:    t,
		totalRunTime: runTime,
	}
}

func (s *statusRunner) handleMessage(msg runResult, w *tabwriter.Writer) {
	if len(s.results) == 0 {
		// This is the first message handled, print the msg
		fmt.Println("\n=== Warm Up ===")
		fmt.Fprintln(w, "Current\tTotal\tAverage\t")
	}
	s.results = append(s.results, msg)

	if msg.completed {
		s.routinesFinished += 1
	}

	s.total += msg.count

	s.perRoutineResults[msg.parallelIndex] = append(s.perRoutineResults[msg.parallelIndex], msg)

	if time.Since(s.lastPrint).Seconds() > 0.9 {
		s.printUpdate(w)
	}
}

func (s *statusRunner) finished() bool {
	return s.routinesFinished == len(s.perRoutineResults)
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

	os.Exit(1)
}

func (s *statusRunner) printUpdate(w *tabwriter.Writer) {
	if s.finished() {
		s.printFinalUpdate()
	}
	// printUntil := int(time.Since(s.start).Seconds())

	// current := 0
	// for _, msgs := range s.perRoutineResults {
	// 	if debug {
	// 		fmt.Println(msgs)
	// 	}
	// 	// Ensure we have the correct number of messages
	// 	if len(msgs) < printUntil-1 {
	// 		current += msgs[printUntil].count
	// 	}
	// }
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

var perSecondCount [][]int
var elapsedTimes [][]float64

var perSecondCountWarmup [][]int
var elapsedTimesWarmup [][]float64

var printedWarmupResults bool = false

var messagePrinter *message.Printer = message.NewPrinter(message.MatchLanguage("en"))

// helper function for handling status updates
func handleMessage(w *tabwriter.Writer, msg runResult) {
	if msg.warmup {
		handleWarmupMessage(w, msg)
		return
	}

	// Check if we need to print out results from warmup. Results come in a channel, so we
	// need to check if all N channels (N = Parallel) have reported final results
	if warmUpDuration > 0 {
		if len(elapsedTimesWarmup[warmUpDuration-1]) == parallelInstances && !printedWarmupResults {
			printFinalResults(elapsedTimesWarmup, perSecondCountWarmup, true)
			printedWarmupResults = true
		}
	}

	if len(perSecondCount) == 0 {
		// Initialize the slice of slices
		for i := 0; i < duration; i++ {
			perSecondCount = append(perSecondCount, []int{})
			elapsedTimes = append(elapsedTimes, []float64{})
		}
	}

	updateSecond := int(msg.timeInSeconds) - 1
	perSecondCount[updateSecond] = append(perSecondCount[updateSecond], msg.count)
	elapsedTimes[updateSecond] = append(elapsedTimes[updateSecond], msg.timeInSeconds)

	if len(perSecondCount[updateSecond]) == parallelInstances {
		if updateSecond == 0 {
			fmt.Println("\n=== Test ===")
			fmt.Fprintln(w, "Current\tTotal\tAverage\t")
		}

		thisCount := sumInts(perSecondCount[updateSecond])
		totalCount := 0
		for _, c := range perSecondCount {
			totalCount += sumInts(c)
		}

		avg := computeAverageOpsPerSecond(perSecondCount, elapsedTimes)

		_, err := fmt.Fprintf(
			w,
			"%s\t%s\t%s\t\n",
			messagePrinter.Sprintf("%d", thisCount),
			messagePrinter.Sprintf("%d", totalCount),
			messagePrinter.Sprintf("%.2f", avg),
		)
		if err != nil {
			panic(err)
		}
		w.Flush()
	}
}

func handleWarmupMessage(w *tabwriter.Writer, msg runResult) {
	if len(perSecondCountWarmup) == 0 {
		// Initialize the slice of slices for warmups
		for i := 0; i < warmUpDuration; i++ {
			perSecondCountWarmup = append(perSecondCountWarmup, []int{})
			elapsedTimesWarmup = append(elapsedTimesWarmup, []float64{})
		}
	}

	updateSecond := int(msg.timeInSeconds) - 1
	perSecondCountWarmup[updateSecond] = append(perSecondCountWarmup[updateSecond], msg.count)
	elapsedTimesWarmup[updateSecond] = append(elapsedTimesWarmup[updateSecond], msg.timeInSeconds)

	if len(perSecondCountWarmup[updateSecond]) == parallelInstances {
		if updateSecond == 0 {
			fmt.Println("\n=== Warmup ===")
			fmt.Fprintln(w, "Current\tTotal\tAverage\t")
		}

		thisCount := sumInts(perSecondCountWarmup[updateSecond])
		totalCount := 0
		for _, c := range perSecondCountWarmup {
			totalCount += sumInts(c)
		}

		avg := computeAverageOpsPerSecond(perSecondCountWarmup, elapsedTimesWarmup)

		p := message.NewPrinter(message.MatchLanguage("en"))

		_, err := fmt.Fprintf(
			w,
			"%s\t%s\t%s\t\n",
			p.Sprintf("%d", thisCount),
			p.Sprintf("%d", totalCount),
			p.Sprintf("%.2f", avg),
		)
		if err != nil {
			panic(err)
		}
		w.Flush()
	}
}

func computeAverageOpsPerSecond(perSecondCount [][]int, elapsedTimes [][]float64) float64 {
	var avg float64

	for p := 0; p < parallelInstances; p++ {
		threadOps := 0
		timeElapsed := 0.0
		for i := 0; i < len(perSecondCount); i++ {
			if len(perSecondCount[i]) == 0 || len(elapsedTimes[i]) == 0 {
				break
			}
			threadOps += perSecondCount[i][p]
			timeElapsed = elapsedTimes[i][p]
		}

		avg += float64(threadOps) / timeElapsed
	}

	return avg
}

func printFinalResults(elapsedTimes [][]float64, perSecondCount [][]int, warmup bool) {
	opsPerRoutine := make([]int, parallelInstances)
	secondsPerRoutine := make([]float64, parallelInstances)
	innerLoop := duration
	if warmup {
		innerLoop = warmUpDuration
	}
	for i := 0; i < parallelInstances; i++ {
		secondsPerRoutine[i] = elapsedTimes[innerLoop-1][i]
		for j := 0; j < innerLoop; j++ {
			opsPerRoutine[i] += perSecondCount[j][i]
		}
	}

	opsPerSecond := 0.0
	for i := 0; i < parallelInstances; i++ {
		opsPerSecond += float64(opsPerRoutine[i]) / secondsPerRoutine[i]
	}

	totalOperations := sumInts(opsPerRoutine)

	fmt.Println("\n=== Results ===")
	secondsPerOp := 1.0 / opsPerSecond
	weightedAvgSec := float64(totalOperations) / opsPerSecond
	fmt.Printf(
		"Completed %s operations in a weighted-average of %ss (%s ops/s, %s s/op)\n",
		messagePrinter.Sprintf("%d", totalOperations),
		messagePrinter.Sprintf("%.2f", weightedAvgSec),
		messagePrinter.Sprintf("%.2f", opsPerSecond),
		messagePrinter.Sprintf("%.3f", secondsPerOp),
	)
}
