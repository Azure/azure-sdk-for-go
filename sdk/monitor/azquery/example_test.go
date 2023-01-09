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

var cred *azidentity.DefaultAzureCredential
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

func ExampleLogsClient_QueryWorkspace() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := azquery.NewLogsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID

	res, err := client.QueryWorkspace(
		context.TODO(),
		workspaceID,
		azquery.Body{
			Query:    to.Ptr("AzureActivity | top 10 by TimeGenerated"), // example Kusto query
			Timespan: to.Ptr(azquery.NewISO8601TimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
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
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := azquery.NewLogsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	workspaceID1 := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	workspaceID2 := "h4bc4471-2e8c-4b1c-8f47-12b9a4d5ac71"
	preferOptions := "include-statistics=true" // Advanced option to include stats. See docs: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery#readme-increase-wait-time-include-statistics-include-render-visualization
	res, err := client.QueryWorkspace(
		context.TODO(),
		workspaceID1,
		azquery.Body{
			Query:                to.Ptr(query),
			Timespan:             to.Ptr(azquery.NewISO8601TimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
			AdditionalWorkspaces: []*string{to.Ptr(workspaceID2)},
		},
		to.Ptr(azquery.LogsClientQueryWorkspaceOptions{Prefer: to.Ptr(preferOptions)}))
	if err != nil {
		//TODO: handle error
	}
	if res.Error != nil {
		//TODO: handle partial error
	}

	// Example of converting table data into a slice of structs
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
}

func ExampleLogsClient_QueryBatch() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := azquery.NewLogsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}

	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	timespan := azquery.NewISO8601TimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery1), Timespan: to.Ptr(timespan)}, CorrelationID: to.Ptr("1"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery2), Timespan: to.Ptr(timespan)}, CorrelationID: to.Ptr("2"), WorkspaceID: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery3), Timespan: to.Ptr(timespan)}, CorrelationID: to.Ptr("3"), WorkspaceID: to.Ptr(workspaceID)},
	}}

	res, err := client.QueryBatch(context.TODO(), batchRequest, nil)
	if err != nil {
		//TODO: handle error
	}

	responses := res.BatchResponse.Responses
	fmt.Println("ID's of successful responses:")
	for _, response := range responses {
		if response.Body.Error == nil {
			fmt.Println(*response.CorrelationID)
		}
	}
}

func ExampleMetricsClient_QueryResource() {
	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}

	if err != nil {
		//TODO: handle error
	}
	res, err := client.QueryResource(context.TODO(), resourceURI,
		&azquery.MetricsClientQueryResourceOptions{
			Timespan:        to.Ptr(azquery.NewISO8601TimeInterval(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC))),
			Interval:        to.Ptr("PT1M"),
			Metricnames:     nil,
			Aggregation:     to.Ptr("Average,count"),
			Top:             to.Ptr[int32](3),
			Orderby:         to.Ptr("Average asc"),
			Filter:          to.Ptr("BlobType eq '*'"),
			ResultType:      nil,
			Metricnamespace: to.Ptr("Microsoft.Storage/storageAccounts/blobServices"),
		})
	if err != nil {
		//TODO: handle error
	}
	_ = res
}

func ExampleMetricsClient_NewListDefinitionsPager() {
	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	pager := client.NewListDefinitionsPager("subscriptions/182c901a-129a-4f5d-86e4-cc6b294590a2/resourceGroups/hyr-log/providers/microsoft.insights/components/f1-bill/providers/microsoft.insights/metricdefinitions", &azquery.MetricsClientListDefinitionsOptions{Metricnamespace: to.Ptr("microsoft.insights/components")})
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
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	pager := client.NewListNamespacesPager("subscriptions/182c901a-129a-4f5d-86e4-cc6b294590a2/resourceGroups/hyr-log/providers/microsoft.insights/components/f1-bill", &azquery.MetricsClientListNamespacesOptions{StartTime: to.Ptr("2020-08-31T15:53:00Z")})
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
