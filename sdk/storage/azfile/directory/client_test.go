//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
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

func (d *DirectoryUnrecordedTestsSuite) TestDirNewDirectoryClient() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, err := testcommon.GetRequiredEnv(testcommon.AccountNameEnvVar)
	_require.NoError(err)

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

func (d *DirectoryUnrecordedTestsSuite) TestDirCreateFileURL() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, err := testcommon.GetRequiredEnv(testcommon.AccountNameEnvVar)
	_require.NoError(err)

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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryCreateUsingSharedKey() {
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
	dirClient, err := directory.NewClientWithSharedKeyCredential(dirURL, cred, nil)
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryCreateUsingConnectionString() {
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
	dirClient, err := directory.NewClientFromConnectionString(*connString, shareName, dirName, nil)
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
	dirClient1, err := directory.NewClientFromConnectionString(*connString, shareName, dirPath, nil)
	_require.NoError(err)

	resp, err = dirClient1.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)

	innerDirName2 := "innerdir2"
	// using '\' as path separator between directories
	dirPath = dirName + "\\" + innerDirName1 + "\\" + innerDirName2
	dirClient2, err := directory.NewClientFromConnectionString(*connString, shareName, dirPath, nil)
	_require.NoError(err)

	resp, err = dirClient2.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
}

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryCreateNegativeMultiLevel() {
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
	dirClient, err := directory.NewClientFromConnectionString(*connString, shareName, dirPath, nil)
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

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := shareClient.NewDirectoryClient(dirName)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

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

	// TODO: create files using dirSASClient
}

func (d *DirectoryUnrecordedTestsSuite) TestDirCreateDeleteDefault() {
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
	dirClient := shareClient.NewDirectoryClient(dirName)

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
	dirClient := shareClient.NewDirectoryClient(dirName)

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

func (d *DirectoryUnrecordedTestsSuite) TestDirCreateDeleteNegativeMultiLevelDir() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirCreateEndWithSlash() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	directoryName := testcommon.GenerateDirectoryName(testName) + "/"
	dirClient := shareClient.NewDirectoryClient(directoryName)

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
