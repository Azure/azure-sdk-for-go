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
	"time"
)

func TestVaultsClient_BeginCreateOrUpdate(t *testing.T) {
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
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)
}

func TestVaultsClient_Get(t *testing.T) {
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
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)

	// get vault
	getResp, err := vaultsClient.Get(ctx, rgName, vaultName, nil)
	require.NoError(t, err)
	require.Equal(t, vaultName, *getResp.Name)
}

func TestVaultsClient_Update(t *testing.T) {
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
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)

	// update vault
	updateResp, err := vaultsClient.Update(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultPatchParameters{
			Tags: map[string]*string{
				"test": to.StringPtr("recording"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, "recording", *updateResp.Tags["test"])
}

func TestVaultsClient_ListDeleted(t *testing.T) {
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
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)

	// list vault deleted
	deletedPager := vaultsClient.ListDeleted(nil)
	require.NoError(t, deletedPager.Err())
}

func TestVaultsClient_CheckNameAvailability(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	subscriptionID := recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	ctx := context.Background()

	// create vault
	vaultsClient := armkeyvault.NewVaultsClient(subscriptionID, cred, opt)
	vaultName, _ := createRandomName(t, "vault")
	resp, err := vaultsClient.CheckNameAvailability(
		ctx,
		armkeyvault.VaultCheckNameAvailabilityParameters{
			Name: to.StringPtr(vaultName),
			Type: to.StringPtr("Microsoft.KeyVault/vaults"),
		},
		nil,
	)
	require.NoError(t, err)
	require.True(t, *resp.NameAvailable)
}

func TestVaultsClient_Delete(t *testing.T) {
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
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)

	// delete vault
	delResp, err := vaultsClient.Delete(ctx, rgName, vaultName, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)
}

func TestVaultsClient_GetDeleted(t *testing.T) {
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
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)

	// delete vault
	delResp, err := vaultsClient.Delete(ctx, rgName, vaultName, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)

	// get deleted vault
	deletedResp, err := vaultsClient.GetDeleted(ctx, vaultName, location, nil)
	require.NoError(t, err)
	require.Equal(t, vaultName, *deletedResp.Name)
}

func TestVaultsClient_BeginPurgeDeleted(t *testing.T) {
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
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)

	// delete vault
	delResp, err := vaultsClient.Delete(ctx, rgName, vaultName, nil)
	require.NoError(t, err)
	require.Equal(t, 200, delResp.RawResponse.StatusCode)

	// purge deleted vault
	purgePollerResp, err := vaultsClient.BeginPurgeDeleted(ctx, vaultName, location, nil)
	require.NoError(t, err)
	var purgeResp armkeyvault.VaultsClientPurgeDeletedResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = purgePollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if purgePollerResp.Poller.Done() {
				purgeResp, err = purgePollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		purgeResp, err = purgePollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, 200, purgeResp.RawResponse.StatusCode)
}

/*
func TestVaultsClient_UpdateAccessPolicy(t *testing.T) {
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

	vaultsClient := armkeyvault.NewVaultsClient(subscriptionID, cred, opt)
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		rgName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsGet.ToPtr(),
								armkeyvault.KeyPermissionsList.ToPtr(),
								armkeyvault.KeyPermissionsCreate.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
								armkeyvault.SecretPermissionsList.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
								armkeyvault.CertificatePermissionsList.ToPtr(),
								armkeyvault.CertificatePermissionsCreate.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(ctx)
			require.NoError(t, err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(ctx)
				require.NoError(t, err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(ctx, 30*time.Second)
		require.NoError(t, err)
	}
	require.Equal(t, vaultName, *vResp.Name)

	accessPolicyResp, err := vaultsClient.UpdateAccessPolicy(
		ctx,
		rgName,
		vaultName,
		armkeyvault.AccessPolicyUpdateKindAdd,
		armkeyvault.VaultAccessPolicyParameters{
			Properties: &armkeyvault.VaultAccessPolicyProperties{
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(tenantID),
						ObjectID: to.StringPtr(objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								armkeyvault.KeyPermissionsEncrypt.ToPtr(),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								armkeyvault.SecretPermissionsGet.ToPtr(),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								armkeyvault.CertificatePermissionsGet.ToPtr(),
							},
						},
					},
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, 200, accessPolicyResp.RawResponse.StatusCode)
}
*/
