// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
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

func TestPutBlockAndPutBlockListWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := StageBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		require.NoError(t, err)
	}

	commitBlockListOptions := CommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	require.NoError(t, err)

	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.LastModified)
	require.Equal(t, *resp.IsServerEncrypted, true)
	require.EqualValues(t, *resp.EncryptionKeySHA256, *(testCPKByValue.EncryptionKeySHA256))

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.Download(ctx, nil)
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	b := bytes.Buffer{}
	reader := getResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue})
	_, err = b.ReadFrom(reader)
	require.NoError(t, err)
	err = reader.Close()
	require.NoError(t, err)
	require.Equal(t, b.String(), "AAA BBB CCC ")
	require.EqualValues(t, *getResp.ETag, *resp.ETag)
	require.EqualValues(t, *getResp.LastModified, *resp.LastModified)
}

func TestPutBlockAndPutBlockListWithCPKByScope(t *testing.T) {
	t.Skip("encryption scope not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := StageBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		require.NoError(t, err)
	}

	commitBlockListOptions := CommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	require.NoError(t, err)
	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.LastModified)
	require.Equal(t, *resp.IsServerEncrypted, true)
	require.EqualValues(t, resp.EncryptionScope, testCPKByScope.EncryptionScope)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	_, err = bbClient.Download(ctx, &downloadBlobOptions)
	require.Error(t, err)

	downloadBlobOptions = DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	b := bytes.Buffer{}
	reader := getResp.Body(nil)
	_, err = b.ReadFrom(reader)
	require.NoError(t, err)
	_ = reader.Close() // The client must close the response body when finished with it
	require.Equal(t, b.String(), "AAA BBB CCC ")
	require.EqualValues(t, *getResp.ETag, *resp.ETag)
	require.EqualValues(t, *getResp.LastModified, *resp.LastModified)
	require.Equal(t, *getResp.IsServerEncrypted, true)
	require.EqualValues(t, *getResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

func TestPutBlockFromURLAndCommitWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")
	uploadResp, err := srcBlob.Upload(ctx, rsc, nil)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)

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
		Offset:  &offset1,
		Count:   &count1,
		CpkInfo: &testCPKByValue,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	require.NoError(t, err)
	require.Equal(t, stageResp1.RawResponse.StatusCode, 201)
	require.NotEqual(t, stageResp1.ContentMD5, "")
	require.NotEqual(t, stageResp1.RequestID, "")
	require.NotEqual(t, stageResp1.Version, "")
	require.Equal(t, stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset:  &offset2,
		Count:   &count2,
		CpkInfo: &testCPKByValue,
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
	commitBlockListOptions := CommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	require.NoError(t, err)
	require.Equal(t, listResp.RawResponse.StatusCode, 201)
	require.NotNil(t, listResp.LastModified)
	require.Equal(t, (*listResp.LastModified).IsZero(), false)
	require.NotNil(t, listResp.ETag)
	require.NotNil(t, listResp.RequestID)
	require.NotNil(t, listResp.Version)
	require.NotNil(t, listResp.Date)
	require.Equal(t, (*listResp.Date).IsZero(), false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Equal(t, blockList.RawResponse.StatusCode, 200)
	require.NotNil(t, blockList.BlockList)
	require.Nil(t, blockList.BlockList.UncommittedBlocks)
	require.NotNil(t, blockList.BlockList.CommittedBlocks)
	require.Len(t, blockList.BlockList.CommittedBlocks, 2)

	// Check data integrity through downloading.
	_, err = destBlob.BlobClient.Download(ctx, nil)
	require.Error(t, err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, content)
	require.EqualValues(t, *downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

func TestPutBlockFromURLAndCommitWithCPKWithScope(t *testing.T) {
	t.Skip("encryption scope not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := internal.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")
	uploadResp, err := srcBlob.Upload(ctx, rsc, nil)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)

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
		Offset:       &offset1,
		Count:        &count1,
		CpkScopeInfo: &testCPKByScope,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	require.NoError(t, err)
	require.Equal(t, stageResp1.RawResponse.StatusCode, 201)
	require.NotEqual(t, stageResp1.ContentMD5, "")
	require.NotEqual(t, stageResp1.RequestID, "")
	require.NotEqual(t, stageResp1.Version, "")
	require.Equal(t, stageResp1.Date.IsZero(), false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset:       &offset2,
		Count:        &count2,
		CpkScopeInfo: &testCPKByScope,
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
	commitBlockListOptions := CommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	require.NoError(t, err)
	require.Equal(t, listResp.RawResponse.StatusCode, 201)
	require.NotNil(t, listResp.LastModified)
	require.Equal(t, (*listResp.LastModified).IsZero(), false)
	require.NotNil(t, listResp.ETag)
	require.NotNil(t, listResp.RequestID)
	require.NotNil(t, listResp.Version)
	require.NotNil(t, listResp.Date)
	require.Equal(t, (*listResp.Date).IsZero(), false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	require.NoError(t, err)
	require.Equal(t, blockList.RawResponse.StatusCode, 200)
	require.NotNil(t, blockList.BlockList)
	require.Nil(t, blockList.BlockList.UncommittedBlocks)
	require.NotNil(t, blockList.BlockList.CommittedBlocks)
	require.Len(t, blockList.BlockList.CommittedBlocks, 2)

	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, content)
	require.EqualValues(t, *downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

func TestUploadBlobWithMD5WithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024
	r, srcData := generateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)
	require.Equal(t, *uploadResp.IsServerEncrypted, true)
	require.EqualValues(t, uploadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.Download(ctx, nil)
	require.Error(t, err)

	_, err = bbClient.Download(ctx, &DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	})
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadResp, err := bbClient.BlobClient.Download(ctx, &DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	})
	require.NoError(t, err)
	require.EqualValues(t, downloadResp.ContentMD5, md5Val[:])
	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
	require.EqualValues(t, downloadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)
}

func TestUploadBlobWithMD5WithCPKScope(t *testing.T) {
	t.Skip("encryption scope not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024
	r, srcData := generateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)
	require.Equal(t, *uploadResp.IsServerEncrypted, true)
	require.EqualValues(t, uploadResp.EncryptionScope, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := bbClient.BlobClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	require.EqualValues(t, downloadResp.ContentMD5, md5Val[:])
	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
	require.EqualValues(t, *downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

func TestAppendBlockWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	abClient := containerClient.NewAppendBlobClient(generateBlobName(testName))

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		resp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		require.NoError(t, err)
		require.Equal(t, resp.RawResponse.StatusCode, 201)
		require.Equal(t, *resp.BlobAppendOffset, strconv.Itoa(index*4))
		require.Equal(t, *resp.BlobCommittedBlockCount, int32(index+1))
		require.NotNil(t, resp.ETag)
		require.NotNil(t, resp.LastModified)
		require.Equal(t, resp.LastModified.IsZero(), false)
		require.NotEqual(t, resp.ContentMD5, "")

		require.NotEqual(t, resp.Version, "")
		require.NotNil(t, resp.Date)
		require.Equal(t, (*resp.Date).IsZero(), false)
		require.Equal(t, *resp.IsServerEncrypted, true)
		require.EqualValues(t, resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)
	}

	// Get blob content without encryption key should fail the request.
	_, err = abClient.Download(ctx, nil)
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := abClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)

	data, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, string(data), "AAA BBB CCC ")
	require.EqualValues(t, *downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

func TestAppendBlockWithCPKScope(t *testing.T) {
	t.Skip("encryption scope not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	abClient := containerClient.NewAppendBlobClient(generateBlobName(testName))

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		resp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		require.NoError(t, err)
		require.Equal(t, resp.RawResponse.StatusCode, 201)
		require.Equal(t, *resp.BlobAppendOffset, strconv.Itoa(index*4))
		require.Equal(t, *resp.BlobCommittedBlockCount, int32(index+1))
		require.NotNil(t, resp.ETag)
		require.NotNil(t, resp.LastModified)
		require.Equal(t, resp.LastModified.IsZero(), false)
		require.NotEqual(t, resp.ContentMD5, "")

		require.NotEqual(t, resp.Version, "")
		require.NotNil(t, resp.Date)
		require.Equal(t, (*resp.Date).IsZero(), false)
		require.Equal(t, *resp.IsServerEncrypted, true)
		require.EqualValues(t, resp.EncryptionScope, testCPKByScope.EncryptionScope)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := abClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)

	data, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, string(data), "AAA BBB CCC ")
	require.EqualValues(t, *downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

func TestAppendBlockFromURLWithCPK(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcABClient := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcABClient.Create(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, cResp1.RawResponse.StatusCode, 201)

	resp, err := srcABClient.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)
	require.Equal(t, *resp.BlobAppendOffset, "0")
	require.Equal(t, *resp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.LastModified)
	require.Equal(t, (*resp.LastModified).IsZero(), false)
	require.Nil(t, resp.ContentMD5)
	require.NotNil(t, resp.RequestID)
	require.NotNil(t, resp.Version)
	require.NotNil(t, resp.Date)
	require.Equal(t, (*resp.Date).IsZero(), false)

	srcBlobParts := NewBlobURLParts(srcABClient.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	require.NoError(t, err)
	require.Equal(t, cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:  &offset,
		Count:   &count,
		CpkInfo: &testCPKByValue,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	require.NoError(t, err)
	require.Equal(t, appendFromURLResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendFromURLResp.BlobAppendOffset, "0")
	require.Equal(t, *appendFromURLResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendFromURLResp.ETag)
	require.NotNil(t, appendFromURLResp.LastModified)
	require.Equal(t, (*appendFromURLResp.LastModified).IsZero(), false)
	require.NotNil(t, appendFromURLResp.ContentMD5)
	require.EqualValues(t, appendFromURLResp.ContentMD5, contentMD5)
	require.NotNil(t, appendFromURLResp.RequestID)
	require.NotNil(t, appendFromURLResp.Version)
	require.NotNil(t, appendFromURLResp.Date)
	require.Equal(t, (*appendFromURLResp.Date).IsZero(), false)
	require.Equal(t, *appendFromURLResp.IsServerEncrypted, true)

	// Get blob content without encryption key should fail the request.
	_, err = destBlob.Download(ctx, nil)
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destBlob.Download(ctx, &downloadBlobOptions)
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)

	require.Equal(t, *downloadResp.IsServerEncrypted, true)
	require.EqualValues(t, *downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
}

func TestAppendBlockFromURLWithCPKScope(t *testing.T) {
	t.Skip("encryption scope not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcClient := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcClient.Create(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, cResp1.RawResponse.StatusCode, 201)

	resp, err := srcClient.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)
	require.Equal(t, *resp.BlobAppendOffset, "0")
	require.Equal(t, *resp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.LastModified)
	require.Equal(t, (*resp.LastModified).IsZero(), false)
	require.Nil(t, resp.ContentMD5)
	require.NotNil(t, resp.RequestID)
	require.NotNil(t, resp.Version)
	require.NotNil(t, resp.Date)
	require.Equal(t, (*resp.Date).IsZero(), false)

	srcBlobParts := NewBlobURLParts(srcClient.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	require.NoError(t, err)
	require.Equal(t, cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:       &offset,
		Count:        &count,
		CpkScopeInfo: &testCPKByScope,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	require.NoError(t, err)
	require.Equal(t, appendFromURLResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendFromURLResp.BlobAppendOffset, "0")
	require.Equal(t, *appendFromURLResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendFromURLResp.ETag)
	require.NotNil(t, appendFromURLResp.LastModified)
	require.Equal(t, (*appendFromURLResp.LastModified).IsZero(), false)
	require.NotNil(t, appendFromURLResp.ContentMD5)
	require.EqualValues(t, appendFromURLResp.ContentMD5, contentMD5)
	require.NotNil(t, appendFromURLResp.RequestID)
	require.NotNil(t, appendFromURLResp.Version)
	require.NotNil(t, appendFromURLResp.Date)
	require.Equal(t, (*appendFromURLResp.Date).IsZero(), false)
	require.Equal(t, *appendFromURLResp.IsServerEncrypted, true)

	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	require.Equal(t, *downloadResp.IsServerEncrypted, true)
	require.EqualValues(t, *downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
}

func TestPageBlockWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := generateData(contentSize)
	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(t, pbName, containerClient, int64(contentSize), &testCPKByValue, nil)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		CpkInfo:   &testCPKByValue,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)
	require.EqualValues(t, uploadResp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	require.NoError(t, err)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), int64(contentSize-1)
	rawStart, rawEnd := pageListResp[0].Raw()
	require.Equal(t, rawStart, start)
	require.Equal(t, rawEnd, end)

	// Get blob content without encryption key should fail the request.
	_, err = pbClient.Download(ctx, nil)
	require.Error(t, err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = pbClient.Download(ctx, &downloadBlobOptions)
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := pbClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)

	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
	require.EqualValues(t, *downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)
}

func TestPageBlockWithCPKScope(t *testing.T) {
	t.Skip("encryption scope not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := generateData(contentSize)
	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(t, pbName, containerClient, int64(contentSize), nil, &testCPKByScope)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange:    &HttpRange{offset, count},
		CpkScopeInfo: &testCPKByScope,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)
	require.EqualValues(t, uploadResp.EncryptionScope, testCPKByScope.EncryptionScope)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	require.NoError(t, err)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), int64(contentSize-1)
	// _assert((*pageListResp)[0], chk.DeepEquals, PageRange{Start: &start, End: &end})
	rawStart, rawEnd := pageListResp[0].Raw()
	require.Equal(t, rawStart, start)
	require.Equal(t, rawEnd, end)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := pbClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)

	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
	require.EqualValues(t, *downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)
}

func TestPageBlockFromURLWithCPK(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024 // 1MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	srcPBName := "src" + generateBlobName(testName)
	bbClient := createNewPageBlobWithSize(t, srcPBName, containerClient, int64(contentSize))
	dstPBName := "dst" + generateBlobName(testName)
	destBlob := createNewPageBlobWithCPK(t, dstPBName, containerClient, int64(contentSize), &testCPKByValue, nil)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := bbClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)
	srcBlobParts := NewBlobURLParts(bbClient.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)
	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.LastModified)
	require.NotNil(t, resp.ContentMD5)
	require.EqualValues(t, resp.ContentMD5, contentMD5)
	require.NotNil(t, resp.RequestID)
	require.NotNil(t, resp.Version)
	require.NotNil(t, resp.Date)
	require.Equal(t, (*resp.Date).IsZero(), false)
	require.Equal(t, *resp.BlobSequenceNumber, int64(0))
	require.Equal(t, *resp.IsServerEncrypted, true)
	require.EqualValues(t, resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	_, err = destBlob.Download(ctx, nil)
	require.Error(t, err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destBlob.Download(ctx, &downloadBlobOptions)
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	require.EqualValues(t, *downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
}

func TestPageBlockFromURLWithCPKScope(t *testing.T) {
	t.Skip("encryption scope not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024 // 1MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	srcPBName := "src" + generateBlobName(testName)
	srcPBClient := createNewPageBlobWithSize(t, srcPBName, containerClient, int64(contentSize))
	dstPBName := "dst" + generateBlobName(testName)
	dstPBBlob := createNewPageBlobWithCPK(t, dstPBName, containerClient, int64(contentSize), nil, &testCPKByScope)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := srcPBClient.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)
	srcBlobParts := NewBlobURLParts(srcPBClient.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkScopeInfo:     &testCPKByScope,
	}
	resp, err := dstPBBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)
	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.LastModified)
	require.NotNil(t, resp.ContentMD5)
	require.EqualValues(t, resp.ContentMD5, contentMD5)
	require.NotNil(t, resp.RequestID)
	require.NotNil(t, resp.Version)
	require.NotNil(t, resp.Date)
	require.Equal(t, (*resp.Date).IsZero(), false)
	require.Equal(t, *resp.BlobSequenceNumber, int64(0))
	require.Equal(t, *resp.IsServerEncrypted, true)
	require.EqualValues(t, resp.EncryptionScope, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := dstPBBlob.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	require.EqualValues(t, *downloadResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)
}

func TestUploadPagesFromURLWithMD5WithCPK(t *testing.T) {
	recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	contentSize := 8 * 1024
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	srcPBName := "src" + generateBlobName(testName)
	srcBlob := createNewPageBlobWithSize(t, srcPBName, containerClient, int64(contentSize))

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := srcBlob.UploadPages(ctx, internal.NopCloser(r), &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)

	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()
	dstPBName := "dst" + generateBlobName(testName)
	destPBClient := createNewPageBlobWithCPK(t, dstPBName, containerClient, int64(contentSize), &testCPKByValue, nil)
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)
	require.NotNil(t, resp.ETag)
	require.NotNil(t, resp.LastModified)
	require.NotNil(t, resp.ContentMD5)
	require.EqualValues(t, resp.ContentMD5, contentMD5)
	require.NotNil(t, resp.RequestID)
	require.NotNil(t, resp.Version)
	require.NotNil(t, resp.Date)
	require.Equal(t, (*resp.Date).IsZero(), false)
	require.Equal(t, *resp.BlobSequenceNumber, int64(0))
	require.Equal(t, *resp.IsServerEncrypted, true)
	require.EqualValues(t, resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	_, err = destPBClient.Download(ctx, nil)
	require.Error(t, err)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = destPBClient.Download(ctx, &downloadBlobOptions)
	require.Error(t, err)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := destPBClient.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	require.EqualValues(t, *downloadResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(&RetryReaderOptions{CpkInfo: &testCPKByValue}))
	require.NoError(t, err)
	require.EqualValues(t, destData, srcData)

	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions1 := UploadPagesFromURLOptions{
		SourceContentMD5: badContentMD5,
	}
	_, err = destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions1)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeMD5Mismatch)
}

func TestClearDiffPagesWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(t, pbName, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)

	contentSize := 2 * 1024
	r := getReaderToGeneratedBytes(contentSize)
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}, CpkInfo: &testCPKByValue}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	require.NoError(t, err)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkInfo: &testCPKByValue,
	}
	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	require.NoError(t, err)

	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}, CpkInfo: &testCPKByValue}
	_, err = pbClient.UploadPages(context.Background(), getReaderToGeneratedBytes(2048), &uploadPagesOptions1)
	require.NoError(t, err)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
	require.NoError(t, err)
	pageRangeResp := pageListResp.PageList.PageRange
	require.NotNil(t, pageRangeResp)
	require.Len(t, pageRangeResp, 1)
	rawStart, rawEnd := pageRangeResp[0].Raw()
	require.Equal(t, rawStart, offset1)
	require.Equal(t, rawEnd, end1)

	clearPagesOptions := ClearPagesOptions{
		CpkInfo: &testCPKByValue,
	}
	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, &clearPagesOptions)
	require.NoError(t, err)
	require.Equal(t, clearResp.RawResponse.StatusCode, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
	require.NoError(t, err)
	require.Nil(t, pageListResp.PageList.PageRange)
}

func TestBlobResizeWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	pbName := generateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(t, pbName, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	_, err = pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	require.NoError(t, err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, _ := pbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.Equal(t, *resp.ContentLength, int64(PageBlobPageBytes))
}

func TestGetSetBlobMetadataWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(t, bbName, containerClient, &testCPKByValue, nil)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.Error(t, err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	require.NoError(t, err)
	require.EqualValues(t, resp.EncryptionKeySHA256, testCPKByValue.EncryptionKeySHA256)

	// Get blob properties without encryption key should fail the request.
	_, err = bbClient.GetProperties(ctx, nil)
	require.Error(t, err)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.NoError(t, err)
	require.NotNil(t, getResp.Metadata)
	require.Len(t, getResp.Metadata, len(basicMetadata))
	require.EqualValues(t, getResp.Metadata, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	require.NoError(t, err)

	getResp, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	require.NoError(t, err)
	require.Nil(t, getResp.Metadata)
}

func TestGetSetBlobMetadataWithCPKScope(t *testing.T) {
	t.Skip("encryption scope is not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(t, bbName, containerClient, nil, &testCPKByScope)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(ctx, basicMetadata, nil)
	require.Error(t, err)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	require.NoError(t, err)
	require.EqualValues(t, resp.EncryptionScope, testCPKByScope.EncryptionScope)

	getResp, err := bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, getResp.Metadata)
	require.Len(t, getResp.Metadata, len(basicMetadata))
	require.EqualValues(t, getResp.Metadata, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	require.NoError(t, err)

	getResp, err = bbClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, getResp.Metadata)
}

func TestBlobSnapshotWithCPK(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(t, bbName, containerClient, &testCPKByValue, nil)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(ctx, nil)
	require.Error(t, err)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	require.Error(t, err)

	createBlobSnapshotOptions1 := CreateBlobSnapshotOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	require.NoError(t, err)
	require.Equal(t, *resp.IsServerEncrypted, false)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	require.EqualValues(t, *dResp.EncryptionKeySHA256, *testCPKByValue.EncryptionKeySHA256)

	_, err = snapshotURL.Delete(ctx, nil)
	require.NoError(t, err)

	// Get blob properties of snapshot without encryption key should fail the request.
	_, err = snapshotURL.GetProperties(ctx, nil)
	require.Error(t, err)
}

func TestBlobSnapshotWithCPKScope(t *testing.T) {
	t.Skip("encryption scope is not available")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName)+"01", svcClient)
	defer deleteContainer(t, containerClient)

	bbName := generateBlobName(testName)
	bbClient := createNewBlockBlobWithCPK(t, bbName, containerClient, nil, &testCPKByScope)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(ctx, nil)
	require.Error(t, err)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkScopeInfo: &testInvalidCPKByScope,
	}
	_, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	require.Error(t, err)

	createBlobSnapshotOptions1 := CreateBlobSnapshotOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	require.NoError(t, err)
	require.Equal(t, *resp.IsServerEncrypted, false)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	require.NoError(t, err)
	require.EqualValues(t, *dResp.EncryptionScope, *testCPKByScope.EncryptionScope)

	_, err = snapshotURL.Delete(ctx, nil)
	require.NoError(t, err)
}
