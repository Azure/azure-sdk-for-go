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

type StorageAccountsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *StorageAccountsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/storage/armstorage/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *StorageAccountsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestStorageAccountsClient(t *testing.T) {
	suite.Run(t, new(StorageAccountsClientTestSuite))
}

func (testsuite *StorageAccountsClientTestSuite) TestStorageAccountsCRUD() {
	// create storage account
	storageAccountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scName := "gotestaccount"
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

	// check name availability
	check, err := storageAccountsClient.CheckNameAvailability(
		testsuite.ctx,
		armstorage.AccountCheckNameAvailabilityParameters{
			Name: to.Ptr(scName),
			Type: to.Ptr("Microsoft.Storage/storageAccounts"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().False(*check.NameAvailable)

	// update
	updateResp, err := storageAccountsClient.Update(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		armstorage.AccountUpdateParameters{
			Properties: &armstorage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: &armstorage.NetworkRuleSet{
					DefaultAction: to.Ptr(armstorage.DefaultActionAllow),
				},
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
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(scName, *updateResp.Name)

	// get properties
	getResp, err := storageAccountsClient.GetProperties(testsuite.ctx, testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(scName, *getResp.Name)

	// list
	listPager := storageAccountsClient.NewListPager(nil)
	testsuite.Require().True(listPager.More())

	// list by resource group
	listByResourceGroup := storageAccountsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	testsuite.Require().True(listByResourceGroup.More())

	// list keys
	keys, err := storageAccountsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Greater(len(keys.Keys), 1)

	// revoke user delegation keys
	_, err = storageAccountsClient.RevokeUserDelegationKeys(testsuite.ctx, testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(err)

	// regenerate key
	regResp, err := storageAccountsClient.RegenerateKey(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		armstorage.AccountRegenerateKeyParameters{
			KeyName: to.Ptr("key2"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Less(1, len(regResp.Keys))

	// delete
	_, err = storageAccountsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(err)
}
