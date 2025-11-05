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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v3"
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
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.tenantID = recording.GetEnvVariable("AZURE_TENANT_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.objectID = recording.GetEnvVariable("AZURE_OBJECT_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *VaultsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestVaultsClient(t *testing.T) {
	suite.Run(t, new(VaultsClientTestSuite))
}

func (testsuite *VaultsClientTestSuite) TestVaultsCRUD() {
	// create vault
	fmt.Println("Call operation: Vaults_CreateOrUpdate")
	vaultsClient, err := armkeyvault.NewVaultsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	vaultName := "go-test-vault-1"
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
	_, err = testutil.PollForTest(testsuite.ctx, vPollerResp)
	testsuite.Require().NoError(err)

	// create vault
	fmt.Println("Call operation: Vaults_CheckNameAvailability")
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
	fmt.Println("Call operation: Vaults_Get")
	_, err = vaultsClient.Get(testsuite.ctx, testsuite.resourceGroupName, vaultName, nil)
	testsuite.Require().NoError(err)

	// update vault
	fmt.Println("Call operation: Vaults_Update")
	_, err = vaultsClient.Update(
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

	// list vault deleted
	fmt.Println("Call operation: Vaults_ListDeleted")
	deletedPager := vaultsClient.NewListDeletedPager(nil)
	testsuite.Require().True(deletedPager.More())

	// delete vault
	fmt.Println("Call operation: Vaults_Delete")
	_, err = vaultsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, vaultName, nil)
	testsuite.Require().NoError(err)

	// get deleted vault
	fmt.Println("Call operation: Vaults_GetDeleted")
	_, err = vaultsClient.GetDeleted(testsuite.ctx, vaultName, testsuite.location, nil)
	testsuite.Require().NoError(err)

	// purge deleted vault
	fmt.Println("Call operation: Vaults_PurgeDeleted")
	purgePollerResp, err := vaultsClient.BeginPurgeDeleted(testsuite.ctx, vaultName, testsuite.location, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, purgePollerResp)
	testsuite.Require().NoError(err)
}
