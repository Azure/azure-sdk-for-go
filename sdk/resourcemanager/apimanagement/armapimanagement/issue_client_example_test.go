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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListIssues.json
func ExampleIssueClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewIssueClient().NewListByServicePager("rg1", "apimService1", &armapimanagement.IssueClientListByServiceOptions{Filter: nil,
		Top:  nil,
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
		// page.IssueCollection = armapimanagement.IssueCollection{
		// 	Count: to.Ptr[int64](1),
		// 	Value: []*armapimanagement.IssueContract{
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0ad01d6cdc"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/issues"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/issues/57d2ef278aa04f0ad01d6cdc"),
		// 			Properties: &armapimanagement.IssueContractProperties{
		// 				APIID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d1f7558aa04f15146d9d8a"),
		// 				CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-02-01T22:21:20.467Z"); return t}()),
		// 				State: to.Ptr(armapimanagement.StateOpen),
		// 				Description: to.Ptr("New API issue description"),
		// 				Title: to.Ptr("New API issue"),
		// 				UserID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/users/1"),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetIssue.json
func ExampleIssueClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewIssueClient().Get(ctx, "rg1", "apimService1", "57d2ef278aa04f0ad01d6cdc", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IssueContract = armapimanagement.IssueContract{
	// 	Name: to.Ptr("57d2ef278aa04f0ad01d6cdc"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/issues"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/issues/57d2ef278aa04f0ad01d6cdc"),
	// 	Properties: &armapimanagement.IssueContractProperties{
	// 		APIID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d1f7558aa04f15146d9d8a"),
	// 		CreatedDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2018-02-01T22:21:20.467Z"); return t}()),
	// 		State: to.Ptr(armapimanagement.StateOpen),
	// 		Description: to.Ptr("New API issue description"),
	// 		Title: to.Ptr("New API issue"),
	// 		UserID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/users/1"),
	// 	},
	// }
}
