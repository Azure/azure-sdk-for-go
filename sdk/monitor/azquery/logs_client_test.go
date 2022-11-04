//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/stretchr/testify/require"
)

var query string = "let dt = datatable (DateTime: datetime, Bool:bool, Guid: guid, Int: int, Long:long, Double: double, String: string, Timespan: timespan, Decimal: decimal, Dynamic: dynamic)\n" + "[datetime(2015-12-31 23:59:59.9), false, guid(74be27de-1e4e-49d9-b579-fe0b331d3642), 12345, 1, 12345.6789, 'string value', 10s, decimal(0.10101), dynamic({\"a\":123, \"b\":\"hello\", \"c\":[1,2,3], \"d\":{}})];" + "range x from 1 to 100 step 1 | extend y=1 | join kind=fullouter dt on $left.y == $right.Long"

type queryTest struct {
	Bool   bool
	Long   int64
	String string
}

func TestQueryWorkspace_BasicQuerySuccess(t *testing.T) {
	client := startLogsTest(t)
	body := azquery.Body{
		Query:    to.Ptr(query),
		Timespan: to.Ptr("2015-12-31/2016-01-01"),
	}
	testSerde(t, &body)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err)
	}
	if res.Error != nil {
		t.Fatalf("expected Error to be nil: %s", res.Error)
	}
	if res.Render != nil {
		t.Fatal("expected Render to be nil")
	}
	if res.Statistics != nil {
		t.Fatal("expected Statistics to be nil")
	}
	if len(res.Tables) != 1 {
		t.Fatal("expected one table")
	}
	if len(res.Tables[0].Rows) != 100 {
		t.Fatal("expected 100 rows")
	}

	var queryResults []queryTest
	for _, table := range res.Tables {
		queryResults = make([]queryTest, len(table.Rows))
		indexLong := table.ColumnIndexLookup["Long"]
		indexString := table.ColumnIndexLookup["String"]
		indexBool := table.ColumnIndexLookup["Bool"]

		for index, row := range table.Rows {
			queryResults[index] = queryTest{
				Long:   int64(row[indexLong].(float64)),
				String: row[indexString].(string),
				Bool:   row[indexBool].(bool),
			}
		}
	}

	if len(queryResults) != 100 {
		t.Fatal("expected 100 structs")
	}
	if queryResults[99].Bool != false {
		t.Fatal("expected Bool to be false")
	}
	if queryResults[99].String != "string value" {
		t.Fatal("expected String to be 'string value")
	}
	if queryResults[99].Long != 1 {
		t.Fatal("expected Long to be 1")
	}

	testSerde(t, &res)
}

func TestQueryWorkspace_BasicQueryFailure(t *testing.T) {
	client := startLogsTest(t)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, azquery.Body{Query: to.Ptr("not a valid query")}, nil)
	if err == nil {
		t.Fatalf("expected an error")
	}
	if res.Error != nil {
		t.Fatalf("expected Error to be nil: %s", res.Error)
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
		t.Fatalf("error with query: %s", err)
	}
	if res.Error == nil {
		t.Fatal("expected an error")
	}
	if res.Error.Code != "PartialError" {
		t.Fatal("expected a partial error")
	}
	if !strings.Contains(res.Error.Error(), "PartialError") {
		t.Fatal("expected error message to contain PartialError")
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
		t.Fatalf("error with query, %s", err)
	}
	if res.Error != nil {
		t.Fatalf("expected Error to be nil: %s", res.Error)
	}
	if res.Tables == nil {
		t.Fatal("expected Tables results")
	}
	if res.Render == nil {
		t.Fatal("expected Render results")
	}
	if res.Statistics == nil {
		t.Fatal("expected Statistics results")
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
		t.Fatalf("error with query, %s", err)
	}
	if res.Error != nil {
		t.Fatalf("expected Error to be nil: %s", res.Error)
	}
	if len(res.Tables[0].Rows) != 100 {
		t.Fatalf("expected 100 results, received")
	}
}

func TestQueryBatch_QuerySuccess(t *testing.T) {
	client := startLogsTest(t)
	query1, query2 := query, query+" | take 2"

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr(query1)}, ID: to.Ptr("1"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(query2)}, ID: to.Ptr("2"), Workspace: to.Ptr(workspaceID)},
	}}
	testSerde(t, &batchRequest)

	res, err := client.QueryBatch(context.Background(), batchRequest, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err)
	}
	if len(res.Responses) != 2 {
		t.Fatal("expected two responses")
	}
	for _, resp := range res.Responses {
		if resp.Body.Error != nil {
			t.Fatalf("expected Error to be nil: %s", resp.Body.Error)
		}
		if resp.Body.Tables == nil {
			t.Fatal("expected a response")
		}
		if *resp.ID == "1" && len(resp.Body.Tables[0].Rows) != 100 {
			t.Fatal("expected 100 rows from batch request 1")
		}
		if *resp.ID == "2" && len(resp.Body.Tables[0].Rows) != 2 {
			t.Fatal("expected 100 rows from batch request 1")
		}
	}
	testSerde(t, &res)
}

func TestQueryBatch_PartialError(t *testing.T) {
	client := startLogsTest(t)

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr("not a valid query")}, ID: to.Ptr("1"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(query)}, ID: to.Ptr("2"), Workspace: to.Ptr(workspaceID)},
	}}

	res, err := client.QueryBatch(context.Background(), batchRequest, nil)
	if err != nil {
		t.Fatalf("error with query, %s", err)
	}
	if len(res.Responses) != 2 {
		t.Fatal("expected two responses")
	}
	for _, resp := range res.Responses {
		if *resp.ID == "1" {
			if resp.Body.Error == nil {
				t.Fatal("expected batch request 1 to fail")
			}
			if resp.Body.Error.Code != "BadArgumentError" {
				t.Fatal("expected BadArgumentError")
			}
			if !strings.Contains(resp.Body.Error.Error(), "BadArgumentError") {
				t.Fatal("expected error message to contain BadArgumentError")
			}
		}
		if *resp.ID == "2" {
			if resp.Body.Error != nil {
				t.Fatalf("expected batch request 2 to succeed: %s", resp.Body.Error)
			}
			if len(resp.Body.Tables[0].Rows) != 100 {
				t.Fatal("expected 100 rows")
			}
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
