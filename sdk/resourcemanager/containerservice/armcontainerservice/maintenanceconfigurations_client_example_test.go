//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcontainerservice_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsList.json
func ExampleMaintenanceConfigurationsClient_NewListByManagedClusterPager_listMaintenanceConfigurationsByManagedCluster() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewMaintenanceConfigurationsClient().NewListByManagedClusterPager("rg1", "clustername1", nil)
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
		// page.MaintenanceConfigurationListResult = armcontainerservice.MaintenanceConfigurationListResult{
		// 	Value: []*armcontainerservice.MaintenanceConfiguration{
		// 		{
		// 			Name: to.Ptr("default"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1/maintenanceConfigurations/default"),
		// 			Properties: &armcontainerservice.MaintenanceConfigurationProperties{
		// 				NotAllowedTime: []*armcontainerservice.TimeSpan{
		// 					{
		// 						End: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-30T12:00:00.000Z"); return t}()),
		// 						Start: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-26T03:00:00.000Z"); return t}()),
		// 				}},
		// 				TimeInWeek: []*armcontainerservice.TimeInWeek{
		// 					{
		// 						Day: to.Ptr(armcontainerservice.WeekDayMonday),
		// 						HourSlots: []*int32{
		// 							to.Ptr[int32](1),
		// 							to.Ptr[int32](2)},
		// 					}},
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsList_MaintenanceWindow.json
func ExampleMaintenanceConfigurationsClient_NewListByManagedClusterPager_listMaintenanceConfigurationsConfiguredWithMaintenanceWindowByManagedCluster() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewMaintenanceConfigurationsClient().NewListByManagedClusterPager("rg1", "clustername1", nil)
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
		// page.MaintenanceConfigurationListResult = armcontainerservice.MaintenanceConfigurationListResult{
		// 	Value: []*armcontainerservice.MaintenanceConfiguration{
		// 		{
		// 			Name: to.Ptr("aksManagedNodeOSUpgradeSchedule"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1/maintenanceConfigurations/aksManagedNodeOSUpgradeSchedule"),
		// 			Properties: &armcontainerservice.MaintenanceConfigurationProperties{
		// 				MaintenanceWindow: &armcontainerservice.MaintenanceWindow{
		// 					DurationHours: to.Ptr[int32](10),
		// 					Schedule: &armcontainerservice.Schedule{
		// 						Daily: &armcontainerservice.DailySchedule{
		// 							IntervalDays: to.Ptr[int32](5),
		// 						},
		// 					},
		// 					StartDate: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-01-01"); return t}()),
		// 					StartTime: to.Ptr("13:30"),
		// 					UTCOffset: to.Ptr("-07:00"),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("aksManagedAutoUpgradeSchedule"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1/maintenanceConfigurations/aksManagedAutoUpgradeSchedule"),
		// 			Properties: &armcontainerservice.MaintenanceConfigurationProperties{
		// 				MaintenanceWindow: &armcontainerservice.MaintenanceWindow{
		// 					DurationHours: to.Ptr[int32](5),
		// 					NotAllowedDates: []*armcontainerservice.DateSpan{
		// 						{
		// 							End: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-25"); return t}()),
		// 							Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-18"); return t}()),
		// 						},
		// 						{
		// 							End: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2024-01-05"); return t}()),
		// 							Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-12-23"); return t}()),
		// 					}},
		// 					Schedule: &armcontainerservice.Schedule{
		// 						AbsoluteMonthly: &armcontainerservice.AbsoluteMonthlySchedule{
		// 							DayOfMonth: to.Ptr[int32](15),
		// 							IntervalMonths: to.Ptr[int32](3),
		// 						},
		// 					},
		// 					StartDate: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-01-01"); return t}()),
		// 					StartTime: to.Ptr("08:30"),
		// 					UTCOffset: to.Ptr("+00:00"),
		// 				},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsGet.json
func ExampleMaintenanceConfigurationsClient_Get_getMaintenanceConfiguration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewMaintenanceConfigurationsClient().Get(ctx, "rg1", "clustername1", "default", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.MaintenanceConfiguration = armcontainerservice.MaintenanceConfiguration{
	// 	Name: to.Ptr("default"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1/maintenanceConfigurations/default"),
	// 	Properties: &armcontainerservice.MaintenanceConfigurationProperties{
	// 		NotAllowedTime: []*armcontainerservice.TimeSpan{
	// 			{
	// 				End: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-30T12:00:00.000Z"); return t}()),
	// 				Start: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-26T03:00:00.000Z"); return t}()),
	// 		}},
	// 		TimeInWeek: []*armcontainerservice.TimeInWeek{
	// 			{
	// 				Day: to.Ptr(armcontainerservice.WeekDayMonday),
	// 				HourSlots: []*int32{
	// 					to.Ptr[int32](1),
	// 					to.Ptr[int32](2)},
	// 			}},
	// 		},
	// 		SystemData: &armcontainerservice.SystemData{
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-01T17:18:19.123Z"); return t}()),
	// 			CreatedBy: to.Ptr("user1"),
	// 			CreatedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-02T17:18:19.123Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("user2"),
	// 			LastModifiedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsGet_MaintenanceWindow.json
func ExampleMaintenanceConfigurationsClient_Get_getMaintenanceConfigurationConfiguredWithMaintenanceWindow() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewMaintenanceConfigurationsClient().Get(ctx, "rg1", "clustername1", "aksManagedNodeOSUpgradeSchedule", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.MaintenanceConfiguration = armcontainerservice.MaintenanceConfiguration{
	// 	Name: to.Ptr("aksManagedNodeOSUpgradeSchedule"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1/maintenanceConfigurations/aksManagedNodeOSUpgradeSchedule"),
	// 	Properties: &armcontainerservice.MaintenanceConfigurationProperties{
	// 		MaintenanceWindow: &armcontainerservice.MaintenanceWindow{
	// 			DurationHours: to.Ptr[int32](4),
	// 			NotAllowedDates: []*armcontainerservice.DateSpan{
	// 				{
	// 					End: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-25"); return t}()),
	// 					Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-18"); return t}()),
	// 				},
	// 				{
	// 					End: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2024-01-05"); return t}()),
	// 					Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-12-23"); return t}()),
	// 			}},
	// 			Schedule: &armcontainerservice.Schedule{
	// 				Daily: &armcontainerservice.DailySchedule{
	// 					IntervalDays: to.Ptr[int32](3),
	// 				},
	// 			},
	// 			StartDate: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-01-01"); return t}()),
	// 			StartTime: to.Ptr("09:30"),
	// 			UTCOffset: to.Ptr("-07:00"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsCreate_Update.json
func ExampleMaintenanceConfigurationsClient_CreateOrUpdate_createUpdateMaintenanceConfiguration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewMaintenanceConfigurationsClient().CreateOrUpdate(ctx, "rg1", "clustername1", "default", armcontainerservice.MaintenanceConfiguration{
		Properties: &armcontainerservice.MaintenanceConfigurationProperties{
			NotAllowedTime: []*armcontainerservice.TimeSpan{
				{
					End:   to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-30T12:00:00.000Z"); return t }()),
					Start: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-26T03:00:00.000Z"); return t }()),
				}},
			TimeInWeek: []*armcontainerservice.TimeInWeek{
				{
					Day: to.Ptr(armcontainerservice.WeekDayMonday),
					HourSlots: []*int32{
						to.Ptr[int32](1),
						to.Ptr[int32](2)},
				}},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.MaintenanceConfiguration = armcontainerservice.MaintenanceConfiguration{
	// 	Name: to.Ptr("default"),
	// 	Type: to.Ptr("Microsoft.ContainerService/managedClusters/maintenanceConfigurations"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1/maintenanceConfigurations/default"),
	// 	Properties: &armcontainerservice.MaintenanceConfigurationProperties{
	// 		NotAllowedTime: []*armcontainerservice.TimeSpan{
	// 			{
	// 				End: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-30T12:00:00.000Z"); return t}()),
	// 				Start: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-11-26T03:00:00.000Z"); return t}()),
	// 		}},
	// 		TimeInWeek: []*armcontainerservice.TimeInWeek{
	// 			{
	// 				Day: to.Ptr(armcontainerservice.WeekDayMonday),
	// 				HourSlots: []*int32{
	// 					to.Ptr[int32](1),
	// 					to.Ptr[int32](2)},
	// 			}},
	// 		},
	// 		SystemData: &armcontainerservice.SystemData{
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-01T17:18:19.123Z"); return t}()),
	// 			CreatedBy: to.Ptr("user1"),
	// 			CreatedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-02T17:18:19.123Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("user2"),
	// 			LastModifiedByType: to.Ptr(armcontainerservice.CreatedByTypeUser),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsCreate_Update_MaintenanceWindow.json
func ExampleMaintenanceConfigurationsClient_CreateOrUpdate_createUpdateMaintenanceConfigurationWithMaintenanceWindow() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewMaintenanceConfigurationsClient().CreateOrUpdate(ctx, "rg1", "clustername1", "aksManagedAutoUpgradeSchedule", armcontainerservice.MaintenanceConfiguration{
		Properties: &armcontainerservice.MaintenanceConfigurationProperties{
			MaintenanceWindow: &armcontainerservice.MaintenanceWindow{
				DurationHours: to.Ptr[int32](10),
				NotAllowedDates: []*armcontainerservice.DateSpan{
					{
						End:   to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-25"); return t }()),
						Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-18"); return t }()),
					},
					{
						End:   to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2024-01-05"); return t }()),
						Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-12-23"); return t }()),
					}},
				Schedule: &armcontainerservice.Schedule{
					RelativeMonthly: &armcontainerservice.RelativeMonthlySchedule{
						DayOfWeek:      to.Ptr(armcontainerservice.WeekDayMonday),
						IntervalMonths: to.Ptr[int32](3),
						WeekIndex:      to.Ptr(armcontainerservice.TypeFirst),
					},
				},
				StartDate: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-01-01"); return t }()),
				StartTime: to.Ptr("08:30"),
				UTCOffset: to.Ptr("+05:30"),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.MaintenanceConfiguration = armcontainerservice.MaintenanceConfiguration{
	// 	Name: to.Ptr("aksManagedAutoUpgradeSchedule"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1/maintenanceConfigurations/aksManagedAutoUpgradeSchedule"),
	// 	Properties: &armcontainerservice.MaintenanceConfigurationProperties{
	// 		MaintenanceWindow: &armcontainerservice.MaintenanceWindow{
	// 			DurationHours: to.Ptr[int32](10),
	// 			NotAllowedDates: []*armcontainerservice.DateSpan{
	// 				{
	// 					End: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-25"); return t}()),
	// 					Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-02-18"); return t}()),
	// 				},
	// 				{
	// 					End: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2024-01-05"); return t}()),
	// 					Start: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-12-23"); return t}()),
	// 			}},
	// 			Schedule: &armcontainerservice.Schedule{
	// 				Weekly: &armcontainerservice.WeeklySchedule{
	// 					DayOfWeek: to.Ptr(armcontainerservice.WeekDayMonday),
	// 					IntervalWeeks: to.Ptr[int32](3),
	// 				},
	// 			},
	// 			StartDate: to.Ptr(func() time.Time { t, _ := time.Parse("2006-01-02", "2023-01-01"); return t}()),
	// 			StartTime: to.Ptr("08:30"),
	// 			UTCOffset: to.Ptr("+05:30"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsDelete.json
func ExampleMaintenanceConfigurationsClient_Delete_deleteMaintenanceConfiguration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewMaintenanceConfigurationsClient().Delete(ctx, "rg1", "clustername1", "default", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/8b5618d760532b69d6f7434f57172ea52e109c79/specification/containerservice/resource-manager/Microsoft.ContainerService/aks/stable/2023-11-01/examples/MaintenanceConfigurationsDelete_MaintenanceWindow.json
func ExampleMaintenanceConfigurationsClient_Delete_deleteMaintenanceConfigurationForNodeOsUpgrade() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcontainerservice.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewMaintenanceConfigurationsClient().Delete(ctx, "rg1", "clustername1", "aksManagedNodeOSUpgradeSchedule", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
