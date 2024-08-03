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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_CreateOrUpdate_MaximumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_BeginCreateOrUpdate_virtualMachineScaleSetExtensionCreateOrUpdateMaximumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().BeginCreateOrUpdate(ctx, "rgcompute", "aaaaaaa", "aaaaaaaaaaaaaaaaaaaaa", armcompute.VirtualMachineScaleSetExtension{
		Name: to.Ptr("{extension-name}"),
		Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
			Type:                    to.Ptr("{extension-Type}"),
			AutoUpgradeMinorVersion: to.Ptr(true),
			EnableAutomaticUpgrade:  to.Ptr(true),
			ForceUpdateTag:          to.Ptr("aaaaaaaaa"),
			ProtectedSettings:       map[string]any{},
			ProvisionAfterExtensions: []*string{
				to.Ptr("aa")},
			Publisher:          to.Ptr("{extension-Publisher}"),
			Settings:           map[string]any{},
			SuppressFailures:   to.Ptr(true),
			TypeHandlerVersion: to.Ptr("{handler-version}"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualMachineScaleSetExtension = armcompute.VirtualMachineScaleSetExtension{
	// 	ID: to.Ptr("aaaaaaaa"),
	// 	Name: to.Ptr("{extension-name}"),
	// 	Type: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaa"),
	// 	Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
	// 		Type: to.Ptr("{extension-Type}"),
	// 		AutoUpgradeMinorVersion: to.Ptr(true),
	// 		EnableAutomaticUpgrade: to.Ptr(true),
	// 		ForceUpdateTag: to.Ptr("aaaaaaaaa"),
	// 		ProtectedSettings: map[string]any{
	// 		},
	// 		ProvisionAfterExtensions: []*string{
	// 			to.Ptr("aa")},
	// 			ProvisioningState: to.Ptr("Succeeded"),
	// 			Publisher: to.Ptr("{extension-Publisher}"),
	// 			Settings: map[string]any{
	// 			},
	// 			SuppressFailures: to.Ptr(true),
	// 			TypeHandlerVersion: to.Ptr("{handler-version}"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_CreateOrUpdate_MinimumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_BeginCreateOrUpdate_virtualMachineScaleSetExtensionCreateOrUpdateMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().BeginCreateOrUpdate(ctx, "rgcompute", "aaaaaaaaaaa", "aaaaaaaaaaa", armcompute.VirtualMachineScaleSetExtension{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualMachineScaleSetExtension = armcompute.VirtualMachineScaleSetExtension{
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_Update_MaximumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_BeginUpdate_virtualMachineScaleSetExtensionUpdateMaximumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().BeginUpdate(ctx, "rgcompute", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "aaaa", armcompute.VirtualMachineScaleSetExtensionUpdate{
		Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
			Type:                    to.Ptr("{extension-Type}"),
			AutoUpgradeMinorVersion: to.Ptr(true),
			EnableAutomaticUpgrade:  to.Ptr(true),
			ForceUpdateTag:          to.Ptr("aaaaaaaaa"),
			ProtectedSettings:       map[string]any{},
			ProvisionAfterExtensions: []*string{
				to.Ptr("aa")},
			Publisher:          to.Ptr("{extension-Publisher}"),
			Settings:           map[string]any{},
			SuppressFailures:   to.Ptr(true),
			TypeHandlerVersion: to.Ptr("{handler-version}"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualMachineScaleSetExtension = armcompute.VirtualMachineScaleSetExtension{
	// 	ID: to.Ptr("aaaaaaaa"),
	// 	Name: to.Ptr("{extension-name}"),
	// 	Type: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaa"),
	// 	Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
	// 		Type: to.Ptr("{extension-Type}"),
	// 		AutoUpgradeMinorVersion: to.Ptr(true),
	// 		EnableAutomaticUpgrade: to.Ptr(true),
	// 		ForceUpdateTag: to.Ptr("aaaaaaaaa"),
	// 		ProtectedSettings: map[string]any{
	// 		},
	// 		ProvisionAfterExtensions: []*string{
	// 			to.Ptr("aa")},
	// 			ProvisioningState: to.Ptr("Succeeded"),
	// 			Publisher: to.Ptr("{extension-Publisher}"),
	// 			Settings: map[string]any{
	// 			},
	// 			SuppressFailures: to.Ptr(true),
	// 			TypeHandlerVersion: to.Ptr("{handler-version}"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_Update_MinimumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_BeginUpdate_virtualMachineScaleSetExtensionUpdateMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().BeginUpdate(ctx, "rgcompute", "aaaaaaaaaaaaaaaaaaaaaaaaaa", "aa", armcompute.VirtualMachineScaleSetExtensionUpdate{}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualMachineScaleSetExtension = armcompute.VirtualMachineScaleSetExtension{
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_Delete_MaximumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_BeginDelete_virtualMachineScaleSetExtensionDeleteMaximumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().BeginDelete(ctx, "rgcompute", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaaaaaaaaaa", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_Delete_MinimumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_BeginDelete_virtualMachineScaleSetExtensionDeleteMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().BeginDelete(ctx, "rgcompute", "aaaa", "aaaaaaaaaaaaaaaaaaaaaaa", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_Get_MaximumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_Get_virtualMachineScaleSetExtensionGetMaximumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().Get(ctx, "rgcompute", "aaaaaaaaaaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaaaaaa", &armcompute.VirtualMachineScaleSetExtensionsClientGetOptions{Expand: to.Ptr("aaaaaaa")})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualMachineScaleSetExtension = armcompute.VirtualMachineScaleSetExtension{
	// 	ID: to.Ptr("aaaaaaaa"),
	// 	Name: to.Ptr("{extension-name}"),
	// 	Type: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaa"),
	// 	Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
	// 		Type: to.Ptr("{extension-Type}"),
	// 		AutoUpgradeMinorVersion: to.Ptr(true),
	// 		EnableAutomaticUpgrade: to.Ptr(true),
	// 		ForceUpdateTag: to.Ptr("aaaaaaaaa"),
	// 		ProtectedSettings: map[string]any{
	// 		},
	// 		ProvisionAfterExtensions: []*string{
	// 			to.Ptr("aa")},
	// 			ProvisioningState: to.Ptr("Succeeded"),
	// 			Publisher: to.Ptr("{extension-Publisher}"),
	// 			Settings: map[string]any{
	// 			},
	// 			SuppressFailures: to.Ptr(true),
	// 			TypeHandlerVersion: to.Ptr("{handler-version}"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_Get_MinimumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_Get_virtualMachineScaleSetExtensionGetMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewVirtualMachineScaleSetExtensionsClient().Get(ctx, "rgcompute", "a", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", &armcompute.VirtualMachineScaleSetExtensionsClientGetOptions{Expand: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.VirtualMachineScaleSetExtension = armcompute.VirtualMachineScaleSetExtension{
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_List_MaximumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_NewListPager_virtualMachineScaleSetExtensionListMaximumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualMachineScaleSetExtensionsClient().NewListPager("rgcompute", "aaaaaaaaaaaaaaaaaaaa", nil)
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
		// page.VirtualMachineScaleSetExtensionListResult = armcompute.VirtualMachineScaleSetExtensionListResult{
		// 	Value: []*armcompute.VirtualMachineScaleSetExtension{
		// 		{
		// 			ID: to.Ptr("aaaaaaaa"),
		// 			Name: to.Ptr("{extension-name}"),
		// 			Type: to.Ptr("aaaaaaaaaaaaaaaaaaaaaaaa"),
		// 			Properties: &armcompute.VirtualMachineScaleSetExtensionProperties{
		// 				Type: to.Ptr("{extension-Type}"),
		// 				AutoUpgradeMinorVersion: to.Ptr(true),
		// 				EnableAutomaticUpgrade: to.Ptr(true),
		// 				ForceUpdateTag: to.Ptr("aaaaaaaaa"),
		// 				ProtectedSettings: map[string]any{
		// 				},
		// 				ProvisionAfterExtensions: []*string{
		// 					to.Ptr("aa")},
		// 					ProvisioningState: to.Ptr("Succeeded"),
		// 					Publisher: to.Ptr("{extension-Publisher}"),
		// 					Settings: map[string]any{
		// 					},
		// 					SuppressFailures: to.Ptr(true),
		// 					TypeHandlerVersion: to.Ptr("{handler-version}"),
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/compute/resource-manager/Microsoft.Compute/ComputeRP/stable/2024-07-01/examples/virtualMachineScaleSetExamples/VirtualMachineScaleSetExtension_List_MinimumSet_Gen.json
func ExampleVirtualMachineScaleSetExtensionsClient_NewListPager_virtualMachineScaleSetExtensionListMinimumSetGen() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcompute.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewVirtualMachineScaleSetExtensionsClient().NewListPager("rgcompute", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil)
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
		// page.VirtualMachineScaleSetExtensionListResult = armcompute.VirtualMachineScaleSetExtensionListResult{
		// 	Value: []*armcompute.VirtualMachineScaleSetExtension{
		// 		{
		// 	}},
		// }
	}
}
