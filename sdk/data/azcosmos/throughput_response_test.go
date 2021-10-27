// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestThroughputResponseParsing(t *testing.T) {
	properties := NewManualThroughputProperties(400)

	etag := azcore.ETag("\"00000000-0000-0000-9b8c-8ea3e19601d7\"")

	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = &etag
	jsonString, err := json.Marshal(&properties)
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
	parsedResponse, err := newThroughputResponse(resp, nil)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if parsedResponse.ThroughputProperties == nil {
		t.Fatal("parsedResponse.ThroughputProperties is nil")
	}

	if parsedResponse.ThroughputProperties.offerId != properties.offerId {
		t.Fatalf("parsedResponse.ThroughputProperties.offerId is %s, expected %s", parsedResponse.ThroughputProperties.offerId, properties.offerId)
	}

	if parsedResponse.ThroughputProperties.offerResourceId != properties.offerResourceId {
		t.Fatalf("parsedResponse.ThroughputProperties.offerResourceId is %s, expected %s", parsedResponse.ThroughputProperties.offerResourceId, properties.offerResourceId)
	}

	if *parsedResponse.ThroughputProperties.ETag != *properties.ETag {
		t.Fatalf("parsedResponse.ThroughputProperties.ETag is %s, expected %s", *parsedResponse.ThroughputProperties.ETag, *properties.ETag)
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

func TestThroughputResponseParsingWithPreviousRU(t *testing.T) {
	var queryRequestCharge float32 = 10.0

	etag := azcore.ETag("\"00000000-0000-0000-9b8c-8ea3e19601d7\"")
	properties := NewManualThroughputProperties(400)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = &etag
	jsonString, err := json.Marshal(&properties)
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
	parsedResponse, err := newThroughputResponse(resp, &queryRequestCharge)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if parsedResponse.ThroughputProperties == nil {
		t.Fatal("parsedResponse.ThroughputProperties is nil")
	}

	if parsedResponse.ThroughputProperties.offerId != properties.offerId {
		t.Fatalf("parsedResponse.ThroughputProperties.offerId is %s, expected %s", parsedResponse.ThroughputProperties.offerId, properties.offerId)
	}

	if parsedResponse.ThroughputProperties.offerResourceId != properties.offerResourceId {
		t.Fatalf("parsedResponse.ThroughputProperties.offerResourceId is %s, expected %s", parsedResponse.ThroughputProperties.offerResourceId, properties.offerResourceId)
	}

	if *parsedResponse.ThroughputProperties.ETag != *properties.ETag {
		t.Fatalf("parsedResponse.ThroughputProperties.ETag is %s, expected %s", *parsedResponse.ThroughputProperties.ETag, *properties.ETag)
	}

	if parsedResponse.ActivityID != "someActivityId" {
		t.Errorf("Expected ActivityId to be %s, but got %s", "someActivityId", parsedResponse.ActivityID)
	}

	if parsedResponse.RequestCharge != 23.42 {
		t.Errorf("Expected RequestCharge to be %f, but got %f", 23.42, parsedResponse.RequestCharge)
	}

	if parsedResponse.ETag != "someEtag" {
		t.Errorf("Expected ETag to be %s, but got %s", "someEtag", parsedResponse.ETag)
	}
}
