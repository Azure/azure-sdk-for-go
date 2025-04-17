//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcompute_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccess_ListBySubscription.json
func ExampleDiskAccessesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDiskAccessesClient().NewListPager(nil)
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
		// page.DiskAccessList = armcompute.DiskAccessList{
		// 	Value: []*armcompute.DiskAccess{
		// 		{
		// 			Name: to.Ptr("myDiskAccess"),
		// 			Type: to.Ptr("Microsoft.Compute/diskAccesses"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess"),
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"department": to.Ptr("Development"),
		// 				"project": to.Ptr("PrivateEndpoints"),
		// 			},
		// 			Properties: &armcompute.DiskAccessProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T04:41:35.079Z"); return t}()),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("myDiskAccess2"),
		// 			Type: to.Ptr("Microsoft.Compute/diskAccesses"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/mySecondResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess2"),
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"department": to.Ptr("Development"),
		// 				"project": to.Ptr("PrivateEndpoints"),
		// 			},
		// 			Properties: &armcompute.DiskAccessProperties{
		// 				PrivateEndpointConnections: []*armcompute.PrivateEndpointConnection{
		// 					{
		// 						Name: to.Ptr("myDiskAccess.d4914cfa-6bc2-4049-a57c-3d1f622d8eef"),
		// 						Type: to.Ptr("Microsoft.Compute/diskAccesses/PrivateEndpointConnections"),
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/mySecondResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess2/privateEndpoinConnections/myDiskAccess2.d4914cfa-6bc2-4049-a57c-3d1f622d8eef"),
		// 						Properties: &armcompute.PrivateEndpointConnectionProperties{
		// 							PrivateEndpoint: &armcompute.PrivateEndpoint{
		// 								ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/mySecondResourceGroup/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint2"),
		// 							},
		// 							PrivateLinkServiceConnectionState: &armcompute.PrivateLinkServiceConnectionState{
		// 								Description: to.Ptr("Auto-Approved"),
		// 								ActionsRequired: to.Ptr("None"),
		// 								Status: to.Ptr(armcompute.PrivateEndpointServiceConnectionStatusApproved),
		// 							},
		// 							ProvisioningState: to.Ptr(armcompute.PrivateEndpointConnectionProvisioningStateSucceeded),
		// 						},
		// 				}},
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T04:41:35.079Z"); return t}()),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccess_ListByResourceGroup.json
func ExampleDiskAccessesClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDiskAccessesClient().NewListByResourceGroupPager("myResourceGroup", nil)
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
		// page.DiskAccessList = armcompute.DiskAccessList{
		// 	Value: []*armcompute.DiskAccess{
		// 		{
		// 			Name: to.Ptr("myDiskAccess"),
		// 			Type: to.Ptr("Microsoft.Compute/diskAccesses"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess"),
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"department": to.Ptr("Development"),
		// 				"project": to.Ptr("PrivateEndpoints"),
		// 			},
		// 			Properties: &armcompute.DiskAccessProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T04:41:35.079Z"); return t}()),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("myDiskAccess2"),
		// 			Type: to.Ptr("Microsoft.Compute/diskAccesses"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess2"),
		// 			Location: to.Ptr("westus"),
		// 			Tags: map[string]*string{
		// 				"department": to.Ptr("Development"),
		// 				"project": to.Ptr("PrivateEndpoints"),
		// 			},
		// 			Properties: &armcompute.DiskAccessProperties{
		// 				PrivateEndpointConnections: []*armcompute.PrivateEndpointConnection{
		// 					{
		// 						Name: to.Ptr("myDiskAccess.d4914cfa-6bc2-4049-a57c-3d1f622d8eef"),
		// 						Type: to.Ptr("Microsoft.Compute/diskAccesses/PrivateEndpointConnections"),
		// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess2/privateEndpoinConnections/myDiskAccess2.d4914cfa-6bc2-4049-a57c-3d1f622d8eef"),
		// 						Properties: &armcompute.PrivateEndpointConnectionProperties{
		// 							PrivateEndpoint: &armcompute.PrivateEndpoint{
		// 								ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint2"),
		// 							},
		// 							PrivateLinkServiceConnectionState: &armcompute.PrivateLinkServiceConnectionState{
		// 								Description: to.Ptr("Auto-Approved"),
		// 								ActionsRequired: to.Ptr("None"),
		// 								Status: to.Ptr(armcompute.PrivateEndpointServiceConnectionStatusApproved),
		// 							},
		// 							ProvisioningState: to.Ptr(armcompute.PrivateEndpointConnectionProvisioningStateSucceeded),
		// 						},
		// 				}},
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T04:41:35.079Z"); return t}()),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccess_Get_WithPrivateEndpoints.json
func ExampleDiskAccessesClient_Get_getInformationAboutADiskAccessResourceWithPrivateEndpoints() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDiskAccessesClient().Get(ctx, "myResourceGroup", "myDiskAccess", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DiskAccess = armcompute.DiskAccess{
	// 	Name: to.Ptr("myDiskAccess"),
	// 	Type: to.Ptr("Microsoft.Compute/diskAccesses"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"department": to.Ptr("Development"),
	// 		"project": to.Ptr("PrivateEndpoints"),
	// 	},
	// 	Properties: &armcompute.DiskAccessProperties{
	// 		PrivateEndpointConnections: []*armcompute.PrivateEndpointConnection{
	// 			{
	// 				Name: to.Ptr("myDiskAccess.d4914cfa-6bc2-4049-a57c-3d1f622d8eef"),
	// 				Type: to.Ptr("Microsoft.Compute/diskAccesses/PrivateEndpointConnections"),
	// 				ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess/privateEndpoinConnections/myDiskAccess.d4914cfa-6bc2-4049-a57c-3d1f622d8eef"),
	// 				Properties: &armcompute.PrivateEndpointConnectionProperties{
	// 					PrivateEndpoint: &armcompute.PrivateEndpoint{
	// 						ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint"),
	// 					},
	// 					PrivateLinkServiceConnectionState: &armcompute.PrivateLinkServiceConnectionState{
	// 						Description: to.Ptr("Auto-Approved"),
	// 						ActionsRequired: to.Ptr("None"),
	// 						Status: to.Ptr(armcompute.PrivateEndpointServiceConnectionStatusApproved),
	// 					},
	// 					ProvisioningState: to.Ptr(armcompute.PrivateEndpointConnectionProvisioningStateSucceeded),
	// 				},
	// 		}},
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T04:41:35.079Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccess_Get.json
func ExampleDiskAccessesClient_Get_getInformationAboutADiskAccessResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDiskAccessesClient().Get(ctx, "myResourceGroup", "myDiskAccess", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DiskAccess = armcompute.DiskAccess{
	// 	Name: to.Ptr("myDiskAccess"),
	// 	Type: to.Ptr("Microsoft.Compute/diskAccesses"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess"),
	// 	Location: to.Ptr("westus"),
	// 	Tags: map[string]*string{
	// 		"department": to.Ptr("Development"),
	// 		"project": to.Ptr("PrivateEndpoints"),
	// 	},
	// 	Properties: &armcompute.DiskAccessProperties{
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T04:41:35.079Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccess_Create.json
func ExampleDiskAccessesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDiskAccessesClient().BeginCreateOrUpdate(ctx, "myResourceGroup", "myDiskAccess", armcompute.DiskAccess{
		Location: to.Ptr("West US"),
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
	// res.DiskAccess = armcompute.DiskAccess{
	// 	Name: to.Ptr("myDiskAccess"),
	// 	Type: to.Ptr("Microsoft.Compute/diskAccesses"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourcegroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess"),
	// 	Location: to.Ptr("West US"),
	// 	Properties: &armcompute.DiskAccessProperties{
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		TimeCreated: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T04:41:35.079Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccess_Update.json
func ExampleDiskAccessesClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDiskAccessesClient().BeginUpdate(ctx, "myResourceGroup", "myDiskAccess", armcompute.DiskAccessUpdate{
		Tags: map[string]*string{
			"department": to.Ptr("Development"),
			"project":    to.Ptr("PrivateEndpoints"),
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
	// res.DiskAccess = armcompute.DiskAccess{
	// 	Name: to.Ptr("myDiskAccess"),
	// 	Type: to.Ptr("Microsoft.Compute/diskAccesses"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourcegroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess"),
	// 	Location: to.Ptr("West US"),
	// 	Tags: map[string]*string{
	// 		"department": to.Ptr("Development"),
	// 		"project": to.Ptr("PrivateEndpoints"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccess_Delete.json
func ExampleDiskAccessesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDiskAccessesClient().BeginDelete(ctx, "myResourceGroup", "myDiskAccess", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccessPrivateEndpointConnection_ListByDiskAccess.json
func ExampleDiskAccessesClient_NewListPrivateEndpointConnectionsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDiskAccessesClient().NewListPrivateEndpointConnectionsPager("myResourceGroup", "myDiskAccess", nil)
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
		// page.PrivateEndpointConnectionListResult = armcompute.PrivateEndpointConnectionListResult{
		// 	Value: []*armcompute.PrivateEndpointConnection{
		// 		{
		// 			Name: to.Ptr("myPrivateEndpointConnection"),
		// 			Type: to.Ptr("Microsoft.Compute/diskAccesses/PrivateEndpointConnections"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess/privateEndpoinConnections/myPrivateEndpointConnection"),
		// 			Properties: &armcompute.PrivateEndpointConnectionProperties{
		// 				PrivateEndpoint: &armcompute.PrivateEndpoint{
		// 					ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint"),
		// 				},
		// 				PrivateLinkServiceConnectionState: &armcompute.PrivateLinkServiceConnectionState{
		// 					Description: to.Ptr("Auto-Approved"),
		// 					ActionsRequired: to.Ptr("None"),
		// 					Status: to.Ptr(armcompute.PrivateEndpointServiceConnectionStatusApproved),
		// 				},
		// 				ProvisioningState: to.Ptr(armcompute.PrivateEndpointConnectionProvisioningStateSucceeded),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccessPrivateEndpointConnection_Get.json
func ExampleDiskAccessesClient_GetAPrivateEndpointConnection() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDiskAccessesClient().GetAPrivateEndpointConnection(ctx, "myResourceGroup", "myDiskAccess", "myPrivateEndpointConnection", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateEndpointConnection = armcompute.PrivateEndpointConnection{
	// 	Name: to.Ptr("myPrivateEndpointConnection"),
	// 	Type: to.Ptr("Microsoft.Compute/diskAccesses/PrivateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess/privateEndpoinConnections/myPrivateEndpointConnection"),
	// 	Properties: &armcompute.PrivateEndpointConnectionProperties{
	// 		PrivateEndpoint: &armcompute.PrivateEndpoint{
	// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint"),
	// 		},
	// 		PrivateLinkServiceConnectionState: &armcompute.PrivateLinkServiceConnectionState{
	// 			Description: to.Ptr("Auto-Approved"),
	// 			ActionsRequired: to.Ptr("None"),
	// 			Status: to.Ptr(armcompute.PrivateEndpointServiceConnectionStatusApproved),
	// 		},
	// 		ProvisioningState: to.Ptr(armcompute.PrivateEndpointConnectionProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccessPrivateEndpointConnection_Approve.json
func ExampleDiskAccessesClient_BeginUpdateAPrivateEndpointConnection() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDiskAccessesClient().BeginUpdateAPrivateEndpointConnection(ctx, "myResourceGroup", "myDiskAccess", "myPrivateEndpointConnection", armcompute.PrivateEndpointConnection{
		Properties: &armcompute.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armcompute.PrivateLinkServiceConnectionState{
				Description: to.Ptr("Approving myPrivateEndpointConnection"),
				Status:      to.Ptr(armcompute.PrivateEndpointServiceConnectionStatusApproved),
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
	// res.PrivateEndpointConnection = armcompute.PrivateEndpointConnection{
	// 	Name: to.Ptr("myPrivateEndpointConnectionName"),
	// 	Type: to.Ptr("Microsoft.Compute/diskAccesses/PrivateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess/privateEndpoinConnections/myPrivateEndpointConnectionName"),
	// 	Properties: &armcompute.PrivateEndpointConnectionProperties{
	// 		PrivateEndpoint: &armcompute.PrivateEndpoint{
	// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myPrivateEndpoint"),
	// 		},
	// 		PrivateLinkServiceConnectionState: &armcompute.PrivateLinkServiceConnectionState{
	// 			Description: to.Ptr("Approving myPrivateEndpointConnection"),
	// 			ActionsRequired: to.Ptr("None"),
	// 			Status: to.Ptr(armcompute.PrivateEndpointServiceConnectionStatusApproved),
	// 		},
	// 		ProvisioningState: to.Ptr(armcompute.PrivateEndpointConnectionProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccessPrivateEndpointConnection_Delete.json
func ExampleDiskAccessesClient_BeginDeleteAPrivateEndpointConnection() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDiskAccessesClient().BeginDeleteAPrivateEndpointConnection(ctx, "myResourceGroup", "myDiskAccess", "myPrivateEndpointConnection", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/0e21cbbb89e3cee48cfb544e1e322cb13a3080da/specification/compute/resource-manager/Microsoft.Compute/DiskRP/stable/2024-03-02/examples/diskAccessExamples/DiskAccessPrivateLinkResources_Get.json
func ExampleDiskAccessesClient_GetPrivateLinkResources() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDiskAccessesClient().GetPrivateLinkResources(ctx, "myResourceGroup", "myDiskAccess", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateLinkResourceListResult = armcompute.PrivateLinkResourceListResult{
	// 	Value: []*armcompute.PrivateLinkResource{
	// 		{
	// 			Name: to.Ptr("disks"),
	// 			Type: to.Ptr("Microsoft.Compute/diskAccesses/privateLinkResources"),
	// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Compute/diskAccesses/myDiskAccess/privateLinkResources/disks"),
	// 			Properties: &armcompute.PrivateLinkResourceProperties{
	// 				GroupID: to.Ptr("disks"),
	// 				RequiredMembers: []*string{
	// 					to.Ptr("diskAccess_1")},
	// 					RequiredZoneNames: []*string{
	// 						to.Ptr("privatelink.blob.core.windows.net")},
	// 					},
	// 			}},
	// 		}
}
