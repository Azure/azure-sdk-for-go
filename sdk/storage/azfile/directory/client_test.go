//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running directory Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &DirectoryRecordedTestsSuite{})
		suite.Run(t, &DirectoryUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &DirectoryRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &DirectoryRecordedTestsSuite{})
	}
}

func (d *DirectoryRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(d.T(), suite, test)
}

func (d *DirectoryRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(d.T(), suite, test)
}

func (d *DirectoryUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (d *DirectoryUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type DirectoryRecordedTestsSuite struct {
	suite.Suite
}

type DirectoryUnrecordedTestsSuite struct {
	suite.Suite
}

func (d *DirectoryRecordedTestsSuite) TestDirNewDirectoryClient() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := shareClient.NewDirectoryClient(dirName)

	subDirName := "inner" + dirName
	subDirClient := dirClient.NewSubdirectoryClient(subDirName)

	correctURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName + "/" + subDirName
	_require.Equal(subDirClient.URL(), correctURL)
}

func (d *DirectoryRecordedTestsSuite) TestDirCreateFileURL() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := svcClient.NewShareClient(shareName)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := shareClient.NewDirectoryClient(dirName)

	fileName := testcommon.GenerateFileName(testName)
	fileClient := dirClient.NewFileClient(fileName)

	correctURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName + "/" + fileName
	_require.Equal(fileClient.URL(), correctURL)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryCreateUsingSharedKey() {
	_require := require.New(d.T())
	testName := d.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClientWithSharedKeyCredential(dirURL, cred, options)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryCreateUsingConnectionString() {
	_require := require.New(d.T())
	testName := d.T().Name()

	connString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	options := &directory.ClientOptions{}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClientFromConnectionString(*connString, shareName, dirName, options)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	innerDirName1 := "innerdir1"
	dirPath := dirName + "/" + innerDirName1
	dirClient1, err := directory.NewClientFromConnectionString(*connString, shareName, dirPath, options)
	_require.NoError(err)

	resp, err = dirClient1.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)

	innerDirName2 := "innerdir2"
	// using '\' as path separator between directories
	dirPath = dirName + "\\" + innerDirName1 + "\\" + innerDirName2
	dirClient2, err := directory.NewClientFromConnectionString(*connString, shareName, dirPath, options)
	_require.NoError(err)

	resp, err = dirClient2.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryCreateNegativeMultiLevel() {
	_require := require.New(d.T())
	testName := d.T().Name()

	connString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	// dirPath where parent dir does not exist
	dirPath := "a/b/c/d/" + dirName
	options := &directory.ClientOptions{}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClientFromConnectionString(*connString, shareName, dirPath, options)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.Error(err)
	_require.Nil(resp.RequestID)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)
}

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryClientUsingSAS() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	permissions := sas.FilePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)

	dirSASURL, err := dirClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	dirSASClient, err := directory.NewClientWithNoCredential(dirSASURL, nil)
	_require.NoError(err)

	_, err = dirSASClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.AuthenticationFailed)

	subDirSASClient := dirSASClient.NewSubdirectoryClient("subdir")
	_, err = subDirSASClient.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.AuthenticationFailed)

	// TODO: directory SAS client unable to do create and get properties on directories.
	// Also unable to do create or get properties on files. Validate this behaviour.
	fileSASClient := dirSASClient.NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fileSASClient.Create(context.Background(), 1024, nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.AuthenticationFailed)

	_, err = fileSASClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.AuthenticationFailed)

	// create file using shared key client
	_, err = dirClient.NewFileClient(testcommon.GenerateFileName(testName)).Create(context.Background(), 1024, nil)
	_require.NoError(err)

	// get properties using SAS client
	_, err = fileSASClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.AuthenticationFailed)
}

func (d *DirectoryRecordedTestsSuite) TestDirCreateDeleteDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := shareClient.NewDirectoryClient(dirName)
	_require.NoError(err)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(cResp.RequestID)
	_require.NotNil(cResp.ETag)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.Equal(cResp.FileCreationTime.IsZero(), false)
	_require.Equal(cResp.FileLastWriteTime.IsZero(), false)
	_require.Equal(cResp.FileChangeTime.IsZero(), false)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.Equal(gResp.FileCreationTime.IsZero(), false)
	_require.Equal(gResp.FileLastWriteTime.IsZero(), false)
	_require.Equal(gResp.FileChangeTime.IsZero(), false)
}

func (d *DirectoryUnrecordedTestsSuite) TestDirSetPropertiesNonDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.GetDirectoryClient(dirName, shareClient)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(cResp.FilePermissionKey)

	creationTime := time.Now().Add(5 * time.Minute).Round(time.Microsecond)
	lastWriteTime := time.Now().Add(10 * time.Minute).Round(time.Millisecond)

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{
				ReadOnly: true,
				System:   true,
			},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
		},
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)
	_require.NotNil(sResp.FileCreationTime)
	_require.NotNil(sResp.FileLastWriteTime)
	_require.NotNil(sResp.FilePermissionKey)
	_require.NotEqual(*sResp.FilePermissionKey, *cResp.FilePermissionKey)
	_require.Equal(*sResp.FileCreationTime, creationTime.UTC())
	_require.Equal(*sResp.FileLastWriteTime, lastWriteTime.UTC())

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.FileCreationTime)
	_require.NotNil(gResp.FileLastWriteTime)
	_require.NotNil(gResp.FilePermissionKey)
	_require.Equal(*gResp.FilePermissionKey, *sResp.FilePermissionKey)
	_require.Equal(*gResp.FileCreationTime, *sResp.FileCreationTime)
	_require.Equal(*gResp.FileLastWriteTime, *sResp.FileLastWriteTime)
	_require.Equal(*gResp.FileAttributes, *sResp.FileAttributes)
}

func (d *DirectoryUnrecordedTestsSuite) TestDirCreateDeleteNonDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.GetDirectoryClient(dirName, shareClient)

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}

	cResp, err := dirClient.Create(context.Background(), &directory.CreateOptions{
		Metadata: md,
		FileSMBProperties: &file.SMBProperties{
			Attributes:    &file.NTFSFileAttributes{None: true},
			CreationTime:  to.Ptr(time.Now().Add(5 * time.Minute)),
			LastWriteTime: to.Ptr(time.Now().Add(10 * time.Minute)),
		},
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FilePermissionKey)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.ETag)
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotNil(cResp.RequestID)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.FilePermissionKey, *cResp.FilePermissionKey)
	_require.EqualValues(gResp.Metadata, md)

	// Creating again will result in 409 and ResourceAlreadyExists.
	_, err = dirClient.Create(context.Background(), &directory.CreateOptions{Metadata: md})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceAlreadyExists)

	dResp, err := dirClient.Delete(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(dResp.Date.IsZero(), false)
	_require.NotNil(dResp.RequestID)
	_require.NotNil(dResp.Version)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirCreateNegativePermissions() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.GetDirectoryClient(dirName, shareClient)
	subDirClient := dirClient.NewSubdirectoryClient("subdir" + dirName)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(cResp.FilePermissionKey)

	// having both Permission and PermissionKey set returns error
	_, err = subDirClient.Create(context.Background(), &directory.CreateOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{None: true},
		},
		FilePermissions: &file.Permissions{
			Permission:    &testcommon.SampleSDDL,
			PermissionKey: cResp.FilePermissionKey,
		},
	})
	_require.Error(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirCreateNegativeAttributes() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := testcommon.GetDirectoryClient(testcommon.GenerateDirectoryName(testName), shareClient)

	// None attribute must be used alone.
	_, err = dirClient.Create(context.Background(), &directory.CreateOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{None: true, ReadOnly: true},
		},
	})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidHeaderValue)
}

func (d *DirectoryRecordedTestsSuite) TestDirCreateDeleteNegativeMultiLevelDir() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	parentDirName := "parent" + testcommon.GenerateDirectoryName(testName)
	parentDirClient := shareClient.NewDirectoryClient(parentDirName)

	subDirName := "subdir" + testcommon.GenerateDirectoryName(testName)
	subDirClient := parentDirClient.NewSubdirectoryClient(subDirName)

	// Directory create with subDirClient
	_, err = subDirClient.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)

	_, err = parentDirClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = subDirClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = subDirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// Delete level by level
	// Delete Non-empty directory should fail
	_, err = parentDirClient.Delete(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.DirectoryNotEmpty)

	_, err = subDirClient.Delete(context.Background(), nil)
	_require.NoError(err)

	_, err = parentDirClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirCreateEndWithSlash() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName) + "/"
	dirClient := testcommon.GetDirectoryClient(dirName, shareClient)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.ETag)
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotNil(cResp.RequestID)
	_require.NotNil(cResp.Version)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirGetSetMetadataDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)
	defer testcommon.DeleteDirectory(context.Background(), _require, dirClient)

	sResp, err := dirClient.SetMetadata(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotNil(sResp.ETag)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.NotNil(gResp.IsServerEncrypted)
	_require.Len(gResp.Metadata, 0)
}

func (d *DirectoryRecordedTestsSuite) TestDirGetSetMetadataNonDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)
	defer testcommon.DeleteDirectory(context.Background(), _require, dirClient)

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}

	sResp, err := dirClient.SetMetadata(context.Background(), &directory.SetMetadataOptions{
		Metadata: md,
	})
	_require.NoError(err)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotNil(sResp.ETag)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.Version)
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.NotNil(gResp.IsServerEncrypted)
	_require.EqualValues(gResp.Metadata, md)
}

func (d *DirectoryRecordedTestsSuite) TestDirSetMetadataNegative() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)
	defer testcommon.DeleteDirectory(context.Background(), _require, dirClient)

	md := map[string]*string{
		"!@#$%^&*()": to.Ptr("!@#$%^&*()"),
	}

	_, err = dirClient.SetMetadata(context.Background(), &directory.SetMetadataOptions{
		Metadata: md,
	})
	_require.Error(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirGetPropertiesNegative() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.GetDirectoryClient(dirName, shareClient)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirGetPropertiesWithBaseDirectory() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := shareClient.NewRootDirectoryClient()

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.Equal(gResp.FileCreationTime.IsZero(), false)
	_require.Equal(gResp.FileLastWriteTime.IsZero(), false)
	_require.Equal(gResp.FileChangeTime.IsZero(), false)
	_require.NotNil(gResp.IsServerEncrypted)
}

func (d *DirectoryRecordedTestsSuite) TestDirGetSetMetadataMergeAndReplace() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)
	defer testcommon.DeleteDirectory(context.Background(), _require, dirClient)

	md := map[string]*string{
		"Color": to.Ptr("RED"),
	}

	sResp, err := dirClient.SetMetadata(context.Background(), &directory.SetMetadataOptions{
		Metadata: md,
	})
	_require.NoError(err)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.NotNil(gResp.IsServerEncrypted)
	_require.EqualValues(gResp.Metadata, md)

	md2 := map[string]*string{
		"Color": to.Ptr("WHITE"),
	}

	sResp, err = dirClient.SetMetadata(context.Background(), &directory.SetMetadataOptions{
		Metadata: md2,
	})
	_require.NoError(err)
	_require.NotNil(sResp.RequestID)
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.ETag)
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotNil(gResp.RequestID)
	_require.NotNil(gResp.Version)
	_require.NotNil(gResp.IsServerEncrypted)
	_require.EqualValues(gResp.Metadata, md2)
}

func (d *DirectoryRecordedTestsSuite) TestSASDirectoryClientNoKey() {
	_require := require.New(d.T())
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	testName := d.T().Name()
	shareName := testcommon.GenerateShareName(testName)
	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient, err := directory.NewClientWithNoCredential(fmt.Sprintf("https://%s.file.core.windows.net/%v/%v", accountName, shareName, dirName), nil)
	_require.NoError(err)

	permissions := sas.FilePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		Create: true,
	}
	expiry := time.Now().Add(time.Hour)

	_, err = dirClient.GetSASURL(permissions, expiry, nil)
	_require.Equal(err, fileerror.MissingSharedKeyCredential)
}

func (d *DirectoryRecordedTestsSuite) TestSASDirectoryClientSignNegative() {
	_require := require.New(d.T())
	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)
	_require.Greater(len(accountKey), 0)

	cred, err := service.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	testName := d.T().Name()
	shareName := testcommon.GenerateShareName(testName)
	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient, err := directory.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.file.core.windows.net/%v%v", accountName, shareName, dirName), cred, nil)
	_require.NoError(err)

	permissions := sas.FilePermissions{
		Read:   true,
		Write:  true,
		Delete: true,
		Create: true,
	}
	expiry := time.Time{}

	_, err = dirClient.GetSASURL(permissions, expiry, nil)
	_require.Equal(err.Error(), "service SAS is missing at least one of these: ExpiryTime or Permissions")
}

// TODO: add tests for listing files and directories after file client is completed

// TODO: add tests for ListHandles and ForceCloseHandles
