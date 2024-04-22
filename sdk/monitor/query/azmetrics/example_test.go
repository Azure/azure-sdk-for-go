// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azmetrics_test

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"github.com/Azure/azure-sdk-for-go/sdk/monitor/query/azmetrics"
)

var client azmetrics.Client

func ExampleNewClient() {
	// The regional endpoint to use. The region should match the region of the requested resources.
	// For global resources, the region should be 'global'
	endpoint := "https://eastus.metrics.monitor.azure.com"

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		//TODO: handle error
	}

	client, err := azmetrics.NewClient(endpoint, cred, nil)
	if err != nil {
		//TODO: handle error
	}
	_ = client
}

func ExampleClient_QueryResources() {
	// This sample uses the Client to retrieve the "Ingress"
	// metric along with the "Average" aggregation type for multiple resources.
	// The query will execute over a timespan of 2 hours with a interval (granularity) of 5 minutes.

	// In this example, storage account resource URIs are queried for metrics.
	resourceURI1 := "/subscriptions/<id>/resourceGroups/<rg>/providers/Microsoft.Storage/storageAccounts/<account-1>"
	resourceURI2 := "/subscriptions/<id>/resourceGroups/<rg>/providers/Microsoft.Storage/storageAccounts/<account-2>"

	res, err := client.QueryResources(
		context.Background(),
		subscriptionID,
		"Microsoft.Storage/storageAccounts",
		[]string{"Ingress"},
		azmetrics.ResourceIDList{ResourceIDs: []string{resourceURI1, resourceURI2}},
		&azmetrics.QueryResourcesOptions{
			Aggregation: to.Ptr("average"),
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
