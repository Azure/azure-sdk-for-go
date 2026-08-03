// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type latencyCollector struct {
	mu        sync.Mutex
	durations []time.Duration
}

func (l *latencyCollector) Add(d time.Duration) {
	// A zero-duration operation is legitimate (an op can complete within the
	// timer's resolution), so only negative samples are discarded.
	if d < 0 {
		return
	}
	l.mu.Lock()
	l.durations = append(l.durations, d)
	l.mu.Unlock()
}

// MergeFrom appends all samples from other into l. It is used to fold a
// per-worker collector (populated lock-free during the hot loop) into the
// shared runner collector after the measurement phase completes.
func (l *latencyCollector) MergeFrom(other *latencyCollector) {
	if other == nil {
		return
	}
	other.mu.Lock()
	samples := append([]time.Duration(nil), other.durations...)
	other.mu.Unlock()

	l.mu.Lock()
	l.durations = append(l.durations, samples...)
	l.mu.Unlock()
}

func (l *latencyCollector) Summary() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.durations) == 0 {
		return "Latency: no data"
	}

	vals := append([]time.Duration(nil), l.durations...)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

	return fmt.Sprintf(
		"Latency (ms): p50=%.2f p90=%.2f p95=%.2f p99=%.2f max=%.2f",
		toMS(percentile(vals, 50)),
		toMS(percentile(vals, 90)),
		toMS(percentile(vals, 95)),
		toMS(percentile(vals, 99)),
		toMS(vals[len(vals)-1]),
	)
}

func percentile(vals []time.Duration, p int) time.Duration {
	if len(vals) == 0 {
		return 0
	}
	if len(vals) == 1 {
		return vals[0]
	}

	idx := (float64(p) / 100.0) * float64(len(vals)-1)
	low := int(idx)
	high := low + 1
	if high >= len(vals) {
		return vals[len(vals)-1]
	}
	frac := idx - float64(low)
	return time.Duration(float64(vals[low])*(1-frac) + float64(vals[high])*frac)
}

func toMS(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}
