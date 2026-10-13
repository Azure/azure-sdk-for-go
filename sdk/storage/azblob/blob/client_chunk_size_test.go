// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests cover the size of the ranges DownloadBuffer and DownloadFile request: the initial range
// defaults to the same size as the subsequent ranges, and an explicit BlockSize applies to both.

// expectedRanges returns the range headers a download of size bytes in blockSize chunks requests.
// The initial range asks for a whole block, as the blob size isn't known yet. The remaining ranges
// stop at the end of the blob, and may be requested in any order so they're sorted.
func expectedRanges(blockSize, size int64) []string {
	ranges := []string{fmt.Sprintf("bytes=0-%d", blockSize-1)}
	for offset := blockSize; offset < size; offset += blockSize {
		ranges = append(ranges, fmt.Sprintf("bytes=%d-%d", offset, min(offset+blockSize, size)-1))
	}
	slices.Sort(ranges[1:])
	return ranges
}

func requireDownloadRanges(t *testing.T, f *fakeLayoutResponder, expected []string) {
	t.Helper()
	actual := f.downloadRangeHeaders()
	require.NotEmpty(t, actual)
	slices.Sort(actual[1:])
	require.Equal(t, expected, actual)
}

func TestDownloadBufferDefaultInitialRangeSize(t *testing.T) {
	for _, size := range []int64{1, DefaultDownloadBlockSize - 1, DefaultDownloadBlockSize, DefaultDownloadBlockSize + 1, 2*DefaultDownloadBlockSize + 1} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			f := &fakeLayoutResponder{size: size, eTag: "etag", noDownloadHint: true}
			client := newFakeLayoutClient(t, f)

			n, err := client.DownloadBuffer(context.Background(), make([]byte, size), nil)
			require.NoError(t, err)
			require.Equal(t, size, n)
			// The initial range is the default block size, the same as every subsequent range, so blobs no
			// larger than one block are downloaded with a single request.
			requireDownloadRanges(t, f, expectedRanges(DefaultDownloadBlockSize, size))
		})
	}
}

func TestDownloadBufferExplicitBlockSizeInitialRange(t *testing.T) {
	const blockSize = 1024 * 1024
	size := int64(2*blockSize + 1)
	f := &fakeLayoutResponder{size: size, eTag: "etag", noDownloadHint: true}
	client := newFakeLayoutClient(t, f)

	n, err := client.DownloadBuffer(context.Background(), make([]byte, size), &DownloadBufferOptions{
		BlockSize: blockSize,
	})
	require.NoError(t, err)
	require.Equal(t, size, n)
	requireDownloadRanges(t, f, expectedRanges(blockSize, size))
}

func TestDownloadFileDefaultInitialRangeSize(t *testing.T) {
	size := 2*DefaultDownloadBlockSize + 1
	f := &fakeLayoutResponder{size: size, eTag: "etag", noDownloadHint: true}
	client := newFakeLayoutClient(t, f)

	file, err := os.CreateTemp(t.TempDir(), "download")
	require.NoError(t, err)
	defer file.Close()

	n, err := client.DownloadFile(context.Background(), file, nil)
	require.NoError(t, err)
	require.Equal(t, size, n)
	requireDownloadRanges(t, f, expectedRanges(DefaultDownloadBlockSize, size))
}
