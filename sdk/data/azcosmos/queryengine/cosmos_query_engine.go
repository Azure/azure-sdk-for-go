// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package queryengine

// QueryEngine is an interface that defines the methods for a query engine.
type QueryEngine interface {
	CreateQueryPipeline(query string, plan string, pkranges string) (QueryPipeline, error)
	CreateReadManyPipeline(partitionKeyRangesData []byte, items []ItemIdentity, pkKind PartitionKeyKind, pkVersion int32) (ReadManyPipeline, error)
	SupportedFeatures() string
}

// QueryRequest describes a request from the pipeline for data from a specific partition key range.
type QueryRequest struct {
	// PartitionKeyRangeID is the ID of the partition key range from which data is requested.
	PartitionKeyRangeID string
	// Continuation is the continuation token to use in the request.
	Continuation string
}

// QueryResult contains the result of a query for a specific partition key range.
type QueryResult struct {
	PartitionKeyRangeID string
	NextContinuation    string
	Data                []byte
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
	// ProvideData provides more data for a given partition key range ID, using data retrieved from the server in response to making a DataRequest.
	ProvideData(data QueryResult) error
	// Close frees the resources associated with the pipeline.
	Close()
}

// QueryPipeline is an interface that defines the methods for a query pipeline.
type ReadManyPipeline interface {
	// Query returns the query text, possibly rewritten by the gateway, which will be used for per-partition queries.
	Query() string
	// IsComplete gets a boolean indicating if the pipeline has concluded
	IsComplete() bool
	// Run executes a single turn of the pipeline, yielding a PipelineResult containing the items and requests for more data.
	Run() (*PipelineResult, error)
	// ProvideData provides more data for a given partition key range ID, using data retrieved from the server in response to making a DataRequest.
	ProvideData(data QueryResult) error
	// Close frees the resources associated with the pipeline.
	Close()
}
