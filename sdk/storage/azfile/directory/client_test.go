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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/lease"
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

func (d *DirectoryRecordedTestsSuite) SetupSuite() {
	d.proxy = testcommon.SetupSuite(&d.Suite)
}

func (d *DirectoryRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&d.Suite, d.proxy)
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
	proxy *recording.TestProxyInstance
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
	dirClient := shareClient.NewDirectoryClient(dirName + "/") // directory name having trailing '/'

	correctDirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName
	_require.Equal(dirClient.URL(), correctDirURL)

	subDirName := "inner" + dirName
	subDirClient := dirClient.NewSubdirectoryClient(subDirName + "/") // subdirectory name having trailing '/'

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

func (d *DirectoryRecordedTestsSuite) TestDirCreateNfs() {
	_require := require.New(d.T())
	testName := d.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)
	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	owner := "345"
	group := "123"
	mode := "6444"

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)
	_require.NoError(err)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := premiumShareClient.NewDirectoryClient(dirName)
	_require.NoError(err)

	cResp, err := dirClient.Create(context.Background(), &directory.CreateOptions{
		Owner:    to.Ptr(owner),
		Group:    to.Ptr(group),
		FileMode: to.Ptr(mode),
	})
	_require.NoError(err)
	_require.NotNil(cResp.ETag)
	_require.Equal(*cResp.Owner, owner)
	_require.Equal(*cResp.Group, group)
	_require.Equal(*cResp.FileMode, mode)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.Owner, owner)
	_require.Equal(*gResp.Group, group)
	_require.Equal(*gResp.FileMode, mode)
	_require.Equal(*gResp.NfsFileType, file.NfsFileTypeDirectory)
}

func (d *DirectoryRecordedTestsSuite) TestDirCreateRenameFilePermissionFormatDefault() {
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

	resp, err := dirClient.Create(context.Background(), &directory.CreateOptions{
		FilePermissionFormat: (*file.PermissionFormat)(to.Ptr(testcommon.FilePermissionFormatSddl)),
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)

	_, err = dirClient.Rename(context.Background(), "testFile", &directory.RenameOptions{
		FilePermissionFormat: (*file.PermissionFormat)(to.Ptr(testcommon.FilePermissionBinary)),
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleBinary,
		},
	})
	_require.NoError(err)
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
	changeTime := currTime.Add(15 * time.Minute).Round(time.Millisecond)

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{
				ReadOnly: true,
				System:   true,
			},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
			ChangeTime:    &changeTime,
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
	_require.Equal(*sResp.FileChangeTime, changeTime.UTC())

	fileAttributes, err := file.ParseNTFSFileAttributes(sResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.ReadOnly)
	_require.True(fileAttributes.System)
	_require.True(fileAttributes.Directory)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.FileCreationTime)
	_require.NotNil(gResp.FileLastWriteTime)
	_require.NotNil(gResp.FilePermissionKey)
	_require.Equal(*gResp.FilePermissionKey, *sResp.FilePermissionKey)
	_require.Equal(*gResp.FileCreationTime, *sResp.FileCreationTime)
	_require.Equal(*gResp.FileLastWriteTime, *sResp.FileLastWriteTime)
	_require.Equal(*gResp.FileChangeTime, *sResp.FileChangeTime)
	_require.Equal(*gResp.FileAttributes, *sResp.FileAttributes)

	fileAttributes2, err := file.ParseNTFSFileAttributes(gResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.ReadOnly)
	_require.True(fileAttributes2.System)
	_require.True(fileAttributes2.Directory)
	_require.EqualValues(fileAttributes2, fileAttributes)
}

func (d *DirectoryRecordedTestsSuite) TestDirSetPropertiesFilePermissionFormat() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.GetDirectoryClient(dirName, shareClient)

	cResp, err := dirClient.Create(context.Background(), &directory.CreateOptions{
		FilePermissionFormat: (*file.PermissionFormat)(to.Ptr(testcommon.FilePermissionFormatSddl)),
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FilePermissionKey)
	_require.NoError(err)

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		FilePermissionFormat: (*directory.FilePermissionFormat)(to.Ptr(testcommon.FilePermissionFormatSddl)),
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)
	_require.NotNil(sResp.FileCreationTime)
	_require.NotNil(sResp.FileLastWriteTime)
	_require.NotNil(sResp.FilePermissionKey)
	_require.NotEqual(*sResp.FilePermissionKey, *cResp.FilePermissionKey)
}

func (d *DirectoryRecordedTestsSuite) TestDirSetPropertiesNfs() {
	_require := require.New(d.T())
	testName := d.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)
	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	owner := "345"
	group := "123"
	mode := "7777"

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, options)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)
	_require.NoError(err)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.GetDirectoryClient(dirName, premiumShareClient)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	// Set the custom permissions
	_, err = dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		Owner:    to.Ptr(owner),
		Group:    to.Ptr(group),
		FileMode: to.Ptr(mode),
	})
	_require.NoError(err)

	response, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*response.FileMode, mode)
	_require.Equal(*response.Group, group)
	_require.Equal(*response.Owner, owner)
	_require.Equal(*response.NfsFileType, file.NfsFileTypeDirectory)

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
			ChangeTime:    to.Ptr(time.Now().Add(10 * time.Minute)),
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

	fileAttributes, err := file.ParseNTFSFileAttributes(cResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.Directory)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*gResp.FilePermissionKey, *cResp.FilePermissionKey)
	_require.EqualValues(gResp.Metadata, md)

	fileAttributes2, err := file.ParseNTFSFileAttributes(gResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.Directory)
	_require.EqualValues(fileAttributes, fileAttributes2)

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

func (d *DirectoryRecordedTestsSuite) TestDirCreateEmptyAttributes() {
	_require := require.New(d.T())
	testName := d.T().Name()
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.GetDirectoryClient(dirName, shareClient)

	cResp, err := dirClient.Create(context.Background(), &directory.CreateOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{},
		},
	})
	_require.NoError(err)
	_require.NotNil(cResp.FilePermissionKey)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.ETag)
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotNil(cResp.RequestID)

	fileAttributes, err := file.ParseNTFSFileAttributes(cResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.Directory)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	fileAttributes2, err := file.ParseNTFSFileAttributes(gResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.Directory)
	_require.EqualValues(fileAttributes, fileAttributes2)
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
			_require.Greater(len(*dir.Name), 0)
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
			_require.Greater(len(*f.Name), 0)
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
		_require.NotNil(resp.Prefix)
		_require.Equal(*resp.Prefix, "1")
		if len(resp.Segment.Directories) > 0 {
			_require.NotNil(resp.Segment.Directories[0].Name)
			_require.Greater(len(*resp.Segment.Directories[0].Name), 0)
			_require.Equal(*resp.Segment.Directories[0].Name, "1"+dirName)
		}
		if len(resp.Segment.Files) > 0 {
			_require.NotNil(resp.Segment.Files[0].Name)
			_require.Greater(len(*resp.Segment.Files[0].Name), 0)
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

func (d *DirectoryRecordedTestsSuite) TestDirectoryCreateWithTrailingSlash() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName) + "/"
	dirClient := shareClient.NewDirectoryClient(dirName)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	subDirClient := dirClient.NewSubdirectoryClient("subdir/")

	_, err = subDirClient.Create(context.Background(), nil)
	_require.NoError(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryCreateDeleteUsingOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	resp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.FileCreationTime)
	_require.NotNil(gResp.FileLastWriteTime)
	_require.NotNil(gResp.FilePermissionKey)

	dResp, err := dirClient.Delete(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(dResp.Date.IsZero(), false)
	_require.NotNil(dResp.RequestID)
	_require.NotNil(dResp.Version)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectorySetPropertiesUsingOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(cResp.FilePermissionKey)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 21:00:00 GMT 2023")
	_require.NoError(err)
	creationTime := currTime.Add(5 * time.Minute).Round(time.Microsecond)
	lastWriteTime := currTime.Add(10 * time.Minute).Round(time.Millisecond)
	changeTime := currTime.Add(15 * time.Minute).Round(time.Millisecond)

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{
				ReadOnly: true,
				System:   true,
			},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
			ChangeTime:    &changeTime,
		},
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)
	_require.NotNil(sResp.FileCreationTime)
	_require.NotNil(sResp.FileLastWriteTime)
	_require.NotNil(sResp.FileChangeTime)
	_require.NotNil(sResp.FilePermissionKey)
	_require.NotEqual(*sResp.FilePermissionKey, *cResp.FilePermissionKey)
	_require.Equal(*sResp.FileCreationTime, creationTime.UTC())
	_require.Equal(*sResp.FileLastWriteTime, lastWriteTime.UTC())
	_require.Equal(*sResp.FileChangeTime, changeTime.UTC())

	fileAttributes, err := file.ParseNTFSFileAttributes(sResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.ReadOnly)
	_require.True(fileAttributes.System)
	_require.True(fileAttributes.Directory)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.FileCreationTime)
	_require.NotNil(gResp.FileLastWriteTime)
	_require.NotNil(gResp.FileChangeTime)
	_require.NotNil(gResp.FilePermissionKey)
	_require.Equal(*gResp.FilePermissionKey, *sResp.FilePermissionKey)
	_require.Equal(*gResp.FileCreationTime, *sResp.FileCreationTime)
	_require.Equal(*gResp.FileLastWriteTime, *sResp.FileLastWriteTime)
	_require.Equal(*gResp.FileChangeTime, *sResp.FileChangeTime)
	_require.Equal(*gResp.FileAttributes, *sResp.FileAttributes)

	fileAttributes2, err := file.ParseNTFSFileAttributes(gResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.ReadOnly)
	_require.True(fileAttributes2.System)
	_require.True(fileAttributes2.Directory)
	_require.EqualValues(fileAttributes2, fileAttributes)
}

func (d *DirectoryRecordedTestsSuite) TestDirectorySetMetadataUsingOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

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

func (d *DirectoryRecordedTestsSuite) TestDirectoryListHandlesUsingOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	resp, err := dirClient.ListHandles(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Handles, 0)
	_require.NotNil(resp.NextMarker)
	_require.Equal(*resp.NextMarker, "")
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryForceCloseHandlesUsingOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	resp, err := dirClient.ForceCloseHandles(context.Background(), "*", nil)
	_require.NoError(err)
	_require.EqualValues(*resp.NumberOfHandlesClosed, 0)
	_require.EqualValues(*resp.NumberOfHandlesFailedToClose, 0)
	_require.Nil(resp.Marker)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryListUsingOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)

	// create directories
	for i := 0; i < 10; i++ {
		subDirClient := dirClient.NewSubdirectoryClient(dirName + fmt.Sprintf("%v", i))
		_, err = subDirClient.Create(context.Background(), nil)
		_require.NoError(err)
	}

	// create files
	for i := 0; i < 5; i++ {
		fileClient := dirClient.NewFileClient(fileName + fmt.Sprintf("%v", i))
		_, err = fileClient.Create(context.Background(), 2048, nil)
		_require.NoError(err)
	}

	subDirCtr, fileCtr := 0, 0
	pager := dirClient.NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		subDirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
		for _, dir := range resp.Segment.Directories {
			_require.NotNil(dir.Name)
			_require.Greater(len(*dir.Name), 0)
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
			_require.Greater(len(*f.Name), 0)
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
	_require.Equal(subDirCtr, 10)
	_require.Equal(fileCtr, 5)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	fileName := testcommon.GenerateFileName(testName)
	fClient := dirClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	resp, err := dirClient.Rename(context.Background(), "newDirName", nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)

	_, err = fClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameUsingOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup)}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := testcommon.GenerateFileName(testName)
	fClient := dirClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	resp, err := dirClient.Rename(context.Background(), "dir1", nil)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.RequestID)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)

	_, err = fClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameParentNotFound() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	_, err = dirClient.Rename(context.Background(), "dir1/dir2/", nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameDifferentDir() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcParentDirCl := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)
	srcDirCl := srcParentDirCl.NewSubdirectoryClient("subDir1")

	_, err = srcDirCl.Create(context.Background(), nil)
	_require.NoError(err)

	_ = testcommon.CreateNewDirectory(context.Background(), _require, "destDir", shareClient)

	_, err = srcDirCl.Rename(context.Background(), "destDir/subDir2", nil)
	_require.NoError(err)

	_, err = srcDirCl.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameIgnoreReadOnly() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcDirCl := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	destParentDirCl := testcommon.CreateNewDirectory(context.Background(), _require, "destDir", shareClient)

	_, err = destParentDirCl.NewFileClient("testFile").Create(context.Background(), 2048, &file.CreateOptions{
		SMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{ReadOnly: true},
		},
	})
	_require.NoError(err)

	_, err = srcDirCl.Rename(context.Background(), "destDir/testFile", nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceAlreadyExists)

	_, err = srcDirCl.Rename(context.Background(), "destDir/testFile", &directory.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
	})
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ReadOnlyAttribute)

	_, err = srcDirCl.Rename(context.Background(), "destDir/testFile", &directory.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
		IgnoreReadOnly:  to.Ptr(true),
	})
	_require.NoError(err)

	_, err = srcDirCl.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameNonDefault() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcDirCl := testcommon.CreateNewDirectory(context.Background(), _require, "dir1", shareClient)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 21:00:00 GMT 2023")
	_require.NoError(err)
	creationTime := currTime.Add(5 * time.Minute).Round(time.Microsecond)
	lastWriteTime := currTime.Add(8 * time.Minute).Round(time.Millisecond)
	changeTime := currTime.Add(10 * time.Minute).Round(time.Millisecond)

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}

	resp, err := srcDirCl.Rename(context.Background(), "dir2", &directory.RenameOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{
				ReadOnly: true,
				System:   true,
			},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
			ChangeTime:    &changeTime,
		},
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
		Metadata: md,
	})
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
	_require.True(fileAttributes.Directory)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameDestLease() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	srcDirCl := testcommon.CreateNewDirectory(context.Background(), _require, "dir1", shareClient)

	destParentDirCl := testcommon.CreateNewDirectory(context.Background(), _require, "dir2", shareClient)

	destFileClient := destParentDirCl.NewFileClient("testFile")
	_, err = destFileClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	proposedLeaseID := "c820a799-76d7-4ee2-6e15-546f19325c2c"

	// acquire lease on destFile
	fileLeaseClient, err := lease.NewFileClient(destFileClient, &lease.FileClientOptions{
		LeaseID: &proposedLeaseID,
	})
	_require.NoError(err)
	acqResp, err := fileLeaseClient.Acquire(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(acqResp.LeaseID)
	_require.Equal(*acqResp.LeaseID, proposedLeaseID)

	destPath := "dir2/testFile"
	_, err = srcDirCl.Rename(context.Background(), destPath, &directory.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
	})
	_require.Error(err)

	_, err = srcDirCl.Rename(context.Background(), destPath, &directory.RenameOptions{
		ReplaceIfExists: to.Ptr(true),
		DestinationLeaseAccessConditions: &directory.DestinationLeaseAccessConditions{
			DestinationLeaseID: acqResp.LeaseID,
		},
	})
	_require.NoError(err)

	_, err = srcDirCl.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryRenameUsingSAS() {
	_require := require.New(d.T())
	testName := d.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
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

	srcDirCl, err := directory.NewClientWithNoCredential(shareClient.URL()+"/dir1?"+sasToken, nil)
	_require.NoError(err)
	_, err = srcDirCl.Create(context.Background(), nil)
	_require.NoError(err)

	destPathWithSAS := "dir2?" + sasToken
	_, err = srcDirCl.Rename(context.Background(), destPathWithSAS, nil)
	_require.NoError(err)

	_, err = srcDirCl.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryCreateDeleteTrailingDot() {
	_require := require.New(d.T())
	testName := d.T().Name()

	options := &service.ClientOptions{AllowTrailingDot: to.Ptr(true)}
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := "dir...."
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.Delete(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestDirectorySetPropertiesTrailingDotAndOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := "dir.."
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{
		FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup),
		AllowTrailingDot:  to.Ptr(true),
	}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(cResp.FilePermissionKey)

	currTime, err := time.Parse(time.UnixDate, "Fri Mar 31 21:00:00 GMT 2023")
	_require.NoError(err)
	creationTime := currTime.Add(5 * time.Minute).Round(time.Microsecond)
	lastWriteTime := currTime.Add(10 * time.Minute).Round(time.Millisecond)
	changeTime := currTime.Add(15 * time.Minute).Round(time.Millisecond)

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(context.Background(), &directory.SetPropertiesOptions{
		FileSMBProperties: &file.SMBProperties{
			Attributes: &file.NTFSFileAttributes{
				ReadOnly: true,
				System:   true,
			},
			CreationTime:  &creationTime,
			LastWriteTime: &lastWriteTime,
			ChangeTime:    &changeTime,
		},
		FilePermissions: &file.Permissions{
			Permission: &testcommon.SampleSDDL,
		},
	})
	_require.NoError(err)
	_require.NotNil(sResp.FileCreationTime)
	_require.NotNil(sResp.FileLastWriteTime)
	_require.NotNil(sResp.FileChangeTime)
	_require.NotNil(sResp.FilePermissionKey)
	_require.NotEqual(*sResp.FilePermissionKey, *cResp.FilePermissionKey)
	_require.Equal(*sResp.FileCreationTime, creationTime.UTC())
	_require.Equal(*sResp.FileLastWriteTime, lastWriteTime.UTC())
	_require.Equal(*sResp.FileChangeTime, changeTime.UTC())

	fileAttributes, err := file.ParseNTFSFileAttributes(sResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes)
	_require.True(fileAttributes.ReadOnly)
	_require.True(fileAttributes.System)
	_require.True(fileAttributes.Directory)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.FileCreationTime)
	_require.NotNil(gResp.FileLastWriteTime)
	_require.NotNil(gResp.FileChangeTime)
	_require.NotNil(gResp.FilePermissionKey)
	_require.Equal(*gResp.FilePermissionKey, *sResp.FilePermissionKey)
	_require.Equal(*gResp.FileCreationTime, *sResp.FileCreationTime)
	_require.Equal(*gResp.FileLastWriteTime, *sResp.FileLastWriteTime)
	_require.Equal(*gResp.FileChangeTime, *sResp.FileChangeTime)
	_require.Equal(*gResp.FileAttributes, *sResp.FileAttributes)

	fileAttributes2, err := file.ParseNTFSFileAttributes(gResp.FileAttributes)
	_require.NoError(err)
	_require.NotNil(fileAttributes2)
	_require.True(fileAttributes2.ReadOnly)
	_require.True(fileAttributes2.System)
	_require.True(fileAttributes2.Directory)
	_require.EqualValues(fileAttributes2, fileAttributes)
}

func (d *DirectoryRecordedTestsSuite) TestDirectorySetMetadataTrailingDot() {
	_require := require.New(d.T())
	testName := d.T().Name()

	options := &service.ClientOptions{AllowTrailingDot: to.Ptr(true)}
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := "dir.."
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	md := map[string]*string{
		"Foo": to.Ptr("FooValuE"),
		"Bar": to.Ptr("bArvaLue"),
	}

	_, err = dirClient.SetMetadata(context.Background(), &directory.SetMetadataOptions{
		Metadata: md,
	})
	_require.NoError(err)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(gResp.Metadata, md)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryListTrailingDot() {
	_require := require.New(d.T())
	testName := d.T().Name()

	options := &service.ClientOptions{AllowTrailingDot: to.Ptr(true)}
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName) + "."
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)
	_require.NoError(err)

	// create directories
	for i := 0; i < 10; i++ {
		subDirName := fmt.Sprintf("dir %v", i)
		if i%2 == 0 {
			subDirName += ".."
		}
		subDirClient := dirClient.NewSubdirectoryClient(subDirName)
		_, err = subDirClient.Create(context.Background(), nil)
		_require.NoError(err)
	}

	// create files
	for i := 0; i < 5; i++ {
		fName := fmt.Sprintf("file %v", i)
		if i%2 == 0 {
			fName += ".."
		}
		fileClient := dirClient.NewFileClient(fName)
		_, err = fileClient.Create(context.Background(), 2048, nil)
		_require.NoError(err)
	}

	subDirCtr, subDirDotCtr, fileCtr, fileDotCtr := 0, 0, 0, 0
	pager := dirClient.NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		subDirCtr += len(resp.Segment.Directories)
		fileCtr += len(resp.Segment.Files)
		for _, dir := range resp.Segment.Directories {
			_require.NotNil(dir.Name)
			_require.Greater(len(*dir.Name), 0)
			if strings.HasSuffix(*dir.Name, "..") {
				subDirDotCtr++
			}
		}
		for _, f := range resp.Segment.Files {
			_require.NotNil(f.Name)
			_require.Greater(len(*f.Name), 0)
			if strings.HasSuffix(*f.Name, "..") {
				fileDotCtr++
			}
		}
	}
	_require.Equal(subDirCtr, 10)
	_require.Equal(subDirDotCtr, 5)
	_require.Equal(fileCtr, 5)
	_require.Equal(fileDotCtr, 3)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryRenameNegativeSourceTrailingDot() {
	_require := require.New(d.T())
	testName := d.T().Name()

	options := &service.ClientOptions{
		AllowTrailingDot: to.Ptr(true),
	}
	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, options)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName) + ".."
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	_, err = dirClient.Rename(context.Background(), "dir1..", nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)
}

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryRenameSourceTrailingDotAndOAuth() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName) + ".."
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{
		FileRequestIntent:      to.Ptr(directory.ShareTokenIntentBackup),
		AllowTrailingDot:       to.Ptr(true),
		AllowSourceTrailingDot: to.Ptr(true),
	}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClient, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	fileName := "file.."
	fClient := dirClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), 2048, nil)
	_require.NoError(err)

	resp, err := dirClient.Rename(context.Background(), "dir1..", nil)
	_require.NoError(err)
	_require.Equal(resp.LastModified.IsZero(), false)
	_require.Equal(resp.FileCreationTime.IsZero(), false)
	_require.Equal(resp.FileLastWriteTime.IsZero(), false)
	_require.Equal(resp.FileChangeTime.IsZero(), false)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ResourceNotFound)

	_, err = fClient.GetProperties(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.ParentNotFound)
}

func (d *DirectoryRecordedTestsSuite) TestListFileDirEncoded() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := "directory\uFFFF"
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	fileName := "file\uFFFE"
	fileClient := testcommon.CreateNewFileFromShare(context.Background(), _require, fileName, 2048, shareClient)

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Segment.Directories, 1)
		_require.NotNil(resp.Segment.Directories[0])
		_require.Equal(*resp.Segment.Directories[0].Name, dirName)
		_require.Len(resp.Segment.Files, 1)
		_require.NotNil(resp.Segment.Files[0])
		_require.Equal(*resp.Segment.Files[0].Name, fileName)
	}
}

func (d *DirectoryRecordedTestsSuite) TestListFileDirEncodedContinuationToken() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileName0 := "file0\uFFFE"
	_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fileName0, 2048, shareClient)

	fileName1 := "file1\uFFFE"
	_ = testcommon.CreateNewFileFromShare(context.Background(), _require, fileName1, 2048, shareClient)

	var files []string
	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		MaxResults: to.Ptr(int32(1)),
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Segment.Files, 1)
		_require.NotNil(resp.Segment.Files[0].Name)
		files = append(files, *resp.Segment.Files[0].Name)
	}

	_require.Len(files, 2)
	_require.Equal(files[0], fileName0)
	_require.Equal(files[1], fileName1)
}

func (d *DirectoryRecordedTestsSuite) TestListFileDirEncodedPrefix() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := "directory\uFFFF"
	_ = testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	pager := shareClient.NewRootDirectoryClient().NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		Prefix: &dirName,
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Len(resp.Segment.Directories, 1)
		_require.Len(resp.Segment.Files, 0)
		_require.NotNil(resp.Segment.Directories[0])
		_require.Equal(*resp.Segment.Directories[0].Name, dirName)
		_require.NotNil(resp.Prefix)
		_require.Equal(*resp.Prefix, dirName)
	}
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryClientDefaultAudience() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{
		FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup),
		Audience:          "https://storage.azure.com/",
	}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClientAudience, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClientAudience.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirectoryClientCustomAudience() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{
		FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup),
		Audience:          "https://" + accountName + ".file.core.windows.net",
	}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClientAudience, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClientAudience.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = dirClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryAudienceNegative() {
	_require := require.New(d.T())
	testName := d.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirURL := "https://" + accountName + ".file.core.windows.net/" + shareName + "/" + dirName

	options := &directory.ClientOptions{
		FileRequestIntent: to.Ptr(directory.ShareTokenIntentBackup),
		Audience:          "https://badaudience.file.core.windows.net",
	}
	testcommon.SetClientOptions(d.T(), &options.ClientOptions)
	dirClientAudience, err := directory.NewClient(dirURL, cred, options)
	_require.NoError(err)

	_, err = dirClientAudience.Create(context.Background(), nil)
	_require.Error(err)
	testcommon.ValidateFileErrorCode(_require, err, fileerror.InvalidAuthenticationInfo)
}
