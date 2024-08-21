//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetapp_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/netapp/armnetapp/v7"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-03-01/examples/SnapshotPolicies_List.json
func ExampleSnapshotPoliciesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSnapshotPoliciesClient().NewListPager("myRG", "account1", nil)
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
		// page.SnapshotPoliciesList = armnetapp.SnapshotPoliciesList{
		// 	Value: []*armnetapp.SnapshotPolicy{
		// 		{
		// 			Name: to.Ptr("account1/snapshotPolicy1"),
		// 			Type: to.Ptr("Microsoft.NetApp/netAppAccounts/snapshotPolicies"),
		// 			ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/snapshotPolicies/snapshotPolicy1"),
		// 			Location: to.Ptr("eastus"),
		// 			Properties: &armnetapp.SnapshotPolicyProperties{
		// 				DailySchedule: &armnetapp.DailySchedule{
		// 					Hour: to.Ptr[int32](14),
		// 					Minute: to.Ptr[int32](30),
		// 					SnapshotsToKeep: to.Ptr[int32](4),
		// 				},
		// 				Enabled: to.Ptr(true),
		// 				HourlySchedule: &armnetapp.HourlySchedule{
		// 					Minute: to.Ptr[int32](50),
		// 					SnapshotsToKeep: to.Ptr[int32](2),
		// 				},
		// 				MonthlySchedule: &armnetapp.MonthlySchedule{
		// 					DaysOfMonth: to.Ptr("10,11,12"),
		// 					Hour: to.Ptr[int32](14),
		// 					Minute: to.Ptr[int32](15),
		// 					SnapshotsToKeep: to.Ptr[int32](5),
		// 				},
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				WeeklySchedule: &armnetapp.WeeklySchedule{
		// 					Day: to.Ptr("Wednesday"),
		// 					Hour: to.Ptr[int32](14),
		// 					Minute: to.Ptr[int32](45),
		// 					SnapshotsToKeep: to.Ptr[int32](3),
		// 				},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-03-01/examples/SnapshotPolicies_Get.json
func ExampleSnapshotPoliciesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSnapshotPoliciesClient().Get(ctx, "myRG", "account1", "snapshotPolicyName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SnapshotPolicy = armnetapp.SnapshotPolicy{
	// 	Name: to.Ptr("account1/snapshotPolicy1"),
	// 	Type: to.Ptr("Microsoft.NetApp/netAppAccounts/snapshotPolicies"),
	// 	ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/snapshotPolicies/snapshotPolicy1"),
	// 	Location: to.Ptr("eastus"),
	// 	Properties: &armnetapp.SnapshotPolicyProperties{
	// 		DailySchedule: &armnetapp.DailySchedule{
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](30),
	// 			SnapshotsToKeep: to.Ptr[int32](4),
	// 		},
	// 		Enabled: to.Ptr(true),
	// 		HourlySchedule: &armnetapp.HourlySchedule{
	// 			Minute: to.Ptr[int32](50),
	// 			SnapshotsToKeep: to.Ptr[int32](2),
	// 		},
	// 		MonthlySchedule: &armnetapp.MonthlySchedule{
	// 			DaysOfMonth: to.Ptr("10,11,12"),
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](15),
	// 			SnapshotsToKeep: to.Ptr[int32](5),
	// 		},
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		WeeklySchedule: &armnetapp.WeeklySchedule{
	// 			Day: to.Ptr("Wednesday"),
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](45),
	// 			SnapshotsToKeep: to.Ptr[int32](3),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-03-01/examples/SnapshotPolicies_Create.json
func ExampleSnapshotPoliciesClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSnapshotPoliciesClient().Create(ctx, "myRG", "account1", "snapshotPolicyName", armnetapp.SnapshotPolicy{
		Location: to.Ptr("eastus"),
		Properties: &armnetapp.SnapshotPolicyProperties{
			DailySchedule: &armnetapp.DailySchedule{
				Hour:            to.Ptr[int32](14),
				Minute:          to.Ptr[int32](30),
				SnapshotsToKeep: to.Ptr[int32](4),
			},
			Enabled: to.Ptr(true),
			HourlySchedule: &armnetapp.HourlySchedule{
				Minute:          to.Ptr[int32](50),
				SnapshotsToKeep: to.Ptr[int32](2),
			},
			MonthlySchedule: &armnetapp.MonthlySchedule{
				DaysOfMonth:     to.Ptr("10,11,12"),
				Hour:            to.Ptr[int32](14),
				Minute:          to.Ptr[int32](15),
				SnapshotsToKeep: to.Ptr[int32](5),
			},
			WeeklySchedule: &armnetapp.WeeklySchedule{
				Day:             to.Ptr("Wednesday"),
				Hour:            to.Ptr[int32](14),
				Minute:          to.Ptr[int32](45),
				SnapshotsToKeep: to.Ptr[int32](3),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SnapshotPolicy = armnetapp.SnapshotPolicy{
	// 	Name: to.Ptr("account1/snapshotPolicy1"),
	// 	Type: to.Ptr("Microsoft.NetApp/netAppAccounts/snapshotPolicies"),
	// 	ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/snapshotPolicies/snapshotPolicy1"),
	// 	Location: to.Ptr("eastus"),
	// 	Properties: &armnetapp.SnapshotPolicyProperties{
	// 		DailySchedule: &armnetapp.DailySchedule{
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](30),
	// 			SnapshotsToKeep: to.Ptr[int32](4),
	// 		},
	// 		Enabled: to.Ptr(true),
	// 		HourlySchedule: &armnetapp.HourlySchedule{
	// 			Minute: to.Ptr[int32](50),
	// 			SnapshotsToKeep: to.Ptr[int32](2),
	// 		},
	// 		MonthlySchedule: &armnetapp.MonthlySchedule{
	// 			DaysOfMonth: to.Ptr("10,11,12"),
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](15),
	// 			SnapshotsToKeep: to.Ptr[int32](5),
	// 		},
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		WeeklySchedule: &armnetapp.WeeklySchedule{
	// 			Day: to.Ptr("Wednesday"),
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](45),
	// 			SnapshotsToKeep: to.Ptr[int32](3),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-03-01/examples/SnapshotPolicies_Update.json
func ExampleSnapshotPoliciesClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSnapshotPoliciesClient().BeginUpdate(ctx, "myRG", "account1", "snapshotPolicyName", armnetapp.SnapshotPolicyPatch{
		Location: to.Ptr("eastus"),
		Properties: &armnetapp.SnapshotPolicyProperties{
			DailySchedule: &armnetapp.DailySchedule{
				Hour:            to.Ptr[int32](14),
				Minute:          to.Ptr[int32](30),
				SnapshotsToKeep: to.Ptr[int32](4),
			},
			Enabled: to.Ptr(true),
			HourlySchedule: &armnetapp.HourlySchedule{
				Minute:          to.Ptr[int32](50),
				SnapshotsToKeep: to.Ptr[int32](2),
			},
			MonthlySchedule: &armnetapp.MonthlySchedule{
				DaysOfMonth:     to.Ptr("10,11,12"),
				Hour:            to.Ptr[int32](14),
				Minute:          to.Ptr[int32](15),
				SnapshotsToKeep: to.Ptr[int32](5),
			},
			WeeklySchedule: &armnetapp.WeeklySchedule{
				Day:             to.Ptr("Wednesday"),
				Hour:            to.Ptr[int32](14),
				Minute:          to.Ptr[int32](45),
				SnapshotsToKeep: to.Ptr[int32](3),
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
	// res.SnapshotPolicy = armnetapp.SnapshotPolicy{
	// 	Name: to.Ptr("account1/snapshotPolicy1"),
	// 	Type: to.Ptr("Microsoft.NetApp/netAppAccounts/snapshotPolicies"),
	// 	ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/snapshotPolicies/snapshotPolicy1"),
	// 	Location: to.Ptr("eastus"),
	// 	Properties: &armnetapp.SnapshotPolicyProperties{
	// 		DailySchedule: &armnetapp.DailySchedule{
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](30),
	// 			SnapshotsToKeep: to.Ptr[int32](4),
	// 		},
	// 		Enabled: to.Ptr(true),
	// 		HourlySchedule: &armnetapp.HourlySchedule{
	// 			Minute: to.Ptr[int32](50),
	// 			SnapshotsToKeep: to.Ptr[int32](2),
	// 		},
	// 		MonthlySchedule: &armnetapp.MonthlySchedule{
	// 			DaysOfMonth: to.Ptr("10,11,12"),
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](15),
	// 			SnapshotsToKeep: to.Ptr[int32](5),
	// 		},
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		WeeklySchedule: &armnetapp.WeeklySchedule{
	// 			Day: to.Ptr("Wednesday"),
	// 			Hour: to.Ptr[int32](14),
	// 			Minute: to.Ptr[int32](45),
	// 			SnapshotsToKeep: to.Ptr[int32](3),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-03-01/examples/SnapshotPolicies_Delete.json
func ExampleSnapshotPoliciesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSnapshotPoliciesClient().BeginDelete(ctx, "resourceGroup", "accountName", "snapshotPolicyName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/c7b98b36e4023331545051284d8500adf98f02fe/specification/netapp/resource-manager/Microsoft.NetApp/stable/2024-03-01/examples/SnapshotPolicies_ListVolumes.json
func ExampleSnapshotPoliciesClient_ListVolumes() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetapp.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSnapshotPoliciesClient().ListVolumes(ctx, "myRG", "account1", "snapshotPolicyName", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SnapshotPolicyVolumeList = armnetapp.SnapshotPolicyVolumeList{
	// 	Value: []*armnetapp.Volume{
	// 		{
	// 			Name: to.Ptr("account1/pool1/volume1"),
	// 			Type: to.Ptr("Microsoft.NetApp/netAppAccounts/capacityPools/volumes"),
	// 			ID: to.Ptr("/subscriptions/D633CC2E-722B-4AE1-B636-BBD9E4C60ED9/resourceGroups/myRG/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1"),
	// 			Location: to.Ptr("eastus"),
	// 			Properties: &armnetapp.VolumeProperties{
	// 				CreationToken: to.Ptr("some-amazing-filepath"),
	// 				FileSystemID: to.Ptr("9760acf5-4638-11e7-9bdb-020073ca7778"),
	// 				ProvisioningState: to.Ptr("Succeeded"),
	// 				ServiceLevel: to.Ptr(armnetapp.ServiceLevelPremium),
	// 				SubnetID: to.Ptr("/subscriptions/9760acf5-4638-11e7-9bdb-020073ca7778/resourceGroups/myRP/providers/Microsoft.Network/virtualNetworks/testvnet3/subnets/testsubnet3"),
	// 				ThroughputMibps: to.Ptr[float32](128),
	// 				UsageThreshold: to.Ptr[int64](107374182400),
	// 			},
	// 	}},
	// }
}
