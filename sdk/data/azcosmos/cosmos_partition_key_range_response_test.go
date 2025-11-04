// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestPartitionKeyRangeResponseParsing(t *testing.T) {
	jsonString := []byte(`{
        "_rid": "mockResourceId",
        "PartitionKeyRanges": [
            {
                "id": "0",
                "_rid": "rid1",
                "_etag": "etag1",
                "minInclusive": "FF",
                "maxExclusive": "5A",
                "_ridPrefix": 1001,
                "_self": "self1",
                "throughputFraction": 0.25,
                "status": "online",
                "parents": ["parent1", "parent2"],
                "ownedArchivalPKRangeIds": ["archive1"],
                "_ts": 12345,
                "lsn": 9876
            },
            {
                "id": "1",
                "_rid": "rid2",
                "_etag": "etag2",
                "minInclusive": "5A",
                "maxExclusive": "FF",
                "_ridPrefix": 1002,
                "_self": "self2",
                "throughputFraction": 0.75,
                "status": "online",
                "parents": ["parent3"],
                "ownedArchivalPKRangeIds": ["archive2", "archive3"],
                "_ts": 67890,
                "lsn": 54321
            }
        ],
        "_count": 2
    }`)

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "mockEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "mockActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "15.75"))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	parsedResponse, err := newPartitionKeyRangeResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if parsedResponse.ActivityID != "mockActivityId" {
		t.Errorf("Expected ActivityID to be %s, but got %s", "mockActivityId", parsedResponse.ActivityID)
	}

	if parsedResponse.RequestCharge != 15.75 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 15.75, parsedResponse.RequestCharge)
	}

	if parsedResponse.ResourceID != "mockResourceId" {
		t.Errorf("Expected Rid to be %s, but got %s", "mockResourceId", parsedResponse.ResourceID)
	}

	if parsedResponse.Count != 2 {
		t.Errorf("Expected Count to be %d, but got %d", 2, parsedResponse.Count)
	}

	if len(parsedResponse.PartitionKeyRanges) != 2 {
		t.Fatalf("Expected 2 partition key ranges, but got %d", len(parsedResponse.PartitionKeyRanges))
	}

	parsedPkr1 := parsedResponse.PartitionKeyRanges[0]

	if parsedPkr1.ID != "0" {
		t.Errorf("Expected ID to be %s, but got %s", "0", parsedPkr1.ID)
	}

	if parsedPkr1.MinInclusive != "FF" {
		t.Errorf("Expected MinInclusive to be %s, but got %s", "FF", parsedPkr1.MinInclusive)
	}

	if parsedPkr1.MaxExclusive != "5A" {
		t.Errorf("Expected MaxExclusive to be %s, but got %s", "5A", parsedPkr1.MaxExclusive)
	}

	if parsedPkr1.ThroughputFraction != 0.25 {
		t.Errorf("Expected ThroughputFraction to be %f, but got %f", 0.25, parsedPkr1.ThroughputFraction)
	}

	if len(parsedPkr1.Parents) != 2 || parsedPkr1.Parents[0] != "parent1" || parsedPkr1.Parents[1] != "parent2" {
		t.Errorf("Parents array not parsed correctly")
	}

	if len(parsedPkr1.OwnedArchivalPKRangeIds) != 1 || parsedPkr1.OwnedArchivalPKRangeIds[0] != "archive1" {
		t.Errorf("OwnedArchivalPKRangeIds not parsed correctly")
	}

	parsedPkr2 := parsedResponse.PartitionKeyRanges[1]

	if parsedPkr2.ID != "1" {
		t.Errorf("Expected ID to be %s, but got %s", "1", parsedPkr2.ID)
	}

	if parsedPkr2.MinInclusive != "5A" {
		t.Errorf("Expected MinInclusive to be %s, but got %s", "5A", parsedPkr2.MinInclusive)
	}

	if parsedPkr2.MaxExclusive != "FF" {
		t.Errorf("Expected MaxExclusive to be %s, but got %s", "FF", parsedPkr2.MaxExclusive)
	}
}
