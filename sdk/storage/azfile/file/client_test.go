//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
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
}

type FileUnrecordedTestsSuite struct {
	suite.Suite
}

func (f *FileUnrecordedTestsSuite) TestFileNewFileClient() {
	_require := require.New(f.T())
	testName := f.T().Name()

	accountName, err := testcommon.GetRequiredEnv(testcommon.AccountNameEnvVar)
	_require.NoError(err)

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

func (f *FileUnrecordedTestsSuite) TestFileCreateUsingSharedKey() {
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

	fileClient, err := file.NewClientWithSharedKeyCredential(fileURL, cred, nil)
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
}

func (f *FileUnrecordedTestsSuite) TestFileCreateUsingConnectionString() {
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
	fileClient1, err := file.NewClientFromConnectionString(*connString, shareName, fileName, nil)
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
	fileClient2, err := file.NewClientFromConnectionString(*connString, shareName, filePath, nil)
	_require.NoError(err)

	_, err = fileClient2.Create(context.Background(), 1024, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)

	testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	// using '\' as path separator
	filePath = dirName + "\\" + fileName
	fileClient3, err := file.NewClientFromConnectionString(*connString, shareName, filePath, nil)
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

func (f *FileUnrecordedTestsSuite) TestFileCreateDeleteDefault() {
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

func (f *FileUnrecordedTestsSuite) TestFileCreateNonDefaultMetadataNonEmpty() {
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

func (f *FileUnrecordedTestsSuite) TestFileCreateNonDefaultHTTPHeaders() {
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

func (f *FileUnrecordedTestsSuite) TestFileCreateNegativeMetadataInvalid() {
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
	var testMd5 []byte
	copy(testMd5[:], md5Str)

	creationTime := time.Now().Add(-time.Hour)
	lastWriteTime := time.Now().Add(-time.Minute * 15)

	options := &file.SetHTTPHeadersOptions{
		Permissions: &file.Permissions{Permission: &testcommon.SampleSDDL},
		SMBProperties: &file.SMBProperties{
			Attributes:    &file.NTFSFileAttributes{Hidden: true},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
		},
		HTTPHeaders: &file.HTTPHeaders{
			ContentType:        to.Ptr("text/html"),
			ContentEncoding:    to.Ptr("gzip"),
			ContentLanguage:    to.Ptr("tr,en"),
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

	_require.EqualValues(*getResp.FileCreationTime, creationTime.UTC())
	_require.EqualValues(*getResp.FileLastWriteTime, lastWriteTime.UTC())

	_require.NotNil(getResp.ETag)
	_require.NotNil(getResp.RequestID)
	_require.NotNil(getResp.Version)
	_require.Equal(getResp.Date.IsZero(), false)
	_require.NotNil(getResp.IsServerEncrypted)
}
