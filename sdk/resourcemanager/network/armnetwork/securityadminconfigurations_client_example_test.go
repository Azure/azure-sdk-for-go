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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkManagerSecurityAdminConfigurationList.json
func ExampleSecurityAdminConfigurationsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSecurityAdminConfigurationsClient().NewListPager("rg1", "testNetworkManager", &armnetwork.SecurityAdminConfigurationsClientListOptions{Top: nil,
		SkipToken: nil,
	})
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
		// page.SecurityAdminConfigurationListResult = armnetwork.SecurityAdminConfigurationListResult{
		// 	Value: []*armnetwork.SecurityAdminConfiguration{
		// 		{
		// 			Name: to.Ptr("myTestSecurityConfig"),
		// 			Type: to.Ptr("Microsoft.Network/networkManagers/securityAdminConfigurations"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/testNetworkManager/securityAdminConfigurations/myTestSecurityConfig"),
		// 			Properties: &armnetwork.SecurityAdminConfigurationPropertiesFormat{
		// 				Description: to.Ptr("A sample policy"),
		// 				ApplyOnNetworkIntentPolicyBasedServices: []*armnetwork.NetworkIntentPolicyBasedService{
		// 					to.Ptr(armnetwork.NetworkIntentPolicyBasedServiceNone)},
		// 					ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 					ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				},
		// 				SystemData: &armnetwork.SystemData{
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
		// 					CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
		// 					CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
		// 					LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkManagerSecurityAdminConfigurationGet.json
func ExampleSecurityAdminConfigurationsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecurityAdminConfigurationsClient().Get(ctx, "rg1", "testNetworkManager", "myTestSecurityConfig", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SecurityAdminConfiguration = armnetwork.SecurityAdminConfiguration{
	// 	Name: to.Ptr("myTestSecurityConfig"),
	// 	Type: to.Ptr("Microsoft.Network/networkManagers/securityAdminConfigurations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/testNetworkManager/securityAdminConfigurations/myTestSecurityConfig"),
	// 	Properties: &armnetwork.SecurityAdminConfigurationPropertiesFormat{
	// 		Description: to.Ptr("A sample policy"),
	// 		ApplyOnNetworkIntentPolicyBasedServices: []*armnetwork.NetworkIntentPolicyBasedService{
	// 			to.Ptr(armnetwork.NetworkIntentPolicyBasedServiceNone)},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		},
	// 		SystemData: &armnetwork.SystemData{
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 			CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 			CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 			LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkManagerSecurityAdminConfigurationPut_ManualAggregation.json
func ExampleSecurityAdminConfigurationsClient_CreateOrUpdate_createManualModeSecurityAdminConfiguration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecurityAdminConfigurationsClient().CreateOrUpdate(ctx, "rg1", "testNetworkManager", "myTestSecurityConfig", armnetwork.SecurityAdminConfiguration{
		Properties: &armnetwork.SecurityAdminConfigurationPropertiesFormat{
			Description: to.Ptr("A configuration which will update any network groups ip addresses at commit times."),
			NetworkGroupAddressSpaceAggregationOption: to.Ptr(armnetwork.AddressSpaceAggregationOptionManual),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SecurityAdminConfiguration = armnetwork.SecurityAdminConfiguration{
	// 	Name: to.Ptr("myTestSecurityConfig"),
	// 	Type: to.Ptr("Microsoft.Network/networkManagers/securityAdminConfigurations"),
	// 	ID: to.Ptr("/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/rg1/providers/Microsoft.Network/networkManager/testNetworkManager/securityAdminConfigurations/myTestSecurityConfig"),
	// 	Properties: &armnetwork.SecurityAdminConfigurationPropertiesFormat{
	// 		Description: to.Ptr("A sample policy"),
	// 		NetworkGroupAddressSpaceAggregationOption: to.Ptr(armnetwork.AddressSpaceAggregationOptionManual),
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 	},
	// 	SystemData: &armnetwork.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 		CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 		CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 		LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkManagerSecurityAdminConfigurationPut.json
func ExampleSecurityAdminConfigurationsClient_CreateOrUpdate_createNetworkManagerSecurityAdminConfiguration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSecurityAdminConfigurationsClient().CreateOrUpdate(ctx, "rg1", "testNetworkManager", "myTestSecurityConfig", armnetwork.SecurityAdminConfiguration{
		Properties: &armnetwork.SecurityAdminConfigurationPropertiesFormat{
			Description: to.Ptr("A sample policy"),
			ApplyOnNetworkIntentPolicyBasedServices: []*armnetwork.NetworkIntentPolicyBasedService{
				to.Ptr(armnetwork.NetworkIntentPolicyBasedServiceNone)},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SecurityAdminConfiguration = armnetwork.SecurityAdminConfiguration{
	// 	Name: to.Ptr("myTestSecurityConfig"),
	// 	Type: to.Ptr("Microsoft.Network/networkManagers/securityAdminConfigurations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManager/testNetworkManager/securityAdminConfigurations/myTestSecurityConfig"),
	// 	Properties: &armnetwork.SecurityAdminConfigurationPropertiesFormat{
	// 		Description: to.Ptr("A sample policy"),
	// 		ApplyOnNetworkIntentPolicyBasedServices: []*armnetwork.NetworkIntentPolicyBasedService{
	// 			to.Ptr(armnetwork.NetworkIntentPolicyBasedServiceNone)},
	// 			ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		},
	// 		SystemData: &armnetwork.SystemData{
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 			CreatedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 			CreatedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-01-11T18:52:27.000Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("b69a9388-9488-4534-b470-7ec6d41beef5"),
	// 			LastModifiedByType: to.Ptr(armnetwork.CreatedByTypeUser),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c58fa033619b12c7cfa8a0ec5a9bf03bb18869ab/specification/network/resource-manager/Microsoft.Network/stable/2024-07-01/examples/NetworkManagerSecurityAdminConfigurationDelete.json
func ExampleSecurityAdminConfigurationsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSecurityAdminConfigurationsClient().BeginDelete(ctx, "rg1", "testNetworkManager", "myTestSecurityConfig", &armnetwork.SecurityAdminConfigurationsClientBeginDeleteOptions{Force: to.Ptr(false)})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
