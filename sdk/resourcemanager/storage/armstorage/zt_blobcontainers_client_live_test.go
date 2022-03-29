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
	storageAccountsClient := armstorage.NewAccountsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	scName := "gotestaccount1"
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

	// put container
	blobContainersClient := armstorage.NewBlobContainersClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
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
				PublicAccess: armstorage.PublicAccessContainer.ToPtr(),
				Metadata: map[string]*string{
					"metadata": to.StringPtr("true"),
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
	listPager := blobContainersClient.List(testsuite.resourceGroupName, scName, nil)
	testsuite.Require().NoError(listPager.Err())
	testsuite.Require().True(listPager.NextPage(testsuite.ctx))

	// clear legal hold
	holdResp, err := blobContainersClient.ClearLegalHold(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		blobContainerName,
		armstorage.LegalHold{
			Tags: []*string{
				to.StringPtr("tag1"),
				to.StringPtr("tag2"),
				to.StringPtr("tag3"),
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
				Action:        armstorage.LeaseContainerRequestActionAcquire.ToPtr(),
				LeaseDuration: to.Int32Ptr(-1),
			},
		},
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, leaseResp.RawResponse.StatusCode)

	// break lease
	breakResp, err := blobContainersClient.Lease(
		testsuite.ctx,
		testsuite.resourceGroupName,
		scName,
		blobContainerName,
		&armstorage.BlobContainersClientLeaseOptions{
			Parameters: &armstorage.LeaseContainerRequest{
				Action:  armstorage.LeaseContainerRequestActionBreak.ToPtr(),
				LeaseID: leaseResp.LeaseID,
			},
		},
	)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal("0", *breakResp.LeaseTimeSeconds)

	// delete
	delResp, err := blobContainersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, scName, blobContainerName, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().Equal(200, delResp.RawResponse.StatusCode)
}
