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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/VirtualHubRouteTableV2Get.json
func ExampleVirtualHubRouteTableV2SClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualHubRouteTableV2SClient().Get(ctx, "rg1", "virtualHub1", "virtualHubRouteTable1a", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualHubRouteTableV2 = armnetwork.VirtualHubRouteTableV2{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeTables/virtualHubRouteTable1a"),
	// 	Name: to.Ptr("virtualHubRouteTable1a"),
	// 	Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 	Properties: &armnetwork.VirtualHubRouteTableV2Properties{
	// 		AttachedConnections: []*string{
	// 			to.Ptr("All_Vnets")},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			Routes: []*armnetwork.VirtualHubRouteV2{
	// 				{
	// 					DestinationType: to.Ptr("CIDR"),
	// 					Destinations: []*string{
	// 						to.Ptr("20.10.0.0/16"),
	// 						to.Ptr("20.20.0.0/16")},
	// 						NextHopType: to.Ptr("IPAddress"),
	// 						NextHops: []*string{
	// 							to.Ptr("10.0.0.68")},
	// 						},
	// 						{
	// 							DestinationType: to.Ptr("CIDR"),
	// 							Destinations: []*string{
	// 								to.Ptr("0.0.0.0/0")},
	// 								NextHopType: to.Ptr("IPAddress"),
	// 								NextHops: []*string{
	// 									to.Ptr("10.0.0.68")},
	// 							}},
	// 						},
	// 					}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/VirtualHubRouteTableV2Put.json
func ExampleVirtualHubRouteTableV2SClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualHubRouteTableV2SClient().BeginCreateOrUpdate(ctx, "rg1", "virtualHub1", "virtualHubRouteTable1a", armnetwork.VirtualHubRouteTableV2{
		Properties: &armnetwork.VirtualHubRouteTableV2Properties{
			AttachedConnections: []*string{
				to.Ptr("All_Vnets")},
			Routes: []*armnetwork.VirtualHubRouteV2{
				{
					DestinationType: to.Ptr("CIDR"),
					Destinations: []*string{
						to.Ptr("20.10.0.0/16"),
						to.Ptr("20.20.0.0/16")},
					NextHopType: to.Ptr("IPAddress"),
					NextHops: []*string{
						to.Ptr("10.0.0.68")},
				},
				{
					DestinationType: to.Ptr("CIDR"),
					Destinations: []*string{
						to.Ptr("0.0.0.0/0")},
					NextHopType: to.Ptr("IPAddress"),
					NextHops: []*string{
						to.Ptr("10.0.0.68")},
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualHubRouteTableV2 = armnetwork.VirtualHubRouteTableV2{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeTables/virtualHubRouteTable1a"),
	// 	Name: to.Ptr("virtualHubRouteTable1a"),
	// 	Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 	Properties: &armnetwork.VirtualHubRouteTableV2Properties{
	// 		AttachedConnections: []*string{
	// 			to.Ptr("All_Vnets")},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			Routes: []*armnetwork.VirtualHubRouteV2{
	// 				{
	// 					DestinationType: to.Ptr("CIDR"),
	// 					Destinations: []*string{
	// 						to.Ptr("20.10.0.0/16"),
	// 						to.Ptr("20.20.0.0/16")},
	// 						NextHopType: to.Ptr("IPAddress"),
	// 						NextHops: []*string{
	// 							to.Ptr("10.0.0.68")},
	// 						},
	// 						{
	// 							DestinationType: to.Ptr("CIDR"),
	// 							Destinations: []*string{
	// 								to.Ptr("0.0.0.0/0")},
	// 								NextHopType: to.Ptr("IPAddress"),
	// 								NextHops: []*string{
	// 									to.Ptr("10.0.0.68")},
	// 							}},
	// 						},
	// 					}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/VirtualHubRouteTableV2Delete.json
func ExampleVirtualHubRouteTableV2SClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualHubRouteTableV2SClient().BeginDelete(ctx, "rg1", "virtualHub1", "virtualHubRouteTable1a", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/VirtualHubRouteTableV2List.json
func ExampleVirtualHubRouteTableV2SClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualHubRouteTableV2SClient().NewListPager("rg1", "virtualHub1", nil)
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
		// page.ListVirtualHubRouteTableV2SResult = armnetwork.ListVirtualHubRouteTableV2SResult{
		// }
	}
}
