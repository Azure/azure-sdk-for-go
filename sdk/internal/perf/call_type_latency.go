// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type callTypeLatencyCollector struct {
	mu     sync.Mutex
	values map[string][]time.Duration
}

func newCallTypeLatencyCollector() *callTypeLatencyCollector {
	return &callTypeLatencyCollector{values: map[string][]time.Duration{}}
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
	c.values[callType] = append(c.values[callType], duration)
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
	for k, v := range other.values {
		copied[k] = append([]time.Duration(nil), v...)
	}
	other.mu.Unlock()

	c.mu.Lock()
	for k, v := range copied {
		c.values[k] = append(c.values[k], v...)
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
