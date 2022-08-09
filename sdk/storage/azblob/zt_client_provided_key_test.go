//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
	"io"
	"strconv"
	"strings"
	"time"
)

/*
Azure Storage supports following operations support of sending customer-provided encryption keys on a request:
Put Blob, Put Block List, Put Block, Put Block from URL, Put Page, Put Page from URL, Append Block,
Set Blob Properties, Set Blob Metadata, Get Blob, Get Blob Properties, Get Blob Metadata, Snapshot Blob.
*/
var testEncryptedKey = "MDEyMzQ1NjcwMTIzNDU2NzAxMjM0NTY3MDEyMzQ1Njc="
var testEncryptedHash = "3QFFFpRA5+XANHqwwbT4yXDmrT/2JaLt/FKHjzhOdoE="
var testEncryptionAlgorithm = EncryptionAlgorithmTypeAES256
var testCPKByValue = CpkInfo{
	EncryptionKey:       &testEncryptedKey,
	EncryptionKeySHA256: &testEncryptedHash,
	EncryptionAlgorithm: &testEncryptionAlgorithm,
}

var testInvalidEncryptedKey = "mumbojumbo"
var testInvalidEncryptedHash = "mumbojumbohash"
var testInvalidCPKByValue = CpkInfo{
	EncryptionKey:       &testInvalidEncryptedKey,
	EncryptionKeySHA256: &testInvalidEncryptedHash,
	EncryptionAlgorithm: &testEncryptionAlgorithm,
}

var testEncryptedScope = "blobgokeytestscope"
var testCPKByScope = CpkScopeInfo{
	EncryptionScope: &testEncryptedScope,
}

var testInvalidEncryptedScope = "mumbojumbo"
var testInvalidCPKByScope = CpkScopeInfo{
	EncryptionScope: &testInvalidEncryptedScope,
}

func (s *azblobTestSuite) TestPutBlockAndPutBlockListWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	bbClient, _ := containerClient.NewBlockBlobClient(generateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := BlockBlobStageBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_require.Nil(err)
	}

	commitBlockListOptions := BlockBlobCommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_require.Nil(err)

	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(*resp.EncryptionKeySHA256, *(testCPKByValue.EncryptionKeySHA256))

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.Download(ctx, nil)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	b := bytes.Buffer{}
	reader := getResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue})
	_, _ = b.ReadFrom(reader)
	_ = reader.Close()
	_require.Equal(b.String(), "AAA BBB CCC ")
	_require.EqualValues(*getResp.ETag, *resp.ETag)
	_require.EqualValues(*getResp.LastModified, *resp.LastModified)
}

func (s *azblobTestSuite) TestPutBlockAndPutBlockListWithCPKByScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	bbClient, _ := containerClient.NewBlockBlobClient(generateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := BlockBlobStageBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_require.Nil(err)
	}

	commitBlockListOptions := BlockBlobCommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_require.Nil(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)

	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	_, err = bbClient.Download(ctx, &downloadBlobOptions)
	_require.NotNil(err)

	downloadBlobOptions = BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	b := bytes.Buffer{}
	reader := getResp.Body(nil)
	_, err = b.ReadFrom(reader)
	_require.Nil(err)
	_ = reader.Close() // The client must close the response body when finished with it
	_require.Equal(b.String(), "AAA BBB CCC ")
	_require.EqualValues(*getResp.ETag, *resp.ETag)
	_require.EqualValues(*getResp.LastModified, *resp.LastModified)
	_require.Equal(*getResp.IsServerEncrypted, true)
	_require.EqualValues(*getResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestPutBlockFromURLAndCommitWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob, _ := containerClient.NewBlockBlobClient("srcblob")
	destBlob, _ := containerClient.NewBlockBlobClient("destblob")
	uploadResp, err := srcBlob.Upload(ctx, rsc, nil)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
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
	options1 := BlockBlobStageBlockFromURLOptions{
		Offset:  &offset1,
		Count:   &count1,
		CpkInfo: &testCPKByValue,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	_require.Nil(err)
	_require.Equal(stageResp1.RawResponse.StatusCode, 201)
	_require.NotEqual(stageResp1.ContentMD5, "")
	_require.NotEqual(stageResp1.RequestID, "")
	_require.NotEqual(stageResp1.Version, "")
	_require.Equal(stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := BlockBlobStageBlockFromURLOptions{
		Offset:  &offset2,
		Count:   &count2,
		CpkInfo: &testCPKByValue,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	_require.Nil(err)
	_require.Equal(stageResp2.RawResponse.StatusCode, 201)
	_require.NotEqual(stageResp2.ContentMD5, "")
	_require.NotEqual(stageResp2.RequestID, "")
	_require.NotEqual(stageResp2.Version, "")
	_require.Equal(stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Equal(blockList.RawResponse.StatusCode, 200)
	_require.NotNil(blockList.BlockList)
	_require.Nil(blockList.BlockList.CommittedBlocks)
	_require.NotNil(blockList.BlockList.UncommittedBlocks)
	_require.Len(blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	commitBlockListOptions := BlockBlobCommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	_require.Nil(err)
	_require.Equal(listResp.RawResponse.StatusCode, 201)
	_require.NotNil(listResp.LastModified)
	_require.Equal((*listResp.LastModified).IsZero(), false)
	_require.NotNil(listResp.ETag)
	_require.NotNil(listResp.RequestID)
	_require.NotNil(listResp.Version)
	_require.NotNil(listResp.Date)
	_require.Equal((*listResp.Date).IsZero(), false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Equal(blockList.RawResponse.StatusCode, 200)
	_require.NotNil(blockList.BlockList)
	_require.Nil(blockList.BlockList.UncommittedBlocks)
	_require.NotNil(blockList.BlockList.CommittedBlocks)
	_require.Len(blockList.BlockList.CommittedBlocks, 2)

	// Check data integrity through downloading.
	_, err = destBlob.BlobClient.Download(ctx, nil)
	_require.NotNil(err)

	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	destData, err := io.ReadAll(downloadResp.Body(nil))
	_require.Nil(err)
	_require.EqualValues(destData, content)
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestPutBlockFromURLAndCommitWithCPKWithScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob, _ := containerClient.NewBlockBlobClient("srcblob")
	destBlob, _ := containerClient.NewBlockBlobClient("destblob")
	uploadResp, err := srcBlob.Upload(ctx, rsc, nil)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
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
	options1 := BlockBlobStageBlockFromURLOptions{
		Offset:       &offset1,
		Count:        &count1,
		CpkScopeInfo: &testCPKByScope,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	_require.Nil(err)
	_require.Equal(stageResp1.RawResponse.StatusCode, 201)
	_require.NotEqual(stageResp1.ContentMD5, "")
	_require.NotEqual(stageResp1.RequestID, "")
	_require.NotEqual(stageResp1.Version, "")
	_require.Equal(stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := BlockBlobStageBlockFromURLOptions{
		Offset:       &offset2,
		Count:        &count2,
		CpkScopeInfo: &testCPKByScope,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	_require.Nil(err)
	_require.Equal(stageResp2.RawResponse.StatusCode, 201)
	_require.NotEqual(stageResp2.ContentMD5, "")
	_require.NotEqual(stageResp2.RequestID, "")
	_require.NotEqual(stageResp2.Version, "")
	_require.Equal(stageResp2.Date.IsZero(), false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Equal(blockList.RawResponse.StatusCode, 200)
	_require.NotNil(blockList.BlockList)
	_require.Nil(blockList.BlockList.CommittedBlocks)
	_require.NotNil(blockList.BlockList.UncommittedBlocks)
	_require.Len(blockList.BlockList.UncommittedBlocks, 2)

	// Commit block list.
	commitBlockListOptions := BlockBlobCommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	_require.Nil(err)
	_require.Equal(listResp.RawResponse.StatusCode, 201)
	_require.NotNil(listResp.LastModified)
	_require.Equal((*listResp.LastModified).IsZero(), false)
	_require.NotNil(listResp.ETag)
	_require.NotNil(listResp.RequestID)
	_require.NotNil(listResp.Version)
	_require.NotNil(listResp.Date)
	_require.Equal((*listResp.Date).IsZero(), false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Equal(blockList.RawResponse.StatusCode, 200)
	_require.NotNil(blockList.BlockList)
	_require.Nil(blockList.BlockList.UncommittedBlocks)
	_require.NotNil(blockList.BlockList.CommittedBlocks)
	_require.Len(blockList.BlockList.CommittedBlocks, 2)

	downloadBlobOptions := BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	destData, err := io.ReadAll(downloadResp.Body(nil))
	_require.Nil(err)
	_require.EqualValues(destData, content)
	_require.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestUploadBlobWithMD5WithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 8 * 1024
	r, srcData := generateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient, _ := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadBlockBlobOptions := BlockBlobUploadOptions{
		CpkInfo: &testCPKByValue,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.Equal(*uploadResp.IsServerEncrypted, true)
	_require.EqualValues(uploadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.Download(ctx, nil)
	_require.NotNil(err)

	_, err = bbClient.Download(ctx, &BlobDownloadOptions{
		CpkInfo: &testInvalidCPKByValue,
	})
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadResp, err := bbClient.BlobClient.Download(ctx, &BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	})
	_require.Nil(err)
	_require.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(downloadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)
}

func (s *azblobTestSuite) TestUploadBlobWithMD5WithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 8 * 1024
	r, srcData := generateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient, _ := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadBlockBlobOptions := BlockBlobUploadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.Equal(*uploadResp.IsServerEncrypted, true)
	_require.EqualValues(uploadResp.EncryptionScope, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := bbClient.BlobClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

func (s *azblobTestSuite) TestAppendBlockWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	abClient, _ := containerClient.NewAppendBlobClient(generateBlobName(testName))

	createAppendBlobOptions := AppendBlobCreateOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlobAppendBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		resp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.Nil(err)
		_require.Equal(resp.RawResponse.StatusCode, 201)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal((*resp.Date).IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		_require.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)
	}

	// Get blob content without encryption key should fail the request.
	_, err = abClient.Download(ctx, nil)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := abClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)

	data, err := io.ReadAll(downloadResp.Body(nil))
	_require.Nil(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

func (s *azblobTestSuite) TestAppendBlockWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	abClient, _ := containerClient.NewAppendBlobClient(generateBlobName(testName))

	createAppendBlobOptions := AppendBlobCreateOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlobAppendBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		resp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.Nil(err)
		_require.Equal(resp.RawResponse.StatusCode, 201)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal((*resp.Date).IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		_require.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := abClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)

	data, err := io.ReadAll(downloadResp.Body(nil))
	_require.Nil(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	_require.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURLWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcABClient, _ := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob, _ := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcABClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(cResp1.RawResponse.StatusCode, 201)

	resp, err := srcABClient.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)
	_require.Equal(*resp.BlobAppendOffset, "0")
	_require.Equal(*resp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal((*resp.LastModified).IsZero(), false)
	_require.Nil(resp.ContentMD5)
	_require.NotNil(resp.RequestID)
	_require.NotNil(resp.Version)
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)

	srcBlobParts, _ := NewBlobURLParts(srcABClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	createAppendBlobOptions := AppendBlobCreateOptions{
		CpkInfo: &testCPKByValue,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	_require.Equal(cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
		Offset:  &offset,
		Count:   &count,
		CpkInfo: &testCPKByValue,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_require.Nil(err)
	_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendFromURLResp.ETag)
	_require.NotNil(appendFromURLResp.LastModified)
	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_require.NotNil(appendFromURLResp.ContentMD5)
	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
	_require.NotNil(appendFromURLResp.RequestID)
	_require.NotNil(appendFromURLResp.Version)
	_require.NotNil(appendFromURLResp.Date)
	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
	_require.Equal(*appendFromURLResp.IsServerEncrypted, true)

	// Get blob content without encryption key should fail the request.
	_, err = destBlob.Download(ctx, nil)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destBlob.Download(ctx, &downloadBlobOptions)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)

	_require.Equal(*downloadResp.IsServerEncrypted, true)
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURLWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcClient, _ := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob, _ := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcClient.Create(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(cResp1.RawResponse.StatusCode, 201)

	resp, err := srcClient.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)
	_require.Equal(*resp.BlobAppendOffset, "0")
	_require.Equal(*resp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal((*resp.LastModified).IsZero(), false)
	_require.Nil(resp.ContentMD5)
	_require.NotNil(resp.RequestID)
	_require.NotNil(resp.Version)
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)

	srcBlobParts, _ := NewBlobURLParts(srcClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	createAppendBlobOptions := AppendBlobCreateOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	_require.Nil(err)
	_require.Equal(cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
		Offset:       &offset,
		Count:        &count,
		CpkScopeInfo: &testCPKByScope,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_require.Nil(err)
	_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendFromURLResp.ETag)
	_require.NotNil(appendFromURLResp.LastModified)
	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_require.NotNil(appendFromURLResp.ContentMD5)
	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
	_require.NotNil(appendFromURLResp.RequestID)
	_require.NotNil(appendFromURLResp.Version)
	_require.NotNil(appendFromURLResp.Date)
	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
	_require.Equal(*appendFromURLResp.IsServerEncrypted, true)

	downloadBlobOptions := BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.Equal(*downloadResp.IsServerEncrypted, true)
	_require.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName), svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := generateData(contentSize)
	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(_require, pbName, containerClient, int64(contentSize), &testCPKByValue, nil)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		CpkInfo:   &testCPKByValue,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.EqualValues(uploadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, CountToEnd)})
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
		resp := pager.PageResponse()
		pageListResp := resp.PageList.PageRange
		start, end := int64(0), int64(contentSize-1)
		rawStart, rawEnd := pageListResp[0].Raw()
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
	}

	// Get blob content without encryption key should fail the request.
	_, err = pbClient.Download(ctx, nil)
	_require.NotNil(err)

	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = pbClient.Download(ctx, &downloadBlobOptions)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := pbClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)

	destData, err := io.ReadAll(downloadResp.Body(nil))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := generateData(contentSize)
	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(_require, pbName, containerClient, int64(contentSize), nil, &testCPKByScope)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange:    &HttpRange{offset, count},
		CpkScopeInfo: &testCPKByScope,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.EqualValues(uploadResp.EncryptionScope, testCPKByScope.EncryptionScope)

	pager := pbClient.GetPageRanges(&PageBlobGetPageRangesOptions{PageRange: NewHttpRange(0, CountToEnd)})
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())
		resp := pager.PageResponse()
		pageListResp := resp.PageList.PageRange
		start, end := int64(0), int64(contentSize-1)
		rawStart, rawEnd := pageListResp[0].Raw()
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := pbClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)

	destData, err := io.ReadAll(downloadResp.Body(nil))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockFromURLWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 8 * 1024 // 1MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	srcPBName := "src" + generateBlobName(testName)
	bbClient := createNewPageBlobWithSize(_require, srcPBName, containerClient, int64(contentSize))
	dstPBName := "dst" + generateBlobName(testName)
	destBlob := createNewPageBlobWithCPK(_require, dstPBName, containerClient, int64(contentSize), &testCPKByValue, nil)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := bbClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)
	srcBlobParts, _ := NewBlobURLParts(bbClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	uploadPagesFromURLOptions := PageBlobUploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.NotNil(resp.ContentMD5)
	_require.EqualValues(resp.ContentMD5, contentMD5)
	_require.NotNil(resp.RequestID)
	_require.NotNil(resp.Version)
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)
	_require.Equal(*resp.BlobSequenceNumber, int64(0))
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	_, err = destBlob.Download(ctx, nil)
	_require.NotNil(err)

	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destBlob.Download(ctx, &downloadBlobOptions)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockFromURLWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 8 * 1024 // 1MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	srcPBName := "src" + generateBlobName(testName)
	srcPBClient := createNewPageBlobWithSize(_require, srcPBName, containerClient, int64(contentSize))
	dstPBName := "dst" + generateBlobName(testName)
	dstPBBlob := createNewPageBlobWithCPK(_require, dstPBName, containerClient, int64(contentSize), nil, &testCPKByScope)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := srcPBClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)
	srcBlobParts, _ := NewBlobURLParts(srcPBClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	uploadPagesFromURLOptions := PageBlobUploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkScopeInfo:     &testCPKByScope,
	}
	resp, err := dstPBBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.NotNil(resp.ContentMD5)
	_require.EqualValues(resp.ContentMD5, contentMD5)
	_require.NotNil(resp.RequestID)
	_require.NotNil(resp.Version)
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)
	_require.Equal(*resp.BlobSequenceNumber, int64(0))
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := dstPBBlob.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURLWithMD5WithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	contentSize := 8 * 1024
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	srcPBName := "src" + generateBlobName(testName)
	srcBlob := createNewPageBlobWithSize(_require, srcPBName, containerClient, int64(contentSize))

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := PageBlobUploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_require.Nil(err)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)

	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_require.Nil(err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	dstPBName := "dst" + generateBlobName(testName)
	destPBClient := createNewPageBlobWithCPK(_require, dstPBName, containerClient, int64(contentSize), &testCPKByValue, nil)
	uploadPagesFromURLOptions := PageBlobUploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.Nil(err)
	_require.Equal(resp.RawResponse.StatusCode, 201)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.NotNil(resp.ContentMD5)
	_require.EqualValues(resp.ContentMD5, contentMD5)
	_require.NotNil(resp.RequestID)
	_require.NotNil(resp.Version)
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)
	_require.Equal(*resp.BlobSequenceNumber, int64(0))
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	_, err = destPBClient.Download(ctx, nil)
	_require.NotNil(err)

	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destPBClient.Download(ctx, &downloadBlobOptions)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destPBClient.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := io.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_require.Nil(err)
	_require.EqualValues(destData, srcData)

	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions1 := PageBlobUploadPagesFromURLOptions{
		SourceContentMD5: badContentMD5,
	}
	_, err = destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions1)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeMD5Mismatch)
}

//func (s *azblobTestSuite) TestClearDiffPagesWithCPK() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	pbName := generateBlobName(testName)
//	pbClient := createNewPageBlobWithCPK(_require, pbName, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)
//
//	contentSize := 2 * 1024
//	r := getReaderToGeneratedBytes(contentSize)
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := PageBlobUploadPagesOptions{Range: &HttpRange{offset, count}, CpkInfo: &testCPKByValue}
//	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
//	_require.Nil(err)
//
//	createBlobSnapshotOptions := BlobCreateSnapshotOptions{
//		CpkInfo: &testCPKByValue,
//	}
//	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
//	_require.Nil(err)
//
//	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
//	uploadPagesOptions1 := PageBlobUploadPagesOptions{Range: &HttpRange{offset1, count1}, CpkInfo: &testCPKByValue}
//	_, err = pbClient.UploadPages(context.Background(), getReaderToGeneratedBytes(2048), &uploadPagesOptions1)
//	_require.Nil(err)
//
//	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
//	_require.Nil(err)
//	pageRangeResp := pageListResp.PageList.Range
//	_require.NotNil(pageRangeResp)
//	_require.Len(pageRangeResp, 1)
//	rawStart, rawEnd := pageRangeResp[0].Raw()
//	_require.Equal(rawStart, offset1)
//	_require.Equal(rawEnd, end1)
//
//	clearPagesOptions := PageBlobClearPagesOptions{
//		CpkInfo: &testCPKByValue,
//	}
//	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, &clearPagesOptions)
//	_require.Nil(err)
//	_require.Equal(clearResp.RawResponse.StatusCode, 201)
//
//	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
//	_require.Nil(err)
//	_require.Nil(pageListResp.PageList.Range)
//}

func (s *azblobTestSuite) TestBlobResizeWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(_require, pbName, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)

	resizePageBlobOptions := PageBlobResizeOptions{
		CpkInfo: &testCPKByValue,
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, _ := pbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.Equal(*resp.ContentLength, int64(PageBlobPageBytes))
}

func (s *azblobTestSuite) TestGetSetBlobMetadataWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_require, bbName, containerClient, &testCPKByValue, nil)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.NotNil(err)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)
	_require.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	// Get blob properties without encryption key should fail the request.
	_, err = bbClient.GetProperties(ctx, nil)
	_require.NotNil(err)

	getBlobPropertiesOptions := BlobGetPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.NotNil(getResp.Metadata)
	_require.Len(getResp.Metadata, len(basicMetadata))
	_require.EqualValues(getResp.Metadata, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	_require.Nil(err)

	getResp, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_require.Nil(err)
	_require.Nil(getResp.Metadata)
}

func (s *azblobTestSuite) TestGetSetBlobMetadataWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_require, bbName, containerClient, nil, &testCPKByScope)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_require.NotNil(err)

	setBlobMetadataOptions := BlobSetMetadataOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_require.Nil(err)
	_require.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)

	getResp, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(getResp.Metadata)
	_require.Len(getResp.Metadata, len(basicMetadata))
	_require.EqualValues(getResp.Metadata, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	_require.Nil(err)

	getResp, err = bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(getResp.Metadata)
}

func (s *azblobTestSuite) TestBlobSnapshotWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_require, bbName, containerClient, &testCPKByValue, nil)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(ctx, nil)
	_require.NotNil(err)

	createBlobSnapshotOptions := BlobCreateSnapshotOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_require.NotNil(err)

	createBlobSnapshotOptions1 := BlobCreateSnapshotOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	_require.Nil(err)
	_require.Equal(*resp.IsServerEncrypted, false)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := BlobDownloadOptions{
		CpkInfo: &testCPKByValue,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(*dResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	_, err = snapshotURL.Delete(ctx, nil)
	_require.Nil(err)

	// Get blob properties of snapshot without encryption key should fail the request.
	_, err = snapshotURL.GetProperties(ctx, nil)
	_require.NotNil(err)

	//_assert(err.(StorageError).Response().StatusCode, chk.Equals, 404)
}

func (s *azblobTestSuite) TestBlobSnapshotWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_require, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_require, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_require, bbName, containerClient, nil, &testCPKByScope)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(ctx, nil)
	_require.NotNil(err)

	createBlobSnapshotOptions := BlobCreateSnapshotOptions{
		CpkScopeInfo: &testInvalidCPKByScope,
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_require.NotNil(err)

	createBlobSnapshotOptions1 := BlobCreateSnapshotOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	_require.Nil(err)
	_require.Equal(*resp.IsServerEncrypted, false)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := BlobDownloadOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(*dResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	_, err = snapshotURL.Delete(ctx, nil)
	_require.Nil(err)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestUploadStreamToBlobBlobPropertiesWithCPKKey() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	_require.NoError(err)

	blobSize := 1024
	bufferSize := 8 * 1024
	maxBuffers := 3

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	// Set up test blob
	blobName := generateBlobName(testName)
	bbClient, err := getBlockBlobClient(blobName, containerClient)
	_require.Nil(err)

	// Create some data to test the upload stream
	blobContentReader, blobData := generateData(blobSize)

	// Perform UploadStream
	uploadResp, err := bbClient.UploadStream(ctx, blobContentReader,
		UploadStreamOptions{
			BufferSize:  bufferSize,
			MaxBuffers:  maxBuffers,
			Metadata:    basicMetadata,
			BlobTagsMap: basicBlobTagsMap,
			HTTPHeaders: &basicHeaders,
			CpkInfo:     &testCPKByValue,
		})

	// Assert that upload was successful
	_require.Equal(err, nil)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)

	getPropertiesResp, err := bbClient.GetProperties(ctx, &BlobGetPropertiesOptions{CpkInfo: &testCPKByValue})
	_require.NoError(err)
	_require.EqualValues(getPropertiesResp.Metadata, basicMetadata)
	_require.Equal(*getPropertiesResp.TagCount, int64(len(basicBlobTagsMap)))
	_require.Equal(getPropertiesResp.GetHTTPHeaders(), basicHeaders)

	getTagsResp, err := bbClient.GetTags(ctx, nil)
	_require.NoError(err)
	_require.Len(getTagsResp.BlobTagSet, 3)
	for _, blobTag := range getTagsResp.BlobTagSet {
		_require.Equal(basicBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	// Download the blob to verify
	downloadResponse, err := bbClient.Download(ctx, &BlobDownloadOptions{CpkInfo: &testCPKByValue})
	_require.NoError(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResponse.RawResponse.Body)
	_require.NoError(err)
	_require.Equal(len(actualBlobData), blobSize)
	_require.EqualValues(actualBlobData, blobData)
}

// nolint
func (s *azblobUnrecordedTestSuite) TestUploadStreamToBlobBlobPropertiesWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	_require.NoError(err)

	blobSize := 1024
	bufferSize := 8 * 1024
	maxBuffers := 3

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	// Set up test blob
	blobName := generateBlobName(testName)
	bbClient, err := getBlockBlobClient(blobName, containerClient)
	_require.NoError(err)

	// Create some data to test the upload stream
	blobContentReader, blobData := generateData(blobSize)

	// Perform UploadStream
	uploadResp, err := bbClient.UploadStream(ctx, blobContentReader,
		UploadStreamOptions{
			BufferSize:   bufferSize,
			MaxBuffers:   maxBuffers,
			Metadata:     basicMetadata,
			BlobTagsMap:  basicBlobTagsMap,
			HTTPHeaders:  &basicHeaders,
			CpkScopeInfo: &testCPKByScope,
		})

	// Assert that upload was successful
	_require.Equal(err, nil)
	_require.Equal(uploadResp.RawResponse.StatusCode, 201)

	getPropertiesResp, err := bbClient.GetProperties(ctx, nil)
	_require.NoError(err)
	_require.EqualValues(getPropertiesResp.Metadata, basicMetadata)
	_require.Equal(*getPropertiesResp.TagCount, int64(len(basicBlobTagsMap)))
	_require.Equal(getPropertiesResp.GetHTTPHeaders(), basicHeaders)

	getTagsResp, err := bbClient.GetTags(ctx, nil)
	_require.NoError(err)
	_require.Len(getTagsResp.BlobTagSet, 3)
	for _, blobTag := range getTagsResp.BlobTagSet {
		_require.Equal(basicBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	// Download the blob to verify
	downloadResponse, err := bbClient.Download(ctx, &BlobDownloadOptions{CpkScopeInfo: &testCPKByScope})
	_require.NoError(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResponse.RawResponse.Body)
	_require.NoError(err)
	_require.Equal(len(actualBlobData), blobSize)
	_require.EqualValues(actualBlobData, blobData)
}
