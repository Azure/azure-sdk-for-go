// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var callTypeCollectorSeed int64

type callTypeLatencyCollector struct {
	mu     sync.Mutex
	values map[string][]time.Duration
	seen   map[string]int64
	max    int
	rng    *rand.Rand
}

func newCallTypeLatencyCollector() *callTypeLatencyCollector {
	return newBoundedCallTypeLatencyCollector(0)
}

func newBoundedCallTypeLatencyCollector(max int) *callTypeLatencyCollector {
	return &callTypeLatencyCollector{
		values: map[string][]time.Duration{},
		seen:   map[string]int64{},
		max:    max,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano() + atomic.AddInt64(&callTypeCollectorSeed, 1))),
	}
}

func (c *callTypeLatencyCollector) Add(callType string, duration time.Duration) {
	// A zero-duration operation is legitimate (an op can complete within the
	// timer's resolution), so only negative samples are discarded.
	if duration < 0 {
		return
	}
	if callType == "" {
		callType = "operation"
	}
	c.mu.Lock()
	c.seen[callType]++
	if c.max < 0 {
		c.mu.Unlock()
		return
	}
	values := c.values[callType]
	if c.max == 0 || len(values) < c.max {
		c.values[callType] = append(values, duration)
	} else {
		index := c.rng.Int63n(c.seen[callType])
		if index < int64(c.max) {
			values[index] = duration
		}
	}
	c.mu.Unlock()
}

// MergeFrom folds all per-call-type samples from other into c. It is used to
// merge a per-worker collector into the shared runner collector after the
// measurement phase completes.
func (c *callTypeLatencyCollector) MergeFrom(other *callTypeLatencyCollector) {
	if other == nil {
		return
	}
	other.mu.Lock()
	copied := make(map[string][]time.Duration, len(other.values))
	copiedSeen := make(map[string]int64, len(other.seen))
	for k, v := range other.values {
		copied[k] = append([]time.Duration(nil), v...)
		copiedSeen[k] = other.seen[k]
	}
	other.mu.Unlock()

	c.mu.Lock()
	for k, v := range copied {
		c.values[k], c.seen[k] = mergeReservoirs(c.values[k], c.seen[k], v, copiedSeen[k], c.max, c.rng)
	}
	c.mu.Unlock()
}

func (c *callTypeLatencyCollector) Summary() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.values) == 0 {
		return "Latency by call type: no data"
	}

	callTypes := make([]string, 0, len(c.values))
	for key := range c.values {
		callTypes = append(callTypes, key)
	}
	sort.Strings(callTypes)

	lines := make([]string, 0, len(callTypes)+1)
	lines = append(lines, "Latency by call type (ms):")
	for _, key := range callTypes {
		vals := append([]time.Duration(nil), c.values[key]...)
		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
		line := fmt.Sprintf("  %s: p50=%.2f p95=%.2f p99=%.2f", key, toMS(percentile(vals, 50)), toMS(percentile(vals, 95)), toMS(percentile(vals, 99)))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
