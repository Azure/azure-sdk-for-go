//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/lease"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

// var headersToIgnoreForLease = []string {"X-Ms-Proposed-Lease-Id", "X-Ms-Lease-Id"}
var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func (s *LeaseRecordedTestsSuite) TestFilesystemAcquireLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	filesystemLeaseClient, _ := lease.NewFileSystemClient(filesystemClient, &lease.FileSystemClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := filesystemLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *filesystemLeaseClient.LeaseID())

	_, err = filesystemLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFilesystemDeleteFilesystemWithoutLeaseId() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	filesystemLeaseClient, _ := lease.NewFileSystemClient(filesystemClient, &lease.FileSystemClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := filesystemLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *filesystemLeaseClient.LeaseID())

	_, err = filesystemClient.Delete(ctx, nil)
	_require.Error(err)

	leaseID := filesystemLeaseClient.LeaseID()
	_, err = filesystemClient.Delete(ctx, &filesystem.DeleteOptions{
		AccessConditions: &filesystem.AccessConditions{
			LeaseAccessConditions: &filesystem.LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFilesystemReleaseLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	filesystemLeaseClient, _ := lease.NewFileSystemClient(filesystemClient, &lease.FileSystemClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := filesystemLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *filesystemLeaseClient.LeaseID())

	_, err = filesystemClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = filesystemLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)

	_, err = filesystemClient.Delete(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFilesystemRenewLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	filesystemLeaseClient, _ := lease.NewFileSystemClient(filesystemClient, &lease.FileSystemClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := filesystemLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *filesystemLeaseClient.LeaseID())

	_, err = filesystemLeaseClient.RenewLease(ctx, nil)
	_require.NoError(err)

	_, err = filesystemLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFilesystemChangeLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, fsName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	fsLeaseClient, _ := lease.NewFileSystemClient(fsClient, &lease.FileSystemClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := fsLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *fsLeaseClient.LeaseID())

	changeLeaseResp, err := fsLeaseClient.ChangeLease(ctx, *proposedLeaseIDs[1], nil)
	_require.NoError(err)
	_require.EqualValues(changeLeaseResp.LeaseID, proposedLeaseIDs[1])
	_require.EqualValues(fsLeaseClient.LeaseID(), proposedLeaseIDs[1])

	_, err = fsLeaseClient.RenewLease(ctx, nil)
	_require.NoError(err)

	_, err = fsLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFileAcquireLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := testcommon.CreateNewFile(context.Background(), _require, fileName, filesystemClient)
	fileLeaseClient, err := lease.NewPathClient(fileClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := fileLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestDeleteFileWithoutLeaseId() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := testcommon.CreateNewFile(context.Background(), _require, fileName, filesystemClient)
	fileLeaseClient, err := lease.NewPathClient(fileClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := fileLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	leaseID := fileLeaseClient.LeaseID()
	_, err = fileClient.Delete(ctx, &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			LeaseAccessConditions: &file.LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFileReleaseLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := testcommon.CreateNewFile(context.Background(), _require, fileName, filesystemClient)
	fileLeaseClient, _ := lease.NewPathClient(fileClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := fileLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = fileLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)

	_, err = fileClient.Delete(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFileRenewLease() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := testcommon.CreateNewFile(context.Background(), _require, fileName, filesystemClient)
	fileLeaseClient, _ := lease.NewPathClient(fileClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := fileLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, fileLeaseClient.LeaseID())

	_, err = fileLeaseClient.RenewLease(ctx, nil)
	_require.NoError(err)

	_, err = fileLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestFileChangeLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := testcommon.CreateNewFile(context.Background(), _require, fileName, filesystemClient)
	fileLeaseClient, _ := lease.NewPathClient(fileClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := fileLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.Equal(*acquireLeaseResponse.LeaseID, *proposedLeaseIDs[0])

	changeLeaseResp, err := fileLeaseClient.ChangeLease(ctx, *proposedLeaseIDs[1], nil)
	_require.NoError(err)
	_require.Equal(*changeLeaseResp.LeaseID, *proposedLeaseIDs[1])

	_, err = fileLeaseClient.RenewLease(ctx, nil)
	_require.NoError(err)

	_, err = fileLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestDirAcquireLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	dirName := testcommon.GenerateDirName(testName)
	dirClient := testcommon.CreateNewDir(context.Background(), _require, dirName, filesystemClient)
	dirLeaseClient, err := lease.NewPathClient(dirClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := dirLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, dirLeaseClient.LeaseID())

	_, err = dirLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestDeleteDirWithoutLeaseId() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	dirName := testcommon.GenerateDirName(testName)
	dirClient := testcommon.CreateNewDir(context.Background(), _require, dirName, filesystemClient)
	dirLeaseClient, err := lease.NewPathClient(dirClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := dirLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, dirLeaseClient.LeaseID())

	_, err = dirClient.Delete(ctx, nil)
	_require.Error(err)

	leaseID := dirLeaseClient.LeaseID()
	_, err = dirClient.Delete(ctx, &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			LeaseAccessConditions: &file.LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestDirReleaseLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	DirName := testcommon.GenerateDirName(testName)
	DirClient := testcommon.CreateNewDir(context.Background(), _require, DirName, filesystemClient)
	DirLeaseClient, _ := lease.NewPathClient(DirClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := DirLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, DirLeaseClient.LeaseID())

	_, err = DirClient.Delete(ctx, nil)
	_require.Error(err)

	_, err = DirLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)

	_, err = DirClient.Delete(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestDirRenewLease() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	dirName := testcommon.GenerateDirName(testName)
	dirClient := testcommon.CreateNewDir(context.Background(), _require, dirName, filesystemClient)
	dirLeaseClient, _ := lease.NewPathClient(dirClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := dirLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, dirLeaseClient.LeaseID())

	_, err = dirLeaseClient.RenewLease(ctx, nil)
	_require.NoError(err)

	_, err = dirLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}

func (s *LeaseRecordedTestsSuite) TestDirChangeLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	filesystemClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, filesystemClient)

	dirName := testcommon.GenerateDirName(testName)
	dirClient := testcommon.CreateNewDir(context.Background(), _require, dirName, filesystemClient)
	dirLeaseClient, _ := lease.NewPathClient(dirClient, &lease.PathClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := dirLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.Equal(*acquireLeaseResponse.LeaseID, *proposedLeaseIDs[0])

	changeLeaseResp, err := dirLeaseClient.ChangeLease(ctx, *proposedLeaseIDs[1], nil)
	_require.NoError(err)
	_require.Equal(*changeLeaseResp.LeaseID, *proposedLeaseIDs[1])

	_, err = dirLeaseClient.RenewLease(ctx, nil)
	_require.NoError(err)

	_, err = dirLeaseClient.ReleaseLease(ctx, nil)
	_require.NoError(err)
}
