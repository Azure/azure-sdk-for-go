//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEncryptionScopesClient_Put(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create storage account
	storageAccountsClient := armstorage.NewStorageAccountsClient(subscriptionID,cred,opt)
	scName,err := createRandomName(t,"account")
	require.NoError(t, err)
	pollerResp,err := storageAccountsClient.BeginCreate(
		ctx,
		rgName,
		scName,
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind: armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr("westus"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
					},
					KeySource: armstorage.KeySourceMicrosoftStorage.ToPtr(),
				},
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
		)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scName,*resp.Name)

	// put file services
	encryptionScopesClient := armstorage.NewEncryptionScopesClient(subscriptionID,cred,opt)
	encryptionScopeName,err := createRandomName(t,"encryption")
	putResp,err := encryptionScopesClient.Put(
		ctx,
		rgName,
		scName,
		encryptionScopeName,
		armstorage.EncryptionScope{
			EncryptionScopeProperties: &armstorage.EncryptionScopeProperties{
				Source: armstorage.EncryptionScopeSourceMicrosoftStorage.ToPtr(),
				State: armstorage.EncryptionScopeStateEnabled.ToPtr(),
			},
		},
		nil,
		)
	require.NoError(t, err)
	require.Equal(t, scName,putResp.Name)
}

func TestEncryptionScopesClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create storage account
	storageAccountsClient := armstorage.NewStorageAccountsClient(subscriptionID,cred,opt)
	scName,err := createRandomName(t,"account")
	require.NoError(t, err)
	pollerResp,err := storageAccountsClient.BeginCreate(
		ctx,
		rgName,
		scName,
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind: armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr("westus"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
					},
					KeySource: armstorage.KeySourceMicrosoftStorage.ToPtr(),
				},
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scName,*resp.Name)

	// put file services
	encryptionScopesClient := armstorage.NewEncryptionScopesClient(subscriptionID,cred,opt)
	encryptionScopeName,err := createRandomName(t,"encryption")
	putResp,err := encryptionScopesClient.Put(
		ctx,
		rgName,
		scName,
		encryptionScopeName,
		armstorage.EncryptionScope{
			EncryptionScopeProperties: &armstorage.EncryptionScopeProperties{
				Source: armstorage.EncryptionScopeSourceMicrosoftStorage.ToPtr(),
				State: armstorage.EncryptionScopeStateEnabled.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, scName,putResp.Name)

	// get
	getResp,err := encryptionScopesClient.Get(ctx,rgName,scName,encryptionScopeName,nil)
	require.NoError(t, err)
	require.Equal(t,encryptionScopeName,*getResp.Name)
}

func TestEncryptionScopesClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create storage account
	storageAccountsClient := armstorage.NewStorageAccountsClient(subscriptionID,cred,opt)
	scName,err := createRandomName(t,"account")
	require.NoError(t, err)
	pollerResp,err := storageAccountsClient.BeginCreate(
		ctx,
		rgName,
		scName,
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind: armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr("westus"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
					},
					KeySource: armstorage.KeySourceMicrosoftStorage.ToPtr(),
				},
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scName,*resp.Name)

	// put file services
	encryptionScopesClient := armstorage.NewEncryptionScopesClient(subscriptionID,cred,opt)
	encryptionScopeName,err := createRandomName(t,"encryption")
	putResp,err := encryptionScopesClient.Put(
		ctx,
		rgName,
		scName,
		encryptionScopeName,
		armstorage.EncryptionScope{
			EncryptionScopeProperties: &armstorage.EncryptionScopeProperties{
				Source: armstorage.EncryptionScopeSourceMicrosoftStorage.ToPtr(),
				State: armstorage.EncryptionScopeStateEnabled.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, scName,putResp.Name)

	// list
	listPager := encryptionScopesClient.List(rgName,scName,nil)
	require.NoError(t, err)
	require.True(t,listPager.NextPage(ctx))
}

func TestEncryptionScopesClient_Patch(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"move","westus")
	defer clean()
	rgName := *rg.Name

	// create storage account
	storageAccountsClient := armstorage.NewStorageAccountsClient(subscriptionID,cred,opt)
	scName,err := createRandomName(t,"account")
	require.NoError(t, err)
	pollerResp,err := storageAccountsClient.BeginCreate(
		ctx,
		rgName,
		scName,
		armstorage.StorageAccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind: armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr("westus"),
			Properties: &armstorage.StorageAccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: armstorage.KeyTypeAccount.ToPtr(),
							Enabled: to.BoolPtr(true),
						},
					},
					KeySource: armstorage.KeySourceMicrosoftStorage.ToPtr(),
				},
			},
			Tags: map[string]*string{
				"key1": to.StringPtr("value1"),
				"key2": to.StringPtr("value2"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	resp,err := pollerResp.PollUntilDone(ctx,10*time.Second)
	require.NoError(t, err)
	require.Equal(t, scName,*resp.Name)

	// put file services
	encryptionScopesClient := armstorage.NewEncryptionScopesClient(subscriptionID,cred,opt)
	encryptionScopeName,err := createRandomName(t,"encryption")
	putResp,err := encryptionScopesClient.Put(
		ctx,
		rgName,
		scName,
		encryptionScopeName,
		armstorage.EncryptionScope{
			EncryptionScopeProperties: &armstorage.EncryptionScopeProperties{
				Source: armstorage.EncryptionScopeSourceMicrosoftStorage.ToPtr(),
				State: armstorage.EncryptionScopeStateEnabled.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, scName,putResp.Name)

	// patch
	patchResp,err := encryptionScopesClient.Patch(
		ctx,
		rgName,
		scName,
		encryptionScopeName,
		armstorage.EncryptionScope{
			EncryptionScopeProperties: &armstorage.EncryptionScopeProperties{
				Source: armstorage.EncryptionScopeSourceMicrosoftStorage.ToPtr(),
				State: armstorage.EncryptionScopeStateEnabled.ToPtr(),
			},
		},
		nil,
		)
	require.NoError(t, err)
	require.Equal(t,encryptionScopeName,*patchResp.Name)
}