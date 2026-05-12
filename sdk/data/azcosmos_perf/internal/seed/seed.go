// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package seed creates and shares seeded Cosmos DB items for perf operations.
package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

// PerfItem is the JSON document used by the perf tool.
type PerfItem struct {
	ID           string `json:"id"`
	PartitionKey string `json:"partition_key"`
	Value        uint64 `json:"value"`
	Payload      string `json:"payload"`
}

// SeededItem identifies a seeded item by ID and partition key.
type SeededItem struct {
	ID           string
	PartitionKey string
}

// SharedItems is a thread-safe, capacity-capped pool of item identities.
type SharedItems struct {
	mu       sync.Mutex
	items    []SeededItem
	capacity int
}

// NewSharedItems creates a shared pool with capacity 2N.
func NewSharedItems(items []SeededItem) *SharedItems {
	copied := make([]SeededItem, len(items), len(items)*2)
	copy(copied, items)
	return &SharedItems{items: copied, capacity: cap(copied)}
}

// Random returns a random item from the pool.
func (s *SharedItems) Random() SeededItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.items[rand.Intn(len(s.items))]
}

// Sample returns up to n random items from the pool. Items may repeat.
func (s *SharedItems) Sample(n int) []SeededItem {
	if n < 1 {
		n = 1
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]SeededItem, n)
	for i := range result {
		result[i] = s.items[rand.Intn(len(s.items))]
	}
	return result
}

// Push adds an item, replacing a random existing entry when the pool is full.
func (s *SharedItems) Push(item SeededItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) < s.capacity {
		s.items = append(s.items, item)
		return
	}
	s.items[rand.Intn(len(s.items))] = item
}

// SeedContainer seeds count items concurrently and aborts on the first error.
func SeedContainer(ctx context.Context, container *azcosmos.ContainerClient, count, concurrency int) ([]SeededItem, error) {
	fmt.Printf("Seeding %d items (concurrency: %d)...\n", count, concurrency)

	items := make([]SeededItem, count)
	for i := 0; i < count; i++ {
		items[i] = SeededItem{ID: uuid.NewString(), PartitionKey: uuid.NewString()}
	}

	seedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan int)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	var completed atomic.Uint64

	worker := func() {
		defer wg.Done()
		for idx := range jobs {
			if seedCtx.Err() != nil {
				return
			}
			seeded := items[idx]
			item := PerfItem{
				ID:           seeded.ID,
				PartitionKey: seeded.PartitionKey,
				Value:        uint64(idx),
				Payload:      "perf-test-seed-payload",
			}
			body, err := json.Marshal(item)
			if err == nil {
				_, err = container.UpsertItem(seedCtx, azcosmos.NewPartitionKeyString(item.PartitionKey), body, nil)
			}
			if err != nil {
				select {
				case errCh <- fmt.Errorf("seed error for item %d: %w", idx, err):
					cancel()
				default:
				}
				return
			}
			done := completed.Add(1)
			if done%200 == 0 || int(done) == count {
				fmt.Printf("  Seeded %d/%d items\n", done, count)
			}
		}
	}

	if concurrency > count {
		concurrency = count
	}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		defer close(jobs)
		for i := range items {
			select {
			case <-seedCtx.Done():
				return
			case jobs <- i:
			}
		}
	}()

	wg.Wait()
	select {
	case err := <-errCh:
		return nil, err
	default:
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	fmt.Println("Seeding complete.")
	return items, nil
}
