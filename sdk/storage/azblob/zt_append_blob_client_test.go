//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"time"
)

// nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlock() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	resp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)

	appendResp, err := abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	_require.Equal(appendResp.RawResponse.StatusCode, 201)
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

	appendResp, err = abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	_require.Equal(*appendResp.BlobAppendOffset, "1024")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(2))
}

// nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	// set up abClient to test
	abClient, _ := containerClient.NewAppendBlobClient(generateBlobName(testName))
	resp, err := abClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)

	// test append block with valid MD5 value
	readerToBody, body := getRandomDataAndReader(1024)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]
	appendBlockOptions := AppendBlobAppendBlockOptions{
		TransactionalContentMD5: contentMD5,
	}
	appendResp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(readerToBody), &appendBlockOptions)
	_require.Nil(err)
	_require.Equal(appendResp.RawResponse.StatusCode, 201)
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
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	_ = body
	appendBlockOptions = AppendBlobAppendBlockOptions{
		TransactionalContentMD5: badMD5,
	}
	appendResp, err = abClient.AppendBlock(context.Background(), internal.NopCloser(readerToBody), &appendBlockOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeMD5Mismatch)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	//ctx := context.Background()
	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	contentMD5 := md5.Sum(sourceData)
	srcBlob, _ := containerClient.NewAppendBlobClient(generateName("appendsrc"))
	destBlob, _ := containerClient.NewAppendBlobClient(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp1.RawResponse.StatusCode, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_require.Nil(err)
	_require.Nil(err)
	_require.Equal(appendResp.RawResponse.StatusCode, 201)
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

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Append block from URL.
	cResp2, err := destBlob.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp2.RawResponse.StatusCode, 201)

	//ctx context.Context, source url.URL, contentLength int64, options *AppendBlobAppendBlockFromURLOptions)
	offset := int64(0)
	count := int64(CountToEnd)
	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
		Offset: &offset,
		Count:  &count,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_require.Nil(err)
	_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendFromURLResp.ETag)
	_require.NotNil(appendFromURLResp.LastModified)
	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_require.NotNil(appendFromURLResp.ContentMD5)
	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5[:])
	_require.NotNil(appendFromURLResp.RequestID)
	_require.NotNil(appendFromURLResp.Version)
	_require.NotNil(appendFromURLResp.Date)
	_require.Equal((*appendFromURLResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	_require.Nil(err)

	destData, err := io.ReadAll(downloadResp.RawResponse.Body)
	_require.Nil(err)
	_require.Equal(destData, sourceData)
	_ = downloadResp.Body(nil).Close()
}

// nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURLWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	ctx := context.Background() // Use default Background context
	srcBlob, _ := containerClient.NewAppendBlobClient(generateName("appendsrc"))
	destBlob, _ := containerClient.NewAppendBlobClient(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(cResp1.RawResponse.StatusCode, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_require.Nil(err)
	_require.Equal(appendResp.RawResponse.StatusCode, 201)
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

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Append block from URL.
	cResp2, err := destBlob.Create(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	contentMD5 := md5Value[:]
	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: contentMD5,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_require.Nil(err)
	_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendFromURLResp.ETag)
	_require.NotNil(appendFromURLResp.LastModified)
	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_require.NotNil(appendFromURLResp.ContentMD5)
	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
	_require.NotNil(appendFromURLResp.RequestID)
	_require.NotNil(appendFromURLResp.Version)
	_require.NotNil(appendFromURLResp.Date)
	_require.Equal((*appendFromURLResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	_require.Nil(err)
	destData, err := io.ReadAll(downloadResp.Body(nil))
	_require.Nil(err)
	_require.EqualValues(destData, sourceData)

	// Test append block from URL with bad MD5 value
	_, badMD5 := getRandomDataAndReader(16)
	appendBlockURLOptions = AppendBlobAppendBlockFromURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: badMD5,
	}
	_, err = destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobCreateAppendMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(ctx, &AppendBlobCreateOptions{
		Metadata: basicMetadata,
	})
	_require.Nil(err)

	resp, err := abClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobCreateAppendMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := AppendBlobCreateOptions{
		Metadata: map[string]string{},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)

	resp, err := abClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestBlobCreateAppendMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := AppendBlobCreateOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.NotNil(err)
	_require.Contains(err.Error(), invalidHeaderErrorSubstring)
}

func (s *azblobTestSuite) TestBlobCreateAppendHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)

	resp, err := abClient.GetProperties(ctx, nil)
	_require.Nil(err)
	h := resp.GetHTTPHeaders()
	_require.EqualValues(h, basicHeaders)
}

func validateAppendBlobPut(_require *require.Assertions, abClient *AppendBlobClient) {
	resp, err := abClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, basicMetadata)
	_require.EqualValues(resp.GetHTTPHeaders(), basicHeaders)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: to.Ptr("garbage"),
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	eTag := "garbage"
	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := AppendBlobCreateOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockNilBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(bytes.NewReader(nil)), nil)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobAppendBlockEmptyBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader("")), nil)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobAppendBlockNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeBlobNotFound)
}

func validateBlockAppended(_require *require.Assertions, abClient *AppendBlobClient, expectedSize int) {
	resp, err := abClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(expectedSize))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient, _ := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_require.Nil(err)
	_require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: to.Ptr("garbage"),
			},
		},
	})
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: to.Ptr("garbage"),
			},
		},
	})
	_require.Nil(err)
	validateBlockAppended(_require, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

//// TODO: Fix this
////func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNegOne() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_require, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	appendPosition := int64(-1)
////	appendBlockOptions := AppendBlobAppendBlockOptions{
////		AppendPositionAccessConditions: &AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err := abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions) // This will cause the library to set the value of the header to 0
////	_require.NotNil(err)
////
////	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
////}
//
////func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchZero() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_require, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	_, err := abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil) // The position will not match, but the condition should be ignored
////	_require.Nil(err)
////
////	appendPosition := int64(0)
////	appendBlockOptions := AppendBlobAppendBlockOptions{
////		AppendPositionAccessConditions: &AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
////	_require.Nil(err)
////
////	validateBlockAppended(c, abClient, 2*len(blockBlobDefaultData))
////}

func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNonZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Ptr(int64(len(blockBlobDefaultData))),
		},
	})
	_require.Nil(err)

	validateBlockAppended(_require, abClient, len(blockBlobDefaultData)*2)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Ptr[int64](-1),
		},
	})
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNonZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Ptr[int64](12),
		},
	})
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeAppendPositionConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMaxSizeTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: to.Ptr(int64(len(blockBlobDefaultData) + 1)),
		},
	})
	_require.Nil(err)
	validateBlockAppended(_require, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMaxSizeFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlobAppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: to.Ptr(int64(len(blockBlobDefaultData) - 1)),
		},
	})
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeMaxBlobSizeConditionNotMet)
}

func (s *azblobTestSuite) TestSealAppendBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	appendResp, err := abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	_require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))

	sealResp, err := abClient.SealAppendBlob(ctx, nil)
	_require.Nil(err)
	_require.Equal(*sealResp.IsSealed, true)

	appendResp, err = abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_require.NotNil(err)
	validateStorageError(_require, err, "BlobIsSealed")

	getPropResp, err := abClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getPropResp.IsSealed, true)
}

// TODO: Learn about the behaviour of AppendPosition
// nolint
//func (s *azblobUnrecordedTestSuite) TestSealAppendBlobWithAppendConditions() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	abName := generateBlobName(testName)
//	abClient := createNewAppendBlob(_require, abName, containerClient)
//
//	sealResp, err := abClient.SealAppendBlob(ctx, &AppendBlobSealOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: to.Ptr(1),
//		},
//	})
//	_require.NotNil(err)
//	_ = sealResp
//
//	sealResp, err = abClient.SealAppendBlob(ctx, &AppendBlobSealOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: to.Ptr(0),
//		},
//	})
//}

func (s *azblobTestSuite) TestCopySealedBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	_, err = abClient.SealAppendBlob(ctx, nil)
	_require.Nil(err)

	copiedBlob1, _ := getAppendBlobClient("copy1"+abName, containerClient)
	// copy sealed blob will get a sealed blob
	_, err = copiedBlob1.StartCopyFromURL(ctx, abClient.URL(), nil)
	_require.Nil(err)

	getResp1, err := copiedBlob1.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp1.IsSealed, true)

	_, err = copiedBlob1.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_require.NotNil(err)
	validateStorageError(_require, err, "BlobIsSealed")

	copiedBlob2, _ := getAppendBlobClient("copy2"+abName, containerClient)
	_, err = copiedBlob2.StartCopyFromURL(ctx, abClient.URL(), &BlobStartCopyOptions{
		SealBlob: to.Ptr(true),
	})
	_require.Nil(err)

	getResp2, err := copiedBlob2.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp2.IsSealed, true)

	_, err = copiedBlob2.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_require.NotNil(err)
	validateStorageError(_require, err, "BlobIsSealed")

	copiedBlob3, _ := getAppendBlobClient("copy3"+abName, containerClient)
	_, err = copiedBlob3.StartCopyFromURL(ctx, abClient.URL(), &BlobStartCopyOptions{
		SealBlob: to.Ptr(false),
	})
	_require.Nil(err)

	getResp3, err := copiedBlob3.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(getResp3.IsSealed)

	appendResp3, err := copiedBlob3.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_require.Nil(err)
	_require.Equal(appendResp3.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp3.BlobAppendOffset, "0")
	_require.Equal(*appendResp3.BlobCommittedBlockCount, int32(1))
}

func (s *azblobTestSuite) TestCopyUnsealedBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_require, abName, containerClient)

	copiedBlob, _ := getAppendBlobClient("copy"+abName, containerClient)
	_, err = copiedBlob.StartCopyFromURL(ctx, abClient.URL(), &BlobStartCopyOptions{
		SealBlob: to.Ptr(true),
	})
	_require.Nil(err)

	getResp, err := copiedBlob.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp.IsSealed, true)
}
