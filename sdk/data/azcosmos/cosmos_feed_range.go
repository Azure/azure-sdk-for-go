// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/epk"
)

// ErrFeedRangeUnresolved is returned when a customer-supplied FeedRange does not
// overlap any current physical partition key range, even after a forced cache
// refresh. Callers can use errors.Is to detect this and fall back to
// re-deriving FeedRanges from GetFeedRanges.
//
// This typically indicates one of:
//   - The FeedRange was constructed for a different container.
//   - The container was deleted and recreated with different partitioning.
//   - The FeedRange boundaries are malformed (outside [00, FF) or empty range).
var ErrFeedRangeUnresolved = errors.New("feed range did not overlap any current partition key range")

// feedRangeUnresolvedError wraps ErrFeedRangeUnresolved with diagnostic detail
// about the unresolvable FeedRange so customers can identify which range failed.
type feedRangeUnresolvedError struct {
	feedRange FeedRange
}

func (e *feedRangeUnresolvedError) Error() string {
	return fmt.Sprintf("%s: [%s, %s)", ErrFeedRangeUnresolved.Error(), e.feedRange.MinInclusive, e.feedRange.MaxExclusive)
}

func (e *feedRangeUnresolvedError) Unwrap() error {
	return ErrFeedRangeUnresolved
}

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

// findOverlappingPartitionKeyRangeIDs returns the IDs of every PK range that
// overlaps the given FeedRange. Used by GetChangeFeed to expand a parent feed
// range into one entry per child after a split (or fold a child into a parent
// after a merge). The returned slice preserves the routing-map order.
//
// Empty feed-range boundaries follow Cosmos convention: "" on MinInclusive
// means "open at the bottom" (lowest possible key), "" on MaxExclusive means
// "open at the top" (highest possible key — normalized to "FF").
//
// Returns ErrFeedRangeUnresolved (wrapped) when the input doesn't overlap
// any range — typically a malformed range or a wrong-container token.
func findOverlappingPartitionKeyRangeIDs(feedRange FeedRange, partitionKeyRanges []partitionKeyRange) ([]string, error) {
	if len(partitionKeyRanges) == 0 {
		return nil, &feedRangeUnresolvedError{feedRange: feedRange}
	}

	feedMin := feedRange.MinInclusive
	feedMax := normalizeMaxBoundary(feedRange.MaxExclusive)

	// Sanity: feedMin must be < feedMax. An equal/inverted range can never overlap any
	// well-formed routing map and would silently match nothing — surface as Unresolved.
	if epk.CompareEPK(feedMin, feedMax) >= 0 {
		return nil, &feedRangeUnresolvedError{feedRange: feedRange}
	}

	ids := make([]string, 0, 2)
	for _, pkr := range partitionKeyRanges {
		if rangesOverlap(feedMin, feedMax, pkr.MinInclusive, normalizeMaxBoundary(pkr.MaxExclusive)) {
			ids = append(ids, pkr.ID)
		}
	}
	if len(ids) == 0 {
		return nil, &feedRangeUnresolvedError{feedRange: feedRange}
	}
	return ids, nil
}

// overlappingPartitionKeyRanges returns the subset of partitionKeyRanges whose
// boundaries overlap the given feedRange, preserving input order. Single-pass
// equivalent of findOverlappingPartitionKeyRangeIDs that returns the ranges
// themselves rather than IDs. Returns nil on no overlap (no error).
//
// Note: O(n) linear scan. For routing-map-backed lookups, prefer
// collectionRoutingMap.getOverlappingRanges (binary search). This helper exists
// for paths that operate on a flat snapshot returned by getPartitionKeyRanges.
func overlappingPartitionKeyRanges(feedRange FeedRange, partitionKeyRanges []partitionKeyRange) []partitionKeyRange {
	if len(partitionKeyRanges) == 0 {
		return nil
	}

	feedMin := feedRange.MinInclusive
	feedMax := normalizeMaxBoundary(feedRange.MaxExclusive)
	if epk.CompareEPK(feedMin, feedMax) >= 0 {
		return nil
	}

	out := make([]partitionKeyRange, 0, 2)
	for _, pkr := range partitionKeyRanges {
		if rangesOverlap(feedMin, feedMax, pkr.MinInclusive, normalizeMaxBoundary(pkr.MaxExclusive)) {
			out = append(out, pkr)
		}
	}
	return out
}

// rangesOverlap reports whether [aMin, aMax) intersects [bMin, bMax).
// All four boundaries must be normalized hex EPK strings (with "FF" used
// for the upper sentinel rather than ""); rangesOverlap does NOT treat
// empty strings as open boundaries.
func rangesOverlap(aMin, aMax, bMin, bMax string) bool {
	// Standard half-open interval overlap test:
	//   intersection is empty iff aMax <= bMin or bMax <= aMin.
	if epk.CompareEPK(aMax, bMin) <= 0 {
		return false
	}
	if epk.CompareEPK(bMax, aMin) <= 0 {
		return false
	}
	return true
}

// normalizeMaxBoundary converts the open-top "" sentinel to "FF" so length-aware
// EPK comparisons against finite ranges produce sensible answers. PK range
// snapshots from the service usually use "FF" but customer-supplied FeedRanges
// often use "". Min boundaries are NOT normalized — "" already sorts lowest.
func normalizeMaxBoundary(maxExclusive string) string {
	if maxExclusive == "" {
		return "FF"
	}
	return maxExclusive
}
