//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running datalake Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &RecordedTestSuite{})
		suite.Run(t, &UnrecordedTestSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &RecordedTestSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &RecordedTestSuite{})
	}
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
}

type UnrecordedTestSuite struct {
	suite.Suite
}

func (s *RecordedTestSuite) TestCreateDirAndDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	createDirOpts := &directory.CreateOptions{
		AccessConditions: nil,
	}

	resp, err = dirClient.Create(context.Background(), createDirOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	DirName := testcommon.GenerateDirName(testName)
	DirClient, err := testcommon.GetDirClient(filesystemName, DirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, DirClient)

	resp, err := DirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	DirName := testcommon.GenerateDirName(testName)
	DirClient, err := testcommon.GetDirClient(filesystemName, DirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, DirClient)

	resp, err := DirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateDirIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateFileIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateFileWithMetadataNotNil() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createDirOpts := &directory.CreateOptions{
		Metadata: testcommon.BasicMetadata,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), createDirOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithEmptyMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createDirOpts := &directory.CreateOptions{
		Metadata: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), createDirOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithNilHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createDirOpts := &directory.CreateOptions{
		HTTPHeaders: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), createDirOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createDirOpts := &directory.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), createDirOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirWithLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createFileOpts := &directory.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	// should fail since leased
	_, err = dirClient.Create(context.Background(), createFileOpts)
	_require.NotNil(err)

	time.Sleep(time.Second * 15)
	resp, err = dirClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateDirWithPermissions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	perms := "0777"
	umask := "0000"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createFileOpts := &directory.CreateOptions{
		Permissions: &perms,
		Umask:       &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	resp2, err := dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp2)
	_require.Equal("rwxrwxrwx", *resp2.Permissions)
}

func (s *RecordedTestSuite) TestCreateDirWithOwnerGroupACLUmask() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	umask := "0000"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createFileOpts := &directory.CreateOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		Umask: &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

}

func (s *RecordedTestSuite) TestDeleteDirWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.Nil(err)

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: nil,
	}

	resp, err := dirClient.Delete(context.Background(), deleteOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteDirIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteOpts := &directory.DeleteOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp1, err := dirClient.Delete(context.Background(), deleteOpts)
	_require.Nil(err)
	_require.NotNil(resp1)
}

func (s *RecordedTestSuite) TestDeleteDirIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteDirIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteDirIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteDirIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteDirIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestSetAccessControlNil() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = dirClient.SetAccessControl(context.Background(), nil)
	_require.NotNil(err)

	_require.Equal(err, datalakeerror.MissingParameters)
}

// TODO: write test that fails if you provide permissions and acls
func (s *RecordedTestSuite) TestSetAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	opts := &directory.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetAccessControlWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	opts := &directory.SetAccessControlOptions{
		Owner:            &owner,
		Group:            &group,
		ACL:              &acl,
		AccessConditions: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = dirClient.SetAccessControl(context.Background(), opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetAccessControlIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetAccessControlIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestSetAccessControlIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetAccessControlIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestSetAccessControlIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetAccessControlIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestGetAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	getACLResp, err := dirClient.GetAccessControl(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestGetAccessControlWithSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
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
	_require.Nil(err)

	dirClient2, _ := directory.NewClientWithNoCredential(sasURL, nil)

	getACLResp, err := dirClient2.GetAccessControl(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestDeleteWithSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)

	dirClient2, _ := directory.NewClientWithNoCredential(sasURL, nil)

	_, err = dirClient2.Delete(context.Background(), nil)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestGetAccessControlWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &directory.GetAccessControlOptions{
		AccessConditions: nil,
	}

	getACLResp, err := dirClient.GetAccessControl(context.Background(), opts)
	_require.Nil(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestGetAccessControlIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
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
	_require.Nil(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestGetAccessControlIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestGetAccessControlIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)
	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		}}

	getACLResp, err := dirClient.GetAccessControl(context.Background(), opts)
	_require.Nil(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestGetAccessControlIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
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
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestGetAccessControlIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
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
	_require.Nil(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestGetAccessControlIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag
	opts := &directory.GetAccessControlOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		}}

	_, err = dirClient.GetAccessControl(context.Background(), opts)
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

///=====================================================================

func (s *RecordedTestSuite) TestUpdateAccessControlRecursive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	acl1 := "user::rwx,group::r--,other::r--"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &directory.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	pager := dirClient.NewUpdateAccessControlRecursivePager(acl1, nil)

	for pager.More() {
		resp1, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Equal(*resp1.SetAccessControlRecursiveResponse.DirectoriesSuccessful, int32(1))
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestRemoveAccessControlRecursive() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "mask," + "default:user,default:group," +
		"user:ec3595d6-2c17-4696-8caa-7e139758d24a,group:ec3595d6-2c17-4696-8caa-7e139758d24a," +
		"default:user:ec3595d6-2c17-4696-8caa-7e139758d24a,default:group:ec3595d6-2c17-4696-8caa-7e139758d24a"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	pager := dirClient.NewRemoveAccessControlRecursivePager(acl, nil)

	for pager.More() {
		resp1, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Equal(*resp1.SetAccessControlRecursiveResponse.DirectoriesSuccessful, int32(1))
		if err != nil {
			break
		}
	}
}

func (s *RecordedTestSuite) TestSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = dirClient.SetMetadata(context.Background(), nil)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetMetadataWithEmptyOpts() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &directory.SetMetadataOptions{
		Metadata: nil,
	}
	_, err = dirClient.SetMetadata(context.Background(), opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetMetadataWithBasicMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &directory.SetMetadataOptions{
		Metadata: testcommon.BasicMetadata,
	}
	_, err = dirClient.SetMetadata(context.Background(), opts)
	_require.Nil(err)
}

func (s *RecordedTestSuite) TestSetMetadataWithAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteDir(context.Background(), _require, dirClient)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &directory.SetMetadataOptions{
		Metadata: testcommon.BasicMetadata,
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = dirClient.SetMetadata(context.Background(), opts)
	_require.Nil(err)
}

func validatePropertiesSet(_require *require.Assertions, dirClient *directory.Client, disposition string) {
	resp, err := dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentDisposition, disposition)
}

func (s *RecordedTestSuite) TestSetHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, nil)
	_require.Nil(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestSetHTTPHeadersWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: nil,
	}

	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Nil(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestSetHTTPHeadersIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestSetHTTPHeadersIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestSetHTTPHeadersIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.Nil(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestSetHTTPHeadersIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestSetHTTPHeadersIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	opts := &directory.SetHTTPHeadersOptions{
		AccessConditions: &directory.AccessConditions{
			ModifiedAccessConditions: &directory.ModifiedAccessConditions{
				IfMatch: etag,
			},
		}}
	_, err = dirClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Nil(err)
	validatePropertiesSet(_require, dirClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestSetHTTPHeadersIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestRenameNoOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	resp1, err := dirClient.Rename(context.Background(), "newName", nil)
	_require.Nil(err)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewDirectoryClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	renameFileOpts := &directory.RenameOptions{
		AccessConditions: nil,
	}

	resp1, err := dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.Nil(err)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewDirectoryClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfModifiedSince: &currentTime,
			},
		},
	}
	resp1, err := dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.Nil(err)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewDirectoryClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestRenameFileIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfUnmodifiedSince: &currentTime,
			},
		},
	}

	resp1, err := dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewDirectoryClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestRenameFileIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	renameFileOpts := &directory.RenameOptions{
		SourceAccessConditions: &directory.SourceAccessConditions{
			SourceModifiedAccessConditions: &directory.SourceModifiedAccessConditions{
				SourceIfMatch: etag,
			},
		},
	}

	resp1, err := dirClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewDirectoryClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	dirName := testcommon.GenerateDirName(testName)
	dirClient, err := testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
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
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}
