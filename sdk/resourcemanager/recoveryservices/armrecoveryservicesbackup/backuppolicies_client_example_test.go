//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armrecoveryservicesbackup_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/80c21c17b4a7aa57f637ee594f7cfd653255a7e0/specification/recoveryservicesbackup/resource-manager/Microsoft.RecoveryServices/stable/2023-04-01/examples/AzureIaasVm/BackupPolicies_List.json
func ExampleBackupPoliciesClient_NewListPager_listProtectionPoliciesWithBackupManagementTypeFilterAsAzureIaasVm() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armrecoveryservicesbackup.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewBackupPoliciesClient().NewListPager("NetSDKTestRsVault", "SwaggerTestRg", &armrecoveryservicesbackup.BackupPoliciesClientListOptions{Filter: to.Ptr("backupManagementType eq 'AzureIaasVM'")})
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
		// page.ProtectionPolicyResourceList = armrecoveryservicesbackup.ProtectionPolicyResourceList{
		// 	Value: []*armrecoveryservicesbackup.ProtectionPolicyResource{
		// 		{
		// 			Name: to.Ptr("DefaultPolicy"),
		// 			Type: to.Ptr("Microsoft.RecoveryServices/vaults/backupPolicies"),
		// 			ID: to.Ptr("/Subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SwaggerTestRg/providers/Microsoft.RecoveryServices/vaults/NetSDKTestRsVault/backupPolicies/DefaultPolicy"),
		// 			Properties: &armrecoveryservicesbackup.AzureIaaSVMProtectionPolicy{
		// 				BackupManagementType: to.Ptr("AzureIaasVM"),
		// 				ProtectedItemsCount: to.Ptr[int32](0),
		// 				RetentionPolicy: &armrecoveryservicesbackup.LongTermRetentionPolicy{
		// 					RetentionPolicyType: to.Ptr("LongTermRetentionPolicy"),
		// 					DailySchedule: &armrecoveryservicesbackup.DailyRetentionSchedule{
		// 						RetentionDuration: &armrecoveryservicesbackup.RetentionDuration{
		// 							Count: to.Ptr[int32](30),
		// 							DurationType: to.Ptr(armrecoveryservicesbackup.RetentionDurationTypeDays),
		// 						},
		// 						RetentionTimes: []*time.Time{
		// 							to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-05T19:00:00Z"); return t}())},
		// 						},
		// 					},
		// 					SchedulePolicy: &armrecoveryservicesbackup.SimpleSchedulePolicy{
		// 						SchedulePolicyType: to.Ptr("SimpleSchedulePolicy"),
		// 						ScheduleRunFrequency: to.Ptr(armrecoveryservicesbackup.ScheduleRunTypeDaily),
		// 						ScheduleRunTimes: []*time.Time{
		// 							to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-05T19:00:00Z"); return t}())},
		// 							ScheduleWeeklyFrequency: to.Ptr[int32](0),
		// 						},
		// 					},
		// 				},
		// 				{
		// 					Name: to.Ptr("testPolicy1"),
		// 					Type: to.Ptr("Microsoft.RecoveryServices/vaults/backupPolicies"),
		// 					ID: to.Ptr("/Subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SwaggerTestRg/providers/Microsoft.RecoveryServices/vaults/NetSDKTestRsVault/backupPolicies/testPolicy1"),
		// 					Properties: &armrecoveryservicesbackup.AzureIaaSVMProtectionPolicy{
		// 						BackupManagementType: to.Ptr("AzureIaasVM"),
		// 						ProtectedItemsCount: to.Ptr[int32](0),
		// 						RetentionPolicy: &armrecoveryservicesbackup.LongTermRetentionPolicy{
		// 							RetentionPolicyType: to.Ptr("LongTermRetentionPolicy"),
		// 							DailySchedule: &armrecoveryservicesbackup.DailyRetentionSchedule{
		// 								RetentionDuration: &armrecoveryservicesbackup.RetentionDuration{
		// 									Count: to.Ptr[int32](1),
		// 									DurationType: to.Ptr(armrecoveryservicesbackup.RetentionDurationTypeDays),
		// 								},
		// 								RetentionTimes: []*time.Time{
		// 									to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-01-24T02:00:00Z"); return t}())},
		// 								},
		// 							},
		// 							SchedulePolicy: &armrecoveryservicesbackup.SimpleSchedulePolicy{
		// 								SchedulePolicyType: to.Ptr("SimpleSchedulePolicy"),
		// 								ScheduleRunFrequency: to.Ptr(armrecoveryservicesbackup.ScheduleRunTypeDaily),
		// 								ScheduleRunTimes: []*time.Time{
		// 									to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-01-24T02:00:00Z"); return t}())},
		// 									ScheduleWeeklyFrequency: to.Ptr[int32](0),
		// 								},
		// 								TimeZone: to.Ptr("Pacific Standard Time"),
		// 							},
		// 					}},
		// 				}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/80c21c17b4a7aa57f637ee594f7cfd653255a7e0/specification/recoveryservicesbackup/resource-manager/Microsoft.RecoveryServices/stable/2023-04-01/examples/AzureIaasVm/V2Policy/v2-List-Policies.json
func ExampleBackupPoliciesClient_NewListPager_listProtectionPoliciesWithBackupManagementTypeFilterAsAzureIaasVmWithBothV1AndV2Policies() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armrecoveryservicesbackup.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewBackupPoliciesClient().NewListPager("NetSDKTestRsVault", "SwaggerTestRg", &armrecoveryservicesbackup.BackupPoliciesClientListOptions{Filter: to.Ptr("backupManagementType eq 'AzureIaasVM'")})
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
		// page.ProtectionPolicyResourceList = armrecoveryservicesbackup.ProtectionPolicyResourceList{
		// 	Value: []*armrecoveryservicesbackup.ProtectionPolicyResource{
		// 		{
		// 			Name: to.Ptr("DefaultPolicy"),
		// 			Type: to.Ptr("Microsoft.RecoveryServices/vaults/backupPolicies"),
		// 			ID: to.Ptr("/Subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SwaggerTestRg/providers/Microsoft.RecoveryServices/vaults/NetSDKTestRsVault/backupPolicies/DefaultPolicy"),
		// 			Properties: &armrecoveryservicesbackup.AzureIaaSVMProtectionPolicy{
		// 				BackupManagementType: to.Ptr("AzureIaasVM"),
		// 				ProtectedItemsCount: to.Ptr[int32](0),
		// 				RetentionPolicy: &armrecoveryservicesbackup.LongTermRetentionPolicy{
		// 					RetentionPolicyType: to.Ptr("LongTermRetentionPolicy"),
		// 					DailySchedule: &armrecoveryservicesbackup.DailyRetentionSchedule{
		// 						RetentionDuration: &armrecoveryservicesbackup.RetentionDuration{
		// 							Count: to.Ptr[int32](30),
		// 							DurationType: to.Ptr(armrecoveryservicesbackup.RetentionDurationTypeDays),
		// 						},
		// 						RetentionTimes: []*time.Time{
		// 							to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-05T19:00:00Z"); return t}())},
		// 						},
		// 					},
		// 					SchedulePolicy: &armrecoveryservicesbackup.SimpleSchedulePolicy{
		// 						SchedulePolicyType: to.Ptr("SimpleSchedulePolicy"),
		// 						ScheduleRunFrequency: to.Ptr(armrecoveryservicesbackup.ScheduleRunTypeDaily),
		// 						ScheduleRunTimes: []*time.Time{
		// 							to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-05T19:00:00Z"); return t}())},
		// 							ScheduleWeeklyFrequency: to.Ptr[int32](0),
		// 						},
		// 					},
		// 				},
		// 				{
		// 					Name: to.Ptr("testPolicy1"),
		// 					Type: to.Ptr("Microsoft.RecoveryServices/vaults/backupPolicies"),
		// 					ID: to.Ptr("/Subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SwaggerTestRg/providers/Microsoft.RecoveryServices/vaults/NetSDKTestRsVault/backupPolicies/testPolicy1"),
		// 					Properties: &armrecoveryservicesbackup.AzureIaaSVMProtectionPolicy{
		// 						BackupManagementType: to.Ptr("AzureIaasVM"),
		// 						ProtectedItemsCount: to.Ptr[int32](0),
		// 						RetentionPolicy: &armrecoveryservicesbackup.LongTermRetentionPolicy{
		// 							RetentionPolicyType: to.Ptr("LongTermRetentionPolicy"),
		// 							DailySchedule: &armrecoveryservicesbackup.DailyRetentionSchedule{
		// 								RetentionDuration: &armrecoveryservicesbackup.RetentionDuration{
		// 									Count: to.Ptr[int32](1),
		// 									DurationType: to.Ptr(armrecoveryservicesbackup.RetentionDurationTypeDays),
		// 								},
		// 								RetentionTimes: []*time.Time{
		// 									to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-01-24T02:00:00Z"); return t}())},
		// 								},
		// 							},
		// 							SchedulePolicy: &armrecoveryservicesbackup.SimpleSchedulePolicy{
		// 								SchedulePolicyType: to.Ptr("SimpleSchedulePolicy"),
		// 								ScheduleRunFrequency: to.Ptr(armrecoveryservicesbackup.ScheduleRunTypeDaily),
		// 								ScheduleRunTimes: []*time.Time{
		// 									to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-01-24T02:00:00Z"); return t}())},
		// 									ScheduleWeeklyFrequency: to.Ptr[int32](0),
		// 								},
		// 								TimeZone: to.Ptr("Pacific Standard Time"),
		// 							},
		// 						},
		// 						{
		// 							Name: to.Ptr("v2-daily-policy"),
		// 							Type: to.Ptr("Microsoft.RecoveryServices/vaults/backupPolicies"),
		// 							ID: to.Ptr("/Subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SwaggerTestRg/providers/Microsoft.RecoveryServices/vaults/NetSDKTestRsVault/backupPolicies/v2-daily-policy"),
		// 							Properties: &armrecoveryservicesbackup.AzureIaaSVMProtectionPolicy{
		// 								BackupManagementType: to.Ptr("AzureIaasVM"),
		// 								ProtectedItemsCount: to.Ptr[int32](0),
		// 								RetentionPolicy: &armrecoveryservicesbackup.LongTermRetentionPolicy{
		// 									RetentionPolicyType: to.Ptr("LongTermRetentionPolicy"),
		// 									DailySchedule: &armrecoveryservicesbackup.DailyRetentionSchedule{
		// 										RetentionDuration: &armrecoveryservicesbackup.RetentionDuration{
		// 											Count: to.Ptr[int32](1),
		// 											DurationType: to.Ptr(armrecoveryservicesbackup.RetentionDurationTypeDays),
		// 										},
		// 										RetentionTimes: []*time.Time{
		// 											to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-01-24T02:00:00Z"); return t}())},
		// 										},
		// 									},
		// 									SchedulePolicy: &armrecoveryservicesbackup.SimpleSchedulePolicyV2{
		// 										SchedulePolicyType: to.Ptr("SimpleSchedulePolicyV2"),
		// 										DailySchedule: &armrecoveryservicesbackup.DailySchedule{
		// 											ScheduleRunTimes: []*time.Time{
		// 												to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-01-24T10:00:00Z"); return t}())},
		// 											},
		// 											ScheduleRunFrequency: to.Ptr(armrecoveryservicesbackup.ScheduleRunTypeDaily),
		// 										},
		// 										TimeZone: to.Ptr("Pacific Standard Time"),
		// 									},
		// 							}},
		// 						}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/80c21c17b4a7aa57f637ee594f7cfd653255a7e0/specification/recoveryservicesbackup/resource-manager/Microsoft.RecoveryServices/stable/2023-04-01/examples/AzureWorkload/BackupPolicies_List.json
func ExampleBackupPoliciesClient_NewListPager_listProtectionPoliciesWithBackupManagementTypeFilterAsAzureWorkload() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armrecoveryservicesbackup.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewBackupPoliciesClient().NewListPager("NetSDKTestRsVault", "SwaggerTestRg", &armrecoveryservicesbackup.BackupPoliciesClientListOptions{Filter: to.Ptr("backupManagementType eq 'AzureWorkload'")})
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
		// page.ProtectionPolicyResourceList = armrecoveryservicesbackup.ProtectionPolicyResourceList{
		// 	Value: []*armrecoveryservicesbackup.ProtectionPolicyResource{
		// 		{
		// 			Name: to.Ptr("HourlyLogBackup"),
		// 			Type: to.Ptr("Microsoft.RecoveryServices/vaults/backupPolicies"),
		// 			ID: to.Ptr("/Subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SwaggerTestRg/providers/Microsoft.RecoveryServices/vaults/NetSDKTestRsVault/backupPolicies/HourlyLogBackup"),
		// 			Properties: &armrecoveryservicesbackup.AzureVMWorkloadProtectionPolicy{
		// 				BackupManagementType: to.Ptr("AzureWorkload"),
		// 				ProtectedItemsCount: to.Ptr[int32](0),
		// 				Settings: &armrecoveryservicesbackup.Settings{
		// 					Issqlcompression: to.Ptr(false),
		// 					TimeZone: to.Ptr("UTC"),
		// 				},
		// 				SubProtectionPolicy: []*armrecoveryservicesbackup.SubProtectionPolicy{
		// 					{
		// 						PolicyType: to.Ptr(armrecoveryservicesbackup.PolicyTypeFull),
		// 						RetentionPolicy: &armrecoveryservicesbackup.LongTermRetentionPolicy{
		// 							RetentionPolicyType: to.Ptr("LongTermRetentionPolicy"),
		// 							DailySchedule: &armrecoveryservicesbackup.DailyRetentionSchedule{
		// 								RetentionDuration: &armrecoveryservicesbackup.RetentionDuration{
		// 									Count: to.Ptr[int32](30),
		// 									DurationType: to.Ptr(armrecoveryservicesbackup.RetentionDurationTypeDays),
		// 								},
		// 								RetentionTimes: []*time.Time{
		// 									to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-05T19:00:00Z"); return t}())},
		// 								},
		// 							},
		// 							SchedulePolicy: &armrecoveryservicesbackup.SimpleSchedulePolicy{
		// 								SchedulePolicyType: to.Ptr("SimpleSchedulePolicy"),
		// 								ScheduleRunFrequency: to.Ptr(armrecoveryservicesbackup.ScheduleRunTypeDaily),
		// 								ScheduleRunTimes: []*time.Time{
		// 									to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-12-05T19:00:00Z"); return t}())},
		// 									ScheduleWeeklyFrequency: to.Ptr[int32](0),
		// 								},
		// 							},
		// 							{
		// 								PolicyType: to.Ptr(armrecoveryservicesbackup.PolicyTypeLog),
		// 								RetentionPolicy: &armrecoveryservicesbackup.SimpleRetentionPolicy{
		// 									RetentionPolicyType: to.Ptr("SimpleRetentionPolicy"),
		// 									RetentionDuration: &armrecoveryservicesbackup.RetentionDuration{
		// 										Count: to.Ptr[int32](30),
		// 										DurationType: to.Ptr(armrecoveryservicesbackup.RetentionDurationTypeDays),
		// 									},
		// 								},
		// 								SchedulePolicy: &armrecoveryservicesbackup.LogSchedulePolicy{
		// 									SchedulePolicyType: to.Ptr("LogSchedulePolicy"),
		// 									ScheduleFrequencyInMins: to.Ptr[int32](60),
		// 								},
		// 						}},
		// 						WorkLoadType: to.Ptr(armrecoveryservicesbackup.WorkloadTypeSQLDataBase),
		// 					},
		// 			}},
		// 		}
	}
}
