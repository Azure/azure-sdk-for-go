// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package mock

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/unstable/queryengine"
)

type PartitionKeyRange struct {
	ID           string `json:"id"`
	MinInclusive string `json:"minInclusive"`
	MaxExclusive string `json:"maxExclusive"`
}

type pkRanges struct {
	PartitionKeyRanges []PartitionKeyRange `json:"PartitionKeyRanges"`
}

type documentPayload[T any] struct {
	Documents []T `json:"Documents"`
}

// MockItem is the type of an item in the "mock" pipeline.
// The "mock" pipeline just merges items from each partition according to their MergeOrder value.
type MockItem struct {
	// ID is the ID of the item.
	ID string `json:"id"`

	// PKRangeId is the partition key range ID of the item.
	PartitionKey string `json:"partitionKey"`

	// MergeOrder is the universal cross-partition order in which the item should be merged.
	MergeOrder int `json:"mergeOrder"`
}

// MockQueryEngine is a mock implementation of the QueryEngine interface.
// This is a VERY rudimentary implementation that emulates the handling of the following query:
// `SELECT * FROM c ORDER BY c.mergeOrder`
// The intent here is to test how the Go SDK interacts with the query engine, not to test the query engine itself.
type MockQueryEngine struct {
	CreateError error
}

// NewMockQueryEngine creates a new MockQueryEngine.
func NewMockQueryEngine() *MockQueryEngine {
	return &MockQueryEngine{}
}

// CreateQueryPipeline creates a new query pipeline for the specified query and partition topology.
func (m *MockQueryEngine) CreateQueryPipeline(query string, plan string, pkranges string) (queryengine.QueryPipeline, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}

	var ranges pkRanges
	if err := json.Unmarshal([]byte(pkranges), &ranges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal partition key ranges: %w", err)
	}
	return newMockQueryPipeline(query, ranges.PartitionKeyRanges), nil
}

func (m *MockQueryEngine) SupportedFeatures() string {
	// We need to return whatever is necessary for the gateway to return a query plan for the query `SELECT * FROM c ORDER BY c.mergeOrder`
	return "OrderBy"
}

type partitionState struct {
	PartitionKeyRange
	started          bool
	queue            []MockItem
	nextContinuation string
}

// IsExhausted returns true if the partition is exhausted.
// A partition is considered exhausted if all the below are true:
// 1. The queue is empty (no more items to return).
// 2. The partition has started (we've received at least one response for it from the server).
// 3. The next continuation token is empty (the last response indicated that there are no more items on the server).
func (m *partitionState) IsExhausted() bool {
	return len(m.queue) == 0 && m.started && m.nextContinuation == ""
}

// ProvideData inserts new items into the queue, and updates the current continuation token.
func (p *partitionState) ProvideData(items []MockItem, continuation string) {
	p.started = true
	p.nextContinuation = continuation
	p.queue = append(p.queue, items...)
}

// PopItem removes the first item from the queue and returns it as a serialized JSON object.
func (p *partitionState) PopItem() ([]byte, error) {
	if len(p.queue) == 0 {
		return nil, fmt.Errorf("no items in queue")
	}
	item := p.queue[0]
	p.queue = p.queue[1:]
	serialized, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize item: %w", err)
	}
	return serialized, nil
}

type MockQueryPipeline struct {
	query          string
	completed      bool
	IsClosed       bool
	partitionState []partitionState
}

func newMockQueryPipeline(query string, partitions []PartitionKeyRange) *MockQueryPipeline {
	partState := make([]partitionState, 0, len(partitions))
	for _, partition := range partitions {
		partState = append(partState, partitionState{
			PartitionKeyRange: partition,
			started:           false,
			queue:             nil,
			nextContinuation:  "",
		})
	}

	return &MockQueryPipeline{
		query:          query,
		IsClosed:       false,
		partitionState: partState,
	}
}

func (m *MockQueryPipeline) Close() {
	m.IsClosed = true
}

func (m *MockQueryPipeline) IsComplete() bool {
	return m.completed
}

// NextBatch returns the next batch of items from the pipeline, as well as any requests needed to collect more data.
func (m *MockQueryPipeline) NextBatch() ([][]byte, []queryengine.QueryRequest, error) {
	if m.IsClosed {
		return nil, nil, fmt.Errorf("pipeline is closed")
	}

	items := make([][]byte, 0)

	// Loop, merging items from each partition, until all partitions are exhausted, or we need more data to conitnue.
	for {
		// Iterate through each partition to find the item with the lowest MergeOrder.
		var lowestMergeOrder int
		var lowestPartition *partitionState
		for i := range m.partitionState {
			// If any partition hasn't started yet, we can't return any items.
			if !m.partitionState[i].started {
				return nil, m.getRequests(), nil
			}

			if m.partitionState[i].IsExhausted() {
				// If this partition is exhausted, it won't contribute any more items, so we can skip it.
				continue
			}

			if len(m.partitionState[i].queue) > 0 && (lowestPartition == nil || m.partitionState[i].queue[0].MergeOrder < lowestMergeOrder) {
				lowestMergeOrder = m.partitionState[i].queue[0].MergeOrder
				lowestPartition = &m.partitionState[i]
			}
		}

		if lowestPartition == nil {
			// All partitions are either exhausted or have no items in the queue, so we need to make requests to get more data.
			break
		} else {
			// Add the item to the result set and remove it from the queue.
			item, err := lowestPartition.PopItem()
			if err != nil {
				return nil, nil, err
			}
			items = append(items, item)
		}

		// If we got here, we added an item to the result set, and we need to go back and check all the partitions again.
	}

	requests := m.getRequests()

	if len(items) == 0 && len(requests) == 0 {
		// If we didn't get any items and have no requests, we're done.
		m.completed = true
	}

	return items, requests, nil
}

// getRequests returns a list of all the QueryRequests that are needed to get the next batch of items.
func (m *MockQueryPipeline) getRequests() []queryengine.QueryRequest {
	requests := make([]queryengine.QueryRequest, 0, len(m.partitionState))
	for i := range m.partitionState {
		if m.partitionState[i].IsExhausted() {
			// If this partition is exhausted, we can't return any items.
			continue
		}

		continuation := ""
		if m.partitionState[i].started {
			continuation = m.partitionState[i].nextContinuation
		}

		requests = append(requests, queryengine.QueryRequest{
			PartitionKeyRangeID: m.partitionState[i].ID,
			Continuation:        continuation,
		})
	}
	return requests
}

// ProvideData is used by the SDK to provide incoming single-partition results to the pipeline.
// The items are expected to be ordered by the query's ORDER BY clause.
func (m *MockQueryPipeline) ProvideData(data queryengine.QueryResult) error {
	if m.IsClosed {
		return fmt.Errorf("pipeline is closed")
	}

	// Parse the items
	var payload documentPayload[MockItem]
	if err := json.Unmarshal(data.Data, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal items: %w", err)
	}

	// Find the partition state for the given partition key range ID and insert the items.
	for i := range m.partitionState {
		if m.partitionState[i].ID == data.PartitionKeyRangeID {
			m.partitionState[i].ProvideData(payload.Documents, data.NextContinuation)
			return nil
		}
	}

	// If we didn't find the partition key range ID, return an error.
	return fmt.Errorf("no partition found with ID %s", data.PartitionKeyRangeID)
}

func (m *MockQueryPipeline) Query() string {
	return m.query
}

var _ queryengine.QueryPipeline = &MockQueryPipeline{}
