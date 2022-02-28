// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
)

func TestStageGetBlocks(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		io.NopCloser(strings.NewReader("hello world"))
		putResp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], internal.NopCloser(strings.NewReader(d)), nil)
		require.NoError(t, err)
		require.Equal(t, putResp.RawResponse.StatusCode, 201)
		require.Nil(t, putResp.ContentMD5)
		require.NotNil(t, putResp.RequestID)
		require.NotNil(t, putResp.Version)
		require.NotNil(t, putResp.Date)
		require.Equal(t, (*putResp.Date).IsZero(), false)
	}

	blockList, err := bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Equal(t, blockList.RawResponse.StatusCode, 200)
	require.Nil(t, blockList.LastModified)
	require.Nil(t, blockList.ETag)
	require.NotNil(t, blockList.ContentType)
	require.Nil(t, blockList.BlobContentLength)
	require.NotNil(t, blockList.RequestID)
	require.NotNil(t, blockList.Version)
	require.NotNil(t, blockList.Date)
	require.Equal(t, (*blockList.Date).IsZero(), false)
	require.NotNil(t, blockList.BlockList)
	require.Nil(t, blockList.BlockList.CommittedBlocks)
	require.NotNil(t, blockList.BlockList.UncommittedBlocks)
	require.Len(t, blockList.BlockList.UncommittedBlocks, len(data))

	listResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
	require.NoError(t, err)
	require.Equal(t, listResp.RawResponse.StatusCode, 201)
	require.NotNil(t, listResp.LastModified)
	require.Equal(t, (*listResp.LastModified).IsZero(), false)
	require.NotNil(t, listResp.ETag)
	require.NotNil(t, listResp.RequestID)
	require.NotNil(t, listResp.Version)
	require.NotNil(t, listResp.Date)
	require.Equal(t, (*listResp.Date).IsZero(), false)

	blockList, err = bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Equal(t, blockList.RawResponse.StatusCode, 200)
	require.NotNil(t, blockList.LastModified)
	require.Equal(t, (*blockList.LastModified).IsZero(), false)
	require.NotNil(t, blockList.ETag)
	require.NotNil(t, blockList.ContentType)
	require.Equal(t, *blockList.BlobContentLength, int64(25))
	require.NotNil(t, blockList.RequestID)
	require.NotNil(t, blockList.Version)
	require.NotNil(t, blockList.Date)
	require.Equal(t, (*blockList.Date).IsZero(), false)
	require.NotNil(t, blockList.BlockList)
	require.NotNil(t, blockList.BlockList.CommittedBlocks)
	require.Nil(t, blockList.BlockList.UncommittedBlocks)
	require.Len(t, blockList.BlockList.CommittedBlocks, len(data))
}

func TestStageBlockFromURL(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	ctx := context.Background() // Use default Background context
	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))

	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, nil)
	require.NoError(t, err)
	require.Equal(t, uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockIDs := generateBlockIDsList(2)

	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockIDs[0], srcBlobURLWithSAS, 0, &StageBlockFromURLOptions{
		Offset: to.Int64Ptr(0),
		Count:  to.Int64Ptr(int64(contentSize / 2)),
	})
	require.NoError(t, err)
	require.Equal(t, stageResp1.RawResponse.StatusCode, 201)
	require.NotEqual(t, stageResp1.ContentMD5, "")
	require.NotEqual(t, stageResp1.RequestID, "")
	require.NotEqual(t, stageResp1.Version, "")
	require.Equal(t, stageResp1.Date.IsZero(), false)

	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockIDs[1], srcBlobURLWithSAS, 0, &StageBlockFromURLOptions{
		Offset: to.Int64Ptr(int64(contentSize / 2)),
		Count:  to.Int64Ptr(int64(CountToEnd)),
	})
	require.NoError(t, err)
	require.Equal(t, stageResp2.RawResponse.StatusCode, 201)
	require.NotEqual(t, stageResp2.ContentMD5, "")
	require.NotEqual(t, stageResp2.RequestID, "")
	require.NotEqual(t, stageResp2.Version, "")
	require.Equal(t, stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Equal(t, blockList.RawResponse.StatusCode, 200)
	require.NotNil(t, blockList.BlockList)
	require.Nil(t, blockList.BlockList.CommittedBlocks)
	require.NotNil(t, blockList.BlockList.UncommittedBlocks)
	require.Len(t, blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	listResp, err := destBlob.CommitBlockList(context.Background(), blockIDs, nil)
	require.NoError(t, err)
	require.Equal(t, listResp.RawResponse.StatusCode, 201)
	require.NotNil(t, listResp.LastModified)
	require.Equal(t, (*listResp.LastModified).IsZero(), false)
	require.NotNil(t, listResp.ETag)
	require.NotNil(t, listResp.RequestID)
	require.NotNil(t, listResp.Version)
	require.NotNil(t, listResp.Date)
	require.Equal(t, (*listResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, content)
}

func TestCopyBlockBlobFromURL(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	contentMD5 := md5.Sum(content)
	body := bytes.NewReader(content)
	ctx := context.Background()

	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, internal.NopCloser(body), nil)
	require.NoError(t, err)
	require.Equal(t, uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Invoke copy bbClient from URL.
	sourceContentMD5 := contentMD5[:]
	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &CopyBlockBlobFromURLOptions{
		Metadata:         map[string]string{"foo": "bar"},
		SourceContentMD5: sourceContentMD5,
	})
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 202)
	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.RequestID)
	require.NotNil(t, resp.Version)
	require.NotNil(t, resp.Date)
	require.Equal(t, (*resp.Date).IsZero(), false)
	require.NotNil(t, resp.CopyID)
	require.EqualValues(t, resp.ContentMD5, sourceContentMD5)
	require.Equal(t, *resp.CopyStatus, "success")

	// Make sure the metadata got copied over
	getPropResp, err := destBlob.GetProperties(ctx, nil)
	require.NoError(t, err)
	metadata := getPropResp.Metadata
	require.NotNil(t, metadata)
	require.Len(t, metadata, 1)
	require.EqualValues(t, metadata, map[string]string{"Foo": "bar"})

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, content)

	// Edge case 1: Provide bad MD5 and make sure the copy fails
	_, badMD5 := getRandomDataAndReader(16)
	copyBlockBlobFromURLOptions1 := CopyBlockBlobFromURLOptions{
		SourceContentMD5: badMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
	require.Error(t, err)

	// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
	copyBlockBlobFromURLOptions2 := CopyBlockBlobFromURLOptions{
		SourceContentMD5: sourceContentMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 202)
	require.EqualValues(t, *resp.CopyStatus, "success")
}

func TestBlobSASQueryParamOverrideResponseHeaders(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	//contentMD5 := md5.Sum(content)

	ctx := context.Background()

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadSrcResp, err := bbClient.Upload(ctx, internal.NopCloser(body), nil)
	require.NoError(t, err)
	require.Equal(t, uploadSrcResp.RawResponse.StatusCode, 201)

	// Get blob url with SAS.
	blobParts := NewBlobURLParts(bbClient.URL())

	cacheControlVal := "cache-control-override"
	contentDispositionVal := "content-disposition-override"
	contentEncodingVal := "content-encoding-override"
	contentLanguageVal := "content-language-override"
	contentTypeVal := "content-type-override"

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	// Append User Delegation SAS token to URL
	blobParts.SAS, err = BlobSASSignatureValues{
		Protocol:           SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:         time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName:      blobParts.ContainerName,
		BlobName:           blobParts.BlobName,
		Permissions:        BlobSASPermissions{Read: true}.String(),
		CacheControl:       cacheControlVal,
		ContentDisposition: contentDispositionVal,
		ContentEncoding:    contentEncodingVal,
		ContentLanguage:    contentLanguageVal,
		ContentType:        contentTypeVal,
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	// Generate new bbClient client
	blobURLWithSAS := blobParts.URL()
	require.NotNil(t, blobURLWithSAS)

	blobClientWithSAS, err := NewBlockBlobClientWithNoCredential(blobURLWithSAS, nil)
	require.NoError(t, err)

	gResp, err := blobClientWithSAS.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *gResp.CacheControl, cacheControlVal)
	require.Equal(t, *gResp.ContentDisposition, contentDispositionVal)
	require.Equal(t, *gResp.ContentEncoding, contentEncodingVal)
	require.Equal(t, *gResp.ContentLanguage, contentLanguageVal)
	require.Equal(t, *gResp.ContentType, contentTypeVal)
}

func TestStageBlockWithMD5(t *testing.T) {
	t.Skip("authentication failed")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	// test put block with valid MD5 value
	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	putResp, err := bbClient.StageBlock(context.Background(), blockID1, rsc, &StageBlockOptions{
		BlockBlobStageBlockOptions: &BlockBlobStageBlockOptions{
			TransactionalContentMD5: contentMD5,
		},
	})
	require.NoError(t, err)
	require.Equal(t, putResp.RawResponse.StatusCode, 201)
	require.EqualValues(t, putResp.ContentMD5, contentMD5)
	require.NotNil(t, putResp.RequestID)
	require.NotNil(t, putResp.Version)
	require.NotNil(t, putResp.Date)
	require.Equal(t, (*putResp.Date).IsZero(), false)

	// test put block with bad MD5 value
	_, badContent := getRandomDataAndReader(contentSize)
	badMD5Value := md5.Sum(badContent)
	badContentMD5 := badMD5Value[:]

	_, _ = rsc.Seek(0, io.SeekStart)
	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	_, err = bbClient.StageBlock(context.Background(), blockID2, rsc, &StageBlockOptions{
		BlockBlobStageBlockOptions: &BlockBlobStageBlockOptions{
			TransactionalContentMD5: badContentMD5,
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), StorageErrorCodeMD5Mismatch)
}

func TestBlobPutBlobHTTPHeaders(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &UploadBlockBlobOptions{
		HTTPHeaders: &basicHeaders,
	})
	require.NoError(t, err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	h := resp.GetHTTPHeaders()
	h.BlobContentMD5 = nil // the service generates a MD5 value, omit before comparing
	require.EqualValues(t, h, basicHeaders)
}

func TestBlobPutBlobMetadataNotEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &UploadBlockBlobOptions{
		Metadata: basicMetadata,
	})
	require.NoError(t, err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	actualMetadata := resp.Metadata
	require.NotNil(t, actualMetadata)
	require.EqualValues(t, actualMetadata, basicMetadata)
}

func TestBlobPutBlobMetadataEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, nil)
	require.NoError(t, err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.Metadata)
}

func TestBlobPutBlobMetadataInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &UploadBlockBlobOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), invalidHeaderErrorSubstring)
}

func TestBlobPutBlobIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, createResp.RawResponse.StatusCode, 201)
	require.NotNil(t, createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, -10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	require.NoError(t, err)
	validateUpload(t, bbClient.BlobClient)
}

func TestBlobPutBlobIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, createResp.RawResponse.StatusCode, 201)
	require.NotNil(t, createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, 10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutBlobIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, createResp.RawResponse.StatusCode, 201)
	require.NotNil(t, createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, 10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	require.NoError(t, err)

	validateUpload(t, bbClient.BlobClient)
}

func TestBlobPutBlobIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	require.Equal(t, createResp.RawResponse.StatusCode, 201)
	require.NotNil(t, createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, -10)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(ctx, internal.NopCloser(bytes.NewReader(nil)), &uploadBlockBlobOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutBlobIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	require.NoError(t, err)

	validateUpload(t, bbClient.BlobClient)
}

func TestBlobPutBlobIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)

	ifMatch := "garbage"
	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &ifMatch,
			},
		},
	}
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &uploadBlockBlobOptions)
	require.Error(t, err)
	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutBlobIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	ifNoneMatch := "garbage"
	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &ifNoneMatch,
			},
		},
	}

	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	require.NoError(t, err)

	validateUpload(t, bbClient.BlobClient)
}

func TestBlobPutBlobIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func validateBlobCommitted(t *testing.T, bbClient BlockBlobClient) {
	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Len(t, resp.BlockList.CommittedBlocks, 1)
}

func setupPutBlockListTest(t *testing.T, testName string) (ContainerClient, BlockBlobClient, []string) {
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)

	blockIDs := generateBlockIDsList(1)
	_, err = bbClient.StageBlock(ctx, blockIDs[0], internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)
	return containerClient, bbClient, blockIDs
}

func TestBlobPutBlockListHTTPHeadersEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobHTTPHeaders: &BlobHTTPHeaders{BlobContentDisposition: &blobContentDisposition},
	})
	require.NoError(t, err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, nil)
	require.NoError(t, err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.ContentDisposition)
}

func TestBlobPutBlockListIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(t, err)
	require.NotNil(t, commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	require.NoError(t, err)

	validateBlobCommitted(t, bbClient)
}

func TestBlobPutBlockListIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	getPropertyResp, err := containerClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, getPropertyResp.Date)

	currentTime := getRelativeTimeFromAnchor(getPropertyResp.Date, 10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_ = err

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutBlockListIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(t, err)
	require.NotNil(t, commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, 10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	require.NoError(t, err)

	validateBlobCommitted(t, bbClient)
}

func TestBlobPutBlockListIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(t, err)
	require.NotNil(t, commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutBlockListIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(t, err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}},
	})
	require.NoError(t, err)

	validateBlobCommitted(t, bbClient)
}

func TestBlobPutBlockListIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(t, err)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutBlockListIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(t, err)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	require.NoError(t, err)

	validateBlobCommitted(t, bbClient)
}

func TestBlobPutBlockListIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(t, err)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestBlobPutBlockListValidateData(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	require.NoError(t, err)

	resp, err := bbClient.Download(ctx, nil)
	require.NoError(t, err)
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, string(data), blockBlobDefaultData)
}

func TestBlobPutBlockListModifyBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(t, testName)
	defer deleteContainer(t, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	require.NoError(t, err)

	_, err = bbClient.StageBlock(ctx, "0001", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(t, err)
	_, err = bbClient.StageBlock(ctx, "0010", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(t, err)
	_, err = bbClient.StageBlock(ctx, "0011", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(t, err)
	_, err = bbClient.StageBlock(ctx, "0100", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(t, err)

	_, err = bbClient.CommitBlockList(ctx, []string{"0001", "0011"}, nil)
	require.NoError(t, err)

	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Len(t, resp.BlockList.CommittedBlocks, 2)
	committed := resp.BlockList.CommittedBlocks
	require.Equal(t, *(committed[0].Name), "0001")
	require.Equal(t, *(committed[1].Name), "0011")
	require.Nil(t, resp.BlockList.UncommittedBlocks)
}

func TestSetTierOnBlobUpload(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	for _, tier := range []AccessTier{AccessTierArchive, AccessTierCool, AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		uploadBlockBlobOptions := UploadBlockBlobOptions{
			HTTPHeaders: &basicHeaders,
			Tier:        &tier,
		}
		_, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &uploadBlockBlobOptions)
		require.NoError(t, err)

		resp, err := bbClient.GetProperties(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, *resp.AccessTier, string(tier))
	}
}

func TestBlobSetTierOnCommit(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := "test" + generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	for _, tier := range []AccessTier{AccessTierCool, AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
		_, err := bbClient.StageBlock(ctx, blockID, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
		require.NoError(t, err)

		_, err = bbClient.CommitBlockList(ctx, []string{blockID}, &CommitBlockListOptions{
			Tier: &tier,
		})
		require.NoError(t, err)

		resp, err := bbClient.GetBlockList(ctx, BlockListTypeCommitted, nil)
		require.NoError(t, err)
		require.NotNil(t, resp.BlockList)
		require.NotNil(t, resp.BlockList.CommittedBlocks)
		require.Nil(t, resp.BlockList.UncommittedBlocks)
		require.Len(t, resp.BlockList.CommittedBlocks, 1)
	}
}

func TestSetTierOnCopyBlockBlobFromURL(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	const contentSize = 4 * 1024 * 1024 // 4 MB
	contentReader, _ := getRandomDataAndReader(contentSize)

	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient(generateBlobName(testName))

	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, internal.NopCloser(contentReader), &UploadBlockBlobOptions{Tier: &tier})
	require.NoError(t, err)
	require.Equal(t, uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	expiryTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	require.NoError(t, err)

	credential, err := getCredential(testAccountDefault)
	if err != nil {
		t.Fatal("Couldn't fetch credential because " + err.Error())
	}
	sasQueryParams, err := AccountSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   AccountSASPermissions{Read: true, List: true}.String(),
		Services:      AccountSASServices{Blob: true}.String(),
		ResourceTypes: AccountSASResourceTypes{Container: true, Object: true}.String(),
	}.Sign(credential)
	require.NoError(t, err)

	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.URL()

	for _, tier := range []AccessTier{AccessTierArchive, AccessTierCool, AccessTierHot} {
		destBlobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		destBlob := containerClient.NewBlockBlobClient(generateBlobName(destBlobName))

		copyBlockBlobFromURLOptions := CopyBlockBlobFromURLOptions{
			Tier:     &tier,
			Metadata: map[string]string{"foo": "bar"},
		}
		resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions)
		require.NoError(t, err)
		require.Equal(t, resp.RawResponse.StatusCode, 202)
		require.Equal(t, *resp.CopyStatus, "success")

		destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, *destBlobPropResp.AccessTier, string(tier))
	}
}

func TestSetTierOnStageBlockFromURL(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))
	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))
	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, &UploadBlockBlobOptions{Tier: &tier})
	require.NoError(t, err)
	require.Equal(t, uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	offset1, count1 := int64(0), int64(4*1024)
	options1 := StageBlockFromURLOptions{
		Offset: &offset1,
		Count:  &count1,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	require.NoError(t, err)
	require.Equal(t, stageResp1.RawResponse.StatusCode, 201)
	require.Nil(t, stageResp1.ContentMD5)
	require.NotEqual(t, *stageResp1.RequestID, "")
	require.NotEqual(t, *stageResp1.Version, "")
	require.Equal(t, stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset: &offset2,
		Count:  &count2,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	require.NoError(t, err)
	require.Equal(t, stageResp2.RawResponse.StatusCode, 201)
	require.NotEqual(t, stageResp2.ContentMD5, "")
	require.NotEqual(t, stageResp2.RequestID, "")
	require.NotEqual(t, stageResp2.Version, "")
	require.Equal(t, stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Equal(t, blockList.RawResponse.StatusCode, 200)
	require.NotNil(t, blockList.BlockList)
	require.Nil(t, blockList.BlockList.CommittedBlocks)
	require.NotNil(t, blockList.BlockList.UncommittedBlocks)
	require.Len(t, blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &CommitBlockListOptions{
		Tier: &tier,
	})
	require.NoError(t, err)
	require.Equal(t, listResp.RawResponse.StatusCode, 201)
	require.NotNil(t, listResp.LastModified)
	require.Equal(t, (*listResp.LastModified).IsZero(), false)
	require.NotNil(t, listResp.ETag)
	require.NotNil(t, listResp.RequestID)
	require.NotNil(t, listResp.Version)
	require.NotNil(t, listResp.Date)
	require.Equal(t, (*listResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, content)

	// Get properties to validate the tier
	destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *destBlobPropResp.AccessTier, string(tier))
}

func TestSetStandardBlobTierWithRehydratePriority(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	standardTier, rehydrateTier, rehydratePriority := AccessTierArchive, AccessTierCool, RehydratePriorityStandard
	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, bbName, containerClient)

	_, err = bbClient.SetTier(ctx, standardTier, &SetTierOptions{
		RehydratePriority: &rehydratePriority,
	})
	require.NoError(t, err)

	getResp1, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp1.AccessTier, string(standardTier))

	_, err = bbClient.SetTier(ctx, rehydrateTier, nil)
	require.NoError(t, err)

	getResp2, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToCool))
}

func TestRehydrateStatus(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName1 := "rehydration_test_blob_1"
	blobName2 := "rehydration_test_blob_2"

	bbClient1 := getBlockBlobClient(blobName1, containerClient)
	reader1, _ := generateData(1024)
	_, err = bbClient1.Upload(ctx, reader1, nil)
	require.NoError(t, err)
	_, err = bbClient1.SetTier(ctx, AccessTierArchive, nil)
	require.NoError(t, err)
	_, err = bbClient1.SetTier(ctx, AccessTierCool, nil)
	require.NoError(t, err)

	getResp1, err := bbClient1.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp1.AccessTier, string(AccessTierArchive))
	require.Equal(t, *getResp1.ArchiveStatus, string(ArchiveStatusRehydratePendingToCool))

	pager := containerClient.ListBlobsFlat(nil)
	var blobs []*BlobItemInternal
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		blobs = append(blobs, resp.ListBlobsFlatSegmentResponse.Segment.BlobItems...)
	}
	require.Nil(t, pager.Err())
	require.GreaterOrEqual(t, len(blobs), 1)
	require.Equal(t, *blobs[0].Properties.AccessTier, AccessTierArchive)
	require.Equal(t, *blobs[0].Properties.ArchiveStatus, ArchiveStatusRehydratePendingToCool)

	// ------------------------------------------

	bbClient2 := getBlockBlobClient(blobName2, containerClient)
	reader2, _ := generateData(1024)
	_, err = bbClient2.Upload(ctx, reader2, nil)
	require.NoError(t, err)
	_, err = bbClient2.SetTier(ctx, AccessTierArchive, nil)
	require.NoError(t, err)
	_, err = bbClient2.SetTier(ctx, AccessTierHot, nil)
	require.NoError(t, err)

	getResp2, err := bbClient2.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp2.AccessTier, string(AccessTierArchive))
	require.Equal(t, *getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToHot))
}

func TestCopyBlobWithRehydratePriority(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	sourceBlobName := generateBlobName(testName)
	sourceBBClient := createNewBlockBlob(t, sourceBlobName, containerClient)

	blobTier, rehydratePriority := AccessTierArchive, RehydratePriorityHigh

	copyBlobName := "copy" + sourceBlobName
	destBBClient := getBlockBlobClient(copyBlobName, containerClient)
	_, err = destBBClient.StartCopyFromURL(ctx, sourceBBClient.URL(), &StartCopyBlobOptions{
		RehydratePriority: &rehydratePriority,
		Tier:              &blobTier,
	})
	require.NoError(t, err)

	getResp1, err := destBBClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp1.AccessTier, string(blobTier))

	_, err = destBBClient.SetTier(ctx, AccessTierHot, nil)
	require.NoError(t, err)

	getResp2, err := destBBClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToHot))
}
