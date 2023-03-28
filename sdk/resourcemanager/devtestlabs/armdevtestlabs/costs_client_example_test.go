//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdevtestlabs_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/devtestlabs/armdevtestlabs/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/Costs_Get.json
func ExampleCostsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCostsClient().Get(ctx, "resourceGroupName", "{labName}", "targetCost", &armdevtestlabs.CostsClientGetOptions{Expand: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.LabCost = armdevtestlabs.LabCost{
	// 	Properties: &armdevtestlabs.LabCostProperties{
	// 		CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-23T22:43:54.7253204+00:00"); return t}()),
	// 		CurrencyCode: to.Ptr("USD"),
	// 		EndDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-31T23:59:59Z"); return t}()),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		StartDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-01T00:00:00Z"); return t}()),
	// 		TargetCost: &armdevtestlabs.TargetCostProperties{
	// 			CostThresholds: []*armdevtestlabs.CostThresholdProperties{
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](25),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusEnabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](50),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusEnabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](75),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](100),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](125),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 			}},
	// 			CycleEndDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-31T23:59:59+00:00"); return t}()),
	// 			CycleStartDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-01T00:00:00+00:00"); return t}()),
	// 			CycleType: to.Ptr(armdevtestlabs.ReportingCycleTypeCalendarMonth),
	// 			Status: to.Ptr(armdevtestlabs.TargetCostStatusEnabled),
	// 			Target: to.Ptr[int32](100),
	// 		},
	// 		UniqueIdentifier: to.Ptr("{uniqueIdentifier}"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/d55b8005f05b040b852c15e74a0f3e36494a15e1/specification/devtestlabs/resource-manager/Microsoft.DevTestLab/stable/2018-09-15/examples/Costs_CreateOrUpdate.json
func ExampleCostsClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdevtestlabs.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewCostsClient().CreateOrUpdate(ctx, "resourceGroupName", "{labName}", "targetCost", armdevtestlabs.LabCost{
		Properties: &armdevtestlabs.LabCostProperties{
			CurrencyCode:  to.Ptr("USD"),
			EndDateTime:   to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-31T23:59:59Z"); return t }()),
			StartDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-01T00:00:00Z"); return t }()),
			TargetCost: &armdevtestlabs.TargetCostProperties{
				CostThresholds: []*armdevtestlabs.CostThresholdProperties{
					{
						DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
							ThresholdValue: to.Ptr[float64](25),
						},
						SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						ThresholdID:                  to.Ptr("00000000-0000-0000-0000-000000000001"),
					},
					{
						DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusEnabled),
						PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
							ThresholdValue: to.Ptr[float64](50),
						},
						SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusEnabled),
						ThresholdID:                  to.Ptr("00000000-0000-0000-0000-000000000002"),
					},
					{
						DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
							ThresholdValue: to.Ptr[float64](75),
						},
						SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						ThresholdID:                  to.Ptr("00000000-0000-0000-0000-000000000003"),
					},
					{
						DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
							ThresholdValue: to.Ptr[float64](100),
						},
						SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						ThresholdID:                  to.Ptr("00000000-0000-0000-0000-000000000004"),
					},
					{
						DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
							ThresholdValue: to.Ptr[float64](125),
						},
						SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
						ThresholdID:                  to.Ptr("00000000-0000-0000-0000-000000000005"),
					}},
				CycleEndDateTime:   to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-31T00:00:00.000Z"); return t }()),
				CycleStartDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-01T00:00:00.000Z"); return t }()),
				CycleType:          to.Ptr(armdevtestlabs.ReportingCycleTypeCalendarMonth),
				Status:             to.Ptr(armdevtestlabs.TargetCostStatusEnabled),
				Target:             to.Ptr[int32](100),
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.LabCost = armdevtestlabs.LabCost{
	// 	Properties: &armdevtestlabs.LabCostProperties{
	// 		CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-23T22:43:54.7253204+00:00"); return t}()),
	// 		CurrencyCode: to.Ptr("USD"),
	// 		EndDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-31T23:59:59Z"); return t}()),
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		StartDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-01T00:00:00Z"); return t}()),
	// 		TargetCost: &armdevtestlabs.TargetCostProperties{
	// 			CostThresholds: []*armdevtestlabs.CostThresholdProperties{
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](25),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusEnabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](50),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusEnabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](75),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](100),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 				},
	// 				{
	// 					DisplayOnChart: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					NotificationSent: to.Ptr("0001-01-01T00:00:00.0000000"),
	// 					PercentageThreshold: &armdevtestlabs.PercentageCostThresholdProperties{
	// 						ThresholdValue: to.Ptr[float64](125),
	// 					},
	// 					SendNotificationWhenExceeded: to.Ptr(armdevtestlabs.CostThresholdStatusDisabled),
	// 					ThresholdID: to.Ptr("00000000-0000-0000-0000-000000000001"),
	// 			}},
	// 			CycleEndDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-31T23:59:59+00:00"); return t}()),
	// 			CycleStartDateTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-12-01T00:00:00+00:00"); return t}()),
	// 			CycleType: to.Ptr(armdevtestlabs.ReportingCycleTypeCalendarMonth),
	// 			Status: to.Ptr(armdevtestlabs.TargetCostStatusEnabled),
	// 			Target: to.Ptr[int32](100),
	// 		},
	// 		UniqueIdentifier: to.Ptr("{uniqueIdentifier}"),
	// 	},
	// }
}
