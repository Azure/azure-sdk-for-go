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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkVirtualApplianceConnectionPut.json
func ExampleVirtualApplianceConnectionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualApplianceConnectionsClient().BeginCreateOrUpdate(ctx, "rg1", "nva1", "connection1", armnetwork.VirtualApplianceConnection{
		Properties: &armnetwork.VirtualApplianceConnectionProperties{
			Name: to.Ptr("connection1"),
			Asn:  to.Ptr[int64](64512),
			BgpPeerAddress: []*string{
				to.Ptr("169.254.16.13"),
				to.Ptr("169.254.16.14")},
			EnableInternetSecurity: to.Ptr(false),
			RoutingConfiguration: &armnetwork.RoutingConfiguration{
				AssociatedRouteTable: &armnetwork.SubResource{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
				},
				InboundRouteMap: &armnetwork.SubResource{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1"),
				},
				OutboundRouteMap: &armnetwork.SubResource{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap2"),
				},
				PropagatedRouteTables: &armnetwork.PropagatedRouteTable{
					IDs: []*armnetwork.SubResource{
						{
							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
						}},
					Labels: []*string{
						to.Ptr("label1")},
				},
			},
			TunnelIdentifier: to.Ptr[int64](0),
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
	// res.VirtualApplianceConnection = armnetwork.VirtualApplianceConnection{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkVirtualAppliances/nva1/networkVirtualApplianceConnections/connection1"),
	// 	Name: to.Ptr("connection1"),
	// 	Properties: &armnetwork.VirtualApplianceConnectionProperties{
	// 		Name: to.Ptr("connection1"),
	// 		Asn: to.Ptr[int64](64512),
	// 		BgpPeerAddress: []*string{
	// 			to.Ptr("169.254.16.13"),
	// 			to.Ptr("169.254.16.14")},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			RoutingConfiguration: &armnetwork.RoutingConfiguration{
	// 				AssociatedRouteTable: &armnetwork.SubResource{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
	// 				},
	// 				InboundRouteMap: &armnetwork.SubResource{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1"),
	// 				},
	// 				OutboundRouteMap: &armnetwork.SubResource{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap2"),
	// 				},
	// 				PropagatedRouteTables: &armnetwork.PropagatedRouteTable{
	// 					IDs: []*armnetwork.SubResource{
	// 						{
	// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
	// 					}},
	// 					Labels: []*string{
	// 						to.Ptr("label1")},
	// 					},
	// 				},
	// 				TunnelIdentifier: to.Ptr[int64](0),
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkVirtualApplianceConnectionGet.json
func ExampleVirtualApplianceConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualApplianceConnectionsClient().Get(ctx, "rg1", "nva1", "connection1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualApplianceConnection = armnetwork.VirtualApplianceConnection{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/NetworkVirtualAppliances/nva1/NetworkVirtualApplianceConnections/connection1"),
	// 	Name: to.Ptr("connection1"),
	// 	Properties: &armnetwork.VirtualApplianceConnectionProperties{
	// 		Name: to.Ptr("connection1"),
	// 		Asn: to.Ptr[int64](64512),
	// 		BgpPeerAddress: []*string{
	// 			to.Ptr("169.254.16.13"),
	// 			to.Ptr("169.254.16.14")},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			RoutingConfiguration: &armnetwork.RoutingConfiguration{
	// 				AssociatedRouteTable: &armnetwork.SubResource{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
	// 				},
	// 				InboundRouteMap: &armnetwork.SubResource{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1"),
	// 				},
	// 				OutboundRouteMap: &armnetwork.SubResource{
	// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap2"),
	// 				},
	// 				PropagatedRouteTables: &armnetwork.PropagatedRouteTable{
	// 					IDs: []*armnetwork.SubResource{
	// 						{
	// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
	// 					}},
	// 					Labels: []*string{
	// 						to.Ptr("label1")},
	// 					},
	// 				},
	// 				TunnelIdentifier: to.Ptr[int64](0),
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkVirtualApplianceConnectionDelete.json
func ExampleVirtualApplianceConnectionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualApplianceConnectionsClient().BeginDelete(ctx, "rg1", "nva1", "connection1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkVirtualApplianceConnectionList.json
func ExampleVirtualApplianceConnectionsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualApplianceConnectionsClient().NewListPager("rg1", "nva1", nil)
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
		// page.VirtualApplianceConnectionList = armnetwork.VirtualApplianceConnectionList{
		// 	Value: []*armnetwork.VirtualApplianceConnection{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/NetworkVirtualAppliances/nva1/NetworkVirtualApplianceConnections/connection1"),
		// 			Name: to.Ptr("connection1"),
		// 			Properties: &armnetwork.VirtualApplianceConnectionProperties{
		// 				Name: to.Ptr("connection1"),
		// 				Asn: to.Ptr[int64](64512),
		// 				BgpPeerAddress: []*string{
		// 					to.Ptr("169.254.16.13"),
		// 					to.Ptr("169.254.16.14")},
		// 					ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 					RoutingConfiguration: &armnetwork.RoutingConfiguration{
		// 						AssociatedRouteTable: &armnetwork.SubResource{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
		// 						},
		// 						InboundRouteMap: &armnetwork.SubResource{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap1"),
		// 						},
		// 						OutboundRouteMap: &armnetwork.SubResource{
		// 							ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/routeMaps/routeMap2"),
		// 						},
		// 						PropagatedRouteTables: &armnetwork.PropagatedRouteTable{
		// 							IDs: []*armnetwork.SubResource{
		// 								{
		// 									ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/hubRouteTable1"),
		// 							}},
		// 							Labels: []*string{
		// 								to.Ptr("label1")},
		// 							},
		// 						},
		// 						TunnelIdentifier: to.Ptr[int64](0),
		// 					},
		// 			}},
		// 		}
	}
}
