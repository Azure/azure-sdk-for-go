//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package appendblob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	suite.Run(t, &AppendBlobRecordedTestsSuite{})
	//suite.Run(t, &AppendBlobUnrecordedTestsSuite{})
}

// nolint
func (s *AppendBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

// nolint
func (s *AppendBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

// nolint
func (s *AppendBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

// nolint
func (s *AppendBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type AppendBlobRecordedTestsSuite struct {
	suite.Suite
}

type AppendBlobUnrecordedTestsSuite struct {
	suite.Suite
}

func getAppendBlobClient(appendBlobName string, containerClient *container.Client) *appendblob.Client {
	return containerClient.NewAppendBlobClient(appendBlobName)
}

func createNewAppendBlob(ctx context.Context, _require *require.Assertions, appendBlobName string, containerClient *container.Client) *appendblob.Client {
	abClient := getAppendBlobClient(appendBlobName, containerClient)

	_, err := abClient.Create(ctx, nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	return abClient
}

// nolint
func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlock() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	appendResp, err := abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendResp.ETag)
	_require.NotNil(appendResp.LastModified)
	_require.Equal((*appendResp.LastModified).IsZero(), false)
	_require.Nil(appendResp.ContentMD5)
	_require.NotNil(appendResp.RequestID)
	_require.NotNil(appendResp.Version)
	_require.NotNil(appendResp.Date)
	_require.Equal((*appendResp.Date).IsZero(), false)

	appendResp, err = abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	_require.Equal(*appendResp.BlobAppendOffset, "1024")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(2))
}

// nolint
func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))
	_, err = abClient.Create(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	// test append block with valid MD5 value
	readerToBody, body := testcommon.GetRandomDataAndReader(1024)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]
	appendBlockOptions := appendblob.AppendBlockOptions{
		TransactionalContentMD5: contentMD5,
	}
	appendResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.Nil(err)
	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendResp.ETag)
	_require.NotNil(appendResp.LastModified)
	_require.Equal((*appendResp.LastModified).IsZero(), false)
	_require.EqualValues(appendResp.ContentMD5, contentMD5)
	_require.NotNil(appendResp.RequestID)
	_require.NotNil(appendResp.Version)
	_require.NotNil(appendResp.Date)
	_require.Equal((*appendResp.Date).IsZero(), false)

	// test append block with bad MD5 value
	readerToBody, body = testcommon.GetRandomDataAndReader(1024)
	_, badMD5 := testcommon.GetRandomDataAndReader(16)
	_ = body
	appendBlockOptions = appendblob.AppendBlockOptions{
		TransactionalContentMD5: badMD5,
	}
	appendResp, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
}

//nolint
//func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	//ctx := ctx
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := testcommon.GetRandomDataAndReader(contentSize)
//	contentMD5 := md5.Sum(sourceData)
//	srcBlob := containerClient.NewAppendBlobClient(generateName("appendsrc"))
//	destBlob := containerClient.NewAppendBlobClient(generateName("appenddest"))
//
//	// Prepare source abClient for copy.
//	_, err = srcBlob.Create(context.Background(), nil)
//	_require.Nil(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	appendResp, err := srcBlob.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
//	_require.Nil(err)
//	_require.Nil(err)
//	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendResp.BlobAppendOffset, "0")
//	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendResp.ETag)
//	_require.NotNil(appendResp.LastModified)
//	_require.Equal((*appendResp.LastModified).IsZero(), false)
//	_require.Nil(appendResp.ContentMD5)
//	_require.NotNil(appendResp.RequestID)
//	_require.NotNil(appendResp.Version)
//	_require.NotNil(appendResp.Date)
//	_require.Equal((*appendResp.Date).IsZero(), false)
//
//	// Get source abClient URL with SAS for AppendBlockFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Append block from URL.
//	_, err = destBlob.Create(context.Background(), nil)
//	_require.Nil(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	//ctx context.Context, source url.URL, contentLength int64, options *appendblob.AppendBlockFromURLOptions)
//	offset := int64(0)
//	count := int64(CountToEnd)
//	appendBlockURLOptions := appendblob.AppendBlockFromURLOptions{
//		Offset: &offset,
//		Count:  &count,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.Nil(err)
//	//_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
//	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendFromURLResp.ETag)
//	_require.NotNil(appendFromURLResp.LastModified)
//	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
//	_require.NotNil(appendFromURLResp.ContentMD5)
//	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5[:])
//	_require.NotNil(appendFromURLResp.RequestID)
//	_require.NotNil(appendFromURLResp.Version)
//	_require.NotNil(appendFromURLResp.Date)
//	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(context.Background(), nil)
//	_require.Nil(err)
//
//	destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//	_require.Nil(err)
//	_require.Equal(destData, sourceData)
//	_ = downloadresp.BodyReader(nil).Close()
//}

//nolint
//func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithMD5() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := testcommon.GetRandomDataAndReader(contentSize)
//	md5Value := md5.Sum(sourceData)
//	ctx := ctx // Use default Background context
//	srcBlob := containerClient.NewAppendBlobClient(generateName("appendsrc"))
//	destBlob := containerClient.NewAppendBlobClient(generateName("appenddest"))
//
//	// Prepare source abClient for copy.
//	_, err = srcBlob.Create(context.Background(), nil)
//	_require.Nil(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	appendResp, err := srcBlob.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
//	_require.Nil(err)
//	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendResp.BlobAppendOffset, "0")
//	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendResp.ETag)
//	_require.NotNil(appendResp.LastModified)
//	_require.Equal((*appendResp.LastModified).IsZero(), false)
//	_require.Nil(appendResp.ContentMD5)
//	_require.NotNil(appendResp.RequestID)
//	_require.NotNil(appendResp.Version)
//	_require.NotNil(appendResp.Date)
//	_require.Equal((*appendResp.Date).IsZero(), false)
//
//	// Get source abClient URL with SAS for AppendBlockFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Append block from URL.
//	_, err = destBlob.Create(context.Background(), nil)
//	_require.Nil(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	offset := int64(0)
//	count := int64(contentSize)
//	contentMD5 := md5Value[:]
//	appendBlockURLOptions := appendblob.AppendBlockFromURLOptions{
//		Offset:           &offset,
//		Count:            &count,
//		SourceContentMD5: contentMD5,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.Nil(err)
//	//_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
//	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendFromURLResp.ETag)
//	_require.NotNil(appendFromURLResp.LastModified)
//	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
//	_require.NotNil(appendFromURLResp.ContentMD5)
//	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
//	_require.NotNil(appendFromURLResp.RequestID)
//	_require.NotNil(appendFromURLResp.Version)
//	_require.NotNil(appendFromURLResp.Date)
//	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(context.Background(), nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(nil))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
//
//	// Test append block from URL with bad MD5 value
//	_, badMD5 := testcommon.GetRandomDataAndReader(16)
//	appendBlockURLOptions = appendblob.AppendBlockFromURLOptions{
//		Offset:           &offset,
//		Count:            &count,
//		SourceContentMD5: badMD5,
//	}
//	_, err = destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.NotNil(err)
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
//}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(context.Background(), &appendblob.CreateOptions{
		Metadata: testcommon.BasicMetadata,
	})
	_require.Nil(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		Metadata: map[string]string{},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NotNil(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	h := blob.ParseHTTPHeaders(resp)
	_require.EqualValues(h, testcommon.BasicHeaders)
}

func validateAppendBlobPut(_require *require.Assertions, abClient *appendblob.Client) {
	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
	_require.EqualValues(blob.ParseHTTPHeaders(resp), testcommon.BasicHeaders)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)

	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr(azcore.ETag("garbage")),
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	eTag := azcore.ETag("garbage")
	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockNilBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(bytes.NewReader(nil)), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockEmptyBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("")), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func validateBlockAppended(_require *require.Assertions, abClient *appendblob.Client, expectedSize int) {
	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(expectedSize))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr(azcore.ETag("garbage")),
			},
		},
	})
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: to.Ptr(azcore.ETag("garbage")),
			},
		},
	})
	_require.Nil(err)
	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

//// TODO: Fix this
////func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNegOne() {
////	bsu := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	appendPosition := int64(-1)
////	appendBlockOptions := appendblob.AppendBlockOptions{
////		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions) // This will cause the library to set the value of the header to 0
////	_require.NotNil(err)
////
////	validateBlockAppended(c, abClient, len(testcommon.BlockBlobDefaultData))
////}
//
////func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchZero() {
////	bsu := testcommon.GetServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	_, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil) // The position will not match, but the condition should be ignored
////	_require.Nil(err)
////
////	appendPosition := int64(0)
////	appendBlockOptions := appendblob.AppendBlockOptions{
////		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
////	_require.Nil(err)
////
////	validateBlockAppended(c, abClient, 2*len(testcommon.BlockBlobDefaultData))
////}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNonZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			AppendPosition: to.Ptr(int64(len(testcommon.BlockBlobDefaultData))),
		},
	})
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData)*2)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			AppendPosition: to.Ptr[int64](-1),
		},
	})
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNonZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			AppendPosition: to.Ptr[int64](12),
		},
	})
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.AppendPositionConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMaxSizeTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			MaxSize: to.Ptr(int64(len(testcommon.BlockBlobDefaultData) + 1)),
		},
	})
	_require.Nil(err)
	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMaxSizeFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			MaxSize: to.Ptr(int64(len(testcommon.BlockBlobDefaultData) - 1)),
		},
	})
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MaxBlobSizeConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestSeal() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	appendResp, err := abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))

	sealResp, err := abClient.Seal(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*sealResp.IsSealed, true)

	appendResp, err = abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, "BlobIsSealed")

	getPropResp, err := abClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getPropResp.IsSealed, true)
}

// TODO: Learn about the behaviour of AppendPosition
// nolint
//func (s *AppendBlobUnrecordedTestsSuite) TestSealWithAppendConditions() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	abName := testcommon.GenerateBlobName(testName)
//	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)
//
//	sealResp, err := abClient.Seal(context.Background(), &AppendBlobSealOptions{
//		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
//			AppendPosition: to.Ptr(1),
//		},
//	})
//	_require.NotNil(err)
//	_ = sealResp
//
//	sealResp, err = abClient.Seal(context.Background(), &AppendBlobSealOptions{
//		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
//			AppendPosition: to.Ptr(0),
//		},
//	})
//}

func (s *AppendBlobRecordedTestsSuite) TestCopySealedBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.Seal(context.Background(), nil)
	_require.Nil(err)

	copiedBlob1 := getAppendBlobClient("copy1"+abName, containerClient)
	// copy sealed blob will get a sealed blob
	_, err = copiedBlob1.StartCopyFromURL(context.Background(), abClient.URL(), nil)
	_require.Nil(err)

	getResp1, err := copiedBlob1.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp1.IsSealed, true)

	_, err = copiedBlob1.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, "BlobIsSealed")

	copiedBlob2 := getAppendBlobClient("copy2"+abName, containerClient)
	_, err = copiedBlob2.StartCopyFromURL(context.Background(), abClient.URL(), &blob.StartCopyFromURLOptions{
		SealBlob: to.Ptr(true),
	})
	_require.Nil(err)

	getResp2, err := copiedBlob2.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp2.IsSealed, true)

	_, err = copiedBlob2.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, "BlobIsSealed")

	copiedBlob3 := getAppendBlobClient("copy3"+abName, containerClient)
	_, err = copiedBlob3.StartCopyFromURL(context.Background(), abClient.URL(), &blob.StartCopyFromURLOptions{
		SealBlob: to.Ptr(false),
	})
	_require.Nil(err)

	getResp3, err := copiedBlob3.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(getResp3.IsSealed)

	appendResp3, err := copiedBlob3.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	//_require.Equal(appendResp3.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp3.BlobAppendOffset, "0")
	_require.Equal(*appendResp3.BlobCommittedBlockCount, int32(1))
}

func (s *AppendBlobRecordedTestsSuite) TestCopyUnsealedBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	copiedBlob := getAppendBlobClient("copy"+abName, containerClient)
	_, err = copiedBlob.StartCopyFromURL(context.Background(), abClient.URL(), &blob.StartCopyFromURLOptions{
		SealBlob: to.Ptr(true),
	})
	_require.Nil(err)

	getResp, err := copiedBlob.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp.IsSealed, true)
}

func (s *AppendBlobUnrecordedTestsSuite) TestCreateAppendBlobWithTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	// _context := getTestContext(testName)
	// ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		Tags: testcommon.SpecialCharBlobTagsMap,
	}
	createResp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	_require.NotNil(createResp.VersionID)

	_, err = abClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	blobGetTagsResponse, err := abClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.SpecialCharBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.SpecialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	resp, err := abClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	snapshotURL, _ := abClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp2.TagCount, int64(len(testcommon.SpecialCharBlobTagsMap)))

	blobGetTagsResponse, err = abClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet = blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.SpecialCharBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.SpecialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobGetPropertiesUsingVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	abClient := createNewAppendBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	blobProp, _ := abClient.GetProperties(context.Background(), nil)

	createAppendBlobOptions := appendblob.CreateOptions{
		Metadata: testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	createResp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	_require.NotNil(createResp.VersionID)
	blobProp, _ = abClient.GetProperties(context.Background(), nil)
	_require.EqualValues(createResp.VersionID, blobProp.VersionID)
	_require.EqualValues(createResp.LastModified, blobProp.LastModified)
	_require.Equal(*createResp.ETag, *blobProp.ETag)
	_require.Equal(*blobProp.IsCurrentVersion, true)
}

func (s *AppendBlobUnrecordedTestsSuite) TestSetBlobMetadataReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)

	metadata := map[string]string{"test_key_1": "test_value_1", "test_key_2": "2019"}
	resp, err := bbClient.SetMetadata(context.Background(), metadata, nil)
	_require.Nil(err)
	_require.NotNil(resp.VersionID)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Metadata: true},
	})

	if pager.More() {
		pageResp, err := pager.NextPage(context.Background())
		_require.Nil(err) // check for an error first
		//s.T().Fail()      // no page was gotten

		_require.NotNil(pageResp.Segment.BlobItems)
		blobList := pageResp.Segment.BlobItems
		_require.Len(blobList, 1)
		blobResp1 := blobList[0]
		_require.Equal(*blobResp1.Name, bbName)
		_require.NotNil(blobResp1.Metadata)
		_require.Len(blobResp1.Metadata, 2)
	}
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	createAppendBlobOptions := appendblob.CreateOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := appendblob.AppendBlockOptions{
			CpkInfo: &testcommon.TestCPKByValue,
		}
		resp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal((*resp.Date).IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		_require.EqualValues(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}

	// Get blob content without encryption key should fail the request.
	_, err = abClient.DownloadStream(context.Background(), nil)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	downloadResp, err := abClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)

	data, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	createAppendBlobOptions := appendblob.CreateOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := appendblob.AppendBlockOptions{
			CpkScopeInfo: &testcommon.TestCPKByScope,
		}
		resp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal((*resp.Date).IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		_require.EqualValues(resp.EncryptionScope, testcommon.TestCPKByScope.EncryptionScope)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	downloadResp, err := abClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)

	data, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
}

//nolint
//func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCPK() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Sum := md5.Sum(srcData)
//	contentMD5 := md5Sum[:]
//	ctx := ctx
//	srcABClient := containerClient.NewAppendBlobClient(generateName("src"))
//	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))
//
//	_, err = srcABClient.Create(ctx, nil)
//	_require.Nil(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	resp, err := srcABClient.AppendBlock(ctx, streaming.NopCloser(r), nil)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 201)
//	_require.Equal(*resp.BlobAppendOffset, "0")
//	_require.Equal(*resp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(resp.ETag)
//	_require.NotNil(resp.LastModified)
//	_require.Equal((*resp.LastModified).IsZero(), false)
//	_require.Nil(resp.ContentMD5)
//	_require.NotNil(resp.RequestID)
//	_require.NotNil(resp.Version)
//	_require.NotNil(resp.Date)
//	_require.Equal((*resp.Date).IsZero(), false)
//
//	srcBlobParts, _ := NewBlobURLParts(srcABClient.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	createAppendBlobOptions := appendblob.CreateOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	_, err = destBlob.Create(ctx, &createAppendBlobOptions)
//	_require.Nil(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	offset := int64(0)
//	count := int64(contentSize)
//	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
//		Offset:  &offset,
//		Count:   &count,
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.Nil(err)
//	//_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
//	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendFromURLResp.ETag)
//	_require.NotNil(appendFromURLResp.LastModified)
//	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
//	_require.NotNil(appendFromURLResp.ContentMD5)
//	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
//	_require.NotNil(appendFromURLResp.RequestID)
//	_require.NotNil(appendFromURLResp.Version)
//	_require.NotNil(appendFromURLResp.Date)
//	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
//	_require.Equal(*appendFromURLResp.IsServerEncrypted, true)
//
//	// Get blob content without encryption key should fail the request.
//	_, err = destBlob.DownloadStream(ctx, nil)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestInvalidCPKByValue,
//	}
//	_, err = destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions = blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//
//	_require.Equal(*downloadResp.IsServerEncrypted, true)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CpkInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
//}

//nolint
//func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCPKScope() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Sum := md5.Sum(srcData)
//	contentMD5 := md5Sum[:]
//	ctx := ctx
//	srcClient := containerClient.NewAppendBlobClient(generateName("src"))
//	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))
//
//	_, err = srcClient.Create(ctx, nil)
//	_require.Nil(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	resp, err := srcClient.AppendBlock(ctx, streaming.NopCloser(r), nil)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 201)
//	_require.Equal(*resp.BlobAppendOffset, "0")
//	_require.Equal(*resp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(resp.ETag)
//	_require.NotNil(resp.LastModified)
//	_require.Equal((*resp.LastModified).IsZero(), false)
//	_require.Nil(resp.ContentMD5)
//	_require.NotNil(resp.RequestID)
//	_require.NotNil(resp.Version)
//	_require.NotNil(resp.Date)
//	_require.Equal((*resp.Date).IsZero(), false)
//
//	srcBlobParts, _ := NewBlobURLParts(srcClient.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	createAppendBlobOptions := appendblob.CreateOptions{
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	_, err = destBlob.Create(ctx, &createAppendBlobOptions)
//	_require.Nil(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	offset := int64(0)
//	count := int64(contentSize)
//	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
//		Offset:       &offset,
//		Count:        &count,
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.Nil(err)
//	//_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
//	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendFromURLResp.ETag)
//	_require.NotNil(appendFromURLResp.LastModified)
//	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
//	_require.NotNil(appendFromURLResp.ContentMD5)
//	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
//	_require.NotNil(appendFromURLResp.RequestID)
//	_require.NotNil(appendFromURLResp.Version)
//	_require.NotNil(appendFromURLResp.Date)
//	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
//	_require.Equal(*appendFromURLResp.IsServerEncrypted, true)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.Equal(*downloadResp.IsServerEncrypted, true)
//	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CpkInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
//}
