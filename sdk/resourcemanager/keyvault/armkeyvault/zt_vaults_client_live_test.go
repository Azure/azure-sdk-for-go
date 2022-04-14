//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/stretchr/testify/suite"
)

type VaultsClientTestSuite struct {
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

func (testsuite *VaultsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.tenantID = testutil.GetEnv("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.objectID = testutil.GetEnv("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/keyvault/armkeyvault/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *VaultsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestVaultsClient(t *testing.T) {
	suite.Run(t, new(VaultsClientTestSuite))
}

func (testsuite *VaultsClientTestSuite) TestVaultsCRUD() {
	// create vault
	vaultsClient, err := armkeyvault.NewVaultsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vaultName := "go-test-vault1"
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

	// create vault
	check, err := vaultsClient.CheckNameAvailability(
		testsuite.ctx,
		armkeyvault.VaultCheckNameAvailabilityParameters{
			Name: to.Ptr(vaultName),
			Type: to.Ptr("Microsoft.KeyVault/vaults"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().False(*check.NameAvailable)

	// get vault
	getResp, err := vaultsClient.Get(testsuite.ctx, testsuite.resourceGroupName, vaultName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vaultName, *getResp.Name)

	// update vault
	updateResp, err := vaultsClient.Update(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vaultName,
		armkeyvault.VaultPatchParameters{
			Tags: map[string]*string{
				"test": to.Ptr("recording"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("recording", *updateResp.Tags["test"])

	// list vault deleted
	deletedPager := vaultsClient.ListDeleted(nil)
	testsuite.Require().True(deletedPager.More())

	// delete vault
	_, err = vaultsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, vaultName, nil)
	testsuite.Require().NoError(err)

	// get deleted vault
	deletedResp, err := vaultsClient.GetDeleted(testsuite.ctx, vaultName, testsuite.location, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vaultName, *deletedResp.Name)

	// purge deleted vault
	purgePollerResp, err := vaultsClient.BeginPurgeDeleted(testsuite.ctx, vaultName, testsuite.location, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, purgePollerResp)
	testsuite.Require().NoError(err)
}
