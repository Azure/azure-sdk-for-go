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

func TestNewChangeFeedResponse(t *testing.T) {
	jsonString := []byte(`{
		"_rid": "ubgwAI1+zvg=",
		"Documents": [
			{
				"id": "Erewhon",
				"license": "GHAS",
				"partitionKey": "33333",
				"_rid": "ubgwAI1+zvgDAAAAAAAAAA==",
				"_self": "dbs/ubgwAA==/colls/ubgwAI1+zvg=/docs/ubgwAI1+zvgDAAAAAAAAAA==/",
				"_etag": "\"e1015e15-0000-0700-0000-6859bda10000\"",
				"_attachments": "attachments/",
				"_ts": 1750711713
			},
			{
				"id": "TraderJoes",
				"license": "Copilot",
				"partitionKey": "44444",
				"_rid": "ubgwAI1+zvgBAAAAAAAACA==",
				"_self": "dbs/ubgwAA==/colls/ubgwAI1+zvg=/docs/ubgwAI1+zvgBAAAAAAAACA==/",
				"_etag": "\"9701c68b-0000-0700-0000-6859c38b0000\"",
				"_attachments": "attachments/",
				"_ts": 1750713227
			}
		],
		"_count": 2
	}`)

	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "someEtag"),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "13.42"),
		mock.WithHeader("Content-Type", "application/json"),
	)

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatal(err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	parsedResponse, err := newChangeFeedResponse(resp)
	if err != nil {
		t.Fatalf("newChangeFeedResponse error: %v", err)
	}

	if parsedResponse.Response.RawResponse == nil {
		t.Fatal("parsedResponse.Response.RawResponse is nil")
	}

	if parsedResponse.ResourceID != "ubgwAI1+zvg=" {
		t.Fatalf("unexpected ResourceID: got %q, want %q", parsedResponse.ResourceID, "ubgwAI1+zvg=")
	}

	if parsedResponse.Count != 2 {
		t.Fatalf("unexpected Count: got %d, want 2", parsedResponse.Count)
	}

	if len(parsedResponse.Documents) != 2 {
		t.Fatalf("unexpected number of Documents: got %d, want 2", len(parsedResponse.Documents))
	}

	var doc0, doc1 map[string]interface{}
	if err := json.Unmarshal(parsedResponse.Documents[0], &doc0); err != nil {
		t.Fatalf("failed to unmarshal first document: %v", err)
	}
	if doc0["id"] != "Erewhon" {
		t.Errorf("unexpected first document id: got %v, want Erewhon", doc0["id"])
	}

	if err := json.Unmarshal(parsedResponse.Documents[1], &doc1); err != nil {
		t.Fatalf("failed to unmarshal second document: %v", err)
	}
	if doc1["id"] != "TraderJoes" {
		t.Errorf("unexpected second document id: got %v, want TraderJoes", doc1["id"])
	}
}

func TestChangeFeedResponseWithFeedRange(t *testing.T) {
	jsonString := []byte(`{
        "_rid": "testResourceId",
        "Documents": [{"id": "doc1"}],
        "_count": 1
    }`)

	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "\"test-etag-123\""),
		mock.WithHeader(cosmosHeaderActivityId, "someActivityId"),
		mock.WithHeader(cosmosHeaderRequestCharge, "2.5"),
		mock.WithStatusCode(200),
	)

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	parsedResponse, err := newChangeFeedResponse(resp)
	if err != nil {
		t.Fatalf("failed to create ChangeFeedResponse: %v", err)
	}

	parsedResponse.FeedRange = &FeedRange{
		MinInclusive: "00",
		MaxExclusive: "FF",
	}

	parsedResponse.PopulateCompositeContinuationToken()

	if parsedResponse.ContinuationToken == "" {
		t.Fatal("expected CompositeContinuationToken to be populated, but it was empty")
	}

	var compositeToken compositeContinuationToken
	err = json.Unmarshal([]byte(parsedResponse.ContinuationToken), &compositeToken)
	if err != nil {
		t.Fatalf("failed to unmarshal composite token: %v", err)
	}

	if compositeToken.Version != cosmosCompositeContinuationTokenVersion {
		t.Errorf("expected Version %d, got %d", cosmosCompositeContinuationTokenVersion, compositeToken.Version)
	}

	if compositeToken.ResourceID != "testResourceId" {
		t.Errorf("unexpected ResourceID: got %q, want %q", compositeToken.ResourceID, "testResourceId")
	}

	if len(compositeToken.Continuation) != 1 {
		t.Fatalf("unexpected number of continuation ranges: got %d, want 1", len(compositeToken.Continuation))
	}

	if compositeToken.Continuation[0].MinInclusive != "00" {
		t.Errorf("unexpected MinInclusive: got %q, want %q", compositeToken.Continuation[0].MinInclusive, "00")
	}

	if compositeToken.Continuation[0].MaxExclusive != "FF" {
		t.Errorf("unexpected MaxExclusive: got %q, want %q", compositeToken.Continuation[0].MaxExclusive, "FF")
	}

	if compositeToken.Continuation[0].ContinuationToken == nil {
		t.Fatal("expected ContinuationToken to be set, but it was nil")
	}

	if *compositeToken.Continuation[0].ContinuationToken != azcore.ETag("\"test-etag-123\"") {
		t.Errorf("unexpected ContinuationToken: got %q, want %q", *compositeToken.Continuation[0].ContinuationToken, "\"test-etag-123\"")
	}
}

func TestChangeFeedResponseWithoutFeedRange(t *testing.T) {
	jsonString := []byte(`{
        "_rid": "testResourceId",
        "Documents": [{"id": "doc1"}],
        "_count": 1
    }`)

	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()
	srv.SetResponse(
		mock.WithBody(jsonString),
		mock.WithHeader(cosmosHeaderEtag, "\"test-etag-456\""),
		mock.WithStatusCode(200),
	)

	req, err := azruntime.NewRequest(context.Background(), http.MethodGet, srv.URL())
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	resp, err := pl.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	parsedResponse, err := newChangeFeedResponse(resp)
	if err != nil {
		t.Fatalf("failed to create ChangeFeedResponse: %v", err)
	}

	parsedResponse.PopulateCompositeContinuationToken()

	if parsedResponse.ContinuationToken != "" {
		t.Errorf("expected CompositeContinuationToken to be empty, but got: %q", parsedResponse.ContinuationToken)
	}
}
