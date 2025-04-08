//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armeventgrid_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/PrivateEndpointConnections_Get.json
func ExamplePrivateEndpointConnectionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateEndpointConnectionsClient().Get(ctx, "examplerg", armeventgrid.PrivateEndpointConnectionsParentTypeTopics, "exampletopic1", "BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateEndpointConnection = armeventgrid.PrivateEndpointConnection{
	// 	Name: to.Ptr("BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B"),
	// 	Type: to.Ptr("Microsoft.EventGrid/topics/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/5B4B650E-28B9-4790-B3AB-DDBD88D727C4/resourceGroups/examplerg/providers/Microsoft.EventGrid/topics/exampletopic1/privateEndpointConnections/BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B"),
	// 	Properties: &armeventgrid.PrivateEndpointConnectionProperties{
	// 		GroupIDs: []*string{
	// 			to.Ptr("topic")},
	// 			PrivateEndpoint: &armeventgrid.PrivateEndpoint{
	// 				ID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.Network/privateEndpoints/bmtpe5"),
	// 			},
	// 			PrivateLinkServiceConnectionState: &armeventgrid.ConnectionState{
	// 				Description: to.Ptr("Test"),
	// 				ActionsRequired: to.Ptr("None"),
	// 				Status: to.Ptr(armeventgrid.PersistedConnectionStatusPending),
	// 			},
	// 			ProvisioningState: to.Ptr(armeventgrid.ResourceProvisioningStateSucceeded),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/PrivateEndpointConnections_Update.json
func ExamplePrivateEndpointConnectionsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPrivateEndpointConnectionsClient().BeginUpdate(ctx, "examplerg", armeventgrid.PrivateEndpointConnectionsParentTypeTopics, "exampletopic1", "BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B", armeventgrid.PrivateEndpointConnection{
		Properties: &armeventgrid.PrivateEndpointConnectionProperties{
			PrivateLinkServiceConnectionState: &armeventgrid.ConnectionState{
				Description:     to.Ptr("approving connection"),
				ActionsRequired: to.Ptr("None"),
				Status:          to.Ptr(armeventgrid.PersistedConnectionStatusApproved),
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
	// res.PrivateEndpointConnection = armeventgrid.PrivateEndpointConnection{
	// 	Name: to.Ptr("BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B"),
	// 	Type: to.Ptr("Microsoft.EventGrid/topics/privateEndpointConnections"),
	// 	ID: to.Ptr("/subscriptions/5B4B650E-28B9-4790-B3AB-DDBD88D727C4/resourceGroups/examplerg/providers/Microsoft.EventGrid/topics/exampletopic1/privateEndpointConnections/BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B"),
	// 	Properties: &armeventgrid.PrivateEndpointConnectionProperties{
	// 		GroupIDs: []*string{
	// 			to.Ptr("topic")},
	// 			PrivateEndpoint: &armeventgrid.PrivateEndpoint{
	// 				ID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.Network/privateEndpoints/bmtpe5"),
	// 			},
	// 			PrivateLinkServiceConnectionState: &armeventgrid.ConnectionState{
	// 				Description: to.Ptr("approving connection"),
	// 				ActionsRequired: to.Ptr("None"),
	// 				Status: to.Ptr(armeventgrid.PersistedConnectionStatusApproved),
	// 			},
	// 			ProvisioningState: to.Ptr(armeventgrid.ResourceProvisioningStateSucceeded),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/PrivateEndpointConnections_Delete.json
func ExamplePrivateEndpointConnectionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPrivateEndpointConnectionsClient().BeginDelete(ctx, "examplerg", armeventgrid.PrivateEndpointConnectionsParentTypeTopics, "exampletopic1", "BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/ee1eec42dcc710ff88db2d1bf574b2f9afe3d654/specification/eventgrid/resource-manager/Microsoft.EventGrid/stable/2025-02-15/examples/PrivateEndpointConnections_ListByResource.json
func ExamplePrivateEndpointConnectionsClient_NewListByResourcePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewPrivateEndpointConnectionsClient().NewListByResourcePager("examplerg", armeventgrid.PrivateEndpointConnectionsParentTypeTopics, "exampletopic1", &armeventgrid.PrivateEndpointConnectionsClientListByResourceOptions{Filter: nil,
		Top: nil,
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
		// page.PrivateEndpointConnectionListResult = armeventgrid.PrivateEndpointConnectionListResult{
		// 	Value: []*armeventgrid.PrivateEndpointConnection{
		// 		{
		// 			Name: to.Ptr("BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B"),
		// 			Type: to.Ptr("Microsoft.EventGrid/topics/privateEndpointConnections"),
		// 			ID: to.Ptr("/subscriptions/5B4B650E-28B9-4790-B3AB-DDBD88D727C4/resourceGroups/examplerg/providers/Microsoft.EventGrid/topics/exampletopic1/privateEndpointConnections/BMTPE5.8A30D251-4C61-489D-A1AA-B37C4A329B8B"),
		// 			Properties: &armeventgrid.PrivateEndpointConnectionProperties{
		// 				GroupIDs: []*string{
		// 					to.Ptr("topic")},
		// 					PrivateEndpoint: &armeventgrid.PrivateEndpoint{
		// 						ID: to.Ptr("/subscriptions/5b4b650e-28b9-4790-b3ab-ddbd88d727c4/resourceGroups/examplerg/providers/Microsoft.Network/privateEndpoints/bmtpe5"),
		// 					},
		// 					PrivateLinkServiceConnectionState: &armeventgrid.ConnectionState{
		// 						Description: to.Ptr("Test"),
		// 						ActionsRequired: to.Ptr("None"),
		// 						Status: to.Ptr(armeventgrid.PersistedConnectionStatusPending),
		// 					},
		// 					ProvisioningState: to.Ptr(armeventgrid.ResourceProvisioningStateSucceeded),
		// 				},
		// 		}},
		// 	}
	}
}
