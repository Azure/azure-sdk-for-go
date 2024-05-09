//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/stretchr/testify/suite"
)

type SecretsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
	tenantID          string
	objectID          string
}

func (testsuite *SecretsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.tenantID = recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.objectID = recording.GetEnvVariable("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/keyvault/armkeyvault/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *SecretsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSecretsClient(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	suite.Run(t, new(SecretsClientTestSuite))
}

func (testsuite *SecretsClientTestSuite) TestSecretsCRUD() {
	// create vault
	fmt.Println("Call operation: Vaults_CreateOrUpdate")
	vaultsClient, err := armkeyvault.NewVaultsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vaultName := "go-test-vault-22"
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.Ptr(testsuite.location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: to.Ptr(armkeyvault.SKUFamilyA),
					Name:   to.Ptr(armkeyvault.SKUNameStandard),
				},
				TenantID: to.Ptr(testsuite.tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.Ptr(testsuite.tenantID),
						ObjectID: to.Ptr(testsuite.objectID),
						Permissions: &armkeyvault.Permissions{
							Keys: []*armkeyvault.KeyPermissions{
								to.Ptr(armkeyvault.KeyPermissionsGet),
								to.Ptr(armkeyvault.KeyPermissionsList),
								to.Ptr(armkeyvault.KeyPermissionsCreate),
							},
							Secrets: []*armkeyvault.SecretPermissions{
								to.Ptr(armkeyvault.SecretPermissionsGet),
								to.Ptr(armkeyvault.SecretPermissionsList),
							},
							Certificates: []*armkeyvault.CertificatePermissions{
								to.Ptr(armkeyvault.CertificatePermissionsGet),
								to.Ptr(armkeyvault.CertificatePermissionsList),
								to.Ptr(armkeyvault.CertificatePermissionsCreate),
							},
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	vResp, err := testutil.PollForTest(testsuite.ctx, vPollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vaultName, *vResp.Name)

	// create secret
	fmt.Println("Call operation: Secrets_CreateOrUpdate")
	secretsClient, err := armkeyvault.NewSecretsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	secretName := "go-test-secret2"
	secretResp, err := secretsClient.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vaultName,
		secretName,
		armkeyvault.SecretCreateOrUpdateParameters{
			Properties: &armkeyvault.SecretProperties{
				Attributes: &armkeyvault.SecretAttributes{
					Enabled: to.Ptr(true),
				},
				Value: to.Ptr("sample-secret-value"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(secretName, *secretResp.Name)

	// update secret
	fmt.Println("Call operation: Secrets_Update")
	updateResp, err := secretsClient.Update(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vaultName,
		secretName,
		armkeyvault.SecretPatchParameters{
			Tags: map[string]*string{
				"test": to.Ptr("recording"),
			},
			Properties: &armkeyvault.SecretPatchProperties{
				Value: to.Ptr("sample-secret-value-update"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("recording", *updateResp.Tags["test"])

	// get secret
	fmt.Println("Call operation: Secrets_Get")
	getResp, err := secretsClient.Get(testsuite.ctx, testsuite.resourceGroupName, vaultName, secretName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(secretName, *getResp.Name)

	// list secret
	fmt.Println("Call operation: Secrets_List")
	secretPager := secretsClient.NewListPager(testsuite.resourceGroupName, vaultName, nil)
	testsuite.Require().True(secretPager.More())
}
