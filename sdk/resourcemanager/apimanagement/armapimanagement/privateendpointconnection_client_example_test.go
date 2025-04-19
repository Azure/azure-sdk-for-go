//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armapimanagement_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListPrivateEndpointConnections.json
func ExamplePrivateEndpointConnectionClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewPrivateEndpointConnectionClient().NewListByServicePager("rg1", "apimService1", nil)
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
		// page.PrivateEndpointConnectionListResult = armapimanagement.PrivateEndpointConnectionListResult{
		// 	Value: []*armapimanagement.PrivateEndpointConnection{
		// 		{
		// 			Name: to.Ptr("privateEndpointProxyName"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/privateEndpointConnections"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/privateEndpointConnections/connectionName"),
		// 			Properties: &armapimanagement.PrivateEndpointConnectionProperties{
		// 				PrivateEndpoint: &armapimanagement.PrivateEndpoint{
		// 					ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/privateEndpointName"),
		// 				},
		// 				PrivateLinkServiceConnectionState: &armapimanagement.PrivateLinkServiceConnectionState{
		// 					Description: to.Ptr("Please approve my request, thanks"),
		// 					ActionsRequired: to.Ptr("None"),
		// 					Status: to.Ptr(armapimanagement.PrivateEndpointServiceConnectionStatusPending),
		// 				},
		// 				ProvisioningState: to.Ptr(armapimanagement.PrivateEndpointConnectionProvisioningStateSucceeded),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("privateEndpointProxyName2"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/privateEndpointConnections"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/privateEndpointConnections/privateEndpointProxyName2"),
		// 			Properties: &armapimanagement.PrivateEndpointConnectionProperties{
		// 				PrivateEndpoint: &armapimanagement.PrivateEndpoint{
		// 					ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/privateEndpointName2"),
		// 				},
		// 				PrivateLinkServiceConnectionState: &armapimanagement.PrivateLinkServiceConnectionState{
		// 					Description: to.Ptr("Please approve my request, thanks"),
		// 					ActionsRequired: to.Ptr("None"),
		// 					Status: to.Ptr(armapimanagement.PrivateEndpointServiceConnectionStatusPending),
		// 				},
		// 				ProvisioningState: to.Ptr(armapimanagement.PrivateEndpointConnectionProvisioningStateSucceeded),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetPrivateEndpointConnection.json
func ExamplePrivateEndpointConnectionClient_GetByName() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateEndpointConnectionClient().GetByName(ctx, "rg1", "apimService1", "privateEndpointConnectionName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateEndpointConnection = armapimanagement.PrivateEndpointConnection{
	// 	Name: to.Ptr("privateEndpointProxyName"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/privateEndpointConnections/privateEndpointConnectionName"),
	// 	Properties: &armapimanagement.PrivateEndpointConnectionProperties{
	// 		PrivateEndpoint: &armapimanagement.PrivateEndpoint{
	// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/privateEndpointName"),
	// 		},
	// 		PrivateLinkServiceConnectionState: &armapimanagement.PrivateLinkServiceConnectionState{
	// 			Description: to.Ptr("Please approve my request, thanks"),
	// 			ActionsRequired: to.Ptr("None"),
	// 			Status: to.Ptr(armapimanagement.PrivateEndpointServiceConnectionStatusPending),
	// 		},
	// 		ProvisioningState: to.Ptr(armapimanagement.PrivateEndpointConnectionProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementApproveOrRejectPrivateEndpointConnection.json
func ExamplePrivateEndpointConnectionClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPrivateEndpointConnectionClient().BeginCreateOrUpdate(ctx, "rg1", "apimService1", "privateEndpointConnectionName", armapimanagement.PrivateEndpointConnectionRequest{
		ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/privateEndpointConnections/connectionName"),
		Properties: &armapimanagement.PrivateEndpointConnectionRequestProperties{
			PrivateLinkServiceConnectionState: &armapimanagement.PrivateLinkServiceConnectionState{
				Description: to.Ptr("The Private Endpoint Connection is approved."),
				Status:      to.Ptr(armapimanagement.PrivateEndpointServiceConnectionStatusApproved),
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
	// res.PrivateEndpointConnection = armapimanagement.PrivateEndpointConnection{
	// 	Name: to.Ptr("privateEndpointConnectionName"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/privateEndpointConnections/privateEndpointConnectionName"),
	// 	Properties: &armapimanagement.PrivateEndpointConnectionProperties{
	// 		PrivateEndpoint: &armapimanagement.PrivateEndpoint{
	// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/privateEndpointName"),
	// 		},
	// 		PrivateLinkServiceConnectionState: &armapimanagement.PrivateLinkServiceConnectionState{
	// 			Description: to.Ptr("The request has been approved."),
	// 			Status: to.Ptr(armapimanagement.PrivateEndpointServiceConnectionStatus("Succeeded")),
	// 		},
	// 		ProvisioningState: to.Ptr(armapimanagement.PrivateEndpointConnectionProvisioningStateSucceeded),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeletePrivateEndpointConnection.json
func ExamplePrivateEndpointConnectionClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPrivateEndpointConnectionClient().BeginDelete(ctx, "rg1", "apimService1", "privateEndpointConnectionName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListPrivateLinkGroupResources.json
func ExamplePrivateEndpointConnectionClient_ListPrivateLinkResources() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateEndpointConnectionClient().ListPrivateLinkResources(ctx, "rg1", "apimService1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateLinkResourceListResult = armapimanagement.PrivateLinkResourceListResult{
	// 	Value: []*armapimanagement.PrivateLinkResource{
	// 		{
	// 			Name: to.Ptr("Gateway"),
	// 			Type: to.Ptr("Microsoft.ApiManagement/service/privateLinkResources"),
	// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/privateLinkResources/Gateway"),
	// 			Properties: &armapimanagement.PrivateLinkResourceProperties{
	// 				GroupID: to.Ptr("Gateway"),
	// 				RequiredMembers: []*string{
	// 					to.Ptr("Gateway")},
	// 					RequiredZoneNames: []*string{
	// 						to.Ptr("privateLink.azure-api.net")},
	// 					},
	// 			}},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetPrivateLinkGroupResource.json
func ExamplePrivateEndpointConnectionClient_GetPrivateLinkResource() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateEndpointConnectionClient().GetPrivateLinkResource(ctx, "rg1", "apimService1", "privateLinkSubResourceName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateLinkResource = armapimanagement.PrivateLinkResource{
	// 	Name: to.Ptr("Gateway"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/privateLinkResources"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/privateLinkResources/Gateway"),
	// 	Properties: &armapimanagement.PrivateLinkResourceProperties{
	// 		GroupID: to.Ptr("Gateway"),
	// 		RequiredMembers: []*string{
	// 			to.Ptr("Gateway")},
	// 			RequiredZoneNames: []*string{
	// 				to.Ptr("privateLink.azure-api.net")},
	// 			},
	// 		}
}
