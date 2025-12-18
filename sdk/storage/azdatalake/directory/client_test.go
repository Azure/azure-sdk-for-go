//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running datalake Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &RecordedTestSuite{})
		suite.Run(t, &UnrecordedTestSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &RecordedTestSuite{})
	case recording.RecordingMode:
		suite.Run(t, &RecordedTestSuite{})
	}
}

func (s *RecordedTestSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *RecordedTestSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
}

func (s *RecordedTestSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *RecordedTestSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *UnrecordedTestSuite) BeforeTest(suite string, test string) {

}

func (s *UnrecordedTestSuite) AfterTest(suite string, test string) {

}

type RecordedTestSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type UnrecordedTestSuite struct {
	suite.Suite
}

func (s *UnrecordedTestSuite) TestCreateDirAndDeleteWithConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()

	connectionString, _ := testcommon.GetGenericConnectionString(testcommon.TestAccountDatalake)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := directory.NewClientFromConnectionString(*connectionString, dirName, filesystemName, nil)

	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestBlobURLAndDFSURL() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_require.Contains(dirClient.DFSURL(), ".dfs.core.windows.net/"+filesystemName+"/"+dirName)
	_require.Contains(dirClient.BlobURL(), ".blob.core.windows.net/"+filesystemName+"/"+dirName)
}

func (s *RecordedTestSuite) TestCreateDirAndDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirUsingCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	dirOpts := &directory.CreateOptions{CPKInfo: &testcommon.TestCPKByValue}
	resp, err := dirClient.Create(context.Background(), dirOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	_require.Equal(true, *(resp.IsServerEncrypted))
	// run the below check if the test is not in playback mode
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.Equal(testcommon.TestCPKByValue.EncryptionKeySHA256, resp.EncryptionKeySHA256)
	}
}

func (s *RecordedTestSuite) TestGetAndCreateFileClient() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	fileClient, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName))
	_require.NoError(err)
	_require.NotNil(fileClient)

	_, err = fileClient.Create(context.Background(), nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestGetAndCreateFileClientWithSpecialNames() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirClient := fsClient.NewDirectoryClient("#,%,?")
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	fileClient, err := dirClient.NewFileClient("?%#,")
	_require.NoError(err)
	_require.NotNil(fileClient)

	_, err = fileClient.Create(context.Background(), nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestCreateNewSubdirectoryClient() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := svcClient.NewFileSystemClient(fsName)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	dirName := testcommon.GenerateDirName(testName)
	dirClient := fsClient.NewDirectoryClient(dirName)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	subdirName := testcommon.GenerateSubDirName(testName)
	subdirClient, err := dirClient.NewSubdirectoryClient(subdirName)
	_require.NoError(err)

	perm := "r-xr-x---"

	_, err = subdirClient.Create(context.Background(), &directory.CreateOptions{
		Permissions: &perm,
		CPKInfo:     &testcommon.TestCPKByValue,
	})
	_require.NoError(err)

	resp, err := subdirClient.GetProperties(context.Background(), &directory.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.NotNil(resp.Permissions)
	_require.Equal(perm, *resp.Permissions)
	_require.Equal(*(resp.IsServerEncrypted), true)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.Equal(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}
	// Create a file under the new directory just to make sure we're not secretly targeting the parent
	fileName := testcommon.GenerateFileName("newFile")
	subdirFileClient, err := subdirClient.NewFileClient(fileName)
	_require.NoError(err)

	_, err = subdirFileClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileContents := []byte(testcommon.DefaultData)
	err = subdirFileClient.UploadBuffer(context.Background(), fileContents, nil)
	_require.NoError(err)

	// check for the file at the parent directory
	dirFileClient, err := dirClient.NewFileClient(fileName)
	_require.NoError(err)

	_, err = dirFileClient.GetProperties(context.Background(), nil)
	_require.Error(err) // we should get back a 404
	_require.True(datalakeerror.HasCode(err, datalakeerror.PathNotFound))
}

func (s *RecordedTestSuite) TestCreateNewSubdirectoryClientWithSpecialName() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := svcClient.NewFileSystemClient(fsName)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	dirClient := fsClient.NewDirectoryClient("?#%")

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	subdirClient, err := dirClient.NewSubdirectoryClient(",%,#,?")
	_require.NoError(err)

	perm := "r-xr-x---"

	_, err = subdirClient.Create(context.Background(), &directory.CreateOptions{
		Permissions: &perm,
		CPKInfo:     &testcommon.TestCPKByValue,
	})
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestCreateDirWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	createDirOpts := &directory.CreateOptions{
		AccessConditions: nil,
	}

	resp, err = dirClient.Create(context.Background(), createDirOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	DirName := testcommon.GenerateDirName(testName)
	DirClient, err := testcommon.GetDirClient(filesystemName, DirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, DirClient)

	resp, err := DirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	createDirOpts := &directory.CreateOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp, err = DirClient.Create(context.Background(), createDirOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	DirName := testcommon.GenerateDirName(testName)
	DirClient, err := testcommon.GetDirClient(filesystemName, DirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, DirClient)

	resp, err := DirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	createDirOpts := &directory.CreateOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp, err = DirClient.Create(context.Background(), createDirOpts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateDirIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	createDirOpts := &directory.CreateOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	resp, err = dirClient.Create(context.Background(), createDirOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	createDirOpts := &directory.CreateOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	resp, err = dirClient.Create(context.Background(), createDirOpts)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateDirIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	createDirOpts := &directory.CreateOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	resp, err = dirClient.Create(context.Background(), createDirOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	createDirOpts := &directory.CreateOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		},
	}

	resp, err = dirClient.Create(context.Background(), createDirOpts)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateDirWithNilHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createDirOpts := &directory.CreateOptions{
		HTTPHeaders: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), createDirOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirWithHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createDirOpts := &directory.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), createDirOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirWithLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &directory.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	// should fail since leased
	_, err = dirClient.Create(context.Background(), createFileOpts)
	_require.Error(err)

	time.Sleep(time.Second * 15)
	resp, err = dirClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirWithPermissions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	perms := "0777"
	umask := "0000"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &directory.CreateOptions{
		Permissions: &perms,
		Umask:       &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	resp2, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.Equal("rwxrwxrwx", *resp2.Permissions)
}

func (s *RecordedTestSuite) TestCreateDirWithOwnerGroupACLUmask() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	umask := "0000"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &directory.CreateOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		Umask: &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

}

func (s *RecordedTestSuite) TestDeleteDirWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: nil,
	}

	resp, err := dirClient.Delete(context.Background(), deleteOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

// To run this test, the NamespaceTenant AAD info needs to be set to an AAD app that does not have any RBAC permissions,
// and entityId needs to be set to the entity ID of the application.
func (s *RecordedTestSuite) TestDeleteDirWithPaginatedDelete() {

	s.T().Skip("AAD app not configured for this test, this will be skipped")
	_require := require.New(s.T())
	testName := s.T().Name()
	user := "user"
	readWriteExecutePermission := "rwx"

	objectId := "" // object ID of an AAD app which has no RBAC permissions
	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	rootDirectory := fsClient.NewDirectoryClient("/")

	dirName := testcommon.GenerateDirName(testName)
	dirURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName + "/" + dirName
	credential, err := azdatalake.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	dirClient, err := directory.NewClientWithSharedKeyCredential(dirURL, credential, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	for i := 0; i < 5020; i++ {
		fileClient, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName) + strconv.Itoa(i))
		_require.NoError(err)
		_require.NotNil(fileClient)

		_, err = fileClient.Create(context.Background(), nil)
		_require.NoError(err)
	}

	accessControlResp, err := rootDirectory.GetAccessControl(context.Background(), nil)
	_require.NoError(err)

	newAcl := *accessControlResp.ACL + "," + user + ":" + objectId + ":" + readWriteExecutePermission

	_, err = rootDirectory.SetAccessControlRecursive(context.Background(), newAcl, nil)
	_require.NoError(err)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	directoryURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName + "/" + dirName

	newDirClient, err := directory.NewClient(directoryURL, cred, nil)
	_require.NoError(err)

	deleteOpts := &directory.DeleteOptions{
		Paginated: to.Ptr(true),
	}

	response, err := newDirClient.Delete(context.Background(), deleteOpts)
	_require.NoError(err)
	_require.Nil(response.Continuation)
	_require.NotNil(response)
}

func (s *RecordedTestSuite) TestDeleteDirIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp1, err := dirClient.Delete(context.Background(), deleteOpts)
	_require.NoError(err)
	_require.NotNil(resp1)
}

func (s *RecordedTestSuite) TestDeleteDirIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.Delete(context.Background(), deleteOpts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteDirIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.Delete(context.Background(), deleteOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteDirIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.Delete(context.Background(), deleteOpts)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteDirIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	_, err = dirClient.Delete(context.Background(), deleteOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteDirIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		},
	}

	_, err = dirClient.Delete(context.Background(), deleteOpts)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirSetAccessControlNil() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = dirClient.SetAccessControl(context.Background(), nil)
	_require.Error(err)

	_require.Equal(err, datalakeerror.MissingParameters)
}

// TODO: write test that fails if you provide permissions and acls
func (s *RecordedTestSuite) TestDirSetAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetAccessControlWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	opts := &directory.SetAccessControlOptions{
		Owner:            &owner,
		Group:            &group,
		ACL:              &acl,
		AccessConditions: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetAccessControlIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)
	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetAccessControlIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirSetAccessControlIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)
	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		}}

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetAccessControlIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirSetAccessControlIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
	etag := resp.ETag

	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetAccessControlIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag
	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		}}

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirGetAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	getACLResp, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *UnrecordedTestSuite) TestDirGetAccessControlWithSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	// Adding SAS and options
	permissions := sas.DirectoryPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	sasURL, err := dirClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	dirClient2, _ := directory.NewClientWithNoCredential(sasURL, nil)

	getACLResp, err := dirClient2.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *UnrecordedTestSuite) TestDeleteWithSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// Adding SAS and options
	permissions := sas.DirectoryPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	sasURL, err := dirClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	dirClient2, _ := directory.NewClientWithNoCredential(sasURL, nil)

	_, err = dirClient2.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirGetAccessControlWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &directory.GetAccessControlOptions{
		AccessConditions: nil,
	}

	getACLResp, err := dirClient.GetAccessControl(context.Background(), opts)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirGetAccessControlIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)
	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	getACLResp, err := dirClient.GetAccessControl(context.Background(), opts)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirGetAccessControlIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)
	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.GetAccessControl(context.Background(), opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirGetAccessControlIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)
	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		}}

	getACLResp, err := dirClient.GetAccessControl(context.Background(), opts)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirGetAccessControlIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.GetAccessControl(context.Background(), opts)
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirGetAccessControlIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)
	etag := resp.ETag

	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	getACLResp, err := dirClient.GetAccessControl(context.Background(), opts)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirGetAccessControlIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag
	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		}}

	_, err = dirClient.GetAccessControl(context.Background(), opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

///=====================================================================

func (s *RecordedTestSuite) TestDirSetAccessControlRecursive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := dirClient.SetAccessControlRecursive(context.Background(), acl, nil)
	_require.NoError(err)

	_require.Equal(resp1.DirectoriesSuccessful, to.Ptr(int32(1)))
	_require.Equal(resp1.FilesSuccessful, to.Ptr(int32(0)))
	_require.Equal(resp1.FailureCount, to.Ptr(int32(0)))

	getACLResp, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirSetAccessControlRecursiveWithBadContinuation() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &directory.SetAccessControlRecursiveOptions{Marker: to.Ptr("garbage")}
	resp1, err := dirClient.SetAccessControlRecursive(context.Background(), acl, opts)
	_require.NoError(err)

	_require.Equal(resp1.DirectoriesSuccessful, to.Ptr(int32(0)))
	_require.Equal(resp1.FilesSuccessful, to.Ptr(int32(0)))
	_require.Equal(resp1.FailureCount, to.Ptr(int32(0)))
}

func (s *RecordedTestSuite) TestDirSetAccessControlRecursiveWithEmptyOpts() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &directory.SetAccessControlRecursiveOptions{}
	resp1, err := dirClient.SetAccessControlRecursive(context.Background(), acl, opts)
	_require.NoError(err)

	_require.Equal(resp1.DirectoriesSuccessful, to.Ptr(int32(1)))
	_require.Equal(resp1.FilesSuccessful, to.Ptr(int32(0)))
	_require.Equal(resp1.FailureCount, to.Ptr(int32(0)))

	getACLResp, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirSetAccessControlRecursiveWithMaxResults() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp1)

	fileClient, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName))
	_require.NoError(err)
	_require.NotNil(fileClient)

	_, err = fileClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileClient1, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName + "1"))
	_require.NoError(err)
	_require.NotNil(fileClient1)

	_, err = fileClient1.Create(context.Background(), nil)
	_require.NoError(err)

	opts := &directory.SetAccessControlRecursiveOptions{BatchSize: to.Ptr(int32(2)), MaxBatches: to.Ptr(int32(1)), ContinueOnFailure: to.Ptr(true), Marker: nil}
	resp2, err := dirClient.SetAccessControlRecursive(context.Background(), acl, opts)
	_require.NoError(err)

	// we expect only one file to have been updated not both since our batch size is 2 and max batches is 1
	_require.Equal(resp2.DirectoriesSuccessful, to.Ptr(int32(1)))
	_require.Equal(resp2.FilesSuccessful, to.Ptr(int32(1)))
	_require.Equal(resp2.FailureCount, to.Ptr(int32(0)))

	getACLResp, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirSetAccessControlRecursiveWithMaxResults2() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp1)

	fileClient, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName))
	_require.NoError(err)
	_require.NotNil(fileClient)

	_, err = fileClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileClient1, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName + "1"))
	_require.NoError(err)
	_require.NotNil(fileClient1)

	_, err = fileClient1.Create(context.Background(), nil)
	_require.NoError(err)

	opts := &directory.SetAccessControlRecursiveOptions{ContinueOnFailure: to.Ptr(true), Marker: nil}
	resp2, err := dirClient.SetAccessControlRecursive(context.Background(), acl, opts)
	_require.NoError(err)

	// we expect only one file to have been updated not both since our batch size is 2 and max batches is 1
	_require.Equal(resp2.DirectoriesSuccessful, to.Ptr(int32(1)))
	_require.Equal(resp2.FilesSuccessful, to.Ptr(int32(2)))
	_require.Equal(resp2.FailureCount, to.Ptr(int32(0)))

	getACLResp, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirSetAccessControlRecursiveWithMaxResults3() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp1)

	fileClient, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName))
	_require.NoError(err)
	_require.NotNil(fileClient)

	_, err = fileClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileClient1, err := dirClient.NewFileClient(testcommon.GenerateFileName(testName + "1"))
	_require.NoError(err)
	_require.NotNil(fileClient1)

	_, err = fileClient1.Create(context.Background(), nil)
	_require.NoError(err)

	opts := &directory.SetAccessControlRecursiveOptions{BatchSize: to.Ptr(int32(1)), ContinueOnFailure: to.Ptr(true), Marker: nil}
	resp2, err := dirClient.SetAccessControlRecursive(context.Background(), acl, opts)
	_require.NoError(err)

	// we expect only one file to have been updated not both since our batch size is 2 and max batches is 1
	_require.Equal(resp2.DirectoriesSuccessful, to.Ptr(int32(1)))
	_require.Equal(resp2.FilesSuccessful, to.Ptr(int32(2)))
	_require.Equal(resp2.FailureCount, to.Ptr(int32(0)))

	getACLResp, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDirUpdateAccessControlRecursive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	acl1 := "user::rwx,group::r--,other::r--"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := dirClient.UpdateAccessControlRecursive(context.Background(), acl1, nil)
	_require.NoError(err)
	_require.Equal(resp1.DirectoriesSuccessful, to.Ptr(int32(1)))

	resp2, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl1, *resp2.ACL)

}

func (s *RecordedTestSuite) TestDirRemoveAccessControlRecursive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "mask," + "default:user,default:group," +
		"user:ec3595d6-2c17-4696-8caa-7e139758d24a,group:ec3595d6-2c17-4696-8caa-7e139758d24a," +
		"default:user:ec3595d6-2c17-4696-8caa-7e139758d24a,default:group:ec3595d6-2c17-4696-8caa-7e139758d24a"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := dirClient.RemoveAccessControlRecursive(context.Background(), acl, nil)
	_require.NoError(err)
	_require.Equal(resp1.DirectoriesSuccessful, to.Ptr(int32(1)))
}

func (s *RecordedTestSuite) TestDirSetMetadataWithBasicMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = dirClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetMetadataWithAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &directory.SetMetadataOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = dirClient.SetMetadata(context.Background(), testcommon.BasicMetadata, opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetMetadataWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), &directory.CreateOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &directory.SetMetadataOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	_, err = dirClient.SetMetadata(context.Background(), testcommon.BasicMetadata, opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirSetMetadataWithCPKNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &directory.SetMetadataOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	_, err = dirClient.SetMetadata(context.Background(), testcommon.BasicMetadata, opts)
	_require.Error(err)
	_require.ErrorContains(err, "BlobDoesNotUseCustomerSpecifiedEncryption")

}

func validatePropertiesSet(_require *require.Assertions, dirClient *directory.Client, disposition string) {
	resp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.ContentDisposition, disposition)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, nil)
	_require.NoError(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeadersWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: nil,
	}

	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.NoError(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeadersIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.NoError(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeadersIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeadersIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.NoError(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeadersIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeadersIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfMatch: etag,
			},
		}}
	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.NoError(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestDirSetHTTPHeadersIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		},
	}
	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *UnrecordedTestSuite) TestDirectoryRenameUsingSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDatalake)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	perms := sas.DirectoryPermissions{Read: true, Create: true, Write: true, Move: true, Delete: true, List: true}
	sasQueryParams, err := sas.DatalakeSignatureValues{
		Protocol:       sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:     time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		FileSystemName: filesystemName,
		Permissions:    perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	sasToken := sasQueryParams.Encode()

	srcDirClient, err := directory.NewClientWithNoCredential(fsClient.DFSURL()+"/dir1?"+sasToken, nil)
	_require.NoError(err)

	_, err = srcDirClient.Create(context.Background(), nil)
	_require.NoError(err)

	destPathWithSAS := "dir2?" + sasToken
	_, err = srcDirClient.Rename(context.Background(), destPathWithSAS, nil)
	_require.NoError(err)

	_, err = srcDirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)
}

func (s *RecordedTestSuite) TestDirRenameNoOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = dirClient.Rename(context.Background(), "newName", nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestDirRenameRequestWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	createOpts := &directory.CreateOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	renameFileOpts := &directory.RenameOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestRenameDirWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	renameFileOpts := &directory.RenameOptions{
		AccessConditions: nil,
	}

	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestRenameDirIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfModifiedSince: &currentTime,
			},
		},
	}
	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestRenameDirIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfModifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestRenameDirIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestRenameDirIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestRenameDirIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfMatch: etag,
			},
		},
	}

	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestRenameDirIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfNoneMatch: etag,
			},
		},
	}

	_, err = dirClient.Rename(context.Background(), "newName", renameFileOpts)

	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestDirGetPropertiesResponseCapture() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// This tests directory.NewClient
	var respFromCtxDir *http.Response
	ctxWithRespDir := runtime.WithCaptureResponse(context.Background(), &respFromCtxDir)
	resp2, err := dirClient.GetProperties(ctxWithRespDir, nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.NotNil(respFromCtxDir) // validate that the respFromCtx is actually populated
	_require.Equal("directory", *resp2.ResourceType)

	// This tests filesystem.NewClient
	dirClient = fsClient.NewDirectoryClient(dirName)
	var respFromCtxFs *http.Response
	ctxWithRespFs := runtime.WithCaptureResponse(context.Background(), &respFromCtxFs)
	resp2, err = dirClient.GetProperties(ctxWithRespFs, nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.NotNil(respFromCtxFs) // validate that the respFromCtx is actually populated
	_require.Equal("directory", *resp2.ResourceType)

	// This tests service.NewClient
	serviceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	fsClient = serviceClient.NewFileSystemClient(filesystemName)
	dirClient = fsClient.NewDirectoryClient(dirName)
	var respFromCtxService *http.Response
	ctxWithRespService := runtime.WithCaptureResponse(context.Background(), &respFromCtxService)
	resp2, err = dirClient.GetProperties(ctxWithRespService, nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.NotNil(respFromCtxService) // validate that the respFromCtx is actually populated
	_require.Equal("directory", *resp2.ResourceType)
}

func (s *RecordedTestSuite) TestDirGetPropertiesWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	dirOpts := &directory.CreateOptions{CPKInfo: &testcommon.TestCPKByValue}
	resp, err := dirClient.Create(context.Background(), dirOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	getPropertiesOpts := &directory.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue}
	response, err := dirClient.GetProperties(context.Background(), getPropertiesOpts)
	_require.NoError(err)
	_require.NotNil(response)
	_require.Equal(*(resp.IsServerEncrypted), true)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.Equal(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}
}

func (s *UnrecordedTestSuite) TestDirCreateDeleteUsingOAuth() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	dirName := testcommon.GenerateDirName(testName)
	dirURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName + "/" + dirName

	dirClient, err := directory.NewClient(dirURL, cred, nil)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestCreateDirectoryClientDefaultAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	dirName := testcommon.GenerateDirName(testName)
	dirURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName + "/" + dirName

	options := &directory.ClientOptions{Audience: "https://storage.azure.com/"}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

}

func (s *RecordedTestSuite) TestCreateDirectoryClientCustomAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	dirName := testcommon.GenerateDirName(testName)
	dirURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName + "/" + dirName

	options := &directory.ClientOptions{Audience: "https://" + accountName + ".blob.core.windows.net"}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)

	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

}

func (s *UnrecordedTestSuite) TestDirectoryClientOnAuthenticationFailure() {
	_require := require.New(s.T())
	testName := s.T().Name()

	tenantID := "xxx"
	clientID := "xxx"
	clientSecret := "xxx"

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	_require.NoError(err)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	url := "https://" + accountName + ".dfs.core.windows.net/"

	srvClient, err := service.NewClient(url, cred, nil)
	_require.NoError(err)

	fsName := testcommon.GenerateFileSystemName(testName)
	fsClient := srvClient.NewFileSystemClient(fsName)
	dirClient := fsClient.NewDirectoryClient("testdir")

	// This should return an error, not panic
	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err, "Expected authentication error")
	_require.Contains(err.Error(), "ClientSecretCredential")
}
