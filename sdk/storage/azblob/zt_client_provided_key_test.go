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
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
var testEncryptionAlgorithm = "AES256"
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
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := StageBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_assert.Nil(err)
	}

	commitBlockListOptions := CommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_assert.Nil(err)

	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.LastModified)
	_assert.Equal(*resp.IsServerEncrypted, true)
	_assert.EqualValues(*resp.EncryptionKeySHA256, *(testCPKByValue.EncryptionKeySHA256))

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.Download(ctx, nil)
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	b := bytes.Buffer{}
	reader := getResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue})
	_, _ = b.ReadFrom(reader)
	_ = reader.Close()
	_assert.Equal(b.String(), "AAA BBB CCC ")
	_assert.EqualValues(*getResp.ETag, *resp.ETag)
	_assert.EqualValues(*getResp.LastModified, *resp.LastModified)
}

func (s *azblobTestSuite) TestPutBlockAndPutBlockListWithCPKByScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := StageBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_assert.Nil(err)
	}

	commitBlockListOptions := CommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_assert.Nil(err)
	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.LastModified)
	_assert.Equal(*resp.IsServerEncrypted, true)
	_assert.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	_, err = bbClient.Download(ctx, &downloadBlobOptions)
	_assert.NotNil(err)

	downloadBlobOptions = DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	b := bytes.Buffer{}
	reader := getResp.Body(RetryReaderOptions{})
	_, err = b.ReadFrom(reader)
	_assert.Nil(err)
	_ = reader.Close() // The client must close the response body when finished with it
	_assert.Equal(b.String(), "AAA BBB CCC ")
	_assert.EqualValues(*getResp.ETag, *resp.ETag)
	_assert.EqualValues(*getResp.LastModified, *resp.LastModified)
	_assert.Equal(*getResp.IsServerEncrypted, true)
	_assert.EqualValues(*getResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPutBlockFromURLAndCommitWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")
	uploadResp, err := srcBlob.Upload(ctx, rsc, nil)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)

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
		Offset:  &offset1,
		Count:   &count1,
		CpkInfo: &testCPKByValue,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	_assert.Nil(err)
	_assert.Equal(stageResp1.RawResponse.StatusCode, 201)
	_assert.NotEqual(stageResp1.ContentMD5, "")
	_assert.NotEqual(stageResp1.RequestID, "")
	_assert.NotEqual(stageResp1.Version, "")
	_assert.Equal(stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset:  &offset2,
		Count:   &count2,
		CpkInfo: &testCPKByValue,
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
	commitBlockListOptions := CommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	_assert.Nil(err)
	_assert.Equal(listResp.RawResponse.StatusCode, 201)
	_assert.NotNil(listResp.LastModified)
	_assert.Equal((*listResp.LastModified).IsZero(), false)
	_assert.NotNil(listResp.ETag)
	_assert.NotNil(listResp.RequestID)
	_assert.NotNil(listResp.Version)
	_assert.NotNil(listResp.Date)
	_assert.Equal((*listResp.Date).IsZero(), false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Equal(blockList.RawResponse.StatusCode, 200)
	_assert.NotNil(blockList.BlockList)
	_assert.Nil(blockList.BlockList.UncommittedBlocks)
	_assert.NotNil(blockList.BlockList.CommittedBlocks)
	_assert.Len(blockList.BlockList.CommittedBlocks, 2)

	// Check data integrity through downloading.
	_, err = destBlob.BlobClient.Download(ctx, nil)
	_assert.NotNil(err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, content)
	_assert.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPutBlockFromURLAndCommitWithCPKWithScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")
	uploadResp, err := srcBlob.Upload(ctx, rsc, nil)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)

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
		Offset:       &offset1,
		Count:        &count1,
		CpkScopeInfo: &testCPKByScope,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	_assert.Nil(err)
	_assert.Equal(stageResp1.RawResponse.StatusCode, 201)
	_assert.NotEqual(stageResp1.ContentMD5, "")
	_assert.NotEqual(stageResp1.RequestID, "")
	_assert.NotEqual(stageResp1.Version, "")
	_assert.Equal(stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset:       &offset2,
		Count:        &count2,
		CpkScopeInfo: &testCPKByScope,
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
	commitBlockListOptions := CommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	_assert.Nil(err)
	_assert.Equal(listResp.RawResponse.StatusCode, 201)
	_assert.NotNil(listResp.LastModified)
	_assert.Equal((*listResp.LastModified).IsZero(), false)
	_assert.NotNil(listResp.ETag)
	_assert.NotNil(listResp.RequestID)
	_assert.NotNil(listResp.Version)
	_assert.NotNil(listResp.Date)
	_assert.Equal((*listResp.Date).IsZero(), false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	_assert.Nil(err)
	_assert.Equal(blockList.RawResponse.StatusCode, 200)
	_assert.NotNil(blockList.BlockList)
	_assert.Nil(blockList.BlockList.UncommittedBlocks)
	_assert.NotNil(blockList.BlockList.CommittedBlocks)
	_assert.Len(blockList.BlockList.CommittedBlocks, 2)

	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, content)
	_assert.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadBlobWithMD5WithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024
	r, srcData := generateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)
	_assert.Equal(*uploadResp.IsServerEncrypted, true)
	_assert.EqualValues(uploadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.Download(ctx, nil)
	_assert.NotNil(err)

	_, err = bbClient.Download(ctx, &DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	})
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadResp, err := bbClient.BlobClient.Download(ctx, &DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	})
	_assert.Nil(err)
	_assert.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
	_assert.EqualValues(downloadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)
}

func (s *azblobTestSuite) TestUploadBlobWithMD5WithCPKScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024
	r, srcData := generateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)
	_assert.Equal(*uploadResp.IsServerEncrypted, true)
	_assert.EqualValues(uploadResp.EncryptionScope, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := bbClient.BlobClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	_assert.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
	_assert.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

func (s *azblobTestSuite) TestAppendBlockWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	abClient := containerClient.NewAppendBlobClient(generateBlobName(testName))

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		resp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_assert.Nil(err)
		_assert.Equal(resp.RawResponse.StatusCode, 201)
		_assert.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_assert.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_assert.NotNil(resp.ETag)
		_assert.NotNil(resp.LastModified)
		_assert.Equal(resp.LastModified.IsZero(), false)
		_assert.NotEqual(resp.ContentMD5, "")

		_assert.NotEqual(resp.Version, "")
		_assert.NotNil(resp.Date)
		_assert.Equal((*resp.Date).IsZero(), false)
		_assert.Equal(*resp.IsServerEncrypted, true)
		_assert.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)
	}

	// Get blob content without encryption key should fail the request.
	_, err = abClient.Download(ctx, nil)
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := abClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(string(data), "AAA BBB CCC ")
	_assert.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

func (s *azblobTestSuite) TestAppendBlockWithCPKScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	abClient := containerClient.NewAppendBlobClient(generateBlobName(testName))

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		resp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_assert.Nil(err)
		_assert.Equal(resp.RawResponse.StatusCode, 201)
		_assert.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_assert.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_assert.NotNil(resp.ETag)
		_assert.NotNil(resp.LastModified)
		_assert.Equal(resp.LastModified.IsZero(), false)
		_assert.NotEqual(resp.ContentMD5, "")

		_assert.NotEqual(resp.Version, "")
		_assert.NotNil(resp.Date)
		_assert.Equal((*resp.Date).IsZero(), false)
		_assert.Equal(*resp.IsServerEncrypted, true)
		_assert.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := abClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)

	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(string(data), "AAA BBB CCC ")
	_assert.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURLWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcABClient := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcABClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp1.RawResponse.StatusCode, 201)

	resp, err := srcABClient.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)
	_assert.Equal(*resp.BlobAppendOffset, "0")
	_assert.Equal(*resp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.LastModified)
	_assert.Equal((*resp.LastModified).IsZero(), false)
	_assert.Nil(resp.ContentMD5)
	_assert.NotNil(resp.RequestID)
	_assert.NotNil(resp.Version)
	_assert.NotNil(resp.Date)
	_assert.Equal((*resp.Date).IsZero(), false)

	srcBlobParts := NewBlobURLParts(srcABClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
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

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	_assert.Nil(err)
	_assert.Equal(cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:  &offset,
		Count:   &count,
		CpkInfo: &testCPKByValue,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_assert.Nil(err)
	_assert.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_assert.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendFromURLResp.ETag)
	_assert.NotNil(appendFromURLResp.LastModified)
	_assert.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_assert.NotNil(appendFromURLResp.ContentMD5)
	_assert.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
	_assert.NotNil(appendFromURLResp.RequestID)
	_assert.NotNil(appendFromURLResp.Version)
	_assert.NotNil(appendFromURLResp.Date)
	_assert.Equal((*appendFromURLResp.Date).IsZero(), false)
	_assert.Equal(*appendFromURLResp.IsServerEncrypted, true)

	// Get blob content without encryption key should fail the request.
	_, err = destBlob.Download(ctx, nil)
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destBlob.Download(ctx, &downloadBlobOptions)
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)

	_assert.Equal(*downloadResp.IsServerEncrypted, true)
	_assert.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestAppendBlockFromURLWithCPKScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcClient := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcClient.Create(context.Background(), nil)
	_assert.Nil(err)
	_assert.Equal(cResp1.RawResponse.StatusCode, 201)

	resp, err := srcClient.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)
	_assert.Equal(*resp.BlobAppendOffset, "0")
	_assert.Equal(*resp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.LastModified)
	_assert.Equal((*resp.LastModified).IsZero(), false)
	_assert.Nil(resp.ContentMD5)
	_assert.NotNil(resp.RequestID)
	_assert.NotNil(resp.Version)
	_assert.NotNil(resp.Date)
	_assert.Equal((*resp.Date).IsZero(), false)

	srcBlobParts := NewBlobURLParts(srcClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
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

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	_assert.Nil(err)
	_assert.Equal(cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:       &offset,
		Count:        &count,
		CpkScopeInfo: &testCPKByScope,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	_assert.Nil(err)
	_assert.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
	_assert.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_assert.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_assert.NotNil(appendFromURLResp.ETag)
	_assert.NotNil(appendFromURLResp.LastModified)
	_assert.Equal((*appendFromURLResp.LastModified).IsZero(), false)
	_assert.NotNil(appendFromURLResp.ContentMD5)
	_assert.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
	_assert.NotNil(appendFromURLResp.RequestID)
	_assert.NotNil(appendFromURLResp.Version)
	_assert.NotNil(appendFromURLResp.Date)
	_assert.Equal((*appendFromURLResp.Date).IsZero(), false)
	_assert.Equal(*appendFromURLResp.IsServerEncrypted, true)

	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	_assert.Equal(*downloadResp.IsServerEncrypted, true)
	_assert.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := generateData(contentSize)
	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(_assert, pbName, containerClient, int64(contentSize), &testCPKByValue, nil)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		CpkInfo:   &testCPKByValue,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)
	_assert.EqualValues(uploadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	_assert.Nil(err)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), int64(contentSize-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	_assert.Equal(rawStart, start)
	_assert.Equal(rawEnd, end)

	// Get blob content without encryption key should fail the request.
	_, err = pbClient.Download(ctx, nil)
	_assert.NotNil(err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = pbClient.Download(ctx, &downloadBlobOptions)
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := pbClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
	_assert.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockWithCPKScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := generateData(contentSize)
	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(_assert, pbName, containerClient, int64(contentSize), nil, &testCPKByScope)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange:    &HttpRange{offset, count},
		CpkScopeInfo: &testCPKByScope,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)
	_assert.EqualValues(uploadResp.EncryptionScope, testCPKByScope.EncryptionScope)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	_assert.Nil(err)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), int64(contentSize-1)
	// _assert((*pageListResp)[0], chk.DeepEquals, PageRange{Start: &start, End: &end})
	rawStart, rawEnd := pageListResp[0].Raw()
	_assert.Equal(rawStart, start)
	_assert.Equal(rawEnd, end)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := pbClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
	_assert.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockFromURLWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024 // 1MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	srcPBName := "src" + generateBlobName(testName)
	bbClient := createNewPageBlobWithSize(_assert, srcPBName, containerClient, int64(contentSize))
	dstPBName := "dst" + generateBlobName(testName)
	destBlob := createNewPageBlobWithCPK(_assert, dstPBName, containerClient, int64(contentSize), &testCPKByValue, nil)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := bbClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)
	srcBlobParts := NewBlobURLParts(bbClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
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
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)
	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.LastModified)
	_assert.NotNil(resp.ContentMD5)
	_assert.EqualValues(resp.ContentMD5, contentMD5)
	_assert.NotNil(resp.RequestID)
	_assert.NotNil(resp.Version)
	_assert.NotNil(resp.Date)
	_assert.Equal((*resp.Date).IsZero(), false)
	_assert.Equal(*resp.BlobSequenceNumber, int64(0))
	_assert.Equal(*resp.IsServerEncrypted, true)
	_assert.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	_, err = destBlob.Download(ctx, nil)
	_assert.NotNil(err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destBlob.Download(ctx, &downloadBlobOptions)
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	_assert.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestPageBlockFromURLWithCPKScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024 // 1MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	srcPBName := "src" + generateBlobName(testName)
	srcPBClient := createNewPageBlobWithSize(_assert, srcPBName, containerClient, int64(contentSize))
	dstPBName := "dst" + generateBlobName(testName)
	dstPBBlob := createNewPageBlobWithCPK(_assert, dstPBName, containerClient, int64(contentSize), nil, &testCPKByScope)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := srcPBClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)
	srcBlobParts := NewBlobURLParts(srcPBClient.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
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
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkScopeInfo:     &testCPKByScope,
	}
	resp, err := dstPBBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)
	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.LastModified)
	_assert.NotNil(resp.ContentMD5)
	_assert.EqualValues(resp.ContentMD5, contentMD5)
	_assert.NotNil(resp.RequestID)
	_assert.NotNil(resp.Version)
	_assert.NotNil(resp.Date)
	_assert.Equal((*resp.Date).IsZero(), false)
	_assert.Equal(*resp.BlobSequenceNumber, int64(0))
	_assert.Equal(*resp.IsServerEncrypted, true)
	_assert.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := dstPBBlob.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	_assert.EqualValues(*downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadPagesFromURLWithMD5WithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 8 * 1024
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	srcPBName := "src" + generateBlobName(testName)
	srcBlob := createNewPageBlobWithSize(_assert, srcPBName, containerClient, int64(contentSize))

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)

	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	_assert.Nil(err)
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
	destPBClient := createNewPageBlobWithCPK(_assert, dstPBName, containerClient, int64(contentSize), &testCPKByValue, nil)
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 201)
	_assert.NotNil(resp.ETag)
	_assert.NotNil(resp.LastModified)
	_assert.NotNil(resp.ContentMD5)
	_assert.EqualValues(resp.ContentMD5, contentMD5)
	_assert.NotNil(resp.RequestID)
	_assert.NotNil(resp.Version)
	_assert.NotNil(resp.Date)
	_assert.Equal((*resp.Date).IsZero(), false)
	_assert.Equal(*resp.BlobSequenceNumber, int64(0))
	_assert.Equal(*resp.IsServerEncrypted, true)
	_assert.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	_, err = destPBClient.Download(ctx, nil)
	_assert.NotNil(err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destPBClient.Download(ctx, &downloadBlobOptions)
	_assert.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destPBClient.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	_assert.EqualValues(*downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	_assert.Nil(err)
	_assert.EqualValues(destData, srcData)

	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions1 := UploadPagesFromURLOptions{
		SourceContentMD5: badContentMD5,
	}
	_, err = destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions1)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

func (s *azblobTestSuite) TestClearDiffPagesWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(_assert, pbName, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)

	contentSize := 2 * 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}, CpkInfo: &testCPKByValue}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_assert.Nil(err)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkInfo: &testCPKByValue,
	}
	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	_assert.Nil(err)

	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}, CpkInfo: &testCPKByValue}
	_, err = pbClient.UploadPages(context.Background(), getReaderToGeneratedBytes(2048), &uploadPagesOptions1)
	_assert.Nil(err)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
	_assert.Nil(err)
	pageRangeResp := pageListResp.PageList.PageRange
	_assert.NotNil(pageRangeResp)
	_assert.Len(pageRangeResp, 1)
	rawStart, rawEnd := pageRangeResp[0].Raw()
	_assert.Equal(rawStart, offset1)
	_assert.Equal(rawEnd, end1)

	clearPagesOptions := ClearPagesOptions{
		CpkInfo: &testCPKByValue,
	}
	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, &clearPagesOptions)
	_assert.Nil(err)
	_assert.Equal(clearResp.RawResponse.StatusCode, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
	_assert.Nil(err)
	_assert.Nil(pageListResp.PageList.PageRange)
}

func (s *azblobTestSuite) TestBlobResizeWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(_assert, pbName, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	_assert.Nil(err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, _ := pbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.Equal(*resp.ContentLength, int64(PageBlobPageBytes))
}

func (s *azblobTestSuite) TestGetSetBlobMetadataWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_assert, bbName, containerClient, &testCPKByValue, nil)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.NotNil(err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_assert.Nil(err)
	_assert.EqualValues(resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	// Get blob properties without encryption key should fail the request.
	_, err = bbClient.GetProperties(ctx, nil)
	_assert.NotNil(err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.Nil(err)
	_assert.NotNil(getResp.Metadata)
	_assert.Len(getResp.Metadata, len(basicMetadata))
	_assert.EqualValues(getResp.Metadata, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	_assert.Nil(err)

	getResp, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	_assert.Nil(err)
	_assert.Nil(getResp.Metadata)
}

func (s *azblobTestSuite) TestGetSetBlobMetadataWithCPKScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_assert, bbName, containerClient, nil, &testCPKByScope)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	_assert.NotNil(err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	_assert.Nil(err)
	_assert.EqualValues(resp.EncryptionScope, testCPKByScope.EncryptionScope)

	getResp, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(getResp.Metadata)
	_assert.Len(getResp.Metadata, len(basicMetadata))
	_assert.EqualValues(getResp.Metadata, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	_assert.Nil(err)

	getResp, err = bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(getResp.Metadata)
}

func (s *azblobTestSuite) TestBlobSnapshotWithCPK() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_assert, bbName, containerClient, &testCPKByValue, nil)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(ctx, nil)
	_assert.NotNil(err)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_assert.NotNil(err)

	createBlobSnapshotOptions1 := CreateBlobSnapshotOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	_assert.Nil(err)
	_assert.Equal(*resp.IsServerEncrypted, false)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	_assert.EqualValues(*dResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	_, err = snapshotURL.Delete(ctx, nil)
	_assert.Nil(err)

	// Get blob properties of snapshot without encryption key should fail the request.
	_, err = snapshotURL.GetProperties(ctx, nil)
	_assert.NotNil(err)

	//_assert(err.(StorageError).Response().StatusCode, chk.Equals, 404)
}

func (s *azblobTestSuite) TestBlobSnapshotWithCPKScope() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(_assert, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(_assert, bbName, containerClient, nil, &testCPKByScope)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(ctx, nil)
	_assert.NotNil(err)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkScopeInfo: &testInvalidCPKByScope,
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	_assert.NotNil(err)

	createBlobSnapshotOptions1 := CreateBlobSnapshotOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	_assert.Nil(err)
	_assert.Equal(*resp.IsServerEncrypted, false)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	_assert.Nil(err)
	_assert.EqualValues(*dResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	_, err = snapshotURL.Delete(ctx, nil)
	_assert.Nil(err)
}
