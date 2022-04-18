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

type BlobContainersClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *BlobContainersClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = testutil.GetEnv("LOCATION", "eastus")
	testsuite.subscriptionID = testutil.GetEnv("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), "sdk/resourcemanager/storage/armstorage/testdata")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *BlobContainersClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestBlobContainersClient(t *testing.T) {
	suite.Run(t, new(BlobContainersClientTestSuite))
}

func (testsuite *BlobContainersClientTestSuite) TestBlobContainersCRUD() {
	// create storage account
	storageAccountsClient, err := armstorage.NewAccountsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	scName := "gotestaccount1"
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

	// put container
	blobContainersClient, err := armstorage.NewBlobContainersClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	blobContainerName := "go-test-container"
	blobResp, err := blobContainersClient.Create(testsuite.ctx, testsuite.resourceGroupName, scName, blobContainerName, armstorage.BlobContainer{}, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(blobContainerName, *blobResp.Name)

	// update
	updateResp, err := blobContainersClient.Update(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		blobContainerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{
				PublicAccess: to.Ptr(armstorage.PublicAccessContainer),
				Metadata: map[string]*string{
					"metadata": to.Ptr("true"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(blobContainerName, *updateResp.Name)

	// get
	getResp, err := blobContainersClient.Get(testsuite.ctx, testsuite.resourceGroupName, scName, blobContainerName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(blobContainerName, *getResp.Name)

	// list
	listPager := blobContainersClient.NewListPager(testsuite.resourceGroupName, scName, nil)
	testsuite.Require().True(listPager.More())

	// clear legal hold
	holdResp, err := blobContainersClient.ClearLegalHold(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		blobContainerName,
		armstorage.LegalHold{
			Tags: []*string{
				to.Ptr("tag1"),
				to.Ptr("tag2"),
				to.Ptr("tag3"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(0, len(holdResp.Tags))

	// acquire lease
	leaseResp, err := blobContainersClient.Lease(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersClientLeaseOptions{
			Parameters: &armstorage.LeaseContainerRequest{
				Action:        to.Ptr(armstorage.LeaseContainerRequestActionAcquire),
				LeaseDuration: to.Ptr[int32](-1),
			},
		},
	)
	testsuite.Require().NoError(err)

	// break lease
	breakResp, err := blobContainersClient.Lease(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersClientLeaseOptions{
			Parameters: &armstorage.LeaseContainerRequest{
				Action:  to.Ptr(armstorage.LeaseContainerRequestActionBreak),
				LeaseID: leaseResp.LeaseID,
			},
		},
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("0", *breakResp.LeaseTimeSeconds)

	// delete
	_, err = blobContainersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, scName, blobContainerName, nil)
	testsuite.Require().NoError(err)
}
