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

func TestTransactionalBatchResponseParsing(t *testing.T) {
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

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)
	parsedResponse, err := newTransactionalBatchResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if !parsedResponse.Success {
		t.Errorf("Expected Success to be true, but got false")
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

	if len(parsedResponse.OperationResults) != 2 {
		t.Errorf("Expected 2 documents, but got %d", len(parsedResponse.OperationResults))
	}

	for index, item := range parsedResponse.OperationResults {
		if index == 0 && item.ETag != "someETag" {
			t.Errorf("Expected ETag to be %s, but got %s", "someETag", item.ETag)
		}

		if index == 1 && item.ETag != "someETag2" {
			t.Errorf("Expected ETag to be %s, but got %s", "someETag2", item.ETag)
		}

		if index == 0 && item.StatusCode != http.StatusOK {
			t.Errorf("Expected StatusCode to be %d, but got %d", http.StatusOK, item.StatusCode)
		}

		if index == 1 && item.StatusCode != http.StatusCreated {
			t.Errorf("Expected StatusCode to be %d, but got %d", http.StatusCreated, item.StatusCode)
		}

		if index == 0 && item.RequestCharge != 10 {
			t.Errorf("Expected RequestCharge to be %f, but got %f", 10.0, item.RequestCharge)
		}

		if index == 1 && item.RequestCharge != 11 {
			t.Errorf("Expected RequestCharge to be %f, but got %f", 11.0, item.RequestCharge)
		}

		if index == 0 && string(item.ResourceBody) != "\"someBody\"" {
			t.Errorf("Expected ResourceBody to be %s, but got %s", "someBody", item.ResourceBody)
		}

		if index == 1 && item.ResourceBody != nil {
			t.Errorf("Expected ResourceBody to be nil, but got %s", item.ResourceBody)
		}

	}
}

func TestTransactionalBatchResponseParsing_Failed(t *testing.T) {
	batchResponseRaw := []map[string]interface{}{
		{"statusCode": 424},
		{"statusCode": 409},
	}

	jsonString, err := json.Marshal(batchResponseRaw)
	if err != nil {
		t.Fatal(err)
	}

	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithStatusCode(http.StatusMultiStatus),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"))

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, _ := pl.Do(req)
	parsedResponse, err := newTransactionalBatchResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.Success {
		t.Errorf("Expected Success to be false, but got true")
	}
}
