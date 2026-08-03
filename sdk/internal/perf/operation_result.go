// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// defaultMaxOperationResults bounds how many per-operation result records are
// retained in memory when --results-file is used. Once the cap is reached the
// collector keeps a uniform random sample of the stream (reservoir sampling)
// so a long, high-throughput run cannot grow memory without bound. A value of
// 0 disables the cap (unbounded retention). Exposed via --max-results.
const defaultMaxOperationResults = 1_000_000

// operationResultsSeed offsets the per-collector PRNG seed so reservoirs in
// concurrently-created worker collectors don't sample identically.
var operationResultsSeed int64

type operationResult struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	LatencyMs float64   `json:"latencyMs"`
	SizeBytes int64     `json:"sizeBytes"`
}

type operationResultsCollector struct {
	mu      sync.Mutex
	results []operationResult
	// max is the maximum number of retained records; 0 means unbounded.
	max int
	// seen is the total number of records observed (including those evicted
	// by reservoir sampling) and drives the replacement probability.
	seen int64
	// rng is a per-collector PRNG used only on the (off-network) sampling
	// path. A dedicated generator avoids the global rand mutex in the hot loop.
	rng *rand.Rand
}

// newOperationResultsCollector returns a collector that retains at most max
// records (0 = unbounded), keeping a uniform random sample beyond that.
func newOperationResultsCollector(max int) *operationResultsCollector {
	seed := time.Now().UnixNano() + atomic.AddInt64(&operationResultsSeed, 1)
	return &operationResultsCollector{
		max: max,
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (c *operationResultsCollector) Add(operation string, latency time.Duration, sizeBytes int64) {
	// A zero-duration operation is legitimate (an op can complete within the
	// timer's resolution), so only negative samples are discarded.
	if latency < 0 {
		return
	}
	if operation == "" {
		operation = "operation"
	}
	rec := operationResult{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		LatencyMs: toMS(latency),
		SizeBytes: sizeBytes,
	}
	c.mu.Lock()
	c.addLocked(rec)
	c.mu.Unlock()
}

// addLocked appends rec, or once the cap is reached applies reservoir sampling
// (Vitter's Algorithm R) to keep a uniform random sample in O(1) time and
// bounded memory. It performs no I/O or blocking, so the measured hot loop is
// not perturbed. The caller must hold c.mu.
func (c *operationResultsCollector) addLocked(rec operationResult) {
	c.seen++
	if c.max <= 0 || len(c.results) < c.max {
		c.results = append(c.results, rec)
		return
	}
	// Replace a random retained element with probability max/seen.
	j := c.rng.Int63n(c.seen)
	if j < int64(c.max) {
		c.results[j] = rec
	}
}

// MergeFrom folds a per-worker collector into the shared runner collector
// after the measurement phase completes. When a cap is set the merge stays
// bounded: the worker's already-sampled records are fed through this
// reservoir, and seen is advanced by the worker's full observed count so
// workers that processed more operations are weighted accordingly.
func (c *operationResultsCollector) MergeFrom(other *operationResultsCollector) {
	if other == nil {
		return
	}
	other.mu.Lock()
	copied := append([]operationResult(nil), other.results...)
	otherSeen := other.seen
	other.mu.Unlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.max <= 0 {
		c.results = append(c.results, copied...)
		c.seen += otherSeen
		return
	}

	for _, rec := range copied {
		c.addLocked(rec)
	}
	if extra := otherSeen - int64(len(copied)); extra > 0 {
		c.seen += extra
	}
}

func (c *operationResultsCollector) WriteJSON(path string) error {
	c.mu.Lock()
	data := append([]operationResult(nil), c.results...)
	c.mu.Unlock()

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal operation results: %w", err)
	}
	if err = os.WriteFile(path, b, 0o600); err != nil {
		return fmt.Errorf("failed to write operation results file %s: %w", path, err)
	}
	return nil
}
