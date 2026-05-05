// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"sort"
)

// collectionRoutingMap holds an immutable snapshot of partition key ranges for a
// container, sorted for efficient EPK lookups. It supports incremental merging
// when partition splits or merges occur.
type collectionRoutingMap struct {
	// orderedRanges are the partition key ranges sorted by MinInclusive ascending.
	orderedRanges []partitionKeyRange
	// rangeByID provides O(1) lookups of ranges by their ID.
	rangeByID map[string]partitionKeyRange
	// goneRanges tracks parent range IDs that have been replaced by children after splits.
	goneRanges map[string]bool
	// changeFeedETag is the ETag for incremental change-feed refreshes.
	changeFeedETag string
}

// newCollectionRoutingMap creates a new collectionRoutingMap from a set of ranges.
// It filters out "gone" parent ranges (identified via the Parents field on child ranges)
// and sorts the remaining ranges by MinInclusive.
func newCollectionRoutingMap(ranges []partitionKeyRange, changeFeedETag string) *collectionRoutingMap {
	goneRanges := make(map[string]bool)
	for _, r := range ranges {
		for _, parent := range r.Parents {
			goneRanges[parent] = true
		}
	}

	// Filter out gone ranges
	filtered := make([]partitionKeyRange, 0, len(ranges))
	for _, r := range ranges {
		if !goneRanges[r.ID] {
			filtered = append(filtered, r)
		}
	}

	// Sort by MinInclusive
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].MinInclusive < filtered[j].MinInclusive
	})

	rangeByID := make(map[string]partitionKeyRange, len(filtered))
	for _, r := range filtered {
		rangeByID[r.ID] = r
	}

	return &collectionRoutingMap{
		orderedRanges:  filtered,
		rangeByID:      rangeByID,
		goneRanges:     goneRanges,
		changeFeedETag: changeFeedETag,
	}
}

// tryCombine merges new ranges (from an incremental change-feed refresh) into
// the existing routing map. Returns a new collectionRoutingMap if the merge
// succeeds (produces a complete covering), or nil if the result is incomplete
// (indicating a full refresh is needed).
func (crm *collectionRoutingMap) tryCombine(newRanges []partitionKeyRange, newETag string) *collectionRoutingMap {
	// Accumulate gone ranges from both existing and new ranges
	combinedGone := make(map[string]bool, len(crm.goneRanges))
	for id := range crm.goneRanges {
		combinedGone[id] = true
	}
	for _, r := range newRanges {
		for _, parent := range r.Parents {
			combinedGone[parent] = true
		}
	}

	// Build a combined set: existing ranges (minus gone) plus new ranges (minus gone)
	combinedByID := make(map[string]partitionKeyRange, len(crm.rangeByID)+len(newRanges))
	for id, r := range crm.rangeByID {
		if !combinedGone[id] {
			combinedByID[id] = r
		}
	}
	for _, r := range newRanges {
		if !combinedGone[r.ID] {
			combinedByID[r.ID] = r
		}
	}

	// Build sorted slice
	combined := make([]partitionKeyRange, 0, len(combinedByID))
	for _, r := range combinedByID {
		combined = append(combined, r)
	}
	sort.Slice(combined, func(i, j int) bool {
		return combined[i].MinInclusive < combined[j].MinInclusive
	})

	// Validate completeness: ranges must form a contiguous covering
	if !isCompleteSetOfRanges(combined) {
		return nil
	}

	return &collectionRoutingMap{
		orderedRanges:  combined,
		rangeByID:      combinedByID,
		goneRanges:     combinedGone,
		changeFeedETag: newETag,
	}
}

// isGone returns true if the given range ID has been replaced (by a split/merge).
func (crm *collectionRoutingMap) isGone(rangeID string) bool {
	return crm.goneRanges[rangeID]
}

// isCompleteSetOfRanges validates that the sorted ranges form a contiguous
// partition covering with no gaps or overlaps. The first range should start
// at "" and each subsequent range should start where the previous one ends.
func isCompleteSetOfRanges(ranges []partitionKeyRange) bool {
	if len(ranges) == 0 {
		return false
	}

	// First range must start at ""
	if ranges[0].MinInclusive != "" {
		return false
	}

	// Each range's MinInclusive must equal the previous range's MaxExclusive
	for i := 1; i < len(ranges); i++ {
		if ranges[i].MinInclusive != ranges[i-1].MaxExclusive {
			return false
		}
	}

	// Last range must end at "FF" (the maximum EPK boundary) or be unbounded ("")
	lastMax := ranges[len(ranges)-1].MaxExclusive
	if lastMax != "FF" && lastMax != "" {
		return false
	}

	return true
}
