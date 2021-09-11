// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"time"
)

func (s *azblobTestSuite) TestPutGetPages() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	contentSize := 1024
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	reader, _ := generateData(1024)
	putResp, err := pbClient.UploadPages(context.Background(), reader, &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(putResp.RawResponse.StatusCode, 201)
	_assert.NotNil(putResp.LastModified)
	_assert.Equal((*putResp.LastModified).IsZero(), false)
	_assert.NotNil(putResp.ETag)
	_assert.Nil(putResp.ContentMD5)
	_assert.Equal(*putResp.BlobSequenceNumber, int64(0))
	_assert.NotNil(*putResp.RequestID)
	_assert.NotNil(*putResp.Version)
	_assert.NotNil(putResp.Date)
	_assert.Equal((*putResp.Date).IsZero(), false)

	pageList, err := pbClient.GetPageRanges(context.Background(), HttpRange{0, 1023}, nil)
	_assert.Nil(err)
	_assert.Equal(pageList.RawResponse.StatusCode, 200)
	_assert.NotNil(pageList.LastModified)
	_assert.Equal((*pageList.LastModified).IsZero(), false)
	_assert.NotNil(pageList.ETag)
	_assert.Equal(*pageList.BlobContentLength, int64(512*10))
	_assert.NotNil(*pageList.RequestID)
	_assert.NotNil(*pageList.Version)
	_assert.NotNil(pageList.Date)
	_assert.Equal((*pageList.Date).IsZero(), false)
	_assert.NotNil(pageList.PageList)
	pageRangeResp := pageList.PageList.PageRange
	_assert.Len(pageRangeResp, 1)
	rawStart, rawEnd := (pageRangeResp)[0].Raw()
	_assert.Equal(rawStart, offset)
	_assert.Equal(rawEnd, count-1)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(_assert, "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(_assert, "dstblob", containerClient, int64(contentSize))

	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp1.RawResponse.StatusCode, 201)
	_assert.NotNil(uploadSrcResp1.LastModified)
	_assert.Equal((*uploadSrcResp1.LastModified).IsZero(), false)
	_assert.NotNil(uploadSrcResp1.ETag)
	_assert.Nil(uploadSrcResp1.ContentMD5)
	_assert.Equal(*uploadSrcResp1.BlobSequenceNumber, int64(0))
	_assert.NotNil(*uploadSrcResp1.RequestID)
	_assert.NotNil(*uploadSrcResp1.Version)
	_assert.NotNil(uploadSrcResp1.Date)
	_assert.Equal((*uploadSrcResp1.Date).IsZero(), false)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		_assert.Error(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL.
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
	_assert.Nil(err)
	_assert.Equal(pResp1.RawResponse.StatusCode, 201)
	_assert.NotNil(pResp1.ETag)
	_assert.NotNil(pResp1.LastModified)
	_assert.NotNil(pResp1.ContentMD5)
	_assert.NotNil(pResp1.RequestID)
	_assert.NotNil(pResp1.Version)
	_assert.NotNil(pResp1.Date)
	_assert.Equal((*pResp1.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, sourceData)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURLWithMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	contentMD5 := md5Value[:]
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(_assert, "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(_assert, "dstblob", containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp1.RawResponse.StatusCode, 201)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		_assert.Error(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL with MD5.
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
	}
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_assert.Nil(err)
	_assert.Equal(pResp1.RawResponse.StatusCode, 201)
	_assert.NotNil(pResp1.ETag)
	_assert.NotNil(pResp1.LastModified)
	_assert.NotNil(pResp1.ContentMD5)
	_assert.EqualValues(pResp1.ContentMD5, contentMD5)
	_assert.NotNil(pResp1.RequestID)
	_assert.NotNil(pResp1.Version)
	_assert.NotNil(pResp1.Date)
	_assert.Equal((*pResp1.Date).IsZero(), false)
	_assert.Equal(*pResp1.BlobSequenceNumber, int64(0))

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, sourceData)

	// Upload page from URL with bad MD5
	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions = UploadPagesFromURLOptions{
		SourceContentMD5: badContentMD5,
	}
	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestClearDiffPages() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	contentSize := 2 * 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_assert.Nil(err)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_assert.Nil(err)

	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}}
	_, err = pbClient.UploadPages(context.Background(), getReaderToGeneratedBytes(2048), &uploadPagesOptions1)
	_assert.Nil(err)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
	_assert.Nil(err)
	pageRangeResp := pageListResp.PageList.PageRange
	_assert.NotNil(pageRangeResp)
	_assert.Len(pageRangeResp, 1)
	// _assert.((pageRangeResp)[0], chk.DeepEquals, PageRange{Start: &offset1, End: &end1})
	rawStart, rawEnd := (pageRangeResp)[0].Raw()
	_assert.Equal(rawStart, offset1)
	_assert.Equal(rawEnd, end1)

	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, nil)
	_assert.Nil(err)
	_assert.Equal(clearResp.RawResponse.StatusCode, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
	_assert.Nil(err)
	_assert.Nil(pageListResp.PageList.PageRange)
}

//nolint
func waitForIncrementalCopy(_assert *assert.Assertions, copyBlobClient PageBlobClient, blobCopyResponse *PageBlobCopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult GetBlobPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(ctx, nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_assert.Fail("")
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

//nolint
func (s *azblobUnrecordedTestSuite) TestIncrementalCopy() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	accessType := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_assert.Nil(err)

	srcBlob := createNewPageBlob(_assert, "src"+generateBlobName(testName), containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = srcBlob.UploadPages(context.Background(), r, &uploadPagesOptions)
	_assert.Nil(err)

	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
	_assert.Nil(err)

	dstBlob := containerClient.NewPageBlobClient("dst" + generateBlobName(testName))

	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 202)
	_assert.NotNil(resp.LastModified)
	_assert.Equal((*resp.LastModified).IsZero(), false)
	_assert.NotNil(resp.ETag)
	_assert.NotEqual(*resp.RequestID, "")
	_assert.NotEqual(*resp.Version, "")
	_assert.NotNil(resp.Date)
	_assert.Equal((*resp.Date).IsZero(), false)
	_assert.NotEqual(*resp.CopyID, "")
	_assert.Equal(*resp.CopyStatus, CopyStatusTypePending)

	waitForIncrementalCopy(_assert, dstBlob, &resp)
}

func (s *azblobTestSuite) TestResizePageBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	resp, err := pbClient.Resize(context.Background(), 2048, nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)

	resp, err = pbClient.Resize(context.Background(), 8192, nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)

	resp2, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp2.ContentLength, int64(8192))
}

func (s *azblobTestSuite) TestPageSequenceNumbers() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	sequenceNumber := int64(0)
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(7)
	actionType = SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(11)
	actionType = SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPutPagesWithMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.Nil(err)
	_assert.Equal(putResp.RawResponse.StatusCode, 201)
	_assert.NotNil(putResp.LastModified)
	_assert.Equal((*putResp.LastModified).IsZero(), false)
	_assert.NotNil(putResp.ETag)
	_assert.NotNil(putResp.ContentMD5)
	_assert.EqualValues(putResp.ContentMD5, contentMD5)
	_assert.Equal(*putResp.BlobSequenceNumber, int64(0))
	_assert.NotNil(*putResp.RequestID)
	_assert.NotNil(*putResp.Version)
	_assert.NotNil(putResp.Date)
	_assert.Equal((*putResp.Date).IsZero(), false)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobCreatePageSizeInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, 1, &createPageBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobCreatePageSequenceInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataNonEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Metadata)
	_assert.EqualValues(resp.Metadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestBlobCreatePageMetadataInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{"In valid1": "bar"},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.NotNil(err)
	_assert.Contains(err.Error(), invalidHeaderErrorSubstring)

}

func (s *azblobTestSuite) TestBlobCreatePageHTTPHeaders() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		HTTPHeaders:        &basicHeaders,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	h := resp.GetHTTPHeaders()
	_assert.EqualValues(h, basicHeaders)
}

func validatePageBlobPut(_assert *assert.Assertions, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Metadata)
	_assert.EqualValues(resp.Metadata, basicMetadata)
	_assert.EqualValues(resp.GetHTTPHeaders(), basicHeaders)
}

func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_assert.Nil(err)

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
	_assert.Nil(err)

	validatePageBlobPut(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_assert.Nil(err)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_assert.Nil(err)

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
	_assert.Nil(err)

	validatePageBlobPut(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	_assert.Nil(err)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

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
	_assert.Nil(err)

	validatePageBlobPut(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.Nil(err)

	validatePageBlobPut(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobPutPagesInvalidRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize/2)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.NotNil(err)
}

//// Body cannot be nil check already added in the request preparer
////func (s *azblobTestSuite) TestBlobPutPagesNilBody() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	pbClient, _ := createNewPageBlob(c, containerClient)
////
////	_, err := pbClient.UploadPages(ctx, nil, nil)
////	_assert.NotNil(err)
////}

func (s *azblobTestSuite) TestBlobPutPagesEmptyBody() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	r := bytes.NewReader([]byte{})
	offset, count := int64(0), int64(0)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobPutPagesNonExistentBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeBlobNotFound)
}

func validateUploadPages(_assert *assert.Assertions, pbClient PageBlobClient) {
	// This will only validate a single put page at 0-PageBlobPageBytes-1
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	_assert.Nil(err)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	_assert.Equal(rawStart, start)
	_assert.Equal(rawEnd, end)
}

func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.Nil(err)

	validateUploadPages(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.Nil(err)

	validateUploadPages(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.Nil(err)

	validateUploadPages(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.Nil(err)

	validateUploadPages(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.Nil(err)

	validateUploadPages(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTETrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

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
	_assert.Nil(err)

	validateUploadPages(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTENegOne() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

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
	_assert.Nil(err)

	validateUploadPages(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		_assert.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := createNewPageBlob(_assert, blobName, containerClient)
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
//	_assert.Nil(err)
//}

func setupClearPagesTest(_assert *assert.Assertions, testName string) (ContainerClient, PageBlobClient) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)

	return containerClient, pbClient
}

func validateClearPagesTest(_assert *assert.Assertions, pbClient PageBlobClient) {
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	_assert.Nil(err)
	pageListResp := resp.PageList.PageRange
	_assert.Nil(pageListResp)
}

func (s *azblobTestSuite) TestBlobClearPagesInvalidRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes + 1}, nil)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.Nil(err)

	validateClearPagesTest(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.Nil(err)

	validateClearPagesTest(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: getPropertiesResp.ETag,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.Nil(err)

	validateClearPagesTest(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.Nil(err)

	validateClearPagesTest(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.Nil(err)

	validateClearPagesTest(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTETrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.Nil(err)

	validateClearPagesTest(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTEFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTENegOne() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.Nil(err)

	validateClearPagesTest(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.Nil(err)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidInput)
}

func setupGetPageRangesTest(_assert *assert.Assertions, testName string) (containerClient ContainerClient, pbClient PageBlobClient) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient = createNewContainer(_assert, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient = createNewPageBlob(_assert, blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)
	return
}

func validateBasicGetPageRanges(_assert *assert.Assertions, resp PageList, err error) {
	_assert.Nil(err)
	_assert.NotNil(resp.PageRange)
	_assert.Len(resp.PageRange, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := (resp.PageRange)[0].Raw()
	_assert.Equal(rawStart, start)
	_assert.Equal(rawEnd, end)
}

func (s *azblobTestSuite) TestBlobGetPageRangesEmptyBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	_assert.Nil(err)
	_assert.Nil(resp.PageList.PageRange)
}

func (s *azblobTestSuite) TestBlobGetPageRangesEmptyRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	_assert.Nil(err)
	validateBasicGetPageRanges(_assert, resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesInvalidRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	_, err := pbClient.GetPageRanges(ctx, HttpRange{-2, 500}, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesNonContiguousRanges() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(2*PageBlobPageBytes), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	_assert.Nil(err)
	pageListResp := resp.PageList.PageRange
	_assert.NotNil(pageListResp)
	_assert.Len(pageListResp, 2)

	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	_assert.Equal(rawStart, start)
	_assert.Equal(rawEnd, end)

	start, end = int64(PageBlobPageBytes*2), int64((PageBlobPageBytes*3)-1)
	rawStart, rawEnd = pageListResp[1].Raw()
	_assert.Equal(rawStart, start)
	_assert.Equal(rawEnd, end)
}

func (s *azblobTestSuite) TestBlobGetPageRangesNotPageAligned() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 2000}, nil)
	_assert.Nil(err)
	validateBasicGetPageRanges(_assert, resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Snapshot)

	snapshotURL := pbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	_assert.Nil(err)
	validateBasicGetPageRanges(_assert, resp2.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.Nil(err)
	validateBasicGetPageRanges(_assert, resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.NotNil(err)

	//serr := err.(StorageError)
	//_assert.(serr.RawResponse.StatusCode, chk.Equals, 304)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.Nil(err)
	validateBasicGetPageRanges(_assert, resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.Nil(err)
	validateBasicGetPageRanges(_assert, resp2.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.Nil(err)
	validateBasicGetPageRanges(_assert, resp.PageList, err)
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	_assert.NotNil(err)
	//serr := err.(StorageError)
	//_assert.(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
}

//nolint
func setupDiffPageRangesTest(_assert *assert.Assertions, testName string) (containerClient ContainerClient,
	pbClient PageBlobClient, snapshot string) {
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient = createNewContainer(_assert, containerName, svcClient)

	blobName := generateName(testName)
	pbClient = createNewPageBlob(_assert, blobName, containerClient)

	r := getReaderToGeneratedBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	snapshot = *resp.Snapshot

	r = getReaderToGeneratedBytes(PageBlobPageBytes)
	offset, count = int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions = UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)
	return
}

//nolint
func validateDiffPageRanges(_assert *assert.Assertions, resp PageList, err error) {
	_assert.Nil(err)
	pageListResp := resp.PageRange
	_assert.NotNil(pageListResp)
	_assert.Len(resp.PageRange, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	_assert.EqualValues(rawStart, start)
	_assert.EqualValues(rawEnd, end)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangesNonExistentSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	snapshotTime, _ := time.Parse(SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshotTime.Format(SnapshotTimeFormat), nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodePreviousSnapshotNotFound)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeInvalidRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{-22, 14}, snapshot, nil)
	_assert.Nil(err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.NotNil(err)

	//stgErr := err.(StorageError)
	//_assert.(stgErr.Response().StatusCode, chk.Equals, 304)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp2.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp.PageList, err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobResizeZero() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	// The default pbClient is created with size > 0, so this should actually update
	_, err = pbClient.Resize(ctx, 0, nil)
	_assert.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.ContentLength, int64(0))
}

func (s *azblobTestSuite) TestBlobResizeInvalidSizeNegative() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	_, err = pbClient.Resize(ctx, -4, nil)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobResizeInvalidSizeMisaligned() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	_, err = pbClient.Resize(ctx, 12, nil)
	_assert.NotNil(err)
}

func validateResize(_assert *assert.Assertions, pbClient PageBlobClient) {
	resp, _ := pbClient.GetProperties(ctx, nil)
	_assert.Equal(*resp.ContentLength, int64(PageBlobPageBytes))
}

func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.Nil(err)

	validateResize(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.Nil(err)

	validateResize(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.Nil(err)

	validateResize(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobResizeIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.Nil(err)

	validateResize(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobResizeIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberActionTypeInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberSequenceNumberInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func validateSequenceNumberSet(_assert *assert.Assertions, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.BlobSequenceNumber, int64(1))
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.Nil(err)

	validateSequenceNumberSet(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.Nil(err)

	validateSequenceNumberSet(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	_assert.Nil(err)
	_assert.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_assert.NotNil(pageBlobCreateResponse.Date)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.Nil(err)

	validateSequenceNumberSet(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, blobName, containerClient)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, "src"+blobName, containerClient)

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
	_assert.Nil(err)

	validateSequenceNumberSet(_assert, pbClient)
}

func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_assert, "src"+blobName, containerClient)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//func setupStartIncrementalCopyTest(_assert *assert.Assertions, testName string) (containerClient ContainerClient,
//	pbClient PageBlobClient, copyPBClient PageBlobClient, snapshot string) {
//	_context := getTestContext(testName)
//	var recording *testframework.Recording
//	if _context != nil {
//		recording = _context.recording
//	}
//	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
//	if err != nil {
//		_assert.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient = createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	accessType := PublicAccessTypeBlob
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
//	}
//	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	_assert.Nil(err)
//
//	pbClient = createNewPageBlob(_assert, generateBlobName(testName), containerClient)
//	resp, _ := pbClient.CreateSnapshot(ctx, nil)
//
//	copyPBClient = getPageBlobClient("copy"+generateBlobName(testName), containerClient)
//
//	// Must create the incremental copy pbClient so that the access conditions work on it
//	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), *resp.Snapshot, nil)
//	_assert.Nil(err)
//	waitForIncrementalCopy(_assert, copyPBClient, &resp2)
//
//	resp, _ = pbClient.CreateSnapshot(ctx, nil) // Take a new snapshot so the next copy will succeed
//	snapshot = *resp.Snapshot
//	return
//}

//func validateIncrementalCopy(_assert *assert.Assertions, copyPBClient PageBlobClient, resp *PageBlobCopyIncrementalResponse) {
//	t := waitForIncrementalCopy(_assert, copyPBClient, resp)
//
//	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
//	copySnapshotURL := copyPBClient.WithSnapshot(*t)
//	_, err := copySnapshotURL.GetProperties(ctx, nil)
//	_assert.Nil(err)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopySnapshotNotExist() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		_assert.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := createNewPageBlob(_assert, "src" + blobName, containerClient)
//	copyPBClient := getPageBlobClient("dst" + blobName, containerClient)
//
//	snapshot := time.Now().UTC().Format(SnapshotTimeFormat)
//	_, err = copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, nil)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeCannotVerifyCopySource)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.Nil(err)
//
//	validateIncrementalCopy(_assert, copyPBClient, &resp)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.Nil(err)
//
//	validateIncrementalCopy(_assert, copyPBClient, &resp)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchTrue() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//	resp, _ := copyPBClient.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfMatch: resp.ETag,
//		},
//	}
//	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.Nil(err)
//
//	validateIncrementalCopy(_assert, copyPBClient, &resp2)
//	defer deleteContainer(_assert, containerClient)
//}
//

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchFalse() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//
//	defer deleteContainer(_assert, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfMatch: &eTag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeTargetConditionNotMet)
//}
//

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//	defer deleteContainer(_assert, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfNoneMatch: &eTag,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.Nil(err)
//
//	validateIncrementalCopy(_assert, copyPBClient, &resp)
//}
//

//nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_assert, testName)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, _ := copyPBClient.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfNoneMatch: resp.ETag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}
