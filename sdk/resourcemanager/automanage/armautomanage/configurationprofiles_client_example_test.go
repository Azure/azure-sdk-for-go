//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armautomanage_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/2dcad6d6e9a96882eb6d317e7500a94be007a9c6/specification/automanage/resource-manager/Microsoft.Automanage/stable/2022-05-04/examples/createOrUpdateConfigurationProfile.json
func ExampleConfigurationProfilesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armautomanage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewConfigurationProfilesClient().CreateOrUpdate(ctx, "customConfigurationProfile", "myResourceGroupName", armautomanage.ConfigurationProfile{
		Location: to.Ptr("East US"),
		Tags: map[string]*string{
			"Organization": to.Ptr("Administration"),
		},
		Properties: &armautomanage.ConfigurationProfileProperties{
			Configuration: map[string]any{
				"Antimalware/Enable":                false,
				"AzureSecurityCenter/Enable":        true,
				"Backup/Enable":                     false,
				"BootDiagnostics/Enable":            true,
				"ChangeTrackingAndInventory/Enable": true,
				"GuestConfiguration/Enable":         true,
				"LogAnalytics/Enable":               true,
				"UpdateManagement/Enable":           true,
				"VMInsights/Enable":                 true,
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ConfigurationProfile = armautomanage.ConfigurationProfile{
	// 	Name: to.Ptr("customConfigurationProfile"),
	// 	Type: to.Ptr("Microsoft.Automanage/configurationProfiles"),
	// 	ID: to.Ptr("/subscriptions/subscriptionId/resourceGroups/myResourceGroupName/providers/Microsoft.Automanage/configurationProfiles/customConfigurationProfile"),
	// 	Location: to.Ptr("East US"),
	// 	Tags: map[string]*string{
	// 		"Organization": to.Ptr("Administration"),
	// 	},
	// 	Properties: &armautomanage.ConfigurationProfileProperties{
	// 		Configuration: map[string]any{
	// 			"Antimalware/Enable": false,
	// 			"AzureSecurityCenter/Enable": true,
	// 			"Backup/Enable": false,
	// 			"BootDiagnostics/Enable": true,
	// 			"ChangeTrackingAndInventory/Enable": true,
	// 			"GuestConfiguration/Enable": true,
	// 			"LogAnalytics/Enable": true,
	// 			"UpdateManagement/Enable": true,
	// 			"VMInsights/Enable": true,
	// 		},
	// 	},
	// 	SystemData: &armautomanage.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.1075056Z"); return t}()),
	// 		CreatedBy: to.Ptr("user1@outlook.com"),
	// 		CreatedByType: to.Ptr(armautomanage.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.1974346Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user2@outlook.com"),
	// 		LastModifiedByType: to.Ptr(armautomanage.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/2dcad6d6e9a96882eb6d317e7500a94be007a9c6/specification/automanage/resource-manager/Microsoft.Automanage/stable/2022-05-04/examples/getConfigurationProfile.json
func ExampleConfigurationProfilesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armautomanage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewConfigurationProfilesClient().Get(ctx, "customConfigurationProfile", "myResourceGroupName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ConfigurationProfile = armautomanage.ConfigurationProfile{
	// 	Name: to.Ptr("customConfigurationProfile"),
	// 	Type: to.Ptr("Microsoft.Automanage/ConfigurationProfiles"),
	// 	ID: to.Ptr("/subscriptions/subscriptionId/resourceGroups/myResourceGroupName/providers/Microsoft.Automanage/configurationProfiles/customConfigurationProfile"),
	// 	Location: to.Ptr("East US"),
	// 	Tags: map[string]*string{
	// 		"Organization": to.Ptr("Administration"),
	// 	},
	// 	Properties: &armautomanage.ConfigurationProfileProperties{
	// 		Configuration: map[string]any{
	// 			"Antimalware/Enable": false,
	// 			"AzureSecurityCenter/Enable": true,
	// 			"Backup/Enable": false,
	// 			"BootDiagnostics/Enable": true,
	// 			"ChangeTrackingAndInventory/Enable": true,
	// 			"GuestConfiguration/Enable": true,
	// 			"LogAnalytics/Enable": true,
	// 			"UpdateManagement/Enable": true,
	// 			"VMInsights/Enable": true,
	// 		},
	// 	},
	// 	SystemData: &armautomanage.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.1075056Z"); return t}()),
	// 		CreatedBy: to.Ptr("user1@outlook.com"),
	// 		CreatedByType: to.Ptr(armautomanage.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.1974346Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user2@outlook.com"),
	// 		LastModifiedByType: to.Ptr(armautomanage.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/2dcad6d6e9a96882eb6d317e7500a94be007a9c6/specification/automanage/resource-manager/Microsoft.Automanage/stable/2022-05-04/examples/deleteConfigurationProfile.json
func ExampleConfigurationProfilesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armautomanage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewConfigurationProfilesClient().Delete(ctx, "rg", "customConfigurationProfile", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/2dcad6d6e9a96882eb6d317e7500a94be007a9c6/specification/automanage/resource-manager/Microsoft.Automanage/stable/2022-05-04/examples/updateConfigurationProfile.json
func ExampleConfigurationProfilesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armautomanage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewConfigurationProfilesClient().Update(ctx, "customConfigurationProfile", "myResourceGroupName", armautomanage.ConfigurationProfileUpdate{
		Tags: map[string]*string{
			"Organization": to.Ptr("Administration"),
		},
		Properties: &armautomanage.ConfigurationProfileProperties{
			Configuration: map[string]any{
				"Antimalware/Enable":                false,
				"AzureSecurityCenter/Enable":        true,
				"Backup/Enable":                     false,
				"BootDiagnostics/Enable":            true,
				"ChangeTrackingAndInventory/Enable": true,
				"GuestConfiguration/Enable":         true,
				"LogAnalytics/Enable":               true,
				"UpdateManagement/Enable":           true,
				"VMInsights/Enable":                 true,
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ConfigurationProfile = armautomanage.ConfigurationProfile{
	// 	Name: to.Ptr("customConfigurationProfile"),
	// 	Type: to.Ptr("Microsoft.Automanage/configurationProfiles"),
	// 	ID: to.Ptr("/subscriptions/subscriptionId/resourceGroups/myResourceGroupName/providers/Microsoft.Automanage/configurationProfiles/customConfigurationProfile"),
	// 	Location: to.Ptr("East US"),
	// 	Tags: map[string]*string{
	// 		"Organization": to.Ptr("Administration"),
	// 	},
	// 	Properties: &armautomanage.ConfigurationProfileProperties{
	// 		Configuration: map[string]any{
	// 			"Antimalware/Enable": false,
	// 			"AzureSecurityCenter/Enable": true,
	// 			"Backup/Enable": false,
	// 			"BootDiagnostics/Enable": true,
	// 			"ChangeTrackingAndInventory/Enable": true,
	// 			"GuestConfiguration/Enable": true,
	// 			"LogAnalytics/Enable": true,
	// 			"UpdateManagement/Enable": true,
	// 			"VMInsights/Enable": true,
	// 		},
	// 	},
	// 	SystemData: &armautomanage.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.1075056Z"); return t}()),
	// 		CreatedBy: to.Ptr("user1@outlook.com"),
	// 		CreatedByType: to.Ptr(armautomanage.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.1974346Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("user2@outlook.com"),
	// 		LastModifiedByType: to.Ptr(armautomanage.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/2dcad6d6e9a96882eb6d317e7500a94be007a9c6/specification/automanage/resource-manager/Microsoft.Automanage/stable/2022-05-04/examples/listConfigurationProfilesByResourceGroup.json
func ExampleConfigurationProfilesClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armautomanage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewConfigurationProfilesClient().NewListByResourceGroupPager("myResourceGroupName", nil)
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
		// page.ConfigurationProfileList = armautomanage.ConfigurationProfileList{
		// 	Value: []*armautomanage.ConfigurationProfile{
		// 		{
		// 			Name: to.Ptr("customConfigurationProfile"),
		// 			Type: to.Ptr("Microsoft.Automanage/ConfigurationProfiles"),
		// 			ID: to.Ptr("/subscriptions/subscriptionId/resourceGroups/myResourceGroupName/providers/Microsoft.Automanage/configurationProfiles/customConfigurationProfile"),
		// 			Location: to.Ptr("East US"),
		// 			Tags: map[string]*string{
		// 				"Organization": to.Ptr("Administration"),
		// 			},
		// 			Properties: &armautomanage.ConfigurationProfileProperties{
		// 				Configuration: map[string]any{
		// 					"Antimalware/Enable": false,
		// 					"AzureSecurityCenter/Enable": true,
		// 					"Backup/Enable": false,
		// 					"BootDiagnostics/Enable": true,
		// 					"ChangeTrackingAndInventory/Enable": true,
		// 					"GuestConfiguration/Enable": true,
		// 					"LogAnalytics/Enable": true,
		// 					"UpdateManagement/Enable": true,
		// 					"VMInsights/Enable": true,
		// 				},
		// 			},
		// 			SystemData: &armautomanage.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("user1@outlook.com"),
		// 				CreatedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.1974346Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user2@outlook.com"),
		// 				LastModifiedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("customConfigurationProfile2"),
		// 			Type: to.Ptr("Microsoft.Automanage/ConfigurationProfiles"),
		// 			ID: to.Ptr("/subscriptions/subscriptionId/resourceGroups/myResourceGroupName/providers/Microsoft.Automanage/configurationProfiles/customConfigurationProfile2"),
		// 			Location: to.Ptr("East US"),
		// 			Tags: map[string]*string{
		// 				"Organization": to.Ptr("Administration"),
		// 			},
		// 			Properties: &armautomanage.ConfigurationProfileProperties{
		// 				Configuration: map[string]any{
		// 					"Antimalware/Enable": false,
		// 					"AzureSecurityCenter/Enable": true,
		// 					"Backup/Enable": false,
		// 					"BootDiagnostics/Enable": true,
		// 					"ChangeTrackingAndInventory/Enable": true,
		// 					"GuestConfiguration/Enable": true,
		// 					"LogAnalytics/Enable": true,
		// 					"UpdateManagement/Enable": true,
		// 					"VMInsights/Enable": true,
		// 				},
		// 			},
		// 			SystemData: &armautomanage.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("user1@outlook.com"),
		// 				CreatedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.1974346Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user2@outlook.com"),
		// 				LastModifiedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/2dcad6d6e9a96882eb6d317e7500a94be007a9c6/specification/automanage/resource-manager/Microsoft.Automanage/stable/2022-05-04/examples/listConfigurationProfilesBySubscription.json
func ExampleConfigurationProfilesClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armautomanage.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewConfigurationProfilesClient().NewListBySubscriptionPager(nil)
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
		// page.ConfigurationProfileList = armautomanage.ConfigurationProfileList{
		// 	Value: []*armautomanage.ConfigurationProfile{
		// 		{
		// 			Name: to.Ptr("customConfigurationProfile"),
		// 			Type: to.Ptr("Microsoft.Automanage/ConfigurationProfiles"),
		// 			ID: to.Ptr("/subscriptions/subscriptionId/resourceGroups/myResourceGroupName/providers/Microsoft.Automanage/configurationProfiles/customConfigurationProfile"),
		// 			Location: to.Ptr("East US"),
		// 			Tags: map[string]*string{
		// 				"Organization": to.Ptr("Administration"),
		// 			},
		// 			Properties: &armautomanage.ConfigurationProfileProperties{
		// 				Configuration: map[string]any{
		// 					"Antimalware/Enable": false,
		// 					"AzureSecurityCenter/Enable": true,
		// 					"Backup/Enable": false,
		// 					"BootDiagnostics/Enable": true,
		// 					"ChangeTrackingAndInventory/Enable": true,
		// 					"GuestConfiguration/Enable": true,
		// 					"LogAnalytics/Enable": true,
		// 					"UpdateManagement/Enable": true,
		// 					"VMInsights/Enable": true,
		// 				},
		// 			},
		// 			SystemData: &armautomanage.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("user1@outlook.com"),
		// 				CreatedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.1974346Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user2@outlook.com"),
		// 				LastModifiedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("customConfigurationProfile2"),
		// 			Type: to.Ptr("Microsoft.Automanage/ConfigurationProfiles"),
		// 			ID: to.Ptr("/subscriptions/subscriptionId/resourceGroups/myResourceGroupName/providers/Microsoft.Automanage/configurationProfiles/customConfigurationProfile2"),
		// 			Location: to.Ptr("East US"),
		// 			Tags: map[string]*string{
		// 				"Organization": to.Ptr("Administration"),
		// 			},
		// 			Properties: &armautomanage.ConfigurationProfileProperties{
		// 				Configuration: map[string]any{
		// 					"Antimalware/Enable": false,
		// 					"AzureSecurityCenter/Enable": true,
		// 					"Backup/Enable": false,
		// 					"BootDiagnostics/Enable": true,
		// 					"ChangeTrackingAndInventory/Enable": true,
		// 					"GuestConfiguration/Enable": true,
		// 					"LogAnalytics/Enable": true,
		// 					"UpdateManagement/Enable": true,
		// 					"VMInsights/Enable": true,
		// 				},
		// 			},
		// 			SystemData: &armautomanage.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-03T01:01:01.1075056Z"); return t}()),
		// 				CreatedBy: to.Ptr("user1@outlook.com"),
		// 				CreatedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-02-04T02:03:01.1974346Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("user2@outlook.com"),
		// 				LastModifiedByType: to.Ptr(armautomanage.CreatedByTypeUser),
		// 			},
		// 	}},
		// }
	}
}
