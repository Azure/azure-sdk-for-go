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

func TestContainerResponseParsing(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("etag")
	properties := ContainerProperties{
		ID:           "someId",
		ETag:         &etag,
		SelfLink:     "someSelfLink",
		ResourceID:   "someResourceId",
		LastModified: nowAsUnix,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"somePath"},
			Version: 2,
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
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", []policy.Policy{}, []policy.Policy{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)
	parsedResponse, err := newContainerResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if parsedResponse.ContainerProperties == nil {
		t.Fatal("parsedResponse.ContainerProperties is nil")
	}

	if properties.ID != parsedResponse.ContainerProperties.ID {
		t.Errorf("Expected Id to be %s, but got %s", properties.ID, parsedResponse.ContainerProperties.ID)
	}

	if *properties.ETag != *parsedResponse.ContainerProperties.ETag {
		t.Errorf("Expected ETag to be %s, but got %s", *properties.ETag, *parsedResponse.ContainerProperties.ETag)
	}

	if properties.SelfLink != parsedResponse.ContainerProperties.SelfLink {
		t.Errorf("Expected SelfLink to be %s, but got %s", properties.SelfLink, parsedResponse.ContainerProperties.SelfLink)
	}

	if properties.ResourceID != parsedResponse.ContainerProperties.ResourceID {
		t.Errorf("Expected ResourceId to be %s, but got %s", properties.ResourceID, parsedResponse.ContainerProperties.ResourceID)
	}

	if properties.LastModified != parsedResponse.ContainerProperties.LastModified {
		t.Errorf("Expected LastModified.Time to be %v, but got %v", properties.LastModified, parsedResponse.ContainerProperties.LastModified)
	}

	if properties.PartitionKeyDefinition.Paths[0] != parsedResponse.ContainerProperties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Expected PartitionKeyDefinition.Paths[0] to be %s, but got %s", properties.PartitionKeyDefinition.Paths[0], parsedResponse.ContainerProperties.PartitionKeyDefinition.Paths[0])
	}

	if properties.PartitionKeyDefinition.Version != parsedResponse.ContainerProperties.PartitionKeyDefinition.Version {
		t.Errorf("Expected PartitionKeyDefinition.Version to be %d, but got %d", properties.PartitionKeyDefinition.Version, parsedResponse.ContainerProperties.PartitionKeyDefinition.Version)
	}

	if parsedResponse.ActivityID != "someActivityId" {
		t.Errorf("Expected ActivityId to be %s, but got %s", "someActivityId", parsedResponse.ActivityID)
	}

	if parsedResponse.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, parsedResponse.RequestCharge)
	}

	if parsedResponse.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", parsedResponse.ETag)
	}
}
