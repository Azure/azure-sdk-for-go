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

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4517f89a8ebd2f6a94e107e5ee60fff9886f3612/specification/compute/resource-manager/Microsoft.Compute/CloudserviceRP/stable/2024-11-04/examples/CloudServiceRole_Get.json
func ExampleCloudServiceRolesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCloudServiceRolesClient().Get(ctx, "{role-name}", "ConstosoRG", "{cs-name}", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CloudServiceRole = armcompute.CloudServiceRole{
	// 	Name: to.Ptr("{role-name}"),
	// 	Type: to.Ptr("Microsoft.Compute/cloudServices/roles"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/ConstosoRG/providers/Microsoft.Compute/cloudServices/{cs-name}/roles/{role-name}"),
	// 	Location: to.Ptr("eastus2euap"),
	// 	Properties: &armcompute.CloudServiceRoleProperties{
	// 		UniqueID: to.Ptr("b03bc269-766b-4921-b91a-7dffaae4d03b:{role-name}"),
	// 	},
	// 	SKU: &armcompute.CloudServiceRoleSKU{
	// 		Name: to.Ptr("Standard_D1_v2"),
	// 		Capacity: to.Ptr[int64](2),
	// 		Tier: to.Ptr("Standard"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4517f89a8ebd2f6a94e107e5ee60fff9886f3612/specification/compute/resource-manager/Microsoft.Compute/CloudserviceRP/stable/2024-11-04/examples/CloudServiceRole_List.json
func ExampleCloudServiceRolesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewCloudServiceRolesClient().NewListPager("ConstosoRG", "{cs-name}", nil)
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
		// page.CloudServiceRoleListResult = armcompute.CloudServiceRoleListResult{
		// 	Value: []*armcompute.CloudServiceRole{
		// 		{
		// 			Name: to.Ptr("ContosoFrontend"),
		// 			Type: to.Ptr("Microsoft.Compute/cloudServices/roles"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/ConstosoRG/providers/Microsoft.Compute/cloudServices/{cs-name}/roles/ContosoFrontend"),
		// 			Location: to.Ptr("eastus2euap"),
		// 			Properties: &armcompute.CloudServiceRoleProperties{
		// 				UniqueID: to.Ptr("b03bc269-766b-4921-b91a-7dffaae4d03b:ContosoFrontend"),
		// 			},
		// 			SKU: &armcompute.CloudServiceRoleSKU{
		// 				Name: to.Ptr("Standard_D1_v2"),
		// 				Capacity: to.Ptr[int64](2),
		// 				Tier: to.Ptr("Standard"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("ContosoBackend"),
		// 			Type: to.Ptr("Microsoft.Compute/cloudServices/roles"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/ConstosoRG/providers/Microsoft.Compute/cloudServices/{cs-name}/roles/ContosoBackend"),
		// 			Location: to.Ptr("eastus2euap"),
		// 			Properties: &armcompute.CloudServiceRoleProperties{
		// 				UniqueID: to.Ptr("b03bc269-766b-4921-b91a-7dffaae4d03b:ContosoBackend"),
		// 			},
		// 			SKU: &armcompute.CloudServiceRoleSKU{
		// 				Name: to.Ptr("Standard_D1_v2"),
		// 				Capacity: to.Ptr[int64](2),
		// 				Tier: to.Ptr("Standard"),
		// 			},
		// 	}},
		// }
	}
}
