// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
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

	// Add EUAP redirect policy
	euapOptions := GetEUAPClientOptions()
	testsuite.options.PerCallPolicies = append(testsuite.options.PerCallPolicies, euapOptions.PerCallPolicies...)

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
func (testsuite *StorageAssetsTestSuite) SkipTestStorageAssetsListByStorageContainer() {
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
func (testsuite *StorageAssetsTestSuite) SkipTestStorageAssetsGet() {
	fmt.Println("Call operation: StorageAssets_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewStorageAssetsClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageContainerName, testsuite.storageAssetName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a storage asset
func (testsuite *StorageAssetsTestSuite) SkipTestStorageAssetsCreateOrUpdate() {
	fmt.Println("Call operation: StorageAssets_CreateOrUpdate")
	// Requires storage container with proper configuration
}

// Test updating a storage asset
func (testsuite *StorageAssetsTestSuite) SkipTestStorageAssetsUpdate() {
	fmt.Println("Call operation: StorageAssets_Update")
	// Requires existing storage asset
}

// Test deleting a storage asset
func (testsuite *StorageAssetsTestSuite) SkipTestStorageAssetsDelete() {
	fmt.Println("Call operation: StorageAssets_Delete")
	// Requires existing storage asset
}
