//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstorage_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/86c6306649b02e542117adb46c61e8019dbd78e9/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/DeletedAccountList.json
func ExampleDeletedAccountsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDeletedAccountsClient().NewListPager(nil)
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
		// page.DeletedAccountListResult = armstorage.DeletedAccountListResult{
		// 	Value: []*armstorage.DeletedAccount{
		// 		{
		// 			Name: to.Ptr("sto1125"),
		// 			Type: to.Ptr("Microsoft.Storage/deletedAccounts"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Storage/locations/eastus/deletedAccounts/sto1125"),
		// 			Properties: &armstorage.DeletedAccountProperties{
		// 				CreationTime: to.Ptr("2020-08-17T03:35:37.4588848Z"),
		// 				DeletionTime: to.Ptr("2020-08-17T04:41:37.3442475Z"),
		// 				Location: to.Ptr("eastus"),
		// 				RestoreReference: to.Ptr("sto1125|2020-08-17T03:35:37.4588848Z"),
		// 				StorageAccountResourceID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/sto/providers/Microsoft.Storage/storageAccounts/sto1125"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("sto1126"),
		// 			Type: to.Ptr("Microsoft.Storage/deletedAccounts"),
		// 			ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Storage/locations/eastus/deletedAccounts/sto1126"),
		// 			Properties: &armstorage.DeletedAccountProperties{
		// 				CreationTime: to.Ptr("2020-08-19T15:10:21.5902165Z"),
		// 				DeletionTime: to.Ptr("2020-08-20T06:11:55.1957302Z"),
		// 				Location: to.Ptr("eastus"),
		// 				RestoreReference: to.Ptr("sto1126|2020-08-17T03:35:37.4588848Z"),
		// 				StorageAccountResourceID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/sto/providers/Microsoft.Storage/storageAccounts/sto1126"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/86c6306649b02e542117adb46c61e8019dbd78e9/specification/storage/resource-manager/Microsoft.Storage/stable/2024-01-01/examples/DeletedAccountGet.json
func ExampleDeletedAccountsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDeletedAccountsClient().Get(ctx, "sto1125", "eastus", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DeletedAccount = armstorage.DeletedAccount{
	// 	Name: to.Ptr("sto1125"),
	// 	Type: to.Ptr("Microsoft.Storage/deletedAccounts"),
	// 	ID: to.Ptr("/subscriptions/{subscription-id}/providers/Microsoft.Storage/locations/eastus/deletedAccounts/sto1125"),
	// 	Properties: &armstorage.DeletedAccountProperties{
	// 		CreationTime: to.Ptr("2020-08-17T03:35:37.4588848Z"),
	// 		DeletionTime: to.Ptr("2020-08-17T04:41:37.3442475Z"),
	// 		Location: to.Ptr("eastus"),
	// 		RestoreReference: to.Ptr("sto1125|2020-08-17T03:35:37.4588848Z"),
	// 		StorageAccountResourceID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/sto/providers/Microsoft.Storage/storageAccounts/sto1125"),
	// 	},
	// }
}
