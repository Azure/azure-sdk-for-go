//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
)

//	func (s *azblobTestSuite) TestStageGetBlocks() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//	//		svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := generateContainerName(testName)
//		containerClient := createNewContainer(_require, containerName, svcClient)
//		defer deleteContainer(_require, containerClient)
//
//		blobName := generateBlobName(testName)
//		bbClient := containerClient.NewBlockBlobClient(blobName)
//
//		data := []string{"Azure ", "Storage ", "Block ", "Blob."}
//		base64BlockIDs := make([]string, len(data))
//
//		for index, d := range data {
//			base64BlockIDs[index] = blockIDIntToBase64(index)
//			io.NopCloser(strings.NewReader("hello world"))
//			putResp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(d)), nil)
//			_require.Nil(err)
//			//_require.Equal(putResp.RawResponse.StatusCode, 201)
//			_require.Nil(putResp.ContentMD5)
//			_require.NotNil(putResp.RequestID)
//			_require.NotNil(putResp.Version)
//			_require.NotNil(putResp.Date)
//			_require.Equal((*putResp.Date).IsZero(), false)
//		}
//
//		blockList, err := bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
//		_require.Nil(err)
//		// _require.Equal(blockList.RawResponse.StatusCode, 200)
//		_require.Nil(blockList.LastModified)
//		_require.Nil(blockList.ETag)
//		_require.NotNil(blockList.ContentType)
//		_require.Nil(blockList.BlobContentLength)
//		_require.NotNil(blockList.RequestID)
//		_require.NotNil(blockList.Version)
//		_require.NotNil(blockList.Date)
//		_require.Equal((*blockList.Date).IsZero(), false)
//		_require.NotNil(blockList.BlockList)
//		_require.Nil(blockList.BlockList.CommittedBlocks)
//		_require.NotNil(blockList.BlockList.UncommittedBlocks)
//		_require.Len(blockList.BlockList.UncommittedBlocks, len(data))
//
//		listResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
//		_require.Nil(err)
//		// _require.Equal(listResp.RawResponse.StatusCode,  201)
//		_require.NotNil(listResp.LastModified)
//		_require.Equal((*listResp.LastModified).IsZero(), false)
//		_require.NotNil(listResp.ETag)
//		_require.NotNil(listResp.RequestID)
//		_require.NotNil(listResp.Version)
//		_require.NotNil(listResp.Date)
//		_require.Equal((*listResp.Date).IsZero(), false)
//
//		blockList, err = bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
//		_require.Nil(err)
//		// _require.Equal(blockList.RawResponse.StatusCode, 200)
//		_require.NotNil(blockList.LastModified)
//		_require.Equal((*blockList.LastModified).IsZero(), false)
//		_require.NotNil(blockList.ETag)
//		_require.NotNil(blockList.ContentType)
//		_require.Equal(*blockList.BlobContentLength, int64(25))
//		_require.NotNil(blockList.RequestID)
//		_require.NotNil(blockList.Version)
//		_require.NotNil(blockList.Date)
//		_require.Equal((*blockList.Date).IsZero(), false)
//		_require.NotNil(blockList.BlockList)
//		_require.NotNil(blockList.BlockList.CommittedBlocks)
//		_require.Nil(blockList.BlockList.UncommittedBlocks)
//		_require.Len(blockList.BlockList.CommittedBlocks, len(data))
//	}
//
// //nolint
//
//	func (s *azblobUnrecordedTestSuite) TestStageBlockFromURL() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := generateContainerName(testName)
//		containerClient := createNewContainer(_require, containerName, svcClient)
//		defer deleteContainer(_require, containerClient)
//
//		contentSize := 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		body := bytes.NewReader(content)
//		rsc := streaming.NopCloser(body)
//
//		ctx := context.Background() // Use default Background context
//		srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))
//
//		destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))
//
//		// Prepare source bbClient for copy.
//		_, err = srcBlob.Upload(ctx, rsc, nil)
//		_require.Nil(err)
//		//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get source blob url with SAS for StageFromURL.
//		srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//		credential, err := getGenericCredential(nil, testAccountDefault)
//		_require.Nil(err)
//
//		srcBlobParts.SAS, err = BlobSASSignatureValues{
//			Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//			ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//			ContainerName: srcBlobParts.ContainerName,
//			BlobName:      srcBlobParts.BlobName,
//			Permissions:   BlobSASPermissions{Read: true}.String(),
//		}.Sign(credential)
//		_require.Nil(err)
//
//		srcBlobURLWithSAS := srcBlobParts.URL()
//
//		// Stage blocks from URL.
//		blockIDs := generateBlockIDsList(2)
//
//		stageResp1, err := destBlob.StageBlockFromURL(ctx, blockIDs[0], srcBlobURLWithSAS, 0, &BlockBlobStageBlockFromURLOptions{
//			Offset: to.Ptr[int64](0),
//			Count:  to.Ptr(int64(contentSize / 2)),
//		})
//		_require.Nil(err)
//		// _require.Equal(stageResp1.RawResponse.StatusCode, 201)
//		_require.NotEqual(stageResp1.ContentMD5, "")
//		_require.NotEqual(stageResp1.RequestID, "")
//		_require.NotEqual(stageResp1.Version, "")
//		_require.Equal(stageResp1.Date.IsZero(), false)
//
//		stageResp2, err := destBlob.StageBlockFromURL(ctx, blockIDs[1], srcBlobURLWithSAS, 0, &BlockBlobStageBlockFromURLOptions{
//			Offset: to.Ptr(int64(contentSize / 2)),
//			Count:  to.Ptr(int64(CountToEnd)),
//		})
//		_require.Nil(err)
//		// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//		_require.NotEqual(stageResp2.ContentMD5, "")
//		_require.NotEqual(stageResp2.RequestID, "")
//		_require.NotEqual(stageResp2.Version, "")
//		_require.Equal(stageResp2.Date.IsZero(), false)
//
//		// Check block list.
//		blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
//		_require.Nil(err)
//		// _require.Equal(blockList.RawResponse.StatusCode, 200)
//		_require.NotNil(blockList.BlockList)
//		_require.Nil(blockList.BlockList.CommittedBlocks)
//		_require.NotNil(blockList.BlockList.UncommittedBlocks)
//		_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//		// Commit block list.
//		listResp, err := destBlob.CommitBlockList(context.Background(), blockIDs, nil)
//		_require.Nil(err)
//		// _require.Equal(listResp.RawResponse.StatusCode,  201)
//		_require.NotNil(listResp.LastModified)
//		_require.Equal((*listResp.LastModified).IsZero(), false)
//		_require.NotNil(listResp.ETag)
//		_require.NotNil(listResp.RequestID)
//		_require.NotNil(listResp.Version)
//		_require.NotNil(listResp.Date)
//		_require.Equal((*listResp.Date).IsZero(), false)
//
//		// Check data integrity through downloading.
//		downloadResp, err := destBlob.BlobClient.DownloadStream(ctx, nil)
//		_require.Nil(err)
//		destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//		_require.Nil(err)
//		_require.EqualValues(destData, content)
//	}
//
// //nolint
//
//	func (s *azblobUnrecordedTestSuite) TestCopyBlockBlobFromURL() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := generateContainerName(testName)
//		containerClient := createNewContainer(_require, containerName, svcClient)
//		defer deleteContainer(_require, containerClient)
//
//		const contentSize = 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		contentMD5 := md5.Sum(content)
//		body := bytes.NewReader(content)
//		ctx := context.Background()
//
//		srcBlob := containerClient.NewBlockBlobClient("srcblob")
//		destBlob := containerClient.NewBlockBlobClient("destblob")
//
//		// Prepare source bbClient for copy.
//		_, err = srcBlob.Upload(ctx, streaming.NopCloser(body), nil)
//		_require.Nil(err)
//		//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get source blob url with SAS for StageFromURL.
//		srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//		credential, err := getGenericCredential(nil, testAccountDefault)
//		_require.Nil(err)
//
//		srcBlobParts.SAS, err = BlobSASSignatureValues{
//			Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//			ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//			ContainerName: srcBlobParts.ContainerName,
//			BlobName:      srcBlobParts.BlobName,
//			Permissions:   BlobSASPermissions{Read: true}.String(),
//		}.Sign(credential)
//		if err != nil {
//			s.T().Fatal(err)
//		}
//
//		srcBlobURLWithSAS := srcBlobParts.URL()
//
//		// Invoke copy bbClient from URL.
//		sourceContentMD5 := contentMD5[:]
//		resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &BlockBlobCopyFromURLOptions{
//			Metadata:         map[string]string{"foo": "bar"},
//			SourceContentMD5: sourceContentMD5,
//		})
//		_require.Nil(err)
//		// _require.Equal(resp.RawResponse.StatusCode, 202)
//		_require.NotNil(resp.ETag)
//		_require.NotNil(resp.RequestID)
//		_require.NotNil(resp.Version)
//		_require.NotNil(resp.Date)
//		_require.Equal((*resp.Date).IsZero(), false)
//		_require.NotNil(resp.CopyID)
//		_require.EqualValues(resp.ContentMD5, sourceContentMD5)
//		_require.Equal(*resp.CopyStatus, "success")
//
//		// Make sure the metadata got copied over
//		getPropResp, err := destBlob.GetProperties(ctx, nil)
//		_require.Nil(err)
//		metadata := getPropResp.Metadata
//		_require.NotNil(metadata)
//		_require.Len(metadata, 1)
//		_require.EqualValues(metadata, map[string]string{"Foo": "bar"})
//
//		// Check data integrity through downloading.
//		downloadResp, err := destBlob.DownloadStream(ctx, nil)
//		_require.Nil(err)
//		destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//		_require.Nil(err)
//		_require.EqualValues(destData, content)
//
//		// Edge case 1: Provide bad MD5 and make sure the copy fails
//		_, badMD5 := getRandomDataAndReader(16)
//		copyBlockBlobFromURLOptions1 := BlockBlobCopyFromURLOptions{
//			SourceContentMD5: badMD5,
//		}
//		resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
//		_require.NotNil(err)
//
//		// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
//		copyBlockBlobFromURLOptions2 := BlockBlobCopyFromURLOptions{
//			SourceContentMD5: sourceContentMD5,
//		}
//		resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
//		_require.Nil(err)
//		// _require.Equal(resp.RawResponse.StatusCode, 202)
//		_require.EqualValues(*resp.CopyStatus, "success")
//	}
//
// //nolint
//
//	func (s *azblobUnrecordedTestSuite) TestBlobSASQueryParamOverrideResponseHeaders() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := generateContainerName(testName)
//		containerClient := createNewContainer(_require, containerName, svcClient)
//		defer deleteContainer(_require, containerClient)
//
//		const contentSize = 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		body := bytes.NewReader(content)
//		//contentMD5 := md5.Sum(content)
//
//		ctx := context.Background()
//
//		bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))
//
//		_, err = bbClient.Upload(ctx, streaming.NopCloser(body), nil)
//		_require.Nil(err)
//		//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get blob url with SAS.
//		blobParts, _ := NewBlobURLParts(bbClient.URL())
//
//		cacheControlVal := "cache-control-override"
//		contentDispositionVal := "content-disposition-override"
//		contentEncodingVal := "content-encoding-override"
//		contentLanguageVal := "content-language-override"
//		contentTypeVal := "content-type-override"
//
//		credential, err := getGenericCredential(nil, testAccountDefault)
//		_require.Nil(err)
//		// Append User Delegation SAS token to URL
//		blobParts.SAS, err = BlobSASSignatureValues{
//			Protocol:           SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//			ExpiryTime:         time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//			ContainerName:      blobParts.ContainerName,
//			BlobName:           blobParts.BlobName,
//			Permissions:        BlobSASPermissions{Read: true}.String(),
//			CacheControl:       cacheControlVal,
//			ContentDisposition: contentDispositionVal,
//			ContentEncoding:    contentEncodingVal,
//			ContentLanguage:    contentLanguageVal,
//			ContentType:        contentTypeVal,
//		}.Sign(credential)
//		_require.Nil(err)
//
//		// Generate new bbClient client
//		blobURLWithSAS := blobParts.URL()
//		_require.NotNil(blobURLWithSAS)
//
//		blobClientWithSAS, err := NewBlockBlobClientWithNoCredential(blobURLWithSAS, nil)
//		_require.Nil(err)
//
//		gResp, err := blobClientWithSAS.GetProperties(ctx, nil)
//		_require.Nil(err)
//		_require.Equal(*gResp.CacheControl, cacheControlVal)
//		_require.Equal(*gResp.ContentDisposition, contentDispositionVal)
//		_require.Equal(*gResp.ContentEncoding, contentEncodingVal)
//		_require.Equal(*gResp.ContentLanguage, contentLanguageVal)
//		_require.Equal(*gResp.ContentType, contentTypeVal)
//	}
//
// nolint
func (s *azblobUnrecordedTestSuite) TestStageBlockWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName := generateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	// test put block with valid MD5 value
	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	putResp, err := bbClient.StageBlock(context.Background(), blockID1, rsc, &blockblob.StageBlockOptions{
		TransactionalContentMD5: contentMD5,
	})
	_require.Nil(err)
	//_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.EqualValues(putResp.ContentMD5, contentMD5)
	_require.NotNil(putResp.RequestID)
	_require.NotNil(putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	// test put block with bad MD5 value
	_, badContent := getRandomDataAndReader(contentSize)
	badMD5Value := md5.Sum(badContent)
	badContentMD5 := badMD5Value[:]

	_, _ = rsc.Seek(0, io.SeekStart)
	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	_, err = bbClient.StageBlock(context.Background(), blockID2, rsc, &blockblob.StageBlockOptions{
		TransactionalContentMD5: badContentMD5,
	})
	_require.NotNil(err)
	_require.Contains(err.Error(), bloberror.MD5Mismatch)
}

func (s *azblobTestSuite) TestBlobPutBlobHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, streaming.NopCloser(body), &blockblob.UploadOptions{
		HTTPHeaders: &basicHeaders,
	})
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	h := blob.ParseHTTPHeaders(resp)
	h.BlobContentMD5 = nil // the service generates a MD5 value, omit before comparing
	_require.EqualValues(h, basicHeaders)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataNotEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, streaming.NopCloser(body), &blockblob.UploadOptions{
		Metadata: basicMetadata,
	})
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	actualMetadata := resp.Metadata
	_require.NotNil(actualMetadata)
	_require.EqualValues(actualMetadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &blockblob.UploadOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	})
	_require.NotNil(err)
	_require.Contains(err.Error(), invalidHeaderErrorSubstring)
}

func (s *azblobTestSuite) TestBlobPutBlobIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, streaming.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, -10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, streaming.NopCloser(body), &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)
	validateUpload(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, streaming.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, 10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, streaming.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, 10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	_require.Nil(err)

	validateUpload(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, streaming.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, -10)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader(nil)), &uploadBlockBlobOptions)
	_ = err

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.Nil(err)

	validateUpload(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)

	ifMatch := "garbage"
	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &ifMatch,
			},
		},
	}
	_, err = bbClient.Upload(ctx, streaming.NopCloser(body), &uploadBlockBlobOptions)
	_require.NotNil(err)
	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	ifNoneMatch := "garbage"
	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &ifNoneMatch,
			},
		},
	}

	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	_require.Nil(err)

	validateUpload(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func validateBlobCommitted(_require *require.Assertions, bbClient *blockblob.Client) {
	resp, err := bbClient.GetBlockList(ctx, blockblob.BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Len(resp.BlockList.CommittedBlocks, 1)
}

func setupPutBlockListTest(t *testing.T, _require *require.Assertions, testName string) (*container.Client, *blockblob.Client, []string) {
	svcClient, err := getServiceClient(t, testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)

	blockIDs := generateBlockIDsList(1)
	_, err = bbClient.StageBlock(ctx, blockIDs[0], streaming.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_require.Nil(err)
	return containerClient, bbClient, blockIDs
}

func (s *azblobTestSuite) TestBlobPutBlockListHTTPHeadersEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, &blockblob.CommitBlockListOptions{
		HTTPHeaders: &blob.HTTPHeaders{BlobContentDisposition: &blobContentDisposition},
	})
	_require.Nil(err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.ContentDisposition)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	getPropertyResp, err := containerClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(getPropertyResp.Date)

	currentTime := getRelativeTimeFromAnchor(getPropertyResp.Date, 10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_ = err

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, 10)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag}},
	})
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	eTag := "garbage"
	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	eTag := "garbage"
	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListValidateData() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	_require.Nil(err)

	resp, err := bbClient.DownloadStream(ctx, nil)
	_require.Nil(err)
	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobPutBlockListModifyBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer deleteContainer(_require, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	_require.Nil(err)

	_, err = bbClient.StageBlock(ctx, "0001", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)
	_, err = bbClient.StageBlock(ctx, "0010", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)
	_, err = bbClient.StageBlock(ctx, "0011", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)
	_, err = bbClient.StageBlock(ctx, "0100", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)

	_, err = bbClient.CommitBlockList(ctx, []string{"0001", "0011"}, nil)
	_require.Nil(err)

	resp, err := bbClient.GetBlockList(ctx, blockblob.BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Len(resp.BlockList.CommittedBlocks, 2)
	committed := resp.BlockList.CommittedBlocks
	_require.Equal(*(committed[0].Name), "0001")
	_require.Equal(*(committed[1].Name), "0011")
	_require.Nil(resp.BlockList.UncommittedBlocks)
}

func (s *azblobTestSuite) TestSetTierOnBlobUpload() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	for _, tier := range []blob.AccessTier{blob.AccessTierArchive, blob.AccessTierCool, blob.AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		uploadBlockBlobOptions := blockblob.UploadOptions{
			HTTPHeaders: &basicHeaders,
			Tier:        &tier,
		}
		_, err := bbClient.Upload(ctx, streaming.NopCloser(strings.NewReader(blockBlobDefaultData)), &uploadBlockBlobOptions)
		_require.Nil(err)

		resp, err := bbClient.GetProperties(ctx, nil)
		_require.Nil(err)
		_require.Equal(*resp.AccessTier, string(tier))
	}
}

func (s *azblobTestSuite) TestBlobSetTierOnCommit() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := "test" + generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	for _, tier := range []blob.AccessTier{blob.AccessTierCool, blob.AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
		_, err := bbClient.StageBlock(ctx, blockID, streaming.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
		_require.Nil(err)

		_, err = bbClient.CommitBlockList(ctx, []string{blockID}, &blockblob.CommitBlockListOptions{
			Tier: &tier,
		})
		_require.Nil(err)

		resp, err := bbClient.GetBlockList(ctx, blockblob.BlockListTypeCommitted, nil)
		_require.Nil(err)
		_require.NotNil(resp.BlockList)
		_require.NotNil(resp.BlockList.CommittedBlocks)
		_require.Nil(resp.BlockList.UncommittedBlocks)
		_require.Len(resp.BlockList.CommittedBlocks, 1)
	}
}

// nolint
func (s *azblobUnrecordedTestSuite) TestSetTierOnCopyBlockBlobFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	const contentSize = 4 * 1024 * 1024 // 4 MB
	contentReader, _ := getRandomDataAndReader(contentSize)

	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient(generateBlobName(testName))

	tier := blob.AccessTierCool
	_, err = srcBlob.Upload(ctx, streaming.NopCloser(contentReader), &blockblob.UploadOptions{Tier: &tier})
	_require.Nil(err)
	//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	expiryTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)

	credential, err := getGenericCredential(testAccountDefault)
	if err != nil {
		s.T().Fatal("Couldn't fetch credential because " + err.Error())
	}
	sasQueryParams, err := service.SASSignatureValues{
		Protocol:      service.SASProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   to.Ptr(service.SASPermissions{Read: true, List: true}).String(),
		Services:      to.Ptr(service.SASServices{Blob: true}).String(),
		ResourceTypes: to.Ptr(service.SASResourceTypes{Container: true, Object: true}).String(),
	}.Sign(credential)
	_require.Nil(err)

	srcBlobParts, _ := azblob.ParseURL(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	for _, tier := range []blob.AccessTier{blob.AccessTierArchive, blob.AccessTierCool, blob.AccessTierHot} {
		destBlobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		destBlob := containerClient.NewBlockBlobClient(generateBlobName(destBlobName))

		copyBlockBlobFromURLOptions := blob.CopyFromURLOptions{
			Tier:     &tier,
			Metadata: map[string]string{"foo": "bar"},
		}
		resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 202)
		_require.Equal(*resp.CopyStatus, "success")

		destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
		_require.Nil(err)
		_require.Equal(*destBlobPropResp.AccessTier, string(tier))
	}
}

////nolint
//func (s *azblobUnrecordedTestSuite) TestSetTierOnStageBlockFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	contentSize := 8 * 1024 // 8 KB
//	content := make([]byte, contentSize)
//	body := bytes.NewReader(content)
//	rsc := streaming.NopCloser(body)
//	ctx := context.Background()
//	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))
//	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))
//	tier := AccessTierCool
//	_, err = srcBlob.Upload(ctx, rsc, &blockblob.UploadOptions{Tier: &tier})
//	_require.Nil(err)
//	//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//	credential, err := getGenericCredential(nil, testAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = blob.SASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Stage blocks from URL.
//	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
//	offset1, count1 := int64(0), int64(4*1024)
//	options1 := blockblob.StageBlockFromURLOptions{
//		Offset: &offset1,
//		Count:  &count1,
//	}
//	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.Nil(err)
//	// _require.Equal(stageResp1.RawResponse.StatusCode, 201)
//	_require.Nil(stageResp1.ContentMD5)
//	_require.NotEqual(*stageResp1.RequestID, "")
//	_require.NotEqual(*stageResp1.Version, "")
//	_require.Equal(stageResp1.Date.IsZero(), false)
//
//	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
//	offset2, count2 := int64(4*1024), int64(blob.CountToEnd)
//	options2 := blockblob.StageBlockFromURLOptions{
//		Offset: &offset2,
//		Count:  &count2,
//	}
//	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.Nil(err)
//	// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp2.ContentMD5, "")
//	_require.NotEqual(stageResp2.RequestID, "")
//	_require.NotEqual(stageResp2.Version, "")
//	_require.Equal(stageResp2.Date.IsZero(), false)
//
//	// Check block list.
//	blockList, err := destBlob.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
//	_require.Nil(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.NotNil(blockList.BlockList.UncommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	// Commit block list.
//	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &blockblob.CommitBlockListOptions{
//		Tier: &tier,
//	})
//	_require.Nil(err)
//	// _require.Equal(listResp.RawResponse.StatusCode,  201)
//	_require.NotNil(listResp.LastModified)
//	_require.Equal((*listResp.LastModified).IsZero(), false)
//	_require.NotNil(listResp.ETag)
//	_require.NotNil(listResp.RequestID)
//	_require.NotNil(listResp.Version)
//	_require.NotNil(listResp.Date)
//	_require.Equal((*listResp.Date).IsZero(), false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.DownloadStream(ctx, nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(nil))
//	_require.Nil(err)
//	_require.EqualValues(destData, content)
//
//	// Get properties to validate the tier
//	destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.Equal(*destBlobPropResp.AccessTier, string(tier))
//}

func (s *azblobTestSuite) TestSetStandardBlobTierWithRehydratePriority() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	standardTier, rehydrateTier, rehydratePriority := blob.AccessTierArchive, blob.AccessTierCool, blob.RehydratePriorityStandard
	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, bbName, containerClient)

	_, err = bbClient.SetTier(ctx, standardTier, &blob.SetTierOptions{
		RehydratePriority: &rehydratePriority,
	})
	_require.Nil(err)

	getResp1, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp1.AccessTier, string(standardTier))

	_, err = bbClient.SetTier(ctx, rehydrateTier, nil)
	_require.Nil(err)

	getResp2, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToCool))
}

func (s *azblobTestSuite) TestRehydrateStatus() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blobName1 := "rehydration_test_blob_1"
	blobName2 := "rehydration_test_blob_2"

	bbClient1 := getBlockBlobClient(blobName1, containerClient)
	reader1, _ := generateData(1024)
	_, err = bbClient1.Upload(ctx, reader1, nil)
	_require.Nil(err)
	_, err = bbClient1.SetTier(ctx, blob.AccessTierArchive, nil)
	_require.Nil(err)
	_, err = bbClient1.SetTier(ctx, blob.AccessTierCool, nil)
	_require.Nil(err)

	getResp1, err := bbClient1.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp1.AccessTier, string(blob.AccessTierArchive))
	_require.Equal(*getResp1.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToCool))

	pager := containerClient.NewListBlobsFlatPager(nil)
	var blobs []*container.BlobItem
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		blobs = append(blobs, resp.ListBlobsFlatSegmentResponse.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.GreaterOrEqual(len(blobs), 1)
	_require.Equal(*blobs[0].Properties.AccessTier, blob.AccessTierArchive)
	_require.Equal(*blobs[0].Properties.ArchiveStatus, blob.ArchiveStatusRehydratePendingToCool)

	// ------------------------------------------

	bbClient2 := getBlockBlobClient(blobName2, containerClient)
	reader2, _ := generateData(1024)
	_, err = bbClient2.Upload(ctx, reader2, nil)
	_require.Nil(err)
	_, err = bbClient2.SetTier(ctx, blob.AccessTierArchive, nil)
	_require.Nil(err)
	_, err = bbClient2.SetTier(ctx, blob.AccessTierHot, nil)
	_require.Nil(err)

	getResp2, err := bbClient2.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp2.AccessTier, string(blob.AccessTierArchive))
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToHot))
}

func (s *azblobTestSuite) TestCopyBlobWithRehydratePriority() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	sourceBlobName := generateBlobName(testName)
	sourceBBClient := createNewBlockBlob(_require, sourceBlobName, containerClient)

	blobTier, rehydratePriority := blob.AccessTierArchive, blob.RehydratePriorityHigh

	copyBlobName := "copy" + sourceBlobName
	destBBClient := getBlockBlobClient(copyBlobName, containerClient)
	_, err = destBBClient.StartCopyFromURL(ctx, sourceBBClient.URL(), &blob.StartCopyFromURLOptions{
		RehydratePriority: &rehydratePriority,
		Tier:              &blobTier,
	})
	_require.Nil(err)

	getResp1, err := destBBClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp1.AccessTier, string(blobTier))

	_, err = destBBClient.SetTier(ctx, blob.AccessTierHot, nil)
	_require.Nil(err)

	getResp2, err := destBBClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToHot))
}

func (s *azblobTestSuite) TestBlobServiceClientDelete() {
	_require := require.New(s.T())
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	code := 404
	runTestRequiringServiceProperties(_require, svcClient, string(rune(code)), enableSoftDelete, testBlobServiceClientDeleteImpl, disableSoftDelete)
}

func setAndCheckBlockBlobTier(_require *require.Assertions, bbClient *blockblob.Client, tier blob.AccessTier) {
	_, err := bbClient.SetTier(ctx, tier, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.AccessTier, string(tier))
}

func (s *azblobTestSuite) TestBlobSetTierAllTiersOnBlockBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(s.T(), testAccountDefault, nil)
	_require.NoError(err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_require, blockBlobName, containerClient)

	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierHot)
	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierCool)
	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierArchive)

}
