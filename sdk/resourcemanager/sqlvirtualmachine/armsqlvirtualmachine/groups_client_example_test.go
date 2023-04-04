//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsqlvirtualmachine_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sqlvirtualmachine/armsqlvirtualmachine"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/sqlvirtualmachine/resource-manager/Microsoft.SqlVirtualMachine/preview/2022-07-01-preview/examples/GetSqlVirtualMachineGroup.json
func ExampleGroupsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsqlvirtualmachine.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewGroupsClient().Get(ctx, "testrg", "testvmgroup", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.Group = armsqlvirtualmachine.Group{
	// 	Name: to.Ptr("testvmgroup"),
	// 	Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
	// 	ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/testvmgroup"),
	// 	Location: to.Ptr("northeurope"),
	// 	Tags: map[string]*string{
	// 		"mytag": to.Ptr("myval"),
	// 	},
	// 	Properties: &armsqlvirtualmachine.GroupProperties{
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		SQLImageOffer: to.Ptr("SQL2016-WS2016"),
	// 		SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
	// 		WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
	// 			ClusterBootstrapAccount: to.Ptr("testrpadmin"),
	// 			ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
	// 			DomainFqdn: to.Ptr("testdomain.com"),
	// 			OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
	// 			SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
	// 			StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/sqlvirtualmachine/resource-manager/Microsoft.SqlVirtualMachine/preview/2022-07-01-preview/examples/CreateOrUpdateSqlVirtualMachineGroup.json
func ExampleGroupsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsqlvirtualmachine.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGroupsClient().BeginCreateOrUpdate(ctx, "testrg", "testvmgroup", armsqlvirtualmachine.Group{
		Location: to.Ptr("northeurope"),
		Tags: map[string]*string{
			"mytag": to.Ptr("myval"),
		},
		Properties: &armsqlvirtualmachine.GroupProperties{
			SQLImageOffer: to.Ptr("SQL2016-WS2016"),
			SQLImageSKU:   to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
			WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
				ClusterBootstrapAccount:  to.Ptr("testrpadmin"),
				ClusterOperatorAccount:   to.Ptr("testrp@testdomain.com"),
				ClusterSubnetType:        to.Ptr(armsqlvirtualmachine.ClusterSubnetTypeMultiSubnet),
				DomainFqdn:               to.Ptr("testdomain.com"),
				OuPath:                   to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
				SQLServiceAccount:        to.Ptr("sqlservice@testdomain.com"),
				StorageAccountPrimaryKey: to.Ptr("<primary storage access key>"),
				StorageAccountURL:        to.Ptr("https://storgact.blob.core.windows.net/"),
			},
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
	// res.Group = armsqlvirtualmachine.Group{
	// 	Name: to.Ptr("testvmgroup"),
	// 	Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
	// 	ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/testvmgroup"),
	// 	Location: to.Ptr("northeurope"),
	// 	Tags: map[string]*string{
	// 		"mytag": to.Ptr("myval"),
	// 	},
	// 	Properties: &armsqlvirtualmachine.GroupProperties{
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		SQLImageOffer: to.Ptr("SQL2016-WS2016"),
	// 		SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
	// 		WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
	// 			ClusterBootstrapAccount: to.Ptr("testrpadmin"),
	// 			ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
	// 			ClusterSubnetType: to.Ptr(armsqlvirtualmachine.ClusterSubnetTypeMultiSubnet),
	// 			DomainFqdn: to.Ptr("testdomain.com"),
	// 			OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
	// 			SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
	// 			StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/sqlvirtualmachine/resource-manager/Microsoft.SqlVirtualMachine/preview/2022-07-01-preview/examples/DeleteSqlVirtualMachineGroup.json
func ExampleGroupsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsqlvirtualmachine.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGroupsClient().BeginDelete(ctx, "testrg", "testvmgroup", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/sqlvirtualmachine/resource-manager/Microsoft.SqlVirtualMachine/preview/2022-07-01-preview/examples/UpdateSqlVirtualMachineGroup.json
func ExampleGroupsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsqlvirtualmachine.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewGroupsClient().BeginUpdate(ctx, "testrg", "testvmgroup", armsqlvirtualmachine.GroupUpdate{
		Tags: map[string]*string{
			"mytag": to.Ptr("myval"),
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
	// res.Group = armsqlvirtualmachine.Group{
	// 	Name: to.Ptr("testvmgroup"),
	// 	Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
	// 	ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/testvm"),
	// 	Location: to.Ptr("northeurope"),
	// 	Tags: map[string]*string{
	// 		"mytag": to.Ptr("myval"),
	// 	},
	// 	Properties: &armsqlvirtualmachine.GroupProperties{
	// 		ProvisioningState: to.Ptr("Succeeded"),
	// 		SQLImageOffer: to.Ptr("SQL2017-WS2016"),
	// 		SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
	// 		WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
	// 			StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/sqlvirtualmachine/resource-manager/Microsoft.SqlVirtualMachine/preview/2022-07-01-preview/examples/ListByResourceGroupSqlVirtualMachineGroup.json
func ExampleGroupsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsqlvirtualmachine.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGroupsClient().NewListByResourceGroupPager("testrg", nil)
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
		// page.GroupListResult = armsqlvirtualmachine.GroupListResult{
		// 	Value: []*armsqlvirtualmachine.Group{
		// 		{
		// 			Name: to.Ptr("testvmgroup"),
		// 			Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
		// 			ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/testvmgroup"),
		// 			Location: to.Ptr("northeurope"),
		// 			Tags: map[string]*string{
		// 				"mytag": to.Ptr("myval"),
		// 			},
		// 			Properties: &armsqlvirtualmachine.GroupProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				SQLImageOffer: to.Ptr("SQL2017-WS2016"),
		// 				SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
		// 				WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
		// 					ClusterBootstrapAccount: to.Ptr("testrpadmin"),
		// 					ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
		// 					DomainFqdn: to.Ptr("testdomain.com"),
		// 					OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
		// 					SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
		// 					StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("testvmgroup1"),
		// 			Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
		// 			ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/testvmgroup1"),
		// 			Location: to.Ptr("northeurope"),
		// 			Tags: map[string]*string{
		// 				"mytag": to.Ptr("myval"),
		// 			},
		// 			Properties: &armsqlvirtualmachine.GroupProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				SQLImageOffer: to.Ptr("SQL2016-WS2016"),
		// 				SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
		// 				WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
		// 					ClusterBootstrapAccount: to.Ptr("testrpadmin"),
		// 					ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
		// 					DomainFqdn: to.Ptr("testdomain.com"),
		// 					OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
		// 					SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
		// 					StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("testvmgroup2"),
		// 			Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
		// 			ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/testvmgroup2"),
		// 			Location: to.Ptr("northeurope"),
		// 			Tags: map[string]*string{
		// 				"mytag": to.Ptr("myval"),
		// 			},
		// 			Properties: &armsqlvirtualmachine.GroupProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				SQLImageOffer: to.Ptr("SQL2016-WS2016"),
		// 				SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
		// 				WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
		// 					ClusterBootstrapAccount: to.Ptr("testrpadmin"),
		// 					ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
		// 					DomainFqdn: to.Ptr("testdomain.com"),
		// 					OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
		// 					SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
		// 					StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
		// 				},
		// 			},
		// 	}},
		// }
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/sqlvirtualmachine/resource-manager/Microsoft.SqlVirtualMachine/preview/2022-07-01-preview/examples/ListSubscriptionSqlVirtualMachineGroup.json
func ExampleGroupsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsqlvirtualmachine.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewGroupsClient().NewListPager(nil)
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
		// page.GroupListResult = armsqlvirtualmachine.GroupListResult{
		// 	Value: []*armsqlvirtualmachine.Group{
		// 		{
		// 			Name: to.Ptr("testvmgroup"),
		// 			Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
		// 			ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/"),
		// 			Location: to.Ptr("northeurope"),
		// 			Tags: map[string]*string{
		// 				"mytag": to.Ptr("myval"),
		// 			},
		// 			Properties: &armsqlvirtualmachine.GroupProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				SQLImageOffer: to.Ptr("SQL2017-WS2016"),
		// 				SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
		// 				WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
		// 					ClusterBootstrapAccount: to.Ptr("testrpadmin"),
		// 					ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
		// 					ClusterSubnetType: to.Ptr(armsqlvirtualmachine.ClusterSubnetTypeMultiSubnet),
		// 					DomainFqdn: to.Ptr("testdomain.com"),
		// 					OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
		// 					SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
		// 					StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("testvmgroup1"),
		// 			Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
		// 			ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/"),
		// 			Location: to.Ptr("northeurope"),
		// 			Tags: map[string]*string{
		// 				"mytag": to.Ptr("myval"),
		// 			},
		// 			Properties: &armsqlvirtualmachine.GroupProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				SQLImageOffer: to.Ptr("SQL2016-WS2016"),
		// 				SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
		// 				WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
		// 					ClusterBootstrapAccount: to.Ptr("testrpadmin"),
		// 					ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
		// 					ClusterSubnetType: to.Ptr(armsqlvirtualmachine.ClusterSubnetTypeMultiSubnet),
		// 					DomainFqdn: to.Ptr("testdomain.com"),
		// 					OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
		// 					SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
		// 					StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
		// 				},
		// 			},
		// 		},
		// 		{
		// 			Name: to.Ptr("testvmgroup2"),
		// 			Type: to.Ptr("Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups"),
		// 			ID: to.Ptr("/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/testrg2/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/"),
		// 			Location: to.Ptr("northeurope"),
		// 			Tags: map[string]*string{
		// 				"mytag": to.Ptr("myval"),
		// 			},
		// 			Properties: &armsqlvirtualmachine.GroupProperties{
		// 				ProvisioningState: to.Ptr("Succeeded"),
		// 				SQLImageOffer: to.Ptr("SQL2016-WS2016"),
		// 				SQLImageSKU: to.Ptr(armsqlvirtualmachine.SQLVMGroupImageSKUEnterprise),
		// 				WsfcDomainProfile: &armsqlvirtualmachine.WsfcDomainProfile{
		// 					ClusterBootstrapAccount: to.Ptr("testrpadmin"),
		// 					ClusterOperatorAccount: to.Ptr("testrp@testdomain.com"),
		// 					ClusterSubnetType: to.Ptr(armsqlvirtualmachine.ClusterSubnetTypeMultiSubnet),
		// 					DomainFqdn: to.Ptr("testdomain.com"),
		// 					OuPath: to.Ptr("OU=WSCluster,DC=testdomain,DC=com"),
		// 					SQLServiceAccount: to.Ptr("sqlservice@testdomain.com"),
		// 					StorageAccountURL: to.Ptr("https://storgact.blob.core.windows.net/"),
		// 				},
		// 			},
		// 	}},
		// }
	}
}
