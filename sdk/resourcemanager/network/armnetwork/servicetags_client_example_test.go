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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/ServiceTagsList.json
func ExampleServiceTagsClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceTagsClient().List(ctx, "westcentralus", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ServiceTagsListResult = armnetwork.ServiceTagsListResult{
	// 	Name: to.Ptr("public"),
	// 	Type: to.Ptr("Microsoft.Network/serviceTags"),
	// 	ChangeNumber: to.Ptr("63"),
	// 	Cloud: to.Ptr("Public"),
	// 	ID: to.Ptr("/subscriptions/subId/providers/Microsoft.Network/serviceTags/public"),
	// 	Values: []*armnetwork.ServiceTagInformation{
	// 		{
	// 			Name: to.Ptr("ApiManagement"),
	// 			ID: to.Ptr("ApiManagement"),
	// 			Properties: &armnetwork.ServiceTagInformationPropertiesFormat{
	// 				AddressPrefixes: []*string{
	// 					to.Ptr("13.64.39.16/32"),
	// 					to.Ptr("40.74.146.80/31"),
	// 					to.Ptr("40.74.147.32/28")},
	// 					ChangeNumber: to.Ptr("7"),
	// 					Region: to.Ptr(""),
	// 					SystemService: to.Ptr("AzureApiManagement"),
	// 				},
	// 			},
	// 			{
	// 				Name: to.Ptr("ApiManagement.AustraliaCentral"),
	// 				ID: to.Ptr("ApiManagement.AustraliaCentral"),
	// 				Properties: &armnetwork.ServiceTagInformationPropertiesFormat{
	// 					AddressPrefixes: []*string{
	// 						to.Ptr("20.36.106.68/31"),
	// 						to.Ptr("20.36.107.176/28")},
	// 						ChangeNumber: to.Ptr("2"),
	// 						Region: to.Ptr("australiacentral"),
	// 						SystemService: to.Ptr("AzureApiManagement"),
	// 					},
	// 				},
	// 				{
	// 					Name: to.Ptr("AppService"),
	// 					ID: to.Ptr("AppService"),
	// 					Properties: &armnetwork.ServiceTagInformationPropertiesFormat{
	// 						AddressPrefixes: []*string{
	// 							to.Ptr("13.64.73.110/32"),
	// 							to.Ptr("191.235.208.12/32"),
	// 							to.Ptr("191.235.215.184/32")},
	// 							ChangeNumber: to.Ptr("13"),
	// 							Region: to.Ptr(""),
	// 							SystemService: to.Ptr("AzureAppService"),
	// 						},
	// 					},
	// 					{
	// 						Name: to.Ptr("ServiceBus"),
	// 						ID: to.Ptr("ServiceBus"),
	// 						Properties: &armnetwork.ServiceTagInformationPropertiesFormat{
	// 							AddressPrefixes: []*string{
	// 								to.Ptr("23.98.82.96/29"),
	// 								to.Ptr("40.68.127.68/32"),
	// 								to.Ptr("40.70.146.64/29")},
	// 								ChangeNumber: to.Ptr("10"),
	// 								Region: to.Ptr(""),
	// 								SystemService: to.Ptr("AzureServiceBus"),
	// 							},
	// 						},
	// 						{
	// 							Name: to.Ptr("ServiceBus.EastUS2"),
	// 							ID: to.Ptr("ServiceBus.EastUS2"),
	// 							Properties: &armnetwork.ServiceTagInformationPropertiesFormat{
	// 								AddressPrefixes: []*string{
	// 									to.Ptr("13.68.110.36/32"),
	// 									to.Ptr("40.70.146.64/29")},
	// 									ChangeNumber: to.Ptr("1"),
	// 									Region: to.Ptr("eastus2"),
	// 									SystemService: to.Ptr("AzureServiceBus"),
	// 								},
	// 						}},
	// 					}
}
