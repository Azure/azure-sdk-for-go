// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !azcosmos_native_preview || !cgo || !((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFeedRangesRequireMatchingNativeExtension(t *testing.T) {
	value, err := NewFeedRange("", "FF")
	require.Zero(t, value)
	require.ErrorContains(t, err, "unpublished")
	value, err = ParseFeedRange("opaque")
	require.Zero(t, value)
	require.ErrorContains(t, err, "native v0.3.0")
	page, err := newTestContainer(t).GetFeedRanges(t.Context(), nil)
	require.Zero(t, page)
	require.ErrorContains(t, err, "unpublished")
	response, err := newTestContainer(t).GetFeedRangeForPartitionKey(t.Context(), NewPartitionKeyString("synthetic"))
	require.Zero(t, response)
	require.ErrorContains(t, err, "unpublished")
}
