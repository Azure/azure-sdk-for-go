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

type StorageContainersTestSuite struct {
	suite.Suite
	ctx                  context.Context
	cred                 azcore.TokenCredential
	options              *arm.ClientOptions
	location             string
	resourceGroupName    string
	subscriptionId       string
	storageContainerName string
}

func (testsuite *StorageContainersTestSuite) SetupSuite() {
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
}

func (testsuite *StorageContainersTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestStorageContainersTestSuite(t *testing.T) {
	suite.Run(t, new(StorageContainersTestSuite))
}

// Test listing storage containers by subscription
func (testsuite *StorageContainersTestSuite) SkipTestStorageContainersListBySubscription() {
	fmt.Println("Call operation: StorageContainers_ListBySubscription")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewStorageContainersClient().NewListBySubscriptionPager(nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test listing storage containers by resource group
func (testsuite *StorageContainersTestSuite) SkipTestStorageContainersListByResourceGroup() {
	fmt.Println("Call operation: StorageContainers_ListByResourceGroup")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	pager := clientFactory.NewStorageContainersClient().NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for pager.More() {
		result, err := pager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		testsuite.Require().NotNil(result.Value)
		break
	}
}

// Test getting a storage container
func (testsuite *StorageContainersTestSuite) SkipTestStorageContainersGet() {
	fmt.Println("Call operation: StorageContainers_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewStorageContainersClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageContainerName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a storage container
func (testsuite *StorageContainersTestSuite) SkipTestStorageContainersCreateOrUpdate() {
	fmt.Println("Call operation: StorageContainers_CreateOrUpdate")
	// Requires proper storage configuration
}

// Test updating a storage container
func (testsuite *StorageContainersTestSuite) SkipTestStorageContainersUpdate() {
	fmt.Println("Call operation: StorageContainers_Update")
	// Requires existing storage container
}

// Test deleting a storage container
func (testsuite *StorageContainersTestSuite) SkipTestStorageContainersDelete() {
	fmt.Println("Call operation: StorageContainers_Delete")
	// Requires existing storage container
}
