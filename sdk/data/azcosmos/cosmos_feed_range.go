// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "fmt"

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

// findPartitionKeyRangeID finds the partition key range ID that matches the given FeedRange.
// Returns the ID if found, or an error if no match exists.
func findPartitionKeyRangeID(feedRange FeedRange, partitionKeyRanges []partitionKeyRange) (string, error) {
	for _, pkr := range partitionKeyRanges {
		if feedRange.MinInclusive == pkr.MinInclusive && feedRange.MaxExclusive == pkr.MaxExclusive {
			return pkr.ID, nil
		}
	}
	return "", fmt.Errorf("no matching partition key range found for feed range [%s, %s)", feedRange.MinInclusive, feedRange.MaxExclusive)
}
