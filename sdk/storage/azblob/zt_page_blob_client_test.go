//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"bytes"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/stretchr/testify/require"
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

	offset, count := int64(0), int64(1024)
	reader, _ := generateData(1024)
	putResp, err := pbClient.UploadPages(ctx, reader, &pageblob.UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
	})
	_require.Nil(err)
	_require.NotNil(putResp.LastModified)
	_require.Equal((*putResp.LastModified).IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.Nil(putResp.ContentMD5)
	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
	_require.NotNil(*putResp.RequestID)
	_require.NotNil(*putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Offset: to.Ptr(int64(0)),
		Count:  to.Ptr(int64(1023)),
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(ctx)
		_require.Nil(err)
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
		if err != nil {
			break
		}
	}
}

//nolint
//func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(contentSize)
//	srcBlob := createNewPageBlobWithSize(_require, "srcblob", containerClient, int64(contentSize))
//	destBlob := createNewPageBlobWithSize(_require, "dstblob", containerClient, int64(contentSize))
//
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadSrcResp1, err := srcBlob.UploadPages(ctx, NopCloser(r), &pageblob.UploadPagesOptions{
//		Offset: to.Ptr(offset),
//		Count: to.Ptr(count),
//	}
//	_require.Nil(err)
//	_require.NotNil(uploadSrcResp1.LastModified)
//	_require.Equal((*uploadSrcResp1.LastModified).IsZero(), false)
//	_require.NotNil(uploadSrcResp1.ETag)
//	_require.Nil(uploadSrcResp1.ContentMD5)
//	_require.Equal(*uploadSrcResp1.BlobSequenceNumber, int64(0))
//	_require.NotNil(*uploadSrcResp1.RequestID)
//	_require.NotNil(*uploadSrcResp1.Version)
//	_require.NotNil(uploadSrcResp1.Date)
//	_require.Equal((*uploadSrcResp1.Date).IsZero(), false)
//
//	// Get source pbClient URL with SAS for UploadPagesFromURL.
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	_require.Nil(err)
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		_require.Error(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Upload page from URL.
//	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
//	_require.Nil(err)
//	// _require.Equal(pResp1.RawResponse.StatusCode, 201)
//	_require.NotNil(pResp1.ETag)
//	_require.NotNil(pResp1.LastModified)
//	_require.NotNil(pResp1.ContentMD5)
//	_require.NotNil(pResp1.RequestID)
//	_require.NotNil(pResp1.Version)
//	_require.NotNil(pResp1.Date)
//	_require.Equal((*pResp1.Date).IsZero(), false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	_require.Nil(err)
//	destData, err := ioutil.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{}))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
//}
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURLWithMD5() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(contentSize)
//	md5Value := md5.Sum(sourceData)
//	contentMD5 := md5Value[:]
//	ctx := ctx // Use default Background context
//	srcBlob := createNewPageBlobWithSize(_require, "srcblob", containerClient, int64(contentSize))
//	destBlob := createNewPageBlobWithSize(_require, "dstblob", containerClient, int64(contentSize))
//
//	// Prepare source pbClient for copy.
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),}
//	_, err = srcBlob.UploadPages(ctx, NopCloser(r), &uploadPagesOptions)
//	_require.Nil(err)
//	// _require.Equal(uploadSrcResp1.RawResponse.StatusCode, 201)
//
//	// Get source pbClient URL with SAS for UploadPagesFromURL.
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	_require.Nil(err)
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = azblob.BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		_require.Error(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Upload page from URL with MD5.
//	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: contentMD5,
//	}
//	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_require.Nil(err)
//	// _require.Equal(pResp1.RawResponse.StatusCode, 201)
//	_require.NotNil(pResp1.ETag)
//	_require.NotNil(pResp1.LastModified)
//	_require.NotNil(pResp1.ContentMD5)
//	_require.EqualValues(pResp1.ContentMD5, contentMD5)
//	_require.NotNil(pResp1.RequestID)
//	_require.NotNil(pResp1.Version)
//	_require.NotNil(pResp1.Date)
//	_require.Equal((*pResp1.Date).IsZero(), false)
//	_require.Equal(*pResp1.BlobSequenceNumber, int64(0))
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	_require.Nil(err)
//	destData, err := ioutil.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{}))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
//
//	// Upload page from URL with bad MD5
//	_, badMD5 := getRandomDataAndReader(16)
//	badContentMD5 := badMD5[:]
//	uploadPagesFromURLOptions = pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: badContentMD5,
//	}
//	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_require.NotNil(err)
//
//	validateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
//}

//nolint
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
	_, err = pbClient.UploadPages(ctx, r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(0)),
		Count:  to.Ptr(int64(contentSize)),
	})
	_require.Nil(err)

	snapshotResp, err := pbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)

	r1 := getReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(ctx, r1, &pageblob.UploadPagesOptions{Offset: to.Ptr(int64(contentSize)), Count: to.Ptr(int64(contentSize))})
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset:       to.Ptr(int64(0)),
		Count:        to.Ptr(int64(4096)),
		PrevSnapshot: snapshotResp.Snapshot,
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(ctx)
		_require.Nil(err)

		pageRangeResp := pageListResp.PageList.PageRange
		_require.NotNil(pageRangeResp)
		_require.Len(pageRangeResp, 1)
		rawStart, rawEnd := (pageRangeResp)[0].Raw()
		_require.Equal(rawStart, int64(2048))
		_require.Equal(rawEnd, int64(4095))
		if err != nil {
			break
		}
	}

	_, err = pbClient.ClearPages(ctx, int64(2048), int64(2048), nil)
	_require.Nil(err)

	pager = pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset:       to.Ptr(int64(0)),
		Count:        to.Ptr(int64(4096)),
		PrevSnapshot: snapshotResp.Snapshot,
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(ctx)
		_require.Nil(err)
		pageRangeResp := pageListResp.PageList.PageRange
		_require.Len(pageRangeResp, 0)
		if err != nil {
			break
		}
	}
}

//nolint
func waitForIncrementalCopy(_require *require.Assertions, copyBlobClient *pageblob.Client, blobCopyResponse *pageblob.CopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult blob.GetPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != blob.CopyStatusTypeSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(ctx, nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_require.Fail("")
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

//nolint
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

	_, err = containerClient.SetAccessPolicy(ctx, &container.SetAccessPolicyOptions{Access: to.Ptr(container.PublicAccessTypeBlob)})
	_require.Nil(err)

	srcBlob := createNewPageBlob(_require, "src"+generateBlobName(testName), containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize)
	_, err = srcBlob.UploadPages(ctx, r, &pageblob.UploadPagesOptions{Offset: to.Ptr(offset), Count: to.Ptr(count)})
	_require.Nil(err)

	snapshotResp, err := srcBlob.CreateSnapshot(ctx, nil)
	_require.Nil(err)

	dstBlob := containerClient.NewPageBlobClient("dst" + generateBlobName(testName))

	resp, err := dstBlob.StartCopyIncremental(ctx, srcBlob.URL(), *snapshotResp.Snapshot, nil)
	_require.Nil(err)
	_require.NotNil(resp.LastModified)
	_require.Equal((*resp.LastModified).IsZero(), false)
	_require.NotNil(resp.ETag)
	_require.NotEqual(*resp.RequestID, "")
	_require.NotEqual(*resp.Version, "")
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)
	_require.NotEqual(*resp.CopyID, "")
	_require.Equal(*resp.CopyStatus, blob.CopyStatusTypePending)

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

	_, err = pbClient.Resize(ctx, 2048, nil)
	_require.Nil(err)

	_, err = pbClient.Resize(ctx, 8192, nil)
	_require.Nil(err)

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
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	sequenceNumber = int64(7)
	actionType = pageblob.SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	updateSequenceNumberPageBlob = pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: to.Ptr(int64(11)),
		ActionType:     to.Ptr(pageblob.SequenceNumberActionTypeUpdate),
	}

	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)
}

//nolint
//func (s *azblobUnrecordedTestSuite) TestPutPagesWithMD5() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := createNewPageBlob(_require, blobName, containerClient)
//
//	// put page with valid MD5
//	contentSize := 1024
//	readerToBody, body := getRandomDataAndReader(contentSize)
//	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
//	md5Value := md5.Sum(body)
//	_ = body
//	contentMD5 := md5Value[:]
//
//	putResp, err := pbClient.UploadPages(ctx, NopCloser(readerToBody), &pageblob.UploadPagesOptions{
//		Offset:                  to.Ptr(offset),
//		Count:                   to.Ptr(count),
//		TransactionalContentMD5: contentMD5,
//	})
//	_require.Nil(err)
//	// _require.Equal(putResp.RawResponse.StatusCode, 201)
//	_require.NotNil(putResp.LastModified)
//	_require.Equal((*putResp.LastModified).IsZero(), false)
//	_require.NotNil(putResp.ETag)
//	_require.NotNil(putResp.ContentMD5)
//	_require.EqualValues(putResp.ContentMD5, contentMD5)
//	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
//	_require.NotNil(*putResp.RequestID)
//	_require.NotNil(*putResp.Version)
//	_require.NotNil(putResp.Date)
//	_require.Equal((*putResp.Date).IsZero(), false)
//
//	// put page with bad MD5
//	readerToBody, _ = getRandomDataAndReader(1024)
//	_, badMD5 := getRandomDataAndReader(16)
//	basContentMD5 := badMD5[:]
//	putResp, err = pbClient.UploadPages(ctx, NopCloser(readerToBody), &pageblob.UploadPagesOptions{
//		Offset:                  to.Ptr(offset),
//		Count:                   to.Ptr(count),
//		TransactionalContentMD5: basContentMD5,
//	})
//	_require.NotNil(err)
//
//	validateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
//}

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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, 1, &createPageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       map[string]string{},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       map[string]string{"In valid1": "bar"},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		HTTPHeaders:    &basicHeaders,
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	h := blob.ParseHTTPHeaders(resp)
	_require.EqualValues(h, basicHeaders)
}

func validatePageBlobPut(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, basicMetadata)
	_require.EqualValues(blob.ParseHTTPHeaders(resp), basicHeaders)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
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
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       basicMetadata,
		HTTPHeaders:    &basicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(ctx, pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

//nolint
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
	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count))}
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
	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count))}
	_, err = pbClient.UploadPages(ctx, NopCloser(r), &uploadPagesOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count))}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func validateUploadPages(_require *require.Assertions, pbClient *pageblob.Client) {
	// This will only validate a single put page at 0-PageBlobPageBytes-1
	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Offset: to.Ptr(int64(0)),
		Count:  to.Ptr(int64(blob.CountToEnd)),
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(ctx)
		_require.Nil(err)

		start, end := int64(0), int64(pageblob.PageBytes-1)
		rawStart, rawEnd := *(pageListResp.PageList.PageRange[0].Start), *(pageListResp.PageList.PageRange[0].End)
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(ctx, r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(ctx, r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(ctx, r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(ctx, r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(ctx, r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	eTag := "garbage"
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	eTag := "garbage"
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThan := int64(10)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
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
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThan := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
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

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}

	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidInput)
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
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
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
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
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

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
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
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
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

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := createNewPageBlob(_require, blobName, containerClient)
//
//	r, _ := generateData(pageblob.PageBytes)
//	offset, count := int64(0), int64(pageblob.PageBytes)
//	ifSequenceNumberEqualTo := int64(-1)
//	uploadPagesOptions := pageblob.UploadPagesOptions{
//		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
//		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions) // This will cause the library to set the value of the header to 0
//	_require.Nil(err)
//}

func setupClearPagesTest(_require *require.Assertions, testName string) (*container.Client, *pageblob.Client) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	return containerClient, pbClient
}

func validateClearPagesTest(_require *require.Assertions, pbClient *pageblob.Client) {
	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		pageListResp, err := pager.NextPage(ctx)
		_require.Nil(err)
		_require.Nil(pageListResp.PageRange)
		if err != nil {
			break
		}
	}

}

func (s *azblobTestSuite) TestBlobClearPagesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes+1), nil)
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

	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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

	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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

	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: getPropertiesResp.ETag,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	eTag := "garbage"
	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	eTag := "garbage"
	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTETrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTEFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTENegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, int64(0), int64(pageblob.PageBytes), &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func setupGetPageRangesTest(_require *require.Assertions, testName string) (containerClient *container.Client, pbClient *pageblob.Client) {
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_require.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient = createNewContainer(_require, containerName, svcClient)

	blobName := generateBlobName(testName)
	pbClient = createNewPageBlob(_require, blobName, containerClient)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)
	return
}

func validateBasicGetPageRanges(_require *require.Assertions, resp pageblob.PageList, err error) {
	_require.Nil(err)
	_require.NotNil(resp.PageRange)
	_require.Len(resp.PageRange, 1)
	start, end := int64(0), int64(pageblob.PageBytes-1)
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

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		_require.Nil(resp.PageRange)
		if err != nil {
			break
		}
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesEmptyRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(-2)), Count: to.Ptr(int64(500))})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.Nil(err)
		if err != nil {
			break
		}
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesNonContiguousRanges() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	r, _ := generateData(pageblob.PageBytes)
	offset, count := int64(2*pageblob.PageBytes), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		pageListResp := resp.PageList.PageRange
		_require.NotNil(pageListResp)
		_require.Len(pageListResp, 2)

		start, end := int64(0), int64(pageblob.PageBytes-1)
		rawStart, rawEnd := pageListResp[0].Raw()
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)

		start, end = int64(pageblob.PageBytes*2), int64((pageblob.PageBytes*3)-1)
		rawStart, rawEnd = pageListResp[1].Raw()
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesNotPageAligned() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(2000))})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
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
	pager := snapshotURL.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp2, err := pager.NextPage(ctx)
		_require.Nil(err)

		validateBasicGetPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
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

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
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

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
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

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}})
	for pager.More() {
		resp2, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfMatch: to.Ptr("garbage"),
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfNoneMatch: to.Ptr("garbage"),
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		if err != nil {
			break
		}
	}

	//serr := err.(StorageError)
	//_require.(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
}

//nolint
func setupDiffPageRangesTest(_require *require.Assertions, testName string) (containerClient *container.Client, pbClient *pageblob.Client, snapshot string) {
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

	r := getReaderToGeneratedBytes(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(offset),
		Count:  to.Ptr(count),
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	snapshot = *resp.Snapshot

	r = getReaderToGeneratedBytes(pageblob.PageBytes)
	offset, count = int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions = pageblob.UploadPagesOptions{Offset: to.Ptr(offset), Count: to.Ptr(count)}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)
	return
}

//nolint
func validateDiffPageRanges(_require *require.Assertions, resp pageblob.PageList, err error) {
	_require.Nil(err)
	_require.NotNil(resp.PageRange)
	_require.Len(resp.PageRange, 1)
	rawStart, rawEnd := resp.PageRange[0].Raw()
	_require.EqualValues(rawStart, int64(0))
	_require.EqualValues(rawEnd, int64(pageblob.PageBytes-1))
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangesNonExistentSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	snapshotTime, _ := time.Parse(blob.SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		PrevSnapshot: to.Ptr(snapshotTime.Format(blob.SnapshotTimeFormat))})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		validateBlobErrorCode(_require, err, bloberror.PreviousSnapshotNotFound)
		if err != nil {
			break
		}
	}

}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)
	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{Offset: to.Ptr(int64(-22)), Count: to.Ptr(int64(14)), Snapshot: &snapshot})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.Nil(err)
		if err != nil {
			break
		}
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	for pager.More() {
		resp2, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateDiffPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateDiffPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

////nolint
//func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
//	defer deleteContainer(_require, containerClient)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	_require.Nil(err)
//
//	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
//		Snapshot: to.Ptr(snapshot),
//		LeaseAccessConditions: &blob.LeaseAccessConditions{
//			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	})
//	for pager.More() {
//		resp2, err := pager.NextPage(ctx)
//		_require.Nil(err)
//		validateDiffPageRanges(_require, resp2.PageList, err)
//		if err != nil {
//			break
//		}
//	}
//}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshotStr),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr("garbage"),
			},
		}})

	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}

	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		PrevSnapshot: to.Ptr(snapshotStr),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: to.Ptr("garbage"),
			},
		}})

	for pager.More() {
		resp2, err := pager.NextPage(ctx)
		_require.Nil(err)
		validateDiffPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_require, testName)
	defer deleteContainer(_require, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		PrevSnapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	})

	for pager.More() {
		_, err := pager.NextPage(ctx)
		_require.NotNil(err)
		if err != nil {
			break
		}
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

func validateResize(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, _ := pbClient.GetProperties(ctx, nil)
	_require.Equal(*resp.ContentLength, int64(pageblob.PageBytes))
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
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
	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
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

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	actionType := pageblob.SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
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
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func validateSequenceNumberSet(_require *require.Assertions, pbClient *pageblob.Client) {
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
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
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
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

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

//func setupStartIncrementalCopyTest(_require *require.Assertions, testName string) (containerClient *container.Client,
//	pbClient *pageblob.Client, copyPBClient *pageblob.Client, snapshot string) {
//	_context := getTestContext(testName)
//	var recording *testframework.Recording
//	if _context != nil {
//		recording = _context.recording
//	}
//	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient = createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	accessType := container.PublicAccessTypeBlob
//	setAccessPolicyOptions := container.SetAccessPolicyOptions{
//		Access: &accessType,
//	}
//	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
//	_require.Nil(err)
//
//	pbClient = createNewPageBlob(_require, generateBlobName(testName), containerClient)
//	resp, _ := pbClient.CreateSnapshot(ctx, nil)
//
//	copyPBClient = getPageBlobClient("copy"+generateBlobName(testName), containerClient)
//
//	// Must create the incremental copy pbClient so that the access conditions work on it
//	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), *resp.Snapshot, nil)
//	_require.Nil(err)
//	waitForIncrementalCopy(_require, copyPBClient, &resp2)
//
//	resp, _ = pbClient.CreateSnapshot(ctx, nil) // Take a new snapshot so the next copy will succeed
//	snapshot = *resp.Snapshot
//	return
//}

//func validateIncrementalCopy(_require *require.Assertions, copyPBClient *pageblob.Client, resp *pageblob.CopyIncrementalResponse) {
//	t := waitForIncrementalCopy(_require, copyPBClient, resp)
//
//	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
//	copySnapshotURL, err := copyPBClient.WithSnapshot(*t)
//	_require.Nil(err)
//	_, err = copySnapshotURL.GetProperties(ctx, nil)
//	_require.Nil(err)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopySnapshotNotExist() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := createNewPageBlob(_require, "src"+blobName, containerClient)
//	copyPBClient := getPageBlobClient("dst"+blobName, containerClient)
//
//	snapshot := time.Now().UTC().Format(blob.SnapshotTimeFormat)
//	_, err = copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, nil)
//	_require.NotNil(err)
//
//	validateBlobErrorCode(_require, err, bloberror.CannotVerifyCopySource)
//}

//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer deleteContainer(_require, containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer deleteContainer(_require, containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer deleteContainer(_require, containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer deleteContainer(_require, containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
//}
//
//// nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	resp, _ := copyPBClient.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfMatch: resp.ETag,
//		},
//	}
//	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp2)
//	defer deleteContainer(_require, containerClient)
//}
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfMatchFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer deleteContainer(_require, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfMatch: &eTag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	validateBlobErrorCode(_require, err, bloberror.TargetConditionNotMet)
//}

////nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	defer deleteContainer(_require, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfNoneMatch: &eTag,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
//}

////nolint
//func (s *azblobUnrecordedTestSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	defer deleteContainer(_require, containerClient)
//
//	resp, _ := copyPBClient.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfNoneMatch: resp.ETag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
//}

func setAndCheckPageBlobTier(_require *require.Assertions, pbClient *pageblob.Client, tier blob.AccessTier) {
	_, err := pbClient.SetTier(ctx, tier, nil)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.AccessTier, string(tier))
}

func (s *azblobTestSuite) TestBlobSetTierAllTiersOnPageBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	premiumServiceClient, err := getServiceClient(_context.recording, testAccountPremium, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	premContainerName := "prem" + generateContainerName(testName)
	premContainerClient := createNewContainer(_require, premContainerName, premiumServiceClient)
	defer deleteContainer(_require, premContainerClient)

	pbName := generateBlobName(testName)
	pbClient := createNewPageBlob(_require, pbName, premContainerClient)

	possibleTiers := []blob.AccessTier{
		blob.AccessTierP4,
		blob.AccessTierP6,
		blob.AccessTierP10,
		blob.AccessTierP20,
		blob.AccessTierP30,
		blob.AccessTierP40,
		blob.AccessTierP50,
	}
	for _, possibleTier := range possibleTiers {
		setAndCheckPageBlobTier(_require, pbClient, possibleTier)
	}
}
