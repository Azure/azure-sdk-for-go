// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)


func TestContainerRead(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("etag")
	properties := ContainerProperties{
		ID:           "containerId",
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
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}

	database, _:=newDatabase("databaseId", client)
	container, _:= newContainer("containerId",database)

	if container.ID() != "containerId" {
		t.Errorf("Expected container ID to be %s, but got %s", "containerId", container.ID())
	}

	resp, err := container.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	if resp.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if resp.ContainerProperties == nil {
		t.Fatal("parsedResponse.ContainerProperties is nil")
	}

	if properties.ID != resp.ContainerProperties.ID {
		t.Errorf("Expected Id to be %s, but got %s", properties.ID, resp.ContainerProperties.ID)
	}

	if *properties.ETag != *resp.ContainerProperties.ETag {
		t.Errorf("Expected ETag to be %s, but got %s", *properties.ETag, *resp.ContainerProperties.ETag)
	}

	if properties.SelfLink != resp.ContainerProperties.SelfLink {
		t.Errorf("Expected SelfLink to be %s, but got %s", properties.SelfLink, resp.ContainerProperties.SelfLink)
	}

	if properties.ResourceID != resp.ContainerProperties.ResourceID {
		t.Errorf("Expected ResourceId to be %s, but got %s", properties.ResourceID, resp.ContainerProperties.ResourceID)
	}

	if properties.LastModified != resp.ContainerProperties.LastModified {
		t.Errorf("Expected LastModified.Time to be %v, but got %v", properties.LastModified, resp.ContainerProperties.LastModified)
	}

	if properties.PartitionKeyDefinition.Paths[0] != resp.ContainerProperties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Expected PartitionKeyDefinition.Paths[0] to be %s, but got %s", properties.PartitionKeyDefinition.Paths[0], resp.ContainerProperties.PartitionKeyDefinition.Paths[0])
	}

	if properties.PartitionKeyDefinition.Version != resp.ContainerProperties.PartitionKeyDefinition.Version {
		t.Errorf("Expected PartitionKeyDefinition.Version to be %d, but got %d", properties.PartitionKeyDefinition.Version, resp.ContainerProperties.PartitionKeyDefinition.Version)
	}

	if resp.ActivityID != "someActivityId" {
		t.Errorf("Expected ActivityId to be %s, but got %s", "someActivityId", resp.ActivityID)
	}

	if resp.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, resp.RequestCharge)
	}

	if resp.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", resp.ETag)
	}
}
