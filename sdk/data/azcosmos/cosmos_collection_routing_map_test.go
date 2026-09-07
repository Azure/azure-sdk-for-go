// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newCollectionRoutingMap_basic(t *testing.T) {
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}

	rm := newCollectionRoutingMap(ranges, "etag1")
	require.NotNil(t, rm)
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "1", rm.orderedRanges[1].ID)
	require.Equal(t, "etag1", rm.changeFeedETag)
	require.False(t, rm.isGone("0"))
	require.False(t, rm.isGone("1"))
}

func Test_newCollectionRoutingMap_sortsRanges(t *testing.T) {
	// Provide ranges in reverse order
	ranges := []partitionKeyRange{
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
	}

	rm := newCollectionRoutingMap(ranges, "")
	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "1", rm.orderedRanges[1].ID)
}

func Test_newCollectionRoutingMap_filtersGoneParents(t *testing.T) {
	// Simulate a split: range "0" split into "2" and "3"
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
		{ID: "2", MinInclusive: "", MaxExclusive: "02E0", Parents: []string{"0"}},
		{ID: "3", MinInclusive: "02E0", MaxExclusive: "05C1CFFFFFFFF8", Parents: []string{"0"}},
	}

	rm := newCollectionRoutingMap(ranges, "etag2")
	require.Equal(t, 3, len(rm.orderedRanges))
	require.True(t, rm.isGone("0"))
	require.False(t, rm.isGone("1"))
	require.False(t, rm.isGone("2"))

	// Verify order
	require.Equal(t, "2", rm.orderedRanges[0].ID)
	require.Equal(t, "3", rm.orderedRanges[1].ID)
	require.Equal(t, "1", rm.orderedRanges[2].ID)
}

func Test_newCollectionRoutingMap_rangeByID(t *testing.T) {
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}

	rm := newCollectionRoutingMap(ranges, "")
	r, ok := rm.rangeByID["0"]
	require.True(t, ok)
	require.Equal(t, "", r.MinInclusive)
	require.Equal(t, "05C1CFFFFFFFF8", r.MaxExclusive)
}

func Test_newCollectionRoutingMap_deduplicatesRepeatedRangeID(t *testing.T) {
	// A change-feed drain accumulates every page, so a range that is touched while
	// a split is in progress is re-delivered on a later page. The routing map must
	// collapse the revisions instead of reporting that the range overlaps itself.
	// See https://github.com/Azure/azure-sdk-for-go/issues/27246
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}

	rm := newCollectionRoutingMap(ranges, "etag1")
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "1", rm.orderedRanges[1].ID)
	require.True(t, isCompleteSetOfRanges(rm.orderedRanges))
}

func Test_newCollectionRoutingMap_duplicateRangeID_lastRevisionWins(t *testing.T) {
	// Records arrive in change-feed order, so a later revision of a range ID
	// supersedes an earlier one. Matches the last-write-wins semantics of the
	// .NET, Java, Python and Rust SDKs.
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
		{ID: "0", MinInclusive: "", MaxExclusive: "0BC1"},
		{ID: "1", MinInclusive: "0BC1", MaxExclusive: "FF"},
	}

	rm := newCollectionRoutingMap(ranges, "etag1")
	require.Equal(t, 2, len(rm.orderedRanges))
	require.True(t, isCompleteSetOfRanges(rm.orderedRanges))

	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "0BC1", rm.orderedRanges[0].MaxExclusive, "later revision of range 0 should win")
	require.Equal(t, "1", rm.orderedRanges[1].ID)
	require.Equal(t, "0BC1", rm.orderedRanges[1].MinInclusive, "later revision of range 1 should win")

	// rangeByID must agree with orderedRanges
	require.Equal(t, "0BC1", rm.rangeByID["0"].MaxExclusive)
	require.Equal(t, "0BC1", rm.rangeByID["1"].MinInclusive)
}

func Test_newCollectionRoutingMap_duplicateRangeID_goneStillFiltered(t *testing.T) {
	// A split that completes mid-drain re-delivers the parent with an updated
	// status alongside its children, and unrelated ranges can be re-delivered in
	// the same drain. Every revision of the parent must be dropped, while the
	// live duplicate is collapsed to its latest revision.
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "3", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
		{ID: "1", MinInclusive: "", MaxExclusive: "02E0", Parents: []string{"0"}},
		{ID: "2", MinInclusive: "02E0", MaxExclusive: "05C1CFFFFFFFF8", Parents: []string{"0"}},
		// Parent re-delivered as Offline, and live range "3" re-delivered as Online.
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8", Status: "Offline"},
		{ID: "3", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF", Status: "Online"},
	}

	rm := newCollectionRoutingMap(ranges, "etag1")
	require.Equal(t, 3, len(rm.orderedRanges))
	require.Equal(t, "1", rm.orderedRanges[0].ID)
	require.Equal(t, "2", rm.orderedRanges[1].ID)
	require.Equal(t, "3", rm.orderedRanges[2].ID)
	// Both revisions of the gone parent are dropped...
	require.True(t, rm.isGone("0"))
	require.NotContains(t, rm.rangeByID, "0")
	// ...and the live duplicate collapses to its latest revision.
	require.Equal(t, "Online", rm.orderedRanges[2].Status)
	require.True(t, isCompleteSetOfRanges(rm.orderedRanges))
}

func Test_newCollectionRoutingMap_deterministicOrderingForInvalidSnapshot(t *testing.T) {
	// Two distinct IDs sharing a boundary can only appear in an invalid snapshot,
	// which is exactly when the discontinuity diagnostic ends up in a customer's
	// error message. The ordering must be reproducible so the diagnostic is too.
	//
	// The tied pair sits at index 7 of 13 records: that is large enough for
	// sort.Slice to leave its stable insertion-sort path, and this particular
	// arrangement is one where the unstable sort swaps the pair.
	ranges := []partitionKeyRange{
		{ID: "p00", MinInclusive: "04", MaxExclusive: "08"},
		{ID: "p01", MinInclusive: "08", MaxExclusive: "0C"},
		{ID: "p02", MinInclusive: "0C", MaxExclusive: "10"},
		{ID: "p03", MinInclusive: "10", MaxExclusive: "14"},
		{ID: "p04", MinInclusive: "14", MaxExclusive: "18"},
		{ID: "p05", MinInclusive: "18", MaxExclusive: "1C"},
		{ID: "p06", MinInclusive: "1C", MaxExclusive: "20"},
		{ID: "parent", MinInclusive: "", MaxExclusive: "FF"},
		{ID: "child", MinInclusive: "", MaxExclusive: "04"},
		{ID: "p07", MinInclusive: "20", MaxExclusive: "24"},
		{ID: "p08", MinInclusive: "24", MaxExclusive: "28"},
		{ID: "p09", MinInclusive: "28", MaxExclusive: "2C"},
		{ID: "p10", MinInclusive: "2C", MaxExclusive: "FF"},
	}

	first := describeRangeDiscontinuity(newCollectionRoutingMap(ranges, "").orderedRanges)
	require.NotEmpty(t, first)
	for i := 0; i < 50; i++ {
		rm := newCollectionRoutingMap(ranges, "")
		require.False(t, isCompleteSetOfRanges(rm.orderedRanges))
		// Change-feed arrival order is preserved for the tied boundary.
		require.Equal(t, "parent", rm.orderedRanges[0].ID)
		require.Equal(t, "child", rm.orderedRanges[1].ID)
		require.Equal(t, first, describeRangeDiscontinuity(rm.orderedRanges))
	}
}

func Test_tryCombine_duplicateRangeID_lastRevisionWins(t *testing.T) {
	// Guards the combined-slice construction: a range ID that appears both in the
	// existing map and repeatedly in the incremental batch must be appended once,
	// carrying its latest revision. Appending per occurrence would leave duplicate
	// entries that fail the continuity check and collapse the merge to nil.
	initial := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}, "etag1")

	// The same range ID is re-delivered across the accumulated incremental pages.
	newRanges := []partitionKeyRange{
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF", Status: "Splitting"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF", Status: "Online"},
	}

	merged := initial.tryCombine(newRanges, "etag2")
	require.NotNil(t, merged)
	require.Equal(t, 2, len(merged.orderedRanges))
	require.Equal(t, "0", merged.orderedRanges[0].ID)
	require.Equal(t, "1", merged.orderedRanges[1].ID)
	require.Equal(t, "Online", merged.orderedRanges[1].Status)
	require.Equal(t, len(merged.rangeByID), len(merged.orderedRanges))
}

func Test_tryCombine_successfulSplit(t *testing.T) {
	initial := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}, "etag1")

	// Split: range "0" splits into "2" and "3"
	newRanges := []partitionKeyRange{
		{ID: "2", MinInclusive: "", MaxExclusive: "02E0", Parents: []string{"0"}},
		{ID: "3", MinInclusive: "02E0", MaxExclusive: "05C1CFFFFFFFF8", Parents: []string{"0"}},
	}

	merged := initial.tryCombine(newRanges, "etag2")
	require.NotNil(t, merged)
	require.Equal(t, 3, len(merged.orderedRanges))
	require.Equal(t, "etag2", merged.changeFeedETag)
	require.True(t, merged.isGone("0"))

	// Verify ranges are sorted correctly
	require.Equal(t, "2", merged.orderedRanges[0].ID)
	require.Equal(t, "3", merged.orderedRanges[1].ID)
	require.Equal(t, "1", merged.orderedRanges[2].ID)
}

func Test_tryCombine_incompleteCovering(t *testing.T) {
	initial := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{ID: "1", MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}, "etag1")

	// Provide only one child — covering is incomplete
	newRanges := []partitionKeyRange{
		{ID: "2", MinInclusive: "", MaxExclusive: "02E0", Parents: []string{"0"}},
	}

	merged := initial.tryCombine(newRanges, "etag2")
	require.Nil(t, merged, "tryCombine should return nil for incomplete covering")
}

func Test_isCompleteSetOfRanges_valid(t *testing.T) {
	ranges := []partitionKeyRange{
		{MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
		{MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}
	require.True(t, isCompleteSetOfRanges(ranges))
}

func Test_isCompleteSetOfRanges_empty(t *testing.T) {
	require.False(t, isCompleteSetOfRanges(nil))
	require.False(t, isCompleteSetOfRanges([]partitionKeyRange{}))
}

func Test_isCompleteSetOfRanges_doesNotStartAtEmpty(t *testing.T) {
	ranges := []partitionKeyRange{
		{MinInclusive: "05C1CFFFFFFFF8", MaxExclusive: "FF"},
	}
	require.False(t, isCompleteSetOfRanges(ranges))
}

func Test_isCompleteSetOfRanges_gap(t *testing.T) {
	ranges := []partitionKeyRange{
		{MinInclusive: "", MaxExclusive: "03"},
		{MinInclusive: "05", MaxExclusive: "FF"}, // gap between 03 and 05
	}
	require.False(t, isCompleteSetOfRanges(ranges))
}

func Test_isCompleteSetOfRanges_doesNotEndAtFF(t *testing.T) {
	ranges := []partitionKeyRange{
		{MinInclusive: "", MaxExclusive: "05C1CFFFFFFFF8"},
	}
	require.False(t, isCompleteSetOfRanges(ranges))
}

func Test_isCompleteSetOfRanges_singleRange(t *testing.T) {
	ranges := []partitionKeyRange{
		{MinInclusive: "", MaxExclusive: "FF"},
	}
	require.True(t, isCompleteSetOfRanges(ranges))
}

func Test_isCompleteSetOfRanges_emptyMaxExclusive(t *testing.T) {
	// Some implementations use "" as unbounded end
	ranges := []partitionKeyRange{
		{MinInclusive: "", MaxExclusive: ""},
	}
	require.True(t, isCompleteSetOfRanges(ranges))
}

func Test_getOverlappingRanges_singleRange(t *testing.T) {
	rm := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "")

	result := rm.getOverlappingRanges("0000", "3FFF")
	require.Len(t, result, 1)
	require.Equal(t, "0", result[0].ID)
}

func Test_getOverlappingRanges_multipleRanges(t *testing.T) {
	rm := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "10"},
		{ID: "1", MinInclusive: "10", MaxExclusive: "20"},
		{ID: "2", MinInclusive: "20", MaxExclusive: "30"},
		{ID: "3", MinInclusive: "30", MaxExclusive: "FF"},
	}, "")

	// Range spanning partitions 1 and 2
	result := rm.getOverlappingRanges("10", "30")
	require.Len(t, result, 2)
	require.Equal(t, "1", result[0].ID)
	require.Equal(t, "2", result[1].ID)
}

func Test_getOverlappingRanges_allRanges(t *testing.T) {
	rm := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "10"},
		{ID: "1", MinInclusive: "10", MaxExclusive: "20"},
		{ID: "2", MinInclusive: "20", MaxExclusive: "FF"},
	}, "")

	result := rm.getOverlappingRanges("", "FF")
	require.Len(t, result, 3)
}

func Test_getOverlappingRanges_pointInMiddle(t *testing.T) {
	rm := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "10"},
		{ID: "1", MinInclusive: "10", MaxExclusive: "20"},
		{ID: "2", MinInclusive: "20", MaxExclusive: "FF"},
	}, "")

	// EPK range that starts and ends within range "1"
	result := rm.getOverlappingRanges("15", "18")
	require.Len(t, result, 1)
	require.Equal(t, "1", result[0].ID)
}

func Test_getOverlappingRanges_mixedLengthBoundaries(t *testing.T) {
	// Simulate HPK container with mixed-length EPK boundaries
	partial := "06AB34CFE4E482236BCACBBF50E234AB"
	fullZero := partial + "00000000000000000000000000000000"

	rm := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: partial},
		{ID: "1", MinInclusive: fullZero, MaxExclusive: "FF"},
	}, "")

	// A query range spanning both should find both
	result := rm.getOverlappingRanges("", "FF")
	require.Len(t, result, 2)
}
