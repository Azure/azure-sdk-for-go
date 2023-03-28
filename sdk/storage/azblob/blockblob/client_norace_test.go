//go:build go1.18 && !race
// +build go1.18,!race

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// these tests are excluded from the race detector for various reasons (e.g. causes OOM)

package blockblob_test

import (
	"context"

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
