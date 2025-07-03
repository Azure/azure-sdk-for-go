// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azcosmos

import "fmt"

// FindPartitionKeyRangeIDForFeedRange finds the partition key range ID that matches the given FeedRange.
// Returns the ID if found, or an error if no match exists.
func FindPartitionKeyRangeId(feedRange FeedRange, partitionKeyRanges []PartitionKeyRangeProperties) (string, error) {
	for _, pkr := range partitionKeyRanges {
		if feedRange.MinInclusive == pkr.MinInclusive && feedRange.MaxExclusive == pkr.MaxExclusive {
			return pkr.ID, nil
		}
	}
	return "", fmt.Errorf("no matching partition key range found for feed range [%s, %s)", feedRange.MinInclusive, feedRange.MaxExclusive)
}

// TODO: Modify this function to use the partitionKeyRangeCache
// FindPartitionKeyRangeIdWithCache finds the partition key range ID for a FeedRange using a partitionKeyRangeCache.
// func FindPartitionKeyRangeIdWithCache(
//     feedRange FeedRange,
//     cache *partitionKeyRangeCache,
// ) (string, error) {
//     partitionKeyRanges := cache.getAll()
//     return FindPartitionKeyRangeId(feedRange, partitionKeyRanges)
// }
