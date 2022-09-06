//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcosmos_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountGet.json
func ExampleDatabaseAccountsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"rg1",
		"ddb1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountPatch.json
func ExampleDatabaseAccountsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginUpdate(ctx,
		"rg1",
		"ddb1",
		armcosmos.DatabaseAccountUpdateParameters{
			Identity: &armcosmos.ManagedServiceIdentity{
				Type: to.Ptr(armcosmos.ResourceIdentityTypeSystemAssignedUserAssigned),
				UserAssignedIdentities: map[string]*armcosmos.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
					"/subscriptions/fa5fc227-a624-475e-b696-cdd604c735bc/resourceGroups/eu2cgroup/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": {},
				},
			},
			Location: to.Ptr("westus"),
			Properties: &armcosmos.DatabaseAccountUpdateProperties{
				AnalyticalStorageConfiguration: &armcosmos.AnalyticalStorageConfiguration{
					SchemaType: to.Ptr(armcosmos.AnalyticalStorageSchemaTypeWellDefined),
				},
				BackupPolicy: &armcosmos.PeriodicModeBackupPolicy{
					Type: to.Ptr(armcosmos.BackupPolicyTypePeriodic),
					PeriodicModeProperties: &armcosmos.PeriodicModeProperties{
						BackupIntervalInMinutes:        to.Ptr[int32](240),
						BackupRetentionIntervalInHours: to.Ptr[int32](720),
						BackupStorageRedundancy:        to.Ptr(armcosmos.BackupStorageRedundancyGeo),
					},
				},
				Capacity: &armcosmos.Capacity{
					TotalThroughputLimit: to.Ptr[int32](2000),
				},
				ConsistencyPolicy: &armcosmos.ConsistencyPolicy{
					DefaultConsistencyLevel: to.Ptr(armcosmos.DefaultConsistencyLevelBoundedStaleness),
					MaxIntervalInSeconds:    to.Ptr[int32](10),
					MaxStalenessPrefix:      to.Ptr[int64](200),
				},
				DefaultIdentity: to.Ptr("FirstPartyIdentity"),
				DiagnosticLogSettings: &armcosmos.DiagnosticLogSettings{
					EnableFullTextQuery: to.Ptr(armcosmos.EnableFullTextQueryTrue),
				},
				EnableAnalyticalStorage: to.Ptr(true),
				EnableFreeTier:          to.Ptr(false),
				IPRules: []*armcosmos.IPAddressOrRange{
					{
						IPAddressOrRange: to.Ptr("23.43.230.120"),
					},
					{
						IPAddressOrRange: to.Ptr("110.12.240.0/12"),
					}},
				IsVirtualNetworkFilterEnabled: to.Ptr(true),
				NetworkACLBypass:              to.Ptr(armcosmos.NetworkACLBypassAzureServices),
				NetworkACLBypassResourceIDs: []*string{
					to.Ptr("/subscriptions/subId/resourcegroups/rgName/providers/Microsoft.Synapse/workspaces/workspaceName")},
				VirtualNetworkRules: []*armcosmos.VirtualNetworkRule{
					{
						ID:                               to.Ptr("/subscriptions/subId/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
						IgnoreMissingVNetServiceEndpoint: to.Ptr(false),
					}},
			},
			Tags: map[string]*string{
				"dept": to.Ptr("finance"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountCreateMax.json
func ExampleDatabaseAccountsClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreateOrUpdate(ctx,
		"rg1",
		"ddb1",
		armcosmos.DatabaseAccountCreateUpdateParameters{
			Identity: &armcosmos.ManagedServiceIdentity{
				Type: to.Ptr(armcosmos.ResourceIdentityTypeSystemAssignedUserAssigned),
				UserAssignedIdentities: map[string]*armcosmos.Components1Jq1T4ISchemasManagedserviceidentityPropertiesUserassignedidentitiesAdditionalproperties{
					"/subscriptions/fa5fc227-a624-475e-b696-cdd604c735bc/resourceGroups/eu2cgroup/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": {},
				},
			},
			Location: to.Ptr("westus"),
			Tags:     map[string]*string{},
			Kind:     to.Ptr(armcosmos.DatabaseAccountKindMongoDB),
			Properties: &armcosmos.DatabaseAccountCreateUpdateProperties{
				AnalyticalStorageConfiguration: &armcosmos.AnalyticalStorageConfiguration{
					SchemaType: to.Ptr(armcosmos.AnalyticalStorageSchemaTypeWellDefined),
				},
				APIProperties: &armcosmos.APIProperties{
					ServerVersion: to.Ptr(armcosmos.ServerVersionThree2),
				},
				BackupPolicy: &armcosmos.PeriodicModeBackupPolicy{
					Type: to.Ptr(armcosmos.BackupPolicyTypePeriodic),
					PeriodicModeProperties: &armcosmos.PeriodicModeProperties{
						BackupIntervalInMinutes:        to.Ptr[int32](240),
						BackupRetentionIntervalInHours: to.Ptr[int32](8),
						BackupStorageRedundancy:        to.Ptr(armcosmos.BackupStorageRedundancyGeo),
					},
				},
				Capacity: &armcosmos.Capacity{
					TotalThroughputLimit: to.Ptr[int32](2000),
				},
				ConsistencyPolicy: &armcosmos.ConsistencyPolicy{
					DefaultConsistencyLevel: to.Ptr(armcosmos.DefaultConsistencyLevelBoundedStaleness),
					MaxIntervalInSeconds:    to.Ptr[int32](10),
					MaxStalenessPrefix:      to.Ptr[int64](200),
				},
				Cors: []*armcosmos.CorsPolicy{
					{
						AllowedOrigins: to.Ptr("https://test"),
					}},
				CreateMode:               to.Ptr(armcosmos.CreateModeDefault),
				DatabaseAccountOfferType: to.Ptr("Standard"),
				DefaultIdentity:          to.Ptr("FirstPartyIdentity"),
				EnableAnalyticalStorage:  to.Ptr(true),
				EnableFreeTier:           to.Ptr(false),
				EnableMaterializedViews:  to.Ptr(false),
				IPRules: []*armcosmos.IPAddressOrRange{
					{
						IPAddressOrRange: to.Ptr("23.43.230.120"),
					},
					{
						IPAddressOrRange: to.Ptr("110.12.240.0/12"),
					}},
				IsVirtualNetworkFilterEnabled: to.Ptr(true),
				KeyVaultKeyURI:                to.Ptr("https://myKeyVault.vault.azure.net"),
				Locations: []*armcosmos.Location{
					{
						FailoverPriority: to.Ptr[int32](0),
						IsZoneRedundant:  to.Ptr(false),
						LocationName:     to.Ptr("southcentralus"),
					},
					{
						FailoverPriority: to.Ptr[int32](1),
						IsZoneRedundant:  to.Ptr(false),
						LocationName:     to.Ptr("eastus"),
					}},
				NetworkACLBypass: to.Ptr(armcosmos.NetworkACLBypassAzureServices),
				NetworkACLBypassResourceIDs: []*string{
					to.Ptr("/subscriptions/subId/resourcegroups/rgName/providers/Microsoft.Synapse/workspaces/workspaceName")},
				PublicNetworkAccess: to.Ptr(armcosmos.PublicNetworkAccessEnabled),
				VirtualNetworkRules: []*armcosmos.VirtualNetworkRule{
					{
						ID:                               to.Ptr("/subscriptions/subId/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
						IgnoreMissingVNetServiceEndpoint: to.Ptr(false),
					}},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountDelete.json
func ExampleDatabaseAccountsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"rg1",
		"ddb1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountFailoverPriorityChange.json
func ExampleDatabaseAccountsClient_BeginFailoverPriorityChange() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginFailoverPriorityChange(ctx,
		"rg1",
		"ddb1-failover",
		armcosmos.FailoverPolicies{
			FailoverPolicies: []*armcosmos.FailoverPolicy{
				{
					FailoverPriority: to.Ptr[int32](0),
					LocationName:     to.Ptr("eastus"),
				},
				{
					FailoverPriority: to.Ptr[int32](1),
					LocationName:     to.Ptr("westus"),
				}},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountList.json
func ExampleDatabaseAccountsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountListByResourceGroup.json
func ExampleDatabaseAccountsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByResourceGroupPager("rg1",
		nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountListKeys.json
func ExampleDatabaseAccountsClient_ListKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListKeys(ctx,
		"rg1",
		"ddb1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountListConnectionStrings.json
func ExampleDatabaseAccountsClient_ListConnectionStrings() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListConnectionStrings(ctx,
		"rg1",
		"ddb1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountOfflineRegion.json
func ExampleDatabaseAccountsClient_BeginOfflineRegion() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginOfflineRegion(ctx,
		"rg1",
		"ddb1",
		armcosmos.RegionForOnlineOffline{},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountOnlineRegion.json
func ExampleDatabaseAccountsClient_BeginOnlineRegion() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginOnlineRegion(ctx,
		"rg1",
		"ddb1",
		armcosmos.RegionForOnlineOffline{},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountListReadOnlyKeys.json
func ExampleDatabaseAccountsClient_GetReadOnlyKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetReadOnlyKeys(ctx,
		"rg1",
		"ddb1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountListReadOnlyKeys.json
func ExampleDatabaseAccountsClient_ListReadOnlyKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListReadOnlyKeys(ctx,
		"rg1",
		"ddb1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountRegenerateKey.json
func ExampleDatabaseAccountsClient_BeginRegenerateKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginRegenerateKey(ctx,
		"rg1",
		"ddb1",
		armcosmos.DatabaseAccountRegenerateKeyParameters{
			KeyKind: to.Ptr(armcosmos.KeyKindPrimary),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountCheckNameExists.json
func ExampleDatabaseAccountsClient_CheckNameExists() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.CheckNameExists(ctx,
		"ddb1",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountGetMetrics.json
func ExampleDatabaseAccountsClient_NewListMetricsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListMetricsPager("rg1",
		"ddb1",
		"$filter=(name.value eq 'Total Requests') and timeGrain eq duration'PT5M' and startTime eq '2017-11-19T23:53:55.2780000Z' and endTime eq '2017-11-20T00:13:55.2780000Z",
		nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountGetUsages.json
func ExampleDatabaseAccountsClient_NewListUsagesPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListUsagesPager("rg1",
		"ddb1",
		&armcosmos.DatabaseAccountsClientListUsagesOptions{Filter: to.Ptr("$filter=name.value eq 'Storage'")})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/cosmos-db/resource-manager/Microsoft.DocumentDB/preview/2022-05-15-preview/examples/CosmosDBDatabaseAccountGetMetricDefinitions.json
func ExampleDatabaseAccountsClient_NewListMetricDefinitionsPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armcosmos.NewDatabaseAccountsClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListMetricDefinitionsPager("rg1",
		"ddb1",
		nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}
