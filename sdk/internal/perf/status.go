// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"text/tabwriter"
)

var perSecondCount [][]int
var elapsedTimes [][]float64

var opsPerRoutine []int
var secondsPerRoutine []float64

// helper function for handling status updates
func handleMessage(w *tabwriter.Writer, msg runResult) {
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
		return
	}

	updateSecond := int(msg.timeInSeconds) - 1
	perSecondCount[updateSecond] = append(perSecondCount[updateSecond], msg.count)
	elapsedTimes[updateSecond] = append(elapsedTimes[updateSecond], msg.timeInSeconds)

	if len(perSecondCount[updateSecond]) == Parallel {
		thisCount := sumInts(perSecondCount[updateSecond])
		totalCount := 0
		for _, c := range perSecondCount {
			totalCount += sumInts(c)
		}

		avg := computeAverageOpsPerSecond()

		_, err := fmt.Fprintf(
			w,
			"%d\t%s\t%s\t%.2f\t\n",
			updateSecond+1,
			commaIze(thisCount),
			commaIze(totalCount),
			avg,
		)
		if err != nil {
			panic(err)
		}
		w.Flush()
	}
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

func printFinalResults() {
	opsPerSecond := 0.0
	for i := 0; i < Parallel; i++ {
		opsPerSecond += float64(opsPerRoutine[i]) / secondsPerRoutine[i]
	}

	totalOperations := sumInts(opsPerRoutine)

	fmt.Println("\n=== Results ===")
	secondsPerOp := 1.0 / opsPerSecond
	weightedAvgSec := float64(totalOperations) / opsPerSecond
	fmt.Printf(
		"Completed %s operations in a weighted-average of %.2fs (%s ops/s, %.3f s/op)\n",
		commaIze(totalOperations),
		weightedAvgSec,
		commaIze(int(opsPerSecond)),
		secondsPerOp,
	)
}
