// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armmongocluster_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mongocluster/armmongocluster"
	"log"
)

// Generated from example definition: /mnt/vss/_work/1/s/eng/azure-sdk-for-go/sdk/resourcemanager/mongocluster/armmongocluster/TempTypeSpecFiles/DocumentDB.MongoCluster.Management/examples/2024-07-01/MongoClusters_FirewallRuleCreate.json
func ExampleFirewallRulesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmongocluster.NewClientFactory("ffffffff-ffff-ffff-ffff-ffffffffffff", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFirewallRulesClient().BeginCreateOrUpdate(ctx, "TestGroup", "myMongoCluster", "rule1", armmongocluster.FirewallRule{
		Properties: &armmongocluster.FirewallRuleProperties{
			StartIPAddress: to.Ptr("0.0.0.0"),
			EndIPAddress:   to.Ptr("255.255.255.255"),
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
	// res = armmongocluster.FirewallRulesClientCreateOrUpdateResponse{
	// 	FirewallRule: &armmongocluster.FirewallRule{
	// 		ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/TestResourceGroup/providers/Microsoft.DocumentDB/mongoClusters/myMongoCluster/firewallRules/rule1"),
	// 		Name: to.Ptr("rule1"),
	// 		Type: to.Ptr("/Microsoft.DocumentDB/mongoClusters/firewallRules"),
	// 		SystemData: &armmongocluster.SystemData{
	// 			CreatedBy: to.Ptr("user1"),
	// 			CreatedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-01T17:18:19.1234567Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("user2"),
	// 			LastModifiedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-02T17:18:19.1234567Z"); return t}()),
	// 		},
	// 		Properties: &armmongocluster.FirewallRuleProperties{
	// 			ProvisioningState: to.Ptr(armmongocluster.ProvisioningStateSucceeded),
	// 			StartIPAddress: to.Ptr("0.0.0.0"),
	// 			EndIPAddress: to.Ptr("255.255.255.255"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: /mnt/vss/_work/1/s/eng/azure-sdk-for-go/sdk/resourcemanager/mongocluster/armmongocluster/TempTypeSpecFiles/DocumentDB.MongoCluster.Management/examples/2024-07-01/MongoClusters_FirewallRuleDelete.json
func ExampleFirewallRulesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmongocluster.NewClientFactory("ffffffff-ffff-ffff-ffff-ffffffffffff", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewFirewallRulesClient().BeginDelete(ctx, "TestGroup", "myMongoCluster", "rule1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: /mnt/vss/_work/1/s/eng/azure-sdk-for-go/sdk/resourcemanager/mongocluster/armmongocluster/TempTypeSpecFiles/DocumentDB.MongoCluster.Management/examples/2024-07-01/MongoClusters_FirewallRuleGet.json
func ExampleFirewallRulesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmongocluster.NewClientFactory("ffffffff-ffff-ffff-ffff-ffffffffffff", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewFirewallRulesClient().Get(ctx, "TestGroup", "myMongoCluster", "rule1", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armmongocluster.FirewallRulesClientGetResponse{
	// 	FirewallRule: &armmongocluster.FirewallRule{
	// 		ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/TestResourceGroup/providers/Microsoft.DocumentDB/mongoClusters/myMongoCluster/firewallRules/rule1"),
	// 		Name: to.Ptr("rule1"),
	// 		Type: to.Ptr("/Microsoft.DocumentDB/mongoClusters/firewallRules"),
	// 		SystemData: &armmongocluster.SystemData{
	// 			CreatedBy: to.Ptr("user1"),
	// 			CreatedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
	// 			CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-01T17:18:19.1234567Z"); return t}()),
	// 			LastModifiedBy: to.Ptr("user2"),
	// 			LastModifiedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
	// 			LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-02T17:18:19.1234567Z"); return t}()),
	// 		},
	// 		Properties: &armmongocluster.FirewallRuleProperties{
	// 			ProvisioningState: to.Ptr(armmongocluster.ProvisioningStateSucceeded),
	// 			StartIPAddress: to.Ptr("0.0.0.0"),
	// 			EndIPAddress: to.Ptr("255.255.255.255"),
	// 		},
	// 	},
	// }
}

// Generated from example definition: /mnt/vss/_work/1/s/eng/azure-sdk-for-go/sdk/resourcemanager/mongocluster/armmongocluster/TempTypeSpecFiles/DocumentDB.MongoCluster.Management/examples/2024-07-01/MongoClusters_FirewallRuleList.json
func ExampleFirewallRulesClient_NewListByMongoClusterPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armmongocluster.NewClientFactory("ffffffff-ffff-ffff-ffff-ffffffffffff", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewFirewallRulesClient().NewListByMongoClusterPager("TestGroup", "myMongoCluster", nil)
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
		// page = armmongocluster.FirewallRulesClientListByMongoClusterResponse{
		// 	FirewallRuleListResult: armmongocluster.FirewallRuleListResult{
		// 		Value: []*armmongocluster.FirewallRule{
		// 			{
		// 				ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/TestResourceGroup/providers/Microsoft.DocumentDB/mongoClusters/myMongoCluster/firewallRules/rule1"),
		// 				Name: to.Ptr("rule1"),
		// 				Type: to.Ptr("/Microsoft.DocumentDB/mongoClusters/firewallRules"),
		// 				SystemData: &armmongocluster.SystemData{
		// 					CreatedBy: to.Ptr("user1"),
		// 					CreatedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-01T17:18:19.1234567Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("user2"),
		// 					LastModifiedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-02T17:18:19.1234567Z"); return t}()),
		// 				},
		// 				Properties: &armmongocluster.FirewallRuleProperties{
		// 					ProvisioningState: to.Ptr(armmongocluster.ProvisioningStateSucceeded),
		// 					StartIPAddress: to.Ptr("0.0.0.0"),
		// 					EndIPAddress: to.Ptr("255.255.255.255"),
		// 				},
		// 			},
		// 			{
		// 				ID: to.Ptr("/subscriptions/ffffffff-ffff-ffff-ffff-ffffffffffff/resourceGroups/TestResourceGroup/providers/Microsoft.DocumentDB/mongoClusters/myMongoCluster/firewallRules/rule2"),
		// 				Name: to.Ptr("rule2"),
		// 				Type: to.Ptr("/Microsoft.DocumentDB/mongoClusters/firewallRules"),
		// 				SystemData: &armmongocluster.SystemData{
		// 					CreatedBy: to.Ptr("user1"),
		// 					CreatedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
		// 					CreatedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-01T17:18:19.1234567Z"); return t}()),
		// 					LastModifiedBy: to.Ptr("user2"),
		// 					LastModifiedByType: to.Ptr(armmongocluster.CreatedByTypeUser),
		// 					LastModifiedAt: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2020-01-02T17:18:19.1234567Z"); return t}()),
		// 				},
		// 				Properties: &armmongocluster.FirewallRuleProperties{
		// 					ProvisioningState: to.Ptr(armmongocluster.ProvisioningStateSucceeded),
		// 					StartIPAddress: to.Ptr("1.0.0.0"),
		// 					EndIPAddress: to.Ptr("255.0.0.0"),
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}
