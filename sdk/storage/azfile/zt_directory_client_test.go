//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
)

func (s *azfileLiveTestSuite) TestDirNewDirectoryClient() {
	_require := require.New(s.T())
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient, err := svcClient.NewShareClient(sharePrefix)
	_require.Nil(err)
	dirClient, err := srClient.NewDirectoryClient(directoryPrefix)
	_require.Nil(err)
	dirClient, err = dirClient.NewDirectoryClient("inner" + directoryPrefix)
	_require.Nil(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/" + sharePrefix + "/" + directoryPrefix + "/inner" + directoryPrefix
	_require.Equal(dirClient.URL(), correctURL)
}

func (s *azfileLiveTestSuite) TestDirCreateFileURL() {
	_require := require.New(s.T())
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient, err := svcClient.NewShareClient(sharePrefix)
	_require.Nil(err)
	dirClient, err := srClient.NewDirectoryClient(directoryPrefix)
	_require.Nil(err)
	fClient, err := dirClient.NewFileClient(filePrefix)
	_require.Nil(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/" + sharePrefix + "/" + directoryPrefix + "/" + filePrefix
	_require.Equal(fClient.URL(), correctURL)
}

func (s *azfileLiveTestSuite) TestDirCreateDeleteDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)

	defer delShare(_require, srClient, nil)

	directoryName := generateDirectoryName(testName)
	dirClient, err := srClient.NewDirectoryClient(directoryName)
	_require.Nil(err)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotEqual(*cResp.ETag, "")
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotEqual(*cResp.RequestID, "")
	_require.NotEqual(*cResp.Version, "")

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)

	defer delDirectory(_require, dirClient)
}

func (s *azfileLiveTestSuite) TestDirSetProperties() {

	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)

	defer delShare(_require, srClient, nil)

	directoryName := generateDirectoryName(testName)
	dirClient, err := srClient.NewDirectoryClient(directoryName)

	cResp, err := dirClient.Create(ctx, nil)
	_require.Nil(err)
	key := *cResp.FilePermissionKey

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(ctx, &DirectorySetPropertiesOptions{FilePermissions: &FilePermissions{PermissionStr: &sampleSDDL}})
	_require.Nil(err)
	_require.NotEqual(*sResp.FilePermissionKey, key)
	key = *sResp.FilePermissionKey

	gResp, err := dirClient.GetProperties(ctx, nil)
	_require.Nil(err)
	// Ensure the new key is present when we GetProperties
	_require.Equal(*gResp.FilePermissionKey, key)
}

func (s *azfileLiveTestSuite) TestDirCreateDeleteNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)

	defer delShare(_require, srClient, nil)

	directoryName := generateDirectoryName(testName)
	dirClient, err := srClient.NewDirectoryClient(directoryName)
	_require.Nil(err)

	md := map[string]string{
		"Foo": "FooValuE",
		"Bar": "bArvaLue",
	}

	cResp, err := dirClient.Create(context.Background(), &DirectoryCreateOptions{Metadata: md, FilePermissions: &FilePermissions{PermissionStr: &sampleSDDL}})
	_require.Nil(err)
	// Ensure that the file key isn't empty, but don't worry about checking the permission. We just need to know it exists.
	_require.NotEqual(*cResp.FilePermissionKey, "")
	// _require.Equal(cResp.RawResponse.StatusCode, 201)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotEqual(*cResp.ETag, "")
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotEqual(*cResp.RequestID, "")
	_require.NotEqual(*cResp.Version, "")

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)

	// Creating again will result in 409 and ResourceAlreadyExists.
	cResp, err = dirClient.Create(context.Background(), &DirectoryCreateOptions{Metadata: md})
	_require.NotNil(err)
	//serr := err.(*ShareError)
	//_require.Equal(serr.Response().StatusCode,409)
	validateShareErrorCode(_require, err, ShareErrorCodeResourceAlreadyExists)

	dResp, err := dirClient.Delete(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(dResp.RawResponse.StatusCode, 202)
	_require.Equal(dResp.Date.IsZero(), false)
	_require.NotEqual(*dResp.RequestID, "")
	_require.NotEqual(*dResp.Version, "")

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.NotNil(err)
	validateShareErrorCode(_require, err, ShareErrorCodeResourceNotFound)
}

func (s *azfileLiveTestSuite) TestDirCreateDeleteNegativeMultiLevelDir() {

	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)

	defer delShare(_require, srClient, nil)

	parentDirName := "parent" + generateDirectoryName(testName)
	subDirName := "subdir" + generateDirectoryName(testName)
	parentDirClient, err := srClient.NewDirectoryClient(parentDirName)
	_require.Nil(err)

	subDirClient, err := parentDirClient.NewDirectoryClient(subDirName)

	// Directory create with subDirClient
	_, err = subDirClient.Create(context.Background(), nil)
	_require.NotNil(err)
	validateShareErrorCode(_require, err, ShareErrorCodeParentNotFound)

	_, err = parentDirClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	_, err = subDirClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)

	_, err = subDirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)

	// Delete level by level
	// Delete Non-empty directory should fail
	_, err = parentDirClient.Delete(context.Background(), nil)
	_require.NotNil(err)
	validateShareErrorCode(_require, err, ShareErrorCodeDirectoryNotEmpty)

	_, err = subDirClient.Delete(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(dResp.RawResponse.StatusCode, 202)

	_, err = parentDirClient.Delete(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(dResp.RawResponse.StatusCode, 202)
}

func (s *azfileLiveTestSuite) TestDirCreateEndWithSlash() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)

	defer delShare(_require, srClient, nil)

	directoryName := generateDirectoryName(testName) + "/"
	dirClient, err := srClient.NewDirectoryClient(directoryName)
	_require.Nil(err)

	defer delDirectory(_require, dirClient)

	cResp, err := dirClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(cResp.RawResponse.StatusCode, 201)
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotEqual(*cResp.ETag, "")
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotEqual(*cResp.RequestID, "")
	_require.NotEqual(*cResp.Version, "")

	_, err = dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)
}

func (s *azfileLiveTestSuite) TestDirGetSetMetadataDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	dirClient := createNewDirectoryFromShare(_require, generateDirectoryName(testName), srClient)
	defer delDirectory(_require, dirClient)

	sResp, err := dirClient.SetMetadata(context.Background(), map[string]string{}, nil)
	_require.Nil(err)
	// _require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.NotEqual(*sResp.RequestID, "")
	_require.NotEqual(*sResp.Version, "")
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_require.NotNil(gResp.IsServerEncrypted)
	_require.Len(gResp.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestDirGetSetMetadataNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	dirClient := createNewDirectoryFromShare(_require, generateDirectoryName(testName), srClient)
	defer delDirectory(_require, dirClient)

	md := map[string]string{
		"Foo": "FooValuE",
		"Bar": "bArvaLue",
	}

	sResp, err := dirClient.SetMetadata(context.Background(), md, nil)
	_require.Nil(err)
	// _require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.NotEqual(*sResp.RequestID, "")
	_require.NotEqual(*sResp.Version, "")
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_require.NotNil(gResp.IsServerEncrypted)
	nmd := gResp.Metadata
	_require.EqualValues(nmd, md)
}

func (s *azfileLiveTestSuite) TestDirSetMetadataNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	dirClient := createNewDirectoryFromShare(_require, generateDirectoryName(testName), srClient)
	defer delDirectory(_require, dirClient)

	md := map[string]string{
		"!@#$%^&*()": "!@#$%^&*()",
	}

	_, err := dirClient.SetMetadata(context.Background(), md, nil)
	_require.NotNil(err)
}

func (s *azfileLiveTestSuite) TestDirGetPropertiesNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	dirClient := getDirectoryClientFromShare(_require, generateDirectoryName(testName), srClient)

	_, err := dirClient.GetProperties(ctx, nil)
	_require.NotNil(err)
	validateShareErrorCode(_require, err, ShareErrorCodeResourceNotFound)
}

func (s *azfileLiveTestSuite) TestDirGetPropertiesWithBaseDirectory() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	dirClient, err := srClient.NewRootDirectoryClient()
	_require.Nil(err)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotNil(gResp.IsServerEncrypted)
}

// Merge is not supported, as the key of metadata would be canonicalized
func (s *azfileLiveTestSuite) TestDirGetSetMetadataMergeAndReplace() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient := getServiceClient(nil, nil, testAccountDefault, nil)

	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	dirClient := createNewDirectoryFromShare(_require, generateDirectoryName(testName), srClient)
	defer delDirectory(_require, dirClient)

	md := map[string]string{
		"Color": "RED",
	}

	sResp, err := dirClient.SetMetadata(context.Background(), md, nil)
	_require.Nil(err)
	// _require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.NotEqual(*sResp.RequestID, "")
	_require.NotEqual(*sResp.Version, "")
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_require.NotNil(gResp.IsServerEncrypted)
	_require.EqualValues(gResp.Metadata, md)

	md2 := map[string]string{
		"Color": "WHITE",
	}

	sResp, err = dirClient.SetMetadata(context.Background(), md2, nil)
	_require.Nil(err)
	// _require.Equal(sResp.RawResponse.StatusCode, 200)
	_require.Equal(sResp.Date.IsZero(), false)
	_require.NotEqual(*sResp.ETag, "")
	_require.NotEqual(*sResp.RequestID, "")
	_require.NotEqual(*sResp.Version, "")
	_require.NotNil(sResp.IsServerEncrypted)

	gResp, err = dirClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(gResp.RawResponse.StatusCode, 200)
	_require.Equal(gResp.Date.IsZero(), false)
	_require.NotEqual(*gResp.ETag, "")
	_require.Equal(gResp.LastModified.IsZero(), false)
	_require.NotEqual(*gResp.RequestID, "")
	_require.NotEqual(*gResp.Version, "")
	_require.NotNil(gResp.IsServerEncrypted)
	nmd2 := gResp.Metadata
	_require.EqualValues(nmd2, md2)
}

//func (s *azfileLiveTestSuite) TestDirListDefault() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)
//	if err != nil {
//
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//
//	defer delShare(_require, srClient, nil)
//
//	dirName := generateDirectoryName(testName)
//	dirClient := createNewDirectoryFromShare(_require, dirName, srClient)
//
//	defer delDirectory(_require, dirClient)
//
//	// Empty directory
//	lResp, err := dirClient.ListFilesAndDirectories(context.Background(), Marker{}, DirectoryListFilesAndDirectoriesOptions{})
//	_require.Nil(err)
//	_require(lResp.Response().StatusCode, chk.Equals, 200)
//	_require(lResp.StatusCode(), chk.Equals, 200)
//	_require(lResp.Status(), chk.Not(chk.Equals), "")
//	_require(lResp.ContentType(), chk.Not(chk.Equals), "")
//	_require(lResp.Date().IsZero(), chk.Equals, false)
//	_require(lResp.RequestID(), chk.Not(chk.Equals), "")
//	_require(lResp.Version(), chk.Not(chk.Equals), "")
//	_require(lResp.DirectoryPath, chk.Equals, dirName)
//	_require(lResp.ShareName, chk.Equals, shareName)
//	_require(lResp.ServiceEndpoint, chk.NotNil)
//	_require(lResp.ShareSnapshot, chk.IsNil)
//	_require(lResp.Prefix, chk.Equals, "")
//	_require(lResp.MaxResults, chk.IsNil)
//	_require(lResp.FileItems, chk.HasLen, 0)
//	_require(lResp.DirectoryItems, chk.HasLen, 0)
//
//	innerDir, innerDirName := createNewDirectoryWithPrefix(c, dirClient, "111")
//	defer delDirectory(c, innerDir)
//
//	innerFile, innerFileName := createNewFileWithPrefix(c, dirClient, "111", 0)
//	defer delFile(_require, innerFile)
//
//	// List 1 file, 1 directory
//	lResp2, err := dirClient.ListFilesAndDirectories(context.Background(), Marker{}, DirectoryListFilesAndDirectoriesOptions{})
//	_require.Nil(err)
//	_require(lResp2.Response().StatusCode, chk.Equals, 200)
//	_require(lResp2.DirectoryItems, chk.HasLen, 1)
//	_require(lResp2.DirectoryItems[0].Name, chk.Equals, innerDirName)
//	_require(lResp2.FileItems, chk.HasLen, 1)
//	_require(lResp2.FileItems[0].Name, chk.Equals, innerFileName)
//	_require(lResp2.FileItems[0].Properties.ContentLength, chk.Equals, int64(0))
//
//	innerDir2, innerDirName2 := createNewDirectoryWithPrefix(c, dirClient, "222")
//	defer delDirectory(c, innerDir2)
//
//	innerFile2, innerFileName2 := createNewFileWithPrefix(c, dirClient, "222", 2)
//	defer delFile(_require, innerFile2)
//
//	// List 2 files and 2 directories
//	lResp3, err := dirClient.ListFilesAndDirectories(context.Background(), Marker{}, DirectoryListFilesAndDirectoriesOptions{})
//	_require.Nil(err)
//	_require(lResp3.Response().StatusCode, chk.Equals, 200)
//	_require(lResp3.DirectoryItems, chk.HasLen, 2)
//	_require(lResp3.DirectoryItems[0].Name, chk.Equals, innerDirName)
//	_require(lResp3.DirectoryItems[1].Name, chk.Equals, innerDirName2)
//	_require(lResp3.FileItems, chk.HasLen, 2)
//	_require(lResp3.FileItems[0].Name, chk.Equals, innerFileName)
//	_require(lResp3.FileItems[0].Properties.ContentLength, chk.Equals, int64(0))
//	_require(lResp3.FileItems[1].Name, chk.Equals, innerFileName2)
//	_require(lResp3.FileItems[1].Properties.ContentLength, chk.Equals, int64(2))
//}

//func (s *azfileLiveTestSuite) TestDirListNonDefault() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)
//	if err != nil {
//
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//
//	defer delShare(_require, srClient, nil)
//
//	dirClient, _ := createNewDirectoryFromShare(_require, generateDirectoryName(testName), srClient)
//
//	defer delDirectory(c, dirClient)
//
//	const testPrefix = "pagedprefix"
//	const maxResultsPerPage = 2
//
//	dir1, dir1Name := createNewDirectoryWithPrefix(c, dirClient, testPrefix+"1")
//	file1, file1Name := createNewFileWithPrefix(c, dirClient, testPrefix+"2", 0)
//	dir2, dir2Name := createNewDirectoryWithPrefix(c, dirClient, testPrefix+"3")
//	file2, file2Name := createNewFileWithPrefix(c, dirClient, testPrefix+"4", 0)
//
//	defer func() {
//		delDirectory(c, dir1)
//		delDirectory(c, dir2)
//		delFile(c, file1)
//		delFile(c, file2)
//	}()
//
//	marker := Marker{}
//
//	lResp, err := dirClient.ListFilesAndDirectories(context.Background(), marker, DirectoryListFilesAndDirectoriesOptions{MaxResults: maxResultsPerPage, Prefix: testPrefix})
//	_require.Nil(err)
//	_require(lResp.FileItems, chk.HasLen, 1)
//	_require(lResp.FileItems[0].Name, chk.Equals, file1Name)
//	_require(lResp.DirectoryItems, chk.HasLen, 1)
//	_require(lResp.DirectoryItems[0].Name, chk.Equals, dir1Name)
//
//	_require(lResp.NextMarker.NotDone(), chk.Equals, true)
//	marker = lResp.NextMarker
//
//	lResp, err = dirClient.ListFilesAndDirectories(context.Background(), marker, DirectoryListFilesAndDirectoriesOptions{MaxResults: maxResultsPerPage, Prefix: testPrefix})
//	_require.Nil(err)
//	_require(lResp.FileItems, chk.HasLen, 1)
//	_require(lResp.FileItems[0].Name, chk.Equals, file2Name)
//	_require(lResp.DirectoryItems, chk.HasLen, 1)
//	_require(lResp.DirectoryItems[0].Name, chk.Equals, dir2Name)
//
//	_require(lResp.NextMarker.NotDone(), chk.Equals, false)
//}

//func (s *azfileLiveTestSuite) TestDirListNegativeNonexistantPrefix() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)
//	if err != nil {
//
//	}
//	shareURL, _ := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	createNewFileFromShare(c, shareURL, 0)
//
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	resp, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, DirectoryListFilesAndDirectoriesOptions{Prefix: filePrefix + filePrefix})
//
//	_require.Nil(err)
//	_require(resp.FileItems, chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestDirListNegativeMaxResults() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)
//	if err != nil {
//
//	}
//	shareURL, _ := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	createNewFileFromShare(c, shareURL, 0)
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	_, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, DirectoryListFilesAndDirectoriesOptions{MaxResults: -2})
//	_require.NotNil(err)
//	_require(strings.Contains(err.Error(), "validation failed"), chk.Equals, true)
//}
//
//func (s *azfileLiveTestSuite) TestDirListNonDefaultMaxResultsZero() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)
//	if err != nil {
//
//	}
//	shareURL, _ := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	createNewFileFromShare(c, shareURL, 0)
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	resp, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, DirectoryListFilesAndDirectoriesOptions{MaxResults: 0})
//
//	_require.Nil(err)
//	_require(resp.FileItems, chk.HasLen, 1)
//}
//
//func (s *azfileLiveTestSuite) TestDirListNonDefaultMaxResultsExact() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)
//	if err != nil {
//
//	}
//	shareURL, _ := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	additionalPrefix := strconv.Itoa(time.Now().Nanosecond())
//	_, dirName1 := createNewDirectoryWithPrefix(c, dirURL, additionalPrefix+"a")
//	_, dirName2 := createNewDirectoryWithPrefix(c, dirURL, additionalPrefix+"b")
//
//	resp, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, DirectoryListFilesAndDirectoriesOptions{MaxResults: 2, Prefix: additionalPrefix})
//
//	_require.Nil(err)
//	_require(resp.DirectoryItems, chk.HasLen, 2)
//	_require(resp.DirectoryItems[0].Name, chk.Equals, dirName1)
//	_require(resp.DirectoryItems[1].Name, chk.Equals, dirName2)
//}
//
//// Test list directories with SAS
//func (s *azfileLiveTestSuite) TestDirListWithShareSAS() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient := getServiceClient(_require, nil, testAccountDefault, nil)
//	if err != nil {
//
//	}
//	credential, accountName := getCredential()
//	share, shareName := createNewShare(_require, generateShareName(testName), svcClient)
//
//	defer delShare(_require, srClient, nil)
//
//	dirClient, dirName := createNewDirectoryFromShare(_require, generateDirectoryName(testName), srClient)
//
//	defer delDirectory(c, dirClient)
//
//	// Create share service SAS
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:    SASProtocolHTTPS,              // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ShareName:   shareName,
//		FilePermissions: ShareSASPermissions{Read: true, Write: true, List: true}.String(),
//	}.Sign(credential)
//	_require.Nil(err)
//
//	// Create the URL of the resource you wish to access and append the SAS query parameters.
//	// Since this is a file SAS, the URL is to the Azure storage file.
//	qp := sasQueryParams.Encode()
//	urlToSendToSomeone := fmt.Sprintf("https://%s.file.core.windows.net/%s/%s?%s",
//		accountName, shareName, dirName, qp)
//
//	u, _ := url.Parse(urlToSendToSomeone)
//	dirURL := NewDirectoryClient(*u, NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//
//	marker := Marker{}
//	lResp, err := dirURL.ListFilesAndDirectories(context.Background(), marker, DirectoryListFilesAndDirectoriesOptions{})
//	_require.Nil(err)
//	_require(lResp.NextMarker.NotDone(), chk.Equals, false)
//}
