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

func TestDatabaseResponseParsing(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := &CosmosDatabaseProperties{
		Id:           "someId",
		ETag:         "someEtag",
		SelfLink:     "someSelfLink",
		ResourceId:   "someResourceId",
		LastModified: &now,
	}

	database := &CosmosDatabase{
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
	parsedResponse, err := newCosmosDatabaseResponse(resp, database)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	if parsedResponse.DatabaseProperties == nil {
		t.Fatal("parsedResponse.DatabaseProperties is nil")
	}

	if properties.Id != parsedResponse.DatabaseProperties.Id {
		t.Errorf("Expected properties.Id to be %s, but got %s", properties.Id, parsedResponse.DatabaseProperties.Id)
	}

	if properties.ETag != parsedResponse.DatabaseProperties.ETag {
		t.Errorf("Expected properties.ETag to be %s, but got %s", properties.ETag, parsedResponse.DatabaseProperties.ETag)
	}

	if properties.SelfLink != parsedResponse.DatabaseProperties.SelfLink {
		t.Errorf("Expected properties.SelfLink to be %s, but got %s", properties.SelfLink, parsedResponse.DatabaseProperties.SelfLink)
	}

	if properties.ResourceId != parsedResponse.DatabaseProperties.ResourceId {
		t.Errorf("Expected properties.ResourceId to be %s, but got %s", properties.ResourceId, parsedResponse.DatabaseProperties.ResourceId)
	}

	if properties.LastModified.Time != parsedResponse.DatabaseProperties.LastModified.Time {
		t.Errorf("Expected properties.LastModified.Time to be %s, but got %s", properties.LastModified.Time.UTC(), parsedResponse.DatabaseProperties.LastModified.Time.UTC())
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

	if parsedResponse.DatabaseProperties.Database != database {
		t.Errorf("Expected database to be %v, but got %v", database, parsedResponse.DatabaseProperties.Database)
	}
}
