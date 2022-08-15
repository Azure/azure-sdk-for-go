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
	_ = res
}

func TestExecute_BasicQueryFailure(t *testing.T) {

}

func TestExecute_AdvancedQuerySuccess(t *testing.T) {
	// special options: timeout, multiple workspaces, statistics, visualization
}

// query with partial correctness??

// batch query tests
func TestBatch_QuerySuccess(t *testing.T) {

}

func TestBatch_QueryFailure(t *testing.T) {

}
