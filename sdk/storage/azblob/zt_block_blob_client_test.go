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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
)

func (s *azblobTestSuite) TestStageGetBlocks() {
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName := generateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		io.NopCloser(strings.NewReader("hello world"))
		putResp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], internal.NopCloser(strings.NewReader(d)), nil)
		require.NoError(s.T(), err)
		require.Equal(s.T(), putResp.RawResponse.StatusCode, 201)
		require.Nil(s.T(), putResp.ContentMD5)
		require.NotNil(s.T(), putResp.RequestID)
		require.NotNil(s.T(), putResp.Version)
		require.NotNil(s.T(), putResp.Date)
		require.Equal(s.T(), (*putResp.Date).IsZero(), false)
	}

	blockList, err := bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), blockList.RawResponse.StatusCode, 200)
	require.Nil(s.T(), blockList.LastModified)
	require.Nil(s.T(), blockList.ETag)
	require.NotNil(s.T(), blockList.ContentType)
	require.Nil(s.T(), blockList.BlobContentLength)
	require.NotNil(s.T(), blockList.RequestID)
	require.NotNil(s.T(), blockList.Version)
	require.NotNil(s.T(), blockList.Date)
	require.Equal(s.T(), (*blockList.Date).IsZero(), false)
	require.NotNil(s.T(), blockList.BlockList)
	require.Nil(s.T(), blockList.BlockList.CommittedBlocks)
	require.NotNil(s.T(), blockList.BlockList.UncommittedBlocks)
	require.Len(s.T(), blockList.BlockList.UncommittedBlocks, len(data))

	listResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), listResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), listResp.LastModified)
	require.Equal(s.T(), (*listResp.LastModified).IsZero(), false)
	require.NotNil(s.T(), listResp.ETag)
	require.NotNil(s.T(), listResp.RequestID)
	require.NotNil(s.T(), listResp.Version)
	require.NotNil(s.T(), listResp.Date)
	require.Equal(s.T(), (*listResp.Date).IsZero(), false)

	blockList, err = bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), blockList.RawResponse.StatusCode, 200)
	require.NotNil(s.T(), blockList.LastModified)
	require.Equal(s.T(), (*blockList.LastModified).IsZero(), false)
	require.NotNil(s.T(), blockList.ETag)
	require.NotNil(s.T(), blockList.ContentType)
	require.Equal(s.T(), *blockList.BlobContentLength, int64(25))
	require.NotNil(s.T(), blockList.RequestID)
	require.NotNil(s.T(), blockList.Version)
	require.NotNil(s.T(), blockList.Date)
	require.Equal(s.T(), (*blockList.Date).IsZero(), false)
	require.NotNil(s.T(), blockList.BlockList)
	require.NotNil(s.T(), blockList.BlockList.CommittedBlocks)
	require.Nil(s.T(), blockList.BlockList.UncommittedBlocks)
	require.Len(s.T(), blockList.BlockList.CommittedBlocks, len(data))
}

//nolint
func (s *azblobUnrecordedTestSuite) TestStageBlockFromURL() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	ctx := context.Background() // Use default Background context
	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))

	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	require.NoError(s.T(), err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(s.T(), err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockIDs := generateBlockIDsList(2)

	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockIDs[0], srcBlobURLWithSAS, 0, &StageBlockFromURLOptions{
		Offset: to.Int64Ptr(0),
		Count:  to.Int64Ptr(int64(contentSize / 2)),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), stageResp1.RawResponse.StatusCode, 201)
	require.NotEqual(s.T(), stageResp1.ContentMD5, "")
	require.NotEqual(s.T(), stageResp1.RequestID, "")
	require.NotEqual(s.T(), stageResp1.Version, "")
	require.Equal(s.T(), stageResp1.Date.IsZero(), false)

	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockIDs[1], srcBlobURLWithSAS, 0, &StageBlockFromURLOptions{
		Offset: to.Int64Ptr(int64(contentSize / 2)),
		Count:  to.Int64Ptr(int64(CountToEnd)),
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), stageResp2.RawResponse.StatusCode, 201)
	require.NotEqual(s.T(), stageResp2.ContentMD5, "")
	require.NotEqual(s.T(), stageResp2.RequestID, "")
	require.NotEqual(s.T(), stageResp2.Version, "")
	require.Equal(s.T(), stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), blockList.RawResponse.StatusCode, 200)
	require.NotNil(s.T(), blockList.BlockList)
	require.Nil(s.T(), blockList.BlockList.CommittedBlocks)
	require.NotNil(s.T(), blockList.BlockList.UncommittedBlocks)
	require.Len(s.T(), blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	listResp, err := destBlob.CommitBlockList(context.Background(), blockIDs, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), listResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), listResp.LastModified)
	require.Equal(s.T(), (*listResp.LastModified).IsZero(), false)
	require.NotNil(s.T(), listResp.ETag)
	require.NotNil(s.T(), listResp.RequestID)
	require.NotNil(s.T(), listResp.Version)
	require.NotNil(s.T(), listResp.Date)
	require.Equal(s.T(), (*listResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	require.NoError(s.T(), err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), destData, content)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestCopyBlockBlobFromURL() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	contentMD5 := md5.Sum(content)
	body := bytes.NewReader(content)
	ctx := context.Background()

	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, internal.NopCloser(body), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	require.NoError(s.T(), err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Invoke copy bbClient from URL.
	sourceContentMD5 := contentMD5[:]
	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &CopyBlockBlobFromURLOptions{
		Metadata:         map[string]string{"foo": "bar"},
		SourceContentMD5: sourceContentMD5,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 202)
	require.NotNil(s.T(), resp.ETag)
	require.NotNil(s.T(), resp.RequestID)
	require.NotNil(s.T(), resp.Version)
	require.NotNil(s.T(), resp.Date)
	require.Equal(s.T(), (*resp.Date).IsZero(), false)
	require.NotNil(s.T(), resp.CopyID)
	require.EqualValues(s.T(), resp.ContentMD5, sourceContentMD5)
	require.Equal(s.T(), *resp.CopyStatus, "success")

	// Make sure the metadata got copied over
	getPropResp, err := destBlob.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	metadata := getPropResp.Metadata
	require.NotNil(s.T(), metadata)
	require.Len(s.T(), metadata, 1)
	require.EqualValues(s.T(), metadata, map[string]string{"Foo": "bar"})

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	require.NoError(s.T(), err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), destData, content)

	// Edge case 1: Provide bad MD5 and make sure the copy fails
	_, badMD5 := getRandomDataAndReader(16)
	copyBlockBlobFromURLOptions1 := CopyBlockBlobFromURLOptions{
		SourceContentMD5: badMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
	require.Error(s.T(), err)

	// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
	copyBlockBlobFromURLOptions2 := CopyBlockBlobFromURLOptions{
		SourceContentMD5: sourceContentMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
	require.NoError(s.T(), err)
	require.Equal(s.T(), resp.RawResponse.StatusCode, 202)
	require.EqualValues(s.T(), *resp.CopyStatus, "success")
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobSASQueryParamOverrideResponseHeaders() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	//contentMD5 := md5.Sum(content)

	ctx := context.Background()

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadSrcResp, err := bbClient.Upload(ctx, internal.NopCloser(body), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uploadSrcResp.RawResponse.StatusCode, 201)

	// Get blob url with SAS.
	blobParts := NewBlobURLParts(bbClient.URL())

	cacheControlVal := "cache-control-override"
	contentDispositionVal := "content-disposition-override"
	contentEncodingVal := "content-encoding-override"
	contentLanguageVal := "content-language-override"
	contentTypeVal := "content-type-override"

	credential, err := getGenericCredential(nil, testAccountDefault)
	require.NoError(s.T(), err)
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
	require.NoError(s.T(), err)

	// Generate new bbClient client
	blobURLWithSAS := blobParts.URL()
	require.NotNil(s.T(), blobURLWithSAS)

	blobClientWithSAS, err := NewBlockBlobClientWithNoCredential(blobURLWithSAS, nil)
	require.NoError(s.T(), err)

	gResp, err := blobClientWithSAS.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *gResp.CacheControl, cacheControlVal)
	require.Equal(s.T(), *gResp.ContentDisposition, contentDispositionVal)
	require.Equal(s.T(), *gResp.ContentEncoding, contentEncodingVal)
	require.Equal(s.T(), *gResp.ContentLanguage, contentLanguageVal)
	require.Equal(s.T(), *gResp.ContentType, contentTypeVal)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestStageBlockWithMD5() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

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
	require.NoError(s.T(), err)
	require.Equal(s.T(), putResp.RawResponse.StatusCode, 201)
	require.EqualValues(s.T(), putResp.ContentMD5, contentMD5)
	require.NotNil(s.T(), putResp.RequestID)
	require.NotNil(s.T(), putResp.Version)
	require.NotNil(s.T(), putResp.Date)
	require.Equal(s.T(), (*putResp.Date).IsZero(), false)

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
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobPutBlobHTTPHeaders() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &UploadBlockBlobOptions{
		HTTPHeaders: &basicHeaders,
	})
	require.NoError(s.T(), err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	h := resp.GetHTTPHeaders()
	h.BlobContentMD5 = nil // the service generates a MD5 value, omit before comparing
	require.EqualValues(s.T(), h, basicHeaders)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataNotEmpty() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &UploadBlockBlobOptions{
		Metadata: basicMetadata,
	})
	require.NoError(s.T(), err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	actualMetadata := resp.Metadata
	require.NotNil(s.T(), actualMetadata)
	require.EqualValues(s.T(), actualMetadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataEmpty() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, nil)
	require.NoError(s.T(), err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Nil(s.T(), resp.Metadata)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataInvalid() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &UploadBlockBlobOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	})
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), invalidHeaderErrorSubstring)
}

func (s *azblobTestSuite) TestBlobPutBlobIfModifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), createResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), createResp.Date)

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
	require.NoError(s.T(), err)
	validateUpload(s.T(), bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfModifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), createResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), createResp.Date)

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
	require.Error(s.T(), err)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), createResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), createResp.Date)

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
	require.NoError(s.T(), err)

	validateUpload(s.T(), bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), createResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), createResp.Date)

	currentTime := getRelativeTimeFromAnchor(createResp.Date, -10)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(ctx, internal.NopCloser(bytes.NewReader(nil)), &uploadBlockBlobOptions)
	_ = err

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

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
	require.NoError(s.T(), err)

	validateUpload(s.T(), bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

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
	require.Error(s.T(), err)
	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfNoneMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

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
	require.NoError(s.T(), err)

	validateUpload(s.T(), bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfNoneMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)

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

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func validateBlobCommitted(t *testing.T, bbClient BlockBlobClient) {
	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Len(t, resp.BlockList.CommittedBlocks, 1)
}

func setupPutBlockListTest(t *testing.T, _context *testContext, testName string) (ContainerClient, BlockBlobClient, []string) {

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
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

func (s *azblobTestSuite) TestBlobPutBlockListHTTPHeadersEmpty() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobHTTPHeaders: &BlobHTTPHeaders{BlobContentDisposition: &blobContentDisposition},
	})
	require.NoError(s.T(), err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, nil)
	require.NoError(s.T(), err)

	resp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Nil(s.T(), resp.ContentDisposition)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfModifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(s.T(), err)
	require.NotNil(s.T(), commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	require.NoError(s.T(), err)

	validateBlobCommitted(s.T(), bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfModifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	getPropertyResp, err := containerClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), getPropertyResp.Date)

	currentTime := getRelativeTimeFromAnchor(getPropertyResp.Date, 10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_ = err

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfUnmodifiedSinceTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(s.T(), err)
	require.NotNil(s.T(), commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, 10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	require.NoError(s.T(), err)

	validateBlobCommitted(s.T(), bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfUnmodifiedSinceFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(s.T(), err)
	require.NotNil(s.T(), commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(s.T(), err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}},
	})
	require.NoError(s.T(), err)

	validateBlobCommitted(s.T(), bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(s.T(), err)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfNoneMatchTrue() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(s.T(), err)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	require.NoError(s.T(), err)

	validateBlobCommitted(s.T(), bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfNoneMatchFalse() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	require.NoError(s.T(), err)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(s.T(), err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListValidateData() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	require.NoError(s.T(), err)

	resp, err := bbClient.Download(ctx, nil)
	require.NoError(s.T(), err)
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	require.NoError(s.T(), err)
	require.Equal(s.T(), string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobPutBlockListModifyBlob() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _context, testName)
	defer deleteContainer(s.T(), containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	require.NoError(s.T(), err)

	_, err = bbClient.StageBlock(ctx, "0001", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(s.T(), err)
	_, err = bbClient.StageBlock(ctx, "0010", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(s.T(), err)
	_, err = bbClient.StageBlock(ctx, "0011", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(s.T(), err)
	_, err = bbClient.StageBlock(ctx, "0100", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	require.NoError(s.T(), err)

	_, err = bbClient.CommitBlockList(ctx, []string{"0001", "0011"}, nil)
	require.NoError(s.T(), err)

	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	require.NoError(s.T(), err)
	require.Len(s.T(), resp.BlockList.CommittedBlocks, 2)
	committed := resp.BlockList.CommittedBlocks
	require.Equal(s.T(), *(committed[0].Name), "0001")
	require.Equal(s.T(), *(committed[1].Name), "0011")
	require.Nil(s.T(), resp.BlockList.UncommittedBlocks)
}

func (s *azblobTestSuite) TestSetTierOnBlobUpload() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	for _, tier := range []AccessTier{AccessTierArchive, AccessTierCool, AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		uploadBlockBlobOptions := UploadBlockBlobOptions{
			HTTPHeaders: &basicHeaders,
			Tier:        &tier,
		}
		_, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &uploadBlockBlobOptions)
		require.NoError(s.T(), err)

		resp, err := bbClient.GetProperties(ctx, nil)
		require.NoError(s.T(), err)
		require.Equal(s.T(), *resp.AccessTier, string(tier))
	}
}

func (s *azblobTestSuite) TestBlobSetTierOnCommit() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := "test" + generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	for _, tier := range []AccessTier{AccessTierCool, AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
		_, err := bbClient.StageBlock(ctx, blockID, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
		require.NoError(s.T(), err)

		_, err = bbClient.CommitBlockList(ctx, []string{blockID}, &CommitBlockListOptions{
			Tier: &tier,
		})
		require.NoError(s.T(), err)

		resp, err := bbClient.GetBlockList(ctx, BlockListTypeCommitted, nil)
		require.NoError(s.T(), err)
		require.NotNil(s.T(), resp.BlockList)
		require.NotNil(s.T(), resp.BlockList.CommittedBlocks)
		require.Nil(s.T(), resp.BlockList.UncommittedBlocks)
		require.Len(s.T(), resp.BlockList.CommittedBlocks, 1)
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetTierOnCopyBlockBlobFromURL() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	const contentSize = 4 * 1024 * 1024 // 4 MB
	contentReader, _ := getRandomDataAndReader(contentSize)

	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient(generateBlobName(testName))

	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, internal.NopCloser(contentReader), &UploadBlockBlobOptions{Tier: &tier})
	require.NoError(s.T(), err)
	require.Equal(s.T(), uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	expiryTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	require.NoError(s.T(), err)

	credential, err := getGenericCredential(nil, testAccountDefault)
	if err != nil {
		s.T().Fatal("Couldn't fetch credential because " + err.Error())
	}
	sasQueryParams, err := AccountSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   AccountSASPermissions{Read: true, List: true}.String(),
		Services:      AccountSASServices{Blob: true}.String(),
		ResourceTypes: AccountSASResourceTypes{Container: true, Object: true}.String(),
	}.Sign(credential)
	require.NoError(s.T(), err)

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
		require.NoError(s.T(), err)
		require.Equal(s.T(), resp.RawResponse.StatusCode, 202)
		require.Equal(s.T(), *resp.CopyStatus, "success")

		destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
		require.NoError(s.T(), err)
		require.Equal(s.T(), *destBlobPropResp.AccessTier, string(tier))
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetTierOnStageBlockFromURL() {

	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))
	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))
	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, &UploadBlockBlobOptions{Tier: &tier})
	require.NoError(s.T(), err)
	require.Equal(s.T(), uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	credential, err := getGenericCredential(nil, testAccountDefault)
	require.NoError(s.T(), err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	offset1, count1 := int64(0), int64(4*1024)
	options1 := StageBlockFromURLOptions{
		Offset: &offset1,
		Count:  &count1,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), stageResp1.RawResponse.StatusCode, 201)
	require.Nil(s.T(), stageResp1.ContentMD5)
	require.NotEqual(s.T(), *stageResp1.RequestID, "")
	require.NotEqual(s.T(), *stageResp1.Version, "")
	require.Equal(s.T(), stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset: &offset2,
		Count:  &count2,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	require.NoError(s.T(), err)
	require.Equal(s.T(), stageResp2.RawResponse.StatusCode, 201)
	require.NotEqual(s.T(), stageResp2.ContentMD5, "")
	require.NotEqual(s.T(), stageResp2.RequestID, "")
	require.NotEqual(s.T(), stageResp2.Version, "")
	require.Equal(s.T(), stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), blockList.RawResponse.StatusCode, 200)
	require.NotNil(s.T(), blockList.BlockList)
	require.Nil(s.T(), blockList.BlockList.CommittedBlocks)
	require.NotNil(s.T(), blockList.BlockList.UncommittedBlocks)
	require.Len(s.T(), blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &CommitBlockListOptions{
		Tier: &tier,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), listResp.RawResponse.StatusCode, 201)
	require.NotNil(s.T(), listResp.LastModified)
	require.Equal(s.T(), (*listResp.LastModified).IsZero(), false)
	require.NotNil(s.T(), listResp.ETag)
	require.NotNil(s.T(), listResp.RequestID)
	require.NotNil(s.T(), listResp.Version)
	require.NotNil(s.T(), listResp.Date)
	require.Equal(s.T(), (*listResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	require.NoError(s.T(), err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(s.T(), err)
	require.EqualValues(s.T(), destData, content)

	// Get properties to validate the tier
	destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *destBlobPropResp.AccessTier, string(tier))
}

func (s *azblobTestSuite) TestSetStandardBlobTierWithRehydratePriority() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	standardTier, rehydrateTier, rehydratePriority := AccessTierArchive, AccessTierCool, RehydratePriorityStandard
	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlob(s.T(), bbName, containerClient)

	_, err = bbClient.SetTier(ctx, standardTier, &SetTierOptions{
		RehydratePriority: &rehydratePriority,
	})
	require.NoError(s.T(), err)

	getResp1, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *getResp1.AccessTier, string(standardTier))

	_, err = bbClient.SetTier(ctx, rehydrateTier, nil)
	require.NoError(s.T(), err)

	getResp2, err := bbClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToCool))
}

func (s *azblobTestSuite) TestRehydrateStatus() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	blobName1 := "rehydration_test_blob_1"
	blobName2 := "rehydration_test_blob_2"

	bbClient1 := getBlockBlobClient(blobName1, containerClient)
	reader1, _ := generateData(1024)
	_, err = bbClient1.Upload(ctx, reader1, nil)
	require.NoError(s.T(), err)
	_, err = bbClient1.SetTier(ctx, AccessTierArchive, nil)
	require.NoError(s.T(), err)
	_, err = bbClient1.SetTier(ctx, AccessTierCool, nil)
	require.NoError(s.T(), err)

	getResp1, err := bbClient1.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *getResp1.AccessTier, string(AccessTierArchive))
	require.Equal(s.T(), *getResp1.ArchiveStatus, string(ArchiveStatusRehydratePendingToCool))

	pager := containerClient.ListBlobsFlat(nil)
	var blobs []*BlobItemInternal
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		blobs = append(blobs, resp.ListBlobsFlatSegmentResponse.Segment.BlobItems...)
	}
	require.Nil(s.T(), pager.Err())
	require.GreaterOrEqual(s.T(), len(blobs), 1)
	require.Equal(s.T(), *blobs[0].Properties.AccessTier, AccessTierArchive)
	require.Equal(s.T(), *blobs[0].Properties.ArchiveStatus, ArchiveStatusRehydratePendingToCool)

	// ------------------------------------------

	bbClient2 := getBlockBlobClient(blobName2, containerClient)
	reader2, _ := generateData(1024)
	_, err = bbClient2.Upload(ctx, reader2, nil)
	require.NoError(s.T(), err)
	_, err = bbClient2.SetTier(ctx, AccessTierArchive, nil)
	require.NoError(s.T(), err)
	_, err = bbClient2.SetTier(ctx, AccessTierHot, nil)
	require.NoError(s.T(), err)

	getResp2, err := bbClient2.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *getResp2.AccessTier, string(AccessTierArchive))
	require.Equal(s.T(), *getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToHot))
}

func (s *azblobTestSuite) TestCopyBlobWithRehydratePriority() {

	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(s.T(), containerName, svcClient)
	defer deleteContainer(s.T(), containerClient)

	sourceBlobName := generateBlobName(testName)
	sourceBBClient := createNewBlockBlob(s.T(), sourceBlobName, containerClient)

	blobTier, rehydratePriority := AccessTierArchive, RehydratePriorityHigh

	copyBlobName := "copy" + sourceBlobName
	destBBClient := getBlockBlobClient(copyBlobName, containerClient)
	_, err = destBBClient.StartCopyFromURL(ctx, sourceBBClient.URL(), &StartCopyBlobOptions{
		RehydratePriority: &rehydratePriority,
		Tier:              &blobTier,
	})
	require.NoError(s.T(), err)

	getResp1, err := destBBClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *getResp1.AccessTier, string(blobTier))

	_, err = destBBClient.SetTier(ctx, AccessTierHot, nil)
	require.NoError(s.T(), err)

	getResp2, err := destBBClient.GetProperties(ctx, nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToHot))
}
