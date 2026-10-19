// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
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
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running directory Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &DirectoryRecordedTestsSuite{})
		suite.Run(t, &DirectoryUnrecordedTestsSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &DirectoryRecordedTestsSuite{})
	case recording.RecordingMode:
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
		FileNFSProperties: &file.NFSProperties{
			Owner:    to.Ptr(owner),
			Group:    to.Ptr(group),
			FileMode: to.Ptr(mode),
		},
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
	_require.Equal(*gResp.NFSFileType, file.NFSFileTypeDirectory)
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
		FileNFSProperties: &file.NFSProperties{
			Owner:    to.Ptr(owner),
			Group:    to.Ptr(group),
			FileMode: to.Ptr(mode),
		},
	})
	_require.NoError(err)

	response, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*response.FileMode, mode)
	_require.Equal(*response.Group, group)
	_require.Equal(*response.Owner, owner)
	_require.Equal(*response.NFSFileType, file.NFSFileTypeDirectory)

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

func (d *DirectoryRecordedTestsSuite) TestDirGetPropertiesByFileID() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirName := testcommon.GenerateDirectoryName(testName)
	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, dirName, shareClient)

	pathProps, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(pathProps.FileName)

	dirIDClient := shareClient.NewDirectoryClientByFileID(*pathProps.ID)
	_require.Equal(*pathProps.ID, dirIDClient.FileID())

	resp, err := dirIDClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*pathProps.ID, *resp.ID)
	_require.Equal(*pathProps.ParentID, *resp.ParentID)
	_require.Equal(dirName, *resp.FileName)
	_require.Equal(*pathProps.FileAttributes, *resp.FileAttributes)

	// unsupported operations fail before a request is sent, so the directory is left untouched
	_, err = dirIDClient.Delete(context.Background(), nil)
	_require.ErrorIs(err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (d *DirectoryRecordedTestsSuite) TestDirGetPropertiesByFileIDError() {
	_require := require.New(d.T())

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(d.T().Name()), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.NewDirectoryClientByFileID("11111111111111111111").GetProperties(context.Background(), nil)
	_require.Error(err)
}

const (
	fakeFileID    = "12384898975283830"
	fakeParentID  = "13835128424026341376"
	fakeShareURL  = "https://account.file.core.windows.net/share"
	fakeFileIDURL = fakeShareURL + "?fileid=" + fakeFileID
)

// fileIDTransport records the requests it handles and answers each with the configured response.
type fileIDTransport struct {
	requests []*http.Request
	status   int
	header   http.Header
}

func (f *fileIDTransport) Do(req *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, req)
	status := f.status
	if status == 0 {
		status = http.StatusOK
	}
	header := f.header
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		Request:    req,
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     header.Clone(),
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func newFakeShareClient(t *testing.T, transport policy.Transporter) *share.Client {
	shareClient, err := share.NewClientWithNoCredential(fakeShareURL, &share.ClientOptions{ClientOptions: policy.ClientOptions{
		Transport: transport,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	}})
	require.NoError(t, err)
	return shareClient
}

func TestNewDirectoryClientByFileIDURL(t *testing.T) {
	dirClient := newFakeShareClient(t, &fileIDTransport{}).NewDirectoryClientByFileID(fakeFileID)

	require.Equal(t, fakeFileIDURL, dirClient.URL())
	require.Equal(t, fakeFileID, dirClient.FileID())

	urlParts, err := sas.ParseURL(dirClient.URL())
	require.NoError(t, err)
	require.Equal(t, "share", urlParts.ShareName)
	require.Equal(t, fakeFileID, urlParts.FileID)
	require.Empty(t, urlParts.DirectoryOrFilePath)
}

func TestDirectoryClientFileIDFromConstructorURL(t *testing.T) {
	dirClient, err := directory.NewClientWithNoCredential(fakeFileIDURL, nil)
	require.NoError(t, err)
	require.Equal(t, fakeFileID, dirClient.FileID())
}

func TestDirectoryClientFileIDEmptyForPathBasedClient(t *testing.T) {
	dirClient, err := directory.NewClientWithNoCredential(fakeShareURL+"/directory", nil)
	require.NoError(t, err)
	require.Empty(t, dirClient.FileID())

	require.Empty(t, newFakeShareClient(t, &fileIDTransport{}).NewRootDirectoryClient().FileID())
}

func TestDirectoryClientByFileIDGetProperties(t *testing.T) {
	_require := require.New(t)
	transport := &fileIDTransport{header: http.Header{
		"X-Ms-File-Id":        []string{fakeFileID},
		"X-Ms-File-Parent-Id": []string{fakeParentID},
		"X-Ms-File-Name":      []string{"directory"},
	}}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID(fakeFileID)

	resp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(fakeFileID, *resp.ID)
	_require.Equal(fakeParentID, *resp.ParentID)
	_require.Equal("directory", *resp.FileName)

	_require.Len(transport.requests, 1)
	req := transport.requests[0]
	_require.Equal("/share", req.URL.Path)
	_require.Equal(fakeFileID, req.URL.Query().Get("fileid"))
	_require.Equal("directory", req.URL.Query().Get("restype"))
}

func TestDirectoryClientByFileIDGetPropertiesError(t *testing.T) {
	transport := &fileIDTransport{
		status: http.StatusNotFound,
		header: http.Header{"X-Ms-Error-Code": []string{string(fileerror.ResourceNotFound)}},
	}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID("11111111111111111111")

	_, err := dirClient.GetProperties(context.Background(), nil)
	require.Error(t, err)
	require.True(t, fileerror.HasCode(err, fileerror.ResourceNotFound))
}

func TestDirectoryClientByFileIDUnsupportedOperations(t *testing.T) {
	ctx := context.Background()
	transport := &fileIDTransport{}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID(fakeFileID)

	operations := map[string]func() error{
		"Create":        func() error { _, err := dirClient.Create(ctx, nil); return err },
		"Delete":        func() error { _, err := dirClient.Delete(ctx, nil); return err },
		"Rename":        func() error { _, err := dirClient.Rename(ctx, "destination", nil); return err },
		"SetProperties": func() error { _, err := dirClient.SetProperties(ctx, nil); return err },
		"SetMetadata":   func() error { _, err := dirClient.SetMetadata(ctx, nil); return err },
		"ForceCloseHandles": func() error {
			_, err := dirClient.ForceCloseHandles(ctx, "*", nil)
			return err
		},
		"ListHandles": func() error { _, err := dirClient.ListHandles(ctx, nil); return err },
		"NewListFilesAndDirectoriesPager": func() error {
			_, err := dirClient.NewListFilesAndDirectoriesPager(nil).NextPage(ctx)
			return err
		},
	}

	for name, operation := range operations {
		t.Run(name, func(t *testing.T) {
			err := operation()
			require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
			require.Contains(t, err.Error(), "directory.Client."+name)
		})
	}
	require.Empty(t, transport.requests, "unsupported operations must fail before a request is sent")
}

func TestDirectoryClientByFileIDNavigation(t *testing.T) {
	// the path of a directory addressed by its file ID is not known, so clients navigated to from it are unusable
	ctx := context.Background()
	transport := &fileIDTransport{}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID(fakeFileID)

	_, err := dirClient.NewFileClient("file").GetProperties(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.NewFileClient("file").GetFileLinks(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.NewSubdirectoryClient("subdirectory").GetProperties(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.NewSubdirectoryClient("subdirectory").Create(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)

	require.Empty(t, transport.requests)
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

func newMockDirectoryClient(t *testing.T) (*directory.Client, *mock.Server) {
	srv, clientOptions := testcommon.NewMockServer(t)
	client, err := directory.NewClientWithNoCredential("https://account.file.core.windows.net/share/dir", &directory.ClientOptions{
		ClientOptions: clientOptions,
	})
	require.NoError(t, err)
	return client, srv
}

func TestListFilesAndDirectoriesPagerMocked(t *testing.T) {
	_require := require.New(t)
	client, srv := newMockDirectoryClient(t)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="share" DirectoryPath="dir">
	<Entries>
		<File><Name>file.txt</Name><Properties><Content-Length>5</Content-Length></Properties></File>
	</Entries>
	<NextMarker>marker1</NextMarker>
</EnumerationResults>`)))
	// the second page is only returned for a request that continues from the marker of the first page
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(`<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="share" DirectoryPath="dir">
	<Entries>
		<Directory><Name>subdir</Name><Properties /></Directory>
	</Entries>
	<NextMarker />
</EnumerationResults>`)), mock.WithPredicate(func(req *http.Request) bool {
		return req.URL.Query().Get("marker") == "marker1"
	}))
	srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))

	var files, dirs []string
	pager := client.NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, f := range page.Segment.Files {
			files = append(files, *f.Name)
		}
		for _, d := range page.Segment.Directories {
			dirs = append(dirs, *d.Name)
		}
	}
	_require.Equal([]string{"file.txt"}, files)
	_require.Equal([]string{"subdir"}, dirs)
	_require.Equal(2, srv.Requests())
}

func TestListFilesAndDirectoriesPagerErrorMocked(t *testing.T) {
	_require := require.New(t)
	client, srv := newMockDirectoryClient(t)
	srv.AppendResponse(mock.WithStatusCode(http.StatusNotFound),
		mock.WithHeader("x-ms-error-code", string(fileerror.ParentNotFound)),
		mock.WithBody([]byte(`<?xml version="1.0" encoding="utf-8"?><Error><Code>ParentNotFound</Code><Message>The specified parent path does not exist.</Message></Error>`)))

	_, err := client.NewListFilesAndDirectoriesPager(nil).NextPage(context.Background())
	var respErr *azcore.ResponseError
	_require.True(errors.As(err, &respErr))
	_require.Equal(http.StatusNotFound, respErr.StatusCode)
	_require.True(fileerror.HasCode(err, fileerror.ParentNotFound))
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryCreateDeleteUsingOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectorySetPropertiesUsingOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectorySetMetadataUsingOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryListHandlesUsingOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryForceCloseHandlesUsingOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryListUsingOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryRenameUsingOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectorySetPropertiesTrailingDotAndOAuth() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryClientDefaultAudience() {
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

func (d *DirectoryUnrecordedTestsSuite) TestDirectoryClientCustomAudience() {
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

func (d *DirectoryRecordedTestsSuite) TestDirCreateDirPropertySemantics() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareName := testcommon.GenerateShareName(testName)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, shareName, svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	propertySemanticsValues := []*directory.PropertySemantics{
		nil,
		to.Ptr(directory.FilePropertySemanticsNew),
		to.Ptr(directory.FilePropertySemanticsRestore),
	}

	for i, ps := range propertySemanticsValues {
		dirName := testcommon.GenerateDirectoryName(testName) + strconv.Itoa(i)
		dirClient := shareClient.NewDirectoryClient(dirName)

		opts := &directory.CreateOptions{
			FilePropertySemantics: ps,
		}

		if ps != nil && *ps == directory.FilePropertySemanticsRestore {
			opts.FilePermissions = &file.Permissions{
				Permission: to.Ptr("O:S-1-5-21-2127521184-1604012920-1887927527-21560751G:S-1-5-21-2127521184-1604012920-1887927527-513D:AI(A;;FA;;;SY)(A;;FA;;;BA)(A;;0x1200a9;;;S-1-5-21-397955417-626881126-188441444-3053964)"),
			}
		}

		cResp, err := dirClient.Create(context.Background(), opts)
		_require.NoError(err)
		_require.NotNil(cResp.ETag)
		_require.NotNil(cResp.RequestID)

		fileAttributes, err := file.ParseNTFSFileAttributes(cResp.FileAttributes)
		_require.NoError(err)
		_require.NotNil(fileAttributes)

		if ps != nil && *ps == directory.FilePropertySemanticsRestore {
			// Restore mode should not automatically add Directory attribute
			_require.False(fileAttributes.Archive)
		} else {
			// New mode (or nil/default) should behave as normal create
			_require.True(fileAttributes.Directory)
		}
	}
}

// NFS List V2 requires service version 2027-03-07. This test is unrecorded until it can be recorded against an NFS
// share in a canary environment, so it runs in live mode only.
func (d *DirectoryUnrecordedTestsSuite) TestDirListFilesAndDirsNFS() {
	_require := require.New(d.T())
	testName := d.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)
	shareName := testcommon.GenerateShareName(testName)
	shareURL := "https://" + cred.AccountName() + ".file.core.windows.net/" + shareName

	premiumShareClient, err := share.NewClientWithSharedKeyCredential(shareURL, cred, nil)
	_require.NoError(err)

	_, err = premiumShareClient.Create(context.Background(), &share.CreateOptions{
		EnabledProtocols: to.Ptr("NFS"),
	})
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, premiumShareClient)

	parentName := testcommon.GenerateDirectoryName(testName)
	parentClient := premiumShareClient.NewDirectoryClient(parentName)
	_, err = parentClient.Create(context.Background(), nil)
	_require.NoError(err)

	owner := "345"
	group := "123"
	mode := "0777"
	nfsProperties := &file.NFSProperties{Owner: to.Ptr(owner), Group: to.Ptr(group), FileMode: to.Ptr(mode)}

	fileName := testcommon.GenerateFileName(testName)
	fileClient := parentClient.NewFileClient(fileName)
	_, err = fileClient.Create(context.Background(), 1024, &file.CreateOptions{NFSProperties: nfsProperties})
	_require.NoError(err)

	subdirName := testcommon.GenerateDirectoryName("subdir")
	_, err = parentClient.NewSubdirectoryClient(subdirName).Create(context.Background(), &directory.CreateOptions{FileNFSProperties: nfsProperties})
	_require.NoError(err)

	symlinkName := testcommon.GenerateFileName("symlink")
	linkText := "/" + parentName + "/" + fileName
	_, err = parentClient.NewFileClient(symlinkName).CreateSymbolicLink(context.Background(), linkText, &file.CreateSymbolicLinkOptions{
		FileNFSProperties: &file.NFSProperties{Owner: to.Ptr(owner), Group: to.Ptr(group)},
	})
	_require.NoError(err)

	pager := parentClient.NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		Include:             directory.ListFilesInclude{Timestamps: true, ETag: true, Permissions: true, LinkCount: true, NfsAttributes: true},
		IncludeExtendedInfo: to.Ptr(true),
	})
	var files []*directory.File
	var dirs []*directory.Directory
	var symlinks []*directory.SymLinkItem
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		files = append(files, resp.Segment.Files...)
		dirs = append(dirs, resp.Segment.Directories...)
		symlinks = append(symlinks, resp.Segment.SymLinkItems...)
	}

	_require.Len(files, 1)
	f := files[0]
	_require.Equal(fileName, *f.Name)
	_require.Equal(directory.FileTypeRegular, f.FileType())
	_require.NotNil(f.LinkCount)
	_require.NotNil(f.Properties)
	_require.Equal(owner, *f.Properties.Owner)
	_require.Equal(group, *f.Properties.Group)
	_require.Equal(mode, *f.Properties.FileMode)

	_require.Len(dirs, 1)
	dir := dirs[0]
	_require.Equal(subdirName, *dir.Name)
	_require.Equal(directory.FileTypeDirectory, dir.FileType())
	_require.NotNil(dir.LinkCount)
	_require.NotNil(dir.Properties)
	_require.Equal(owner, *dir.Properties.Owner)
	_require.Equal(group, *dir.Properties.Group)
	_require.Equal(mode, *dir.Properties.FileMode)

	_require.Len(symlinks, 1)
	symlink := symlinks[0]
	_require.Equal(symlinkName, *symlink.Name.Content)
	_require.Equal(directory.FileTypeSymLink, symlink.FileType())
	_require.NotNil(symlink.LinkCount)
	_require.NotEmpty(*symlink.LinkText)
	_require.NotNil(symlink.Properties)
	_require.Equal(owner, *symlink.Properties.Owner)
	_require.Equal(group, *symlink.Properties.Group)
	_require.NotNil(symlink.Properties.FileMode)
}

func (d *DirectoryRecordedTestsSuite) TestDirListFilesAndDirsIncludeAllSMB() {
	_require := require.New(d.T())
	testName := d.T().Name()

	svcClient, err := testcommon.GetServiceClient(d.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	parentClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	fileName := testcommon.GenerateFileName(testName)
	_, err = parentClient.NewFileClient(fileName).Create(context.Background(), 1024, nil)
	_require.NoError(err)

	subdirName := testcommon.GenerateDirectoryName("subdir")
	_, err = parentClient.NewSubdirectoryClient(subdirName).Create(context.Background(), nil)
	_require.NoError(err)

	// ListFilesInclude.All is sent to the service as the literal "All".
	pager := parentClient.NewListFilesAndDirectoriesPager(&directory.ListFilesAndDirectoriesOptions{
		Include:             directory.ListFilesInclude{All: true},
		IncludeExtendedInfo: to.Ptr(true),
	})
	var files []*directory.File
	var dirs []*directory.Directory
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		files = append(files, resp.Segment.Files...)
		dirs = append(dirs, resp.Segment.Directories...)
	}

	_require.Len(files, 1)
	f := files[0]
	_require.Equal(fileName, *f.Name)
	_require.Equal(directory.FileTypeRegular, f.FileType())
	_require.NotNil(f.ID)
	_require.Equal("Archive", *f.Attributes)
	_require.NotNil(f.PermissionKey)
	_require.NotNil(f.Properties.CreationTime)
	_require.NotNil(f.Properties.LastAccessTime)
	_require.NotNil(f.Properties.LastWriteTime)
	_require.NotNil(f.Properties.ChangeTime)
	_require.NotNil(f.Properties.LastModified)
	_require.NotNil(f.Properties.ETag)
	// NFS-only fields stay nil on SMB
	_require.Nil(f.LinkCount)
	_require.Nil(f.Properties.Owner)
	_require.Nil(f.Properties.Group)
	_require.Nil(f.Properties.FileMode)

	_require.Len(dirs, 1)
	dir := dirs[0]
	_require.Equal(subdirName, *dir.Name)
	_require.Equal(directory.FileTypeDirectory, dir.FileType())
	_require.NotNil(dir.ID)
	_require.Equal("Directory", *dir.Attributes)
	_require.NotNil(dir.PermissionKey)
	_require.NotNil(dir.Properties.CreationTime)
	_require.NotNil(dir.Properties.LastAccessTime)
	_require.NotNil(dir.Properties.LastWriteTime)
	_require.NotNil(dir.Properties.ChangeTime)
	_require.NotNil(dir.Properties.LastModified)
	_require.NotNil(dir.Properties.ETag)
	// NFS-only fields stay nil on SMB
	_require.Nil(dir.LinkCount)
	_require.Nil(dir.Properties.Owner)
	_require.Nil(dir.Properties.Group)
	_require.Nil(dir.Properties.FileMode)
}

// newMockListClient returns a Client whose mock server returns body as a single 200 OK List Directories and Files response.
func newMockListClient(t *testing.T, body string) (*directory.Client, *mock.Server) {
	client, srv := newMockDirectoryClient(t)
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(body)))
	return client, srv
}

func listFirstPage(t *testing.T, client *directory.Client, options *directory.ListFilesAndDirectoriesOptions) directory.ListFilesAndDirectoriesResponse {
	pager := client.NewListFilesAndDirectoriesPager(options)
	require.True(t, pager.More())
	resp, err := pager.NextPage(context.Background())
	require.NoError(t, err)
	return resp
}

func TestListFilesAndDirectoriesNFSAllItemTypesMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, nfsListFilesAndDirectoriesXML)

	resp := listFirstPage(t, client, nil)
	segment := resp.Segment
	_require.NotNil(segment)
	_require.Equal("9223442405599084545", *resp.DirectoryID)

	// 1 Directory + 4 Files + 1 SymLink + 1 BlockDevice + 1 CharDevice + 1 Fifo + 1 Socket = 10.
	_require.Len(segment.Directories, 1)
	_require.Len(segment.Files, 4)
	_require.Len(segment.SymLinkItems, 1)
	_require.Len(segment.BlockDeviceItems, 1)
	_require.Len(segment.CharDeviceItems, 1)
	_require.Len(segment.FifoItems, 1)
	_require.Len(segment.SocketItems, 1)

	// Directory
	dir := segment.Directories[0]
	_require.Equal(directory.FileTypeDirectory, dir.FileType())
	_require.Equal("subdir", *dir.Name)
	_require.Equal("12682206919419625485", *dir.ID)
	_require.Equal(int64(2), *dir.LinkCount)
	_require.Nil(dir.Properties.ContentLength)
	_require.Equal("1000", *dir.Properties.Owner)
	_require.Equal("1000", *dir.Properties.Group)
	_require.Equal("0755", *dir.Properties.FileMode)
	_require.Equal(azcore.ETag(`"0x8DEAF1479E1C087"`), *dir.Properties.ETag)
	_require.NotNil(dir.Properties.CreationTime)
	_require.NotNil(dir.Properties.LastModified)

	// Files
	expectedFiles := []struct {
		name      string
		id        string
		size      int64
		linkCount int64
	}{
		{"hardlink.txt", "13835128424026472451", 80, 2},
		{"regular.txt", "13835128424026472451", 80, 2},
		{"regularClose.txt", "9799903157902508049", 14, 1},
		{"regularOpen.txt", "17293892937847013391", 14, 1},
	}
	for i, expected := range expectedFiles {
		f := segment.Files[i]
		_require.Equal(directory.FileTypeRegular, f.FileType())
		_require.Equal(expected.name, *f.Name)
		_require.Equal(expected.id, *f.ID)
		_require.Equal(expected.size, *f.Properties.ContentLength)
		_require.Equal(expected.linkCount, *f.LinkCount)
		_require.Equal("1000", *f.Properties.Owner)
		_require.Equal("1000", *f.Properties.Group)
		_require.Equal("0644", *f.Properties.FileMode)
		_require.Nil(f.Attributes)
		_require.Nil(f.PermissionKey)
	}

	// SymLink
	symlink := segment.SymLinkItems[0]
	_require.Equal(directory.FileTypeSymLink, symlink.FileType())
	_require.Equal("symlink.txt", *symlink.Name.Content)
	_require.Nil(symlink.Name.Encoded)
	_require.Equal("10376363910205931529", *symlink.FileID)
	_require.Equal(int64(1), *symlink.LinkCount)
	_require.Equal("/mnt/s2/dir2/regular.txt", *symlink.LinkText)
	_require.Equal("1000", *symlink.Properties.Owner)
	_require.Equal("1000", *symlink.Properties.Group)
	_require.Equal("0777", *symlink.Properties.FileMode)

	// BlockDevice
	block := segment.BlockDeviceItems[0]
	_require.Equal(directory.FileTypeBlockDevice, block.FileType())
	_require.Equal("block_device", *block.Name.Content)
	_require.Equal("10952824662509355033", *block.FileID)
	_require.Equal(int64(1), *block.LinkCount)
	_require.Equal(int64(8), *block.DeviceMajor)
	_require.Equal(int64(0), *block.DeviceMinor)
	_require.Equal("0", *block.Properties.Owner)
	_require.Equal("0", *block.Properties.Group)
	_require.Equal("0640", *block.Properties.FileMode)

	// CharDevice
	charDev := segment.CharDeviceItems[0]
	_require.Equal(directory.FileTypeCharacterDevice, charDev.FileType())
	_require.Equal("char_device", *charDev.Name.Content)
	_require.Equal("16717432185543589911", *charDev.FileID)
	_require.Equal(int64(1), *charDev.LinkCount)
	_require.Equal(int64(1), *charDev.DeviceMajor)
	_require.Equal(int64(7), *charDev.DeviceMinor)
	_require.Equal("0", *charDev.Properties.Owner)
	_require.Equal("0", *charDev.Properties.Group)
	_require.Equal("0644", *charDev.Properties.FileMode)

	// Fifo
	fifo := segment.FifoItems[0]
	_require.Equal(directory.FileTypeFifo, fifo.FileType())
	_require.Equal("fifo_pipe", *fifo.Name.Content)
	_require.Equal("14988049928633319435", *fifo.FileID)
	_require.Equal(int64(1), *fifo.LinkCount)
	_require.Equal("1000", *fifo.Properties.Owner)
	_require.Equal("1000", *fifo.Properties.Group)
	_require.Equal("0644", *fifo.Properties.FileMode)

	// Socket
	socket := segment.SocketItems[0]
	_require.Equal(directory.FileTypeSocket, socket.FileType())
	_require.Equal("unix_socket", *socket.Name.Content)
	_require.Equal("16429201809391878183", *socket.FileID)
	_require.Equal(int64(1), *socket.LinkCount)
	_require.Equal("0", *socket.Properties.Owner)
	_require.Equal("0", *socket.Properties.Group)
	_require.Equal("0755", *socket.Properties.FileMode)
}

func TestListFilesAndDirectoriesNFSEncodedNameMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="nfsshare" DirectoryPath="dir" Encoded="true">
	<Entries>
		<BlockDevice>
			<Name Encoded="true">block%EF%BF%BF</Name>
			<Properties />
		</BlockDevice>
		<CharDevice>
			<Name Encoded="true">char%EF%BF%BF</Name>
			<Properties />
		</CharDevice>
		<Fifo>
			<Name Encoded="true">fifo%EF%BF%BF</Name>
			<Properties />
		</Fifo>
		<SymLink>
			<Name Encoded="true">link%EF%BF%BF.txt</Name>
			<Properties />
			<LinkText>target.txt</LinkText>
		</SymLink>
		<Socket>
			<Name Encoded="false">socket%20name</Name>
			<Properties />
		</Socket>
	</Entries>
	<NextMarker />
</EnumerationResults>`)

	resp := listFirstPage(t, client, nil)
	segment := resp.Segment

	assertDecoded := func(name *directory.StringEncoded, expected string) {
		_require.Equal(expected, *name.Content)
		_require.False(*name.Encoded)
	}
	assertDecoded(segment.BlockDeviceItems[0].Name, "block￿")
	assertDecoded(segment.CharDeviceItems[0].Name, "char￿")
	assertDecoded(segment.FifoItems[0].Name, "fifo￿")
	assertDecoded(segment.SymLinkItems[0].Name, "link￿.txt")
	_require.Equal("target.txt", *segment.SymLinkItems[0].LinkText)
	// A name that is not encoded is returned unchanged.
	assertDecoded(segment.SocketItems[0].Name, "socket%20name")
}

func TestListFilesAndDirectoriesNFSMalformedEncodedNameMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="nfsshare" DirectoryPath="dir" Encoded="true">
	<Entries>
		<SymLink>
			<Name Encoded="true">bad%zzname</Name>
			<Properties />
			<LinkText>target.txt</LinkText>
		</SymLink>
		<Fifo>
			<Name Encoded="true">fifo%EF%BF%BF</Name>
			<Properties />
		</Fifo>
	</Entries>
	<NextMarker />
</EnumerationResults>`)

	// A malformed encoding does not fail the page; the name is returned as sent with Encoded still true.
	resp := listFirstPage(t, client, nil)
	symlinkName := resp.Segment.SymLinkItems[0].Name
	_require.Equal("bad%zzname", *symlinkName.Content)
	_require.True(*symlinkName.Encoded)
	_require.Equal("target.txt", *resp.Segment.SymLinkItems[0].LinkText)

	// Other items on the page are still decoded.
	fifoName := resp.Segment.FifoItems[0].Name
	_require.Equal("fifo￿", *fifoName.Content)
	_require.False(*fifoName.Encoded)
}

func TestListFilesAndDirectoriesSMBMalformedEncodedNameMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="smbshare" DirectoryPath="dir">
	<Prefix Encoded="true">bad%zzprefix</Prefix>
	<Entries>
		<File>
			<Name Encoded="true">bad%zzfile</Name>
			<Properties><Content-Length>0</Content-Length></Properties>
		</File>
		<Directory>
			<Name Encoded="true">bad%zzdir</Name>
			<Properties />
		</Directory>
		<File>
			<Name Encoded="true">file%EF%BF%BF</Name>
			<Properties><Content-Length>0</Content-Length></Properties>
		</File>
	</Entries>
	<NextMarker />
</EnumerationResults>`)

	// A malformed encoding does not fail the page; the value is returned as sent.
	resp := listFirstPage(t, client, nil)
	_require.Equal("bad%zzprefix", *resp.Prefix)
	_require.Len(resp.Segment.Files, 2)
	_require.Equal("bad%zzfile", *resp.Segment.Files[0].Name)
	_require.Len(resp.Segment.Directories, 1)
	_require.Equal("bad%zzdir", *resp.Segment.Directories[0].Name)

	// Other items on the page are still decoded.
	_require.Equal("file￿", *resp.Segment.Files[1].Name)
}

func TestListFilesAndDirectoriesSMBMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockListClient(t, `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="smbshare" DirectoryPath="dir">
	<Entries>
		<File>
			<Name>file.txt</Name>
			<FileId>13835128424026341376</FileId>
			<Properties>
				<Content-Length>1024</Content-Length>
				<CreationTime>2026-05-11T04:19:21.0288880Z</CreationTime>
				<Etag>"0x8DEAF15A711D34B"</Etag>
			</Properties>
			<Attributes>Archive</Attributes>
			<PermissionKey>4066528134148476695*1</PermissionKey>
		</File>
		<Directory>
			<Name>subdir</Name>
			<FileId>13835093239654252544</FileId>
			<Properties />
			<Attributes>Directory</Attributes>
			<PermissionKey>4066528134148476695*1</PermissionKey>
		</Directory>
	</Entries>
	<NextMarker />
</EnumerationResults>`)

	resp := listFirstPage(t, client, nil)
	segment := resp.Segment
	_require.Len(segment.Files, 1)
	_require.Len(segment.Directories, 1)
	_require.Empty(segment.SymLinkItems)
	_require.Empty(segment.BlockDeviceItems)
	_require.Empty(segment.CharDeviceItems)
	_require.Empty(segment.FifoItems)
	_require.Empty(segment.SocketItems)

	f := segment.Files[0]
	_require.Equal(directory.FileTypeRegular, f.FileType())
	_require.Equal("file.txt", *f.Name)
	_require.Equal("Archive", *f.Attributes)
	_require.Equal("4066528134148476695*1", *f.PermissionKey)
	_require.Equal(int64(1024), *f.Properties.ContentLength)
	// NFS-only fields stay nil on SMB
	_require.Nil(f.LinkCount)
	_require.Nil(f.Properties.Owner)
	_require.Nil(f.Properties.Group)
	_require.Nil(f.Properties.FileMode)

	dir := segment.Directories[0]
	_require.Equal(directory.FileTypeDirectory, dir.FileType())
	_require.Equal("subdir", *dir.Name)
	_require.Equal("Directory", *dir.Attributes)
	_require.Nil(dir.LinkCount)
	_require.Nil(dir.Properties.Owner)
	_require.Nil(dir.Properties.Group)
	_require.Nil(dir.Properties.FileMode)
}

func TestListFilesAndDirectoriesIncludeQueryMocked(t *testing.T) {
	testCases := []struct {
		name     string
		include  directory.ListFilesInclude
		expected string
	}{
		{"None", directory.ListFilesInclude{}, ""},
		{"SMB", directory.ListFilesInclude{Timestamps: true, ETag: true, Attributes: true, PermissionKey: true}, "Timestamps,Etag,Attributes,PermissionKey"},
		{"NFS", directory.ListFilesInclude{Timestamps: true, ETag: true, Permissions: true, LinkCount: true, NfsAttributes: true}, "Timestamps,Etag,Permissions,LinkCount,NfsAttributes"},
		{"All", directory.ListFilesInclude{All: true}, "All"},
		// All takes precedence and is sent as the literal "All" rather than being expanded.
		{"AllWithOthers", directory.ListFilesInclude{All: true, Timestamps: true, LinkCount: true}, "All"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, srv := newMockDirectoryClient(t)
			// the list response is only returned when the include parameter is sent as expected
			srv.AppendResponse(mock.WithStatusCode(http.StatusOK), mock.WithBody([]byte(emptyListFilesAndDirectoriesXML)),
				mock.WithPredicate(func(req *http.Request) bool {
					_, ok := req.URL.Query()["include"]
					return req.URL.Query().Get("include") == tc.expected && ok == (tc.expected != "")
				}))
			srv.AppendResponse(mock.WithStatusCode(http.StatusBadRequest))
			_ = listFirstPage(t, client, &directory.ListFilesAndDirectoriesOptions{Include: tc.include})
			require.Equal(t, 1, srv.Requests())
		})
	}
}

func TestFileTypeValues(t *testing.T) {
	require.Equal(t, []directory.FileType{
		directory.FileTypeRegular,
		directory.FileTypeDirectory,
		directory.FileTypeSymLink,
		directory.FileTypeBlockDevice,
		directory.FileTypeCharacterDevice,
		directory.FileTypeSocket,
		directory.FileTypeFifo,
	}, directory.PossibleFileTypeValues())
}

const emptyListFilesAndDirectoriesXML = `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="myshare" DirectoryPath="mydir">
	<Entries />
	<NextMarker />
</EnumerationResults>`

// nfsListFilesAndDirectoriesXML is a List Directories and Files response containing every NFS item type.
const nfsListFilesAndDirectoriesXML = `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="http://localhost:10004/testaccount1/" ShareName="nfsshare" DirectoryPath="dir2">
	<DirectoryId>9223442405599084545</DirectoryId>
	<Entries>
		<BlockDevice>
			<Name>block_device</Name>
			<FileId>10952824662509355033</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:33.1472879Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:33.1472879Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:33.1472879Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:33.1472879Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:33 GMT</Last-Modified>
				<Etag>"0x8DEAF1480FD31EF"</Etag>
				<Uid>0</Uid>
				<Gid>0</Gid>
				<Mode>0640</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
			<DeviceMajor>8</DeviceMajor>
			<DeviceMinor>0</DeviceMinor>
		</BlockDevice>
		<CharDevice>
			<Name>char_device</Name>
			<FileId>16717432185543589911</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:33.1120018Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:33.1120018Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:33.1120018Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:33.1120018Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:33 GMT</Last-Modified>
				<Etag>"0x8DEAF1480F7CF92"</Etag>
				<Uid>0</Uid>
				<Gid>0</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
			<DeviceMajor>1</DeviceMajor>
			<DeviceMinor>7</DeviceMinor>
		</CharDevice>
		<Fifo>
			<Name>fifo_pipe</Name>
			<FileId>14988049928633319435</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:21.1982170Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.1982170Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:21.1982170Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:21.1982170Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:21 GMT</Last-Modified>
				<Etag>"0x8DEAF1479DDE95A"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</Fifo>
		<File>
			<Name>hardlink.txt</Name>
			<FileId>13835128424026472451</FileId>
			<Properties>
				<Content-Length>80</Content-Length>
				<CreationTime>2026-05-11T04:19:21.0288880Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.0288880Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.5326411Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.5326411Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A711D34B"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>2</LinkCount>
		</File>
		<File>
			<Name>regular.txt</Name>
			<FileId>13835128424026472451</FileId>
			<Properties>
				<Content-Length>80</Content-Length>
				<CreationTime>2026-05-11T04:19:21.0288880Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.0288880Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.5326411Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.5326411Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A711D34B"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>2</LinkCount>
		</File>
		<File>
			<Name>regularClose.txt</Name>
			<FileId>9799903157902508049</FileId>
			<Properties>
				<Content-Length>14</Content-Length>
				<CreationTime>2026-05-11T04:19:21.2920068Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.2920068Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.6790214Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.6790214Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A7282946"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</File>
		<File>
			<Name>regularOpen.txt</Name>
			<FileId>17293892937847013391</FileId>
			<Properties>
				<Content-Length>14</Content-Length>
				<CreationTime>2026-05-11T04:19:21.2612277Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.2612277Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:51.5075180Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:51.5075180Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:51 GMT</Last-Modified>
				<Etag>"0x8DEAF15AA08EE6C"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0644</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</File>
		<Directory>
			<Name>subdir</Name>
			<FileId>12682206919419625485</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:19:21.2233863Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.2233863Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:21.2233863Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:21.2233863Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:21 GMT</Last-Modified>
				<Etag>"0x8DEAF1479E1C087"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0755</Mode>
			</Properties>
			<LinkCount>2</LinkCount>
		</Directory>
		<SymLink>
			<Name>symlink.txt</Name>
			<FileId>10376363910205931529</FileId>
			<Properties>
				<Content-Length>24</Content-Length>
				<CreationTime>2026-05-11T04:19:21.1661608Z</CreationTime>
				<LastAccessTime>2026-05-11T04:19:21.1661608Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:19:21.1661608Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:19:21.1661608Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:19:21 GMT</Last-Modified>
				<Etag>"0x8DEAF1479D90528"</Etag>
				<Uid>1000</Uid>
				<Gid>1000</Gid>
				<Mode>0777</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
			<LinkText>/mnt/s2/dir2/regular.txt</LinkText>
		</SymLink>
		<Socket>
			<Name>unix_socket</Name>
			<FileId>16429201809391878183</FileId>
			<Properties>
				<CreationTime>2026-05-11T04:27:46.7308143Z</CreationTime>
				<LastAccessTime>2026-05-11T04:27:46.7308143Z</LastAccessTime>
				<LastWriteTime>2026-05-11T04:27:46.7308143Z</LastWriteTime>
				<ChangeTime>2026-05-11T04:27:46.7308143Z</ChangeTime>
				<Last-Modified>Mon, 11 May 2026 04:27:46 GMT</Last-Modified>
				<Etag>"0x8DEAF15A730106F"</Etag>
				<Uid>0</Uid>
				<Gid>0</Gid>
				<Mode>0755</Mode>
			</Properties>
			<LinkCount>1</LinkCount>
		</Socket>
	</Entries>
	<NextMarker />
</EnumerationResults>`
