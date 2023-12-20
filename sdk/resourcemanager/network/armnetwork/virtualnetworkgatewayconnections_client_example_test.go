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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionCreate.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginCreateOrUpdate(ctx, "rg1", "connS2S", armnetwork.VirtualNetworkGatewayConnection{
		Location: to.Ptr("centralus"),
		Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
			ConnectionMode:     to.Ptr(armnetwork.VirtualNetworkGatewayConnectionModeDefault),
			ConnectionProtocol: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
			ConnectionType:     to.Ptr(armnetwork.VirtualNetworkGatewayConnectionTypeIPsec),
			DpdTimeoutSeconds:  to.Ptr[int32](30),
			EgressNatRules: []*armnetwork.SubResource{
				{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule2"),
				}},
			EnableBgp: to.Ptr(false),
			GatewayCustomBgpIPAddresses: []*armnetwork.GatewayCustomBgpIPAddressIPConfiguration{
				{
					CustomBgpIPAddress: to.Ptr("169.254.21.1"),
					IPConfigurationID:  to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/default"),
				},
				{
					CustomBgpIPAddress: to.Ptr("169.254.21.3"),
					IPConfigurationID:  to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/ActiveActive"),
				}},
			IngressNatRules: []*armnetwork.SubResource{
				{
					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule1"),
				}},
			IPSecPolicies: []*armnetwork.IPSecPolicy{},
			LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
				ID:       to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw"),
				Location: to.Ptr("centralus"),
				Tags:     map[string]*string{},
				Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
					GatewayIPAddress: to.Ptr("x.x.x.x"),
					LocalNetworkAddressSpace: &armnetwork.AddressSpace{
						AddressPrefixes: []*string{
							to.Ptr("10.1.0.0/16")},
					},
				},
			},
			RoutingWeight:                  to.Ptr[int32](0),
			SharedKey:                      to.Ptr("Abc123"),
			TrafficSelectorPolicies:        []*armnetwork.TrafficSelectorPolicy{},
			UsePolicyBasedTrafficSelectors: to.Ptr(false),
			VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
				ID:       to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw"),
				Location: to.Ptr("centralus"),
				Tags:     map[string]*string{},
				Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
					Active: to.Ptr(false),
					BgpSettings: &armnetwork.BgpSettings{
						Asn:               to.Ptr[int64](65514),
						BgpPeeringAddress: to.Ptr("10.0.1.30"),
						PeerWeight:        to.Ptr[int32](0),
					},
					EnableBgp:   to.Ptr(false),
					GatewayType: to.Ptr(armnetwork.VirtualNetworkGatewayTypeVPN),
					IPConfigurations: []*armnetwork.VirtualNetworkGatewayIPConfiguration{
						{
							ID:   to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/gwipconfig1"),
							Name: to.Ptr("gwipconfig1"),
							Properties: &armnetwork.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
								PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
								PublicIPAddress: &armnetwork.SubResource{
									ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/publicIPAddresses/gwpip"),
								},
								Subnet: &armnetwork.SubResource{
									ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/GatewaySubnet"),
								},
							},
						}},
					SKU: &armnetwork.VirtualNetworkGatewaySKU{
						Name: to.Ptr(armnetwork.VirtualNetworkGatewaySKUNameVPNGw1),
						Tier: to.Ptr(armnetwork.VirtualNetworkGatewaySKUTierVPNGw1),
					},
					VPNType: to.Ptr(armnetwork.VPNTypeRouteBased),
				},
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
	// res.VirtualNetworkGatewayConnection = armnetwork.VirtualNetworkGatewayConnection{
	// 	Name: to.Ptr("connS2S"),
	// 	Type: to.Ptr("Microsoft.Network/connections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/connections/connS2S"),
	// 	Location: to.Ptr("centralus"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
	// 		ConnectionMode: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionModeDefault),
	// 		ConnectionProtocol: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
	// 		ConnectionType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionTypeIPsec),
	// 		DpdTimeoutSeconds: to.Ptr[int32](30),
	// 		EgressBytesTransferred: to.Ptr[int64](0),
	// 		EgressNatRules: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule2"),
	// 		}},
	// 		EnableBgp: to.Ptr(false),
	// 		GatewayCustomBgpIPAddresses: []*armnetwork.GatewayCustomBgpIPAddressIPConfiguration{
	// 			{
	// 				CustomBgpIPAddress: to.Ptr("169.254.21.1"),
	// 				IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/default"),
	// 			},
	// 			{
	// 				CustomBgpIPAddress: to.Ptr("169.254.21.3"),
	// 				IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/ActiveActive"),
	// 		}},
	// 		IngressBytesTransferred: to.Ptr[int64](0),
	// 		IngressNatRules: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule1"),
	// 		}},
	// 		IPSecPolicies: []*armnetwork.IPSecPolicy{
	// 		},
	// 		LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw"),
	// 			Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
	// 			},
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		RoutingWeight: to.Ptr[int32](0),
	// 		SharedKey: to.Ptr("Abc123"),
	// 		UseLocalAzureIPAddress: to.Ptr(false),
	// 		UsePolicyBasedTrafficSelectors: to.Ptr(false),
	// 		VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw"),
	// 			Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionGet.json
func ExampleVirtualNetworkGatewayConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().Get(ctx, "rg1", "connS2S", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualNetworkGatewayConnection = armnetwork.VirtualNetworkGatewayConnection{
	// 	Name: to.Ptr("connS2S"),
	// 	Type: to.Ptr("Microsoft.Network/connections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/connections/connS2S"),
	// 	Location: to.Ptr("centralus"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
	// 		ConnectionMode: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionModeDefault),
	// 		ConnectionProtocol: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
	// 		ConnectionStatus: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionStatusConnecting),
	// 		ConnectionType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionTypeIPsec),
	// 		DpdTimeoutSeconds: to.Ptr[int32](30),
	// 		EgressBytesTransferred: to.Ptr[int64](0),
	// 		EgressNatRules: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule2"),
	// 		}},
	// 		EnableBgp: to.Ptr(false),
	// 		GatewayCustomBgpIPAddresses: []*armnetwork.GatewayCustomBgpIPAddressIPConfiguration{
	// 			{
	// 				CustomBgpIPAddress: to.Ptr("169.254.21.1"),
	// 				IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/default"),
	// 			},
	// 			{
	// 				CustomBgpIPAddress: to.Ptr("169.254.21.3"),
	// 				IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/ActiveActive"),
	// 		}},
	// 		IngressBytesTransferred: to.Ptr[int64](0),
	// 		IngressNatRules: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule1"),
	// 		}},
	// 		IPSecPolicies: []*armnetwork.IPSecPolicy{
	// 		},
	// 		LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw"),
	// 			Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
	// 			},
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		RoutingWeight: to.Ptr[int32](0),
	// 		SharedKey: to.Ptr("Abc123"),
	// 		TrafficSelectorPolicies: []*armnetwork.TrafficSelectorPolicy{
	// 		},
	// 		UseLocalAzureIPAddress: to.Ptr(false),
	// 		UsePolicyBasedTrafficSelectors: to.Ptr(false),
	// 		VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw"),
	// 			Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionDelete.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginDelete(ctx, "rg1", "conn1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionUpdateTags.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginUpdateTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginUpdateTags(ctx, "rg1", "test", armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
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
	// res.VirtualNetworkGatewayConnection = armnetwork.VirtualNetworkGatewayConnection{
	// 	Name: to.Ptr("test"),
	// 	Type: to.Ptr("Microsoft.Network/connections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/connections/test"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"tag1": to.Ptr("value1"),
	// 		"tag2": to.Ptr("value2"),
	// 	},
	// 	Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
	// 		ConnectionStatus: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionStatusUnknown),
	// 		ConnectionType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionTypeIPsec),
	// 		EgressBytesTransferred: to.Ptr[int64](0),
	// 		EgressNatRules: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule2"),
	// 		}},
	// 		EnableBgp: to.Ptr(false),
	// 		GatewayCustomBgpIPAddresses: []*armnetwork.GatewayCustomBgpIPAddressIPConfiguration{
	// 			{
	// 				CustomBgpIPAddress: to.Ptr("169.254.21.1"),
	// 				IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/default"),
	// 			},
	// 			{
	// 				CustomBgpIPAddress: to.Ptr("169.254.21.3"),
	// 				IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/ActiveActive"),
	// 		}},
	// 		IngressBytesTransferred: to.Ptr[int64](0),
	// 		IngressNatRules: []*armnetwork.SubResource{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/natRules/natRule1"),
	// 		}},
	// 		IPSecPolicies: []*armnetwork.IPSecPolicy{
	// 		},
	// 		LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/lgw"),
	// 			Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
	// 			},
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		RoutingWeight: to.Ptr[int32](0),
	// 		SharedKey: to.Ptr("temp1234"),
	// 		TrafficSelectorPolicies: []*armnetwork.TrafficSelectorPolicy{
	// 		},
	// 		UsePolicyBasedTrafficSelectors: to.Ptr(false),
	// 		VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw"),
	// 			Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
	// 			},
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionSetSharedKey.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginSetSharedKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginSetSharedKey(ctx, "rg1", "connS2S", armnetwork.ConnectionSharedKey{
		Value: to.Ptr("AzureAbc123"),
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
	// res.ConnectionSharedKey = armnetwork.ConnectionSharedKey{
	// 	ID: to.Ptr(""),
	// 	Value: to.Ptr("AzureAbc123"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionGetSharedKey.json
func ExampleVirtualNetworkGatewayConnectionsClient_GetSharedKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().GetSharedKey(ctx, "rg1", "connS2S", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ConnectionSharedKey = armnetwork.ConnectionSharedKey{
	// 	ID: to.Ptr(""),
	// 	Value: to.Ptr("AzureAbc123"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionsList.json
func ExampleVirtualNetworkGatewayConnectionsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualNetworkGatewayConnectionsClient().NewListPager("rg1", nil)
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
		// page.VirtualNetworkGatewayConnectionListResult = armnetwork.VirtualNetworkGatewayConnectionListResult{
		// 	Value: []*armnetwork.VirtualNetworkGatewayConnection{
		// 		{
		// 			Name: to.Ptr("conn1"),
		// 			Type: to.Ptr("Microsoft.Network/connections"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/connections/conn1"),
		// 			Location: to.Ptr("centralus"),
		// 			Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 			Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
		// 				ConnectionMode: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionModeDefault),
		// 				ConnectionProtocol: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv1),
		// 				ConnectionType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionTypeIPsec),
		// 				DpdTimeoutSeconds: to.Ptr[int32](30),
		// 				EgressBytesTransferred: to.Ptr[int64](0),
		// 				EgressNatRules: []*armnetwork.SubResource{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw1/natRules/natRule2"),
		// 				}},
		// 				EnableBgp: to.Ptr(false),
		// 				GatewayCustomBgpIPAddresses: []*armnetwork.GatewayCustomBgpIPAddressIPConfiguration{
		// 					{
		// 						CustomBgpIPAddress: to.Ptr("169.254.21.1"),
		// 						IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/default"),
		// 					},
		// 					{
		// 						CustomBgpIPAddress: to.Ptr("169.254.21.3"),
		// 						IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw/ipConfigurations/ActiveActive"),
		// 				}},
		// 				IngressBytesTransferred: to.Ptr[int64](0),
		// 				IngressNatRules: []*armnetwork.SubResource{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw1/natRules/natRule1"),
		// 				}},
		// 				IPSecPolicies: []*armnetwork.IPSecPolicy{
		// 				},
		// 				LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw1"),
		// 					Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
		// 					},
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				RoutingWeight: to.Ptr[int32](0),
		// 				TrafficSelectorPolicies: []*armnetwork.TrafficSelectorPolicy{
		// 				},
		// 				UseLocalAzureIPAddress: to.Ptr(false),
		// 				UsePolicyBasedTrafficSelectors: to.Ptr(false),
		// 				VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw1"),
		// 					Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
		// 					},
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("conn2"),
		// 			Type: to.Ptr("Microsoft.Network/connections"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/connections/conn2"),
		// 			Location: to.Ptr("eastus"),
		// 			Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 			Properties: &armnetwork.VirtualNetworkGatewayConnectionPropertiesFormat{
		// 				ConnectionMode: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionModeDefault),
		// 				ConnectionProtocol: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionProtocolIKEv2),
		// 				ConnectionType: to.Ptr(armnetwork.VirtualNetworkGatewayConnectionTypeIPsec),
		// 				DpdTimeoutSeconds: to.Ptr[int32](20),
		// 				EgressBytesTransferred: to.Ptr[int64](0),
		// 				EgressNatRules: []*armnetwork.SubResource{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw2/natRules/natRule2"),
		// 				}},
		// 				EnableBgp: to.Ptr(false),
		// 				GatewayCustomBgpIPAddresses: []*armnetwork.GatewayCustomBgpIPAddressIPConfiguration{
		// 					{
		// 						CustomBgpIPAddress: to.Ptr("169.254.21.4"),
		// 						IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw2/ipConfigurations/default"),
		// 					},
		// 					{
		// 						CustomBgpIPAddress: to.Ptr("169.254.21.6"),
		// 						IPConfigurationID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw2/ipConfigurations/ActiveActive"),
		// 				}},
		// 				IngressBytesTransferred: to.Ptr[int64](0),
		// 				IngressNatRules: []*armnetwork.SubResource{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw2/natRules/natRule1"),
		// 				}},
		// 				IPSecPolicies: []*armnetwork.IPSecPolicy{
		// 				},
		// 				LocalNetworkGateway2: &armnetwork.LocalNetworkGateway{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw2"),
		// 					Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
		// 					},
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				RoutingWeight: to.Ptr[int32](0),
		// 				TrafficSelectorPolicies: []*armnetwork.TrafficSelectorPolicy{
		// 				},
		// 				UseLocalAzureIPAddress: to.Ptr(true),
		// 				UsePolicyBasedTrafficSelectors: to.Ptr(false),
		// 				VirtualNetworkGateway1: &armnetwork.VirtualNetworkGateway{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworkGateways/vpngw2"),
		// 					Properties: &armnetwork.VirtualNetworkGatewayPropertiesFormat{
		// 					},
		// 				},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionResetSharedKey.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginResetSharedKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginResetSharedKey(ctx, "rg1", "conn1", armnetwork.ConnectionResetSharedKey{
		KeyLength: to.Ptr[int32](128),
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
	// res.ConnectionResetSharedKey = armnetwork.ConnectionResetSharedKey{
	// 	KeyLength: to.Ptr[int32](128),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionStartPacketCaptureFilterData.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginStartPacketCapture_startPacketCaptureOnVirtualNetworkGatewayConnectionWithFilter() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginStartPacketCapture(ctx, "rg1", "vpngwcn1", &armnetwork.VirtualNetworkGatewayConnectionsClientBeginStartPacketCaptureOptions{Parameters: &armnetwork.VPNPacketCaptureStartParameters{
		FilterData: to.Ptr("{'TracingFlags': 11,'MaxPacketBufferSize': 120,'MaxFileSize': 200,'Filters': [{'SourceSubnets': ['20.1.1.0/24'],'DestinationSubnets': ['10.1.1.0/24'],'SourcePort': [500],'DestinationPort': [4500],'Protocol': 6,'TcpFlags': 16,'CaptureSingleDirectionTrafficOnly': true}]}"),
	},
	})
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
	// res.Value = "\"{\"Status\":\"Successful\",\"Data\":null}\""
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionStartPacketCapture.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginStartPacketCapture_startPacketCaptureOnVirtualNetworkGatewayConnectionWithoutFilter() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginStartPacketCapture(ctx, "rg1", "vpngwcn1", &armnetwork.VirtualNetworkGatewayConnectionsClientBeginStartPacketCaptureOptions{Parameters: nil})
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
	// res.Value = "\"{\"Status\":\"Successful\",\"Data\":null}\""
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionStopPacketCapture.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginStopPacketCapture() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginStopPacketCapture(ctx, "rg1", "vpngwcn1", armnetwork.VPNPacketCaptureStopParameters{
		SasURL: to.Ptr("https://teststorage.blob.core.windows.net/?sv=2018-03-28&ss=bfqt&srt=sco&sp=rwdlacup&se=2019-09-13T07:44:05Z&st=2019-09-06T23:44:05Z&spr=https&sig=V1h9D1riltvZMI69d6ihENnFo%2FrCvTqGgjO2lf%2FVBhE%3D"),
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
	// res.Value = "\"{\"Status\":\"Successful\",\"Data\":null}\""
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionGetIkeSas.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginGetIkeSas() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginGetIkeSas(ctx, "rg1", "vpngwcn1", nil)
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
	// res.Value = "\"{\"Status\":\"Successful\",\"Data\":null}\""
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/VirtualNetworkGatewayConnectionReset.json
func ExampleVirtualNetworkGatewayConnectionsClient_BeginResetConnection() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualNetworkGatewayConnectionsClient().BeginResetConnection(ctx, "rg1", "conn1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
