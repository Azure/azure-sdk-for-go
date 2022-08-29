//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"
)

const workspaceID1 = "d2d0e126-fa1e-4b0a-b647-250cdd471e68"
const workspaceID2 = "9dad0092-fd13-403a-b367-a189a090a541"
const query = "let dt = datatable (DateTime: datetime, Bool:bool, Guid: guid, Int: int, Long:long, Double: double, String: string, Timespan: timespan, Decimal: decimal, Dynamic: dynamic)\n" + "[datetime(2015-12-31 23:59:59.9), false, guid(74be27de-1e4e-49d9-b579-fe0b331d3642), 12345, 1, 12345.6789, 'string value', 10s, decimal(0.10101), dynamic({\"a\":123, \"b\":\"hello\", \"c\":[1,2,3], \"d\":{}})];" + "range x from 1 to 100 step 1 | extend y=1 | join kind=fullouter dt on $left.y == $right.Long"

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)
}

func getClient(t *testing.T) *azquery.LogsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatal("error constructing credential")
	}

	return azquery.NewLogsClient(cred, nil)
}

func TestQueryWorkspace_BasicQuerySuccess(t *testing.T) {
	client := getClient(t)
	query := query
	timespan := azquery.QueryTimeInterval(time.Now().Add(-12*time.Hour), time.Now())
	body := azquery.Body{
		Query:    &query,
		Timespan: &timespan,
	}
	testSerde(t, &body)

	res, err := client.QueryWorkspace(context.Background(), workspaceID1, body, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}

	if res.Results.Error != nil {
		t.Fatal("expended Error to be nil")
	}
	if res.Results.Render != nil {
		t.Fatal("expended Render to be nil")
	}
	if res.Results.Statistics != nil {
		t.Fatal("expended Statistics to be nil")
	}
	if len(res.Results.Tables) != 1 {
		t.Fatal("expected one table")
	}
	if len(res.Results.Tables[0].Rows) != 100 {
		t.Fatal("expected 100 rows")
	}

	testSerde(t, &res.Results)
}

func TestQueryWorkspace_BasicQueryFailure(t *testing.T) {
	client := getClient(t)
	query := "not a valid query"
	body := azquery.Body{
		Query: &query,
	}

	res, err := client.QueryWorkspace(context.Background(), workspaceID1, body, nil)
	if err == nil {
		t.Fatalf("expected BadArgumentError")
	}
	if res.Results.Tables != nil {
		t.Fatalf("expected no results")
	}
	testSerde(t, &res.Results)
}

func TestQueryWorkspace_PartialError(t *testing.T) {
	client := getClient(t)
	query := "let Weight = 92233720368547758; range x from 1 to 3 step 1 | summarize percentilesw(x, Weight * 100, 50)"
	body := azquery.Body{
		Query: &query,
	}

	res, err := client.QueryWorkspace(context.Background(), workspaceID1, body, nil)
	if err != nil {
		t.Fatal("error with query")
	}
	if *res.Results.Error.Code != "PartialError" {
		t.Fatal("expected a partial error")
	}

	testSerde(t, &res.Results)
}

// tests for special options: timeout, statistics, visualization
func TestQueryWorkspace_AdvancedQuerySuccess(t *testing.T) {
	client := getClient(t)
	query := query
	body := azquery.Body{
		Query: &query,
	}
	prefer := "wait=180,include-statistics=true,include-render=true"
	options := &azquery.LogsClientQueryWorkspaceOptions{Prefer: &prefer}

	res, err := client.QueryWorkspace(context.Background(), workspaceID1, body, options)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}
	if res.Results.Tables == nil {
		t.Fatal("expected Tables results")
	}
	if res.Results.Error != nil {
		t.Fatal("expended Error to be nil")
	}
	if res.Results.Render == nil {
		t.Fatal("expended Render results")
	}
	if res.Results.Statistics == nil {
		t.Fatal("expended Statistics results")
	}
}

func TestQueryWorkspace_MultipleWorkspaces(t *testing.T) {
	client := getClient(t)
	query := "union * | where TimeGenerated > ago(100d) | project TenantId | summarize count() by TenantId"
	workspaceID2 := workspaceID2
	workspaces := []*string{&workspaceID2}
	body := azquery.Body{
		Query:      &query,
		Workspaces: workspaces,
	}
	testSerde(t, &body)

	res, err := client.QueryWorkspace(context.Background(), workspaceID1, body, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}
	if res.Results.Error != nil {
		t.Fatal("result error should be nil")
	}
	if len(res.Results.Tables[0].Rows) != 2 {
		t.Fatal("expected 2 results")
	}
}

//TODO fix path and method attributes
/*func TestBatch_QuerySuccess(t *testing.T) {
	client := getClient(t)
	query1, query2 := query, query+" | take 2"
	id1, id2 := "1", "2"
	workspaceID := workspaceID1
	body1 := azquery.Body{
		Query: &query1,
	}
	body2 := azquery.Body{
		Query: &query2,
	}
	req1 := azquery.BatchQueryRequest{Body: &body1, ID: &id1, Workspace: &workspaceID}
	req2 := azquery.BatchQueryRequest{Body: &body2, ID: &id2, Workspace: &workspaceID}
	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{&req1, &req2}}
	testSerde(t, &batchRequest)

	res, err := client.Batch(context.Background(), batchRequest, nil)
	if err != nil {
		t.Fatalf("expected non nil error: %s", err.Error())
	}
	if len(res.BatchResponse.Responses) != 2 {
		t.Fatal("expected two responses")
	}
	testSerde(t, &res.BatchResponse)
}*/

func TestBatch_QueryFailure(t *testing.T) {

}
