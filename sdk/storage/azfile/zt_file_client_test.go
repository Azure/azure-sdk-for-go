// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"bytes"
	"crypto/md5"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	testFileRangeSize         = 512                  // Use this number considering clear range's function
	fileShareMaxQuota         = int32(5120)          // Size is in GB (Service Version 2020-02-10)
	fileMaxAllowedSizeInBytes = int64(4398046511104) // 4 TiB (Service Version 2020-02-10)
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
	defer delShare(_assert, srClient, nil)

	// Create and delete fClient in root directory.
	fileName := generateFileName(testName)
	dirClient1, err := srClient.NewRootDirectoryClient()
	_assert.Nil(err)
	fClient, err := dirClient1.NewFileClient(fileName)
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

	// Create and delete fClient in named directory.
	afClient, err := dirClient.NewFileClient(generateFileName(testName))
	_assert.Nil(err)

	cResp, err = afClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	_assert.NotEqual(cResp.ETag, "")
	_assert.Equal(cResp.LastModified.IsZero(), false)
	_assert.NotEqual(cResp.RequestID, "")
	_assert.NotEqual(cResp.Version, "")
	_assert.Equal(cResp.Date.IsZero(), false)
	_assert.NotNil(cResp.IsServerEncrypted)

	delResp, err = afClient.Delete(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(delResp.RawResponse.StatusCode, 202)
	_assert.NotEqual(delResp.RequestID, "")
	_assert.NotEqual(delResp.Version, "")
	_assert.Equal(delResp.Date.IsZero(), false)
}

func (s *azfileLiveTestSuite) TestFileCreateNonDefaultMetadataNonEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	_, err = fClient.Create(ctx, &CreateFileOptions{
		Metadata: basicMetadata,
	})
	_assert.Nil(err)

	resp, err := fClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.Metadata, len(basicMetadata))
}

func (s *azfileLiveTestSuite) TestFileCreateNonDefaultHTTPHeaders() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
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
	defer delShare(_assert, srClient, nil)
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
//	defer delShare(_assert, srClient, nil)
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
//	options := &SetFileHTTPHeadersOptions{
//		Permissions: &Permissions{ PermissionStr: &sampleSDDL},
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
//	_assert.NotEqual(setResp.ETag, "")
//	_assert.Equal(setResp.LastModified.IsZero(),  false)
//	_assert.NotEqual(setResp.RequestID, "")
//	_assert.NotEqual(setResp.Version, "")
//	_assert.Equal(setResp.Date.IsZero(),  false)
//	_assert.NotNil(setResp.IsServerEncrypted)
//
//	getResp, err := fClient.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert.Equal(getResp.RawResponse.StatusCode,  200)
//	_assert.Equal(setResp.LastModified.IsZero(),  false)
//	_assert.Equal(*getResp.FileType,  "File")
//
//	_assert.EqualValues(getResp.ContentType,  options.FileHTTPHeaders.FileContentType)
//	_assert.EqualValues(getResp.ContentEncoding,  options.FileHTTPHeaders.FileContentEncoding)
//	_assert.EqualValues(getResp.ContentLanguage,  options.FileHTTPHeaders.FileContentLanguage)
//	_assert.EqualValues(getResp.ContentMD5, options.FileHTTPHeaders.FileContentMD5)
//	_assert.EqualValues(getResp.CacheControl,  options.FileHTTPHeaders.FileCacheControl)
//	_assert.EqualValues(getResp.ContentDisposition,  options.FileHTTPHeaders.FileContentDisposition)
//	_assert.Equal(*getResp.ContentLength,  int64(0))
//	// We'll just ensure a permission exists, no need to test overlapping functionality.
//	_assert.NotEqual(*getResp.PermissionKey, "")
//	// Ensure our attributes and other properties (after parsing) are equivalent to our original
//	// There's an overlapping test for this in ntfs_property_bitflags_test.go, but it doesn't hurt to test it alongside other things.
//	_assert.EqualValues(ParseFileAttributeFlagsString(*getResp.FileAttributes),  *options.SMBProperties.FileAttributes)
//
//	fct, _ := time.Parse(ISO8601, *getResp.FileCreationTime)
//	_assert.EqualValues(fct, creationTime)
//	fwt, _ := time.Parse(ISO8601, *getResp.FileLastWriteTime)
//	_assert.Equal(fwt, lastWriteTime)
//
//	_assert.NotEqual(getResp.ETag, "")
//	_assert.NotEqual(getResp.RequestID, "")
//	_assert.NotEqual(getResp.Version, "")
//	_assert.Equal(getResp.Date.IsZero(),  false)
//	_assert.NotNil(getResp.IsServerEncrypted)
//}

//func (s *azfileLiveTestSuite) TestFilePreservePermissions() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
//
//	fClient, _ := createNewFileFromShareWithPermissions(c, srClient, 0)
//	defer delFile(_assert, fClient)
//
//	// Grab the original perm key before we set file headers.
//	getResp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//
//	oKey := getResp.PermissionKey()
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
//		FileContentType:        to.StringPtr("text/html"),
//		FileContentEncoding:        to.StringPtr("gzip"),
//		FileContentLanguage:        to.StringPtr("tr,en"),
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
//	_assert.(setResp.ETag, chk.Not(chk.Equals), "")
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
//	_assert.(getResp.ContentEncoding, chk.Equals, properties.ContentEncoding)
//	_assert.(getResp.ContentLanguage, chk.Equals, properties.ContentLanguage)
//	_assert.(getResp.ContentMD5,, chk.DeepEquals, properties.ContentMD5)
//	_assert.(getResp.CacheControl, chk.Equals, properties.CacheControl)
//	_assert.(getResp.ContentDisposition, chk.Equals, properties.ContentDisposition)
//	_assert.(getResp.ContentLength, chk.Equals, int64(0))
//	// Ensure that the permission key gets preserved
//	_assert.(getResp.PermissionKey(), chk.Equals, oKey)
//	timeAdapter = SMBPropertyAdapter{PropertySource: getResp}
//	c.Log("Original last write time: ", lwTime, " new time: ", timeAdapter.FileLastWriteTime())
//	_assert.(timeAdapter.FileLastWriteTime().Equal(lwTime), chk.Equals, true)
//	c.Log("Original creation time: ", cTime, " new time: ", timeAdapter.FileCreationTime())
//	_assert.(timeAdapter.FileCreationTime().Equal(cTime), chk.Equals, true)
//	_assert.(getResp.FileAttributes(), chk.Equals, attribs)
//
//	_assert.(getResp.ETag, chk.Not(chk.Equals), "")
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
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	md5Str := "MDAwMDAwMDA="
//	var testMd5 []byte
//	copy(testMd5[:], md5Str)
//
//	properties := FileHTTPHeaders{
//		FileContentType:        to.StringPtr("text/html"),
//		FileContentEncoding:        to.StringPtr("gzip"),
//		FileContentLanguage:        to.StringPtr("tr,en"),
//		FileContentMD5:         testMd5,
//		FileCacheControl:       "no-transform",
//		ContentDisposition: "attachment",
//	}
//	setResp, err := fClient.SetHTTPHeaders(ctx, properties)
//	_assert.Nil(err)
//	_assert.(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.ETag, chk.Not(chk.Equals), "")
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
//	resp, _ := srClient.CreateSnapshot(ctx, map[string]string)
//	snapshotURL := fClient.WithSnapshot(resp.Snapshot())
//
//	getResp, err := snapshotURL.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(getResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(getResp.FileType(), chk.Equals, "File")
//
//	_assert.(getResp.ContentType, chk.Equals, properties.ContentType)
//	_assert.(getResp.ContentEncoding, chk.Equals, properties.ContentEncoding)
//	_assert.(getResp.ContentLanguage, chk.Equals, properties.ContentLanguage)
//	_assert.(getResp.ContentMD5,, chk.DeepEquals, properties.ContentMD5)
//	_assert.(getResp.CacheControl, chk.Equals, properties.CacheControl)
//	_assert.(getResp.ContentDisposition, chk.Equals, properties.ContentDisposition)
//	_assert.(getResp.ContentLength, chk.Equals, int64(0))
//
//	_assert.(getResp.ETag, chk.Not(chk.Equals), "")
//	_assert.(getResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(getResp.Version, chk.Not(chk.Equals), "")
//	_assert.(getResp.Date.IsZero(), chk.Equals, false)
//	_assert.(getResp.IsServerEncrypted, chk.NotNil)
//	_assert.(getResp.Metadata, chk.DeepEquals, metadata)
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
//	defer delShare(_assert, srClient, nil)
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	metadata := Metadata{
//		"foo": "foovalue",
//		"bar": "barvalue",
//	}
//	setResp, err := fClient.SetMetadata(ctx, metadata)
//	_assert.Nil(err)
//	_assert.(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(setResp.ETag, chk.Not(chk.Equals), "")
//	_assert.(setResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(setResp.Version, chk.Not(chk.Equals), "")
//	_assert.(setResp.Date.IsZero(), chk.Equals, false)
//	_assert.(setResp.IsServerEncrypted, chk.NotNil)
//
//	getResp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(getResp.ETag, chk.Not(chk.Equals), "")
//	_assert.(getResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(getResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(getResp.Version, chk.Not(chk.Equals), "")
//	_assert.(getResp.Date.IsZero(), chk.Equals, false)
//	md := getResp.Metadata
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
//	defer delShare(_assert, srClient, nil)
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"not": "nil"})
//	_assert.Nil(err)
//
//	_, err = fClient.SetMetadata(ctx, nil)
//	_assert.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp.Metadata, chk.HasLen, 0)
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
//	defer delShare(_assert, srClient, nil)
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"not": "nil"})
//	_assert.Nil(err)
//
//	_, err = fClient.SetMetadata(ctx, map[string]string)
//	_assert.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_assert.Nil(err)
//	_assert.(resp.Metadata, chk.HasLen, 0)
//}

//func (s *azfileLiveTestSuite) TestFileSetMetadataInvalidField() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
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
//	defer delShare(_assert, srClient, nil)
//
//	srcFile, _ := createNewFileFromShare(_assert, generateFileName(testName), 2048, srClient)
//	defer delFile(c, srcFile)
//
//	destFile, _ := getFileClientFromShare(_assert, generateFileName(testName), srClient)
//	defer delFile(c, destFile)
//
//	_, err := srcFile.UploadRange(ctx, 0, generateData(2048), nil)
//	_assert.Nil(err)
//
//	copyResp, err := destFile.StartCopy(ctx, srcFile.URL(), nil)
//	_assert.Nil(err)
//	_assert.(copyResp.RawResponse.StatusCode, chk.Equals, 202)
//	_assert.(copyResp.ETag, chk.Not(chk.Equals), "")
//	_assert.(copyResp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(copyResp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(copyResp.Version, chk.Not(chk.Equals), "")
//	_assert.(copyResp.Date.IsZero(), chk.Equals, false)
//	_assert.(copyResp.CopyID, chk.Not(chk.Equals), "")
//	_assert.(copyResp.CopyStatus, chk.Not(chk.Equals), "")
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
//		_assert.(getResp.CopyID, chk.Equals, copyResp.CopyID())
//		_assert.(getResp.CopyStatus, chk.Not(chk.Equals), "")
//		_assert.(getResp.CopySource, chk.Equals, srcFile.String())
//		copyStatus = getResp.CopyStatus
//
//		time.Sleep(time.Duration(5) * time.Second)
//	}
//
//	if getResp != nil && getResp.CopyStatus == CopyStatusSuccess {
//		// Abort will fail after copy finished
//		abortResp, err := destFile.AbortCopy(ctx, copyResp.CopyID())
//		_assert.NotNil(err)
//		_assert.(abortResp, chk.IsNil)
//		se, ok := err.(StorageError)
//		_assert.(ok, chk.Equals, true)
//		_assert.(se.RawResponse.StatusCode, chk.Equals, http.StatusConflict)
//	}
//}

func waitForCopy(_assert *assert.Assertions, copyfClient FileClient, fileCopyResponse FileStartCopyResponse) {
	status := fileCopyResponse.CopyStatus
	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
	start := time.Now()
	for *status != CopyStatusTypeSuccess {
		GetPropertiesResult, _ := copyfClient.GetProperties(ctx, nil)
		status = GetPropertiesResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_assert.Fail("")
		}
	}
}

func (s *azfileLiveTestSuite) TestFileStartCopyDestEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_assert, "src"+generateFileName(testName), fileDefaultData, srClient)
	copyfClient := getFileClientFromShare(_assert, "dst"+generateFileName(testName), srClient)

	fileCopyResponse, err := copyfClient.StartCopy(ctx, fClient.URL(), nil)
	_assert.Nil(err)
	waitForCopy(_assert, copyfClient, fileCopyResponse)

	resp, err := copyfClient.Download(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)

	// Read the file data to verify the copy
	data, _ := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Equal(*resp.ContentLength, int64(len(fileDefaultData)))
	_assert.Equal(string(data), fileDefaultData)
	resp.RawResponse.Body.Close()
}

//func (s *azfileLiveTestSuite) TestFileStartCopyMetadata() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
//	fClient := createNewFileFromShare(_assert, "src" + generateFileName(testName), 0, srClient)
//	copyfClient := getFileClientFromShare(_assert, "dst" + generateFileName(testName), srClient)
//
//	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), &StartFileCopyOptions{Metadata: basicMetadata})
//	_assert.Nil(err)
//	waitForCopy(_assert, copyfClient, resp)
//
//	resp2, err := copyfClient.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert.EqualValues(resp2.Metadata, basicMetadata)
//}

func (s *azfileLiveTestSuite) TestFileStartCopyMetadataNil() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShare(_assert, "src"+generateFileName(testName), 0, srClient)
	copyfClient := getFileClientFromShare(_assert, "dst"+generateFileName(testName), srClient)

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err = copyfClient.Create(ctx, &CreateFileOptions{Metadata: basicMetadata})
	_assert.Nil(err)

	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), nil)
	_assert.Nil(err)

	waitForCopy(_assert, copyfClient, resp)

	resp2, err := copyfClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp2.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestFileStartCopyMetadataEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShare(_assert, "src"+generateFileName(testName), 0, srClient)
	copyfClient := getFileClientFromShare(_assert, "dst"+generateFileName(testName), srClient)

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err = copyfClient.Create(ctx, &CreateFileOptions{Metadata: basicMetadata})
	_assert.Nil(err)

	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), &StartFileCopyOptions{Metadata: map[string]string{}})
	_assert.Nil(err)

	waitForCopy(_assert, copyfClient, resp)

	resp2, err := copyfClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp2.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestFileStartCopyNegativeMetadataInvalidField() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShare(_assert, "src"+generateFileName(testName), 0, srClient)
	copyfClient := getFileClientFromShare(_assert, "dst"+generateFileName(testName), srClient)

	_, err = copyfClient.StartCopy(ctx, fClient.URL(), &StartFileCopyOptions{Metadata: map[string]string{"!@#$%^&*()": "!@#$%^&*()"}})
	_assert.NotNil(err)
}

func (s *azfileLiveTestSuite) TestFileStartCopySourceNonExistent() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := getFileClientFromShare(_assert, "src"+generateFileName(testName), srClient)
	copyfClient := getFileClientFromShare(_assert, "dst"+generateFileName(testName), srClient)

	_, err = copyfClient.StartCopy(ctx, fClient.URL(), nil)
	validateStorageError(_assert, err, StorageErrorCodeResourceNotFound)
}

//func (s *azfileLiveTestSuite) TestFileStartCopyUsingSASSrc() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
//	fClient, fileName := createNewFileFromShareWithGivenData(c, srClient)
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
//	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
//	_assert.(resp2.ContentLength, chk.Equals, int64(len(fileDefaultData)))
//	_assert.(string(data), chk.Equals, fileDefaultData)
//	resp2.RawResponse.Body.Close()
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyUsingSASDest() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
//	fClient, fileName := createNewFileFromShareWithGivenData(c, srClient)
//	_ = fClient
//
//	// Generate SAS on the source
//	serviceSASValues := FileSASSignatureValues{ExpiryTime: time.Now().Add(time.Hour).UTC(),
//		Permissions: FileSASPermissions{Read: true, Write: true, Create: true}.String(), ShareName: shareName, FilePath: fileName}
//	credentials, _ := getCredential()
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credentials)
//	_assert.Nil(err)
//
//	copysrClient, copyShareName := createNewShare(_assert, generateShareName(testName), svcClient)
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
//	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
//	_, err = resp2.Body(RetryReaderOptions{}).Read(data)
//	_assert.(resp2.ContentLength, chk.Equals, int64(len(fileDefaultData)))
//	_assert.(string(data), chk.Equals, fileDefaultData)
//	resp2.Body(RetryReaderOptions{}).Close()
//}

//func (s *azfileLiveTestSuite) TestFileAbortCopyInProgress() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_assert, shareName, svcClient)
//	defer delShare(_assert, srClient, nil)
//	fileName := generateFileName(testName)
//	fClient := getFileClientFromShare(_assert, fileName, srClient)
//
//	// Create a large file that takes time to copy
//	fileSize := 12 * 1024 * 1024
//	fileData := make([]byte, fileSize, fileSize)
//	for i := range fileData {
//		fileData[i] = byte('a' + i%26)
//	}
//	_, err = fClient.Create(ctx, &CreateFileOptions{FileContentLength: to.Int64Ptr(int64(fileSize)), FileHTTPHeaders: &FileHTTPHeaders{}})
//	_assert.Nil(err)
//
//	_, err = fClient.UploadRange(ctx, 0, internal.NopCloser(bytes.NewReader(fileData[0:4*1024*1024])), nil)
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 4*1024*1024, internal.NopCloser(bytes.NewReader(fileData[4*1024*1024:8*1024*1024])), nil)
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 8*1024*1024, internal.NopCloser(bytes.NewReader(fileData[8*1024*1024:])), nil)
//	_assert.Nil(err)
//	serviceSASValues := FileSASSignatureValues{ExpiryTime: time.Now().Add(time.Hour).UTC(),
//		Permissions: FileSASPermissions{Read: true, Write: true, Create: true}.String(), ShareName: shareName, FilePath: fileName}
//	credentials, _ := getGenericCredential(nil, testAccountDefault)
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credentials)
//	_assert.Nil(err)
//	srcFileWithSasURL := fClient.URL()
//	srcFileWithSasURL.RawQuery = queryParams.Encode()
//
//	fsu2, err := getGenericCredential(nil, testAccountSecondary)
//	_assert.Nil(err)
//	copysrClient, _ := createNewShare(_assert, fsu2)
//	copyfClient, _ := getFileClientFromShare(c, copysrClient)
//
//	defer delShare(c, copysrClient, DeleteSnapshotsOptionNone)
//
//	resp, err := copyfClient.StartCopy(ctx, srcFileWithSasURL, nil)
//	_assert.Nil(err)
//	_assert.(resp.CopyStatus, chk.Equals, CopyStatusPending)
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
//	_assert.(resp2.CopyStatus, chk.Equals, CopyStatusAborted)
//}

func (s *azfileLiveTestSuite) TestFileAbortCopyNoCopyStarted() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)

	defer delShare(_assert, srClient, nil)

	copyfClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)
	_, err = copyfClient.AbortCopy(ctx, "copynotstarted", nil)
	validateStorageError(_assert, err, StorageErrorCodeInvalidQueryParameterValue)
}

//func (s *azfileLiveTestSuite) TestResizeFile() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
//
//	fClient := createNewFileFromShare(_assert, srClient, 1234)
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
//	defer delShare(_assert, srClient, nil)
//	fClient := createNewFileFromShare(_assert, srClient, 10)
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
//	defer delShare(_assert, srClient, nil)
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.Resize(ctx, -4)
//	_assert.NotNil(err)
//	sErr := err.(StorageError)
//	_assert.(sErr.RawResponse.StatusCode, chk.Equals, http.StatusBadRequest)
//}
//
//func (f *azfileLivetestSuite) TestServiceSASShareSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
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
//	_, err = fClient.Create(ctx, int64(len(s)), FileHTTPHeaders{}, map[string]string)
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 0, bytes.NewReader([]byte(s)), nil)
//	_assert.Nil(err)
//	_, err = fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_, err = fClient.Delete(ctx)
//	_assert.Nil(err)
//
//	_, err = dirURL.Create(ctx, map[string]string, SMBProperties{})
//	_assert.Nil(err)
//
//	_, err = dirURL.ListFilesAndDirectoriesSegment(ctx, Marker{}, ListFilesAndDirectoriesOptions{})
//	_assert.Nil(err)
//}
//
//func (f *azfileLivetestSuite) TestServiceSASFileSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
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
//	_, err = fClient.Create(ctx, int64(len(s)), FileHTTPHeaders{}, map[string]string)
//	_assert.Nil(err)
//	_, err = fClient.UploadRange(ctx, 0, bytes.NewReader([]byte(s)), nil)
//	_assert.Nil(err)
//	dResp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_assert.(dResp.CacheControl, chk.Equals, cacheControlVal)
//	_assert.(dResp.ContentDisposition, chk.Equals, contentDispositionVal)
//	_assert.(dResp.ContentEncoding, chk.Equals, contentEncodingVal)
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
//	defer delShare(_assert, srClient, nil)
//
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//	defer delFile(_assert, fClient)
//
//	// Download entire fClient, check status code 200.
//	resp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_assert.Nil(err)
//	_assert.(resp.RawResponse.StatusCode, chk.Equals, http.StatusOK)
//	_assert.(resp.ContentLength, chk.Equals, int64(0))
//	_assert.(resp.FileContentMD5, chk.IsNil) // Note: FileContentMD5 is returned, only when range is specified explicitly.
//
//	download, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_assert.Nil(err)
//	_assert.(download, chk.HasLen, 0)
//	_assert.(resp.AcceptRanges, chk.Equals, "bytes")
//	_assert.(resp.CacheControl, chk.Equals, "")
//	_assert.(resp.ContentDisposition, chk.Equals, "")
//	_assert.(resp.ContentEncoding, chk.Equals, "")
//	_assert.(resp.ContentRange, chk.Equals, "") // Note: ContentRange is returned, only when range is specified explicitly.
//	_assert.(resp.ContentType, chk.Equals, "application/octet-stream")
//	_assert.(resp.CopyCompletionTime.IsZero(), chk.Equals, true)
//	_assert.(resp.CopyID, chk.Equals, "")
//	_assert.(resp.CopyProgress, chk.Equals, "")
//	_assert.(resp.CopySource, chk.Equals, "")
//	_assert.(resp.CopyStatus, chk.Equals, "")
//	_assert.(resp.CopyStatusDescription, chk.Equals, "")
//	_assert.(resp.Date.IsZero(), chk.Equals, false)
//	_assert.(resp.ETag, chk.Not(chk.Equals), "")
//	_assert.(resp.LastModified.IsZero(), chk.Equals, false)
//	_assert.(resp.Metadata, chk.DeepEquals, map[string]string)
//	_assert.(resp.RequestID, chk.Not(chk.Equals), "")
//	_assert.(resp.Version, chk.Not(chk.Equals), "")
//	_assert.(resp.IsServerEncrypted, chk.NotNil)
//}

//func (s *azfileLiveTestSuite) TestUploadDownloadDefaultNonDefaultMD5() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//	defer delShare(_assert, srClient, nil)
//
//	fClient := createNewFileFromShare(_assert, generateFileName(testName), 2048, srClient)
//	defer delFile(_assert, fClient)
//
//	contentR, contentD := generateData(2048)
//
//	pResp, err := fClient.UploadRange(ctx, 0, contentR, nil)
//	_assert.Nil(err)
//	_assert.NotNil(pResp.ContentMD5)
//	_assert.Equal(pResp.RawResponse.StatusCode, http.StatusCreated)
//	_assert.NotNil(pResp.IsServerEncrypted)
//	_assert.NotEqual(pResp.ETag, "")
//	_assert.Equal(pResp.LastModified.IsZero(), false)
//	_assert.NotEqual(pResp.RequestID,  "")
//	_assert.NotEqual(pResp.Version, "")
//	_assert.Equal(pResp.Date.IsZero(), false)
//
//	// Get with rangeGetContentMD5 enabled.
//	// Partial data, check status code 206.
//	resp, err := fClient.Download(ctx, 0, 1024, &DownloadFileOptions{RangeGetContentMD5: to.BoolPtr(true)})
//	_assert.Nil(err)
//	_assert.Equal(resp.RawResponse.StatusCode, http.StatusPartialContent)
//	_assert.Equal(*resp.ContentLength, int64(1024))
//	_assert.NotNil(resp.ContentMD5)
//	_assert.Equal(*resp.ContentType, "application/octet-stream")
//	_assert.NotEqual(resp.RawResponse.Status, "")
//
//	downloadedData, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_assert.Nil(err)
//	_assert.EqualValues(downloadedData, contentD[:1024])
//
//	// Set ContentMD5 for the entire file.
//	_, err = fClient.SetHTTPHeaders(ctx,
//		&SetFileHTTPHeadersOptions{
//		FileHTTPHeaders: &FileHTTPHeaders{
//			FileContentMD5: pResp.ContentMD5,
//			FileContentLanguage: to.StringPtr("test")}})
//	_assert.Nil(err)
//
//	// Test get with another type of range index, and validate if FileContentMD5 can be got correct.
//	resp, err = fClient.Download(ctx, 1024, CountToEnd, nil)
//	_assert.Nil(err)
//	_assert.Equal(resp.RawResponse.StatusCode, http.StatusPartialContent)
//	_assert.Equal(*resp.ContentLength, int64(1024))
//	_assert.Nil(resp.ContentMD5)
//	_assert.EqualValues(resp.FileContentMD5, pResp.ContentMD5)
//	_assert.Equal(*resp.ContentLanguage, "test")
//	// Note: when it's downloading range, range's MD5 is returned, when set rangeGetContentMD5=true, currently set it to false, so should be empty
//
//	downloadedData, err = ioutil.ReadAll(resp.RawResponse.Body)
//	_assert.Nil(err)
//	_assert.EqualValues(downloadedData, contentD[1024:])
//
//	_assert.Equal(*resp.AcceptRanges, "bytes")
//	_assert.Nil(resp.CacheControl)
//	_assert.Nil(resp.ContentDisposition)
//	_assert.Nil(resp.ContentEncoding)
//	_assert.Equal(*resp.ContentRange, "bytes 1024-2047/2048")
//	_assert.Nil(resp.ContentType) // Note ContentType is set to empty during SetHTTPHeaders
//	_assert.Equal(resp.CopyCompletionTime.IsZero(), true)
//	_assert.Nil(resp.CopyID)
//	_assert.Nil(resp.CopyProgress)
//	_assert.Nil(resp.CopySource)
//	_assert.Nil(resp.CopyStatus)
//	_assert.Nil(resp.CopyStatusDescription)
//	_assert.Equal(resp.Date.IsZero(), false)
//	_assert.NotEqual(*resp.ETag, "")
//	_assert.Equal(resp.LastModified.IsZero(), false)
//	_assert.EqualValues(resp.Metadata, map[string]string{})
//	_assert.NotEqual(*resp.RequestID, "")
//	_assert.NotEqual(*resp.Version, "")
//	_assert.NotNil(resp.IsServerEncrypted)
//
//	// Get entire fClient, check status code 200.
//	resp, err = fClient.Download(ctx, 0, CountToEnd, nil)
//	_assert.Nil(err)
//	_assert.Equal(resp.RawResponse.StatusCode, http.StatusOK)
//	_assert.Equal(*resp.ContentLength, int64(2048))
//	_assert.EqualValues(resp.ContentMD5, pResp.ContentMD5) // Note: This case is inted to get entire fClient, entire file's MD5 will be returned.
//	_assert.Nil(resp.FileContentMD5)                      // Note: FileContentMD5 is returned, only when range is specified explicitly.
//
//	downloadedData, err = ioutil.ReadAll(resp.RawResponse.Body)
//	_assert.Nil(err)
//	_assert.EqualValues(downloadedData, contentD[:])
//
//	_assert.Equal(*resp.AcceptRanges, "bytes")
//	_assert.Equal(*resp.CacheControl, "")
//	_assert.Equal(*resp.ContentDisposition, "")
//	_assert.Equal(*resp.ContentEncoding, "")
//	_assert.Equal(*resp.ContentRange, "") // Note: ContentRange is returned, only when range is specified explicitly.
//	_assert.Equal(*resp.ContentType, "")
//	_assert.Equal(resp.CopyCompletionTime.IsZero(), true)
//	_assert.Equal(*resp.CopyID, "")
//	_assert.Equal(*resp.CopyProgress, "")
//	_assert.Equal(*resp.CopySource, "")
//	_assert.Equal(*resp.CopyStatus, "")
//	_assert.Equal(*resp.CopyStatusDescription, "")
//	_assert.Equal(resp.Date.IsZero(), false)
//	_assert.NotEqual(*resp.ETag, "")
//	_assert.Equal(resp.LastModified.IsZero(), false)
//	_assert.EqualValues(resp.Metadata, map[string]string{})
//	_assert.NotEqual(*resp.RequestID, "")
//	_assert.NotEqual(*resp.Version, "")
//	_assert.NotNil(resp.IsServerEncrypted)
//}

func (s *azfileLiveTestSuite) TestFileDownloadDataNonExistentFile() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	_, err = fClient.Download(ctx, 0, CountToEnd, nil)
	validateStorageError(_assert, err, StorageErrorCodeResourceNotFound)
}

//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileDownloadDataNegativeOffset() {
//// 	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//// 	defer delShare(_assert, srClient, nil)
//// 	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.Download(ctx, -1, CountToEnd, false)
//// 	_assert.NotNil(err)
//// 	_assert.(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }

//---------------------------------------------------------------------------------------------------------------------

func (s *azfileLiveTestSuite) TestFileDownloadDataOffsetOutOfRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)

	_, err = fClient.Download(ctx, int64(len(fileDefaultData)), CountToEnd, nil)
	validateStorageError(_assert, err, StorageErrorCodeInvalidRange)
}

//// Don't check count by design.
//// func (s *azfileLiveTestSuite) TestFileDownloadDataInvalidCount() {
//// 	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//// 	defer delShare(_assert, srClient, nil)
//// 	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.Download(ctx, 0, -100, false)
//// 	_assert.NotNil(err)
//// 	_assert.(strings.Contains(err.Error(), "count must be >= 0"), chk.Equals, true)
//// }

func (s *azfileLiveTestSuite) TestFileDownloadDataEntireFile() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_assert, generateFileName(testName), fileDefaultData, srClient)

	resp, err := fClient.Download(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)

	// Specifying a count of 0 results in the value being ignored
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.EqualValues(string(data), fileDefaultData)
}

func (s *azfileLiveTestSuite) TestFileDownloadDataCountExact() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_assert, generateFileName(testName), fileDefaultData, srClient)

	resp, err := fClient.Download(ctx, 0, int64(len(fileDefaultData)), nil)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.EqualValues(string(data), fileDefaultData)
}

func (s *azfileLiveTestSuite) TestFileDownloadDataCountOutOfRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_assert, generateFileName(testName), fileDefaultData, srClient)

	resp, err := fClient.Download(ctx, 0, int64(len(fileDefaultData))*2, nil)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.EqualValues(string(data), fileDefaultData)
}

//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileUploadRangeNegativeInvalidOffset() {
//// 	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
//// 	defer delShare(_assert, srClient, nil)
//// 	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.UploadRange(ctx, -2, strings.NewReader(fileDefaultData), nil)
//// 	_assert.NotNil(err)
//// 	_assert.(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }

func (s *azfileLiveTestSuite) TestFileUploadRangeNilBody() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)

	_, err = fClient.UploadRange(ctx, 0, nil, nil)
	_assert.NotNil(err)
	_assert.Contains(err.Error(), "body must not be nil")
}

func (s *azfileLiveTestSuite) TestFileUploadRangeEmptyBody() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)

	_, err = fClient.UploadRange(ctx, 0, internal.NopCloser(bytes.NewReader([]byte{})), nil)
	_assert.NotNil(err)
	_assert.Contains(err.Error(), "body must contain readable data whose size is > 0")
}

func (s *azfileLiveTestSuite) TestFileUploadRangeNonExistentFile() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	rsc, _ := generateData(12)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	validateStorageError(_assert, err, StorageErrorCodeResourceNotFound)
}

func (s *azfileLiveTestSuite) TestFileUploadRangeTransactionalMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := createNewFileFromShare(_assert, generateFileName(testName), 2048, srClient)
	defer delFile(_assert, fClient)

	contentR, contentD := generateData(2048)
	md5 := md5.Sum(contentD)

	// Upload range with correct transactional MD5
	pResp, err := fClient.UploadRange(ctx, 0, contentR, &UploadFileRangeOptions{ContentMD5: md5[:]})
	_assert.Nil(err)
	_assert.NotNil(pResp.ContentMD5)
	_assert.Equal(pResp.RawResponse.StatusCode, http.StatusCreated)
	_assert.NotEqual(pResp.ETag, "")
	_assert.Equal(pResp.LastModified.IsZero(), false)
	_assert.NotEqual(pResp.RequestID, "")
	_assert.NotEqual(pResp.Version, "")
	_assert.Equal(pResp.Date.IsZero(), false)
	_assert.EqualValues(pResp.ContentMD5, md5[:])

	// Upload range with empty MD5, nil MD5 is covered by other cases.
	pResp, err = fClient.UploadRange(ctx, 1024, internal.NopCloser(bytes.NewReader(contentD[1024:])), nil)
	_assert.Nil(err)
	_assert.NotNil(pResp.ContentMD5)
	_assert.Equal(pResp.RawResponse.StatusCode, http.StatusCreated)

	resp, err := fClient.Download(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, http.StatusOK)
	_assert.Equal(*resp.ContentLength, int64(2048))

	downloadedData, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.EqualValues(downloadedData, contentD[:])
}

func (s *azfileLiveTestSuite) TestFileUploadRangeIncorrectTransactionalMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := createNewFileFromShare(_assert, generateFileName(testName), 2048, srClient)
	defer delFile(_assert, fClient)

	contentR, _ := generateData(2048)
	_, incorrectMD5 := generateData(16)

	// Upload range with incorrect transactional MD5
	_, err = fClient.UploadRange(ctx, 0, contentR, &UploadFileRangeOptions{ContentMD5: incorrectMD5[:]})
	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

//func (s *azfileLiveTestSuite) TestUploadRangeFromURL() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_assert, shareName, svcClient)
//	defer delShare(_assert, srClient, nil)
//
//	// create the source file and populate it with random data at a specific offset
//	expectedDataSize := 2048
//	totalFileSize := 4096
//	srcOffset := 999
//	expectedDataReader, expectedData := generateData(expectedDataSize)
//	srcfClient := createNewFileFromShare(_assert, "src" + generateFileName(testName), int64(totalFileSize), srClient)
//
//	_, err = srcfClient.UploadRange(ctx, int64(srcOffset), expectedDataReader, nil)
//	_assert.Nil(err)
//
//	// generate a URL with SAS pointing to the source file
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	_assert.Nil(err)
//
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:    SASProtocolHTTPS,
//		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour),
//		ShareName:   shareName,
//		Permissions: FileSASPermissions{Create: true, Read: true, Write: true, Delete: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_assert.Nil(err)
//	srcfClient.u.RawQuery = sasQueryParams.Encode()
//	rawSrcURL := srcfClient.URL()
//
//	// create the destination file
//	dstfClient := createNewFileFromShare(_assert, "dst" + generateFileName(testName), int64(totalFileSize), srClient)
//
//	// invoke UploadRange on dstfClient and put the data at a random range
//	// source and destination have different offsets, so we can test both values at the same time
//	dstOffset := 100
//	uploadFromURLResp, err := dstfClient.UploadRangeFromURL(ctx, rawSrcURL, int64(srcOffset), int64(dstOffset), int64(expectedDataSize), nil)
//	_assert.Nil(err)
//	_assert.Equal(uploadFromURLResp.RawResponse.StatusCode, 201)
//
//	// verify the destination
//	resp, err := dstfClient.Download(ctx, int64(dstOffset), int64(expectedDataSize), &DownloadFileOptions{RangeGetContentMD5: to.BoolPtr(false)})
//	_assert.Nil(err)
//	downloadedData, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_assert.Nil(err)
//	_assert.EqualValues(downloadedData, expectedData)
//}

// Testings for GetRangeList and ClearRange
func (s *azfileLiveTestSuite) TestGetRangeListNonDefaultExact() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	fileSize := int64(512 * 10)

	_, err = fClient.Create(ctx, &CreateFileOptions{FileContentLength: to.Int64Ptr(fileSize), FileHTTPHeaders: &FileHTTPHeaders{}})
	_assert.Nil(err)

	defer delFile(_assert, fClient)

	rsc, _ := generateData(1024)
	putResp, err := fClient.UploadRange(ctx, 0, rsc, nil)
	_assert.Nil(err)
	_assert.Equal(putResp.RawResponse.StatusCode, 201)
	_assert.Equal(putResp.LastModified.IsZero(), false)
	_assert.NotEqual(putResp.ETag, "")
	_assert.NotNil(putResp.ContentMD5)
	_assert.NotEqual(putResp.RequestID, "")
	_assert.NotEqual(putResp.Version, "")
	_assert.Equal(putResp.Date.IsZero(), false)

	rangeList, err := fClient.GetRangeList(ctx, 0, 1023, nil)
	_assert.Nil(err)
	_assert.Equal(rangeList.RawResponse.StatusCode, 200)
	_assert.Equal(rangeList.LastModified.IsZero(), false)
	_assert.NotEqual(rangeList.ETag, "")
	_assert.Equal(*rangeList.FileContentLength, fileSize)
	_assert.NotEqual(rangeList.RequestID, "")
	_assert.NotEqual(rangeList.Version, "")
	_assert.Equal(rangeList.Date.IsZero(), false)
	_assert.Len(rangeList.Ranges, 1)
	_assert.EqualValues(*rangeList.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(1023)})
}

// Default means clear the entire file's range
func (s *azfileLiveTestSuite) TestClearRangeDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := createNewFileFromShare(_assert, generateFileName(testName), 2048, srClient)
	defer delFile(_assert, fClient)

	rsc, _ := generateData(2048)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	_assert.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 0, 2048, nil)
	_assert.Nil(err)
	_assert.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	_assert.Len(rangeList.Ranges, 0)
}

func (s *azfileLiveTestSuite) TestClearRangeNonDefault() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := createNewFileFromShare(_assert, generateFileName(testName), 4096, srClient)
	defer delFile(_assert, fClient)

	rsc, _ := generateData(2048)
	_, err = fClient.UploadRange(ctx, 2048, rsc, nil)
	_assert.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 2048, 2048, nil)
	_assert.Nil(err)
	_assert.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	_assert.Len(rangeList.Ranges, 0)
}

func (s *azfileLiveTestSuite) TestClearRangeMultipleRanges() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := createNewFileFromShare(_assert, generateFileName(testName), 2048, srClient)
	defer delFile(_assert, fClient)

	rsc, _ := generateData(2048)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	_assert.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 1024, 1024, nil)
	_assert.Nil(err)
	_assert.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	_assert.Len(rangeList.Ranges, 1)
	_assert.EqualValues(*rangeList.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(1023)})
}

// When not 512 aligned, clear range will set 0 the non-512 aligned range, and will not eliminate the range.
func (s *azfileLiveTestSuite) TestClearRangeNonDefaultCount() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)

	fClient := createNewFileFromShare(_assert, generateFileName(testName), int64(1), srClient)
	defer delFile(_assert, fClient)

	d := []byte{1}
	_, err = fClient.UploadRange(ctx, 0, internal.NopCloser(bytes.NewReader(d)), nil)
	_assert.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 0, 1, nil)
	_assert.Nil(err)
	_assert.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	_assert.Len(rangeList.Ranges, 1)
	_assert.EqualValues(*rangeList.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(0)})

	dResp, err := fClient.Download(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)

	_bytes, err := ioutil.ReadAll(dResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(_bytes, []byte{0})
}

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

func (s *azfileLiveTestSuite) TestFileClearRangeNegativeInvalidCount() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	srClient, _ := getShareClient(generateShareName(testName), svcClient)
	fClient := getFileClientFromShare(_assert, generateFileName(testName), srClient)

	_, err = fClient.ClearRange(ctx, 0, 0, nil)
	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), "invalid argument: either offset is < 0 or count <= 0"), true)
}

func setupGetRangeListTest(_assert *assert.Assertions, testName string) (srClient ShareClient, fClient FileClient) {
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	_assert.Nil(err)
	srClient = createNewShare(_assert, generateShareName(testName), svcClient)
	fClient = createNewFileFromShare(_assert, generateFileName(testName), int64(testFileRangeSize), srClient)

	rsc, _ := generateData(testFileRangeSize)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	_assert.Nil(err)
	return
}

func validateBasicGetRangeList(_assert *assert.Assertions, resp FileGetRangeListResponse, err error) {
	_assert.Nil(err)
	_assert.Len(resp.Ranges, 1)
	_assert.EqualValues(*resp.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(int64(testFileRangeSize - 1))})
}

func (s *azfileLiveTestSuite) TestFileGetRangeListDefaultEmptyFile() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer delShare(_assert, srClient, nil)
	fClient := createNewFileFromShare(_assert, generateFileName(testName), 0, srClient)

	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	_assert.Len(resp.Ranges, 0)
}

func (s *azfileLiveTestSuite) TestFileGetRangeListDefaultRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_assert, testName)
	defer delShare(_assert, srClient, nil)

	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	validateBasicGetRangeList(_assert, resp, err)
}

func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRanges() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_assert, testName)
	defer delShare(_assert, srClient, nil)

	_, err := fClient.Resize(ctx, int64(testFileRangeSize*3), nil)
	_assert.Nil(err)

	rsc, _ := generateData(testFileRangeSize)
	_, err = fClient.UploadRange(ctx, int64(testFileRangeSize*2), rsc, nil)
	_assert.Nil(err)
	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	_assert.Len(resp.Ranges, 2)
	_assert.EqualValues(*resp.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(int64(testFileRangeSize - 1))})
	_assert.EqualValues(*resp.Ranges[1], FileRange{Start: to.Int64Ptr(int64(testFileRangeSize * 2)), End: to.Int64Ptr(int64((testFileRangeSize * 3) - 1))})
}

func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRangesCountLess() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_assert, testName)
	defer delShare(_assert, srClient, nil)

	resp, err := fClient.GetRangeList(ctx, 0, int64(testFileRangeSize-1), nil)
	_assert.Nil(err)
	_assert.Len(resp.Ranges, 1)
	_assert.EqualValues(*resp.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(int64(testFileRangeSize - 1))})
}

func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRangesCountExceed() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_assert, testName)
	defer delShare(_assert, srClient, nil)

	resp, err := fClient.GetRangeList(ctx, int64(0), int64(testFileRangeSize+1), nil)
	_assert.Nil(err)
	validateBasicGetRangeList(_assert, resp, err)
}

func (s *azfileLiveTestSuite) TestFileGetRangeListSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_assert, testName)
	delShareInclude := DeleteSnapshotsOptionTypeInclude
	defer delShare(_assert, srClient, &DeleteShareOptions{DeleteSnapshots: &delShareInclude})

	resp, _ := srClient.CreateSnapshot(ctx, &CreateShareSnapshotOptions{Metadata: map[string]string{}})
	_assert.NotNil(resp.Snapshot)
	fClientWithSnapshot := fClient.WithSnapshot(*resp.Snapshot)

	resp2, err := fClientWithSnapshot.GetRangeList(ctx, 0, CountToEnd, nil)
	_assert.Nil(err)
	validateBasicGetRangeList(_assert, resp2, err)
}

func (s *azfileLiveTestSuite) TestUnexpectedEOFRecovery() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_assert, generateShareName(testName), svcClient)

	defer delShare(_assert, srClient, &DeleteShareOptions{
		DeleteSnapshots: &deleteSnapshotsInclude,
	})

	fClient := createNewFileFromShare(_assert, generateFileName(testName), 2048, srClient)

	contentR, contentD := generateData(2048)

	resp, err := fClient.UploadRange(ctx, 0, contentR, nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, http.StatusCreated)
	_assert.NotEqual(resp.RequestID, "")

	dlResp, err := fClient.Download(ctx, 0, 2048, nil)
	_assert.Nil(err)

	// Verify that we can inject errors first.
	reader := dlResp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))

	_, err = ioutil.ReadAll(reader)
	_assert.NotNil(err)
	_assert.Equal(err.Error(), "unrecoverable error")

	// Then inject the retryable error.
	reader = dlResp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))

	buf, err := ioutil.ReadAll(reader)
	_assert.Nil(err)
	_assert.EqualValues(buf, contentD)
}

func (s *azfileLiveTestSuite) TestCreateMaximumSizeFileShare() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	srClient, err := getShareClient(generateShareName(testName), svcClient)
	cResp, err := srClient.Create(ctx, &CreateShareOptions{
		Quota: &fileShareMaxQuota,
	})
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)

	defer delShare(_assert, srClient, &DeleteShareOptions{DeleteSnapshots: &deleteSnapshotsInclude})
	dirClient, err := srClient.NewRootDirectoryClient()
	_assert.Nil(err)
	fClient := getFileClientFromDirectory(_assert, generateFileName(testName), dirClient)
	_, err = fClient.Create(ctx, &CreateFileOptions{
		FileContentLength: &fileMaxAllowedSizeInBytes,
		FileHTTPHeaders:   &FileHTTPHeaders{},
	})
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
}
