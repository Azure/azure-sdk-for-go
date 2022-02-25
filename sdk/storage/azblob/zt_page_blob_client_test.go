// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"io/ioutil"
	"testing"
	"time"

	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
)

func (s *azblobTestSuite) TestPutGetPages() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	contentSize := 1024
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	reader, _ := generateData(1024)
	putResp, err := pbClient.UploadPages(context.Background(), reader, &uploadPagesOptions)
	require.NoError(s.T(), err)
	require.Equal(s.T(), putResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), putResp.LastModified)
	require.Equal(s.T(), (*putResp.LastModified).IsZero(), false)
	require.NotNil(s.T(), putResp.ETag)
	require.Nil(s.T(), putResp.ContentMD5)
	require.Equal(s.T(), *putResp.BlobSequenceNumber, int64(0))
	require.NotNil(s.T(), *putResp.RequestID)
	require.NotNil(s.T(), *putResp.Version)
	require.NotNil(s.T(), putResp.Date)
	require.Equal(s.T(), (*putResp.Date).IsZero(), false)

	pageList, err := pbClient.GetPageRanges(context.Background(), HttpRange{0, 1023}, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageList.RawResponse.StatusCode, 200)
	require.NotNil(s.T(), pageList.LastModified)
	require.Equal(s.T(), (*pageList.LastModified).IsZero(), false)
	require.NotNil(s.T(), pageList.ETag)
	require.Equal(s.T(), *pageList.BlobContentLength, int64(512*10))
	require.NotNil(s.T(), *pageList.RequestID)
	require.NotNil(s.T(), *pageList.Version)
	require.NotNil(s.T(), pageList.Date)
	require.Equal(s.T(), (*pageList.Date).IsZero(), false)
	require.NotNil(s.T(), pageList.PageList)
	pageRangeResp := pageList.PageList.PageRange
	require.Len(s.T(), pageRangeResp, 1)
	rawStart, rawEnd := (pageRangeResp)[0].Raw()
	require.Equal(s.T(), rawStart, offset)
	require.Equal(s.T(), rawEnd, count-1)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURL() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(s.T(), "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(s.T(), "dstblob", containerClient, int64(contentSize))

	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uploadSrcResp1.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), uploadSrcResp1.LastModified)
	require.Equal(s.T(), (*uploadSrcResp1.LastModified).IsZero(), false)
	require.NotNil(s.T(), uploadSrcResp1.ETag)
	require.Nil(s.T(), uploadSrcResp1.ContentMD5)
	require.Equal(s.T(), *uploadSrcResp1.BlobSequenceNumber, int64(0))
	require.NotNil(s.T(), *uploadSrcResp1.RequestID)
	require.NotNil(s.T(), *uploadSrcResp1.Version)
	require.NotNil(s.T(), uploadSrcResp1.Date)
	require.Equal(s.T(), (*uploadSrcResp1.Date).IsZero(), false)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getGenericCredential(nil, testAccountDefault)
	require.NoError(s.T(), err)
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		require.Error(s.T(), err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL.
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pResp1.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pResp1.ETag)
	require.NotNil(s.T(), pResp1.LastModified)
	require.NotNil(s.T(), pResp1.ContentMD5)
	require.NotNil(s.T(), pResp1.RequestID)
	require.NotNil(s.T(), pResp1.Version)
	require.NotNil(s.T(), pResp1.Date)
	require.Equal(s.T(), (*pResp1.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	require.NoError(s.T(), err)
	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{}))
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), destData, sourceData)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURLWithMD5() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	contentMD5 := md5Value[:]
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(s.T(), "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(s.T(), "dstblob", containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uploadSrcResp1.RawResponse.StatusCode, 201)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getGenericCredential(nil, testAccountDefault)
	require.NoError(s.T(), err)
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		require.Error(s.T(), err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL with MD5.
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
	}
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pResp1.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pResp1.ETag)
	require.NotNil(s.T(), pResp1.LastModified)
	require.NotNil(s.T(), pResp1.ContentMD5)
	require.EqualValues(s.T(), pResp1.ContentMD5, contentMD5)
	require.NotNil(s.T(), pResp1.RequestID)
	require.NotNil(s.T(), pResp1.Version)
	require.NotNil(s.T(), pResp1.Date)
	require.Equal(s.T(), (*pResp1.Date).IsZero(), false)
	require.Equal(s.T(), *pResp1.BlobSequenceNumber, int64(0))

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	require.NoError(s.T(), err)
	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{}))
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), destData, sourceData)

	// Upload page from URL with bad MD5
	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions = UploadPagesFromURLOptions{
		SourceContentMD5: badContentMD5,
	}
	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeMD5Mismatch)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestClearDiffPages() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	contentSize := 2 * 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	require.NoError(s.T(), err)

	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}}
	_, err = pbClient.UploadPages(context.Background(), getReaderToGeneratedBytes(2048), &uploadPagesOptions1)
	require.NoError(s.T(), err)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
	require.NoError(s.T(), err)
	pageRangeResp := pageListResp.PageList.PageRange
	require.NotNil(s.T(), pageRangeResp)
	require.Len(s.T(), pageRangeResp, 1)
	// require.(s.T(), (pageRangeResp)[0], chk.DeepEquals, PageRange{Start: &offset1, End: &end1})
	rawStart, rawEnd := (pageRangeResp)[0].Raw()
	require.Equal(s.T(), rawStart, offset1)
	require.Equal(s.T(), rawEnd, end1)

	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), clearResp.RawResponse.StatusCode, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
	require.NoError(s.T(), err)
	require.Nil(s.T(), pageListResp.PageList.PageRange)
}

//nolint
func waitForIncrementalCopy(t *testing.T, copyBlobClient PageBlobClient, blobCopyResponse *PageBlobCopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult GetBlobPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(ctx, nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			require.Fail(t, "")
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

//nolint
func (s *azblobUnrecordedTestSuite) TestIncrementalCopy() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	accessType := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	require.NoError(s.T(), err)

	srcBlob := createNewPageBlob(s.T(), "src"+generateBlobName(testName), containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = srcBlob.UploadPages(context.Background(), r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
	require.NoError(s.T(), err)

	dstBlob := containerClient.NewPageBlobClient("dst" + generateBlobName(testName))

	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 202)
	require.NotNil(s.T(), resp.LastModified)
	require.Equal(s.T(), (*resp.LastModified).IsZero(), false)
	require.NotNil(s.T(), resp.ETag)
	require.NotEqual(s.T(), *resp.RequestID, "")
	require.NotEqual(s.T(), *resp.Version, "")
	require.NotNil(s.T(), resp.Date)
	require.Equal(s.T(), (*resp.Date).IsZero(), false)
	require.NotEqual(s.T(), *resp.CopyID, "")
	require.Equal(s.T(), *resp.CopyStatus, CopyStatusTypePending)

	waitForIncrementalCopy(s.T(), dstBlob, &resp)
}

func (s *azblobTestSuite) TestResizePageBlob() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	resp, err := pbClient.Resize(context.Background(), 2048, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 200)

	resp, err = pbClient.Resize(context.Background(), 8192, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 200)

	resp2, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *resp2.ContentLength, int64(8192))
}

func (s *azblobTestSuite) TestPageSequenceNumbers() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(0)
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(7)
	actionType = SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(11)
	actionType = SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 200)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPutPagesWithMD5() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	// put page with valid MD5
	contentSize := 1024
	readerToBody, body := getRandomDataAndReader(contentSize)
	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]
	uploadPagesOptions := UploadPagesOptions{
		PageRange:               &HttpRange{offset, count},
		TransactionalContentMD5: contentMD5,
	}

	putResp, err := pbClient.UploadPages(context.Background(), internal.NopCloser(readerToBody), &uploadPagesOptions)
	require.NoError(s.T(), err)
	require.Equal(s.T(), putResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), putResp.LastModified)
	require.Equal(s.T(), (*putResp.LastModified).IsZero(), false)
	require.NotNil(s.T(), putResp.ETag)
	require.NotNil(s.T(), putResp.ContentMD5)
	require.EqualValues(s.T(), putResp.ContentMD5, contentMD5)
	require.Equal(s.T(), *putResp.BlobSequenceNumber, int64(0))
	require.NotNil(s.T(), *putResp.RequestID)
	require.NotNil(s.T(), *putResp.Version)
	require.NotNil(s.T(), putResp.Date)
	require.Equal(s.T(), (*putResp.Date).IsZero(), false)

	// put page with bad MD5
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	basContentMD5 := badMD5[:]
	_ = body
	uploadPagesOptions = UploadPagesOptions{
		PageRange:               &HttpRange{offset, count},
		TransactionalContentMD5: basContentMD5,
	}
	putResp, err = pbClient.UploadPages(context.Background(), internal.NopCloser(readerToBody), &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobCreatePageSizeInvalid() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, 1, &createPageBlobOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobCreatePageSequenceInvalid() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.Error(s.T(), err)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataNonEmpty() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.NoError(s.T(), err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp.Metadata)
	require.EqualValues(s.T(), resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataEmpty() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.NoError(s.T(), err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Nil(s.T(), resp.Metadata)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataInvalid() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{"In valid1": "bar"},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), invalidHeaderErrorSubstring)

}

func (s *azblobTestSuite) TestBlobCreatePageHTTPHeaders() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		HTTPHeaders:        &basicHeaders,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.NoError(s.T(), err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	h := resp.GetHTTPHeaders()
	require.EqualValues(s.T(), h, basicHeaders)
}

func validatePageBlobPut(t *testing.T, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Metadata)
	require.EqualValues(t, resp.Metadata, basicMetadata)
	require.EqualValues(t, resp.GetHTTPHeaders(), basicHeaders)
}

func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.NoError(s.T(), err)

	validatePageBlobPut(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.NoError(s.T(), err)

	validatePageBlobPut(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.NoError(s.T(), err)

	validatePageBlobPut(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.NoError(s.T(), err)

	validatePageBlobPut(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
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
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobPutPagesInvalidRange() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize/2)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)
}

//// Body cannot be nil check already added in the request preparer
////func (s *azblobTestSuite) TestBlobPutPagesNilBody() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(s.T(), containerClient)
////	pbClient, _ := createNewPageBlob(c, containerClient)
////
////	_, err := pbClient.UploadPages(ctx, nil, nil)
////	require.Error(s.T(), err)
////}

func (s *azblobTestSuite) TestBlobPutPagesEmptyBody() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	r := bytes.NewReader([]byte{})
	offset, count := int64(0), int64(0)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.Error(s.T(), err)
}

func (s *azblobTestSuite) TestBlobPutPagesNonExistentBlob() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeBlobNotFound)
}

func validateUploadPages(t *testing.T, pbClient PageBlobClient) {
	// This will only validate a single put page at 0-PageBlobPageBytes-1
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	require.NoError(t, err)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	require.Equal(t, rawStart, start)
	require.Equal(t, rawEnd, end)
}

func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	validateUploadPages(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	validateUploadPages(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	validateUploadPages(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	eTag := "garbage"
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	eTag := "garbage"
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	validateUploadPages(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThan := int64(10)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	validateUploadPages(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThan := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}

	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTETrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	validateUploadPages(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTENegOne() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	validateUploadPages(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeSequenceNumberConditionNotMet)
}

//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne() {
//
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		require.Fail(s.T(), "Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(s.T(), containerName, svcClient)
//	defer deleteContainer(s.T(), containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := createNewPageBlob(s.T(), blobName, containerClient)
//
//	r, _ := generateData(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberEqualTo := int64(-1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions) // This will cause the library to set the value of the header to 0
//	require.NoError(s.T(), err)
//}

func setupClearPagesTest(t *testing.T, testName string) (ContainerClient, PageBlobClient) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(t, err)

	return containerClient, pbClient
}

func validateClearPagesTest(t *testing.T, pbClient PageBlobClient) {
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(t, err)
	pageListResp := resp.PageList.PageRange
	require.Nil(t, pageListResp)
}

func (s *azblobTestSuite) TestBlobClearPagesInvalidRange() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes + 1}, nil)
	require.Error(s.T(), err)
}

func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(s.T(), err)

	validateClearPagesTest(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(s.T(), err)

	validateClearPagesTest(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: getPropertiesResp.ETag,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(s.T(), err)

	validateClearPagesTest(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(s.T(), err)

	validateClearPagesTest(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(s.T(), err)

	validateClearPagesTest(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTETrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(s.T(), err)

	validateClearPagesTest(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTEFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTENegOne() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(s.T(), err)

	validateClearPagesTest(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne() {

	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeInvalidInput)
}

func setupGetPageRangesTest(t *testing.T, testName string) (containerClient ContainerClient, pbClient PageBlobClient) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient = createNewContainer(t, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient = createNewPageBlob(t, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(t, err)
	return
}

func validateBasicGetPageRanges(t *testing.T, resp PageList, err error) {
	require.NoError(t, err)
	require.NotNil(t, resp.PageRange)
	require.Len(t, resp.PageRange, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := (resp.PageRange)[0].Raw()
	require.Equal(t, rawStart, start)
	require.Equal(t, rawEnd, end)
}

func (s *azblobTestSuite) TestBlobGetPageRangesEmptyBlob() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(s.T(), err)
	require.Nil(s.T(), resp.PageList.PageRange)
}

func (s *azblobTestSuite) TestBlobGetPageRangesEmptyRange() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(s.T(), err)
	validateBasicGetPageRanges(s.T(), resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesInvalidRange() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	_, err := pbClient.GetPageRanges(ctx, HttpRange{-2, 500}, nil)
	require.NoError(s.T(), err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesNonContiguousRanges() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(2*PageBlobPageBytes), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(s.T(), err)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(s.T(), err)
	pageListResp := resp.PageList.PageRange
	require.NotNil(s.T(), pageListResp)
	require.Len(s.T(), pageListResp, 2)

	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	require.Equal(s.T(), rawStart, start)
	require.Equal(s.T(), rawEnd, end)

	start, end = int64(PageBlobPageBytes*2), int64((PageBlobPageBytes*3)-1)
	rawStart, rawEnd = pageListResp[1].Raw()
	require.Equal(s.T(), rawStart, start)
	require.Equal(s.T(), rawEnd, end)
}

func (s *azblobTestSuite) TestBlobGetPageRangesNotPageAligned() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 2000}, nil)
	require.NoError(s.T(), err)
	validateBasicGetPageRanges(s.T(), resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesSnapshot() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp.Snapshot)

	snapshotURL := pbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(s.T(), err)
	validateBasicGetPageRanges(s.T(), resp2.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateBasicGetPageRanges(s.T(), resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(s.T(), err)

	//serr := err.(StorageError)
	//require.(s.T(), serr.RawResponse.StatusCode, chk.Equals, 304)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateBasicGetPageRanges(s.T(), resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateBasicGetPageRanges(s.T(), resp2.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchTrue() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateBasicGetPageRanges(s.T(), resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchFalse() {

	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(s.T(), err)
	//serr := err.(StorageError)
	//require.(s.T(), serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
}

//nolint
func setupDiffPageRangesTest(t *testing.T, testName string) (containerClient ContainerClient,
	pbClient PageBlobClient, snapshot string) {
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient = createNewContainer(t, containerName, svcClient)

	blobName := generateName(testName)
	pbClient = createNewPageBlob(t, blobName, containerClient)

	r := getReaderToGeneratedBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(t, err)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	require.NoError(t, err)
	snapshot = *resp.Snapshot

	r = getReaderToGeneratedBytes(PageBlobPageBytes)
	offset, count = int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions = UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(t, err)
	return
}

//nolint
func validateDiffPageRanges(t *testing.T, resp PageList, err error) {
	require.NoError(t, err)
	pageListResp := resp.PageRange
	require.NotNil(t, pageListResp)
	require.Len(t, resp.PageRange, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	require.EqualValues(t, rawStart, start)
	require.EqualValues(t, rawEnd, end)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangesNonExistentSnapshot() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	snapshotTime, _ := time.Parse(SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshotTime.Format(SnapshotTimeFormat), nil)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodePreviousSnapshotNotFound)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeInvalidRange() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{-22, 14}, snapshot, nil)
	require.NoError(s.T(), err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceTrue() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateDiffPageRanges(s.T(), resp.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceFalse() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(s.T(), err)

	//stgErr := err.(StorageError)
	//require.(s.T(), stgErr.Response().StatusCode, chk.Equals, 304)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateDiffPageRanges(s.T(), resp.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchTrue() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateDiffPageRanges(s.T(), resp2.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchFalse() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchTrue() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(s.T(), err)
	validateDiffPageRanges(s.T(), resp.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchFalse() {

	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), testName)
	defer deleteContainer(s.T(), containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(s.T(), err)
}

func (s *azblobTestSuite) TestBlobResizeZero() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	// The default pbClient is created with size > 0, so this should actually update
	_, err = pbClient.Resize(ctx, 0, nil)
	require.NoError(s.T(), err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *resp.ContentLength, int64(0))
}

func (s *azblobTestSuite) TestBlobResizeInvalidSizeNegative() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	_, err = pbClient.Resize(ctx, -4, nil)
	require.Error(s.T(), err)
}

func (s *azblobTestSuite) TestBlobResizeInvalidSizeMisaligned() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	_, err = pbClient.Resize(ctx, 12, nil)
	require.Error(s.T(), err)
}

func validateResize(t *testing.T, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentLength, int64(PageBlobPageBytes))
}

func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(s.T(), err)

	validateResize(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(s.T(), err)

	validateResize(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(s.T(), err)

	validateResize(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfNoneMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(s.T(), err)

	validateResize(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfNoneMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberActionTypeInvalid() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberSequenceNumberInvalid() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	defer func() { // Invalid sequence number should panic
		_ = recover()
	}()

	sequenceNumber := int64(-1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeInvalidHeaderValue)
}

func validateSequenceNumberSet(t *testing.T, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.BlobSequenceNumber, int64(1))
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	validateSequenceNumberSet(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	validateSequenceNumberSet(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	validateSequenceNumberSet(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), blobName, containerClient)

	eTag := "garbage"
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), "src"+blobName, containerClient)

	eTag := "garbage"
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(s.T(), err)

	validateSequenceNumberSet(s.T(), pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		require.Fail(s.T(), "Unable to fetch service client because "+err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(s.T(), "src"+blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

//func setupStartIncrementalCopyTest(t *testing.T, testName string) (containerClient ContainerClient,
//	pbClient PageBlobClient, copyPBClient PageBlobClient, snapshot string) {
//	_context := getTestContext(testName)
//	var recording *testframework.Recording
//	if _context != nil {
//		recording = _context.recording
//	}
//	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
//	if err != nil {
//		require.Fail(s.T(), "Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient = createNewContainer(s.T(), containerName, svcClient)
//	defer deleteContainer(s.T(), containerClient)
//
//	accessType := PublicAccessTypeBlob
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
//	}
//	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	require.NoError(s.T(), err)
//
//	pbClient = createNewPageBlob(s.T(), generateBlobName(testName), containerClient)
//	resp, _ := pbClient.CreateSnapshot(ctx, nil)
//
//	copyPBClient = getPageBlobClient("copy"+generateBlobName(testName), containerClient)
//
//	// Must create the incremental copy pbClient so that the access conditions work on it
//	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), *resp.Snapshot, nil)
//	require.NoError(s.T(), err)
//	waitForIncrementalCopy(s.T(), copyPBClient, &resp2)
//
//	resp, _ = pbClient.CreateSnapshot(ctx, nil) // Take a new snapshot so the next copy will succeed
//	snapshot = *resp.Snapshot
//	return
//}

//func validateIncrementalCopy(t *testing.T, copyPBClient PageBlobClient, resp *PageBlobCopyIncrementalResponse) {
//	t := waitForIncrementalCopy(s.T(), copyPBClient, resp)
//
//	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
//	copySnapshotURL := copyPBClient.WithSnapshot(*t)
//	_, err := copySnapshotURL.GetProperties(ctx, nil)
//	require.NoError(s.T(), err)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopySnapshotNotExist() {
//
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		require.Fail(s.T(), "Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(s.T(), containerName, svcClient)
//	defer deleteContainer(s.T(), containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := createNewPageBlob(s.T(), "src" + blobName, containerClient)
//	copyPBClient := getPageBlobClient("dst" + blobName, containerClient)
//
//	snapshot := time.Now().UTC().Format(SnapshotTimeFormat)
//	_, err = copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, nil)
//	require.Error(s.T(), err)
//
//	validateStorageError(s.T(), err, StorageErrorCodeCannotVerifyCopySource)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//
//	defer deleteContainer(s.T(), containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.NoError(s.T(), err)
//
//	validateIncrementalCopy(s.T(), copyPBClient, &resp)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//
//	defer deleteContainer(s.T(), containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.Error(s.T(), err)
//
//	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//
//	defer deleteContainer(s.T(), containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.NoError(s.T(), err)
//
//	validateIncrementalCopy(s.T(), copyPBClient, &resp)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//
//	defer deleteContainer(s.T(), containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.Error(s.T(), err)
//
//	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchTrue() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//	resp, _ := copyPBClient.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfMatch: resp.ETag,
//		},
//	}
//	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.NoError(s.T(), err)
//
//	validateIncrementalCopy(s.T(), copyPBClient, &resp2)
//	defer deleteContainer(s.T(), containerClient)
//}
//

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchFalse() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//
//	defer deleteContainer(s.T(), containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfMatch: &eTag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.Error(s.T(), err)
//
//	validateStorageError(s.T(), err, StorageErrorCodeTargetConditionNotMet)
//}
//

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//	defer deleteContainer(s.T(), containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfNoneMatch: &eTag,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.NoError(s.T(), err)
//
//	validateIncrementalCopy(s.T(), copyPBClient, &resp)
//}
//

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse() {
//
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(s.T(), testName)
//	defer deleteContainer(s.T(), containerClient)
//
//	resp, _ := copyPBClient.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfNoneMatch: resp.ETag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	require.Error(s.T(), err)
//
//	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
//}
