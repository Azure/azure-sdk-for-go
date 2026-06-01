// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type operationResult struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	LatencyMs float64   `json:"latencyMs"`
	SizeBytes int64     `json:"sizeBytes"`
}

type operationResultsCollector struct {
	mu      sync.Mutex
	results []operationResult
}

func (c *operationResultsCollector) Add(operation string, latency time.Duration, sizeBytes int64) {
	if latency <= 0 {
		return
	}
	if operation == "" {
		operation = "operation"
	}
	c.mu.Lock()
	c.results = append(c.results, operationResult{
		Timestamp: time.Now().UTC(),
		Operation: operation,
		LatencyMs: toMS(latency),
		SizeBytes: sizeBytes,
	})
	c.mu.Unlock()
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
