//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnotificationhubs_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/notificationhubs/armnotificationhubs"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubCheckNameAvailability.json
func ExampleClient_CheckNotificationHubAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().CheckNotificationHubAvailability(ctx, "5ktrial", "locp-newns", armnotificationhubs.CheckAvailabilityParameters{
		Name:     to.Ptr("sdktest"),
		Location: to.Ptr("West Europe"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.CheckAvailabilityResult = armnotificationhubs.CheckAvailabilityResult{
	// 	Name: to.Ptr("sdktest"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/checkNotificationHubAvailability"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourcegroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/locp-newns/CheckNotificationHubAvailability"),
	// 	Location: to.Ptr("West Europe"),
	// 	IsAvailiable: to.Ptr(true),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubCreate.json
func ExampleClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().CreateOrUpdate(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", armnotificationhubs.NotificationHubCreateOrUpdateParameters{
		Location:   to.Ptr("eastus"),
		Properties: &armnotificationhubs.NotificationHubProperties{},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NotificationHubResource = armnotificationhubs.NotificationHubResource{
	// 	Name: to.Ptr("nh-sdk-hub"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/sdkresourceGroup/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/nh-sdk-hub"),
	// 	Location: to.Ptr("eastus"),
	// 	Properties: &armnotificationhubs.NotificationHubProperties{
	// 		AuthorizationRules: []*armnotificationhubs.SharedAccessAuthorizationRuleProperties{
	// 		},
	// 		RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubPatch.json
func ExampleClient_Patch() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().Patch(ctx, "sdkresourceGroup", "nh-sdk-ns", "sdk-notificationHubs-8708", &armnotificationhubs.ClientPatchOptions{Parameters: &armnotificationhubs.NotificationHubPatchParameters{}})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NotificationHubResource = armnotificationhubs.NotificationHubResource{
	// 	Name: to.Ptr("nh-sdk-hub"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/sdkresourceGroup/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/nh-sdk-hub"),
	// 	Location: to.Ptr("South Central US"),
	// 	Properties: &armnotificationhubs.NotificationHubProperties{
	// 		AuthorizationRules: []*armnotificationhubs.SharedAccessAuthorizationRuleProperties{
	// 		},
	// 		RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubDelete.json
func ExampleClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewClient().Delete(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubGet.json
func ExampleClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().Get(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NotificationHubResource = armnotificationhubs.NotificationHubResource{
	// 	Name: to.Ptr("nh-sdk-hub"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/sdkresourceGroup/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/nh-sdk-hub"),
	// 	Location: to.Ptr("South Central US"),
	// 	Properties: &armnotificationhubs.NotificationHubProperties{
	// 		AuthorizationRules: []*armnotificationhubs.SharedAccessAuthorizationRuleProperties{
	// 		},
	// 		RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubDebugSend.json
func ExampleClient_DebugSend() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewClient().DebugSend(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", &armnotificationhubs.ClientDebugSendOptions{Parameters: map[string]any{
		"data": map[string]any{
			"message": "Hello",
		},
	},
	})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubAuthorizationRuleCreate.json
func ExampleClient_CreateOrUpdateAuthorizationRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().CreateOrUpdateAuthorizationRule(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", "DefaultListenSharedAccessSignature", armnotificationhubs.SharedAccessAuthorizationRuleCreateOrUpdateParameters{
		Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
			Rights: []*armnotificationhubs.AccessRights{
				to.Ptr(armnotificationhubs.AccessRightsListen),
				to.Ptr(armnotificationhubs.AccessRightsSend)},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SharedAccessAuthorizationRuleResource = armnotificationhubs.SharedAccessAuthorizationRuleResource{
	// 	Name: to.Ptr("DefaultListenSharedAccessSignature"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/NotificationHubs/nh-sdk-hub/AuthorizationRules/DefaultListenSharedAccessSignature"),
	// 	Location: to.Ptr("West Europe"),
	// 	Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
	// 		ClaimType: to.Ptr("SharedAccessKey"),
	// 		ClaimValue: to.Ptr("None"),
	// 		CreatedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
	// 		KeyName: to.Ptr("DefaultListenSharedAccessSignature"),
	// 		ModifiedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
	// 		PrimaryKey: to.Ptr("#################################"),
	// 		Rights: []*armnotificationhubs.AccessRights{
	// 			to.Ptr(armnotificationhubs.AccessRightsListen)},
	// 			SecondaryKey: to.Ptr("#################################"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubAuthorizationRuleDelete.json
func ExampleClient_DeleteAuthorizationRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewClient().DeleteAuthorizationRule(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", "DefaultListenSharedAccessSignature", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubAuthorizationRuleGet.json
func ExampleClient_GetAuthorizationRule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().GetAuthorizationRule(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", "DefaultListenSharedAccessSignature", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SharedAccessAuthorizationRuleResource = armnotificationhubs.SharedAccessAuthorizationRuleResource{
	// 	Name: to.Ptr("DefaultListenSharedAccessSignature"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/NotificationHubs/nh-sdk-hub/AuthorizationRules/DefaultListenSharedAccessSignature"),
	// 	Location: to.Ptr("West Europe"),
	// 	Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
	// 		ClaimType: to.Ptr("SharedAccessKey"),
	// 		ClaimValue: to.Ptr("None"),
	// 		CreatedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
	// 		KeyName: to.Ptr("DefaultListenSharedAccessSignature"),
	// 		ModifiedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
	// 		PrimaryKey: to.Ptr("#################################"),
	// 		Rights: []*armnotificationhubs.AccessRights{
	// 			to.Ptr(armnotificationhubs.AccessRightsListen)},
	// 			SecondaryKey: to.Ptr("#################################"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubListByNameSpace.json
func ExampleClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewClient().NewListPager("5ktrial", "nh-sdk-ns", nil)
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
		// page.NotificationHubListResult = armnotificationhubs.NotificationHubListResult{
		// 	Value: []*armnotificationhubs.NotificationHubResource{
		// 		{
		// 			Name: to.Ptr("nh-sdk-hub"),
		// 			Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
		// 			ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/NotificationHubs/nh-sdk-hub"),
		// 			Location: to.Ptr("South Central US"),
		// 			Properties: &armnotificationhubs.NotificationHubProperties{
		// 				AuthorizationRules: []*armnotificationhubs.SharedAccessAuthorizationRuleProperties{
		// 				},
		// 				RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubAuthorizationRuleListAll.json
func ExampleClient_NewListAuthorizationRulesPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewClient().NewListAuthorizationRulesPager("5ktrial", "nh-sdk-ns", "nh-sdk-hub", nil)
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
		// page.SharedAccessAuthorizationRuleListResult = armnotificationhubs.SharedAccessAuthorizationRuleListResult{
		// 	Value: []*armnotificationhubs.SharedAccessAuthorizationRuleResource{
		// 		{
		// 			Name: to.Ptr("DefaultListenSharedAccessSignature"),
		// 			Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules"),
		// 			ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/NotificationHubs/nh-sdk-hub/AuthorizationRules/DefaultListenSharedAccessSignature"),
		// 			Location: to.Ptr("West Europe"),
		// 			Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
		// 				ClaimType: to.Ptr("SharedAccessKey"),
		// 				ClaimValue: to.Ptr("None"),
		// 				CreatedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
		// 				KeyName: to.Ptr("DefaultListenSharedAccessSignature"),
		// 				ModifiedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
		// 				PrimaryKey: to.Ptr("#################################"),
		// 				Rights: []*armnotificationhubs.AccessRights{
		// 					to.Ptr(armnotificationhubs.AccessRightsListen)},
		// 					SecondaryKey: to.Ptr("#################################"),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("DefaultFullSharedAccessSignature"),
		// 				Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules"),
		// 				ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/NotificationHubs/nh-sdk-hub/AuthorizationRules/DefaultFullSharedAccessSignature"),
		// 				Location: to.Ptr("West Europe"),
		// 				Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
		// 					ClaimType: to.Ptr("SharedAccessKey"),
		// 					ClaimValue: to.Ptr("None"),
		// 					CreatedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
		// 					KeyName: to.Ptr("DefaultFullSharedAccessSignature"),
		// 					ModifiedTime: to.Ptr("2018-05-02T00:45:22.0150024Z"),
		// 					PrimaryKey: to.Ptr("#################################"),
		// 					Rights: []*armnotificationhubs.AccessRights{
		// 						to.Ptr(armnotificationhubs.AccessRightsListen),
		// 						to.Ptr(armnotificationhubs.AccessRightsManage),
		// 						to.Ptr(armnotificationhubs.AccessRightsSend)},
		// 						SecondaryKey: to.Ptr("#################################"),
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubAuthorizationRuleListKey.json
func ExampleClient_ListKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().ListKeys(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", "sdk-AuthRules-5800", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ResourceListKeys = armnotificationhubs.ResourceListKeys{
	// 	KeyName: to.Ptr("sdk-AuthRules-5800"),
	// 	PrimaryConnectionString: to.Ptr("Endpoint=sb://sdk-namespace-7982.servicebus.windows-int.net/;SharedAccessKeyName=sdk-AuthRules-5800;SharedAccessKey=############################################;EntityPath=sdk-notificationHubs-2317"),
	// 	PrimaryKey: to.Ptr("############################################"),
	// 	SecondaryConnectionString: to.Ptr("Endpoint=sb://sdk-namespace-7982.servicebus.windows-int.net/;SharedAccessKeyName=sdk-AuthRules-5800;SharedAccessKey=############################################;EntityPath=sdk-notificationHubs-2317"),
	// 	SecondaryKey: to.Ptr("############################################"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubAuthorizationRuleRegenrateKey.json
func ExampleClient_RegenerateKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().RegenerateKeys(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", "DefaultListenSharedAccessSignature", armnotificationhubs.PolicykeyResource{
		PolicyKey: to.Ptr("PrimaryKey"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ResourceListKeys = armnotificationhubs.ResourceListKeys{
	// 	KeyName: to.Ptr("DefaultListenSharedAccessSignature"),
	// 	PrimaryConnectionString: to.Ptr("Endpoint=sb://nh-sdk-ns.servicebus.windows.net/;SharedAccessKeyName=DefaultListenSharedAccessSignature;SharedAccessKey=#################################"),
	// 	PrimaryKey: to.Ptr("#################################"),
	// 	SecondaryConnectionString: to.Ptr("Endpoint=sb://nh-sdk-ns.servicebus.windows.net/;SharedAccessKeyName=DefaultListenSharedAccessSignature;SharedAccessKey=#################################"),
	// 	SecondaryKey: to.Ptr("#################################"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/7a2ac91de424f271cf91cc8009f3fe9ee8249086/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/stable/2017-04-01/examples/NotificationHubs/NotificationHubPnsCredentials.json
func ExampleClient_GetPnsCredentials() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().GetPnsCredentials(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PnsCredentialsResource = armnotificationhubs.PnsCredentialsResource{
	// 	Name: to.Ptr("nh-sdk-hub"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/pnsCredentials"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/nh-sdk-hub/pnsCredentials"),
	// 	Location: to.Ptr("West Europe"),
	// 	Properties: &armnotificationhubs.PnsCredentialsProperties{
	// 		MpnsCredential: &armnotificationhubs.MpnsCredential{
	// 			Properties: &armnotificationhubs.MpnsCredentialProperties{
	// 				Thumbprint: to.Ptr("#################################"),
	// 			},
	// 		},
	// 	},
	// }
}
