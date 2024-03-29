//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdashboard_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dashboard/armdashboard"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/78eac0bd58633028293cb1ec1709baa200bed9e2/specification/dashboard/resource-manager/Microsoft.Dashboard/stable/2023-09-01/examples/ManagedPrivateEndpoints_List.json
func ExampleManagedPrivateEndpointsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdashboard.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewManagedPrivateEndpointsClient().NewListPager("myResourceGroup", "myWorkspace", nil)
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
		// page.ManagedPrivateEndpointModelListResponse = armdashboard.ManagedPrivateEndpointModelListResponse{
		// 	Value: []*armdashboard.ManagedPrivateEndpointModel{
		// 		{
		// 			Name: to.Ptr("myMPEName"),
		// 			Type: to.Ptr("Microsoft.Dashboard/grafana/managedPrivateEndpoint"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/Microsoft.Dashboard/grafana/myWorkspace/managedPrivateEndpoints/myPrivateEndpointName"),
		// 			Location: to.Ptr("West US"),
		// 			Properties: &armdashboard.ManagedPrivateEndpointModelProperties{
		// 				ConnectionState: &armdashboard.ManagedPrivateEndpointConnectionState{
		// 					Description: to.Ptr("Auto-Approved"),
		// 					Status: to.Ptr(armdashboard.ManagedPrivateEndpointConnectionStatusApproved),
		// 				},
		// 				GroupIDs: []*string{
		// 					to.Ptr("grafana")},
		// 					PrivateLinkResourceID: to.Ptr("/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-000000000000/resourceGroups/xx-rg/providers/Microsoft.Kusto/Clusters/sampleKustoResource1"),
		// 					PrivateLinkResourceRegion: to.Ptr("West US"),
		// 					PrivateLinkServicePrivateIP: to.Ptr("10.0.0.5"),
		// 					PrivateLinkServiceURL: to.Ptr("my-self-hosted-influxdb.westus.mydomain.com"),
		// 					ProvisioningState: to.Ptr(armdashboard.ProvisioningStateSucceeded),
		// 					RequestMessage: to.Ptr("Example Request Message"),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("myMPEName2"),
		// 				Type: to.Ptr("Microsoft.Dashboard/grafana/managedPrivateEndpoint"),
		// 				ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/Microsoft.Dashboard/grafana/myWorkspace/managedPrivateEndpoints/myPrivateEndpointName2"),
		// 				Location: to.Ptr("West US"),
		// 				Properties: &armdashboard.ManagedPrivateEndpointModelProperties{
		// 					ConnectionState: &armdashboard.ManagedPrivateEndpointConnectionState{
		// 						Description: to.Ptr("Example Reject Reason"),
		// 						Status: to.Ptr(armdashboard.ManagedPrivateEndpointConnectionStatusRejected),
		// 					},
		// 					GroupIDs: []*string{
		// 						to.Ptr("grafana")},
		// 						PrivateLinkResourceID: to.Ptr("/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-000000000000/resourceGroups/xx-rg/providers/Microsoft.Kusto/Clusters/sampleKustoResource2"),
		// 						PrivateLinkResourceRegion: to.Ptr("West US"),
		// 						ProvisioningState: to.Ptr(armdashboard.ProvisioningStateSucceeded),
		// 						RequestMessage: to.Ptr("Example Request Message 2"),
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/78eac0bd58633028293cb1ec1709baa200bed9e2/specification/dashboard/resource-manager/Microsoft.Dashboard/stable/2023-09-01/examples/ManagedPrivateEndpoints_Refresh.json
func ExampleManagedPrivateEndpointsClient_BeginRefresh() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdashboard.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewManagedPrivateEndpointsClient().BeginRefresh(ctx, "myResourceGroup", "myWorkspace", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/78eac0bd58633028293cb1ec1709baa200bed9e2/specification/dashboard/resource-manager/Microsoft.Dashboard/stable/2023-09-01/examples/ManagedPrivateEndpoints_Get.json
func ExampleManagedPrivateEndpointsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdashboard.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewManagedPrivateEndpointsClient().Get(ctx, "myResourceGroup", "myWorkspace", "myMPEName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ManagedPrivateEndpointModel = armdashboard.ManagedPrivateEndpointModel{
	// 	Name: to.Ptr("myMPEName"),
	// 	Type: to.Ptr("Microsoft.Dashboard/grafana/managedPrivateEndpoint"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/Microsoft.Dashboard/grafana/myWorkspace/managedPrivateEndpoints/myPrivateEndpointName"),
	// 	SystemData: &armdashboard.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armdashboard.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.197Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armdashboard.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("West US"),
	// 	Properties: &armdashboard.ManagedPrivateEndpointModelProperties{
	// 		ConnectionState: &armdashboard.ManagedPrivateEndpointConnectionState{
	// 			Description: to.Ptr("Auto-Approved"),
	// 			Status: to.Ptr(armdashboard.ManagedPrivateEndpointConnectionStatusApproved),
	// 		},
	// 		GroupIDs: []*string{
	// 			to.Ptr("grafana")},
	// 			PrivateLinkResourceID: to.Ptr("/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-000000000000/resourceGroups/xx-rg/providers/Microsoft.Kusto/Clusters/sampleKustoResource"),
	// 			PrivateLinkResourceRegion: to.Ptr("West US"),
	// 			ProvisioningState: to.Ptr(armdashboard.ProvisioningStateSucceeded),
	// 			RequestMessage: to.Ptr("Example Request Message"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/78eac0bd58633028293cb1ec1709baa200bed9e2/specification/dashboard/resource-manager/Microsoft.Dashboard/stable/2023-09-01/examples/ManagedPrivateEndpoints_Create.json
func ExampleManagedPrivateEndpointsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdashboard.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewManagedPrivateEndpointsClient().BeginCreate(ctx, "myResourceGroup", "myWorkspace", "myMPEName", armdashboard.ManagedPrivateEndpointModel{
		Location: to.Ptr("West US"),
		Properties: &armdashboard.ManagedPrivateEndpointModelProperties{
			GroupIDs: []*string{
				to.Ptr("grafana")},
			PrivateLinkResourceID:     to.Ptr("/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-000000000000/resourceGroups/xx-rg/providers/Microsoft.Kusto/Clusters/sampleKustoResource"),
			PrivateLinkResourceRegion: to.Ptr("West US"),
			PrivateLinkServiceURL:     to.Ptr("my-self-hosted-influxdb.westus.mydomain.com"),
			RequestMessage:            to.Ptr("Example Request Message"),
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
	// res.ManagedPrivateEndpointModel = armdashboard.ManagedPrivateEndpointModel{
	// 	Name: to.Ptr("myMPEName"),
	// 	Type: to.Ptr("Microsoft.Dashboard/grafana/managedPrivateEndpoint"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/Microsoft.Dashboard/grafana/myWorkspace/managedPrivateEndpoints/myPrivateEndpointName"),
	// 	SystemData: &armdashboard.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armdashboard.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.197Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armdashboard.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("West US"),
	// 	Properties: &armdashboard.ManagedPrivateEndpointModelProperties{
	// 		GroupIDs: []*string{
	// 			to.Ptr("grafana")},
	// 			PrivateLinkResourceID: to.Ptr("/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-000000000000/resourceGroups/xx-rg/providers/Microsoft.Kusto/Clusters/sampleKustoResource"),
	// 			PrivateLinkResourceRegion: to.Ptr("West US"),
	// 			PrivateLinkServiceURL: to.Ptr("my-self-hosted-influxdb.westus.mydomain.com"),
	// 			ProvisioningState: to.Ptr(armdashboard.ProvisioningStateSucceeded),
	// 			RequestMessage: to.Ptr("Example Request Message"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/78eac0bd58633028293cb1ec1709baa200bed9e2/specification/dashboard/resource-manager/Microsoft.Dashboard/stable/2023-09-01/examples/ManagedPrivateEndpoints_Patch.json
func ExampleManagedPrivateEndpointsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdashboard.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewManagedPrivateEndpointsClient().BeginUpdate(ctx, "myResourceGroup", "myWorkspace", "myMPEName", armdashboard.ManagedPrivateEndpointUpdateParameters{
		Tags: map[string]*string{
			"Environment": to.Ptr("Dev 2"),
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
	// res.ManagedPrivateEndpointModel = armdashboard.ManagedPrivateEndpointModel{
	// 	Name: to.Ptr("myMPEName"),
	// 	Type: to.Ptr("Microsoft.Dashboard/grafana/managedPrivateEndpoint"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/Microsoft.Dashboard/grafana/myWorkspace/managedPrivateEndpoints/myPrivateEndpointName"),
	// 	SystemData: &armdashboard.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armdashboard.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.197Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armdashboard.CreatedByTypeUser),
	// 	},
	// 	Location: to.Ptr("West US"),
	// 	Properties: &armdashboard.ManagedPrivateEndpointModelProperties{
	// 		GroupIDs: []*string{
	// 			to.Ptr("grafana")},
	// 			PrivateLinkResourceID: to.Ptr("/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-000000000000/resourceGroups/xx-rg/providers/Microsoft.Kusto/Clusters/sampleKustoResource"),
	// 			PrivateLinkResourceRegion: to.Ptr("West US"),
	// 			PrivateLinkServicePrivateIP: to.Ptr("10.0.0.5"),
	// 			PrivateLinkServiceURL: to.Ptr("my-self-hosted-influxdb.westus.mydomain.com"),
	// 			ProvisioningState: to.Ptr(armdashboard.ProvisioningStateSucceeded),
	// 			RequestMessage: to.Ptr("Example Request Message"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/78eac0bd58633028293cb1ec1709baa200bed9e2/specification/dashboard/resource-manager/Microsoft.Dashboard/stable/2023-09-01/examples/ManagedPrivateEndpoints_Delete.json
func ExampleManagedPrivateEndpointsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdashboard.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewManagedPrivateEndpointsClient().BeginDelete(ctx, "myResourceGroup", "myWorkspace", "myMPEName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
