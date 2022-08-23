//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

const workspaceID = "d2d0e126-fa1e-4b0a-b647-250cdd471e68"

func TestQueryWorkspace_BasicQuerySuccess(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	client := azquery.NewClient(cred, nil)

	query := "search * | take 5"
	timespan := azquery.QueryTimeInterval(time.Now().Add(-12*time.Hour), time.Now())

	body := azquery.Body{
		Query:    &query,
		Timespan: &timespan,
	}

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}

	// test for correctness
	if res.Results.Error != nil {
		t.Fatal("expended Error to be nil")
	}
	if res.Results.Render != nil {
		t.Fatal("expended Render to be nil")
	}
	if res.Results.Statistics != nil {
		t.Fatal("expended Statistics to be nil")
	}
	if *res.Results.Tables[0].Name != "PrimaryResult" {
		t.Fatal("expended table name to be PrimaryResult")
	}
	if len(res.Results.Tables) != 1 {
		t.Fatal("expected one table")
	}
	if len(res.Results.Tables[0].Rows) != 5 {
		t.Fatal("expected 5 rows")
	}
}

func TestExecute_BasicQueryFailure(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	client := azquery.NewClient(cred, nil)

	var strPointer = new(string)
	*strPointer = "not a valid query"
	body := azquery.Body{
		Query: strPointer,
	}

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, nil)
	if err == nil {
		t.Fatalf("expected BadArgumentError")
	}
	if res.Results.Tables != nil {
		t.Fatalf("expected no results")
	}
}

func TestQueryWorkspace_AdvancedQuerySuccess(t *testing.T) {
	// special options: timeout, multiple workspaces, statistics, visualization
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	client := azquery.NewClient(cred, nil)

	query := "search * | take 5"
	body := azquery.Body{
		Query: &query,
	}

	prefer := "wait=180,include-statistics=true,include-render=true"
	options := &azquery.ClientQueryWorkspaceOptions{Prefer: &prefer}

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, options)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}

	if res.Results.Error != nil {
		t.Fatal("expended Error to be nil")
	}
	if res.Results.Render == nil {
		t.Fatal("expended Render have value")
	}
	if res.Results.Statistics == nil {
		t.Fatal("expended Statistics have value")
	}
}

// query with partial correctness??

// batch query tests
func TestBatch_QuerySuccess(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	client := azquery.NewClient(cred, nil)
	query := "search * | take 5"
	id := "1"
	workspaceID := "d2d0e126-fa1e-4b0a-b647-250cdd471e68"

	body := azquery.Body{
		Query: &query,
	}

	path := azquery.BatchQueryRequestPathQuery
	method := azquery.BatchQueryRequestMethodPOST

	req1 := azquery.BatchQueryRequest{Body: &body, ID: &id, Workspace: &workspaceID, Path: &path, Method: &method}

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{&req1}}

	res, err := client.Batch(context.Background(), batchRequest, &azquery.ClientBatchOptions{})
	if err != nil {
		t.Fatalf("expected non nil error: %s", err.Error())
	}
	if *(res.Responses[0].ID) != "1" {
		t.Fatalf("error, wrong response id")
	}
}

func TestBatch_QueryFailure(t *testing.T) {

}
