//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"
)

var query string = "let dt = datatable (DateTime: datetime, Bool:bool, Guid: guid, Int: int, Long:long, Double: double, String: string, Timespan: timespan, Decimal: decimal, Dynamic: dynamic)\n" + "[datetime(2015-12-31 23:59:59.9), false, guid(74be27de-1e4e-49d9-b579-fe0b331d3642), 12345, 1, 12345.6789, 'string value', 10s, decimal(0.10101), dynamic({\"a\":123, \"b\":\"hello\", \"c\":[1,2,3], \"d\":{}})];" + "range x from 1 to 100 step 1 | extend y=1 | join kind=fullouter dt on $left.y == $right.Long"

func TestQueryWorkspace_BasicQuerySuccess(t *testing.T) {
	client := startLogsTest(t)

	body := azquery.Body{
		Query:    to.Ptr(query),
		Timespan: to.Ptr("2015-12-31/2016-01-01"),
	}
	testSerde(t, &body)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}

	if res.Error != nil {
		t.Fatal("expended Error to be nil")
	}
	if res.Render != nil {
		t.Fatal("expended Render to be nil")
	}
	if res.Statistics != nil {
		t.Fatal("expended Statistics to be nil")
	}
	if len(res.Tables) != 1 {
		t.Fatal("expected one table")
	}
	if len(res.Tables[0].Rows) != 100 {
		t.Fatal("expected 100 rows")
	}

	testSerde(t, &res)
}

func TestQueryWorkspace_BasicQueryFailure(t *testing.T) {
	client := startLogsTest(t)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, azquery.Body{Query: to.Ptr("not a valid query")}, nil)
	if err == nil {
		t.Fatalf("expected an error")
	}
	if res.Tables != nil {
		t.Fatalf("expected no results")
	}

	var httpErr *azcore.ResponseError
	if !errors.As(err, &httpErr) {
		t.Fatal("expected an azcore.ResponseError")
	}
	if httpErr.ErrorCode != "BadArgumentError" {
		t.Fatal("expected a BadArgumentError")
	}
	if httpErr.StatusCode != 400 {
		t.Fatal("expected a 400 error")
	}

	testSerde(t, &res)
}

func TestQueryWorkspace_PartialError(t *testing.T) {
	client := startLogsTest(t)
	query := "let Weight = 92233720368547758; range x from 1 to 3 step 1 | summarize percentilesw(x, Weight * 100, 50)"

	res, err := client.QueryWorkspace(context.Background(), workspaceID, azquery.Body{Query: &query}, nil)
	if err != nil {
		t.Fatal("error with query")
	}
	if *res.Error.Code != "PartialError" {
		t.Fatal("expected a partial error")
	}
	if res.Error.RawMessage == nil {
		t.Fatal("expected rawmessage to contain error info")
	}

	testSerde(t, &res)
}

// tests for special options: timeout, statistics, visualization
func TestQueryWorkspace_AdvancedQuerySuccess(t *testing.T) {
	client := startLogsTest(t)
	query := query
	body := azquery.Body{
		Query: &query,
	}
	prefer := "wait=180,include-statistics=true,include-render=true"
	options := &azquery.LogsClientQueryWorkspaceOptions{Prefer: &prefer}

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, options)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}
	if res.Tables == nil {
		t.Fatal("expected Tables results")
	}
	if res.Error != nil {
		t.Fatal("expended Error to be nil")
	}
	if res.Render == nil {
		t.Fatal("expended Render results")
	}
	if res.Statistics == nil {
		t.Fatal("expended Statistics results")
	}
}

func TestQueryWorkspace_MultipleWorkspaces(t *testing.T) {
	client := startLogsTest(t)
	workspaces := []*string{&workspaceID2}
	body := azquery.Body{
		Query:      &query,
		Workspaces: workspaces,
	}
	testSerde(t, &body)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err.Error())
	}
	if res.Error != nil {
		t.Fatal("result error should be nil")
	}
	if len(res.Tables[0].Rows) != 100 {
		t.Fatalf("expected 100 results, received")
	}
}

func TestBatch_QuerySuccess(t *testing.T) {
	client := startLogsTest(t)
	query1, query2 := query, query+" | take 2"

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr(query1)}, ID: to.Ptr("1"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(query2)}, ID: to.Ptr("2"), Workspace: to.Ptr(workspaceID)},
	}}
	testSerde(t, &batchRequest)

	res, err := client.Batch(context.Background(), batchRequest, nil)
	if err != nil {
		t.Fatalf("expected non nil error: %s", err.Error())
	}
	if len(res.Responses) != 2 {
		t.Fatal("expected two responses")
	}
	for _, resp := range res.Responses {
		if resp.Body.Error != nil {
			t.Fatal("expected a successful response")
		}
	}
	testSerde(t, &res)
}

func TestBatch_PartialError(t *testing.T) {
	client := startLogsTest(t)

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr("not a valid query")}, ID: to.Ptr("1"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(query)}, ID: to.Ptr("2"), Workspace: to.Ptr(workspaceID)},
	}}

	res, err := client.Batch(context.Background(), batchRequest, nil)
	if err != nil {
		t.Fatalf("expected non nil error: %s", err.Error())
	}
	if len(res.Responses) != 2 {
		t.Fatal("expected two responses")
	}
	for _, resp := range res.Responses {
		if *resp.ID == "1" && resp.Body.Error == nil {
			t.Fatal("expected batch request 1 to fail")
		}
		if *resp.ID == "2" && resp.Body.Error != nil {
			t.Fatal("expected batch request 2 to succeed")
		}
	}
}

func TestLogConstants(t *testing.T) {
	batchMethod := []azquery.BatchQueryRequestMethod{azquery.BatchQueryRequestMethodPOST}
	batchMethodRes := azquery.PossibleBatchQueryRequestMethodValues()
	require.Equal(t, batchMethod, batchMethodRes)

	batchPath := []azquery.BatchQueryRequestPath{azquery.BatchQueryRequestPathQuery}
	batchPathRes := azquery.PossibleBatchQueryRequestPathValues()
	require.Equal(t, batchPath, batchPathRes)

	logsColumnType := []azquery.LogsColumnType{azquery.LogsColumnTypeBool, azquery.LogsColumnTypeDatetime, azquery.LogsColumnTypeDecimal, azquery.LogsColumnTypeDynamic, azquery.LogsColumnTypeGUID, azquery.LogsColumnTypeInt, azquery.LogsColumnTypeLong, azquery.LogsColumnTypeReal, azquery.LogsColumnTypeString, azquery.LogsColumnTypeTimespan}
	logsColumnTypeRes := azquery.PossibleLogsColumnTypeValues()
	require.Equal(t, logsColumnType, logsColumnTypeRes)
}
