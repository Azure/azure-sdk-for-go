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

type DiskAccessTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	diskAccessName    string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *DiskAccessTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.diskAccessName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "diskaccess", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *DiskAccessTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDiskAccessTestSuite(t *testing.T) {
	suite.Run(t, new(DiskAccessTestSuite))
}

// Microsoft.Compute/diskAccesses/{diskAccessName}
func (testsuite *DiskAccessTestSuite) TestDiskAccesses() {
	var err error
	// From step DiskAccesses_CreateOrUpdate
	fmt.Println("Call operation: DiskAccesses_CreateOrUpdate")
	diskAccessesClient, err := armcompute.NewDiskAccessesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	diskAccessesClientCreateOrUpdateResponsePoller, err := diskAccessesClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskAccessName, armcompute.DiskAccess{
		Location: to.Ptr(testsuite.location),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, diskAccessesClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DiskAccesses_List
	fmt.Println("Call operation: DiskAccesses_List")
	diskAccessesClientNewListPager := diskAccessesClient.NewListPager(nil)
	for diskAccessesClientNewListPager.More() {
		_, err := diskAccessesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DiskAccesses_ListPrivateEndpointConnections
	fmt.Println("Call operation: DiskAccesses_ListPrivateEndpointConnections")
	diskAccessesClientNewListPrivateEndpointConnectionsPager := diskAccessesClient.NewListPrivateEndpointConnectionsPager(testsuite.resourceGroupName, testsuite.diskAccessName, nil)
	for diskAccessesClientNewListPrivateEndpointConnectionsPager.More() {
		_, err := diskAccessesClientNewListPrivateEndpointConnectionsPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DiskAccesses_GetPrivateLinkResources
	fmt.Println("Call operation: DiskAccesses_GetPrivateLinkResources")
	_, err = diskAccessesClient.GetPrivateLinkResources(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskAccessName, nil)
	testsuite.Require().NoError(err)

	// From step DiskAccesses_ListByResourceGroup
	fmt.Println("Call operation: DiskAccesses_ListByResourceGroup")
	diskAccessesClientNewListByResourceGroupPager := diskAccessesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for diskAccessesClientNewListByResourceGroupPager.More() {
		_, err := diskAccessesClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step DiskAccesses_Get
	fmt.Println("Call operation: DiskAccesses_Get")
	_, err = diskAccessesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskAccessName, nil)
	testsuite.Require().NoError(err)

	// From step DiskAccesses_Update
	fmt.Println("Call operation: DiskAccesses_Update")
	diskAccessesClientUpdateResponsePoller, err := diskAccessesClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskAccessName, armcompute.DiskAccessUpdate{
		Tags: map[string]*string{
			"department": to.Ptr("Development"),
			"project":    to.Ptr("PrivateEndpoints"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, diskAccessesClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step DiskAccesses_Delete
	fmt.Println("Call operation: DiskAccesses_Delete")
	diskAccessesClientDeleteResponsePoller, err := diskAccessesClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskAccessName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, diskAccessesClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
