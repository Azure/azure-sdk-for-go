// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armfileshares_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/fileshares/armfileshares"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/suite"
)

type FileSharesTestSuite struct {
	suite.Suite

	ctx                context.Context
	cred               azcore.TokenCredential
	options            *arm.ClientOptions
	fileShareName      string
	storageAccountName string
	location           string
	resourceGroupName  string
	subscriptionId     string
}

func (testsuite *FileSharesTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.fileShareName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "fileshare", 15, true)
	testsuite.storageAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "storageacc", 16, true)
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
}

func (testsuite *FileSharesTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestFileSharesTestSuite(t *testing.T) {
	suite.Run(t, new(FileSharesTestSuite))
}

// Test creating a file share snapshot
func (testsuite *FileSharesTestSuite) TestFileShareCreate() {
	var err error
	snapshotName := "snapshot1"

	// First, create the file share
	fmt.Println("Call operation: FileShares_CreateOrUpdate")
	fileShareClient, err := armfileshares.NewClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	fileShareCreatePoller, err := fileShareClient.BeginCreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.fileShareName, armfileshares.FileShare{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"environment": to.Ptr("test"),
		},
		Properties: &armfileshares.FileShareProperties{
			MountName:                      to.Ptr("fileshare"),
			MediaTier:                      to.Ptr(armfileshares.MediaTierSSD),
			Redundancy:                     to.Ptr(armfileshares.RedundancyLocal),
			Protocol:                       to.Ptr(armfileshares.ProtocolNFS),
			ProvisionedStorageGiB:          to.Ptr[int32](32),
			ProvisionedIOPerSec:            to.Ptr[int32](3000),
			ProvisionedThroughputMiBPerSec: to.Ptr[int32](125),
			NfsProtocolProperties: &armfileshares.NfsProtocolProperties{
				RootSquash: to.Ptr(armfileshares.ShareRootSquashNoRootSquash),
			},
			PublicNetworkAccess: to.Ptr(armfileshares.PublicNetworkAccessEnabled),
		},
	}, nil)
	testsuite.Require().NoError(err)

	_, err = testutil.PollForTest(testsuite.ctx, fileShareCreatePoller)
	testsuite.Require().NoError(err)

	// From step FileShareSnapshots_CreateOrUpdate
	fmt.Println("Call operation: FileShareSnapshots_CreateOrUpdateFileShareSnapshot")
	fileShareSnapshotsClient, err := armfileshares.NewFileShareSnapshotsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)

	fileShareSnapshotsClientCreateResponsePoller, err := fileShareSnapshotsClient.BeginCreateOrUpdateFileShareSnapshot(testsuite.ctx, testsuite.resourceGroupName, testsuite.fileShareName, snapshotName, armfileshares.FileShareSnapshot{
		Properties: &armfileshares.FileShareSnapshotProperties{
			Metadata: map[string]*string{
				"test": to.Ptr("snapshot"),
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	_, err = testutil.PollForTest(testsuite.ctx, fileShareSnapshotsClientCreateResponsePoller)
	testsuite.Require().NoError(err)

	// From step FileShareSnapshots_Get
	fmt.Println("Call operation: FileShareSnapshots_GetFileShareSnapshot")
	_, err = fileShareSnapshotsClient.GetFileShareSnapshot(testsuite.ctx, testsuite.resourceGroupName, testsuite.fileShareName, snapshotName, nil)
	testsuite.Require().NoError(err)

	// From step FileShareSnapshots_List
	fmt.Println("Call operation: FileShareSnapshots_ListByFileShare")
	fileShareSnapshotsClientNewListPager := fileShareSnapshotsClient.NewListByFileSharePager(testsuite.resourceGroupName, testsuite.fileShareName, nil)
	for fileShareSnapshotsClientNewListPager.More() {
		_, err := fileShareSnapshotsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step FileShareSnapshots_Delete
	fmt.Println("Call operation: FileShareSnapshots_DeleteFileShareSnapshot")
	fileShareSnapshotsClientDeleteResponsePoller, err := fileShareSnapshotsClient.BeginDeleteFileShareSnapshot(testsuite.ctx, testsuite.resourceGroupName, testsuite.fileShareName, snapshotName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fileShareSnapshotsClientDeleteResponsePoller)
	testsuite.Require().NoError(err)

	// Clean up - delete the file share
	fmt.Println("Call operation: FileShares_Delete")
	fileShareDeletePoller, err := fileShareClient.BeginDelete(testsuite.ctx, testsuite.resourceGroupName, testsuite.fileShareName, nil)
	testsuite.Require().NoError(err)
	_, err = testutil.PollForTest(testsuite.ctx, fileShareDeletePoller)
	testsuite.Require().NoError(err)
}
