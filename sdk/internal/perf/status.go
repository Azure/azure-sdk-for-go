// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"text/tabwriter"

	"golang.org/x/text/message"
)

var perSecondCount [][]int
var elapsedTimes [][]float64

var opsPerRoutine []int
var secondsPerRoutine []float64

var printedWarmupResults bool = false

// helper function for handling status updates
func handleMessage(w *tabwriter.Writer, msg runResult) {
	if msg.warmup {
		handleWarmupMessage(w, msg)
		return
	}

	// Check if we need to print out results from warmup
	if len(elapsedTimesWarmup[WarmUp-1]) == Parallel && !printedWarmupResults {
		printFinalResultsWarmup()
		printedWarmupResults = true
	}

	if len(perSecondCount) == 0 {
		// Initialize the slice of slices
		for i := 0; i < Duration; i++ {
			perSecondCount = append(perSecondCount, []int{})
			elapsedTimes = append(elapsedTimes, []float64{})
		}
		opsPerRoutine = make([]int, Parallel)
		secondsPerRoutine = make([]float64, Parallel)
	}

	if msg.completed {
		opsPerRoutine[msg.parallelIndex] = msg.count
		secondsPerRoutine[msg.parallelIndex] = msg.timeInSeconds
	}

	updateSecond := int(msg.timeInSeconds) - 1
	perSecondCount[updateSecond] = append(perSecondCount[updateSecond], msg.count)
	elapsedTimes[updateSecond] = append(elapsedTimes[updateSecond], msg.timeInSeconds)

	if len(perSecondCount[updateSecond]) == Parallel {
		if updateSecond == 0 {
			fmt.Println("\n=== Test ===")
			fmt.Fprintln(w, "Current\tTotal\tAverage\t")
		}

		thisCount := sumInts(perSecondCount[updateSecond])
		totalCount := 0
		for _, c := range perSecondCount {
			totalCount += sumInts(c)
		}

		avg := computeAverageOpsPerSecond()

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

var perSecondCountWarmup [][]int
var elapsedTimesWarmup [][]float64

var opsPerRoutineWarmup []int
var secondsPerRoutineWarmup []float64

func handleWarmupMessage(w *tabwriter.Writer, msg runResult) {
	if len(perSecondCountWarmup) == 0 {
		// Initialize the slice of slices for warmups
		for i := 0; i < WarmUp; i++ {
			perSecondCountWarmup = append(perSecondCountWarmup, []int{})
			elapsedTimesWarmup = append(elapsedTimesWarmup, []float64{})
		}
		opsPerRoutineWarmup = make([]int, Parallel)
		secondsPerRoutineWarmup = make([]float64, Parallel)
	}

	updateSecond := int(msg.timeInSeconds) - 1
	perSecondCountWarmup[updateSecond] = append(perSecondCountWarmup[updateSecond], msg.count)
	elapsedTimesWarmup[updateSecond] = append(elapsedTimesWarmup[updateSecond], msg.timeInSeconds)

	if len(perSecondCountWarmup[updateSecond]) == Parallel {
		if updateSecond == 0 {
			fmt.Println("\n=== Warmup ===")
			fmt.Fprintln(w, "Current\tTotal\tAverage\t")
		}

		thisCount := sumInts(perSecondCountWarmup[updateSecond])
		totalCount := 0
		for _, c := range perSecondCountWarmup {
			totalCount += sumInts(c)
		}

		avg := computeAverageOpsPerSecondWarmup()

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

func computeAverageOpsPerSecondWarmup() float64 {
	var avg float64

	for p := 0; p < Parallel; p++ {
		threadOps := 0
		timeElapsed := 0.0
		for i := 0; i < len(perSecondCountWarmup); i++ {
			if len(perSecondCountWarmup[i]) == 0 || len(elapsedTimesWarmup[i]) == 0 {
				break
			}
			threadOps += perSecondCountWarmup[i][p]
			timeElapsed = elapsedTimesWarmup[i][p]
		}

		avg += float64(threadOps) / timeElapsed
	}

	return avg
}

func computeAverageOpsPerSecond() float64 {
	var avg float64

	for p := 0; p < Parallel; p++ {
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

func printFinalResultsWarmup() {
	opsPerRoutineWarmup := make([]int, Parallel)
	secondsPerRoutineWarmup := make([]float64, Parallel)
	for i := 0; i < Parallel; i++ {
		secondsPerRoutineWarmup[i] = elapsedTimesWarmup[WarmUp-1][i]
		for j := 0; j < WarmUp; j++ {
			opsPerRoutineWarmup[i] += perSecondCountWarmup[j][i]
		}
	}

	opsPerSecond := 0.0
	for i := 0; i < Parallel; i++ {
		opsPerSecond += float64(opsPerRoutineWarmup[i]) / secondsPerRoutineWarmup[i]
	}

	totalOperations := sumInts(opsPerRoutineWarmup)

	p := message.NewPrinter(message.MatchLanguage("en"))

	fmt.Println("\n=== Results ===")
	secondsPerOp := 1.0 / opsPerSecond
	weightedAvgSec := float64(totalOperations) / opsPerSecond
	fmt.Printf(
		"Completed %s operations in a weighted-average of %.2fs (%s ops/s, %.3f s/op)\n",
		p.Sprintf("%d", totalOperations),
		weightedAvgSec,
		p.Sprintf("%d", int(opsPerSecond)),
		secondsPerOp,
	)
}

func printFinalResults() {
	opsPerSecond := 0.0
	for i := 0; i < Parallel; i++ {
		opsPerSecond += float64(opsPerRoutine[i]) / secondsPerRoutine[i]
	}

	totalOperations := sumInts(opsPerRoutine)
	p := message.NewPrinter(message.MatchLanguage("en"))

	fmt.Println("\n=== Results ===")
	secondsPerOp := 1.0 / opsPerSecond
	weightedAvgSec := float64(totalOperations) / opsPerSecond
	fmt.Printf(
		"Completed %s operations in a weighted-average of %.2fs (%s ops/s, %.3f s/op)\n",
		p.Sprintf("%d", totalOperations),
		weightedAvgSec,
		p.Sprintf("%d", int(opsPerSecond)),
		secondsPerOp,
	)
}
