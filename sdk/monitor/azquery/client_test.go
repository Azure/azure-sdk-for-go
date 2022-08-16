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

	var strPointer = new(string)
	*strPointer = "let dt = datatable (DateTime: datetime, Bool:bool, Guid: guid, Int: int, Long:long, Double: double, String: string, Timespan: timespan, Decimal: decimal, Dynamic: dynamic)\n" + "[datetime(2015-12-31 23:59:59.9), false, guid(74be27de-1e4e-49d9-b579-fe0b331d3642), 12345, 1, 12345.6789, 'string value', 10s, decimal(0.10101), dynamic({\"a\":123, \"b\":\"hello\", \"c\":[1,2,3], \"d\":{}})];" + "range x from 1 to 100 step 1 | extend y=1 | join kind=fullouter dt on $left.y == $right.Long"
	body := Body{
		Query: strPointer,
	}

	res, err := client.Execute(context.Background(), workspaceID, body, nil)
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
	if len(res.Results.Tables[0].Rows) != 100 {
		t.Fatal("expected 100 rows")
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

	res, err := client.Execute(context.Background(), workspaceID, body, nil)
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

	var queryPointer = new(string)
	*queryPointer = "let dt = datatable (DateTime: datetime, Bool:bool, Guid: guid, Int: int, Long:long, Double: double, String: string, Timespan: timespan, Decimal: decimal, Dynamic: dynamic)\n" + "[datetime(2015-12-31 23:59:59.9), false, guid(74be27de-1e4e-49d9-b579-fe0b331d3642), 12345, 1, 12345.6789, 'string value', 10s, decimal(0.10101), dynamic({\"a\":123, \"b\":\"hello\", \"c\":[1,2,3], \"d\":{}})];" + "range x from 1 to 100 step 1 | extend y=1 | join kind=fullouter dt on $left.y == $right.Long"
	body := Body{
		Query: queryPointer,
	}
	var optionsPointer = new(string)
	*optionsPointer = "wait=180,include-statistics=true,include-render=true"
	options := &ClientExecuteOptions{Prefer: optionsPointer}

	res, err := client.Execute(context.Background(), workspaceID, body, options)
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

}

func TestBatch_QueryFailure(t *testing.T) {

}
