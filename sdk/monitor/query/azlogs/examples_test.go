// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azlogs_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azlogs"
)

var (
	client      azlogs.Client
	kustoQuery1 string
	kustoQuery2 string
	kustoQuery3 string
)

type queryResult struct {
	Bool   bool
	Long   int64
	Double float64
	String string
}

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client, err := azlogs.NewClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	_ = client
}

func ExampleClient_QueryWorkspace() {
	// QueryWorkspace allows users to query log data.

	// A workspace ID is required to query logs. To find the workspace ID:
	// 1. If not already made, create a Log Analytics workspace (https://learn.microsoft.com/azure/azure-monitor/logs/quick-create-workspace).
	// 2. Navigate to your workspace's page in the Azure portal.
	// 3. From the **Overview** blade, copy the value of the ***Workspace ID*** property.

	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID

	res, err := client.QueryWorkspace(
		context.TODO(),
		workspaceID,
		azlogs.QueryBody{
			Query:    to.Ptr("AzureActivity | top 10 by TimeGenerated"), // example Kusto query
			Timespan: to.Ptr(azlogs.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
		},
		nil)
	if err != nil {
		//TODO: handle error
	}
	if res.Error != nil {
		//TODO: handle partial error
	}

	// Print Rows
	for _, table := range res.Tables {
		for _, row := range table.Rows {
			fmt.Println(row)
		}
	}
}

func ExampleClient_QueryWorkspace_second() {
	// `QueryWorkspace` also has more advanced options, including querying multiple workspaces
	// and QueryOptions (including statistics and visualization information and increasing default timeout).

	// When multiple workspaces are included in the query, the logs in the result table are not grouped
	// according to the workspace from which it was retrieved.
	workspaceID1 := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	workspaceID2 := "h4bc4471-2e8c-4b1c-8f47-12b9a4d5ac71"
	additionalWorkspaces := []string{workspaceID2}

	// Advanced query options
	// Setting Statistics to true returns stats information in Results.Statistics
	// Setting Visualization to true returns visualization information in Results.Visualization
	options := &azlogs.QueryWorkspaceOptions{
		Options: &azlogs.QueryOptions{
			Statistics:    to.Ptr(true),
			Visualization: to.Ptr(true),
			Wait:          to.Ptr(600),
		},
	}

	res, err := client.QueryWorkspace(
		context.TODO(),
		workspaceID1,
		azlogs.QueryBody{
			Query:                to.Ptr(query),
			Timespan:             to.Ptr(azlogs.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
			AdditionalWorkspaces: additionalWorkspaces,
		},
		options)
	if err != nil {
		//TODO: handle error
	}
	if res.Error != nil {
		//TODO: handle partial error
	}

	// Example of converting table data into a slice of structs.
	// Query results are returned in Table Rows and are of type any.
	// Type assertion is required to access the underlying value of each index in a Row.
	var QueryResults []queryResult
	for _, table := range res.Tables {
		QueryResults = make([]queryResult, len(table.Rows))
		for index, row := range table.Rows {
			QueryResults[index] = queryResult{
				Bool:   row[0].(bool),
				Long:   int64(row[1].(float64)),
				Double: float64(row[2].(float64)),
				String: row[3].(string),
			}
		}
	}

	fmt.Println(QueryResults)

	// Print out Statistics
	fmt.Printf("Statistics: %s", string(res.Statistics))

	// Print out Visualization information
	fmt.Printf("Visualization: %s", string(res.Visualization))

}

func ExampleClient_QueryResource() {
	// Instead of requiring a Log Analytics workspace,
	// QueryResource allows users to query logs directly from an Azure resource through a resource ID.

	// To find the resource ID:
	// 1. Navigate to your resource's page in the Azure portal.
	// 2. From the **Overview** blade, select the **JSON View** link.
	// 3. In the resulting JSON, copy the value of the `id` property.

	resourceID := "/subscriptions/fajfkx93-c1d8-40ad-9cce-e49c10ca8qe6/resourceGroups/testgroup/providers/Microsoft.Storage/storageAccounts/mystorageacount" // example resource ID

	res, err := client.QueryResource(
		context.TODO(),
		resourceID,
		azlogs.QueryBody{
			Query:    to.Ptr("StorageBlobLogs | where TimeGenerated > ago(3d)"), // example Kusto query
			Timespan: to.Ptr(azlogs.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
		},
		nil)
	if err != nil {
		//TODO: handle error
	}
	if res.Error != nil {
		//TODO: handle partial error
	}

	// Print Rows
	for _, table := range res.Tables {
		for _, row := range table.Rows {
			fmt.Println(row)
		}
	}
}

func ExampleClient_QueryBatch() {
	// `QueryBatch` is an advanced method allowing users to execute multiple log queries in a single request.

	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	timespan := azlogs.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))

	batchRequest := azlogs.BatchRequest{[]azlogs.BatchQueryRequest{
		{Body: &azlogs.QueryBody{Query: to.Ptr(kustoQuery1), Timespan: to.Ptr(timespan)}, ID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azlogs.QueryBody{Query: to.Ptr(kustoQuery2), Timespan: to.Ptr(timespan)}, ID: to.Ptr("2"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azlogs.QueryBody{Query: to.Ptr(kustoQuery3), Timespan: to.Ptr(timespan)}, ID: to.Ptr("3"), WorkspaceID: to.Ptr(workspaceID)},
	}}

	res, err := client.QueryBatch(context.TODO(), batchRequest, nil)
	if err != nil {
		//TODO: handle error
	}

	// `QueryBatch` can return results in any order, usually by time it takes each individual query to complete.
	// Use the `ID` field to identify the correct response.
	responses := res.Responses
	fmt.Println("ID's of successful responses:")
	for _, response := range responses {
		if response.Body.Error == nil {
			fmt.Println(*response.ID)
		}
	}
}
