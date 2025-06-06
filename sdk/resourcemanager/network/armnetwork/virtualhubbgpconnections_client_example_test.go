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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VirtualHubBgpConnectionList.json
func ExampleVirtualHubBgpConnectionsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualHubBgpConnectionsClient().NewListPager("rg1", "hub1", nil)
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
		// page.ListVirtualHubBgpConnectionResults = armnetwork.ListVirtualHubBgpConnectionResults{
		// 	Value: []*armnetwork.BgpConnection{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/bgpConnections/conn1"),
		// 			Name: to.Ptr("conn1"),
		// 			Etag: to.Ptr("W/\"72090554-7e3b-43f2-80ad-99a9020dcb11\""),
		// 			Properties: &armnetwork.BgpConnectionProperties{
		// 				ConnectionState: to.Ptr(armnetwork.HubBgpConnectionStatusConnected),
		// 				HubVirtualNetworkConnection: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualHubs/hub1/hubVirtualNetworkConnections/hubVnetConn1"),
		// 				},
		// 				PeerAsn: to.Ptr[int64](20000),
		// 				PeerIP: to.Ptr("192.168.1.5"),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VirtualRouterPeerListLearnedRoute.json
func ExampleVirtualHubBgpConnectionsClient_BeginListLearnedRoutes() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualHubBgpConnectionsClient().BeginListLearnedRoutes(ctx, "rg1", "virtualRouter1", "peer1", nil)
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
	// res.Value = map[string][]*armnetwork.PeerRoute{
	// 	"RouteServiceRole_IN_0": []*armnetwork.PeerRoute{
	// 		{
	// 			AsPath: to.Ptr("65002-65001"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.101.0.0/16"),
	// 			NextHop: to.Ptr("10.85.4.4"),
	// 			Origin: to.Ptr("EBgp"),
	// 			SourcePeer: to.Ptr("10.85.4.4"),
	// 			Weight: to.Ptr[int32](32768),
	// 		},
	// 		{
	// 			AsPath: to.Ptr("65002-65001"),
	// 			LocalAddress: to.Ptr("10.85.3.5"),
	// 			Network: to.Ptr("10.101.0.0/16"),
	// 			NextHop: to.Ptr("10.85.4.4"),
	// 			Origin: to.Ptr("EBgp"),
	// 			SourcePeer: to.Ptr("10.85.4.4"),
	// 			Weight: to.Ptr[int32](32768),
	// 	}},
	// 	"RouteServiceRole_IN_1": []*armnetwork.PeerRoute{
	// 		{
	// 			AsPath: to.Ptr("65002-65001"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.101.0.0/16"),
	// 			NextHop: to.Ptr("10.85.4.4"),
	// 			Origin: to.Ptr("EBgp"),
	// 			SourcePeer: to.Ptr("10.85.4.4"),
	// 			Weight: to.Ptr[int32](32768),
	// 		},
	// 		{
	// 			AsPath: to.Ptr("65002-65001"),
	// 			LocalAddress: to.Ptr("10.85.3.5"),
	// 			Network: to.Ptr("10.101.0.0/16"),
	// 			NextHop: to.Ptr("10.85.4.4"),
	// 			Origin: to.Ptr("EBgp"),
	// 			SourcePeer: to.Ptr("10.85.4.4"),
	// 			Weight: to.Ptr[int32](32768),
	// 	}},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/VirtualRouterPeerListAdvertisedRoute.json
func ExampleVirtualHubBgpConnectionsClient_BeginListAdvertisedRoutes() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualHubBgpConnectionsClient().BeginListAdvertisedRoutes(ctx, "rg1", "virtualRouter1", "peer1", nil)
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
	// res.Value = map[string][]*armnetwork.PeerRoute{
	// 	"RouteServiceRole_IN_0": []*armnetwork.PeerRoute{
	// 		{
	// 			AsPath: to.Ptr("65515"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.45.0.0/16"),
	// 			NextHop: to.Ptr("10.85.3.4"),
	// 			Origin: to.Ptr("Igp"),
	// 			SourcePeer: to.Ptr("10.85.3.4"),
	// 			Weight: to.Ptr[int32](0),
	// 		},
	// 		{
	// 			AsPath: to.Ptr("65515"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.85.0.0/16"),
	// 			NextHop: to.Ptr("10.85.3.4"),
	// 			Origin: to.Ptr("Igp"),
	// 			SourcePeer: to.Ptr("10.85.3.4"),
	// 			Weight: to.Ptr[int32](0),
	// 		},
	// 		{
	// 			AsPath: to.Ptr("65515"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.100.0.0/16"),
	// 			NextHop: to.Ptr("10.85.3.4"),
	// 			Origin: to.Ptr("Igp"),
	// 			SourcePeer: to.Ptr("10.85.3.4"),
	// 			Weight: to.Ptr[int32](0),
	// 	}},
	// 	"RouteServiceRole_IN_1": []*armnetwork.PeerRoute{
	// 		{
	// 			AsPath: to.Ptr("65515"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.45.0.0/16"),
	// 			NextHop: to.Ptr("10.85.3.4"),
	// 			Origin: to.Ptr("Igp"),
	// 			SourcePeer: to.Ptr("10.85.3.4"),
	// 			Weight: to.Ptr[int32](0),
	// 		},
	// 		{
	// 			AsPath: to.Ptr("65515"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.85.0.0/16"),
	// 			NextHop: to.Ptr("10.85.3.4"),
	// 			Origin: to.Ptr("Igp"),
	// 			SourcePeer: to.Ptr("10.85.3.4"),
	// 			Weight: to.Ptr[int32](0),
	// 		},
	// 		{
	// 			AsPath: to.Ptr("65515"),
	// 			LocalAddress: to.Ptr("10.85.3.4"),
	// 			Network: to.Ptr("10.100.0.0/16"),
	// 			NextHop: to.Ptr("10.85.3.4"),
	// 			Origin: to.Ptr("Igp"),
	// 			SourcePeer: to.Ptr("10.85.3.4"),
	// 			Weight: to.Ptr[int32](0),
	// 	}},
	// }
}
