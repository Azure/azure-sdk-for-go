package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	chk "gopkg.in/check.v1"
	"io/ioutil"
	"strings"
	"time"
)

func (s *aztestsSuite) TestPutGetPages(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	pbClient, _ := createNewPageBlob(c, containerClient)

	contentSize := 1024
	offset, end, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	putResp, err := pbClient.UploadPages(context.Background(), getReaderToRandomBytes(1024), &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(putResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(putResp.LastModified, chk.NotNil)
	c.Assert((*putResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(putResp.ETag, chk.NotNil)
	c.Assert(putResp.ContentMD5, chk.IsNil)
	c.Assert(*putResp.BlobSequenceNumber, chk.Equals, int64(0))
	c.Assert(*putResp.RequestID, chk.NotNil)
	c.Assert(*putResp.Version, chk.NotNil)
	c.Assert(putResp.Date, chk.NotNil)
	c.Assert((*putResp.Date).IsZero(), chk.Equals, false)

	pageList, err := pbClient.GetPageRanges(context.Background(), 0, 1023, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(pageList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(pageList.LastModified, chk.NotNil)
	c.Assert((*pageList.LastModified).IsZero(), chk.Equals, false)
	c.Assert(pageList.ETag, chk.NotNil)
	c.Assert(*pageList.BlobContentLength, chk.Equals, int64(512*10))
	c.Assert(*pageList.RequestID, chk.NotNil)
	c.Assert(*pageList.Version, chk.NotNil)
	c.Assert(pageList.Date, chk.NotNil)
	c.Assert((*pageList.Date).IsZero(), chk.Equals, false)
	c.Assert(pageList.PageList, chk.NotNil)
	pageRangeResp := pageList.PageList.PageRange
	c.Assert(*pageRangeResp, chk.HasLen, 1)
	c.Assert((*pageRangeResp)[0], chk.DeepEquals, PageRange{Start: &offset, End: &end})
}

func (s *aztestsSuite) TestUploadPagesFromURL(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := 8 * 1024 // 8KB
	r, sourceData := getRandomDataAndReader(contentSize)
	ctx := context.Background() // Use default Background context
	srcBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
	destBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))

	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp1.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(uploadSrcResp1.LastModified, chk.NotNil)
	c.Assert((*uploadSrcResp1.LastModified).IsZero(), chk.Equals, false)
	c.Assert(uploadSrcResp1.ETag, chk.NotNil)
	c.Assert(uploadSrcResp1.ContentMD5, chk.IsNil)
	c.Assert(*uploadSrcResp1.BlobSequenceNumber, chk.Equals, int64(0))
	c.Assert(*uploadSrcResp1.RequestID, chk.NotNil)
	c.Assert(*uploadSrcResp1.Version, chk.NotNil)
	c.Assert(uploadSrcResp1.Date, chk.NotNil)
	c.Assert((*uploadSrcResp1.Date).IsZero(), chk.Equals, false)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL.
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(pResp1.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(pResp1.ETag, chk.NotNil)
	c.Assert(pResp1.LastModified, chk.NotNil)
	c.Assert(pResp1.ContentMD5, chk.NotNil)
	c.Assert(pResp1.RequestID, chk.NotNil)
	c.Assert(pResp1.Version, chk.NotNil)
	c.Assert(pResp1.Date, chk.NotNil)
	c.Assert((*pResp1.Date).IsZero(), chk.Equals, false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)
}

func (s *aztestsSuite) TestUploadPagesFromURLWithMD5(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := 8 * 1024 // 8KB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	contentMD5 := md5Value[:]
	ctx := context.Background() // Use default Background context
	srcBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
	destBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp1.RawResponse.StatusCode, chk.Equals, 201)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL with MD5.
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMd5: &contentMD5,
	}
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(pResp1.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(pResp1.ETag, chk.NotNil)
	c.Assert(pResp1.LastModified, chk.NotNil)
	c.Assert(pResp1.ContentMD5, chk.NotNil)
	c.Assert(*pResp1.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(pResp1.RequestID, chk.NotNil)
	c.Assert(pResp1.Version, chk.NotNil)
	c.Assert(pResp1.Date, chk.NotNil)
	c.Assert((*pResp1.Date).IsZero(), chk.Equals, false)
	c.Assert(*pResp1.BlobSequenceNumber, chk.Equals, int64(0))

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)

	// Upload page from URL with bad MD5
	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions = UploadPagesFromURLOptions{
		SourceContentMd5: &badContentMD5,
	}
	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeMd5Mismatch)
}

func (s *aztestsSuite) TestClearDiffPages(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	pbClient, _ := createNewPageBlob(c, containerClient)

	contentSize := 2 * 1024
	r := getReaderToRandomBytes(contentSize)
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	_, err := pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	c.Assert(err, chk.IsNil)

	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
	uploadPagesOptions1 := UploadPagesOptions{Offset: &offset1, Count: &count1}
	_, err = pbClient.UploadPages(context.Background(), getReaderToRandomBytes(2048), &uploadPagesOptions1)
	c.Assert(err, chk.IsNil)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), 0, 4096, *snapshotResp.Snapshot, nil)
	c.Assert(err, chk.IsNil)
	pageRangeResp := pageListResp.PageList.PageRange
	c.Assert(pageRangeResp, chk.NotNil)
	c.Assert(*pageRangeResp, chk.HasLen, 1)
	c.Assert((*pageRangeResp)[0], chk.DeepEquals, PageRange{Start: &offset1, End: &end1})

	clearResp, err := pbClient.ClearPages(context.Background(), 2048, 2048, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(clearResp.RawResponse.StatusCode, chk.Equals, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), 0, 4095, *snapshotResp.Snapshot, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(pageListResp.PageList.PageRange, chk.IsNil)
}

func waitForIncrementalCopy(c *chk.C, copyBlobClient PageBlobClient, blobCopyResponse *PageBlobCopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult BlobGetPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(ctx, nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			c.Fail()
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

func (s *aztestsSuite) TestIncrementalCopy(c *chk.C) {
	bsu := getBSU()

	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	accessType := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
	}
	_, err := containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	srcBlob, _ := createNewPageBlob(c, containerClient)
	contentSize := 1024
	r := getReaderToRandomBytes(contentSize)
	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	_, err = srcBlob.UploadPages(context.Background(), r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
	c.Assert(err, chk.IsNil)

	dstBlob := containerClient.NewPageBlobClient(generateBlobName())

	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert((*resp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(*resp.RequestID, chk.Not(chk.Equals), "")
	c.Assert(*resp.Version, chk.Not(chk.Equals), "")
	c.Assert(resp.Date, chk.NotNil)
	c.Assert((*resp.Date).IsZero(), chk.Equals, false)
	c.Assert(*resp.CopyID, chk.Not(chk.Equals), "")
	c.Assert(*resp.CopyStatus, chk.Equals, CopyStatusTypePending)

	waitForIncrementalCopy(c, dstBlob, &resp)
}

func (s *aztestsSuite) TestResizePageBlob(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	pbClient, _ := createNewPageBlob(c, containerClient)
	resp, err := pbClient.Resize(context.Background(), 2048, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)

	resp, err = pbClient.Resize(context.Background(), 8192, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)

	resp2, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp2.ContentLength, chk.Equals, int64(8192))
}

func (s *aztestsSuite) TestPageSequenceNumbers(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	pbClient, _ := createNewPageBlob(c, containerClient)

	defer deleteContainer(c, containerClient)

	sequenceNumber := int64(0)
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)

	sequenceNumber = int64(7)
	actionType = SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)

	sequenceNumber = int64(11)
	actionType = SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 200)
}

func (s *aztestsSuite) TestPutPagesWithMD5(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	pbClient, _ := createNewPageBlob(c, containerClient)

	// put page with valid MD5
	contentSize := 1024
	readerToBody, body := getRandomDataAndReader(contentSize)
	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
	md5Value := md5.Sum(body)
	contentMD5 := md5Value[:]
	uploadPagesOptions := UploadPagesOptions{
		Offset:                  &offset,
		Count:                   &count,
		TransactionalContentMd5: &contentMD5,
	}

	putResp, err := pbClient.UploadPages(context.Background(), readerToBody, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(putResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(putResp.LastModified, chk.NotNil)
	c.Assert((*putResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(putResp.ETag, chk.NotNil)
	c.Assert(putResp.ContentMD5, chk.NotNil)
	c.Assert(*putResp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(*putResp.BlobSequenceNumber, chk.Equals, int64(0))
	c.Assert(*putResp.RequestID, chk.NotNil)
	c.Assert(*putResp.Version, chk.NotNil)
	c.Assert(putResp.Date, chk.NotNil)
	c.Assert((*putResp.Date).IsZero(), chk.Equals, false)

	// put page with bad MD5
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	basContentMD5 := badMD5[:]
	uploadPagesOptions = UploadPagesOptions{
		Offset:                  &offset,
		Count:                   &count,
		TransactionalContentMd5: &basContentMD5,
	}
	putResp, err = pbClient.UploadPages(context.Background(), readerToBody, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeMd5Mismatch)
}

func (s *aztestsSuite) TestBlobCreatePageSizeInvalid(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := getPageBlobClient(c, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err := pbClient.Create(ctx, 1, &createPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeInvalidHeaderValue)
}

func (s *aztestsSuite) TestBlobCreatePageSequenceInvalid(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := getPageBlobClient(c, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobCreatePageMetadataNonEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := getPageBlobClient(c, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.IsNil)

	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.NotNil)
	c.Assert(*resp.Metadata, chk.DeepEquals, basicMetadata)
}

func (s *aztestsSuite) TestBlobCreatePageMetadataEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := getPageBlobClient(c, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &map[string]string{},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)

	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.IsNil)
}

func (s *aztestsSuite) TestBlobCreatePageMetadataInvalid(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := getPageBlobClient(c, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &map[string]string{"In valid1": "bar"},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.NotNil)
	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)

}

func (s *aztestsSuite) TestBlobCreatePageHTTPHeaders(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := getPageBlobClient(c, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		BlobHttpHeaders:    &basicHeaders,
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.IsNil)

	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	h := resp.NewHTTPHeaders()
	c.Assert(h, chk.DeepEquals, basicHeaders)
}

func validatePageBlobPut(c *chk.C, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.NotNil)
	c.Assert(*resp.Metadata, chk.DeepEquals, basicMetadata)
	c.Assert(resp.NewHTTPHeaders(), chk.DeepEquals, basicHeaders)
}

func (s *aztestsSuite) TestBlobCreatePageIfModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	currentTime := getRelativeTimeGMT(-10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validatePageBlobPut(c, pbClient)
}

func (s *aztestsSuite) TestBlobCreatePageIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	currentTime := getRelativeTimeGMT(10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobCreatePageIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	currentTime := getRelativeTimeGMT(10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validatePageBlobPut(c, pbClient)
}

func (s *aztestsSuite) TestBlobCreatePageIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	currentTime := getRelativeTimeGMT(-10)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobCreatePageIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validatePageBlobPut(c, pbClient)
}

func (s *aztestsSuite) TestBlobCreatePageIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobCreatePageIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validatePageBlobPut(c, pbClient)
}

func (s *aztestsSuite) TestBlobCreatePageIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient) // Originally created without metadata

	resp, _ := pbClient.GetProperties(ctx, nil)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           &basicMetadata,
		BlobHttpHeaders:    &basicHeaders,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesInvalidRange(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	contentSize := 1024
	r := getReaderToRandomBytes(contentSize)
	offset, count := int64(0), int64(contentSize/2)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.Not(chk.IsNil))
}

// Body cannot be nil check already added in the request preparer
//func (s *aztestsSuite) TestBlobPutPagesNilBody(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	pbClient, _ := createNewPageBlob(c, containerClient)
//
//	_, err := pbClient.UploadPages(ctx, nil, nil)
//	c.Assert(err, chk.NotNil)
//}

func (s *aztestsSuite) TestBlobPutPagesEmptyBody(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := bytes.NewReader([]byte{})
	offset, count := int64(0), int64(0)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobPutPagesNonExistentBlob(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := getPageBlobClient(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{Offset: &offset, Count: &count}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeBlobNotFound)
}

func validateUploadPages(c *chk.C, pbClient PageBlobClient) {
	// This will only validate a single put page at 0-PageBlobPageBytes-1
	resp, err := pbClient.GetPageRanges(ctx, 0, CountToEnd, nil)
	c.Assert(err, chk.IsNil)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), int64(PageBlobPageBytes-1)
	c.Assert((*pageListResp)[0], chk.DeepEquals, PageRange{Start: &start, End: &end})
}

func (s *aztestsSuite) TestBlobPutPagesIfModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	validateUploadPages(c, pbClient)
}

func (s *aztestsSuite) TestBlobPutPagesIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)
	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	validateUploadPages(c, pbClient)
}

func (s *aztestsSuite) TestBlobPutPagesIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	validateUploadPages(c, pbClient)
}

func (s *aztestsSuite) TestBlobPutPagesIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	eTag := "garbage"
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	eTag := "garbage"
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	validateUploadPages(c, pbClient)
}

func (s *aztestsSuite) TestBlobPutPagesIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThan := int64(10)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	validateUploadPages(c, pbClient)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThan := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeSequenceNumberConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeInvalidInput)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberLTETrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	validateUploadPages(c, pbClient)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeSequenceNumberConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberLTENegOne(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	//validateStorageError(c, err, )
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberEqualTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	validateUploadPages(c, pbClient)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberEqualFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeSequenceNumberConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	ifSequenceNumberEqualTo := int64(-1)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions) // This will cause the library to set the value of the header to 0
	c.Assert(err, chk.NotNil)

	//validateStorageError(c, err, )
}

func setupClearPagesTest(c *chk.C) (ContainerClient, PageBlobClient) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	pbClient, _ := createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	return containerClient, pbClient
}

func validateClearPagesTest(c *chk.C, pbClient PageBlobClient) {
	resp, err := pbClient.GetPageRanges(ctx, 0, 0, nil)
	c.Assert(err, chk.IsNil)
	pageListResp := resp.PageList.PageRange
	c.Assert(pageListResp, chk.IsNil)
}

func (s *aztestsSuite) TestBlobClearPagesInvalidRange(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes+1, nil)
	c.Assert(err, chk.Not(chk.IsNil))
}

func (s *aztestsSuite) TestBlobClearPagesIfModifiedSinceTrue(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.IsNil)

	validateClearPagesTest(c, pbClient)
}

func (s *aztestsSuite) TestBlobClearPagesIfModifiedSinceFalse(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobClearPagesIfUnmodifiedSinceTrue(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.IsNil)

	validateClearPagesTest(c, pbClient)
}

func (s *aztestsSuite) TestBlobClearPagesIfUnmodifiedSinceFalse(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobClearPagesIfMatchTrue(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.IsNil)

	validateClearPagesTest(c, pbClient)
}

func (s *aztestsSuite) TestBlobClearPagesIfMatchFalse(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobClearPagesIfNoneMatchTrue(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.IsNil)

	validateClearPagesTest(c, pbClient)
}

func (s *aztestsSuite) TestBlobClearPagesIfNoneMatchFalse(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.IsNil)

	validateClearPagesTest(c, pbClient)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeSequenceNumberConditionNotMet)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeInvalidInput)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberLTETrue(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.IsNil)

	validateClearPagesTest(c, pbClient)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberLTEFalse(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeSequenceNumberConditionNotMet)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberLTENegOne(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions) // This will cause the library to set the value of the header to 0
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeInvalidInput)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberEqualTrue(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.IsNil)

	validateClearPagesTest(c, pbClient)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberEqualFalse(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeSequenceNumberConditionNotMet)
}

func (s *aztestsSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne(c *chk.C) {
	containerClient, pbClient := setupClearPagesTest(c)
	defer deleteContainer(c, containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, 0, PageBlobPageBytes, &clearPageOptions) // This will cause the library to set the value of the header to 0
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeInvalidInput)
}

func setupGetPageRangesTest(c *chk.C) (containerClient ContainerClient, pbClient PageBlobClient) {
	bsu := getBSU()
	containerClient, _ = createNewContainer(c, bsu)
	pbClient, _ = createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	return
}

func validateBasicGetPageRanges(c *chk.C, resp *PageList, err error) {
	c.Assert(err, chk.IsNil)
	c.Assert(resp.PageRange, chk.NotNil)
	c.Assert(*resp.PageRange, chk.HasLen, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	c.Assert((*resp.PageRange)[0], chk.DeepEquals, PageRange{Start: &start, End: &end})
}

func (s *aztestsSuite) TestBlobGetPageRangesEmptyBlob(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, 0, 0, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.PageList.PageRange, chk.IsNil)
}

func (s *aztestsSuite) TestBlobGetPageRangesEmptyRange(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, 0, 0, nil)
	c.Assert(err, chk.IsNil)
	validateBasicGetPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobGetPageRangesInvalidRange(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	_, err := pbClient.GetPageRanges(ctx, -2, 500, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobGetPageRangesNonContiguousRanges(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(2*PageBlobPageBytes), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	resp, err := pbClient.GetPageRanges(ctx, 0, 0, nil)
	c.Assert(err, chk.IsNil)
	pageListResp := resp.PageList.PageRange
	c.Assert(pageListResp, chk.NotNil)
	c.Assert(*pageListResp, chk.HasLen, 2)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	c.Assert((*pageListResp)[0], chk.DeepEquals, PageRange{Start: &start, End: &end})
	start, end = int64(PageBlobPageBytes*2), int64((PageBlobPageBytes*3)-1)
	c.Assert((*pageListResp)[1], chk.DeepEquals, PageRange{Start: &start, End: &end})
}

func (s *aztestsSuite) TestBlobGetPageRangesNotPageAligned(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, 0, 2000, nil)
	c.Assert(err, chk.IsNil)
	validateBasicGetPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobGetPageRangesSnapshot(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Snapshot, chk.NotNil)

	snapshotURL := pbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetPageRanges(ctx, 0, 0, nil)
	c.Assert(err, chk.IsNil)
	validateBasicGetPageRanges(c, resp2.PageList, err)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfModifiedSinceTrue(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateBasicGetPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfModifiedSinceFalse(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)

	//serr := err.(StorageError)
	//c.Assert(serr.RawResponse.StatusCode, chk.Equals, 304)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateBasicGetPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfMatchTrue(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateBasicGetPageRanges(c, resp2.PageList, err)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfMatchFalse(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfNoneMatchTrue(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateBasicGetPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobGetPageRangesIfNoneMatchFalse(c *chk.C) {
	containerClient, pbClient := setupGetPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, 0, 0, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)
	//serr := err.(StorageError)
	//c.Assert(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
}

func setupDiffPageRangesTest(c *chk.C) (containerClient ContainerClient, pbClient PageBlobClient, snapshot string) {
	bsu := getBSU()
	containerClient, _ = createNewContainer(c, bsu)
	pbClient, _ = createNewPageBlob(c, containerClient)

	r := getReaderToRandomBytes(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.IsNil)
	snapshot = *resp.Snapshot

	r = getReaderToRandomBytes(PageBlobPageBytes)
	offset, count = int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions = UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
	}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	return
}

func validateDiffPageRanges(c *chk.C, resp *PageList, err error) {
	c.Assert(err, chk.IsNil)
	pageListResp := resp.PageRange
	c.Assert(pageListResp, chk.NotNil)
	c.Assert(*resp.PageRange, chk.HasLen, 1)
	start, end := int64(0), int64(PageBlobPageBytes-1)
	c.Assert((*pageListResp)[0], chk.DeepEquals, PageRange{Start: &start, End: &end})
}

func (s *aztestsSuite) TestBlobDiffPageRangesNonExistentSnapshot(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	snapshotTime, _ := time.Parse(SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	_, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshotTime.Format(SnapshotTimeFormat), nil)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodePreviousSnapshotNotFound)
}

func (s *aztestsSuite) TestBlobDiffPageRangeInvalidRange(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)
	_, err := pbClient.GetPageRangesDiff(ctx, -22, 14, snapshot, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfModifiedSinceTrue(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateDiffPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfModifiedSinceFalse(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)

	//stgErr := err.(StorageError)
	//c.Assert(stgErr.Response().StatusCode, chk.Equals, 304)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateDiffPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfMatchTrue(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateDiffPageRanges(c, resp2.PageList, err)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfMatchFalse(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfNoneMatchTrue(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.IsNil)
	validateDiffPageRanges(c, resp.PageList, err)
}

func (s *aztestsSuite) TestBlobDiffPageRangeIfNoneMatchFalse(c *chk.C) {
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(c)
	defer deleteContainer(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, 0, 0, snapshot, &getPageRangesOptions)
	c.Assert(err, chk.NotNil)

	//serr := err.(StorageError)
	//c.Assert(serr.Response().StatusCode, chk.Equals, 304)
}

func (s *aztestsSuite) TestBlobResizeZero(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	// The default pbClient is created with size > 0, so this should actually update
	_, err := pbClient.Resize(ctx, 0, nil)
	c.Assert(err, chk.IsNil)

	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.ContentLength, chk.Equals, int64(0))
}

func (s *aztestsSuite) TestBlobResizeInvalidSizeNegative(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	_, err := pbClient.Resize(ctx, -4, nil)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobResizeInvalidSizeMisaligned(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	_, err := pbClient.Resize(ctx, 12, nil)
	c.Assert(err, chk.NotNil)
}

func validateResize(c *chk.C, pbClient PageBlobClient) {
	resp, _ := pbClient.GetProperties(ctx, nil)
	c.Assert(*resp.ContentLength, chk.Equals, int64(PageBlobPageBytes))
}

func (s *aztestsSuite) TestBlobResizeIfModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateResize(c, pbClient)
}

func (s *aztestsSuite) TestBlobResizeIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobResizeIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateResize(c, pbClient)
}

func (s *aztestsSuite) TestBlobResizeIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobResizeIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateResize(c, pbClient)
}

func (s *aztestsSuite) TestBlobResizeIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobResizeIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateResize(c, pbClient)
}

func (s *aztestsSuite) TestBlobResizeIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberActionTypeInvalid(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeInvalidHeaderValue)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberSequenceNumberInvalid(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	defer func() { // Invalid sequence number should panic
		recover()
	}()

	sequenceNumber := int64(-1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}

	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeInvalidHeaderValue)
}

func validateSequenceNumberSet(c *chk.C, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.BlobSequenceNumber, chk.Equals, int64(1))
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	validateSequenceNumberSet(c, pbClient)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	validateSequenceNumberSet(c, pbClient)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	validateSequenceNumberSet(c, pbClient)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	eTag := "garbage"
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	eTag := "garbage"
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.IsNil)

	validateSequenceNumberSet(c, pbClient)
}

func (s *aztestsSuite) TestBlobSetSequenceNumberIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		ActionType: &actionType,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func setupStartIncrementalCopyTest(c *chk.C) (containerClient ContainerClient, pbClient PageBlobClient, copyPBClient PageBlobClient, snapshot string) {
	bsu := getBSU()
	containerClient, _ = createNewContainer(c, bsu)

	accessType := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
	}
	_, err := containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	pbClient, _ = createNewPageBlob(c, containerClient)
	resp, _ := pbClient.CreateSnapshot(ctx, nil)
	copyPBClient, _ = getPageBlobClient(c, containerClient)

	// Must create the incremental copy pbClient so that the access conditions work on it
	resp2, err := copyPBClient.StartCopyIncremental(ctx, pbClient.URL(), *resp.Snapshot, nil)
	c.Assert(err, chk.IsNil)
	waitForIncrementalCopy(c, copyPBClient, &resp2)

	resp, _ = pbClient.CreateSnapshot(ctx, nil) // Take a new snapshot so the next copy will succeed
	snapshot = *resp.Snapshot
	return
}

func validateIncrementalCopy(c *chk.C, copyBlobURL PageBlobClient, resp *PageBlobCopyIncrementalResponse) {
	t := waitForIncrementalCopy(c, copyBlobURL, resp)

	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
	copySnapshotURL := copyBlobURL.WithSnapshot(*t)
	_, err := copySnapshotURL.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopySnapshotNotExist(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlob(c, containerClient)
	copyBlobURL, _ := getPageBlobClient(c, containerClient)

	snapshot := time.Now().UTC().Format(SnapshotTimeFormat)
	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, nil)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeCannotVerifyCopySource)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-20)

	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	resp, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateIncrementalCopy(c, copyBlobURL, &resp)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(20)

	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(20)

	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	resp, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateIncrementalCopy(c, copyBlobURL, &resp)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(-20)

	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfMatchTrue(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	resp, _ := copyBlobURL.GetProperties(ctx, nil)

	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	resp2, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateIncrementalCopy(c, copyBlobURL, &resp2)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfMatchFalse(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: &eTag,
		},
	}
	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeTargetConditionNotMet)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	eTag := "garbage"
	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: &eTag,
		},
	}
	resp, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.IsNil)

	validateIncrementalCopy(c, copyBlobURL, &resp)
}

func (s *aztestsSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse(c *chk.C) {
	containerClient, pbClient, copyBlobURL, snapshot := setupStartIncrementalCopyTest(c)

	defer deleteContainer(c, containerClient)

	resp, _ := copyBlobURL.GetProperties(ctx, nil)

	copyIncrementalPageBlobOptions := CopyIncrementalPageBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}
	_, err := copyBlobURL.StartCopyIncremental(ctx, pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, ServiceCodeConditionNotMet)
}
