// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testframework"
	"github.com/stretchr/testify/assert"
	"time"
)

//func (s *azblobTestSuite) TestPutGetPages() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	contentSize := 1024
//	offset, end, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	putResp, err := pbClient.UploadPages(context.Background(), getReaderToRandomBytes(1024), &uploadPagesOptions)
//	_assert.Nil(err)
//	_assert.(putResp.RawResponse.StatusCode, chk.Equals, 201)
//	_assert.(putResp.LastModified, chk.NotNil)
//	_assert.((*putResp.LastModified).IsZero(), chk.Equals, false)
//	_assert.(putResp.ETag, chk.NotNil)
//	_assert.(putResp.ContentMD5, chk.IsNil)
//	_assert.(*putResp.BlobSequenceNumber, chk.Equals, int64(0))
//	_assert.(*putResp.RequestID, chk.NotNil)
//	_assert.(*putResp.Version, chk.NotNil)
//	_assert.(putResp.Date, chk.NotNil)
//	_assert.((*putResp.Date).IsZero(), chk.Equals, false)
//
//	pageList, err := pbClient.GetPageRanges(context.Background(), HttpRange{0, 1023}, nil)
//	_assert.Nil(err)
//	_assert.(pageList.RawResponse.StatusCode, chk.Equals, 200)
//	_assert.(pageList.LastModified, chk.NotNil)
//	_assert.((*pageList.LastModified).IsZero(), chk.Equals, false)
//	_assert.(pageList.ETag, chk.NotNil)
//	_assert.(*pageList.BlobContentLength, chk.Equals, int64(512*10))
//	_assert.(*pageList.RequestID, chk.NotNil)
//	_assert.(*pageList.Version, chk.NotNil)
//	_assert.(pageList.Date, chk.NotNil)
//	_assert.((*pageList.Date).IsZero(), chk.Equals, false)
//	_assert.(pageList.PageList, chk.NotNil)
//	pageRangeResp := pageList.PageList.PageRange
//	_assert.(*pageRangeResp, chk.HasLen, 1)
//	rawStart, rawEnd := (*pageRangeResp)[0].Raw()
//	_assert.(rawStart, chk.Equals, offset)
//	_assert.(rawEnd, chk.Equals, end)
//}
//
//func (s *azblobTestSuite) TestUploadPagesFromURL() {
//	svcClient := getServiceClient(nil)
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(contentSize)
//	ctx := context.Background() // Use default Background context
//	srcBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
//	destBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
//
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	uploadSrcResp1, err := srcBlob.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//	_assert.(uploadSrcResp1.RawResponse.StatusCode, chk.Equals, 201)
//	_assert.(uploadSrcResp1.LastModified, chk.NotNil)
//	_assert.((*uploadSrcResp1.LastModified).IsZero(), chk.Equals, false)
//	_assert.(uploadSrcResp1.ETag, chk.NotNil)
//	_assert.(uploadSrcResp1.ContentMD5, chk.IsNil)
//	_assert.(*uploadSrcResp1.BlobSequenceNumber, chk.Equals, int64(0))
//	_assert.(*uploadSrcResp1.RequestID, chk.NotNil)
//	_assert.(*uploadSrcResp1.Version, chk.NotNil)
//	_assert.(uploadSrcResp1.Date, chk.NotNil)
//	_assert.((*uploadSrcResp1.Date).IsZero(), chk.Equals, false)
//
//	// Get source pbClient URL with SAS for UploadPagesFromURL.
//	srcBlobParts := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Upload page from URL.
//	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
//	_assert.Nil(err)
//	_assert.(pResp1.RawResponse.StatusCode, chk.Equals, 201)
//	_assert.(pResp1.ETag, chk.NotNil)
//	_assert.(pResp1.LastModified, chk.NotNil)
//	_assert.(pResp1.ContentMD5, chk.NotNil)
//	_assert.(pResp1.RequestID, chk.NotNil)
//	_assert.(pResp1.Version, chk.NotNil)
//	_assert.(pResp1.Date, chk.NotNil)
//	_assert.((*pResp1.Date).IsZero(), chk.Equals, false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	_assert.Nil(err)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	_assert.Nil(err)
//	_assert.(destData, chk.DeepEquals, sourceData)
//}
//
//func (s *azblobTestSuite) TestUploadPagesFromURLWithMD5() {
//	svcClient := getServiceClient(nil)
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(contentSize)
//	md5Value := md5.Sum(sourceData)
//	contentMD5 := md5Value[:]
//	ctx := context.Background() // Use default Background context
//	srcBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
//	destBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
//
//	// Prepare source pbClient for copy.
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	uploadSrcResp1, err := srcBlob.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//	_assert.(uploadSrcResp1.RawResponse.StatusCode, chk.Equals, 201)
//
//	// Get source pbClient URL with SAS for UploadPagesFromURL.
//	srcBlobParts := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Upload page from URL with MD5.
//	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
//		SourceContentMD5: &contentMD5,
//	}
//	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_assert.Nil(err)
//	_assert.(pResp1.RawResponse.StatusCode, chk.Equals, 201)
//	_assert.(pResp1.ETag, chk.NotNil)
//	_assert.(pResp1.LastModified, chk.NotNil)
//	_assert.(pResp1.ContentMD5, chk.NotNil)
//	_assert.(*pResp1.ContentMD5, chk.DeepEquals, contentMD5)
//	_assert.(pResp1.RequestID, chk.NotNil)
//	_assert.(pResp1.Version, chk.NotNil)
//	_assert.(pResp1.Date, chk.NotNil)
//	_assert.((*pResp1.Date).IsZero(), chk.Equals, false)
//	_assert.(*pResp1.BlobSequenceNumber, chk.Equals, int64(0))
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	_assert.Nil(err)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	_assert.Nil(err)
//	_assert.(destData, chk.DeepEquals, sourceData)
//
//	// Upload page from URL with bad MD5
//	_, badMD5 := getRandomDataAndReader(16)
//	badContentMD5 := badMD5[:]
//	uploadPagesFromURLOptions = UploadPagesFromURLOptions{
//		SourceContentMD5: &badContentMD5,
//	}
//	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeMD5Mismatch)
//}

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
	r := getReaderToRandomBytes(contentSize)
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_assert.Nil(err)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_assert.Nil(err)

	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}}
	_, err = pbClient.UploadPages(context.Background(), getReaderToRandomBytes(2048), &uploadPagesOptions1)
	_assert.Nil(err)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
	_assert.Nil(err)
	pageRangeResp := pageListResp.PageList.PageRange
	_assert.NotNil(pageRangeResp)
	_assert.Len(*pageRangeResp, 1)
	// _assert.((*pageRangeResp)[0], chk.DeepEquals, PageRange{Start: &offset1, End: &end1})
	rawStart, rawEnd := (*pageRangeResp)[0].Raw()
	_assert.Equal(rawStart, offset1)
	_assert.Equal(rawEnd, end1)

	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, nil)
	_assert.Nil(err)
	_assert.Equal(clearResp.RawResponse.StatusCode, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
	_assert.Nil(err)
	_assert.Nil(pageListResp.PageList.PageRange)
}

func waitForIncrementalCopy(_assert *assert.Assertions, copyBlobClient PageBlobClient, blobCopyResponse *PageBlobCopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult BlobGetPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != CopyStatusSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(ctx, nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_assert.Fail("")
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

//func (s *azblobTestSuite) TestIncrementalCopy() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
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
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	accessType := PublicAccessBlob
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
//	}
//	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	_assert.Nil(err)
//
//	srcBlob := createNewPageBlob(_assert, "src"+generateBlobName(testName), containerClient)
//
//	contentSize := 1024
//	r := getReaderToRandomBytes(contentSize)
//	offset, count := int64(0), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	_, err = srcBlob.UploadPages(context.Background(), r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
//	_assert.Nil(err)
//
//	dstBlob := containerClient.NewPageBlobClient("dst" + generateBlobName(testName))
//
//	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
//	_assert.Nil(err)
//	_assert.Equal(resp.RawResponse.StatusCode, 202)
//	_assert.NotNil(resp.LastModified)
//	_assert.Equal((*resp.LastModified).IsZero(), false)
//	_assert.NotNil(resp.ETag)
//	_assert.NotEqual(*resp.RequestID, "")
//	_assert.NotEqual(*resp.Version, "")
//	_assert.NotNil(resp.Date)
//	_assert.Equal((*resp.Date).IsZero(), false)
//	_assert.NotEqual(*resp.CopyID, "")
//	_assert.Equal(*resp.CopyStatus, CopyStatusPending)
//
//	waitForIncrementalCopy(_assert, dstBlob, &resp)
//}

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
	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(7)
	actionType = SequenceNumberActionMax
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(11)
	actionType = SequenceNumberActionUpdate
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 200)
}

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
	contentMD5 := md5Value[:]
	uploadPagesOptions := UploadPagesOptions{
		PageRange:               &HttpRange{offset, count},
		TransactionalContentMD5: &contentMD5,
	}

	putResp, err := pbClient.UploadPages(context.Background(), readerToBody, &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(putResp.RawResponse.StatusCode, 201)
	_assert.NotNil(putResp.LastModified)
	_assert.Equal((*putResp.LastModified).IsZero(), false)
	_assert.NotNil(putResp.ETag)
	_assert.NotNil(putResp.ContentMD5)
	_assert.EqualValues(*putResp.ContentMD5, contentMD5)
	_assert.Equal(*putResp.BlobSequenceNumber, int64(0))
	_assert.NotNil(*putResp.RequestID)
	_assert.NotNil(*putResp.Version)
	_assert.NotNil(putResp.Date)
	_assert.Equal((*putResp.Date).IsZero(), false)

	// put page with bad MD5
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	basContentMD5 := badMD5[:]
	uploadPagesOptions = UploadPagesOptions{
		PageRange:               &HttpRange{offset, count},
		TransactionalContentMD5: &basContentMD5,
	}
	putResp, err = pbClient.UploadPages(context.Background(), readerToBody, &uploadPagesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobCreatePageSizeInvalid() {
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
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
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
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
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
		Metadata:           &basicMetadata,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.Nil(err)

	resp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(resp.Metadata)
	_assert.EqualValues(resp.Metadata, basicMetadata)
}

//func (s *azblobTestSuite) TestBlobCreatePageMetadataEmpty() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		_assert.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	blobName := generateBlobName(testName)
//	pbClient := getPageBlobClient(blobName, containerClient)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &map[string]string{},
//	}
//	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert.NotNil(resp.Metadata)
//}

func (s *azblobTestSuite) TestBlobCreatePageMetadataInvalid() {
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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &map[string]string{"In valid1": "bar"},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.NotNil(err)
	_assert.Contains(err.Error(), invalidHeaderErrorSubstring)

}

func (s *azblobTestSuite) TestBlobCreatePageHTTPHeaders() {
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
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		BlobHTTPHeaders:    &basicHeaders,
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

//func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	_assert.Nil(err)
//
//	validatePageBlobPut(c, pbClient)
//}

//func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	currentTime := getRelativeTimeGMT(10)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}

//func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	currentTime := getRelativeTimeGMT(10)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	_assert.Nil(err)
//
//	validatePageBlobPut(c, pbClient)
//}

//func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}

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
		Metadata:           &basicMetadata,
		BlobHTTPHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
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
		Metadata:           &basicMetadata,
		BlobHTTPHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
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
		Metadata:           &basicMetadata,
		BlobHTTPHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
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
		Metadata:           &basicMetadata,
		BlobHTTPHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

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
	r := getReaderToRandomBytes(contentSize)
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
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.NotNil(err)
}

func (s *azblobTestSuite) TestBlobPutPagesNonExistentBlob() {
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
	pbClient := getPageBlobClient(blobName, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
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
	rawStart, rawEnd := (*pageListResp)[0].Raw()
	_assert.Equal(rawStart, start)
	_assert.Equal(rawEnd, end)
}

//func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	validateUploadPages(c, pbClient)
//}

//func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}

//func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	validateUploadPages(c, pbClient)
//}

//func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}

//func (s *azblobTestSuite) TestBlobPutPagesIfMatchTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	validateUploadPages(c, pbClient)
//}

//func (s *azblobTestSuite) TestBlobPutPagesIfMatchFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	eTag := "garbage"
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	eTag := "garbage"
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberLessThan := int64(10)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberLessThan := int64(1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberLessThanOrEqualTo := int64(-1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTETrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(1)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberLessThanOrEqualTo := int64(1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberLessThanOrEqualTo := int64(1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTENegOne() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberLessThanOrEqualTo := int64(-1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	//validateStorageError(c, err, )
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualTrue() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(1)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberEqualTo := int64(1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualFalse() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberEqualTo := int64(1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	ifSequenceNumberEqualTo := int64(-1)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions) // This will cause the library to set the value of the header to 0
//	_assert.NotNil(err)
//
//	//validateStorageError(c, err, )
//}
//
//func setupClearPagesTest() (ContainerClient, PageBlobClient) {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	return containerClient, pbClient
//}
//
//func validateClearPagesTest(, pbClient PageBlobClient) {
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	_assert.Nil(err)
//	pageListResp := resp.PageList.PageRange
//	_assert.(pageListResp, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesInvalidRange() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes + 1}, nil)
//	_assert.(err, chk.Not(chk.IsNil))
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceTrue() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.Nil(err)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceFalse() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceTrue() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.Nil(err)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceFalse() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfMatchTrue() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.Nil(err)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfMatchFalse() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	eTag := "garbage"
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchTrue() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	eTag := "garbage"
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.Nil(err)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchFalse() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	clearPageOptions := ClearPagesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	ifSequenceNumberLessThan := int64(10)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.Nil(err)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	ifSequenceNumberLessThan := int64(1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	ifSequenceNumberLessThan := int64(-1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTETrue() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	ifSequenceNumberLessThanOrEqualTo := int64(10)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.Nil(err)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTEFalse() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	ifSequenceNumberLessThanOrEqualTo := int64(1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTENegOne() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	ifSequenceNumberLessThanOrEqualTo := int64(-1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualTrue() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	ifSequenceNumberEqualTo := int64(10)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.Nil(err)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualFalse() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	_assert.Nil(err)
//
//	ifSequenceNumberEqualTo := int64(1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne() {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	ifSequenceNumberEqualTo := int64(-1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func setupGetPageRangesTest() (containerClient ContainerClient, pbClient PageBlobClient) {
//	svcClient := getServiceClient(nil)
//	containerClient, _ = createNewContainer(c, svcClient)
//	pbClient, _ = createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	return
//}
//
//func validateBasicGetPageRanges(, resp *PageList, err error) {
//	_assert.Nil(err)
//	_assert.(resp.PageRange, chk.NotNil)
//	_assert.(*resp.PageRange, chk.HasLen, 1)
//	start, end := int64(0), int64(PageBlobPageBytes-1)
//	rawStart, rawEnd := (*resp.PageRange)[0].Raw()
//	_assert.(rawStart, chk.Equals, start)
//	_assert.(rawEnd, chk.Equals, end)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesEmptyBlob() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	_assert.Nil(err)
//	_assert.(resp.PageList.PageRange, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesEmptyRange() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	_assert.Nil(err)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesInvalidRange() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	_, err := pbClient.GetPageRanges(ctx, HttpRange{-2, 500}, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesNonContiguousRanges() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(2*PageBlobPageBytes), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_assert.Nil(err)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	_assert.Nil(err)
//	pageListResp := resp.PageList.PageRange
//	_assert.(pageListResp, chk.NotNil)
//	_assert.(*pageListResp, chk.HasLen, 2)
//
//	start, end := int64(0), int64(PageBlobPageBytes-1)
//	rawStart, rawEnd := (*pageListResp)[0].Raw()
//	_assert.(rawStart, chk.Equals, start)
//	_assert.(rawEnd, chk.Equals, end)
//
//	start, end = int64(PageBlobPageBytes*2), int64((PageBlobPageBytes*3)-1)
//	rawStart, rawEnd = (*pageListResp)[1].Raw()
//	_assert.(rawStart, chk.Equals, start)
//	_assert.(rawEnd, chk.Equals, end)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesNotPageAligned() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 2000}, nil)
//	_assert.Nil(err)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesSnapshot() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, err := pbClient.CreateSnapshot(ctx, nil)
//	_assert.Nil(err)
//	_assert.(resp.Snapshot, chk.NotNil)
//
//	snapshotURL := pbClient.WithSnapshot(*resp.Snapshot)
//	resp2, err := snapshotURL.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	_assert.Nil(err)
//	validateBasicGetPageRanges(c, resp2.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceTrue() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.Nil(err)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceFalse() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.NotNil(err)
//
//	//serr := err.(StorageError)
//	//_assert.(serr.RawResponse.StatusCode, chk.Equals, 304)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.Nil(err)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchTrue() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	_assert.Nil(err)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	resp2, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.Nil(err)
//	validateBasicGetPageRanges(c, resp2.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchFalse() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	eTag := "garbage"
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchTrue() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	eTag := "garbage"
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: &eTag,
//			},
//		},
//	}
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.Nil(err)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchFalse() {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(_assert, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	_assert.NotNil(err)
//	//serr := err.(StorageError)
//	//_assert.(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
//}
//

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

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)
	snapshot = *resp.Snapshot

	r = getReaderToRandomBytes(PageBlobPageBytes)
	offset, count = int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions = UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)
	return
}

func validateDiffPageRanges(_assert *assert.Assertions, resp *PageList, err error) {
	_assert.Nil(err)
	pageListResp := resp.PageRange
	_assert.NotNil(pageListResp)
	_assert.Len(*resp.PageRange, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := (*pageListResp)[0].Raw()
	_assert.EqualValues(rawStart, start)
	_assert.EqualValues(rawEnd, end)
}

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

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeInvalidRange() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{-22, 14}, snapshot, nil)
	_assert.Nil(err)
}

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp.PageList, err)
}

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
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

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp.PageList, err)
}

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp2.PageList, err)
}

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	_assert.Nil(err)
	validateDiffPageRanges(_assert, resp.PageList, err)
}

func (s *azblobUnrecordedTestSuite) TestBlobDiffPageRangeIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(_assert, testName)
	defer deleteContainer(_assert, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
		BlobAccessConditions: BlobAccessConditions{
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
	actionType := SequenceNumberAction("garbage")
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
		recover()
	}()

	sequenceNumber := int64(-1)
	actionType := SequenceNumberActionUpdate
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

	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
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

	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
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

	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
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

	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
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

	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
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
	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
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
	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
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

	actionType := SequenceNumberActionIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func setupStartIncrementalCopyTest(_assert *assert.Assertions, testName string) (containerClient ContainerClient,
	pbClient PageBlobClient, copyPBClient PageBlobClient, snapshot string) {
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
	defer deleteContainer(_assert, containerClient)

	accessType := PublicAccessBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	_assert.Nil(err)

	pbClient = createNewPageBlob(_assert, generateBlobName(testName), containerClient)
	resp, _ := pbClient.CreateSnapshot(ctx, nil)

	copyPBClient = getPageBlobClient("copy"+generateBlobName(testName), containerClient)

	// Must create the incremental copy pbClient so that the access conditions work on it
	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), *resp.Snapshot, nil)
	_assert.Nil(err)
	waitForIncrementalCopy(_assert, copyPBClient, &resp2)

	resp, _ = pbClient.CreateSnapshot(ctx, nil) // Take a new snapshot so the next copy will succeed
	snapshot = *resp.Snapshot
	return
}

func validateIncrementalCopy(_assert *assert.Assertions, copyPBClient PageBlobClient, resp *PageBlobCopyIncrementalResponse) {
	t := waitForIncrementalCopy(_assert, copyPBClient, resp)

	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
	copySnapshotURL := copyPBClient.WithSnapshot(*t)
	_, err := copySnapshotURL.GetProperties(ctx, nil)
	_assert.Nil(err)
}

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
