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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementListGroupUsers.json
func ExampleGroupUserClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGroupUserClient().NewListPager("rg1", "apimService1", "57d2ef278aa04f0888cba3f3", &armapimanagement.GroupUserClientListOptions{Filter: nil,
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
		// page.UserCollection = armapimanagement.UserCollection{
		// 	Count: to.Ptr[int64](1),
		// 	Value: []*armapimanagement.UserContract{
		// 		{
		// 			Name: to.Ptr("armTemplateUser1"),
		// 			Type: to.Ptr("Microsoft.ApiManagement/service/groups/users"),
		// 			ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/users/kjoshiarmTemplateUser1"),
		// 			Properties: &armapimanagement.UserContractProperties{
		// 				Identities: []*armapimanagement.UserIdentityContract{
		// 					{
		// 						ID: to.Ptr("user1@live.com"),
		// 						Provider: to.Ptr("Basic"),
		// 				}},
		// 				Note: to.Ptr("note for user 1"),
		// 				State: to.Ptr(armapimanagement.UserStateActive),
		// 				Email: to.Ptr("user1@live.com"),
		// 				FirstName: to.Ptr("user1"),
		// 				LastName: to.Ptr("lastname1"),
		// 				RegistrationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-05-31T18:54:41.447Z"); return t}()),
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementHeadGroupUser.json
func ExampleGroupUserClient_CheckEntityExists() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewGroupUserClient().CheckEntityExists(ctx, "rg1", "apimService1", "59306a29e4bbd510dc24e5f9", "5931a75ae4bbd512a88c680b", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementCreateGroupUser.json
func ExampleGroupUserClient_Create() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGroupUserClient().Create(ctx, "rg1", "apimService1", "tempgroup", "59307d350af58404d8a26300", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.UserContract = armapimanagement.UserContract{
	// 	Name: to.Ptr("59307d350af58404d8a26300"),
	// 	Type: to.Ptr("Microsoft.ApiManagement/service/groups/users"),
	// 	ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ApiManagement/service/apimService1/users/59307d350af58404d8a26300"),
	// 	Properties: &armapimanagement.UserContractProperties{
	// 		Identities: []*armapimanagement.UserIdentityContract{
	// 		},
	// 		State: to.Ptr(armapimanagement.UserStateActive),
	// 		Email: to.Ptr("testuser1@live.com"),
	// 		FirstName: to.Ptr("test"),
	// 		Groups: []*armapimanagement.GroupContractProperties{
	// 		},
	// 		LastName: to.Ptr("user"),
	// 		RegistrationDate: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-06-01T20:46:45.437Z"); return t}()),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/3db6867b8e524ea6d1bc7a3bbb989fe50dd2f184/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2024-05-01/examples/ApiManagementDeleteGroupUser.json
func ExampleGroupUserClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armapimanagement.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewGroupUserClient().Delete(ctx, "rg1", "apimService1", "templategroup", "59307d350af58404d8a26300", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
