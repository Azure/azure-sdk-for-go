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
	if duration <= 0 {
		return
	}
	if callType == "" {
		callType = "operation"
	}
	c.mu.Lock()
	c.values[callType] = append(c.values[callType], duration)
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
