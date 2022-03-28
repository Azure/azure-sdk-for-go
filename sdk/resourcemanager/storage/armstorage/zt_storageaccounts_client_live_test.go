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
	storageAccountsClient := armstorage.NewAccountsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	scName := "gotestaccount"
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

	// check name availability
	check, err := storageAccountsClient.CheckNameAvailability(
		testsuite.ctx,
		armstorage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(scName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
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
					DefaultAction: armstorage.DefaultActionAllow.ToPtr(),
				},
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
	listPager := storageAccountsClient.List(nil)
	testsuite.Require().NoError(listPager.Err())
	testsuite.Require().True(listPager.NextPage(testsuite.ctx))

	// list by resource group
	listByResourceGroup := storageAccountsClient.ListByResourceGroup(testsuite.resourceGroupName, nil)
	testsuite.Require().NoError(listByResourceGroup.Err())
	testsuite.Require().True(listByResourceGroup.NextPage(testsuite.ctx))

	// list keys
	keys, err := storageAccountsClient.ListKeys(testsuite.ctx, testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Greater(len(keys.Keys), 1)

	// revoke user delegation keys
	revokeResp, err := storageAccountsClient.RevokeUserDelegationKeys(testsuite.ctx, testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, revokeResp.RawResponse.StatusCode)

	// regenerate key
	regResp, err := storageAccountsClient.RegenerateKey(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		armstorage.AccountRegenerateKeyParameters{
			KeyName: to.StringPtr("key2"),
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Less(1, len(regResp.Keys))

	// delete
	delResp, err := storageAccountsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}
