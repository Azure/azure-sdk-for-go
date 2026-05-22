// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestContainerGetFeedRanges(t *testing.T) {
	expectedJsonResponse := []byte(`{
	"_rid": "lypXAMSZ-Cs=",
	"PartitionKeyRanges": [
        {
            "_rid": "lypXAMSZ-CuZAAAAAAAAUA==",
            "id": "151",
            "_etag": "\"0000cc70-0000-0100-0000-682306240000\"",
            "minInclusive": "05C1E18D2D7F08",
            "maxExclusive": "05C1E18D2D83FA",
            "ridPrefix": 151,
            "_self": "dbs/lypXAA==/colls/lypXAMSZ-Cs=/pkranges/lypXAMSZ-CuZAAAAAAAAUA==/",
            "throughputFraction": 0.0125,
            "status": "online",
            "parents": [
                "5",
                "10",
                "31"
            ],
            "ownedArchivalPKRangeIds": [
                "31"
            ],
            "_ts": 1747125796,
            "lsn": 22874
        },
        {
            "_rid": "lypXAMSZ-CulAAAAAAAAUA==",
            "id": "163",
            "_etag": "\"0000dd1b-0000-0100-0000-67f6d6a70000\"",
            "minInclusive": "05C1C7FF3903F8",
            "maxExclusive": "05C1C9CD673390",
            "ridPrefix": 163,
            "_self": "dbs/lypXAA==/colls/lypXAMSZ-Cs=/pkranges/lypXAMSZ-CulAAAAAAAAUA==/",
            "throughputFraction": 0.0125,
            "status": "online",
            "parents": [
                "1",
                "19",
                "39"
            ],
            "ownedArchivalPKRangeIds": [
                "39"
            ],
            "_ts": 1744230055,
            "lsn": 22599
        }
	],
	"_count": 2
	}`)

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(expectedJsonResponse),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200),
	)

	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	feedRanges, err := container.GetFeedRanges(context.TODO())
	if err != nil {
		t.Fatalf("GetFeedRanges failed: %v", err)
	}

	if len(feedRanges) != 2 {
		t.Fatalf("Expected 2 feed ranges, got %d", len(feedRanges))
	}

	if feedRanges[0].MinInclusive != "05C1E18D2D7F08" {
		t.Errorf("Expected MinInclusive to be 05C1E18D2D7F08, got %s", feedRanges[0].MinInclusive)
	}

	if feedRanges[0].MaxExclusive != "05C1E18D2D83FA" {
		t.Errorf("Expected MaxExclusive to be 05C1E18D2D83FA, got %s", feedRanges[0].MaxExclusive)
	}

	if feedRanges[1].MinInclusive != "05C1C7FF3903F8" {
		t.Errorf("Expected MinInclusive to be 05C1C7FF3903F8, got %s", feedRanges[1].MinInclusive)
	}

	if feedRanges[1].MaxExclusive != "05C1C9CD673390" {
		t.Errorf("Expected MaxExclusive to be 05C1C9CD673390, got %s", feedRanges[1].MaxExclusive)
	}
}

func TestContainerGetFeedRangesEmpty(t *testing.T) {
	expectedJsonResponse := `{
    "_rid": "lypXAMSZ-Cs=",
    "PartitionKeyRanges": [],
    "_count": 0
	}`

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody([]byte(expectedJsonResponse)),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200),
	)

	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	feedRanges, err := container.GetFeedRanges(context.TODO())
	if err != nil {
		t.Fatalf("GetFeedRanges failed: %v", err)
	}

	if len(feedRanges) != 0 {
		t.Fatalf("Expected 0 feed ranges, got %d", len(feedRanges))
	}
}

func TestContainerGetFeedRanges_UsesCache(t *testing.T) {
	containerResponse := []byte(`{
		"id": "containerId",
		"_rid": "testRID",
		"_self": "dbs/db1/colls/containerId/",
		"partitionKey": {
			"paths": ["/pk"],
			"kind": "Hash",
			"version": 2
		}
	}`)

	pkRangeResponse := []byte(`{
		"_rid": "testRID",
		"PartitionKeyRanges": [
			{
				"_rid": "testRID_range0",
				"id": "0",
				"_etag": "\"etag0\"",
				"minInclusive": "",
				"maxExclusive": "05C1E18D2D7F08",
				"status": "online",
				"parents": []
			},
			{
				"_rid": "testRID_range1",
				"id": "1",
				"_etag": "\"etag1\"",
				"minInclusive": "05C1E18D2D7F08",
				"maxExclusive": "FF",
				"status": "online",
				"parents": []
			}
		],
		"_count": 2
	}`)

	srv, close := mock.NewTLSServer()
	defer close()

	// First call will need: container read (for RID) + PK range fetch
	srv.AppendResponse(
		mock.WithBody(containerResponse),
		mock.WithStatusCode(200),
	)
	srv.AppendResponse(
		mock.WithBody(pkRangeResponse),
		mock.WithHeader(cosmosHeaderEtag, "changeFeedEtag1"),
		mock.WithStatusCode(200),
	)
	// 304 Not Modified — terminates the change-feed loop
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "changeFeedEtag1"),
	)

	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{
		endpoint:    srv.URL(),
		endpointUrl: defaultEndpoint,
		internal:    internalClient,
		gem:         gem,
		caches: &sharedCacheSet{
			pkRangeCache:   newPartitionKeyRangeCache(),
			containerCache: newContainerPropertiesCache(),
		},
	}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	// First call: populates caches, makes 2 HTTP requests
	feedRanges, err := container.GetFeedRanges(context.TODO())
	require.NoError(t, err)
	require.Equal(t, 2, len(feedRanges))
	require.Equal(t, "", feedRanges[0].MinInclusive)
	require.Equal(t, "05C1E18D2D7F08", feedRanges[0].MaxExclusive)
	require.Equal(t, "05C1E18D2D7F08", feedRanges[1].MinInclusive)
	require.Equal(t, "FF", feedRanges[1].MaxExclusive)

	requestsAfterFirstCall := srv.Requests()
	require.Equal(t, 3, requestsAfterFirstCall, "first call should make 3 HTTP requests (container read + PK ranges + 304)")

	// Second call: should use caches, no additional HTTP requests
	// (no more responses queued — would panic if a request was made)
	feedRanges2, err := container.GetFeedRanges(context.TODO())
	require.NoError(t, err)
	require.Equal(t, 2, len(feedRanges2))
	require.Equal(t, feedRanges[0], feedRanges2[0])
	require.Equal(t, feedRanges[1], feedRanges2[1])

	require.Equal(t, requestsAfterFirstCall, srv.Requests(), "second call should make 0 HTTP requests (cache hit)")
}

// TestNormalizeMaxBoundary verifies the open-top "" sentinel becomes "FF",
// matching the Cosmos convention used everywhere routing math is performed.
// Min boundaries are NOT normalized — "" already sorts lowest.
func TestNormalizeMaxBoundary(t *testing.T) {
	require.Equal(t, "FF", normalizeMaxBoundary(""), "empty Max must normalize to FF")
	require.Equal(t, "FF", normalizeMaxBoundary("FF"), "explicit FF must round-trip")
	require.Equal(t, "80", normalizeMaxBoundary("80"), "non-empty Max must pass through unchanged")
	require.Equal(t, "00", normalizeMaxBoundary("00"))
	require.Equal(t, "AABBCCDD", normalizeMaxBoundary("AABBCCDD"), "long EPK must pass through unchanged")
}

// TestRangesOverlap_TableDriven exercises the half-open interval overlap
// predicate. All four boundaries must be normalized hex EPK strings.
func TestRangesOverlap_TableDriven(t *testing.T) {
	cases := []struct {
		name                   string
		aMin, aMax, bMin, bMax string
		want                   bool
	}{
		{"identical", "10", "20", "10", "20", true},
		{"a-contains-b", "00", "FF", "20", "30", true},
		{"b-contains-a", "20", "30", "00", "FF", true},
		{"a-overlaps-b-left", "00", "30", "20", "FF", true},
		{"a-overlaps-b-right", "20", "FF", "00", "30", true},
		{"adjacent-a-then-b", "00", "20", "20", "FF", false},
		{"adjacent-b-then-a", "20", "FF", "00", "20", false},
		{"disjoint-a-below", "00", "10", "20", "30", false},
		{"disjoint-a-above", "30", "FF", "00", "20", false},
		{"single-byte-overlap", "10", "11", "10", "FF", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, rangesOverlap(c.aMin, c.aMax, c.bMin, c.bMax))
		})
	}
}

// TestOverlappingPartitionKeyRanges_TableDriven covers the boundary
// behaviors of overlap-resolution against a snapshot routing map. This is
// the single most critical helper in the F1 fix — it replaces the exact-
// match findPartitionKeyRangeID that was the original-bug source.
func TestOverlappingPartitionKeyRanges_TableDriven(t *testing.T) {
	threeRanges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "55"},
		{ID: "1", MinInclusive: "55", MaxExclusive: "AA"},
		{ID: "2", MinInclusive: "AA", MaxExclusive: "FF"},
	}

	cases := []struct {
		name    string
		feed    FeedRange
		ranges  []partitionKeyRange
		wantIDs []string
	}{
		{
			name:    "exact-match-single-range",
			feed:    FeedRange{MinInclusive: "55", MaxExclusive: "AA"},
			ranges:  threeRanges,
			wantIDs: []string{"1"},
		},
		{
			name:    "spans-all-three-ranges",
			feed:    FeedRange{MinInclusive: "", MaxExclusive: "FF"},
			ranges:  threeRanges,
			wantIDs: []string{"0", "1", "2"},
		},
		{
			name:    "spans-all-three-with-empty-max",
			feed:    FeedRange{MinInclusive: "", MaxExclusive: ""},
			ranges:  threeRanges,
			wantIDs: []string{"0", "1", "2"},
		},
		{
			name:    "strict-sub-range-of-one-physical",
			feed:    FeedRange{MinInclusive: "20", MaxExclusive: "40"},
			ranges:  threeRanges,
			wantIDs: []string{"0"},
		},
		{
			name:    "straddles-two-ranges",
			feed:    FeedRange{MinInclusive: "30", MaxExclusive: "70"},
			ranges:  threeRanges,
			wantIDs: []string{"0", "1"},
		},
		{
			name:    "preserves-input-order",
			feed:    FeedRange{MinInclusive: "10", MaxExclusive: "C0"},
			ranges:  threeRanges,
			wantIDs: []string{"0", "1", "2"},
		},
		{
			name:    "empty-snapshot-returns-nil",
			feed:    FeedRange{MinInclusive: "", MaxExclusive: "FF"},
			ranges:  nil,
			wantIDs: nil,
		},
		{
			name:    "inverted-feed-range-returns-nil",
			feed:    FeedRange{MinInclusive: "FF", MaxExclusive: "00"},
			ranges:  threeRanges,
			wantIDs: nil,
		},
		{
			name:    "equal-bounds-returns-nil",
			feed:    FeedRange{MinInclusive: "55", MaxExclusive: "55"},
			ranges:  threeRanges,
			wantIDs: nil,
		},
		{
			name: "no-overlap-with-non-empty-snapshot",
			feed: FeedRange{MinInclusive: "10", MaxExclusive: "20"},
			ranges: []partitionKeyRange{
				{ID: "x", MinInclusive: "30", MaxExclusive: "40"},
				{ID: "y", MinInclusive: "60", MaxExclusive: "70"},
			},
			wantIDs: nil,
		},
		{
			name: "boundary-touch-is-not-overlap",
			feed: FeedRange{MinInclusive: "55", MaxExclusive: "AA"},
			ranges: []partitionKeyRange{
				{ID: "left", MinInclusive: "00", MaxExclusive: "55"}, // touches but doesn't overlap
				{ID: "match", MinInclusive: "55", MaxExclusive: "AA"},
				{ID: "right", MinInclusive: "AA", MaxExclusive: "FF"}, // touches but doesn't overlap
			},
			wantIDs: []string{"match"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := overlappingPartitionKeyRanges(c.feed, c.ranges)
			gotIDs := make([]string, 0, len(out))
			for _, r := range out {
				gotIDs = append(gotIDs, r.ID)
			}
			if len(c.wantIDs) == 0 {
				require.Empty(t, gotIDs)
			} else {
				require.Equal(t, c.wantIDs, gotIDs)
			}
		})
	}
}
