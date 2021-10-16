package azfile

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
)

const (
	testFileRangeSize         = 512           // Use this number considering clear range's function
	fileShareMaxQuota         = 5120          // Size is in GB (Service Version 2020-02-10)
	fileMaxAllowedSizeInBytes = 4398046511104 // 4 TiB (Service Version 2020-02-10)
)

func (s *azfileLiveTestSuite) TestFileCreateDeleteDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareName := generateShareName(testName)
	srClient := createNewShare(_assert, shareName, svcClient)
	defer delShare(_assert, srClient)

	// Create and delete fClient in root directory.
	fileName := generateFileName(testName)
	fClient, err := srClient.NewRootDirectoryClient().NewFileClient(fileName)
	_assert.Nil(err)

	cResp, err := fClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	_assert.NotEqual(cResp.ETag, azcore.ETag(""))
	_assert.Equal(cResp.LastModified.IsZero(), false)
	_assert.NotEqual(cResp.RequestID, "")
	_assert.NotEqual(cResp.Version, "")
	_assert.Equal(cResp.Date.IsZero(), false)
	_assert.NotNil(cResp.IsServerEncrypted)

	delResp, err := fClient.Delete(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(delResp.RawResponse.StatusCode, 202)
	_assert.NotEqual(delResp.RequestID, "")
	_assert.NotEqual(delResp.Version, "")
	_assert.Equal(delResp.Date.IsZero(), false)

	dirName := generateDirectoryName(testName)
	dirClient := createNewDirectoryFromShare(_assert, dirName, srClient)
	defer delDirectory(_assert, dirClient)

	//// Create and delete fClient in named directory.
	//fClient = dir.NewfClient(generateFileName())
	//
	//cResp, err = fClient.Create(ctx, 0, FileHTTPHeaders{}, nil)
	//_assert.Nil(err)
	//_assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	//_assert(cResp.ETag, chk.Not(chk.Equals), ETagNone)
	//_assert(cResp.LastModified.IsZero(), chk.Equals, false)
	//_assert(cResp.RequestID, chk.Not(chk.Equals), "")
	//_assert(cResp.Version, chk.Not(chk.Equals), "")
	//_assert(cResp.Date.IsZero(), chk.Equals, false)
	//_assert(cResp.IsServerEncrypted, chk.NotNil)
	//
	//delResp, err = fClient.Delete(ctx)
	//_assert.Nil(err)
	//_assert(delResp.RawResponse.StatusCode, chk.Equals, 202)
	//_assert(delResp.RequestID, chk.Not(chk.Equals), "")
	//_assert(delResp.Version, chk.Not(chk.Equals), "")
	//_assert(delResp.Date.IsZero(), chk.Equals, false)
}

func (s *azfileLiveTestSuite) TestFileCreateNonDefaultMetadataNonEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient)

	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	_, err = fClient.Create(ctx, &CreateFileOptions{
		Metadata: basicMetadata,
	})
	_assert.Nil(err)

	resp, err := fClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azfileLiveTestSuite) TestFileCreateNonDefaultHTTPHeaders() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient)
	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	_, err = fClient.Create(ctx, &CreateFileOptions{FileHTTPHeaders: &basicHeaders})
	_assert.Nil(err)

	_, err = fClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	//h := resp.RawResponse.Header
	//_assert.EqualValues(h, basicHeaders)
}

func (s *azfileLiveTestSuite) TestFileCreateNegativeMetadataInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient)
	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	_, err = fClient.Create(ctx, &CreateFileOptions{
		Metadata:        map[string]string{"!@#$%^&*()": "!@#$%^&*()"},
		FileHTTPHeaders: &FileHTTPHeaders{},
	})
	_assert.NotNil(err)
}

//func (s *azfileLiveTestSuite) TestFileGetSetPropertiesNonDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//	defer delFile(_assert, fClient)
//
//	md5Str := "MDAwMDAwMDA="
//	var testMd5 []byte
//	copy(testMd5[:], md5Str)
//
//	attribs := FileAttributeTemporary.Add(FileAttributeHidden)
//	creationTime := time.Now().Add(-time.Hour)
//	lastWriteTime := time.Now().Add(-time.Minute * 15)
//
//	// Format and re-parse the times so we have the same precision
//	creationTime, err = time.Parse(ISO8601, creationTime.Format(ISO8601))
//	_assert.Nil(err)
//	lastWriteTime, err = time.Parse(ISO8601, lastWriteTime.Format(ISO8601))
//	_assert.Nil(err)
//
//
//	options := &SetFileHTTPHeadersOptions{
//		FilePermissions: &FilePermissions{ FilePermissionStr: &sampleSDDL},
//		SMBProperties: &SMBProperties{
//			FileAttributes:    &attribs,
//			FileCreationTime:  &creationTime,
//			FileLastWriteTime: &lastWriteTime,
//		},
//		FileHTTPHeaders: &FileHTTPHeaders{
//			FileContentType:        to.StringPtr("text/html"),
//			FileContentEncoding:    to.StringPtr("gzip"),
//			FileContentLanguage:    to.StringPtr("tr,en"),
//			FileContentMD5:         testMd5,
//			FileCacheControl:       to.StringPtr("no-transform"),
//			FileContentDisposition: to.StringPtr("attachment"),
//		},
//	}
//	setResp, err := fClient.SetHTTPHeaders(ctx, options)
//	_assert.Nil(err)
//	_assert.Equal(setResp.RawResponse.StatusCode, 200)
//	_assert.NotEqual(setResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.Equal(setResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.NotEqual(setResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.NotEqual(setResp.Version, chk.Not(chk.Equals), "")
//	_assert.Equal(setResp.Date.IsZero(), chk.Equals, false)
//	_assert.NotNil(setResp.IsServerEncrypted, chk.NotNil)
//
//	getResp, err := fClient.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert.Equal(getResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.Equal(setResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.Equal(*getResp.FileType, chk.Equals, "File")
//
//	_assert.(*getResp.ContentType, chk.Equals, properties.ContentType)
//	_assert.(*getResp.ContentEncoding, chk.Equals, properties.ContentEncoding)
//	_assert.(*getResp.ContentLanguage, chk.Equals, properties.ContentLanguage)
//	_assert.(*getResp.ContentMD5,, chk.DeepEquals, properties.ContentMD5)
//	_assert.(*getResp.CacheControl, chk.Equals, properties.CacheControl)
//	_assert.(*getResp.ContentDisposition, chk.Equals, properties.ContentDisposition)
//	_assert.(*getResp.ContentLength, chk.Equals, int64(0))
//	// We'll just ensure a permission exists, no need to test overlapping functionality.
//	_assert.(*getResp.FilePermissionKey, chk.Not(chk.Equals), "")
//	// Ensure our attributes and other properties (after parsing) are equivalent to our original
//	// There's an overlapping test for this in ntfs_property_bitflags_test.go, but it doesn't hurt to test it alongside other things.
//	_assert.(ParseFileAttributeFlagsString(getResp.FileAttributes()), chk.Equals, attribs)
//	// Adapt to time.Time
//	adapter := SMBPropertyAdapter{PropertySource: getResp}
//	c.Log("Original last write time: ", lastWriteTime, " new time: ", adapter.FileLastWriteTime())
//	_assert.(adapter.FileLastWriteTime().Equal(lastWriteTime), chk.Equals, true)
//	c.Log("Original creation time: ", creationTime, " new time: ", adapter.FileCreationTime())
//	_assert.(adapter.FileCreationTime().Equal(creationTime), chk.Equals, true)
//
//	_assert.(getResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(getResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(getResp.Version, chk.Not(chk.Equals), "")
//	_assert.(getResp.Date.IsZero(), chk.Equals, false)
//	_assert.(getResp.IsServerEncrypted, chk.NotNil)
//}
//
//func (s *azfileLiveTestSuite) TestFilePreservePermissions() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShareWithPermissions(c, srClient, 0)
//	defer delFile(c, fClient)
//
//	// Grab the original perm key before we set file headers.
//	getResp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//
//	oKey := getResp.FilePermissionKey()
//	timeAdapter := SMBPropertyAdapter{PropertySource: getResp}
//	cTime := timeAdapter.FileCreationTime()
//	lwTime := timeAdapter.FileLastWriteTime()
//	attribs := getResp.FileAttributes()
//
//	md5Str := "MDAwMDAwMDA="
//	var testMd5 []byte
//	copy(testMd5[:], md5Str)
//
//	properties := FileHTTPHeaders{
//		ContentType:        "text/html",
//		ContentEncoding:    "gzip",
//		ContentLanguage:    "tr,en",
//		ContentMD5:         testMd5,
//		CacheControl:       "no-transform",
//		ContentDisposition: "attachment",
//		SMBProperties:      SMBProperties{
//			// SMBProperties, when options are left nil, leads to preserving.
//		},
//	}
//
//	setResp, err := fClient.SetHTTPHeaders(ctx, properties)
//	_assert.Nil(err)
//	_assert.(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(setResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(setResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(setResp.Version, chk.Not(chk.Equals), "")
//	_assert.(setResp.Date.IsZero(), chk.Equals, false)
//	_assert.(setResp.IsServerEncrypted, chk.NotNil)
//
//	getResp, err = fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(getResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(getResp.FileType(), chk.Equals, "File")
//
//	_assert.(getResp.ContentType, chk.Equals, properties.ContentType)
//	_assert.(getResp.ContentEncoding(), chk.Equals, properties.ContentEncoding)
//	_assert.(getResp.ContentLanguage, chk.Equals, properties.ContentLanguage)
//	_assert.(getResp.ContentMD5,, chk.DeepEquals, properties.ContentMD5)
//	_assert.(getResp.CacheControl, chk.Equals, properties.CacheControl)
//	_assert.(getResp.ContentDisposition, chk.Equals, properties.ContentDisposition)
//	_assert.(getResp.ContentLength, chk.Equals, int64(0))
//	// Ensure that the permission key gets preserved
//	_assert.(getResp.FilePermissionKey(), chk.Equals, oKey)
//	timeAdapter = SMBPropertyAdapter{PropertySource: getResp}
//	c.Log("Original last write time: ", lwTime, " new time: ", timeAdapter.FileLastWriteTime())
//	_assert.(timeAdapter.FileLastWriteTime().Equal(lwTime), chk.Equals, true)
//	c.Log("Original creation time: ", cTime, " new time: ", timeAdapter.FileCreationTime())
//	_assert.(timeAdapter.FileCreationTime().Equal(cTime), chk.Equals, true)
//	_assert.(getResp.FileAttributes(), chk.Equals, attribs)
//
//	_assert.(getResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(getResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(getResp.Version, chk.Not(chk.Equals), "")
//	_assert.(getResp.Date.IsZero(), chk.Equals, false)
//	_assert.(getResp.IsServerEncrypted, chk.NotNil)
//}
//
//func (s *azfileLiveTestSuite) TestFileGetSetPropertiesSnapshot() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(c, srClient, DeleteSnapshotsOptionInclude)
//
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	md5Str := "MDAwMDAwMDA="
//	var testMd5 []byte
//	copy(testMd5[:], md5Str)
//
//	properties := FileHTTPHeaders{
//		ContentType:        "text/html",
//		ContentEncoding:    "gzip",
//		ContentLanguage:    "tr,en",
//		ContentMD5:         testMd5,
//		CacheControl:       "no-transform",
//		ContentDisposition: "attachment",
//	}
//	setResp, err := fClient.SetHTTPHeaders(ctx, properties)
//	_assert.Nil(err)
//	_assert.(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(setResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(setResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(setResp.Version, chk.Not(chk.Equals), "")
//	_assert.(setResp.Date.IsZero(), chk.Equals, false)
//	_assert.(setResp.IsServerEncrypted, chk.NotNil)
//
//	metadata := Metadata{
//		"foo": "foovalue",
//		"bar": "barvalue",
//	}
//	setResp2, err := fClient.SetMetadata(ctx, metadata)
//	_assert.Nil(err)
//	_assert.(setResp2.RawResponse.StatusCode, chk.Equals, 200)
//
//	resp, _ := srClient.CreateSnapshot(ctx, Metadata{})
//	snapshotURL := fClient.WithSnapshot(resp.Snapshot())
//
//	getResp, err := snapshotURL.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(getResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(getResp.FileType(), chk.Equals, "File")
//
//	_assert.(getResp.ContentType, chk.Equals, properties.ContentType)
//	_assert.(getResp.ContentEncoding(), chk.Equals, properties.ContentEncoding)
//	_assert.(getResp.ContentLanguage, chk.Equals, properties.ContentLanguage)
//	_assert.(getResp.ContentMD5,, chk.DeepEquals, properties.ContentMD5)
//	_assert.(getResp.CacheControl, chk.Equals, properties.CacheControl)
//	_assert.(getResp.ContentDisposition, chk.Equals, properties.ContentDisposition)
//	_assert.(getResp.ContentLength, chk.Equals, int64(0))
//
//	_assert.(getResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(getResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(getResp.Version, chk.Not(chk.Equals), "")
//	_assert.(getResp.Date.IsZero(), chk.Equals, false)
//	_assert.(getResp.IsServerEncrypted, chk.NotNil)
//	_assert.(getResp.NewMetadata(), chk.DeepEquals, metadata)
//}
//
//func (s *azfileLiveTestSuite) TestGetSetMetadataNonDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	metadata := Metadata{
//		"foo": "foovalue",
//		"bar": "barvalue",
//	}
//	setResp, err := fClient.SetMetadata(ctx, metadata)
//	_assert.Nil(err)
//	_assert.(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(setResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(setResp.Version, chk.Not(chk.Equals), "")
//	_assert.(setResp.Date.IsZero(), chk.Equals, false)
//	_assert.(setResp.IsServerEncrypted, chk.NotNil)
//
//	getResp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(getResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(getResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(getResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(getResp.Version, chk.Not(chk.Equals), "")
//	_assert.(getResp.Date.IsZero(), chk.Equals, false)
//	md := getResp.NewMetadata()
//	_assert.(md, chk.DeepEquals, metadata)
//}
//
//func (s *azfileLiveTestSuite) TestFileSetMetadataNil() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"not": "nil"})
//	_assert.Nil(err)
//
//	_, err = fClient.SetMetadata(ctx, nil)
//	_assert.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp.NewMetadata(), chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestFileSetMetadataDefaultEmpty() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"not": "nil"})
//	_assert.Nil(err)
//
//	_, err = fClient.SetMetadata(ctx, Metadata{})
//	_assert.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp.NewMetadata(), chk.HasLen, 0)
//}

//func (s *azfileLiveTestSuite) TestFileSetMetadataInvalidField() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"!@#$%^&*()": "!@#$%^&*()"})
//	_assert.NotNil(err)
//}
//
//func (s *azfileLiveTestSuite) TestStartCopyDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	srcFile, _ := createNewFileFromShare(c, srClient, 2048)
//	defer delFile(c, srcFile)
//
//	destFile, _ := getFileClientFromShare(c, srClient)
//	defer delFile(c, destFile)
//
//	_, err := srcFile.UploadRange(ctx, 0, getReaderToRandomBytes(2048), nil)
//	_assert.Nil(err)
//
//	copyResp, err := destFile.StartCopy(ctx, srcFile.URL(), nil)
//	_assert.Nil(err)
//	_assert.(copyResp.RawResponse.StatusCode, chk.Equals, 202)
//	_assert.(copyResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(copyResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(copyResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(copyResp.Version, chk.Not(chk.Equals), "")
//	_assert.(copyResp.Date.IsZero(), chk.Equals, false)
//	_assert.(copyResp.CopyID(), chk.Not(chk.Equals), "")
//	_assert.(copyResp.CopyStatus(), chk.Not(chk.Equals), "")
//
//	var copyStatus CopyStatusType
//	timeout := time.Duration(2) * time.Minute
//	start := time.Now()
//
//	var getResp *FileGetPropertiesResponse
//
//	for copyStatus != CopyStatusSuccess && time.Now().Sub(start) < timeout {
//		getResp, err = destFile.GetProperties(ctx)
//		_assert.Nil(err)
//		_assert.(getResp.CopyID(), chk.Equals, copyResp.CopyID())
//		_assert.(getResp.CopyStatus(), chk.Not(chk.Equals), CopyStatusNone)
//		_assert.(getResp.CopySource(), chk.Equals, srcFile.String())
//		copyStatus = getResp.CopyStatus()
//
//		time.Sleep(time.Duration(5) * time.Second)
//	}
//
//	if getResp != nil && getResp.CopyStatus() == CopyStatusSuccess {
//		// Abort will fail after copy finished
//		abortResp, err := destFile.AbortCopy(ctx, copyResp.CopyID())
//		_assert.NotNil(err)
//		_assert.(abortResp, chk.IsNil)
//		se, ok := err.(StorageError)
//		_assert.(ok, chk.Equals, true)
//		_assert.(se.RawResponse.StatusCode, chk.Equals, http.StatusConflict)
//	}
//}
//
//func waitForCopy(, copyfClient fClient, fileCopyResponse *FileStartCopyResponse) {
//	status := fileCopyResponse.CopyStatus()
//	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
//	start := time.Now()
//	for status != CopyStatusSuccess {
//		GetPropertiesResult, _ := copyfClient.GetProperties(ctx)
//		status = GetPropertiesResult.CopyStatus()
//		currentTime := time.Now()
//		if currentTime.Sub(start) >= time.Minute {
//			c.Fail()
//		}
//	}
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyDestEmpty() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShareWithDefaultData(c, srClient)
//	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	fileCopyResponse, err := copyfClient.StartCopy(ctx, fClient.URL(), nil)
//	_assert.Nil(err)
//	waitForCopy(c, copyfClient, fileCopyResponse)
//
//	resp, err := copyfClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//
//	// Read the file data to verify the copy
//	data, _ := ioutil.ReadAll(resp.Response().Body)
//	_assert.(resp.ContentLength, chk.Equals, int64(len(fileDefaultData)))
//	_assert.(string(data), chk.Equals, fileDefaultData)
//	resp.Response().Body.Close()
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyMetadata() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), basicMetadata)
//	_assert.Nil(err)
//	waitForCopy(c, copyfClient, resp)
//
//	resp2, err := copyfClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp2.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyMetadataNil() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
//	_, err := copyfClient.Create(ctx, 0, FileHTTPHeaders{}, basicMetadata)
//	_assert.Nil(err)
//
//	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), nil)
//	_assert.Nil(err)
//
//	waitForCopy(c, copyfClient, resp)
//
//	resp2, err := copyfClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp2.NewMetadata(), chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyMetadataEmpty() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	// Have the destination start with metadata so we ensure the empty metadata passed later takes effect
//	_, err := copyfClient.Create(ctx, 0, FileHTTPHeaders{}, basicMetadata)
//	_assert.Nil(err)
//
//	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), Metadata{})
//	_assert.Nil(err)
//
//	waitForCopy(c, copyfClient, resp)
//
//	resp2, err := copyfClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp2.NewMetadata(), chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyNegativeMetadataInvalidField() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	_, err := copyfClient.StartCopy(ctx, fClient.URL(), Metadata{"!@#$%^&*()": "!@#$%^&*()"})
//	_assert.NotNil(err)
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopySourceNonExistant() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	_, err := copyfClient.StartCopy(ctx, fClient.URL(), nil)
//	validateStorageError(c, err, ServiceCodeResourceNotFound)
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyUsingSASSrc() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(c, fsu)
//	defer delShare(_assert, srClient)
//	fClient, fileName := createNewFileFromShareWithDefaultData(c, srClient)
//
//	// Create sas values for the source file
//	credential, _ := getCredential()
//	serviceSASValues := FileSASSignatureValues{Version: "2015-04-05", StartTime: time.Now().Add(-1 * time.Hour).UTC(),
//		ExpiryTime: time.Now().Add(time.Hour).UTC(), Permissions: FileSASPermissions{Read: true, Write: true, Create: true, Delete: true}.String(),
//		ShareName: shareName, FilePath: fileName}
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
//	_assert.Nil(err)
//
//	// Create URLs to the destination file with sas parameters
//	sasURL := fClient.URL()
//	sasURL.RawQuery = queryParams.Encode()
//
//	// Create a new container for the destination
//	copysrClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(c, copysrClient, DeleteSnapshotsOptionNone)
//	copyfClient, _ := getFileClientFromShare(c, copysrClient)
//
//	resp, err := copyfClient.StartCopy(ctx, sasURL, nil)
//	_assert.Nil(err)
//
//	waitForCopy(c, copyfClient, resp)
//
//	resp2, err := copyfClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//
//	data, err := ioutil.ReadAll(resp2.Response().Body)
//	_assert.(resp2.ContentLength, chk.Equals, int64(len(fileDefaultData)))
//	_assert.(string(data), chk.Equals, fileDefaultData)
//	resp2.Response().Body.Close()
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyUsingSASDest() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(c, fsu)
//	defer delShare(_assert, srClient)
//	fClient, fileName := createNewFileFromShareWithDefaultData(c, srClient)
//	_ = fClient
//
//	// Generate SAS on the source
//	serviceSASValues := FileSASSignatureValues{ExpiryTime: time.Now().Add(time.Hour).UTC(),
//		Permissions: FileSASPermissions{Read: true, Write: true, Create: true}.String(), ShareName: shareName, FilePath: fileName}
//	credentials, _ := getCredential()
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credentials)
//	_assert.Nil(err)
//
//	copysrClient, copyShareName := createNewShare(c, fsu)
//	defer delShare(c, copysrClient, DeleteSnapshotsOptionNone)
//	copyfClient, copyFileName := getFileClientFromShare(c, copysrClient)
//
//	// Generate Sas for the destination
//	copyServiceSASvalues := FileSASSignatureValues{StartTime: time.Now().Add(-1 * time.Hour).UTC(),
//		ExpiryTime: time.Now().Add(time.Hour).UTC(), Permissions: FileSASPermissions{Read: true, Write: true}.String(),
//		ShareName: copyShareName, FilePath: copyFileName}
//	copyQueryParams, err := copyServiceSASvalues.NewSASQueryParameters(credentials)
//	_assert.Nil(err)
//
//	// Generate anonymous URL to destination with SAS
//	anonURL := fsu.URL()
//	anonURL.RawQuery = copyQueryParams.Encode()
//	anonPipeline := NewPipeline(NewAnonymousCredential(), PipelineOptions{})
//	anonFSU := NewServiceURL(anonURL, anonPipeline)
//	anonfClient := anonFSU.NewsrClient(copyShareName)
//	anonfClient := anonfClient.NewRootdirClient().NewfClient(copyFileName)
//
//	// Apply sas to source
//	srcFileWithSasURL := fClient.URL()
//	srcFileWithSasURL.RawQuery = queryParams.Encode()
//
//	resp, err := anonfClient.StartCopy(ctx, srcFileWithSasURL, nil)
//	_assert.Nil(err)
//
//	// Allow copy to happen
//	waitForCopy(c, anonfClient, resp)
//
//	resp2, err := copyfClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//
//	data, err := ioutil.ReadAll(resp2.Response().Body)
//	_, err = resp2.Body(RetryReaderOptions{}).Read(data)
//	_assert.(resp2.ContentLength, chk.Equals, int64(len(fileDefaultData)))
//	_assert.(string(data), chk.Equals, fileDefaultData)
//	resp2.Body(RetryReaderOptions{}).Close()
//}
//
//func (s *azfileLiveTestSuite) TestFileAbortCopyInProgress() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(c, fsu)
//	defer delShare(_assert, srClient)
//	fClient, fileName := getFileClientFromShare(c, srClient)
//
//	// Create a large file that takes time to copy
//	fileSize := 12 * 1024 * 1024
//	fileData := make([]byte, fileSize, fileSize)
//	for i := range fileData {
//		fileData[i] = byte('a' + i%26)
//	}
//	_, err := fClient.Create(ctx, int64(fileSize), FileHTTPHeaders{}, nil)
//	_assert.Nil(err)
//
//	_, err = fClient.UploadRange(ctx, 0, bytes.NewReader(fileData[0:4*1024*1024]), nil)
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 4*1024*1024, bytes.NewReader(fileData[4*1024*1024:8*1024*1024]), nil)
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 8*1024*1024, bytes.NewReader(fileData[8*1024*1024:]), nil)
//	_assert.Nil(err)
//	serviceSASValues := FileSASSignatureValues{ExpiryTime: time.Now().Add(time.Hour).UTC(),
//		Permissions: FileSASPermissions{Read: true, Write: true, Create: true}.String(), ShareName: shareName, FilePath: fileName}
//	credentials, _ := getCredential()
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credentials)
//	_assert.Nil(err)
//	srcFileWithSasURL := fClient.URL()
//	srcFileWithSasURL.RawQuery = queryParams.Encode()
//
//	fsu2, err := getAlternateFSU()
//	_assert.Nil(err)
//	copysrClient, _ := createNewShare(c, fsu2)
//	copyfClient, _ := getFileClientFromShare(c, copysrClient)
//
//	defer delShare(c, copysrClient, DeleteSnapshotsOptionNone)
//
//	resp, err := copyfClient.StartCopy(ctx, srcFileWithSasURL, nil)
//	_assert.Nil(err)
//	_assert.(resp.CopyStatus(), chk.Equals, CopyStatusPending)
//
//	_, err = copyfClient.AbortCopy(ctx, resp.CopyID())
//	if err != nil {
//		// If the error is nil, the test continues as normal.
//		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
//		_assert.((err.(StorageError)).RawResponse.StatusCode, chk.Equals, 409)
//		c.Error("The test failed because the copy completed because it was aborted")
//	}
//
//	resp2, _ := copyfClient.GetProperties(ctx)
//	_assert.(resp2.CopyStatus(), chk.Equals, CopyStatusAborted)
//}
//
//func (s *azfileLiveTestSuite) TestFileAbortCopyNoCopyStarted() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//
//	defer delShare(_assert, srClient)
//
//	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//	_, err := copyfClient.AbortCopy(ctx, "copynotstarted")
//	validateStorageError(c, err, ServiceCodeInvalidQueryParameterValue)
//}
//
//func (s *azfileLiveTestSuite) TestResizeFile() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 1234)
//
//	gResp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(gResp.ContentLength, chk.Equals, int64(1234))
//
//	rResp, err := fClient.Resize(ctx, 4096)
//	_assert.Nil(err)
//	_assert.(rResp.RawResponse.StatusCode, chk.Equals, 200)
//
//	gResp, err = fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(gResp.ContentLength, chk.Equals, int64(4096))
//}
//
//func (s *azfileLiveTestSuite) TestFileResizeZero() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(c, srClient, 10)
//
//	// The default file is created with size > 0, so this should actually update
//	_, err := fClient.Resize(ctx, 0)
//	_assert.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp.ContentLength, chk.Equals, int64(0))
//}
//
//func (s *azfileLiveTestSuite) TestFileResizeInvalidSizeNegative() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.Resize(ctx, -4)
//	_assert.NotNil(err)
//	sErr := err.(StorageError)
//	_assert.(sErr.RawResponse.StatusCode, chk.Equals, http.StatusBadRequest)
//}
//
//func (f *fClientSuite) TestServiceSASShareSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(c, fsu)
//	defer delShare(_assert, srClient)
//
//	credential, accountName := getCredential()
//
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:    SASProtocolHTTPS,
//		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour),
//		ShareName:   shareName,
//		Permissions: ShareSASPermissions{Create: true, Read: true, Write: true, Delete: true, List: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_assert.Nil(err)
//
//	qp := sasQueryParams.Encode()
//
//	fileName := "testFile"
//	dirName := "testDir"
//	fClientStr := fmt.Sprintf("https://%s.file.core.windows.net/%s/%s?%s",
//		accountName, shareName, fileName, qp)
//	fu, _ := url.Parse(fClientStr)
//
//	dirUrlStr := fmt.Sprintf("https://%s.file.core.windows.net/%s/%s?%s",
//		accountName, shareName, dirName, qp)
//	du, _ := url.Parse(dirUrlStr)
//
//	fClient := NewfClient(*fu, NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//	dirURL := NewdirClient(*du, NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//
//	s := "Hello"
//	_, err = fClient.Create(ctx, int64(len(s)), FileHTTPHeaders{}, Metadata{})
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 0, bytes.NewReader([]byte(s)), nil)
//	_assert.Nil(err)
//	_, err = fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_, err = fClient.Delete(ctx)
//	_assert.Nil(err)
//
//	_, err = dirURL.Create(ctx, Metadata{}, SMBProperties{})
//	_assert.Nil(err)
//
//	_, err = dirURL.ListFilesAndDirectoriesSegment(ctx, Marker{}, ListFilesAndDirectoriesOptions{})
//	_assert.Nil(err)
//}
//
//func (f *fClientSuite) TestServiceSASFileSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(c, fsu)
//	defer delShare(_assert, srClient)
//
//	credential, accountName := getCredential()
//
//	cacheControlVal := "cache-control-override"
//	contentDispositionVal := "content-disposition-override"
//	contentEncodingVal := "content-encoding-override"
//	contentLanguageVal := "content-language-override"
//	contentTypeVal := "content-type-override"
//
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:           SASProtocolHTTPS,
//		ExpiryTime:         time.Now().UTC().Add(48 * time.Hour),
//		ShareName:          shareName,
//		Permissions:        FileSASPermissions{Create: true, Read: true, Write: true, Delete: true}.String(),
//		CacheControl:       cacheControlVal,
//		ContentDisposition: contentDispositionVal,
//		ContentEncoding:    contentEncodingVal,
//		ContentLanguage:    contentLanguageVal,
//		ContentType:        contentTypeVal,
//	}.NewSASQueryParameters(credential)
//	_assert.Nil(err)
//
//	qp := sasQueryParams.Encode()
//
//	fileName := "testFile"
//	urlWithSAS := fmt.Sprintf("https://%s.file.core.windows.net/%s/%s?%s",
//		accountName, shareName, fileName, qp)
//	u, _ := url.Parse(urlWithSAS)
//
//	fClient := NewfClient(*u, NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//
//	s := "Hello"
//	_, err = fClient.Create(ctx, int64(len(s)), FileHTTPHeaders{}, Metadata{})
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 0, bytes.NewReader([]byte(s)), nil)
//	_assert.Nil(err)
//	dResp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_assert.(dResp.CacheControl, chk.Equals, cacheControlVal)
//	_assert.(dResp.ContentDisposition, chk.Equals, contentDispositionVal)
//	_assert.(dResp.ContentEncoding(), chk.Equals, contentEncodingVal)
//	_assert.(dResp.ContentLanguage, chk.Equals, contentLanguageVal)
//	_assert.(dResp.ContentType, chk.Equals, contentTypeVal)
//	_, err = fClient.Delete(ctx)
//	_assert.Nil(err)
//}
//
//func (s *azfileLiveTestSuite) TestDownloadEmptyZeroSizeFile() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//	defer delFile(c, fClient)
//
//	// Download entire fClient, check status code 200.
//	resp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_assert.(resp.RawResponse.StatusCode, chk.Equals, http.StatusOK)
//	_assert.(resp.ContentLength, chk.Equals, int64(0))
//	_assert.(resp.FileContentMD5, chk.IsNil) // Note: FileContentMD5 is returned, only when range is specified explicitly.
//
//	download, err := ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(download, chk.HasLen, 0)
//	_assert.(resp.AcceptRanges(), chk.Equals, "bytes")
//	_assert.(resp.CacheControl, chk.Equals, "")
//	_assert.(resp.ContentDisposition, chk.Equals, "")
//	_assert.(resp.ContentEncoding(), chk.Equals, "")
//	_assert.(resp.ContentRange(), chk.Equals, "") // Note: ContentRange is returned, only when range is specified explicitly.
//	_assert.(resp.ContentType, chk.Equals, "application/octet-stream")
//	_assert.(resp.CopyCompletionTime().IsZero(), chk.Equals, true)
//	_assert.(resp.CopyID(), chk.Equals, "")
//	_assert.(resp.CopyProgress(), chk.Equals, "")
//	_assert.(resp.CopySource(), chk.Equals, "")
//	_assert.(resp.CopyStatus(), chk.Equals, CopyStatusNone)
//	_assert.(resp.CopyStatusDescription(), chk.Equals, "")
//	_assert.(resp.Date.IsZero(), chk.Equals, false)
//	_assert.(resp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(resp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(resp.NewMetadata(), chk.DeepEquals, Metadata{})
//	_assert.(resp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(resp.Version, chk.Not(chk.Equals), "")
//	_assert.(resp.IsServerEncrypted, chk.NotNil)
//}
//
//func (s *azfileLiveTestSuite) TestUploadDownloadDefaultNonDefaultMD5() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 2048)
//	defer delFile(c, fClient)
//
//	contentR, contentD := getRandomDataAndReader(2048)
//
//	pResp, err := fClient.UploadRange(ctx, 0, contentR, nil)
//	_assert.Nil(err)
//	_assert.NotNil(pResp.ContentMD5,, chk.NotNil)
//	_assert.(pResp.RawResponse.StatusCode, chk.Equals, http.StatusCreated)
//	_assert.(pResp.IsServerEncrypted, chk.NotNil)
//	_assert.(pResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(pResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(pResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(pResp.Version, chk.Not(chk.Equals), "")
//	_assert.(pResp.Date.IsZero(), chk.Equals, false)
//
//	// Get with rangeGetContentMD5 enabled.
//	// Partial data, check status code 206.
//	resp, err := fClient.Download(ctx, 0, 1024, true)
//	_assert.Nil(err)
//	_assert.(resp.RawResponse.StatusCode, chk.Equals, http.StatusPartialContent)
//	_assert.(resp.ContentLength, chk.Equals, int64(1024))
//	_assert.(resp.ContentMD5, chk.NotNil)
//	_assert.(resp.ContentType, chk.Equals, "application/octet-stream")
//	_assert.(resp.Status(), chk.Not(chk.Equals), "")
//
//	download, err := ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(download, chk.DeepEquals, contentD[:1024])
//
//	// Set ContentMD5 for the entire file.
//	_, err = fClient.SetHTTPHeaders(ctx, FileHTTPHeaders{ContentMD5: pResp.ContentMD5,, ContentLanguage: "test"})
//	_assert.Nil(err)
//
//	// Test get with another type of range index, and validate if FileContentMD5 can be get correclty.
//	resp, err = fClient.Download(ctx, 1024, CountToEnd, false)
//	_assert.Nil(err)
//	_assert.Equal(resp.StatusCode(), chk.Equals, http.StatusPartialContent)
//	_assert.Equal(resp.ContentLength(), chk.Equals, int64(1024))
//	_assert.Nil(resp.ContentMD5(), chk.IsNil)
//	_assert.EqualValues(resp.FileContentMD5(), chk.DeepEquals, pResp.ContentMD5,)
//	_assert.(resp.ContentLanguage, chk.Equals, "test")
//	// Note: when it's downloading range, range's MD5 is returned, when set rangeGetContentMD5=true, currently set it to false, so should be empty
//	_assert.EqualValues(resp.NewMetadata(), chk.DeepEquals, FileHTTPHeaders{ContentLanguage: "test"})
//
//	download, err = ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(download, chk.DeepEquals, contentD[1024:])
//
//	_assert.(resp.AcceptRanges(), chk.Equals, "bytes")
//	_assert.(resp.CacheControl, chk.Equals, "")
//	_assert.(resp.ContentDisposition, chk.Equals, "")
//	_assert.(resp.ContentEncoding(), chk.Equals, "")
//	_assert.(resp.ContentRange(), chk.Equals, "bytes 1024-2047/2048")
//	_assert.(resp.ContentType, chk.Equals, "") // Note ContentType is set to empty during SetHTTPHeaders
//	_assert.(resp.CopyCompletionTime().IsZero(), chk.Equals, true)
//	_assert.(resp.CopyID(), chk.Equals, "")
//	_assert.(resp.CopyProgress(), chk.Equals, "")
//	_assert.(resp.CopySource(), chk.Equals, "")
//	_assert.(resp.CopyStatus(), chk.Equals, CopyStatusNone)
//	_assert.(resp.CopyStatusDescription(), chk.Equals, "")
//	_assert.(resp.Date.IsZero(), chk.Equals, false)
//	_assert.(resp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(resp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(resp.NewMetadata(), chk.DeepEquals, Metadata{})
//	_assert.(resp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(resp.Version, chk.Not(chk.Equals), "")
//	_assert.(resp.IsServerEncrypted, chk.NotNil)
//
//	// Get entire fClient, check status code 200.
//	resp, err = fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_assert.(resp.RawResponse.StatusCode, chk.Equals, http.StatusOK)
//	_assert.(resp.ContentLength, chk.Equals, int64(2048))
//	_assert.(resp.ContentMD5, chk.DeepEquals, pResp.ContentMD5,) // Note: This case is inted to get entire fClient, entire file's MD5 will be returned.
//	_assert.(resp.FileContentMD5, chk.IsNil)                      // Note: FileContentMD5 is returned, only when range is specified explicitly.
//
//	download, err = ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(download, chk.DeepEquals, contentD[:])
//
//	_assert.(resp.AcceptRanges(), chk.Equals, "bytes")
//	_assert.(resp.CacheControl, chk.Equals, "")
//	_assert.(resp.ContentDisposition, chk.Equals, "")
//	_assert.(resp.ContentEncoding(), chk.Equals, "")
//	_assert.(resp.ContentRange(), chk.Equals, "") // Note: ContentRange is returned, only when range is specified explicitly.
//	_assert.(resp.ContentType, chk.Equals, "")
//	_assert.(resp.CopyCompletionTime().IsZero(), chk.Equals, true)
//	_assert.(resp.CopyID(), chk.Equals, "")
//	_assert.(resp.CopyProgress(), chk.Equals, "")
//	_assert.(resp.CopySource(), chk.Equals, "")
//	_assert.(resp.CopyStatus(), chk.Equals, CopyStatusNone)
//	_assert.(resp.CopyStatusDescription(), chk.Equals, "")
//	_assert.(resp.Date.IsZero(), chk.Equals, false)
//	_assert.(resp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(resp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(resp.NewMetadata(), chk.DeepEquals, Metadata{})
//	_assert.(resp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(resp.Version, chk.Not(chk.Equals), "")
//	_assert.(resp.IsServerEncrypted, chk.NotNil)
//}
//
//func (s *azfileLiveTestSuite) TestFileDownloadDataNonExistantFile() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	_, err := fClient.Download(ctx, 0, CountToEnd, false)
//	validateStorageError(c, err, ServiceCodeResourceNotFound)
//}
//
//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileDownloadDataNegativeOffset() {
//// 	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//// 	defer delShare(_assert, srClient)
//// 	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.Download(ctx, -1, CountToEnd, false)
//// 	_assert.NotNil(err)
//// 	_assert.(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }
//
//func (s *azfileLiveTestSuite) TestFileDownloadDataOffsetOutOfRange() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.Download(ctx, int64(len(fileDefaultData)), CountToEnd, false)
//	validateStorageError(c, err, ServiceCodeInvalidRange)
//}
//
//// Don't check count by design.
//// func (s *azfileLiveTestSuite) TestFileDownloadDataInvalidCount() {
//// 	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//// 	defer delShare(_assert, srClient)
//// 	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.Download(ctx, 0, -100, false)
//// 	_assert.NotNil(err)
//// 	_assert.(strings.Contains(err.Error(), "count must be >= 0"), chk.Equals, true)
//// }
//
//func (s *azfileLiveTestSuite) TestFileDownloadDataEntireFile() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShareWithDefaultData(c, srClient)
//
//	resp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//
//	// Specifying a count of 0 results in the value being ignored
//	data, err := ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(string(data), chk.Equals, fileDefaultData)
//}
//
//func (s *azfileLiveTestSuite) TestFileDownloadDataCountExact() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShareWithDefaultData(c, srClient)
//
//	resp, err := fClient.Download(ctx, 0, int64(len(fileDefaultData)), false)
//	_assert.Nil(err)
//
//	data, err := ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(string(data), chk.Equals, fileDefaultData)
//}
//
//func (s *azfileLiveTestSuite) TestFileDownloadDataCountOutOfRange() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShareWithDefaultData(c, srClient)
//
//	resp, err := fClient.Download(ctx, 0, int64(len(fileDefaultData))*2, false)
//	_assert.Nil(err)
//
//	data, err := ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(string(data), chk.Equals, fileDefaultData)
//}
//
//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileUploadRangeNegativeInvalidOffset() {
//// 	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//// 	defer delShare(_assert, srClient)
//// 	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.UploadRange(ctx, -2, strings.NewReader(fileDefaultData), nil)
//// 	_assert.NotNil(err)
//// 	_assert.(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }
//
//func (s *azfileLiveTestSuite) TestFileUploadRangeNilBody() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.UploadRange(ctx, 0, nil, nil)
//	_assert.NotNil(err)
//	_assert.(strings.Contains(err.Error(), "body must not be nil"), chk.Equals, true)
//}
//
//func (s *azfileLiveTestSuite) TestFileUploadRangeEmptyBody() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.UploadRange(ctx, 0, bytes.NewReader([]byte{}), nil)
//	_assert.NotNil(err)
//	_assert.(strings.Contains(err.Error(), "body must contain readable data whose size is > 0"), chk.Equals, true)
//}
//
//func (s *azfileLiveTestSuite) TestFileUploadRangeNonExistantFile() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	_, err := fClient.UploadRange(ctx, 0, getReaderToRandomBytes(12), nil)
//	validateStorageError(c, err, ServiceCodeResourceNotFound)
//}
//
//func (s *azfileLiveTestSuite) TestFileUploadRangeTransactionalMD5() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 2048)
//	defer delFile(c, fClient)
//
//	contentR, contentD := getRandomDataAndReader(2048)
//	md5 := md5.Sum(contentD)
//
//	// Upload range with correct transactional MD5
//	pResp, err := fClient.UploadRange(ctx, 0, contentR, md5[:])
//	_assert.Nil(err)
//	_assert.(pResp.ContentMD5,, chk.NotNil)
//	_assert.(pResp.RawResponse.StatusCode, chk.Equals, http.StatusCreated)
//	_assert.(pResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(pResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(pResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(pResp.Version, chk.Not(chk.Equals), "")
//	_assert.(pResp.Date.IsZero(), chk.Equals, false)
//	_assert.(pResp.ContentMD5,, chk.DeepEquals, md5[:])
//
//	// Upload range with empty MD5, nil MD5 is covered by other cases.
//	pResp, err = fClient.UploadRange(ctx, 1024, bytes.NewReader(contentD[1024:]), nil)
//	_assert.Nil(err)
//	_assert.(pResp.ContentMD5,, chk.NotNil)
//	_assert.(pResp.RawResponse.StatusCode, chk.Equals, http.StatusCreated)
//
//	resp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_assert.(resp.RawResponse.StatusCode, chk.Equals, http.StatusOK)
//	_assert.(resp.ContentLength, chk.Equals, int64(2048))
//
//	download, err := ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(download, chk.DeepEquals, contentD[:])
//}
//
//func (s *azfileLiveTestSuite) TestFileUploadRangeIncorrectTransactionalMD5() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 2048)
//	defer delFile(c, fClient)
//
//	contentR, _ := getRandomDataAndReader(2048)
//	_, incorrectMD5 := getRandomDataAndReader(16)
//
//	// Upload range with incorrect transactional MD5
//	_, err := fClient.UploadRange(ctx, 0, contentR, incorrectMD5[:])
//	validateStorageError(c, err, ServiceCodeMd5Mismatch)
//}
//
//func (f *fClientSuite) TestUploadRangeFromURL() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(c, fsu)
//	defer delShare(_assert, srClient)
//
//	// create the source file and populate it with random data at a specific offset
//	expectedDataSize := 2048
//	totalFileSize := 4096
//	srcOffset := 999
//	expectedDataReader, expectedData := getRandomDataAndReader(expectedDataSize)
//	srcfClient, _ := createNewFileFromShare(c, srClient, int64(totalFileSize))
//	_, err := srcfClient.UploadRange(ctx, int64(srcOffset), expectedDataReader, nil)
//	_assert.Nil(err)
//
//	// generate a URL with SAS pointing to the source file
//	credential, _ := getCredential()
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:    SASProtocolHTTPS,
//		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour),
//		ShareName:   shareName,
//		Permissions: FileSASPermissions{Create: true, Read: true, Write: true, Delete: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_assert.Nil(err)
//	rawSrcURL := srcfClient.URL()
//	rawSrcURL.RawQuery = sasQueryParams.Encode()
//
//	// create the destination file
//	dstfClient, _ := createNewFileFromShare(c, srClient, int64(totalFileSize))
//
//	// invoke UploadRange on dstfClient and put the data at a random range
//	// source and destination have different offsets so we can test both values at the same time
//	dstOffset := 100
//	uploadFromURLResp, err := dstfClient.UploadRangeFromURL(ctx, rawSrcURL, int64(srcOffset),
//		int64(dstOffset), int64(expectedDataSize))
//	_assert.Nil(err)
//	_assert.(uploadFromURLResp.RawResponse.StatusCode, chk.Equals, 201)
//
//	// verify the destination
//	resp, err := dstfClient.Download(ctx, int64(dstOffset), int64(expectedDataSize), false)
//	_assert.Nil(err)
//	download, err := ioutil.ReadAll(resp.Response().Body)
//	_assert.Nil(err)
//	_assert.(download, chk.DeepEquals, expectedData)
//}
//
//// Testings for GetRangeList and ClearRange
//func (s *azfileLiveTestSuite) TestGetRangeListNonDefaultExact() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	fileSize := int64(512 * 10)
//
//	fClient.Create(ctx, fileSize, FileHTTPHeaders{}, nil)
//
//	defer delFile(c, fClient)
//
//	putResp, err := fClient.UploadRange(ctx, 0, getReaderToRandomBytes(1024), nil)
//	_assert.Nil(err)
//	_assert.(putResp.RawResponse.StatusCode, chk.Equals, 201)
//	_assert.(putResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(putResp.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(putResp.ContentMD5,, chk.NotNil)
//	_assert.(putResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(putResp.Version, chk.Not(chk.Equals), "")
//	_assert.(putResp.Date.IsZero(), chk.Equals, false)
//
//	rangeList, err := fClient.GetRangeList(ctx, 0, 1023)
//	_assert.Nil(err)
//	_assert.(rangeList.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(rangeList.LastModified.IsZero(), chk.Equals, false)
//	_assert.(rangeList.ETag, chk.Not(chk.Equals), ETagNone)
//	_assert.(rangeList.FileContentLength, chk.Equals, fileSize)
//	_assert.(rangeList.RequestID, chk.Not(chk.Equals), "")
//	_assert.(rangeList.Version, chk.Not(chk.Equals), "")
//	_assert.(rangeList.Date.IsZero(), chk.Equals, false)
//	_assert.(rangeList.Ranges, chk.HasLen, 1)
//	_assert.(rangeList.Ranges[0], chk.DeepEquals, FileRange{XMLName: xml.Name{Space: "", Local: "Range"}, Start: 0, End: 1023})
//}
//
//// Default means clear the entire file's range
//func (s *azfileLiveTestSuite) TestClearRangeDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 2048)
//	defer delFile(c, fClient)
//
//	_, err := fClient.UploadRange(ctx, 0, getReaderToRandomBytes(2048), nil)
//	_assert.Nil(err)
//
//	clearResp, err := fClient.ClearRange(ctx, 0, 2048)
//	_assert.Nil(err)
//	_assert.(clearResp.RawResponse.StatusCode, chk.Equals, 201)
//
//	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd)
//	_assert.Nil(err)
//	_assert.(rangeList.Ranges, chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestClearRangeNonDefault() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 4096)
//	defer delFile(c, fClient)
//
//	_, err := fClient.UploadRange(ctx, 2048, getReaderToRandomBytes(2048), nil)
//	_assert.Nil(err)
//
//	clearResp, err := fClient.ClearRange(ctx, 2048, 2048)
//	_assert.Nil(err)
//	_assert.(clearResp.RawResponse.StatusCode, chk.Equals, 201)
//
//	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd)
//	_assert.Nil(err)
//	_assert.(rangeList.Ranges, chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestClearRangeMultipleRanges() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 2048)
//	defer delFile(c, fClient)
//
//	_, err := fClient.UploadRange(ctx, 0, getReaderToRandomBytes(2048), nil)
//	_assert.Nil(err)
//
//	clearResp, err := fClient.ClearRange(ctx, 1024, 1024)
//	_assert.Nil(err)
//	_assert.(clearResp.RawResponse.StatusCode, chk.Equals, 201)
//
//	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd)
//	_assert.Nil(err)
//	_assert.(rangeList.Ranges, chk.HasLen, 1)
//	_assert.(rangeList.Ranges[0], chk.DeepEquals, FileRange{XMLName: xml.Name{Space: "", Local: "Range"}, Start: 0, End: 1023})
//}
//
//// When not 512 aligned, clear range will set 0 the non-512 aligned range, and will not eliminate the range.
//func (s *azfileLiveTestSuite) TestClearRangeNonDefault1Count() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//
//	fClient, _ := createNewFileFromShare(c, srClient, 1)
//	defer delFile(c, fClient)
//
//	d := []byte{1}
//	_, err := fClient.UploadRange(ctx, 0, bytes.NewReader(d), nil)
//	_assert.Nil(err)
//
//	clearResp, err := fClient.ClearRange(ctx, 0, 1)
//	_assert.Nil(err)
//	_assert.(clearResp.RawResponse.StatusCode, chk.Equals, 201)
//
//	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd)
//	_assert.Nil(err)
//	_assert.(rangeList.Ranges, chk.HasLen, 1)
//	_assert.(rangeList.Ranges[0], chk.DeepEquals, FileRange{XMLName: xml.Name{Space: "", Local: "Range"}, Start: 0, End: 0})
//
//	dResp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	bytes, err := ioutil.ReadAll(dResp.Body(RetryReaderOptions{}))
//	_assert.Nil(err)
//	_assert.(bytes, chk.DeepEquals, []byte{0})
//}
//
//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileClearRangeNegativeInvalidOffset() {
//// 	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient, _ := getsrClient(c, fsu)
//// 	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//// 	_, err := fClient.ClearRange(ctx, -1, 1)
//// 	_assert.NotNil(err)
//// 	_assert.(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }
//
//func (s *azfileLiveTestSuite) TestFileClearRangeNegativeInvalidCount() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, _ := getsrClient(c, fsu)
//	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//
//	_, err := fClient.ClearRange(ctx, 0, 0)
//	_assert.NotNil(err)
//	_assert.(strings.Contains(err.Error(), "count cannot be CountToEnd, and must be > 0"), chk.Equals, true)
//}
//
//func setupGetRangeListTest() (srClient srClient, fClient fClient) {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, _ = createNewShare(c, fsu)
//	fClient, _ = createNewFileFromShare(c, srClient, int64(testFileRangeSize))
//
//	_, err := fClient.UploadRange(ctx, 0, getReaderToRandomBytes(testFileRangeSize), nil)
//	_assert.Nil(err)
//
//	return
//}
//
//func validateBasicGetRangeList(, resp *ShareFileRangeList, err error) {
//	_assert.Nil(err)
//	_assert.(resp.Ranges, chk.HasLen, 1)
//	_assert.(resp.Ranges[0], chk.Equals, FileRange{XMLName: xml.Name{Space: "", Local: "Range"}, Start: 0, End: testFileRangeSize - 1})
//}
//
//func (s *azfileLiveTestSuite) TestFileGetRangeListDefaultEmptyFile() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient)
//	fClient, _ := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd)
//	_assert.Nil(err)
//	_assert.(resp.Ranges, chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestFileGetRangeListDefault1Range() {
//	srClient, fClient := setupGetRangeListTest(c)
//	defer delShare(_assert, srClient)
//
//	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd)
//	validateBasicGetRangeList(c, resp, err)
//}
//
//func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRanges() {
//	srClient, fClient := setupGetRangeListTest(c)
//	defer delShare(_assert, srClient)
//
//	_, err := fClient.Resize(ctx, int64(testFileRangeSize*3))
//	_assert.Nil(err)
//
//	_, err = fClient.UploadRange(ctx, testFileRangeSize*2, getReaderToRandomBytes(testFileRangeSize), nil)
//	_assert.Nil(err)
//	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd)
//	_assert.Nil(err)
//	_assert.(resp.Ranges, chk.HasLen, 2)
//	_assert.(resp.Ranges[0], chk.Equals, FileRange{XMLName: xml.Name{Space: "", Local: "Range"}, Start: 0, End: testFileRangeSize - 1})
//	_assert.(resp.Ranges[1], chk.Equals, FileRange{XMLName: xml.Name{Space: "", Local: "Range"}, Start: testFileRangeSize * 2, End: (testFileRangeSize * 3) - 1})
//}
//
//func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRangesCountLess() {
//	srClient, fClient := setupGetRangeListTest(c)
//	defer delShare(_assert, srClient)
//
//	resp, err := fClient.GetRangeList(ctx, 0, testFileRangeSize-1)
//	_assert.Nil(err)
//	_assert.(resp.Ranges, chk.HasLen, 1)
//	_assert.(resp.Ranges[0], chk.Equals, FileRange{XMLName: xml.Name{Space: "", Local: "Range"}, Start: 0, End: testFileRangeSize - 1})
//}
//
//func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRangesCountExceed() {
//	srClient, fClient := setupGetRangeListTest(c)
//	defer delShare(_assert, srClient)
//
//	resp, err := fClient.GetRangeList(ctx, 0, testFileRangeSize+1)
//	_assert.Nil(err)
//	validateBasicGetRangeList(c, resp, err)
//}
//
//func (s *azfileLiveTestSuite) TestFileGetRangeListSnapshot() {
//	srClient, fClient := setupGetRangeListTest(c)
//	defer delShare(c, srClient, DeleteSnapshotsOptionInclude)
//
//	resp, _ := srClient.CreateSnapshot(ctx, Metadata{})
//	snapshotURL := fClient.WithSnapshot(resp.Snapshot())
//	resp2, err := snapshotURL.GetRangeList(ctx, 0, CountToEnd)
//	_assert.Nil(err)
//	validateBasicGetRangeList(c, resp2, err)
//}
//
//func (s *azfileLiveTestSuite) TestUnexpectedEOFRecovery() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	share, _ := createNewShare(c, fsu)
//	defer delShare(c, share, DeleteSnapshotsOptionInclude)
//
//	fClient, _ := createNewFileFromShare(c, share, 2048)
//
//	contentR, contentD := getRandomDataAndReader(2048)
//
//	resp, err := fClient.UploadRange(ctx, 0, contentR, nil)
//	_assert.Nil(err)
//	_assert.(resp.RawResponse.StatusCode, chk.Equals, http.StatusCreated)
//	_assert.(resp.RequestID, chk.Not(chk.Equals), "")
//
//	dlResp, err := fClient.Download(ctx, 0, 2048, false)
//	_assert.Nil(err)
//
//	// Verify that we can inject errors first.
//	reader := dlResp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))
//
//	_, err = ioutil.ReadAll(reader)
//	_assert.NotNil(err)
//	_assert.(err.Error(), chk.Equals, "unrecoverable error")
//
//	// Then inject the retryable error.
//	reader = dlResp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))
//
//	buf, err := ioutil.ReadAll(reader)
//	_assert.Nil(err)
//	_assert.(buf, chk.DeepEquals, contentD)
//}
//
//func (s *azfileLiveTestSuite) TestCreateMaximumSizeFileShare() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	srClient, err := getsrClient(generateShareName(testName), svcClient)
//	cResp, err := srClient.Create(ctx, nil, fileShareMaxQuota)
//	_assert.Nil(err)
//	_assert.(cResp.RawResponse.StatusCode, chk.Equals, 201)
//	defer delShare(c, share, DeleteSnapshotsOptionInclude)
//	dir := share.NewRootdirClient()
//
//	fClient := getFileClientFromDirectory(_assert, generateFileName(testName), dir)
//
//	_, err = fClient.Create(ctx, fileMaxAllowedSizeInBytes, FileHTTPHeaders{}, nil)
//	_assert.Nil(err)
//	_assert.(cResp.RawResponse.StatusCode, chk.Equals, 201)
//}
