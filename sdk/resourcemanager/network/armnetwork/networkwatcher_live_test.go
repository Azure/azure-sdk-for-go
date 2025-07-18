//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armnetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v8"
	"github.com/stretchr/testify/suite"
)

type NetworkWatcherTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	networkWatcherName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *NetworkWatcherTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.networkWatcherName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "networkwat", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *NetworkWatcherTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TTestNetworkWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkWatcherTestSuite))
}

// Microsoft.Network/networkWatchers/{networkWatcherName}
func (testsuite *NetworkWatcherTestSuite) TestNetworkWatchers() {
	var err error

	// From step NetworkWatchers_CreateOrUpdate
	watchersClient, err := armnetwork.NewWatchersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	fmt.Println("Call operation: NetworkWatchers_CreateOrUpdate")
	_, err = watchersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkWatcherName, armnetwork.Watcher{
		Location:   to.Ptr(testsuite.location),
		Properties: &armnetwork.WatcherPropertiesFormat{},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NetworkWatchers_ListAll
	fmt.Println("Call operation: NetworkWatchers_ListAll")
	watchersClientNewListAllPager := watchersClient.NewListAllPager(nil)
	for watchersClientNewListAllPager.More() {
		_, err := watchersClientNewListAllPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkWatchers_List
	fmt.Println("Call operation: NetworkWatchers_List")
	watchersClientNewListPager := watchersClient.NewListPager(testsuite.resourceGroupName, nil)
	for watchersClientNewListPager.More() {
		_, err := watchersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step NetworkWatchers_Get
	fmt.Println("Call operation: NetworkWatchers_Get")
	_, err = watchersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkWatcherName, nil)
	testsuite.Require().NoError(err)

	// From step NetworkWatchers_UpdateTags
	fmt.Println("Call operation: NetworkWatchers_UpdateTags")
	_, err = watchersClient.UpdateTags(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkWatcherName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"tag1": to.Ptr("value1"),
			"tag2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step NetworkWatchers_Delete
	fmt.Println("Call operation: NetworkWatchers_Delete")
	watchersClientDeleteResponsePoller, err := watchersClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.networkWatcherName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, watchersClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
