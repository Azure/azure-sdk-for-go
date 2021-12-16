// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

func (s *azblobTestSuite) TestStageGetBlocks() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		io.NopCloser(strings.NewReader("hello world"))
		putResp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], internal.NopCloser(strings.NewReader(d)), nil)
		_assert.Nil(err)
		_assert.Equal(putResp.RawResponse.StatusCode, 201)
		_assert.Nil(putResp.ContentMD5)
		_assert.NotNil(putResp.RequestID)
		_assert.NotNil(putResp.Version)
		_assert.NotNil(putResp.Date)
		_assert.Equal((*putResp.Date).IsZero(), false)
	}

	blockList, err := bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Equal(blockList.RawResponse.StatusCode, 200)
	_assert.Nil(blockList.LastModified)
	_assert.Nil(blockList.ETag)
	_assert.NotNil(blockList.ContentType)
	_assert.Nil(blockList.BlobContentLength)
	_assert.NotNil(blockList.RequestID)
	_assert.NotNil(blockList.Version)
	_assert.NotNil(blockList.Date)
	_assert.Equal((*blockList.Date).IsZero(), false)
	_assert.NotNil(blockList.BlockList)
	_assert.Nil(blockList.BlockList.CommittedBlocks)
	_assert.NotNil(blockList.BlockList.UncommittedBlocks)
	_assert.Len(blockList.BlockList.UncommittedBlocks, len(data))

	listResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
	_assert.Nil(err)
	_assert.Equal(listResp.RawResponse.StatusCode, 201)
	_assert.NotNil(listResp.LastModified)
	_assert.Equal((*listResp.LastModified).IsZero(), false)
	_assert.NotNil(listResp.ETag)
	_assert.NotNil(listResp.RequestID)
	_assert.NotNil(listResp.Version)
	_assert.NotNil(listResp.Date)
	_assert.Equal((*listResp.Date).IsZero(), false)

	blockList, err = bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Equal(blockList.RawResponse.StatusCode, 200)
	_assert.NotNil(blockList.LastModified)
	_assert.Equal((*blockList.LastModified).IsZero(), false)
	_assert.NotNil(blockList.ETag)
	_assert.NotNil(blockList.ContentType)
	_assert.Equal(*blockList.BlobContentLength, int64(25))
	_assert.NotNil(blockList.RequestID)
	_assert.NotNil(blockList.Version)
	_assert.NotNil(blockList.Date)
	_assert.Equal((*blockList.Date).IsZero(), false)
	_assert.NotNil(blockList.BlockList)
	_assert.NotNil(blockList.BlockList.CommittedBlocks)
	_assert.Nil(blockList.BlockList.UncommittedBlocks)
	_assert.Len(blockList.BlockList.CommittedBlocks, len(data))
}

//nolint
func (s *azblobUnrecordedTestSuite) TestStageBlockFromURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	ctx := context.Background() // Use default Background context
	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))

	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, nil)
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	_assert.Nil(err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockIDs := generateBlockIDsList(2)

	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockIDs[0], srcBlobURLWithSAS, 0, &StageBlockFromURLOptions{
		Offset: to.Int64Ptr(0),
		Count:  to.Int64Ptr(int64(contentSize / 2)),
	})
	_assert.Nil(err)
	_assert.Equal(stageResp1.RawResponse.StatusCode, 201)
	_assert.NotEqual(stageResp1.ContentMD5, "")
	_assert.NotEqual(stageResp1.RequestID, "")
	_assert.NotEqual(stageResp1.Version, "")
	_assert.Equal(stageResp1.Date.IsZero(), false)

	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockIDs[1], srcBlobURLWithSAS, 0, &StageBlockFromURLOptions{
		Offset: to.Int64Ptr(int64(contentSize / 2)),
		Count:  to.Int64Ptr(int64(CountToEnd)),
	})
	_assert.Nil(err)
	_assert.Equal(stageResp2.RawResponse.StatusCode, 201)
	_assert.NotEqual(stageResp2.ContentMD5, "")
	_assert.NotEqual(stageResp2.RequestID, "")
	_assert.NotEqual(stageResp2.Version, "")
	_assert.Equal(stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Equal(blockList.RawResponse.StatusCode, 200)
	_assert.NotNil(blockList.BlockList)
	_assert.Nil(blockList.BlockList.CommittedBlocks)
	_assert.NotNil(blockList.BlockList.UncommittedBlocks)
	_assert.Len(blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	listResp, err := destBlob.CommitBlockList(context.Background(), blockIDs, nil)
	_assert.Nil(err)
	_assert.Equal(listResp.RawResponse.StatusCode, 201)
	_assert.NotNil(listResp.LastModified)
	_assert.Equal((*listResp.LastModified).IsZero(), false)
	_assert.NotNil(listResp.ETag)
	_assert.NotNil(listResp.RequestID)
	_assert.NotNil(listResp.Version)
	_assert.NotNil(listResp.Date)
	_assert.Equal((*listResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, content)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestCopyBlockBlobFromURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	contentMD5 := md5.Sum(content)
	body := bytes.NewReader(content)
	ctx := context.Background()

	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, internal.NopCloser(body), nil)
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)

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
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 202)
	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.RequestID)
	_assert.NotNil(resp.Version)
	_assert.NotNil(resp.Date)
	_assert.Equal((*resp.Date).IsZero(), false)
	_assert.NotNil(resp.CopyID)
	_assert.EqualValues(resp.ContentMD5, sourceContentMD5)
	_assert.Equal(*resp.CopyStatus, "success")

	// Make sure the metadata got copied over
	getPropResp, err := destBlob.GetProperties(ctx, nil)
	_assert.Nil(err)
	metadata := getPropResp.Metadata
	_assert.NotNil(metadata)
	_assert.Len(metadata, 1)
	_assert.EqualValues(metadata, map[string]string{"Foo": "bar"})

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, content)

	// Edge case 1: Provide bad MD5 and make sure the copy fails
	_, badMD5 := getRandomDataAndReader(16)
	copyBlockBlobFromURLOptions1 := CopyBlockBlobFromURLOptions{
		SourceContentMD5: badMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
	_assert.NotNil(err)

	// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
	copyBlockBlobFromURLOptions2 := CopyBlockBlobFromURLOptions{
		SourceContentMD5: sourceContentMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 202)
	_assert.EqualValues(*resp.CopyStatus, "success")
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBlobSASQueryParamOverrideResponseHeaders() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	//contentMD5 := md5.Sum(content)

	ctx := context.Background()

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadSrcResp, err := bbClient.Upload(ctx, internal.NopCloser(body), nil)
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get blob url with SAS.
	blobParts := NewBlobURLParts(bbClient.URL())

	cacheControlVal := "cache-control-override"
	contentDispositionVal := "content-disposition-override"
	contentEncodingVal := "content-encoding-override"
	contentLanguageVal := "content-language-override"
	contentTypeVal := "content-type-override"

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
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
	_assert.Nil(err)

	// Generate new bbClient client
	blobURLWithSAS := blobParts.URL()
	_assert.NotNil(blobURLWithSAS)

	blobClientWithSAS, err := NewBlockBlobClientWithNoCredential(blobURLWithSAS, nil)
	_assert.Nil(err)

	gResp, err := blobClientWithSAS.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*gResp.CacheControl, cacheControlVal)
	_assert.Equal(*gResp.ContentDisposition, contentDispositionVal)
	_assert.Equal(*gResp.ContentEncoding, contentEncodingVal)
	_assert.Equal(*gResp.ContentLanguage, contentLanguageVal)
	_assert.Equal(*gResp.ContentType, contentTypeVal)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestStageBlockWithMD5() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

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
	_assert.Nil(err)
	_assert.Equal(putResp.RawResponse.StatusCode, 201)
	_assert.EqualValues(putResp.ContentMD5, contentMD5)
	_assert.NotNil(putResp.RequestID)
	_assert.NotNil(putResp.Version)
	_assert.NotNil(putResp.Date)
	_assert.Equal((*putResp.Date).IsZero(), false)

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
	_assert.NotNil(err)
	_assert.Contains(err.Error(), StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestBlobPutBlobHTTPHeaders() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &UploadBlockBlobOptions{
		HTTPHeaders: &basicHeaders,
	})
	_assert.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	h := resp.GetHTTPHeaders()
	h.BlobContentMD5 = nil // the service generates a MD5 value, omit before comparing
	_assert.EqualValues(h, basicHeaders)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataNotEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(ctx, internal.NopCloser(body), &UploadBlockBlobOptions{
		Metadata: basicMetadata,
	})
	_assert.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	actualMetadata := resp.Metadata
	_assert.NotNil(actualMetadata)
	_assert.EqualValues(actualMetadata, basicMetadata)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, nil)
	_assert.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestBlobPutBlobMetadataInvalid() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)

	_, err = bbClient.Upload(ctx, rsc, &UploadBlockBlobOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	})
	_assert.NotNil(err)
	_assert.Contains(err.Error(), invalidHeaderErrorSubstring)
}

func (s *azblobTestSuite) TestBlobPutBlobIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(createResp.RawResponse.StatusCode, 201)
	_assert.NotNil(createResp.Date)

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
	_assert.Nil(err)
	validateUpload(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(createResp.RawResponse.StatusCode, 201)
	_assert.NotNil(createResp.Date)

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
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(createResp.RawResponse.StatusCode, 201)
	_assert.NotNil(createResp.Date)

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
	_assert.Nil(err)

	validateUpload(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	_assert.Equal(createResp.RawResponse.StatusCode, 201)
	_assert.NotNil(createResp.Date)

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

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

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
	_assert.Nil(err)

	validateUpload(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

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
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlobIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

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
	_assert.Nil(err)

	validateUpload(_assert, bbClient.BlobClient)
}

func (s *azblobTestSuite) TestBlobPutBlobIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blockBlobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)

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

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func validateBlobCommitted(_assert *assert.Assertions, bbClient BlockBlobClient) {
	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Len(resp.BlockList.CommittedBlocks, 1)
}

func setupPutBlockListTest(_assert *assert.Assertions, _context *testContext,
	testName string) (ContainerClient, BlockBlobClient, []string) {

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)

	blobName := generateBlobName(testName)
	bbClient := getBlockBlobClient(blobName, containerClient)

	blockIDs := generateBlockIDsList(1)
	_, err = bbClient.StageBlock(ctx, blockIDs[0], internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	_assert.Nil(err)
	return containerClient, bbClient, blockIDs
}

func (s *azblobTestSuite) TestBlobPutBlockListHTTPHeadersEmpty() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobHTTPHeaders: &BlobHTTPHeaders{BlobContentDisposition: &blobContentDisposition},
	})
	_assert.Nil(err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, nil)
	_assert.Nil(err)

	resp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.ContentDisposition)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_assert.Nil(err)
	_assert.NotNil(commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_assert.Nil(err)

	validateBlobCommitted(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	getPropertyResp, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(getPropertyResp.Date)

	currentTime := getRelativeTimeFromAnchor(getPropertyResp.Date, 10)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_ = err

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfUnmodifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_assert.Nil(err)
	_assert.NotNil(commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, 10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	_assert.Nil(err)

	validateBlobCommitted(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfUnmodifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_assert.Nil(err)
	_assert.NotNil(commitBlockListResp.Date)

	currentTime := getRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_assert.Nil(err)

	_, err = bbClient.CommitBlockList(ctx, blockIDs, &CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}},
	})
	_assert.Nil(err)

	validateBlobCommitted(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_assert.Nil(err)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfNoneMatchTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_assert.Nil(err)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)
	_assert.Nil(err)

	validateBlobCommitted(_assert, bbClient)
}

func (s *azblobTestSuite) TestBlobPutBlockListIfNoneMatchFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	resp, err := bbClient.CommitBlockList(ctx, blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_assert.Nil(err)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: &BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(ctx, blockIDs, &commitBlockListOptions)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestBlobPutBlockListValidateData() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	_assert.Nil(err)

	resp, err := bbClient.Download(ctx, nil)
	_assert.Nil(err)
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(string(data), blockBlobDefaultData)
}

func (s *azblobTestSuite) TestBlobPutBlockListModifyBlob() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	containerClient, bbClient, blockIDs := setupPutBlockListTest(_assert, _context, testName)
	defer deleteContainer(_assert, containerClient)

	_, err := bbClient.CommitBlockList(ctx, blockIDs, nil)
	_assert.Nil(err)

	_, err = bbClient.StageBlock(ctx, "0001", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_assert.Nil(err)
	_, err = bbClient.StageBlock(ctx, "0010", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_assert.Nil(err)
	_, err = bbClient.StageBlock(ctx, "0011", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_assert.Nil(err)
	_, err = bbClient.StageBlock(ctx, "0100", internal.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_assert.Nil(err)

	_, err = bbClient.CommitBlockList(ctx, []string{"0001", "0011"}, nil)
	_assert.Nil(err)

	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Len(resp.BlockList.CommittedBlocks, 2)
	committed := resp.BlockList.CommittedBlocks
	_assert.Equal(*(committed[0].Name), "0001")
	_assert.Equal(*(committed[1].Name), "0011")
	_assert.Nil(resp.BlockList.UncommittedBlocks)
}

func (s *azblobTestSuite) TestSetTierOnBlobUpload() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	for _, tier := range []AccessTier{AccessTierArchive, AccessTierCool, AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		uploadBlockBlobOptions := UploadBlockBlobOptions{
			HTTPHeaders: &basicHeaders,
			Tier:        &tier,
		}
		_, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &uploadBlockBlobOptions)
		_assert.Nil(err)

		resp, err := bbClient.GetProperties(ctx, nil)
		_assert.Nil(err)
		_assert.Equal(*resp.AccessTier, string(tier))
	}
}

func (s *azblobTestSuite) TestBlobSetTierOnCommit() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := "test" + generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	for _, tier := range []AccessTier{AccessTierCool, AccessTierHot} {
		blobName := strings.ToLower(string(tier)) + generateBlobName(testName)
		bbClient := getBlockBlobClient(blobName, containerClient)

		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
		_, err := bbClient.StageBlock(ctx, blockID, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
		_assert.Nil(err)

		_, err = bbClient.CommitBlockList(ctx, []string{blockID}, &CommitBlockListOptions{
			Tier: &tier,
		})
		_assert.Nil(err)

		resp, err := bbClient.GetBlockList(ctx, BlockListTypeCommitted, nil)
		_assert.Nil(err)
		_assert.NotNil(resp.BlockList)
		_assert.NotNil(resp.BlockList.CommittedBlocks)
		_assert.Nil(resp.BlockList.UncommittedBlocks)
		_assert.Len(resp.BlockList.CommittedBlocks, 1)
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetTierOnCopyBlockBlobFromURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	const contentSize = 4 * 1024 * 1024 // 4 MB
	contentReader, _ := getRandomDataAndReader(contentSize)

	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient(generateBlobName(testName))

	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, internal.NopCloser(contentReader), &UploadBlockBlobOptions{Tier: &tier})
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	expiryTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_assert.Nil(err)

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
	_assert.Nil(err)

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
		_assert.Nil(err)
		_assert.Equal(resp.RawResponse.StatusCode, 202)
		_assert.Equal(*resp.CopyStatus, "success")

		destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
		_assert.Nil(err)
		_assert.Equal(*destBlobPropResp.AccessTier, string(tier))
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetTierOnStageBlockFromURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName(testName))
	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName(testName))
	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, &UploadBlockBlobOptions{Tier: &tier})
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
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
	_assert.Nil(err)
	_assert.Equal(stageResp1.RawResponse.StatusCode, 201)
	_assert.Nil(stageResp1.ContentMD5)
	_assert.NotEqual(*stageResp1.RequestID, "")
	_assert.NotEqual(*stageResp1.Version, "")
	_assert.Equal(stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset: &offset2,
		Count:  &count2,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	_assert.Nil(err)
	_assert.Equal(stageResp2.RawResponse.StatusCode, 201)
	_assert.NotEqual(stageResp2.ContentMD5, "")
	_assert.NotEqual(stageResp2.RequestID, "")
	_assert.NotEqual(stageResp2.Version, "")
	_assert.Equal(stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Equal(blockList.RawResponse.StatusCode, 200)
	_assert.NotNil(blockList.BlockList)
	_assert.Nil(blockList.BlockList.CommittedBlocks)
	_assert.NotNil(blockList.BlockList.UncommittedBlocks)
	_assert.Len(blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &CommitBlockListOptions{
		Tier: &tier,
	})
	_assert.Nil(err)
	_assert.Equal(listResp.RawResponse.StatusCode, 201)
	_assert.NotNil(listResp.LastModified)
	_assert.Equal((*listResp.LastModified).IsZero(), false)
	_assert.NotNil(listResp.ETag)
	_assert.NotNil(listResp.RequestID)
	_assert.NotNil(listResp.Version)
	_assert.NotNil(listResp.Date)
	_assert.Equal((*listResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, content)

	// Get properties to validate the tier
	destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*destBlobPropResp.AccessTier, string(tier))
}

func (s *azblobTestSuite) TestSetStandardBlobTierWithRehydratePriority() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	standardTier, rehydrateTier, rehydratePriority := AccessTierArchive, AccessTierCool, RehydratePriorityStandard
	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, bbName, containerClient)

	_, err = bbClient.SetTier(ctx, standardTier, &SetTierOptions{
		RehydratePriority: &rehydratePriority,
	})
	_assert.Nil(err)

	getResp1, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp1.AccessTier, string(standardTier))

	_, err = bbClient.SetTier(ctx, rehydrateTier, nil)
	_assert.Nil(err)

	getResp2, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToCool))
}

func (s *azblobTestSuite) TestRehydrateStatus() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName1 := "rehydration_test_blob_1"
	blobName2 := "rehydration_test_blob_2"

	bbClient1 := getBlockBlobClient(blobName1, containerClient)
	reader1, _ := generateData(1024)
	_, err = bbClient1.Upload(ctx, reader1, nil)
	_assert.Nil(err)
	_, err = bbClient1.SetTier(ctx, AccessTierArchive, nil)
	_assert.Nil(err)
	_, err = bbClient1.SetTier(ctx, AccessTierCool, nil)
	_assert.Nil(err)

	getResp1, err := bbClient1.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp1.AccessTier, string(AccessTierArchive))
	_assert.Equal(*getResp1.ArchiveStatus, string(ArchiveStatusRehydratePendingToCool))

	pager := containerClient.ListBlobsFlat(nil)
	var blobs []*BlobItemInternal
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		blobs = append(blobs, resp.ListBlobsFlatSegmentResponse.Segment.BlobItems...)
	}
	_assert.Nil(pager.Err())
	_assert.GreaterOrEqual(len(blobs), 1)
	_assert.Equal(*blobs[0].Properties.AccessTier, AccessTierArchive)
	_assert.Equal(*blobs[0].Properties.ArchiveStatus, ArchiveStatusRehydratePendingToCool)

	// ------------------------------------------

	bbClient2 := getBlockBlobClient(blobName2, containerClient)
	reader2, _ := generateData(1024)
	_, err = bbClient2.Upload(ctx, reader2, nil)
	_assert.Nil(err)
	_, err = bbClient2.SetTier(ctx, AccessTierArchive, nil)
	_assert.Nil(err)
	_, err = bbClient2.SetTier(ctx, AccessTierHot, nil)
	_assert.Nil(err)

	getResp2, err := bbClient2.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp2.AccessTier, string(AccessTierArchive))
	_assert.Equal(*getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToHot))
}

func (s *azblobTestSuite) TestCopyBlobWithRehydratePriority() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	sourceBlobName := generateBlobName(testName)
	sourceBBClient := createNewBlockBlob(_assert, sourceBlobName, containerClient)

	blobTier, rehydratePriority := AccessTierArchive, RehydratePriorityHigh

	copyBlobName := "copy" + sourceBlobName
	destBBClient := getBlockBlobClient(copyBlobName, containerClient)
	_, err = destBBClient.StartCopyFromURL(ctx, sourceBBClient.URL(), &StartCopyBlobOptions{
		RehydratePriority: &rehydratePriority,
		Tier:              &blobTier,
	})
	_assert.Nil(err)

	getResp1, err := destBBClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp1.AccessTier, string(blobTier))

	_, err = destBBClient.SetTier(ctx, AccessTierHot, nil)
	_assert.Nil(err)

	getResp2, err := destBBClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getResp2.ArchiveStatus, string(ArchiveStatusRehydratePendingToHot))
}
