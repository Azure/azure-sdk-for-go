// +build !emulator

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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

func TestThroughputResponseParsing(t *testing.T) {
	nowAsUnix := time.Now().Unix()

	now := UnixTime{
		Time: time.Unix(nowAsUnix, 0),
	}

	properties := NewManualThroughputProperties(400)
	properties.offerId = "HFln"
	properties.offerResourceId = "4SRTANCD3Dw="
	properties.ETag = "\"00000000-0000-0000-9b8c-8ea3e19601d7\""
	properties.LastModified = &now

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

	req, err := azcore.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azcore.NewPipeline(srv)
	resp, _ := pl.Do(req)
	parsedResponse, err := newThroughputResponse(resp)
	if err != nil {
		t.Fatal(err)
	}

	if parsedResponse.RawResponse == nil {
		t.Fatal("parsedResponse.RawResponse is nil")
	}

	otherProperties := parsedResponse.ThroughputProperties

	if otherProperties.offerType != "" {
		t.Errorf("OfferType mismatch %v", otherProperties.offerType)
	}

	if otherProperties.offerResourceId != "4SRTANCD3Dw=" {
		t.Errorf("OfferResourceId mismatch %v", otherProperties.offerResourceId)
	}

	if otherProperties.version != offerVersion2 {
		t.Errorf("OfferVersion mismatch %v", otherProperties.version)
	}

	if otherProperties.offerId != "HFln" {
		t.Errorf("OfferId mismatch %v", otherProperties.offerId)
	}

	if otherProperties.ETag != "\"00000000-0000-0000-9b8c-8ea3e19601d7\"" {
		t.Errorf("Etag mismatch %v", otherProperties.ETag)
	}

	if otherProperties.LastModified.Time != properties.LastModified.Time {
		t.Errorf("Timestamp mismatch %v", otherProperties.LastModified.Time)
	}

	mt, err := otherProperties.ManualThroughput()
	if err != nil {
		t.Fatal(err)
	}

	if mt != 400 {
		t.Errorf("ManualThroughput mismatch %v", mt)
	}
}
