// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// changeFeedRange represents a range of partition key values for a Cosmos container's change feed.
// It is used to identify a specific range of documents for change feed processing.
type changeFeedRange struct {
	// MinInclusive contains the minimum inclusive value of the partition key range.
	MinInclusive string `json:"minInclusive"`
	// MaxExclusive contains the maximum exclusive value of the partition key range.
	MaxExclusive string `json:"maxExclusive"`
	// ContinuationToken is used to continue reading the change feed from a specific point.
	ContinuationToken *azcore.ETag `json:"continuationToken,omitempty"`
	// epkMinHeader is the header for the minimum inclusive value of the partition key range.
	// This is used internally to set the headers for change feed requests.
	epkMinHeader string `json:"-"`
	// epkMaxHeader is the header for the maximum exclusive value of the partition key range.
	// This is used internally to set the headers for change feed requests.
	epkMaxHeader string `json:"-"`
}

// ChangeFeedRangeOptions includes options for creating a new change feed range.
type ChangeFeedRangeOptions struct {
	// ContinuationToken is used to continue reading the change feed from a specific point.
	ContinuationToken *azcore.ETag
	// EpkMinHeader is the header for the minimum inclusive value of the partition key range.
	EpkMinHeader *string
	// EpkMaxHeader is the header for the maximum exclusive value of the partition key range.
	EpkMaxHeader *string
}

// newChangeFeedRange creates a new changeFeedRange with the specified minimum inclusive and maximum exclusive values.
// Acts as a FeedRange for which change feed is being requested.
// Designed for internal use only for creating change feed ranges.
func newChangeFeedRange(minInclusive, maxExclusive string, options *ChangeFeedRangeOptions) changeFeedRange {
	result := changeFeedRange{
		MinInclusive: minInclusive,
		MaxExclusive: maxExclusive,
	}

	if options != nil {
		if options.ContinuationToken != nil {
			continuationETag := *options.ContinuationToken
			result.ContinuationToken = &continuationETag
		}
		if options.EpkMinHeader != nil {
			result.epkMinHeader = *options.EpkMinHeader
		}
		if options.EpkMaxHeader != nil {
			result.epkMaxHeader = *options.EpkMaxHeader
		}
	}

	return result
}
