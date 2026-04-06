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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type DiskTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	diskName          string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *DiskTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.diskName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "diskname", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *DiskTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestDiskTestSuite(t *testing.T) {
	suite.Run(t, new(DiskTestSuite))
}

// Microsoft.Compute/disks/{diskName}
func (testsuite *DiskTestSuite) TestDisks() {
	var err error
	// From step Disks_CreateOrUpdate
	fmt.Println("Call operation: Disks_CreateOrUpdate")
	disksClient, err := armcompute.NewDisksClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	disksClientCreateOrUpdateResponsePoller, err := disksClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskName, armcompute.Disk{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.DiskProperties{
			CreationData: &armcompute.CreationData{
				CreateOption: to.Ptr(armcompute.DiskCreateOptionEmpty),
			},
			DiskSizeGB: to.Ptr[int32](200),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, disksClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Disks_List
	fmt.Println("Call operation: Disks_List")
	disksClientNewListPager := disksClient.NewListPager(nil)
	for disksClientNewListPager.More() {
		_, err := disksClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Disks_ListByResourceGroup
	fmt.Println("Call operation: Disks_ListByResourceGroup")
	disksClientNewListByResourceGroupPager := disksClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for disksClientNewListByResourceGroupPager.More() {
		_, err := disksClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Disks_Get
	fmt.Println("Call operation: Disks_Get")
	_, err = disksClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskName, nil)
	testsuite.Require().NoError(err)

	// From step Disks_Update
	fmt.Println("Call operation: Disks_Update")
	disksClientUpdateResponsePoller, err := disksClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskName, armcompute.DiskUpdate{
		Properties: &armcompute.DiskUpdateProperties{
			NetworkAccessPolicy: to.Ptr(armcompute.NetworkAccessPolicyAllowAll),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, disksClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Disks_GrantAccess
	fmt.Println("Call operation: Disks_GrantAccess")
	disksClientGrantAccessResponsePoller, err := disksClient.BeginGrantAccess(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskName, armcompute.GrantAccessData{
		Access:            to.Ptr(armcompute.AccessLevelRead),
		DurationInSeconds: to.Ptr[int32](300),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, disksClientGrantAccessResponsePoller)
	testsuite.Require().NoError(err)

	// From step Disks_RevokeAccess
	fmt.Println("Call operation: Disks_RevokeAccess")
	disksClientRevokeAccessResponsePoller, err := disksClient.BeginRevokeAccess(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, disksClientRevokeAccessResponsePoller)
	testsuite.Require().NoError(err)

	// From step Disks_Delete
	fmt.Println("Call operation: Disks_Delete")
	disksClientDeleteResponsePoller, err := disksClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.diskName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, disksClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
