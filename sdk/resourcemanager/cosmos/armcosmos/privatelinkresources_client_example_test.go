//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcosmos_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/011ecc5633300a5eefe43dde748f269d39e96458/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/stable/2025-04-15/examples/CosmosDBPrivateLinkResourceListGet.json
func ExamplePrivateLinkResourcesClient_NewListByDatabaseAccountPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcosmos.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewPrivateLinkResourcesClient().NewListByDatabaseAccountPager("rg1", "ddb1", nil)
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
		// page.PrivateLinkResourceListResult = armcosmos.PrivateLinkResourceListResult{
		// 	Value: []*armcosmos.PrivateLinkResource{
		// 		{
		// 			Name: to.Ptr("sql"),
		// 			Type: to.Ptr("Microsoft.DocumentDB/databaseAccounts/privateLinkResources"),
		// 			ID: to.Ptr("subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/Default/providers/Microsoft.DocumentDb/databaseAccounts/ddb1/privateLinkResources/sql"),
		// 			Properties: &armcosmos.PrivateLinkResourceProperties{
		// 				GroupID: to.Ptr("sql"),
		// 				RequiredMembers: []*string{
		// 					to.Ptr("ddb1"),
		// 					to.Ptr("ddb1-westus")},
		// 					RequiredZoneNames: []*string{
		// 						to.Ptr("privatelink.documents.azure.net")},
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/011ecc5633300a5eefe43dde748f269d39e96458/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/stable/2025-04-15/examples/CosmosDBPrivateLinkResourceGet.json
func ExamplePrivateLinkResourcesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armcosmos.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewPrivateLinkResourcesClient().Get(ctx, "rg1", "ddb1", "sql", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.PrivateLinkResource = armcosmos.PrivateLinkResource{
	// 	Name: to.Ptr("sql"),
	// 	Type: to.Ptr("Microsoft.DocumentDB/databaseAccounts/privateLinkResources"),
	// 	ID: to.Ptr("subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/Default/providers/Microsoft.DocumentDb/databaseAccounts/ddb1/privateLinkResources/sql"),
	// 	Properties: &armcosmos.PrivateLinkResourceProperties{
	// 		GroupID: to.Ptr("sql"),
	// 		RequiredMembers: []*string{
	// 			to.Ptr("ddb1"),
	// 			to.Ptr("ddb1-westus")},
	// 			RequiredZoneNames: []*string{
	// 				to.Ptr("privatelink.documents.azure.net")},
	// 			},
	// 		}
}
