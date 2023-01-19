//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armrecoveryservicesbackup_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3751704f5318f1175875c94b66af769db917f2d3/specification/recoveryservicesbackup/resource-manager/Microsoft.RecoveryServices/stable/2023-01-01/examples/ResourceGuardProxyCRUD/ListResourceGuardProxy.json
func ExampleResourceGuardProxiesClient_NewGetPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armrecoveryservicesbackup.NewResourceGuardProxiesClient("0b352192-dcac-4cc7-992e-a96190ccc68c", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewGetPager("sampleVault", "SampleResourceGroup", nil)
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
		// page.ResourceGuardProxyBaseResourceList = armrecoveryservicesbackup.ResourceGuardProxyBaseResourceList{
		// 	Value: []*armrecoveryservicesbackup.ResourceGuardProxyBaseResource{
		// 		{
		// 			Name: to.Ptr("swaggerExample"),
		// 			Type: to.Ptr("Microsoft.RecoveryServices/vaults/backupResourceGuardProxies"),
		// 			ID: to.Ptr("/backupmanagement/resources/sampleVault/backupResourceGuardProxies/swaggerExample"),
		// 			Properties: &armrecoveryservicesbackup.ResourceGuardProxyBase{
		// 				Description: to.Ptr("Please take JIT access before performing any of the critical operation"),
		// 				LastUpdatedTime: to.Ptr("2021-02-11T12:20:47.8210031Z"),
		// 				ResourceGuardOperationDetails: []*armrecoveryservicesbackup.ResourceGuardOperationDetail{
		// 					{
		// 						DefaultResourceRequest: to.Ptr("/subscriptions/c999d45b-944f-418c-a0d8-c3fcfd1802c8/resourceGroups/vaultguardRGNew/providers/Microsoft.DataProtection/resourceGuards/VaultGuardTestNew/deleteResourceGuardProxyRequests/default"),
		// 						VaultCriticalOperation: to.Ptr("Microsoft.DataProtection/resourceGuards/deleteResourceGuardProxyRequests"),
		// 					},
		// 					{
		// 						DefaultResourceRequest: to.Ptr("/subscriptions/c999d45b-944f-418c-a0d8-c3fcfd1802c8/resourceGroups/vaultguardRGNew/providers/Microsoft.DataProtection/resourceGuards/VaultGuardTestNew/disableSoftDeleteRequests/default"),
		// 						VaultCriticalOperation: to.Ptr("Microsoft.DataProtection/resourceGuards/disableSoftDeleteRequests"),
		// 				}},
		// 				ResourceGuardResourceID: to.Ptr("/subscriptions/c999d45b-944f-418c-a0d8-c3fcfd1802c8/resourceGroups/vaultguardRGNew/providers/Microsoft.DataProtection/resourceGuards/VaultGuardTestNew"),
		// 			},
		// 	}},
		// }
	}
}
