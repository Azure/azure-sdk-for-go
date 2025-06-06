//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/OperationList.json
func ExampleOperationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
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
		// page.OperationListResult = armnetwork.OperationListResult{
		// 	Value: []*armnetwork.Operation{
		// 		{
		// 			Name: to.Ptr("Microsoft.Network/localnetworkgateways/read"),
		// 			Display: &armnetwork.OperationDisplay{
		// 				Description: to.Ptr("Gets LocalNetworkGateway"),
		// 				Operation: to.Ptr("Get LocalNetworkGateway"),
		// 				Provider: to.Ptr("Microsoft Network"),
		// 				Resource: to.Ptr("LocalNetworkGateway"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Network/localnetworkgateways/write"),
		// 			Display: &armnetwork.OperationDisplay{
		// 				Description: to.Ptr("Creates or updates an existing LocalNetworkGateway"),
		// 				Operation: to.Ptr("Create or update LocalNetworkGateway"),
		// 				Provider: to.Ptr("Microsoft Network"),
		// 				Resource: to.Ptr("LocalNetworkGateway"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Network/localnetworkgateways/delete"),
		// 			Display: &armnetwork.OperationDisplay{
		// 				Description: to.Ptr("Deletes LocalNetworkGateway"),
		// 				Operation: to.Ptr("Delete LocalNetworkGateway"),
		// 				Provider: to.Ptr("Microsoft Network"),
		// 				Resource: to.Ptr("LocalNetworkGateway"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Network/networkInterfaces/providers/Microsoft.Insights/metricDefinitions/read"),
		// 			Display: &armnetwork.OperationDisplay{
		// 				Description: to.Ptr("Gets available metrics for the Network Interface"),
		// 				Operation: to.Ptr("Read Network Interface metric definitions"),
		// 				Provider: to.Ptr("Microsoft Network"),
		// 				Resource: to.Ptr("Network Interface metric definition"),
		// 			},
		// 			Origin: to.Ptr("system"),
		// 			Properties: &armnetwork.OperationPropertiesFormat{
		// 				ServiceSpecification: &armnetwork.OperationPropertiesFormatServiceSpecification{
		// 					MetricSpecifications: []*armnetwork.MetricSpecification{
		// 						{
		// 							Name: to.Ptr("BytesSentRate"),
		// 							AggregationType: to.Ptr("Total"),
		// 							Availabilities: []*armnetwork.Availability{
		// 								{
		// 									BlobDuration: to.Ptr("01:00:00"),
		// 									Retention: to.Ptr("00:00:00"),
		// 									TimeGrain: to.Ptr("00:01:00"),
		// 								},
		// 								{
		// 									BlobDuration: to.Ptr("1.00:00:00"),
		// 									Retention: to.Ptr("00:00:00"),
		// 									TimeGrain: to.Ptr("01:00:00"),
		// 							}},
		// 							Dimensions: []*armnetwork.Dimension{
		// 							},
		// 							DisplayDescription: to.Ptr("Number of bytes the Network Interface sent"),
		// 							DisplayName: to.Ptr("Bytes Sent"),
		// 							EnableRegionalMdmAccount: to.Ptr(false),
		// 							FillGapWithZero: to.Ptr(false),
		// 							IsInternal: to.Ptr(false),
		// 							MetricFilterPattern: to.Ptr("^__Ready__$"),
		// 							Unit: to.Ptr("Count"),
		// 						},
		// 						{
		// 							Name: to.Ptr("BytesReceivedRate"),
		// 							AggregationType: to.Ptr("Total"),
		// 							Availabilities: []*armnetwork.Availability{
		// 								{
		// 									BlobDuration: to.Ptr("01:00:00"),
		// 									Retention: to.Ptr("00:00:00"),
		// 									TimeGrain: to.Ptr("00:01:00"),
		// 								},
		// 								{
		// 									BlobDuration: to.Ptr("1.00:00:00"),
		// 									Retention: to.Ptr("00:00:00"),
		// 									TimeGrain: to.Ptr("01:00:00"),
		// 							}},
		// 							Dimensions: []*armnetwork.Dimension{
		// 							},
		// 							DisplayDescription: to.Ptr("Number of bytes the Network Interface received"),
		// 							DisplayName: to.Ptr("Bytes Received"),
		// 							EnableRegionalMdmAccount: to.Ptr(false),
		// 							FillGapWithZero: to.Ptr(false),
		// 							IsInternal: to.Ptr(false),
		// 							MetricFilterPattern: to.Ptr("^__Ready__$"),
		// 							Unit: to.Ptr("Count"),
		// 					}},
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.Network/networksecuritygroups/providers/Microsoft.Insights/logDefinitions/read"),
		// 			Display: &armnetwork.OperationDisplay{
		// 				Description: to.Ptr("Gets the events for network security group"),
		// 				Operation: to.Ptr("Get Network Security Group Event Log Definitions"),
		// 				Provider: to.Ptr("Microsoft Network"),
		// 				Resource: to.Ptr("Network Security Groups Log Definitions"),
		// 			},
		// 			Origin: to.Ptr("system"),
		// 			Properties: &armnetwork.OperationPropertiesFormat{
		// 				ServiceSpecification: &armnetwork.OperationPropertiesFormatServiceSpecification{
		// 					LogSpecifications: []*armnetwork.LogSpecification{
		// 						{
		// 							Name: to.Ptr("NetworkSecurityGroupEvent"),
		// 							BlobDuration: to.Ptr("PT1H"),
		// 							DisplayName: to.Ptr("Network Security Group Event"),
		// 						},
		// 						{
		// 							Name: to.Ptr("NetworkSecurityGroupRuleCounter"),
		// 							BlobDuration: to.Ptr("PT1H"),
		// 							DisplayName: to.Ptr("Network Security Group Rule Counter"),
		// 						},
		// 						{
		// 							Name: to.Ptr("NetworkSecurityGroupFlowEvent"),
		// 							BlobDuration: to.Ptr("PT1H"),
		// 							DisplayName: to.Ptr("Network Security Group Rule Flow Event"),
		// 					}},
		// 				},
		// 			},
		// 	}},
		// }
	}
}
