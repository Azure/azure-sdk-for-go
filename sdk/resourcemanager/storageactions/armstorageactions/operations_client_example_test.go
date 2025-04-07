//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstorageactions_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storageactions/armstorageactions"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/storageactions/resource-manager/Microsoft.StorageActions/stable/2023-01-01/examples/misc/OperationsList.json
func ExampleOperationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstorageactions.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewOperationsClient().NewListPager(nil)
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
		// page.OperationListResult = armstorageactions.OperationListResult{
		// 	Value: []*armstorageactions.Operation{
		// 		{
		// 			Name: to.Ptr("Microsoft.StorageActions/storageTasks/read"),
		// 			Display: &armstorageactions.OperationDisplay{
		// 				Description: to.Ptr("Gets or Lists existing StorageTask resource(s)."),
		// 				Operation: to.Ptr("Get or List StorageTask resource(s)."),
		// 				Provider: to.Ptr("Microsoft StorageActions"),
		// 				Resource: to.Ptr("StorageTasks"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.StorageActions/storageTasks/write"),
		// 			Display: &armstorageactions.OperationDisplay{
		// 				Description: to.Ptr("Creates or Updates StorageTask resource."),
		// 				Operation: to.Ptr("Create or Update StorageTask resource."),
		// 				Provider: to.Ptr("Microsoft StorageActions"),
		// 				Resource: to.Ptr("StorageTasks"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 		},
		// 		{
		// 			Name: to.Ptr("Microsoft.StorageActions/storageTasks/delete"),
		// 			Display: &armstorageactions.OperationDisplay{
		// 				Description: to.Ptr("Deletes StorageTask resource."),
		// 				Operation: to.Ptr("Delete StorageTask resource."),
		// 				Provider: to.Ptr("Microsoft StorageActions"),
		// 				Resource: to.Ptr("StorageTasks"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 	}},
		// }
	}
}
