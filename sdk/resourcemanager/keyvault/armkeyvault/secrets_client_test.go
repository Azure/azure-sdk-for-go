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

func TestSecretsClient_CreateOrUpdate(t *testing.T) {
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

	// create secret
	secretsClient := armkeyvault.NewSecretsClient(subscriptionID, cred, opt)
	secretName, _ := createRandomName(t, "secret")
	secretResp, err := secretsClient.CreateOrUpdate(
		ctx,
		rgName,
		*vault.Name,
		secretName,
		armkeyvault.SecretCreateOrUpdateParameters{
			Properties: &armkeyvault.SecretProperties{
				Attributes: &armkeyvault.SecretAttributes{
					Enabled: to.BoolPtr(true),
				},
				Value: to.StringPtr("sample-secret-value"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, secretName, *secretResp.Name)
}

func TestSecretsClient_Update(t *testing.T) {
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

	// create secret
	secretsClient := armkeyvault.NewSecretsClient(subscriptionID, cred, opt)
	secretName, _ := createRandomName(t, "secret")
	secretResp, err := secretsClient.CreateOrUpdate(
		ctx,
		rgName,
		*vault.Name,
		secretName,
		armkeyvault.SecretCreateOrUpdateParameters{
			Properties: &armkeyvault.SecretProperties{
				Attributes: &armkeyvault.SecretAttributes{
					Enabled: to.BoolPtr(true),
				},
				Value: to.StringPtr("sample-secret-value"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, secretName, *secretResp.Name)

	// update secret
	updateResp, err := secretsClient.Update(
		ctx,
		rgName,
		*vault.Name,
		secretName,
		armkeyvault.SecretPatchParameters{
			Tags: map[string]*string{
				"test": to.StringPtr("recording"),
			},
			Properties: &armkeyvault.SecretPatchProperties{
				Value: to.StringPtr("sample-secret-value-update"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, "recording", *updateResp.Tags["test"])
}

func TestSecretsClient_Get(t *testing.T) {
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

	// create secret
	secretsClient := armkeyvault.NewSecretsClient(subscriptionID, cred, opt)
	secretName, _ := createRandomName(t, "secret")
	secretResp, err := secretsClient.CreateOrUpdate(
		ctx,
		rgName,
		*vault.Name,
		secretName,
		armkeyvault.SecretCreateOrUpdateParameters{
			Properties: &armkeyvault.SecretProperties{
				Attributes: &armkeyvault.SecretAttributes{
					Enabled: to.BoolPtr(true),
				},
				Value: to.StringPtr("sample-secret-value"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, secretName, *secretResp.Name)

	// get secret
	getResp, err := secretsClient.Get(ctx, rgName, *vault.Name, secretName, nil)
	require.NoError(t, err)
	require.Equal(t, secretName, *getResp.Name)
}

func TestSecretsClient_List(t *testing.T) {
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

	// create secret
	secretsClient := armkeyvault.NewSecretsClient(subscriptionID, cred, opt)
	secretName, _ := createRandomName(t, "secret")
	secretResp, err := secretsClient.CreateOrUpdate(
		ctx,
		rgName,
		*vault.Name,
		secretName,
		armkeyvault.SecretCreateOrUpdateParameters{
			Properties: &armkeyvault.SecretProperties{
				Attributes: &armkeyvault.SecretAttributes{
					Enabled: to.BoolPtr(true),
				},
				Value: to.StringPtr("sample-secret-value"),
			},
		},
		nil,
	)
	require.NoError(t, err)
	require.Equal(t, secretName, *secretResp.Name)

	// get secret
	secretPager := secretsClient.List(rgName, *vault.Name, nil)
	require.NoError(t, secretPager.Err())
}

func createVault(t *testing.T, ctx context.Context, vaultsClient *armkeyvault.VaultsClient, resourceGroupName, location, tenantID, objectID string) *armkeyvault.Vault {
	vaultName, _ := createRandomName(t, "vault")
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
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
	vResp, err := vPollerResp.PollUntilDone(ctx, 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, vaultName, *vResp.Name)
	return &vResp.Vault
}
