// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azlogs_test

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azlogs"
	"github.com/stretchr/testify/require"
)

var query string = "let dt = datatable (DateTime: datetime, Bool:bool, Guid: guid, Int: int, Long:long, Double: double, String: string, Timespan: timespan, Decimal: decimal, Dynamic: dynamic)\n" + "[datetime(2015-12-31 23:59:59.9), false, guid(74be27de-1e4e-49d9-b579-fe0b331d3642), 12345, 1, 12345.6789, 'string value', 10s, decimal(0.10101), dynamic({\"a\":123, \"b\":\"hello\", \"c\":[1,2,3], \"d\":{}})];" + "range x from 1 to 100 step 1 | extend y=1 | join kind=fullouter dt on $left.y == $right.Long"

type queryTest struct {
	Bool   bool
	Long   int64
	String string
}

func TestClient(t *testing.T) {
	client, err := azlogs.NewClient(credential, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	c := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "https://...",
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			azlogs.ServiceName: {
				Audience: "",
				Endpoint: "",
			},
		},
	}
	opts := azcore.ClientOptions{Cloud: c}
	cloudClient, err := azlogs.NewClient(credential, &azlogs.ClientOptions{ClientOptions: opts})
	require.Error(t, err)
	require.Equal(t, err.Error(), "provided Cloud field is missing Azure Monitor Logs configuration")
	require.Nil(t, cloudClient)
}

func TestQueryWorkspace_BasicQuerySuccess(t *testing.T) {
	client := startTest(t)
	timespan := azlogs.NewTimeInterval(time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 2, 0, 0, 0, 0, time.UTC))
	body := azlogs.QueryBody{
		Query:    to.Ptr(query),
		Timespan: to.Ptr(timespan),
	}
	testSerde(t, &body)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, nil)
	require.NoError(t, err)
	require.Nil(t, res.Error)
	require.Nil(t, res.Visualization)
	require.Nil(t, res.Statistics)
	require.Len(t, res.Tables, 1)
	require.Len(t, res.Tables[0].Rows, 100)

	var queryResults []queryTest
	for _, table := range res.Tables {
		queryResults = make([]queryTest, len(table.Rows))

		for index, row := range table.Rows {
			queryResults[index] = queryTest{
				Long:   int64(row[6].(float64)),
				String: row[8].(string),
				Bool:   row[3].(bool),
			}
		}
	}

	require.Len(t, queryResults, 100)
	require.False(t, queryResults[99].Bool)
	require.Equal(t, queryResults[99].String, "string value")
	require.Equal(t, queryResults[99].Long, int64(1))

	testSerde(t, &res)
}

func TestQueryWorkspace_BasicQueryFailure(t *testing.T) {
	client := startTest(t)

	res, err := client.QueryWorkspace(
		context.Background(),
		workspaceID,
		azlogs.QueryBody{
			Query:    to.Ptr("not a valid query"),
			Timespan: to.Ptr(azlogs.TimeInterval("PT2H")),
		},
		nil,
	)
	require.Error(t, err)
	require.Nil(t, res.Error)
	require.Nil(t, res.Tables)

	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "BadArgumentError")
	require.Equal(t, httpErr.StatusCode, 400)

	testSerde(t, &res)
}

func TestQueryWorkspace_PartialError(t *testing.T) {
	client := startTest(t)
	query := "let Weight = 92233720368547758; range x from 1 to 3 step 1 | summarize percentilesw(x, Weight * 100, 50)"

	res, err := client.QueryWorkspace(context.Background(), workspaceID, azlogs.QueryBody{Query: &query}, nil)
	require.NoError(t, err)
	require.NotNil(t, res.Error)
	require.Equal(t, res.Error.Code, "PartialError")
	require.Contains(t, res.Error.Error(), "PartialError")

	testSerde(t, &res)
}

// tests for special options: timeout, statistics, visualization
func TestQueryWorkspace_AdvancedQuerySuccess(t *testing.T) {
	client := startTest(t)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, azlogs.QueryBody{Query: &query},
		&azlogs.QueryWorkspaceOptions{Options: &azlogs.QueryOptions{Statistics: to.Ptr(true), Visualization: to.Ptr(true), Wait: to.Ptr(600)}})
	require.NoError(t, err)
	require.Nil(t, res.Error)
	require.NotNil(t, res.Tables)
	require.NotNil(t, res.Visualization)
	require.NotNil(t, res.Statistics)
	testSerde(t, &res)
}

func TestQueryWorkspace_MultipleWorkspaces(t *testing.T) {
	client := startTest(t)
	workspaces := []string{workspaceID2}
	body := azlogs.QueryBody{
		Query:                &query,
		AdditionalWorkspaces: workspaces,
	}
	testSerde(t, &body)

	res, err := client.QueryWorkspace(context.Background(), workspaceID, body, nil)
	require.NoError(t, err)
	require.Nil(t, res.Error)
	require.Len(t, res.Tables[0].Rows, 100)
}

func TestQueryResource(t *testing.T) {
	client := startTest(t)
	timespan := azlogs.NewTimeInterval(time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 2, 0, 0, 0, 0, time.UTC))
	body := azlogs.QueryBody{
		Query:    to.Ptr(query),
		Timespan: to.Ptr(timespan),
	}
	testSerde(t, &body)

	res, err := client.QueryResource(context.Background(), resourceURI, body, nil)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Nil(t, res.Error)
	require.Nil(t, res.Visualization)
	require.Nil(t, res.Statistics)
	require.Len(t, res.Tables, 1)
	require.Len(t, res.Tables[0].Rows, 100)
	testSerde(t, &res)
}

func TestQueryResource_Fail(t *testing.T) {
	client := startTest(t)

	res, err := client.QueryResource(
		context.Background(),
		resourceURI,
		azlogs.QueryBody{
			Query:    to.Ptr("not a valid query"),
			Timespan: to.Ptr(azlogs.TimeInterval("PT2H")),
		},
		nil,
	)
	require.Error(t, err)
	require.Nil(t, res.Error)
	require.Nil(t, res.Tables)

	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "BadArgumentError")
	require.Equal(t, httpErr.StatusCode, 400)

	testSerde(t, &res)
}

func TestQueryResource_Advanced(t *testing.T) {
	client := startTest(t)

	res, err := client.QueryResource(context.Background(), resourceURI, azlogs.QueryBody{Query: &query},
		&azlogs.QueryResourceOptions{Options: &azlogs.QueryOptions{Statistics: to.Ptr(true), Visualization: to.Ptr(true), Wait: to.Ptr(600)}})
	require.NoError(t, err)
	require.Nil(t, res.Error)
	require.NotNil(t, res.Tables)
	require.NotNil(t, res.Visualization)
	require.NotNil(t, res.Statistics)
	testSerde(t, &res)
}

func TestQueryBatch_QuerySuccess(t *testing.T) {
	client := startTest(t)
	query1, query2 := query, query+" | take 2"
	timespan := azlogs.NewTimeInterval(time.Date(2022, 3, 2, 0, 0, 0, 0, time.UTC), time.Date(2022, 3, 3, 0, 0, 0, 0, time.UTC))

	batchRequest := azlogs.BatchRequest{[]azlogs.BatchQueryRequest{
		{Body: &azlogs.QueryBody{Query: to.Ptr(query1), Timespan: to.Ptr(timespan)}, ID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azlogs.QueryBody{Query: to.Ptr(query2), Timespan: to.Ptr(timespan)}, ID: to.Ptr("2"), WorkspaceID: to.Ptr(workspaceID)},
	}}
	testSerde(t, &batchRequest)

	res, err := client.QueryBatch(context.Background(), batchRequest, nil)
	require.NoError(t, err)
	require.Len(t, res.Responses, 2)
	for _, resp := range res.Responses {
		require.Nil(t, resp.Body.Error)
		require.NotNil(t, resp.Body.Tables)
		if *resp.ID == "1" && len(resp.Body.Tables[0].Rows) != 100 {
			t.Fatal("expected 100 rows from batch request 1")
		}
		if *resp.ID == "2" && len(resp.Body.Tables[0].Rows) != 2 {
			t.Fatal("expected 2 rows from batch request 2")
		}
	}
	testSerde(t, &res)
}

func TestQueryBatch_BasicQueryFailure(t *testing.T) {
	client := startTest(t)

	batchRequest := azlogs.BatchRequest{[]azlogs.BatchQueryRequest{
		{Body: &azlogs.QueryBody{Query: to.Ptr(query)}, ID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azlogs.QueryBody{Query: to.Ptr(query)}, ID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID)},
	}}
	testSerde(t, &batchRequest)

	res, err := client.QueryBatch(context.Background(), batchRequest, nil)
	require.Error(t, err)
	require.Nil(t, res.Responses)

	var httpErr *azcore.ResponseError
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, httpErr.ErrorCode, "BadArgumentError")
	require.Equal(t, httpErr.StatusCode, 400)
}

func TestQueryBatch_AdvancedQuerySuccess(t *testing.T) {
	client := startTest(t)
	timespan := azlogs.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))
	batchPrefer1 := "wait=600,include-statistics=true,include-render=true"
	headers1 := map[string]*string{"prefer": &batchPrefer1}
	batchPrefer2 := "wait=180,include-statistics=true,include-render=true"
	headers2 := map[string]*string{"prefer": &batchPrefer2}

	batchRequestAdvanced := azlogs.BatchRequest{[]azlogs.BatchQueryRequest{
		{Body: &azlogs.QueryBody{Query: to.Ptr(query), Timespan: to.Ptr(timespan)}, ID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID2), Headers: headers1},
		{Body: &azlogs.QueryBody{Query: to.Ptr(query), Timespan: to.Ptr(timespan)}, ID: to.Ptr("2"), WorkspaceID: to.Ptr(workspaceID2), Headers: headers2},
	}}
	testSerde(t, &batchRequestAdvanced)

	res, err := client.QueryBatch(context.Background(), batchRequestAdvanced, nil)
	require.NoError(t, err)
	require.Len(t, res.Responses, 2)
	for _, resp := range res.Responses {
		require.Nil(t, resp.Body.Error)
		require.NotNil(t, resp.Body.Tables)
		require.NotNil(t, resp.Body.Visualization)
		require.NotNil(t, resp.Body.Statistics)
		require.Len(t, resp.Body.Tables[0].Rows, 100)
	}
	testSerde(t, &res)
}

func TestQueryBatch_PartialError(t *testing.T) {
	client := startTest(t)

	batchRequest := azlogs.BatchRequest{[]azlogs.BatchQueryRequest{
		{Body: &azlogs.QueryBody{Query: to.Ptr("not a valid query")}, ID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azlogs.QueryBody{Query: to.Ptr(query)}, ID: to.Ptr("2"), WorkspaceID: to.Ptr(workspaceID)},
	}}

	res, err := client.QueryBatch(context.Background(), batchRequest, nil)
	require.NoError(t, err)
	require.Len(t, res.Responses, 2)
	for _, resp := range res.Responses {
		if *resp.ID == "1" {
			require.NotNil(t, resp.Body.Error)
			require.Equal(t, resp.Body.Error.Code, "BadArgumentError")
			require.Contains(t, resp.Body.Error.Error(), "BadArgumentError")
		}
		if *resp.ID == "2" {
			require.Nil(t, resp.Body.Error)
			require.Len(t, resp.Body.Tables[0].Rows, 100)
		}
	}
}

func TestTimeInterval(t *testing.T) {
	timespan := azlogs.NewTimeInterval(time.Date(2022, 3, 2, 1, 2, 3, 0, time.UTC), time.Date(2022, 3, 3, 0, 0, 0, 0, time.UTC))
	require.Equal(t, timespan, azlogs.TimeInterval("2022-03-02T01:02:03Z/2022-03-03T00:00:00Z"))

	start, end, err := timespan.Values()
	require.NoError(t, err)
	require.Equal(t, start, time.Date(2022, 3, 2, 1, 2, 3, 0, time.UTC))
	require.Equal(t, end, time.Date(2022, 3, 3, 0, 0, 0, 0, time.UTC))

	_, _, err = to.Ptr(azlogs.TimeInterval("hi")).Values()
	require.Error(t, err)
}
