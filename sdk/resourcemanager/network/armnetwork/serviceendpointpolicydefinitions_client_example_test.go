//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetwork_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/f4c6c8697c59f966db0d1e36b62df3af3bca9065/specification/network/resource-manager/Microsoft.Network/stable/2023-11-01/examples/ServiceEndpointPolicyDefinitionDelete.json
func ExampleServiceEndpointPolicyDefinitionsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServiceEndpointPolicyDefinitionsClient().BeginDelete(ctx, "rg1", "testPolicy", "testDefinition", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/f4c6c8697c59f966db0d1e36b62df3af3bca9065/specification/network/resource-manager/Microsoft.Network/stable/2023-11-01/examples/ServiceEndpointPolicyDefinitionGet.json
func ExampleServiceEndpointPolicyDefinitionsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewServiceEndpointPolicyDefinitionsClient().Get(ctx, "rg1", "testPolicy", "testDefinition", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.ServiceEndpointPolicyDefinition = armnetwork.ServiceEndpointPolicyDefinition{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/serviceEndpointPolicies/testPolicy/serviceEndpointPolicyDefinitions/testDefinition"),
	// 	Name: to.Ptr("testDefinition"),
	// 	Properties: &armnetwork.ServiceEndpointPolicyDefinitionPropertiesFormat{
	// 		Description: to.Ptr("Storage Service EndpointPolicy Definition"),
	// 		Service: to.Ptr("Microsoft.Storage"),
	// 		ServiceResources: []*string{
	// 			to.Ptr("/subscriptions/subid1"),
	// 			to.Ptr("/subscriptions/subid1/resourceGroups/storageRg"),
	// 			to.Ptr("/subscriptions/subid1/resourceGroups/storageRg/providers/Microsoft.Storage/storageAccounts/stAccount")},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/f4c6c8697c59f966db0d1e36b62df3af3bca9065/specification/network/resource-manager/Microsoft.Network/stable/2023-11-01/examples/ServiceEndpointPolicyDefinitionCreate.json
func ExampleServiceEndpointPolicyDefinitionsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewServiceEndpointPolicyDefinitionsClient().BeginCreateOrUpdate(ctx, "rg1", "testPolicy", "testDefinition", armnetwork.ServiceEndpointPolicyDefinition{
		Properties: &armnetwork.ServiceEndpointPolicyDefinitionPropertiesFormat{
			Description: to.Ptr("Storage Service EndpointPolicy Definition"),
			Service:     to.Ptr("Microsoft.Storage"),
			ServiceResources: []*string{
				to.Ptr("/subscriptions/subid1"),
				to.Ptr("/subscriptions/subid1/resourceGroups/storageRg"),
				to.Ptr("/subscriptions/subid1/resourceGroups/storageRg/providers/Microsoft.Storage/storageAccounts/stAccount")},
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
	// res.ServiceEndpointPolicyDefinition = armnetwork.ServiceEndpointPolicyDefinition{
	// 	ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/serviceEndpointPolicies/testPolicy/serviceEndpointPolicyDefinitions/testDefinition"),
	// 	Name: to.Ptr("testDefinition"),
	// 	Properties: &armnetwork.ServiceEndpointPolicyDefinitionPropertiesFormat{
	// 		Description: to.Ptr("Storage Service EndpointPolicy Definition"),
	// 		Service: to.Ptr("Microsoft.Storage"),
	// 		ServiceResources: []*string{
	// 			to.Ptr("/subscriptions/subid1"),
	// 			to.Ptr("/subscriptions/subid1/resourceGroups/storageRg"),
	// 			to.Ptr("/subscriptions/subid1/resourceGroups/storageRg/providers/Microsoft.Storage/storageAccounts/stAccount")},
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/f4c6c8697c59f966db0d1e36b62df3af3bca9065/specification/network/resource-manager/Microsoft.Network/stable/2023-11-01/examples/ServiceEndpointPolicyDefinitionList.json
func ExampleServiceEndpointPolicyDefinitionsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armnetwork.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewServiceEndpointPolicyDefinitionsClient().NewListByResourceGroupPager("rg1", "testPolicy", nil)
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
		// page.ServiceEndpointPolicyDefinitionListResult = armnetwork.ServiceEndpointPolicyDefinitionListResult{
		// 	Value: []*armnetwork.ServiceEndpointPolicyDefinition{
		// 		{
		// 			ID: to.Ptr("/subscriptions/subid/resourceGroups/rg1/providers/Microsoft.Network/serviceEndpointPolicies/testPolicy/serviceEndpointPolicyDefinitions/testDef"),
		// 			Name: to.Ptr("testDef"),
		// 			Properties: &armnetwork.ServiceEndpointPolicyDefinitionPropertiesFormat{
		// 				Description: to.Ptr("Storage Service EndpointPolicy Definition"),
		// 				Service: to.Ptr("Microsoft.Storage"),
		// 				ServiceResources: []*string{
		// 					to.Ptr("/subscriptions/subid1"),
		// 					to.Ptr("/subscriptions/subid1/resourceGroups/storageRg"),
		// 					to.Ptr("/subscriptions/subid1/resourceGroups/storageRg/providers/Microsoft.Storage/storageAccounts/stAccount")},
		// 				},
		// 		}},
		// 	}
	}
}
