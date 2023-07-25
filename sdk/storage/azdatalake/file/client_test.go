//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"hash/crc64"
	"io"
	"math/rand"
	"net/http"
	"os"
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

func validateFileDeleted(_require *require.Assertions, fileClient *file.Client) {
	_, err := fileClient.GetAccessControl(context.Background(), nil)
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)
}

func (s *RecordedTestSuite) TestCreateFileAndDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	createFileOpts := &file.CreateOptions{
		AccessConditions: nil,
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	createFileOpts := &file.CreateOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	createFileOpts := &file.CreateOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateFileIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	createFileOpts := &file.CreateOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	createFileOpts := &file.CreateOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	createFileOpts := &file.CreateOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	createFileOpts := &file.CreateOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		},
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
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

	createFileOpts := &file.CreateOptions{
		Metadata: testcommon.BasicMetadata,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), createFileOpts)
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

	createFileOpts := &file.CreateOptions{
		Metadata: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), createFileOpts)
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

	createFileOpts := &file.CreateOptions{
		HTTPHeaders: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), createFileOpts)
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

	createFileOpts := &file.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithExpiryAbsolute() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	expiryTimeAbsolute := time.Now().Add(8 * time.Second)
	expiry := file.CreationExpiryTypeAbsolute(expiryTimeAbsolute)
	createFileOpts := &file.CreateOptions{
		Expiry: expiry,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp1.ExpiresOn)
	_require.Equal(expiryTimeAbsolute.UTC().Format(http.TimeFormat), (*resp1.ExpiresOn).UTC().Format(http.TimeFormat))
}

func (s *RecordedTestSuite) TestCreateFileWithExpiryRelativeToNow() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	expiry := file.CreationExpiryTypeRelativeToNow(8 * time.Second)
	createFileOpts := &file.CreateOptions{
		Expiry: expiry,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp1.ExpiresOn)

	time.Sleep(time.Second * 10)
	_, err = fClient.GetProperties(context.Background(), nil)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)
}

func (s *RecordedTestSuite) TestCreateFileWithNeverExpire() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		Expiry: file.CreationExpiryTypeNever{},
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp1.ExpiresOn)
}

func (s *RecordedTestSuite) TestCreateFileWithLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	// should fail since leased
	_, err = fClient.Create(context.Background(), createFileOpts)
	_require.NotNil(err)

	time.Sleep(time.Second * 15)
	resp, err = fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithPermissions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	perms := "0777"
	umask := "0000"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		Permissions: &perms,
		Umask:       &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	resp2, err := fClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp2)
	_require.Equal("rwxrwxrwx", *resp2.Permissions)
}

func (s *RecordedTestSuite) TestCreateFileWithOwnerGroupACLUmask() {
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

	createFileOpts := &file.CreateOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		Umask: &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.Nil(err)
	_require.NotNil(resp)

}

func (s *RecordedTestSuite) TestDeleteFileWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fClient.Create(context.Background(), nil)
	_require.Nil(err)

	deleteOpts := &file.DeleteOptions{
		AccessConditions: nil,
	}

	resp, err := fClient.Delete(context.Background(), deleteOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteFileIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteOpts := &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp1, err := fClient.Delete(context.Background(), deleteOpts)
	_require.Nil(err)
	_require.NotNil(resp1)
}

func (s *RecordedTestSuite) TestDeleteFileIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteOpts := &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.Delete(context.Background(), deleteOpts)
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteFileIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	deleteOpts := &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.Delete(context.Background(), deleteOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteFileIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteOpts := &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.Delete(context.Background(), deleteOpts)
	_require.NotNil(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteFileIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	deleteOpts := &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	_, err = fClient.Delete(context.Background(), deleteOpts)
	_require.Nil(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteFileIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	deleteOpts := &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		},
	}

	_, err = fClient.Delete(context.Background(), deleteOpts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = fClient.SetAccessControl(context.Background(), nil)
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

	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	opts := &file.SetAccessControlOptions{
		Owner:            &owner,
		Group:            &group,
		ACL:              &acl,
		AccessConditions: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)
	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)
	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)
	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		}}

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)
	etag := resp.ETag

	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag
	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		}}

	_, err = fClient.SetAccessControl(context.Background(), opts)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	getACLResp, err := fClient.GetAccessControl(context.Background(), nil)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	// Adding SAS and options
	permissions := sas.FilePermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	sasURL, err := fClient.GetSASURL(permissions, expiry, nil)
	_require.Nil(err)

	fClient2, _ := file.NewClientWithNoCredential(sasURL, nil)

	getACLResp, err := fClient2.GetAccessControl(context.Background(), nil)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	// Adding SAS and options
	permissions := sas.FilePermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(time.Hour)

	sasURL, err := fClient.GetSASURL(permissions, expiry, nil)
	_require.Nil(err)

	fClient2, _ := file.NewClientWithNoCredential(sasURL, nil)

	_, err = fClient2.Delete(context.Background(), nil)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &file.GetAccessControlOptions{
		AccessConditions: nil,
	}

	getACLResp, err := fClient.GetAccessControl(context.Background(), opts)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)
	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	getACLResp, err := fClient.GetAccessControl(context.Background(), opts)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)
	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.GetAccessControl(context.Background(), opts)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)
	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		}}

	getACLResp, err := fClient.GetAccessControl(context.Background(), opts)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.GetAccessControl(context.Background(), opts)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)
	etag := resp.ETag

	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}

	getACLResp, err := fClient.GetAccessControl(context.Background(), opts)
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

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag
	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		}}

	_, err = fClient.GetAccessControl(context.Background(), opts)
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestUpdateAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	acl1 := "user::rwx,group::r--,other::r--"
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = fClient.UpdateAccessControl(context.Background(), acl1, nil)
	_require.Nil(err)

	getACLResp, err := fClient.GetAccessControl(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(acl1, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestRemoveAccessControl() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = fClient.RemoveAccessControl(context.Background(), acl, nil)
	_require.Nil(err)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = fClient.SetMetadata(context.Background(), nil)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &file.SetMetadataOptions{
		Metadata: nil,
	}
	_, err = fClient.SetMetadata(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &file.SetMetadataOptions{
		Metadata: testcommon.BasicMetadata,
	}
	_, err = fClient.SetMetadata(context.Background(), opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.SetMetadataOptions{
		Metadata: testcommon.BasicMetadata,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fClient.SetMetadata(context.Background(), opts)
	_require.Nil(err)
}

func validatePropertiesSet(_require *require.Assertions, fileClient *file.Client, disposition string) {
	resp, err := fileClient.GetProperties(context.Background(), nil)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, nil)
	_require.Nil(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: nil,
	}

	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Nil(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Nil(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Nil(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfMatch: etag,
			},
		}}
	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.Nil(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		},
	}
	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	resp1, err := fClient.Rename(context.Background(), "newName", nil)
	_require.Nil(err)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	renameFileOpts := &file.RenameOptions{
		AccessConditions: nil,
	}

	resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.Nil(err)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfModifiedSince: &currentTime,
			},
		},
	}
	resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.Nil(err)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfModifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfUnmodifiedSince: &currentTime,
			},
		},
	}

	resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfMatch: etag,
			},
		},
	}

	resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NotNil(resp1)
	_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	etag := resp.ETag

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfNoneMatch: etag,
			},
		},
	}

	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NotNil(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestFileUploadStream() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 100 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileUploadSmallStream() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileUploadTinyStream() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 4
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileUploadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 200 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	// create local file
	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	err = os.WriteFile("testFile", content, 0644)
	_require.NoError(err)

	defer func() {
		err = os.Remove("testFile")
		_require.NoError(err)
	}()

	fh, err := os.Open("testFile")
	_require.NoError(err)

	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	hash := md5.New()
	_, err = io.Copy(hash, fh)
	_require.NoError(err)
	//contentMD5 := hash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
	//dResp, err := fClient.DownloadStream(context.Background(), nil)
	//_require.NoError(err)
	//
	//data, err := io.ReadAll(dResp.Body)
	//_require.NoError(err)
	//
	//downloadedMD5Value := md5.Sum(data)
	//downloadedContentMD5 := downloadedMD5Value[:]
	//
	//_require.EqualValues(downloadedContentMD5, contentMD5)
	//
	//gResp2, err := fClient.GetProperties(context.Background(), nil)
	//_require.NoError(err)
	//_require.Equal(*gResp2.ContentLength, fileSize)
	//
	//rangeList, err := fClient.GetRangeList(context.Background(), nil)
	//_require.NoError(err)
	//_require.Len(rangeList.Ranges, 1)
	//_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	//_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (s *RecordedTestSuite) TestSmallFileUploadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	// create local file
	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	err = os.WriteFile("testFile", content, 0644)
	_require.NoError(err)

	defer func() {
		err = os.Remove("testFile")
		_require.NoError(err)
	}()

	fh, err := os.Open("testFile")
	_require.NoError(err)

	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	hash := md5.New()
	_, err = io.Copy(hash, fh)
	_require.NoError(err)
	//contentMD5 := hash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
	//dResp, err := fClient.DownloadStream(context.Background(), nil)
	//_require.NoError(err)
	//
	//data, err := io.ReadAll(dResp.Body)
	//_require.NoError(err)
	//
	//downloadedMD5Value := md5.Sum(data)
	//downloadedContentMD5 := downloadedMD5Value[:]
	//
	//_require.EqualValues(downloadedContentMD5, contentMD5)
	//
	//gResp2, err := fClient.GetProperties(context.Background(), nil)
	//_require.NoError(err)
	//_require.Equal(*gResp2.ContentLength, fileSize)
	//
	//rangeList, err := fClient.GetRangeList(context.Background(), nil)
	//_require.NoError(err)
	//_require.Len(rangeList.Ranges, 1)
	//_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	//_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (s *RecordedTestSuite) TestTinyFileUploadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 10
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	// create local file
	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	err = os.WriteFile("testFile", content, 0644)
	_require.NoError(err)

	defer func() {
		err = os.Remove("testFile")
		_require.NoError(err)
	}()

	fh, err := os.Open("testFile")
	_require.NoError(err)

	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	hash := md5.New()
	_, err = io.Copy(hash, fh)
	_require.NoError(err)
	//contentMD5 := hash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		ChunkSize: 2,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
	//dResp, err := fClient.DownloadStream(context.Background(), nil)
	//_require.NoError(err)
	//
	//data, err := io.ReadAll(dResp.Body)
	//_require.NoError(err)
	//
	//downloadedMD5Value := md5.Sum(data)
	//downloadedContentMD5 := downloadedMD5Value[:]
	//
	//_require.EqualValues(downloadedContentMD5, contentMD5)
	//
	//gResp2, err := fClient.GetProperties(context.Background(), nil)
	//_require.NoError(err)
	//_require.Equal(*gResp2.ContentLength, fileSize)
	//
	//rangeList, err := fClient.GetRangeList(context.Background(), nil)
	//_require.NoError(err)
	//_require.Len(rangeList.Ranges, 1)
	//_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	//_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (s *RecordedTestSuite) TestFileUploadBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 100 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	//md5Value := md5.Sum(content)
	//contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	//dResp, err := fClient.DownloadStream(context.Background(), nil)
	//_require.NoError(err)
	//
	//data, err := io.ReadAll(dResp.Body)
	//_require.NoError(err)
	//
	//downloadedMD5Value := md5.Sum(data)
	//downloadedContentMD5 := downloadedMD5Value[:]
	//
	//_require.EqualValues(downloadedContentMD5, contentMD5)
	//
	//gResp2, err := fClient.GetProperties(context.Background(), nil)
	//_require.NoError(err)
	//_require.Equal(*gResp2.ContentLength, fileSize)
	//
	//rangeList, err := fClient.GetRangeList(context.Background(), nil)
	//_require.NoError(err)
	//_require.Len(rangeList.Ranges, 1)
	//_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	//_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (s *RecordedTestSuite) TestFileUploadSmallBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	//md5Value := md5.Sum(content)
	//contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	//dResp, err := fClient.DownloadStream(context.Background(), nil)
	//_require.NoError(err)
	//
	//data, err := io.ReadAll(dResp.Body)
	//_require.NoError(err)
	//
	//downloadedMD5Value := md5.Sum(data)
	//downloadedContentMD5 := downloadedMD5Value[:]
	//
	//_require.EqualValues(downloadedContentMD5, contentMD5)
	//
	//gResp2, err := fClient.GetProperties(context.Background(), nil)
	//_require.NoError(err)
	//_require.Equal(*gResp2.ContentLength, fileSize)
	//
	//rangeList, err := fClient.GetRangeList(context.Background(), nil)
	//_require.NoError(err)
	//_require.Len(rangeList.Ranges, 1)
	//_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	//_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (s *RecordedTestSuite) TestFileAppendAndFlushData() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), nil)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFileAppendAndFlushDataWithValidation() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFilesystemName(testName)
	fsClient, err := testcommon.GetFilesystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFilesystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.Nil(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	opts := &file.AppendDataOptions{
		TransactionalValidation: file.TransferValidationTypeComputeCRC64(),
	}
	putResp, err := srcFClient.AppendData(context.Background(), 0, rsc, opts)
	_require.Nil(err)
	// _require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.NotNil(putResp.ContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(putResp.ContentCRC64), contentCRC64)

	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), nil)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}
