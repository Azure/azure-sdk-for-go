//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
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
	attribs := getResp.FileAttributes

	md5Str := "MDAwMDAwMDA="
	var testMd5 []byte
	copy(testMd5[:], md5Str)

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
	var testMd5 []byte
	copy(testMd5[:], md5Str)

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
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

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
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			CreationTime:       file.SourceCopyFileCreationTime{},
			LastWriteTime:      file.SourceCopyFileLastWriteTime{},
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
	_require.EqualValues(resp2.FileAttributes, cResp.FileAttributes)
	_require.EqualValues(resp2.FilePermissionKey, cResp.FilePermissionKey)
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
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FileCreationTime)
	_require.NotNil(cResp.FileLastWriteTime)
	_require.NotNil(cResp.FileAttributes)
	_require.NotNil(cResp.FilePermissionKey)

	destCreationTime := currTime.Add(5 * time.Minute)
	destLastWriteTIme := currTime.Add(6 * time.Minute)
	_, err = copyFClient.StartCopyFromURL(context.Background(), fClient.URL(), &file.StartCopyFromURLOptions{
		CopyFileSMBInfo: &file.CopyFileSMBInfo{
			CreationTime:  file.DestinationCopyFileCreationTime(destCreationTime),
			LastWriteTime: file.DestinationCopyFileLastWriteTime(destLastWriteTIme),
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

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
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

	_, err = fileClient.GetSASURL(permissions, expiry, nil)
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")
}

// TODO: Add tests for different options of StartCopyFromURL()

// TODO: Add tests for upload and download methods

// TODO: Add tests for GetRangeList, ListHandles and ForceCloseHandles
