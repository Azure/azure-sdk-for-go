package armcompute_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDisksClient_CreateOrUpdate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	cred, opt := authenticateTest(t)
	conn := arm.NewDefaultConnection(cred, opt)
	subscriptionID := recording.GetEnvVariable(t, "AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	tenantID := recording.GetEnvVariable(t, "AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")

	// create resource group
	rgName, err := createRandomName(t, "testRP")
	require.NoError(t, err)
	rgClient := armresources.NewResourceGroupsClient(conn, subscriptionID)
	_, err = rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
		Location: to.StringPtr("westus2"),
	}, nil)
	defer cleanup(t, rgClient, rgName)
	require.NoError(t, err)

	// create vault
	vClient := armkeyvault.NewVaultsClient(conn, subscriptionID)
	vName, err := createRandomName(t, "vault")
	require.NoError(t, err)
	vPoller, err := vClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		vName,
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
						ObjectID: to.StringPtr("00000000-0000-0000-0000-000000000000"),
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
	vResp, err := vPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *vResp.Name, vName)

	// create key
	keyClient := armkeyvault.NewKeysClient(conn, subscriptionID)
	keyName, err := createRandomName(t, "key")
	require.NoError(t, err)
	keyResp, err := keyClient.CreateIfNotExist(
		context.Background(),
		rgName,
		vName,
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
	require.Equal(t, keyResp.Name, keyName)

	// create disk
	diskClient := armcompute.NewDisksClient(conn, subscriptionID)
	diskName, err := createRandomName(t, "disk")
	require.NoError(t, err)
	diskPoller, err := diskClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		diskName,
		armcompute.Disk{
			Resource: armcompute.Resource{
				Location: to.StringPtr(location),
			},
			SKU: &armcompute.DiskSKU{
				Name: armcompute.DiskStorageAccountTypesStandardLRS.ToPtr(),
			},
			Properties: &armcompute.DiskProperties{
				CreationData: &armcompute.CreationData{
					CreateOption: armcompute.DiskCreateOptionEmpty.ToPtr(),
				},
				DiskSizeGB: to.Int32Ptr(64),
			},
		},
		nil,
	)
	require.NoError(t, err)
	diskResp, err := diskPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *diskResp.Name, diskName)

	// create disk encryption set
	desClient := armcompute.NewDiskEncryptionSetsClient(conn, subscriptionID)
	desName, err := createRandomName(t, "set")
	require.NoError(t, err)
	desPoller, err := desClient.BeginCreateOrUpdate(
		context.Background(),
		rgName,
		desName,
		armcompute.DiskEncryptionSet{
			Resource: armcompute.Resource{
				Location: to.StringPtr(location),
			},
			Identity: &armcompute.EncryptionSetIdentity{
				Type: armcompute.DiskEncryptionSetIdentityTypeSystemAssigned.ToPtr(),
			},
			Properties: &armcompute.EncryptionSetProperties{
				ActiveKey: &armcompute.KeyForDiskEncryptionSet{
					SourceVault: &armcompute.SourceVault{
						ID: to.StringPtr(*vResp.ID),
					},
					KeyURL: to.StringPtr(*keyResp.Properties.KeyURIWithVersion),
				},
			},
		},
		nil,
	)
	require.NoError(t, err)
	desResp, err := desPoller.PollUntilDone(context.Background(), 10*time.Second)
	require.NoError(t, err)
	require.Equal(t, *desResp.Name, desName)
}
