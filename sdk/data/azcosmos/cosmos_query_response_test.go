// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestQueryResponseParsing(t *testing.T) {
	queryResponseRaw := map[string][]map[string]string{
		"Documents": {
			{"id": "id1", "name": "name"},
			{"id": "id2", "name": "name"},
		},
	}

	jsonString, err := json.Marshal(queryResponseRaw)
	if err != nil {
		t.Fatal(err)
	}

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderQueryMetrics, "someQueryMetrics"),
		mock.WithHeader(cosmosHeaderIndexUtilization, "indexUtilization"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)
	parsedResponse, err := newQueryResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
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

	if *parsedResponse.QueryMetrics != "someQueryMetrics" {
		t.Errorf("Expected IndexMetrics to be %s, but got %s", "someQueryMetrics", *parsedResponse.IndexMetrics)
	}

	if *parsedResponse.IndexMetrics != "indexUtilization" {
		t.Errorf("Expected IndexUtilization to be %s, but got %s", "indexUtilization", *parsedResponse.IndexMetrics)
	}

	if len(parsedResponse.Items) != 2 {
		t.Errorf("Expected 2 documents, but got %d", len(parsedResponse.Items))
	}

	for index, item := range parsedResponse.Items {
		var itemResponseBody map[string]interface{}
		err = json.Unmarshal(item, &itemResponseBody)
		if err != nil {
			t.Fatalf("Failed to unmarshal item response: %v", err)
		}

		if itemResponseBody["id"] != ("id" + strconv.Itoa(index+1)) {
			t.Errorf("Expected id to be %s, but got %s", "id"+strconv.Itoa(index+1), itemResponseBody["id"])
		}

		if itemResponseBody["name"] != "name" {
			t.Errorf("Expected name to be %s, but got %s", "name", itemResponseBody["name"])
		}
	}
}

func TestQueryResponseValueParsing(t *testing.T) {
	queryResponseRaw := map[string][]string{
		"Documents": {"id1", "id2"},
	}

	jsonString, err := json.Marshal(queryResponseRaw)
	if err != nil {
		t.Fatal(err)
	}

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderQueryMetrics, "someQueryMetrics"),
		mock.WithHeader(cosmosHeaderIndexUtilization, "indexUtilization"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)
	parsedResponse, err := newQueryResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
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

	if *parsedResponse.QueryMetrics != "someQueryMetrics" {
		t.Errorf("Expected IndexMetrics to be %s, but got %s", "someQueryMetrics", *parsedResponse.IndexMetrics)
	}

	if *parsedResponse.IndexMetrics != "indexUtilization" {
		t.Errorf("Expected IndexUtilization to be %s, but got %s", "indexUtilization", *parsedResponse.IndexMetrics)
	}

	if len(parsedResponse.Items) != 2 {
		t.Errorf("Expected 2 documents, but got %d", len(parsedResponse.Items))
	}

	for index, item := range parsedResponse.Items {
		var itemResponseBody string
		err = json.Unmarshal(item, &itemResponseBody)
		if err != nil {
			t.Fatalf("Failed to unmarshal item response: %v", err)
		}

		if itemResponseBody != ("id" + strconv.Itoa(index+1)) {
			t.Errorf("Expected id to be %s, but got %s", "id"+strconv.Itoa(index+1), itemResponseBody)
		}
	}
}

func TestQueryContainersResponseParsing(t *testing.T) {
	queryResponseRaw := map[string][]map[string]string{
		"DocumentCollections": {
			{"id": "id1"},
			{"id": "id2"},
		},
	}

	jsonString, err := json.Marshal(queryResponseRaw)
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

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)
	parsedResponse, err := newContainersQueryResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
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

	if len(parsedResponse.Containers) != 2 {
		t.Errorf("Expected 2 containers, but got %d", len(parsedResponse.Containers))
	}

	for index, container := range parsedResponse.Containers {
		if err != nil {
			t.Fatalf("Failed to unmarshal item response: %v", err)
		}

		if container.ID != ("id" + strconv.Itoa(index+1)) {
			t.Errorf("Expected id to be %s, but got %s", "id"+strconv.Itoa(index+1), container.ID)
		}
	}
}

func TestQueryDatabasesResponseParsing(t *testing.T) {
	queryResponseRaw := map[string][]map[string]string{
		"Databases": {
			{"id": "id1"},
			{"id": "id2"},
		},
	}

	jsonString, err := json.Marshal(queryResponseRaw)
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

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)
	parsedResponse, err := newDatabasesQueryResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
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

	if len(parsedResponse.Databases) != 2 {
		t.Errorf("Expected 2 containers, but got %d", len(parsedResponse.Databases))
	}

	for index, db := range parsedResponse.Databases {
		if err != nil {
			t.Fatalf("Failed to unmarshal item response: %v", err)
		}

		if db.ID != ("id" + strconv.Itoa(index+1)) {
			t.Errorf("Expected id to be %s, but got %s", "id"+strconv.Itoa(index+1), db.ID)
		}
	}
}
