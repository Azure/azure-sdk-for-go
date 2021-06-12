// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

//
//import (
//	"bytes"
//	"context"
//	"crypto/md5"
//	chk "gopkg.in/check.v1"
//	"io/ioutil"
//	"strings"
//	"time"
//)
//
//func (s *azblobTestSuite) TestPutGetPages(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	contentSize := 1024
//	offset, end, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	putResp, err := pbClient.UploadPages(context.Background(), getReaderToRandomBytes(1024), &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//	c.Assert(putResp.RawResponse.StatusCode, chk.Equals, 201)
//	c.Assert(putResp.LastModified, chk.NotNil)
//	c.Assert((*putResp.LastModified).IsZero(), chk.Equals, false)
//	c.Assert(putResp.ETag, chk.NotNil)
//	c.Assert(putResp.ContentMD5, chk.IsNil)
//	c.Assert(*putResp.BlobSequenceNumber, chk.Equals, int64(0))
//	c.Assert(*putResp.RequestID, chk.NotNil)
//	c.Assert(*putResp.Version, chk.NotNil)
//	c.Assert(putResp.Date, chk.NotNil)
//	c.Assert((*putResp.Date).IsZero(), chk.Equals, false)
//
//	pageList, err := pbClient.GetPageRanges(context.Background(), HttpRange{0, 1023}, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(pageList.RawResponse.StatusCode, chk.Equals, 200)
//	c.Assert(pageList.LastModified, chk.NotNil)
//	c.Assert((*pageList.LastModified).IsZero(), chk.Equals, false)
//	c.Assert(pageList.ETag, chk.NotNil)
//	c.Assert(*pageList.BlobContentLength, chk.Equals, int64(512*10))
//	c.Assert(*pageList.RequestID, chk.NotNil)
//	c.Assert(*pageList.Version, chk.NotNil)
//	c.Assert(pageList.Date, chk.NotNil)
//	c.Assert((*pageList.Date).IsZero(), chk.Equals, false)
//	c.Assert(pageList.PageList, chk.NotNil)
//	pageRangeResp := pageList.PageList.PageRange
//	c.Assert(*pageRangeResp, chk.HasLen, 1)
//	rawStart, rawEnd := (*pageRangeResp)[0].Raw()
//	c.Assert(rawStart, chk.Equals, offset)
//	c.Assert(rawEnd, chk.Equals, end)
//}
//
//func (s *azblobTestSuite) TestUploadPagesFromURL(c *chk.C) {
//	bsu := getServiceClient(nil)
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadSrcResp1.RawResponse.StatusCode, chk.Equals, 201)
//	c.Assert(uploadSrcResp1.LastModified, chk.NotNil)
//	c.Assert((*uploadSrcResp1.LastModified).IsZero(), chk.Equals, false)
//	c.Assert(uploadSrcResp1.ETag, chk.NotNil)
//	c.Assert(uploadSrcResp1.ContentMD5, chk.IsNil)
//	c.Assert(*uploadSrcResp1.BlobSequenceNumber, chk.Equals, int64(0))
//	c.Assert(*uploadSrcResp1.RequestID, chk.NotNil)
//	c.Assert(*uploadSrcResp1.Version, chk.NotNil)
//	c.Assert(uploadSrcResp1.Date, chk.NotNil)
//	c.Assert((*uploadSrcResp1.Date).IsZero(), chk.Equals, false)
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
//	c.Assert(err, chk.IsNil)
//	c.Assert(pResp1.RawResponse.StatusCode, chk.Equals, 201)
//	c.Assert(pResp1.ETag, chk.NotNil)
//	c.Assert(pResp1.LastModified, chk.NotNil)
//	c.Assert(pResp1.ContentMD5, chk.NotNil)
//	c.Assert(pResp1.RequestID, chk.NotNil)
//	c.Assert(pResp1.Version, chk.NotNil)
//	c.Assert(pResp1.Date, chk.NotNil)
//	c.Assert((*pResp1.Date).IsZero(), chk.Equals, false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, sourceData)
//}
//
//func (s *azblobTestSuite) TestUploadPagesFromURLWithMD5(c *chk.C) {
//	bsu := getServiceClient(nil)
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadSrcResp1.RawResponse.StatusCode, chk.Equals, 201)
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
//	c.Assert(err, chk.IsNil)
//	c.Assert(pResp1.RawResponse.StatusCode, chk.Equals, 201)
//	c.Assert(pResp1.ETag, chk.NotNil)
//	c.Assert(pResp1.LastModified, chk.NotNil)
//	c.Assert(pResp1.ContentMD5, chk.NotNil)
//	c.Assert(*pResp1.ContentMD5, chk.DeepEquals, contentMD5)
//	c.Assert(pResp1.RequestID, chk.NotNil)
//	c.Assert(pResp1.Version, chk.NotNil)
//	c.Assert(pResp1.Date, chk.NotNil)
//	c.Assert((*pResp1.Date).IsZero(), chk.Equals, false)
//	c.Assert(*pResp1.BlobSequenceNumber, chk.Equals, int64(0))
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, sourceData)
//
//	// Upload page from URL with bad MD5
//	_, badMD5 := getRandomDataAndReader(16)
//	badContentMD5 := badMD5[:]
//	uploadPagesFromURLOptions = UploadPagesFromURLOptions{
//		SourceContentMD5: &badContentMD5,
//	}
//	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeMD5Mismatch)
//}
//
//func (s *azblobTestSuite) TestClearDiffPages(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	contentSize := 2 * 1024
//	r := getReaderToRandomBytes(contentSize)
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	_, err := pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//
//	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
//	c.Assert(err, chk.IsNil)
//
//	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
//	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}}
//	_, err = pbClient.UploadPages(context.Background(), getReaderToRandomBytes(2048), &uploadPagesOptions1)
//	c.Assert(err, chk.IsNil)
//
//	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
//	c.Assert(err, chk.IsNil)
//	pageRangeResp := pageListResp.PageList.PageRange
//	c.Assert(pageRangeResp, chk.NotNil)
//	c.Assert(*pageRangeResp, chk.HasLen, 1)
//	// c.Assert((*pageRangeResp)[0], chk.DeepEquals, PageRange{Start: &offset1, End: &end1})
//	rawStart, rawEnd := (*pageRangeResp)[0].Raw()
//	c.Assert(rawStart, chk.Equals, offset1)
//	c.Assert(rawEnd, chk.Equals, end1)
//
//	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(clearResp.RawResponse.StatusCode, chk.Equals, 201)
//
//	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(pageListResp.PageList.PageRange, chk.IsNil)
//}
//
//func waitForIncrementalCopy(c *chk.C, copyBlobClient PageBlobClient, blobCopyResponse *PageBlobCopyIncrementalResponse) *string {
//	status := *blobCopyResponse.CopyStatus
//	var getPropertiesAndMetadataResult BlobGetPropertiesResponse
//	// Wait for the copy to finish
//	start := time.Now()
//	for status != CopyStatusSuccess {
//		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(ctx, nil)
//		status = *getPropertiesAndMetadataResult.CopyStatus
//		currentTime := time.Now()
//		if currentTime.Sub(start) >= time.Minute {
//			c.Fail()
//		}
//	}
//	return getPropertiesAndMetadataResult.DestinationSnapshot
//}
//
//func (s *azblobTestSuite) TestIncrementalCopy(c *chk.C) {
//	bsu := getServiceClient(nil)
//
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	accessType := PublicAccessBlob
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
//	}
//	_, err := containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	c.Assert(err, chk.IsNil)
//
//	srcBlob, _ := createNewPageBlob(c, containerClient)
//	contentSize := 1024
//	r := getReaderToRandomBytes(contentSize)
//	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	_, err = srcBlob.UploadPages(context.Background(), r, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//
//	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
//	c.Assert(err, chk.IsNil)
//
//	dstBlob := containerClient.NewPageBlobClient(generateBlobName())
//
//	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
//	c.Assert(resp.LastModified, chk.NotNil)
//	c.Assert((*resp.LastModified).IsZero(), chk.Equals, false)
//	c.Assert(resp.ETag, chk.NotNil)
//	c.Assert(*resp.RequestID, chk.Not(chk.Equals), "")
//	c.Assert(*resp.Version, chk.Not(chk.Equals), "")
//	c.Assert(resp.Date, chk.NotNil)
//	c.Assert((*resp.Date).IsZero(), chk.Equals, false)
//	c.Assert(*resp.CopyID, chk.Not(chk.Equals), "")
//	c.Assert(*resp.CopyStatus, chk.Equals, CopyStatusPending)
//
//	waitForIncrementalCopy(c, dstBlob, &resp)
//}
//
//func (s *azblobTestSuite) TestResizePageBlob(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	pbClient, _ := createNewPageBlob(c, containerClient)
//	resp, err := pbClient.Resize(context.Background(), 2048, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)
//
//	resp, err = pbClient.Resize(context.Background(), 8192, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)
//
//	resp2, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(*resp2.ContentLength, chk.Equals, int64(8192))
//}
//
//func (s *azblobTestSuite) TestPageSequenceNumbers(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	defer deleteContainer(containerClient)
//
//	sequenceNumber := int64(0)
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	resp, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)
//
//	sequenceNumber = int64(7)
//	actionType = SequenceNumberActionMax
//	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)
//
//	sequenceNumber = int64(11)
//	actionType = SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)
//}
//
//func (s *azblobTestSuite) TestPutPagesWithMD5(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	// put page with valid MD5
//	contentSize := 1024
//	readerToBody, body := getRandomDataAndReader(contentSize)
//	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
//	md5Value := md5.Sum(body)
//	contentMD5 := md5Value[:]
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange:               &HttpRange{offset, count},
//		TransactionalContentMD5: &contentMD5,
//	}
//
//	putResp, err := pbClient.UploadPages(context.Background(), readerToBody, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//	c.Assert(putResp.RawResponse.StatusCode, chk.Equals, 201)
//	c.Assert(putResp.LastModified, chk.NotNil)
//	c.Assert((*putResp.LastModified).IsZero(), chk.Equals, false)
//	c.Assert(putResp.ETag, chk.NotNil)
//	c.Assert(putResp.ContentMD5, chk.NotNil)
//	c.Assert(*putResp.ContentMD5, chk.DeepEquals, contentMD5)
//	c.Assert(*putResp.BlobSequenceNumber, chk.Equals, int64(0))
//	c.Assert(*putResp.RequestID, chk.NotNil)
//	c.Assert(*putResp.Version, chk.NotNil)
//	c.Assert(putResp.Date, chk.NotNil)
//	c.Assert((*putResp.Date).IsZero(), chk.Equals, false)
//
//	// put page with bad MD5
//	readerToBody, body = getRandomDataAndReader(1024)
//	_, badMD5 := getRandomDataAndReader(16)
//	basContentMD5 := badMD5[:]
//	uploadPagesOptions = UploadPagesOptions{
//		PageRange:               &HttpRange{offset, count},
//		TransactionalContentMD5: &basContentMD5,
//	}
//	putResp, err = pbClient.UploadPages(context.Background(), readerToBody, &uploadPagesOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeMD5Mismatch)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageSizeInvalid(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := getPageBlobClient(c, containerClient)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//	}
//	_, err := pbClient.Create(ctx, 1, &createPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageSequenceInvalid(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := getPageBlobClient(c, containerClient)
//
//	sequenceNumber := int64(-1)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageMetadataNonEmpty(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := getPageBlobClient(c, containerClient)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Metadata, chk.NotNil)
//	c.Assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageMetadataEmpty(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := getPageBlobClient(c, containerClient)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &map[string]string{},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Metadata, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageMetadataInvalid(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := getPageBlobClient(c, containerClient)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &map[string]string{"In valid1": "bar"},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
//
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageHTTPHeaders(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := getPageBlobClient(c, containerClient)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		BlobHTTPHeaders:    &basicHeaders,
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	h := resp.NewHTTPHeaders()
//	c.Assert(h, chk.DeepEquals, basicHeaders)
//}
//
//func validatePageBlobPut(c *chk.C, pbClient PageBlobClient) {
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Metadata, chk.NotNil)
//	c.Assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//	c.Assert(resp.NewHTTPHeaders(), chk.DeepEquals, basicHeaders)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validatePageBlobPut(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfModifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validatePageBlobPut(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfUnmodifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validatePageBlobPut(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	sequenceNumber := int64(0)
//	eTag := "garbage"
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	sequenceNumber := int64(0)
//	eTag := "garbage"
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validatePageBlobPut(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobCreatePageIfNoneMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	sequenceNumber := int64(0)
//	createPageBlobOptions := CreatePageBlobOptions{
//		BlobSequenceNumber: &sequenceNumber,
//		Metadata:           &basicMetadata,
//		BlobHTTPHeaders:    &basicHeaders,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesInvalidRange(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	contentSize := 1024
//	r := getReaderToRandomBytes(contentSize)
//	offset, count := int64(0), int64(contentSize/2)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.Not(chk.IsNil))
//}
//
//// Body cannot be nil check already added in the request preparer
////func (s *azblobTestSuite) TestBlobPutPagesNilBody(c *chk.C) {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	pbClient, _ := createNewPageBlob(c, containerClient)
////
////	_, err := pbClient.UploadPages(ctx, nil, nil)
////	c.Assert(err, chk.NotNil)
////}
//
//func (s *azblobTestSuite) TestBlobPutPagesEmptyBody(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := bytes.NewReader([]byte{})
//	offset, count := int64(0), int64(0)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.NotNil)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesNonExistentBlob(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := getPageBlobClient(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeBlobNotFound)
//}
//
//func validateUploadPages(c *chk.C, pbClient PageBlobClient) {
//	// This will only validate a single put page at 0-PageBlobPageBytes-1
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
//	c.Assert(err, chk.IsNil)
//	pageListResp := resp.PageList.PageRange
//	start, end := int64(0), int64(PageBlobPageBytes-1)
//	rawStart, rawEnd := (*pageListResp)[0].Raw()
//	c.Assert(rawStart, chk.Equals, start)
//	c.Assert(rawEnd, chk.Equals, end)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfModifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfUnmodifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfNoneMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTETrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(1)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
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
//	c.Assert(err, chk.IsNil)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberLTENegOne(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	//validateStorageError(c, err, )
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(1)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
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
//	c.Assert(err, chk.IsNil)
//
//	validateUploadPages(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	//validateStorageError(c, err, )
//}
//
//func setupClearPagesTest(c *chk.C) (ContainerClient, PageBlobClient) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//
//	return containerClient, pbClient
//}
//
//func validateClearPagesTest(c *chk.C, pbClient PageBlobClient) {
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	c.Assert(err, chk.IsNil)
//	pageListResp := resp.PageList.PageRange
//	c.Assert(pageListResp, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesInvalidRange(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes + 1}, nil)
//	c.Assert(err, chk.Not(chk.IsNil))
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfModifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfUnmodifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfMatchTrue(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfMatchFalse(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchTrue(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfNoneMatchFalse(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	ifSequenceNumberLessThan := int64(10)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	ifSequenceNumberLessThan := int64(1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	ifSequenceNumberLessThan := int64(-1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTETrue(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	ifSequenceNumberLessThanOrEqualTo := int64(10)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTEFalse(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	ifSequenceNumberLessThanOrEqualTo := int64(1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberLTENegOne(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	ifSequenceNumberLessThanOrEqualTo := int64(-1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualTrue(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	ifSequenceNumberEqualTo := int64(10)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateClearPagesTest(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualFalse(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	sequenceNumber := int64(10)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	ifSequenceNumberEqualTo := int64(1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeSequenceNumberConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne(c *chk.C) {
//	containerClient, pbClient := setupClearPagesTest(c)
//	defer deleteContainer(containerClient)
//
//	ifSequenceNumberEqualTo := int64(-1)
//	clearPageOptions := ClearPagesOptions{
//		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidInput)
//}
//
//func setupGetPageRangesTest(c *chk.C) (containerClient ContainerClient, pbClient PageBlobClient) {
//	bsu := getServiceClient(nil)
//	containerClient, _ = createNewContainer(c, bsu)
//	pbClient, _ = createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//
//	return
//}
//
//func validateBasicGetPageRanges(c *chk.C, resp *PageList, err error) {
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.PageRange, chk.NotNil)
//	c.Assert(*resp.PageRange, chk.HasLen, 1)
//	start, end := int64(0), int64(PageBlobPageBytes-1)
//	rawStart, rawEnd := (*resp.PageRange)[0].Raw()
//	c.Assert(rawStart, chk.Equals, start)
//	c.Assert(rawEnd, chk.Equals, end)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesEmptyBlob(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.PageList.PageRange, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesEmptyRange(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	c.Assert(err, chk.IsNil)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesInvalidRange(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	_, err := pbClient.GetPageRanges(ctx, HttpRange{-2, 500}, nil)
//	c.Assert(err, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesNonContiguousRanges(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(2*PageBlobPageBytes), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	c.Assert(err, chk.IsNil)
//	pageListResp := resp.PageList.PageRange
//	c.Assert(pageListResp, chk.NotNil)
//	c.Assert(*pageListResp, chk.HasLen, 2)
//
//	start, end := int64(0), int64(PageBlobPageBytes-1)
//	rawStart, rawEnd := (*pageListResp)[0].Raw()
//	c.Assert(rawStart, chk.Equals, start)
//	c.Assert(rawEnd, chk.Equals, end)
//
//	start, end = int64(PageBlobPageBytes*2), int64((PageBlobPageBytes*3)-1)
//	rawStart, rawEnd = (*pageListResp)[1].Raw()
//	c.Assert(rawStart, chk.Equals, start)
//	c.Assert(rawEnd, chk.Equals, end)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesNotPageAligned(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 2000}, nil)
//	c.Assert(err, chk.IsNil)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesSnapshot(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	resp, err := pbClient.CreateSnapshot(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Snapshot, chk.NotNil)
//
//	snapshotURL := pbClient.WithSnapshot(*resp.Snapshot)
//	resp2, err := snapshotURL.GetPageRanges(ctx, HttpRange{0, 0}, nil)
//	c.Assert(err, chk.IsNil)
//	validateBasicGetPageRanges(c, resp2.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfModifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	//serr := err.(StorageError)
//	//c.Assert(serr.RawResponse.StatusCode, chk.Equals, 304)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchTrue(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	resp2, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
//	c.Assert(err, chk.IsNil)
//	validateBasicGetPageRanges(c, resp2.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfMatchFalse(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchTrue(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.IsNil)
//	validateBasicGetPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobGetPageRangesIfNoneMatchFalse(c *chk.C) {
//	containerClient, pbClient := setupGetPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	c.Assert(err, chk.NotNil)
//	//serr := err.(StorageError)
//	//c.Assert(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
//}
//
//func setupDiffPageRangesTest(c *chk.C) (containerClient ContainerClient, pbClient PageBlobClient, snapshot string) {
//	bsu := getServiceClient(nil)
//	containerClient, _ = createNewContainer(c, bsu)
//	pbClient, _ = createNewPageBlob(c, containerClient)
//
//	r := getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count := int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := pbClient.CreateSnapshot(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	snapshot = *resp.Snapshot
//
//	r = getReaderToRandomBytes(PageBlobPageBytes)
//	offset, count = int64(0), int64(PageBlobPageBytes)
//	uploadPagesOptions = UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	c.Assert(err, chk.IsNil)
//	return
//}
//
//func validateDiffPageRanges(c *chk.C, resp *PageList, err error) {
//	c.Assert(err, chk.IsNil)
//	pageListResp := resp.PageRange
//	c.Assert(pageListResp, chk.NotNil)
//	c.Assert(*resp.PageRange, chk.HasLen, 1)
//	start, end := int64(0), int64(PageBlobPageBytes-1)
//	rawStart, rawEnd := (*pageListResp)[0].Raw()
//	c.Assert(rawStart, chk.DeepEquals, start)
//	c.Assert(rawEnd, chk.DeepEquals, end)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangesNonExistentSnapshot(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	snapshotTime, _ := time.Parse(SnapshotTimeFormat, snapshot)
//	snapshotTime = snapshotTime.Add(time.Minute)
//	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshotTime.Format(SnapshotTimeFormat), nil)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodePreviousSnapshotNotFound)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeInvalidRange(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
//	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{-22, 14}, snapshot, nil)
//	c.Assert(err, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfModifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.IsNil)
//	validateDiffPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfModifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.NotNil)
//
//	//stgErr := err.(StorageError)
//	//c.Assert(stgErr.Response().StatusCode, chk.Equals, 304)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.IsNil)
//	validateDiffPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfMatchTrue(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	resp2, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.IsNil)
//	validateDiffPageRanges(c, resp2.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfMatchFalse(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	eTag := "garbage"
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfNoneMatchTrue(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
//
//	eTag := "garbage"
//	getPageRangesOptions := GetPageRangesOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: &eTag,
//			},
//		},
//	}
//	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.IsNil)
//	validateDiffPageRanges(c, resp.PageList, err)
//}
//
//func (s *azblobTestSuite) TestBlobDiffPageRangeIfNoneMatchFalse(c *chk.C) {
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
//	defer deleteContainer(containerClient)
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
//	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
//	c.Assert(err, chk.NotNil)
//
//	//serr := err.(StorageError)
//	//c.Assert(serr.Response().StatusCode, chk.Equals, 304)
//}
//
//func (s *azblobTestSuite) TestBlobResizeZero(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	// The default pbClient is created with size > 0, so this should actually update
//	_, err := pbClient.Resize(ctx, 0, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(*resp.ContentLength, chk.Equals, int64(0))
//}
//
//func (s *azblobTestSuite) TestBlobResizeInvalidSizeNegative(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	_, err := pbClient.Resize(ctx, -4, nil)
//	c.Assert(err, chk.NotNil)
//}
//
//func (s *azblobTestSuite) TestBlobResizeInvalidSizeMisaligned(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	_, err := pbClient.Resize(ctx, 12, nil)
//	c.Assert(err, chk.NotNil)
//}
//
//func validateResize(c *chk.C, pbClient PageBlobClient) {
//	resp, _ := pbClient.GetProperties(ctx, nil)
//	c.Assert(*resp.ContentLength, chk.Equals, int64(PageBlobPageBytes))
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateResize(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfModifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateResize(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfUnmodifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateResize(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	eTag := "garbage"
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfNoneMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	eTag := "garbage"
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateResize(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobResizeIfNoneMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	resizePageBlobOptions := ResizePageBlobOptions{
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberActionTypeInvalid(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	sequenceNumber := int64(1)
//	actionType := SequenceNumberAction("garbage")
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberSequenceNumberInvalid(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	defer func() { // Invalid sequence number should panic
//		recover()
//	}()
//
//	sequenceNumber := int64(-1)
//	actionType := SequenceNumberActionUpdate
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		BlobSequenceNumber: &sequenceNumber,
//		ActionType:         &actionType,
//	}
//
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
//}
//
//func validateSequenceNumberSet(c *chk.C, pbClient PageBlobClient) {
//	resp, err := pbClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(*resp.BlobSequenceNumber, chk.Equals, int64(1))
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	validateSequenceNumberSet(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfModifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfModifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	validateSequenceNumberSet(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfUnmodifiedSince: &currentTime,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	validateSequenceNumberSet(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	eTag := "garbage"
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchTrue(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	eTag := "garbage"
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: &eTag,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.IsNil)
//
//	validateSequenceNumberSet(c, pbClient)
//}
//
//func (s *azblobTestSuite) TestBlobSetSequenceNumberIfNoneMatchFalse(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, _ := pbClient.GetProperties(ctx, nil)
//
//	actionType := SequenceNumberActionIncrement
//	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
//		ActionType: &actionType,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{
//				IfNoneMatch: resp.ETag,
//			},
//		},
//	}
//	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func setupStartIncrementalCopyTest(c *chk.C) (containerClient ContainerClient, pbClient PageBlobClient, copyPBClient PageBlobClient, snapshot string) {
//	bsu := getServiceClient(nil)
//	containerClient, _ = createNewContainer(c, bsu)
//
//	accessType := PublicAccessBlob
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
//	}
//	_, err := containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	c.Assert(err, chk.IsNil)
//
//	pbClient, _ = createNewPageBlob(c, containerClient)
//	resp, _ := pbClient.CreateSnapshot(ctx, nil)
//	copyPBClient, _ = getPageBlobClient(c, containerClient)
//
//	// Must create the incremental copy pbClient so that the access conditions work on it
//	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), *resp.Snapshot, nil)
//	c.Assert(err, chk.IsNil)
//	waitForIncrementalCopy(c, copyPBClient, &resp2)
//
//	resp, _ = pbClient.CreateSnapshot(ctx, nil) // Take a new snapshot so the next copy will succeed
//	snapshot = *resp.Snapshot
//	return
//}
//
//func validateIncrementalCopy(c *chk.C, copyBlobURL PageBlobClient, resp *PageBlobCopyIncrementalResponse) {
//	t := waitForIncrementalCopy(c, copyBlobURL, resp)
//
//	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
//	copySnapshotURL := copyBlobURL.WithSnapshot(*t)
//	_, err := copySnapshotURL.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopySnapshotNotExist(c *chk.C) {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//	copyBlobURL, _ := getPageBlobClient(c, containerClient)
//
//	snapshot := time.Now().UTC().Format(SnapshotTimeFormat)
//	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, nil)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeCannotVerifyCopySource)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateIncrementalCopy(c, copyBlobURL, &resp)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	currentTime := getRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateIncrementalCopy(c, copyBlobURL, &resp)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	currentTime := getRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfMatchTrue(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	resp, _ := copyBlobURL.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfMatch: resp.ETag,
//		},
//	}
//	resp2, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateIncrementalCopy(c, copyBlobURL, &resp2)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfMatchFalse(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfMatch: &eTag,
//		},
//	}
//	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeTargetConditionNotMet)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfNoneMatch: &eTag,
//		},
//	}
//	resp, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateIncrementalCopy(c, copyBlobURL, &resp)
//}
//
//func (s *azblobTestSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse(c *chk.C) {
//	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)
//
//	defer deleteContainer(containerClient)
//
//	resp, _ := copyBlobURL.GetProperties(ctx, nil)
//
//	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfNoneMatch: resp.ETag,
//		},
//	}
//	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	c.Assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
