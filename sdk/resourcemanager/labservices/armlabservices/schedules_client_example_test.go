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

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/labservices/armlabservices"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/Schedules/listSchedule.json
func ExampleSchedulesClient_NewListByLabPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewSchedulesClient().NewListByLabPager("testrg123", "testlab", &armlabservices.SchedulesClientListByLabOptions{Filter: nil})
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
		// page.PagedSchedules = armlabservices.PagedSchedules{
		// 	Value: []*armlabservices.Schedule{
		// 		{
		// 			Name: to.Ptr("schedule1"),
		// 			Type: to.Ptr("Microsoft.LabServices/Schedule"),
		// 			ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labs/testlab/schedules/schedule1"),
		// 			Properties: &armlabservices.ScheduleProperties{
		// 				Notes: to.Ptr("Schedule 1 for students"),
		// 				RecurrencePattern: &armlabservices.RecurrencePattern{
		// 					ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t}()),
		// 					Frequency: to.Ptr(armlabservices.RecurrenceFrequencyDaily),
		// 					Interval: to.Ptr[int32](1),
		// 				},
		// 				StartAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T12:00:00.000Z"); return t}()),
		// 				StopAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T18:00:00.000Z"); return t}()),
		// 				TimeZoneID: to.Ptr("America/Los_Angeles"),
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
		// 		},
		// 		{
		// 			Name: to.Ptr("schedule2"),
		// 			Type: to.Ptr("Microsoft.LabServices/Schedule"),
		// 			ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labs/testlab/schedules/schedule2"),
		// 			Properties: &armlabservices.ScheduleProperties{
		// 				Notes: to.Ptr("Schedule 2 for students"),
		// 				RecurrencePattern: &armlabservices.RecurrencePattern{
		// 					ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t}()),
		// 					Frequency: to.Ptr(armlabservices.RecurrenceFrequencyWeekly),
		// 					WeekDays: []*armlabservices.WeekDay{
		// 						to.Ptr(armlabservices.WeekDayMonday),
		// 						to.Ptr(armlabservices.WeekDayTuesday),
		// 						to.Ptr(armlabservices.WeekDayWednesday),
		// 						to.Ptr(armlabservices.WeekDayThursday),
		// 						to.Ptr(armlabservices.WeekDayFriday)},
		// 					},
		// 					StartAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T12:00:00.000Z"); return t}()),
		// 					StopAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T18:00:00.000Z"); return t}()),
		// 					TimeZoneID: to.Ptr("America/Los_Angeles"),
		// 					ProvisioningState: to.Ptr(armlabservices.ProvisioningStateSucceeded),
		// 				},
		// 				SystemData: &armlabservices.SystemData{
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-01T10:00:00.000Z"); return t}()),
		// 					CreatedBy: to.Ptr("identity123"),
		// 					CreatedByType: to.Ptr(armlabservices.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-06-01T09:12:28.000Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("identity123"),
		// 					LastModifiedByType: to.Ptr(armlabservices.CreatedByTypeUser),
		// 				},
		// 		}},
		// 	}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/Schedules/getSchedule.json
func ExampleSchedulesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSchedulesClient().Get(ctx, "testrg123", "testlab", "schedule1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Schedule = armlabservices.Schedule{
	// 	Name: to.Ptr("schedule1"),
	// 	Type: to.Ptr("Microsoft.LabServices/Schedule"),
	// 	ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labs/testlab/schedules/schedule1"),
	// 	Properties: &armlabservices.ScheduleProperties{
	// 		Notes: to.Ptr("Schedule 1 for students"),
	// 		RecurrencePattern: &armlabservices.RecurrencePattern{
	// 			ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t}()),
	// 			Frequency: to.Ptr(armlabservices.RecurrenceFrequencyDaily),
	// 			Interval: to.Ptr[int32](1),
	// 		},
	// 		StartAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T12:00:00.000Z"); return t}()),
	// 		StopAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T18:00:00.000Z"); return t}()),
	// 		TimeZoneID: to.Ptr("America/Los_Angeles"),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/Schedules/putSchedule.json
func ExampleSchedulesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSchedulesClient().CreateOrUpdate(ctx, "testrg123", "testlab", "schedule1", armlabservices.Schedule{
		Properties: &armlabservices.ScheduleProperties{
			Notes: to.Ptr("Schedule 1 for students"),
			RecurrencePattern: &armlabservices.RecurrencePattern{
				ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t }()),
				Frequency:      to.Ptr(armlabservices.RecurrenceFrequencyDaily),
				Interval:       to.Ptr[int32](2),
			},
			StartAt:    to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T12:00:00.000Z"); return t }()),
			StopAt:     to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T18:00:00.000Z"); return t }()),
			TimeZoneID: to.Ptr("America/Los_Angeles"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Schedule = armlabservices.Schedule{
	// 	Name: to.Ptr("schedule1"),
	// 	Type: to.Ptr("Microsoft.LabServices/Schedule"),
	// 	ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labs/testlab/schedules/schedule1"),
	// 	Properties: &armlabservices.ScheduleProperties{
	// 		Notes: to.Ptr("Schedule 1 for students"),
	// 		RecurrencePattern: &armlabservices.RecurrencePattern{
	// 			ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t}()),
	// 			Frequency: to.Ptr(armlabservices.RecurrenceFrequencyDaily),
	// 			Interval: to.Ptr[int32](2),
	// 		},
	// 		StartAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T12:00:00.000Z"); return t}()),
	// 		StopAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T18:00:00.000Z"); return t}()),
	// 		TimeZoneID: to.Ptr("America/Los_Angeles"),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/Schedules/patchSchedule.json
func ExampleSchedulesClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewSchedulesClient().Update(ctx, "testrg123", "testlab", "schedule1", armlabservices.ScheduleUpdate{
		Properties: &armlabservices.ScheduleUpdateProperties{
			RecurrencePattern: &armlabservices.RecurrencePattern{
				ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t }()),
				Frequency:      to.Ptr(armlabservices.RecurrenceFrequencyDaily),
				Interval:       to.Ptr[int32](2),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Schedule = armlabservices.Schedule{
	// 	Name: to.Ptr("schedule1"),
	// 	Type: to.Ptr("Microsoft.LabServices/Schedule"),
	// 	ID: to.Ptr("/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.LabServices/labs/testlab/schedules/schedule1"),
	// 	Properties: &armlabservices.ScheduleProperties{
	// 		Notes: to.Ptr("Schedule 1 for students"),
	// 		RecurrencePattern: &armlabservices.RecurrencePattern{
	// 			ExpirationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-08-14T23:59:59.000Z"); return t}()),
	// 			Frequency: to.Ptr(armlabservices.RecurrenceFrequencyDaily),
	// 			Interval: to.Ptr[int32](2),
	// 		},
	// 		StartAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T12:00:00.000Z"); return t}()),
	// 		StopAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-05-26T18:00:00.000Z"); return t}()),
	// 		TimeZoneID: to.Ptr("America/Los_Angeles"),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/4c2cdccf6ca3281dd50ed8788ce1de2e0d480973/specification/labservices/resource-manager/Microsoft.LabServices/stable/2022-08-01/examples/Schedules/deleteSchedule.json
func ExampleSchedulesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armlabservices.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewSchedulesClient().BeginDelete(ctx, "testrg123", "testlab", "schedule1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}
