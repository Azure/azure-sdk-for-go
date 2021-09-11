// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"time"
)

//nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlock() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	resp, err := abClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)

	appendResp, err := abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_assert.Nil(err)
	_assert.Equal(appendResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendResp.BlobAppendOffset, "0")
	_assert.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendResp.ETag)
	_assert.NotNil(appendResp.LastModified)
	_assert.Equal((*appendResp.LastModified).IsZero(), false)
	_assert.Nil(appendResp.ContentMD5)
	_assert.NotNil(appendResp.RequestID)
	_assert.NotNil(appendResp.Version)
	_assert.NotNil(appendResp.Date)
	_assert.Equal((*appendResp.Date).IsZero(), false)

	appendResp, err = abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_assert.Nil(err)
	_assert.Equal(*appendResp.BlobAppendOffset, "1024")
	_assert.Equal(*appendResp.BlobCommittedBlockCount, int32(2))
}

//nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockWithMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobClient(generateBlobName(testName))
	resp, err := abClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)

	// test append block with valid MD5 value
	readerToBody, body := getRandomDataAndReader(1024)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]
	appendBlockOptions := AppendBlockOptions{
		TransactionalContentMD5: contentMD5,
	}
	appendResp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(readerToBody), &appendBlockOptions)
	_assert.Nil(err)
	_assert.Equal(appendResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendResp.BlobAppendOffset, "0")
	_assert.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendResp.ETag)
	_assert.NotNil(appendResp.LastModified)
	_assert.Equal((*appendResp.LastModified).IsZero(), false)
	_assert.EqualValues(appendResp.ContentMD5, contentMD5)
	_assert.NotNil(appendResp.RequestID)
	_assert.NotNil(appendResp.Version)
	_assert.NotNil(appendResp.Date)
	_assert.Equal((*appendResp.Date).IsZero(), false)

	// test append block with bad MD5 value
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	_ = body
	appendBlockOptions = AppendBlockOptions{
		TransactionalContentMD5: badMD5,
	}
	appendResp, err = abClient.AppendBlock(context.Background(), internal.NopCloser(readerToBody), &appendBlockOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	//ctx := context.Background()
	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	contentMD5 := md5.Sum(sourceData)
	srcBlob := containerClient.NewAppendBlobClient(generateName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp1.RawResponse.StatusCode, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_assert.Nil(err)
	_assert.Nil(err)
	_assert.Equal(appendResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendResp.BlobAppendOffset, "0")
	_assert.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendResp.ETag)
	_assert.NotNil(appendResp.LastModified)
	_assert.Equal((*appendResp.LastModified).IsZero(), false)
	_assert.Nil(appendResp.ContentMD5)
	_assert.NotNil(appendResp.RequestID)
	_assert.NotNil(appendResp.Version)
	_assert.NotNil(appendResp.Date)
	_assert.Equal((*appendResp.Date).IsZero(), false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)

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
	_assert.Nil(err)
	_assert.Equal(cResp2.RawResponse.StatusCode, 201)

	//ctx context.Context, source url.URL, contentLength int64, options *AppendBlockURLOptions)
	offset := int64(0)
	count := int64(CountToEnd)
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset: &offset,
		Count:  &count,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_assert.Nil(err)
	_assert.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_assert.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendFromURLResp.ETag)
	_assert.NotNil(appendFromURLResp.LastModified)
	_assert.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_assert.NotNil(appendFromURLResp.ContentMD5)
	_assert.EqualValues(appendFromURLResp.ContentMD5, contentMD5[:])
	_assert.NotNil(appendFromURLResp.RequestID)
	_assert.NotNil(appendFromURLResp.Version)
	_assert.NotNil(appendFromURLResp.Date)
	_assert.Equal((*appendFromURLResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	_assert.Nil(err)

	destData, err := ioutil.ReadAll(downloadResp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(destData, sourceData)
	_ = downloadResp.Body(RetryReaderOptions{}).Close()
}

//nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURLWithMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	ctx := context.Background() // Use default Background context
	srcBlob := containerClient.NewAppendBlobClient(generateName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp1.RawResponse.StatusCode, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_assert.Nil(err)
	_assert.Equal(appendResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendResp.BlobAppendOffset, "0")
	_assert.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendResp.ETag)
	_assert.NotNil(appendResp.LastModified)
	_assert.Equal((*appendResp.LastModified).IsZero(), false)
	_assert.Nil(appendResp.ContentMD5)
	_assert.NotNil(appendResp.RequestID)
	_assert.NotNil(appendResp.Version)
	_assert.NotNil(appendResp.Date)
	_assert.Equal((*appendResp.Date).IsZero(), false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)

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
	_assert.Nil(err)
	_assert.Equal(cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	contentMD5 := md5Value[:]
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: contentMD5,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_assert.Nil(err)
	_assert.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_assert.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendFromURLResp.ETag)
	_assert.NotNil(appendFromURLResp.LastModified)
	_assert.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_assert.NotNil(appendFromURLResp.ContentMD5)
	_assert.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
	_assert.NotNil(appendFromURLResp.RequestID)
	_assert.NotNil(appendFromURLResp.Version)
	_assert.NotNil(appendFromURLResp.Date)
	_assert.Equal((*appendFromURLResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, sourceData)

	// Test append block from URL with bad MD5 value
	_, badMD5 := getRandomDataAndReader(16)
	appendBlockURLOptions = AppendBlockURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: badMD5,
	}
	_, err = destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobCreateAppendMetadataNonEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(ctx, &CreateAppendBlobOptions{
		Metadata: basicMetadata,
	})
	_assert.Nil(err)

	resp, err := abClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Metadata)
	_assert.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobCreateAppendMetadataEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: map[string]string{},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)

	resp, err := abClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestBlobCreateAppendMetadataInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.NotNil(err)
	_assert.Contains(err.Error(), invalidHeaderErrorSubstring)
}

func (s *azblobTestSuite) TestBlobCreateAppendHTTPHeaders() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)

	resp, err := abClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	h := resp.GetHTTPHeaders()
	_assert.EqualValues(h, basicHeaders)
}

func validateAppendBlobPut(_assert *assert.Assertions, abClient AppendBlobClient) {
	resp, err := abClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Metadata)
	_assert.EqualValues(resp.Metadata, basicMetadata)
	_assert.EqualValues(resp.GetHTTPHeaders(), basicHeaders)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)

	validateAppendBlobPut(_assert, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)

	validateAppendBlobPut(_assert, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)

	validateAppendBlobPut(_assert, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: to.StringPtr("garbage"),
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	eTag := "garbage"
	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)

	validateAppendBlobPut(_assert, abClient)
}

func (s *azblobTestSuite) TestBlobCreateAppendIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockNilBody() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(bytes.NewReader(nil)), nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobAppendBlockEmptyBody() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader("")), nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobAppendBlockNonExistentBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeBlobNotFound)
}

func validateBlockAppended(_assert *assert.Assertions, abClient AppendBlobClient, expectedSize int) {
	resp, err := abClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.ContentLength, int64(expectedSize))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_assert.Nil(err)

	validateBlockAppended(_assert, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_assert.Nil(err)

	validateBlockAppended(_assert, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_assert.NotNil(appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_assert.Nil(err)

	validateBlockAppended(_assert, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: to.StringPtr("garbage"),
			},
		},
	})
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: to.StringPtr("garbage"),
			},
		},
	})
	_assert.Nil(err)
	validateBlockAppended(_assert, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//// TODO: Fix this
////func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNegOne() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	appendPosition := int64(-1)
////	appendBlockOptions := AppendBlockOptions{
////		AppendPositionAccessConditions: &AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err := abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions) // This will cause the library to set the value of the header to 0
////	_assert.NotNil(err)
////
////	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
////}
//
////func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchZero() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	_, err := abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil) // The position will not match, but the condition should be ignored
////	_assert.Nil(err)
////
////	appendPosition := int64(0)
////	appendBlockOptions := AppendBlockOptions{
////		AppendPositionAccessConditions: &AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
////	_assert.Nil(err)
////
////	validateBlockAppended(c, abClient, 2*len(blockBlobDefaultData))
////}

func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNonZero() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Int64Ptr(int64(len(blockBlobDefaultData))),
		},
	})
	_assert.Nil(err)

	validateBlockAppended(_assert, abClient, len(blockBlobDefaultData)*2)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNegOne() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Int64Ptr(-1),
		},
	})
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNonZero() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Int64Ptr(12),
		},
	})
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeAppendPositionConditionNotMet)
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMaxSizeTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: to.Int64Ptr(int64(len(blockBlobDefaultData) + 1)),
		},
	})
	_assert.Nil(err)
	validateBlockAppended(_assert, abClient, len(blockBlobDefaultData))
}

func (s *azblobTestSuite) TestBlobAppendBlockIfMaxSizeFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: to.Int64Ptr(int64(len(blockBlobDefaultData) - 1)),
		},
	})
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeMaxBlobSizeConditionNotMet)
}

func (s *azblobTestSuite) TestSealAppendBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	appendResp, err := abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_assert.Nil(err)
	_assert.Equal(appendResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendResp.BlobAppendOffset, "0")
	_assert.Equal(*appendResp.BlobCommittedBlockCount, int32(1))

	sealResp, err := abClient.SealAppendBlob(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*sealResp.IsSealed, true)

	appendResp, err = abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, "BlobIsSealed")

	getPropResp, err := abClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getPropResp.IsSealed, true)
}

// TODO: Learn about the behaviour of AppendPosition
// nolint
//func (s *azblobUnrecordedTestSuite) TestSealAppendBlobWithAppendConditions() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	abName := generateBlobName(testName)
//	abClient := createNewAppendBlob(_assert, abName, containerClient)
//
//	sealResp, err := abClient.SealAppendBlob(ctx, &SealAppendBlobOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: to.Int64Ptr(1),
//		},
//	})
//	_assert.NotNil(err)
//	_ = sealResp
//
//	sealResp, err = abClient.SealAppendBlob(ctx, &SealAppendBlobOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: to.Int64Ptr(0),
//		},
//	})
//}

func (s *azblobTestSuite) TestCopySealedBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.SealAppendBlob(ctx, nil)
	_assert.Nil(err)

	copiedBlob1 := getAppendBlobClient("copy1"+abName, containerClient)
	// copy sealed blob will get a sealed blob
	_, err = copiedBlob1.StartCopyFromURL(ctx, abClient.URL(), nil)
	_assert.Nil(err)

	getResp1, err := copiedBlob1.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp1.IsSealed, true)

	_, err = copiedBlob1.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, "BlobIsSealed")

	copiedBlob2 := getAppendBlobClient("copy2"+abName, containerClient)
	_, err = copiedBlob2.StartCopyFromURL(ctx, abClient.URL(), &StartCopyBlobOptions{
		SealBlob: to.BoolPtr(true),
	})
	_assert.Nil(err)

	getResp2, err := copiedBlob2.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp2.IsSealed, true)

	_, err = copiedBlob2.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_assert.NotNil(err)
	validateStorageError(_assert, err, "BlobIsSealed")

	copiedBlob3 := getAppendBlobClient("copy3"+abName, containerClient)
	_, err = copiedBlob3.StartCopyFromURL(ctx, abClient.URL(), &StartCopyBlobOptions{
		SealBlob: to.BoolPtr(false),
	})
	_assert.Nil(err)

	getResp3, err := copiedBlob3.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(getResp3.IsSealed)

	appendResp3, err := copiedBlob3.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	_assert.Nil(err)
	_assert.Equal(appendResp3.RawResponse.StatusCode, 201)
	_assert.Equal(*appendResp3.BlobAppendOffset, "0")
	_assert.Equal(*appendResp3.BlobCommittedBlockCount, int32(1))
}

func (s *azblobTestSuite) TestCopyUnsealedBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	copiedBlob := getAppendBlobClient("copy"+abName, containerClient)
	_, err = copiedBlob.StartCopyFromURL(ctx, abClient.URL(), &StartCopyBlobOptions{
		SealBlob: to.BoolPtr(true),
	})
	_assert.Nil(err)

	getResp, err := copiedBlob.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp.IsSealed, true)
}
