//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeysClient_CreateIfNotExist(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	tenantID := recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	objectID := recording.GetEnvVariable("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create vault
	vaultsClient := armkeyvault.NewVaultsClient(subscriptionID, cred, opt)
	vault := createVault(t, ctx, vaultsClient, rgName, location, tenantID, objectID)
	vaultName := *vault.Name

	// create key
	keysClient := armkeyvault.NewKeysClient(subscriptionID, cred, opt)
	keyName, _ := createRandomName(t, "key")
	createResp, err := keysClient.CreateIfNotExist(
		ctx,
		rgName,
		vaultName,
		keyName,
		armkeyvault.KeyCreateParameters{
			Properties: &armkeyvault.KeyProperties{
				Attributes: &armkeyvault.KeyAttributes{
					Enabled: to.BoolPtr(true),
				},
				KeySize: to.Int32Ptr(2048),
				KeyOps: []*armkeyvault.JSONWebKeyOperation{
					armkeyvault.JSONWebKeyOperationEncrypt.ToPtr(),
					armkeyvault.JSONWebKeyOperationDecrypt.ToPtr(),
				},
				Kty: armkeyvault.JSONWebKeyTypeRSA.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, keyName, *createResp.Name)
}

func TestKeysClient_Get(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	tenantID := recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	objectID := recording.GetEnvVariable("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create vault
	vaultsClient := armkeyvault.NewVaultsClient(subscriptionID, cred, opt)
	vault := createVault(t, ctx, vaultsClient, rgName, location, tenantID, objectID)
	vaultName := *vault.Name

	// create key
	keysClient := armkeyvault.NewKeysClient(subscriptionID, cred, opt)
	keyName, _ := createRandomName(t, "key")
	createResp, err := keysClient.CreateIfNotExist(
		ctx,
		rgName,
		vaultName,
		keyName,
		armkeyvault.KeyCreateParameters{
			Properties: &armkeyvault.KeyProperties{
				Attributes: &armkeyvault.KeyAttributes{
					Enabled: to.BoolPtr(true),
				},
				KeySize: to.Int32Ptr(2048),
				KeyOps: []*armkeyvault.JSONWebKeyOperation{
					armkeyvault.JSONWebKeyOperationEncrypt.ToPtr(),
					armkeyvault.JSONWebKeyOperationDecrypt.ToPtr(),
				},
				Kty: armkeyvault.JSONWebKeyTypeRSA.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, keyName, *createResp.Name)

	// get key
	getResp, err := keysClient.Get(ctx, rgName, vaultName, keyName, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, *getResp.Name)
}

func TestKeysClient_GetVersion(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	tenantID := recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	objectID := recording.GetEnvVariable("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg", location)
	defer clean()
	rgName := *rg.Name

	// create vault
	vaultsClient := armkeyvault.NewVaultsClient(subscriptionID, cred, opt)
	vault := createVault(t, ctx, vaultsClient, rgName, location, tenantID, objectID)
	vaultName := *vault.Name

	// create key
	keysClient := armkeyvault.NewKeysClient(subscriptionID, cred, opt)
	keyName, _ := createRandomName(t, "key")
	createResp, err := keysClient.CreateIfNotExist(
		ctx,
		rgName,
		vaultName,
		keyName,
		armkeyvault.KeyCreateParameters{
			Properties: &armkeyvault.KeyProperties{
				Attributes: &armkeyvault.KeyAttributes{
					Enabled: to.BoolPtr(true),
				},
				KeySize: to.Int32Ptr(2048),
				KeyOps: []*armkeyvault.JSONWebKeyOperation{
					armkeyvault.JSONWebKeyOperationEncrypt.ToPtr(),
					armkeyvault.JSONWebKeyOperationDecrypt.ToPtr(),
				},
				Kty: armkeyvault.JSONWebKeyTypeRSA.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, keyName, *createResp.Name)

	// get key
	getResp, err := keysClient.GetVersion(ctx, rgName, vaultName, keyName, "7.2", nil)
	require.NoError(t, err)
	require.Equal(t, keyName, *getResp.Name)
}

func TestKeysClient_List(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	tenantID := recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	objectID := recording.GetEnvVariable("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, _ := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	//defer clean()
	rgName := *rg.Name

	// create vault
	vaultsClient := armkeyvault.NewVaultsClient(subscriptionID, cred, opt)
	vault := createVault(t, ctx, vaultsClient, rgName, location, tenantID, objectID)
	vaultName := *vault.Name

	// create key
	keysClient := armkeyvault.NewKeysClient(subscriptionID, cred, opt)
	keyName, _ := createRandomName(t, "key")
	createResp, err := keysClient.CreateIfNotExist(
		ctx,
		rgName,
		vaultName,
		keyName,
		armkeyvault.KeyCreateParameters{
			Properties: &armkeyvault.KeyProperties{
				Attributes: &armkeyvault.KeyAttributes{
					Enabled: to.BoolPtr(true),
				},
				KeySize: to.Int32Ptr(2048),
				KeyOps: []*armkeyvault.JSONWebKeyOperation{
					armkeyvault.JSONWebKeyOperationEncrypt.ToPtr(),
					armkeyvault.JSONWebKeyOperationDecrypt.ToPtr(),
				},
				Kty: armkeyvault.JSONWebKeyTypeRSA.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, keyName, *createResp.Name)

	// list
	versions := keysClient.List(rgName, vaultName, nil)
	require.NoError(t, versions.Err())
}

func TestKeysClient_ListVersions(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	tenantID := recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	objectID := recording.GetEnvVariable("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()
	location := "westus"

	// create resource group
	rg, clean := createResourceGroup(t, cred, opt, subscriptionID, "rg2", location)
	defer clean()
	rgName := *rg.Name

	// create vault
	vaultsClient := armkeyvault.NewVaultsClient(subscriptionID, cred, opt)
	vault := createVault(t, ctx, vaultsClient, rgName, location, tenantID, objectID)
	vaultName := *vault.Name

	// create key
	keysClient := armkeyvault.NewKeysClient(subscriptionID, cred, opt)
	keyName, _ := createRandomName(t, "key")
	createResp, err := keysClient.CreateIfNotExist(
		ctx,
		rgName,
		vaultName,
		keyName,
		armkeyvault.KeyCreateParameters{
			Properties: &armkeyvault.KeyProperties{
				Attributes: &armkeyvault.KeyAttributes{
					Enabled: to.BoolPtr(true),
				},
				KeySize: to.Int32Ptr(2048),
				KeyOps: []*armkeyvault.JSONWebKeyOperation{
					armkeyvault.JSONWebKeyOperationEncrypt.ToPtr(),
					armkeyvault.JSONWebKeyOperationDecrypt.ToPtr(),
				},
				Kty: armkeyvault.JSONWebKeyTypeRSA.ToPtr(),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, keyName, *createResp.Name)

	// list versions
	versions := keysClient.ListVersions(rgName, vaultName, keyName, nil)
	require.NoError(t, versions.Err())
}
