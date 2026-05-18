// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdiscovery_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
func (testsuite *StorageContainersTestSuite) TestStorageContainersListBySubscription() {
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
func (testsuite *StorageContainersTestSuite) TestStorageContainersListByResourceGroup() {
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
func (testsuite *StorageContainersTestSuite) TestStorageContainersGet() {
	fmt.Println("Call operation: StorageContainers_Get")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	_, err = clientFactory.NewStorageContainersClient().Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.storageContainerName, nil)
	testsuite.Require().NoError(err)
}

// Test creating a storage container
func (testsuite *StorageContainersTestSuite) TestStorageContainersCreateOrUpdate() {
	fmt.Println("Call operation: StorageContainers_CreateOrUpdate")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	storageContainersClient := clientFactory.NewStorageContainersClient()
	poller, err := storageContainersClient.BeginCreateOrUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.storageContainerName,
		armdiscovery.StorageContainer{
			Location: to.Ptr(testsuite.location),
			Properties: &armdiscovery.StorageContainerProperties{
				StorageStore: &armdiscovery.AzureStorageBlobStore{
					Kind:             to.Ptr(armdiscovery.StorageStoreTypeAzureStorageBlob),
					StorageAccountID: to.Ptr("/subscriptions/" + testsuite.subscriptionId + "/resourceGroups/" + testsuite.resourceGroupName + "/providers/Microsoft.Storage/storageAccounts/mytststr"),
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	sc, err := poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(sc.ID)
	fmt.Println("Created storage container:", *sc.Name)
}

// Test updating a storage container
func (testsuite *StorageContainersTestSuite) TestStorageContainersUpdate() {
	fmt.Println("Call operation: StorageContainers_Update")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	poller, err := clientFactory.NewStorageContainersClient().BeginUpdate(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.storageContainerName,
		armdiscovery.StorageContainer{
			Tags: map[string]*string{
				"updated": to.Ptr("true"),
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)
	result, err := poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	testsuite.Require().NotNil(result.ID)
	fmt.Println("Updated storage container:", *result.Name)
}

// Test deleting a storage container
func (testsuite *StorageContainersTestSuite) TestStorageContainersDelete() {
	fmt.Println("Call operation: StorageContainers_Delete")
	clientFactory, err := armdiscovery.NewClientFactory(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	poller, err := clientFactory.NewStorageContainersClient().BeginDelete(
		testsuite.ctx,
		testsuite.resourceGroupName,
		testsuite.storageContainerName,
		nil,
	)
	testsuite.Require().NoError(err)
	_, err = poller.PollUntilDone(testsuite.ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	testsuite.Require().NoError(err)
	fmt.Println("Deleted storage container:", testsuite.storageContainerName)
}
