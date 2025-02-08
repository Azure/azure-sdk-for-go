// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

type QueryEngine interface {
	CreateQueryPipeline(plan []byte, pkranges []byte) (QueryPipeline, error)
}

type DataRequest struct {
	PartitionKeyRangeID string
	Continuation        *string
}

func (r *DataRequest) toHeaders() *map[string]string {
	headers := make(map[string]string)
	if r.Continuation != nil {
		headers[cosmosHeaderContinuationToken] = *r.Continuation
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
	// GetRewrittenQuery returns the query text, possibly rewritten by the gateway, which will be used for per-partition queries.
	GetRewrittenQuery() string
	// IsComplete gets a boolean indicating if the pipeline has concluded
	IsComplete() bool
	// NextBatch gets the next batch of items, which will be empty if there are no more items in the buffer.
	// The number of items retrieved will be capped by the provided maxPageSize if it is positive.
	// Any remaining items will be returned by the next call to NextBatch.
	NextBatch(maxPageSize int32) ([][]byte, error)
	// NextRequests gets the next batch of DataRequests which must be fulfilled, using ProvideData, before more data can be retrieved using NextBatch.
	NextRequests() ([]DataRequest, error)
	// ProvideData provides more data for a given partition key range ID, using data retrieved from the server in response to making a DataRequest.
	ProvideData(partitionKeyRangeId string, data []byte, continuation string) error
}
