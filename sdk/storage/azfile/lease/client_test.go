//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running lease Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &LeaseRecordedTestsSuite{})
		suite.Run(t, &LeaseUnrecordedTestsSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &LeaseRecordedTestsSuite{})
	case recording.RecordingMode:
		suite.Run(t, &LeaseRecordedTestsSuite{})
	}
}

func (s *LeaseRecordedTestsSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *LeaseRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
}

func (s *LeaseRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *LeaseRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *LeaseUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *LeaseUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type LeaseRecordedTestsSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type LeaseUnrecordedTestsSuite struct {
	suite.Suite
}

var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func (l *LeaseRecordedTestsSuite) TestShareAcquireLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	_, err = shareClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = shareLeaseClient.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestNegativeShareAcquireMultipleLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient0, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	shareLeaseClient1, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[1],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse0, err := shareLeaseClient0.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse0.LeaseID)
	_require.EqualValues(*acquireLeaseResponse0.LeaseID, *shareLeaseClient0.LeaseID())

	// acquiring lease for the second time returns LeaseAlreadyPresent error
	_, err = shareLeaseClient1.Acquire(ctx, int32(60), nil)
	_require.Error(err)

	_, err = shareLeaseClient0.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestShareDeleteShareWithoutLeaseId() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	_, err = shareClient.Delete(ctx, nil)
	_require.Error(err)

	leaseID := shareLeaseClient.LeaseID()
	_, err = shareClient.Delete(ctx, &share.DeleteOptions{
		LeaseAccessConditions: &share.LeaseAccessConditions{LeaseID: leaseID},
	})
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestShareReleaseLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	_, err = shareClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = shareLeaseClient.Release(ctx, nil)
	_require.NoError(err)

	_, err = shareClient.Delete(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestShareRenewLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(15), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	_, err = shareLeaseClient.Renew(ctx, nil)
	_require.NoError(err)

	_, err = shareLeaseClient.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestShareBreakLeaseDefault() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	bResp, err := shareLeaseClient.Break(ctx, nil)
	_require.NoError(err)
	_require.NotNil(bResp.LeaseTime)

	_, err = shareClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = shareLeaseClient.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestShareBreakLeaseNonDefault() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	bResp, err := shareLeaseClient.Break(ctx, &lease.ShareBreakOptions{
		BreakPeriod: to.Ptr((int32)(5)),
	})
	_require.NoError(err)
	_require.NotNil(bResp.LeaseTime)

	_, err = shareClient.Delete(ctx, nil)
	_require.Error(err)

	// wait for lease to expire
	time.Sleep(6 * time.Second)

	_, err = shareClient.Delete(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestNegativeShareBreakRenewLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	bResp, err := shareLeaseClient.Break(ctx, &lease.ShareBreakOptions{
		BreakPeriod: to.Ptr((int32)(5)),
	})
	_require.NoError(err)
	_require.NotNil(bResp.LeaseTime)

	// renewing broken lease returns error
	_, err = shareLeaseClient.Renew(ctx, nil)
	_require.Error(err)

	_, err = shareLeaseClient.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestShareChangeLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	shareLeaseClient, err := lease.NewShareClient(shareClient, &lease.ShareClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := shareLeaseClient.Acquire(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *shareLeaseClient.LeaseID())

	oldLeaseID := shareLeaseClient.LeaseID()

	changeLeaseResp, err := shareLeaseClient.Change(ctx, *proposedLeaseIDs[1], nil)
	_require.NoError(err)
	_require.EqualValues(changeLeaseResp.LeaseID, proposedLeaseIDs[1])
	_require.EqualValues(shareLeaseClient.LeaseID(), proposedLeaseIDs[1])

	_, err = shareClient.Delete(ctx, &share.DeleteOptions{
		LeaseAccessConditions: &share.LeaseAccessConditions{
			LeaseID: oldLeaseID,
		},
	})
	_require.Error(err)

	_, err = shareLeaseClient.Renew(ctx, nil)
	_require.NoError(err)

	_, err = shareLeaseClient.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestFileAcquireLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	ctx := context.Background()
	fileName := testcommon.GenerateFileName(testName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(ctx, 0, nil)
	_require.NoError(err)

	fileLeaseClient, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	acquireLeaseResponse, err := fileLeaseClient.Acquire(ctx, nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = fileLeaseClient.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestNegativeFileAcquireMultipleLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	ctx := context.Background()
	fileName := testcommon.GenerateFileName(testName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(ctx, 0, nil)
	_require.NoError(err)

	fileLeaseClient0, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	fileLeaseClient1, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[1],
	})
	_require.NoError(err)

	acquireLeaseResponse, err := fileLeaseClient0.Acquire(ctx, nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient0.LeaseID())

	// acquiring lease for the second time returns LeaseAlreadyPresent error
	_, err = fileLeaseClient1.Acquire(ctx, nil)
	_require.Error(err)

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = fileLeaseClient0.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestDeleteFileWithoutLeaseId() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	ctx := context.Background()
	fileName := testcommon.GenerateFileName(testName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(ctx, 0, nil)
	_require.NoError(err)

	fileLeaseClient, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	acquireLeaseResponse, err := fileLeaseClient.Acquire(ctx, nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	leaseID := fileLeaseClient.LeaseID()
	_, err = fileClient.Delete(ctx, &file.DeleteOptions{
		LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: leaseID,
		},
	})
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestFileReleaseLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	ctx := context.Background()
	fileName := testcommon.GenerateFileName(testName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(ctx, 0, nil)
	_require.NoError(err)

	fileLeaseClient, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	acquireLeaseResponse, err := fileLeaseClient.Acquire(ctx, nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = fileLeaseClient.Release(ctx, nil)
	_require.NoError(err)

	_, err = fileClient.Delete(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestFileChangeLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	ctx := context.Background()
	fileName := testcommon.GenerateFileName(testName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(ctx, 0, nil)
	_require.NoError(err)

	fileLeaseClient, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	acquireLeaseResponse, err := fileLeaseClient.Acquire(ctx, nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.Equal(*acquireLeaseResponse.LeaseID, *proposedLeaseIDs[0])

	oldLeaseID := fileLeaseClient.LeaseID()

	changeLeaseResp, err := fileLeaseClient.Change(ctx, *proposedLeaseIDs[1], nil)
	_require.NoError(err)
	_require.Equal(*changeLeaseResp.LeaseID, *proposedLeaseIDs[1])

	_, err = fileClient.Delete(ctx, &file.DeleteOptions{
		LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: oldLeaseID,
		},
	})
	_require.Error(err)

	_, err = fileLeaseClient.Release(ctx, nil)
	_require.NoError(err)
}

func (l *LeaseRecordedTestsSuite) TestNegativeFileDeleteAfterReleaseLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	ctx := context.Background()
	fileName := testcommon.GenerateFileName(testName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(ctx, 0, nil)
	_require.NoError(err)

	fileLeaseClient, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	acquireLeaseResponse, err := fileLeaseClient.Acquire(ctx, nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = fileLeaseClient.Release(ctx, nil)
	_require.NoError(err)

	// deleting file after its lease has expired or released returns error.
	_, err = fileClient.Delete(ctx, &file.DeleteOptions{
		LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: fileLeaseClient.LeaseID(),
		},
	})
	_require.Error(err)
}

func (l *LeaseRecordedTestsSuite) TestFileBreakLease() {
	_require := require.New(l.T())
	testName := l.T().Name()

	svcClient, err := testcommon.GetServiceClient(l.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	ctx := context.Background()
	fileName := testcommon.GenerateFileName(testName)
	fileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)
	_, err = fileClient.Create(ctx, 0, nil)
	_require.NoError(err)

	fileLeaseClient, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	acquireLeaseResponse, err := fileLeaseClient.Acquire(ctx, nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = fileLeaseClient.Break(ctx, nil)
	_require.NoError(err)

	_, err = fileClient.Delete(ctx, nil)
	_require.NoError(err)
}
