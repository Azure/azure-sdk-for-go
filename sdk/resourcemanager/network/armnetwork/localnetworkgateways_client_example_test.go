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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/LocalNetworkGatewayCreate.json
func ExampleLocalNetworkGatewaysClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetwork.NewLocalNetworkGatewaysClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreateOrUpdate(ctx, "rg1", "localgw", armnetwork.LocalNetworkGateway{
		Location: to.Ptr("Central US"),
		Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
			Fqdn:             to.Ptr("site1.contoso.com"),
			GatewayIPAddress: to.Ptr("11.12.13.14"),
			LocalNetworkAddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.1.0.0/16")},
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
	// res.LocalNetworkGateway = armnetwork.LocalNetworkGateway{
	// 	Name: to.Ptr("localgw"),
	// 	Type: to.Ptr("Microsoft.Network/localNetworkGateways"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw"),
	// 	Location: to.Ptr("centralus"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
	// 		GatewayIPAddress: to.Ptr("11.12.13.14"),
	// 		LocalNetworkAddressSpace: &armnetwork.AddressSpace{
	// 			AddressPrefixes: []*string{
	// 				to.Ptr("10.1.0.0/16")},
	// 			},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/LocalNetworkGatewayGet.json
func ExampleLocalNetworkGatewaysClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetwork.NewLocalNetworkGatewaysClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx, "rg1", "localgw", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.LocalNetworkGateway = armnetwork.LocalNetworkGateway{
	// 	Name: to.Ptr("localgw"),
	// 	Type: to.Ptr("Microsoft.Network/localNetworkGateways"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw"),
	// 	Location: to.Ptr("centralus"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
	// 		GatewayIPAddress: to.Ptr("x.x.x.x"),
	// 		LocalNetworkAddressSpace: &armnetwork.AddressSpace{
	// 			AddressPrefixes: []*string{
	// 				to.Ptr("10.1.0.0/16")},
	// 			},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/LocalNetworkGatewayDelete.json
func ExampleLocalNetworkGatewaysClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetwork.NewLocalNetworkGatewaysClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx, "rg1", "localgw", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/LocalNetworkGatewayUpdateTags.json
func ExampleLocalNetworkGatewaysClient_UpdateTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetwork.NewLocalNetworkGatewaysClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.UpdateTags(ctx, "rg1", "lgw", armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.LocalNetworkGateway = armnetwork.LocalNetworkGateway{
	// 	Name: to.Ptr("lgw"),
	// 	Type: to.Ptr("Microsoft.Network/localNetworkGateways"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/lgw"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"tag1": to.Ptr("value1"),
	// 		"tag2": to.Ptr("value2"),
	// 	},
	// 	Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
	// 		GatewayIPAddress: to.Ptr("12.0.0.1"),
	// 		LocalNetworkAddressSpace: &armnetwork.AddressSpace{
	// 			AddressPrefixes: []*string{
	// 				to.Ptr("12.0.0.0/8")},
	// 			},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/a60468a0c5e2beb054680ae488fb9f92699f0a0d/specification/network/resource-manager/Microsoft.Network/stable/2022-09-01/examples/LocalNetworkGatewayList.json
func ExampleLocalNetworkGatewaysClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armnetwork.NewLocalNetworkGatewaysClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager("rg1", nil)
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
		// page.LocalNetworkGatewayListResult = armnetwork.LocalNetworkGatewayListResult{
		// 	Value: []*armnetwork.LocalNetworkGateway{
		// 		{
		// 			Name: to.Ptr("localgw1"),
		// 			Type: to.Ptr("Microsoft.Network/localNetworkGateways"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw1"),
		// 			Location: to.Ptr("centralus"),
		// 			Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 			Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
		// 				GatewayIPAddress: to.Ptr("x.x.x.x"),
		// 				LocalNetworkAddressSpace: &armnetwork.AddressSpace{
		// 					AddressPrefixes: []*string{
		// 						to.Ptr("10.1.0.0/16")},
		// 					},
		// 					ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 					ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("localgw2"),
		// 				Type: to.Ptr("Microsoft.Network/localNetworkGateways"),
		// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/localNetworkGateways/localgw2"),
		// 				Location: to.Ptr("eastus"),
		// 				Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 				Properties: &armnetwork.LocalNetworkGatewayPropertiesFormat{
		// 					GatewayIPAddress: to.Ptr("x.x.x.x"),
		// 					LocalNetworkAddressSpace: &armnetwork.AddressSpace{
		// 						AddressPrefixes: []*string{
		// 							to.Ptr("10.2.0.0/16")},
		// 						},
		// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 						ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 					},
		// 			}},
		// 		}
	}
}
