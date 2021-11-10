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

func TestPrivateEndpointConnectionsClient_Put(t *testing.T) {
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

	pecClient := armstorage.NewPrivateEndpointConnectionsClient(subscriptionID,cred,opt)
	privateEndpointConnectionName,err := createRandomName(t,"connection")
	putResp,err := pecClient.Put(
		ctx,
		rgName,
		scName,
		privateEndpointConnectionName,
		armstorage.PrivateEndpointConnection{
			Properties: &armstorage.PrivateEndpointConnectionProperties{
				PrivateLinkServiceConnectionState: &armstorage.PrivateLinkServiceConnectionState{
					Status: armstorage.PrivateEndpointServiceConnectionStatusRejected.ToPtr(),
					Description: to.StringPtr("Auto-Approved"),
				},
			},
		},
		nil,
		)
	require.NoError(t, err)
	require.Equal(t, privateEndpointConnectionName,*putResp.Name)
}

func TestPrivateEndpointConnectionsClient_Get(t *testing.T) {
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

	// put
	pecClient := armstorage.NewPrivateEndpointConnectionsClient(subscriptionID,cred,opt)
	privateEndpointConnectionName,err := createRandomName(t,"connection")
	putResp,err := pecClient.Put(
		ctx,
		rgName,
		scName,
		privateEndpointConnectionName,
		armstorage.PrivateEndpointConnection{
			Properties: &armstorage.PrivateEndpointConnectionProperties{
				PrivateLinkServiceConnectionState: &armstorage.PrivateLinkServiceConnectionState{
					Status: armstorage.PrivateEndpointServiceConnectionStatusRejected.ToPtr(),
					Description: to.StringPtr("Auto-Approved"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, privateEndpointConnectionName,*putResp.Name)

	// get
	getResp,err := pecClient.Get(ctx,rgName,scName,privateEndpointConnectionName,nil)
	require.NoError(t, err)
	require.Equal(t, privateEndpointConnectionName,*getResp.Name)
}

func TestPrivateEndpointConnectionsClient_List(t *testing.T) {
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

	// put
	pecClient := armstorage.NewPrivateEndpointConnectionsClient(subscriptionID,cred,opt)
	privateEndpointConnectionName,err := createRandomName(t,"connection")
	putResp,err := pecClient.Put(
		ctx,
		rgName,
		scName,
		privateEndpointConnectionName,
		armstorage.PrivateEndpointConnection{
			Properties: &armstorage.PrivateEndpointConnectionProperties{
				PrivateLinkServiceConnectionState: &armstorage.PrivateLinkServiceConnectionState{
					Status: armstorage.PrivateEndpointServiceConnectionStatusRejected.ToPtr(),
					Description: to.StringPtr("Auto-Approved"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, privateEndpointConnectionName,*putResp.Name)

	// list
	listResp,err := pecClient.List(ctx,rgName,scName,nil)
	require.NoError(t, err)
	require.Greater(t,len(listResp.Value), 1)
}

func TestPrivateEndpointConnectionsClient_Delete(t *testing.T) {
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

	// put
	pecClient := armstorage.NewPrivateEndpointConnectionsClient(subscriptionID,cred,opt)
	privateEndpointConnectionName,err := createRandomName(t,"connection")
	putResp,err := pecClient.Put(
		ctx,
		rgName,
		scName,
		privateEndpointConnectionName,
		armstorage.PrivateEndpointConnection{
			Properties: &armstorage.PrivateEndpointConnectionProperties{
				PrivateLinkServiceConnectionState: &armstorage.PrivateLinkServiceConnectionState{
					Status: armstorage.PrivateEndpointServiceConnectionStatusRejected.ToPtr(),
					Description: to.StringPtr("Auto-Approved"),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, privateEndpointConnectionName,*putResp.Name)

	// list
	delResp,err := pecClient.Delete(ctx,rgName,scName,privateEndpointConnectionName,nil)
	require.NoError(t, err)
	require.Equal(t, 200,delResp.RawResponse.StatusCode)
}