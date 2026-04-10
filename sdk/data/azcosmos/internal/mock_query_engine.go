// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package mock

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
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

// QueryRequestConfig controls what QueryRequest values the pipeline should return.
type QueryRequestConfig struct {
	// Optional query override to return in the per-partition QueryRequest
	Query             *string
	IncludeParameters bool
}

// MockQueryEngine is a mock implementation of the QueryEngine interface.
// This is a VERY rudimentary implementation that emulates the handling of the following query:
// `SELECT * FROM c ORDER BY c.mergeOrder`
// The intent here is to test how the Go SDK interacts with the query engine, not to test the query engine itself.
type MockQueryEngine struct {
	CreateError        error
	QueryRequestConfig *QueryRequestConfig
}

// NewMockQueryEngine creates a new MockQueryEngine.
func NewMockQueryEngine() *MockQueryEngine {
	return &MockQueryEngine{}
}

// WithQueryRequestConfig returns an engine preconfigured to return the specified query request override.
func WithQueryRequestConfig(cfg *QueryRequestConfig) *MockQueryEngine {
	return &MockQueryEngine{QueryRequestConfig: cfg}
}

// CreateQueryPipeline creates a new query pipeline for the specified query and partition topology.
func (m *MockQueryEngine) CreateQueryPipeline(query string, plan string, pkranges string) (queryengine.QueryPipeline, error) {
	// capture config for this pipeline
	var cfg *QueryRequestConfig
	if m.QueryRequestConfig != nil {
		c := *m.QueryRequestConfig
		cfg = &c
	}

	var ranges pkRanges
	if err := json.Unmarshal([]byte(pkranges), &ranges); err != nil {
		return nil, fmt.Errorf("failed to unmarshal partition key ranges: %w", err)
	}
	return newMockQueryPipeline(query, ranges.PartitionKeyRanges, cfg), nil
}

// CreateReadManyPipeline creates a read-many pipeline which returns the provided item identities
// serialized as JSON documents. This is a simplified pipeline used by tests to exercise the
// SDK's ReadMany->QueryEngine glue without making network calls for each item.
func (m *MockQueryEngine) CreateReadManyPipeline(items []queryengine.ItemIdentity, pkranges string, pkKind string, pkVersion uint8, pkPaths []string) (queryengine.QueryPipeline, error) {
	return &MockReadManyPipeline{items: items, completed: false, resultingItems: make([][]byte, 0, len(items))}, nil
}

// MockReadManyPipeline is a minimal QueryPipeline implementation for ReadMany tests.
type MockReadManyPipeline struct {
	items          []queryengine.ItemIdentity
	completed      bool
	resultingItems [][]byte
}

func (m *MockReadManyPipeline) Close() {
	m.completed = true
}

func (m *MockReadManyPipeline) IsComplete() bool {
	return m.completed
}

func (m *MockReadManyPipeline) Run() (*queryengine.PipelineResult, error) {
	if m.IsComplete() {
		return &queryengine.PipelineResult{IsCompleted: true, Items: m.resultingItems, Requests: nil}, nil
	}
	// first run return queries to execute
	requests := make([]queryengine.QueryRequest, 0, len(m.items))
	for i := range m.items {
		pk := m.items[i].PartitionKeyValue
		create_query := fmt.Sprintf("Select * from c where c.id = '%s' and c.pk = '%s'", m.items[i].ID, pk)
		requests = append(requests, queryengine.QueryRequest{
			Query: create_query,
		})
	}

	// second run return result
	m.completed = true
	return &queryengine.PipelineResult{IsCompleted: true, Items: nil, Requests: requests}, nil
}

func (m *MockReadManyPipeline) ProvideData(data []queryengine.QueryResult) error {
	for _, res := range data {
		m.resultingItems = append(m.resultingItems, res.Data)
	}
	return nil
}

func (m *MockReadManyPipeline) Query() string {
	return ""
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
	nextIndex        uint64
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
	query              string
	completed          bool
	IsClosed           bool
	partitionState     []partitionState
	queryRequestConfig *QueryRequestConfig
}

func newMockQueryPipeline(query string, partitions []PartitionKeyRange, cfg *QueryRequestConfig) *MockQueryPipeline {
	partState := make([]partitionState, 0, len(partitions))
	for _, partition := range partitions {
		partState = append(partState, partitionState{
			PartitionKeyRange: partition,
			started:           false,
			queue:             nil,
			nextContinuation:  "",
			nextIndex:         0,
		})
	}

	return &MockQueryPipeline{
		query:              query,
		IsClosed:           false,
		partitionState:     partState,
		queryRequestConfig: cfg,
	}
}

func (m *MockQueryPipeline) Close() {
	m.IsClosed = true
}

func (m *MockQueryPipeline) IsComplete() bool {
	return m.completed
}

// NextBatch returns the next batch of items from the pipeline, as well as any requests needed to collect more data.
func (m *MockQueryPipeline) Run() (*queryengine.PipelineResult, error) {
	if m.IsClosed {
		return nil, fmt.Errorf("pipeline is closed")
	}

	items := make([][]byte, 0)

	// Loop, merging items from each partition, until all partitions are exhausted, or we need more data to continue.
	for {
		// Iterate through each partition to find the item with the lowest MergeOrder.
		var lowestMergeOrder int
		var lowestPartition *partitionState
		for i := range m.partitionState {
			// If any partition hasn't started yet, we can't return any items.
			if !m.partitionState[i].started {
				return &queryengine.PipelineResult{
					IsCompleted: false,
					Items:       nil,
					Requests:    m.getRequests(),
				}, nil
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
				return nil, err
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

	return &queryengine.PipelineResult{
		IsCompleted: m.completed,
		Items:       items,
		Requests:    requests,
	}, nil
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

		// Respect any per-pipeline override for the query and include-parameters flag.
		q := ""
		includeParams := false
		if m.queryRequestConfig != nil {
			if m.queryRequestConfig.Query != nil {
				q = *m.queryRequestConfig.Query
			}
			includeParams = m.queryRequestConfig.IncludeParameters
		}

		requests = append(requests, queryengine.QueryRequest{
			PartitionKeyRangeID: m.partitionState[i].ID,
			Id:                  m.partitionState[i].nextIndex,
			Continuation:        continuation,
			Query:               q,
			IncludeParameters:   includeParams,
			Drain:               false,
		})
	}
	return requests
}

// ProvideData is used by the SDK to provide incoming single-partition results to the pipeline.
// The items are expected to be ordered by the query's ORDER BY clause.
func (m *MockQueryPipeline) ProvideData(data []queryengine.QueryResult) error {
	if m.IsClosed {
		return fmt.Errorf("pipeline is closed")
	}

	// Parse the items
	var payload documentPayload[MockItem]
	if err := json.Unmarshal(data[0].Data, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal items: %w", err)
	}

	// Find the partition state for the given partition key range ID and insert the items.
	for i := range m.partitionState {
		if m.partitionState[i].ID == data[0].PartitionKeyRangeID {
			// Validate request ordering: the provided result must match the expected nextIndex.
			if m.partitionState[i].nextIndex != data[0].RequestId {
				return fmt.Errorf("out of order data provided for partition key range %s: expected index %d, got %d", data[0].PartitionKeyRangeID, m.partitionState[i].nextIndex, data[0].RequestId)
			}
			// advance expected index for next request
			m.partitionState[i].nextIndex++
			m.partitionState[i].ProvideData(payload.Documents, data[0].NextContinuation)
			return nil
		}
	}

	// If we didn't find the partition key range ID, return an error.
	return fmt.Errorf("no partition found with ID %s", data[0].PartitionKeyRangeID)
}

func (m *MockQueryPipeline) Query() string {
	return m.query
}

var _ queryengine.QueryPipeline = &MockQueryPipeline{}
