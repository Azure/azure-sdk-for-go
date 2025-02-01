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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/PeerExpressRouteCircuitConnectionGet.json
func ExamplePeerExpressRouteCircuitConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPeerExpressRouteCircuitConnectionsClient().Get(ctx, "rg1", "ExpressRouteARMCircuitA", "AzurePrivatePeering", "60aee347-e889-4a42-8c1b-0aae8b1e4013", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PeerExpressRouteCircuitConnection = armnetwork.PeerExpressRouteCircuitConnection{
	// 	ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering/peerConnections/60aee347-e889-4a42-8c1b-0aae8b1e4013"),
	// 	Name: to.Ptr("60aee347-e889-4a42-8c1b-0aae8b1e4013"),
	// 	Etag: to.Ptr("W/\"6ffbbb06-da20-44ca-a34f-280c4653b1e9\""),
	// 	Properties: &armnetwork.PeerExpressRouteCircuitConnectionPropertiesFormat{
	// 		AddressPrefix: to.Ptr("20.0.0.0/29"),
	// 		AuthResourceGUID: to.Ptr(""),
	// 		CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
	// 		ConnectionName: to.Ptr("circuitConnectionWestusEastus"),
	// 		ExpressRouteCircuitPeering: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering"),
	// 		},
	// 		PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitB/peerings/AzurePrivatePeering"),
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/PeerExpressRouteCircuitConnectionList.json
func ExamplePeerExpressRouteCircuitConnectionsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewPeerExpressRouteCircuitConnectionsClient().NewListPager("rg1", "ExpressRouteARMCircuitA", "AzurePrivatePeering", nil)
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
		// page.PeerExpressRouteCircuitConnectionListResult = armnetwork.PeerExpressRouteCircuitConnectionListResult{
		// 	Value: []*armnetwork.PeerExpressRouteCircuitConnection{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering/peerConnections/60aee347-e889-4a42-8c1b-0aae8b1e4013"),
		// 			Name: to.Ptr("60aee347-e889-4a42-8c1b-0aae8b1e4013"),
		// 			Etag: to.Ptr("W/\"6ffbbb06-da20-44ca-a34f-280c4653b1e9\""),
		// 			Properties: &armnetwork.PeerExpressRouteCircuitConnectionPropertiesFormat{
		// 				AddressPrefix: to.Ptr("20.0.0.0/29"),
		// 				AuthResourceGUID: to.Ptr(""),
		// 				CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
		// 				ConnectionName: to.Ptr("circuitConnectionWestusEastus"),
		// 				ExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering"),
		// 				},
		// 				PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitB/peerings/AzurePrivatePeering"),
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 			},
		// 		},
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering/peerConnections/c8b17193-8dd3-4f61-866d-8cdd2e2e268e"),
		// 			Name: to.Ptr("c8b17193-8dd3-4f61-866d-8cdd2e2e268e"),
		// 			Etag: to.Ptr("W/\"6ffbbb06-da20-44ca-a34f-280c4653b1e9\""),
		// 			Properties: &armnetwork.PeerExpressRouteCircuitConnectionPropertiesFormat{
		// 				AddressPrefix: to.Ptr("30.0.0.0/29"),
		// 				AuthResourceGUID: to.Ptr("64283012-d377-421d-8398-f6aeb2ac7ea0"),
		// 				CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
		// 				ConnectionName: to.Ptr("circuitConnectionCentralusEastus"),
		// 				ExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering"),
		// 				},
		// 				PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid2/resourceGroups/rg1/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitC/peerings/AzurePrivatePeering"),
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 			},
		// 	}},
		// }
	}
}
