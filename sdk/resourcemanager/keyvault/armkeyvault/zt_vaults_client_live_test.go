//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armkeyvault_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
	vaultsClient := armkeyvault.NewVaultsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	vaultName := "go-test-vault1"
	vPollerResp, err := vaultsClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		vaultName,
		armkeyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(testsuite.location),
			Properties: &armkeyvault.VaultProperties{
				SKU: &armkeyvault.SKU{
					Family: armkeyvault.SKUFamilyA.ToPtr(),
					Name:   armkeyvault.SKUNameStandard.ToPtr(),
				},
				TenantID: to.StringPtr(testsuite.tenantID),
				AccessPolicies: []*armkeyvault.AccessPolicyEntry{
					{
						TenantID: to.StringPtr(testsuite.tenantID),
						ObjectID: to.StringPtr(testsuite.objectID),
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
	testsuite.Require().NoError(err)
	var vResp armkeyvault.VaultsClientCreateOrUpdateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = vPollerResp.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if vPollerResp.Poller.Done() {
				vResp, err = vPollerResp.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		vResp, err = vPollerResp.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(vaultName, *vResp.Name)

	// create vault
	check, err := vaultsClient.CheckNameAvailability(
		testsuite.ctx,
		armkeyvault.VaultCheckNameAvailabilityParameters{
			Name: to.StringPtr(vaultName),
			Type: to.StringPtr("Microsoft.KeyVault/vaults"),
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
				"test": to.StringPtr("recording"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("recording", *updateResp.Tags["test"])

	// list vault deleted
	deletedPager := vaultsClient.ListDeleted(nil)
	testsuite.Require().NoError(deletedPager.Err())
	testsuite.Require().True(deletedPager.NextPage(testsuite.ctx))

	// delete vault
	delResp, err := vaultsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, vaultName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)

	// get deleted vault
	deletedResp, err := vaultsClient.GetDeleted(testsuite.ctx, vaultName, testsuite.location, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(vaultName, *deletedResp.Name)

	// purge deleted vault
	purgePollerResp, err := vaultsClient.BeginPurgeDeleted(testsuite.ctx, vaultName, testsuite.location, nil)
	testsuite.Require().NoError(err)
	var purgeResp armkeyvault.VaultsClientPurgeDeletedResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = purgePollerResp.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if purgePollerResp.Poller.Done() {
				purgeResp, err = purgePollerResp.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		purgeResp, err = purgePollerResp.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(200, purgeResp.RawResponse.StatusCode)
}
