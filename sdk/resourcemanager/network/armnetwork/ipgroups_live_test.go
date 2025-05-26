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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"github.com/stretchr/testify/suite"
)

type IpGroupsTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	ipGroupsName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *IpGroupsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.ipGroupsName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "ipgroupsna", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *IpGroupsTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestIpGroupsTestSuite(t *testing.T) {
	suite.Run(t, new(IpGroupsTestSuite))
}

// Microsoft.Network/ipGroups/{ipGroupsName}
func (testsuite *IpGroupsTestSuite) TestIpGroups() {
	var err error
	// From step IpGroups_CreateOrUpdate
	fmt.Println("Call operation: IPGroups_CreateOrUpdate")
	iPGroupsClient, err := armnetwork.NewIPGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	iPGroupsClientCreateOrUpdateResponsePoller, err := iPGroupsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.ipGroupsName, armnetwork.IPGroup{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
		},
		Properties: &armnetwork.IPGroupPropertiesFormat{
			IPAddresses: []*string{
				to.Ptr("13.64.39.16/32"),
				to.Ptr("40.74.146.80/31"),
				to.Ptr("40.74.147.32/28")},
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, iPGroupsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step IpGroups_List
	fmt.Println("Call operation: IPGroups_List")
	iPGroupsClientNewListPager := iPGroupsClient.NewListPager(nil)
	for iPGroupsClientNewListPager.More() {
		_, err := iPGroupsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IpGroups_ListByResourceGroup
	fmt.Println("Call operation: IPGroups_ListByResourceGroup")
	iPGroupsClientNewListByResourceGroupPager := iPGroupsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for iPGroupsClientNewListByResourceGroupPager.More() {
		_, err := iPGroupsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IpGroups_Get
	fmt.Println("Call operation: IPGroups_Get")
	_, err = iPGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.ipGroupsName, &armnetwork.IPGroupsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step IpGroups_UpdateGroups
	fmt.Println("Call operation: IPGroups_UpdateGroups")
	_, err = iPGroupsClient.UpdateGroups(testsuite.ctx, testsuite.resourceGroupName, testsuite.ipGroupsName, armnetwork.TagsObject{
		Tags: map[string]*string{
			"key1": to.Ptr("value1"),
			"key2": to.Ptr("value2"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IpGroups_Delete
	fmt.Println("Call operation: IPGroups_Delete")
	iPGroupsClientDeleteResponsePoller, err := iPGroupsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.ipGroupsName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, iPGroupsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
