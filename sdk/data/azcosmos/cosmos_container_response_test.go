// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestContainerResponseParsing(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := &CosmosContainerProperties{
		Id:           "someId",
		ETag:         "someEtag",
		SelfLink:     "someSelfLink",
		ResourceId:   "someResourceId",
		LastModified: &now,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"somePath"},
			Version: PartitionKeyDefinitionVersion2,
		},
	}

	container := &CosmosContainer{
		Id: "someId",
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

	pl := azruntime.NewPipeline(srv)
	resp, _ := pl.Do(req)
	parsedResponse, err := newCosmosContainerResponse(resp, container)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if parsedResponse.ContainerProperties == nil {
		t.Fatal("parsedResponse.ContainerProperties is nil")
	}

	if properties.Id != parsedResponse.ContainerProperties.Id {
		t.Errorf("Expected Id to be %s, but got %s", properties.Id, parsedResponse.ContainerProperties.Id)
	}

	if properties.ETag != parsedResponse.ContainerProperties.ETag {
		t.Errorf("Expected ETag to be %s, but got %s", properties.ETag, parsedResponse.ContainerProperties.ETag)
	}

	if properties.SelfLink != parsedResponse.ContainerProperties.SelfLink {
		t.Errorf("Expected SelfLink to be %s, but got %s", properties.SelfLink, parsedResponse.ContainerProperties.SelfLink)
	}

	if properties.ResourceId != parsedResponse.ContainerProperties.ResourceId {
		t.Errorf("Expected ResourceId to be %s, but got %s", properties.ResourceId, parsedResponse.ContainerProperties.ResourceId)
	}

	if properties.LastModified.Time != parsedResponse.ContainerProperties.LastModified.Time {
		t.Errorf("Expected LastModified.Time to be %s, but got %s", properties.LastModified.Time.UTC(), parsedResponse.ContainerProperties.LastModified.Time.UTC())
	}

	if properties.PartitionKeyDefinition.Paths[0] != parsedResponse.ContainerProperties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Expected PartitionKeyDefinition.Paths[0] to be %s, but got %s", properties.PartitionKeyDefinition.Paths[0], parsedResponse.ContainerProperties.PartitionKeyDefinition.Paths[0])
	}

	if properties.PartitionKeyDefinition.Version != parsedResponse.ContainerProperties.PartitionKeyDefinition.Version {
		t.Errorf("Expected PartitionKeyDefinition.Version to be %d, but got %d", properties.PartitionKeyDefinition.Version, parsedResponse.ContainerProperties.PartitionKeyDefinition.Version)
	}

	if parsedResponse.ActivityId != "someActivityId" {
		t.Errorf("Expected ActivityId to be %s, but got %s", "someActivityId", parsedResponse.ActivityId)
	}

	if parsedResponse.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, parsedResponse.RequestCharge)
	}

	if parsedResponse.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", parsedResponse.ETag)
	}

	if parsedResponse.ContainerProperties.Container != container {
		t.Errorf("Expected Container to be %v, but got %v", container, parsedResponse.ContainerProperties.Container)
	}
}
