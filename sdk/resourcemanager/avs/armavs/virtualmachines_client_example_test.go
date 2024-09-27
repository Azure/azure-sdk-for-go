//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armavs_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/avs/armavs/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e7678455846b1000fd31db27596e4ca3d299a872/specification/vmware/resource-manager/Microsoft.AVS/stable/2023-09-01/examples/VirtualMachines_List.json
func ExampleVirtualMachinesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualMachinesClient().NewListPager("group1", "cloud1", "cluster1", nil)
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
		// page.VirtualMachinesList = armavs.VirtualMachinesList{
		// 	Value: []*armavs.VirtualMachine{
		// 		{
		// 			Name: to.Ptr("vm-209"),
		// 			Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/virtualMachines"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-209"),
		// 			Properties: &armavs.VirtualMachineProperties{
		// 				DisplayName: to.Ptr("contoso-vm1"),
		// 				FolderPath: to.Ptr("vm/folder-1"),
		// 				MoRefID: to.Ptr("vm-209"),
		// 				RestrictMovement: to.Ptr(armavs.VirtualMachineRestrictMovementStateDisabled),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("vm-128"),
		// 			Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/virtualMachines"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-128"),
		// 			Properties: &armavs.VirtualMachineProperties{
		// 				DisplayName: to.Ptr("contoso-vm2"),
		// 				FolderPath: to.Ptr("vm"),
		// 				MoRefID: to.Ptr("vm-128"),
		// 				RestrictMovement: to.Ptr(armavs.VirtualMachineRestrictMovementStateEnabled),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e7678455846b1000fd31db27596e4ca3d299a872/specification/vmware/resource-manager/Microsoft.AVS/stable/2023-09-01/examples/VirtualMachines_Get.json
func ExampleVirtualMachinesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachinesClient().Get(ctx, "group1", "cloud1", "cluster1", "vm-209", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualMachine = armavs.VirtualMachine{
	// 	Name: to.Ptr("vm-209"),
	// 	Type: to.Ptr("Microsoft.AVS/privateClouds/clusters/virtualMachines"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/cloud1/clusters/cluster1/virtualMachines/vm-209"),
	// 	Properties: &armavs.VirtualMachineProperties{
	// 		DisplayName: to.Ptr("contoso-vm"),
	// 		FolderPath: to.Ptr("vm/folder-1"),
	// 		MoRefID: to.Ptr("vm-209"),
	// 		RestrictMovement: to.Ptr(armavs.VirtualMachineRestrictMovementStateDisabled),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e7678455846b1000fd31db27596e4ca3d299a872/specification/vmware/resource-manager/Microsoft.AVS/stable/2023-09-01/examples/VirtualMachines_RestrictMovement.json
func ExampleVirtualMachinesClient_BeginRestrictMovement() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armavs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachinesClient().BeginRestrictMovement(ctx, "group1", "cloud1", "cluster1", "vm-209", armavs.VirtualMachineRestrictMovement{
		RestrictMovement: to.Ptr(armavs.VirtualMachineRestrictMovementStateEnabled),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
