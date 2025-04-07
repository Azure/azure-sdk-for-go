//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsearch_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/search/armsearch"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/search/resource-manager/Microsoft.Search/preview/2024-06-01-preview/examples/UpdatePrivateEndpointConnection.json
func ExamplePrivateEndpointConnectionsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsearch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateEndpointConnectionsClient().Update(ctx, "rg1", "mysearchservice", "testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546", armsearch.PrivateEndpointConnection{
		Properties: &armsearch.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armsearch.PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState{
				Description: to.Ptr("Rejected for some reason."),
				Status:      to.Ptr(armsearch.PrivateLinkServiceConnectionStatusRejected),
			},
		},
	}, &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateEndpointConnection = armsearch.PrivateEndpointConnection{
	// 	Name: to.Ptr("testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
	// 	Type: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Search/searchServices/mysearchservice/privateEndpointConnections/testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
	// 	Properties: &armsearch.PrivateEndpointConnectionProperties{
	// 		PrivateEndpoint: &armsearch.PrivateEndpointConnectionPropertiesPrivateEndpoint{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/testEndpoint"),
	// 		},
	// 		PrivateLinkServiceConnectionState: &armsearch.PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState{
	// 			Description: to.Ptr("Rejected for some reason."),
	// 			ActionsRequired: to.Ptr("None"),
	// 			Status: to.Ptr(armsearch.PrivateLinkServiceConnectionStatusRejected),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/search/resource-manager/Microsoft.Search/preview/2024-06-01-preview/examples/GetPrivateEndpointConnection.json
func ExamplePrivateEndpointConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsearch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateEndpointConnectionsClient().Get(ctx, "rg1", "mysearchservice", "testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546", &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateEndpointConnection = armsearch.PrivateEndpointConnection{
	// 	Name: to.Ptr("testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
	// 	Type: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Search/searchServices/mysearchservice/privateEndpointConnections/testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
	// 	Properties: &armsearch.PrivateEndpointConnectionProperties{
	// 		PrivateEndpoint: &armsearch.PrivateEndpointConnectionPropertiesPrivateEndpoint{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/testEndpoint"),
	// 		},
	// 		PrivateLinkServiceConnectionState: &armsearch.PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState{
	// 			Description: to.Ptr(""),
	// 			ActionsRequired: to.Ptr("None"),
	// 			Status: to.Ptr(armsearch.PrivateLinkServiceConnectionStatusApproved),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/search/resource-manager/Microsoft.Search/preview/2024-06-01-preview/examples/DeletePrivateEndpointConnection.json
func ExamplePrivateEndpointConnectionsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsearch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateEndpointConnectionsClient().Delete(ctx, "rg1", "mysearchservice", "testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546", &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateEndpointConnection = armsearch.PrivateEndpointConnection{
	// 	Name: to.Ptr("testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
	// 	Type: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Search/searchServices/mysearchservice/privateEndpointConnections/testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
	// 	Properties: &armsearch.PrivateEndpointConnectionProperties{
	// 		PrivateEndpoint: &armsearch.PrivateEndpointConnectionPropertiesPrivateEndpoint{
	// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/testEndpoint"),
	// 		},
	// 		PrivateLinkServiceConnectionState: &armsearch.PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState{
	// 			Description: to.Ptr(""),
	// 			ActionsRequired: to.Ptr("None"),
	// 			Status: to.Ptr(armsearch.PrivateLinkServiceConnectionStatusDisconnected),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/search/resource-manager/Microsoft.Search/preview/2024-06-01-preview/examples/ListPrivateEndpointConnectionsByService.json
func ExamplePrivateEndpointConnectionsClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsearch.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewPrivateEndpointConnectionsClient().NewListByServicePager("rg1", "mysearchservice", &armsearch.SearchManagementRequestOptions{ClientRequestID: nil}, nil)
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
		// page.PrivateEndpointConnectionListResult = armsearch.PrivateEndpointConnectionListResult{
		// 	Value: []*armsearch.PrivateEndpointConnection{
		// 		{
		// 			Name: to.Ptr("testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
		// 			Type: to.Ptr("Microsoft.Search/searchServices/privateEndpointConnections"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Search/searchServices/mysearchservice/privateEndpointConnections/testEndpoint.50bf4fbe-d7c1-4b48-a642-4f5892642546"),
		// 			Properties: &armsearch.PrivateEndpointConnectionProperties{
		// 				PrivateEndpoint: &armsearch.PrivateEndpointConnectionPropertiesPrivateEndpoint{
		// 					ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/privateEndpoints/testEndpoint"),
		// 				},
		// 				PrivateLinkServiceConnectionState: &armsearch.PrivateEndpointConnectionPropertiesPrivateLinkServiceConnectionState{
		// 					Description: to.Ptr(""),
		// 					ActionsRequired: to.Ptr("None"),
		// 					Status: to.Ptr(armsearch.PrivateLinkServiceConnectionStatusApproved),
		// 				},
		// 			},
		// 	}},
		// }
	}
}
