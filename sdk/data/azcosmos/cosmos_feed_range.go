// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// FeedRange represents a range of partition key values for a Cosmos container.
// It is used to identify a specific range of documents for change feed processing.
type FeedRange struct {
	// MinInclusive contains the minimum inclusive value of the partition key range.
	MinInclusive string
	// MaxExclusive contains the maximum exclusive value of the partition key range.
	MaxExclusive string
}

// NewFeedRange creates a new FeedRange with the specified minimum inclusive and maximum exclusive values.
func NewFeedRange(minInclusive, maxExclusive string) FeedRange {
	return FeedRange{
		MinInclusive: minInclusive,
		MaxExclusive: maxExclusive,
	}
}
