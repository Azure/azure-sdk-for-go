// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	// cosmosDefaultMaxItemCount represents unlimited items in the response
	cosmosDefaultMaxItemCount = -1
	cosmosBaseTen             = 10
)

// ChangeFeedResourceType represents the resource type for change feed operations.
type ChangeFeedResourceType int

const (
	ChangeFeedResourceTypeContainer ChangeFeedResourceType = iota
	ChangeFeedResourceTypePartitionKey
	ChangeFeedResourceTypeFeedRange
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

	// ResourceType tracks the resource type based on FeedRange, PartitionKey, or neither.
	resourceType ChangeFeedResourceType
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
		var continuationTokenForPartitionKey continuationTokenForPartitionKey
		if err := json.Unmarshal([]byte(*options.Continuation), &compositeToken); err == nil && len(compositeToken.Continuation) > 0 {
			// If the continuation token is a composite token, parse it
			if compositeToken.Continuation[0].ContinuationToken != nil {
				headers[headerIfNoneMatch] = string(*compositeToken.Continuation[0].ContinuationToken)
			}
			if options.FeedRange == nil {
				options.FeedRange = &FeedRange{
					MinInclusive: compositeToken.Continuation[0].MinInclusive,
					MaxExclusive: compositeToken.Continuation[0].MaxExclusive,
				}
			}
		} else if err := json.Unmarshal([]byte(*options.Continuation), &continuationTokenForPartitionKey); err == nil {
			// If the continuation token is for a partition key, parse it
			if continuationTokenForPartitionKey.Continuation != nil {
				headers[headerIfNoneMatch] = string(*continuationTokenForPartitionKey.Continuation)
			}
			if options.PartitionKey == nil && continuationTokenForPartitionKey.PartitionKey != nil {
				pkBytes, err := json.Marshal(continuationTokenForPartitionKey.PartitionKey)
				if err == nil {
					var pk PartitionKey
					err = json.Unmarshal(pkBytes, &pk)
					if err == nil {
						options.PartitionKey = &pk
					}
				}
			}
		} else {
			// If the continuation token is a simple ETag, use it directly
			headers[headerIfNoneMatch] = *options.Continuation
		}
	}

	if options.FeedRange != nil && len(partitionKeyRanges) > 0 {
		if id, err := findPartitionKeyRangeID(*options.FeedRange, partitionKeyRanges); err == nil {
			headers[headerXmsDocumentDbPartitionKeyRangeId] = id
		} else {
			return nil
		}
	}

	if options.PartitionKey != nil {
		partitionKeyJSON, err := json.Marshal(options.PartitionKey)
		if err == nil {
			headers[cosmosHeaderPartitionKey] = string(partitionKeyJSON)
		}
	}

	if len(headers) == 0 {
		return nil
	}

	return &headers
}

// SetResourceType sets the ResourceType based on which option is set.
func (options *ChangeFeedOptions) SetResourceType() {
	if options.FeedRange != nil {
		options.resourceType = ChangeFeedResourceTypeFeedRange
	} else if options.PartitionKey != nil {
		options.resourceType = ChangeFeedResourceTypePartitionKey
	} else {
		options.resourceType = ChangeFeedResourceTypeContainer
	}
}
