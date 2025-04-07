//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsearch_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/search/armsearch"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/search/resource-manager/Microsoft.Search/preview/2024-06-01-preview/examples/SearchListOperations.json
func ExampleOperationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsearch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.OperationListResult = armsearch.OperationListResult{
		// 	Value: []*armsearch.Operation{
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/operations/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Lists all of the available operations of the Microsoft.Search provider."),
		// 				Operation: to.Ptr("List all available operations"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/register/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Registers the subscription for the search resource provider and enables the creation of search services."),
		// 				Operation: to.Ptr("Register the Search Resource Provider"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Creates or updates the search service."),
		// 				Operation: to.Ptr("Set Search Service"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Reads the search service."),
		// 				Operation: to.Ptr("Get Search Service"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Deletes the search service."),
		// 				Operation: to.Ptr("Delete Search Service"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/start/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Starts the search service."),
		// 				Operation: to.Ptr("Start Search Service"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/stop/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Stops the search service."),
		// 				Operation: to.Ptr("Stop Search Service"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/listAdminKeys/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Reads the admin keys."),
		// 				Operation: to.Ptr("Get Admin Key"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/regenerateAdminKey/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Regenerates the admin key."),
		// 				Operation: to.Ptr("Regenerate Admin Key"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/listQueryKeys/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Returns the list of query API keys for the given Azure AI Search service."),
		// 				Operation: to.Ptr("Get Query Keys"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("API Keys"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/createQueryKey/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Creates the query key."),
		// 				Operation: to.Ptr("Create Query Key"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Search Services"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/deleteQueryKey/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Deletes the query key."),
		// 				Operation: to.Ptr("Delete Query Key"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("API Keys"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/checkNameAvailability/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Checks availability of the service name."),
		// 				Operation: to.Ptr("Check Service Name Availability"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Service Name Availability"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/diagnosticSettings/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Gets the diganostic setting for the resource."),
		// 				Operation: to.Ptr("Get Diagnostic Setting"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Diagnostic Settings"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/diagnosticSettings/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Creates or updates the diganostic setting for the resource."),
		// 				Operation: to.Ptr("Set Diagnostic Setting"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Diagnostic Settings"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/metricDefinitions/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Gets the available metrics for the search service."),
		// 				Operation: to.Ptr("Read search service metric definitions"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("The metric definitions for the search service"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("system"),
		// 			Properties: &armsearch.OperationProperties{
		// 				ServiceSpecification: &armsearch.OperationServiceSpecification{
		// 					MetricSpecifications: []*armsearch.OperationMetricsSpecification{
		// 						{
		// 							Name: to.Ptr("SearchLatency"),
		// 							AggregationType: to.Ptr("Average"),
		// 							Availabilities: []*armsearch.OperationAvailability{
		// 								{
		// 									BlobDuration: to.Ptr("PT1H"),
		// 									TimeGrain: to.Ptr("PT1M"),
		// 							}},
		// 							DisplayDescription: to.Ptr("Average search latency for the search service"),
		// 							DisplayName: to.Ptr("Search Latency"),
		// 							Unit: to.Ptr("Seconds"),
		// 						},
		// 						{
		// 							Name: to.Ptr("SearchQueriesPerSecond"),
		// 							AggregationType: to.Ptr("Average"),
		// 							Availabilities: []*armsearch.OperationAvailability{
		// 								{
		// 									BlobDuration: to.Ptr("PT1H"),
		// 									TimeGrain: to.Ptr("PT1M"),
		// 							}},
		// 							DisplayDescription: to.Ptr("Search queries per second for the search service."),
		// 							DisplayName: to.Ptr("Search queries per second"),
		// 							Unit: to.Ptr("CountPerSecond"),
		// 						},
		// 						{
		// 							Name: to.Ptr("ThrottledSearchQueriesPercentage"),
		// 							AggregationType: to.Ptr("Average"),
		// 							Availabilities: []*armsearch.OperationAvailability{
		// 								{
		// 									BlobDuration: to.Ptr("PT1H"),
		// 									TimeGrain: to.Ptr("PT1M"),
		// 							}},
		// 							DisplayDescription: to.Ptr("Percentage of search queries that were throttled for the search service."),
		// 							DisplayName: to.Ptr("Throttled search queries percentage"),
		// 							Unit: to.Ptr("Percent"),
		// 					}},
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/logDefinitions/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Gets the available logs for the search service."),
		// 				Operation: to.Ptr("Read search service log definitions"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("The log definition for the search service"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("system"),
		// 			Properties: &armsearch.OperationProperties{
		// 				ServiceSpecification: &armsearch.OperationServiceSpecification{
		// 					LogSpecifications: []*armsearch.OperationLogsSpecification{
		// 						{
		// 							Name: to.Ptr("OperationLogs"),
		// 							BlobDuration: to.Ptr("PT1H"),
		// 							DisplayName: to.Ptr("Operation Logs"),
		// 					}},
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnectionProxies/validate/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Validates a private endpoint connection create call from NRP (Microsoft.Network Resource Provider) side."),
		// 				Operation: to.Ptr("Validate Private Endpoint Connection Proxy"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Private Endpoint Connection Proxy"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnectionProxies/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Creates a private endpoint connection proxy with the specified parameters or updates the properties or tags for the specified private endpoint connection proxy."),
		// 				Operation: to.Ptr("Create Private Endpoint Connection Proxy"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Private Endpoint Connection Proxy"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnectionProxies/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Returns the list of private endpoint connection proxies or gets the properties for the specified private endpoint connection proxy."),
		// 				Operation: to.Ptr("Get Private Endpoint Connection Proxy"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Private Endpoint Connection Proxy"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnectionProxies/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Deletes an existing private endpoint connection proxy."),
		// 				Operation: to.Ptr("Delete Private Endpoint Connection Proxy"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Private Endpoint Connection Proxy"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnections/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Creates a private endpoint connection with the specified parameters or updates the properties or tags for the specified private endpoint connections."),
		// 				Operation: to.Ptr("Create Private Endpoint Connection"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Private Endpoint Connection"),
		// 			},
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnections/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Returns the list of private endpoint connections or gets the properties for the specified private endpoint connection."),
		// 				Operation: to.Ptr("Get Private Endpoint Connection"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Private Endpoint Connection"),
		// 			},
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnections/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Deletes an existing private endpoint connection."),
		// 				Operation: to.Ptr("Delete Private Endpoint Connection"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Private Endpoint Connection"),
		// 			},
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/sharedPrivateLinkResources/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Creates a new shared private link resource with the specified parameters or updates the properties for the specified shared private link resource."),
		// 				Operation: to.Ptr("Create Shared Private Link Resource"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Shared Private Link Resource"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/sharedPrivateLinkResources/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Returns the list of shared private link resources or gets the properties for the specified shared private link resource."),
		// 				Operation: to.Ptr("Get Shared Private Link Resource"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Shared Private Link Resource"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/sharedPrivateLinkResources/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Deletes an existing shared private link resource."),
		// 				Operation: to.Ptr("Delete Shared Private Link Resource"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Shared Private Link Resource"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/sharedPrivateLinkResources/operationStatuses/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Get the details of a long running shared private link resource operation."),
		// 				Operation: to.Ptr("Get Operation Status"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Shared Private Link Resource"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexes/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Return an index or its statistics, return a list of indexes or their statistics, or test the lexical analysis components of an index."),
		// 				Operation: to.Ptr("Read Index"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Indexes"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexes/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Create an index or modify its properties."),
		// 				Operation: to.Ptr("Create or Update Index"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Indexes"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexes/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Delete an index."),
		// 				Operation: to.Ptr("Delete Index"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Indexes"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/synonymMaps/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Return a synonym map or a list of synonym maps."),
		// 				Operation: to.Ptr("Read Synonym Map"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Synonym Maps"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/synonymMaps/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Create a synonym map or modify its properties."),
		// 				Operation: to.Ptr("Create or Update Synonym Map"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Synonym Maps"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/synonymMaps/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Delete a synonym map."),
		// 				Operation: to.Ptr("Delete Synonym Map"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Synonym Maps"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/dataSources/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Return a data source or a list of data sources."),
		// 				Operation: to.Ptr("Read Data Source"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Data Sources"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/dataSources/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Create a data source or modify its properties."),
		// 				Operation: to.Ptr("Create or Update Data Source"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Data Sources"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/dataSources/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Delete a data source."),
		// 				Operation: to.Ptr("Delete Data Source"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Data Sources"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/skillsets/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Return a skillset or a list of skillsets."),
		// 				Operation: to.Ptr("Read Skillset"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Skillsets"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/skillsets/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Create a skillset or modify its properties."),
		// 				Operation: to.Ptr("Create or Update Skillset"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Skillsets"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/skillsets/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Delete a skillset."),
		// 				Operation: to.Ptr("Delete Skillset"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Skillsets"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexers/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Return an indexer or its status, or return a list of indexers or their statuses."),
		// 				Operation: to.Ptr("Read Indexer"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Indexers"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexers/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Create an indexer, modify its properties, or manage its execution."),
		// 				Operation: to.Ptr("Create or Manage Indexer"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Indexers"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexers/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Delete an indexer."),
		// 				Operation: to.Ptr("Delete Indexer"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Indexers"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/debugSessions/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Return a debug session or a list of debug sessions."),
		// 				Operation: to.Ptr("Read Debug Session"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Debug Sessions"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/debugSessions/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Create a debug session or modify its properties."),
		// 				Operation: to.Ptr("Create or Update Debug Session"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Debug Sessions"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/debugSessions/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Delete a debug session."),
		// 				Operation: to.Ptr("Delete Debug Session"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Debug Sessions"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/debugSessions/execute/action"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Use a debug session, get execution data, or evaluate expressions on it."),
		// 				Operation: to.Ptr("Execute Debug Session"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Debug Sessions"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexes/documents/read"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Read documents or suggested query terms from an index."),
		// 				Operation: to.Ptr("Read Documents"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Documents"),
		// 			},
		// 			IsDataAction: to.Ptr(true),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexes/documents/write"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Upload documents to an index or modify existing documents."),
		// 				Operation: to.Ptr("Write Documents"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Documents"),
		// 			},
		// 			IsDataAction: to.Ptr(true),
		// 			Origin: to.Ptr("user,system"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Search/searchServices/indexes/documents/delete"),
		// 			Display: &armsearch.OperationDisplay{
		// 				Description: to.Ptr("Delete documents from an index."),
		// 				Operation: to.Ptr("Delete Documents"),
		// 				Provider: to.Ptr("Microsoft Search"),
		// 				Resource: to.Ptr("Documents"),
		// 			},
		// 			IsDataAction: to.Ptr(true),
		// 			Origin: to.Ptr("user,system"),
		// 	}},
		// }
	}
}
