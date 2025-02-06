//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"hash/crc64"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running file Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &FileRecordedTestsSuite{})
		suite.Run(t, &FileUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &FileRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &FileRecordedTestsSuite{})
	}
}

func (f *FileRecordedTestsSuite) SetupSuite() {
	f.proxy = testcommon.SetupSuite(&f.Suite)
}

func (f *FileRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&f.Suite, f.proxy)
}

func (f *FileRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(f.T(), suite, test)
}

func (f *FileRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(f.T(), suite, test)
}

func (f *FileUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (f *FileUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type FileRecordedTestsSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type FileUnrecordedTestsSuite struct {
	suite.Suite
}

func (f *FileRecordedTestsSuite) TestFileNewFileClient() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := shareClient.NewDirectoryClient(dirName)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := dirClient.NewFileClient(fileName)

	correctURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName + "/" + fileName
	_require.Equal(fileClient.URL(), correctURL)

	rootFileClient := shareClient.NewRootDirectoryClient().NewFileClient(fileName)

	correctURL = "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName
	_require.Equal(rootFileClient.URL(), correctURL)
}

func (f *FileRecordedTestsSuite) TestFileCreateUsingSharedKey() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName + "/" + dirName + "/" + fileName

	options := &file.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClient, err := file.NewClientWithSharedKeyCredential(fileURL, cred, options)
	_require.NoError(err)

	// creating file where directory does not exist gives ParentNotFound error
	_, err = fileClient.Create(context.Background(), 1024, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)

	testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	resp, err := fileClient.Create(context.Background(), 1024, nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	fileAttributes, err := file.ParseNTFSFileAttributes(resp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
}

func (f *FileRecordedTestsSuite) TestFileCreateUsingConnectionString() {
	_require := require.New(f.T())
	testName := f.T().Name()

	connString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)
	options := &file.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClient1, err := file.NewClientFromConnectionString(*connString, shareName, fileName, options)
	_require.NoError(err)

	resp, err := fileClient1.Create(context.Background(), 1024, nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	filePath := dirName + "/" + fileName
	fileClient2, err := file.NewClientFromConnectionString(*connString, shareName, filePath, options)
	_require.NoError(err)

	_, err = fileClient2.Create(context.Background(), 1024, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)

	testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	// using '\' as path separator
	filePath = dirName + "\\" + fileName
	fileClient3, err := file.NewClientFromConnectionString(*connString, shareName, filePath, options)
	_require.NoError(err)

	resp, err = fileClient3.Create(context.Background(), 1024, nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
}

func (f *FileUnrecordedTestsSuite) TestFileClientUsingSAS() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := dirClient.NewFileClient(fileName)

	permissions := sas.FilePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)

	fileSASURL, err := fileClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	fileSASClient, err := file.NewClientWithNoCredential(fileSASURL, nil)
	_require.NoError(err)

	_, err = fileSASClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)

	_, err = fileSASClient.Create(context.Background(), 1024, nil)
	_require.NoError(err)

	resp, err := fileSASClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
}

func (f *FileRecordedTestsSuite) TestFileCreateDeleteDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	rootDirClient := shareClient.NewRootDirectoryClient()
	_require.NoError(err)

	fClient := rootDirClient.NewFileClient(fileName)

	// Create and delete file in root directory.
	cResp, err := fClient.Create(context.Background(), 1024, nil)
	_require.NoError(err)
	_require.NotNil(cResp.ETag)
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotNil(cResp.RequestID)
	_require.NotNil(cResp.Version)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.IsServerEncrypted)

	delResp, err := fClient.Delete(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(delResp.RequestID)
	_require.NotNil(delResp.Version)
	_require.Equal(delResp.Date.IsZero(), false)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	// Create and delete file in named directory.
	afClient := dirClient.NewFileClient(fileName)

	cResp, err = afClient.Create(context.Background(), 1024, nil)
	_require.NoError(err)
	_require.NotNil(cResp.ETag)
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotNil(cResp.RequestID)
	_require.NotNil(cResp.Version)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.IsServerEncrypted)

	delResp, err = afClient.Delete(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(delResp.RequestID)
	_require.NotNil(delResp.Version)
	_require.Equal(delResp.Date.IsZero(), false)
}

func (f *FileRecordedTestsSuite) TestFileCreateNonDefaultMetadataNonEmpty() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 1024, &file.CreateOptions{
		Metadata: testcommon.BasicMetadata,
	})
	_require.NoError(err)

	resp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Metadata, len(testcommon.BasicMetadata))
	for k, v := range resp.Metadata {
		val := testcommon.BasicMetadata[strings.ToLower(k)]
		_require.NotNil(val)
		_require.Equal(*v, *val)
	}
}

func (f *FileRecordedTestsSuite) TestFileCreateRenameFilePermissionFormatDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 1024, &file.CreateOptions{
		FilePermissionFormat: (*file.PermissionFormat)(to.Ptr(testcommon.FilePermissionFormatSddl)),
		Permissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)

	_, err = fClient.Rename(context.Background(), "file2", &file.RenameOptions{
		FilePermissionFormat: (*file.PermissionFormat)(to.Ptr(testcommon.FilePermissionBinary)),
		Permissions: &file.Permissions{
			Permission: &testcommon.SampleBinary,
		},
	})
	_require.NoError(err)

}

func (f *FileRecordedTestsSuite) TestFileCreateNonDefaultHTTPHeaders() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))

	httpHeaders := file.HTTPHeaders{
		ContentType:        to.Ptr("my_type"),
		ContentDisposition: to.Ptr("my_disposition"),
		CacheControl:       to.Ptr("control"),
		ContentMD5:         nil,
		ContentLanguage:    to.Ptr("my_language"),
		ContentEncoding:    to.Ptr("my_encoding"),
	}

	_, err = fClient.Create(context.Background(), 1024, &file.CreateOptions{
		HTTPHeaders: &httpHeaders,
	})
	_require.NoError(err)

	resp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp.ContentType, httpHeaders.ContentType)
	_require.EqualValues(resp.ContentDisposition, httpHeaders.ContentDisposition)
	_require.EqualValues(resp.CacheControl, httpHeaders.CacheControl)
	_require.EqualValues(resp.ContentLanguage, httpHeaders.ContentLanguage)
	_require.EqualValues(resp.ContentEncoding, httpHeaders.ContentEncoding)
	_require.Nil(resp.ContentMD5)
}

func (f *FileRecordedTestsSuite) TestFileCreateNegativeMetadataInvalid() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 1024, &file.CreateOptions{
		Metadata:    map[string]*string{"!@#$%^&*()": to.Ptr("!@#$%^&*()")},
		HTTPHeaders: &file.HTTPHeaders{},
	})
	_require.Error(err)
}

func (f *FileRecordedTestsSuite) TestCreateFileNFS() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)
	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	owner := "345"
	group := "123"
	fileMode := "7777"

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)
	_require.NoError(err)
	fClient := premiumShareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	resp, err := fClient.Create(context.Background(), 1024, &file.CreateOptions{
		Owner:    to.Ptr(owner),
		Group:    to.Ptr(group),
		FileMode: to.Ptr(fileMode),
	})
	_require.NoError(err)
	_require.Equal(*resp.FileMode, fileMode)
	_require.Equal(*resp.Group, group)
	_require.Equal(*resp.Owner, owner)
	_require.Equal(*resp.NfsFileType, file.NfsFileType("Regular"))
}

func (f *FileUnrecordedTestsSuite) TestFileGetSetPropertiesNonDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	creationTime := time.Now().Add(-time.Hour)
	lastWriteTime := time.Now().Add(-time.Minute * 15)
	changeTime := time.Now().Add(-time.Minute * 30)

	options := &file.SetHTTPHeadersOptions{
		Permissions: &file.Permissions{Permission: &testcommon.SampleSDDL},
		SMBProperties: &file.SMBProperties{
			Attributes:    &file.NTFSFileAttributes{Hidden: true},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
			ChangeTime:    &changeTime,
		},
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
	}
	setResp, err := fClient.SetHTTPHeaders(context.Background(), options)
	_require.NoError(err)
	_require.NotNil(setResp.ETag)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotNil(setResp.RequestID)
	_require.NotNil(setResp.Version)
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)

	fileAttributes, err := file.ParseNTFSFileAttributes(setResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.Hidden)

	getResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.Equal(*getResp.FileType, "File")

	_require.EqualValues(getResp.ContentType, options.HTTPHeaders.ContentType)
	_require.EqualValues(getResp.ContentEncoding, options.HTTPHeaders.ContentEncoding)
	_require.EqualValues(getResp.ContentLanguage, options.HTTPHeaders.ContentLanguage)
	_require.EqualValues(getResp.ContentMD5, options.HTTPHeaders.ContentMD5)
	_require.EqualValues(getResp.CacheControl, options.HTTPHeaders.CacheControl)
	_require.EqualValues(getResp.ContentDisposition, options.HTTPHeaders.ContentDisposition)
	_require.Equal(*getResp.ContentLength, int64(0))
	// We'll just ensure a permission exists, no need to test overlapping functionality.
	_require.NotEqual(getResp.FilePermissionKey, "")
	_require.Equal(*getResp.FileAttributes, options.SMBProperties.Attributes.String())

	fileAttributes2, err := file.ParseNTFSFileAttributes(getResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.Hidden)
	_require.EqualValues(fileAttributes, fileAttributes2)

	_require.EqualValues(getResp.FileCreationTime.Format(testcommon.ISO8601), creationTime.UTC().Format(testcommon.ISO8601))
	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
}

func (f *FileRecordedTestsSuite) TestFileGetSetPropertiesDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 0, shareClient)

	setResp, err := fClient.SetHTTPHeaders(context.Background(), nil)
	_require.NoError(err)
	_require.NotEqual(*setResp.ETag, "")
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotEqual(setResp.RequestID, "")
	_require.NotEqual(setResp.Version, "")
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)

	metadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}
	_, err = fClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: metadata,
	})
	_require.NoError(err)

	getResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.Equal(*getResp.FileType, "File")

	_require.Nil(getResp.ContentType)
	_require.Nil(getResp.ContentEncoding)
	_require.Nil(getResp.ContentLanguage)
	_require.Nil(getResp.ContentMD5)
	_require.Nil(getResp.CacheControl)
	_require.Nil(getResp.ContentDisposition)
	_require.Equal(*getResp.ContentLength, int64(0))

	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
	_require.EqualValues(getResp.Metadata, metadata)
}

func (f *FileUnrecordedTestsSuite) TestFileSetHTTPHeaders() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, &file.CreateOptions{
		FilePermissionFormat: (*file.PermissionFormat)(to.Ptr(testcommon.FilePermissionFormatSddl)),
		Permissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	creationTime := time.Now().Add(-time.Hour)
	lastWriteTime := time.Now().Add(-time.Minute * 15)
	changeTime := time.Now().Add(-time.Minute * 30)

	options := &file.SetHTTPHeadersOptions{
		FilePermissionFormat: (*file.PermissionFormat)(to.Ptr(testcommon.FilePermissionBinary)),
		Permissions:          &file.Permissions{Permission: &testcommon.SampleBinary},
		SMBProperties: &file.SMBProperties{
			Attributes:    &file.NTFSFileAttributes{Hidden: true},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
			ChangeTime:    &changeTime,
		},
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
	}
	setResp, err := fClient.SetHTTPHeaders(context.Background(), options)
	_require.NoError(err)
	_require.NotNil(setResp.ETag)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotNil(setResp.RequestID)
	_require.NotNil(setResp.Version)
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)
}

func (f *FileRecordedTestsSuite) TestFileSetHTTPHeadersNfs() {
	_require := require.New(f.T())
	testName := f.T().Name()
	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)
	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	owner := "345"
	group := "123"
	fileMode := "7777"

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)
	_require.NoError(err)
	fClient := premiumShareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	opts := &file.SetHTTPHeadersOptions{
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
		Owner:    to.Ptr(owner),
		Group:    to.Ptr(group),
		FileMode: to.Ptr(fileMode),
	}
	setResp, err := fClient.SetHTTPHeaders(context.Background(), opts)
	_require.NoError(err)
	_require.NotNil(setResp.ETag)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotNil(setResp.RequestID)
	_require.NotNil(setResp.Version)
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)
	_require.NotNil(setResp.LinkCount)
	_require.Equal(*setResp.FileMode, fileMode)
	_require.Equal(*setResp.Group, group)
	_require.Equal(*setResp.Owner, owner)

	getResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getResp.LinkCount)
	_require.NotNil(getResp.FileType)
	_require.Equal(*getResp.FileMode, fileMode)
	_require.Equal(*getResp.Group, group)
	_require.Equal(*getResp.Owner, owner)
}

func (f *FileRecordedTestsSuite) TestFilePreservePermissions() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, &file.CreateOptions{
		Permissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)

	// Grab the original perm key before we set file headers.
	getResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	pKey := getResp.FilePermissionKey
	cTime := getResp.FileCreationTime
	lwTime := getResp.FileLastWriteTime
	changeTime := getResp.FileChangeTime
	attribs := getResp.FileAttributes

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	properties := file.SetHTTPHeadersOptions{
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
		// SMBProperties, when options are left nil, leads to preserving.
		SMBProperties: &file.SMBProperties{},
	}

	setResp, err := fClient.SetHTTPHeaders(context.Background(), &properties)
	_require.NoError(err)
	_require.NotNil(setResp.ETag)
	_require.NotNil(setResp.RequestID)
	_require.NotNil(setResp.LastModified)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotNil(setResp.Version)
	_require.Equal(setResp.Date.IsZero(), false)

	getResp, err = fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(setResp.LastModified)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.Equal(*getResp.FileType, "File")

	_require.EqualValues(getResp.ContentType, properties.HTTPHeaders.ContentType)
	_require.EqualValues(getResp.ContentEncoding, properties.HTTPHeaders.ContentEncoding)
	_require.EqualValues(getResp.ContentLanguage, properties.HTTPHeaders.ContentLanguage)
	_require.EqualValues(getResp.ContentMD5, properties.HTTPHeaders.ContentMD5)
	_require.EqualValues(getResp.CacheControl, properties.HTTPHeaders.CacheControl)
	_require.EqualValues(getResp.ContentDisposition, properties.HTTPHeaders.ContentDisposition)
	_require.Equal(*getResp.ContentLength, int64(0))
	// Ensure that the permission key gets preserved
	_require.EqualValues(getResp.FilePermissionKey, pKey)
	_require.EqualValues(cTime, getResp.FileCreationTime)
	_require.EqualValues(lwTime, getResp.FileLastWriteTime)
	_require.NotEqualValues(changeTime, getResp.FileChangeTime) // default value is "now" for file change time
	_require.EqualValues(attribs, getResp.FileAttributes)

	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
}

func (f *FileRecordedTestsSuite) TestFileGetSetPropertiesSnapshot() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer func() {
		_, err := shareClient.Delete(context.Background(), &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})
		_require.NoError(err)
	}()

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	fileSetHTTPHeadersOptions := file.SetHTTPHeadersOptions{
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
	}
	setResp, err := fClient.SetHTTPHeaders(context.Background(), &fileSetHTTPHeadersOptions)
	_require.NoError(err)
	_require.NotEqual(*setResp.ETag, "")
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotEqual(setResp.RequestID, "")
	_require.NotEqual(setResp.Version, "")
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)

	metadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}
	_, err = fClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: metadata,
	})
	_require.NoError(err)

	resp, err := shareClient.CreateSnapshot(context.Background(), &share.CreateSnapshotOptions{Metadata: map[string]*string{}})
	_require.NoError(err)
	_require.NotNil(resp.Snapshot)

	// get properties on the share snapshot
	getResp, err := fClient.GetProperties(context.Background(), &file.GetPropertiesOptions{
		ShareSnapshot: resp.Snapshot,
	})
	_require.NoError(err)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.Equal(*getResp.FileType, "File")

	_require.EqualValues(getResp.ContentType, fileSetHTTPHeadersOptions.HTTPHeaders.ContentType)
	_require.EqualValues(getResp.ContentEncoding, fileSetHTTPHeadersOptions.HTTPHeaders.ContentEncoding)
	_require.EqualValues(getResp.ContentLanguage, fileSetHTTPHeadersOptions.HTTPHeaders.ContentLanguage)
	_require.EqualValues(getResp.ContentMD5, fileSetHTTPHeadersOptions.HTTPHeaders.ContentMD5)
	_require.EqualValues(getResp.CacheControl, fileSetHTTPHeadersOptions.HTTPHeaders.CacheControl)
	_require.EqualValues(getResp.ContentDisposition, fileSetHTTPHeadersOptions.HTTPHeaders.ContentDisposition)
	_require.Equal(*getResp.ContentLength, int64(0))

	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
	_require.EqualValues(getResp.Metadata, metadata)
}

func (f *FileRecordedTestsSuite) TestGetSetMetadataNonDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	metadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}
	setResp, err := fClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: metadata,
	})
	_require.NoError(err)
	_require.NotNil(setResp.ETag)
	_require.NotNil(setResp.RequestID)
	_require.NotNil(setResp.Version)
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)

	getResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
	_require.EqualValues(getResp.Metadata, metadata)
}

func (f *FileRecordedTestsSuite) TestFileSetMetadataNil() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	md := map[string]*string{"Not": to.Ptr("nil")}

	_, err = fClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: md,
	})
	_require.NoError(err)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp1.Metadata, md)

	_, err = fClient.SetMetadata(context.Background(), nil)
	_require.NoError(err)

	resp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp2.Metadata, 0)
}

func (f *FileRecordedTestsSuite) TestFileSetMetadataDefaultEmpty() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	md := map[string]*string{"Not": to.Ptr("nil")}

	_, err = fClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: md,
	})
	_require.NoError(err)

	resp1, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp1.Metadata, md)

	_, err = fClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: map[string]*string{},
	})
	_require.NoError(err)

	resp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp2.Metadata, 0)
}

func (f *FileRecordedTestsSuite) TestFileSetMetadataInvalidField() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	_, err = fClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: map[string]*string{"!@#$%^&*()": to.Ptr("!@#$%^&*()")},
	})
	_require.Error(err)
}

func (f *FileUnrecordedTestsSuite) TestFileDelete() {
	if recording.GetRecordMode() == recording.LiveMode {
		f.T().Skip("This test cannot be made live")
	}
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)

	response, err := shareClient.Delete(context.Background(), &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})
	_require.NoError(err)
	_require.NotNil(response.FileShareUsageBytes)
	_require.NotNil(response.FileShareSnapshotUsageBytes)
}

func (f *FileRecordedTestsSuite) TestStartCopyDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcFile := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	destFile := shareClient.NewRootDirectoryClient().NewFileClient("dest" + testcommon.GenerateFileName(testName))

	fileSize := int64(2048)
	_, err = srcFile.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	contentR, srcContent := testcommon.GenerateData(int(fileSize))
	srcContentMD5 := md5.Sum(srcContent)

	_, err = srcFile.UploadRange(context.Background(), 0, contentR, nil)
	_require.NoError(err)

	copyResp, err := destFile.StartCopyFromURL(context.Background(), srcFile.URL(), nil)
	_require.NoError(err)
	_require.NotNil(copyResp.ETag)
	_require.Equal(copyResp.LastModified.IsZero(), false)
	_require.NotNil(copyResp.RequestID)
	_require.NotNil(copyResp.Version)
	_require.Equal(copyResp.Date.IsZero(), false)
	_require.NotEqual(copyResp.CopyStatus, "")

	time.Sleep(time.Duration(5) * time.Second)

	getResp, err := destFile.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(getResp.CopyID, copyResp.CopyID)
	_require.NotEqual(*getResp.CopyStatus, "")
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.Equal(*getResp.CopySource, srcFile.URL())
	}
	_require.Equal(*getResp.CopyStatus, file.CopyStatusTypeSuccess)

	// Abort will fail after copy finished
	_, err = destFile.AbortCopy(context.Background(), *copyResp.CopyID, nil)
	_require.Error(err)
	testcommon.ValidateHTTPErrorCode(_require, err, http.StatusConflict)

	// validate data copied
	dResp, err := destFile.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range:              file.HTTPRange{Offset: 0, Count: fileSize},
		RangeGetContentMD5: to.Ptr(true),
	})
	_require.NoError(err)

	destContent, err := io.ReadAll(dResp.Body)
	_require.NoError(err)
	_require.EqualValues(srcContent, destContent)
	_require.Equal(dResp.ContentMD5, srcContentMD5[:])

	fileAttributes, err := file.ParseNTFSFileAttributes(dResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyDestEmpty() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShareWithData(context.Background(), _require, "src"+testcommon.GenerateFileName(testName), shareClient)
	copyFClient := testcommon.GetFileClientFromShare("dest"+testcommon.GenerateFileName(testName), shareClient)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), nil)
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp, err := copyFClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	// Read the file data to verify the copy
	data, err := io.ReadAll(resp.Body)
	defer func() {
		err = resp.Body.Close()
		_require.NoError(err)
	}()

	_require.NoError(err)
	_require.Equal(*resp.ContentLength, int64(len(testcommon.FileDefaultData)))
	_require.Equal(string(data), testcommon.FileDefaultData)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyMetadata() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	basicMetadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}
	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{Metadata: basicMetadata})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp2.Metadata, basicMetadata)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyMetadataNil() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	basicMetadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err = copyFClient.Create(context.Background(), 0, &file.CreateOptions{Metadata: basicMetadata})
	_require.NoError(err)

	gResp, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(gResp.Metadata, basicMetadata)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), nil)
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp2.Metadata, 0)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyMetadataEmpty() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	basicMetadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err = copyFClient.Create(context.Background(), 0, &file.CreateOptions{Metadata: basicMetadata})
	_require.NoError(err)

	gResp, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(gResp.Metadata, basicMetadata)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{Metadata: map[string]*string{}})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp2.Metadata, 0)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyNegativeMetadataInvalidField() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		Metadata: map[string]*string{"!@#$%^&*()": to.Ptr("!@#$%^&*()")},
	})
	_require.Error(err)
}

func (f *FileRecordedTestsSuite) TestFileStartCopySourceCreationTime() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 21:00:00 GMT 2023")
	_require.NoError(err)

	cResp, err := fClient.Create(context.Background(), 0, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes:    &file.NTFSFileAttributes{ReadOnly: true, Hidden: true},
			CreationTime:  to.Ptr(currTime.Add(5 * time.Minute)),
			LastWriteTime: to.Ptr(currTime.Add(2 * time.Minute)),
			ChangeTime:    to.Ptr(currTime.Add(8 * time.Minute)),
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	fileAttributes, err := file.ParseNTFSFileAttributes(cResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.ReadOnly)
	_require.True(fileAttributes.Hidden)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			CreationTime: file.SourceCopyFileCreationTime{},
		},
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp2.FileCreationTime, cResp.FileCreationTime)
	_require.NotEqualValues(resp2.FileLastWriteTime, cResp.FileLastWriteTime)
	_require.NotEqualValues(resp2.FileChangeTime, cResp.FileChangeTime)
	_require.NotEqualValues(resp2.FileAttributes, cResp.FileAttributes)
}

func (f *FileRecordedTestsSuite) TestFileStartCopySourceProperties() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)

	cResp, err := fClient.Create(context.Background(), 0, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes:    &file.NTFSFileAttributes{System: true},
			CreationTime:  to.Ptr(currTime.Add(1 * time.Minute)),
			LastWriteTime: to.Ptr(currTime.Add(2 * time.Minute)),
			ChangeTime:    to.Ptr(currTime.Add(5 * time.Minute)),
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	fileAttributes, err := file.ParseNTFSFileAttributes(cResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.System)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			CreationTime:       file.SourceCopyFileCreationTime{},
			LastWriteTime:      file.SourceCopyFileLastWriteTime{},
			ChangeTime:         file.SourceCopyFileChangeTime{},
			Attributes:         file.SourceCopyFileAttributes{},
			PermissionCopyMode: to.Ptr(file.PermissionCopyModeTypeSource),
		},
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp2.FileCreationTime, cResp.FileCreationTime)
	_require.EqualValues(resp2.FileLastWriteTime, cResp.FileLastWriteTime)
	_require.EqualValues(resp2.FileChangeTime, cResp.FileChangeTime)
	_require.EqualValues(resp2.FileAttributes, cResp.FileAttributes)
	_require.EqualValues(resp2.FilePermissionKey, cResp.FilePermissionKey)

	fileAttributes2, err := file.ParseNTFSFileAttributes(resp2.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.System)
	_require.EqualValues(fileAttributes, fileAttributes2)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyDifferentProperties() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 20:00:00 GMT 2023")
	_require.NoError(err)

	cResp, err := fClient.Create(context.Background(), 0, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes:    &file.NTFSFileAttributes{System: true},
			CreationTime:  to.Ptr(currTime.Add(1 * time.Minute)),
			LastWriteTime: to.Ptr(currTime.Add(2 * time.Minute)),
			ChangeTime:    to.Ptr(currTime.Add(3 * time.Minute)),
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	destCreationTime := currTime.Add(5 * time.Minute)
	destLastWriteTIme := currTime.Add(6 * time.Minute)
	destChangeTime := currTime.Add(7 * time.Minute)
	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			CreationTime:  file.DestinationCopyFileCreationTime(destCreationTime),
			LastWriteTime: file.DestinationCopyFileLastWriteTime(destLastWriteTIme),
			ChangeTime:    file.DestinationCopyFileChangeTime(destChangeTime),
			Attributes:    file.DestinationCopyFileAttributes{ReadOnly: true},
		},
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotEqualValues(resp2.FileCreationTime, cResp.FileCreationTime)
	_require.EqualValues(*resp2.FileCreationTime, destCreationTime.UTC())
	_require.NotEqualValues(resp2.FileLastWriteTime, cResp.FileLastWriteTime)
	_require.EqualValues(*resp2.FileLastWriteTime, destLastWriteTIme.UTC())
	_require.NotEqualValues(resp2.FileChangeTime, cResp.FileChangeTime)
	_require.EqualValues(*resp2.FileChangeTime, destChangeTime.UTC())
	_require.NotEqualValues(resp2.FileAttributes, cResp.FileAttributes)
	_require.EqualValues(resp2.FilePermissionKey, cResp.FilePermissionKey)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyOverrideMode() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	cResp, err := fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		Permissions: &file.Permissions{
			Permission: to.Ptr(testcommon.SampleSDDL),
		},
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			PermissionCopyMode: to.Ptr(file.PermissionCopyModeTypeOverride),
		},
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotEqualValues(resp2.FileCreationTime, cResp.FileCreationTime)
	_require.NotEqualValues(resp2.FileLastWriteTime, cResp.FileLastWriteTime)
	_require.NotEqualValues(resp2.FileChangeTime, cResp.FileChangeTime)
	_require.NotEqualValues(resp2.FilePermissionKey, cResp.FilePermissionKey)
}

func (f *FileRecordedTestsSuite) TestNegativeFileStartCopyOverrideMode() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	cResp, err := fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	// permission or permission key is required when the PermissionCopyMode is override.
	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			PermissionCopyMode: to.Ptr(file.PermissionCopyModeTypeOverride),
		},
	})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.MissingRequiredHeader)
}

func (f *FileRecordedTestsSuite) TestFileStartCopySetArchiveAttributeTrue() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	cResp, err := fClient.Create(context.Background(), 0, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{ReadOnly: true, Hidden: true},
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			Attributes:          file.DestinationCopyFileAttributes{System: true, ReadOnly: true},
			SetArchiveAttribute: to.Ptr(true),
		},
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotEqualValues(resp2.FileCreationTime, cResp.FileCreationTime)
	_require.NotEqualValues(resp2.FileLastWriteTime, cResp.FileLastWriteTime)
	_require.NotEqualValues(resp2.FileChangeTime, cResp.FileChangeTime)
	_require.Contains(*resp2.FileAttributes, "Archive")
}

func (f *FileRecordedTestsSuite) TestFileStartCopySetArchiveAttributeFalse() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	cResp, err := fClient.Create(context.Background(), 0, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{ReadOnly: true, Hidden: true},
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			Attributes:          file.DestinationCopyFileAttributes{System: true, ReadOnly: true},
			SetArchiveAttribute: to.Ptr(false),
		},
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotEqualValues(resp2.FileCreationTime, cResp.FileCreationTime)
	_require.NotEqualValues(resp2.FileLastWriteTime, cResp.FileLastWriteTime)
	_require.NotEqualValues(resp2.FileChangeTime, cResp.FileChangeTime)
	_require.NotContains(*resp2.FileAttributes, "Archive")
}

func (f *FileRecordedTestsSuite) TestFileStartCopyDestReadOnly() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	cResp, err := fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	_, err = copyFClient.Create(context.Background(), 0, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{ReadOnly: true},
		},
	})
	_require.NoError(err)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			IgnoreReadOnly: to.Ptr(true),
		},
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp2, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotEqualValues(resp2.FileCreationTime, cResp.FileCreationTime)
	_require.NotEqualValues(resp2.FileLastWriteTime, cResp.FileLastWriteTime)
	_require.NotEqualValues(resp2.FileChangeTime, cResp.FileChangeTime)
}

func (f *FileRecordedTestsSuite) TestNegativeFileStartCopyDestReadOnly() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	cResp, err := fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileChangeTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	_, err = copyFClient.Create(context.Background(), 0, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{ReadOnly: true},
		},
	})
	_require.NoError(err)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ReadOnlyAttribute)
}

func (f *FileRecordedTestsSuite) TestFileStartCopySourceNonExistent() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), nil)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileUnrecordedTestsSuite) TestFileStartCopyUsingSASSrc() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, "src"+shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fClient := testcommon.CreateNewFileFromShareWithData(context.Background(), _require, "src"+fileName, shareClient)

	fileURLWithSAS, err := fClient.GetSASURL(sas.FilePermissions{Read: true, Write: true, Create: true, Delete: true}, time.Now().Add(5*time.Minute).UTC(), nil)
	_require.NoError(err)

	// Create a new share for the destination
	copyShareClient := testcommon.CreateNewShare(context.Background(), _require, "dest"+shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, copyShareClient)

	copyFileClient := testcommon.GetFileClientFromShare("dst"+fileName, copyShareClient)

	_, err = copyFileClient.StartCopyFromURL(context.Background(), fileURLWithSAS, nil)
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	dResp, err := copyFileClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	defer func() {
		err = dResp.Body.Close()
		_require.NoError(err)
	}()

	_require.NoError(err)
	_require.Equal(*dResp.ContentLength, int64(len(testcommon.FileDefaultData)))
	_require.Equal(string(data), testcommon.FileDefaultData)
}

func (f *FileRecordedTestsSuite) TestFileStartCopyModeCopyModeNfs() {

	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)
	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	owner := "345"
	group := "123"
	mode := "6444"

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)
	_require.NoError(err)

	fClient := premiumShareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	copyFClient := premiumShareClient.NewRootDirectoryClient().NewFileClient("dst" + testcommon.GenerateFileName(testName))

	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		Owner:             to.Ptr(owner),
		Group:             to.Ptr(group),
		FileMode:          to.Ptr(mode),
		FileOwnerCopyMode: to.Ptr(file.OwnerCopyModeOverride),
		FileModeCopyMode:  to.Ptr(file.ModeCopyModeOverride),
	})
	_require.NoError(err)

	time.Sleep(4 * time.Second)

	resp, err := copyFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.Group, group)
	_require.Equal(*resp.Owner, owner)
}

func (f *FileRecordedTestsSuite) TestFileAbortCopyNoCopyStarted() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	copyFClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = copyFClient.AbortCopy(context.Background(), "copynotstarted", nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidQueryParameterValue)
}

func (f *FileRecordedTestsSuite) TestResizeFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 1234, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, int64(1234))

	_, err = fClient.Resize(context.Background(), 4096, nil)
	_require.NoError(err)

	gResp, err = fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, int64(4096))
}

func (f *FileRecordedTestsSuite) TestFileResizeZero() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 10, nil)
	_require.NoError(err)

	_, err = fClient.Resize(context.Background(), 0, nil)
	_require.NoError(err)

	resp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.ContentLength, int64(0))
}

func (f *FileRecordedTestsSuite) TestFileResizeInvalidSizeNegative() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 0, nil)
	_require.NoError(err)

	_, err = fClient.Resize(context.Background(), -4, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.OutOfRangeInput)
}

func (f *FileRecordedTestsSuite) TestNegativeFileSizeMoreThanShareQuota() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	var fileShareMaxQuota int32 = 1024                  // share size in GiB which is 1TiB
	var fileMaxAllowedSizeInBytes int64 = 4398046511104 // file size in bytes which is 4 TiB

	shareClient := testcommon.GetShareClient(testcommon.GenerateShareName(testName), svcClient)
	_, err = shareClient.Create(context.Background(), &share.CreateOptions{
		Quota: &fileShareMaxQuota,
	})
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileMaxAllowedSizeInBytes, &file.CreateOptions{
		HTTPHeaders: &file.HTTPHeaders{},
	})
	_require.Error(err)
}

func (f *FileRecordedTestsSuite) TestCreateMaximumSizeFileShare() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	var fileShareMaxQuota int32 = 5120                  // share size in GiB which is 5TiB
	var fileMaxAllowedSizeInBytes int64 = 4398046511104 // file size in bytes which is 4 TiB

	shareClient := testcommon.GetShareClient(testcommon.GenerateShareName(testName), svcClient)
	_, err = shareClient.Create(context.Background(), &share.CreateOptions{
		Quota: &fileShareMaxQuota,
	})
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileMaxAllowedSizeInBytes, &file.CreateOptions{
		HTTPHeaders: &file.HTTPHeaders{},
	})
	_require.NoError(err)
}

func (f *FileRecordedTestsSuite) TestSASFileClientNoKey() {
	_require := require.New(f.T())
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	testName := f.T().Name()
	shareName := testcommon.GenerateShareName(testName)
	fileName := testcommon.GenerateFileName(testName)
	fileClient, err := file.NewClientWithNoCredential(fmt.Sprintf("https://%s.file.core.windows.net/%v/%v", accountName, shareName, fileName), nil)
	_require.NoError(err)

	permissions := sas.FilePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)

	_, err = fileClient.GetSASURL(permissions, expiry, nil)
	_require.Equal(err, fileerror.MissingSharedKeyCredential)
}

func (f *FileRecordedTestsSuite) TestSASFileClientSignNegative() {
	_require := require.New(f.T())
	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)
	_require.Greater(len(accountKey), 0)

	cred, err := file.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	testName := f.T().Name()
	shareName := testcommon.GenerateShareName(testName)
	fileName := testcommon.GenerateFileName(testName)
	fileClient, err := file.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.file.core.windows.net/%v%v", accountName, shareName, fileName), cred, nil)
	_require.NoError(err)

	permissions := sas.FilePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		Create: true,
	}
	expiry := time.Time{}

	// zero expiry time
	_, err = fileClient.GetSASURL(permissions, expiry, &file.GetSASURLOptions{StartTime: to.Ptr(time.Now())})
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")

	// zero start and expiry time
	_, err = fileClient.GetSASURL(permissions, expiry, &file.GetSASURLOptions{})
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")

	// empty permissions
	_, err = fileClient.GetSASURL(sas.FilePermissions{}, expiry, nil)
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")
}

func (f *FileRecordedTestsSuite) TestFileUploadClearListRange() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 10
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

	contentSize := 1024 * 2 // 2KB
	contentR, contentD := testcommon.GenerateData(contentSize)
	md5Value := md5.Sum(contentD)
	contentMD5 := md5Value[:]

	uResp, err := fClient.UploadRange(context.Background(), 0, contentR, &file.UploadRangeOptions{
		TransactionalValidation: file.TransferValidationTypeMD5(contentMD5),
	})
	_require.NoError(err)
	_require.NotNil(uResp.ContentMD5)
	_require.EqualValues(uResp.ContentMD5, contentMD5)

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(int64(contentSize - 1))})

	cResp, err := fClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: int64(contentSize)}, nil)
	_require.NoError(err)
	_require.Nil(cResp.ContentMD5)

	rangeList2, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList2.Ranges, 0)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURL() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 20
	srcFileName := "src" + testcommon.GenerateFileName(testName)
	srcFClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	perms := sas.FilePermissions{Read: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ShareName:   shareName,
		FilePath:    srcFileName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	srcFileSAS := srcFClient.URL() + "?" + sasQueryParams.Encode()

	destFClient := shareClient.NewRootDirectoryClient().NewFileClient("dest" + testcommon.GenerateFileName(testName))
	_, err = destFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)
	uResp, err := destFClient.UploadRangeFromURL(context.Background(), srcFileSAS, 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
	})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(uResp.XMSContentCRC64), contentCRC64)

	rangeList, err := destFClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, int64(contentSize-1))

	cResp, err := destFClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: int64(contentSize)}, nil)
	_require.NoError(err)
	_require.Nil(cResp.ContentMD5)

	rangeList2, err := destFClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList2.Ranges, 0)
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeFromURLNegative() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 20
	srcFileName := "src" + testcommon.GenerateFileName(testName)
	srcFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, srcFileName, fileSize, shareClient)

	gResp, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

	contentSize := 1024 * 8 // 8KB
	rsc, content := testcommon.GenerateData(contentSize)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	destFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "dest"+testcommon.GenerateFileName(testName), fileSize, shareClient)

	_, err = destFClient.UploadRangeFromURL(context.Background(), srcFClient.URL(), 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
	})
	_require.Error(err)
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeFromURLOffsetNegative() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 20
	srcFileName := "src" + testcommon.GenerateFileName(testName)
	srcFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, srcFileName, fileSize, shareClient)

	gResp, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

	contentSize := 1024 * 8 // 8KB
	destFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "dest"+testcommon.GenerateFileName(testName), fileSize, shareClient)

	// error is returned when source offset is negative
	_, err = destFClient.UploadRangeFromURL(context.Background(), srcFClient.URL(), -1, 0, int64(contentSize), nil)
	_require.Error(err)
	_require.Equal(err.Error(), "invalid argument: source and destination offsets must be >= 0")
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURLCopySourceAuthBlob() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	// Getting token
	accessToken, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 10
	contentSize := 1024 * 8 // 8KB
	_, content := testcommon.GenerateData(contentSize)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	// create source block blob
	blobClient, err := azblob.NewClient("https://"+accountName+".blob.core.windows.net/", cred, nil)
	_require.NoError(err)

	containerName := "goc" + testcommon.GenerateEntityName(testName)
	blobName := "blob" + testcommon.GenerateEntityName(testName)
	_, err = blobClient.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)
	defer func() {
		_, err := blobClient.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	_, err = blobClient.UploadBuffer(context.Background(), containerName, blobName, content, nil)
	_require.NoError(err)

	destFClient := shareClient.NewRootDirectoryClient().NewFileClient("dest" + testcommon.GenerateFileName(testName))
	_, err = destFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	blobURL := blobClient.ServiceClient().NewContainerClient(containerName).NewBlockBlobClient(blobName).URL()
	uResp, err := destFClient.UploadRangeFromURL(context.Background(), blobURL, 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
		CopySourceAuthorization: to.Ptr("Bearer " + accessToken.Token),
	})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(uResp.XMSContentCRC64), contentCRC64)

	// validate the content uploaded
	dResp, err := destFClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{Offset: 0, Count: int64(contentSize)},
	})
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	defer func() {
		err = dResp.Body.Close()
		_require.NoError(err)
	}()

	_require.EqualValues(data, content)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURLCopySourceAuthFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	// Getting token
	accessToken, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 10
	contentSize := 1024 * 8 // 8KB
	_, content := testcommon.GenerateData(contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	srcFClient := shareClient.NewRootDirectoryClient().NewFileClient("src" + testcommon.GenerateFileName(testName))
	_, err = srcFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	destFileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/dest" + testcommon.GenerateFileName(testName)
	destFClient, err := file.NewClient(destFileURL, cred, &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)})
	_require.NoError(err)

	_, err = destFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	uResp, err := destFClient.UploadRangeFromURL(context.Background(), srcFClient.URL(), 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
		CopySourceAuthorization: to.Ptr("Bearer " + accessToken.Token),
	})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(uResp.XMSContentCRC64), contentCRC64)

	// validate the content uploaded
	dResp, err := destFClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{Offset: 0, Count: int64(contentSize)},
	})
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	defer func() {
		err = dResp.Body.Close()
		_require.NoError(err)
	}()

	_require.EqualValues(data, content)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURLWithEmptyUploadRangeFromURLOptions() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 20
	srcFileName := "src" + testcommon.GenerateFileName(testName)
	srcFClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := srcFClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	perms := sas.FilePermissions{Read: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ShareName:   shareName,
		FilePath:    srcFileName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	srcFileSAS := srcFClient.URL() + "?" + sasQueryParams.Encode()

	destFClient := shareClient.NewRootDirectoryClient().NewFileClient("dest" + testcommon.GenerateFileName(testName))
	_, err = destFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)
	uResp, err := destFClient.UploadRangeFromURL(context.Background(), srcFileSAS, 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)

	rangeList, err := destFClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, int64(contentSize-1))

	cResp, err := destFClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: int64(contentSize)}, nil)
	_require.NoError(err)
	_require.Nil(cResp.ContentMD5)

	rangeList2, err := destFClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList2.Ranges, 0)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadBuffer() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 100 * 1024 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 200 * 1024 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadStream() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 100 * 1024 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileUnrecordedTestsSuite) TestFileDownloadBuffer() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 100 * 1024 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileUnrecordedTestsSuite) TestFileDownloadFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 100 * 1024 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileRecordedTestsSuite) TestUploadDownloadDefaultNonDefaultMD5() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "src"+testcommon.GenerateFileName(testName), 2048, shareClient)
	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	contentR, contentD := testcommon.GenerateData(2048)

	pResp, err := fClient.UploadRange(context.Background(), 0, contentR, nil)
	_require.NoError(err)
	_require.NotNil(pResp.ContentMD5)
	_require.NotNil(pResp.IsServerEncrypted)
	_require.NotNil(pResp.ETag)
	_require.Equal(pResp.LastModified.IsZero(), false)
	_require.NotNil(pResp.RequestID)
	_require.NotNil(pResp.Version)
	_require.Equal(pResp.Date.IsZero(), false)

	// Get with rangeGetContentMD5 enabled.
	// Partial data, check status code 206.
	resp, err := fClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range:              file.HTTPRange{Offset: 0, Count: 1024},
		RangeGetContentMD5: to.Ptr(true),
	})
	_require.NoError(err)
	_require.Equal(*resp.ContentLength, int64(1024))
	_require.NotNil(resp.ContentMD5)
	_require.Equal(*resp.ContentType, "application/octet-stream")

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.EqualValues(downloadedData, contentD[:1024])

	// Set ContentMD5 for the entire file.
	_, err = fClient.SetHTTPHeaders(context.Background(), &file.SetHTTPHeadersOptions{
		HTTPHeaders: &file.HTTPHeaders{
			ContentMD5:      pResp.ContentMD5,
			ContentLanguage: to.Ptr("test")},
	})
	_require.NoError(err)

	// Test get with another type of range index, and validate if FileContentMD5 can be got correct.
	resp, err = fClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{Offset: 1024, Count: file.CountToEnd},
	})
	_require.NoError(err)
	_require.Equal(*resp.ContentLength, int64(1024))
	_require.Nil(resp.ContentMD5)
	_require.EqualValues(resp.FileContentMD5, pResp.ContentMD5)
	_require.Equal(*resp.ContentLanguage, "test")
	// Note: when it's downloading range, range's MD5 is returned, when set rangeGetContentMD5=true, currently set it to false, so should be empty

	downloadedData, err = io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.EqualValues(downloadedData, contentD[1024:])

	_require.Equal(*resp.AcceptRanges, "bytes")
	_require.Nil(resp.CacheControl)
	_require.Nil(resp.ContentDisposition)
	_require.Nil(resp.ContentEncoding)
	_require.Equal(*resp.ContentRange, "bytes 1024-2047/2048")
	_require.Nil(resp.ContentType) // Note ContentType is set to empty during SetHTTPHeaders
	_require.Nil(resp.CopyID)
	_require.Nil(resp.CopyProgress)
	_require.Nil(resp.CopySource)
	_require.Nil(resp.CopyStatus)
	_require.Nil(resp.CopyStatusDescription)
	_require.Equal(resp.Date.IsZero(), false)
	_require.NotEqual(*resp.ETag, "")
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Nil(resp.Metadata)
	_require.NotEqual(*resp.RequestID, "")
	_require.NotEqual(*resp.Version, "")
	_require.NotNil(resp.IsServerEncrypted)

	// Get entire fClient, check status code 200.
	resp, err = fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.ContentLength, int64(2048))
	_require.EqualValues(resp.ContentMD5, pResp.ContentMD5) // Note: This case is inted to get entire fClient, entire file's MD5 will be returned.
	_require.Nil(resp.FileContentMD5)                       // Note: FileContentMD5 is returned, only when range is specified explicitly.

	downloadedData, err = io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.EqualValues(downloadedData, contentD)

	_require.Equal(*resp.AcceptRanges, "bytes")
	_require.Nil(resp.CacheControl)
	_require.Nil(resp.ContentDisposition)
	_require.Nil(resp.ContentEncoding)
	_require.Nil(resp.ContentRange) // Note: ContentRange is returned, only when range is specified explicitly.
	_require.Nil(resp.ContentType)
	_require.Nil(resp.CopyCompletionTime)
	_require.Nil(resp.CopyID)
	_require.Nil(resp.CopyProgress)
	_require.Nil(resp.CopySource)
	_require.Nil(resp.CopyStatus)
	_require.Nil(resp.CopyStatusDescription)
	_require.Equal(resp.Date.IsZero(), false)
	_require.NotEqual(*resp.ETag, "")
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Nil(resp.Metadata)
	_require.NotEqual(*resp.RequestID, "")
	_require.NotEqual(*resp.Version, "")
	_require.NotNil(resp.IsServerEncrypted)
}

func (f *FileRecordedTestsSuite) TestFileDownloadDataNonExistentFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.GetFileClientFromShare(testcommon.GenerateFileName(testName), shareClient)

	_, err = fClient.DownloadStream(context.Background(), nil)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileDownloadDataOffsetOutOfRange() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 0, shareClient)

	_, err = fClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{
			Offset: int64(len(testcommon.FileDefaultData)),
			Count:  file.CountToEnd,
		},
	})
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidRange)
}

func (f *FileRecordedTestsSuite) TestFileDownloadDataEntireFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShareWithData(context.Background(), _require, testcommon.GenerateFileName(testName), shareClient)

	resp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	// Specifying a count of 0 results in the value being ignored
	data, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.EqualValues(string(data), testcommon.FileDefaultData)
}

func (f *FileRecordedTestsSuite) TestFileDownloadDataCountExact() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShareWithData(context.Background(), _require, testcommon.GenerateFileName(testName), shareClient)

	resp, err := fClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{
			Offset: 0,
			Count:  int64(len(testcommon.FileDefaultData)),
		},
	})
	_require.NoError(err)

	data, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.EqualValues(string(data), testcommon.FileDefaultData)
}

func (f *FileRecordedTestsSuite) TestFileDownloadDataCountOutOfRange() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShareWithData(context.Background(), _require, testcommon.GenerateFileName(testName), shareClient)

	resp, err := fClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{
			Offset: 0,
			Count:  int64(len(testcommon.FileDefaultData)) * 2,
		},
	})
	_require.NoError(err)

	data, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.EqualValues(string(data), testcommon.FileDefaultData)
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeNilBody() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "src"+testcommon.GenerateFileName(testName), 0, shareClient)

	_, err = fClient.UploadRange(context.Background(), 0, nil, nil)
	_require.Error(err)
	_require.Contains(err.Error(), "body must not be nil")
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeEmptyBody() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 0, shareClient)

	_, err = fClient.UploadRange(context.Background(), 0, streaming.NopCloser(bytes.NewReader([]byte{})), nil)
	_require.Error(err)
	_require.Contains(err.Error(), "body must contain readable data whose size is > 0")
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeNonExistentFile() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.GetFileClientFromShare(testcommon.GenerateFileName(testName), shareClient)

	rsc, _ := testcommon.GenerateData(12)
	_, err = fClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeTransactionalMD5() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClient)

	contentR, contentD := testcommon.GenerateData(2048)
	_md5 := md5.Sum(contentD)

	// Upload range with correct transactional MD5
	pResp, err := fClient.UploadRange(context.Background(), 0, contentR, &file.UploadRangeOptions{
		TransactionalValidation: file.TransferValidationTypeMD5(_md5[:]),
	})
	_require.NoError(err)
	_require.NotNil(pResp.ContentMD5)
	_require.NotNil(pResp.ETag)
	_require.Equal(pResp.LastModified.IsZero(), false)
	_require.NotNil(pResp.RequestID)
	_require.NotNil(pResp.Version)
	_require.Equal(pResp.Date.IsZero(), false)
	_require.EqualValues(pResp.ContentMD5, _md5[:])

	// Upload range with empty MD5, nil MD5 is covered by other cases.
	pResp, err = fClient.UploadRange(context.Background(), 1024, streaming.NopCloser(bytes.NewReader(contentD[1024:])), nil)
	_require.NoError(err)
	_require.NotNil(pResp.ContentMD5)

	resp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.ContentLength, int64(2048))

	downloadedData, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.EqualValues(downloadedData, contentD)
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeIncorrectTransactionalMD5() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClient)

	contentR, _ := testcommon.GenerateData(2048)
	_, incorrectMD5 := testcommon.GenerateData(16)

	// Upload range with incorrect transactional MD5
	_, err = fClient.UploadRange(context.Background(), 0, contentR, &file.UploadRangeOptions{
		TransactionalValidation: file.TransferValidationTypeMD5(incorrectMD5),
	})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.MD5Mismatch)
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeLastWrittenModePreserve() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.GetFileClientFromShare(testcommon.GenerateFileName(testName), shareClient)
	cResp, err := fClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileLastWriteTime)

	contentR, _ := testcommon.GenerateData(2048)

	// Upload range with correct transactional MD5
	pResp, err := fClient.UploadRange(context.Background(), 0, contentR, &file.UploadRangeOptions{
		LastWrittenMode: to.Ptr(file.LastWrittenModePreserve),
	})
	_require.NoError(err)
	_require.NotNil(pResp.FileLastWriteTime)
	_require.EqualValues(*pResp.FileLastWriteTime, *cResp.FileLastWriteTime)
}

func (f *FileRecordedTestsSuite) TestFileUploadRangeLastWrittenModeNow() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.GetFileClientFromShare(testcommon.GenerateFileName(testName), shareClient)
	cResp, err := fClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileLastWriteTime)

	contentR, _ := testcommon.GenerateData(2048)

	// Upload range with correct transactional MD5
	pResp, err := fClient.UploadRange(context.Background(), 0, contentR, &file.UploadRangeOptions{
		LastWrittenMode: to.Ptr(file.LastWrittenModeNow),
	})
	_require.NoError(err)
	_require.NotNil(pResp.FileLastWriteTime)
	_require.NotEqualValues(*pResp.FileLastWriteTime, *cResp.FileLastWriteTime)
}

// Testings for GetRangeList and ClearRange
func (f *FileRecordedTestsSuite) TestGetRangeListNonDefaultExact() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.GetFileClientFromShare(testcommon.GenerateFileName(testName), shareClient)

	fileSize := int64(5 * 1024)
	_, err = fClient.Create(context.Background(), fileSize, &file.CreateOptions{HTTPHeaders: &file.HTTPHeaders{}})
	_require.NoError(err)
	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	rsc, _ := testcommon.GenerateData(1024)
	putResp, err := fClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.NotNil(putResp.ContentMD5)
	_require.NotNil(putResp.RequestID)
	_require.NotNil(putResp.Version)
	_require.Equal(putResp.Date.IsZero(), false)

	rangeList, err := fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range: file.HTTPRange{
			Offset: 0,
			Count:  fileSize,
		},
	})
	_require.NoError(err)
	_require.Equal(rangeList.LastModified.IsZero(), false)
	_require.NotNil(rangeList.ETag)
	_require.Equal(*rangeList.FileContentLength, fileSize)
	_require.NotNil(rangeList.RequestID)
	_require.NotNil(rangeList.Version)
	_require.Equal(rangeList.Date.IsZero(), false)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, int64(1023))
}

// Default means clear the entire file's range
func (f *FileRecordedTestsSuite) TestClearRangeDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClient)
	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	rsc, _ := testcommon.GenerateData(2048)
	_, err = fClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	_, err = fClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: 2048}, nil)
	_require.NoError(err)

	rangeList, err := fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range: file.HTTPRange{Offset: 0, Count: file.CountToEnd},
	})
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 0)
}

func (f *FileRecordedTestsSuite) TestClearRangeNonDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 4096, shareClient)
	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	rsc, _ := testcommon.GenerateData(2048)
	_, err = fClient.UploadRange(context.Background(), 2048, rsc, nil)
	_require.NoError(err)

	_, err = fClient.ClearRange(context.Background(), file.HTTPRange{Offset: 2048, Count: 2048}, nil)
	_require.NoError(err)

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 0)
}

func (f *FileRecordedTestsSuite) TestClearRangeMultipleRanges() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClient)
	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	rsc, _ := testcommon.GenerateData(2048)
	_, err = fClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	_, err = fClient.ClearRange(context.Background(), file.HTTPRange{Offset: 1024, Count: 1024}, nil)
	_require.NoError(err)

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(int64(1023))})
}

// When not 512 aligned, clear range will set 0 the non-512 aligned range, and will not eliminate the range.
func (f *FileRecordedTestsSuite) TestClearRangeNonDefaultCount() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 1, shareClient)
	defer testcommon.DeleteFile(context.Background(), _require, fClient)

	d := []byte{65}
	_, err = fClient.UploadRange(context.Background(), 0, streaming.NopCloser(bytes.NewReader(d)), nil)
	_require.NoError(err)

	rangeList, err := fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range: file.HTTPRange{Offset: 0, Count: file.CountToEnd},
	})
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(int64(0))})

	_, err = fClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: 1}, nil)
	_require.NoError(err)

	rangeList, err = fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range: file.HTTPRange{Offset: 0, Count: file.CountToEnd},
	})
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(int64(0))})

	dResp, err := fClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	_bytes, err := io.ReadAll(dResp.Body)
	_require.NoError(err)
	_require.EqualValues(_bytes, []byte{0})
}

func (f *FileRecordedTestsSuite) TestFileClearRangeNegativeInvalidCount() {
	_require := require.New(f.T())
	testName := f.T().Name()
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.GetShareClient(testcommon.GenerateShareName(testName), svcClient)
	fClient := testcommon.GetFileClientFromShare(testcommon.GenerateFileName(testName), shareClient)

	_, err = fClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: 0}, nil)
	_require.Error(err)
	_require.Contains(err.Error(), "invalid argument: either offset is < 0 or count <= 0")
}

func (f *FileRecordedTestsSuite) TestFileGetRangeListDefaultEmptyFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 0, shareClient)

	resp, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Ranges, 0)
}

func setupGetRangeListTest(_require *require.Assertions, testName string, fileSize int64, shareClient *share.Client) *file.Client {
	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), fileSize, shareClient)
	rsc, _ := testcommon.GenerateData(int(fileSize))
	_, err := fClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)
	return fClient
}

func (f *FileRecordedTestsSuite) TestFileGetRangeListDefaultRange() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileSize := int64(512)
	fClient := setupGetRangeListTest(_require, testName, fileSize, shareClient)

	resp, err := fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range: file.HTTPRange{Offset: 0, Count: file.CountToEnd},
	})
	_require.NoError(err)
	_require.Len(resp.Ranges, 1)
	_require.EqualValues(*resp.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(fileSize - 1)})
}

func (f *FileRecordedTestsSuite) TestFileGetRangeListNonContiguousRanges() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileSize := int64(512)
	fClient := setupGetRangeListTest(_require, testName, fileSize, shareClient)

	_, err = fClient.Resize(context.Background(), fileSize*3, nil)
	_require.NoError(err)

	rsc, _ := testcommon.GenerateData(int(fileSize))
	_, err = fClient.UploadRange(context.Background(), fileSize*2, rsc, nil)
	_require.NoError(err)

	resp, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Ranges, 2)
	_require.EqualValues(*resp.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(fileSize - 1)})
	_require.EqualValues(*resp.Ranges[1], file.ShareFileRange{Start: to.Ptr(fileSize * 2), End: to.Ptr((fileSize * 3) - 1)})
}

func (f *FileRecordedTestsSuite) TestFileGetRangeListNonContiguousRangesCountLess() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileSize := int64(512)
	fClient := setupGetRangeListTest(_require, testName, fileSize, shareClient)

	resp, err := fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range: file.HTTPRange{Offset: 0, Count: fileSize},
	})
	_require.NoError(err)
	_require.Len(resp.Ranges, 1)
	_require.EqualValues(int64(0), *(resp.Ranges[0].Start))
	_require.EqualValues(fileSize-1, *(resp.Ranges[0].End))
}

func (f *FileRecordedTestsSuite) TestFileGetRangeListNonContiguousRangesCountExceed() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileSize := int64(512)
	fClient := setupGetRangeListTest(_require, testName, fileSize, shareClient)

	resp, err := fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range: file.HTTPRange{Offset: 0, Count: fileSize + 1},
	})
	_require.NoError(err)
	_require.NoError(err)
	_require.Len(resp.Ranges, 1)
	_require.EqualValues(*resp.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(fileSize - 1)})
}

func (f *FileRecordedTestsSuite) TestFileGetRangeListSnapshot() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer func() {
		_, err := shareClient.Delete(context.Background(), &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})
		_require.NoError(err)
	}()

	fileSize := int64(512)
	fClient := setupGetRangeListTest(_require, testName, fileSize, shareClient)

	resp, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.Snapshot)

	resp2, err := fClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		Range:         file.HTTPRange{Offset: 0, Count: file.CountToEnd},
		ShareSnapshot: resp.Snapshot,
	})
	_require.NoError(err)
	_require.Len(resp2.Ranges, 1)
	_require.EqualValues(*resp2.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(fileSize - 1)})
}

func (f *FileRecordedTestsSuite) TestFileGetRangeListSupportRename() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer func() {
		_, err := shareClient.Delete(context.Background(), &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})
		_require.NoError(err)
	}()

	fileSize := int64(512)
	fClient := setupGetRangeListTest(_require, testName, fileSize, shareClient)

	snapshotResponse1, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(snapshotResponse1.Snapshot)

	rsc, _ := testcommon.GenerateData(int(fileSize))
	_, err = fClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	renameResponse, err := fClient.Rename(context.Background(), "file2", nil)
	_require.NoError(err)

	newFileClient := testcommon.GetFileClientFromShare("file2", shareClient)
	_require.NoError(err)

	snapshotResponse2, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(snapshotResponse2.Snapshot)

	_, err = newFileClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		PrevShareSnapshot: snapshotResponse1.Snapshot,
		ShareSnapshot:     snapshotResponse2.Snapshot,
		SupportRename:     to.Ptr(false),
	})
	_require.Error(err, "PreviousSnapshotNotFound")

	_, err = newFileClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		PrevShareSnapshot: snapshotResponse1.Snapshot,
		ShareSnapshot:     snapshotResponse2.Snapshot,
		SupportRename:     nil,
	})
	_require.Error(err, "PreviousSnapshotNotFound")

	resp, err := newFileClient.GetRangeList(context.Background(), &file.GetRangeListOptions{
		PrevShareSnapshot: snapshotResponse1.Snapshot,
		ShareSnapshot:     snapshotResponse2.Snapshot,
		SupportRename:     to.Ptr(true),
	})
	_require.NoError(err)
	_require.Len(resp.Ranges, 1)
	_require.EqualValues(*resp.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(fileSize - 1)})
	_require.Equal(resp.ETag, renameResponse.ETag)
}

func (f *FileRecordedTestsSuite) TestFileUploadDownloadSmallBuffer() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 10 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

	_, content := testcommon.GenerateData(int(fileSize))
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	err = fClient.UploadBuffer(context.Background(), content, &file.UploadBufferOptions{
		Concurrency: 5,
		ChunkSize:   2 * 1024,
	})
	_require.NoError(err)

	destBuffer := make([]byte, fileSize)
	cnt, err := fClient.DownloadBuffer(context.Background(), destBuffer, &file.DownloadBufferOptions{
		ChunkSize:   2 * 1024,
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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileRecordedTestsSuite) TestFileUploadDownloadSmallFile() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 10 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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
	defer func(destFile *os.File) {
		err = destFile.Close()
		_require.NoError(err)
	}(destFile)

	cnt, err := fClient.DownloadFile(context.Background(), destFile, &file.DownloadFileOptions{
		ChunkSize:   2 * 1024,
		Concurrency: 5,
	})
	_require.NoError(err)
	_require.Equal(cnt, fileSize)

	destHash := md5.New()
	_, err = io.Copy(destHash, destFile)
	_require.NoError(err)
	downloadedContentMD5 := destHash.Sum(nil)

	_require.EqualValues(downloadedContentMD5, contentMD5)

	gResp2, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp2.ContentLength, fileSize)

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileRecordedTestsSuite) TestFileUploadDownloadSmallStream() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 10 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileRecordedTestsSuite) TestFileUploadDownloadWithProgress() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 10 * 1024
	fClient := shareClient.NewRootDirectoryClient().NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	gResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.ContentLength, fileSize)

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

	rangeList, err := fClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.Equal(*rangeList.Ranges[0].Start, int64(0))
	_require.Equal(*rangeList.Ranges[0].End, fileSize-1)
}

func (f *FileRecordedTestsSuite) TestFileListHandlesDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClient)

	resp, err := fClient.ListHandles(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Handles, 0)
	_require.NotNil(resp.NextMarker)
	_require.Equal(*resp.NextMarker, "")
}

func (f *FileRecordedTestsSuite) TestFileListHandlesClientNameFieldCheck() {
	if recording.GetRecordMode() == recording.LiveMode {
		f.T().Skip("This test cannot be made live")
	}
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.GetShareClient(testcommon.GenerateShareName(testName), svcClient)

	fClient := testcommon.GetFileClientFromShare(testcommon.GenerateFileName(testName), shareClient)
	resp, err := fClient.ListHandles(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Handles, 1)
	_require.NotNil(*resp.Handles[0].ClientName)
}

func (f *FileRecordedTestsSuite) TestFileForceCloseHandlesDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClient)

	resp, err := fClient.ForceCloseHandles(context.Background(), "*", nil)
	_require.NoError(err)
	_require.EqualValues(*resp.NumberOfHandlesClosed, 0)
	_require.EqualValues(*resp.NumberOfHandlesFailedToClose, 0)
	_require.Nil(resp.Marker)
}

func (f *FileRecordedTestsSuite) TestFileCreateDeleteUsingOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClient, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	resp, err := fileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	gResp, err := fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.FileCreationTime)
	_require.NotNil(gResp.FileLastWriteTime)
	_require.NotNil(gResp.FilePermissionKey)
	_require.Equal(*gResp.ContentLength, int64(2048))

	dResp, err := fileClient.Delete(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(dResp.Date.IsZero(), false)
	_require.NotNil(dResp.RequestID)
	_require.NotNil(dResp.Version)

	_, err = fileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileGetSetPropertiesUsingOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	clOptions := &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(f.T(), &clOptions.ClientOptions)
	fClient, err := file.NewClient(fileURL, cred, clOptions)
	_require.NoError(err)

	_, err = fClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	options := &file.SetHTTPHeadersOptions{
		Permissions: &file.Permissions{Permission: &testcommon.SampleSDDL},
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{Hidden: true, ReadOnly: true},
		},
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
	}
	setResp, err := fClient.SetHTTPHeaders(context.Background(), options)
	_require.NoError(err)
	_require.NotNil(setResp.ETag)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotNil(setResp.RequestID)
	_require.NotNil(setResp.Version)
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)

	fileAttributes, err := file.ParseNTFSFileAttributes(setResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.Hidden)
	_require.True(fileAttributes.ReadOnly)

	getResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(getResp.LastModified.IsZero(), false)
	_require.Equal(*getResp.FileType, "File")

	_require.EqualValues(getResp.ContentType, options.HTTPHeaders.ContentType)
	_require.EqualValues(getResp.ContentEncoding, options.HTTPHeaders.ContentEncoding)
	_require.EqualValues(getResp.ContentLanguage, options.HTTPHeaders.ContentLanguage)
	_require.EqualValues(getResp.ContentMD5, options.HTTPHeaders.ContentMD5)
	_require.EqualValues(getResp.CacheControl, options.HTTPHeaders.CacheControl)
	_require.EqualValues(getResp.ContentDisposition, options.HTTPHeaders.ContentDisposition)
	_require.Equal(*getResp.ContentLength, int64(2048))
	_require.NotEqual(getResp.FilePermissionKey, "")

	fileAttributes2, err := file.ParseNTFSFileAttributes(getResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.Hidden)
	_require.True(fileAttributes2.ReadOnly)
	_require.EqualValues(fileAttributes, fileAttributes2)
	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
}

func (f *FileRecordedTestsSuite) TestFileSetMetadataUsingOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClient, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = fileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	metadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}
	_, err = fileClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: metadata,
	})
	_require.NoError(err)

	getResp, err := fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp.ContentLength, int64(2048))
	_require.EqualValues(getResp.Metadata, metadata)
}

func (f *FileRecordedTestsSuite) TestFileUploadClearListRangeUsingOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClient, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = fileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	contentSize := 1024 * 2 // 2KB
	contentR, contentD := testcommon.GenerateData(contentSize)
	md5Value := md5.Sum(contentD)
	contentMD5 := md5Value[:]

	uResp, err := fileClient.UploadRange(context.Background(), 0, contentR, &file.UploadRangeOptions{
		TransactionalValidation: file.TransferValidationTypeMD5(contentMD5),
	})
	_require.NoError(err)
	_require.NotNil(uResp.ContentMD5)
	_require.EqualValues(uResp.ContentMD5, contentMD5)

	rangeList, err := fileClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(int64(contentSize - 1))})

	cResp, err := fileClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: int64(contentSize)}, nil)
	_require.NoError(err)
	_require.Nil(cResp.ContentMD5)

	rangeList2, err := fileClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList2.Ranges, 0)
}

func (f *FileRecordedTestsSuite) TestFileRenameDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcFileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, testcommon.GenerateFileName(testName), 2048, shareClient)

	resp, err := srcFileClient.Rename(context.Background(), "testFile", nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileUnrecordedTestsSuite) TestFileRenameUsingOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	srcFileClient, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = srcFileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	resp, err := srcFileClient.Rename(context.Background(), "testFile", nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURLUsingOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	sharedKeyCred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	tokenCred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 10
	contentSize := 1024 * 8 // 8KB
	_, content := testcommon.GenerateData(contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	srcFileName := "src" + testcommon.GenerateFileName(testName)
	srcFClient := shareClient.NewRootDirectoryClient().NewFileClient(srcFileName)
	_, err = srcFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	perms := sas.FilePermissions{Read: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ShareName:   shareName,
		FilePath:    srcFileName,
		Permissions: perms.String(),
	}.SignWithSharedKey(sharedKeyCred)
	_require.NoError(err)

	srcFileSAS := srcFClient.URL() + "?" + sasQueryParams.Encode()

	destFileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/dest" + testcommon.GenerateFileName(testName)
	destFClient, err := file.NewClient(destFileURL, tokenCred, &file.ClientOptions{FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup)})
	_require.NoError(err)

	_, err = destFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	uResp, err := destFClient.UploadRangeFromURL(context.Background(), srcFileSAS, 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
	})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(uResp.XMSContentCRC64), contentCRC64)

	// validate the content uploaded
	dResp, err := destFClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{Offset: 0, Count: int64(contentSize)},
	})
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	defer func() {
		err = dResp.Body.Close()
		_require.NoError(err)
	}()

	_require.EqualValues(data, content)
}

func (f *FileRecordedTestsSuite) TestFileRenameDifferentDir() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcDirCl := testcommon.CreateNewDirectory(context.Background(), _require, "dir1", shareClient)

	srcFileClient := srcDirCl.NewFileClient("file1")
	_, err = srcFileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	_ = testcommon.CreateNewDirectory(context.Background(), _require, "dir2", shareClient)

	_, err = srcFileClient.Rename(context.Background(), "dir2/file2/", nil)
	_require.NoError(err)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileRenameIgnoreReadOnly() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcFileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "file1", 2048, shareClient)

	_, err = shareClient.NewRootDirectoryClient().NewFileClient("file2").Create(context.Background(), 1024, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{ReadOnly: true},
		},
	})
	_require.NoError(err)

	_, err = srcFileClient.Rename(context.Background(), "file2", nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceAlreadyExists)

	_, err = srcFileClient.Rename(context.Background(), "file2", &file.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
	})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ReadOnlyAttribute)

	_, err = srcFileClient.Rename(context.Background(), "file2", &file.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
		IgnoreReadOnly:  to.Ptr(true),
	})
	_require.NoError(err)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileRenameNonDefault() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcFileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "file1", 2048, shareClient)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 21:00:00 GMT 2023")
	_require.NoError(err)
	creationTime := currTime.Add(5 * time.Minute).Round(time.Microsecond)
	lastWriteTime := currTime.Add(8 * time.Minute).Round(time.Millisecond)
	changeTime := currTime.Add(10 * time.Minute).Round(time.Millisecond)

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}

	renameOptions := file.RenameOptions{
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{
				ReadOnly: true,
				System:   true,
			},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
			ChangeTime:    &changeTime,
		},
		Permissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
		Metadata:    md,
		ContentType: to.Ptr("my_type"),
	}

	resp, err := srcFileClient.Rename(context.Background(), "file2", &renameOptions)
	_require.NoError(err)
	_require.NotNil(resp.FileCreationTime)
	_require.Equal(*resp.FileCreationTime, creationTime.UTC())
	_require.NotNil(resp.FileLastWriteTime)
	_require.Equal(*resp.FileLastWriteTime, lastWriteTime.UTC())
	_require.NotNil(resp.FileChangeTime)
	_require.Equal(*resp.FileChangeTime, changeTime.UTC())
	_require.NotNil(resp.FilePermissionKey)

	fileAttributes, err := file.ParseNTFSFileAttributes(resp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.ReadOnly)
	_require.True(fileAttributes.System)
}

func (f *FileRecordedTestsSuite) TestFileRenameSrcDestLease() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcFileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "file1", 2048, shareClient)
	destFileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "file2", 2048, shareClient)

	var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

	// acquire lease on source file
	srcFileLeaseClient, err := lease.NewFileClient(srcFileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)
	srcAcqResp, err := srcFileLeaseClient.Acquire(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(srcAcqResp.LeaseID)
	_require.Equal(*srcAcqResp.LeaseID, *proposedLeaseIDs[0])

	// acquire lease on destination file
	destFileLeaseClient, err := lease.NewFileClient(destFileClient, &lease.FileClientOptions{
		LeaseID: proposedLeaseIDs[1],
	})
	_require.NoError(err)
	destAcqResp, err := destFileLeaseClient.Acquire(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(destAcqResp.LeaseID)
	_require.Equal(*destAcqResp.LeaseID, *proposedLeaseIDs[1])

	_, err = srcFileClient.Rename(context.Background(), "file2", &file.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
	})
	_require.Error(err)

	_, err = srcFileClient.Rename(context.Background(), "file2", &file.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
		SourceLeaseAccessConditions: &file.SourceLeaseAccessConditions{
			SourceLeaseID: srcAcqResp.LeaseID,
		},
		DestinationLeaseAccessConditions: &file.DestinationLeaseAccessConditions{
			DestinationLeaseID: destAcqResp.LeaseID,
		},
	})
	_require.NoError(err)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileUnrecordedTestsSuite) TestFileRenameUsingSAS() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	perms := sas.FilePermissions{Read: true, Create: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ShareName:   shareName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	sasToken := sasQueryParams.Encode()

	srcFileClient, err := file.NewClientWithNoCredential(shareClient.URL()+"/file1?"+sasToken, nil)
	_require.NoError(err)

	_, err = srcFileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	destPathWithSAS := "file2?" + sasToken
	_, err = srcFileClient.Rename(context.Background(), destPathWithSAS, nil)
	_require.NoError(err)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileCreateDeleteTrailingDot() {
	_require := require.New(f.T())
	testName := f.T().Name()

	options := &service.ClientOptions{AllowTrailingDot: to.Ptr(true)}
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := "file.."
	fileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, fileName, 2048, shareClient)

	_, err = fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = fileClient.Delete(context.Background(), nil)
	_require.NoError(err)

	_, err = fileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileGetSetPropertiesTrailingDotOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := "file.."
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	clOptions := &file.ClientOptions{
		FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup),
		AllowTrailingDot:  to.Ptr(true),
	}
	testcommon.SetClientOptions(f.T(), &clOptions.ClientOptions)
	fClient, err := file.NewClient(fileURL, cred, clOptions)
	_require.NoError(err)

	_, err = fClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	md5Str := "MDAwMDAwMDA="
	testMd5 := []byte(md5Str)

	options := &file.SetHTTPHeadersOptions{
		Permissions: &file.Permissions{Permission: &testcommon.SampleSDDL},
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{Hidden: true, ReadOnly: true},
		},
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("en"),
			ContentMD5:         testMd5,
			CacheControl:       to.Ptr("no-transform"),
			ContentDisposition: to.Ptr("attachment"),
		},
	}
	setResp, err := fClient.SetHTTPHeaders(context.Background(), options)
	_require.NoError(err)
	_require.NotNil(setResp.ETag)
	_require.Equal(setResp.LastModified.IsZero(), false)
	_require.NotNil(setResp.RequestID)
	_require.NotNil(setResp.Version)
	_require.Equal(setResp.Date.IsZero(), false)
	_require.NotNil(setResp.IsServerEncrypted)

	fileAttributes, err := file.ParseNTFSFileAttributes(setResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.Hidden)
	_require.True(fileAttributes.ReadOnly)

	getResp, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(getResp.LastModified.IsZero(), false)
	_require.Equal(*getResp.FileType, "File")

	_require.EqualValues(getResp.ContentType, options.HTTPHeaders.ContentType)
	_require.EqualValues(getResp.ContentEncoding, options.HTTPHeaders.ContentEncoding)
	_require.EqualValues(getResp.ContentLanguage, options.HTTPHeaders.ContentLanguage)
	_require.EqualValues(getResp.ContentMD5, options.HTTPHeaders.ContentMD5)
	_require.EqualValues(getResp.CacheControl, options.HTTPHeaders.CacheControl)
	_require.EqualValues(getResp.ContentDisposition, options.HTTPHeaders.ContentDisposition)
	_require.Equal(*getResp.ContentLength, int64(2048))
	_require.NotEqual(getResp.FilePermissionKey, "")

	fileAttributes2, err := file.ParseNTFSFileAttributes(getResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.Hidden)
	_require.True(fileAttributes2.ReadOnly)
	_require.EqualValues(fileAttributes, fileAttributes2)
	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
}

func (f *FileRecordedTestsSuite) TestFileSetMetadataTrailingDot() {
	_require := require.New(f.T())
	testName := f.T().Name()

	options := &service.ClientOptions{AllowTrailingDot: to.Ptr(true)}
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := "file.."
	fileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, fileName, 2048, shareClient)

	metadata := map[string]*string{
		"Foo": to.Ptr("Foovalue"),
		"Bar": to.Ptr("Barvalue"),
	}
	_, err = fileClient.SetMetadata(context.Background(), &file.SetMetadataOptions{
		Metadata: metadata,
	})
	_require.NoError(err)

	getResp, err := fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp.ContentLength, int64(2048))
	_require.EqualValues(getResp.Metadata, metadata)
}

func (f *FileRecordedTestsSuite) TestFileUploadClearListRangeTrailingDotOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName) + ".."
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{
		FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup),
		AllowTrailingDot:  to.Ptr(true),
	}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClient, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = fileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	contentSize := 1024 * 2 // 2KB
	contentR, contentD := testcommon.GenerateData(contentSize)
	md5Value := md5.Sum(contentD)
	contentMD5 := md5Value[:]

	uResp, err := fileClient.UploadRange(context.Background(), 0, contentR, &file.UploadRangeOptions{
		TransactionalValidation: file.TransferValidationTypeMD5(contentMD5),
	})
	_require.NoError(err)
	_require.NotNil(uResp.ContentMD5)
	_require.EqualValues(uResp.ContentMD5, contentMD5)

	rangeList, err := fileClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], file.ShareFileRange{Start: to.Ptr(int64(0)), End: to.Ptr(int64(contentSize - 1))})

	cResp, err := fileClient.ClearRange(context.Background(), file.HTTPRange{Offset: 0, Count: int64(contentSize)}, nil)
	_require.NoError(err)
	_require.Nil(cResp.ContentMD5)

	rangeList2, err := fileClient.GetRangeList(context.Background(), nil)
	_require.NoError(err)
	_require.Len(rangeList2.Ranges, 0)
}

func (f *FileUnrecordedTestsSuite) TestFileRenameTrailingDotOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName) + ".."
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{
		FileRequestIntent:      to.Ptr(file.ShareTokenIntentBackup),
		AllowTrailingDot:       to.Ptr(true),
		AllowSourceTrailingDot: to.Ptr(true),
	}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	srcFileClient, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = srcFileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	_, err = srcFileClient.Rename(context.Background(), "file..", nil)
	_require.NoError(err)

	_, err = srcFileClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileRecordedTestsSuite) TestFileRenameNegativeSourceTrailingDot() {
	_require := require.New(f.T())
	testName := f.T().Name()

	options := &service.ClientOptions{
		AllowTrailingDot: to.Ptr(true),
	}
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName) + ".."
	srcFileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, fileName, 2048, shareClient)

	_, err = srcFileClient.Rename(context.Background(), "file..", nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURLTrailingDot() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	options := &service.ClientOptions{
		AllowTrailingDot:       to.Ptr(true),
		AllowSourceTrailingDot: to.Ptr(true),
	}
	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 10
	contentSize := 1024 * 8 // 8KB
	rsc, content := testcommon.GenerateData(contentSize)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	srcFileName := "srcFile.."
	srcFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, srcFileName, fileSize, shareClient)
	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	perms := sas.FilePermissions{Read: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		//Protocol:    sas.ProtocolHTTPS,
		ExpiryTime:  time.Now().UTC().Add(1 * time.Hour),
		ShareName:   shareName,
		FilePath:    srcFileName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	srcFileSAS := srcFClient.URL() + "?" + sasQueryParams.Encode()

	destFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, "destFile..", fileSize, shareClient)

	uResp, err := destFClient.UploadRangeFromURL(context.Background(), srcFileSAS, 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
	})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(uResp.XMSContentCRC64), contentCRC64)

	// validate the content uploaded
	dResp, err := destFClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range: file.HTTPRange{Offset: 0, Count: int64(contentSize)},
	})
	_require.NoError(err)

	data, err := io.ReadAll(dResp.Body)
	defer func() {
		err = dResp.Body.Close()
		_require.NoError(err)
	}()

	_require.EqualValues(data, content)
}

func (f *FileRecordedTestsSuite) TestStartCopyTrailingDotOAuth() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	clOptions := &file.ClientOptions{
		FileRequestIntent:      to.Ptr(file.ShareTokenIntentBackup),
		AllowTrailingDot:       to.Ptr(true),
		AllowSourceTrailingDot: to.Ptr(true),
	}
	testcommon.SetClientOptions(f.T(), &clOptions.ClientOptions)

	srcFileName, destFileName := "srcFile..", "destFile.."
	srcFileClient, err := file.NewClient(fmt.Sprintf("https://%s.file.core.windows.net/%s/%s", accountName, shareName, srcFileName), cred, clOptions)
	_require.NoError(err)
	destFileClient, err := file.NewClient(fmt.Sprintf("https://%s.file.core.windows.net/%s/%s", accountName, shareName, destFileName), cred, clOptions)
	_require.NoError(err)

	fileSize := int64(2048)
	_, err = srcFileClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)

	contentR, srcContent := testcommon.GenerateData(int(fileSize))
	srcContentMD5 := md5.Sum(srcContent)

	_, err = srcFileClient.UploadRange(context.Background(), 0, contentR, nil)
	_require.NoError(err)

	copyResp, err := destFileClient.StartCopyFromURL(context.Background(), srcFileClient.URL(), nil)
	_require.NoError(err)
	_require.NotEqual(copyResp.CopyStatus, "")

	time.Sleep(time.Duration(5) * time.Second)

	getResp, err := destFileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(getResp.CopyID, copyResp.CopyID)
	_require.NotEqual(*getResp.CopyStatus, "")
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.Equal(*getResp.CopySource, srcFileClient.URL())
	}
	_require.Equal(*getResp.CopyStatus, file.CopyStatusTypeSuccess)

	// validate data copied
	dResp, err := destFileClient.DownloadStream(context.Background(), &file.DownloadStreamOptions{
		Range:              file.HTTPRange{Offset: 0, Count: fileSize},
		RangeGetContentMD5: to.Ptr(true),
	})
	_require.NoError(err)

	destContent, err := io.ReadAll(dResp.Body)
	_require.NoError(err)
	_require.EqualValues(srcContent, destContent)
	_require.Equal(dResp.ContentMD5, srcContentMD5[:])

	fileAttributes, err := file.ParseNTFSFileAttributes(dResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURLPreserve() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 20
	srcFileName := "src" + testcommon.GenerateFileName(testName)
	srcFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, srcFileName, fileSize, shareClient)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	perms := sas.FilePermissions{Read: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,
		ExpiryTime:  time.Now().UTC().Add(1 * time.Hour),
		ShareName:   shareName,
		FilePath:    srcFileName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	srcFileSAS := srcFClient.URL() + "?" + sasQueryParams.Encode()

	destFClient := testcommon.GetFileClientFromShare("dest"+testcommon.GenerateFileName(testName), shareClient)
	cResp, err := destFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileLastWriteTime)

	uResp, err := destFClient.UploadRangeFromURL(context.Background(), srcFileSAS, 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
		LastWrittenMode:         to.Ptr(file.LastWrittenModePreserve),
	})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(uResp.XMSContentCRC64), contentCRC64)
	_require.NotNil(uResp.FileLastWriteTime)
	_require.EqualValues(*uResp.FileLastWriteTime, *cResp.FileLastWriteTime)
}

func (f *FileUnrecordedTestsSuite) TestFileUploadRangeFromURLNow() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	var fileSize int64 = 1024 * 20
	srcFileName := "src" + testcommon.GenerateFileName(testName)
	srcFClient := testcommon.CreateNewFileFromShare(context.Background(), _require, srcFileName, fileSize, shareClient)

	contentSize := 1024 * 8 // 8KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	contentCRC64 := crc64.Checksum(content, shared.CRC64Table)

	_, err = srcFClient.UploadRange(context.Background(), 0, rsc, nil)
	_require.NoError(err)

	perms := sas.FilePermissions{Read: true, Write: true}
	sasQueryParams, err := sas.SignatureValues{
		Protocol:    sas.ProtocolHTTPS,
		ExpiryTime:  time.Now().UTC().Add(1 * time.Hour),
		ShareName:   shareName,
		FilePath:    srcFileName,
		Permissions: perms.String(),
	}.SignWithSharedKey(cred)
	_require.NoError(err)

	srcFileSAS := srcFClient.URL() + "?" + sasQueryParams.Encode()

	destFClient := testcommon.GetFileClientFromShare("dest"+testcommon.GenerateFileName(testName), shareClient)
	cResp, err := destFClient.Create(context.Background(), fileSize, nil)
	_require.NoError(err)
	_require.NotNil(cResp.FileLastWriteTime)

	uResp, err := destFClient.UploadRangeFromURL(context.Background(), srcFileSAS, 0, 0, int64(contentSize), &file.UploadRangeFromURLOptions{
		SourceContentValidation: file.SourceContentValidationTypeCRC64(contentCRC64),
		LastWrittenMode:         to.Ptr(file.LastWrittenModeNow),
	})
	_require.NoError(err)
	_require.NotNil(uResp.XMSContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(uResp.XMSContentCRC64), contentCRC64)
	_require.NotNil(uResp.FileLastWriteTime)
	_require.NotEqualValues(*uResp.FileLastWriteTime, *cResp.FileLastWriteTime)
}

type serviceVersionTest struct{}

// newServiceVersionTestPolicy returns a policy that checks the x-ms-version header
func newServiceVersionTestPolicy() policy.Policy {
	return &serviceVersionTest{}
}

func (m serviceVersionTest) Do(req *policy.Request) (*http.Response, error) {
	const versionHeader = "x-ms-version"
	currentVersion := map[string][]string(req.Raw().Header)[versionHeader]
	if currentVersion[0] != generated.ServiceVersion {
		return nil, fmt.Errorf("%s service version doesn't match expected version: %s", currentVersion[0], generated.ServiceVersion)
	}

	return &http.Response{
		Request:    req.Raw(),
		Status:     "Created",
		StatusCode: http.StatusCreated,
		Header:     http.Header{},
		Body:       http.NoBody,
	}, nil
}

func TestServiceVersion(t *testing.T) {
	client, err := file.NewClientWithNoCredential("https://fake/file/testpath", &file.ClientOptions{
		ClientOptions: policy.ClientOptions{
			PerCallPolicies: []policy.Policy{newServiceVersionTestPolicy()},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	_, err = client.Create(context.Background(), 1024, nil)
	require.NoError(t, err)
}

func (f *FileRecordedTestsSuite) TestFileClientDefaultAudience() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{
		FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup),
		Audience:          "https://storage.azure.com/",
	}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClientAudience, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = fileClientAudience.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	_, err = fileClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (f *FileRecordedTestsSuite) TestFileClientCustomAudience() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{
		FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup),
		Audience:          "https://" + accountName + ".file.core.windows.net",
	}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClientAudience, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = fileClientAudience.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	_, err = fileClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (f *FileUnrecordedTestsSuite) TestFileClientAudienceNegative() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fileURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + fileName

	options := &file.ClientOptions{
		FileRequestIntent: to.Ptr(file.ShareTokenIntentBackup),
		Audience:          "https://badaudience.file.core.windows.net",
	}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	fileClientAudience, err := file.NewClient(fileURL, cred, options)
	_require.NoError(err)

	_, err = fileClientAudience.Create(context.Background(), 2048, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidAuthenticationInfo)
}

type fakeDownloadFile struct {
	contentSize int64
	numChunks   uint64
}

// nolint
func (f *fakeDownloadFile) Do(req *http.Request) (*http.Response, error) {
	// check how many times range based get file is called
	if _, ok := req.Header["x-ms-range"]; ok {
		atomic.AddUint64(&f.numChunks, 1)
	}
	return &http.Response{
		Request:    req,
		Status:     "Created",
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Length": []string{fmt.Sprintf("%v", f.contentSize)}},
		Body:       http.NoBody,
	}, nil
}

func TestDownloadSmallChunkSize(t *testing.T) {
	_require := require.New(t)

	fileSize := int64(100 * 1024 * 1024)
	chunkSize := int64(1024)
	numChunks := uint64(((fileSize - 1) / chunkSize) + 1)
	fbb := &fakeDownloadFile{
		contentSize: fileSize,
	}

	log.SetListener(nil) // no logging

	fileClient, err := file.NewClientWithNoCredential("https://fake/file/path", &file.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: fbb,
		},
	})
	_require.NoError(err)
	_require.NotNil(fileClient)

	// download to a temp file and verify contents
	tmp, err := os.CreateTemp("", "")
	_require.NoError(err)
	defer tmp.Close()

	_, err = fileClient.DownloadFile(context.Background(), tmp, &file.DownloadFileOptions{ChunkSize: chunkSize})
	_require.NoError(err)

	_require.Equal(atomic.LoadUint64(&fbb.numChunks), numChunks)

	// reset counter
	atomic.StoreUint64(&fbb.numChunks, 0)

	buff := make([]byte, fileSize)
	_, err = fileClient.DownloadBuffer(context.Background(), buff, &file.DownloadBufferOptions{ChunkSize: chunkSize})
	_require.NoError(err)

	_require.Equal(atomic.LoadUint64(&fbb.numChunks), numChunks)
}

// TODO: Add tests for retry header options

func (f *FileRecordedTestsSuite) TestCreateHardLinkNFS() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)

	directoryName := testcommon.GenerateDirectoryName(testName)
	directoryClient := premiumShareClient.NewRootDirectoryClient().NewSubdirectoryClient(directoryName)
	_, err = directoryClient.Create(context.Background(), nil)
	_require.NoError(err)

	// Create a source file
	sourceFileName := testcommon.GenerateFileName("file1")
	sourceFileClient := directoryClient.NewFileClient(sourceFileName)
	_, err = sourceFileClient.Create(context.Background(), int64(1024), nil)
	_require.NoError(err)

	// Create a hard link to the source file
	hardLinkFileName := testcommon.GenerateFileName("file2")
	hardLinkFileClient := directoryClient.NewFileClient(hardLinkFileName)

	targetFilePath := fmt.Sprintf("/%s/%s", directoryName, sourceFileName)
	resp, err := hardLinkFileClient.CreateHardLink(context.Background(), &file.CreateHardLinkOptions{
		TargetFile: targetFilePath,
	})
	_require.NoError(err)
	_require.NotNil(resp)

	_require.Equal(*resp.NfsFileType, file.NfsFileType("Regular"))
	_require.Equal(resp.Owner, to.Ptr("0"))
	_require.Equal(resp.Group, to.Ptr("0"))
	_require.Equal(resp.FileMode, to.Ptr("0664"))
	_require.Equal(resp.LinkCount, to.Ptr(int64(2)))

	_require.NotNil(resp.FileCreationTime)
	_require.NotNil(resp.FileLastWriteTime)
	_require.NotNil(resp.FileChangeTime)
}

func (f *FileRecordedTestsSuite) TestCreateHardLinkNFSWithLease() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)

	directoryName := testcommon.GenerateDirectoryName(testName)
	directoryClient := premiumShareClient.NewRootDirectoryClient().NewSubdirectoryClient(directoryName)
	_, err = directoryClient.Create(context.Background(), nil)
	_require.NoError(err)

	// Create a source file
	sourceFileName := testcommon.GenerateFileName("file1")
	sourceFileClient := directoryClient.NewFileClient(sourceFileName)
	_, err = sourceFileClient.Create(context.Background(), int64(1024), nil)
	_require.NoError(err)

	leaseId := to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c")
	leaseClient, err := lease.NewShareClient(premiumShareClient, &lease.ShareClientOptions{LeaseID: leaseId})
	_require.NoError(err)

	_, err = leaseClient.Acquire(context.Background(), int32(60), nil)
	_require.NoError(err)

	// Create a hard link to the source file
	hardLinkFileName := testcommon.GenerateFileName("file2")
	hardLinkFileClient := directoryClient.NewFileClient(hardLinkFileName)

	targetFilePath := fmt.Sprintf("/%s/%s", directoryName, sourceFileName)
	resp, err := hardLinkFileClient.CreateHardLink(context.Background(), &file.CreateHardLinkOptions{
		TargetFile:            targetFilePath,
		LeaseAccessConditions: &file.LeaseAccessConditions{LeaseID: leaseId},
	})
	_require.NoError(err)
	_require.NotNil(resp)

	_require.Equal(*resp.NfsFileType, file.NfsFileType("Regular"))
	_require.Equal(resp.Owner, to.Ptr("0"))
	_require.Equal(resp.Group, to.Ptr("0"))
	_require.Equal(resp.FileMode, to.Ptr("0664"))
	_require.Equal(resp.LinkCount, to.Ptr(int64(2)))

	_require.NotNil(resp.FileCreationTime)
	_require.NotNil(resp.FileLastWriteTime)
	_require.NotNil(resp.FileChangeTime)
	_, err = leaseClient.Release(context.Background(), nil)
	_require.NoError(err)
}

func (f *FileRecordedTestsSuite) TestCreateHardLinkNFSNilOptions() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)

	directoryName := testcommon.GenerateDirectoryName(testName)
	directoryClient := premiumShareClient.NewRootDirectoryClient().NewSubdirectoryClient(directoryName)
	_, err = directoryClient.Create(context.Background(), nil)
	_require.NoError(err)

	// Create a source file
	sourceFileName := testcommon.GenerateFileName("file1")
	sourceFileClient := directoryClient.NewFileClient(sourceFileName)
	_, err = sourceFileClient.Create(context.Background(), int64(1024), nil)
	_require.NoError(err)

	// Create a hard link to the source file
	hardLinkFileName := testcommon.GenerateFileName("file2")
	hardLinkFileClient := directoryClient.NewFileClient(hardLinkFileName)

	resp, err := hardLinkFileClient.CreateHardLink(context.Background(), nil)
	_require.Error(err)
	_require.Error(err, "targetFile cannot be nil")
	_require.NotNil(resp)
}
