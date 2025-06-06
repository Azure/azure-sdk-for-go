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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/PerimeterAssociableResourcesList.json
func ExampleSecurityPerimeterAssociableResourceTypesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSecurityPerimeterAssociableResourceTypesClient().NewListPager("westus", nil)
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
		// page.PerimeterAssociableResourcesListResult = armnetwork.PerimeterAssociableResourcesListResult{
		// 	Value: []*armnetwork.PerimeterAssociableResource{
		// 		{
		// 			Name: to.Ptr("Microsoft.Sql.servers"),
		// 			Type: to.Ptr("Microsoft.Network/PerimeterAssociableResourceTypes"),
		// 			ID: to.Ptr("/subscriptions/{subscriptionID}/providers/Microsoft.Network/perimeterAssociableResourceTypes/Microsoft.Sql.servers"),
		// 			Properties: &armnetwork.PerimeterAssociableResourceProperties{
		// 				DisplayName: to.Ptr("Microsoft.Sql/servers"),
		// 				PublicDNSZones: []*string{
		// 					to.Ptr("database.windows.net")},
		// 					ResourceType: to.Ptr("Microsoft.Sql/servers"),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("Microsoft.Storage.accounts"),
		// 				Type: to.Ptr("Microsoft.Network/PerimeterAssociableResourceTypes"),
		// 				ID: to.Ptr("/subscriptions/{subscriptionId}/providers/Microsoft.Network/perimeterAssociableResourceTypes/Microsoft.Storage.storageAccounts"),
		// 				Properties: &armnetwork.PerimeterAssociableResourceProperties{
		// 					DisplayName: to.Ptr("Microsoft.Storage/accounts"),
		// 					PublicDNSZones: []*string{
		// 						to.Ptr("blob.core.windows.net"),
		// 						to.Ptr("table.core.windows.net"),
		// 						to.Ptr("queue.core.windows.net"),
		// 						to.Ptr("file.core.windows.net")},
		// 						ResourceType: to.Ptr("Microsoft.Storage/accounts"),
		// 					},
		// 			}},
		// 		}
	}
}
