//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

var logsClient azquery.LogsClient
var metricsClient azquery.MetricsClient
var metricsBatchClient azquery.MetricsBatchClient
var kustoQuery1 string
var kustoQuery2 string
var kustoQuery3 string

type queryResult struct {
	Bool   bool
	Long   int64
	Double float64
	String string
}

func ExampleNewLogsClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client, err := azquery.NewLogsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	_ = client
}

func ExampleNewMetricsClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	_ = client
}

func ExampleNewMetricsBatchClient() {
	// The regional endpoint to use. The region should match the region of the requested resources.
	// For global resources, the region should be 'global'
	endpoint := "https://eastus.metrics.monitor.azure.com"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client, err := azquery.NewMetricsBatchClient(endpoint, cred, nil)
	if err != nil {
		//TODO: handle error
	}
	_ = client
}

func ExampleLogsClient_QueryWorkspace() {
	// QueryWorkspace allows users to query log data.

	// A workspace ID is required to query logs. To find the workspace ID:
	// 1. If not already made, create a Log Analytics workspace (https://learn.microsoft.com/azure/azure-monitor/logs/quick-create-workspace).
	// 2. Navigate to your workspace's page in the Azure portal.
	// 3. From the **Overview** blade, copy the value of the ***Workspace ID*** property.

	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID

	res, err := logsClient.QueryWorkspace(
		context.TODO(),
		workspaceID,
		azquery.Body{
			Query:    to.Ptr("AzureActivity | top 10 by TimeGenerated"), // example Kusto query
			Timespan: to.Ptr(azquery.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
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

func ExampleLogsClient_QueryWorkspace_second() {
	// `QueryWorkspace` also has more advanced options, including querying multiple workspaces
	// and LogsQueryOptions (including statistics and visualization information and increasing default timeout).

	// When multiple workspaces are included in the query, the logs in the result table are not grouped
	// according to the workspace from which it was retrieved.
	workspaceID1 := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	workspaceID2 := "h4bc4471-2e8c-4b1c-8f47-12b9a4d5ac71"
	additionalWorkspaces := []*string{to.Ptr(workspaceID2)}

	// Advanced query options
	// Setting Statistics to true returns stats information in Results.Statistics
	// Setting Visualization to true returns visualization information in Results.Visualization
	options := &azquery.LogsClientQueryWorkspaceOptions{
		Options: &azquery.LogsQueryOptions{
			Statistics:    to.Ptr(true),
			Visualization: to.Ptr(true),
			Wait:          to.Ptr(600),
		},
	}

	res, err := logsClient.QueryWorkspace(
		context.TODO(),
		workspaceID1,
		azquery.Body{
			Query:                to.Ptr(query),
			Timespan:             to.Ptr(azquery.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
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

func ExampleLogsClient_QueryResource() {
	// Instead of requiring a Log Analytics workspace,
	// QueryResource allows users to query logs directly from an Azure resource through a resource ID.

	// To find the resource ID:
	// 1. Navigate to your resource's page in the Azure portal.
	// 2. From the **Overview** blade, select the **JSON View** link.
	// 3. In the resulting JSON, copy the value of the `id` property.

	resourceID := "/subscriptions/fajfkx93-c1d8-40ad-9cce-e49c10ca8qe6/resourceGroups/testgroup/providers/Microsoft.Storage/storageAccounts/mystorageacount" // example resource ID

	res, err := logsClient.QueryResource(
		context.TODO(),
		resourceID,
		azquery.Body{
			Query:    to.Ptr("StorageBlobLogs | where TimeGenerated > ago(3d)"), // example Kusto query
			Timespan: to.Ptr(azquery.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
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

func ExampleLogsClient_QueryBatch() {
	// `QueryBatch` is an advanced method allowing users to execute multiple log queries in a single request.
	// For help formatting a `BatchRequest`, please use the method `NewBatchQueryRequest`.

	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	timespan := azquery.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery1), Timespan: to.Ptr(timespan)}, CorrelationID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery2), Timespan: to.Ptr(timespan)}, CorrelationID: to.Ptr("2"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery3), Timespan: to.Ptr(timespan)}, CorrelationID: to.Ptr("3"), WorkspaceID: to.Ptr(workspaceID)},
	}}

	res, err := logsClient.QueryBatch(context.TODO(), batchRequest, nil)
	if err != nil {
		//TODO: handle error
	}

	// `QueryBatch` can return results in any order, usually by time it takes each individual query to complete.
	// Use the `CorrelationID` field to identify the correct response.
	responses := res.BatchResponse.Responses
	fmt.Println("ID's of successful responses:")
	for _, response := range responses {
		if response.Body.Error == nil {
			fmt.Println(*response.CorrelationID)
		}
	}
}

func ExampleMetricsClient_QueryResource() {
	// QueryResource is used to query metrics on an Azure resource.
	// For each requested metric, a set of aggregated values is returned inside the `TimeSeries` collection.

	// resource ID is required to query metrics. To find the resource ID:
	// 1. Navigate to your resource's page in the Azure portal.
	// 2. From the **Overview** blade, select the **JSON View** link.
	// 3. In the resulting JSON, copy the value of the `id` property.
	resourceURI := "subscriptions/182c901a-129a-4f5d-86e4-afdsb294590a2/resourceGroups/test-log/providers/microsoft.insights/components/f1-bill/providers/microsoft.insights/metricdefinitions"

	res, err := metricsClient.QueryResource(context.TODO(), resourceURI,
		&azquery.MetricsClientQueryResourceOptions{
			Timespan:        to.Ptr(azquery.NewTimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
			Interval:        to.Ptr("PT1M"),
			Aggregation:     to.SliceOfPtrs(azquery.AggregationTypeAverage, azquery.AggregationTypeCount),
			Top:             to.Ptr[int32](3),
			OrderBy:         to.Ptr("Average asc"),
			Filter:          to.Ptr("BlobType eq '*'"),
			MetricNamespace: to.Ptr("Microsoft.Storage/storageAccounts/blobServices"),
		})
	if err != nil {
		//TODO: handle error
	}

	// Print out metric name and the time stamps of each metric data point
	for _, metric := range res.Value {
		fmt.Println(*metric.Name.Value)
		for _, timeSeriesElement := range metric.TimeSeries {
			for _, metricValue := range timeSeriesElement.Data {
				fmt.Println(metricValue.TimeStamp)
			}
		}
	}
}

func ExampleMetricsClient_QueryResource_second() {
	// `QueryResource` with a filter specified.
	// This example queries an Azure Application Insights resource
	// for all the server requests over the last day that failed.
	_, err := metricsClient.QueryResource(context.Background(), resourceURI,
		&azquery.MetricsClientQueryResourceOptions{
			Aggregation: to.SliceOfPtrs(azquery.AggregationTypeCount),
			Filter:      to.Ptr("request/success eq '0'"),
			Interval:    to.Ptr("PT5M"),
			MetricNames: to.Ptr("requests/count"),
			Timespan:    to.Ptr(azquery.TimeInterval("P1D")),
		})
	if err != nil {
		//TODO: handle error
	}
}

func ExampleMetricsClient_NewListDefinitionsPager() {
	pager := metricsClient.NewListDefinitionsPager(resourceURI, &azquery.MetricsClientListDefinitionsOptions{MetricNamespace: to.Ptr("microsoft.insights/components")})
	for pager.More() {
		nextResult, err := pager.NextPage(context.TODO())
		if err != nil {
			//TODO: handle error
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

func ExampleMetricsClient_NewListNamespacesPager() {
	pager := metricsClient.NewListNamespacesPager(resourceURI, &azquery.MetricsClientListNamespacesOptions{StartTime: to.Ptr("2020-08-31T15:53:00Z")})
	for pager.More() {
		nextResult, err := pager.NextPage(context.TODO())
		if err != nil {
			//TODO: handle error
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

func ExampleMetricsBatchClient_QueryBatch() {
	// This sample uses the MetricsBatchClient to retrieve the "Ingress"
	// metric along with the "Average" aggregation type for multiple resources.
	// The query will execute over a timespan of 2 hours with a interval (granularity) of 5 minutes.

	// In this example, storage account resource URIs are queried for metrics.
	resourceURI1 := "/subscriptions/<id>/resourceGroups/<rg>/providers/Microsoft.Storage/storageAccounts/<account-1>"
	resourceURI2 := "/subscriptions/<id>/resourceGroups/<rg>/providers/Microsoft.Storage/storageAccounts/<account-2>"

	res, err := metricsBatchClient.QueryBatch(
		context.Background(),
		subscriptionID,
		"Microsoft.Storage/storageAccounts",
		[]string{"Ingress"},
		azquery.ResourceIDList{ResourceIDs: to.SliceOfPtrs(resourceURI1, resourceURI2)},
		&azquery.MetricsBatchClientQueryBatchOptions{
			Aggregation: to.SliceOfPtrs(azquery.AggregationTypeAverage),
			StartTime:   to.Ptr("2023-11-15"),
			EndTime:     to.Ptr("2023-11-16"),
			Interval:    to.Ptr("PT5M"),
		},
	)
	if err != nil {
		//TODO: handle error
	}

	// Print out results
	for _, result := range res.Values {
		for _, metric := range result.Values {
			fmt.Println(*metric.Name.Value + ": " + *metric.DisplayDescription)
			for _, timeSeriesElement := range metric.TimeSeries {
				for _, metricValue := range timeSeriesElement.Data {
					fmt.Printf("The ingress at %v is %v.\n", metricValue.TimeStamp.String(), *metricValue.Average)
				}
			}
		}
	}
}
