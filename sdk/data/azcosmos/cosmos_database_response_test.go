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

func TestDatabaseResponseParsing(t *testing.T) {
	nowAsUnix := time.Unix(time.Now().Unix(), 0)

	etag := azcore.ETag("someETag")
	properties := DatabaseProperties{
		ID:           "someId",
		ETag:         &etag,
		SelfLink:     "someSelfLink",
		ResourceID:   "someResourceId",
		LastModified: nowAsUnix,
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
	parsedResponse, err := newDatabaseResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if parsedResponse.DatabaseProperties == nil {
		t.Fatal("parsedResponse.DatabaseProperties is nil")
	}

	if properties.ID != parsedResponse.DatabaseProperties.ID {
		t.Errorf("Expected properties.Id to be %s, but got %s", properties.ID, parsedResponse.DatabaseProperties.ID)
	}

	if *properties.ETag != *parsedResponse.DatabaseProperties.ETag {
		t.Errorf("Expected properties.ETag to be %s, but got %s", *properties.ETag, *parsedResponse.DatabaseProperties.ETag)
	}

	if properties.SelfLink != parsedResponse.DatabaseProperties.SelfLink {
		t.Errorf("Expected properties.SelfLink to be %s, but got %s", properties.SelfLink, parsedResponse.DatabaseProperties.SelfLink)
	}

	if properties.ResourceID != parsedResponse.DatabaseProperties.ResourceID {
		t.Errorf("Expected properties.ResourceId to be %s, but got %s", properties.ResourceID, parsedResponse.DatabaseProperties.ResourceID)
	}

	if properties.LastModified != parsedResponse.DatabaseProperties.LastModified {
		t.Errorf("Expected properties.LastModified.Time to be %v, but got %v", properties.LastModified, parsedResponse.DatabaseProperties.LastModified)
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
