//go:build go1.18 && !race
// +build go1.18,!race

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// these tests are excluded from the race detector for various reasons (e.g. causes OOM)

package blockblob_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/require"
)

func (s *BlockBlobUnrecordedTestsSuite) TestLargeBlockBufferedUploadInParallel() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	var largeBlockSize, numberOfBlocks int64 = 2500 * 1024 * 1024, 2
	content := make([]byte, numberOfBlocks*largeBlockSize)

	_, err = bbClient.UploadBuffer(context.Background(), content, &blockblob.UploadBufferOptions{
		BlockSize:   largeBlockSize,
		Concurrency: 2,
	})
	_require.Nil(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Len(resp.BlockList.CommittedBlocks, 2)
	_require.Equal(*resp.BlobContentLength, numberOfBlocks*largeBlockSize)
	committed := resp.BlockList.CommittedBlocks
	_require.Equal(*(committed[0].Size), largeBlockSize)
	_require.Equal(*(committed[1].Size), largeBlockSize)
}

func (s *BlockBlobUnrecordedTestsSuite) TestLargeBlockStreamUploadWithDifferentBlockSize() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	var firstBlockSize, secondBlockSize int64 = 2500 * 1024 * 1024, 10 * 1024 * 1024
	content := make([]byte, firstBlockSize+secondBlockSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.UploadStream(context.Background(), rsc, &blockblob.UploadStreamOptions{
		BlockSize:   firstBlockSize,
		Concurrency: 2,
	})
	_require.Nil(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Len(resp.BlockList.CommittedBlocks, 2)
	_require.Equal(*resp.BlobContentLength, firstBlockSize+secondBlockSize)
	committed := resp.BlockList.CommittedBlocks
	_require.Equal(*(committed[0].Size), firstBlockSize)
	_require.Equal(*(committed[1].Size), secondBlockSize)
}

func (s *BlockBlobUnrecordedTestsSuite) TestLargeBlockBlobStage() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	const largeBlockSize = blockblob.MaxStageBlockBytes
	content := make([]byte, largeBlockSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	_, err = bbClient.StageBlock(context.Background(), blockID, rsc, nil)
	_require.Nil(err)

	_, err = bbClient.CommitBlockList(context.Background(), []string{blockID}, nil)
	_require.Nil(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.Nil(err)
	_require.Len(resp.BlockList.CommittedBlocks, 1)
	committed := resp.BlockList.CommittedBlocks
	_require.Equal(*(committed[0].Name), blockID)
	_require.Equal(*(committed[0].Size), largeBlockSize)
	_require.Nil(resp.BlockList.UncommittedBlocks)
}

func TestUploadBufferUnevenBlockSize(t *testing.T) {
	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/path", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: fbb,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// create fake source that's not evenly divisible by 50000 (max number of blocks)
	// and greater than MaxUploadBlobBytes (256MB) so that it doesn't fit into a single upload.

	buffer := make([]byte, 263*1024*1024)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 1
	}

	_, err = client.UploadBuffer(context.Background(), buffer, &blockblob.UploadBufferOptions{
		Concurrency: 1,
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(buffer)), fbb.totalStaged)
}

func TestUploadBufferEvenBlockSize(t *testing.T) {
	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/path", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: fbb,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// create fake source that's evenly divisible by 50000 (max number of blocks)
	// and greater than MaxUploadBlobBytes (256MB) so that it doesn't fit into a single upload.

	buffer := make([]byte, 270000000)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 1
	}

	_, err = client.UploadBuffer(context.Background(), buffer, &blockblob.UploadBufferOptions{
		Concurrency: 1,
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(buffer)), fbb.totalStaged)
}

func TestUploadLogEvent(t *testing.T) {
	listnercalled := false
	log.SetEvents(azblob.EventUpload, log.EventRequest, log.EventResponse)
	log.SetListener(func(cls log.Event, msg string) {
		t.Logf("%s: %s\n", cls, msg)
		if cls == azblob.EventUpload {
			listnercalled = true
			require.Equal(t, msg, "blob name path1/path2 actual size 270000000 block-size 4194304 block-count 65")
		}
	})

	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/path1/path2", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: fbb,
			Telemetry: policy.TelemetryOptions{ApplicationID: "testApp/1.0.0-preview.2"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// create fake source that's evenly divisible by 50000 (max number of blocks)
	// and greater than MaxUploadBlobBytes (256MB) so that it doesn't fit into a single upload.

	buffer := make([]byte, 270000000)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 1
	}

	_, err = client.UploadBuffer(context.Background(), buffer, &blockblob.UploadBufferOptions{
		Concurrency: 1,
		Progress: func(bytesTransferred int64) {
			t.Logf("%v percent job done", (bytesTransferred*100)/270000000)
		},
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(buffer)), fbb.totalStaged)
	require.Equal(t, listnercalled, true)
}
