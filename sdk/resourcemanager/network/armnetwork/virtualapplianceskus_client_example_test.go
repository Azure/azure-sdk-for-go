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

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/NetworkVirtualApplianceSkuList.json
func ExampleVirtualApplianceSKUsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualApplianceSKUsClient().NewListPager(nil)
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
		// page.VirtualApplianceSKUListResult = armnetwork.VirtualApplianceSKUListResult{
		// 	Value: []*armnetwork.VirtualApplianceSKU{
		// 		{
		// 			Name: to.Ptr("ciscoSdwan"),
		// 			Type: to.Ptr("Microsoft.Network/networkVirtualApplianceSkus"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkVirtualApplianceSkus/ciscoSdwan"),
		// 			Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
		// 			Properties: &armnetwork.VirtualApplianceSKUPropertiesFormat{
		// 				AvailableScaleUnits: []*armnetwork.VirtualApplianceSKUInstances{
		// 					{
		// 						InstanceCount: to.Ptr[int32](2),
		// 						ScaleUnit: to.Ptr("1"),
		// 					},
		// 					{
		// 						InstanceCount: to.Ptr[int32](2),
		// 						ScaleUnit: to.Ptr("2"),
		// 				}},
		// 				AvailableVersions: []*string{
		// 					to.Ptr("11.12")},
		// 					Vendor: to.Ptr("Cisco"),
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/639ecfad68419328658bd4cfe7094af4ce472be2/specification/network/resource-manager/Microsoft.Network/stable/2023-06-01/examples/NetworkVirtualApplianceSkuGet.json
func ExampleVirtualApplianceSKUsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualApplianceSKUsClient().Get(ctx, "ciscoSdwan", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualApplianceSKU = armnetwork.VirtualApplianceSKU{
	// 	Name: to.Ptr("ciscoSdwan"),
	// 	Type: to.Ptr("Microsoft.Network/networkVirtualApplianceSkus"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkVirtualApplianceSkus/ciscoSdwan"),
	// 	Etag: to.Ptr("w/\\00000000-0000-0000-0000-000000000000\\"),
	// 	Properties: &armnetwork.VirtualApplianceSKUPropertiesFormat{
	// 		AvailableScaleUnits: []*armnetwork.VirtualApplianceSKUInstances{
	// 			{
	// 				InstanceCount: to.Ptr[int32](2),
	// 				ScaleUnit: to.Ptr("1"),
	// 			},
	// 			{
	// 				InstanceCount: to.Ptr[int32](2),
	// 				ScaleUnit: to.Ptr("2"),
	// 		}},
	// 		AvailableVersions: []*string{
	// 			to.Ptr("11.12")},
	// 			Vendor: to.Ptr("Cisco"),
	// 		},
	// 	}
}
