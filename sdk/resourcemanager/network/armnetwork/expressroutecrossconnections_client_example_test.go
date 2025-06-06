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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionList.json
func ExampleExpressRouteCrossConnectionsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewExpressRouteCrossConnectionsClient().NewListPager(&armnetwork.ExpressRouteCrossConnectionsClientListOptions{Filter: nil})
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
		// page.ExpressRouteCrossConnectionListResult = armnetwork.ExpressRouteCrossConnectionListResult{
		// 	Value: []*armnetwork.ExpressRouteCrossConnection{
		// 		{
		// 			Name: to.Ptr("<circuitServiceKey>"),
		// 			Type: to.Ptr("Microsoft.Network/expressRouteCrossConnections"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/CrossConnectionSiliconValley/providers/Microsoft.Network/expressRouteCrossConnections/<circuitServiceKey>"),
		// 			Location: to.Ptr("brazilsouth"),
		// 			Properties: &armnetwork.ExpressRouteCrossConnectionProperties{
		// 				BandwidthInMbps: to.Ptr[int32](1000),
		// 				ExpressRouteCircuit: &armnetwork.ExpressRouteCircuitReference{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/ertest/providers/Microsoft.Network/expressRouteCircuits/er1"),
		// 				},
		// 				PeeringLocation: to.Ptr("SiliconValley"),
		// 				Peerings: []*armnetwork.ExpressRouteCrossConnectionPeering{
		// 				},
		// 				PrimaryAzurePort: to.Ptr("bvtazureixp01"),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				STag: to.Ptr[int32](2),
		// 				SecondaryAzurePort: to.Ptr("bvtazureixp01"),
		// 				ServiceProviderProvisioningState: to.Ptr(armnetwork.ServiceProviderProvisioningStateNotProvisioned),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionListByResourceGroup.json
func ExampleExpressRouteCrossConnectionsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewExpressRouteCrossConnectionsClient().NewListByResourceGroupPager("CrossConnection-SiliconValley", nil)
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
		// page.ExpressRouteCrossConnectionListResult = armnetwork.ExpressRouteCrossConnectionListResult{
		// 	Value: []*armnetwork.ExpressRouteCrossConnection{
		// 		{
		// 			Name: to.Ptr("<circuitServiceKey>"),
		// 			Type: to.Ptr("Microsoft.Network/expressRouteCrossConnections"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/CrossConnectionSilicon-Valley/providers/Microsoft.Network/expressRouteCrossConnections/<circuitServiceKey>"),
		// 			Location: to.Ptr("brazilsouth"),
		// 			Properties: &armnetwork.ExpressRouteCrossConnectionProperties{
		// 				BandwidthInMbps: to.Ptr[int32](1000),
		// 				ExpressRouteCircuit: &armnetwork.ExpressRouteCircuitReference{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/ertest/providers/Microsoft.Network/expressRouteCircuits/er1"),
		// 				},
		// 				PeeringLocation: to.Ptr("SiliconValley"),
		// 				Peerings: []*armnetwork.ExpressRouteCrossConnectionPeering{
		// 				},
		// 				PrimaryAzurePort: to.Ptr("bvtazureixp01"),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				STag: to.Ptr[int32](2),
		// 				SecondaryAzurePort: to.Ptr("bvtazureixp01"),
		// 				ServiceProviderProvisioningState: to.Ptr(armnetwork.ServiceProviderProvisioningStateNotProvisioned),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionGet.json
func ExampleExpressRouteCrossConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewExpressRouteCrossConnectionsClient().Get(ctx, "CrossConnection-SiliconValley", "<circuitServiceKey>", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ExpressRouteCrossConnection = armnetwork.ExpressRouteCrossConnection{
	// 	Name: to.Ptr("<circuitServiceKey>"),
	// 	Type: to.Ptr("Microsoft.Network/expressRouteCrossConnections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/CrossConnection-SiliconValley/providers/Microsoft.Network/expressRouteCrossConnections/<circuitServiceKey>"),
	// 	Location: to.Ptr("brazilsouth"),
	// 	Etag: to.Ptr("W/\"c0e6477e-8150-4d4f-9bf6-bb10e6acb63a\""),
	// 	Properties: &armnetwork.ExpressRouteCrossConnectionProperties{
	// 		BandwidthInMbps: to.Ptr[int32](1000),
	// 		ExpressRouteCircuit: &armnetwork.ExpressRouteCircuitReference{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/ertest/providers/Microsoft.Network/expressRouteCircuits/er1"),
	// 		},
	// 		PeeringLocation: to.Ptr("SiliconValley"),
	// 		Peerings: []*armnetwork.ExpressRouteCrossConnectionPeering{
	// 		},
	// 		PrimaryAzurePort: to.Ptr("bvtazureixp01"),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		STag: to.Ptr[int32](2),
	// 		SecondaryAzurePort: to.Ptr("bvtazureixp01"),
	// 		ServiceProviderProvisioningState: to.Ptr(armnetwork.ServiceProviderProvisioningStateNotProvisioned),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionUpdate.json
func ExampleExpressRouteCrossConnectionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExpressRouteCrossConnectionsClient().BeginCreateOrUpdate(ctx, "CrossConnection-SiliconValley", "<circuitServiceKey>", armnetwork.ExpressRouteCrossConnection{
		Properties: &armnetwork.ExpressRouteCrossConnectionProperties{
			ServiceProviderProvisioningState: to.Ptr(armnetwork.ServiceProviderProvisioningStateNotProvisioned),
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
	// res.ExpressRouteCrossConnection = armnetwork.ExpressRouteCrossConnection{
	// 	Name: to.Ptr("<circuitServiceKey>"),
	// 	Type: to.Ptr("Microsoft.Network/expressRouteCrossConnections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/CrossConnectionSiliconValley/providers/Microsoft.Network/expressRouteCrossConnections/<circuitServiceKey>"),
	// 	Location: to.Ptr("brazilsouth"),
	// 	Properties: &armnetwork.ExpressRouteCrossConnectionProperties{
	// 		BandwidthInMbps: to.Ptr[int32](1000),
	// 		ExpressRouteCircuit: &armnetwork.ExpressRouteCircuitReference{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/ertest/providers/Microsoft.Network/expressRouteCircuits/er1"),
	// 		},
	// 		PeeringLocation: to.Ptr("SiliconValley"),
	// 		Peerings: []*armnetwork.ExpressRouteCrossConnectionPeering{
	// 		},
	// 		PrimaryAzurePort: to.Ptr("bvtazureixp01"),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		STag: to.Ptr[int32](2),
	// 		SecondaryAzurePort: to.Ptr("bvtazureixp01"),
	// 		ServiceProviderProvisioningState: to.Ptr(armnetwork.ServiceProviderProvisioningStateNotProvisioned),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionUpdateTags.json
func ExampleExpressRouteCrossConnectionsClient_UpdateTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewExpressRouteCrossConnectionsClient().UpdateTags(ctx, "CrossConnection-SiliconValley", "<circuitServiceKey>", armnetwork.TagsObject{
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
	// res.ExpressRouteCrossConnection = armnetwork.ExpressRouteCrossConnection{
	// 	Name: to.Ptr("er1"),
	// 	Type: to.Ptr("Microsoft.Network/expressRouteCrossConnections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/CrossConnectionSiliconValley/providers/Microsoft.Network/expressRouteCrossConnections/<circuitServiceKey>"),
	// 	Location: to.Ptr("brazilsouth"),
	// 	Tags: map[string]*string{
	// 		"tag1": to.Ptr("value1"),
	// 		"tag2": to.Ptr("value2"),
	// 	},
	// 	Properties: &armnetwork.ExpressRouteCrossConnectionProperties{
	// 		BandwidthInMbps: to.Ptr[int32](1000),
	// 		ExpressRouteCircuit: &armnetwork.ExpressRouteCircuitReference{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/ertest/providers/Microsoft.Network/expressRouteCircuits/er1"),
	// 		},
	// 		PeeringLocation: to.Ptr("SiliconValley"),
	// 		Peerings: []*armnetwork.ExpressRouteCrossConnectionPeering{
	// 		},
	// 		PrimaryAzurePort: to.Ptr("bvtazureixp01"),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		STag: to.Ptr[int32](2),
	// 		SecondaryAzurePort: to.Ptr("bvtazureixp01"),
	// 		ServiceProviderProvisioningState: to.Ptr(armnetwork.ServiceProviderProvisioningStateNotProvisioned),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionsArpTable.json
func ExampleExpressRouteCrossConnectionsClient_BeginListArpTable() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExpressRouteCrossConnectionsClient().BeginListArpTable(ctx, "CrossConnection-SiliconValley", "<circuitServiceKey>", "AzurePrivatePeering", "primary", nil)
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
	// res.ExpressRouteCircuitsArpTableListResult = armnetwork.ExpressRouteCircuitsArpTableListResult{
	// 	Value: []*armnetwork.ExpressRouteCircuitArpTable{
	// 		{
	// 			Age: to.Ptr[int32](0),
	// 			Interface: to.Ptr("Microsoft"),
	// 			IPAddress: to.Ptr("192.116.14.254"),
	// 			MacAddress: to.Ptr("885a.9269.9110"),
	// 	}},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionsRouteTableSummary.json
func ExampleExpressRouteCrossConnectionsClient_BeginListRoutesTableSummary() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExpressRouteCrossConnectionsClient().BeginListRoutesTableSummary(ctx, "CrossConnection-SiliconValley", "<circuitServiceKey>", "AzurePrivatePeering", "primary", nil)
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
	// res.ExpressRouteCrossConnectionsRoutesTableSummaryListResult = armnetwork.ExpressRouteCrossConnectionsRoutesTableSummaryListResult{
	// 	Value: []*armnetwork.ExpressRouteCrossConnectionRoutesTableSummary{
	// 		{
	// 			Asn: to.Ptr[int32](65514),
	// 			Neighbor: to.Ptr("10.6.1.112"),
	// 			StateOrPrefixesReceived: to.Ptr("Active"),
	// 			UpDown: to.Ptr("1d14h"),
	// 		},
	// 		{
	// 			Asn: to.Ptr[int32](65514),
	// 			Neighbor: to.Ptr("10.6.1.113"),
	// 			StateOrPrefixesReceived: to.Ptr("1"),
	// 			UpDown: to.Ptr("1d14h"),
	// 	}},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/ExpressRouteCrossConnectionsRouteTable.json
func ExampleExpressRouteCrossConnectionsClient_BeginListRoutesTable() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewExpressRouteCrossConnectionsClient().BeginListRoutesTable(ctx, "CrossConnection-SiliconValley", "<circuitServiceKey>", "AzurePrivatePeering", "primary", nil)
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
	// res.ExpressRouteCircuitsRoutesTableListResult = armnetwork.ExpressRouteCircuitsRoutesTableListResult{
	// 	Value: []*armnetwork.ExpressRouteCircuitRoutesTable{
	// 		{
	// 			Path: to.Ptr("65514"),
	// 			LocPrf: to.Ptr(""),
	// 			Network: to.Ptr("10.6.0.0/16"),
	// 			NextHop: to.Ptr("10.6.1.12"),
	// 			Weight: to.Ptr[int32](0),
	// 		},
	// 		{
	// 			Path: to.Ptr("65514"),
	// 			LocPrf: to.Ptr(""),
	// 			Network: to.Ptr("10.7.0.0/16"),
	// 			NextHop: to.Ptr("10.7.1.13"),
	// 			Weight: to.Ptr[int32](0),
	// 	}},
	// }
}
