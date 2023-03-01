//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdatafactory_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/ManagedVirtualNetworks_ListByFactory.json
func ExampleManagedVirtualNetworksClient_NewListByFactoryPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewManagedVirtualNetworksClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByFactoryPager("exampleResourceGroup", "exampleFactoryName", nil)
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
		// page.ManagedVirtualNetworkListResponse = armdatafactory.ManagedVirtualNetworkListResponse{
		// 	Value: []*armdatafactory.ManagedVirtualNetworkResource{
		// 		{
		// 			Name: to.Ptr("exampleManagedVirtualNetworkName"),
		// 			Type: to.Ptr("Microsoft.DataFactory/factories/managedVirtualNetworks"),
		// 			Etag: to.Ptr("0400f1a1-0000-0000-0000-5b2188640000"),
		// 			ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/managedVirtualNetworks/exampleManagedVirtualNetworkName"),
		// 			Properties: &armdatafactory.ManagedVirtualNetwork{
		// 				Alias: to.Ptr("exampleFactoryName"),
		// 				VNetID: to.Ptr("5a7bd944-87e6-454a-8d4d-9fba446514fd"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/ManagedVirtualNetworks_Create.json
func ExampleManagedVirtualNetworksClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewManagedVirtualNetworksClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleManagedVirtualNetworkName", armdatafactory.ManagedVirtualNetworkResource{
		Properties: &armdatafactory.ManagedVirtualNetwork{},
	}, &armdatafactory.ManagedVirtualNetworksClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ManagedVirtualNetworkResource = armdatafactory.ManagedVirtualNetworkResource{
	// 	Name: to.Ptr("exampleManagedVirtualNetworkName"),
	// 	Type: to.Ptr("Microsoft.DataFactory/factories/managedVirtualNetworks"),
	// 	Etag: to.Ptr("000046c4-0000-0000-0000-5b2198bf0000"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/managedVirtualNetworks/exampleManagedVirtualNetworkName"),
	// 	Properties: &armdatafactory.ManagedVirtualNetwork{
	// 		Alias: to.Ptr("exampleFactoryName"),
	// 		VNetID: to.Ptr("12345678-1234-1234-1234-12345678123"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/5c9459305484e0456b4a922e3d31a61e2ddd3c99/specification/datafactory/resource-manager/Microsoft.DataFactory/stable/2018-06-01/examples/ManagedVirtualNetworks_Get.json
func ExampleManagedVirtualNetworksClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatafactory.NewManagedVirtualNetworksClient("12345678-1234-1234-1234-12345678abc", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx, "exampleResourceGroup", "exampleFactoryName", "exampleManagedVirtualNetworkName", &armdatafactory.ManagedVirtualNetworksClientGetOptions{IfNoneMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ManagedVirtualNetworkResource = armdatafactory.ManagedVirtualNetworkResource{
	// 	Name: to.Ptr("exampleManagedVirtualNetworkName"),
	// 	Type: to.Ptr("Microsoft.DataFactory/factories/managedVirtualNetworks"),
	// 	Etag: to.Ptr("15003c4f-0000-0200-0000-5cbe090b0000"),
	// 	ID: to.Ptr("/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/managedVirtualNetworks/exampleManagedVirtualNetworkName"),
	// 	Properties: &armdatafactory.ManagedVirtualNetwork{
	// 		Alias: to.Ptr("exampleFactoryName"),
	// 		VNetID: to.Ptr("5a7bd944-87e6-454a-8d4d-9fba446514fd"),
	// 	},
	// }
}
