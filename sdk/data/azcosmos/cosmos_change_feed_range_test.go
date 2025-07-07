// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// newChangeFeedRange creates a new ChangeFeedRange with the specified minimum inclusive and maximum exclusive values.
func TestNewChangeFeedRange_Basic(t *testing.T) {
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

// TestNewChangeFeedRange_NilContinuationToken tests the case where we
// will only use the min and max values, without any continuation token or headers.
func TestNewChangeFeedRange_NilOptions(t *testing.T) {
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
