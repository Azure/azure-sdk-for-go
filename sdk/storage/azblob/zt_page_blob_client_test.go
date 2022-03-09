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

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
)

func TestPutGetPages(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	contentSize := 1024
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	reader, _ := generateData(1024)
	putResp, err := pbClient.UploadPages(context.Background(), reader, &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, putResp.RawResponse.StatusCode, 201)
	require.NotNil(t, putResp.LastModified)
	require.Equal(t, (*putResp.LastModified).IsZero(), false)
	require.NotNil(t, putResp.ETag)
	require.Nil(t, putResp.ContentMD5)
	require.Equal(t, *putResp.BlobSequenceNumber, int64(0))
	require.NotNil(t, *putResp.RequestID)
	require.NotNil(t, *putResp.Version)
	require.NotNil(t, putResp.Date)
	require.Equal(t, (*putResp.Date).IsZero(), false)

	pageList, err := pbClient.GetPageRanges(context.Background(), HttpRange{0, 1023}, nil)
	require.NoError(t, err)
	require.Equal(t, pageList.RawResponse.StatusCode, 200)
	require.NotNil(t, pageList.LastModified)
	require.Equal(t, (*pageList.LastModified).IsZero(), false)
	require.NotNil(t, pageList.ETag)
	require.Equal(t, *pageList.BlobContentLength, int64(512*10))
	require.NotNil(t, *pageList.RequestID)
	require.NotNil(t, *pageList.Version)
	require.NotNil(t, pageList.Date)
	require.Equal(t, (*pageList.Date).IsZero(), false)
	require.NotNil(t, pageList.PageList)
	pageRangeResp := pageList.PageList.PageRange
	require.Len(t, pageRangeResp, 1)
	rawStart, rawEnd := (pageRangeResp)[0].Raw()
	require.Equal(t, rawStart, offset)
	require.Equal(t, rawEnd, count-1)
}

func TestUploadPagesFromURL(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(t, "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(t, "dstblob", containerClient, int64(contentSize))

	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, uploadSrcResp1.RawResponse.StatusCode, 201)
	require.NotNil(t, uploadSrcResp1.LastModified)
	require.Equal(t, (*uploadSrcResp1.LastModified).IsZero(), false)
	require.NotNil(t, uploadSrcResp1.ETag)
	require.Nil(t, uploadSrcResp1.ContentMD5)
	require.Equal(t, *uploadSrcResp1.BlobSequenceNumber, int64(0))
	require.NotNil(t, *uploadSrcResp1.RequestID)
	require.NotNil(t, *uploadSrcResp1.Version)
	require.NotNil(t, uploadSrcResp1.Date)
	require.Equal(t, (*uploadSrcResp1.Date).IsZero(), false)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL.
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
	require.NoError(t, err)
	require.Equal(t, pResp1.RawResponse.StatusCode, 201)
	require.NotNil(t, pResp1.ETag)
	require.NotNil(t, pResp1.LastModified)
	require.NotNil(t, pResp1.ContentMD5)
	require.NotNil(t, pResp1.RequestID)
	require.NotNil(t, pResp1.Version)
	require.NotNil(t, pResp1.Date)
	require.Equal(t, (*pResp1.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{}))
	require.NoError(t, err)
	require.EqualValues(t, destData, sourceData)
}

func TestUploadPagesFromURLWithMD5(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	contentMD5 := md5Value[:]
	ctx := context.Background() // Use default Background context
	srcBlob := createNewPageBlobWithSize(t, "srcblob", containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(t, "dstblob", containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, uploadSrcResp1.RawResponse.StatusCode, 201)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Upload page from URL with MD5.
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
	}
	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	require.NoError(t, err)
	require.Equal(t, pResp1.RawResponse.StatusCode, 201)
	require.NotNil(t, pResp1.ETag)
	require.NotNil(t, pResp1.LastModified)
	require.NotNil(t, pResp1.ContentMD5)
	require.EqualValues(t, pResp1.ContentMD5, contentMD5)
	require.NotNil(t, pResp1.RequestID)
	require.NotNil(t, pResp1.Version)
	require.NotNil(t, pResp1.Date)
	require.Equal(t, (*pResp1.Date).IsZero(), false)
	require.Equal(t, *pResp1.BlobSequenceNumber, int64(0))

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{}))
	require.NoError(t, err)
	require.EqualValues(t, destData, sourceData)

	// Upload page from URL with bad MD5
	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions = UploadPagesFromURLOptions{
		SourceContentMD5: badContentMD5,
	}
	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeMD5Mismatch)
}

func TestClearDiffPages(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	contentSize := 2 * 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	require.NoError(t, err)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	require.NoError(t, err)

	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}}
	_, err = pbClient.UploadPages(context.Background(), getReaderToGeneratedBytes(2048), &uploadPagesOptions1)
	require.NoError(t, err)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
	require.NoError(t, err)
	pageRangeResp := pageListResp.PageList.PageRange
	require.NotNil(t, pageRangeResp)
	require.Len(t, pageRangeResp, 1)
	// require.(t, (pageRangeResp)[0], chk.DeepEquals, PageRange{Start: &offset1, End: &end1})
	rawStart, rawEnd := (pageRangeResp)[0].Raw()
	require.Equal(t, rawStart, offset1)
	require.Equal(t, rawEnd, end1)

	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, nil)
	require.NoError(t, err)
	require.Equal(t, clearResp.RawResponse.StatusCode, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
	require.NoError(t, err)
	require.Nil(t, pageListResp.PageList.PageRange)
}

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

func TestIncrementalCopy(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	accessType := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{Access: &accessType},
	}
	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
	require.NoError(t, err)

	srcBlob := createNewPageBlob(t, "src"+generateBlobName(testName), containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = srcBlob.UploadPages(context.Background(), r, &uploadPagesOptions)
	require.NoError(t, err)

	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
	require.NoError(t, err)

	dstBlob := containerClient.NewPageBlobClient("dst" + generateBlobName(testName))

	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 202)
	require.NotNil(t, resp.LastModified)
	require.Equal(t, (*resp.LastModified).IsZero(), false)
	require.NotNil(t, resp.ETag)
	require.NotEqual(t, *resp.RequestID, "")
	require.NotEqual(t, *resp.Version, "")
	require.NotNil(t, resp.Date)
	require.Equal(t, (*resp.Date).IsZero(), false)
	require.NotEqual(t, *resp.CopyID, "")
	require.Equal(t, *resp.CopyStatus, CopyStatusTypePending)

	waitForIncrementalCopy(t, dstBlob, &resp)
}

func TestResizePageBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	resp, err := pbClient.Resize(context.Background(), 2048, nil)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 200)

	resp, err = pbClient.Resize(context.Background(), 8192, nil)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 200)

	resp2, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp2.ContentLength, int64(8192))
}

func TestPageSequenceNumbers(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	sequenceNumber := int64(0)
	actionType := SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(7)
	actionType = SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 200)

	sequenceNumber = int64(11)
	actionType = SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob = UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	resp, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 200)
}

func TestPutPagesWithMD5(t *testing.T) {
	t.Skip("AuthenticationFailed")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.NoError(t, err)
	require.Equal(t, putResp.RawResponse.StatusCode, 201)
	require.NotNil(t, putResp.LastModified)
	require.Equal(t, (*putResp.LastModified).IsZero(), false)
	require.NotNil(t, putResp.ETag)
	require.NotNil(t, putResp.ContentMD5)
	require.EqualValues(t, putResp.ContentMD5, contentMD5)
	require.Equal(t, *putResp.BlobSequenceNumber, int64(0))
	require.NotNil(t, *putResp.RequestID)
	require.NotNil(t, *putResp.Version)
	require.NotNil(t, putResp.Date)
	require.Equal(t, (*putResp.Date).IsZero(), false)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeMD5Mismatch)
}

func TestBlobCreatePageSizeInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, 1, &createPageBlobOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidHeaderValue)
}

func TestBlobCreatePageSequenceInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.Error(t, err)
}

func TestBlobCreatePageMetadataNonEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           basicMetadata,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.NoError(t, err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Metadata)
	require.EqualValues(t, resp.Metadata, basicMetadata)
}

func TestBlobCreatePageMetadataEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.NoError(t, err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.Metadata)
}

func TestBlobCreatePageMetadataInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		Metadata:           map[string]string{"In valid1": "bar"},
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.Error(t, err)
	require.Contains(t, err.Error(), invalidHeaderErrorSubstring)

}

func TestBlobCreatePageHTTPHeaders(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := CreatePageBlobOptions{
		BlobSequenceNumber: &sequenceNumber,
		HTTPHeaders:        &basicHeaders,
	}
	_, err = pbClient.Create(ctx, PageBlobPageBytes, &createPageBlobOptions)
	require.NoError(t, err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	h := resp.GetHTTPHeaders()
	require.EqualValues(t, h, basicHeaders)
}

func validatePageBlobPut(t *testing.T, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Metadata)
	require.EqualValues(t, resp.Metadata, basicMetadata)
	require.EqualValues(t, resp.GetHTTPHeaders(), basicHeaders)
}

func TestBlobCreatePageIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(t, err)

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
	require.NoError(t, err)

	validatePageBlobPut(t, pbClient)
}

func TestBlobCreatePageIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(t, err)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobCreatePageIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(t, err)

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
	require.NoError(t, err)

	validatePageBlobPut(t, pbClient)
}

func TestBlobCreatePageIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(ctx, PageBlobPageBytes, nil)
	require.NoError(t, err)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobCreatePageIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

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
	require.NoError(t, err)

	validatePageBlobPut(t, pbClient)
}

func TestBlobCreatePageIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobCreatePageIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.NoError(t, err)

	validatePageBlobPut(t, pbClient)
}

func TestBlobCreatePageIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutPagesInvalidRange(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	contentSize := 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize/2)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(t, err)
}

func TestBlobPutPagesEmptyBody(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	r := bytes.NewReader([]byte{})
	offset, count := int64(0), int64(0)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.Error(t, err)
}

func TestBlobPutPagesNonExistentBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(0), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}}
	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeBlobNotFound)
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

func TestBlobPutPagesIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.NoError(t, err)

	validateUploadPages(t, pbClient)
}

func TestBlobPutPagesIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutPagesIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.NoError(t, err)

	validateUploadPages(t, pbClient)
}

func TestBlobPutPagesIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutPagesIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.NoError(t, err)

	validateUploadPages(t, pbClient)
}

func TestBlobPutPagesIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutPagesIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.NoError(t, err)

	validateUploadPages(t, pbClient)
}

func TestBlobPutPagesIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutPagesIfSequenceNumberLessThanTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.NoError(t, err)

	validateUploadPages(t, pbClient)
}

func TestBlobPutPagesIfSequenceNumberLessThanFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func TestBlobPutPagesIfSequenceNumberLessThanNegOne(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidInput)
}

func TestBlobPutPagesIfSequenceNumberLTETrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

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
	require.NoError(t, err)

	validateUploadPages(t, pbClient)
}

func TestBlobPutPagesIfSequenceNumberLTEqualFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func TestBlobPutPagesIfSequenceNumberLTENegOne(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.Error(t, err)
}

func TestBlobPutPagesIfSequenceNumberEqualTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

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
	require.NoError(t, err)

	validateUploadPages(t, pbClient)
}

func TestBlobPutPagesIfSequenceNumberEqualFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func setupClearPagesTest(t *testing.T, testName string) (ContainerClient, PageBlobClient) {
	svcClient, err := createServiceClient(t, testAccountDefault)
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

func TestBlobClearPagesInvalidRange(t *testing.T) {
	stop := start(t)
	defer stop()

	containerClient, pbClient := setupClearPagesTest(t, t.Name())
	defer deleteContainer(t, containerClient)

	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes + 1}, nil)
	require.Error(t, err)
}

func TestBlobClearPagesIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(t, err)

	validateClearPagesTest(t, pbClient)
}

func TestBlobClearPagesIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()

	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobClearPagesIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(t, err)

	validateClearPagesTest(t, pbClient)
}

func TestBlobClearPagesIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobClearPagesIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: getPropertiesResp.ETag,
			},
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(t, err)

	validateClearPagesTest(t, pbClient)
}

func TestBlobClearPagesIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobClearPagesIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	eTag := "garbage"
	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(t, err)

	validateClearPagesTest(t, pbClient)
}

func TestBlobClearPagesIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	clearPageOptions := ClearPagesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobClearPagesIfSequenceNumberLessThanTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(t, err)

	validateClearPagesTest(t, pbClient)
}

func TestBlobClearPagesIfSequenceNumberLessThanFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func TestBlobClearPagesIfSequenceNumberLessThanNegOne(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidInput)
}

func TestBlobClearPagesIfSequenceNumberLTETrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(t, err)

	validateClearPagesTest(t, pbClient)
}

func TestBlobClearPagesIfSequenceNumberLTEFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func TestBlobClearPagesIfSequenceNumberLTENegOne(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidInput)
}

func TestBlobClearPagesIfSequenceNumberEqualTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.NoError(t, err)

	validateClearPagesTest(t, pbClient)
}

func TestBlobClearPagesIfSequenceNumberEqualFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	sequenceNumber := int64(10)
	actionType := SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.NoError(t, err)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeSequenceNumberConditionNotMet)
}

func TestBlobClearPagesIfSequenceNumberEqualNegOne(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupClearPagesTest(t, testName)
	defer deleteContainer(t, containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := ClearPagesOptions{
		SequenceNumberAccessConditions: &SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(ctx, HttpRange{0, PageBlobPageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidInput)
}

func setupGetPageRangesTest(t *testing.T, testName string) (containerClient ContainerClient, pbClient PageBlobClient) {
	svcClient, err := createServiceClient(t, testAccountDefault)
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

func TestBlobGetPageRangesEmptyBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(t, err)
	require.Nil(t, resp.PageList.PageRange)
}

func TestBlobGetPageRangesEmptyRange(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(t, err)
	validateBasicGetPageRanges(t, resp.PageList, err)
}

func TestBlobGetPageRangesInvalidRange(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	_, err := pbClient.GetPageRanges(ctx, HttpRange{-2, 500}, nil)
	require.NoError(t, err)
}

func TestBlobGetPageRangesNonContiguousRanges(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	r, _ := generateData(PageBlobPageBytes)
	offset, count := int64(2*PageBlobPageBytes), int64(PageBlobPageBytes)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	_, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(t, err)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(t, err)
	pageListResp := resp.PageList.PageRange
	require.NotNil(t, pageListResp)
	require.Len(t, pageListResp, 2)

	start, end := int64(0), int64(PageBlobPageBytes-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	require.Equal(t, rawStart, start)
	require.Equal(t, rawEnd, end)

	start, end = int64(PageBlobPageBytes*2), int64((PageBlobPageBytes*3)-1)
	rawStart, rawEnd = pageListResp[1].Raw()
	require.Equal(t, rawStart, start)
	require.Equal(t, rawEnd, end)
}

func TestBlobGetPageRangesNotPageAligned(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 2000}, nil)
	require.NoError(t, err)
	validateBasicGetPageRanges(t, resp.PageList, err)
}

func TestBlobGetPageRangesSnapshot(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Snapshot)

	snapshotURL := pbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetPageRanges(ctx, HttpRange{0, 0}, nil)
	require.NoError(t, err)
	validateBasicGetPageRanges(t, resp2.PageList, err)
}

func TestBlobGetPageRangesIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(t, err)
	validateBasicGetPageRanges(t, resp.PageList, err)
}

func TestBlobGetPageRangesIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(t, err)
}

func TestBlobGetPageRangesIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(t, err)
	validateBasicGetPageRanges(t, resp.PageList, err)
}

func TestBlobGetPageRangesIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	currentTime := getRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobGetPageRangesIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(t, err)
	validateBasicGetPageRanges(t, resp2.PageList, err)
}

func TestBlobGetPageRangesIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobGetPageRangesIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.NoError(t, err)
	validateBasicGetPageRanges(t, resp.PageList, err)
}

func TestBlobGetPageRangesIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient := setupGetPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRanges(ctx, HttpRange{0, 0}, &getPageRangesOptions)
	require.Error(t, err)
}

func setupDiffPageRangesTest(t *testing.T, testName string) (containerClient ContainerClient, pbClient PageBlobClient, snapshot string) {
	svcClient, err := createServiceClient(t, testAccountDefault)
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

func TestBlobDiffPageRangesNonExistentSnapshot(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	snapshotTime, _ := time.Parse(SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshotTime.Format(SnapshotTimeFormat), nil)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodePreviousSnapshotNotFound)
}

func TestBlobDiffPageRangeInvalidRange(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{-22, 14}, snapshot, nil)
	require.NoError(t, err)
}

func TestBlobDiffPageRangeIfModifiedSinceTrue(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(t, err)
	validateDiffPageRanges(t, resp.PageList, err)
}

func TestBlobDiffPageRangeIfModifiedSinceFalse(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(t, err)
}

func TestBlobDiffPageRangeIfUnmodifiedSinceTrue(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeGMT(10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(t, err)
	validateDiffPageRanges(t, resp.PageList, err)
}

func TestBlobDiffPageRangeIfUnmodifiedSinceFalse(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobDiffPageRangeIfMatchTrue(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	resp2, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(t, err)
	validateDiffPageRanges(t, resp2.PageList, err)
}

func TestBlobDiffPageRangeIfMatchFalse(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobDiffPageRangeIfNoneMatchTrue(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	eTag := "garbage"
	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	resp, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.NoError(t, err)
	validateDiffPageRanges(t, resp.PageList, err)
}

func TestBlobDiffPageRangeIfNoneMatchFalse(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	getPageRangesOptions := GetPageRangesOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.GetPageRangesDiff(ctx, HttpRange{0, 0}, snapshot, &getPageRangesOptions)
	require.Error(t, err)
}

func TestBlobResizeZero(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	// The default pbClient is created with size > 0, so this should actually update
	_, err = pbClient.Resize(ctx, 0, nil)
	require.NoError(t, err)

	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentLength, int64(0))
}

func TestBlobResizeInvalidSizeNegative(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	_, err = pbClient.Resize(ctx, -4, nil)
	require.Error(t, err)
}

func TestBlobResizeInvalidSizeMisaligned(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	_, err = pbClient.Resize(ctx, 12, nil)
	require.Error(t, err)
}

func validateResize(t *testing.T, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentLength, int64(PageBlobPageBytes))
}

func TestBlobResizeIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(t, err)

	validateResize(t, pbClient)
}

func TestBlobResizeIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobResizeIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(t, err)

	validateResize(t, pbClient)
}

func TestBlobResizeIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

	currentTime := getRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobResizeIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(t, err)

	validateResize(t, pbClient)
}

func TestBlobResizeIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobResizeIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(t, err)

	validateResize(t, pbClient)
}

func TestBlobResizeIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	resp, _ := pbClient.GetProperties(ctx, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobSetSequenceNumberActionTypeInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := UpdateSequenceNumberPageBlob{
		BlobSequenceNumber: &sequenceNumber,
		ActionType:         &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(ctx, &updateSequenceNumberPageBlob)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidHeaderValue)
}

func TestBlobSetSequenceNumberSequenceNumberInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidHeaderValue)
}

func validateSequenceNumberSet(t *testing.T, pbClient PageBlobClient) {
	resp, err := pbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.BlobSequenceNumber, int64(1))
}

func TestBlobSetSequenceNumberIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.NoError(t, err)

	validateSequenceNumberSet(t, pbClient)
}

func TestBlobSetSequenceNumberIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobSetSequenceNumberIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.NoError(t, err)

	validateSequenceNumberSet(t, pbClient)
}

func TestBlobSetSequenceNumberIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(ctx, PageBlobPageBytes*10, nil)
	require.NoError(t, err)
	require.Equal(t, pageBlobCreateResponse.RawResponse.StatusCode, 201)
	require.NotNil(t, pageBlobCreateResponse.Date)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobSetSequenceNumberIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.NoError(t, err)

	validateSequenceNumberSet(t, pbClient)
}

func TestBlobSetSequenceNumberIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, blobName, containerClient)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobSetSequenceNumberIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, "src"+blobName, containerClient)

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
	require.NoError(t, err)

	validateSequenceNumberSet(t, pbClient)
}

func TestBlobSetSequenceNumberIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	pbClient := createNewPageBlob(t, "src"+blobName, containerClient)

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
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}
