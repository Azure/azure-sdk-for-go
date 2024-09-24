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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/NetworkManagerGet.json
func ExampleManagersClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagersClient().Get(ctx, "rg1", "testNetworkManager", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Manager = armnetwork.Manager{
	// 	Name: to.Ptr("testNetworkManager"),
	// 	Type: to.Ptr("Microsoft.Network/networkManagers"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/rg1/providers/Microsoft.Network/networkManagers/testNetworkManager"),
	// 	Properties: &armnetwork.ManagerProperties{
	// 		Description: to.Ptr("My Test Network Manager"),
	// 		NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
	// 			to.Ptr(armnetwork.ConfigurationTypeSecurityUser)},
	// 			NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
	// 				ManagementGroups: []*string{
	// 				},
	// 				Subscriptions: []*string{
	// 					to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000")},
	// 				},
	// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 				ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 			},
	// 			SystemData: &armnetwork.SystemData{
	// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 				CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 				CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 				LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 				LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/NetworkManagerPut.json
func ExampleManagersClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagersClient().CreateOrUpdate(ctx, "rg1", "TestNetworkManager", armnetwork.Manager{
		Properties: &armnetwork.ManagerProperties{
			Description: to.Ptr("My Test Network Manager"),
			NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
				to.Ptr(armnetwork.ConfigurationTypeConnectivity)},
			NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
				ManagementGroups: []*string{
					to.Ptr("/Microsoft.Management/testmg")},
				Subscriptions: []*string{
					to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000")},
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Manager = armnetwork.Manager{
	// 	Name: to.Ptr("TestNetworkManager"),
	// 	Type: to.Ptr("Microsoft.Network/networkManagers"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/rg1/providers/Microsoft.Network/networkManagers/TestNetworkManager"),
	// 	Etag: to.Ptr("sadf-asdf-asdf-asdf"),
	// 	Properties: &armnetwork.ManagerProperties{
	// 		Description: to.Ptr("My Test Network Manager"),
	// 		NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
	// 			to.Ptr(armnetwork.ConfigurationTypeConnectivity)},
	// 			NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
	// 				ManagementGroups: []*string{
	// 					to.Ptr("Microsoft.Management/managementGroups/testMg")},
	// 					Subscriptions: []*string{
	// 						to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000")},
	// 					},
	// 					ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 					ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 				},
	// 				SystemData: &armnetwork.SystemData{
	// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 					CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 					CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 					LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 					LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 				},
	// 			}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/NetworkManagerDelete.json
func ExampleManagersClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewManagersClient().BeginDelete(ctx, "rg1", "testNetworkManager", &armnetwork.ManagersClientBeginDeleteOptions{Force: to.Ptr(false)})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/NetworkManagerPatch.json
func ExampleManagersClient_Patch() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagersClient().Patch(ctx, "rg1", "testNetworkManager", armnetwork.PatchObject{
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
	// res.Manager = armnetwork.Manager{
	// 	Name: to.Ptr("testNetworkManager"),
	// 	Type: to.Ptr("Microsoft.Network/networkManager"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/rg1/providers/Microsoft.Network/networkManagers/testNetworkManager"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"tag1": to.Ptr("value1"),
	// 		"tag2": to.Ptr("value2"),
	// 	},
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Properties: &armnetwork.ManagerProperties{
	// 		Description: to.Ptr("My Test Network Manager"),
	// 		NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
	// 			to.Ptr(armnetwork.ConfigurationTypeConnectivity)},
	// 			NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
	// 				ManagementGroups: []*string{
	// 				},
	// 				Subscriptions: []*string{
	// 					to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000001")},
	// 				},
	// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 				ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 			},
	// 			SystemData: &armnetwork.SystemData{
	// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 				CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 				CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 				LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 				LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/NetworkManagerListAll.json
func ExampleManagersClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewManagersClient().NewListBySubscriptionPager(&armnetwork.ManagersClientListBySubscriptionOptions{Top: nil,
		SkipToken: nil,
	})
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
		// page.ManagerListResult = armnetwork.ManagerListResult{
		// 	Value: []*armnetwork.Manager{
		// 		{
		// 			Name: to.Ptr("testNetworkManager"),
		// 			Type: to.Ptr("Microsoft.Network/networkManagers"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/rg1/providers/Microsoft.Network/networkManagers/testNetworkManager"),
		// 			Etag: to.Ptr("sadf-asdf-asdf-asdf"),
		// 			Properties: &armnetwork.ManagerProperties{
		// 				Description: to.Ptr("My Test Network Manager"),
		// 				NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
		// 					to.Ptr(armnetwork.ConfigurationTypeSecurityUser)},
		// 					NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
		// 						ManagementGroups: []*string{
		// 						},
		// 						Subscriptions: []*string{
		// 							to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000")},
		// 						},
		// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 						ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 					},
		// 					SystemData: &armnetwork.SystemData{
		// 						CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
		// 						CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
		// 						CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 						LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
		// 						LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
		// 						LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4883fa5dbf6f2c9093fac8ce334547e9dfac68fa/specification/network/resource-manager/Microsoft.Network/stable/2024-03-01/examples/NetworkManagerList.json
func ExampleManagersClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewManagersClient().NewListPager("rg1", &armnetwork.ManagersClientListOptions{Top: nil,
		SkipToken: nil,
	})
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
		// page.ManagerListResult = armnetwork.ManagerListResult{
		// 	Value: []*armnetwork.Manager{
		// 		{
		// 			Name: to.Ptr("testNetworkManager"),
		// 			Type: to.Ptr("Microsoft.Network/networkManagers"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroup/rg1/providers/Microsoft.Network/networkManagers/testNetworkManager"),
		// 			Etag: to.Ptr("sadf-asdf-asdf-asdf"),
		// 			Properties: &armnetwork.ManagerProperties{
		// 				Description: to.Ptr("My Test Network Manager"),
		// 				NetworkManagerScopeAccesses: []*armnetwork.ConfigurationType{
		// 					to.Ptr(armnetwork.ConfigurationTypeConnectivity)},
		// 					NetworkManagerScopes: &armnetwork.ManagerPropertiesNetworkManagerScopes{
		// 						ManagementGroups: []*string{
		// 						},
		// 						Subscriptions: []*string{
		// 							to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000")},
		// 						},
		// 						ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 						ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 					},
		// 					SystemData: &armnetwork.SystemData{
		// 						CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
		// 						CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
		// 						CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 						LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
		// 						LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
		// 						LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 					},
		// 			}},
		// 		}
	}
}
