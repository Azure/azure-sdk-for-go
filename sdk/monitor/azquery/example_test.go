//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azquery_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

var cred *azidentity.DefaultAzureCredential
var kustoQuery1 string
var kustoQuery2 string
var kustoQuery3 string

func ExampleNewLogsClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client := azquery.NewLogsClient(cred, nil)
	_ = client
}

func ExampleNewMetricsClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client := azquery.NewMetricsClient(cred, nil)
	_ = client
}

func ExampleLogsClient_QueryWorkspace() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client := azquery.NewLogsClient(cred, nil)
	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	query := "AzureActivity | top 10 by TimeGenerated"     // Kusto query
	timespan := "2022-08-30/2022-08-31"                    // ISO8601 Standard timespan

	res, err := client.QueryWorkspace(context.TODO(), workspaceID, azquery.Body{Query: to.Ptr(query), Timespan: to.Ptr(timespan)}, nil)
	if err != nil {
		//TODO: handle error
	}
	if res.Results.Error != nil {
		//TODO: handle partial error
	}

	table := res.Results.Tables[0]
	fmt.Println("Response rows:")
	for _, row := range table.Rows {
		fmt.Println(row)
	}
}

func ExampleLogsClient_Batch() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}
	client := azquery.NewLogsClient(cred, nil)
	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65" // example Azure Log Analytics Workspace ID
	timespan := "2022-08-30/2022-08-31"                    // ISO8601 Standard Timespan

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery1), Timespan: to.Ptr(timespan)}, ID: to.Ptr("1"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery2), Timespan: to.Ptr(timespan)}, ID: to.Ptr("2"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery3), Timespan: to.Ptr(timespan)}, ID: to.Ptr("3"), Workspace: to.Ptr(workspaceID)},
	}}

	res, err := client.Batch(context.TODO(), batchRequest, nil)
	if err != nil {
		//TODO: handle error
	}

	responses := res.BatchResponse.Responses
	fmt.Println("ID's of successful responses:")
	for _, response := range responses {
		if response.Body.Error == nil {
			fmt.Println(*response.ID)
		}
	}
}

func ExampleMetricsClient_QueryResource() {
	client := azquery.NewMetricsClient(cred, nil)
	res, err := client.QueryResource(context.Background(), resourceURI,
		&azquery.MetricsClientQueryResourceOptions{Timespan: to.Ptr("2017-04-14T02:20:00Z/2017-04-14T04:20:00Z"),
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
