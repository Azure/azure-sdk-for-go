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

func Test_tryCombine_successfulSplit(t *testing.T) {
	// Initial state: two ranges
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
