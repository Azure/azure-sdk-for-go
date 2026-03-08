// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/discovery/armdiscovery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type StorageAssetsTestSuite struct {
	suite.Suite
	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	location             string
	resourceGroupName    string
	subscriptionId       string
	storageContainerName string
	storageAssetName     string
}

func (testsuite *StorageAssetsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.location = recording.GetEnvVariable("LOCATION", ResourceLocation)
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "newapiversiontest")
	testsuite.storageContainerName = "test-storage-container"
	testsuite.storageAssetName = "test-storage-asset"
}

func (testsuite *StorageAssetsTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestStorageAssetsTestSuite(t *testing.T) {
	suite.Run(t, new(StorageAssetsTestSuite))
}

// Test listing storage assets by storage container
func (testsuite *StorageAssetsTestSuite) TestStorageAssetsListByStorageContainer() {
	fmt.Println("Call operation: StorageAssets_ListByStorageContainer")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewStorageAssetsClient().NewListByStorageContainerPager(testsuite.resourceGroupName, testsuite.storageContainerName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a storage asset
func (testsuite *StorageAssetsTestSuite) TestStorageAssetsGet() {
	fmt.Println("Call operation: StorageAssets_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewStorageAssetsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageContainerName, testsuite.storageAssetName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a storage asset
func (testsuite *StorageAssetsTestSuite) TestStorageAssetsCreateOrUpdate() {
	fmt.Println("Call operation: StorageAssets_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	storageAssetsClient := clientFactory.NewStorageAssetsClient()
	poller, err := storageAssetsClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.storageContainerName,
		testsuite.storageAssetName,
		armdiscovery.StorageAsset{
			Location: to.Ptr(testsuite.location),
			Properties: &armdiscovery.StorageAssetProperties{
				Description: to.Ptr("Test storage asset for SDK validation"),
				Path:        to.Ptr("data/test-assets"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	sa, err := poller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(sa.ID)
	fmt.Println("Created storage asset:", *sa.Name)
}

// Test updating a storage asset
func (testsuite *StorageAssetsTestSuite) TestStorageAssetsUpdate() {
	fmt.Println("Call operation: StorageAssets_Update")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	poller, err := clientFactory.NewStorageAssetsClient().BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.storageContainerName,
		testsuite.storageAssetName,
		armdiscovery.StorageAsset{
			Tags: map[string]*string{
				"updated": to.Ptr("true"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	result, err := poller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(result.ID)
	fmt.Println("Updated storage asset:", *result.Name)
}

// Test deleting a storage asset
func (testsuite *StorageAssetsTestSuite) TestStorageAssetsDelete() {
	fmt.Println("Call operation: StorageAssets_Delete")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	poller, err := clientFactory.NewStorageAssetsClient().BeginDelete(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.storageContainerName,
		testsuite.storageAssetName,
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = poller.PollUntilDone(testsuite.ctx, nil)
	testsuite.Require().NoError(err)
	fmt.Println("Deleted storage asset:", testsuite.storageAssetName)
}
