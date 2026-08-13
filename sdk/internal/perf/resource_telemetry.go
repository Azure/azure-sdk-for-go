// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"runtime"
)

type resourceTelemetrySnapshot struct {
	allocMiB      float64
	totalAllocMiB float64
	sysMiB        float64
	numGC         uint32
	goroutines    int
}

func captureResourceTelemetry() resourceTelemetrySnapshot {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return resourceTelemetrySnapshot{
		allocMiB:      bytesToMiB(mem.Alloc),
		totalAllocMiB: bytesToMiB(mem.TotalAlloc),
		sysMiB:        bytesToMiB(mem.Sys),
		numGC:         mem.NumGC,
		goroutines:    runtime.NumGoroutine(),
	}
}

func (s resourceTelemetrySnapshot) DiffSummary(after resourceTelemetrySnapshot) string {
	return fmt.Sprintf(
		"Resource telemetry: alloc(MiB)=%.2f totalAlloc(MiB)=%.2f sys(MiB)=%.2f gc=%d goroutines=%d",
		after.allocMiB-s.allocMiB,
		after.totalAllocMiB-s.totalAllocMiB,
		after.sysMiB-s.sysMiB,
		after.numGC-s.numGC,
		after.goroutines-s.goroutines,
	)
}

func bytesToMiB(v uint64) float64 {
	return float64(v) / (1024.0 * 1024.0)
}
