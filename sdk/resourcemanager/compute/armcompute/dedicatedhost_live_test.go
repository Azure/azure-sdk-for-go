//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v8"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type DedicatedHostTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	hostGroupName     string
	hostName          string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *DedicatedHostTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.hostGroupName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "hostgroupn", 16, false)
	testsuite.hostName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "hostname", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *DedicatedHostTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDedicatedHostTestSuite(t *testing.T) {
	suite.Run(t, new(DedicatedHostTestSuite))
}

func (testsuite *DedicatedHostTestSuite) Prepare() {
	var err error
	// From step DedicatedHostGroups_CreateOrUpdate
	fmt.Println("Call operation: DedicatedHostGroups_CreateOrUpdate")
	dedicatedHostGroupsClient, err := armcompute.NewDedicatedHostGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = dedicatedHostGroupsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, armcompute.DedicatedHostGroup{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"department": to.Ptr("finance"),
		},
		Properties: &armcompute.DedicatedHostGroupProperties{
			PlatformFaultDomainCount:  to.Ptr[int32](3),
			SupportAutomaticPlacement: to.Ptr(true),
		},
		Zones: []*string{
			to.Ptr("1")},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/hostGroups
func (testsuite *DedicatedHostTestSuite) TestDedicatedHostGroups() {
	var err error
	// From step DedicatedHostGroups_ListBySubscription
	fmt.Println("Call operation: DedicatedHostGroups_ListBySubscription")
	dedicatedHostGroupsClient, err := armcompute.NewDedicatedHostGroupsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dedicatedHostGroupsClientNewListBySubscriptionPager := dedicatedHostGroupsClient.NewListBySubscriptionPager(nil)
	for dedicatedHostGroupsClientNewListBySubscriptionPager.More() {
		_, err := dedicatedHostGroupsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DedicatedHostGroups_ListByResourceGroup
	fmt.Println("Call operation: DedicatedHostGroups_ListByResourceGroup")
	dedicatedHostGroupsClientNewListByResourceGroupPager := dedicatedHostGroupsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for dedicatedHostGroupsClientNewListByResourceGroupPager.More() {
		_, err := dedicatedHostGroupsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DedicatedHostGroups_Get
	fmt.Println("Call operation: DedicatedHostGroups_Get")
	_, err = dedicatedHostGroupsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, &armcompute.DedicatedHostGroupsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step DedicatedHostGroups_Update
	fmt.Println("Call operation: DedicatedHostGroups_Update")
	_, err = dedicatedHostGroupsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, armcompute.DedicatedHostGroupUpdate{}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Compute/hostGroups/hosts
func (testsuite *DedicatedHostTestSuite) TestDedicatedHosts() {
	var err error
	// From step DedicatedHosts_CreateOrUpdate
	fmt.Println("Call operation: DedicatedHosts_CreateOrUpdate")
	dedicatedHostsClient, err := armcompute.NewDedicatedHostsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	dedicatedHostsClientCreateOrUpdateResponsePoller, err := dedicatedHostsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, testsuite.hostName, armcompute.DedicatedHost{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"department": to.Ptr("HR"),
		},
		Properties: &armcompute.DedicatedHostProperties{
			PlatformFaultDomain: to.Ptr[int32](1),
		},
		SKU: &armcompute.SKU{
			Name: to.Ptr("DSv3-Type1"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dedicatedHostsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DedicatedHosts_ListByHostGroup
	fmt.Println("Call operation: DedicatedHosts_ListByHostGroup")
	dedicatedHostsClientNewListByHostGroupPager := dedicatedHostsClient.NewListByHostGroupPager(testsuite.resourceGroupName, testsuite.hostGroupName, nil)
	for dedicatedHostsClientNewListByHostGroupPager.More() {
		_, err := dedicatedHostsClientNewListByHostGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DedicatedHosts_Get
	fmt.Println("Call operation: DedicatedHosts_Get")
	_, err = dedicatedHostsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, testsuite.hostName, &armcompute.DedicatedHostsClientGetOptions{Expand: nil})
	testsuite.Require().NoError(err)

	// From step DedicatedHosts_Update
	fmt.Println("Call operation: DedicatedHosts_Update")
	dedicatedHostsClientUpdateResponsePoller, err := dedicatedHostsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, testsuite.hostName, armcompute.DedicatedHostUpdate{}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dedicatedHostsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DedicatedHosts_Restart
	fmt.Println("Call operation: DedicatedHosts_Restart")
	dedicatedHostsClientRestartResponsePoller, err := dedicatedHostsClient.BeginRestart(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, testsuite.hostName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dedicatedHostsClientRestartResponsePoller)
	testsuite.Require().NoError(err)

	// From step DedicatedHosts_Delete
	fmt.Println("Call operation: DedicatedHosts_Delete")
	dedicatedHostsClientDeleteResponsePoller, err := dedicatedHostsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.hostGroupName, testsuite.hostName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, dedicatedHostsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
