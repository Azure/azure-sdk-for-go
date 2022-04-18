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

type EncryptionScopesClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *EncryptionScopesClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/storage/armstorage/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *EncryptionScopesClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestEncryptionScopesClient(t *testing.T) {
	suite.Run(t, new(EncryptionScopesClientTestSuite))
}

func (testsuite *EncryptionScopesClientTestSuite) TestEncryptionScopesCRUD() {
	// create storage account
	storageAccountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scName := "gotestaccount2"
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

	// put file services
	encryptionScopesClient, err := armstorage.NewEncryptionScopesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	encryptionScopeName := "go-test-encryption"
	putResp, err := encryptionScopesClient.Put(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		encryptionScopeName,
		armstorage.EncryptionScope{
			EncryptionScopeProperties: &armstorage.EncryptionScopeProperties{
				Source: to.Ptr(armstorage.EncryptionScopeSourceMicrosoftStorage),
				State:  to.Ptr(armstorage.EncryptionScopeStateEnabled),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(encryptionScopeName, *putResp.Name)

	// get
	getResp, err := encryptionScopesClient.Get(testsuite.ctx, testsuite.resourceGroupName, scName, encryptionScopeName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(encryptionScopeName, *getResp.Name)

	// list
	listPager := encryptionScopesClient.NewListPager(testsuite.resourceGroupName, scName, nil)
	testsuite.Require().True(listPager.More())

	// patch
	patchResp, err := encryptionScopesClient.Patch(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		encryptionScopeName,
		armstorage.EncryptionScope{
			EncryptionScopeProperties: &armstorage.EncryptionScopeProperties{
				Source: to.Ptr(armstorage.EncryptionScopeSourceMicrosoftStorage),
				State:  to.Ptr(armstorage.EncryptionScopeStateEnabled),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(encryptionScopeName, *patchResp.Name)
}
