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

type SnapshotTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	snapshotName      string
	location          string
	resourceGroupName string
	subscriptionId    string
}

func (testsuite *SnapshotTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.snapshotName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "snapshotna", 16, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *SnapshotTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestSnapshotTestSuite(t *testing.T) {
	suite.Run(t, new(SnapshotTestSuite))
}

// Microsoft.Compute/snapshots/{snapshotName}
func (testsuite *SnapshotTestSuite) TestSnapshots() {
	var err error
	// From step Snapshots_CreateOrUpdate
	fmt.Println("Call operation: Snapshots_CreateOrUpdate")
	snapshotsClient, err := armcompute.NewSnapshotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	snapshotsClientCreateOrUpdateResponsePoller, err := snapshotsClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.snapshotName, armcompute.Snapshot{
		Location: to.Ptr(testsuite.location),
		Properties: &armcompute.SnapshotProperties{
			CreationData: &armcompute.CreationData{
				CreateOption: to.Ptr(armcompute.DiskCreateOptionEmpty),
			},
			DiskSizeGB: to.Ptr[int32](10),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, snapshotsClientCreateOrUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Snapshots_List
	fmt.Println("Call operation: Snapshots_List")
	snapshotsClientNewListPager := snapshotsClient.NewListPager(nil)
	for snapshotsClientNewListPager.More() {
		_, err := snapshotsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Snapshots_ListByResourceGroup
	fmt.Println("Call operation: Snapshots_ListByResourceGroup")
	snapshotsClientNewListByResourceGroupPager := snapshotsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, nil)
	for snapshotsClientNewListByResourceGroupPager.More() {
		_, err := snapshotsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step Snapshots_Get
	fmt.Println("Call operation: Snapshots_Get")
	_, err = snapshotsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.snapshotName, nil)
	testsuite.Require().NoError(err)

	// From step Snapshots_Update
	fmt.Println("Call operation: Snapshots_Update")
	snapshotsClientUpdateResponsePoller, err := snapshotsClient.BeginUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.snapshotName, armcompute.SnapshotUpdate{
		Properties: &armcompute.SnapshotUpdateProperties{
			DiskSizeGB: to.Ptr[int32](20),
		},
		Tags: map[string]*string{
			"department": to.Ptr("Development"),
			"project":    to.Ptr("UpdateSnapshots"),
		},
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, snapshotsClientUpdateResponsePoller)
	testsuite.Require().NoError(err)

	// From step Snapshots_GrantAccess
	fmt.Println("Call operation: Snapshots_GrantAccess")
	snapshotsClientGrantAccessResponsePoller, err := snapshotsClient.BeginGrantAccess(testsuite.ctx, testsuite.resourceGroupName, testsuite.snapshotName, armcompute.GrantAccessData{
		Access:            to.Ptr(armcompute.AccessLevelRead),
		DurationInSeconds: to.Ptr[int32](300),
	}, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, snapshotsClientGrantAccessResponsePoller)
	testsuite.Require().NoError(err)

	// From step Snapshots_RevokeAccess
	fmt.Println("Call operation: Snapshots_RevokeAccess")
	snapshotsClientRevokeAccessResponsePoller, err := snapshotsClient.BeginRevokeAccess(testsuite.ctx, testsuite.resourceGroupName, testsuite.snapshotName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, snapshotsClientRevokeAccessResponsePoller)
	testsuite.Require().NoError(err)

	// From step Snapshots_Delete
	fmt.Println("Call operation: Snapshots_Delete")
	snapshotsClientDeleteResponsePoller, err := snapshotsClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.snapshotName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, snapshotsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)
}
