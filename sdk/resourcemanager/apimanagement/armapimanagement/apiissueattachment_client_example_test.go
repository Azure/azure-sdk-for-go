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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListApiIssueAttachments.json
func ExampleAPIIssueAttachmentClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAPIIssueAttachmentClient().NewListByServicePager("rg1", "apimService1", "57d1f7558aa04f15146d9d8a", "57d2ef278aa04f0ad01d6cdc", &armapimanagement.APIIssueAttachmentClientListByServiceOptions{Filter: nil,
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
		// page.IssueAttachmentCollection = armapimanagement.IssueAttachmentCollection{
		// 	Count: to.Ptr[int64](1),
		// 	Value: []*armapimanagement.IssueAttachmentContract{
		// 		{
		// 			Name: to.Ptr("57d2ef278aa04f0888cba3f3"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/apis/issues/attachments"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d1f7558aa04f15146d9d8a/issues/57d2ef278aa04f0ad01d6cdc/attachments/57d2ef278aa04f0888cba3f3"),
		// 			Properties: &armapimanagement.IssueAttachmentContractProperties{
		// 				Content: to.Ptr("https://.../image.jpg"),
		// 				ContentFormat: to.Ptr("link"),
		// 				Title: to.Ptr("Issue attachment."),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementHeadApiIssueAttachment.json
func ExampleAPIIssueAttachmentClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAPIIssueAttachmentClient().GetEntityTag(ctx, "rg1", "apimService1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", "57d2ef278aa04f0888cba3f3", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementGetApiIssueAttachment.json
func ExampleAPIIssueAttachmentClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIIssueAttachmentClient().Get(ctx, "rg1", "apimService1", "57d2ef278aa04f0888cba3f3", "57d2ef278aa04f0ad01d6cdc", "57d2ef278aa04f0888cba3f3", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IssueAttachmentContract = armapimanagement.IssueAttachmentContract{
	// 	Name: to.Ptr("57d2ef278aa04f0888cba3f3"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/issues/attachments"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d1f7558aa04f15146d9d8a/issues/57d2ef278aa04f0ad01d6cdc/attachments/57d2ef278aa04f0888cba3f3"),
	// 	Properties: &armapimanagement.IssueAttachmentContractProperties{
	// 		Content: to.Ptr("https://.../image.jpg"),
	// 		ContentFormat: to.Ptr("link"),
	// 		Title: to.Ptr("Issue attachment."),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateApiIssueAttachment.json
func ExampleAPIIssueAttachmentClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAPIIssueAttachmentClient().CreateOrUpdate(ctx, "rg1", "apimService1", "57d1f7558aa04f15146d9d8a", "57d2ef278aa04f0ad01d6cdc", "57d2ef278aa04f0888cba3f3", armapimanagement.IssueAttachmentContract{
		Properties: &armapimanagement.IssueAttachmentContractProperties{
			Content:       to.Ptr("IEJhc2U2NA=="),
			ContentFormat: to.Ptr("image/jpeg"),
			Title:         to.Ptr("Issue attachment."),
		},
	}, &armapimanagement.APIIssueAttachmentClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.IssueAttachmentContract = armapimanagement.IssueAttachmentContract{
	// 	Name: to.Ptr("57d2ef278aa04f0888cba3f3"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/apis/issues/attachments"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/apis/57d1f7558aa04f15146d9d8a/issues/57d2ef278aa04f0ad01d6cdc/attachments/57d2ef278aa04f0888cba3f3"),
	// 	Properties: &armapimanagement.IssueAttachmentContractProperties{
	// 		Content: to.Ptr("https://.../image.jpg"),
	// 		ContentFormat: to.Ptr("link"),
	// 		Title: to.Ptr("Issue attachment."),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeleteApiIssueAttachment.json
func ExampleAPIIssueAttachmentClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewAPIIssueAttachmentClient().Delete(ctx, "rg1", "apimService1", "57d1f7558aa04f15146d9d8a", "57d2ef278aa04f0ad01d6cdc", "57d2ef278aa04f0888cba3f3", "*", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
