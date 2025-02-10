// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

type QueryEngine interface {
	CreateQueryPipeline(query string, plan string, pkranges string) (QueryPipeline, error)
	SupportedFeatures() string
}

type DataRequest struct {
	PartitionKeyRangeID string
	Continuation        string
}

func (r *DataRequest) toHeaders() *map[string]string {
	headers := make(map[string]string)
	if r.Continuation != "" {
		headers[cosmosHeaderContinuationToken] = r.Continuation
	}
	headers[cosmosHeaderPartitionKeyRangeId] = r.PartitionKeyRangeID
	return &headers
}

type PipelineResult struct {
	IsCompleted bool
	Items       [][]byte
	Requests    []DataRequest
}

type QueryPipeline interface {
	// Query returns the query text, possibly rewritten by the gateway, which will be used for per-partition queries.
	Query() string
	// IsComplete gets a boolean indicating if the pipeline has concluded
	IsComplete() bool
	// NextBatch gets the next batch of items, which will be empty if there are no more items in the buffer, and the next set of DataRequests which must be fulfilled, which will be empty if there are no more requests.
	// If both the items and requests are empty, the pipeline has concluded.
	NextBatch(maxPageSize int32) ([][]byte, []DataRequest, error)
	// ProvideData provides more data for a given partition key range ID, using data retrieved from the server in response to making a DataRequest.
	ProvideData(partitionKeyRangeId string, data string, continuation string) error
	// Close frees the resources associated with the pipeline.
	Close()
}
