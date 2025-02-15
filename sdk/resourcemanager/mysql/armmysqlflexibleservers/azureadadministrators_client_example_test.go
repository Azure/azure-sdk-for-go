//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armmysqlflexibleservers_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysqlflexibleservers/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e26b89bcbec9eed5026c01416e481408b2a1ca1a/specification/mysql/resource-manager/Microsoft.DBforMySQL/AAD/stable/2023-12-30/examples/AzureADAdministratorCreate.json
func ExampleAzureADAdministratorsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmysqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAzureADAdministratorsClient().BeginCreateOrUpdate(ctx, "testrg", "mysqltestsvc4", armmysqlflexibleservers.AdministratorNameActiveDirectory, armmysqlflexibleservers.AzureADAdministrator{
		Properties: &armmysqlflexibleservers.AdministratorProperties{
			AdministratorType:  to.Ptr(armmysqlflexibleservers.AdministratorTypeActiveDirectory),
			IdentityResourceID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/test-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-umi"),
			Login:              to.Ptr("bob@contoso.com"),
			Sid:                to.Ptr("c6b82b90-a647-49cb-8a62-0d2d3cb7ac7c"),
			TenantID:           to.Ptr("c12b7025-bfe2-46c1-b463-993b5e4cd467"),
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
	// res.AzureADAdministrator = armmysqlflexibleservers.AzureADAdministrator{
	// 	Name: to.Ptr("ActiveDirectory"),
	// 	Type: to.Ptr("Microsoft.DBforMySQL/flexibleServers/administrators"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforMySQL/flexibleServers/mysqltestsvc4/administrators/ActiveDirectory"),
	// 	Properties: &armmysqlflexibleservers.AdministratorProperties{
	// 		AdministratorType: to.Ptr(armmysqlflexibleservers.AdministratorTypeActiveDirectory),
	// 		IdentityResourceID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/test-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-umi"),
	// 		Login: to.Ptr("bob@contoso.com"),
	// 		Sid: to.Ptr("c6b82b90-a647-49cb-8a62-0d2d3cb7ac7c"),
	// 		TenantID: to.Ptr("c12b7025-bfe2-46c1-b463-993b5e4cd467"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e26b89bcbec9eed5026c01416e481408b2a1ca1a/specification/mysql/resource-manager/Microsoft.DBforMySQL/AAD/stable/2023-12-30/examples/AzureADAdministratorDelete.json
func ExampleAzureADAdministratorsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmysqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewAzureADAdministratorsClient().BeginDelete(ctx, "testrg", "mysqltestsvc4", armmysqlflexibleservers.AdministratorNameActiveDirectory, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e26b89bcbec9eed5026c01416e481408b2a1ca1a/specification/mysql/resource-manager/Microsoft.DBforMySQL/AAD/stable/2023-12-30/examples/AzureADAdministratorGet.json
func ExampleAzureADAdministratorsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmysqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewAzureADAdministratorsClient().Get(ctx, "testrg", "mysqltestsvc4", armmysqlflexibleservers.AdministratorNameActiveDirectory, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.AzureADAdministrator = armmysqlflexibleservers.AzureADAdministrator{
	// 	Name: to.Ptr("ActiveDirectory"),
	// 	Type: to.Ptr("Microsoft.DBforMySQL/flexibleServers/administrators"),
	// 	ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforMySQL/flexibleServers/mysqltestsvc4/administrators/ActiveDirectory"),
	// 	Properties: &armmysqlflexibleservers.AdministratorProperties{
	// 		AdministratorType: to.Ptr(armmysqlflexibleservers.AdministratorTypeActiveDirectory),
	// 		IdentityResourceID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/test-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-umi"),
	// 		Login: to.Ptr("bob@contoso.com"),
	// 		Sid: to.Ptr("c6b82b90-a647-49cb-8a62-0d2d3cb7ac7c"),
	// 		TenantID: to.Ptr("c12b7025-bfe2-46c1-b463-993b5e4cd467"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e26b89bcbec9eed5026c01416e481408b2a1ca1a/specification/mysql/resource-manager/Microsoft.DBforMySQL/AAD/stable/2023-12-30/examples/AzureADAdministratorsListByServer.json
func ExampleAzureADAdministratorsClient_NewListByServerPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmysqlflexibleservers.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewAzureADAdministratorsClient().NewListByServerPager("testrg", "mysqltestsvc4", nil)
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
		// page.AdministratorListResult = armmysqlflexibleservers.AdministratorListResult{
		// 	Value: []*armmysqlflexibleservers.AzureADAdministrator{
		// 		{
		// 			Name: to.Ptr("ActiveDirectory"),
		// 			Type: to.Ptr("Microsoft.DBforMySQL/flexibleServers/administrators"),
		// 			ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/testrg/providers/Microsoft.DBforMySQL/flexibleServers/mysqltestsvc4/administrators/ActiveDirectory"),
		// 			Properties: &armmysqlflexibleservers.AdministratorProperties{
		// 				AdministratorType: to.Ptr(armmysqlflexibleservers.AdministratorTypeActiveDirectory),
		// 				IdentityResourceID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/test-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-umi"),
		// 				Login: to.Ptr("bob@contoso.com"),
		// 				Sid: to.Ptr("c6b82b90-a647-49cb-8a62-0d2d3cb7ac7c"),
		// 				TenantID: to.Ptr("c12b7025-bfe2-46c1-b463-993b5e4cd467"),
		// 			},
		// 	}},
		// }
	}
}
