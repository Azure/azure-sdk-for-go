// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestPartitionKeyRangeResponseParsing(t *testing.T) {
	// Create a sample PartitionKeyRangeResponse
	mockResponse := PartitionKeyRangeResponse{
		Rid: "mockResourceId",
		PartitionKeyRanges: []PartitionKeyRange{
			{
				ID:                      "0",
				Rid:                     "rid1",
				ETag:                    "etag1",
				MinInclusive:            "FF",
				MaxExclusive:            "5A",
				RidPrefix:               1001,
				Self:                    "self1",
				ThroughputFraction:      0.25,
				Status:                  "online",
				Parents:                 []string{"parent1", "parent2"},
				OwnedArchivalPKRangeIds: []string{"archive1"},
				Timestamp:               12345,
				LSN:                     9876,
			},
			{
				ID:                      "1",
				Rid:                     "rid2",
				ETag:                    "etag2",
				MinInclusive:            "5A",
				MaxExclusive:            "FF",
				RidPrefix:               1002,
				Self:                    "self2",
				ThroughputFraction:      0.75,
				Status:                  "online",
				Parents:                 []string{"parent3"},
				OwnedArchivalPKRangeIds: []string{"archive2", "archive3"},
				Timestamp:               67890,
				LSN:                     54321,
			},
		},
		Count: 2,
	}

	// Marshal only the data structure (mockResponse) to JSON for the mock server response
	jsonString, err := json.Marshal(mockResponse)
	if err != nil {
		t.Fatal(err)
	}

	// Setup mock server
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "mockEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "mockActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "15.75"))

	// Create request and pipeline
	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)

	// Parse the response
	parsedResponse, err := newPartitionKeyRangeResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	// Validate basic response properties
	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	// Verify response metadata
	if parsedResponse.ActivityID != "mockActivityId" {
		t.Errorf("Expected ActivityID to be %s, but got %s", "mockActivityId", parsedResponse.ActivityID)
	}

	if parsedResponse.RequestCharge != 15.75 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 15.75, parsedResponse.RequestCharge)
	}

	if parsedResponse.ETag != "mockEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "mockEtag", parsedResponse.ETag)
	}

	// Verify PartitionKeyRangeResponse specific fields
	if parsedResponse.Rid != mockResponse.Rid {
		t.Errorf("Expected Rid to be %s, but got %s", mockResponse.Rid, parsedResponse.Rid)
	}

	if parsedResponse.Count != mockResponse.Count {
		t.Errorf("Expected Count to be %d, but got %d", mockResponse.Count, parsedResponse.Count)
	}

	// Verify PartitionKeyRanges array length
	if len(parsedResponse.PartitionKeyRanges) != len(mockResponse.PartitionKeyRanges) {
		t.Fatalf("Expected %d partition key ranges, but got %d",
			len(mockResponse.PartitionKeyRanges), len(parsedResponse.PartitionKeyRanges))
	}

	// Verify first PartitionKeyRange
	pkr1 := mockResponse.PartitionKeyRanges[0]
	parsedPkr1 := parsedResponse.PartitionKeyRanges[0]

	if pkr1.ID != parsedPkr1.ID {
		t.Errorf("Expected ID to be %s, but got %s", pkr1.ID, parsedPkr1.ID)
	}

	if pkr1.Rid != parsedPkr1.Rid {
		t.Errorf("Expected Rid to be %s, but got %s", pkr1.Rid, parsedPkr1.Rid)
	}

	if pkr1.ETag != parsedPkr1.ETag {
		t.Errorf("Expected ETag to be %s, but got %s", pkr1.ETag, parsedPkr1.ETag)
	}

	if pkr1.MinInclusive != parsedPkr1.MinInclusive {
		t.Errorf("Expected MinInclusive to be %s, but got %s", pkr1.MinInclusive, parsedPkr1.MinInclusive)
	}

	if pkr1.MaxExclusive != parsedPkr1.MaxExclusive {
		t.Errorf("Expected MaxExclusive to be %s, but got %s", pkr1.MaxExclusive, parsedPkr1.MaxExclusive)
	}

	if pkr1.RidPrefix != parsedPkr1.RidPrefix {
		t.Errorf("Expected RidPrefix to be %d, but got %d", pkr1.RidPrefix, parsedPkr1.RidPrefix)
	}

	if pkr1.Self != parsedPkr1.Self {
		t.Errorf("Expected Self to be %s, but got %s", pkr1.Self, parsedPkr1.Self)
	}

	if pkr1.ThroughputFraction != parsedPkr1.ThroughputFraction {
		t.Errorf("Expected ThroughputFraction to be %f, but got %f",
			pkr1.ThroughputFraction, parsedPkr1.ThroughputFraction)
	}

	if pkr1.Status != parsedPkr1.Status {
		t.Errorf("Expected Status to be %s, but got %s", pkr1.Status, parsedPkr1.Status)
	}

	// Check arrays
	if len(pkr1.Parents) != len(parsedPkr1.Parents) {
		t.Errorf("Expected Parents array length to be %d, but got %d",
			len(pkr1.Parents), len(parsedPkr1.Parents))
	} else {
		for i, parent := range pkr1.Parents {
			if parent != parsedPkr1.Parents[i] {
				t.Errorf("Expected Parents[%d] to be %s, but got %s", i, parent, parsedPkr1.Parents[i])
			}
		}
	}

	if len(pkr1.OwnedArchivalPKRangeIds) != len(parsedPkr1.OwnedArchivalPKRangeIds) {
		t.Errorf("Expected OwnedArchivalPKRangeIds array length to be %d, but got %d",
			len(pkr1.OwnedArchivalPKRangeIds), len(parsedPkr1.OwnedArchivalPKRangeIds))
	} else {
		for i, id := range pkr1.OwnedArchivalPKRangeIds {
			if id != parsedPkr1.OwnedArchivalPKRangeIds[i] {
				t.Errorf("Expected OwnedArchivalPKRangeIds[%d] to be %s, but got %s",
					i, id, parsedPkr1.OwnedArchivalPKRangeIds[i])
			}
		}
	}

	if pkr1.Timestamp != parsedPkr1.Timestamp {
		t.Errorf("Expected Timestamp to be %d, but got %d", pkr1.Timestamp, parsedPkr1.Timestamp)
	}

	if pkr1.LSN != parsedPkr1.LSN {
		t.Errorf("Expected LSN to be %d, but got %d", pkr1.LSN, parsedPkr1.LSN)
	}

	// Verify second PartitionKeyRange (simplified version)
	pkr2 := mockResponse.PartitionKeyRanges[1]
	parsedPkr2 := parsedResponse.PartitionKeyRanges[1]

	if pkr2.ID != parsedPkr2.ID {
		t.Errorf("Expected ID to be %s, but got %s", pkr2.ID, parsedPkr2.ID)
	}

	if pkr2.MinInclusive != parsedPkr2.MinInclusive {
		t.Errorf("Expected MinInclusive to be %s, but got %s", pkr2.MinInclusive, parsedPkr2.MinInclusive)
	}

	if pkr2.MaxExclusive != parsedPkr2.MaxExclusive {
		t.Errorf("Expected MaxExclusive to be %s, but got %s", pkr2.MaxExclusive, parsedPkr2.MaxExclusive)
	}
}
