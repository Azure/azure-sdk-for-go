//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armauthorization_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/310a0100f5b020c1900c527a6aa70d21992f078a/specification/authorization/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/GetAlertConfigurationById.json
func ExampleAlertConfigurationsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armauthorization.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAlertConfigurationsClient().Get(ctx, "subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f", "TooManyOwnersAssignedToResource", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AlertConfiguration = armauthorization.AlertConfiguration{
	// 	Name: to.Ptr("TooManyOwnersAssignedToResource"),
	// 	Type: to.Ptr("Microsoft.Authorization/roleManagementAlertConfigurations"),
	// 	ID: to.Ptr("/subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f/providers/Microsoft.Authorization/roleManagementAlertConfigurations/TooManyOwnersAssignedToResource"),
	// 	Properties: &armauthorization.TooManyOwnersAssignedToResourceAlertConfigurationProperties{
	// 		AlertConfigurationType: to.Ptr("TooManyOwnersAssignedToResourceAlertConfiguration"),
	// 		AlertDefinitionID: to.Ptr("TooManyOwnersAssignedToResource"),
	// 		IsEnabled: to.Ptr(true),
	// 		Scope: to.Ptr("subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f"),
	// 		ThresholdNumberOfOwners: to.Ptr[int32](3),
	// 		ThresholdPercentageOfOwnersOutOfAllRoleMembers: to.Ptr[int32](40),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/310a0100f5b020c1900c527a6aa70d21992f078a/specification/authorization/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/UpdateAlertConfiguration.json
func ExampleAlertConfigurationsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armauthorization.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAlertConfigurationsClient().Update(ctx, "subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f", "TooManyOwnersAssignedToResource", armauthorization.AlertConfiguration{
		Properties: &armauthorization.TooManyOwnersAssignedToResourceAlertConfigurationProperties{
			AlertConfigurationType:  to.Ptr("TooManyOwnersAssignedToResourceAlertConfiguration"),
			IsEnabled:               to.Ptr(true),
			ThresholdNumberOfOwners: to.Ptr[int32](2),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/310a0100f5b020c1900c527a6aa70d21992f078a/specification/authorization/resource-manager/Microsoft.Authorization/preview/2022-08-01-preview/examples/GetAlertConfigurations.json
func ExampleAlertConfigurationsClient_NewListForScopePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armauthorization.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAlertConfigurationsClient().NewListForScopePager("subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f", nil)
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
		// page.AlertConfigurationListResult = armauthorization.AlertConfigurationListResult{
		// 	Value: []*armauthorization.AlertConfiguration{
		// 		{
		// 			Name: to.Ptr("DuplicateRoleCreated"),
		// 			Type: to.Ptr("Microsoft.Authorization/roleManagementAlertConfigurations"),
		// 			ID: to.Ptr("/subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f/providers/Microsoft.Authorization/roleManagementAlertConfigurations/DuplicateRoleCreated"),
		// 			Properties: &armauthorization.DuplicateRoleCreatedAlertConfigurationProperties{
		// 				AlertConfigurationType: to.Ptr("DuplicateRoleCreatedAlertConfiguration"),
		// 				AlertDefinitionID: to.Ptr("DuplicateRoleCreated"),
		// 				IsEnabled: to.Ptr(true),
		// 				Scope: to.Ptr("subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("TooManyOwnersAssignedToResource"),
		// 			Type: to.Ptr("Microsoft.Authorization/roleManagementAlertConfigurations"),
		// 			ID: to.Ptr("/subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f/providers/Microsoft.Authorization/roleManagementAlertConfigurations/TooManyOwnersAssignedToResource"),
		// 			Properties: &armauthorization.TooManyOwnersAssignedToResourceAlertConfigurationProperties{
		// 				AlertConfigurationType: to.Ptr("TooManyOwnersAssignedToResourceAlertConfiguration"),
		// 				AlertDefinitionID: to.Ptr("TooManyOwnersAssignedToResource"),
		// 				IsEnabled: to.Ptr(true),
		// 				Scope: to.Ptr("subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f"),
		// 				ThresholdNumberOfOwners: to.Ptr[int32](2),
		// 				ThresholdPercentageOfOwnersOutOfAllRoleMembers: to.Ptr[int32](3),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("TooManyPermanentOwnersAssignedToResource"),
		// 			Type: to.Ptr("Microsoft.Authorization/roleManagementAlertConfigurations"),
		// 			ID: to.Ptr("/subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f/providers/Microsoft.Authorization/roleManagementAlertConfigurations/TooManyPermanentOwnersAssignedToResource"),
		// 			Properties: &armauthorization.TooManyPermanentOwnersAssignedToResourceAlertConfigurationProperties{
		// 				AlertConfigurationType: to.Ptr("TooManyPermanentOwnersAssignedToResourceAlertConfiguration"),
		// 				AlertDefinitionID: to.Ptr("TooManyPermanentOwnersAssignedToResource"),
		// 				IsEnabled: to.Ptr(true),
		// 				Scope: to.Ptr("subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f"),
		// 				ThresholdNumberOfPermanentOwners: to.Ptr[int32](10),
		// 				ThresholdPercentageOfPermanentOwnersOutOfAllOwners: to.Ptr[int32](10),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("AzureRolesAssignedOutsidePimAlert"),
		// 			Type: to.Ptr("Microsoft.Authorization/roleManagementAlertConfigurations"),
		// 			ID: to.Ptr("/subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f/providers/Microsoft.Authorization/roleManagementAlertConfigurations/AzureRolesAssignedOutsidePimAlert"),
		// 			Properties: &armauthorization.AzureRolesAssignedOutsidePimAlertConfigurationProperties{
		// 				AlertConfigurationType: to.Ptr("AzureRolesAssignedOutsidePimAlertConfiguration"),
		// 				AlertDefinitionID: to.Ptr("AzureRolesAssignedOutsidePimAlert"),
		// 				IsEnabled: to.Ptr(true),
		// 				Scope: to.Ptr("subscriptions/afa2a084-766f-4003-8ae1-c4aeb893a99f"),
		// 			},
		// 	}},
		// }
	}
}
