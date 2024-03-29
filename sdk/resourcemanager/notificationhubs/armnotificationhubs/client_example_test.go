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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/notificationhubs/armnotificationhubs/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/CheckAvailability.json
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
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/checkNamespaceAvailability"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/sdktest"),
	// 	IsAvailiable: to.Ptr(true),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/Get.json
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
	// 	Name: to.Ptr("test"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test"),
	// 	Location: to.Ptr("East US"),
	// 	Tags: map[string]*string{
	// 		"hubTag1": to.Ptr("hubTagValue1"),
	// 		"hubTag2": to.Ptr("hubTagValue2"),
	// 	},
	// 	Properties: &armnotificationhubs.NotificationHubProperties{
	// 		Name: to.Ptr("test"),
	// 		DailyMaxActiveDevices: to.Ptr[int64](0),
	// 		RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/CreateOrUpdate.json
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
	res, err := clientFactory.NewClient().CreateOrUpdate(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", armnotificationhubs.NotificationHubResource{
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
	// 	Name: to.Ptr("test"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test"),
	// 	Location: to.Ptr("East US"),
	// 	Tags: map[string]*string{
	// 		"hubTag1": to.Ptr("hubTagValue1"),
	// 		"hubTag2": to.Ptr("hubTagValue2"),
	// 	},
	// 	Properties: &armnotificationhubs.NotificationHubProperties{
	// 		Name: to.Ptr("test"),
	// 		DailyMaxActiveDevices: to.Ptr[int64](0),
	// 		RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/Update.json
func ExampleClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnotificationhubs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewClient().Update(ctx, "sdkresourceGroup", "nh-sdk-ns", "sdk-notificationHubs-8708", armnotificationhubs.NotificationHubPatchParameters{
		Properties: &armnotificationhubs.NotificationHubProperties{
			GCMCredential: &armnotificationhubs.GCMCredential{
				Properties: &armnotificationhubs.GCMCredentialProperties{
					GCMEndpoint:  to.Ptr("https://fcm.googleapis.com/fcm/send"),
					GoogleAPIKey: to.Ptr("###################################"),
				},
			},
			RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.NotificationHubResource = armnotificationhubs.NotificationHubResource{
	// 	Name: to.Ptr("test"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test"),
	// 	Location: to.Ptr("East US"),
	// 	Tags: map[string]*string{
	// 		"hubTag1": to.Ptr("hubTagValue1"),
	// 		"hubTag2": to.Ptr("hubTagValue2"),
	// 	},
	// 	Properties: &armnotificationhubs.NotificationHubProperties{
	// 		Name: to.Ptr("test"),
	// 		DailyMaxActiveDevices: to.Ptr[int64](0),
	// 		RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/Delete.json
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/List.json
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
	pager := clientFactory.NewClient().NewListPager("5ktrial", "nh-sdk-ns", &armnotificationhubs.ClientListOptions{SkipToken: nil,
		Top: nil,
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
		// page.NotificationHubListResult = armnotificationhubs.NotificationHubListResult{
		// 	Value: []*armnotificationhubs.NotificationHubResource{
		// 		{
		// 			Name: to.Ptr("test"),
		// 			Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs"),
		// 			ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test"),
		// 			Location: to.Ptr("East US"),
		// 			Tags: map[string]*string{
		// 				"hubTag1": to.Ptr("hubTagValue1"),
		// 				"hubTag2": to.Ptr("hubTagValue2"),
		// 			},
		// 			Properties: &armnotificationhubs.NotificationHubProperties{
		// 				Name: to.Ptr("test"),
		// 				DailyMaxActiveDevices: to.Ptr[int64](0),
		// 				RegistrationTTL: to.Ptr("10675199.02:48:05.4775807"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/DebugSend.json
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
	res, err := clientFactory.NewClient().DebugSend(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DebugSendResponse = armnotificationhubs.DebugSendResponse{
	// 	Name: to.Ptr("test"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/debugSend"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test"),
	// 	Properties: &armnotificationhubs.DebugSendResult{
	// 		Failure: to.Ptr[int64](0),
	// 		Results: []*armnotificationhubs.RegistrationResult{
	// 		},
	// 		Success: to.Ptr[int64](0),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/AuthorizationRuleCreateOrUpdate.json
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
	res, err := clientFactory.NewClient().CreateOrUpdateAuthorizationRule(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", "MyManageSharedAccessKey", armnotificationhubs.SharedAccessAuthorizationRuleResource{
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
	// 	Name: to.Ptr("MyManageSharedAccessKey"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test/authorizationRules/MyManageSharedAccessKey"),
	// 	Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
	// 		CreatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T12:24:40.586Z"); return t}()),
	// 		ModifiedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T12:24:40.586Z"); return t}()),
	// 		Rights: []*armnotificationhubs.AccessRights{
	// 			to.Ptr(armnotificationhubs.AccessRightsListen),
	// 			to.Ptr(armnotificationhubs.AccessRightsSend)},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/AuthorizationRuleDelete.json
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/AuthorizationRuleGet.json
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
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test/authorizationRules/DefaultListenSharedAccessSignature"),
	// 	Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
	// 		CreatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T12:24:40.586Z"); return t}()),
	// 		ModifiedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T12:24:40.586Z"); return t}()),
	// 		Rights: []*armnotificationhubs.AccessRights{
	// 			to.Ptr(armnotificationhubs.AccessRightsListen),
	// 			to.Ptr(armnotificationhubs.AccessRightsSend)},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/AuthorizationRuleList.json
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
		// 			ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test/authorizationRules/DefaultListenSharedAccessSignature"),
		// 			Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
		// 				CreatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T10:43:00.532Z"); return t}()),
		// 				ModifiedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T10:43:00.532Z"); return t}()),
		// 				Rights: []*armnotificationhubs.AccessRights{
		// 					to.Ptr(armnotificationhubs.AccessRightsListen)},
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("DefaultFullSharedAccessSignature"),
		// 				Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/authorizationRules"),
		// 				ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test/authorizationRules/DefaultFullSharedAccessSignature"),
		// 				Properties: &armnotificationhubs.SharedAccessAuthorizationRuleProperties{
		// 					CreatedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T10:43:00.532Z"); return t}()),
		// 					ModifiedTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-04-26T10:43:00.532Z"); return t}()),
		// 					Rights: []*armnotificationhubs.AccessRights{
		// 						to.Ptr(armnotificationhubs.AccessRightsManage),
		// 						to.Ptr(armnotificationhubs.AccessRightsListen),
		// 						to.Ptr(armnotificationhubs.AccessRightsSend)},
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/AuthorizationRuleListKeys.json
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
	// 	PrimaryConnectionString: to.Ptr("Endpoint=sb://nh-sdk-ns.servicebus.windows-int.net/;SharedAccessKeyName=sdk-AuthRules-5800;SharedAccessKey=############################################;EntityPath=sdk-notificationHubs-2317"),
	// 	PrimaryKey: to.Ptr("############################################"),
	// 	SecondaryConnectionString: to.Ptr("Endpoint=sb://nh-sdk-ns.servicebus.windows-int.net/;SharedAccessKeyName=sdk-AuthRules-5800;SharedAccessKey=############################################;EntityPath=sdk-notificationHubs-2317"),
	// 	SecondaryKey: to.Ptr("############################################"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/AuthorizationRuleRegenerateKey.json
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
	res, err := clientFactory.NewClient().RegenerateKeys(ctx, "5ktrial", "nh-sdk-ns", "nh-sdk-hub", "DefaultListenSharedAccessSignature", armnotificationhubs.PolicyKeyResource{
		PolicyKey: to.Ptr(armnotificationhubs.PolicyKeyTypePrimaryKey),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ResourceListKeys = armnotificationhubs.ResourceListKeys{
	// 	KeyName: to.Ptr("DefaultListenSharedAccessSignature"),
	// 	PrimaryConnectionString: to.Ptr("Endpoint=sb://nh-sdk-ns.servicebus.windows-int.net/;SharedAccessKeyName=sdk-AuthRules-5800;SharedAccessKey=############################################;EntityPath=sdk-notificationHubs-2317"),
	// 	PrimaryKey: to.Ptr("############################################"),
	// 	SecondaryConnectionString: to.Ptr("Endpoint=sb://nh-sdk-ns.servicebus.windows-int.net/;SharedAccessKeyName=sdk-AuthRules-5800;SharedAccessKey=############################################;EntityPath=sdk-notificationHubs-2317"),
	// 	SecondaryKey: to.Ptr("############################################"),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/85cfba195a19120f309bd292c4261aa53a586adb/specification/notificationhubs/resource-manager/Microsoft.NotificationHubs/preview/2023-10-01-preview/examples/NotificationHubs/PnsCredentialsGet.json
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
	// 	Name: to.Ptr("test"),
	// 	Type: to.Ptr("Microsoft.NotificationHubs/namespaces/notificationHubs/pnsCredentials"),
	// 	ID: to.Ptr("/subscriptions/29cfa613-cbbc-4512-b1d6-1b3a92c7fa40/resourceGroups/5ktrial/providers/Microsoft.NotificationHubs/namespaces/nh-sdk-ns/notificationHubs/test"),
	// 	Properties: &armnotificationhubs.PnsCredentials{
	// 		GCMCredential: &armnotificationhubs.GCMCredential{
	// 			Properties: &armnotificationhubs.GCMCredentialProperties{
	// 				GCMEndpoint: to.Ptr("https://fcm.googleapis.com/fcm/send"),
	// 				GoogleAPIKey: to.Ptr("###################################"),
	// 			},
	// 		},
	// 	},
	// }
}
