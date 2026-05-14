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
	require.Equal(t, 2, requestsAfterFirstCall, "first call should make 2 HTTP requests (container read + PK ranges)")

	// Second call: should use caches, no additional HTTP requests
	// (no more responses queued — would panic if a request was made)
	feedRanges2, err := container.GetFeedRanges(context.TODO())
	require.NoError(t, err)
	require.Equal(t, 2, len(feedRanges2))
	require.Equal(t, feedRanges[0], feedRanges2[0])
	require.Equal(t, feedRanges[1], feedRanges2[1])

	require.Equal(t, requestsAfterFirstCall, srv.Requests(), "second call should make 0 HTTP requests (cache hit)")
}
