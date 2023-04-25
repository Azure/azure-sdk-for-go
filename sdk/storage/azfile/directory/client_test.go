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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
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

func (d *DirectoryRecordedTestsSuite) TestDirSetPropertiesDefault() {
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

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(sResp.FileCreationTime)
	_require.NotNil(sResp.FileLastWriteTime)
	_require.NotNil(sResp.FilePermissionKey)
	_require.Equal(*sResp.FilePermissionKey, *cResp.FilePermissionKey)

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

func (d *DirectoryRecordedTestsSuite) TestDirSetPropertiesNonDefault() {
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

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 21:00:00 GMT 2023")
	_require.NoError(err)
	creationTime := currTime.Add(5 * time.Minute).Round(time.Microsecond)
	lastWriteTime := currTime.Add(10 * time.Minute).Round(time.Millisecond)

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

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)

	for i := 0; i < 10; i++ {
		_ = testcommon.CreateNewDirectory(context.Background(), _require, dirName+fmt.Sprintf("%v", i), shareClient)
	}

	for i := 0; i < 5; i++ {
		_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fileName+fmt.Sprintf("%v", i), 2048, shareClient)
	}

	dirCtr, fileCtr := 0, 0
	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		dirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
		for _, dir := range resp.Segment.Directories {
			_require.NotNil(dir.Name)
			_require.NotNil(dir.ID)
			_require.Nil(dir.Attributes)
			_require.Nil(dir.PermissionKey)
			_require.Nil(dir.Properties.ETag)
			_require.Nil(dir.Properties.ChangeTime)
			_require.Nil(dir.Properties.CreationTime)
			_require.Nil(dir.Properties.ContentLength)
		}
		for _, f := range resp.Segment.Files {
			_require.NotNil(f.Name)
			_require.NotNil(f.ID)
			_require.Nil(f.Attributes)
			_require.Nil(f.PermissionKey)
			_require.Nil(f.Properties.ETag)
			_require.Nil(f.Properties.ChangeTime)
			_require.Nil(f.Properties.CreationTime)
			_require.NotNil(f.Properties.ContentLength)
			_require.Equal(*f.Properties.ContentLength, int64(2048))
		}
	}
	_require.Equal(dirCtr, 10)
	_require.Equal(fileCtr, 5)
}

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsInclude() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)

	for i := 0; i < 10; i++ {
		_ = testcommon.CreateNewDirectory(context.Background(), _require, dirName+fmt.Sprintf("%v", i), shareClient)
	}

	for i := 0; i < 5; i++ {
		_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fileName+fmt.Sprintf("%v", i), 2048, shareClient)
	}

	dirCtr, fileCtr := 0, 0
	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		Include:             directory.ListFilesInclude{Timestamps: true, ETag: true, Attributes: true, PermissionKey: true},
		IncludeExtendedInfo: to.Ptr(true),
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		dirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
		for _, dir := range resp.Segment.Directories {
			_require.NotNil(dir.Name)
			_require.NotNil(dir.ID)
			_require.NotNil(dir.Attributes)
			_require.NotNil(dir.PermissionKey)
			_require.NotNil(dir.Properties.ETag)
			_require.NotNil(dir.Properties.ChangeTime)
			_require.NotNil(dir.Properties.CreationTime)
			_require.Nil(dir.Properties.ContentLength)
		}
		for _, f := range resp.Segment.Files {
			_require.NotNil(f.Name)
			_require.NotNil(f.ID)
			_require.NotNil(f.Attributes)
			_require.NotNil(f.PermissionKey)
			_require.NotNil(f.Properties.ETag)
			_require.NotNil(f.Properties.ChangeTime)
			_require.NotNil(f.Properties.CreationTime)
			_require.NotNil(f.Properties.ContentLength)
			_require.Equal(*f.Properties.ContentLength, int64(2048))
		}
	}
	_require.Equal(dirCtr, 10)
	_require.Equal(fileCtr, 5)
}

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsMaxResultsAndMarker() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)

	for i := 0; i < 10; i++ {
		_ = testcommon.CreateNewDirectory(context.Background(), _require, dirName+fmt.Sprintf("%v", i), shareClient)
	}

	for i := 0; i < 5; i++ {
		_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fileName+fmt.Sprintf("%v", i), 2048, shareClient)
	}

	dirCtr, fileCtr := 0, 0
	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		MaxResults: to.Ptr(int32(2)),
	})
	resp, err := pager.NextPage(context.Background())
	_require.NoError(err)
	dirCtr += len(resp.Segment.Directories)
	fileCtr += len(resp.Segment.Files)
	_require.Equal(dirCtr+fileCtr, 2)

	pager = shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		Marker:     resp.NextMarker,
		MaxResults: to.Ptr(int32(5)),
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		dirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
	}
	_require.Equal(dirCtr, 10)
	_require.Equal(fileCtr, 5)
}

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsWithPrefix() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)

	for i := 0; i < 10; i++ {
		_ = testcommon.CreateNewDirectory(context.Background(), _require, fmt.Sprintf("%v", i)+dirName, shareClient)
	}

	for i := 0; i < 5; i++ {
		_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fmt.Sprintf("%v", i)+fileName, 2048, shareClient)
	}

	dirCtr, fileCtr := 0, 0
	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		Prefix: to.Ptr("1"),
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		dirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
		if len(resp.Segment.Directories) > 0 {
			_require.NotNil(resp.Segment.Directories[0].Name)
			_require.Equal(*resp.Segment.Directories[0].Name, "1"+dirName)
		}
		if len(resp.Segment.Files) > 0 {
			_require.NotNil(resp.Segment.Files[0].Name)
			_require.Equal(*resp.Segment.Files[0].Name, "1"+fileName)
		}
	}
	_require.Equal(dirCtr, 1)
	_require.Equal(fileCtr, 1)
}

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsMaxResultsNegative() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)

	for i := 0; i < 2; i++ {
		_ = testcommon.CreateNewDirectory(context.Background(), _require, dirName+fmt.Sprintf("%v", i), shareClient)
	}

	for i := 0; i < 2; i++ {
		_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fileName+fmt.Sprintf("%v", i), 2048, shareClient)
	}

	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		MaxResults: to.Ptr(int32(-1)),
	})
	_, err = pager.NextPage(context.Background())
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.OutOfRangeQueryParameterValue)
}

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsSnapshot() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer func() {
		_, err := shareClient.Delete(context.Background(), &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})
		_require.NoError(err)
	}()

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)

	for i := 0; i < 10; i++ {
		_ = testcommon.CreateNewDirectory(context.Background(), _require, dirName+fmt.Sprintf("%v", i), shareClient)
	}

	for i := 0; i < 5; i++ {
		_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fileName+fmt.Sprintf("%v", i), 2048, shareClient)
	}

	snapResp, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(snapResp.Snapshot)

	_, err = shareClient.NewRootDirectoryClient().GetProperties(context.Background(), &directory.GetPropertiesOptions{ShareSnapshot: snapResp.Snapshot})
	_require.NoError(err)

	dirCtr, fileCtr := 0, 0
	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		ShareSnapshot: snapResp.Snapshot,
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		dirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
	}
	_require.Equal(dirCtr, 10)
	_require.Equal(fileCtr, 5)
}

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsInsideDir() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	fileName := testcommon.GenerateFileName(testName)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	for i := 0; i < 5; i++ {
		_, err = dirClient.NewSubdirectoryClient("subdir"+fmt.Sprintf("%v", i)).Create(context.Background(), nil)
		_require.NoError(err)
	}

	for i := 0; i < 5; i++ {
		_, err = dirClient.NewFileClient(fileName+fmt.Sprintf("%v", i)).Create(context.Background(), 0, nil)
		_require.NoError(err)
	}

	dirCtr, fileCtr := 0, 0
	pager := dirClient.NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		Include: directory.ListFilesInclude{Timestamps: true, ETag: true, Attributes: true, PermissionKey: true},
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		dirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
		for _, dir := range resp.Segment.Directories {
			_require.NotNil(dir.Name)
			_require.NotNil(dir.ID)
			_require.NotNil(dir.Attributes)
			_require.NotNil(dir.PermissionKey)
			_require.NotNil(dir.Properties.ETag)
			_require.NotNil(dir.Properties.ChangeTime)
			_require.NotNil(dir.Properties.CreationTime)
			_require.Nil(dir.Properties.ContentLength)
		}
		for _, f := range resp.Segment.Files {
			_require.NotNil(f.Name)
			_require.NotNil(f.ID)
			_require.NotNil(f.Attributes)
			_require.NotNil(f.PermissionKey)
			_require.NotNil(f.Properties.ETag)
			_require.NotNil(f.Properties.ChangeTime)
			_require.NotNil(f.Properties.CreationTime)
			_require.NotNil(f.Properties.ContentLength)
			_require.Equal(*f.Properties.ContentLength, int64(0))
		}
	}
	_require.Equal(dirCtr, 5)
	_require.Equal(fileCtr, 5)
}

func (d *DirectoryRecordedTestsSuite) TestDirListHandlesDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	resp, err := dirClient.ListHandles(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Handles, 0)
	_require.NotNil(resp.NextMarker)
	_require.Equal(*resp.NextMarker, "")
}

func (d *DirectoryRecordedTestsSuite) TestDirForceCloseHandlesDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	resp, err := dirClient.ForceCloseHandles(context.Background(), "*", nil)
	_require.NoError(err)
	_require.EqualValues(*resp.NumberOfHandlesClosed, 0)
	_require.EqualValues(*resp.NumberOfHandlesFailedToClose, 0)
	_require.Nil(resp.Marker)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryCreateNegativeWithoutSAS() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClientWithNoCredential(dirURL, nil)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.Error(err)
}

// TODO: add tests for listing files and directories after file client is completed

// TODO: add tests for ListHandles and ForceCloseHandles
