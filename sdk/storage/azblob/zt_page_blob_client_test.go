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
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
	"io"
	"time"
)

func (s *azblobTestSuite) TestPutGetPages() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	contentSize := 1024
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{PageRange: &HttpRange{offset, count}}
	reader, _ := generateData(1024)
	putResp, err := pbClient.UploadPages(context.Background(), reader, &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.NotNil(putResp.LastModified)
	_require.Equal((*putResp.LastModified).IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.Nil(putResp.ContentMD5)
	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
	_require.NotNil(*putResp.RequestID)
	_require.NotNil(*putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: &HttpRange{Offset: 0, Count: 1023}})

	for pager.NextPage(ctx) {
		pageListResp := pager.PageResponse()
		_require.Nil(pager.Err())

		_require.Equal(pageListResp.RawResponse.StatusCode, 200)
		_require.NotNil(pageListResp.LastModified)
		_require.Equal((*pageListResp.LastModified).IsZero(), false)
		_require.NotNil(pageListResp.ETag)
		_require.Equal(*pageListResp.BlobContentLength, int64(512*10))
		_require.NotNil(*pageListResp.RequestID)
		_require.NotNil(*pageListResp.Version)
		_require.NotNil(pageListResp.Date)
		_require.Equal((*pageListResp.Date).IsZero(), false)
		_require.NotNil(pageListResp.PageList)
		pageRangeResp := pageListResp.PageList.PageRange
		_require.Len(pageRangeResp, 1)
		rawStart, rawEnd := (pageRangeResp)[0].Raw()
		_require.Equal(rawStart, offset)
		_require.Equal(rawEnd, count-1)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(_require, "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(_require, "dstblob", containerClient, int64(contentSize))

	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(uploadSrcResp1.RawResponse.StatusCode, 201)
	_require.NotNil(uploadSrcResp1.LastModified)
	_require.Equal((*uploadSrcResp1.LastModified).IsZero(), false)
	_require.NotNil(uploadSrcResp1.ETag)
	_require.Nil(uploadSrcResp1.ContentMD5)
	_require.Equal(*uploadSrcResp1.BlobSequenceNumber, int64(0))
	_require.NotNil(*uploadSrcResp1.RequestID)
	_require.NotNil(*uploadSrcResp1.Version)
	_require.NotNil(uploadSrcResp1.Date)
	_require.Equal((*uploadSrcResp1.Date).IsZero(), false)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		_require.Error(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL.
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
	_require.Nil(err)
	_require.Equal(pResp1.RawResponse.StatusCode, 201)
	_require.NotNil(pResp1.ETag)
	_require.NotNil(pResp1.LastModified)
	_require.NotNil(pResp1.ContentMD5)
	_require.NotNil(pResp1.RequestID)
	_require.NotNil(pResp1.Version)
	_require.NotNil(pResp1.Date)
	_require.Equal((*pResp1.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	_require.Nil(err)
	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{}))
	_require.Nil(err)
	_require.EqualValues(destData, sourceData)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURLWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	contentMD5 := md5Value[:]
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(_require, "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(_require, "dstblob", containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(uploadSrcResp1.RawResponse.StatusCode, 201)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		_require.Error(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL with MD5.
	uploadPagesFromURLOptions := PageBlobUploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
	}
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.Nil(err)
	_require.Equal(pResp1.RawResponse.StatusCode, 201)
	_require.NotNil(pResp1.ETag)
	_require.NotNil(pResp1.LastModified)
	_require.NotNil(pResp1.ContentMD5)
	_require.EqualValues(pResp1.ContentMD5, contentMD5)
	_require.NotNil(pResp1.RequestID)
	_require.NotNil(pResp1.Version)
	_require.NotNil(pResp1.Date)
	_require.Equal((*pResp1.Date).IsZero(), false)
	_require.Equal(*pResp1.BlobSequenceNumber, int64(0))

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	_require.Nil(err)
	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{}))
	_require.Nil(err)
	_require.EqualValues(destData, sourceData)

	// Upload page from URL with bad MD5
	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions = PageBlobUploadPagesFromURLOptions{
		SourceContentMD5: badContentMD5,
	}
	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeMD5Mismatch)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestClearDiffPages() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	contentSize := 2 * 1024
	r := getReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(context.Background(), r, &PageBlobUploadPagesOptions{PageRange: NewHttpRange(int64(0), int64(contentSize))})
	_require.Nil(err)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	r1 := getReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(context.Background(), r1, &PageBlobUploadPagesOptions{PageRange: NewHttpRange(int64(contentSize), int64(contentSize))})
	_require.Nil(err)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{PageRange: &HttpRange{Offset: 0, Count: int64(4096)}, PrevSnapshot: snapshotResp.Snapshot})

	for pager.NextPage(ctx) {
		pageListResp := pager.PageResponse()
		_require.Nil(pager.Err())
		pageRangeResp := pageListResp.PageList.PageRange
		_require.NotNil(pageRangeResp)
		_require.Len(pageRangeResp, 1)
		rawStart, rawEnd := (pageRangeResp)[0].Raw()
		_require.Equal(rawStart, int64(2048))
		_require.Equal(rawEnd, int64(4095))
	}

	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, nil)
	_require.Nil(err)
	_require.Equal(clearResp.RawResponse.StatusCode, 201)

	pager = pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{PageRange: &HttpRange{Offset: 0, Count: int64(4096)}, PrevSnapshot: snapshotResp.Snapshot})

	for pager.NextPage(ctx) {
		pageListResp := pager.PageResponse()
		_require.Nil(pager.Err())
		pageRangeResp := pageListResp.PageList.PageRange
		_require.Len(pageRangeResp, 0)
	}
}

// nolint
func waitForIncrementalCopy(_require *require.Assertions, copyBlobClient *PageBlobClient, blobCopyResponse *PageBlobCopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult BlobGetPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(ctx, nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_require.Fail("")
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

// nolint
func (s *azblobUnrecordedTestSuite) TestIncrementalCopy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	accessType := PublicAccessTypeBlob

	_, err = containerClient.SetAccessPolicy(context.Background(), &ContainerSetAccessPolicyOptions{Access: &accessType})
	_require.Nil(err)

	srcBlob := createNewPageBlob(_require, "src"+generateBlobName(testName), containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = srcBlob.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	dstBlob, _ := containerClient.NewPageBlobClient("dst" + generateBlobName(testName))

	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 202)
	_require.NotNil(resp.LastModified)
	_require.Equal((*resp.LastModified).IsZero(), false)
	_require.NotNil(resp.ETag)
	_require.NotEqual(*resp.RequestID, "")
	_require.NotEqual(*resp.Version, "")
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)
	_require.NotEqual(*resp.CopyID, "")
	_require.Equal(*resp.CopyStatus, CopyStatusTypePending)

	waitForIncrementalCopy(_require, dstBlob, &resp)
}

func (s *azblobTestSuite) TestResizePageBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	resp, err := pbClient.Resize(context.Background(), 2048, nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 200)

	resp, err = pbClient.Resize(context.Background(), 8192, nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 200)

	resp2, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(8192))
}

func (s *azblobTestSuite) TestPageSequenceNumbers() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(0)
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(7)
	actionType = SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(11)
	actionType = SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob = PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 200)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestPutPagesWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	// put page with valid MD5
	contentSize := 1024
	readerToBody, body := getRandomDataAndReader(contentSize)
	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]

	putResp, err := pbClient.UploadPages(context.Background(), internal.NopCloser(readerToBody), &PageBlobUploadPagesOptions{
		PageRange:               NewHttpRange(offset, count),
		TransactionalContentMD5: contentMD5,
	})
	_require.Nil(err)
	_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.NotNil(putResp.LastModified)
	_require.Equal((*putResp.LastModified).IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.NotNil(putResp.ContentMD5)
	_require.EqualValues(putResp.ContentMD5, contentMD5)
	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
	_require.NotNil(*putResp.RequestID)
	_require.NotNil(*putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	// put page with bad MD5
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	basContentMD5 := badMD5[:]
	_ = body
	putResp, err = pbClient.UploadPages(context.Background(), internal.NopCloser(readerToBody), &PageBlobUploadPagesOptions{
		PageRange:               NewHttpRange(offset, count),
		TransactionalContentMD5: basContentMD5,
	})
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobCreatePageSizeInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, 1, &createPageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobCreatePageSequenceInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{"In valid1": "bar"},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.NotNil(err)
	_require.Contains(err.Error(), invalidHeaderErrorSubstring)

}

func (s *azblobTestSuite) TestBlobCreatePageHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		HTTPHeaders:        &basicHeaders,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	h := resp.GetHTTPHeaders()
	_require.EqualValues(h, basicHeaders)
}

func validatePageBlobPut(_require *require.Assertions, pbClient *PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, basicMetadata)
	_require.EqualValues(resp.GetHTTPHeaders(), basicHeaders)
}

func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	sequenceNumber := int64(0)
	createPageBlobOptions := PageBlobCreateOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
		HTTPHeaders:        &basicHeaders,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobPutPagesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize/2)
	uploadPagesOptions := PageBlobUploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)
}

//// Body cannot be nil check already added in the request preparer
////func (s *azblobTestSuite) TestBlobPutPagesNilBody() {
////  svcClient := getServiceClient()
////  containerClient, _ := createNewContainer(c, svcClient)
////  defer deleteContainer(_require, containerClient)
////  pbClient, _ := createNewPageBlob(c, containerClient)
////
////  _, err := pbClient.UploadPages(ctx, nil, nil)
////  _require.NotNil(err)
////}

func (s *azblobTestSuite) TestBlobPutPagesEmptyBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	r := bytes.NewReader([]byte{})
	offset, count := int64(0), int64(0)
	uploadPagesOptions := PageBlobUploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobPutPagesNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := PageBlobUploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeBlobNotFound)
}

func validateUploadPages(_require *require.Assertions, pbClient *PageBlobClient) {
	// This will only validate a single put page at 0-PageBlobPageBytes-1
	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{
		PageRange: &HttpRange{0, CountToEnd},
	})

	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp := pager.PageResponse()
		pageListResp := resp.PageList.PageRange
		start, end := int64(0), int64(PageBlobPageBytes-1)
		rawStart, rawEnd := pageListResp[0].Raw()
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
	}

}

func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	_, err = pbClient.UploadPages(ctx, r, &PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	_, err = pbClient.UploadPages(ctx, r, &PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	_, err = pbClient.UploadPages(ctx, r, &PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	_, err = pbClient.UploadPages(ctx, r, &PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	_, err = pbClient.UploadPages(ctx, r, &PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	eTag := "garbage"
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	eTag := "garbage"
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThan := int64(10)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThan := int64(1)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}

	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTETrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTENegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne() {
//    _require := require.New(s.T())
//    testName := s.T().Name()
//    _context := getTestContext(testName)
//    svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//    if err != nil {
//        _require.Fail("Unable to fetch service client because " + err.Error())
//    }
//
//    containerName := generateContainerName(testName)
//    containerClient := createNewContainer(_require, containerName, svcClient)
//    defer deleteContainer(_require, containerClient)
//
//    blobName := generateBlobName(testName)
//    pbClient := createNewPageBlob(_require, blobName, containerClient)
//
//    r, _ := generateData(PageBlobPageBytes)
//    offset, count := int64(0), int64(PageBlobPageBytes)
//    ifSequenceNumberEqualTo := int64(-1)
//    uploadPagesOptions := PageBlobUploadPagesOptions{
//        PageRange: NewHttpRange(offset, count),
//        SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//            IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//        },
//    }
//    _, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions) // This will cause the library to set the value of the header to 0
//    _require.Nil(err)
//}

func setupClearPagesTest(_require *require.Assertions, testName string) (*ContainerClient, *PageBlobClient) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	return containerClient, pbClient
}

func validateClearPagesTest(_require *require.Assertions, pbClient *PageBlobClient) {
	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0)})
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
		pageListResp := pager.PageResponse().PageRange
		_require.Nil(pageListResp)
	}

}

func (s *azblobTestSuite) TestBlobClearPagesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes + 1}, nil)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		}})
	_require.Nil(err)
	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	clearPageOptions := PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: getPropertiesResp.ETag,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	eTag := "garbage"
	clearPageOptions := PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	eTag := "garbage"
	clearPageOptions := PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	clearPageOptions := PageBlobClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTETrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTEFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTENegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := PageBlobClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidInput)
}

func setupGetPageRangesTest(_require *require.Assertions, testName string) (containerClient *ContainerClient, pbClient *PageBlobClient) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient = createNewContainer(_require, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient = createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)
	return
}

func validateBasicGetPageRanges(_require *require.Assertions, resp PageList, err error) {
	_require.Nil(err)
	_require.NotNil(resp.PageRange)
	_require.Len(resp.PageRange, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := (resp.PageRange)[0].Raw()
	_require.Equal(rawStart, start)
	_require.Equal(rawEnd, end)
}

func (s *azblobTestSuite) TestBlobGetPageRangesEmptyBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0)})
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
		_require.Nil(pager.PageResponse().PageList.PageRange)
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesEmptyRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0)})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp := pager.PageResponse()
		validateBasicGetPageRanges(_require, resp.PageList, err)
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(-2, 500)})
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesNonContiguousRanges() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(2*PageBlobPageBytes), int64(PageBlobPageBytes)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: NewHttpRange(offset, count),
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0)})
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
		resp := pager.PageResponse()
		pageListResp := resp.PageList.PageRange
		_require.NotNil(pageListResp)
		_require.Len(pageListResp, 2)

		start, end := int64(0), int64(PageBlobPageBytes-1)
		rawStart, rawEnd := pageListResp[0].Raw()
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)

		start, end = int64(PageBlobPageBytes*2), int64((PageBlobPageBytes*3)-1)
		rawStart, rawEnd = pageListResp[1].Raw()
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesNotPageAligned() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 2000)})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp := pager.PageResponse()
		validateBasicGetPageRanges(_require, resp.PageList, err)
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	snapshotURL, _ := pbClient.WithSnapshot(*resp.Snapshot)
	pager := snapshotURL.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0)})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp2 := pager.PageResponse()
		validateBasicGetPageRanges(_require, resp2.PageList, err)
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
		validateBasicGetPageRanges(_require, pager.PageResponse().PageList, pager.Err())
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.NextPage(ctx) {
		_require.Nil(err)
	}

	//serr := err.(StorageError)
	//_require.(serr.RawResponse.StatusCode, chk.Equals, 304)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp := pager.PageResponse()
		validateBasicGetPageRanges(_require, resp.PageList, err)
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.NotNil(err)

		validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp2 := pager.PageResponse()
		validateBasicGetPageRanges(_require, resp2.PageList, err)
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: to.Ptr("garbage"),
		},
	}})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.NotNil(err)
		validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: to.Ptr("garbage"),
		},
	}})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp := pager.PageResponse()
		validateBasicGetPageRanges(_require, resp.PageList, err)
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, 0), BlobAccessConditions: &BlobAccessConditions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}})
	for pager.NextPage(ctx) {
		_require.NotNil(pager.Err())
	}

	//serr := err.(StorageError)
	//_require.(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
}

// nolint
func setupDiffPageRangesTest(_require *require.Assertions, testName string) (containerClient *ContainerClient, pbClient *PageBlobClient, snapshot string) {
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient = createNewContainer(_require, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient = createNewPageBlob(_require, blobName, containerClient)

	r := getReaderToGeneratedBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: &HttpRange{Offset: offset, Count: count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	snapshot = *resp.Snapshot

	r = getReaderToGeneratedBytes(PageBlobPageBytes)
	offset, count = int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions = PageBlobUploadPagesOptions{
		PageRange: &HttpRange{Offset: offset, Count: count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)
	return
}

// nolint
func validateDiffPageRanges(_require *require.Assertions, resp PageList, err error) {
	_require.Nil(err)
	pageListResp := resp.PageRange
	_require.NotNil(pageListResp)
	_require.Len(resp.PageRange, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	_require.EqualValues(rawStart, start)
	_require.EqualValues(rawEnd, end)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangesNonExistentSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	snapshotTime, _ := time.Parse(SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange:    NewHttpRange(0, 0),
		PrevSnapshot: to.Ptr(snapshotTime.Format(SnapshotTimeFormat))})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.NotNil(err)
		validateStorageError(_require, err, StorageErrorCodePreviousSnapshotNotFound)
	}

}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)
	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{PageRange: NewHttpRange(-22, 14), Snapshot: &snapshot})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange: NewHttpRange(0, 0),
		Snapshot:  to.Ptr(snapshot),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	for pager.NextPage(ctx) {
		err := pager.Err()
		resp := pager.PageResponse()
		validateDiffPageRanges(_require, resp.PageList, err)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(10)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange: NewHttpRange(0, 0),
		Snapshot:  to.Ptr(snapshot),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp := pager.PageResponse()
		validateDiffPageRanges(_require, resp.PageList, err)
	}

}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(10)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange: NewHttpRange(0, 0),
		Snapshot:  to.Ptr(snapshot),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp := pager.PageResponse()
		validateDiffPageRanges(_require, resp.PageList, err)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange: NewHttpRange(0, 0),
		Snapshot:  to.Ptr(snapshot),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		_require.NotNil(err)
		validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
	}

}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange: NewHttpRange(0, 0),
		Snapshot:  to.Ptr(snapshot),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	for pager.NextPage(ctx) {
		err := pager.Err()
		_require.Nil(err)
		resp2 := pager.PageResponse()
		validateDiffPageRanges(_require, resp2.PageList, err)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange: NewHttpRange(0, 0),
		Snapshot:  to.Ptr(snapshotStr),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: to.Ptr("garbage"),
			},
		}})

	for pager.NextPage(ctx) {
		_require.NotNil(pager.Err())
		validateStorageError(_require, pager.Err(), StorageErrorCodeConditionNotMet)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange:    NewHttpRange(0, 0),
		PrevSnapshot: to.Ptr(snapshotStr),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: to.Ptr("garbage"),
			},
		}})

	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
		resp := pager.PageResponse()
		validateDiffPageRanges(_require, resp.PageList, pager.Err())
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	pager := pbClient.GetPageRangesDiff(&PageBlobGetPageRangesDiffOptions{
		PageRange:    NewHttpRange(0, 0),
		PrevSnapshot: to.Ptr(snapshot),
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	})

	for pager.NextPage(ctx) {
		_require.NotNil(pager.Err())
	}
}

func (s *azblobTestSuite) TestBlobResizeZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	// The default pbClient is created with size > 0, so this should actually update
	_, err = pbClient.Resize(ctx, 0, nil)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(0))
}

func (s *azblobTestSuite) TestBlobResizeInvalidSizeNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	_, err = pbClient.Resize(ctx, -4, nil)
	_require.NotNil(err)
}

func (s *azblobTestSuite) TestBlobResizeInvalidSizeMisaligned() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	_, err = pbClient.Resize(ctx, 12, nil)
	_require.NotNil(err)
}

func validateResize(_require *require.Assertions, pbClient *PageBlobClient) {
	resp, _ := pbClient.GetProperties(ctx, nil)
	_require.Equal(*resp.ContentLength, int64(PageBlobPageBytes))
}

func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := PageBlobResizeOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberActionTypeInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberSequenceNumberInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	defer func() { // Invalid sequence number should panic
		_ = recover()
	}()

	sequenceNumber := int64(-1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidHeaderValue)
}

func validateSequenceNumberSet(_require *require.Assertions, pbClient *PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.BlobSequenceNumber, int64(1))
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient, _ := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_require.Nil(err)
	_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	eTag := "garbage"
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, "src"+blobName, containerClient)

	eTag := "garbage"
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, "src"+blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := PageBlobUpdateSequenceNumberOptions{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

////func setupStartIncrementalCopyTest(_require *require.Assertions, testName string) (containerClient containerClient,
////    pbClient PageBlobClient, copyPBClient PageBlobClient, snapshot string) {
////    _context := getTestContext(testName)
////    var recording *testframework.Recording
////    if _context != nil {
////        recording = _context.recording
////    }
////    svcClient, err := getServiceClient(recording, testAccountDefault, nil)
////    if err != nil {
////        _require.Fail("Unable to fetch service client because " + err.Error())
////    }
////
////    containerName := generateContainerName(testName)
////    containerClient = createNewContainer(_require, containerName, svcClient)
////    defer deleteContainer(_require, containerClient)
////
////    accessType := PublicAccessTypeBlob
////    setAccessPolicyOptions := ContainerSetAccessPolicyOptions{
////        ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
////    }
////    _, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
////    _require.Nil(err)
////
////    pbClient = createNewPageBlob(_require, generateBlobName(testName), containerClient)
////    resp, _ := pbClient.CreateSnapshot(ctx, nil)
////
////    copyPBClient = getPageBlobClient("copy"+generateBlobName(testName), containerClient)
////
////    // Must create the incremental copy pbClient so that the access conditions work on it
////    resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), *resp.Snapshot, nil)
////    _require.Nil(err)
////    waitForIncrementalCopy(_require, copyPBClient, &resp2)
////
////    resp, _ = pbClient.CreateSnapshot(ctx, nil) // Take a new snapshot so the next copy will succeed
////    snapshot = *resp.Snapshot
////    return
////}
//
////func validateIncrementalCopy(_require *require.Assertions, copyPBClient PageBlobClient, resp *PageBlobCopyIncrementalResponse) {
////    t := waitForIncrementalCopy(_require, copyPBClient, resp)
////
////    // If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
////    copySnapshotURL := copyPBClient.WithSnapshot(*t)
////    _, err := copySnapshotURL.GetProperties(ctx, nil)
////    _require.Nil(err)
////}
//
////func (s *azblobTestSuite) TestBlobStartIncrementalCopySnapshotNotExist() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    _context := getTestContext(testName)
////    svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
////    if err != nil {
////        _require.Fail("Unable to fetch service client because " + err.Error())
////    }
////
////    containerName := generateContainerName(testName)
////    containerClient := createNewContainer(_require, containerName, svcClient)
////    defer deleteContainer(_require, containerClient)
////
////    blobName := generateBlobName(testName)
////    pbClient := createNewPageBlob(_require, "src" + blobName, containerClient)
////    copyPBClient := getPageBlobClient("dst" + blobName, containerClient)
////
////    snapshot := time.Now().UTC().Format(SnapshotTimeFormat)
////    _, err = copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, nil)
////    _require.NotNil(err)
////
////    validateStorageError(_require, err, StorageErrorCodeCannotVerifyCopySource)
////}
//
////func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////
////    defer deleteContainer(_require, containerClient)
////
////    currentTime := getRelativeTimeGMT(-20)
////
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfModifiedSince: &currentTime,
////        },
////    }
////    resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.Nil(err)
////
////    validateIncrementalCopy(_require, copyPBClient, &resp)
////}
//
////func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////
////    defer deleteContainer(_require, containerClient)
////
////    currentTime := getRelativeTimeGMT(20)
////
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfModifiedSince: &currentTime,
////        },
////    }
////    _, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.NotNil(err)
////
////    validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
////}
//
////func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////
////    defer deleteContainer(_require, containerClient)
////
////    currentTime := getRelativeTimeGMT(20)
////
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfUnmodifiedSince: &currentTime,
////        },
////    }
////    resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.Nil(err)
////
////    validateIncrementalCopy(_require, copyPBClient, &resp)
////}
//
////func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////
////    defer deleteContainer(_require, containerClient)
////
////    currentTime := getRelativeTimeGMT(-20)
////
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfUnmodifiedSince: &currentTime,
////        },
////    }
////    _, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.NotNil(err)
////
////    validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
////}
//
////nolint
////func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchTrue() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////    resp, _ := copyPBClient.GetProperties(ctx, nil)
////
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfMatch: resp.ETag,
////        },
////    }
////    resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.Nil(err)
////
////    validateIncrementalCopy(_require, copyPBClient, &resp2)
////    defer deleteContainer(_require, containerClient)
////}
////
//
////nolint
////func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchFalse() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////
////    defer deleteContainer(_require, containerClient)
////
////    eTag := "garbage"
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfMatch: &eTag,
////        },
////    }
////    _, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.NotNil(err)
////
////    validateStorageError(_require, err, StorageErrorCodeTargetConditionNotMet)
////}
////
//
////nolint
////func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////    defer deleteContainer(_require, containerClient)
////
////    eTag := "garbage"
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfNoneMatch: &eTag,
////        },
////    }
////    resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.Nil(err)
////
////    validateIncrementalCopy(_require, copyPBClient, &resp)
////}
////
//
////nolint
////func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse() {
////    _require := require.New(s.T())
////    testName := s.T().Name()
////    containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
////    defer deleteContainer(_require, containerClient)
////
////    resp, _ := copyPBClient.GetProperties(ctx, nil)
////
////    copyIncrementalPageBlobOptions := PageBlobCopyIncrementalOptions{
////        ModifiedAccessConditions: &ModifiedAccessConditions{
////            IfNoneMatch: resp.ETag,
////        },
////    }
////    _, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
////    _require.NotNil(err)
////
////    validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
////}
