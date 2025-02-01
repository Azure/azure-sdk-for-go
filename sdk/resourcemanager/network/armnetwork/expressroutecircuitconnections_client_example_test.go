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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/ExpressRouteCircuitConnectionDelete.json
func ExampleExpressRouteCircuitConnectionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExpressRouteCircuitConnectionsClient().BeginDelete(ctx, "rg1", "ExpressRouteARMCircuitA", "AzurePrivatePeering", "circuitConnectionUSAUS", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/ExpressRouteCircuitConnectionGet.json
func ExampleExpressRouteCircuitConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewExpressRouteCircuitConnectionsClient().Get(ctx, "rg1", "ExpressRouteARMCircuitA", "AzurePrivatePeering", "circuitConnectionUSAUS", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ExpressRouteCircuitConnection = armnetwork.ExpressRouteCircuitConnection{
	// 	ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering/connections/circuitConnectionUSAUS"),
	// 	Name: to.Ptr("circuitConnectionUSAUS"),
	// 	Type: to.Ptr("Microsoft.Network/expressRouteCircuits/peerings/connections"),
	// 	Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 	Properties: &armnetwork.ExpressRouteCircuitConnectionPropertiesFormat{
	// 		AddressPrefix: to.Ptr("10.0.0.0/24"),
	// 		AuthorizationKey: to.Ptr("946a1918-b7a2-4917-b43c-8c4cdaee006a"),
	// 		CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
	// 		ExpressRouteCircuitPeering: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/dedharcktlocal/peerings/AzurePrivatePeering"),
	// 		},
	// 		IPv6CircuitConnectionConfig: &armnetwork.IPv6CircuitConnectionConfig{
	// 			AddressPrefix: to.Ptr("aa:bb::1/125"),
	// 			CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
	// 		},
	// 		PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid2/resourceGroups/dedharcktpeer/providers/Microsoft.Network/expressRouteCircuits/dedharcktremote/peerings/AzurePrivatePeering"),
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/ExpressRouteCircuitConnectionCreate.json
func ExampleExpressRouteCircuitConnectionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExpressRouteCircuitConnectionsClient().BeginCreateOrUpdate(ctx, "rg1", "ExpressRouteARMCircuitA", "AzurePrivatePeering", "circuitConnectionUSAUS", armnetwork.ExpressRouteCircuitConnection{
		Properties: &armnetwork.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix:    to.Ptr("10.0.0.0/29"),
			AuthorizationKey: to.Ptr("946a1918-b7a2-4917-b43c-8c4cdaee006a"),
			ExpressRouteCircuitPeering: &armnetwork.SubResource{
				ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/dedharcktlocal/peerings/AzurePrivatePeering"),
			},
			IPv6CircuitConnectionConfig: &armnetwork.IPv6CircuitConnectionConfig{
				AddressPrefix: to.Ptr("aa:bb::/125"),
			},
			PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
				ID: to.Ptr("/subscriptions/subid2/resourceGroups/dedharcktpeer/providers/Microsoft.Network/expressRouteCircuits/dedharcktremote/peerings/AzurePrivatePeering"),
			},
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
	// res.ExpressRouteCircuitConnection = armnetwork.ExpressRouteCircuitConnection{
	// 	ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering/connections/circuitConnectionUSAUS"),
	// 	Name: to.Ptr("circuitConnectionUSAUS"),
	// 	Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 	Properties: &armnetwork.ExpressRouteCircuitConnectionPropertiesFormat{
	// 		AddressPrefix: to.Ptr("10.0.0.0/24"),
	// 		AuthorizationKey: to.Ptr("946a1918-b7a2-4917-b43c-8c4cdaee006a"),
	// 		CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
	// 		ExpressRouteCircuitPeering: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/dedharcktlocal/peerings/AzurePrivatePeering"),
	// 		},
	// 		IPv6CircuitConnectionConfig: &armnetwork.IPv6CircuitConnectionConfig{
	// 			AddressPrefix: to.Ptr("aa:bb::1/125"),
	// 			CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
	// 		},
	// 		PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
	// 			ID: to.Ptr("/subscriptions/subid2/resourceGroups/dedharcktpeer/providers/Microsoft.Network/expressRouteCircuits/dedharcktremote/peerings/AzurePrivatePeering"),
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/ExpressRouteCircuitConnectionList.json
func ExampleExpressRouteCircuitConnectionsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewExpressRouteCircuitConnectionsClient().NewListPager("rg1", "ExpressRouteARMCircuitA", "AzurePrivatePeering", nil)
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
		// page.ExpressRouteCircuitConnectionListResult = armnetwork.ExpressRouteCircuitConnectionListResult{
		// 	Value: []*armnetwork.ExpressRouteCircuitConnection{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering/connections/circuitConnectionUSAUS"),
		// 			Name: to.Ptr("circuitConnectionUSAUS"),
		// 			Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
		// 			Properties: &armnetwork.ExpressRouteCircuitConnectionPropertiesFormat{
		// 				AddressPrefix: to.Ptr("10.0.0.0/24"),
		// 				AuthorizationKey: to.Ptr("946a1918-b7a2-4917-b43c-8c4cdaee006a"),
		// 				CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
		// 				ExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/dedharcktlocal/peerings/AzurePrivatePeering"),
		// 				},
		// 				IPv6CircuitConnectionConfig: &armnetwork.IPv6CircuitConnectionConfig{
		// 					AddressPrefix: to.Ptr("aa:bb::1/125"),
		// 					CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
		// 				},
		// 				PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid2/resourceGroups/dedharcktpeer/providers/Microsoft.Network/expressRouteCircuits/dedharcktremote/peerings/AzurePrivatePeering"),
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 			},
		// 		},
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/ExpressRouteARMCircuitA/peerings/AzurePrivatePeering/connections/circuitConnectionUSEUR"),
		// 			Name: to.Ptr("circuitConnectionUSEUR"),
		// 			Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
		// 			Properties: &armnetwork.ExpressRouteCircuitConnectionPropertiesFormat{
		// 				AddressPrefix: to.Ptr("20.0.0.0/24"),
		// 				CircuitConnectionStatus: to.Ptr(armnetwork.CircuitConnectionStatusConnected),
		// 				ExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharcktinit/providers/Microsoft.Network/expressRouteCircuits/dedharcktlocal/peerings/AzurePrivatePeering"),
		// 				},
		// 				PeerExpressRouteCircuitPeering: &armnetwork.SubResource{
		// 					ID: to.Ptr("/subscriptions/subid1/resourceGroups/dedharckteurope/providers/Microsoft.Network/expressRouteCircuits/dedharcktams/peerings/AzurePrivatePeering"),
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 			},
		// 	}},
		// }
	}
}
