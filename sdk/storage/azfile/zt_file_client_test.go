//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"bytes"
	"crypto/md5"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
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
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareName := generateShareName(testName)
	srClient := createNewShare(_require, shareName, svcClient)
	defer delShare(_require, srClient, nil)

	// Create and delete fClient in root directory.
	fileName := generateFileName(testName)
	dirClient1, err := srClient.NewRootDirectoryClient()
	_require.Nil(err)
	fClient, err := dirClient1.NewFileClient(fileName)
	_require.Nil(err)

	cResp, err := fClient.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)
	_require.NotEqual(cResp.ETag, azcore.ETag(""))
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotEqual(cResp.RequestID, "")
	_require.NotEqual(cResp.Version, "")
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.IsServerEncrypted)

	delResp, err := fClient.Delete(ctx, nil)
	_require.Nil(err)
	_require.Equal(delResp.RawResponse.StatusCode, 202)
	_require.NotEqual(delResp.RequestID, "")
	_require.NotEqual(delResp.Version, "")
	_require.Equal(delResp.Date.IsZero(), false)

	dirName := generateDirectoryName(testName)
	dirClient := createNewDirectoryFromShare(_require, dirName, srClient)
	defer delDirectory(_require, dirClient)

	// Create and delete fClient in named directory.
	afClient, err := dirClient.NewFileClient(generateFileName(testName))
	_require.Nil(err)

	cResp, err = afClient.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)
	_require.NotEqual(cResp.ETag, "")
	_require.Equal(cResp.LastModified.IsZero(), false)
	_require.NotEqual(cResp.RequestID, "")
	_require.NotEqual(cResp.Version, "")
	_require.Equal(cResp.Date.IsZero(), false)
	_require.NotNil(cResp.IsServerEncrypted)

	delResp, err = afClient.Delete(ctx, nil)
	_require.Nil(err)
	_require.Equal(delResp.RawResponse.StatusCode, 202)
	_require.NotEqual(delResp.RequestID, "")
	_require.NotEqual(delResp.Version, "")
	_require.Equal(delResp.Date.IsZero(), false)
}

func (s *azfileLiveTestSuite) TestFileCreateNonDefaultMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)

	_, err = fClient.Create(ctx, &CreateFileOptions{
		Metadata: basicMetadata,
	})
	_require.Nil(err)

	resp, err := fClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_requireLen(resp.Metadata, len(basicMetadata))
}

func (s *azfileLiveTestSuite) TestFileCreateNonDefaultHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)

	_, err = fClient.Create(ctx, &CreateFileOptions{FileHTTPHeaders: &basicHeaders})
	_require.Nil(err)

	_, err = fClient.GetProperties(ctx, nil)
	_require.Nil(err)

	//h := resp.RawResponse.Header
	//_require.EqualValues(h, basicHeaders)
}

func (s *azfileLiveTestSuite) TestFileCreateNegativeMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)

	_, err = fClient.Create(ctx, &CreateFileOptions{
		Metadata:        map[string]string{"!@#$%^&*()": "!@#$%^&*()"},
		FileHTTPHeaders: &FileHTTPHeaders{},
	})
	_require.NotNil(err)
}

//func (s *azfileLiveTestSuite) TestFileGetSetPropertiesNonDefault() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//	defer delFile(_require, fClient)
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
//	_require.Nil(err)
//	lastWriteTime, err = time.Parse(ISO8601, lastWriteTime.Format(ISO8601))
//	_require.Nil(err)
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
//	_require.Nil(err)
//	_require.Equal(setResp.RawResponse.StatusCode, 200)
//	_require.NotEqual(setResp.ETag, "")
//	_require.Equal(setResp.LastModified.IsZero(),  false)
//	_require.NotEqual(setResp.RequestID, "")
//	_require.NotEqual(setResp.Version, "")
//	_require.Equal(setResp.Date.IsZero(),  false)
//	_require.NotNil(setResp.IsServerEncrypted)
//
//	getResp, err := fClient.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.Equal(getResp.RawResponse.StatusCode,  200)
//	_require.Equal(setResp.LastModified.IsZero(),  false)
//	_require.Equal(*getResp.FileType,  "File")
//
//	_require.EqualValues(getResp.ContentType,  options.FileHTTPHeaders.FileContentType)
//	_require.EqualValues(getResp.ContentEncoding,  options.FileHTTPHeaders.FileContentEncoding)
//	_require.EqualValues(getResp.ContentLanguage,  options.FileHTTPHeaders.FileContentLanguage)
//	_require.EqualValues(getResp.ContentMD5, options.FileHTTPHeaders.FileContentMD5)
//	_require.EqualValues(getResp.CacheControl,  options.FileHTTPHeaders.FileCacheControl)
//	_require.EqualValues(getResp.ContentDisposition,  options.FileHTTPHeaders.FileContentDisposition)
//	_require.Equal(*getResp.ContentLength,  int64(0))
//	// We'll just ensure a permission exists, no need to test overlapping functionality.
//	_require.NotEqual(*getResp.PermissionKey, "")
//	// Ensure our attributes and other properties (after parsing) are equivalent to our original
//	// There's an overlapping test for this in ntfs_property_bitflags_test.go, but it doesn't hurt to test it alongside other things.
//	_require.EqualValues(ParseFileAttributeFlagsString(*getResp.FileAttributes),  *options.SMBProperties.FileAttributes)
//
//	fct, _ := time.Parse(ISO8601, *getResp.FileCreationTime)
//	_require.EqualValues(fct, creationTime)
//	fwt, _ := time.Parse(ISO8601, *getResp.FileLastWriteTime)
//	_require.Equal(fwt, lastWriteTime)
//
//	_require.NotEqual(getResp.ETag, "")
//	_require.NotEqual(getResp.RequestID, "")
//	_require.NotEqual(getResp.Version, "")
//	_require.Equal(getResp.Date.IsZero(),  false)
//	_require.NotNil(getResp.IsServerEncrypted)
//}

//func (s *azfileLiveTestSuite) TestFilePreservePermissions() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//
//	fClient, _ := createNewFileFromShareWithPermissions(c, srClient, 0)
//	defer delFile(_require, fClient)
//
//	// Grab the original perm key before we set file headers.
//	getResp, err := fClient.GetProperties(ctx)
//	_require.Nil(err)
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
//	_require.Nil(err)
//	_require(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_require(setResp.ETag, chk.Not(chk.Equals), "")
//	_require(setResp.LastModified.IsZero(), chk.Equals, false)
//	_require(setResp.RequestID, chk.Not(chk.Equals), "")
//	_require(setResp.Version, chk.Not(chk.Equals), "")
//	_require(setResp.Date.IsZero(), chk.Equals, false)
//	_require(setResp.IsServerEncrypted, chk.NotNil)
//
//	getResp, err = fClient.GetProperties(ctx)
//	_require.Nil(err)
//	_require(getResp.RawResponse.StatusCode, chk.Equals, 200)
//	_require(setResp.LastModified.IsZero(), chk.Equals, false)
//	_require(getResp.FileType(), chk.Equals, "File")
//
//	_require(getResp.ContentType, chk.Equals, properties.ContentType)
//	_require(getResp.ContentEncoding, chk.Equals, properties.ContentEncoding)
//	_require(getResp.ContentLanguage, chk.Equals, properties.ContentLanguage)
//	_require(getResp.ContentMD5,, chk.DeepEquals, properties.ContentMD5)
//	_require(getResp.CacheControl, chk.Equals, properties.CacheControl)
//	_require(getResp.ContentDisposition, chk.Equals, properties.ContentDisposition)
//	_require(getResp.ContentLength, chk.Equals, int64(0))
//	// Ensure that the permission key gets preserved
//	_require(getResp.PermissionKey(), chk.Equals, oKey)
//	timeAdapter = SMBPropertyAdapter{PropertySource: getResp}
//	c.Log("Original last write time: ", lwTime, " new time: ", timeAdapter.FileLastWriteTime())
//	_require(timeAdapter.FileLastWriteTime().Equal(lwTime), chk.Equals, true)
//	c.Log("Original creation time: ", cTime, " new time: ", timeAdapter.FileCreationTime())
//	_require(timeAdapter.FileCreationTime().Equal(cTime), chk.Equals, true)
//	_require(getResp.FileAttributes(), chk.Equals, attribs)
//
//	_require(getResp.ETag, chk.Not(chk.Equals), "")
//	_require(getResp.RequestID, chk.Not(chk.Equals), "")
//	_require(getResp.Version, chk.Not(chk.Equals), "")
//	_require(getResp.Date.IsZero(), chk.Equals, false)
//	_require(getResp.IsServerEncrypted, chk.NotNil)
//}
//
//func (s *azfileLiveTestSuite) TestFileGetSetPropertiesSnapshot() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(c, srClient, DeleteSnapshotsOptionInclude)
//
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
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
//	_require.Nil(err)
//	_require(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_require(setResp.ETag, chk.Not(chk.Equals), "")
//	_require(setResp.LastModified.IsZero(), chk.Equals, false)
//	_require(setResp.RequestID, chk.Not(chk.Equals), "")
//	_require(setResp.Version, chk.Not(chk.Equals), "")
//	_require(setResp.Date.IsZero(), chk.Equals, false)
//	_require(setResp.IsServerEncrypted, chk.NotNil)
//
//	metadata := Metadata{
//		"foo": "foovalue",
//		"bar": "barvalue",
//	}
//	setResp2, err := fClient.SetMetadata(ctx, metadata)
//	_require.Nil(err)
//	_require(setResp2.RawResponse.StatusCode, chk.Equals, 200)
//
//	resp, _ := srClient.CreateSnapshot(ctx, map[string]string)
//	snapshotURL := fClient.WithSnapshot(resp.Snapshot())
//
//	getResp, err := snapshotURL.GetProperties(ctx)
//	_require.Nil(err)
//	_require(getResp.RawResponse.StatusCode, chk.Equals, 200)
//	_require(setResp.LastModified.IsZero(), chk.Equals, false)
//	_require(getResp.FileType(), chk.Equals, "File")
//
//	_require(getResp.ContentType, chk.Equals, properties.ContentType)
//	_require(getResp.ContentEncoding, chk.Equals, properties.ContentEncoding)
//	_require(getResp.ContentLanguage, chk.Equals, properties.ContentLanguage)
//	_require(getResp.ContentMD5,, chk.DeepEquals, properties.ContentMD5)
//	_require(getResp.CacheControl, chk.Equals, properties.CacheControl)
//	_require(getResp.ContentDisposition, chk.Equals, properties.ContentDisposition)
//	_require(getResp.ContentLength, chk.Equals, int64(0))
//
//	_require(getResp.ETag, chk.Not(chk.Equals), "")
//	_require(getResp.RequestID, chk.Not(chk.Equals), "")
//	_require(getResp.Version, chk.Not(chk.Equals), "")
//	_require(getResp.Date.IsZero(), chk.Equals, false)
//	_require(getResp.IsServerEncrypted, chk.NotNil)
//	_require(getResp.Metadata, chk.DeepEquals, metadata)
//}
//
//func (s *azfileLiveTestSuite) TestGetSetMetadataNonDefault() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//	metadata := Metadata{
//		"foo": "foovalue",
//		"bar": "barvalue",
//	}
//	setResp, err := fClient.SetMetadata(ctx, metadata)
//	_require.Nil(err)
//	_require(setResp.RawResponse.StatusCode, chk.Equals, 200)
//	_require(setResp.ETag, chk.Not(chk.Equals), "")
//	_require(setResp.RequestID, chk.Not(chk.Equals), "")
//	_require(setResp.Version, chk.Not(chk.Equals), "")
//	_require(setResp.Date.IsZero(), chk.Equals, false)
//	_require(setResp.IsServerEncrypted, chk.NotNil)
//
//	getResp, err := fClient.GetProperties(ctx)
//	_require.Nil(err)
//	_require(getResp.ETag, chk.Not(chk.Equals), "")
//	_require(getResp.LastModified.IsZero(), chk.Equals, false)
//	_require(getResp.RequestID, chk.Not(chk.Equals), "")
//	_require(getResp.Version, chk.Not(chk.Equals), "")
//	_require(getResp.Date.IsZero(), chk.Equals, false)
//	md := getResp.Metadata
//	_require(md, chk.DeepEquals, metadata)
//}
//
//func (s *azfileLiveTestSuite) TestFileSetMetadataNil() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"not": "nil"})
//	_require.Nil(err)
//
//	_, err = fClient.SetMetadata(ctx, nil)
//	_require.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_require.Nil(err)
//	_require(resp.Metadata, chk.HasLen, 0)
//}
//
//func (s *azfileLiveTestSuite) TestFileSetMetadataDefaultEmpty() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"not": "nil"})
//	_require.Nil(err)
//
//	_, err = fClient.SetMetadata(ctx, map[string]string)
//	_require.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_require.Nil(err)
//	_require(resp.Metadata, chk.HasLen, 0)
//}

//func (s *azfileLiveTestSuite) TestFileSetMetadataInvalidField() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.SetMetadata(ctx, Metadata{"!@#$%^&*()": "!@#$%^&*()"})
//	_require.NotNil(err)
//}
//
//func (s *azfileLiveTestSuite) TestStartCopyDefault() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//
//	srcFile, _ := createNewFileFromShare(_require, generateFileName(testName), 2048, srClient)
//	defer delFile(c, srcFile)
//
//	destFile, _ := getFileClientFromShare(_require, generateFileName(testName), srClient)
//	defer delFile(c, destFile)
//
//	_, err := srcFile.UploadRange(ctx, 0, generateData(2048), nil)
//	_require.Nil(err)
//
//	copyResp, err := destFile.StartCopy(ctx, srcFile.URL(), nil)
//	_require.Nil(err)
//	_require(copyResp.RawResponse.StatusCode, chk.Equals, 202)
//	_require(copyResp.ETag, chk.Not(chk.Equals), "")
//	_require(copyResp.LastModified.IsZero(), chk.Equals, false)
//	_require(copyResp.RequestID, chk.Not(chk.Equals), "")
//	_require(copyResp.Version, chk.Not(chk.Equals), "")
//	_require(copyResp.Date.IsZero(), chk.Equals, false)
//	_require(copyResp.CopyID, chk.Not(chk.Equals), "")
//	_require(copyResp.CopyStatus, chk.Not(chk.Equals), "")
//
//	var copyStatus CopyStatusType
//	timeout := time.Duration(2) * time.Minute
//	start := time.Now()
//
//	var getResp *FileGetPropertiesResponse
//
//	for copyStatus != CopyStatusSuccess && time.Now().Sub(start) < timeout {
//		getResp, err = destFile.GetProperties(ctx)
//		_require.Nil(err)
//		_require(getResp.CopyID, chk.Equals, copyResp.CopyID())
//		_require(getResp.CopyStatus, chk.Not(chk.Equals), "")
//		_require(getResp.CopySource, chk.Equals, srcFile.String())
//		copyStatus = getResp.CopyStatus
//
//		time.Sleep(time.Duration(5) * time.Second)
//	}
//
//	if getResp != nil && getResp.CopyStatus == CopyStatusSuccess {
//		// Abort will fail after copy finished
//		abortResp, err := destFile.AbortCopy(ctx, copyResp.CopyID())
//		_require.NotNil(err)
//		_require(abortResp, chk.IsNil)
//		se, ok := err.(StorageError)
//		_require(ok, chk.Equals, true)
//		_require(se.RawResponse.StatusCode, chk.Equals, http.StatusConflict)
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
			_requireFail("")
		}
	}
}

func (s *azfileLiveTestSuite) TestFileStartCopyDestEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_require, "src"+generateFileName(testName), fileDefaultData, srClient)
	copyfClient := getFileClientFromShare(_require, "dst"+generateFileName(testName), srClient)

	fileCopyResponse, err := copyfClient.StartCopy(ctx, fClient.URL(), nil)
	_require.Nil(err)
	waitForCopy(_require, copyfClient, fileCopyResponse)

	resp, err := copyfClient.Download(ctx, 0, CountToEnd, nil)
	_require.Nil(err)

	// Read the file data to verify the copy
	data, _ := ioutil.ReadAll(resp.RawResponse.Body)
	_require.Equal(*resp.ContentLength, int64(len(fileDefaultData)))
	_require.Equal(string(data), fileDefaultData)
	resp.RawResponse.Body.Close()
}

//func (s *azfileLiveTestSuite) TestFileStartCopyMetadata() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient := createNewFileFromShare(_require, "src" + generateFileName(testName), 0, srClient)
//	copyfClient := getFileClientFromShare(_require, "dst" + generateFileName(testName), srClient)
//
//	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), &StartFileCopyOptions{Metadata: basicMetadata})
//	_require.Nil(err)
//	waitForCopy(_require, copyfClient, resp)
//
//	resp2, err := copyfClient.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.EqualValues(resp2.Metadata, basicMetadata)
//}

func (s *azfileLiveTestSuite) TestFileStartCopyMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShare(_require, "src"+generateFileName(testName), 0, srClient)
	copyfClient := getFileClientFromShare(_require, "dst"+generateFileName(testName), srClient)

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err = copyfClient.Create(ctx, &CreateFileOptions{Metadata: basicMetadata})
	_require.Nil(err)

	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), nil)
	_require.Nil(err)

	waitForCopy(_require, copyfClient, resp)

	resp2, err := copyfClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_requireLen(resp2.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestFileStartCopyMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShare(_require, "src"+generateFileName(testName), 0, srClient)
	copyfClient := getFileClientFromShare(_require, "dst"+generateFileName(testName), srClient)

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err = copyfClient.Create(ctx, &CreateFileOptions{Metadata: basicMetadata})
	_require.Nil(err)

	resp, err := copyfClient.StartCopy(ctx, fClient.URL(), &StartFileCopyOptions{Metadata: map[string]string{}})
	_require.Nil(err)

	waitForCopy(_require, copyfClient, resp)

	resp2, err := copyfClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_requireLen(resp2.Metadata, 0)
}

func (s *azfileLiveTestSuite) TestFileStartCopyNegativeMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShare(_require, "src"+generateFileName(testName), 0, srClient)
	copyfClient := getFileClientFromShare(_require, "dst"+generateFileName(testName), srClient)

	_, err = copyfClient.StartCopy(ctx, fClient.URL(), &StartFileCopyOptions{Metadata: map[string]string{"!@#$%^&*()": "!@#$%^&*()"}})
	_require.NotNil(err)
}

func (s *azfileLiveTestSuite) TestFileStartCopySourceNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := getFileClientFromShare(_require, "src"+generateFileName(testName), srClient)
	copyfClient := getFileClientFromShare(_require, "dst"+generateFileName(testName), srClient)

	_, err = copyfClient.StartCopy(ctx, fClient.URL(), nil)
	validateStorageError(_require, err, StorageErrorCodeResourceNotFound)
}

//func (s *azfileLiveTestSuite) TestFileStartCopyUsingSASSrc() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient, fileName := createNewFileFromShareWithGivenData(c, srClient)
//
//	// Create sas values for the source file
//	credential, _ := getCredential()
//	serviceSASValues := FileSASSignatureValues{Version: "2015-04-05", StartTime: time.Now().Add(-1 * time.Hour).UTC(),
//		ExpiryTime: time.Now().Add(time.Hour).UTC(), Permissions: FileSASPermissions{Read: true, Write: true, Create: true, Delete: true}.String(),
//		ShareName: shareName, FilePath: fileName}
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
//	_require.Nil(err)
//
//	// Create URLs to the destination file with sas parameters
//	sasURL := fClient.URL()
//	sasURL.RawQuery = queryParams.Encode()
//
//	// Create a new container for the destination
//	copysrClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(c, copysrClient, DeleteSnapshotsOptionNone)
//	copyfClient, _ := getFileClientFromShare(c, copysrClient)
//
//	resp, err := copyfClient.StartCopy(ctx, sasURL, nil)
//	_require.Nil(err)
//
//	waitForCopy(c, copyfClient, resp)
//
//	resp2, err := copyfClient.Download(ctx, 0, CountToEnd, false)
//	_require.Nil(err)
//
//	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
//	_require(resp2.ContentLength, chk.Equals, int64(len(fileDefaultData)))
//	_require(string(data), chk.Equals, fileDefaultData)
//	resp2.RawResponse.Body.Close()
//}
//
//func (s *azfileLiveTestSuite) TestFileStartCopyUsingSASDest() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient, fileName := createNewFileFromShareWithGivenData(c, srClient)
//	_ = fClient
//
//	// Generate SAS on the source
//	serviceSASValues := FileSASSignatureValues{ExpiryTime: time.Now().Add(time.Hour).UTC(),
//		Permissions: FileSASPermissions{Read: true, Write: true, Create: true}.String(), ShareName: shareName, FilePath: fileName}
//	credentials, _ := getCredential()
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credentials)
//	_require.Nil(err)
//
//	copysrClient, copyShareName := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(c, copysrClient, DeleteSnapshotsOptionNone)
//	copyfClient, copyFileName := getFileClientFromShare(c, copysrClient)
//
//	// Generate Sas for the destination
//	copyServiceSASvalues := FileSASSignatureValues{StartTime: time.Now().Add(-1 * time.Hour).UTC(),
//		ExpiryTime: time.Now().Add(time.Hour).UTC(), Permissions: FileSASPermissions{Read: true, Write: true}.String(),
//		ShareName: copyShareName, FilePath: copyFileName}
//	copyQueryParams, err := copyServiceSASvalues.NewSASQueryParameters(credentials)
//	_require.Nil(err)
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
//	_require.Nil(err)
//
//	// Allow copy to happen
//	waitForCopy(c, anonfClient, resp)
//
//	resp2, err := copyfClient.Download(ctx, 0, CountToEnd, false)
//	_require.Nil(err)
//
//	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
//	_, err = resp2.Body(RetryReaderOptions{}).Read(data)
//	_require(resp2.ContentLength, chk.Equals, int64(len(fileDefaultData)))
//	_require(string(data), chk.Equals, fileDefaultData)
//	resp2.Body(RetryReaderOptions{}).Close()
//}

//func (s *azfileLiveTestSuite) TestFileAbortCopyInProgress() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_require, shareName, svcClient)
//	defer delShare(_require, srClient, nil)
//	fileName := generateFileName(testName)
//	fClient := getFileClientFromShare(_require, fileName, srClient)
//
//	// Create a large file that takes time to copy
//	fileSize := 12 * 1024 * 1024
//	fileData := make([]byte, fileSize, fileSize)
//	for i := range fileData {
//		fileData[i] = byte('a' + i%26)
//	}
//	_, err = fClient.Create(ctx, &CreateFileOptions{FileContentLength: to.Int64Ptr(int64(fileSize)), FileHTTPHeaders: &FileHTTPHeaders{}})
//	_require.Nil(err)
//
//	_, err = fClient.UploadRange(ctx, 0, internal.NopCloser(bytes.NewReader(fileData[0:4*1024*1024])), nil)
//	_require.Nil(err)
//	_, err = fClient.UploadRange(ctx, 4*1024*1024, internal.NopCloser(bytes.NewReader(fileData[4*1024*1024:8*1024*1024])), nil)
//	_require.Nil(err)
//	_, err = fClient.UploadRange(ctx, 8*1024*1024, internal.NopCloser(bytes.NewReader(fileData[8*1024*1024:])), nil)
//	_require.Nil(err)
//	serviceSASValues := FileSASSignatureValues{ExpiryTime: time.Now().Add(time.Hour).UTC(),
//		Permissions: FileSASPermissions{Read: true, Write: true, Create: true}.String(), ShareName: shareName, FilePath: fileName}
//	credentials, _ := getGenericCredential(nil, testAccountDefault)
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credentials)
//	_require.Nil(err)
//	srcFileWithSasURL := fClient.URL()
//	srcFileWithSasURL.RawQuery = queryParams.Encode()
//
//	fsu2, err := getGenericCredential(nil, testAccountSecondary)
//	_require.Nil(err)
//	copysrClient, _ := createNewShare(_require, fsu2)
//	copyfClient, _ := getFileClientFromShare(c, copysrClient)
//
//	defer delShare(c, copysrClient, DeleteSnapshotsOptionNone)
//
//	resp, err := copyfClient.StartCopy(ctx, srcFileWithSasURL, nil)
//	_require.Nil(err)
//	_require(resp.CopyStatus, chk.Equals, CopyStatusPending)
//
//	_, err = copyfClient.AbortCopy(ctx, resp.CopyID())
//	if err != nil {
//		// If the error is nil, the test continues as normal.
//		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
//		_require((err.(StorageError)).RawResponse.StatusCode, chk.Equals, 409)
//		c.Error("The test failed because the copy completed because it was aborted")
//	}
//
//	resp2, _ := copyfClient.GetProperties(ctx)
//	_require(resp2.CopyStatus, chk.Equals, CopyStatusAborted)
//}

func (s *azfileLiveTestSuite) TestFileAbortCopyNoCopyStarted() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)

	defer delShare(_require, srClient, nil)

	copyfClient := getFileClientFromShare(_require, generateFileName(testName), srClient)
	_, err = copyfClient.AbortCopy(ctx, "copynotstarted", nil)
	validateStorageError(_require, err, StorageErrorCodeInvalidQueryParameterValue)
}

//func (s *azfileLiveTestSuite) TestResizeFile() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//
//	fClient := createNewFileFromShare(_require, srClient, 1234)
//
//	gResp, err := fClient.GetProperties(ctx)
//	_require.Nil(err)
//	_require(gResp.ContentLength, chk.Equals, int64(1234))
//
//	rResp, err := fClient.Resize(ctx, 4096)
//	_require.Nil(err)
//	_require(rResp.RawResponse.StatusCode, chk.Equals, 200)
//
//	gResp, err = fClient.GetProperties(ctx)
//	_require.Nil(err)
//	_require(gResp.ContentLength, chk.Equals, int64(4096))
//}
//
//func (s *azfileLiveTestSuite) TestFileResizeZero() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient := createNewFileFromShare(_require, srClient, 10)
//
//	// The default file is created with size > 0, so this should actually update
//	_, err := fClient.Resize(ctx, 0)
//	_require.Nil(err)
//
//	resp, err := fClient.GetProperties(ctx)
//	_require.Nil(err)
//	_require(resp.ContentLength, chk.Equals, int64(0))
//}
//
//func (s *azfileLiveTestSuite) TestFileResizeInvalidSizeNegative() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//	_, err := fClient.Resize(ctx, -4)
//	_require.NotNil(err)
//	sErr := err.(StorageError)
//	_require(sErr.RawResponse.StatusCode, chk.Equals, http.StatusBadRequest)
//}
//
//func (f *azfileLivetestSuite) TestServiceSASShareSAS() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//
//	credential, accountName := getCredential()
//
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:    SASProtocolHTTPS,
//		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour),
//		ShareName:   shareName,
//		Permissions: ShareSASPermissions{Create: true, Read: true, Write: true, Delete: true, List: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_require.Nil(err)
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
//	_require.Nil(err)
//	_, err = fClient.UploadRange(ctx, 0, bytes.NewReader([]byte(s)), nil)
//	_require.Nil(err)
//	_, err = fClient.Download(ctx, 0, CountToEnd, false)
//	_require.Nil(err)
//	_, err = fClient.Delete(ctx)
//	_require.Nil(err)
//
//	_, err = dirURL.Create(ctx, map[string]string, SMBProperties{})
//	_require.Nil(err)
//
//	_, err = dirURL.ListFilesAndDirectoriesSegment(ctx, Marker{}, ListFilesAndDirectoriesOptions{})
//	_require.Nil(err)
//}
//
//func (f *azfileLivetestSuite) TestServiceSASFileSAS() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient, shareName := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
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
//	_require.Nil(err)
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
//	_require.Nil(err)
//	_, err = fClient.UploadRange(ctx, 0, bytes.NewReader([]byte(s)), nil)
//	_require.Nil(err)
//	dResp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_require.Nil(err)
//	_require(dResp.CacheControl, chk.Equals, cacheControlVal)
//	_require(dResp.ContentDisposition, chk.Equals, contentDispositionVal)
//	_require(dResp.ContentEncoding, chk.Equals, contentEncodingVal)
//	_require(dResp.ContentLanguage, chk.Equals, contentLanguageVal)
//	_require(dResp.ContentType, chk.Equals, contentTypeVal)
//	_, err = fClient.Delete(ctx)
//	_require.Nil(err)
//}
//
//func (s *azfileLiveTestSuite) TestDownloadEmptyZeroSizeFile() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//	defer delFile(_require, fClient)
//
//	// Download entire fClient, check status code 200.
//	resp, err := fClient.Download(ctx, 0, CountToEnd, false)
//	_require.Nil(err)
//	_require(resp.RawResponse.StatusCode, chk.Equals, http.StatusOK)
//	_require(resp.ContentLength, chk.Equals, int64(0))
//	_require(resp.FileContentMD5, chk.IsNil) // Note: FileContentMD5 is returned, only when range is specified explicitly.
//
//	download, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_require.Nil(err)
//	_require(download, chk.HasLen, 0)
//	_require(resp.AcceptRanges, chk.Equals, "bytes")
//	_require(resp.CacheControl, chk.Equals, "")
//	_require(resp.ContentDisposition, chk.Equals, "")
//	_require(resp.ContentEncoding, chk.Equals, "")
//	_require(resp.ContentRange, chk.Equals, "") // Note: ContentRange is returned, only when range is specified explicitly.
//	_require(resp.ContentType, chk.Equals, "application/octet-stream")
//	_require(resp.CopyCompletionTime.IsZero(), chk.Equals, true)
//	_require(resp.CopyID, chk.Equals, "")
//	_require(resp.CopyProgress, chk.Equals, "")
//	_require(resp.CopySource, chk.Equals, "")
//	_require(resp.CopyStatus, chk.Equals, "")
//	_require(resp.CopyStatusDescription, chk.Equals, "")
//	_require(resp.Date.IsZero(), chk.Equals, false)
//	_require(resp.ETag, chk.Not(chk.Equals), "")
//	_require(resp.LastModified.IsZero(), chk.Equals, false)
//	_require(resp.Metadata, chk.DeepEquals, map[string]string)
//	_require(resp.RequestID, chk.Not(chk.Equals), "")
//	_require(resp.Version, chk.Not(chk.Equals), "")
//	_require(resp.IsServerEncrypted, chk.NotNil)
//}

//func (s *azfileLiveTestSuite) TestUploadDownloadDefaultNonDefaultMD5() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//	defer delShare(_require, srClient, nil)
//
//	fClient := createNewFileFromShare(_require, generateFileName(testName), 2048, srClient)
//	defer delFile(_require, fClient)
//
//	contentR, contentD := generateData(2048)
//
//	pResp, err := fClient.UploadRange(ctx, 0, contentR, nil)
//	_require.Nil(err)
//	_require.NotNil(pResp.ContentMD5)
//	_require.Equal(pResp.RawResponse.StatusCode, http.StatusCreated)
//	_require.NotNil(pResp.IsServerEncrypted)
//	_require.NotEqual(pResp.ETag, "")
//	_require.Equal(pResp.LastModified.IsZero(), false)
//	_require.NotEqual(pResp.RequestID,  "")
//	_require.NotEqual(pResp.Version, "")
//	_require.Equal(pResp.Date.IsZero(), false)
//
//	// Get with rangeGetContentMD5 enabled.
//	// Partial data, check status code 206.
//	resp, err := fClient.Download(ctx, 0, 1024, &DownloadFileOptions{RangeGetContentMD5: to.BoolPtr(true)})
//	_require.Nil(err)
//	_require.Equal(resp.RawResponse.StatusCode, http.StatusPartialContent)
//	_require.Equal(*resp.ContentLength, int64(1024))
//	_require.NotNil(resp.ContentMD5)
//	_require.Equal(*resp.ContentType, "application/octet-stream")
//	_require.NotEqual(resp.RawResponse.Status, "")
//
//	downloadedData, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_require.Nil(err)
//	_require.EqualValues(downloadedData, contentD[:1024])
//
//	// Set ContentMD5 for the entire file.
//	_, err = fClient.SetHTTPHeaders(ctx,
//		&SetFileHTTPHeadersOptions{
//		FileHTTPHeaders: &FileHTTPHeaders{
//			FileContentMD5: pResp.ContentMD5,
//			FileContentLanguage: to.StringPtr("test")}})
//	_require.Nil(err)
//
//	// Test get with another type of range index, and validate if FileContentMD5 can be got correct.
//	resp, err = fClient.Download(ctx, 1024, CountToEnd, nil)
//	_require.Nil(err)
//	_require.Equal(resp.RawResponse.StatusCode, http.StatusPartialContent)
//	_require.Equal(*resp.ContentLength, int64(1024))
//	_require.Nil(resp.ContentMD5)
//	_require.EqualValues(resp.FileContentMD5, pResp.ContentMD5)
//	_require.Equal(*resp.ContentLanguage, "test")
//	// Note: when it's downloading range, range's MD5 is returned, when set rangeGetContentMD5=true, currently set it to false, so should be empty
//
//	downloadedData, err = ioutil.ReadAll(resp.RawResponse.Body)
//	_require.Nil(err)
//	_require.EqualValues(downloadedData, contentD[1024:])
//
//	_require.Equal(*resp.AcceptRanges, "bytes")
//	_require.Nil(resp.CacheControl)
//	_require.Nil(resp.ContentDisposition)
//	_require.Nil(resp.ContentEncoding)
//	_require.Equal(*resp.ContentRange, "bytes 1024-2047/2048")
//	_require.Nil(resp.ContentType) // Note ContentType is set to empty during SetHTTPHeaders
//	_require.Equal(resp.CopyCompletionTime.IsZero(), true)
//	_require.Nil(resp.CopyID)
//	_require.Nil(resp.CopyProgress)
//	_require.Nil(resp.CopySource)
//	_require.Nil(resp.CopyStatus)
//	_require.Nil(resp.CopyStatusDescription)
//	_require.Equal(resp.Date.IsZero(), false)
//	_require.NotEqual(*resp.ETag, "")
//	_require.Equal(resp.LastModified.IsZero(), false)
//	_require.EqualValues(resp.Metadata, map[string]string{})
//	_require.NotEqual(*resp.RequestID, "")
//	_require.NotEqual(*resp.Version, "")
//	_require.NotNil(resp.IsServerEncrypted)
//
//	// Get entire fClient, check status code 200.
//	resp, err = fClient.Download(ctx, 0, CountToEnd, nil)
//	_require.Nil(err)
//	_require.Equal(resp.RawResponse.StatusCode, http.StatusOK)
//	_require.Equal(*resp.ContentLength, int64(2048))
//	_require.EqualValues(resp.ContentMD5, pResp.ContentMD5) // Note: This case is inted to get entire fClient, entire file's MD5 will be returned.
//	_require.Nil(resp.FileContentMD5)                      // Note: FileContentMD5 is returned, only when range is specified explicitly.
//
//	downloadedData, err = ioutil.ReadAll(resp.RawResponse.Body)
//	_require.Nil(err)
//	_require.EqualValues(downloadedData, contentD[:])
//
//	_require.Equal(*resp.AcceptRanges, "bytes")
//	_require.Equal(*resp.CacheControl, "")
//	_require.Equal(*resp.ContentDisposition, "")
//	_require.Equal(*resp.ContentEncoding, "")
//	_require.Equal(*resp.ContentRange, "") // Note: ContentRange is returned, only when range is specified explicitly.
//	_require.Equal(*resp.ContentType, "")
//	_require.Equal(resp.CopyCompletionTime.IsZero(), true)
//	_require.Equal(*resp.CopyID, "")
//	_require.Equal(*resp.CopyProgress, "")
//	_require.Equal(*resp.CopySource, "")
//	_require.Equal(*resp.CopyStatus, "")
//	_require.Equal(*resp.CopyStatusDescription, "")
//	_require.Equal(resp.Date.IsZero(), false)
//	_require.NotEqual(*resp.ETag, "")
//	_require.Equal(resp.LastModified.IsZero(), false)
//	_require.EqualValues(resp.Metadata, map[string]string{})
//	_require.NotEqual(*resp.RequestID, "")
//	_require.NotEqual(*resp.Version, "")
//	_require.NotNil(resp.IsServerEncrypted)
//}

func (s *azfileLiveTestSuite) TestFileDownloadDataNonExistentFile() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)

	_, err = fClient.Download(ctx, 0, CountToEnd, nil)
	validateStorageError(_require, err, StorageErrorCodeResourceNotFound)
}

//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileDownloadDataNegativeOffset() {
//// 	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//// 	defer delShare(_require, srClient, nil)
//// 	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.Download(ctx, -1, CountToEnd, false)
//// 	_require.NotNil(err)
//// 	_require(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }

//---------------------------------------------------------------------------------------------------------------------

func (s *azfileLiveTestSuite) TestFileDownloadDataOffsetOutOfRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)

	_, err = fClient.Download(ctx, int64(len(fileDefaultData)), CountToEnd, nil)
	validateStorageError(_require, err, StorageErrorCodeInvalidRange)
}

//// Don't check count by design.
//// func (s *azfileLiveTestSuite) TestFileDownloadDataInvalidCount() {
//// 	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//// 	defer delShare(_require, srClient, nil)
//// 	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.Download(ctx, 0, -100, false)
//// 	_require.NotNil(err)
//// 	_require(strings.Contains(err.Error(), "count must be >= 0"), chk.Equals, true)
//// }

func (s *azfileLiveTestSuite) TestFileDownloadDataEntireFile() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_require, generateFileName(testName), fileDefaultData, srClient)

	resp, err := fClient.Download(ctx, 0, CountToEnd, nil)
	_require.Nil(err)

	// Specifying a count of 0 results in the value being ignored
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.EqualValues(string(data), fileDefaultData)
}

func (s *azfileLiveTestSuite) TestFileDownloadDataCountExact() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_require, generateFileName(testName), fileDefaultData, srClient)

	resp, err := fClient.Download(ctx, 0, int64(len(fileDefaultData)), nil)
	_require.Nil(err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.EqualValues(string(data), fileDefaultData)
}

func (s *azfileLiveTestSuite) TestFileDownloadDataCountOutOfRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShareWithGivenData(_require, generateFileName(testName), fileDefaultData, srClient)

	resp, err := fClient.Download(ctx, 0, int64(len(fileDefaultData))*2, nil)
	_require.Nil(err)

	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.EqualValues(string(data), fileDefaultData)
}

//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileUploadRangeNegativeInvalidOffset() {
//// 	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient := createNewShare(_require, generateShareName(testName), svcClient)
//// 	defer delShare(_require, srClient, nil)
//// 	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)
//
//// 	_, err := fClient.UploadRange(ctx, -2, strings.NewReader(fileDefaultData), nil)
//// 	_require.NotNil(err)
//// 	_require(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }

func (s *azfileLiveTestSuite) TestFileUploadRangeNilBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)

	_, err = fClient.UploadRange(ctx, 0, nil, nil)
	_require.NotNil(err)
	_requireContains(err.Error(), "body must not be nil")
}

func (s *azfileLiveTestSuite) TestFileUploadRangeEmptyBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)

	_, err = fClient.UploadRange(ctx, 0, NopCloser(bytes.NewReader([]byte{})), nil)
	_require.NotNil(err)
	_requireContains(err.Error(), "body must contain readable data whose size is > 0")
}

func (s *azfileLiveTestSuite) TestFileUploadRangeNonExistentFile() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)

	rsc, _ := generateData(12)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	validateStorageError(_require, err, StorageErrorCodeResourceNotFound)
}

func (s *azfileLiveTestSuite) TestFileUploadRangeTransactionalMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := createNewFileFromShare(_require, generateFileName(testName), 2048, srClient)
	defer delFile(_require, fClient)

	contentR, contentD := generateData(2048)
	md5 := md5.Sum(contentD)

	// Upload range with correct transactional MD5
	pResp, err := fClient.UploadRange(ctx, 0, contentR, &UploadFileRangeOptions{ContentMD5: md5[:]})
	_require.Nil(err)
	_require.NotNil(pResp.ContentMD5)
	_require.Equal(pResp.RawResponse.StatusCode, http.StatusCreated)
	_require.NotEqual(pResp.ETag, "")
	_require.Equal(pResp.LastModified.IsZero(), false)
	_require.NotEqual(pResp.RequestID, "")
	_require.NotEqual(pResp.Version, "")
	_require.Equal(pResp.Date.IsZero(), false)
	_require.EqualValues(pResp.ContentMD5, md5[:])

	// Upload range with empty MD5, nil MD5 is covered by other cases.
	pResp, err = fClient.UploadRange(ctx, 1024, NopCloser(bytes.NewReader(contentD[1024:])), nil)
	_require.Nil(err)
	_require.NotNil(pResp.ContentMD5)
	_require.Equal(pResp.RawResponse.StatusCode, http.StatusCreated)

	resp, err := fClient.Download(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, http.StatusOK)
	_require.Equal(*resp.ContentLength, int64(2048))

	downloadedData, err := ioutil.ReadAll(resp.RawResponse.Body)
	_require.Nil(err)
	_require.EqualValues(downloadedData, contentD[:])
}

func (s *azfileLiveTestSuite) TestFileUploadRangeIncorrectTransactionalMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := createNewFileFromShare(_require, generateFileName(testName), 2048, srClient)
	defer delFile(_require, fClient)

	contentR, _ := generateData(2048)
	_, incorrectMD5 := generateData(16)

	// Upload range with incorrect transactional MD5
	_, err = fClient.UploadRange(ctx, 0, contentR, &UploadFileRangeOptions{ContentMD5: incorrectMD5[:]})
	validateStorageError(_require, err, StorageErrorCodeMD5Mismatch)
}

//func (s *azfileLiveTestSuite) TestUploadRangeFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	shareName := generateShareName(testName)
//	srClient := createNewShare(_require, shareName, svcClient)
//	defer delShare(_require, srClient, nil)
//
//	// create the source file and populate it with random data at a specific offset
//	expectedDataSize := 2048
//	totalFileSize := 4096
//	srcOffset := 999
//	expectedDataReader, expectedData := generateData(expectedDataSize)
//	srcfClient := createNewFileFromShare(_require, "src" + generateFileName(testName), int64(totalFileSize), srClient)
//
//	_, err = srcfClient.UploadRange(ctx, int64(srcOffset), expectedDataReader, nil)
//	_require.Nil(err)
//
//	// generate a URL with SAS pointing to the source file
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	_require.Nil(err)
//
//	sasQueryParams, err := FileSASSignatureValues{
//		Protocol:    SASProtocolHTTPS,
//		ExpiryTime:  time.Now().UTC().Add(48 * time.Hour),
//		ShareName:   shareName,
//		Permissions: FileSASPermissions{Create: true, Read: true, Write: true, Delete: true}.String(),
//	}.NewSASQueryParameters(credential)
//	_require.Nil(err)
//	srcfClient.u.RawQuery = sasQueryParams.Encode()
//	rawSrcURL := srcfClient.URL()
//
//	// create the destination file
//	dstfClient := createNewFileFromShare(_require, "dst" + generateFileName(testName), int64(totalFileSize), srClient)
//
//	// invoke UploadRange on dstfClient and put the data at a random range
//	// source and destination have different offsets, so we can test both values at the same time
//	dstOffset := 100
//	uploadFromURLResp, err := dstfClient.UploadRangeFromURL(ctx, rawSrcURL, int64(srcOffset), int64(dstOffset), int64(expectedDataSize), nil)
//	_require.Nil(err)
//	_require.Equal(uploadFromURLResp.RawResponse.StatusCode, 201)
//
//	// verify the destination
//	resp, err := dstfClient.Download(ctx, int64(dstOffset), int64(expectedDataSize), &DownloadFileOptions{RangeGetContentMD5: to.BoolPtr(false)})
//	_require.Nil(err)
//	downloadedData, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_require.Nil(err)
//	_require.EqualValues(downloadedData, expectedData)
//}

// Testings for GetRangeList and ClearRange
func (s *azfileLiveTestSuite) TestGetRangeListNonDefaultExact() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)

	fileSize := int64(512 * 10)

	_, err = fClient.Create(ctx, &CreateFileOptions{FileContentLength: to.Int64Ptr(fileSize), FileHTTPHeaders: &FileHTTPHeaders{}})
	_require.Nil(err)

	defer delFile(_require, fClient)

	rsc, _ := generateData(1024)
	putResp, err := fClient.UploadRange(ctx, 0, rsc, nil)
	_require.Nil(err)
	_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotEqual(putResp.ETag, "")
	_require.NotNil(putResp.ContentMD5)
	_require.NotEqual(putResp.RequestID, "")
	_require.NotEqual(putResp.Version, "")
	_require.Equal(putResp.Date.IsZero(), false)

	rangeList, err := fClient.GetRangeList(ctx, 0, 1023, nil)
	_require.Nil(err)
	_require.Equal(rangeList.RawResponse.StatusCode, 200)
	_require.Equal(rangeList.LastModified.IsZero(), false)
	_require.NotEqual(rangeList.ETag, "")
	_require.Equal(*rangeList.FileContentLength, fileSize)
	_require.NotEqual(rangeList.RequestID, "")
	_require.NotEqual(rangeList.Version, "")
	_require.Equal(rangeList.Date.IsZero(), false)
	_requireLen(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(1023)})
}

// Default means clear the entire file's range
func (s *azfileLiveTestSuite) TestClearRangeDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := createNewFileFromShare(_require, generateFileName(testName), 2048, srClient)
	defer delFile(_require, fClient)

	rsc, _ := generateData(2048)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	_require.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 0, 2048, nil)
	_require.Nil(err)
	_require.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	_requireLen(rangeList.Ranges, 0)
}

func (s *azfileLiveTestSuite) TestClearRangeNonDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := createNewFileFromShare(_require, generateFileName(testName), 4096, srClient)
	defer delFile(_require, fClient)

	rsc, _ := generateData(2048)
	_, err = fClient.UploadRange(ctx, 2048, rsc, nil)
	_require.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 2048, 2048, nil)
	_require.Nil(err)
	_require.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	_requireLen(rangeList.Ranges, 0)
}

func (s *azfileLiveTestSuite) TestClearRangeMultipleRanges() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := createNewFileFromShare(_require, generateFileName(testName), 2048, srClient)
	defer delFile(_require, fClient)

	rsc, _ := generateData(2048)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	_require.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 1024, 1024, nil)
	_require.Nil(err)
	_require.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	_requireLen(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(1023)})
}

// When not 512 aligned, clear range will set 0 the non-512 aligned range, and will not eliminate the range.
func (s *azfileLiveTestSuite) TestClearRangeNonDefaultCount() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)

	fClient := createNewFileFromShare(_require, generateFileName(testName), int64(1), srClient)
	defer delFile(_require, fClient)

	d := []byte{1}
	_, err = fClient.UploadRange(ctx, 0, NopCloser(bytes.NewReader(d)), nil)
	_require.Nil(err)

	clearResp, err := fClient.ClearRange(ctx, 0, 1, nil)
	_require.Nil(err)
	_require.Equal(clearResp.RawResponse.StatusCode, 201)

	rangeList, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	_requireLen(rangeList.Ranges, 1)
	_require.EqualValues(*rangeList.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(0)})

	dResp, err := fClient.Download(ctx, 0, CountToEnd, nil)
	_require.Nil(err)

	_bytes, err := ioutil.ReadAll(dResp.Body(RetryReaderOptions{}))
	_require.Nil(err)
	_require.EqualValues(_bytes, []byte{0})
}

//// Don't check offset by design.
//// func (s *azfileLiveTestSuite) TestFileClearRangeNegativeInvalidOffset() {
//// 	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//// 	srClient, _ := getsrClient(c, fsu)
//// 	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)
//
//// 	_, err := fClient.ClearRange(ctx, -1, 1)
//// 	_require.NotNil(err)
//// 	_require(strings.Contains(err.Error(), "offset must be >= 0"), chk.Equals, true)
//// }

func (s *azfileLiveTestSuite) TestFileClearRangeNegativeInvalidCount() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	srClient, _ := getShareClient(generateShareName(testName), svcClient)
	fClient := getFileClientFromShare(_require, generateFileName(testName), srClient)

	_, err = fClient.ClearRange(ctx, 0, 0, nil)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), "invalid argument: either offset is < 0 or count <= 0"), true)
}

func setupGetRangeListTest(_assert *assert.Assertions, testName string) (srClient ShareClient, fClient FileClient) {
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	_require.Nil(err)
	srClient = createNewShare(_require, generateShareName(testName), svcClient)
	fClient = createNewFileFromShare(_require, generateFileName(testName), int64(testFileRangeSize), srClient)

	rsc, _ := generateData(testFileRangeSize)
	_, err = fClient.UploadRange(ctx, 0, rsc, nil)
	_require.Nil(err)
	return
}

func validateBasicGetRangeList(_assert *assert.Assertions, resp FileGetRangeListResponse, err error) {
	_require.Nil(err)
	_requireLen(resp.Ranges, 1)
	_require.EqualValues(*resp.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(int64(testFileRangeSize - 1))})
}

func (s *azfileLiveTestSuite) TestFileGetRangeListDefaultEmptyFile() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)
	defer delShare(_require, srClient, nil)
	fClient := createNewFileFromShare(_require, generateFileName(testName), 0, srClient)

	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	_requireLen(resp.Ranges, 0)
}

func (s *azfileLiveTestSuite) TestFileGetRangeListDefaultRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_require, testName)
	defer delShare(_require, srClient, nil)

	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	validateBasicGetRangeList(_require, resp, err)
}

func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRanges() {
	_require := require.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_require, testName)
	defer delShare(_require, srClient, nil)

	_, err := fClient.Resize(ctx, int64(testFileRangeSize*3), nil)
	_require.Nil(err)

	rsc, _ := generateData(testFileRangeSize)
	_, err = fClient.UploadRange(ctx, int64(testFileRangeSize*2), rsc, nil)
	_require.Nil(err)
	resp, err := fClient.GetRangeList(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	_requireLen(resp.Ranges, 2)
	_require.EqualValues(*resp.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(int64(testFileRangeSize - 1))})
	_require.EqualValues(*resp.Ranges[1], FileRange{Start: to.Int64Ptr(int64(testFileRangeSize * 2)), End: to.Int64Ptr(int64((testFileRangeSize * 3) - 1))})
}

func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRangesCountLess() {
	_require := require.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_require, testName)
	defer delShare(_require, srClient, nil)

	resp, err := fClient.GetRangeList(ctx, 0, int64(testFileRangeSize-1), nil)
	_require.Nil(err)
	_requireLen(resp.Ranges, 1)
	_require.EqualValues(*resp.Ranges[0], FileRange{Start: to.Int64Ptr(0), End: to.Int64Ptr(int64(testFileRangeSize - 1))})
}

func (s *azfileLiveTestSuite) TestFileGetRangeListNonContiguousRangesCountExceed() {
	_require := require.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_require, testName)
	defer delShare(_require, srClient, nil)

	resp, err := fClient.GetRangeList(ctx, int64(0), int64(testFileRangeSize+1), nil)
	_require.Nil(err)
	validateBasicGetRangeList(_require, resp, err)
}

func (s *azfileLiveTestSuite) TestFileGetRangeListSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	srClient, fClient := setupGetRangeListTest(_require, testName)
	delShareInclude := DeleteSnapshotsOptionTypeInclude
	defer delShare(_require, srClient, &ShareDeleteOptions{DeleteSnapshots: &delShareInclude})

	resp, _ := srClient.CreateSnapshot(ctx, &ShareCreateSnapshotOptions{Metadata: map[string]string{}})
	_require.NotNil(resp.Snapshot)
	fClientWithSnapshot := fClient.WithSnapshot(*resp.Snapshot)

	resp2, err := fClientWithSnapshot.GetRangeList(ctx, 0, CountToEnd, nil)
	_require.Nil(err)
	validateBasicGetRangeList(_require, resp2, err)
}

func (s *azfileLiveTestSuite) TestUnexpectedEOFRecovery() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	srClient := createNewShare(_require, generateShareName(testName), svcClient)

	defer delShare(_require, srClient, &ShareDeleteOptions{
		DeleteSnapshots: &deleteSnapshotsInclude,
	})

	fClient := createNewFileFromShare(_require, generateFileName(testName), 2048, srClient)

	contentR, contentD := generateData(2048)

	resp, err := fClient.UploadRange(ctx, 0, contentR, nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, http.StatusCreated)
	_require.NotEqual(resp.RequestID, "")

	dlResp, err := fClient.Download(ctx, 0, 2048, nil)
	_require.Nil(err)

	// Verify that we can inject errors first.
	reader := dlResp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))

	_, err = ioutil.ReadAll(reader)
	_require.NotNil(err)
	_require.Equal(err.Error(), "unrecoverable error")

	// Then inject the retryable error.
	reader = dlResp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))

	buf, err := ioutil.ReadAll(reader)
	_require.Nil(err)
	_require.EqualValues(buf, contentD)
}

func (s *azfileLiveTestSuite) TestCreateMaximumSizeFileShare() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	srClient, err := getShareClient(generateShareName(testName), svcClient)
	cResp, err := srClient.Create(ctx, &ShareCreateOptions{
		Quota: &fileShareMaxQuota,
	})
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	defer delShare(_require, srClient, &ShareDeleteOptions{DeleteSnapshots: &deleteSnapshotsInclude})
	dirClient, err := srClient.NewRootDirectoryClient()
	_require.Nil(err)
	fClient := getFileClientFromDirectory(_require, generateFileName(testName), dirClient)
	_, err = fClient.Create(ctx, &CreateFileOptions{
		FileContentLength: &fileMaxAllowedSizeInBytes,
		FileHTTPHeaders:   &FileHTTPHeaders{},
	})
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)
}
