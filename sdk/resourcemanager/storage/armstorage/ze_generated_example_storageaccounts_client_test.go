//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstorage_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountCheckNameAvailability.json
func ExampleAccountsClient_CheckNameAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CheckNameAvailability(ctx,
		armstorage.AccountCheckNameAvailabilityParameters{
			Name: to.Ptr("sto3363"),
			Type: to.Ptr("Microsoft.Storage/storageAccounts"),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/NfsV3AccountCreate.json
func ExampleAccountsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreate(ctx,
		"res9101",
		"sto4445",
		armstorage.AccountCreateParameters{
			Kind:     to.Ptr(armstorage.KindBlockBlobStorage),
			Location: to.Ptr("eastus"),
			Properties: &armstorage.AccountPropertiesCreateParameters{
				IsHnsEnabled: to.Ptr(true),
				EnableNfsV3:  to.Ptr(true),
				NetworkRuleSet: &armstorage.NetworkRuleSet{
					Bypass:        to.Ptr(armstorage.BypassAzureServices),
					DefaultAction: to.Ptr(armstorage.DefaultActionAllow),
					IPRules:       []*armstorage.IPRule{},
					VirtualNetworkRules: []*armstorage.VirtualNetworkRule{
						{
							VirtualNetworkResourceID: to.Ptr("/subscriptions/{subscription-id}/resourceGroups/res9101/providers/Microsoft.Network/virtualNetworks/net123/subnets/subnet12"),
						}},
				},
				EnableHTTPSTrafficOnly: to.Ptr(false),
			},
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNamePremiumLRS),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountDelete.json
func ExampleAccountsClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"res4228",
		"sto2434",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountGetAsyncSkuConversionStatus.json
func ExampleAccountsClient_GetProperties() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.GetProperties(ctx,
		"res9407",
		"sto8596",
		&armstorage.AccountsClientGetPropertiesOptions{Expand: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountEnableAD.json
func ExampleAccountsClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"res9407",
		"sto8596",
		armstorage.AccountUpdateParameters{
			Properties: &armstorage.AccountPropertiesUpdateParameters{
				AzureFilesIdentityBasedAuthentication: &armstorage.AzureFilesIdentityBasedAuthentication{
					ActiveDirectoryProperties: &armstorage.ActiveDirectoryProperties{
						AccountType:       to.Ptr(armstorage.ActiveDirectoryPropertiesAccountTypeUser),
						AzureStorageSid:   to.Ptr("S-1-5-21-2400535526-2334094090-2402026252-0012"),
						DomainGUID:        to.Ptr("aebfc118-9fa9-4732-a21f-d98e41a77ae1"),
						DomainName:        to.Ptr("adtest.com"),
						DomainSid:         to.Ptr("S-1-5-21-2400535526-2334094090-2402026252"),
						ForestName:        to.Ptr("adtest.com"),
						NetBiosDomainName: to.Ptr("adtest.com"),
						SamAccountName:    to.Ptr("sam12498"),
					},
					DirectoryServiceOptions: to.Ptr(armstorage.DirectoryServiceOptionsAD),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountList.json
func ExampleAccountsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountListByResourceGroup.json
func ExampleAccountsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByResourceGroupPager("res6117",
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountListKeys.json
func ExampleAccountsClient_ListKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListKeys(ctx,
		"res418",
		"sto2220",
		&armstorage.AccountsClientListKeysOptions{Expand: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountRegenerateKerbKey.json
func ExampleAccountsClient_RegenerateKey() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.RegenerateKey(ctx,
		"res4167",
		"sto3539",
		armstorage.AccountRegenerateKeyParameters{
			KeyName: to.Ptr("kerb1"),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountListAccountSAS.json
func ExampleAccountsClient_ListAccountSAS() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListAccountSAS(ctx,
		"res7985",
		"sto8588",
		armstorage.AccountSasParameters{
			KeyToSign:              to.Ptr("key1"),
			SharedAccessExpiryTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-05-24T11:42:03.1567373Z"); return t }()),
			Permissions:            to.Ptr(armstorage.PermissionsR),
			Protocols:              to.Ptr(armstorage.HTTPProtocolHTTPSHTTP),
			ResourceTypes:          to.Ptr(armstorage.SignedResourceTypesS),
			Services:               to.Ptr(armstorage.ServicesB),
			SharedAccessStartTime:  to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-05-24T10:42:03.1567373Z"); return t }()),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountListServiceSAS.json
func ExampleAccountsClient_ListServiceSAS() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.ListServiceSAS(ctx,
		"res7439",
		"sto1299",
		armstorage.ServiceSasParameters{
			CanonicalizedResource:  to.Ptr("/blob/sto1299/music"),
			SharedAccessExpiryTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2017-05-24T11:32:48.8457197Z"); return t }()),
			Permissions:            to.Ptr(armstorage.PermissionsL),
			Resource:               to.Ptr(armstorage.SignedResourceC),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountFailover.json
func ExampleAccountsClient_BeginFailover() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginFailover(ctx,
		"res4228",
		"sto2434",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountHierarchicalNamespaceMigration.json
func ExampleAccountsClient_BeginHierarchicalNamespaceMigration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginHierarchicalNamespaceMigration(ctx,
		"res4228",
		"sto2434",
		"HnsOnValidationRequest",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountAbortHierarchicalNamespaceMigration.json
func ExampleAccountsClient_BeginAbortHierarchicalNamespaceMigration() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginAbortHierarchicalNamespaceMigration(ctx,
		"res4228",
		"sto2434",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/BlobRangesRestore.json
func ExampleAccountsClient_BeginRestoreBlobRanges() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginRestoreBlobRanges(ctx,
		"res9101",
		"sto4445",
		armstorage.BlobRestoreParameters{
			BlobRanges: []*armstorage.BlobRestoreRange{
				{
					EndRange:   to.Ptr("container/blobpath2"),
					StartRange: to.Ptr("container/blobpath1"),
				},
				{
					EndRange:   to.Ptr(""),
					StartRange: to.Ptr("container2/blobpath3"),
				}},
			TimeToRestore: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2019-04-20T15:30:00.0000000Z"); return t }()),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/storage/resource-manager/Microsoft.Storage/stable/2021-09-01/examples/StorageAccountRevokeUserDelegationKeys.json
func ExampleAccountsClient_RevokeUserDelegationKeys() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armstorage.NewAccountsClient("{subscription-id}", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.RevokeUserDelegationKeys(ctx,
		"res4167",
		"sto3539",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
