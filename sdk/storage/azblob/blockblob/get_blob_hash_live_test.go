//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"context"
	"crypto/sha256"
	"io"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stretchr/testify/require"
)

func TestGetBlobHashLiveKnownCommittedRanges(t *testing.T) {
	blobURL := os.Getenv("AZBLOB_XSTORE_GET_BLOB_HASH_URL")
	if blobURL == "" {
		t.Skip("set AZBLOB_XSTORE_GET_BLOB_HASH_URL to a readable XStore block blob SAS URL")
	}

	client, err := blockblob.NewClientWithNoCredential(blobURL, nil)
	require.NoError(t, err)

	blockList, err := client.GetBlockList(context.Background(), blockblob.BlockListTypeCommitted, &blockblob.GetBlockListOptions{
		Include: []blockblob.BlockListIncludeItem{blockblob.BlockListIncludeItemCrc64},
	})
	require.NoError(t, err)
	require.NotNil(t, blockList.ETag)

	ranges := make([]blockblob.BlobHashRange, 0, 4)
	var total int64
	for _, block := range blockList.CommittedBlocks {
		if block.Offset == nil || block.Size == nil || *block.Size <= 0 || *block.Size > maxBlobHashBytes {
			continue
		}
		if total > maxBlobHashBytes-*block.Size {
			break
		}
		ranges = append(ranges, blockblob.BlobHashRange{Offset: *block.Offset, Count: *block.Size})
		total += *block.Size
		if len(ranges) == cap(ranges) {
			break
		}
	}
	require.NotEmpty(t, ranges)

	etag := azcore.ETag(*blockList.ETag)
	conditions := &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
	}
	hashes, err := client.GetBlobHash(context.Background(), ranges, &blockblob.GetBlobHashOptions{
		AccessConditions: conditions,
	})
	require.NoError(t, err)
	require.Len(t, hashes.RangeHashes, len(ranges))

	returned := make(map[blockblob.BlobHashRange][]byte, len(hashes.RangeHashes))
	for _, result := range hashes.RangeHashes {
		returned[blockblob.BlobHashRange{Offset: result.Offset, Count: result.Count}] = result.SHA256
	}

	for _, rnge := range ranges {
		download, err := client.BlobClient().DownloadStream(context.Background(), &blob.DownloadStreamOptions{
			Range:            blob.HTTPRange{Offset: rnge.Offset, Count: rnge.Count},
			AccessConditions: conditions,
		})
		require.NoError(t, err)
		data, readErr := io.ReadAll(download.Body)
		closeErr := download.Body.Close()
		require.NoError(t, readErr)
		require.NoError(t, closeErr)

		expected := sha256.Sum256(data)
		require.Equal(t, expected[:], returned[rnge])
	}
}
