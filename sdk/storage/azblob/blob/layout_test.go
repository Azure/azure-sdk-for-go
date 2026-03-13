// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetIdealEndpoint_EmptyLayoutRanges(t *testing.T) {
	l := layout{
		layoutRanges:  []layoutRange{},
		contentLength: 100,
	}
	result := getIdealEndpoint(50, l)
	require.Equal(t, "", result)
}

func TestGetIdealEndpoint_SingleRange(t *testing.T) {
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 100, endpoint: "endpoint1"},
		},
		contentLength: 100,
	}

	// Offset at start
	require.Equal(t, "endpoint1", getIdealEndpoint(0, l))

	// Offset in middle
	require.Equal(t, "endpoint1", getIdealEndpoint(50, l))

	// Offset at end
	require.Equal(t, "endpoint1", getIdealEndpoint(100, l))
}

func TestGetIdealEndpoint_MultipleRanges(t *testing.T) {
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 99, endpoint: "endpoint1"},
			{start: 100, end: 199, endpoint: "endpoint2"},
			{start: 200, end: 299, endpoint: "endpoint3"},
		},
		contentLength: 300,
	}

	// Offset in first range
	require.Equal(t, "endpoint1", getIdealEndpoint(0, l))
	require.Equal(t, "endpoint1", getIdealEndpoint(50, l))
	require.Equal(t, "endpoint1", getIdealEndpoint(99, l))

	// Offset in second range
	require.Equal(t, "endpoint2", getIdealEndpoint(100, l))
	require.Equal(t, "endpoint2", getIdealEndpoint(150, l))
	require.Equal(t, "endpoint2", getIdealEndpoint(199, l))

	// Offset in third range
	require.Equal(t, "endpoint3", getIdealEndpoint(200, l))
	require.Equal(t, "endpoint3", getIdealEndpoint(250, l))
	require.Equal(t, "endpoint3", getIdealEndpoint(299, l))
}

func TestGetIdealEndpoint_BinarySearchBoundary(t *testing.T) {
	// Test with more ranges to exercise binary search properly
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 9, endpoint: "ep0"},
			{start: 10, end: 19, endpoint: "ep1"},
			{start: 20, end: 29, endpoint: "ep2"},
			{start: 30, end: 39, endpoint: "ep3"},
			{start: 40, end: 49, endpoint: "ep4"},
			{start: 50, end: 59, endpoint: "ep5"},
			{start: 60, end: 69, endpoint: "ep6"},
		},
		contentLength: 70,
	}

	// Test boundaries at each range
	require.Equal(t, "ep0", getIdealEndpoint(0, l))
	require.Equal(t, "ep0", getIdealEndpoint(9, l))
	require.Equal(t, "ep1", getIdealEndpoint(10, l))
	require.Equal(t, "ep3", getIdealEndpoint(35, l))
	require.Equal(t, "ep6", getIdealEndpoint(65, l))
	require.Equal(t, "ep6", getIdealEndpoint(69, l))
}

func TestGetIdealEndpoint_SameEndpointDifferentRanges(t *testing.T) {
	l := layout{
		layoutRanges: []layoutRange{
			{start: 0, end: 49, endpoint: "endpointA"},
			{start: 50, end: 99, endpoint: "endpointB"},
			{start: 100, end: 149, endpoint: "endpointA"},
		},
		contentLength: 150,
	}

	require.Equal(t, "endpointA", getIdealEndpoint(25, l))
	require.Equal(t, "endpointB", getIdealEndpoint(75, l))
	require.Equal(t, "endpointA", getIdealEndpoint(125, l))
}
