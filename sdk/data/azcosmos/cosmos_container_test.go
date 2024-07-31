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

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

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

func TestContainerDeleteItem(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(204))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.DeleteItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	if err != nil {
		t.Fatalf("Failed to delete item: %v", err)
	}

	if resp.RawResponse == nil {
		t.Fatal("RawResponse is nil")
	}

	if resp.ActivityID == "" {
		t.Fatal("Activity id was not returned")
	}

	if resp.RequestCharge == 0 {
		t.Fatal("Request charge was not returned")
	}

	if resp.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, resp.RequestCharge)
	}

	if resp.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", resp.ETag)
	}

	if verifier.requests[0].method != http.MethodDelete {
		t.Errorf("Expected method to be %s, but got %s", http.MethodDelete, verifier.requests[0].method)
	}

	if verifier.requests[0].url.RequestURI() != "/dbs/databaseId/colls/containerId/docs/doc1" {
		t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs/doc1", verifier.requests[0].url.RequestURI())
	}
}

func TestContainerReadItem(t *testing.T) {
	jsonString := []byte(`{"id":"doc1","foo":"bar"}`)
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.ReadItem(context.TODO(), NewPartitionKeyString("1"), "doc1", nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if string(resp.Value) != string(jsonString) {
		t.Errorf("Expected value to be %s, but got %s", string(jsonString), string(resp.Value))
	}

	if resp.RawResponse == nil {
		t.Fatal("RawResponse is nil")
	}

	if resp.ActivityID == "" {
		t.Fatal("Activity id was not returned")
	}

	if resp.RequestCharge == 0 {
		t.Fatal("Request charge was not returned")
	}

	if resp.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, resp.RequestCharge)
	}

	if resp.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", resp.ETag)
	}

	if verifier.requests[0].method != http.MethodGet {
		t.Errorf("Expected method to be %s, but got %s", http.MethodGet, verifier.requests[0].method)
	}

	if verifier.requests[0].url.RequestURI() != "/dbs/databaseId/colls/containerId/docs/doc1" {
		t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs/doc1", verifier.requests[0].url.RequestURI())
	}
}

func TestContainerReplaceItem(t *testing.T) {
	jsonString := []byte(`{"id":"doc1","foo":"bar"}`)
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.ReplaceItem(context.TODO(), NewPartitionKeyString("1"), "doc1", jsonString, nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if string(resp.Value) != string(jsonString) {
		t.Errorf("Expected value to be %s, but got %s", string(jsonString), string(resp.Value))
	}

	if resp.RawResponse == nil {
		t.Fatal("RawResponse is nil")
	}

	if resp.ActivityID == "" {
		t.Fatal("Activity id was not returned")
	}

	if resp.RequestCharge == 0 {
		t.Fatal("Request charge was not returned")
	}

	if resp.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, resp.RequestCharge)
	}

	if resp.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", resp.ETag)
	}

	if verifier.requests[0].method != http.MethodPut {
		t.Errorf("Expected method to be %s, but got %s", http.MethodPut, verifier.requests[0].method)
	}

	if verifier.requests[0].body != string(jsonString) {
		t.Errorf("Expected body to be %s, but got %s", string(jsonString), string(verifier.requests[0].body))
	}

	if verifier.requests[0].url.RequestURI() != "/dbs/databaseId/colls/containerId/docs/doc1" {
		t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs/doc1", verifier.requests[0].url.RequestURI())
	}
}

func TestContainerUpsertItem(t *testing.T) {
	jsonString := []byte(`{"id":"doc1","foo":"bar"}`)
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.UpsertItem(context.TODO(), NewPartitionKeyString("1"), jsonString, nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if string(resp.Value) != string(jsonString) {
		t.Errorf("Expected value to be %s, but got %s", string(jsonString), string(resp.Value))
	}

	if resp.RawResponse == nil {
		t.Fatal("RawResponse is nil")
	}

	if resp.ActivityID == "" {
		t.Fatal("Activity id was not returned")
	}

	if resp.RequestCharge == 0 {
		t.Fatal("Request charge was not returned")
	}

	if resp.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, resp.RequestCharge)
	}

	if resp.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", resp.ETag)
	}

	if verifier.requests[0].method != http.MethodPost {
		t.Errorf("Expected method to be %s, but got %s", http.MethodPost, verifier.requests[0].method)
	}

	if verifier.requests[0].headers.Get(cosmosHeaderIsUpsert) != "true" {
		t.Errorf("Expected header to be %s, but got %s", cosmosHeaderIsUpsert, verifier.requests[0].headers.Get(cosmosHeaderIsUpsert))
	}

	if verifier.requests[0].body != string(jsonString) {
		t.Errorf("Expected body to be %s, but got %s", string(jsonString), string(verifier.requests[0].body))
	}

	if verifier.requests[0].url.RequestURI() != "/dbs/databaseId/colls/containerId/docs" {
		t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs", verifier.requests[0].url.RequestURI())
	}
}

func TestContainerCreateItem(t *testing.T) {
	jsonString := []byte(`{"id":"doc1","foo":"bar"}`)
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.UpsertItem(context.TODO(), NewPartitionKeyString("1"), jsonString, nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if string(resp.Value) != string(jsonString) {
		t.Errorf("Expected value to be %s, but got %s", string(jsonString), string(resp.Value))
	}

	if resp.RawResponse == nil {
		t.Fatal("RawResponse is nil")
	}

	if resp.ActivityID == "" {
		t.Fatal("Activity id was not returned")
	}

	if resp.RequestCharge == 0 {
		t.Fatal("Request charge was not returned")
	}

	if resp.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, resp.RequestCharge)
	}

	if resp.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", resp.ETag)
	}

	if verifier.requests[0].method != http.MethodPost {
		t.Errorf("Expected method to be %s, but got %s", http.MethodPost, verifier.requests[0].method)
	}

	if verifier.requests[0].headers.Get(cosmosHeaderIsUpsert) == "" {
		t.Errorf("Expected header to be empty, but got %s", verifier.requests[0].headers.Get(cosmosHeaderIsUpsert))
	}

	if verifier.requests[0].body != string(jsonString) {
		t.Errorf("Expected body to be %s, but got %s", string(jsonString), string(verifier.requests[0].body))
	}

	if verifier.requests[0].url.RequestURI() != "/dbs/databaseId/colls/containerId/docs" {
		t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs", verifier.requests[0].url.RequestURI())
	}
}

func TestContainerQueryItems(t *testing.T) {
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

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	receivedIds := []string{}
	queryPager := container.NewQueryItemsPager("select * from c", NewPartitionKeyString("1"), nil)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query items: %v", err)
		}

		for _, item := range queryResponse.Items {
			var itemResponseBody map[string]interface{}
			err = json.Unmarshal(item, &itemResponseBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			receivedIds = append(receivedIds, itemResponseBody["id"].(string))
		}

		if queryPager.More() && (queryResponse.ContinuationToken == nil || *queryResponse.ContinuationToken != "someContinuationToken") {
			t.Errorf("Expected ContinuationToken to be %s, but got %s", "someContinuationToken", *queryResponse.ContinuationToken)
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
		if receivedIds[i] != "doc"+strconv.Itoa(i+1) {
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

		if request.url.RequestURI() != "/dbs/databaseId/colls/containerId/docs" {
			t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs", request.url.RequestURI())
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

func TestContainerExecuteBatch(t *testing.T) {
	batchResponseRaw := []map[string]interface{}{
		{"statusCode": 200, "requestCharge": 10.0, "eTag": "someETag", "resourceBody": "someBody"},
		{"statusCode": 201, "requestCharge": 11.0, "eTag": "someETag2"},
	}

	jsonString, err := json.Marshal(batchResponseRaw)
	if err != nil {
		t.Fatal(err)
	}

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithStatusCode(http.StatusOK),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	pk := NewPartitionKeyString("pk")
	batch := container.NewTransactionalBatch(pk)
	_, err = container.ExecuteTransactionalBatch(context.TODO(), batch, nil)
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}

	batch.ReadItem("someId", nil)

	body := map[string]string{
		"foo": "bar",
	}

	itemMarshall, _ := json.Marshal(body)
	batch.CreateItem(itemMarshall, nil)

	_, err = container.ExecuteTransactionalBatch(context.TODO(), batch, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(verifier.requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(verifier.requests))
	}

	request := verifier.requests[0]

	if request.method != http.MethodPost {
		t.Errorf("Expected method to be %s, but got %s", http.MethodPost, request.method)
	}

	if request.url.RequestURI() != "/dbs/databaseId/colls/containerId/docs" {
		t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs", request.url.RequestURI())
	}

	marshalledOperations, _ := json.Marshal(batch.operations)
	if request.body != string(marshalledOperations) {
		t.Errorf("Expected %v, but got %v", string(marshalledOperations), request.body)
	}
}

func TestContainerPatchItem(t *testing.T) {
	jsonString := []byte(`{"id":"doc1","foo":"bar","hello":"world"}`)
	patchOpt := PatchOperations{}
	patchOpt.AppendSet("/hello", "world")

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), internal: internalClient, gem: gem}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.PatchItem(context.TODO(), NewPartitionKeyString("1"), "doc1", patchOpt, nil)
	if err != nil {
		t.Fatalf("Failed to patch item: %v", err)
	}

	if string(resp.Value) != string(jsonString) {
		t.Errorf("Expected value to be %s, but got %s", string(jsonString), string(resp.Value))
	}

	if resp.RawResponse == nil {
		t.Fatal("RawResponse is nil")
	}

	if resp.ActivityID == "" {
		t.Fatal("Activity id was not returned")
	}

	if resp.RequestCharge == 0 {
		t.Fatal("Request charge was not returned")
	}

	if resp.RequestCharge != 13.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 13.42, resp.RequestCharge)
	}

	if resp.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", resp.ETag)
	}

	if verifier.requests[0].method != http.MethodPatch {
		t.Errorf("Expected method to be %s, but got %s", http.MethodPatch, verifier.requests[0].method)
	}

	if verifier.requests[0].url.RequestURI() != "/dbs/databaseId/colls/containerId/docs/doc1" {
		t.Errorf("Expected url to be %s, but got %s", "/dbs/databaseId/colls/containerId/docs/doc1", verifier.requests[0].url.RequestURI())
	}
}
