//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdnsresolver_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dnsresolver/armdnsresolver"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/dnsresolver/resource-manager/Microsoft.Network/preview/2023-07-01-preview/examples/DnsResolverDomainList_Put.json
func ExampleDomainListsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdnsresolver.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDomainListsClient().BeginCreateOrUpdate(ctx, "sampleResourceGroup", "sampleDnsResolverDomainList", armdnsresolver.DomainList{
		Location: to.Ptr("westus2"),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armdnsresolver.DomainListProperties{
			Domains: []*string{
				to.Ptr("contoso.com")},
		},
	}, &armdnsresolver.DomainListsClientBeginCreateOrUpdateOptions{IfMatch: nil,
		IfNoneMatch: nil,
	})
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
	// res.DomainList = armdnsresolver.DomainList{
	// 	Name: to.Ptr("sampleDnsResolverDomainList"),
	// 	Type: to.Ptr("Microsoft.Network/dnsResolverDomainLists"),
	// 	ID: to.Ptr("/subscriptions/abdd4249-9f34-4cc6-8e42-c2e32110603e/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/dnsResolverDomainLists/sampleDnsResolverDomainList"),
	// 	SystemData: &armdnsresolver.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-01T01:01:01.107Z"); return t}()),
	// 		CreatedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-02T02:03:01.197Z"); return t}()),
	// 		LastModifiedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
	// 	},
	// 	Location: to.Ptr("westus2"),
	// 	Tags: map[string]*string{
	// 		"key1": to.Ptr("value1"),
	// 	},
	// 	Etag: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	Properties: &armdnsresolver.DomainListProperties{
	// 		Domains: []*string{
	// 			to.Ptr("contoso.com")},
	// 			ProvisioningState: to.Ptr(armdnsresolver.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("b6b2d964-8588-4e3a-a7fe-8a5b7fe8eca5"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/dnsresolver/resource-manager/Microsoft.Network/preview/2023-07-01-preview/examples/DnsResolverDomainList_Patch.json
func ExampleDomainListsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdnsresolver.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDomainListsClient().BeginUpdate(ctx, "sampleResourceGroup", "sampleDnsResolverDomainList", armdnsresolver.DomainListPatch{
		Properties: &armdnsresolver.DomainListPatchProperties{
			Domains: []*string{
				to.Ptr("contoso.com")},
		},
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
	}, &armdnsresolver.DomainListsClientBeginUpdateOptions{IfMatch: nil})
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
	// res.DomainList = armdnsresolver.DomainList{
	// 	Name: to.Ptr("sampleDnsResolverDomainList"),
	// 	Type: to.Ptr("Microsoft.Network/dnsResolverDomainLists"),
	// 	ID: to.Ptr("/subscriptions/abdd4249-9f34-4cc6-8e42-c2e32110603e/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/dnsResolverDomainLists/sampleDnsResolverDomainList"),
	// 	SystemData: &armdnsresolver.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-01T01:01:01.107Z"); return t}()),
	// 		CreatedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-02T02:03:01.197Z"); return t}()),
	// 		LastModifiedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
	// 	},
	// 	Location: to.Ptr("westus2"),
	// 	Tags: map[string]*string{
	// 		"key1": to.Ptr("value1"),
	// 	},
	// 	Etag: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	Properties: &armdnsresolver.DomainListProperties{
	// 		Domains: []*string{
	// 			to.Ptr("contoso.com")},
	// 			ProvisioningState: to.Ptr(armdnsresolver.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("b6b2d964-8588-4e3a-a7fe-8a5b7fe8eca5"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/dnsresolver/resource-manager/Microsoft.Network/preview/2023-07-01-preview/examples/DnsResolverDomainList_Delete.json
func ExampleDomainListsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdnsresolver.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewDomainListsClient().BeginDelete(ctx, "sampleResourceGroup", "sampleDnsResolverDomainList", &armdnsresolver.DomainListsClientBeginDeleteOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/dnsresolver/resource-manager/Microsoft.Network/preview/2023-07-01-preview/examples/DnsResolverDomainList_Get.json
func ExampleDomainListsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdnsresolver.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewDomainListsClient().Get(ctx, "sampleResourceGroup", "sampleDnsResolverDomainList", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.DomainList = armdnsresolver.DomainList{
	// 	Name: to.Ptr("sampleDnsResolverDomainList"),
	// 	Type: to.Ptr("Microsoft.Network/dnsResolverDomainLists"),
	// 	ID: to.Ptr("/subscriptions/abdd4249-9f34-4cc6-8e42-c2e32110603e/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/dnsResolverDomainLists/sampleDnsResolverDomainList"),
	// 	SystemData: &armdnsresolver.SystemData{
	// 		CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-03T01:01:01.107Z"); return t}()),
	// 		CreatedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
	// 		LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-04T02:03:01.197Z"); return t}()),
	// 		LastModifiedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
	// 	},
	// 	Location: to.Ptr("westus2"),
	// 	Tags: map[string]*string{
	// 		"key1": to.Ptr("value1"),
	// 	},
	// 	Etag: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	Properties: &armdnsresolver.DomainListProperties{
	// 		Domains: []*string{
	// 			to.Ptr("contoso.com")},
	// 			ProvisioningState: to.Ptr(armdnsresolver.ProvisioningStateSucceeded),
	// 			ResourceGUID: to.Ptr("a7e1a32c-498c-401c-a805-5bc3518257b8"),
	// 		},
	// 	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/dnsresolver/resource-manager/Microsoft.Network/preview/2023-07-01-preview/examples/DnsResolverDomainList_ListByResourceGroup.json
func ExampleDomainListsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdnsresolver.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDomainListsClient().NewListByResourceGroupPager("sampleResourceGroup", &armdnsresolver.DomainListsClientListByResourceGroupOptions{Top: nil})
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
		// page.DomainListResult = armdnsresolver.DomainListResult{
		// 	Value: []*armdnsresolver.DomainList{
		// 		{
		// 			Name: to.Ptr("sampleDnsResolverDomainList1"),
		// 			Type: to.Ptr("Microsoft.Network/dnsResolverDomainLists"),
		// 			ID: to.Ptr("/subscriptions/abdd4249-9f34-4cc6-8e42-c2e32110603e/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/dnsResolverDomainLists/sampleDnsResolverDomainList1"),
		// 			SystemData: &armdnsresolver.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-01T01:01:01.107Z"); return t}()),
		// 				CreatedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-02T02:03:01.197Z"); return t}()),
		// 				LastModifiedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 			},
		// 			Location: to.Ptr("westus2"),
		// 			Tags: map[string]*string{
		// 				"key1": to.Ptr("value1"),
		// 			},
		// 			Etag: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 			Properties: &armdnsresolver.DomainListProperties{
		// 				Domains: []*string{
		// 					to.Ptr("contoso.com")},
		// 					ProvisioningState: to.Ptr(armdnsresolver.ProvisioningStateSucceeded),
		// 					ResourceGUID: to.Ptr("ad9c8da4-3bb2-4821-a878-c2cb07b01fb6"),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("sampleDnsResolverDomainList2"),
		// 				Type: to.Ptr("Microsoft.Network/dnsResolverDomainLists"),
		// 				ID: to.Ptr("/subscriptions/abdd4249-9f34-4cc6-8e42-c2e32110603e/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/dnsResolverDomainLists/sampleDnsResolverDomainList2"),
		// 				SystemData: &armdnsresolver.SystemData{
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-03T01:01:01.107Z"); return t}()),
		// 					CreatedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-04T02:03:01.197Z"); return t}()),
		// 					LastModifiedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 				},
		// 				Location: to.Ptr("westus2"),
		// 				Tags: map[string]*string{
		// 					"key1": to.Ptr("value1"),
		// 				},
		// 				Etag: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				Properties: &armdnsresolver.DomainListProperties{
		// 					Domains: []*string{
		// 						to.Ptr("contoso.com")},
		// 						ProvisioningState: to.Ptr(armdnsresolver.ProvisioningStateSucceeded),
		// 						ResourceGUID: to.Ptr("b6b2d964-8588-4e3a-a7fe-8a5b7fe8eca5"),
		// 					},
		// 			}},
		// 		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/366aaa13cdd218b9adac716680e49473673410c8/specification/dnsresolver/resource-manager/Microsoft.Network/preview/2023-07-01-preview/examples/DnsResolverDomainList_ListBySubscription.json
func ExampleDomainListsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armdnsresolver.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewDomainListsClient().NewListPager(&armdnsresolver.DomainListsClientListOptions{Top: nil})
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
		// page.DomainListResult = armdnsresolver.DomainListResult{
		// 	Value: []*armdnsresolver.DomainList{
		// 		{
		// 			Name: to.Ptr("sampleDnsResolverDomainList1"),
		// 			Type: to.Ptr("Microsoft.Network/dnsResolverDomainLists"),
		// 			ID: to.Ptr("/subscriptions/abdd4249-9f34-4cc6-8e42-c2e32110603e/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/dnsResolverDomainLists/sampleDnsResolverDomainList1"),
		// 			SystemData: &armdnsresolver.SystemData{
		// 				CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-01T01:01:01.107Z"); return t}()),
		// 				CreatedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 				LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-02T02:03:01.197Z"); return t}()),
		// 				LastModifiedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 			},
		// 			Location: to.Ptr("westus2"),
		// 			Tags: map[string]*string{
		// 				"key1": to.Ptr("value1"),
		// 			},
		// 			Etag: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 			Properties: &armdnsresolver.DomainListProperties{
		// 				Domains: []*string{
		// 					to.Ptr("contoso.com")},
		// 					ProvisioningState: to.Ptr(armdnsresolver.ProvisioningStateSucceeded),
		// 					ResourceGUID: to.Ptr("ad9c8da4-3bb2-4821-a878-c2cb07b01fb6"),
		// 				},
		// 			},
		// 			{
		// 				Name: to.Ptr("sampleDnsResolverDomainList2"),
		// 				Type: to.Ptr("Microsoft.Network/dnsResolverDomainLists"),
		// 				ID: to.Ptr("/subscriptions/abdd4249-9f34-4cc6-8e42-c2e32110603e/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/dnsResolverDomainLists/sampleDnsResolverDomainList2"),
		// 				SystemData: &armdnsresolver.SystemData{
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-01T01:01:01.107Z"); return t}()),
		// 					CreatedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2021-04-02T02:03:01.197Z"); return t}()),
		// 					LastModifiedByType: to.Ptr(armdnsresolver.CreatedByTypeApplication),
		// 				},
		// 				Location: to.Ptr("westus2"),
		// 				Tags: map[string]*string{
		// 					"key1": to.Ptr("value1"),
		// 				},
		// 				Etag: to.Ptr("00000000-0000-0000-0000-000000000000"),
		// 				Properties: &armdnsresolver.DomainListProperties{
		// 					Domains: []*string{
		// 						to.Ptr("contoso.com")},
		// 						ProvisioningState: to.Ptr(armdnsresolver.ProvisioningStateSucceeded),
		// 						ResourceGUID: to.Ptr("b6b2d964-8588-4e3a-a7fe-8a5b7fe8eca5"),
		// 					},
		// 			}},
		// 		}
	}
}
