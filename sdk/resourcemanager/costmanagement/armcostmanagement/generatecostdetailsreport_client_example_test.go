//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcostmanagement_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/GenerateCostDetailsReportByBillingAccountEnterpriseAgreementCustomerAndBillingPeriod.json
func ExampleGenerateCostDetailsReportClient_BeginCreateOperation_generateCostDetailsReportByBillingAccountEnterpriseAgreementCustomerAndBillingPeriod() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginCreateOperation(ctx, "providers/Microsoft.Billing/billingAccounts/12345", armcostmanagement.GenerateCostDetailsReportRequestDefinition{
		BillingPeriod: to.Ptr("202205"),
		Metric:        to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				BillingPeriod: to.Ptr("202205"),
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 			},
	// 			RequestScope: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/GenerateCostDetailsReportByBillingProfileAndInvoiceId.json
func ExampleGenerateCostDetailsReportClient_BeginCreateOperation_generateCostDetailsReportByBillingProfileAndInvoiceId() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginCreateOperation(ctx, "providers/Microsoft.Billing/billingAccounts/12345:6789/billingProfiles/13579", armcostmanagement.GenerateCostDetailsReportRequestDefinition{
		InvoiceID: to.Ptr("M1234567"),
		Metric:    to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345:6789/billingProfiles/13579/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				BillingPeriod: to.Ptr("202205"),
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 			},
	// 			RequestScope: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345:6789/billingProfiles/13579"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/GenerateCostDetailsReportByBillingProfileAndInvoiceIdAndCustomerId.json
func ExampleGenerateCostDetailsReportClient_BeginCreateOperation_generateCostDetailsReportByBillingProfileAndInvoiceIdAndCustomerId() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginCreateOperation(ctx, "providers/Microsoft.Billing/billingAccounts/12345:6789/customers/13579", armcostmanagement.GenerateCostDetailsReportRequestDefinition{
		InvoiceID: to.Ptr("M1234567"),
		Metric:    to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345:6789/customers/13579/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				BillingPeriod: to.Ptr("202205"),
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 			},
	// 			RequestScope: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345:6789/customers/13579"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/GenerateCostDetailsReportByCustomerAndTimePeriod.json
func ExampleGenerateCostDetailsReportClient_BeginCreateOperation_generateCostDetailsReportByCustomerAndTimePeriod() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginCreateOperation(ctx, "providers/Microsoft.Billing/billingAccounts/12345:6789/customers/13579", armcostmanagement.GenerateCostDetailsReportRequestDefinition{
		Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
		TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
			End:   to.Ptr("2020-03-15"),
			Start: to.Ptr("2020-03-01"),
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345:6789/customers/13579/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 				TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
	// 					End: to.Ptr("2020-03-15"),
	// 					Start: to.Ptr("2020-03-01"),
	// 				},
	// 			},
	// 			RequestScope: to.Ptr("providers/Microsoft.Billing/billingAccounts/12345:6789/customers/13579"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/GenerateCostDetailsReportByDepartmentsAndTimePeriod.json
func ExampleGenerateCostDetailsReportClient_BeginCreateOperation_generateCostDetailsReportByDepartmentsAndTimePeriod() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginCreateOperation(ctx, "providers/Microsoft.Billing/departments/12345", armcostmanagement.GenerateCostDetailsReportRequestDefinition{
		Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
		TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
			End:   to.Ptr("2020-03-15"),
			Start: to.Ptr("2020-03-01"),
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("providers/Microsoft.Billing/departments/12345/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 				TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
	// 					End: to.Ptr("2020-03-15"),
	// 					Start: to.Ptr("2020-03-01"),
	// 				},
	// 			},
	// 			RequestScope: to.Ptr("providers/Microsoft.Billing/departments/12345"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/GenerateCostDetailsReportByEnrollmentAccountsAndTimePeriod.json
func ExampleGenerateCostDetailsReportClient_BeginCreateOperation_generateCostDetailsReportByEnrollmentAccountsAndTimePeriod() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginCreateOperation(ctx, "providers/Microsoft.Billing/enrollmentAccounts/1234", armcostmanagement.GenerateCostDetailsReportRequestDefinition{
		Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
		TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
			End:   to.Ptr("2020-03-15"),
			Start: to.Ptr("2020-03-01"),
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("providers/Microsoft.Billing/enrollmentAccounts/1234/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 				TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
	// 					End: to.Ptr("2020-03-15"),
	// 					Start: to.Ptr("2020-03-01"),
	// 				},
	// 			},
	// 			RequestScope: to.Ptr("providers/Microsoft.Billing/enrollmentAccounts/1234"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/GenerateCostDetailsReportBySubscriptionAndTimePeriod.json
func ExampleGenerateCostDetailsReportClient_BeginCreateOperation_generateCostDetailsReportBySubscriptionAndTimePeriod() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginCreateOperation(ctx, "subscriptions/00000000-0000-0000-0000-000000000000", armcostmanagement.GenerateCostDetailsReportRequestDefinition{
		Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
		TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
			End:   to.Ptr("2020-03-15"),
			Start: to.Ptr("2020-03-01"),
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 				TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
	// 					End: to.Ptr("2020-03-15"),
	// 					Start: to.Ptr("2020-03-01"),
	// 				},
	// 			},
	// 			RequestScope: to.Ptr("subscriptions/00000000-0000-0000-0000-000000000000"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/17aa6a1314de5aafef059d9aa2229901df506e75/specification/cost-management/resource-manager/Microsoft.CostManagement/stable/2022-10-01/examples/CostDetailsOperationResultsBySubscriptionScope.json
func ExampleGenerateCostDetailsReportClient_BeginGetOperationResults() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcostmanagement.NewClientFactory(cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGenerateCostDetailsReportClient().BeginGetOperationResults(ctx, "subscriptions/00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", nil)
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
	// res.CostDetailsOperationResults = armcostmanagement.CostDetailsOperationResults{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	ID: to.Ptr("subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/costDetailsOperationResults/00000000-0000-0000-0000-000000000000"),
	// 	Manifest: &armcostmanagement.ReportManifest{
	// 		BlobCount: to.Ptr[int32](1),
	// 		Blobs: []*armcostmanagement.BlobInfo{
	// 			{
	// 				BlobLink: to.Ptr("https://ccmreportstorageeastus.blob.core.windows.net/armreports/00000/00000000-0000-0000-0000-000000000000?sv=2020-05-31&sr=b&sig=abcd"),
	// 				ByteCount: to.Ptr[int64](32741),
	// 		}},
	// 		ByteCount: to.Ptr[int64](32741),
	// 		CompressData: to.Ptr(false),
	// 		DataFormat: to.Ptr(armcostmanagement.CostDetailsDataFormatCSVCostDetailsDataFormat),
	// 		ManifestVersion: to.Ptr("2022-10-01"),
	// 		RequestContext: &armcostmanagement.RequestContext{
	// 			RequestBody: &armcostmanagement.GenerateCostDetailsReportRequestDefinition{
	// 				Metric: to.Ptr(armcostmanagement.CostDetailsMetricTypeActualCostCostDetailsMetricType),
	// 				TimePeriod: &armcostmanagement.CostDetailsTimePeriod{
	// 					End: to.Ptr("2020-03-15"),
	// 					Start: to.Ptr("2020-03-01"),
	// 				},
	// 			},
	// 			RequestScope: to.Ptr("subscriptions/00000000-0000-0000-0000-000000000000"),
	// 		},
	// 	},
	// 	Status: to.Ptr(armcostmanagement.CostDetailsStatusTypeCompletedCostDetailsStatusType),
	// 	ValidTill: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2022-05-10T08:08:46.1973252Z"); return t}()),
	// }
}
