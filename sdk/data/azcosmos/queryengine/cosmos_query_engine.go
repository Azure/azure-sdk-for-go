// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package queryengine

// QueryEngine is an interface that defines the methods for a query engine.
type QueryEngine interface {
	CreateQueryPipeline(query string, plan string, pkranges string) (QueryPipeline, error)
	CreateReadManyPipeline(items []ItemIdentity, pkranges string, pkKind string, pkVersion uint8, pkPaths []string) (QueryPipeline, error)
	SupportedFeatures() string
}

// ItemIdentity contains the unique identifiers for an item in a container.
type ItemIdentity struct {
	// json string representation of the partition key value
	PartitionKeyValue string
	// ID of the item to read
	ID string
}

// QueryRequest describes a request from the pipeline for data from a specific partition key range.
type QueryRequest struct {
	// PartitionKeyRangeID is the ID of the partition key range from which data is requested.
	PartitionKeyRangeID string
	// The ID of this request, within the partition key range.
	//
	// Opaque identifier that must be provided back to the pipeline when providing data.
	Id uint64
	// Continuation is the continuation token to use in the request.
	Continuation string
	// The query to execute for this partition key range, if different from the original query.
	Query string
	// If a query is specified, this flag indicates if the query parameters should be included with that query.
	//
	// Sometimes, when an override query is specified, it differs in structure from the original query, and the original parameters are not valid.
	IncludeParameters bool
	// If specified, indicates that the SDK should IMMEDIATELY drain all remaining results from this partition key range, following continuation tokens, until no more results are available.
	// All the data from this partition key range should be provided BEFORE any new items will be made available.
	// The data may be provided in multiple QueryResults, but every result correlated to this request should have the same RequestId value.
	//
	// This allows engines to optimize for non-streaming scenarios, where the entire result set must be provided to the engine before it can make progress.
	Drain bool
}

// QueryResult contains the result of a query for a specific partition key range.
type QueryResult struct {
	// The ID of the partition key range that was queried.
	PartitionKeyRangeID string
	// The ID of the QueryRequest that generated this result.
	RequestId uint64
	// The continuation token to be used for the next request, if any.
	NextContinuation string
	// The raw body of the response from the query.
	Data []byte
}

// NewQueryRequest creates a new QueryRequest with the specified partition key range ID, continuation token, and data.
func NewQueryResult(partitionKeyRangeID string, data []byte, continuation string) QueryResult {
	return QueryResult{
		PartitionKeyRangeID: partitionKeyRangeID,
		Data:                data,
		NextContinuation:    continuation,
	}
}

// NewQueryRequestString creates a new QueryRequest with the specified partition key range ID, continuation token, and data (as a string).
func NewQueryResultString(partitionKeyRangeID string, data string, continuation string) QueryResult {
	return NewQueryResult(partitionKeyRangeID, []byte(data), continuation)
}

// PipelineResult contains the result of running a single turn of the query pipeline.
type PipelineResult struct {
	// IsCompleted indicates if the pipeline has completed processing.
	IsCompleted bool

	// Items contains the items returned by the pipeline.
	Items [][]byte

	// Requests contains the requests made by the pipeline for more data.
	Requests []QueryRequest
}

// QueryPipeline is an interface that defines the methods for a query pipeline.
type QueryPipeline interface {
	// Query returns the query text, possibly rewritten by the gateway, which will be used for per-partition queries.
	Query() string
	// IsComplete gets a boolean indicating if the pipeline has concluded
	IsComplete() bool
	// Run executes a single turn of the pipeline, yielding a PipelineResult containing the items and requests for more data.
	Run() (*PipelineResult, error)
	// Data from multiple partition ranges may be provided at once.
	// However, each page of data must be provided in order.
	// So, for any given partition key range, page n's results must be earlier in the `data` slice than page n+1's results.
	// Data from different partition key ranges may be interleaved, as long as each partition key range's pages are in order.
	//
	// The pipeline will use the QueryResult.RequestId field to validate this.
	//
	// When providing data from a draining request (i.e. a request with Drain set to true), all pages for that draining request can share the same QueryResult.RequestId.
	ProvideData(data []QueryResult) error
	// Close frees the resources associated with the pipeline.
	Close()
}
