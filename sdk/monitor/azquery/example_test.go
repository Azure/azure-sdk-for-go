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
	query := "AzureActivity | top 10 by TimeGenerated"     // Example Kusto query
	timespan := "2022-08-30/2022-08-31"                    // ISO8601 Standard timespan

	res, err := client.QueryWorkspace(context.TODO(), workspaceID, azquery.Body{Query: to.Ptr(query), Timespan: to.Ptr(timespan)}, nil)
	if err != nil {
		//TODO: handle error
	}
	if res.Error != nil {
		//TODO: handle partial error
	}

	// example use case of processing table results
	// creates of slice of all tenantIDs resulting from the 'AzureActivity' query
	tenantIDs := make([]string, len(res.Tables[0].Rows))
	for index, row := range res.Tables[0].Rows {
		tenantIDs[index] = row[0].(string)
	}

	fmt.Println(tenantIDs)
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
	workspaceID := "g4d1e129-fb1e-4b0a-b234-250abc987ea65"                                                                                                                                                                                                                 // example Azure Log Analytics Workspace ID
	query := "let dt = datatable (Bool:bool, Long:long, Double: double, String: string, Decimal: decimal)\n" + "[false, 1, 12345.6789, 'string value', decimal(0.10101)];" + "range x from 1 to 10 step 1 | extend y=1 | join kind=fullouter dt on $left.y == $right.Long" // Example Kusto query
	timespan := "2022-08-30/2022-08-31"                                                                                                                                                                                                                                    // ISO8601 Standard timespan

	res, err := client.QueryWorkspace(context.TODO(), workspaceID, azquery.Body{Query: to.Ptr(query), Timespan: to.Ptr(timespan)}, nil)
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
		indexBool := table.ColumnIndexLookup["Bool"]
		indexLong := table.ColumnIndexLookup["Long"]
		indexDouble := table.ColumnIndexLookup["Double"]
		indexString := table.ColumnIndexLookup["String"]

		for index, row := range table.Rows {
			QueryResults[index] = queryResult{
				Bool:   row[indexBool].(bool),
				Long:   int64(row[indexLong].(float64)),
				Double: float64(row[indexDouble].(float64)),
				String: row[indexString].(string),
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
	timespan := "2022-08-30/2022-08-31"                    // ISO8601 Standard Timespan

	batchRequest := azquery.BatchRequest{[]*azquery.BatchQueryRequest{
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery1), Timespan: to.Ptr(timespan)}, ID: to.Ptr("1"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery2), Timespan: to.Ptr(timespan)}, ID: to.Ptr("2"), Workspace: to.Ptr(workspaceID)},
		{Body: &azquery.Body{Query: to.Ptr(kustoQuery3), Timespan: to.Ptr(timespan)}, ID: to.Ptr("3"), Workspace: to.Ptr(workspaceID)},
	}}

	res, err := client.QueryBatch(context.TODO(), batchRequest, nil)
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
	client, err := azquery.NewMetricsClient(cred, nil)
	if err != nil {
		//TODO: handle error
	}
	res, err := client.QueryResource(context.TODO(), resourceURI,
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
