//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ManagementPoliciesClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *ManagementPoliciesClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/storage/armstorage/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *ManagementPoliciesClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestManagementPoliciesClient(t *testing.T) {
	suite.Run(t, new(ManagementPoliciesClientTestSuite))
}

func (testsuite *ManagementPoliciesClientTestSuite) TestManagementPoliciesCRUD() {
	// create storage account
	storageAccountsClient := armstorage.NewAccountsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	scName := "gotestaccount3"
	pollerResp, err := storageAccountsClient.BeginCreate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		armstorage.AccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: armstorage.SKUNameStandardGRS.ToPtr(),
			},
			Kind:     armstorage.KindStorageV2.ToPtr(),
			Location: to.StringPtr(testsuite.location),
			Properties: &armstorage.AccountPropertiesCreateParameters{
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
	testsuite.Require().NoError(err)
	var resp armstorage.AccountsClientCreateResponse
	if recording.GetRecordMode() == recording.PlaybackMode {
		for {
			_, err = pollerResp.Poller.Poll(testsuite.ctx)
			testsuite.Require().NoError(err)
			if pollerResp.Poller.Done() {
				resp, err = pollerResp.Poller.FinalResponse(testsuite.ctx)
				testsuite.Require().NoError(err)
				break
			}
		}
	} else {
		resp, err = pollerResp.PollUntilDone(testsuite.ctx, 30*time.Second)
		testsuite.Require().NoError(err)
	}
	testsuite.Require().Equal(scName, *resp.Name)

	// create management policy
	managementPoliciesClient := armstorage.NewManagementPoliciesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	mpResp, err := managementPoliciesClient.CreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		armstorage.ManagementPolicyNameDefault,
		armstorage.ManagementPolicy{
			Properties: &armstorage.ManagementPolicyProperties{
				Policy: &armstorage.ManagementPolicySchema{
					Rules: []*armstorage.ManagementPolicyRule{
						{
							Enabled: to.BoolPtr(true),
							Name:    to.StringPtr("olcmtest"),
							Type:    armstorage.RuleTypeLifecycle.ToPtr(),
							Definition: &armstorage.ManagementPolicyDefinition{
								Filters: &armstorage.ManagementPolicyFilter{
									BlobTypes: []*string{
										to.StringPtr("blockBlob"),
									},
									PrefixMatch: []*string{
										to.StringPtr("olcmtestcontainer"),
									},
								},
								Actions: &armstorage.ManagementPolicyAction{
									BaseBlob: &armstorage.ManagementPolicyBaseBlob{
										TierToCool: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(30),
										},
										TierToArchive: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(90),
										},
										Delete: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Float32Ptr(1000),
										},
									},
									Snapshot: &armstorage.ManagementPolicySnapShot{
										Delete: &armstorage.DateAfterCreation{
											DaysAfterCreationGreaterThan: to.Float32Ptr(30),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("DefaultManagementPolicy", *mpResp.Name)

	// get management policy
	getResp, err := managementPoliciesClient.Get(testsuite.ctx, testsuite.resourceGroupName, scName, "default", nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("DefaultManagementPolicy", *getResp.Name)

	// delete management policy
	delResp, err := managementPoliciesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, scName, "default", nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}
