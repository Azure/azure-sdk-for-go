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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/NetworkWatcherFlowLogCreate.json
func ExampleFlowLogsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFlowLogsClient().BeginCreateOrUpdate(ctx, "rg1", "nw1", "fl", armnetwork.FlowLog{
		Location: to.Ptr("centraluseuap"),
		Identity: &armnetwork.ManagedServiceIdentity{
			Type: to.Ptr(armnetwork.ResourceIdentityTypeUserAssigned),
			UserAssignedIdentities: map[string]*armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
				"/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": {},
			},
		},
		Properties: &armnetwork.FlowLogPropertiesFormat{
			Format: &armnetwork.FlowLogFormatParameters{
				Type:    to.Ptr(armnetwork.FlowLogFormatTypeJSON),
				Version: to.Ptr[int32](1),
			},
			Enabled:                  to.Ptr(true),
			EnabledFilteringCriteria: to.Ptr("srcIP=158.255.7.8 || dstPort=56891"),
			StorageID:                to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/nwtest1mgvbfmqsigdxe"),
			TargetResourceID:         to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/desmondcentral-nsg"),
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
	// res.FlowLog = armnetwork.FlowLog{
	// 	Name: to.Ptr("Microsoft.Networkdesmond-rgdesmondcentral-nsg"),
	// 	Type: to.Ptr("Microsoft.Network/networkWatchers/FlowLogs"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkWatchers/nw/FlowLogs/fl"),
	// 	Location: to.Ptr("centraluseuap"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Identity: &armnetwork.ManagedServiceIdentity{
	// 		Type: to.Ptr(armnetwork.ResourceIdentityTypeUserAssigned),
	// 		UserAssignedIdentities: map[string]*armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
	// 			"/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": &armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
	// 				ClientID: to.Ptr("c16d15e1-f60a-40e4-8a05-df3d3f655c14"),
	// 				PrincipalID: to.Ptr("e3858881-e40c-43bd-9cde-88da39c05023"),
	// 			},
	// 		},
	// 	},
	// 	Properties: &armnetwork.FlowLogPropertiesFormat{
	// 		Format: &armnetwork.FlowLogFormatParameters{
	// 			Type: to.Ptr(armnetwork.FlowLogFormatTypeJSON),
	// 			Version: to.Ptr[int32](1),
	// 		},
	// 		Enabled: to.Ptr(true),
	// 		EnabledFilteringCriteria: to.Ptr("srcIP=158.255.7.8 || dstPort=56891"),
	// 		FlowAnalyticsConfiguration: &armnetwork.TrafficAnalyticsProperties{
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		RetentionPolicy: &armnetwork.RetentionPolicyParameters{
	// 			Days: to.Ptr[int32](0),
	// 			Enabled: to.Ptr(false),
	// 		},
	// 		StorageID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/nwtest1mgvbfmqsigdxe"),
	// 		TargetResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		TargetResourceID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/desmondcentral-nsg"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/NetworkWatcherFlowLogUpdateTags.json
func ExampleFlowLogsClient_UpdateTags() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFlowLogsClient().UpdateTags(ctx, "rg1", "nw", "fl", armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.FlowLog = armnetwork.FlowLog{
	// 	Name: to.Ptr("Microsoft.Networkdesmond-rgdesmondcentral-nsg"),
	// 	Type: to.Ptr("Microsoft.Network/networkWatchers/FlowLogs"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkWatchers/nw/FlowLogs/fl"),
	// 	Location: to.Ptr("centralus"),
	// 	Tags: map[string]*string{
	// 		"tag1": to.Ptr("value1"),
	// 		"tag2": to.Ptr("value2"),
	// 	},
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Identity: &armnetwork.ManagedServiceIdentity{
	// 		Type: to.Ptr(armnetwork.ResourceIdentityTypeUserAssigned),
	// 		UserAssignedIdentities: map[string]*armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
	// 			"/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": &armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
	// 				ClientID: to.Ptr("c16d15e1-f60a-40e4-8a05-df3d3f655c14"),
	// 				PrincipalID: to.Ptr("e3858881-e40c-43bd-9cde-88da39c05023"),
	// 			},
	// 		},
	// 	},
	// 	Properties: &armnetwork.FlowLogPropertiesFormat{
	// 		Format: &armnetwork.FlowLogFormatParameters{
	// 			Type: to.Ptr(armnetwork.FlowLogFormatTypeJSON),
	// 			Version: to.Ptr[int32](1),
	// 		},
	// 		Enabled: to.Ptr(true),
	// 		EnabledFilteringCriteria: to.Ptr("srcIP=158.255.7.8 || dstPort=56891"),
	// 		FlowAnalyticsConfiguration: &armnetwork.TrafficAnalyticsProperties{
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		RetentionPolicy: &armnetwork.RetentionPolicyParameters{
	// 			Days: to.Ptr[int32](0),
	// 			Enabled: to.Ptr(false),
	// 		},
	// 		StorageID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Storage/storageAccounts/nwtest1mgvbfmqsigdxe"),
	// 		TargetResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		TargetResourceID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/desmondcentral-nsg"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/NetworkWatcherFlowLogGet.json
func ExampleFlowLogsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFlowLogsClient().Get(ctx, "rg1", "nw1", "flowLog1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.FlowLog = armnetwork.FlowLog{
	// 	Name: to.Ptr("flowLog1"),
	// 	Type: to.Ptr("Microsoft.Network/networkWatchers/FlowLogs"),
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkWatchers/тц1/FlowLogs/flowLog1"),
	// 	Location: to.Ptr("centraluseuap"),
	// 	Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
	// 	Identity: &armnetwork.ManagedServiceIdentity{
	// 		Type: to.Ptr(armnetwork.ResourceIdentityTypeUserAssigned),
	// 		UserAssignedIdentities: map[string]*armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
	// 			"/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": &armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
	// 				ClientID: to.Ptr("c16d15e1-f60a-40e4-8a05-df3d3f655c14"),
	// 				PrincipalID: to.Ptr("e3858881-e40c-43bd-9cde-88da39c05023"),
	// 			},
	// 		},
	// 	},
	// 	Properties: &armnetwork.FlowLogPropertiesFormat{
	// 		Format: &armnetwork.FlowLogFormatParameters{
	// 			Type: to.Ptr(armnetwork.FlowLogFormatTypeJSON),
	// 			Version: to.Ptr[int32](2),
	// 		},
	// 		Enabled: to.Ptr(true),
	// 		EnabledFilteringCriteria: to.Ptr("srcIP=158.255.7.8 || dstPort=56891"),
	// 		FlowAnalyticsConfiguration: &armnetwork.TrafficAnalyticsProperties{
	// 			NetworkWatcherFlowAnalyticsConfiguration: &armnetwork.TrafficAnalyticsConfigurationProperties{
	// 				Enabled: to.Ptr(false),
	// 				TrafficAnalyticsInterval: to.Ptr[int32](60),
	// 				WorkspaceID: to.Ptr("-"),
	// 				WorkspaceRegion: to.Ptr("-"),
	// 			},
	// 		},
	// 		ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
	// 		RetentionPolicy: &armnetwork.RetentionPolicyParameters{
	// 			Days: to.Ptr[int32](0),
	// 			Enabled: to.Ptr(false),
	// 		},
	// 		StorageID: to.Ptr("/subscriptions/subid/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/wzstorage002"),
	// 		TargetResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 		TargetResourceID: to.Ptr("/subscriptions/subid/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/vm5-nsg"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/NetworkWatcherFlowLogDelete.json
func ExampleFlowLogsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFlowLogsClient().BeginDelete(ctx, "rg1", "nw1", "fl", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/b43042075540b8d67cce7d3d9f70b9b9f5a359da/specification/network/resource-manager/Microsoft.Network/stable/2024-05-01/examples/NetworkWatcherFlowLogList.json
func ExampleFlowLogsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewFlowLogsClient().NewListPager("rg1", "nw1", nil)
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
		// page.FlowLogListResult = armnetwork.FlowLogListResult{
		// 	Value: []*armnetwork.FlowLog{
		// 		{
		// 			Name: to.Ptr("flowLog1"),
		// 			Type: to.Ptr("Microsoft.Network/networkWatchers/FlowLogs"),
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/networkWatchers/тц1/FlowLogs/flowLog1"),
		// 			Location: to.Ptr("centraluseuap"),
		// 			Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 			Identity: &armnetwork.ManagedServiceIdentity{
		// 				Type: to.Ptr(armnetwork.ResourceIdentityTypeUserAssigned),
		// 				UserAssignedIdentities: map[string]*armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
		// 					"/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": &armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
		// 						ClientID: to.Ptr("c16d15e1-f60a-40e4-8a05-df3d3f655c14"),
		// 						PrincipalID: to.Ptr("e3858881-e40c-43bd-9cde-88da39c05023"),
		// 					},
		// 				},
		// 			},
		// 			Properties: &armnetwork.FlowLogPropertiesFormat{
		// 				Format: &armnetwork.FlowLogFormatParameters{
		// 					Type: to.Ptr(armnetwork.FlowLogFormatTypeJSON),
		// 					Version: to.Ptr[int32](2),
		// 				},
		// 				Enabled: to.Ptr(true),
		// 				EnabledFilteringCriteria: to.Ptr("srcIP=158.255.7.8 || dstPort=56891"),
		// 				FlowAnalyticsConfiguration: &armnetwork.TrafficAnalyticsProperties{
		// 					NetworkWatcherFlowAnalyticsConfiguration: &armnetwork.TrafficAnalyticsConfigurationProperties{
		// 						Enabled: to.Ptr(false),
		// 						TrafficAnalyticsInterval: to.Ptr[int32](60),
		// 						WorkspaceID: to.Ptr("-"),
		// 						WorkspaceRegion: to.Ptr("-"),
		// 					},
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				RetentionPolicy: &armnetwork.RetentionPolicyParameters{
		// 					Days: to.Ptr[int32](0),
		// 					Enabled: to.Ptr(false),
		// 				},
		// 				StorageID: to.Ptr("/subscriptions/subid/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/wzstorage002"),
		// 				TargetResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				TargetResourceID: to.Ptr("/subscriptions/subid/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/vm5-nsg"),
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("flowLog2"),
		// 			Type: to.Ptr("Microsoft.Network/networkWatchers/FlowLogs"),
		// 			ID: to.Ptr("/subscriptions/96e68903-0a56-4819-9987-8d08ad6a1f99/resourceGroups/NetworkWatcherRG/providers/Microsoft.Network/networkWatchers/NetworkWatcher_centraluseuap/FlowLogs/flowLog2"),
		// 			Location: to.Ptr("centraluseuap"),
		// 			Etag: to.Ptr("W/\"00000000-0000-0000-0000-000000000000\""),
		// 			Identity: &armnetwork.ManagedServiceIdentity{
		// 				Type: to.Ptr(armnetwork.ResourceIdentityTypeUserAssigned),
		// 				UserAssignedIdentities: map[string]*armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
		// 					"/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": &armnetwork.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
		// 						ClientID: to.Ptr("c16d15e1-f60a-40e4-8a05-df3d3f655c14"),
		// 						PrincipalID: to.Ptr("e3858881-e40c-43bd-9cde-88da39c05023"),
		// 					},
		// 				},
		// 			},
		// 			Properties: &armnetwork.FlowLogPropertiesFormat{
		// 				Format: &armnetwork.FlowLogFormatParameters{
		// 					Type: to.Ptr(armnetwork.FlowLogFormatTypeJSON),
		// 					Version: to.Ptr[int32](2),
		// 				},
		// 				Enabled: to.Ptr(true),
		// 				EnabledFilteringCriteria: to.Ptr("srcIP=158.255.7.8 || dstPort=56891"),
		// 				FlowAnalyticsConfiguration: &armnetwork.TrafficAnalyticsProperties{
		// 				},
		// 				ProvisioningState: to.Ptr(armnetwork.ProvisioningStateSucceeded),
		// 				RetentionPolicy: &armnetwork.RetentionPolicyParameters{
		// 					Days: to.Ptr[int32](0),
		// 					Enabled: to.Ptr(false),
		// 				},
		// 				StorageID: to.Ptr("/subscriptions/subid/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/iraflowlogtest2diag"),
		// 				TargetResourceGUID: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				TargetResourceID: to.Ptr("/subscriptions/subid/resourceGroups/rg/providers/Microsoft.Network/networkSecurityGroups/DSCP-test-vm1-nsg"),
		// 			},
		// 	}},
		// }
	}
}
