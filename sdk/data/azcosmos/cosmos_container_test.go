// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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

func TestContainerQueryItemsPerPartitionKey(t *testing.T) {
	jsonStringpage1 := []byte(`{"Documents":[{"id":"doc1","foo":"bar"},{"id":"doc2","foo":"bar"}]}`)
	jsonStringpage2 := []byte(`{"Documents":[{"id":"doc3","foo":"bar"},{"id":"doc4","foo":"bar"},{"id":"doc5","foo":"bar"}]}`)
	
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(
		mock.WithBody(jsonStringpage1),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderQueryMetrics, "someQueryMetrics"),
		mock.WithHeader(cosmosHeaderIndexUtilization, "someIndexUtilization"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithHeader(cosmosHeaderContinuationToken, "someContinuationToken"),
		mock.WithStatusCode(200))
	srv.AppendResponse(
		mock.WithBody(jsonStringpage2),
		mock.WithHeader(cosmosHeaderQueryMetrics, "someQueryMetrics"),
		mock.WithHeader(cosmosHeaderIndexUtilization, "someIndexUtilization"),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client := &Client{endpoint: srv.URL(), pipeline: pl}

	database, _:=newDatabase("databaseId", client)
	container, _:= newContainer("containerId",database)

	receivedIds := []string{}
	queryPager := container.QueryItemsByPartitionKey("select * from c", NewPartitionKeyString("1"), nil)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query items: %v", err)
		}

		for _, item := range queryResponse.Items {
			var itemResponseBody map[string]interface{}
			json.Unmarshal(item, &itemResponseBody)
			receivedIds = append(receivedIds, itemResponseBody["id"].(string))
		}

		if queryPager.More() && queryResponse.ContinuationToken != "someContinuationToken" {
			t.Errorf("Expected ContinuationToken to be %s, but got %s", "someContinuationToken", queryResponse.ContinuationToken)
		}

		if queryResponse.QueryMetrics == nil || *queryResponse.QueryMetrics != "someQueryMetrics" {
			t.Errorf("Expected QueryMetrics to be %s, but got %s", "someQueryMetrics", *queryResponse.QueryMetrics)
		}

		if queryResponse.IndexMetrics == nil || *queryResponse.IndexMetrics != "someIndexUtilization" {
			t.Errorf("Expected IndexMetrics to be %s, but got %s", "someIndexUtilization", *queryResponse.IndexMetrics)
		}

		if queryResponse.ActivityID == "" {
			t.Fatal("Activity id was not returned")
		}

		if queryResponse.RequestCharge == 0 {
			t.Fatal("Request charge was not returned")
		}
	}

	for i := 0; i < 5; i++ {
		if receivedIds[i] != "doc" + strconv.Itoa(i+1) {
			t.Fatalf("Expected id %d, got %s", i, receivedIds[i])
		}
	}

	if len(verifier.requests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(verifier.requests))
	}

	for index, request := range verifier.requests {
		if request.method != http.MethodPost {
			t.Errorf("Expected method to be %s, but got %s", http.MethodPost, request.method)
		}

		if !request.isQuery {
			t.Errorf("Expected request to be a query, but it was not")
		}

		if request.body != "{\"query\":\"select * from c\"}" {
			t.Errorf("Expected %v, but got %v", "{\"query\":\"select * from c\"}", request.body)
		}

		if request.contentType != cosmosHeaderValuesQuery {
			t.Errorf("Expected %v, but got %v", cosmosHeaderValuesQuery, request.contentType)
		}
		
		if index == 0 && request.headers.Get(cosmosHeaderContinuationToken) != "" {
			t.Errorf("Expected ContinuationToken to be %s, but got %s", "", request.headers.Get(cosmosHeaderContinuationToken))
		}

		if index == 1 && request.headers.Get(cosmosHeaderContinuationToken) != "someContinuationToken" {
			t.Errorf("Expected ContinuationToken to be %s, but got %s", "someContinuationToken", request.headers.Get(cosmosHeaderContinuationToken))
		}
	}
}