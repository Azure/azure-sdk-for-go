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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/bf204aab860f2eb58a9d346b00d44760f2a9b0a2/specification/eventgrid/resource-manager/Microsoft.EventGrid/preview/2023-12-15-preview/examples/PermissionBindings_Get.json
func ExamplePermissionBindingsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPermissionBindingsClient().Get(ctx, "examplerg", "exampleNamespaceName1", "examplePermissionBindingName1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PermissionBinding = armeventgrid.PermissionBinding{
	// 	Name: to.Ptr("examplePermissionBindingName1"),
	// 	Type: to.Ptr("Microsoft.EventGrid/namespaces/permissionBindings"),
	// 	ID: to.Ptr("/subscriptions/8f6b6269-84f2-4d09-9e31-1127efcd1e40/resourceGroups/examplerg/providers/Microsoft.EventGrid/namespaces/exampleNamespaceName1/permissionBindings/examplePermissionBindingName1"),
	// 	Properties: &armeventgrid.PermissionBindingProperties{
	// 		ClientGroupName: to.Ptr("exampleClientGroupName1"),
	// 		Permission: to.Ptr(armeventgrid.PermissionTypePublisher),
	// 		ProvisioningState: to.Ptr(armeventgrid.PermissionBindingProvisioningStateSucceeded),
	// 		TopicSpaceName: to.Ptr("exampleTopicSpaceName1"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/bf204aab860f2eb58a9d346b00d44760f2a9b0a2/specification/eventgrid/resource-manager/Microsoft.EventGrid/preview/2023-12-15-preview/examples/PermissionBindings_CreateOrUpdate.json
func ExamplePermissionBindingsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPermissionBindingsClient().BeginCreateOrUpdate(ctx, "examplerg", "exampleNamespaceName1", "examplePermissionBindingName1", armeventgrid.PermissionBinding{
		Properties: &armeventgrid.PermissionBindingProperties{
			ClientGroupName: to.Ptr("exampleClientGroupName1"),
			Permission:      to.Ptr(armeventgrid.PermissionTypePublisher),
			TopicSpaceName:  to.Ptr("exampleTopicSpaceName1"),
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
	// res.PermissionBinding = armeventgrid.PermissionBinding{
	// 	Name: to.Ptr("examplePermissionBindingName1"),
	// 	Type: to.Ptr("Microsoft.EventGrid/namespaces/permissionBindings"),
	// 	ID: to.Ptr("/subscriptions/8f6b6269-84f2-4d09-9e31-1127efcd1e40/resourceGroups/examplerg/providers/Microsoft.EventGrid/namespaces/exampleNamespaceName1/permissionBindings/examplePermissionBindingName1"),
	// 	Properties: &armeventgrid.PermissionBindingProperties{
	// 		ClientGroupName: to.Ptr("exampleClientGroupName1"),
	// 		Permission: to.Ptr(armeventgrid.PermissionTypePublisher),
	// 		ProvisioningState: to.Ptr(armeventgrid.PermissionBindingProvisioningStateSucceeded),
	// 		TopicSpaceName: to.Ptr("exampleTopicSpaceName1"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/bf204aab860f2eb58a9d346b00d44760f2a9b0a2/specification/eventgrid/resource-manager/Microsoft.EventGrid/preview/2023-12-15-preview/examples/PermissionBindings_Delete.json
func ExamplePermissionBindingsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewPermissionBindingsClient().BeginDelete(ctx, "examplerg", "exampleNamespaceName1", "examplePermissionBindingName1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/bf204aab860f2eb58a9d346b00d44760f2a9b0a2/specification/eventgrid/resource-manager/Microsoft.EventGrid/preview/2023-12-15-preview/examples/PermissionBindings_ListByNamespace.json
func ExamplePermissionBindingsClient_NewListByNamespacePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armeventgrid.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewPermissionBindingsClient().NewListByNamespacePager("examplerg", "namespace123", &armeventgrid.PermissionBindingsClientListByNamespaceOptions{Filter: nil,
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
		// page.PermissionBindingsListResult = armeventgrid.PermissionBindingsListResult{
		// 	Value: []*armeventgrid.PermissionBinding{
		// 		{
		// 			Name: to.Ptr("examplePermissionBindingName1"),
		// 			Type: to.Ptr("Microsoft.EventGrid/namespaces/permissionBindings"),
		// 			ID: to.Ptr("/subscriptions/8f6b6269-84f2-4d09-9e31-1127efcd1e40/resourceGroups/examplerg/providers/Microsoft.EventGrid/namespaces/exampleNamespaceName1/permissionBindings/examplePermissionBindingName1"),
		// 			Properties: &armeventgrid.PermissionBindingProperties{
		// 				ClientGroupName: to.Ptr("exampleClientGroupName1"),
		// 				Permission: to.Ptr(armeventgrid.PermissionTypePublisher),
		// 				ProvisioningState: to.Ptr(armeventgrid.PermissionBindingProvisioningStateSucceeded),
		// 				TopicSpaceName: to.Ptr("exampleTopicSpaceName1"),
		// 			},
		// 	}},
		// }
	}
}