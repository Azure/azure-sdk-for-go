//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

const workspaceID = "d2d0e126-fa1e-4b0a-b647-250cdd471e68"

func TestExecute_BasicQuerySuccess(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	client := NewClient(cred, nil)

	query := "search * | take 5"
	time := "2022-08-01/2022-08-02"

	body := Body{
		Query:    &query,
		Timespan: &time,
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

	client := NewClient(cred, nil)

	var strPointer = new(string)
	*strPointer = "not a valid query"
	body := Body{
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

func TestExecute_AdvancedQuerySuccess(t *testing.T) {
	// special options: timeout, multiple workspaces, statistics, visualization
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	client := NewClient(cred, nil)

	query := "search * | take 5"
	body := Body{
		Query: &query,
	}

	prefer := "wait=180,include-statistics=true,include-render=true"
	options := &ClientQueryWorkspaceOptions{Prefer: &prefer}

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

	client := NewClient(cred, nil)
	query := "search * | take 5"
	time := "2022-08-01/2022-08-02"
	id := "189285912908589803580198308859812"
	workspaceID := "d2d0e126-fa1e-4b0a-b647-250cdd471e68"

	body := Body{
		Query:    &query,
		Timespan: &time,
	}

	req1 := BatchQueryRequest{Body: &body, ID: &id, Workspace: &workspaceID}

	batchRequest := BatchRequest{[]*BatchQueryRequest{&req1}}

	res, err := client.Batch(context.Background(), batchRequest, nil)
	_ = res
}

func TestBatch_QueryFailure(t *testing.T) {

}
