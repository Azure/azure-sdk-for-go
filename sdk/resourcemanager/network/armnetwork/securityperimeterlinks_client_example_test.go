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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspLinkGet.json
func ExampleSecurityPerimeterLinksClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecurityPerimeterLinksClient().Get(ctx, "rg1", "nsp1", "link1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NspLink = armnetwork.NspLink{
	// 	Name: to.Ptr("link1"),
	// 	Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/links"),
	// 	ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/links/link1"),
	// 	SystemData: &armnetwork.SecurityPerimeterSystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		CreatedBy: to.Ptr("user"),
	// 		CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user"),
	// 		LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 	},
	// 	Properties: &armnetwork.NspLinkProperties{
	// 		Description: to.Ptr("Auto Approved"),
	// 		AutoApprovedRemotePerimeterResourceID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp2"),
	// 		LocalInboundProfiles: []*string{
	// 			to.Ptr("*")},
	// 			LocalOutboundProfiles: []*string{
	// 				to.Ptr("*")},
	// 				ProvisioningState: to.Ptr(armnetwork.NspLinkProvisioningStateSucceeded),
	// 				RemoteInboundProfiles: []*string{
	// 					to.Ptr("*")},
	// 					RemoteOutboundProfiles: []*string{
	// 						to.Ptr("*")},
	// 						RemotePerimeterGUID: to.Ptr("guid"),
	// 						RemotePerimeterLocation: to.Ptr("westus2"),
	// 						Status: to.Ptr(armnetwork.NspLinkStatusApproved),
	// 					},
	// 				}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspLinkPut.json
func ExampleSecurityPerimeterLinksClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecurityPerimeterLinksClient().CreateOrUpdate(ctx, "rg1", "nsp1", "link1", armnetwork.NspLink{
		Properties: &armnetwork.NspLinkProperties{
			AutoApprovedRemotePerimeterResourceID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp2"),
			LocalInboundProfiles: []*string{
				to.Ptr("*")},
			RemoteInboundProfiles: []*string{
				to.Ptr("*")},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NspLink = armnetwork.NspLink{
	// 	Name: to.Ptr("link1"),
	// 	Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/links"),
	// 	ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/links/link1"),
	// 	SystemData: &armnetwork.SecurityPerimeterSystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		CreatedBy: to.Ptr("user"),
	// 		CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user"),
	// 		LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 	},
	// 	Properties: &armnetwork.NspLinkProperties{
	// 		Description: to.Ptr("Auto Approved"),
	// 		AutoApprovedRemotePerimeterResourceID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp2"),
	// 		LocalInboundProfiles: []*string{
	// 			to.Ptr("*")},
	// 			LocalOutboundProfiles: []*string{
	// 				to.Ptr("*")},
	// 				ProvisioningState: to.Ptr(armnetwork.NspLinkProvisioningStateSucceeded),
	// 				RemoteInboundProfiles: []*string{
	// 					to.Ptr("*")},
	// 					RemoteOutboundProfiles: []*string{
	// 						to.Ptr("*")},
	// 						RemotePerimeterGUID: to.Ptr("guid"),
	// 						RemotePerimeterLocation: to.Ptr("westus2"),
	// 						Status: to.Ptr(armnetwork.NspLinkStatusApproved),
	// 					},
	// 				}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspLinkDelete.json
func ExampleSecurityPerimeterLinksClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSecurityPerimeterLinksClient().BeginDelete(ctx, "rg1", "nsp1", "link1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NspLinkList.json
func ExampleSecurityPerimeterLinksClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSecurityPerimeterLinksClient().NewListPager("rg1", "nsp1", &armnetwork.SecurityPerimeterLinksClientListOptions{Top: nil,
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
		// page.NspLinkListResult = armnetwork.NspLinkListResult{
		// 	Value: []*armnetwork.NspLink{
		// 		{
		// 			Name: to.Ptr("link1"),
		// 			Type: to.Ptr("Microsoft.Network/networkSecurityPerimeters/links"),
		// 			ID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp1/links/link1"),
		// 			SystemData: &armnetwork.SecurityPerimeterSystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				CreatedBy: to.Ptr("user"),
		// 				CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-02-07T18:07:36.344Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user"),
		// 				LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 			},
		// 			Properties: &armnetwork.NspLinkProperties{
		// 				Description: to.Ptr("Auto Approved"),
		// 				AutoApprovedRemotePerimeterResourceID: to.Ptr("/subscriptions/subId/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityPerimeters/nsp2"),
		// 				LocalInboundProfiles: []*string{
		// 					to.Ptr("*")},
		// 					LocalOutboundProfiles: []*string{
		// 						to.Ptr("*")},
		// 						ProvisioningState: to.Ptr(armnetwork.NspLinkProvisioningStateSucceeded),
		// 						RemoteInboundProfiles: []*string{
		// 							to.Ptr("*")},
		// 							RemoteOutboundProfiles: []*string{
		// 								to.Ptr("*")},
		// 								RemotePerimeterGUID: to.Ptr("guid"),
		// 								RemotePerimeterLocation: to.Ptr("westus2"),
		// 								Status: to.Ptr(armnetwork.NspLinkStatusApproved),
		// 							},
		// 					}},
		// 				}
	}
}
