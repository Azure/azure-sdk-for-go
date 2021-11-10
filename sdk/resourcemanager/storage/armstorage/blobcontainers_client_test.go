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

func TestBlobContainersClient_Create(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)
}

func TestBlobContainersClient_CreateOrUpdateImmutabilityPolicy(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// create immutability policy container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.CreateOrUpdateImmutabilityPolicy(
		ctx,
		rgName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersCreateOrUpdateImmutabilityPolicyOptions{
			Parameters: &armstorage.ImmutabilityPolicy{
				Properties: &armstorage.ImmutabilityPolicyProperty{
					ImmutabilityPeriodSinceCreationInDays: to.Int32Ptr(3),
					AllowProtectedAppendWrites: to.BoolPtr(true),
				},
			},
		},
		)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)
}

func TestBlobContainersClient_DeleteImmutabilityPolicy(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// create immutability policy container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.CreateOrUpdateImmutabilityPolicy(
		ctx,
		rgName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersCreateOrUpdateImmutabilityPolicyOptions{
			Parameters: &armstorage.ImmutabilityPolicy{
				Properties: &armstorage.ImmutabilityPolicyProperty{
					ImmutabilityPeriodSinceCreationInDays: to.Int32Ptr(3),
					AllowProtectedAppendWrites: to.BoolPtr(true),
				},
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// delete immutability policy container
	eTag := *blobResp.ETag
	delResp,err := blobContainersClient.DeleteImmutabilityPolicy(ctx, rgName, scName, blobContainerName, eTag, nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,delResp.Name)
}

func TestBlobContainersClient_GetImmutabilityPolicy(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// create immutability policy container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.CreateOrUpdateImmutabilityPolicy(
		ctx,
		rgName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersCreateOrUpdateImmutabilityPolicyOptions{
			Parameters: &armstorage.ImmutabilityPolicy{
				Properties: &armstorage.ImmutabilityPolicyProperty{
					ImmutabilityPeriodSinceCreationInDays: to.Int32Ptr(3),
					AllowProtectedAppendWrites: to.BoolPtr(true),
				},
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// get immutability policy
	getResp,err := blobContainersClient.GetImmutabilityPolicy(ctx,rgName,scName,blobContainerName,nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,*getResp.Name)
}

func TestBlobContainersClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// get
	getResp,err := blobContainersClient.Get(ctx,rgName,scName,blobContainerName,nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,*getResp.Name)
}

func TestBlobContainersClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// list
	listPager := blobContainersClient.List(rgName,scName,nil)
	require.True(t,listPager.NextPage(ctx))
}

func TestBlobContainersClient_SetLegalHold(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// set legal hold
	holdResp,err := blobContainersClient.SetLegalHold(
		ctx,
		rgName,
		scName,
		blobContainerName,
		armstorage.LegalHold{
			Tags: []*string{
				to.StringPtr("tag1"),
				to.StringPtr("tag2"),
				to.StringPtr("tag3"),
			},
		},
		nil,
		)
	require.NoError(t, err)
	require.Equal(t,3,len(holdResp.Tags))
}

func TestBlobContainersClient_ClearLegalHold(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// clear legal hold
	holdResp,err := blobContainersClient.ClearLegalHold(
		ctx,
		rgName,
		scName,
		blobContainerName,
		armstorage.LegalHold{
			Tags: []*string{
				to.StringPtr("tag1"),
				to.StringPtr("tag2"),
				to.StringPtr("tag3"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t,0,len(holdResp.Tags))
}

func TestBlobContainersClient_Lease(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// acquire lease
	leaseResp,err := blobContainersClient.Lease(
		ctx,
		rgName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersLeaseOptions{
			Parameters: &armstorage.LeaseContainerRequest{
				Action: armstorage.LeaseContainerRequestActionAcquire.ToPtr(),
				LeaseDuration: to.Int32Ptr(-1),
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t,1,*leaseResp.LeaseTimeSeconds)

	// break lease
	breakResp,err := blobContainersClient.Lease(
		ctx,
		rgName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersLeaseOptions{
			Parameters: &armstorage.LeaseContainerRequest{
				Action: armstorage.LeaseContainerRequestActionBreak.ToPtr(),
				LeaseID: leaseResp.LeaseID,
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t,1,*breakResp.LeaseTimeSeconds)
}

func TestBlobContainersClient_Update(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// update
	updateResp,err := blobContainersClient.Update(
		ctx,
		rgName,
		scName,
		blobContainerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{
				PublicAccess: armstorage.PublicAccessContainer.ToPtr(),
				Metadata: map[string]*string{
					"metadata": to.StringPtr("true"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t,blobContainerName,*updateResp.Name)
}

func TestBlobContainersClient_LockImmutabilityPolicy(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// create immutability policy container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.CreateOrUpdateImmutabilityPolicy(
		ctx,
		rgName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersCreateOrUpdateImmutabilityPolicyOptions{
			Parameters: &armstorage.ImmutabilityPolicy{
				Properties: &armstorage.ImmutabilityPolicyProperty{
					ImmutabilityPeriodSinceCreationInDays: to.Int32Ptr(3),
					AllowProtectedAppendWrites: to.BoolPtr(true),
				},
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// lock immutability policy
	lockResp,err := blobContainersClient.LockImmutabilityPolicy(ctx,rgName,scName,blobContainerName,*blobResp.ETag,nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,*lockResp.Name)
}

func TestBlobContainersClient_ExtendImmutabilityPolicy(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// create immutability policy container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.CreateOrUpdateImmutabilityPolicy(
		ctx,
		rgName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersCreateOrUpdateImmutabilityPolicyOptions{
			Parameters: &armstorage.ImmutabilityPolicy{
				Properties: &armstorage.ImmutabilityPolicyProperty{
					ImmutabilityPeriodSinceCreationInDays: to.Int32Ptr(3),
					AllowProtectedAppendWrites: to.BoolPtr(true),
				},
			},
		},
	)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// lock immutability policy
	//lockResp,err := blobContainersClient.LockImmutabilityPolicy(ctx,rgName,scName,blobContainerName,*blobResp.ETag,nil)
	//require.NoError(t, err)
	//require.Equal(t, blobContainerName,*lockResp.Name)

	extendResp,err := blobContainersClient.ExtendImmutabilityPolicy(ctx,rgName,scName,blobContainerName,*blobResp.ETag,nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,*extendResp.Name)
}

func TestBlobContainersClient_Delete(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create resource group
	rg,clean := createResourceGroup(t,cred,opt,subscriptionID,"create","westus")
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

	// create container
	blobContainersClient := armstorage.NewBlobContainersClient(subscriptionID,cred,opt)
	blobContainerName,err := createRandomName(t,"container")
	require.NoError(t, err)
	blobResp,err := blobContainersClient.Create(ctx,rgName,scName,blobContainerName,armstorage.BlobContainer{},nil)
	require.NoError(t, err)
	require.Equal(t, blobContainerName,blobResp.Name)

	// delete
	delResp,err := blobContainersClient.Delete(ctx,rgName,scName,blobContainerName,nil)
	require.NoError(t, err)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}