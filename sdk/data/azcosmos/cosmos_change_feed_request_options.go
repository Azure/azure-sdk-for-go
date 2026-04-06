// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
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
}

func (options *ChangeFeedOptions) toHeaders(partitionKeyRanges []partitionKeyRange) *map[string]string {
	headers := make(map[string]string)

	headers[cosmosHeaderChangeFeed] = cosmosHeaderValuesChangeFeed

	if options.MaxItemCount > 0 {
		headers[cosmosHeaderMaxItemCount] = strconv.FormatInt(int64(options.MaxItemCount), 10)
	}

	if options.StartFrom != nil {
		formatted := options.StartFrom.UTC().Format(time.RFC1123)
		headers[cosmosHeaderIfModifiedSince] = formatted
	}

	if options.Continuation != nil && *options.Continuation != "" {
		var compositeToken compositeContinuationToken
		if err := json.Unmarshal([]byte(*options.Continuation), &compositeToken); err == nil && len(compositeToken.Continuation) > 0 {
			if compositeToken.Continuation[0].ContinuationToken != nil {
				headers[headerIfNoneMatch] = string(*compositeToken.Continuation[0].ContinuationToken)
			}
			if options.FeedRange == nil {
				options.FeedRange = &FeedRange{
					MinInclusive: compositeToken.Continuation[0].MinInclusive,
					MaxExclusive: compositeToken.Continuation[0].MaxExclusive,
				}
			}
		} else {
			headers[headerIfNoneMatch] = *options.Continuation
		}
	}

	if options.PartitionKey != nil {
		partitionKeyJSON, err := options.PartitionKey.toJsonString()
		if err == nil {
			headers[cosmosHeaderPartitionKey] = string(partitionKeyJSON)
		}
	}

	if options.FeedRange != nil && len(partitionKeyRanges) > 0 {
		if id, err := findPartitionKeyRangeID(*options.FeedRange, partitionKeyRanges); err == nil {
			headers[headerXmsDocumentDbPartitionKeyRangeId] = id
		} else {
			return nil
		}
	}

	if len(headers) == 0 {
		return nil
	}

	return &headers
}
