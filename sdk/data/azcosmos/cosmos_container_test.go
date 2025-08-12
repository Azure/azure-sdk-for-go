// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	defaultEndpoint, _ := url.Parse(srv.URL())
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(200))

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
	defer close()
	srv.SetResponse(
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithStatusCode(204))

	verifier := pipelineVerifier{}

	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
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
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
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
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
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
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
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
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
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
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
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
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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
	defaultEndpoint, _ := url.Parse(srv.URL())
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
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}

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

func TestContainerReadPartitionKeyRanges(t *testing.T) {
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
	"_count": 100
	}`)

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

	resp, err := container.getPartitionKeyRanges(context.TODO(), nil)
	if err != nil {
		t.Fatalf("GetPartitionKeys failed: %v", err)
	}

	if resp.PartitionKeyRanges == nil {
		t.Fatal("PartitionKeyRanges is nil")
	}
	if len(resp.PartitionKeyRanges) != 2 {
		t.Fatalf("Expected 2 partition key ranges, got %d", len(resp.PartitionKeyRanges))
	}
	high_level_rid := resp.ResourceID
	if high_level_rid != "lypXAMSZ-Cs=" {
		t.Errorf("Expected Rid to be lypXAMSZ-Cs=, got %s", high_level_rid)
	}

	if resp.Count != 100 {
		t.Errorf("Expected Count to be 100, got %d", resp.Count)
	}
	pkr1 := resp.PartitionKeyRanges[0]
	if pkr1.ID != "151" {
		t.Errorf("Expected ID to be 151, got %s", pkr1.ID)
	}
	if pkr1.MinInclusive != "05C1E18D2D7F08" {
		t.Errorf("Expected MinInclusive to be 05C1E18D2D7F08, got %s", pkr1.MinInclusive)
	}
	if pkr1.MaxExclusive != "05C1E18D2D83FA" {
		t.Errorf("Expected MaxExclusive to be 05C1E18D2D83FA, got %s", pkr1.MaxExclusive)
	}
	if len(pkr1.Parents) != 3 || pkr1.Parents[0] != "5" {
		t.Errorf("Expected Parents to be [5 10 31], got %v", pkr1.Parents)
	}

	pkr2 := resp.PartitionKeyRanges[1]
	if pkr2.ID != "163" {
		t.Errorf("Expected ID to be 163, got %s", pkr2.ID)
	}
	if pkr2.MinInclusive != "05C1C7FF3903F8" {
		t.Errorf("Expected MinInclusive to be 05C1C7FF3903F8, got %s", pkr2.MinInclusive)
	}
	if pkr2.MaxExclusive != "05C1C9CD673390" {
		t.Errorf("Expected MaxExclusive to be 05C1C9CD673390, got %s", pkr2.MaxExclusive)
	}
	if len(pkr2.Parents) != 3 || pkr2.Parents[0] != "1" {
		t.Errorf("Expected Parents to be [1 19 39], got %v", pkr2.Parents)
	}
}

func TestContainerReadPartitionKeyRangesEmpty(t *testing.T) {
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

	resp, err := container.getPartitionKeyRanges(context.TODO(), nil)
	if err != nil {
		t.Fatalf("GetPartitionKeys failed: %v", err)
	}

	if resp.PartitionKeyRanges == nil {
		t.Fatal("PartitionKeyRanges is nil")
	}
	if len(resp.PartitionKeyRanges) != 0 {
		t.Fatalf("Expected 0 partition key ranges, got %d", len(resp.PartitionKeyRanges))
	}
}

func TestContainerGetChangeFeedWithStartFrom(t *testing.T) {
	changeFeedBody := []byte(
		`{"_rid":"test-rid",
		"Documents":[{"id":"doc1"},{"id":"doc2"}],
		"_count":2}`)
	srv, close := mock.NewTLSServer()
	defaultEndpoint, _ := url.Parse(srv.URL())
	defer close()
	srv.SetResponse(
		mock.WithBody(changeFeedBody),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "5.5"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	feedRange := &FeedRange{
		MinInclusive: "00",
		MaxExclusive: "FF",
	}

	modifiedSince := time.Now().Add(-time.Hour).UTC()
	options := &ChangeFeedOptions{
		StartFrom: &modifiedSince,
		FeedRange: feedRange,
	}

	resp, err := container.GetChangeFeed(context.TODO(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed returned error: %v", err)
	}
	if resp.ResourceID != "test-rid" {
		t.Errorf("Expected ResourceID 'test-rid', got %v", resp.ResourceID)
	}
	if resp.Count != 2 {
		t.Errorf("Expected Count 2, got %v", resp.Count)
	}
	if len(resp.Documents) != 2 {
		t.Errorf("Expected 2 documents, got %v", len(resp.Documents))
	}

	if len(verifier.requests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(verifier.requests))
	}

	request := verifier.requests[1]
	ifModifiedSinceHeader := request.headers.Get(cosmosHeaderIfModifiedSince)
	expectedIfModifiedSince := modifiedSince.Format(time.RFC1123)

	if ifModifiedSinceHeader == "" {
		t.Errorf("If-Modified-Since header was not set")
	} else if ifModifiedSinceHeader != expectedIfModifiedSince {
		t.Errorf("Expected If-Modified-Since header to be %s, but got %s", expectedIfModifiedSince, ifModifiedSinceHeader)
	}
}

func TestContainerGetChangeFeedWithStartFromFiltering(t *testing.T) {
	// This test verifies that:
	// 1. The If-Modified-Since header is properly set based on the StartFrom parameter
	// 2. We can request and retrieve documents with different timestamps

	// First response: All documents when using beginning of time filter
	allDocumentsBody := []byte(`{
		"_rid": "test-rid",
		"Documents": [
			{"id": "doc1", "_ts": 1730000000},
			{"id": "doc2", "_ts": 1735000000},
			{"id": "doc3", "_ts": 1740000000}
		],
		"_count": 3
	}`)

	// Second response: Only documents after the filter time
	filteredDocumentsBody := []byte(`{
		"_rid": "test-rid",
		"Documents": [
			{"id": "doc3", "_ts": 1740000000}
		],
		"_count": 1
	}`)

	srv, close := mock.NewTLSServer()
	defaultEndpoint, _ := url.Parse(srv.URL())
	defer close()

	// Set up mock responses
	srv.SetResponse(
		mock.WithBody(allDocumentsBody),
		mock.WithHeader(cosmosHeaderEtag, "etagAll"),
		mock.WithHeader(cosmosHeaderActivityId, "activityIdAll"),
		mock.WithHeader(cosmosHeaderRequestCharge, "2.5"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	feedRange := &FeedRange{
		MinInclusive: "00",
		MaxExclusive: "FF",
	}

	// 1. First call: Get all documents (from beginning of time)
	beginningOfTime := time.Unix(0, 0).UTC()
	allDocsOptions := &ChangeFeedOptions{
		StartFrom: &beginningOfTime,
		FeedRange: feedRange,
	}

	allDocsResp, err := container.GetChangeFeed(context.TODO(), allDocsOptions)
	if err != nil {
		t.Fatalf("First GetChangeFeed returned error: %v", err)
	}

	if allDocsResp.Count != 3 {
		t.Errorf("Expected 3 documents in first response, got %d", allDocsResp.Count)
	}
	if len(allDocsResp.Documents) != 3 {
		t.Errorf("Expected 3 documents in first response, got %d", len(allDocsResp.Documents))
	}

	var allDocs []map[string]interface{}
	for i, docBytes := range allDocsResp.Documents {
		var doc map[string]interface{}
		if err := json.Unmarshal(docBytes, &doc); err != nil {
			t.Fatalf("Failed to unmarshal document %d: %v", i, err)
		}
		allDocs = append(allDocs, doc)
	}

	expectedIDs := []string{"doc1", "doc2", "doc3"}
	for i, doc := range allDocs {
		if doc["id"] != expectedIDs[i] {
			t.Errorf("Expected document %d to have ID '%s', got '%s'", i, expectedIDs[i], doc["id"])
		}
	}

	if len(verifier.requests) < 2 {
		t.Fatalf("Expected at least 2 requests, got %d", len(verifier.requests))
	}

	firstRequest := verifier.requests[1]
	firstIfModifiedSinceHeader := firstRequest.headers.Get(cosmosHeaderIfModifiedSince)
	firstExpectedIfModifiedSince := beginningOfTime.Format(time.RFC1123)

	if firstIfModifiedSinceHeader == "" {
		t.Errorf("If-Modified-Since header was not set in first request")
	} else if firstIfModifiedSinceHeader != firstExpectedIfModifiedSince {
		t.Errorf("Expected If-Modified-Since header to be %s in first request, but got %s",
			firstExpectedIfModifiedSince, firstIfModifiedSinceHeader)
	}

	// Reset the mock server and verifier for the second test
	srv.SetResponse(
		mock.WithBody(filteredDocumentsBody),
		mock.WithHeader(cosmosHeaderEtag, "etagFiltered"),
		mock.WithHeader(cosmosHeaderActivityId, "activityIdFiltered"),
		mock.WithHeader(cosmosHeaderRequestCharge, "1.5"),
		mock.WithStatusCode(200))

	verifier = pipelineVerifier{}
	internalClient, _ = azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client = &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ = newDatabase("databaseId", client)
	container, _ = newContainer("containerId", database)

	// 2. Second call: Get only documents after midpoint timestamp
	midpointTime := time.Unix(1737000000, 0).UTC() // This should filter out doc1 and doc2, keep only doc3
	filteredOptions := &ChangeFeedOptions{
		StartFrom: &midpointTime,
		FeedRange: feedRange,
	}

	filteredResp, err := container.GetChangeFeed(context.TODO(), filteredOptions)
	if err != nil {
		t.Fatalf("Second GetChangeFeed returned error: %v", err)
	}

	if filteredResp.Count != 1 {
		t.Errorf("Expected 1 document in filtered response, got %d", filteredResp.Count)
	}
	if len(filteredResp.Documents) != 1 {
		t.Errorf("Expected 1 document in filtered response, got %d", len(filteredResp.Documents))
	}

	var filteredDoc map[string]interface{}
	if err := json.Unmarshal(filteredResp.Documents[0], &filteredDoc); err != nil {
		t.Fatalf("Failed to unmarshal filtered document: %v", err)
	}
	if filteredDoc["id"] != "doc3" {
		t.Errorf("Expected filtered document to have ID 'doc3', got '%s'", filteredDoc["id"])
	}

	if len(verifier.requests) < 2 {
		t.Fatalf("Expected at least 2 requests in second test, got %d", len(verifier.requests))
	}

	secondRequest := verifier.requests[1]
	secondIfModifiedSinceHeader := secondRequest.headers.Get(cosmosHeaderIfModifiedSince)
	secondExpectedIfModifiedSince := midpointTime.Format(time.RFC1123)

	if secondIfModifiedSinceHeader == "" {
		t.Errorf("If-Modified-Since header was not set in second request")
	} else if secondIfModifiedSinceHeader != secondExpectedIfModifiedSince {
		t.Errorf("Expected If-Modified-Since header to be %s in second request, but got %s",
			secondExpectedIfModifiedSince, secondIfModifiedSinceHeader)
	}
}

func TestContainerGetChangeFeedForEPKRange(t *testing.T) {
	changeFeedBody := []byte(`{
        "_rid": "test-resource-id",
        "Documents": [{"id": "doc1"}, {"id": "doc2"}],
        "_count": 2
    }`)

	pkRangesBody := []byte(`{
        "_rid": "test-resource-id",
        "PartitionKeyRanges": [{
            "_rid": "range-rid",
            "id": "0",
            "minInclusive": "00",
            "maxExclusive": "FF"
        }],
        "_count": 1
    }`)

	srv, close := mock.NewTLSServer()
	defaultEndpoint, _ := url.Parse(srv.URL())
	defer close()

	// First response should be for the partition key ranges request
	srv.AppendResponse(
		mock.WithBody(pkRangesBody),
		mock.WithHeader(cosmosHeaderActivityId, "pkRangesActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "1.0"),
		mock.WithStatusCode(200))

	// Second response should be for the change feed request
	srv.AppendResponse(
		mock.WithBody(changeFeedBody),
		mock.WithHeader(cosmosHeaderEtag, "\"etag-12345\""),
		mock.WithHeader(cosmosHeaderActivityId, "changeFeedActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "3.5"),
		mock.WithStatusCode(200))

	verifier := pipelineVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	feedRange := &FeedRange{
		MinInclusive: "00",
		MaxExclusive: "FF",
	}
	options := &ChangeFeedOptions{
		MaxItemCount: 10,
		FeedRange:    feedRange,
	}

	resp, err := container.GetChangeFeed(context.TODO(), options)
	if err != nil {
		t.Fatalf("GetChangeFeedForEPKRange failed: %v", err)
	}

	if resp.ResourceID != "test-resource-id" {
		t.Errorf("unexpected ResourceID: got %q, want %q", resp.ResourceID, "test-resource-id")
	}

	if resp.Count != 2 {
		t.Errorf("unexpected Count: got %d, want 2", resp.Count)
	}

	if len(resp.Documents) != 2 {
		t.Errorf("unexpected number of Documents: got %d, want 2", len(resp.Documents))
	}

	if len(verifier.requests) != 2 {
		t.Fatalf("Expected exactly 2 requests (partition key ranges and change feed), got %d", len(verifier.requests))
	}

	// First request should be to get partition key ranges
	pkRangesRequest := verifier.requests[0]
	if !strings.Contains(pkRangesRequest.url.Path, "pkranges") {
		t.Errorf("Expected first request to be for partition key ranges, got URL path: %s", pkRangesRequest.url.Path)
	}
	expectedPkRangesPath := "/dbs/databaseId/colls/containerId/pkranges"
	if !strings.Contains(pkRangesRequest.url.Path, expectedPkRangesPath) {
		t.Errorf("Expected partition key ranges path to contain %s, got %s",
			expectedPkRangesPath, pkRangesRequest.url.Path)
	}

	// Second request should be the change feed request
	changeFeedRequest := verifier.requests[1]
	if !strings.Contains(changeFeedRequest.url.Path, "/docs") {
		t.Errorf("Expected second request to be for documents, got URL path: %s", changeFeedRequest.url.Path)
	}

	pkRangeHeader := changeFeedRequest.headers.Get(headerXmsDocumentDbPartitionKeyRangeId)
	if pkRangeHeader != "0" {
		t.Errorf("Expected partition key range ID '0' in request header, got %q", pkRangeHeader)
	}

	changeFeedHeader := changeFeedRequest.headers.Get(cosmosHeaderChangeFeed)
	if changeFeedHeader != cosmosHeaderValuesChangeFeed {
		t.Errorf("Expected change feed header to be %q, got %q",
			cosmosHeaderValuesChangeFeed, changeFeedHeader)
	}

	if resp.ContinuationToken == "" {
		t.Fatal("expected ContinuationToken to be populated, but it was empty")
	}

	var compositeToken compositeContinuationToken
	err = json.Unmarshal([]byte(resp.ContinuationToken), &compositeToken)
	if err != nil {
		t.Fatalf("failed to unmarshal composite token: %v", err)
	}

	if compositeToken.Version != cosmosCompositeContinuationTokenVersion {
		t.Errorf("unexpected version in composite token: got %d, want %d",
			compositeToken.Version, cosmosCompositeContinuationTokenVersion)
	}

	if compositeToken.ResourceID != "test-resource-id" {
		t.Errorf("unexpected ResourceID in composite token: got %q, want %q",
			compositeToken.ResourceID, "test-resource-id")
	}

	if len(compositeToken.Continuation) != 1 {
		t.Fatalf("unexpected number of continuation ranges: got %d, want 1",
			len(compositeToken.Continuation))
	}

	if compositeToken.Continuation[0].MinInclusive != "00" {
		t.Errorf("unexpected MinInclusive in continuation token: got %q, want %q",
			compositeToken.Continuation[0].MinInclusive, "00")
	}

	if compositeToken.Continuation[0].MaxExclusive != "FF" {
		t.Errorf("unexpected MaxExclusive in continuation token: got %q, want %q",
			compositeToken.Continuation[0].MaxExclusive, "FF")
	}

	if compositeToken.Continuation[0].ContinuationToken == nil {
		t.Fatal("expected ContinuationToken to be set, but it was nil")
	}

	if *compositeToken.Continuation[0].ContinuationToken != azcore.ETag("\"etag-12345\"") {
		t.Errorf("unexpected ContinuationToken: got %q, want %q",
			*compositeToken.Continuation[0].ContinuationToken, "\"etag-12345\"")
	}

	// Now test using the continuation token in a subsequent request
	options2 := &ChangeFeedOptions{
		MaxItemCount: 10,
		Continuation: &resp.ContinuationToken,
	}

	headers := options2.toHeaders(nil)
	if headers == nil {
		t.Fatal("expected headers to be non-nil")
	}

	h := *headers
	if h[headerIfNoneMatch] != "\"etag-12345\"" {
		t.Errorf("unexpected IfNoneMatch header: got %q, want %q",
			h[headerIfNoneMatch], "\"etag-12345\"")
	}

	if h[cosmosHeaderChangeFeed] != cosmosHeaderValuesChangeFeed {
		t.Errorf("unexpected ChangeFeed header in continuation request: got %q, want %q",
			h[cosmosHeaderChangeFeed], cosmosHeaderValuesChangeFeed)
	}
}

func TestContainerGetChangeFeedForPartitionKey(t *testing.T) {
	changeFeedBody := []byte(`{
        "_rid": "test-resource-id",
        "Documents": [{"id": "doc1", "pk": "test-pk"}, {"id": "doc2", "pk": "test-pk"}],
        "_count": 2
    }`)

	srv, close := mock.NewTLSServer()
	defaultEndpoint, _ := url.Parse(srv.URL())
	defer close()

	// Set up mock response
	srv.SetResponse(
		mock.WithBody(changeFeedBody),
		mock.WithHeader(cosmosHeaderEtag, "\"etag-pk-123\""),
		mock.WithHeader(cosmosHeaderActivityId, "pkChangeFeedActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "2.5"),
		mock.WithStatusCode(200))

	// Set up verifier to capture request details
	verifier := pipelineVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	// Create request with partition key
	partitionKey := NewPartitionKeyString("test-pk")
	options := &ChangeFeedOptions{
		MaxItemCount: 10,
		PartitionKey: &partitionKey,
	}

	resp, err := container.GetChangeFeed(context.TODO(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed with PartitionKey failed: %v", err)
	}

	// Validate response content
	if resp.ResourceID != "test-resource-id" {
		t.Errorf("unexpected ResourceID: got %q, want %q", resp.ResourceID, "test-resource-id")
	}

	if resp.Count != 2 {
		t.Errorf("unexpected document count: got %d, want 2", resp.Count)
	}

	if len(resp.Documents) != 2 {
		t.Errorf("unexpected number of documents: got %d, want 2", len(resp.Documents))
	}

	// Validate the partition key was properly set
	if resp.PartitionKey == nil {
		t.Fatal("expected PartitionKey to be set, but it was nil")
	}

	pkJSON, err := resp.PartitionKey.toJsonString()
	if err != nil {
		t.Fatalf("failed to convert partition key to JSON: %v", err)
	}

	if string(pkJSON) != "[\"test-pk\"]" {
		t.Errorf("unexpected partition key value: got %q, want %q", string(pkJSON), "[\"test-pk\"]")
	}

	// Validate request
	if len(verifier.requests) == 0 {
		t.Fatal("no requests were made")
	}
	request := verifier.requests[0]

	// Verify the partition key is present in the header
	partitionKeyHeader := request.headers.Get(cosmosHeaderPartitionKey)
	expectedPKHeader := "[\"test-pk\"]"
	if partitionKeyHeader != expectedPKHeader {
		t.Errorf("unexpected partition key header: got %q, want %q", partitionKeyHeader, expectedPKHeader)
	}

	// Verify change feed header
	changeFeedHeader := request.headers.Get(cosmosHeaderChangeFeed)
	if changeFeedHeader != cosmosHeaderValuesChangeFeed {
		t.Errorf("incorrect change feed header: got %q, want %q", changeFeedHeader, cosmosHeaderValuesChangeFeed)
	}

	// Verify continuation token is populated correctly
	if resp.ContinuationToken == "" {
		t.Fatal("expected ContinuationToken to be populated, but it was empty")
	}

	// Unmarshal and validate continuation token
	var contToken continuationTokenForPartitionKey
	err = json.Unmarshal([]byte(resp.ContinuationToken), &contToken)
	if err != nil {
		t.Fatalf("failed to unmarshal continuation token: %v", err)
	}

	if contToken.Version != cosmosContinuationTokenForPartitionKeyVersion {
		t.Errorf("incorrect token version: got %q, want %q", contToken.Version, cosmosContinuationTokenForPartitionKeyVersion)
	}

	if contToken.ResourceID != "test-resource-id" {
		t.Errorf("incorrect ResourceID in token: got %q, want %q", contToken.ResourceID, "test-resource-id")
	}

	if contToken.PartitionKey == nil {
		t.Errorf("Expected PartitionKey %v, got nil", partitionKey)
	} else {
		contTokenPKStr, contTokenPKErr := contToken.PartitionKey.toJsonString()
		partitionKeyStr, partitionKeyErr := partitionKey.toJsonString()
		if contTokenPKErr != nil || partitionKeyErr != nil || contTokenPKStr != partitionKeyStr {
			t.Errorf("Expected PartitionKey %v, got %v", partitionKey, contToken.PartitionKey)
		}
	}

	if contToken.Continuation == nil {
		t.Fatal("Continuation should not be nil")
	}

	expectedEtag := azcore.ETag("\"etag-pk-123\"")
	if *contToken.Continuation != expectedEtag {
		t.Errorf("incorrect continuation token: got %q, want %q", *contToken.Continuation, expectedEtag)
	}

	// Verify we can use the token in a subsequent request
	srv.SetResponse(
		mock.WithBody([]byte(`{"_rid": "test-resource-id", "Documents": [], "_count": 0}`)),
		mock.WithHeader(cosmosHeaderEtag, "\"etag-pk-456\""),
		mock.WithStatusCode(200))

	verifier = pipelineVerifier{}
	internalClient, _ = azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client = &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ = newDatabase("databaseId", client)
	container, _ = newContainer("containerId", database)

	// Create request with continuation token
	options2 := &ChangeFeedOptions{
		MaxItemCount: 10,
		Continuation: &resp.ContinuationToken,
	}

	resp2, err := container.GetChangeFeed(context.TODO(), options2)
	if err != nil {
		t.Fatalf("GetChangeFeed with continuation token failed: %v", err)
	}

	if resp2.Count != 0 {
		t.Errorf("unexpected document count in second response: got %d, want 0", resp2.Count)
	}

	// Verify the continuation token is properly sent in headers
	if len(verifier.requests) == 0 {
		t.Fatal("no requests were made in second call")
	}
	request2 := verifier.requests[0]

	ifNoneMatchHeader := request2.headers.Get(headerIfNoneMatch)
	if ifNoneMatchHeader != "\"etag-pk-123\"" {
		t.Errorf("incorrect If-None-Match header: got %q, want %q", ifNoneMatchHeader, "\"etag-pk-123\"")
	}
}

func TestContainerGetChangeFeedForContainer(t *testing.T) {
	changeFeedBody := []byte(`{
        "_rid": "test-container-resource-id",
        "Documents": [{"id": "doc1"}, {"id": "doc2"}, {"id": "doc3"}],
        "_count": 3
    }`)

	srv, close := mock.NewTLSServer()
	defaultEndpoint, _ := url.Parse(srv.URL())
	defer close()

	// Set up mock response
	srv.SetResponse(
		mock.WithBody(changeFeedBody),
		mock.WithHeader(cosmosHeaderEtag, "\"etag-container-123\""),
		mock.WithHeader(cosmosHeaderActivityId, "containerChangeFeedActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "3.5"),
		mock.WithStatusCode(200))

	// Set up verifier to capture request details
	verifier := pipelineVerifier{}
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	// Create request without partition key or feed range (container level)
	options := &ChangeFeedOptions{
		MaxItemCount: 10,
	}

	resp, err := container.GetChangeFeed(context.TODO(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed for container failed: %v", err)
	}

	// Validate response content
	if resp.ResourceID != "test-container-resource-id" {
		t.Errorf("unexpected ResourceID: got %q, want %q", resp.ResourceID, "test-container-resource-id")
	}

	if resp.Count != 3 {
		t.Errorf("unexpected document count: got %d, want 3", resp.Count)
	}

	if len(resp.Documents) != 3 {
		t.Errorf("unexpected number of documents: got %d, want 3", len(resp.Documents))
	}

	// Validate request
	if len(verifier.requests) == 0 {
		t.Fatal("no requests were made")
	}
	request := verifier.requests[0]

	// Verify the partition key is not present in the header
	partitionKeyHeader := request.headers.Get(cosmosHeaderPartitionKey)
	if partitionKeyHeader != "" {
		t.Errorf("unexpected partition key header: should be empty but got %q", partitionKeyHeader)
	}

	// Verify change feed header is present
	changeFeedHeader := request.headers.Get(cosmosHeaderChangeFeed)
	if changeFeedHeader != cosmosHeaderValuesChangeFeed {
		t.Errorf("incorrect change feed header: got %q, want %q", changeFeedHeader, cosmosHeaderValuesChangeFeed)
	}

	// Verify max item count header
	maxItemCountHeader := request.headers.Get(cosmosHeaderMaxItemCount)
	if maxItemCountHeader != "10" {
		t.Errorf("incorrect max item count header: got %q, want %q", maxItemCountHeader, "10")
	}

	// Verify URL path is correct
	expectedPath := "/dbs/databaseId/colls/containerId/docs"
	if !strings.Contains(request.url.Path, expectedPath) {
		t.Errorf("expected URL path to contain %q, got %q", expectedPath, request.url.Path)
	}

	// Verify continuation token behavior - should be empty for first request
	if resp.ContinuationToken != "" {
		t.Errorf("expected ContinuationToken to be empty for container-level request, got %q", resp.ContinuationToken)
	}

	// Test with ETag as continuation token
	srv.SetResponse(
		mock.WithBody([]byte(`{"_rid": "test-container-resource-id", "Documents": [], "_count": 0}`)),
		mock.WithStatusCode(200))

	verifier = pipelineVerifier{}
	internalClient, _ = azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{PerCall: []policy.Policy{&verifier}}, &policy.ClientOptions{Transport: srv})
	client = &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ = newDatabase("databaseId", client)
	container, _ = newContainer("containerId", database)

	// Use the ETag as continuation token
	etag := string(resp.ETag)
	options2 := &ChangeFeedOptions{
		MaxItemCount: 10,
		Continuation: &etag,
	}

	resp2, err := container.GetChangeFeed(context.TODO(), options2)
	if err != nil {
		t.Fatalf("GetChangeFeed with ETag continuation token failed: %v", err)
	}

	if resp2.Count != 0 {
		t.Errorf("unexpected document count in second response: got %d, want 0", resp2.Count)
	}

	// Verify the If-None-Match header was set correctly
	if len(verifier.requests) == 0 {
		t.Fatal("no requests were made in second call")
	}
	request2 := verifier.requests[0]

	ifNoneMatchHeader := request2.headers.Get(headerIfNoneMatch)
	if ifNoneMatchHeader != etag {
		t.Errorf("incorrect If-None-Match header: got %q, want %q", ifNoneMatchHeader, etag)
	}
}
