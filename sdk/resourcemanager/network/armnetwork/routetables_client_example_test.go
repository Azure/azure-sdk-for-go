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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/RouteTableDelete.json
func ExampleRouteTablesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewRouteTablesClient().BeginDelete(ctx, "rg1", "testrt", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/RouteTableGet.json
func ExampleRouteTablesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewRouteTablesClient().Get(ctx, "rg1", "testrt", &armnetwork.RouteTablesClientGetOptions{Expand: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.RouteTable = armnetwork.RouteTable{
	// 	Name: to.Ptr("testrt"),
	// 	Type: to.Ptr("Microsoft.Network/routeTables"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armnetwork.RouteTablePropertiesFormat{
	// 		DisableBgpRoutePropagation: to.Ptr(false),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		Routes: []*armnetwork.Route{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt/routes/route1"),
	// 				Name: to.Ptr("route1"),
	// 				Properties: &armnetwork.RoutePropertiesFormat{
	// 					AddressPrefix: to.Ptr("10.0.3.0/24"),
	// 					NextHopType: to.Ptr(armnetwork.RouteNextHopTypeVirtualNetworkGateway),
	// 					ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 				},
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/RouteTableCreate.json
func ExampleRouteTablesClient_BeginCreateOrUpdate_createRouteTable() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewRouteTablesClient().BeginCreateOrUpdate(ctx, "rg1", "testrt", armnetwork.RouteTable{
		Location: to.Ptr("westus"),
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
	// res.RouteTable = armnetwork.RouteTable{
	// 	Name: to.Ptr("testrt"),
	// 	Type: to.Ptr("Microsoft.Network/routeTables"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armnetwork.RouteTablePropertiesFormat{
	// 		DisableBgpRoutePropagation: to.Ptr(true),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		Routes: []*armnetwork.Route{
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/RouteTableCreateWithRoute.json
func ExampleRouteTablesClient_BeginCreateOrUpdate_createRouteTableWithRoute() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewRouteTablesClient().BeginCreateOrUpdate(ctx, "rg1", "testrt", armnetwork.RouteTable{
		Location: to.Ptr("westus"),
		Properties: &armnetwork.RouteTablePropertiesFormat{
			DisableBgpRoutePropagation: to.Ptr(true),
			Routes: []*armnetwork.Route{
				{
					Name: to.Ptr("route1"),
					Properties: &armnetwork.RoutePropertiesFormat{
						AddressPrefix: to.Ptr("10.0.3.0/24"),
						NextHopType:   to.Ptr(armnetwork.RouteNextHopTypeVirtualNetworkGateway),
					},
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
	// res.RouteTable = armnetwork.RouteTable{
	// 	Name: to.Ptr("testrt"),
	// 	Type: to.Ptr("Microsoft.Network/routeTables"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armnetwork.RouteTablePropertiesFormat{
	// 		DisableBgpRoutePropagation: to.Ptr(true),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		Routes: []*armnetwork.Route{
	// 			{
	// 				ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt/routes/route1"),
	// 				Name: to.Ptr("route1"),
	// 				Properties: &armnetwork.RoutePropertiesFormat{
	// 					AddressPrefix: to.Ptr("10.0.3.0/24"),
	// 					NextHopType: to.Ptr(armnetwork.RouteNextHopTypeVirtualNetworkGateway),
	// 					ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 				},
	// 		}},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/RouteTableUpdateTags.json
func ExampleRouteTablesClient_UpdateTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewRouteTablesClient().UpdateTags(ctx, "rg1", "testrt", armnetwork.TagsObject{
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
	// res.RouteTable = armnetwork.RouteTable{
	// 	Name: to.Ptr("testrt"),
	// 	Type: to.Ptr("Microsoft.Network/routeTables"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"tag1": to.Ptr("value1"),
	// 		"tag2": to.Ptr("value2"),
	// 	},
	// 	Properties: &armnetwork.RouteTablePropertiesFormat{
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		Routes: []*armnetwork.Route{
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/RouteTableList.json
func ExampleRouteTablesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewRouteTablesClient().NewListPager("rg1", nil)
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
		// page.RouteTableListResult = armnetwork.RouteTableListResult{
		// 	Value: []*armnetwork.RouteTable{
		// 		{
		// 			Name: to.Ptr("testrt"),
		// 			Type: to.Ptr("Microsoft.Network/routeTables"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armnetwork.RouteTablePropertiesFormat{
		// 				DisableBgpRoutePropagation: to.Ptr(true),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				Routes: []*armnetwork.Route{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt/routes/route1"),
		// 						Name: to.Ptr("route1"),
		// 						Properties: &armnetwork.RoutePropertiesFormat{
		// 							AddressPrefix: to.Ptr("10.0.3.0/24"),
		// 							NextHopType: to.Ptr(armnetwork.RouteNextHopTypeVirtualNetworkGateway),
		// 							ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 						},
		// 				}},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("testrt2"),
		// 			Type: to.Ptr("Microsoft.Network/routeTables"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt2"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armnetwork.RouteTablePropertiesFormat{
		// 				DisableBgpRoutePropagation: to.Ptr(true),
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				Routes: []*armnetwork.Route{
		// 				},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4be63be2cabdebeb6974b22d89ed6fdb44392541/specification/network/resource-manager/Microsoft.Network/stable/2023-09-01/examples/RouteTableListAll.json
func ExampleRouteTablesClient_NewListAllPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewRouteTablesClient().NewListAllPager(nil)
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
		// page.RouteTableListResult = armnetwork.RouteTableListResult{
		// 	Value: []*armnetwork.RouteTable{
		// 		{
		// 			Name: to.Ptr("testrt"),
		// 			Type: to.Ptr("Microsoft.Network/routeTables"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armnetwork.RouteTablePropertiesFormat{
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				Routes: []*armnetwork.Route{
		// 					{
		// 						ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/routeTables/testrt/routes/route1"),
		// 						Name: to.Ptr("route1"),
		// 						Properties: &armnetwork.RoutePropertiesFormat{
		// 							AddressPrefix: to.Ptr("10.0.3.0/24"),
		// 							NextHopType: to.Ptr(armnetwork.RouteNextHopTypeVirtualNetworkGateway),
		// 							ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 						},
		// 				}},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("testrt3"),
		// 			Type: to.Ptr("Microsoft.Network/routeTables"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg2/providers/Microsoft.Network/routeTables/testrt3"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armnetwork.RouteTablePropertiesFormat{
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				Routes: []*armnetwork.Route{
		// 				},
		// 			},
		// 	}},
		// }
	}
}
