//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armlabservices_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/labservices/armlabservices"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/LabPlans/listLabPlans.json
func ExampleLabPlansClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewLabPlansClient().NewListBySubscriptionPager(&armlabservices.LabPlansClientListBySubscriptionOptions{Filter: nil})
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
		// page.PagedLabPlans = armlabservices.PagedLabPlans{
		// 	Value: []*armlabservices.LabPlan{
		// 		{
		// 			Name: to.Ptr("testlabplan"),
		// 			Type: to.Ptr("Microsoft.LabServices/LabPlan"),
		// 			ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labPlans/testlabplan"),
		// 			Location: to.Ptr("westus"),
		// 			Identity: &armlabservices.Identity{
		// 				Type: to.Ptr("SystemAssigned"),
		// 				PrincipalID: to.Ptr("a799d834-d3c9-4554-ba5d-3771a995aa6a"),
		// 				TenantID: to.Ptr("8c33e124-0581-45e0-84d4-f80d59de7806"),
		// 			},
		// 			Properties: &armlabservices.LabPlanProperties{
		// 				DefaultAutoShutdownProfile: &armlabservices.AutoShutdownProfile{
		// 					DisconnectDelay: to.Ptr("PT5M"),
		// 					IdleDelay: to.Ptr("PT15M"),
		// 					NoConnectDelay: to.Ptr("PT15M"),
		// 					ShutdownOnDisconnect: to.Ptr(armlabservices.EnableStateEnabled),
		// 					ShutdownOnIdle: to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
		// 					ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
		// 				},
		// 				DefaultConnectionProfile: &armlabservices.ConnectionProfile{
		// 					ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
		// 					ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
		// 					WebRdpAccess: to.Ptr(armlabservices.ConnectionTypeNone),
		// 					WebSSHAccess: to.Ptr(armlabservices.ConnectionTypeNone),
		// 				},
		// 				DefaultNetworkProfile: &armlabservices.LabPlanNetworkProfile{
		// 					SubnetID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default"),
		// 				},
		// 				SharedGalleryID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Compute/galleries/testsig"),
		// 				SupportInfo: &armlabservices.SupportInfo{
		// 					Email: to.Ptr("help@contoso.com"),
		// 					Instructions: to.Ptr("Contact support for help."),
		// 					Phone: to.Ptr("+1-202-555-0123"),
		// 					URL: to.Ptr("help.contoso.com"),
		// 				},
		// 				ProvisioningState: to.Ptr(armlabservices.ProvisioningStateSucceeded),
		// 			},
		// 			SystemData: &armlabservices.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T10:00:00.000Z"); return t}()),
		// 				CreatedBy: to.Ptr("identity123"),
		// 				CreatedByType: to.Ptr(armlabservices.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-06-01T09:12:28.000Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("identity123"),
		// 				LastModifiedByType: to.Ptr(armlabservices.CreatedByTypeUser),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/LabPlans/listResourceGroupLabPlans.json
func ExampleLabPlansClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewLabPlansClient().NewListByResourceGroupPager("testrg123", nil)
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
		// page.PagedLabPlans = armlabservices.PagedLabPlans{
		// 	Value: []*armlabservices.LabPlan{
		// 		{
		// 			Name: to.Ptr("testlabplan"),
		// 			Type: to.Ptr("Microsoft.LabServices/LabPlan"),
		// 			ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labPlans/testlabplan"),
		// 			Location: to.Ptr("westus"),
		// 			Properties: &armlabservices.LabPlanProperties{
		// 				DefaultAutoShutdownProfile: &armlabservices.AutoShutdownProfile{
		// 					DisconnectDelay: to.Ptr("PT5M"),
		// 					IdleDelay: to.Ptr("PT5M"),
		// 					NoConnectDelay: to.Ptr("PT5M"),
		// 					ShutdownOnDisconnect: to.Ptr(armlabservices.EnableStateEnabled),
		// 					ShutdownOnIdle: to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
		// 					ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
		// 				},
		// 				DefaultConnectionProfile: &armlabservices.ConnectionProfile{
		// 					ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
		// 					ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
		// 					WebRdpAccess: to.Ptr(armlabservices.ConnectionTypeNone),
		// 					WebSSHAccess: to.Ptr(armlabservices.ConnectionTypeNone),
		// 				},
		// 				DefaultNetworkProfile: &armlabservices.LabPlanNetworkProfile{
		// 					SubnetID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default"),
		// 				},
		// 				SharedGalleryID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Compute/galleries/testsig"),
		// 				SupportInfo: &armlabservices.SupportInfo{
		// 					Email: to.Ptr("help@contoso.com"),
		// 					Instructions: to.Ptr("Contact support for help."),
		// 					Phone: to.Ptr("+1-202-555-0123"),
		// 					URL: to.Ptr("help.contoso.com"),
		// 				},
		// 				ProvisioningState: to.Ptr(armlabservices.ProvisioningStateSucceeded),
		// 			},
		// 			SystemData: &armlabservices.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T10:00:00.000Z"); return t}()),
		// 				CreatedBy: to.Ptr("identity123"),
		// 				CreatedByType: to.Ptr(armlabservices.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-06-01T09:12:28.000Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("identity123"),
		// 				LastModifiedByType: to.Ptr(armlabservices.CreatedByTypeUser),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/LabPlans/getLabPlan.json
func ExampleLabPlansClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewLabPlansClient().Get(ctx, "testrg123", "testlabplan", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.LabPlan = armlabservices.LabPlan{
	// 	Name: to.Ptr("testlabplan"),
	// 	Type: to.Ptr("Microsoft.LabServices/LabPlan"),
	// 	ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labPlans/testlabplan"),
	// 	Location: to.Ptr("westus"),
	// 	Identity: &armlabservices.Identity{
	// 		Type: to.Ptr("SystemAssigned"),
	// 		PrincipalID: to.Ptr("a799d834-d3c9-4554-ba5d-3771a995aa6a"),
	// 		TenantID: to.Ptr("8c33e124-0581-45e0-84d4-f80d59de7806"),
	// 	},
	// 	Properties: &armlabservices.LabPlanProperties{
	// 		DefaultAutoShutdownProfile: &armlabservices.AutoShutdownProfile{
	// 			DisconnectDelay: to.Ptr("PT5M"),
	// 			IdleDelay: to.Ptr("PT5M"),
	// 			NoConnectDelay: to.Ptr("PT5M"),
	// 			ShutdownOnDisconnect: to.Ptr(armlabservices.EnableStateEnabled),
	// 			ShutdownOnIdle: to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
	// 			ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
	// 		},
	// 		DefaultConnectionProfile: &armlabservices.ConnectionProfile{
	// 			ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
	// 			ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
	// 			WebRdpAccess: to.Ptr(armlabservices.ConnectionTypeNone),
	// 			WebSSHAccess: to.Ptr(armlabservices.ConnectionTypeNone),
	// 		},
	// 		DefaultNetworkProfile: &armlabservices.LabPlanNetworkProfile{
	// 			SubnetID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default"),
	// 		},
	// 		SharedGalleryID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Compute/galleries/testsig"),
	// 		SupportInfo: &armlabservices.SupportInfo{
	// 			Email: to.Ptr("help@contoso.com"),
	// 			Instructions: to.Ptr("Contact support for help."),
	// 			Phone: to.Ptr("+1-202-555-0123"),
	// 			URL: to.Ptr("help.contoso.com"),
	// 		},
	// 		ProvisioningState: to.Ptr(armlabservices.ProvisioningStateSucceeded),
	// 	},
	// 	SystemData: &armlabservices.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T10:00:00.000Z"); return t}()),
	// 		CreatedBy: to.Ptr("identity123"),
	// 		CreatedByType: to.Ptr(armlabservices.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-06-01T09:12:28.000Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("identity123"),
	// 		LastModifiedByType: to.Ptr(armlabservices.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/LabPlans/putLabPlan.json
func ExampleLabPlansClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewLabPlansClient().BeginCreateOrUpdate(ctx, "testrg123", "testlabplan", armlabservices.LabPlan{
		Location: to.Ptr("westus"),
		Properties: &armlabservices.LabPlanProperties{
			DefaultAutoShutdownProfile: &armlabservices.AutoShutdownProfile{
				DisconnectDelay:          to.Ptr("PT5M"),
				IdleDelay:                to.Ptr("PT5M"),
				NoConnectDelay:           to.Ptr("PT5M"),
				ShutdownOnDisconnect:     to.Ptr(armlabservices.EnableStateEnabled),
				ShutdownOnIdle:           to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
				ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
			},
			DefaultConnectionProfile: &armlabservices.ConnectionProfile{
				ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				WebRdpAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
				WebSSHAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
			},
			DefaultNetworkProfile: &armlabservices.LabPlanNetworkProfile{
				SubnetID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default"),
			},
			SharedGalleryID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Compute/galleries/testsig"),
			SupportInfo: &armlabservices.SupportInfo{
				Email:        to.Ptr("help@contoso.com"),
				Instructions: to.Ptr("Contact support for help."),
				Phone:        to.Ptr("+1-202-555-0123"),
				URL:          to.Ptr("help.contoso.com"),
			},
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
	// res.LabPlan = armlabservices.LabPlan{
	// 	Name: to.Ptr("testlabplan"),
	// 	Type: to.Ptr("Microsoft.LabServices/LabPlan"),
	// 	ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labPlans/testlabplan"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armlabservices.LabPlanProperties{
	// 		DefaultAutoShutdownProfile: &armlabservices.AutoShutdownProfile{
	// 			DisconnectDelay: to.Ptr("PT5M"),
	// 			IdleDelay: to.Ptr("PT5M"),
	// 			NoConnectDelay: to.Ptr("PT5M"),
	// 			ShutdownOnDisconnect: to.Ptr(armlabservices.EnableStateEnabled),
	// 			ShutdownOnIdle: to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
	// 			ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
	// 		},
	// 		DefaultConnectionProfile: &armlabservices.ConnectionProfile{
	// 			ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
	// 			ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
	// 			WebRdpAccess: to.Ptr(armlabservices.ConnectionTypeNone),
	// 			WebSSHAccess: to.Ptr(armlabservices.ConnectionTypeNone),
	// 		},
	// 		DefaultNetworkProfile: &armlabservices.LabPlanNetworkProfile{
	// 			SubnetID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default"),
	// 		},
	// 		SharedGalleryID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Compute/galleries/testsig"),
	// 		SupportInfo: &armlabservices.SupportInfo{
	// 			Email: to.Ptr("help@contoso.com"),
	// 			Instructions: to.Ptr("Contact support for help."),
	// 			Phone: to.Ptr("+1-202-555-0123"),
	// 			URL: to.Ptr("help.contoso.com"),
	// 		},
	// 		ProvisioningState: to.Ptr(armlabservices.ProvisioningStateSucceeded),
	// 	},
	// 	SystemData: &armlabservices.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T10:00:00.000Z"); return t}()),
	// 		CreatedBy: to.Ptr("identity123"),
	// 		CreatedByType: to.Ptr(armlabservices.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-06-01T09:12:28.000Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("identity123"),
	// 		LastModifiedByType: to.Ptr(armlabservices.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/LabPlans/patchLabPlan.json
func ExampleLabPlansClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewLabPlansClient().BeginUpdate(ctx, "testrg123", "testlabplan", armlabservices.LabPlanUpdate{
		Properties: &armlabservices.LabPlanUpdateProperties{
			DefaultConnectionProfile: &armlabservices.ConnectionProfile{
				ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
				WebRdpAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
				WebSSHAccess:    to.Ptr(armlabservices.ConnectionTypeNone),
			},
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
	// res.LabPlan = armlabservices.LabPlan{
	// 	Name: to.Ptr("testlabplan"),
	// 	Type: to.Ptr("Microsoft.LabServices/LabPlan"),
	// 	ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labPlans/testlabplan"),
	// 	Location: to.Ptr("westus"),
	// 	Properties: &armlabservices.LabPlanProperties{
	// 		DefaultAutoShutdownProfile: &armlabservices.AutoShutdownProfile{
	// 			DisconnectDelay: to.Ptr("PT5M"),
	// 			IdleDelay: to.Ptr("PT5M"),
	// 			NoConnectDelay: to.Ptr("PT5M"),
	// 			ShutdownOnDisconnect: to.Ptr(armlabservices.EnableStateEnabled),
	// 			ShutdownOnIdle: to.Ptr(armlabservices.ShutdownOnIdleModeUserAbsence),
	// 			ShutdownWhenNotConnected: to.Ptr(armlabservices.EnableStateEnabled),
	// 		},
	// 		DefaultConnectionProfile: &armlabservices.ConnectionProfile{
	// 			ClientRdpAccess: to.Ptr(armlabservices.ConnectionTypePublic),
	// 			ClientSSHAccess: to.Ptr(armlabservices.ConnectionTypePublic),
	// 			WebRdpAccess: to.Ptr(armlabservices.ConnectionTypeNone),
	// 			WebSSHAccess: to.Ptr(armlabservices.ConnectionTypeNone),
	// 		},
	// 		DefaultNetworkProfile: &armlabservices.LabPlanNetworkProfile{
	// 			SubnetID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/default"),
	// 		},
	// 		SharedGalleryID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.Compute/galleries/testsig"),
	// 		SupportInfo: &armlabservices.SupportInfo{
	// 			Email: to.Ptr("help@contoso.com"),
	// 			Instructions: to.Ptr("Contact support for help."),
	// 			Phone: to.Ptr("+1-202-555-0123"),
	// 			URL: to.Ptr("help.contoso.com"),
	// 		},
	// 		ProvisioningState: to.Ptr(armlabservices.ProvisioningStateSucceeded),
	// 	},
	// 	SystemData: &armlabservices.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T10:00:00.000Z"); return t}()),
	// 		CreatedBy: to.Ptr("identity123"),
	// 		CreatedByType: to.Ptr(armlabservices.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-06-01T09:12:28.000Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("identity123"),
	// 		LastModifiedByType: to.Ptr(armlabservices.CreatedByTypeUser),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/LabPlans/deleteLabPlan.json
func ExampleLabPlansClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewLabPlansClient().BeginDelete(ctx, "testrg123", "testlabplan", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/LabPlans/saveImageVirtualMachine.json
func ExampleLabPlansClient_BeginSaveImage() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewLabPlansClient().BeginSaveImage(ctx, "testrg123", "testlabplan", armlabservices.SaveImageBody{
		Name:                to.Ptr("Test Image"),
		LabVirtualMachineID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labs/testlab/virtualMachines/template"),
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
