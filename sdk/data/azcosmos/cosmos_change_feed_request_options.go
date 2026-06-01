// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"strconv"
	"time"
)

// ChangeFeedOptions defines the options for retrieving the change feed.
// Incorporate Continuation
type ChangeFeedOptions struct {
	// MaxItemCount limits the number of items returned per page.
	// Valid values are > 0. The service may return fewer items than requested.
	MaxItemCount int32

	// StartFrom is a user-friendly way to specify the time for change feed
	// Will be set to the IfModifiedSince header
	StartFrom *time.Time

	// PartitionKey is the logical partition key value for the request.
	// Use this to read from a specific logical partition.
	PartitionKey *PartitionKey

	// Feed Range specifies the range of pk values that map to a logical partition.
	FeedRange *FeedRange

	// CompositeContinuation is used to continue reading the change feed from a specific point.
	Continuation *string

	// PriorityLevel overrides the client-level default priority for this operation.
	// Valid values are PriorityLevelHigh and PriorityLevelLow.
	PriorityLevel *PriorityLevel
	// ThroughputBucket overrides the client-level default throughput bucket for this operation.
	// For more information, see https://aka.ms/CosmosDB/ThroughputBuckets
	// The valid range is 1 to 5 (inclusive).
	ThroughputBucket *int32
}

// buildRequestHeaders constructs the exact headers needed for a single
// change-feed request against one queue head. Pure builder: callers MUST
// supply the head and the already-resolved PK-range ID (overlap-matched
// via the routing map, NOT exact-matched).
//
// This is the path used by the new queue-driven GetChangeFeed loop. The
// caller-options-level Continuation token is NOT consulted here — the
// queue-head ETag drives If-None-Match because the queue may have been
// split-expanded since the token was issued.
//
// Returns an error when PartitionKey serialization fails — sending a
// change-feed read with a missing PK header would yield an opaque
// server-side error, so we surface the cause to the caller instead.
func (options *ChangeFeedOptions) buildRequestHeaders(head changeFeedRange, resolvedPKRangeID string) (map[string]string, error) {
	headers := make(map[string]string, 6)
	headers[cosmosHeaderChangeFeed] = cosmosHeaderValuesChangeFeed

	if options != nil {
		if options.MaxItemCount > 0 {
			headers[cosmosHeaderMaxItemCount] = strconv.FormatInt(int64(options.MaxItemCount), 10)
		}
		if options.StartFrom != nil {
			headers[cosmosHeaderIfModifiedSince] = options.StartFrom.UTC().Format(time.RFC1123)
		}
		if options.PartitionKey != nil {
			pkJSON, err := options.PartitionKey.toJsonString()
			if err != nil {
				return nil, fmt.Errorf("ChangeFeedOptions: serializing PartitionKey: %w", err)
			}
			headers[cosmosHeaderPartitionKey] = string(pkJSON)
		}
	}

	if head.ContinuationToken != nil && *head.ContinuationToken != "" {
		headers[headerIfNoneMatch] = string(*head.ContinuationToken)
	}

	if resolvedPKRangeID != "" {
		headers[headerXmsDocumentDbPartitionKeyRangeId] = resolvedPKRangeID
	}

	return headers, nil
}
