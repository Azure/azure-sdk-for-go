//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armstoragemover_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storagemover/armstoragemover/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_List_MaximumSet.json
func ExampleAgentsClient_NewListPager_agentsListMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAgentsClient().NewListPager("examples-rg", "examples-storageMoverName", nil)
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
		// page.AgentList = armstoragemover.AgentList{
		// 	Value: []*armstoragemover.Agent{
		// 		{
		// 			Name: to.Ptr("examples-agentName1"),
		// 			Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
		// 			ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName1"),
		// 			SystemData: &armstoragemover.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
		// 				CreatedBy: to.Ptr("string"),
		// 				CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
		// 				LastModifiedBy: to.Ptr("string"),
		// 				LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
		// 			},
		// 			Properties: &armstoragemover.AgentProperties{
		// 				Description: to.Ptr("Example Agent 1 Description"),
		// 				AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
		// 				AgentVersion: to.Ptr("1.0.0"),
		// 				ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName1"),
		// 				ArcVMUUID: to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
		// 				LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
		// 				LocalIPAddress: to.Ptr("192.168.0.0"),
		// 				MemoryInMB: to.Ptr[int64](4096),
		// 				NumberOfCores: to.Ptr[int64](8),
		// 				ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
		// 				TimeZone: to.Ptr("Eastern Standard Time"),
		// 				UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
		// 					WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
		// 						{
		// 							LimitInMbps: to.Ptr[int32](2000),
		// 							EndTime: &armstoragemover.Time{
		// 								Hour: to.Ptr[int32](18),
		// 								Minute: to.Ptr(armstoragemover.Minute(30)),
		// 							},
		// 							StartTime: &armstoragemover.Time{
		// 								Hour: to.Ptr[int32](9),
		// 								Minute: to.Ptr(armstoragemover.Minute(0)),
		// 							},
		// 							Days: []*armstoragemover.DayOfWeek{
		// 								to.Ptr(armstoragemover.DayOfWeekMonday)},
		// 						}},
		// 					},
		// 					UptimeInSeconds: to.Ptr[int64](522),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("examples-agentName2"),
		// 				Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
		// 				ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName2"),
		// 				SystemData: &armstoragemover.SystemData{
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
		// 					CreatedBy: to.Ptr("string"),
		// 					CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("string"),
		// 					LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
		// 				},
		// 				Properties: &armstoragemover.AgentProperties{
		// 					AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
		// 					AgentVersion: to.Ptr("1.0.0"),
		// 					ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName2"),
		// 					ArcVMUUID: to.Ptr("147a1f84-7bbf-4e99-9a6a-a1735a91dfd5"),
		// 					LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
		// 					LocalIPAddress: to.Ptr("192.168.0.0"),
		// 					MemoryInMB: to.Ptr[int64](4096),
		// 					NumberOfCores: to.Ptr[int64](8),
		// 					ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
		// 					TimeZone: to.Ptr("Eastern Standard Time"),
		// 					UptimeInSeconds: to.Ptr[int64](877),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("examples-agentName3"),
		// 				Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
		// 				ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName3"),
		// 				SystemData: &armstoragemover.SystemData{
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
		// 					CreatedBy: to.Ptr("string"),
		// 					CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("string"),
		// 					LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
		// 				},
		// 				Properties: &armstoragemover.AgentProperties{
		// 					AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
		// 					AgentVersion: to.Ptr("1.0.0"),
		// 					ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName3"),
		// 					ArcVMUUID: to.Ptr("648a7958-f99e-4268-b20e-94c96558dc0d"),
		// 					ErrorDetails: &armstoragemover.AgentPropertiesErrorDetails{
		// 						Code: to.Ptr("SampleErrorCode"),
		// 						Message: to.Ptr("Detailed sample error message."),
		// 					},
		// 					LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
		// 					LocalIPAddress: to.Ptr("192.168.0.0"),
		// 					MemoryInMB: to.Ptr[int64](100024),
		// 					NumberOfCores: to.Ptr[int64](32),
		// 					ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
		// 					TimeZone: to.Ptr("Eastern Standard Time"),
		// 					UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
		// 						WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
		// 							{
		// 								LimitInMbps: to.Ptr[int32](5000),
		// 								EndTime: &armstoragemover.Time{
		// 									Hour: to.Ptr[int32](24),
		// 									Minute: to.Ptr(armstoragemover.Minute(0)),
		// 								},
		// 								StartTime: &armstoragemover.Time{
		// 									Hour: to.Ptr[int32](0),
		// 									Minute: to.Ptr(armstoragemover.Minute(0)),
		// 								},
		// 								Days: []*armstoragemover.DayOfWeek{
		// 									to.Ptr(armstoragemover.DayOfWeekSaturday),
		// 									to.Ptr(armstoragemover.DayOfWeekSunday)},
		// 							}},
		// 						},
		// 						UptimeInSeconds: to.Ptr[int64](1025),
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_List_MinimumSet.json
func ExampleAgentsClient_NewListPager_agentsListMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAgentsClient().NewListPager("examples-rg", "examples-storageMoverName", nil)
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
		// page.AgentList = armstoragemover.AgentList{
		// 	Value: []*armstoragemover.Agent{
		// 	},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_Get_MaximumSet.json
func ExampleAgentsClient_Get_agentsGetMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAgentsClient().Get(ctx, "examples-rg", "examples-storageMoverName", "examples-agentName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Agent = armstoragemover.Agent{
	// 	Name: to.Ptr("examples-agentName"),
	// 	Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
	// 	ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName"),
	// 	SystemData: &armstoragemover.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 	},
	// 	Properties: &armstoragemover.AgentProperties{
	// 		Description: to.Ptr("Example Agent Description"),
	// 		AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
	// 		AgentVersion: to.Ptr("1.0.0"),
	// 		ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
	// 		ArcVMUUID: to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
	// 		LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
	// 		LocalIPAddress: to.Ptr("192.168.0.0"),
	// 		MemoryInMB: to.Ptr[int64](4096),
	// 		NumberOfCores: to.Ptr[int64](8),
	// 		ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
	// 		TimeZone: to.Ptr("Eastern Standard Time"),
	// 		UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
	// 			WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
	// 				{
	// 					LimitInMbps: to.Ptr[int32](2000),
	// 					EndTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](18),
	// 						Minute: to.Ptr(armstoragemover.Minute(30)),
	// 					},
	// 					StartTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](9),
	// 						Minute: to.Ptr(armstoragemover.Minute(0)),
	// 					},
	// 					Days: []*armstoragemover.DayOfWeek{
	// 						to.Ptr(armstoragemover.DayOfWeekMonday)},
	// 				}},
	// 			},
	// 			UptimeInSeconds: to.Ptr[int64](522),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_Get_MinimumSet.json
func ExampleAgentsClient_Get_agentsGetMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAgentsClient().Get(ctx, "examples-rg", "examples-storageMoverName", "examples-agentName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Agent = armstoragemover.Agent{
	// 	Name: to.Ptr("examples-agentName"),
	// 	Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
	// 	ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName"),
	// 	SystemData: &armstoragemover.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 	},
	// 	Properties: &armstoragemover.AgentProperties{
	// 		AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
	// 		AgentVersion: to.Ptr("1.0.0"),
	// 		ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
	// 		ArcVMUUID: to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
	// 		LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
	// 		LocalIPAddress: to.Ptr("192.168.0.0"),
	// 		MemoryInMB: to.Ptr[int64](4096),
	// 		NumberOfCores: to.Ptr[int64](8),
	// 		ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
	// 		TimeZone: to.Ptr("Eastern Standard Time"),
	// 		UptimeInSeconds: to.Ptr[int64](522),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_CreateOrUpdate_MaximumSet.json
func ExampleAgentsClient_CreateOrUpdate_agentsCreateOrUpdateMaximumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAgentsClient().CreateOrUpdate(ctx, "examples-rg", "examples-storageMoverName", "examples-agentName", armstoragemover.Agent{
		Properties: &armstoragemover.AgentProperties{
			Description:   to.Ptr("Example Agent Description"),
			ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
			ArcVMUUID:     to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
			UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
				WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
					{
						LimitInMbps: to.Ptr[int32](2000),
						EndTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](18),
							Minute: to.Ptr(armstoragemover.Minute(30)),
						},
						StartTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](9),
							Minute: to.Ptr(armstoragemover.Minute(0)),
						},
						Days: []*armstoragemover.DayOfWeek{
							to.Ptr(armstoragemover.DayOfWeekMonday)},
					}},
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Agent = armstoragemover.Agent{
	// 	Name: to.Ptr("examples-agentName"),
	// 	Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
	// 	ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName"),
	// 	SystemData: &armstoragemover.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 	},
	// 	Properties: &armstoragemover.AgentProperties{
	// 		Description: to.Ptr("Example Agent Description"),
	// 		AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
	// 		AgentVersion: to.Ptr("1.0.0"),
	// 		ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
	// 		ArcVMUUID: to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
	// 		LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
	// 		LocalIPAddress: to.Ptr("192.168.0.0"),
	// 		MemoryInMB: to.Ptr[int64](4096),
	// 		NumberOfCores: to.Ptr[int64](8),
	// 		ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
	// 		TimeZone: to.Ptr("Eastern Standard Time"),
	// 		UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
	// 			WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
	// 				{
	// 					LimitInMbps: to.Ptr[int32](2000),
	// 					EndTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](18),
	// 						Minute: to.Ptr(armstoragemover.Minute(30)),
	// 					},
	// 					StartTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](9),
	// 						Minute: to.Ptr(armstoragemover.Minute(0)),
	// 					},
	// 					Days: []*armstoragemover.DayOfWeek{
	// 						to.Ptr(armstoragemover.DayOfWeekMonday)},
	// 				}},
	// 			},
	// 			UptimeInSeconds: to.Ptr[int64](522),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_CreateOrUpdate_MinimumSet.json
func ExampleAgentsClient_CreateOrUpdate_agentsCreateOrUpdateMinimumSet() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAgentsClient().CreateOrUpdate(ctx, "examples-rg", "examples-storageMoverName", "examples-agentName", armstoragemover.Agent{
		Properties: &armstoragemover.AgentProperties{
			ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
			ArcVMUUID:     to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Agent = armstoragemover.Agent{
	// 	Name: to.Ptr("examples-agentName"),
	// 	Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
	// 	ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName"),
	// 	SystemData: &armstoragemover.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 	},
	// 	Properties: &armstoragemover.AgentProperties{
	// 		AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
	// 		AgentVersion: to.Ptr("1.0.0"),
	// 		ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
	// 		ArcVMUUID: to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
	// 		LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
	// 		LocalIPAddress: to.Ptr("192.168.0.0"),
	// 		MemoryInMB: to.Ptr[int64](4096),
	// 		NumberOfCores: to.Ptr[int64](8),
	// 		ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
	// 		TimeZone: to.Ptr("Eastern Standard Time"),
	// 		UptimeInSeconds: to.Ptr[int64](522),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_CreateOrUpdate_UploadLimitSchedule_Overnight.json
func ExampleAgentsClient_CreateOrUpdate_agentsCreateOrUpdateWithOvernightUploadLimitSchedule() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAgentsClient().CreateOrUpdate(ctx, "examples-rg", "examples-storageMoverName", "examples-agentName", armstoragemover.Agent{
		Properties: &armstoragemover.AgentProperties{
			ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
			ArcVMUUID:     to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
			UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
				WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
					{
						LimitInMbps: to.Ptr[int32](2000),
						EndTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](24),
							Minute: to.Ptr(armstoragemover.Minute(0)),
						},
						StartTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](18),
							Minute: to.Ptr(armstoragemover.Minute(0)),
						},
						Days: []*armstoragemover.DayOfWeek{
							to.Ptr(armstoragemover.DayOfWeekMonday),
							to.Ptr(armstoragemover.DayOfWeekTuesday),
							to.Ptr(armstoragemover.DayOfWeekWednesday),
							to.Ptr(armstoragemover.DayOfWeekThursday),
							to.Ptr(armstoragemover.DayOfWeekFriday),
							to.Ptr(armstoragemover.DayOfWeekSaturday),
							to.Ptr(armstoragemover.DayOfWeekSunday)},
					},
					{
						LimitInMbps: to.Ptr[int32](2000),
						EndTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](9),
							Minute: to.Ptr(armstoragemover.Minute(0)),
						},
						StartTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](0),
							Minute: to.Ptr(armstoragemover.Minute(0)),
						},
						Days: []*armstoragemover.DayOfWeek{
							to.Ptr(armstoragemover.DayOfWeekMonday),
							to.Ptr(armstoragemover.DayOfWeekTuesday),
							to.Ptr(armstoragemover.DayOfWeekWednesday),
							to.Ptr(armstoragemover.DayOfWeekThursday),
							to.Ptr(armstoragemover.DayOfWeekFriday),
							to.Ptr(armstoragemover.DayOfWeekSaturday),
							to.Ptr(armstoragemover.DayOfWeekSunday)},
					}},
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Agent = armstoragemover.Agent{
	// 	Name: to.Ptr("examples-agentName"),
	// 	Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
	// 	ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName"),
	// 	SystemData: &armstoragemover.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 	},
	// 	Properties: &armstoragemover.AgentProperties{
	// 		AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
	// 		AgentVersion: to.Ptr("1.0.0"),
	// 		ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
	// 		ArcVMUUID: to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
	// 		LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
	// 		LocalIPAddress: to.Ptr("192.168.0.0"),
	// 		MemoryInMB: to.Ptr[int64](4096),
	// 		NumberOfCores: to.Ptr[int64](8),
	// 		ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
	// 		TimeZone: to.Ptr("Eastern Standard Time"),
	// 		UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
	// 			WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
	// 				{
	// 					LimitInMbps: to.Ptr[int32](2000),
	// 					EndTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](24),
	// 						Minute: to.Ptr(armstoragemover.Minute(0)),
	// 					},
	// 					StartTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](18),
	// 						Minute: to.Ptr(armstoragemover.Minute(0)),
	// 					},
	// 					Days: []*armstoragemover.DayOfWeek{
	// 						to.Ptr(armstoragemover.DayOfWeekMonday),
	// 						to.Ptr(armstoragemover.DayOfWeekTuesday),
	// 						to.Ptr(armstoragemover.DayOfWeekWednesday),
	// 						to.Ptr(armstoragemover.DayOfWeekThursday),
	// 						to.Ptr(armstoragemover.DayOfWeekFriday),
	// 						to.Ptr(armstoragemover.DayOfWeekSaturday),
	// 						to.Ptr(armstoragemover.DayOfWeekSunday)},
	// 					},
	// 					{
	// 						LimitInMbps: to.Ptr[int32](2000),
	// 						EndTime: &armstoragemover.Time{
	// 							Hour: to.Ptr[int32](9),
	// 							Minute: to.Ptr(armstoragemover.Minute(0)),
	// 						},
	// 						StartTime: &armstoragemover.Time{
	// 							Hour: to.Ptr[int32](0),
	// 							Minute: to.Ptr(armstoragemover.Minute(0)),
	// 						},
	// 						Days: []*armstoragemover.DayOfWeek{
	// 							to.Ptr(armstoragemover.DayOfWeekMonday),
	// 							to.Ptr(armstoragemover.DayOfWeekTuesday),
	// 							to.Ptr(armstoragemover.DayOfWeekWednesday),
	// 							to.Ptr(armstoragemover.DayOfWeekThursday),
	// 							to.Ptr(armstoragemover.DayOfWeekFriday),
	// 							to.Ptr(armstoragemover.DayOfWeekSaturday),
	// 							to.Ptr(armstoragemover.DayOfWeekSunday)},
	// 					}},
	// 				},
	// 				UptimeInSeconds: to.Ptr[int64](522),
	// 			},
	// 		}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_Update.json
func ExampleAgentsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAgentsClient().Update(ctx, "examples-rg", "examples-storageMoverName", "examples-agentName", armstoragemover.AgentUpdateParameters{
		Properties: &armstoragemover.AgentUpdateProperties{
			Description: to.Ptr("Example Agent Description"),
			UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
				WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
					{
						LimitInMbps: to.Ptr[int32](2000),
						EndTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](18),
							Minute: to.Ptr(armstoragemover.Minute(30)),
						},
						StartTime: &armstoragemover.Time{
							Hour:   to.Ptr[int32](9),
							Minute: to.Ptr(armstoragemover.Minute(0)),
						},
						Days: []*armstoragemover.DayOfWeek{
							to.Ptr(armstoragemover.DayOfWeekMonday)},
					}},
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Agent = armstoragemover.Agent{
	// 	Name: to.Ptr("examples-agentName"),
	// 	Type: to.Ptr("Microsoft.StorageMover/storageMovers/agents"),
	// 	ID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.StorageMover/storageMovers/examples-storageMoverName/agents/examples-agentName"),
	// 	SystemData: &armstoragemover.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T01:01:01.107Z"); return t}()),
	// 		CreatedBy: to.Ptr("string"),
	// 		CreatedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:01:01.107Z"); return t}()),
	// 		LastModifiedBy: to.Ptr("string"),
	// 		LastModifiedByType: to.Ptr(armstoragemover.CreatedByTypeUser),
	// 	},
	// 	Properties: &armstoragemover.AgentProperties{
	// 		Description: to.Ptr("Example Agent Description"),
	// 		AgentStatus: to.Ptr(armstoragemover.AgentStatusOnline),
	// 		AgentVersion: to.Ptr("1.0.0"),
	// 		ArcResourceID: to.Ptr("/subscriptions/60bcfc77-6589-4da2-b7fd-f9ec9322cf95/resourceGroups/examples-rg/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"),
	// 		ArcVMUUID: to.Ptr("3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"),
	// 		LastStatusUpdate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2023-07-01T02:21:01.107Z"); return t}()),
	// 		LocalIPAddress: to.Ptr("192.168.0.0"),
	// 		MemoryInMB: to.Ptr[int64](4096),
	// 		NumberOfCores: to.Ptr[int64](8),
	// 		ProvisioningState: to.Ptr(armstoragemover.ProvisioningStateSucceeded),
	// 		TimeZone: to.Ptr("Eastern Standard Time"),
	// 		UploadLimitSchedule: &armstoragemover.UploadLimitSchedule{
	// 			WeeklyRecurrences: []*armstoragemover.UploadLimitWeeklyRecurrence{
	// 				{
	// 					LimitInMbps: to.Ptr[int32](2000),
	// 					EndTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](18),
	// 						Minute: to.Ptr(armstoragemover.Minute(30)),
	// 					},
	// 					StartTime: &armstoragemover.Time{
	// 						Hour: to.Ptr[int32](9),
	// 						Minute: to.Ptr(armstoragemover.Minute(0)),
	// 					},
	// 					Days: []*armstoragemover.DayOfWeek{
	// 						to.Ptr(armstoragemover.DayOfWeekMonday)},
	// 				}},
	// 			},
	// 			UptimeInSeconds: to.Ptr[int64](522),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4ee6d9fd7687d4b67117c5a167c191a7e7e70b53/specification/storagemover/resource-manager/Microsoft.StorageMover/stable/2024-07-01/examples/Agents_Delete.json
func ExampleAgentsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armstoragemover.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAgentsClient().BeginDelete(ctx, "examples-rg", "examples-storageMoverName", "examples-agentName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
