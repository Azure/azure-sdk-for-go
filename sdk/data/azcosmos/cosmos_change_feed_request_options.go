// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

const (
	// cosmosDefaultMaxItemCount represents unlimited items in the response
	cosmosDefaultMaxItemCount = -1
)

// ChangeFeedOptions defines the options for retrieving the change feed.
// Incorporate Continuation
type ChangeFeedOptions struct {
	// MaxItemCount limits the number of items returned per page.
	// Valid values are > 0. The service may return fewer items than requested.
	MaxItemCount int32

	// ChangeFeedStartFrom is a user-friendly way to specify the time for change feed
	// Will be set to the IfModifiedSince header
	ChangeFeedStartFrom *time.Time

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

	// Always setting the AIM header to "Incremental Feed" for change feed requests
	headers[cosmosHeaderChangeFeed] = cosmosHeaderValuesChangeFeed

	// If MaxItemCount is set to a positive value, it will be included in the headers.
	// If it is 0, negative, or not set it will be set to -1 to indicate no limit.
	if options.MaxItemCount > 0 {
		headers[cosmosHeaderMaxItemCount] = strconv.FormatInt(int64(options.MaxItemCount), 10)
	} else {
		headers[cosmosHeaderMaxItemCount] = strconv.FormatInt(cosmosDefaultMaxItemCount, 10)
	}
	// Formats the time as RFC1123, e.g., "Mon, 02 Jan 2006 15:04:05 MST" (e.g., "Thu, 27 Jun 2025 14:30:00 UTC")
	// If ChangeFeedStartFrom is set, will internally map to If-Modified-Since
	if options.ChangeFeedStartFrom != nil {
		formatted := options.ChangeFeedStartFrom.UTC().Format(time.RFC1123)
		formatted = strings.Replace(formatted, "UTC", "GMT", 1)
		headers[cosmosHeaderIfModifiedSince] = formatted
	}

	// If PartitionKey is set, convert it to JSON and add it to the headers.
	if options.PartitionKey != nil {
		partitionKeyJSON, err := options.PartitionKey.toJsonString()
		if err == nil {
			headers[cosmosHeaderPartitionKey] = string(partitionKeyJSON)
		}
	}

	// If FeedRange struct is set, using function FindPartitionKeyRangeId to see if there is a 1:1 match
	if options.FeedRange != nil && len(partitionKeyRanges) > 0 {
		if id, err := FindPartitionKeyRangeID(*options.FeedRange, partitionKeyRanges); err == nil {
			headers[headerXmsDocumentDbPartitionKeyRangeId] = id
		} else {
			return nil
		}
	}

	// Handle composite continuation token
	if options.Continuation != nil && *options.Continuation != "" {
		// Try to parse as composite token first
		var compositeToken compositeContinuationToken
		if err := json.Unmarshal([]byte(*options.Continuation), &compositeToken); err == nil && len(compositeToken.Continuation) > 0 {
			// It's a composite token - extract the ETag from the first range
			if compositeToken.Continuation[0].ContinuationToken != nil {
				headers[headerIfNoneMatch] = string(*compositeToken.Continuation[0].ContinuationToken)
			}
			// Also set the feed range from the composite token if not already set
			if options.FeedRange == nil {
				options.FeedRange = &FeedRange{
					MinInclusive: compositeToken.Continuation[0].MinInclusive,
					MaxExclusive: compositeToken.Continuation[0].MaxExclusive,
				}
			}
		} else {
			// Not a composite token, treat as simple ETag
			headers[headerIfNoneMatch] = *options.Continuation
		}
	}

	if len(headers) == 0 {
		return nil
	}

	return &headers
}
