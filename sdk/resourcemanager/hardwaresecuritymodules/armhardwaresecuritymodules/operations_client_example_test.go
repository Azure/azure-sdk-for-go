//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhardwaresecuritymodules_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hardwaresecuritymodules/armhardwaresecuritymodules/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/18b5c820705ab69735b7e1e2e0da5e37ca6e1969/specification/hardwaresecuritymodules/resource-manager/Microsoft.HardwareSecurityModules/stable/2021-11-30/examples/DedicatedHsm_OperationsList.json
func ExampleOperationsClient_NewListPager_getAListOfDedicatedHsmOperations() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhardwaresecuritymodules.NewClientFactory("<subscription-id>", cred, nil)
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
		// page.DedicatedHsmOperationListResult = armhardwaresecuritymodules.DedicatedHsmOperationListResult{
		// 	Value: []*armhardwaresecuritymodules.DedicatedHsmOperation{
		// 		{
		// 			Name: to.Ptr("hsm1"),
		// 			Display: &armhardwaresecuritymodules.DedicatedHsmOperationDisplay{
		// 				Description: to.Ptr("Update a dedicated HSM in the specified subscription"),
		// 				Operation: to.Ptr("DedicatedHsm_Update"),
		// 				Provider: to.Ptr("Microsoft HardwareSecurityModules"),
		// 				Resource: to.Ptr("Dedicated HSM"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("system"),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/18b5c820705ab69735b7e1e2e0da5e37ca6e1969/specification/hardwaresecuritymodules/resource-manager/Microsoft.HardwareSecurityModules/stable/2021-11-30/examples/PaymentHsm_OperationsList.json
func ExampleOperationsClient_NewListPager_getAListOfPaymentHsmOperations() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armhardwaresecuritymodules.NewClientFactory("<subscription-id>", cred, nil)
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
		// page.DedicatedHsmOperationListResult = armhardwaresecuritymodules.DedicatedHsmOperationListResult{
		// 	Value: []*armhardwaresecuritymodules.DedicatedHsmOperation{
		// 		{
		// 			Name: to.Ptr("hsm1"),
		// 			Display: &armhardwaresecuritymodules.DedicatedHsmOperationDisplay{
		// 				Description: to.Ptr("Update a dedicated HSM in the specified subscription"),
		// 				Operation: to.Ptr("DedicatedHsm_Update"),
		// 				Provider: to.Ptr("Microsoft HardwareSecurityModules"),
		// 				Resource: to.Ptr("Dedicated HSM"),
		// 			},
		// 			IsDataAction: to.Ptr(false),
		// 			Origin: to.Ptr("system"),
		// 	}},
		// }
	}
}
