//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armapimanagement_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsByApi.json
func ExampleReportsClient_NewListByAPIPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListByAPIPager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByAPIOptions{Top: nil,
		Skip:    nil,
		Orderby: nil,
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
		// page.ReportCollection = armapimanagement.ReportCollection{
		// 	Count: to.Ptr[int64](2),
		// 	Value: []*armapimanagement.ReportRecordContract{
		// 		{
		// 			Name: to.Ptr("Echo API"),
		// 			APIID: to.Ptr("/apis/5600b59475ff190048040001"),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 		},
		// 		{
		// 			Name: to.Ptr("httpbin"),
		// 			APIID: to.Ptr("/apis/57a03a13e4bbd5119c8b19e9"),
		// 			APITimeAvg: to.Ptr[float64](1015.7607923076923),
		// 			APITimeMax: to.Ptr[float64](1819.2173),
		// 			APITimeMin: to.Ptr[float64](330.3206),
		// 			Bandwidth: to.Ptr[int64](11019),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](1),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](13),
		// 			CallCountTotal: to.Ptr[int32](14),
		// 			ServiceTimeAvg: to.Ptr[float64](957.094776923077),
		// 			ServiceTimeMax: to.Ptr[float64](1697.3612),
		// 			ServiceTimeMin: to.Ptr[float64](215.24),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsByUser.json
func ExampleReportsClient_NewListByUserPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListByUserPager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByUserOptions{Top: nil,
		Skip:    nil,
		Orderby: nil,
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
		// page.ReportCollection = armapimanagement.ReportCollection{
		// 	Count: to.Ptr[int64](3),
		// 	Value: []*armapimanagement.ReportRecordContract{
		// 		{
		// 			Name: to.Ptr("Administrator"),
		// 			APITimeAvg: to.Ptr[float64](1015.7607923076923),
		// 			APITimeMax: to.Ptr[float64](1819.2173),
		// 			APITimeMin: to.Ptr[float64](330.3206),
		// 			Bandwidth: to.Ptr[int64](11019),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](1),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](13),
		// 			CallCountTotal: to.Ptr[int32](14),
		// 			ServiceTimeAvg: to.Ptr[float64](957.094776923077),
		// 			ServiceTimeMax: to.Ptr[float64](1697.3612),
		// 			ServiceTimeMin: to.Ptr[float64](215.24),
		// 			UserID: to.Ptr("/users/1"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Samir Solanki"),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 			UserID: to.Ptr("/users/56eaec62baf08b06e46d27fd"),
		// 		},
		// 		{
		// 			Name: to.Ptr("Anonymous"),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 			UserID: to.Ptr("/users/54c800b332965a0035030000"),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsByOperation.json
func ExampleReportsClient_NewListByOperationPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListByOperationPager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByOperationOptions{Top: nil,
		Skip:    nil,
		Orderby: nil,
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
		// page.ReportCollection = armapimanagement.ReportCollection{
		// 	Count: to.Ptr[int64](3),
		// 	Value: []*armapimanagement.ReportRecordContract{
		// 		{
		// 			Name: to.Ptr("get"),
		// 			OperationID: to.Ptr("/apis/57a03a13e4bbd5119c8b19e9/operations/57a03a1dd8d14f0a780d7d14"),
		// 			APIID: to.Ptr("/apis/57a03a13e4bbd5119c8b19e9"),
		// 			APITimeAvg: to.Ptr[float64](1015.7607923076923),
		// 			APITimeMax: to.Ptr[float64](1819.2173),
		// 			APITimeMin: to.Ptr[float64](330.3206),
		// 			Bandwidth: to.Ptr[int64](11019),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](1),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](13),
		// 			CallCountTotal: to.Ptr[int32](14),
		// 			ServiceTimeAvg: to.Ptr[float64](957.094776923077),
		// 			ServiceTimeMax: to.Ptr[float64](1697.3612),
		// 			ServiceTimeMin: to.Ptr[float64](215.24),
		// 		},
		// 		{
		// 			Name: to.Ptr("GetWeatherInformation"),
		// 			OperationID: to.Ptr("/apis/57c999d1e4bbd50c988cb2c3/operations/57c999d1e4bbd50df889c93e"),
		// 			APIID: to.Ptr("/apis/57c999d1e4bbd50c988cb2c3"),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 		},
		// 		{
		// 			Name: to.Ptr("GetCityForecastByZIP"),
		// 			OperationID: to.Ptr("/apis/57c999d1e4bbd50c988cb2c3/operations/57c999d1e4bbd50df889c93f"),
		// 			APIID: to.Ptr("/apis/57c999d1e4bbd50c988cb2c3"),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsByProduct.json
func ExampleReportsClient_NewListByProductPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListByProductPager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByProductOptions{Top: nil,
		Skip:    nil,
		Orderby: nil,
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
		// page.ReportCollection = armapimanagement.ReportCollection{
		// 	Count: to.Ptr[int64](2),
		// 	Value: []*armapimanagement.ReportRecordContract{
		// 		{
		// 			Name: to.Ptr("Starter"),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ProductID: to.Ptr("/products/5600b59475ff190048060001"),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 		},
		// 		{
		// 			Name: to.Ptr("Unlimited"),
		// 			APITimeAvg: to.Ptr[float64](1015.7607923076923),
		// 			APITimeMax: to.Ptr[float64](1819.2173),
		// 			APITimeMin: to.Ptr[float64](330.3206),
		// 			Bandwidth: to.Ptr[int64](11019),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](1),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](13),
		// 			CallCountTotal: to.Ptr[int32](14),
		// 			ProductID: to.Ptr("/products/5600b59475ff190048060002"),
		// 			ServiceTimeAvg: to.Ptr[float64](957.094776923077),
		// 			ServiceTimeMax: to.Ptr[float64](1697.3612),
		// 			ServiceTimeMin: to.Ptr[float64](215.24),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsByGeo.json
func ExampleReportsClient_NewListByGeoPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListByGeoPager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByGeoOptions{Top: nil,
		Skip: nil,
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
		// page.ReportCollection = armapimanagement.ReportCollection{
		// 	Value: []*armapimanagement.ReportRecordContract{
		// 		{
		// 			APITimeAvg: to.Ptr[float64](1015.7607923076923),
		// 			APITimeMax: to.Ptr[float64](1819.2173),
		// 			APITimeMin: to.Ptr[float64](330.3206),
		// 			Bandwidth: to.Ptr[int64](11019),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](1),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](13),
		// 			CallCountTotal: to.Ptr[int32](14),
		// 			Country: to.Ptr("US"),
		// 			Region: to.Ptr("WA"),
		// 			ServiceTimeAvg: to.Ptr[float64](957.094776923077),
		// 			ServiceTimeMax: to.Ptr[float64](1697.3612),
		// 			ServiceTimeMin: to.Ptr[float64](215.24),
		// 			Zip: to.Ptr("98052"),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsBySubscription.json
func ExampleReportsClient_NewListBySubscriptionPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListBySubscriptionPager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListBySubscriptionOptions{Top: nil,
		Skip:    nil,
		Orderby: nil,
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
		// page.ReportCollection = armapimanagement.ReportCollection{
		// 	Count: to.Ptr[int64](3),
		// 	Value: []*armapimanagement.ReportRecordContract{
		// 		{
		// 			Name: to.Ptr(""),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ProductID: to.Ptr("/products/5600b59475ff190048060001"),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 			SubscriptionID: to.Ptr("/subscriptions/5600b59475ff190048070001"),
		// 			UserID: to.Ptr("/users/1"),
		// 		},
		// 		{
		// 			Name: to.Ptr(""),
		// 			APITimeAvg: to.Ptr[float64](1015.7607923076923),
		// 			APITimeMax: to.Ptr[float64](1819.2173),
		// 			APITimeMin: to.Ptr[float64](330.3206),
		// 			Bandwidth: to.Ptr[int64](11019),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](1),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](13),
		// 			CallCountTotal: to.Ptr[int32](14),
		// 			ProductID: to.Ptr("/products/5600b59475ff190048060002"),
		// 			ServiceTimeAvg: to.Ptr[float64](957.094776923077),
		// 			ServiceTimeMax: to.Ptr[float64](1697.3612),
		// 			ServiceTimeMin: to.Ptr[float64](215.24),
		// 			SubscriptionID: to.Ptr("/subscriptions/5600b59475ff190048070002"),
		// 			UserID: to.Ptr("/users/1"),
		// 		},
		// 		{
		// 			Name: to.Ptr(""),
		// 			APITimeAvg: to.Ptr[float64](0),
		// 			APITimeMax: to.Ptr[float64](0),
		// 			APITimeMin: to.Ptr[float64](0),
		// 			Bandwidth: to.Ptr[int64](0),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](0),
		// 			CallCountTotal: to.Ptr[int32](0),
		// 			ProductID: to.Ptr("/products/5702e97e5157a50f48dce801"),
		// 			ServiceTimeAvg: to.Ptr[float64](0),
		// 			ServiceTimeMax: to.Ptr[float64](0),
		// 			ServiceTimeMin: to.Ptr[float64](0),
		// 			SubscriptionID: to.Ptr("/subscriptions/5702e97e5157a50a9c733303"),
		// 			UserID: to.Ptr("/users/1"),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsByTime.json
func ExampleReportsClient_NewListByTimePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListByTimePager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", "PT15M", &armapimanagement.ReportsClientListByTimeOptions{Top: nil,
		Skip:    nil,
		Orderby: nil,
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
		// page.ReportCollection = armapimanagement.ReportCollection{
		// 	Count: to.Ptr[int64](2),
		// 	Value: []*armapimanagement.ReportRecordContract{
		// 		{
		// 			APITimeAvg: to.Ptr[float64](1337.46335),
		// 			APITimeMax: to.Ptr[float64](1819.2173),
		// 			APITimeMin: to.Ptr[float64](885.0839000000001),
		// 			Bandwidth: to.Ptr[int64](3243),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](0),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](4),
		// 			CallCountTotal: to.Ptr[int32](4),
		// 			Interval: to.Ptr("PT15M"),
		// 			ServiceTimeAvg: to.Ptr[float64](1255.917425),
		// 			ServiceTimeMax: to.Ptr[float64](1697.3612),
		// 			ServiceTimeMin: to.Ptr[float64](882.8264),
		// 			Timestamp: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-06-03T00:15:00.000Z"); return t}()),
		// 		},
		// 		{
		// 			APITimeAvg: to.Ptr[float64](872.7818777777778),
		// 			APITimeMax: to.Ptr[float64](1093.8407),
		// 			APITimeMin: to.Ptr[float64](330.3206),
		// 			Bandwidth: to.Ptr[int64](7776),
		// 			CacheHitCount: to.Ptr[int32](0),
		// 			CacheMissCount: to.Ptr[int32](0),
		// 			CallCountBlocked: to.Ptr[int32](1),
		// 			CallCountFailed: to.Ptr[int32](0),
		// 			CallCountOther: to.Ptr[int32](0),
		// 			CallCountSuccess: to.Ptr[int32](9),
		// 			CallCountTotal: to.Ptr[int32](10),
		// 			Interval: to.Ptr("PT15M"),
		// 			ServiceTimeAvg: to.Ptr[float64](824.2847111111112),
		// 			ServiceTimeMax: to.Ptr[float64](973.2262000000001),
		// 			ServiceTimeMin: to.Ptr[float64](215.24),
		// 			Timestamp: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-06-03T00:30:00.000Z"); return t}()),
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e436160e64c0f8d7fb20d662be2712f71f0a7ef5/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetReportsByRequest.json
func ExampleReportsClient_NewListByRequestPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewReportsClient().NewListByRequestPager("rg1", "apimService1", "timestamp ge datetime'2017-06-01T00:00:00' and timestamp le datetime'2017-06-04T00:00:00'", &armapimanagement.ReportsClientListByRequestOptions{Top: nil,
		Skip: nil,
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
		// page.RequestReportCollection = armapimanagement.RequestReportCollection{
		// 	Count: to.Ptr[int64](2),
		// 	Value: []*armapimanagement.RequestReportRecordContract{
		// 		{
		// 			OperationID: to.Ptr("/apis/5931a75ae4bbd512a88c680b/operations/-"),
		// 			Method: to.Ptr("GET"),
		// 			APIID: to.Ptr("/apis/5931a75ae4bbd512a88c680b"),
		// 			APIRegion: to.Ptr("East Asia"),
		// 			APITime: to.Ptr[float64](221.1544),
		// 			Cache: to.Ptr("none"),
		// 			IPAddress: to.Ptr("207.xx.155.xx"),
		// 			ProductID: to.Ptr("/products/-"),
		// 			RequestID: to.Ptr("63e7119c-26aa-433c-96d7-f6f3267ff52f"),
		// 			RequestSize: to.Ptr[int32](0),
		// 			ResponseCode: to.Ptr[int32](404),
		// 			ResponseSize: to.Ptr[int32](405),
		// 			ServiceTime: to.Ptr[float64](0),
		// 			SubscriptionID: to.Ptr("/subscriptions/5600b59475ff190048070002"),
		// 			Timestamp: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-06-03T00:17:00.164Z"); return t}()),
		// 			URL: to.Ptr("https://apimService1.azure-api.net/echo/resource?param1=sample"),
		// 			UserID: to.Ptr("/users/1"),
		// 		},
		// 		{
		// 			OperationID: to.Ptr("/apis/5931a75ae4bbd512a88c680b/operations/-"),
		// 			Method: to.Ptr("POST"),
		// 			APIID: to.Ptr("/apis/5931a75ae4bbd512a88c680b"),
		// 			APIRegion: to.Ptr("East Asia"),
		// 			APITime: to.Ptr[float64](6.675400000000001),
		// 			Cache: to.Ptr("none"),
		// 			IPAddress: to.Ptr("207.xx.155.xx"),
		// 			ProductID: to.Ptr("/products/-"),
		// 			RequestID: to.Ptr("e581b7f7-c9ec-4fc6-8ab9-3855d9b00b04"),
		// 			RequestSize: to.Ptr[int32](0),
		// 			ResponseCode: to.Ptr[int32](404),
		// 			ResponseSize: to.Ptr[int32](403),
		// 			ServiceTime: to.Ptr[float64](0),
		// 			SubscriptionID: to.Ptr("/subscriptions/5600b59475ff190048070002"),
		// 			Timestamp: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-06-03T00:17:20.525Z"); return t}()),
		// 			URL: to.Ptr("https://apimService1.azure-api.net/echo/resource"),
		// 			UserID: to.Ptr("/users/1"),
		// 	}},
		// }
	}
}
