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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/service"
	"hash/crc64"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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

//func validateFileDeleted(_require *require.Assertions, fileClient *file.Client) {
//	_, err := fileClient.GetAccessControl(context.Background(), nil)
//	_require.Error(err)
//
//	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)
//}

func (s *UnrecordedTestSuite) TestCreateFileAndDeleteWithConnectionString() {

	_require := require.New(s.T())
	testName := s.T().Name()

	connectionString, _ := testcommon.GetGenericConnectionString(testcommon.TestAccountDatalake)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateDirName(testName)
	fileClient, err := file.NewClientFromConnectionString(*connectionString, fileName, filesystemName, nil)

	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fileClient)

	resp, err := fileClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileAndDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	createFileOpts := &file.CreateOptions{
		AccessConditions: nil,
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	createFileOpts := &file.CreateOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)
	_require.Equal(*(resp.IsServerEncrypted), true)
	_require.Equal(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *RecordedTestSuite) TestCreateFileIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateFileIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateFileIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestCreateFileWithNilHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		HTTPHeaders: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *UnrecordedTestSuite) TestCreateFileWithExpiryAbsolute() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	expiryTimeAbsolute := time.Now().Add(8 * time.Second)
	createFileOpts := &file.CreateOptions{
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeAbsolute,
			ExpiresOn:  time.Now().Add(8 * time.Second).UTC().Format(http.TimeFormat),
		},
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp1.ExpiresOn)
	_require.Equal(expiryTimeAbsolute.UTC().Format(http.TimeFormat), (*resp1.ExpiresOn).UTC().Format(http.TimeFormat))
}

func (s *RecordedTestSuite) TestCreateFileWithExpiryNever() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeNeverExpire,
		},
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.Delete(context.Background(), nil)
	_require.NoError(err)

}

func (s *RecordedTestSuite) TestCreateFileWithExpiryRelativeToNow() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeRelativeToNow,
			ExpiresOn:  strconv.FormatInt((8 * time.Second).Milliseconds(), 10),
		},
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp1.ExpiresOn)

	time.Sleep(time.Second * 10)
	_, err = fClient.GetProperties(context.Background(), nil)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)
}

func (s *RecordedTestSuite) TestCreateFileWithNeverExpire() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		Expiry: file.CreateExpiryValues{
			ExpiryType: file.CreateExpiryTypeNeverExpire,
		},
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp1.ExpiresOn)
}

func (s *RecordedTestSuite) TestCreateFileWithLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	// should fail since leased
	_, err = fClient.Create(context.Background(), createFileOpts)
	_require.Error(err)

	time.Sleep(time.Second * 15)
	resp, err = fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestCreateFileWithPermissions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	perms := "0777"
	umask := "0000"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		Permissions: &perms,
		Umask:       &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	resp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.Equal("rwxrwxrwx", *resp2.Permissions)
}

func (s *RecordedTestSuite) TestCreateFileWithOwnerGroupACLUmask() {
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

	createFileOpts := &file.CreateOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
		Umask: &umask,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

}

func (s *RecordedTestSuite) TestDeleteFileWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)

	deleteOpts := &file.DeleteOptions{
		AccessConditions: nil,
	}

	resp, err := fClient.Delete(context.Background(), deleteOpts)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteFileIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	deleteOpts := &file.DeleteOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	resp1, err := fClient.Delete(context.Background(), deleteOpts)
	_require.NoError(err)
	_require.NotNil(resp1)
}

func (s *RecordedTestSuite) TestDeleteFileIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteFileIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteFileIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestDeleteFileIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *RecordedTestSuite) TestDeleteFileIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileSetExpiry() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.SetExpiry(context.Background(), file.SetExpiryValues{ExpiryType: file.SetExpiryTypeNeverExpire}, nil)
	_require.NoError(err)

	res, err := fClient.GetProperties(context.Background(), nil)
	_require.Nil(res.ExpiresOn)
	_require.NoError(err)

	_, err = fClient.SetExpiry(
		context.Background(),
		file.SetExpiryValues{
			ExpiryType: file.SetExpiryTypeRelativeToCreation,
			ExpiresOn:  strconv.Itoa(10),
		},
		nil,
	)
	_require.NoError(err)

	time.Sleep(time.Second * 12)

	_, err = fClient.GetProperties(context.Background(), nil)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)
}

func (s *UnrecordedTestSuite) TestFileSetExpiryTypeAbsoluteTime() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.SetExpiry(
		context.Background(),
		file.SetExpiryValues{
			ExpiryType: file.SetExpiryTypeAbsolute,
			ExpiresOn:  time.Now().Add(5 * time.Second).UTC().Format(http.TimeFormat),
		},
		nil)
	_require.NoError(err)

	time.Sleep(time.Second * 7)

	_, err = fClient.GetProperties(context.Background(), nil)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)

}

func (s *RecordedTestSuite) TestFileSetAccessControlNil() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.SetAccessControl(context.Background(), nil)
	_require.Error(err)

	_require.Equal(err, datalakeerror.MissingParameters)
}

// TODO: write test that fails if you provide permissions and acls
func (s *RecordedTestSuite) TestFileSetAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	opts := &file.SetAccessControlOptions{
		Owner: &owner,
		Group: &group,
		ACL:   &acl,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.SetAccessControl(context.Background(), opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetAccessControlWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	owner := "4cf4e284-f6a8-4540-b53e-c3469af032dc"
	group := owner
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	opts := &file.SetAccessControlOptions{
		Owner:            &owner,
		Group:            &group,
		ACL:              &acl,
		AccessConditions: nil,
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.SetAccessControl(context.Background(), opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetAccessControlIfModifiedSinceTrue() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetAccessControlIfModifiedSinceFalse() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileSetAccessControlIfUnmodifiedSinceTrue() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetAccessControlIfUnmodifiedSinceFalse() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileSetAccessControlIfETagMatch() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetAccessControlIfETagMatchFalse() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileGetAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	getACLResp, err := fClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *UnrecordedTestSuite) TestFileGetAccessControlWithSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
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
	_require.NoError(err)

	fClient2, _ := file.NewClientWithNoCredential(sasURL, nil)

	getACLResp, err := fClient2.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *UnrecordedTestSuite) TestFileDeleteWithSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)

	fClient2, _ := file.NewClientWithNoCredential(sasURL, nil)

	_, err = fClient2.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *UnrecordedTestSuite) TestFileEncryptionScopeSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.DataLakeEncryptionScopeEnvVar)
	_require.Nil(err)

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDatalake)
	_require.NoError(err)

	perms := sas.FilePermissions{Read: true, Create: true, Write: true, Move: true, Delete: true, List: true}
	sasQueryParams, err := sas.DatalakeSignatureValues{
		Protocol:        sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:      time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		FileSystemName:  filesystemName,
		Permissions:     perms.String(),
		EncryptionScope: encryptionScope,
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	sasToken := sasQueryParams.Encode()

	srcFileClient, err := file.NewClientWithNoCredential(fsClient.DFSURL()+"/file?"+sasToken, nil)
	_require.NoError(err)
	_require.NotNil(srcFileClient)

	_, err = srcFileClient.Create(context.Background(), nil)
	_require.NoError(err)

	response, err := srcFileClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NoError(err)
	_require.Equal(encryptionScope, *response.EncryptionScope)

}

func (s *UnrecordedTestSuite) TestAccountEncryptionScopeSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.DataLakeEncryptionScopeEnvVar)
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDatalake)
	_require.Nil(err)

	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:        sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:      time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:     to.Ptr(sas.AccountPermissions{Read: true, Create: true, Write: true, Delete: true}).String(),
		ResourceTypes:   to.Ptr(sas.AccountResourceTypes{Service: true, Container: true, Object: true}).String(),
		EncryptionScope: encryptionScope,
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	sasToken := sasQueryParams.Encode()

	srcFileClient, err := file.NewClientWithNoCredential(fsClient.DFSURL()+"/file?"+sasToken, nil)
	_require.NoError(err)
	_require.NotNil(srcFileClient)

	resp, err := srcFileClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// create local file
	_, content := generateData(10 * 1024)
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

	// upload the file
	err = srcFileClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	response, err := srcFileClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	testcommon.DeleteFileSystem(context.Background(), _require, fsClient)
	_require.Equal(encryptionScope, *response.EncryptionScope)

	// validate the data downloaded
	downloadedData, err := io.ReadAll(response.Body)
	_require.NoError(err)
	_require.Equal(len(content), len(downloadedData))
	_require.EqualValues(content, downloadedData)
}

func (s *UnrecordedTestSuite) TestGetUserDelegationEncryptionScopeSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDatalake)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := service.NewClient("https://"+accountName+".dfs.core.windows.net/", cred, nil)
	_require.NoError(err)

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient := testcommon.CreateNewFileSystem(context.Background(), _require, filesystemName, svcClient)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	// Set current and past time and create key
	currentTime := time.Now().UTC().Add(-10 * time.Second)
	pastTime := currentTime.Add(48 * time.Hour)
	info := service.KeyInfo{
		Start:  to.Ptr(currentTime.UTC().Format(sas.TimeFormat)),
		Expiry: to.Ptr(pastTime.UTC().Format(sas.TimeFormat)),
	}

	udc, err := svcClient.GetUserDelegationCredential(context.Background(), info, nil)
	_require.NoError(err)

	// get permissions and details for sas
	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.DataLakeEncryptionScopeEnvVar)
	_require.Nil(err)

	// Create Blob Signature Values with desired permissions and sign with user delegation credential
	perms := sas.FilePermissions{Read: true, Create: true, Write: true, Move: true, Delete: true, List: true}
	sasQueryParams, err := sas.DatalakeSignatureValues{
		Protocol:        sas.ProtocolHTTPS, // Users MUST use HTTPS (not HTTP)
		StartTime:       time.Now().UTC().Add(time.Second * -10),
		ExpiryTime:      time.Now().UTC().Add(15 * time.Minute), // 15 minutes before expiration
		FileSystemName:  filesystemName,
		Permissions:     perms.String(),
		EncryptionScope: encryptionScope,
	}.SignWithUserDelegation(udc)
	_require.Nil(err)

	sasURL := fsClient.DFSURL() + "/file?" + sasQueryParams.Encode()
	// This URL can be used to authenticate requests now
	srcFileClient, err := file.NewClientWithNoCredential(sasURL, nil)
	_require.NoError(err)

	_, err = srcFileClient.Create(context.Background(), nil)
	_require.NoError(err)

	response, err := srcFileClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NoError(err)
	_require.Equal(encryptionScope, *response.EncryptionScope)
}

func (s *RecordedTestSuite) TestFileGetAccessControlWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &file.GetAccessControlOptions{
		AccessConditions: nil,
	}

	getACLResp, err := fClient.GetAccessControl(context.Background(), opts)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestFileGetAccessControlIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
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
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestFileGetAccessControlIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileGetAccessControlIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)
	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		}}

	getACLResp, err := fClient.GetAccessControl(context.Background(), opts)
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestFileGetAccessControlIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
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
	_require.Error(err)

	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileGetAccessControlIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
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
	_require.NoError(err)
	_require.Equal(acl, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestFileGetAccessControlIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag
	opts := &file.GetAccessControlOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		}}

	_, err = fClient.GetAccessControl(context.Background(), opts)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileUpdateAccessControl() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	acl := "user::rwx,group::r-x,other::rwx"
	acl1 := "user::rwx,group::r--,other::r--"
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createOpts := &file.CreateOptions{
		ACL: &acl,
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.UpdateAccessControl(context.Background(), acl1, nil)
	_require.NoError(err)

	getACLResp, err := fClient.GetAccessControl(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(acl1, *getACLResp.ACL)
}

func (s *RecordedTestSuite) TestFileRemoveAccessControl() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.RemoveAccessControl(context.Background(), acl, nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetMetadataWithBasicMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetMetadataWithAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.SetMetadataOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = fClient.SetMetadata(context.Background(), testcommon.BasicMetadata, opts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestFileSetMetadataWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	resp, err := fClient.Create(context.Background(), &file.CreateOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &file.SetMetadataOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	res, err := fClient.SetMetadata(context.Background(), testcommon.BasicMetadata, opts)
	_require.NoError(err)
	_require.Equal(*(res.IsServerEncrypted), true)
	_require.Equal(res.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func validatePropertiesSet(_require *require.Assertions, fileClient *file.Client, disposition string) {
	resp, err := fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.ContentDisposition, disposition)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, nil)
	_require.NoError(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeadersWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: nil,
	}

	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.NoError(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeadersIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeadersIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeadersIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeadersIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeadersIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	opts := &file.SetHTTPHeadersOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfMatch: etag,
			},
		}}
	_, err = fClient.SetHTTPHeaders(context.Background(), testcommon.BasicHeaders, opts)
	_require.NoError(err)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestFileSetHTTPHeadersIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.ConditionNotMet)
}

func (s *UnrecordedTestSuite) TestFileRenameUsingSAS() {
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

	perms := sas.FilePermissions{Read: true, Create: true, Write: true, Move: true, Delete: true, List: true}
	sasQueryParams, err := sas.DatalakeSignatureValues{
		Protocol:       sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:     time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		FileSystemName: filesystemName,
		Permissions:    perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	sasToken := sasQueryParams.Encode()

	srcFileClient, err := file.NewClientWithNoCredential(fsClient.DFSURL()+"/file1?"+sasToken, nil)
	_require.NoError(err)

	_, err = srcFileClient.Create(context.Background(), nil)
	_require.NoError(err)

	destPathWithSAS := "file2?" + sasToken
	_, err = srcFileClient.Rename(context.Background(), destPathWithSAS, nil)
	_require.NoError(err)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.PathNotFound)
}

func (s *RecordedTestSuite) TestRenameNoOptions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, err = fClient.Rename(context.Background(), "newName", nil)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestRenameFileWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	renameFileOpts := &file.RenameOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
}

func (s *RecordedTestSuite) TestRenameFileWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	renameFileOpts := &file.RenameOptions{
		AccessConditions: nil,
	}

	//resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
	//_require.NotNil(resp1)
	//_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfModifiedSince: &currentTime,
			},
		},
	}
	//resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
	//_require.NotNil(resp1)
	//_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfModifiedSince: &currentTime,
			},
		},
	}

	//_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)

	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestRenameFileIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfUnmodifiedSince: &currentTime,
			},
		},
	}

	//resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
	//_require.NotNil(resp1)
	//_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfUnmodifiedSince: &currentTime,
			},
		},
	}

	//_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)

	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *RecordedTestSuite) TestRenameFileIfETagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfMatch: etag,
			},
		},
	}

	//resp1, err := fClient.Rename(context.Background(), "newName", renameFileOpts)
	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_require.NoError(err)
	//_require.NotNil(resp1)
	//_require.Contains(resp1.NewFileClient.DFSURL(), "newName")
}

func (s *RecordedTestSuite) TestRenameFileIfETagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	etag := resp.ETag

	renameFileOpts := &file.RenameOptions{
		SourceAccessConditions: &file.SourceAccessConditions{
			SourceModifiedAccessConditions: &file.SourceModifiedAccessConditions{
				SourceIfNoneMatch: etag,
			},
		},
	}

	//_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)
	_, err = fClient.Rename(context.Background(), "newName", renameFileOpts)

	_require.Error(err)
	testcommon.ValidateErrorCode(_require, err, datalakeerror.SourceConditionNotMet)
}

func (s *UnrecordedTestSuite) TestFileUploadDownloadStream() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 100 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

}

func (s *RecordedTestSuite) TestFileUploadDownloadSmallStream() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileUploadTinyStream() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 4
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
}

func (s *UnrecordedTestSuite) TestFileUploadDownloadStreamWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 1 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), &file.CreateOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), &file.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
	_require.Equal(true, *(dResp.IsServerEncrypted))
	_require.Equal(testcommon.TestCPKByValue.EncryptionKeySHA256, dResp.EncryptionKeySHA256)
}

func (s *UnrecordedTestSuite) TestFileUploadDownloadStreamWithCPKNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 1 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), &file.CreateOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)

	err = fClient.UploadStream(context.Background(), streaming.NopCloser(bytes.NewReader(content)), &file.UploadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), &file.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	_, err = fClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{})
	_require.Error(err)
	_require.ErrorContains(err, "PathUsesCustomerSpecifiedEncryption")
}

func (s *UnrecordedTestSuite) TestFileUploadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 100 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	contentMD5 := hash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
}

func (s *RecordedTestSuite) TestSmallFileUploadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// create local file
	_, content := testcommon.GenerateData(int(fileSize))
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
	contentMD5 := hash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
}

func (s *RecordedTestSuite) TestSmallFileUploadFileWithAccessConditionsAndHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	createFileOpts := &file.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	// create local file
	_, content := testcommon.GenerateData(int(fileSize))
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
	contentMD5 := hash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		}},
		HTTPHeaders: &testcommon.BasicHeaders,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), &file.GetPropertiesOptions{
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		},
		}})
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestTinyFileUploadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// create local file
	_, content := testcommon.GenerateData(int(fileSize))
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
	contentMD5 := hash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		ChunkSize: 2,
	})
	_require.NoError(err)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
}

func (s *UnrecordedTestSuite) TestFileUploadBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 100 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
}

func (s *RecordedTestSuite) TestFileUploadSmallBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
}

func (s *RecordedTestSuite) TestFileUploadSmallBufferWithAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}
	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NotNil(err)

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		}},
	})
	_require.Nil(err)

	gResp2, err := fClient.GetProperties(context.Background(), &file.GetPropertiesOptions{
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		},
		}})
	_require.Nil(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
}

func (s *RecordedTestSuite) TestFileUploadSmallBufferWithHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
		HTTPHeaders: &testcommon.BasicHeaders,
	})

	_require.NoError(err)
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(data)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)
	validatePropertiesSet(_require, fClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestDownloadDataContentMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	createFileOpts := &file.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	resp, err := fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
		HTTPHeaders: &testcommon.BasicHeaders,
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		},
		}})

	_require.NoError(err)
	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	options := file.DownloadStreamOptions{
		Range: &file.HTTPRange{
			Count:  3,
			Offset: 10,
		},
		RangeGetContentMD5: to.Ptr(true),
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		},
		}}
	resp1, err := fClient.DownloadStream(context.Background(), &options)
	_require.Nil(err)
	mdf := md5.Sum(content[10:13])
	_require.Equal(resp1.ContentMD5, mdf[:])
}

func (s *RecordedTestSuite) TestFileAppendAndFlushData() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
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

func (s *UnrecordedTestSuite) TestFileAppendAndFlushDataWithValidation() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
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
	_require.NoError(err)
	// _require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.NotNil(putResp.ContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(putResp.ContentCRC64), contentCRC64)

	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), nil)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFileAppendAndFlushDataWithEmptyOpts() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	opts := &file.AppendDataOptions{}
	opts1 := &file.FlushDataOptions{}

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, opts)
	_require.NoError(err)

	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts1)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFileAppendAndFlushDataWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	createOptions := &file.CreateOptions{CPKInfo: &testcommon.TestCPKByValue}
	resp, err := srcFClient.Create(context.Background(), createOptions)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	opts := &file.AppendDataOptions{CPKInfo: &testcommon.TestCPKByValue}
	opts1 := &file.FlushDataOptions{CPKInfo: &testcommon.TestCPKByValue}

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, opts)
	_require.NoError(err)

	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts1)
	_require.NoError(err)
	getPropertiesOptions := &file.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue}
	gResp2, err := srcFClient.GetProperties(context.Background(), getPropertiesOptions)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
	_require.Equal(true, *(gResp2.IsServerEncrypted))
	_require.Equal(testcommon.TestCPKByValue.EncryptionKeySHA256, gResp2.EncryptionKeySHA256)
}

func (s *RecordedTestSuite) TestFileAppendAndFlushDataWithLeasedFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	createFileOpts := &file.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}
	resp, err := srcFClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	opts := &file.AppendDataOptions{LeaseAccessConditions: &file.LeaseAccessConditions{
		LeaseID: proposedLeaseIDs[0],
	}}

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.Error(err)

	_, err = rsc.Seek(0, io.SeekStart)
	_require.NoError(err)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, opts)
	_require.NoError(err)

	opts1 := &file.FlushDataOptions{AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
		LeaseID: proposedLeaseIDs[0],
	}}}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts1)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), &file.GetPropertiesOptions{
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		}}})
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFileAppendAndFlushAndDownloadDataWithLeasedFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	createFileOpts := &file.CreateOptions{
		ProposedLeaseID: proposedLeaseIDs[0],
		LeaseDuration:   to.Ptr(int64(15)),
	}
	resp, err := srcFClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	opts := &file.AppendDataOptions{LeaseAccessConditions: &file.LeaseAccessConditions{
		LeaseID: proposedLeaseIDs[0],
	}}

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, opts)
	_require.NoError(err)

	_, err = rsc.Seek(0, io.SeekStart)
	_require.NoError(err)

	_, err = srcFClient.AppendData(context.Background(), int64(contentSize), rsc, opts)
	_require.NoError(err)

	opts1 := &file.FlushDataOptions{AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
		LeaseID: proposedLeaseIDs[0],
	}}}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize)*2, opts1)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), &file.GetPropertiesOptions{
		AccessConditions: &file.AccessConditions{LeaseAccessConditions: &file.LeaseAccessConditions{
			LeaseID: proposedLeaseIDs[0],
		}}})
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize)*2)

	destBuffer := make([]byte, contentSize*2)
	cnt, err := srcFClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   8 * 1024,
		Concurrency: 5,
	})
	_require.NoError(err)
	_require.Equal(cnt, int64(contentSize*2))
}

func (s *RecordedTestSuite) TestAppendAndFlushFileWithHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	opts := &file.FlushDataOptions{HTTPHeaders: &testcommon.BasicHeaders}

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
	validatePropertiesSet(_require, srcFClient, *testcommon.BasicHeaders.ContentDisposition)
}

func (s *RecordedTestSuite) TestFlushWithNilAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	opts := &file.FlushDataOptions{
		AccessConditions: nil,
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFlushWithEmptyAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	opts := &file.FlushDataOptions{
		AccessConditions: &file.AccessConditions{},
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

// ==========================================

func (s *RecordedTestSuite) TestFlushIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)
	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.FlushDataOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFlushIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &file.FlushDataOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.Error(err)
}

func (s *RecordedTestSuite) TestFlushIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, 10)

	opts := &file.FlushDataOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFlushIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	_, err = srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	opts := &file.FlushDataOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.Error(err)
}

func (s *RecordedTestSuite) TestFlushIfEtagMatch() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	resp1, err := srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)
	etag := resp1.ETag

	opts := &file.FlushDataOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfMatch: etag,
			},
		},
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

func (s *RecordedTestSuite) TestFlushIfEtagMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	srcFileName := "src" + testcommon.GenerateFileName(testName)

	srcFClient, err := testcommon.GetFileClient(filesystemName, srcFileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := srcFClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	contentSize := 1024 * 8 // 8KB
	rsc, _ := testcommon.GenerateData(contentSize)

	resp1, err := srcFClient.AppendData(context.Background(), 0, rsc, nil)
	_require.NoError(err)
	etag := resp1.ETag

	opts := &file.FlushDataOptions{
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfNoneMatch: etag,
			},
		},
	}
	_, err = srcFClient.FlushData(context.Background(), int64(contentSize), opts)
	_require.NoError(err)

	gResp2, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, int64(contentSize))
}

// TODO: test retain uncommitted data

func (s *UnrecordedTestSuite) TestFileDownloadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 100 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	destFileName := "BigFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	_require.NoError(err)
	defer func(name string) {
		err = os.Remove(name)
		_require.NoError(err)
	}(destFileName)
	defer func(destFile *os.File) {
		err = destFile.Close()
		_require.NoError(err)
	}(destFile)

	cnt, err := fClient.DownloadFile(context.Background(), destFile, &file.DownloadFileOptions{
		ChunkSize:   10 * 1024 * 1024,
		Concurrency: 5,
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	hash := md5.New()
	_, err = io.Copy(hash, destFile)
	_require.NoError(err)
	downloadedContentMD5 := hash.Sum(nil)

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileUploadDownloadSmallFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// create local file
	_, content := testcommon.GenerateData(int(fileSize))
	srcFileName := "testFileUpload"
	err = os.WriteFile(srcFileName, content, 0644)
	_require.NoError(err)
	defer func() {
		err = os.Remove(srcFileName)
		_require.NoError(err)
	}()
	fh, err := os.Open(srcFileName)
	_require.NoError(err)
	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	srcHash := md5.New()
	_, err = io.Copy(srcHash, fh)
	_require.NoError(err)
	contentMD5 := srcHash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	destFileName := "SmallFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	_require.NoError(err)
	defer func(name string) {
		err = os.Remove(name)
		_require.NoError(err)
	}(destFileName)

	cnt, err := fClient.DownloadFile(context.Background(), destFile, &file.DownloadFileOptions{
		ChunkSize:   2 * 1024,
		Concurrency: 5,
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	err = destFile.Close()
	_require.NoError(err)

	newDestFileHandle, err := os.Open(destFileName)
	_require.NoError(err)
	defer func(file *os.File) {
		err = file.Close()
		_require.NoError(err)
	}(newDestFileHandle)

	destHash := md5.New()
	_, err = io.Copy(destHash, newDestFileHandle)
	_require.NoError(err)
	downloadedContentMD5 := destHash.Sum(nil)

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

}

func (s *RecordedTestSuite) TestFileUploadDownloadSmallFileWithRange() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// create local file
	_, content := testcommon.GenerateData(int(fileSize))
	srcFileName := "testFileUpload"
	err = os.WriteFile(srcFileName, content, 0644)
	_require.NoError(err)
	defer func() {
		err = os.Remove(srcFileName)
		_require.NoError(err)
	}()
	fh, err := os.Open(srcFileName)
	_require.NoError(err)
	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	srcHash := md5.New()
	_, err = io.Copy(srcHash, fh)
	_require.NoError(err)
	contentMD5 := srcHash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	destFileName := "SmallFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	_require.NoError(err)
	defer func(name string) {
		err = os.Remove(name)
		_require.NoError(err)
	}(destFileName)

	cnt, err := fClient.DownloadFile(context.Background(), destFile, &file.DownloadFileOptions{
		ChunkSize:   2 * 1024,
		Concurrency: 5,
		Range: &file.HTTPRange{
			Offset: 0,
			Count:  10 * 1024,
		},
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	err = destFile.Close()
	_require.NoError(err)

	newDestFileHandle, err := os.Open(destFileName)
	_require.NoError(err)
	defer func(file *os.File) {
		err = file.Close()
		_require.NoError(err)
	}(newDestFileHandle)

	destHash := md5.New()
	_, err = io.Copy(destHash, newDestFileHandle)
	_require.NoError(err)
	downloadedContentMD5 := destHash.Sum(nil)

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileUploadDownloadSmallFileWithAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	// create local file
	_, content := testcommon.GenerateData(int(fileSize))
	srcFileName := "testFileUpload"
	err = os.WriteFile(srcFileName, content, 0644)
	_require.NoError(err)
	defer func() {
		err = os.Remove(srcFileName)
		_require.NoError(err)
	}()
	fh, err := os.Open(srcFileName)
	_require.NoError(err)
	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	srcHash := md5.New()
	_, err = io.Copy(srcHash, fh)
	_require.NoError(err)
	contentMD5 := srcHash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	destFileName := "SmallFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	_require.NoError(err)
	defer func(name string) {
		err = os.Remove(name)
		_require.NoError(err)
	}(destFileName)

	cnt, err := fClient.DownloadFile(context.Background(), destFile, &file.DownloadFileOptions{
		ChunkSize:   2 * 1024,
		Concurrency: 5,
		Range: &file.HTTPRange{
			Offset: 0,
			Count:  10 * 1024,
		},
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	err = destFile.Close()
	_require.NoError(err)

	newDestFileHandle, err := os.Open(destFileName)
	_require.NoError(err)
	defer func(file *os.File) {
		err = file.Close()
		_require.NoError(err)
	}(newDestFileHandle)

	destHash := md5.New()
	_, err = io.Copy(destHash, newDestFileHandle)
	_require.NoError(err)
	downloadedContentMD5 := destHash.Sum(nil)

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileUploadDownloadSmallFileWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), &file.CreateOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	})
	_require.NoError(err)
	_require.NotNil(resp)

	// create local file
	_, content := testcommon.GenerateData(int(fileSize))
	srcFileName := "testFileUpload"
	err = os.WriteFile(srcFileName, content, 0644)
	_require.NoError(err)
	defer func() {
		err = os.Remove(srcFileName)
		_require.NoError(err)
	}()
	fh, err := os.Open(srcFileName)
	_require.NoError(err)
	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	srcHash := md5.New()
	_, err = io.Copy(srcHash, fh)
	_require.NoError(err)
	contentMD5 := srcHash.Sum(nil)

	err = fClient.UploadFile(context.Background(), fh, &file.UploadFileOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
		CPKInfo:     &testcommon.TestCPKByValue,
	})
	_require.NoError(err)

	destFileName := "SmallFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	_require.NoError(err)
	defer func(name string) {
		err = os.Remove(name)
		_require.NoError(err)
	}(destFileName)

	cnt, err := fClient.DownloadFile(context.Background(), destFile, &file.DownloadFileOptions{
		ChunkSize:   2 * 1024,
		Concurrency: 5,
		Range: &file.HTTPRange{
			Offset: 0,
			Count:  10 * 1024,
		},
		CPKInfo: &testcommon.TestCPKByValue,
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	err = destFile.Close()
	_require.NoError(err)

	newDestFileHandle, err := os.Open(destFileName)
	_require.NoError(err)
	defer func(file *os.File) {
		err = file.Close()
		_require.NoError(err)
	}(newDestFileHandle)

	destHash := md5.New()
	_, err = io.Copy(destHash, newDestFileHandle)
	_require.NoError(err)
	downloadedContentMD5 := destHash.Sum(nil)

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), &file.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
	_require.Equal(*(gResp2.IsServerEncrypted), true)
	_require.Equal(gResp2.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *RecordedTestSuite) TestFileUploadDownloadWithProgress() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	bytesUploaded := int64(0)
	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
		Progress: func(bytesTransferred int64) {
			_require.GreaterOrEqual(bytesTransferred, bytesUploaded)
			bytesUploaded = bytesTransferred
		},
	})
	_require.NoError(err)
	_require.Equal(bytesUploaded, fileSize)

	destBuffer := make([]byte, fileSize)
	bytesDownloaded := int64(0)
	cnt, err := fClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   2 * 1024,
		Concurrency: 5,
		Progress: func(bytesTransferred int64) {
			_require.GreaterOrEqual(bytesTransferred, bytesDownloaded)
			bytesDownloaded = bytesTransferred
		},
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)
	_require.Equal(bytesDownloaded, fileSize)

	downloadedMD5Value := md5.Sum(destBuffer)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *UnrecordedTestSuite) TestFileDownloadBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 100 * 1024 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	destBuffer := make([]byte, fileSize)
	cnt, err := fClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   10 * 1024 * 1024,
		Concurrency: 5,
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	downloadedMD5Value := md5.Sum(destBuffer)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileDownloadSmallBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	destBuffer := make([]byte, fileSize)
	cnt, err := fClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   10 * 1024 * 1024,
		Concurrency: 5,
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	downloadedMD5Value := md5.Sum(destBuffer)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileDownloadSmallBufferWithHTTPRange() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content[0 : fileSize/2])
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	destBuffer := make([]byte, fileSize/2)
	cnt, err := fClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   10 * 1024 * 1024,
		Concurrency: 5,
		Range: &file.HTTPRange{
			Offset: 0,
			Count:  5 * 1024,
		},
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize/2)

	downloadedMD5Value := md5.Sum(destBuffer)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileDownloadSmallBufferWithAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err := fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content[0 : fileSize/2])
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(resp.Date, -10)

	destBuffer := make([]byte, fileSize/2)
	cnt, err := fClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   10 * 1024 * 1024,
		Concurrency: 5,
		AccessConditions: &file.AccessConditions{
			ModifiedAccessConditions: &file.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
		Range: &file.HTTPRange{
			Offset: 0,
			Count:  5 * 1024,
		},
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize/2)

	downloadedMD5Value := md5.Sum(destBuffer)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
}

func (s *RecordedTestSuite) TestFileDownloadBufferWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()

	filesystemName := testcommon.GenerateFileSystemName(testName)
	fsClient, err := testcommon.GetFileSystemClient(filesystemName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	defer testcommon.DeleteFileSystem(context.Background(), _require, fsClient)

	_, err = fsClient.Create(context.Background(), nil)
	_require.NoError(err)

	var fileSize int64 = 10 * 1024
	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	fileCreateOpts := &file.CreateOptions{CPKInfo: &testcommon.TestCPKByValue}
	resp, err := fClient.Create(context.Background(), fileCreateOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content[0:fileSize])
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	})
	_require.NoError(err)

	destBuffer := make([]byte, fileSize)
	_, err = fClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	})
	_require.NoError(err)

	downloadedMD5Value := md5.Sum(destBuffer)
	downloadedContentMD5 := downloadedMD5Value[:]

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), &file.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue})
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)
	_require.Equal(*(gResp2.IsServerEncrypted), true)
	_require.Equal(gResp2.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *RecordedTestSuite) TestFileGetPropertiesResponseCapture() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, dirName+"/"+fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	resp, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)

	// This tests file.NewClient
	var respFromCtxFile *http.Response
	ctxWithRespFile := runtime.WithCaptureResponse(context.Background(), &respFromCtxFile)
	resp2, err := fClient.GetProperties(ctxWithRespFile, nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.NotNil(respFromCtxFile) // validate that the respFromCtx is actually populated
	_require.Equal("file", respFromCtxFile.Header.Get("x-ms-resource-type"))

	// This tests filesystem.NewClient
	fClient = fsClient.NewFileClient(dirName + "/" + fileName)
	var respFromCtxFs *http.Response
	ctxWithRespFs := runtime.WithCaptureResponse(context.Background(), &respFromCtxFs)
	resp2, err = fClient.GetProperties(ctxWithRespFs, nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.NotNil(respFromCtxFs) // validate that the respFromCtx is actually populated
	_require.Equal("file", respFromCtxFs.Header.Get("x-ms-resource-type"))

	// This tests service.NewClient
	serviceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	fsClient = serviceClient.NewFileSystemClient(filesystemName)
	dirClient = fsClient.NewDirectoryClient(dirName)
	fClient, err = dirClient.NewFileClient(fileName)
	_require.NoError(err)
	var respFromCtxService *http.Response
	ctxWithRespService := runtime.WithCaptureResponse(context.Background(), &respFromCtxService)
	resp2, err = fClient.GetProperties(ctxWithRespService, nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.NotNil(respFromCtxService) // validate that the respFromCtx is actually populated
	_require.Equal("file", respFromCtxService.Header.Get("x-ms-resource-type"))

	// This tests directory.NewClient
	var respFromCtxDir *http.Response
	ctxWithRespDir := runtime.WithCaptureResponse(context.Background(), &respFromCtxDir)
	dirClient, err = testcommon.GetDirClient(filesystemName, dirName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)
	fClient, err = dirClient.NewFileClient(fileName)
	_require.NoError(err)
	resp2, err = fClient.GetProperties(ctxWithRespDir, nil)
	_require.NoError(err)
	_require.NotNil(resp2)
	_require.NotNil(respFromCtxDir) // validate that the respFromCtx is actually populated
	_require.Equal("file", respFromCtxDir.Header.Get("x-ms-resource-type"))
}

func (s *RecordedTestSuite) TestFileGetPropertiesWithCPK() {
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

	fileName := testcommon.GenerateFileName(testName)
	fClient, err := testcommon.GetFileClient(filesystemName, dirName+"/"+fileName, s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	createFileOpts := &file.CreateOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	resp, err = fClient.Create(context.Background(), createFileOpts)
	_require.NoError(err)
	_require.NotNil(resp)

	GetPropertiesOpts := &file.GetPropertiesOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	// This tests file.NewClient
	var respFromCtxFile *http.Response
	ctxWithRespFile := runtime.WithCaptureResponse(context.Background(), &respFromCtxFile)
	response, err := fClient.GetProperties(ctxWithRespFile, GetPropertiesOpts)
	_require.NoError(err)
	_require.NotNil(response)
	_require.NotNil(respFromCtxFile.Header.Get("x-ms-encryption-key-sha256")) // validate that the x-ms-encryption-key-sha256 is actually populated
	_require.Equal(testcommon.TestCPKByValue.EncryptionKeySHA256, response.EncryptionKeySHA256)
}

func (s *UnrecordedTestSuite) TestFileCreateDeleteUsingOAuth() {
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

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".dfs.core.windows.net/" + filesystemName + "/" + fileName

	fClient, err := file.NewClient(fileURL, cred, nil)
	_require.NoError(err)

	_, err = fClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}
