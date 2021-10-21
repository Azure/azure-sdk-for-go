package azfile

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
)

func (s *azfileLiveTestSuite) TestDirNewDirectoryClient() {
	_assert := assert.New(s.T())
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient, err := svcClient.NewShareClient(sharePrefix)
	_assert.Nil(err)
	dirClient, err := srClient.NewDirectoryClient(directoryPrefix)
	_assert.Nil(err)
	dirClient, err = dirClient.NewDirectoryClient("inner" + directoryPrefix)
	_assert.Nil(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/" + sharePrefix + "/" + directoryPrefix + "/inner" + directoryPrefix
	_assert.Equal(dirClient.URL(), correctURL)
}

func (s *azfileLiveTestSuite) TestDirCreateFileURL() {
	_assert := assert.New(s.T())
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient, err := svcClient.NewShareClient(sharePrefix)
	_assert.Nil(err)
	dirClient, err := srClient.NewDirectoryClient(directoryPrefix)
	_assert.Nil(err)
	fClient, err := dirClient.NewFileClient(filePrefix)
	_assert.Nil(err)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".file.core.windows.net/" + sharePrefix + "/" + directoryPrefix + "/" + filePrefix
	_assert.Equal(fClient.URL(), correctURL)
}

func (s *azfileLiveTestSuite) TestDirCreateDeleteDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)

	defer delShare(_assert, srClient, nil)

	directoryName := generateDirectoryName(testName)
	dirClient, err := srClient.NewDirectoryClient(directoryName)
	_assert.Nil(err)

	cResp, err := dirClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	_assert.Equal(cResp.Date.IsZero(), false)
	_assert.NotEqual(*cResp.ETag, "")
	_assert.Equal(cResp.LastModified.IsZero(), false)
	_assert.NotEqual(*cResp.RequestID, "")
	_assert.NotEqual(*cResp.Version, "")

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)

	defer delDirectory(_assert, dirClient)
}

func (s *azfileLiveTestSuite) TestDirSetProperties() {

	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)

	defer delShare(_assert, srClient, nil)

	directoryName := generateDirectoryName(testName)
	dirClient, err := srClient.NewDirectoryClient(directoryName)

	cResp, err := dirClient.Create(ctx, nil)
	_assert.Nil(err)
	key := *cResp.FilePermissionKey

	// Set the custom permissions
	sResp, err := dirClient.SetProperties(ctx, &SetDirectoryPropertiesOptions{FilePermissions: &Permissions{PermissionStr: &sampleSDDL}})
	_assert.Nil(err)
	_assert.NotEqual(*sResp.FilePermissionKey, key)
	key = *sResp.FilePermissionKey

	gResp, err := dirClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	// Ensure the new key is present when we GetProperties
	_assert.Equal(*gResp.FilePermissionKey, key)
}

func (s *azfileLiveTestSuite) TestDirCreateDeleteNonDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)

	defer delShare(_assert, srClient, nil)

	directoryName := generateDirectoryName(testName)
	dirClient, err := srClient.NewDirectoryClient(directoryName)
	_assert.Nil(err)

	md := map[string]string{
		"Foo": "FooValuE",
		"Bar": "bArvaLue",
	}

	cResp, err := dirClient.Create(context.Background(), &CreateDirectoryOptions{Metadata: md, FilePermissions: &Permissions{PermissionStr: &sampleSDDL}})
	_assert.Nil(err)
	// Ensure that the file key isn't empty, but don't worry about checking the permission. We just need to know it exists.
	_assert.NotEqual(*cResp.FilePermissionKey, "")
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	_assert.Equal(cResp.Date.IsZero(), false)
	_assert.NotEqual(*cResp.ETag, "")
	_assert.Equal(cResp.LastModified.IsZero(), false)
	_assert.NotEqual(*cResp.RequestID, "")
	_assert.NotEqual(*cResp.Version, "")

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)

	// Creating again will result in 409 and ResourceAlreadyExists.
	cResp, err = dirClient.Create(context.Background(), &CreateDirectoryOptions{Metadata: md})
	_assert.NotNil(err)
	//serr := err.(*StorageError)
	//_assert.Equal(serr.Response().StatusCode,409)
	validateStorageError(_assert, err, StorageErrorCodeResourceAlreadyExists)

	dResp, err := dirClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(dResp.RawResponse.StatusCode, 202)
	_assert.Equal(dResp.Date.IsZero(), false)
	_assert.NotEqual(*dResp.RequestID, "")
	_assert.NotEqual(*dResp.Version, "")

	gResp, err = dirClient.GetProperties(context.Background(), nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeResourceNotFound)
}

func (s *azfileLiveTestSuite) TestDirCreateDeleteNegativeMultiLevelDir() {

	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)

	defer delShare(_assert, srClient, nil)

	parentDirName := "parent" + generateDirectoryName(testName)
	subDirName := "subdir" + generateDirectoryName(testName)
	parentDirClient, err := srClient.NewDirectoryClient(parentDirName)
	_assert.Nil(err)

	subDirClient, err := parentDirClient.NewDirectoryClient(subDirName)

	// Directory create with subDirClient
	cResp, err := subDirClient.Create(context.Background(), nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeParentNotFound)

	cResp, err = parentDirClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	cResp, err = subDirClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	gResp, err := subDirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)

	// Delete level by level
	// Delete Non-empty directory should fail
	_, err = parentDirClient.Delete(context.Background(), nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeDirectoryNotEmpty)

	dResp, err := subDirClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(dResp.RawResponse.StatusCode, 202)

	dResp, err = parentDirClient.Delete(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(dResp.RawResponse.StatusCode, 202)
}

func (s *azfileLiveTestSuite) TestDirCreateEndWithSlash() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)

	defer delShare(_assert, srClient, nil)

	directoryName := generateDirectoryName(testName) + "/"
	dirClient, err := srClient.NewDirectoryClient(directoryName)
	_assert.Nil(err)

	defer delDirectory(_assert, dirClient)

	cResp, err := dirClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	_assert.Equal(cResp.Date.IsZero(), false)
	_assert.NotEqual(*cResp.ETag, "")
	_assert.Equal(cResp.LastModified.IsZero(), false)
	_assert.NotEqual(*cResp.RequestID, "")
	_assert.NotEqual(*cResp.Version, "")

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
}

func (s *azfileLiveTestSuite) TestDirGetSetMetadataDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	dirClient := createNewDirectoryFromShare(_assert, generateDirectoryName(testName), srClient)
	defer delDirectory(_assert, dirClient)

	sResp, err := dirClient.SetMetadata(context.Background(), map[string]string{}, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.NotEqual(*sResp.RequestID, "")
	_assert.NotEqual(*sResp.Version, "")
	_assert.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.NotNil(gResp.IsServerEncrypted)
	_assert.Len(gResp.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestDirGetSetMetadataNonDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	dirClient := createNewDirectoryFromShare(_assert, generateDirectoryName(testName), srClient)
	defer delDirectory(_assert, dirClient)

	md := map[string]string{
		"Foo": "FooValuE",
		"Bar": "bArvaLue",
	}

	sResp, err := dirClient.SetMetadata(context.Background(), md, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.NotEqual(*sResp.RequestID, "")
	_assert.NotEqual(*sResp.Version, "")
	_assert.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.NotNil(gResp.IsServerEncrypted)
	nmd := gResp.Metadata
	_assert.EqualValues(nmd, md)
}

func (s *azfileLiveTestSuite) TestDirSetMetadataNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	dirClient := createNewDirectoryFromShare(_assert, generateDirectoryName(testName), srClient)
	defer delDirectory(_assert, dirClient)

	md := map[string]string{
		"!@#$%^&*()": "!@#$%^&*()",
	}

	_, err = dirClient.SetMetadata(context.Background(), md, nil)
	_assert.NotNil(err)
}

func (s *azfileLiveTestSuite) TestDirGetPropertiesNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	dirClient := getDirectoryClientFromShare(_assert, generateDirectoryName(testName), srClient)

	_, err = dirClient.GetProperties(ctx, nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeResourceNotFound)
}

func (s *azfileLiveTestSuite) TestDirGetPropertiesWithBaseDirectory() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	dirClient, err := srClient.NewRootDirectoryClient()
	_assert.Nil(err)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotNil(gResp.IsServerEncrypted)
}

// Merge is not supported, as the key of metadata would be canonicalized
func (s *azfileLiveTestSuite) TestDirGetSetMetadataMergeAndReplace() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	dirClient := createNewDirectoryFromShare(_assert, generateDirectoryName(testName), srClient)
	defer delDirectory(_assert, dirClient)

	md := map[string]string{
		"Color": "RED",
	}

	sResp, err := dirClient.SetMetadata(context.Background(), md, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.NotEqual(*sResp.RequestID, "")
	_assert.NotEqual(*sResp.Version, "")
	_assert.NotNil(sResp.IsServerEncrypted)

	gResp, err := dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.NotNil(gResp.IsServerEncrypted)
	_assert.EqualValues(gResp.Metadata, md)

	md2 := map[string]string{
		"Color": "WHITE",
	}

	sResp, err = dirClient.SetMetadata(context.Background(), md2, nil)
	_assert.Nil(err)
	_assert.Equal(sResp.RawResponse.StatusCode, 200)
	_assert.Equal(sResp.Date.IsZero(), false)
	_assert.NotEqual(*sResp.ETag, "")
	_assert.NotEqual(*sResp.RequestID, "")
	_assert.NotEqual(*sResp.Version, "")
	_assert.NotNil(sResp.IsServerEncrypted)

	gResp, err = dirClient.GetProperties(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(gResp.RawResponse.StatusCode, 200)
	_assert.Equal(gResp.Date.IsZero(), false)
	_assert.NotEqual(*gResp.ETag, "")
	_assert.Equal(gResp.LastModified.IsZero(), false)
	_assert.NotEqual(*gResp.RequestID, "")
	_assert.NotEqual(*gResp.Version, "")
	_assert.NotNil(gResp.IsServerEncrypted)
	nmd2 := gResp.Metadata
	_assert.EqualValues(nmd2, md2)
}

//func (s *azfileLiveTestSuite) TestDirListDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//
//	defer delShare(_assert, srClient, nil)
//
//	dirName := generateDirectoryName(testName)
//	dirClient := createNewDirectoryFromShare(_assert, dirName, srClient)
//
//	defer delDirectory(_assert, dirClient)
//
//	// Empty directory
//	lResp, err := dirClient.ListFilesAndDirectories(context.Background(), Marker{}, ListFilesAndDirectoriesOptions{})
//	_assert.Nil(err)
//	_assert(lResp.Response().StatusCode, chk.Equals, 200)
//	_assert(lResp.StatusCode(), chk.Equals, 200)
//	_assert(lResp.Status(), chk.Not(chk.Equals), "")
//	_assert(lResp.ContentType(), chk.Not(chk.Equals), "")
//	_assert(lResp.Date().IsZero(), chk.Equals, false)
//	_assert(lResp.RequestID(), chk.Not(chk.Equals), "")
//	_assert(lResp.Version(), chk.Not(chk.Equals), "")
//	_assert(lResp.DirectoryPath, chk.Equals, dirName)
//	_assert(lResp.ShareName, chk.Equals, shareName)
//	_assert(lResp.ServiceEndpoint, chk.NotNil)
//	_assert(lResp.ShareSnapshot, chk.IsNil)
//	_assert(lResp.Prefix, chk.Equals, "")
//	_assert(lResp.MaxResults, chk.IsNil)
//	_assert(lResp.FileItems, chk.HasLen, 0)
//	_assert(lResp.DirectoryItems, chk.HasLen, 0)
//
//	innerDir, innerDirName := createNewDirectoryWithPrefix(c, dirClient, "111")
//	defer delDirectory(c, innerDir)
//
//	innerFile, innerFileName := createNewFileWithPrefix(c, dirClient, "111", 0)
//	defer delFile(c, innerFile)
//
//	// List 1 file, 1 directory
//	lResp2, err := dirClient.ListFilesAndDirectories(context.Background(), Marker{}, ListFilesAndDirectoriesOptions{})
//	_assert.Nil(err)
//	_assert(lResp2.Response().StatusCode, chk.Equals, 200)
//	_assert(lResp2.DirectoryItems, chk.HasLen, 1)
//	_assert(lResp2.DirectoryItems[0].Name, chk.Equals, innerDirName)
//	_assert(lResp2.FileItems, chk.HasLen, 1)
//	_assert(lResp2.FileItems[0].Name, chk.Equals, innerFileName)
//	_assert(lResp2.FileItems[0].Properties.ContentLength, chk.Equals, int64(0))
//
//	innerDir2, innerDirName2 := createNewDirectoryWithPrefix(c, dirClient, "222")
//	defer delDirectory(c, innerDir2)
//
//	innerFile2, innerFileName2 := createNewFileWithPrefix(c, dirClient, "222", 2)
//	defer delFile(c, innerFile2)
//
//	// List 2 files and 2 directories
//	lResp3, err := dirClient.ListFilesAndDirectories(context.Background(), Marker{}, ListFilesAndDirectoriesOptions{})
//	_assert.Nil(err)
//	_assert(lResp3.Response().StatusCode, chk.Equals, 200)
//	_assert(lResp3.DirectoryItems, chk.HasLen, 2)
//	_assert(lResp3.DirectoryItems[0].Name, chk.Equals, innerDirName)
//	_assert(lResp3.DirectoryItems[1].Name, chk.Equals, innerDirName2)
//	_assert(lResp3.FileItems, chk.HasLen, 2)
//	_assert(lResp3.FileItems[0].Name, chk.Equals, innerFileName)
//	_assert(lResp3.FileItems[0].Properties.ContentLength, chk.Equals, int64(0))
//	_assert(lResp3.FileItems[1].Name, chk.Equals, innerFileName2)
//	_assert(lResp3.FileItems[1].Properties.ContentLength, chk.Equals, int64(2))
//}

//func (s *azfileLiveTestSuite) TestDirListNonDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//
//	defer delShare(_assert, srClient, nil)
//
//	dirClient, _ := createNewDirectoryFromShare(_assert, generateDirectoryName(testName), srClient)
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
//	lResp, err := dirClient.ListFilesAndDirectories(context.Background(), marker, ListFilesAndDirectoriesOptions{MaxResults: maxResultsPerPage, Prefix: testPrefix})
//	_assert.Nil(err)
//	_assert(lResp.FileItems, chk.HasLen, 1)
//	_assert(lResp.FileItems[0].Name, chk.Equals, file1Name)
//	_assert(lResp.DirectoryItems, chk.HasLen, 1)
//	_assert(lResp.DirectoryItems[0].Name, chk.Equals, dir1Name)
//
//	_assert(lResp.NextMarker.NotDone(), chk.Equals, true)
//	marker = lResp.NextMarker
//
//	lResp, err = dirClient.ListFilesAndDirectories(context.Background(), marker, ListFilesAndDirectoriesOptions{MaxResults: maxResultsPerPage, Prefix: testPrefix})
//	_assert.Nil(err)
//	_assert(lResp.FileItems, chk.HasLen, 1)
//	_assert(lResp.FileItems[0].Name, chk.Equals, file2Name)
//	_assert(lResp.DirectoryItems, chk.HasLen, 1)
//	_assert(lResp.DirectoryItems[0].Name, chk.Equals, dir2Name)
//
//	_assert(lResp.NextMarker.NotDone(), chk.Equals, false)
//}

//func (s *azfileLiveTestSuite) TestDirListNegativeNonexistantPrefix() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareURL, _ := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	createNewFileFromShare(c, shareURL, 0)
//
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	resp, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, ListFilesAndDirectoriesOptions{Prefix: filePrefix + filePrefix})
//
//	_assert.Nil(err)
//	_assert(resp.FileItems, chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestDirListNegativeMaxResults() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareURL, _ := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	createNewFileFromShare(c, shareURL, 0)
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	_, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, ListFilesAndDirectoriesOptions{MaxResults: -2})
//	_assert.NotNil(err)
//	_assert(strings.Contains(err.Error(), "validation failed"), chk.Equals, true)
//}
//
//func (s *azfileLiveTestSuite) TestDirListNonDefaultMaxResultsZero() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareURL, _ := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	createNewFileFromShare(c, shareURL, 0)
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	resp, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, ListFilesAndDirectoriesOptions{MaxResults: 0})
//
//	_assert.Nil(err)
//	_assert(resp.FileItems, chk.HasLen, 1)
//}
//
//func (s *azfileLiveTestSuite) TestDirListNonDefaultMaxResultsExact() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareURL, _ := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(c, shareURL, DeleteSnapshotsOptionNone)
//	dirURL := shareURL.NewRootDirectoryClient()
//
//	additionalPrefix := strconv.Itoa(time.Now().Nanosecond())
//	_, dirName1 := createNewDirectoryWithPrefix(c, dirURL, additionalPrefix+"a")
//	_, dirName2 := createNewDirectoryWithPrefix(c, dirURL, additionalPrefix+"b")
//
//	resp, err := dirURL.ListFilesAndDirectories(ctx, Marker{}, ListFilesAndDirectoriesOptions{MaxResults: 2, Prefix: additionalPrefix})
//
//	_assert.Nil(err)
//	_assert(resp.DirectoryItems, chk.HasLen, 2)
//	_assert(resp.DirectoryItems[0].Name, chk.Equals, dirName1)
//	_assert(resp.DirectoryItems[1].Name, chk.Equals, dirName2)
//}
//
//// Test list directories with SAS
//func (s *azfileLiveTestSuite) TestDirListWithShareSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	credential, accountName := getCredential()
//	share, shareName := createNewShare(_assert, generateShareName(testName), svcClient)
//
//	defer delShare(_assert, srClient, nil)
//
//	dirClient, dirName := createNewDirectoryFromShare(_assert, generateDirectoryName(testName), srClient)
//
//	defer delDirectory(c, dirClient)
//
//	// Create share service SAS
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:    SASProtocolHTTPS,              // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ShareName:   shareName,
//		Permissions: ShareSASPermissions{Read: true, Write: true, List: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_assert.Nil(err)
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
//	lResp, err := dirURL.ListFilesAndDirectories(context.Background(), marker, ListFilesAndDirectoriesOptions{})
//	_assert.Nil(err)
//	_assert(lResp.NextMarker.NotDone(), chk.Equals, false)
//}
