// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestNewChangeFeedRangeBasic(t *testing.T) {
	min := ""
	max := "1FFFFFF"
	token := azcore.ETag("15")
	minHeader := "minHeaderValue"
	maxHeader := "maxHeaderValue"

	options := &ChangeFeedRangeOptions{
		ContinuationToken: &token,
		EpkMinHeader:      &minHeader,
		EpkMaxHeader:      &maxHeader,
	}

	cfr := newChangeFeedRange(min, max, options)

	if cfr.MinInclusive != min {
		t.Errorf("MinInclusive mismatch: got %s, want %s", cfr.MinInclusive, min)
	}
	if cfr.MaxExclusive != max {
		t.Errorf("MaxExclusive mismatch: got %s, want %s", cfr.MaxExclusive, max)
	}
	if cfr.ContinuationToken == nil || *cfr.ContinuationToken != token {
		t.Errorf("ContinuationToken mismatch: got %v, want %s", cfr.ContinuationToken, token)
	}
	if cfr.epkMinHeader != minHeader {
		t.Errorf("epkMinHeader mismatch: got %s, want %s", cfr.epkMinHeader, minHeader)
	}
	if cfr.epkMaxHeader != maxHeader {
		t.Errorf("epkMaxHeader mismatch: got %s, want %s", cfr.epkMaxHeader, maxHeader)
	}
}

func TestOverlappingChangeFeedRanges(t *testing.T) {
	// Arrange
	range1 := newChangeFeedRange("00", "80", nil)
	range2 := newChangeFeedRange("40", "FF", nil)

	// Assert - verify ranges overlap
	if !(range1.MinInclusive < range2.MaxExclusive) {
		t.Errorf("Range1.MinInclusive (%s) should be less than Range2.MaxExclusive (%s)", range1.MinInclusive, range2.MaxExclusive)
	}
	if !(range2.MinInclusive < range1.MaxExclusive) {
		t.Errorf("Range2.MinInclusive (%s) should be less than Range1.MaxExclusive (%s)", range2.MinInclusive, range1.MaxExclusive)
	}

	// Assert - verify the full range they cover
	if range1.MinInclusive != "00" {
		t.Errorf("Expected Range1.MinInclusive to be '00', got %s", range1.MinInclusive)
	}
	if range2.MaxExclusive != "FF" {
		t.Errorf("Expected Range2.MaxExclusive to be 'FF', got %s", range2.MaxExclusive)
	}
}

func TestNewChangeFeedRangeNilOptions(t *testing.T) {
	min := "A"
	max := "Z"
	cfr := newChangeFeedRange(min, max, nil)

	if cfr.MinInclusive != min {
		t.Errorf("MinInclusive mismatch: got %s, want %s", cfr.MinInclusive, min)
	}
	if cfr.MaxExclusive != max {
		t.Errorf("MaxExclusive mismatch: got %s, want %s", cfr.MaxExclusive, max)
	}
	if cfr.ContinuationToken != nil {
		t.Errorf("Expected nil ContinuationToken, got %v", cfr.ContinuationToken)
	}
	if cfr.epkMinHeader != "" {
		t.Errorf("Expected empty epkMinHeader, got %s", cfr.epkMinHeader)
	}
	if cfr.epkMaxHeader != "" {
		t.Errorf("Expected empty epkMaxHeader, got %s", cfr.epkMaxHeader)
	}
}
