//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	suite.Run(t, &BlockBlobRecordedTestsSuite{})
	//suite.Run(t, &BlockBlobUnrecordedTestsSuite{})
}

// nolint
func (s *BlockBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

// nolint
func (s *BlockBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

// nolint
func (s *BlockBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

// nolint
func (s *BlockBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type BlockBlobRecordedTestsSuite struct {
	suite.Suite
}

type BlockBlobUnrecordedTestsSuite struct {
	suite.Suite
}

//	func (s *BlockBlobRecordedTestsSuite) TestStageGetBlocks() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//	//		svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//		blobName := testcommon.GenerateBlobName(testName)
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
//	func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockFromURL() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//		contentSize := 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		body := bytes.NewReader(content)
//		rsc := streaming.NopCloser(body)
//
//		ctx := context.Background() // Use default Background context
//		srcBlob := containerClient.NewBlockBlobClient("src" + testcommon.GenerateBlobName(testName))
//
//		destBlob := containerClient.NewBlockBlobClient("dst" + testcommon.GenerateBlobName(testName))
//
//		// Prepare source bbClient for copy.
//		_, err = srcBlob.Upload(context.Background(), rsc, nil)
//		_require.Nil(err)
//		//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get source blob url with SAS for StageFromURL.
//		srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//		credential, err := testcommon.GetGenericCredential(nil, testcommon.TestAccountDefault)
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
//		blockIDs := testcommon.GenerateBlockIDsList(2)
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
//		downloadResp, err := destBlob.BlobClient.DownloadStream(context.Background(), nil)
//		_require.Nil(err)
//		destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//		_require.Nil(err)
//		_require.EqualValues(destData, content)
//	}
//
// //nolint
//
//	func (s *BlockBlobUnrecordedTestsSuite) TestCopyBlockBlobFromURL() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
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
//		_, err = srcBlob.Upload(context.Background(), streaming.NopCloser(body), nil)
//		_require.Nil(err)
//		//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get source blob url with SAS for StageFromURL.
//		srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//		credential, err := testcommon.GetGenericCredential(nil, testcommon.TestAccountDefault)
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
//		getPropResp, err := destBlob.GetProperties(context.Background(), nil)
//		_require.Nil(err)
//		metadata := getPropResp.Metadata
//		_require.NotNil(metadata)
//		_require.Len(metadata, 1)
//		_require.EqualValues(metadata, map[string]string{"Foo": "bar"})
//
//		// Check data integrity through downloading.
//		downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
//		_require.Nil(err)
//		destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//		_require.Nil(err)
//		_require.EqualValues(destData, content)
//
//		// Edge case 1: Provide bad MD5 and make sure the copy fails
//		_, badMD5 := testcommon.GetRandomDataAndReader(16)
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
//	func (s *BlockBlobUnrecordedTestsSuite) TestBlobSASQueryParamOverrideResponseHeaders() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//		const contentSize = 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		body := bytes.NewReader(content)
//		//contentMD5 := md5.Sum(content)
//
//		ctx := context.Background()
//
//		bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))
//
//		_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), nil)
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
//		credential, err := testcommon.GetGenericCredential(nil, testcommon.TestAccountDefault)
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
//		gResp, err := blobClientWithSAS.GetProperties(context.Background(), nil)
//		_require.Nil(err)
//		_require.Equal(*gResp.CacheControl, cacheControlVal)
//		_require.Equal(*gResp.ContentDisposition, contentDispositionVal)
//		_require.Equal(*gResp.ContentEncoding, contentEncodingVal)
//		_require.Equal(*gResp.ContentLanguage, contentLanguageVal)
//		_require.Equal(*gResp.ContentType, contentTypeVal)
//	}
//
// nolint
func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
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
	_, badContent := testcommon.GetRandomDataAndReader(contentSize)
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

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &blockblob.UploadOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
	})
	_require.Nil(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	h := blob.ParseHTTPHeaders(resp)
	h.BlobContentMD5 = nil // the service generates a MD5 value, omit before comparing
	_require.EqualValues(h, testcommon.BasicHeaders)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobMetadataNotEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &blockblob.UploadOptions{
		Metadata: testcommon.BasicMetadata,
	})
	_require.Nil(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	actualMetadata := resp.Metadata
	_require.NotNil(actualMetadata)
	_require.EqualValues(actualMetadata, testcommon.BasicMetadata)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, &blockblob.UploadOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	})
	_require.NotNil(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, -10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)
	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, 10)

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

	_, err = bbClient.Upload(context.Background(), rsc, &uploadBlockBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, 10)

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
	_, err = bbClient.Upload(context.Background(), rsc, &uploadBlockBlobOptions)
	_require.Nil(err)

	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	//_require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, -10)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(nil)), &uploadBlockBlobOptions)
	_ = err

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.Nil(err)

	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(context.Background(), nil)
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
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &uploadBlockBlobOptions)
	_require.NotNil(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(context.Background(), nil)
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

	_, err = bbClient.Upload(context.Background(), rsc, &uploadBlockBlobOptions)
	_require.Nil(err)

	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func validateBlobCommitted(_require *require.Assertions, bbClient *blockblob.Client) {
	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Len(resp.BlockList.CommittedBlocks, 1)
}

func setupPutBlockListTest(t *testing.T, _require *require.Assertions, testName string) (*container.Client, *blockblob.Client, []string) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	blockIDs := testcommon.GenerateBlockIDsList(1)
	_, err = bbClient.StageBlock(context.Background(), blockIDs[0], streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Nil(err)
	return containerClient, bbClient, blockIDs
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListHTTPHeadersEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		HTTPHeaders: &blob.HTTPHeaders{BlobContentDisposition: &testcommon.BlobContentDisposition},
	})
	_require.Nil(err)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.ContentDisposition)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertyResp, err := containerClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(getPropertyResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertyResp.Date, 10)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_ = err

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(commitBlockListResp.Date, 10)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag}},
	})
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	eTag := "garbage"
	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	eTag := "garbage"
	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)
	_require.Nil(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.Nil(err)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListValidateData() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil)
	_require.Nil(err)

	resp, err := bbClient.DownloadStream(context.Background(), nil)
	_require.Nil(err)
	data, err := io.ReadAll(resp.Body)
	_require.Nil(err)
	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListModifyBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil)
	_require.Nil(err)

	_, err = bbClient.StageBlock(context.Background(), "0001", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)
	_, err = bbClient.StageBlock(context.Background(), "0010", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)
	_, err = bbClient.StageBlock(context.Background(), "0011", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)
	_, err = bbClient.StageBlock(context.Background(), "0100", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.Nil(err)

	_, err = bbClient.CommitBlockList(context.Background(), []string{"0001", "0011"}, nil)
	_require.Nil(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Len(resp.BlockList.CommittedBlocks, 2)
	committed := resp.BlockList.CommittedBlocks
	_require.Equal(*(committed[0].Name), "0001")
	_require.Equal(*(committed[1].Name), "0011")
	_require.Nil(resp.BlockList.UncommittedBlocks)
}

func (s *BlockBlobRecordedTestsSuite) TestSetTierOnBlobUpload() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	for _, tier := range []blob.AccessTier{blob.AccessTierArchive, blob.AccessTierCool, blob.AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + testcommon.GenerateBlobName(testName)
		bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

		uploadBlockBlobOptions := blockblob.UploadOptions{
			HTTPHeaders: &testcommon.BasicHeaders,
			Tier:        &tier,
		}
		_, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &uploadBlockBlobOptions)
		_require.Nil(err)

		resp, err := bbClient.GetProperties(context.Background(), nil)
		_require.Nil(err)
		_require.Equal(*resp.AccessTier, string(tier))
	}
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSetTierOnCommit() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := "test" + testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	for _, tier := range []blob.AccessTier{blob.AccessTierCool, blob.AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + testcommon.GenerateBlobName(testName)
		bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
		_, err := bbClient.StageBlock(context.Background(), blockID, streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
		_require.Nil(err)

		_, err = bbClient.CommitBlockList(context.Background(), []string{blockID}, &blockblob.CommitBlockListOptions{
			Tier: &tier,
		})
		_require.Nil(err)

		resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeCommitted, nil)
		_require.Nil(err)
		_require.NotNil(resp.BlockList)
		_require.NotNil(resp.BlockList.CommittedBlocks)
		_require.Nil(resp.BlockList.UncommittedBlocks)
		_require.Len(resp.BlockList.CommittedBlocks, 1)
	}
}

// nolint
func (s *BlockBlobUnrecordedTestsSuite) TestSetTierOnCopyBlockBlobFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	const contentSize = 4 * 1024 * 1024 // 4 MB
	contentReader, _ := testcommon.GetRandomDataAndReader(contentSize)

	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	tier := blob.AccessTierCool
	_, err = srcBlob.Upload(context.Background(), streaming.NopCloser(contentReader), &blockblob.UploadOptions{Tier: &tier})
	_require.Nil(err)
	//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	expiryTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)

	credential, err := testcommon.GetGenericCredential(testcommon.TestAccountDefault)
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
		destBlobName := strings.ToLower(string(tier)) + testcommon.GenerateBlobName(testName)
		destBlob := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(destBlobName))

		copyBlockBlobFromURLOptions := blob.CopyFromURLOptions{
			Tier:     &tier,
			Metadata: map[string]string{"foo": "bar"},
		}
		resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 202)
		_require.Equal(*resp.CopyStatus, "success")

		destBlobPropResp, err := destBlob.GetProperties(context.Background(), nil)
		_require.Nil(err)
		_require.Equal(*destBlobPropResp.AccessTier, string(tier))
	}
}

////nolint
//func (s *BlockBlobUnrecordedTestsSuite) TestSetTierOnStageBlockFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 8 * 1024 // 8 KB
//	content := make([]byte, contentSize)
//	body := bytes.NewReader(content)
//	rsc := streaming.NopCloser(body)
//	ctx := context.Background()
//	srcBlob := containerClient.NewBlockBlobClient("src" + testcommon.GenerateBlobName(testName))
//	destBlob := containerClient.NewBlockBlobClient("dst" + testcommon.GenerateBlobName(testName))
//	tier := AccessTierCool
//	_, err = srcBlob.Upload(context.Background(), rsc, &blockblob.UploadOptions{Tier: &tier})
//	_require.Nil(err)
//	//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//	credential, err := testcommon.GetGenericCredential(nil, testcommon.TestAccountDefault)
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
//	downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(nil))
//	_require.Nil(err)
//	_require.EqualValues(destData, content)
//
//	// Get properties to validate the tier
//	destBlobPropResp, err := destBlob.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//	_require.Equal(*destBlobPropResp.AccessTier, string(tier))
//}

func (s *BlockBlobRecordedTestsSuite) TestSetStandardBlobTierWithRehydratePriority() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	standardTier, rehydrateTier, rehydratePriority := blob.AccessTierArchive, blob.AccessTierCool, blob.RehydratePriorityStandard
	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)

	_, err = bbClient.SetTier(context.Background(), standardTier, &blob.SetTierOptions{
		RehydratePriority: &rehydratePriority,
	})
	_require.Nil(err)

	getResp1, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp1.AccessTier, string(standardTier))

	_, err = bbClient.SetTier(context.Background(), rehydrateTier, nil)
	_require.Nil(err)

	getResp2, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToCool))
}

func (s *BlockBlobRecordedTestsSuite) TestRehydrateStatus() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName1 := "rehydration_test_blob_1"
	blobName2 := "rehydration_test_blob_2"

	bbClient1 := testcommon.GetBlockBlobClient(blobName1, containerClient)
	reader1, _ := testcommon.GenerateData(1024)
	_, err = bbClient1.Upload(context.Background(), reader1, nil)
	_require.Nil(err)
	_, err = bbClient1.SetTier(context.Background(), blob.AccessTierArchive, nil)
	_require.Nil(err)
	_, err = bbClient1.SetTier(context.Background(), blob.AccessTierCool, nil)
	_require.Nil(err)

	getResp1, err := bbClient1.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp1.AccessTier, string(blob.AccessTierArchive))
	_require.Equal(*getResp1.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToCool))

	pager := containerClient.NewListBlobsFlatPager(nil)
	var blobs []*container.BlobItem
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
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

	bbClient2 := testcommon.GetBlockBlobClient(blobName2, containerClient)
	reader2, _ := testcommon.GenerateData(1024)
	_, err = bbClient2.Upload(context.Background(), reader2, nil)
	_require.Nil(err)
	_, err = bbClient2.SetTier(context.Background(), blob.AccessTierArchive, nil)
	_require.Nil(err)
	_, err = bbClient2.SetTier(context.Background(), blob.AccessTierHot, nil)
	_require.Nil(err)

	getResp2, err := bbClient2.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp2.AccessTier, string(blob.AccessTierArchive))
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToHot))
}

func (s *BlockBlobRecordedTestsSuite) TestCopyBlobWithRehydratePriority() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sourceBlobName := testcommon.GenerateBlobName(testName)
	sourceBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, sourceBlobName, containerClient)

	blobTier, rehydratePriority := blob.AccessTierArchive, blob.RehydratePriorityHigh

	copyBlobName := "copy" + sourceBlobName
	destBBClient := testcommon.GetBlockBlobClient(copyBlobName, containerClient)
	_, err = destBBClient.StartCopyFromURL(context.Background(), sourceBBClient.URL(), &blob.StartCopyFromURLOptions{
		RehydratePriority: &rehydratePriority,
		Tier:              &blobTier,
	})
	_require.Nil(err)

	getResp1, err := destBBClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp1.AccessTier, string(blobTier))

	_, err = destBBClient.SetTier(context.Background(), blob.AccessTierHot, nil)
	_require.Nil(err)

	getResp2, err := destBBClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToHot))
}

func (s *BlockBlobRecordedTestsSuite) TestBlobServiceClientDelete() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	code := 404
	testcommon.RunTestRequiringServiceProperties(context.Background(), _require, svcClient, string(rune(code)),
		testcommon.EnableSoftDelete, func(context.Context, *require.Assertions, *service.Client) error { return nil }, testcommon.DisableSoftDelete)
}

func setAndCheckBlockBlobTier(_require *require.Assertions, bbClient *blockblob.Client, tier blob.AccessTier) {
	_, err := bbClient.SetTier(context.Background(), tier, nil)
	_require.Nil(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.AccessTier, string(tier))
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSetTierAllTiersOnBlockBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierHot)
	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierCool)
	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierArchive)

}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobGetPropertiesUsingVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	blobProp, _ := bbClient.GetProperties(context.Background(), nil)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		Metadata: testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	uploadResp, err := bbClient.Upload(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), &uploadBlockBlobOptions)
	_require.Nil(err)
	_require.NotNil(uploadResp.VersionID)
	blobProp, _ = bbClient.GetProperties(context.Background(), nil)
	_require.EqualValues(uploadResp.VersionID, blobProp.VersionID)
	_require.EqualValues(uploadResp.LastModified, blobProp.LastModified)
	_require.Equal(*uploadResp.ETag, *blobProp.ETag)
	_require.Equal(*blobProp.IsCurrentVersion, true)
}

func (s *BlockBlobRecordedTestsSuite) TestGetSetBlobMetadataWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, &testcommon.TestCPKByValue, nil)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NotNil(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	resp, err := bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)
	_require.EqualValues(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)

	// Get blob properties without encryption key should fail the request.
	_, err = bbClient.GetProperties(context.Background(), nil)
	_require.NotNil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	getResp, err := bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.NotNil(getResp.Metadata)
	_require.Len(getResp.Metadata, len(testcommon.BasicMetadata))
	_require.EqualValues(getResp.Metadata, testcommon.BasicMetadata)

	_, err = bbClient.SetMetadata(context.Background(), map[string]string{}, &setBlobMetadataOptions)
	_require.Nil(err)

	getResp, err = bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.Nil(getResp.Metadata)
}

func (s *BlockBlobRecordedTestsSuite) TestGetSetBlobMetadataWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, nil, &testcommon.TestCPKByScope)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.NotNil(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	resp, err := bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)
	_require.EqualValues(resp.EncryptionScope, testcommon.TestCPKByScope.EncryptionScope)

	getResp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(getResp.Metadata)
	_require.Len(getResp.Metadata, len(testcommon.BasicMetadata))
	_require.EqualValues(getResp.Metadata, testcommon.BasicMetadata)

	_, err = bbClient.SetMetadata(context.Background(), map[string]string{}, &setBlobMetadataOptions)
	_require.Nil(err)

	getResp, err = bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(getResp.Metadata)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSnapshotWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, &testcommon.TestCPKByValue, nil)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(context.Background(), nil)
	_require.NotNil(err)

	createBlobSnapshotOptions := blob.CreateSnapshotOptions{
		CpkInfo: &testcommon.TestInvalidCPKByValue,
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	_require.NotNil(err)

	createBlobSnapshotOptions1 := blob.CreateSnapshotOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	resp, err := bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions1)
	_require.Nil(err)
	_require.Equal(*resp.IsServerEncrypted, false)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	dResp, err := snapshotURL.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(*dResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)

	_, err = snapshotURL.Delete(context.Background(), nil)
	_require.Nil(err)

	// Get blob properties of snapshot without encryption key should fail the request.
	_, err = snapshotURL.GetProperties(context.Background(), nil)
	_require.NotNil(err)

	//_assert(err.(StorageError).Response().StatusCode, chk.Equals, 404)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSnapshotWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, nil, &testcommon.TestCPKByScope)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(context.Background(), nil)
	_require.NotNil(err)

	createBlobSnapshotOptions := blob.CreateSnapshotOptions{
		CpkScopeInfo: &testcommon.TestInvalidCPKByScope,
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	_require.NotNil(err)

	createBlobSnapshotOptions1 := blob.CreateSnapshotOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	resp, err := bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions1)
	_require.Nil(err)
	_require.Equal(*resp.IsServerEncrypted, false)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	dResp, err := snapshotURL.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(*dResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)

	_, err = snapshotURL.Delete(context.Background(), nil)
	_require.Nil(err)
}

func (s *BlockBlobRecordedTestsSuite) TestCreateAndDownloadBlobSpecialCharactersWithVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	data := []rune("-._/()$=',~0123456789")
	for i := 0; i < len(data); i++ {
		blobName := "abc" + string(data[i])
		blobClient := containerClient.NewBlockBlobClient(blobName)
		resp, err := blobClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(string(data[i]))), nil)
		_require.Nil(err)
		_require.NotNil(resp.VersionID)

		blobClientWithVersionID, err := blobClient.WithVersionID(*resp.VersionID)
		_require.Nil(err)
		dResp, err := blobClientWithVersionID.DownloadStream(context.Background(), nil)
		_require.Nil(err)
		d1, err := io.ReadAll(dResp.Body)
		_require.Nil(err)
		_require.NotEqual(*dResp.Version, "")
		_require.EqualValues(string(d1), string(data[i]))
		_require.NotNil(dResp.VersionID)
		_require.Equal(*dResp.VersionID, *resp.VersionID)
	}
}

func (s *BlockBlobRecordedTestsSuite) TestDeleteSpecificBlobVersion() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		uploadResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"+strconv.Itoa(i)))), &blockblob.UploadOptions{
			Metadata: testcommon.BasicMetadata,
		})
		_require.Nil(err)
		_require.NotNil(uploadResp.VersionID)
		versions = append(versions, *uploadResp.VersionID)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemVersions},
	})

	found := make([]*container.BlobItem, 0)
	for listPager.More() {
		resp, err := listPager.NextPage(context.Background())
		_require.Nil(err)
		if err != nil {
			break
		}
		found = append(found, resp.Segment.BlobItems...)
	}
	_require.Len(found, 5)

	// Deleting the 2nd and 3rd versions
	for i := 0; i < 3; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.Nil(err)
		_, err = bbClientWithVersionID.Delete(context.Background(), nil)
		_require.Nil(err)
		//_require.Equal(deleteResp.RawResponse.StatusCode, 202)
	}

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemVersions},
	})

	found = make([]*container.BlobItem, 0)
	for listPager.More() {
		resp, err := listPager.NextPage(context.Background())
		_require.Nil(err)
		if err != nil {
			break
		}
		found = append(found, resp.Segment.BlobItems...)
	}
	_require.Len(found, 2)

	for i := 3; i < 5; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.Nil(err)
		downloadResp, err := bbClientWithVersionID.DownloadStream(context.Background(), nil)
		_require.Nil(err)
		destData, err := io.ReadAll(downloadResp.Body)
		_require.Nil(err)
		_require.EqualValues(destData, "data"+strconv.Itoa(i))
	}
}

func (s *BlockBlobRecordedTestsSuite) TestPutBlockListReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = testcommon.BlockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(d)), nil)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.NotNil(resp.Version)
		_require.NotEqual(*resp.Version, "")
	}

	commitResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
	_require.Nil(err)
	_require.NotNil(commitResp.VersionID)

	contentResp, err := bbClient.DownloadStream(context.Background(), nil)
	_require.Nil(err)
	contentData, err := io.ReadAll(contentResp.Body)
	_require.Nil(err)
	_require.EqualValues(contentData, []uint8(strings.Join(data, "")))
}

// nolint
func (s *BlockBlobUnrecordedTestsSuite) TestCreateBlockBlobReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	testSize := 2 * 1024 * 1024 // 1MB
	r, _ := testcommon.GetRandomDataAndReader(testSize)
	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	// Prepare source blob for copy.
	uploadResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(r), nil)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.NotNil(uploadResp.VersionID)

	csResp, err := bbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	//_require.Equal(csResp.RawResponse.StatusCode, 201)
	_require.NotNil(csResp.VersionID)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemSnapshots},
	})

	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = bbClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.Nil(err)
	//_require.Equal(deleteResp.RawResponse.StatusCode, 202)

	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemSnapshots, container.ListBlobsIncludeItemVersions},
	})

	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.NotEqual(len(found), 0)
}

func (s *BlockBlobUnrecordedTestsSuite) TestORSSource() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(nil, testcommon.TestAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)

	getResp, err := bbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(getResp.ObjectReplicationRules)
}

//func (s *azblobTestSuite) TestSnapshotSAS() {
//	//Generate URLs ----------------------------------------------------------------------------------------------------
//	bsu := getServiceClient(nil)
//	containerClient, containerName := getContainerClient(bsu)
//	blobURL, blobName := getBlockBlobClient(c, containerClient)
//
//	_, err := containerClient.Create(ctx, nil)
//	defer containerClient.Delete(ctx, nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	//Create file in container, download from snapshot to test. --------------------------------------------------------
//	blobClient := containerClient.NewBlockBlobClient(blobName)
//	data := "Hello world!"
//
//	contentType := "text/plain"
//	uploadBlockBlobOptions := BlockBlobUploadOptions{
//		HTTPHeaders: &HTTPHeaders{
//			BlobContentType: &contentType,
//		},
//	}
//	_, err = blobClient.Upload(ctx, strings.NewReader(data), &uploadBlockBlobOptions)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	//Create a snapshot & URL
//	createSnapshot, err := blobClient.CreateSnapshot(ctx, nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//	_assert(createSnapshot.Snapshot, chk.NotNil)
//
//	//Format snapshot time
//	snapTime, err := time.Parse(SnapshotTimeFormat, *createSnapshot.Snapshot)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	//Get credentials & current time
//	currentTime := time.Now().UTC()
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//
//	//Create SAS query
//	snapSASQueryParams, err := BlobSASSignatureValues{
//		StartTime:     currentTime,
//		ExpiryTime:    currentTime.Add(48 * time.Hour),
//		SnapshotTime:  snapTime,
//		Permissions:   "racwd",
//		ContainerName: containerName,
//		BlobName:      blobName,
//		Protocol:      SASProtocolHTTPS,
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//	time.Sleep(time.Second * 2)
//
//	//Attach SAS query to block blob URL
//	snapParts := NewBlobURLParts(blobURL.URL())
//	snapParts.SAS = snapSASQueryParams
//	sbUrl, err := NewBlockBlobClient(snapParts.URL(), azcore.AnonymousCredential(), nil)
//
//	//Test the snapshot
//	downloadResponse, err := sbUrl.Download(ctx, nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	downloadedData := &bytes.Buffer{}
//	reader := downloadResponse.Body(RetryReaderOptions{})
//	downloadedData.ReadFrom(reader)
//	reader.Close()
//
//	_assert(data, chk.Equals, downloadedData.String())
//
//	//Try to delete snapshot -------------------------------------------------------------------------------------------
//	_, err = sbUrl.Delete(ctx, nil)
//	if err != nil { //This shouldn't fail.
//		s.T().Fatal(err)
//	}
//
//	//Create a normal blob and attempt to use the snapshot SAS against it (assuming failure) ---------------------------
//	//If this succeeds, it means a normal SAS token was created.
//
//	uploadBlockBlobOptions1 := BlockBlobUploadOptions{
//		HTTPHeaders: &HTTPHeaders{
//			BlobContentType: &contentType,
//		},
//	}
//	fsbUrl := containerClient.NewBlockBlobClient("failsnap")
//	_, err = fsbUrl.Upload(ctx, strings.NewReader(data), &uploadBlockBlobOptions1)
//	if err != nil {
//		s.T().Fatal(err) //should succeed to create the blob via normal auth means
//	}
//
//	fsbUrlParts := NewBlobURLParts(fsbUrl.URL())
//	fsbUrlParts.SAS = snapSASQueryParams
//	fsbUrl, err = NewBlockBlobClient(fsbUrlParts.URL(), azcore.AnonymousCredential(), nil) //re-use fsbUrl as we don't need the sharedkey version anymore
//
//	resp, err := fsbUrl.Delete(ctx, nil)
//	if err == nil {
//		c.Fatal(resp) //This SHOULD fail. Otherwise we have a normal SAS token...
//	}
//}
