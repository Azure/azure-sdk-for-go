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

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListSKUs-Gateways.json
func ExampleGatewaySKUsClient_NewListAvailableSKUsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGatewaySKUsClient().NewListAvailableSKUsPager("rg1", "apimService1", nil)
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
		// page.GatewayResourceSKUResults = armapimanagement.GatewayResourceSKUResults{
		// 	Value: []*armapimanagement.GatewayResourceSKUResult{
		// 		{
		// 			Capacity: &armapimanagement.GatewaySKUCapacity{
		// 				Default: to.Ptr[int32](1),
		// 				Maximum: to.Ptr[int32](4),
		// 				Minimum: to.Ptr[int32](1),
		// 				ScaleType: to.Ptr(armapimanagement.GatewaySKUCapacityScaleTypeManual),
		// 			},
		// 			ResourceType: to.Ptr("Microsoft.ApiManagement/gateways"),
		// 			SKU: &armapimanagement.GatewaySKU{
		// 				Name: to.Ptr(armapimanagement.APIGatewaySKUTypeWorkspaceGatewayStandard),
		// 			},
		// 		},
		// 		{
		// 			Capacity: &armapimanagement.GatewaySKUCapacity{
		// 				Default: to.Ptr[int32](1),
		// 				Maximum: to.Ptr[int32](12),
		// 				Minimum: to.Ptr[int32](1),
		// 				ScaleType: to.Ptr(armapimanagement.GatewaySKUCapacityScaleTypeManual),
		// 			},
		// 			ResourceType: to.Ptr("Microsoft.ApiManagement/gateways"),
		// 			SKU: &armapimanagement.GatewaySKU{
		// 				Name: to.Ptr(armapimanagement.APIGatewaySKUTypeWorkspaceGatewayPremium),
		// 			},
		// 	}},
		// }
	}
}
