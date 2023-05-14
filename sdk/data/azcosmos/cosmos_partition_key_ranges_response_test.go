// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestPartitionKeyRangesResponseParsing(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("etag")
	properties := &PartitionKeyRangesProperties{
		ResourceID: "someResourceId",
		Count:      1,
		PartitionKeyRanges: []PartitionKeyRange{
			{
				ID:           "someId",
				ETag:         &etag,
				SelfLink:     "someSelfLink",
				LastModified: nowAsUnix,
				MinInclusive: "someMinInclusive",
				MaxExclusive: "someMaxExclusive",
			},
		},
	}

	jsonString, err := json.Marshal(properties)
	if err != nil {
		t.Fatal(err)
	}

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"))

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

	if parsedResponse.PartitionKeyRangesProperties == nil {
		t.Fatal("parsedResponse.PartitionKeyRangesProperties is nil")
	}

	if properties.ResourceID != parsedResponse.PartitionKeyRangesProperties.ResourceID {
		t.Errorf("Expected ResourceID to be %s, but got %s", properties.ResourceID, parsedResponse.PartitionKeyRangesProperties.ResourceID)
	}

	if properties.Count != parsedResponse.PartitionKeyRangesProperties.Count {
		t.Errorf("Expected Count to be %d, but got %d", properties.Count, parsedResponse.PartitionKeyRangesProperties.Count)
	}

	if len(properties.PartitionKeyRanges) != len(parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges) {
		t.Errorf("Expected PartitionKeyRanges length to be %d, but got %d", len(properties.PartitionKeyRanges), len(parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges))
	}

	if properties.PartitionKeyRanges[0].ID != parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].ID {
		t.Errorf("Expected ID to be %s, but got %s", properties.PartitionKeyRanges[0].ID, parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].ID)
	}

	if *properties.PartitionKeyRanges[0].ETag != *parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].ETag {
		t.Errorf("Expected ETag to be %v, but got %v", properties.PartitionKeyRanges[0].ETag, parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].ETag)
	}

	if properties.PartitionKeyRanges[0].SelfLink != parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].SelfLink {
		t.Errorf("Expected SelfLink to be %s, but got %s", properties.PartitionKeyRanges[0].SelfLink, parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].SelfLink)
	}

	if properties.PartitionKeyRanges[0].LastModified != parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].LastModified {
		t.Errorf("Expected LastModified to be %s, but got %s", properties.PartitionKeyRanges[0].LastModified, parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].LastModified)
	}

	if properties.PartitionKeyRanges[0].MinInclusive != parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].MinInclusive {
		t.Errorf("Expected MinInclusive to be %s, but got %s", properties.PartitionKeyRanges[0].MinInclusive, parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].MinInclusive)
	}

	if properties.PartitionKeyRanges[0].MaxExclusive != parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].MaxExclusive {
		t.Errorf("Expected MaxExclusive to be %s, but got %s", properties.PartitionKeyRanges[0].MaxExclusive, parsedResponse.PartitionKeyRangesProperties.PartitionKeyRanges[0].MaxExclusive)
	}

	if parsedResponse.ActivityID != "someActivityId" {
		t.Errorf("Expected ActivityId to be %s, but got %s", "someActivityId", parsedResponse.ActivityID)
	}

	if parsedResponse.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", parsedResponse.ETag)
	}
}
