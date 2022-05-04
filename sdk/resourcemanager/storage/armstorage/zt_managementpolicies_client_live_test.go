//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armstorage_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/stretchr/testify/suite"
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
	storageAccountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scName := "gotestaccount3"
	pollerResp, err := storageAccountsClient.BeginCreate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		armstorage.AccountCreateParameters{
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNameStandardGRS),
			},
			Kind:     to.Ptr(armstorage.KindStorageV2),
			Location: to.Ptr(testsuite.location),
			Properties: &armstorage.AccountPropertiesCreateParameters{
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
						Blob: &armstorage.EncryptionService{
							KeyType: to.Ptr(armstorage.KeyTypeAccount),
							Enabled: to.Ptr(true),
						},
					},
					KeySource: to.Ptr(armstorage.KeySourceMicrosoftStorage),
				},
			},
			Tags: map[string]*string{
				"key1": to.Ptr("value1"),
				"key2": to.Ptr("value2"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	resp, err := testutil.PollForTest(testsuite.ctx, pollerResp)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(scName, *resp.Name)

	// create management policy
	managementPoliciesClient, err := armstorage.NewManagementPoliciesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
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
							Enabled: to.Ptr(true),
							Name:    to.Ptr("olcmtest"),
							Type:    to.Ptr(armstorage.RuleTypeLifecycle),
							Definition: &armstorage.ManagementPolicyDefinition{
								Filters: &armstorage.ManagementPolicyFilter{
									BlobTypes: []*string{
										to.Ptr("blockBlob"),
									},
									PrefixMatch: []*string{
										to.Ptr("olcmtestcontainer"),
									},
								},
								Actions: &armstorage.ManagementPolicyAction{
									BaseBlob: &armstorage.ManagementPolicyBaseBlob{
										TierToCool: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Ptr[float32](30),
										},
										TierToArchive: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Ptr[float32](90),
										},
										Delete: &armstorage.DateAfterModification{
											DaysAfterModificationGreaterThan: to.Ptr[float32](1000),
										},
									},
									Snapshot: &armstorage.ManagementPolicySnapShot{
										Delete: &armstorage.DateAfterCreation{
											DaysAfterCreationGreaterThan: to.Ptr[float32](30),
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
	_, err = managementPoliciesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, scName, "default", nil)
	testsuite.Require().NoError(err)
}
