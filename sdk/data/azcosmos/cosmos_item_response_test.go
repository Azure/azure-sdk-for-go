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

func TestItemResponseParsing(t *testing.T) {
	properties := map[string]string{
		"id":   "id",
		"name": "name",
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
	parsedResponse, err := newItemResponse(resp)
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

	if string(parsedResponse.Value) != string(jsonString) {
		t.Errorf("Expected Value to be %s, but got %s", string(jsonString), string(parsedResponse.Value))
	}
}
