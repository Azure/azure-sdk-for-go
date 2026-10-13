// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

// The STG105 file ID integration tests require service version 2027-03-07. TestFileClientByFileIDSASInURL,
// TestFileClientByFileIDOAuth and TestFileGetFileLinksNFSHardLinks run only in live mode.

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// createFileForFileIDTest creates a directory and a 1 KiB file in a new share, and returns the share client, the
// directory client, the path based file client and the properties of the file. The caller deletes the share.
func createFileForFileIDTest(s *suite.Suite) (*share.Client, *directory.Client, *file.Client, file.GetPropertiesResponse) {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)
	fClient := dirClient.NewFileClient(testcommon.GenerateFileName(testName))
	_, err = fClient.Create(context.Background(), 1024, nil)
	_require.NoError(err)

	props, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(props.ID)
	_require.NotNil(props.ParentID)
	return shareClient, dirClient, fClient, props
}

func fileNameOf(_require *require.Assertions, fClient *file.Client) string {
	urlParts, err := file.ParseURL(fClient.URL())
	_require.NoError(err)
	return path.Base(urlParts.DirectoryOrFilePath)
}

func (f *FileRecordedTestsSuite) TestFileGetPropertiesByFileID() {
	_require := require.New(f.T())
	shareClient, _, fClient, pathProps := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	fileIDClient := shareClient.NewFileClientByFileID(*pathProps.ID)
	_require.Equal(*pathProps.ID, fileIDClient.FileID())

	resp, err := fileIDClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*pathProps.ID, *resp.ID)
	_require.Equal(*pathProps.ParentID, *resp.ParentID)
	_require.Equal(fileNameOf(_require, fClient), *resp.FileName)
	_require.Equal(*pathProps.ContentLength, *resp.ContentLength)

	// the name of the file is only returned when the file is addressed by its file ID
	_require.Nil(pathProps.FileName)
}

func (f *FileRecordedTestsSuite) TestFileGetPropertiesByFileIDError() {
	_require := require.New(f.T())

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(f.T().Name()), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.NewFileClientByFileID("11111111111111111111").GetProperties(context.Background(), nil)
	_require.Error(err)
}

func (f *FileRecordedTestsSuite) TestFileGetFileLinks() {
	_require := require.New(f.T())
	shareClient, _, fClient, pathProps := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	resp, err := shareClient.NewFileClientByFileID(*pathProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)

	_require.Equal(*pathProps.ID, *resp.ID)
	_require.Equal(*pathProps.ParentID, *resp.ParentID)
	_require.Equal(int64(1024), *resp.FileContentLength)
	_require.Equal("application/octet-stream", *resp.FileContentType)
	_require.Equal(*pathProps.FileAttributes, *resp.FileAttributes)
	_require.Equal(*pathProps.FilePermissionKey, *resp.FilePermissionKey)
	_require.Equal(*pathProps.FileCreationTime, *resp.FileCreationTime)
	_require.Equal(*pathProps.FileLastWriteTime, *resp.FileLastWriteTime)
	_require.Equal(*pathProps.FileChangeTime, *resp.FileChangeTime)

	_require.Len(resp.HardLinks, 1)
	_require.Equal(fileNameOf(_require, fClient), *resp.HardLinks[0].FileName)
	_require.Equal(*pathProps.ParentID, *resp.HardLinks[0].ParentID)
}

func (f *FileRecordedTestsSuite) TestFileGetFileLinksRootDirectory() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	rootDirClient := shareClient.NewRootDirectoryClient()
	fileName := testcommon.GenerateFileName(testName)
	fClient := rootDirClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), 1024, nil)
	_require.NoError(err)

	rootProps, err := rootDirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	fileProps, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	resp, err := shareClient.NewFileClientByFileID(*fileProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.HardLinks, 1)
	_require.Equal(fileName, *resp.HardLinks[0].FileName)
	_require.Equal(*rootProps.ID, *resp.HardLinks[0].ParentID)
}

func (f *FileRecordedTestsSuite) TestFileGetFileLinksContentHeaders() {
	_require := require.New(f.T())
	shareClient, _, fClient, _ := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	md5 := []byte("0123456789abcdef")
	httpHeaders := file.HTTPHeaders{
		ContentType:        to.Ptr("text/html"),
		ContentEncoding:    to.Ptr("gzip"),
		ContentLanguage:    to.Ptr("en-US"),
		ContentDisposition: to.Ptr("attachment"),
		CacheControl:       to.Ptr("no-transform"),
		ContentMD5:         md5,
	}
	_, err := fClient.SetHTTPHeaders(context.Background(), &file.SetHTTPHeadersOptions{HTTPHeaders: &httpHeaders})
	_require.NoError(err)

	pathProps, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	resp, err := shareClient.NewFileClientByFileID(*pathProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)

	// content properties are returned with an x-ms- prefix, since the response has a body
	_require.Equal(*httpHeaders.ContentType, *resp.FileContentType)
	_require.Equal(*httpHeaders.ContentEncoding, *resp.FileContentEncoding)
	_require.Equal(*httpHeaders.ContentLanguage, *resp.FileContentLanguage)
	_require.Equal(*httpHeaders.ContentDisposition, *resp.FileContentDisposition)
	_require.Equal(*httpHeaders.CacheControl, *resp.FileCacheControl)
	_require.Equal(md5, resp.FileContentMD5)

	_require.Equal(*pathProps.FileAttributes, *resp.FileAttributes)
	_require.Equal(*pathProps.FilePermissionKey, *resp.FilePermissionKey)
	_require.Equal(*pathProps.FileCreationTime, *resp.FileCreationTime)
	_require.Equal(*pathProps.FileLastWriteTime, *resp.FileLastWriteTime)
	_require.Equal(*pathProps.FileChangeTime, *resp.FileChangeTime)
}

func (f *FileRecordedTestsSuite) TestFileGetFileLinksEncoded() {
	_require := require.New(f.T())
	testName := f.T().Name()

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(testName), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirClient := testcommon.CreateNewDirectory(context.Background(), _require, testcommon.GenerateDirectoryName(testName), shareClient)

	// ￾ cannot be represented in XML, so the service percent-encodes the FileName element and marks it with
	// Encoded="true"
	fileName := testcommon.GenerateFileName(testName) + "￾"
	fClient := dirClient.NewFileClient(fileName)
	_, err = fClient.Create(context.Background(), 1024, nil)
	_require.NoError(err)

	pathProps, err := fClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	resp, err := shareClient.NewFileClientByFileID(*pathProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.HardLinks, 1)
	_require.Equal(fileName, *resp.HardLinks[0].FileName)
	_require.Equal(*pathProps.ParentID, *resp.HardLinks[0].ParentID)
}

func (f *FileRecordedTestsSuite) TestFileGetFileLinksShareSnapshot() {
	_require := require.New(f.T())
	shareClient, _, fClient, pathProps := createFileForFileIDTest(&f.Suite)

	snapshotResp, err := shareClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	defer func() {
		_, _ = shareClient.Delete(context.Background(), &share.DeleteOptions{DeleteSnapshots: to.Ptr(share.DeleteSnapshotsOptionTypeInclude)})
	}()

	snapshotShareClient, err := shareClient.WithSnapshot(*snapshotResp.Snapshot)
	_require.NoError(err)

	resp, err := snapshotShareClient.NewFileClientByFileID(*pathProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.HardLinks, 1)
	_require.Equal(fileNameOf(_require, fClient), *resp.HardLinks[0].FileName)

	resp, err = shareClient.NewFileClientByFileID(*pathProps.ID).GetFileLinks(context.Background(), &file.GetFileLinksOptions{
		ShareSnapshot: snapshotResp.Snapshot,
	})
	_require.NoError(err)
	_require.Len(resp.HardLinks, 1)
}

func (f *FileRecordedTestsSuite) TestFileGetFileLinksLease() {
	_require := require.New(f.T())
	shareClient, _, fClient, pathProps := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	proposedLeaseID := to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c")
	leaseClient, err := lease.NewFileClient(fClient, &lease.FileClientOptions{LeaseID: proposedLeaseID})
	_require.NoError(err)
	_, err = leaseClient.Acquire(context.Background(), nil)
	_require.NoError(err)
	defer func() { _, _ = leaseClient.Release(context.Background(), nil) }()

	resp, err := shareClient.NewFileClientByFileID(*pathProps.ID).GetFileLinks(context.Background(), &file.GetFileLinksOptions{
		LeaseAccessConditions: &file.LeaseAccessConditions{LeaseID: proposedLeaseID},
	})
	_require.NoError(err)
	_require.Len(resp.HardLinks, 1)
	_require.Equal(lease.StateTypeLeased, *resp.LeaseState)
}

func (f *FileRecordedTestsSuite) TestFileGetFileLinksError() {
	_require := require.New(f.T())

	svcClient, err := testcommon.GetServiceClient(f.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	shareClient := testcommon.CreateNewShare(context.Background(), _require, testcommon.GenerateShareName(f.T().Name()), svcClient)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	_, err = shareClient.NewFileClientByFileID("11111111111111111111").GetFileLinks(context.Background(), nil)
	_require.Error(err)
}

func (f *FileUnrecordedTestsSuite) TestFileGetFileLinksNFSHardLinks() {
	_require := require.New(f.T())
	testName := f.T().Name()

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountPremium)
	_require.NoError(err)

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	shareClient, err := share.NewClientWithSharedKeyCredential(
		"https://"+cred.AccountName()+".file.core.windows.net/"+testcommon.GenerateShareName(testName), cred, options)
	_require.NoError(err)
	_, err = shareClient.Create(context.Background(), &share.CreateOptions{EnabledProtocols: to.Ptr("NFS")})
	_require.NoError(err)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	directoryName := testcommon.GenerateDirectoryName(testName)
	dirClient := shareClient.NewRootDirectoryClient().NewSubdirectoryClient(directoryName)
	_, err = dirClient.Create(context.Background(), nil)
	_require.NoError(err)

	sourceFileName := testcommon.GenerateFileName("file1")
	_, err = dirClient.NewFileClient(sourceFileName).Create(context.Background(), 1024, nil)
	_require.NoError(err)

	hardLinkFileName := testcommon.GenerateFileName("file2")
	_, err = dirClient.NewFileClient(hardLinkFileName).CreateHardLink(context.Background(), fmt.Sprintf("/%s/%s", directoryName, sourceFileName), nil)
	_require.NoError(err)

	dirProps, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	fileProps, err := dirClient.NewFileClient(sourceFileName).GetProperties(context.Background(), nil)
	_require.NoError(err)

	resp, err := shareClient.NewFileClientByFileID(*fileProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(int64(2), *resp.LinkCount)
	_require.Len(resp.HardLinks, 2)

	names := make([]string, 0, len(resp.HardLinks))
	for _, link := range resp.HardLinks {
		names = append(names, *link.FileName)
		_require.Equal(*dirProps.ID, *link.ParentID)
	}
	_require.ElementsMatch([]string{sourceFileName, hardLinkFileName}, names)
}

func (f *FileRecordedTestsSuite) TestFileClientByFileIDSharedKeyConstructor() {
	_require := require.New(f.T())
	shareClient, dirClient, fClient, fileProps := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirProps, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	urlParts, err := sas.ParseURL(shareClient.URL())
	_require.NoError(err)

	urlParts.FileID = *fileProps.ID
	fileOptions := &file.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &fileOptions.ClientOptions)
	fileIDClient, err := file.NewClientWithSharedKeyCredential(urlParts.String(), cred, fileOptions)
	_require.NoError(err)

	urlParts.FileID = *dirProps.ID
	dirOptions := &directory.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &dirOptions.ClientOptions)
	dirIDClient, err := directory.NewClientWithSharedKeyCredential(urlParts.String(), cred, dirOptions)
	_require.NoError(err)

	resp, err := fileIDClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*fileProps.ID, *resp.ID)
	_require.Equal(fileNameOf(_require, fClient), *resp.FileName)

	dirResp, err := dirIDClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*dirProps.ID, *dirResp.ID)
	_require.Equal(path.Base(mustDirectoryPath(_require, dirClient)), *dirResp.FileName)
}

func (f *FileRecordedTestsSuite) TestFileClientByFileIDConnectionString() {
	_require := require.New(f.T())
	shareClient, dirClient, fClient, fileProps := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirProps, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	connString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)
	urlParts, err := sas.ParseURL(shareClient.URL())
	_require.NoError(err)

	options := &share.ClientOptions{}
	testcommon.SetClientOptions(f.T(), &options.ClientOptions)
	connShareClient, err := share.NewClientFromConnectionString(*connString, urlParts.ShareName, options)
	_require.NoError(err)

	resp, err := connShareClient.NewFileClientByFileID(*fileProps.ID).GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*fileProps.ID, *resp.ID)
	_require.Equal(fileNameOf(_require, fClient), *resp.FileName)

	dirResp, err := connShareClient.NewDirectoryClientByFileID(*dirProps.ID).GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*dirProps.ID, *dirResp.ID)
	_require.Equal(path.Base(mustDirectoryPath(_require, dirClient)), *dirResp.FileName)
}

func mustDirectoryPath(_require *require.Assertions, dirClient *directory.Client) string {
	urlParts, err := sas.ParseURL(dirClient.URL())
	_require.NoError(err)
	return urlParts.DirectoryOrFilePath
}

func (f *FileUnrecordedTestsSuite) TestFileClientByFileIDSASInURL() {
	_require := require.New(f.T())
	shareClient, _, fClient, fileProps := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	// a file ID addressed request has no path to sign, so a share SAS is used
	shareSASURL, err := shareClient.GetSASURL(sas.SharePermissions{Read: true}, time.Now().Add(time.Hour), nil)
	_require.NoError(err)

	urlParts, err := sas.ParseURL(shareSASURL)
	_require.NoError(err)
	urlParts.FileID = *fileProps.ID

	fileIDClient, err := file.NewClientWithNoCredential(urlParts.String(), nil)
	_require.NoError(err)
	_require.Equal(*fileProps.ID, fileIDClient.FileID())

	resp, err := fileIDClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*fileProps.ID, *resp.ID)
	_require.Equal(fileNameOf(_require, fClient), *resp.FileName)

	// navigation from a share client created from a SAS URL keeps the SAS
	sasShareClient, err := share.NewClientWithNoCredential(shareSASURL, nil)
	_require.NoError(err)
	links, err := sasShareClient.NewFileClientByFileID(*fileProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)
	_require.Len(links.HardLinks, 1)
}

func (f *FileUnrecordedTestsSuite) TestFileClientByFileIDOAuth() {
	_require := require.New(f.T())
	shareClient, dirClient, fClient, fileProps := createFileForFileIDTest(&f.Suite)
	defer testcommon.DeleteShare(context.Background(), _require, shareClient)

	dirProps, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	oauthShareClient, err := share.NewClient(shareClient.URL(), cred, &share.ClientOptions{FileRequestIntent: to.Ptr(share.TokenIntentBackup)})
	_require.NoError(err)

	resp, err := oauthShareClient.NewFileClientByFileID(*fileProps.ID).GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*fileProps.ID, *resp.ID)
	_require.Equal(fileNameOf(_require, fClient), *resp.FileName)

	links, err := oauthShareClient.NewFileClientByFileID(*fileProps.ID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)
	_require.Len(links.HardLinks, 1)
	_require.Equal(fileNameOf(_require, fClient), *links.HardLinks[0].FileName)

	dirResp, err := oauthShareClient.NewDirectoryClientByFileID(*dirProps.ID).GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*dirProps.ID, *dirResp.ID)

	// the file ID client still rejects unsupported operations before sending a request
	_, err = oauthShareClient.NewFileClientByFileID(*fileProps.ID).Delete(context.Background(), nil)
	_require.ErrorIs(err, fileerror.ErrFileIDNotSupported)
}
