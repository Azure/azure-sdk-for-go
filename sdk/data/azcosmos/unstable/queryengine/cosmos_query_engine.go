// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package queryengine

type QueryEngine interface {
	CreateQueryPipeline(query string, plan string, pkranges string) (QueryPipeline, error)
	SupportedFeatures() string
}

type QueryRequest struct {
	PartitionKeyRangeID string
	Continuation        string
}

type QueryResult struct {
	PartitionKeyRangeID string
	NextContinuation    string
	Data                []byte
}

type PipelineResult struct {
	IsCompleted bool
	Items       [][]byte
	Requests    []QueryRequest
}

type QueryPipeline interface {
	// Query returns the query text, possibly rewritten by the gateway, which will be used for per-partition queries.
	Query() string
	// IsComplete gets a boolean indicating if the pipeline has concluded
	IsComplete() bool
	// NextBatch gets the next batch of items, which will be empty if there are no more items in the buffer, and the next set of QueryRequests which must be fulfilled, which will be empty if there are no more requests.
	// If both the items and requests are empty, the pipeline has concluded.
	NextBatch() ([][]byte, []QueryRequest, error)
	// ProvideData provides more data for a given partition key range ID, using data retrieved from the server in response to making a DataRequest.
	ProvideData(data QueryResult) error
	// Close frees the resources associated with the pipeline.
	Close()
}
