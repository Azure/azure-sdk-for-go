// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/stretchr/testify/require"
)

// These tests cover the size of the first chunk of partitioned uploads and downloads: it's the same as
// the size of the subsequent chunks, whether that's the default or an explicitly set ChunkSize.

// fakeChunkedFile serves appends, flushes and ranged reads of a file of the given size.
type fakeChunkedFile struct {
	size int64

	mu sync.Mutex
	// appendSizes records the size of each Append Data request, keyed by its position.
	appendSizes map[int64]int64
	flushCount  int
	// readRanges records the x-ms-range header of each read, in the order sent.
	readRanges []string
}

func (f *fakeChunkedFile) Do(req *http.Request) (*http.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case req.Method == http.MethodPatch && req.URL.Query().Get("action") == "append":
		n, err := io.Copy(io.Discard, req.Body)
		if err != nil {
			return nil, err
		}
		position, err := strconv.ParseInt(req.URL.Query().Get("position"), 10, 64)
		if err != nil {
			return nil, err
		}
		if f.appendSizes == nil {
			f.appendSizes = map[int64]int64{}
		}
		f.appendSizes[position] = n
		return &http.Response{Request: req, StatusCode: http.StatusAccepted, Header: http.Header{}, Body: http.NoBody}, nil
	case req.Method == http.MethodPatch && req.URL.Query().Get("action") == "flush":
		f.flushCount++
		return &http.Response{Request: req, StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody}, nil
	case req.Method == http.MethodGet:
		r := rawHeaderValue(req.Header, "x-ms-range")
		f.readRanges = append(f.readRanges, r)
		var start, end int64
		if _, err := fmt.Sscanf(r, "bytes=%d-%d", &start, &end); err != nil {
			return nil, err
		}
		end = min(end, f.size-1)
		header := http.Header{}
		header.Set("Content-Length", strconv.FormatInt(end-start+1, 10))
		header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, f.size))
		header.Set("ETag", "etag")
		return &http.Response{
			Request:    req,
			StatusCode: http.StatusPartialContent,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader(make([]byte, end-start+1))),
		}, nil
	}
	return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL)
}

// rawHeaderValue reads a header the way the generated clients write it. They assign directly
// into the header map using the wire casing, e.g. Header["x-ms-range"], which http.Header.Get
// cannot see because it canonicalizes the name to "X-Ms-Range".
func rawHeaderValue(h http.Header, name string) string {
	if v := h.Get(name); v != "" {
		return v
	}
	if v := h[name]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func newFakeChunkedFileClient(t *testing.T, f *fakeChunkedFile) *Client {
	t.Helper()
	client, err := NewClientWithNoCredential("https://fake.dfs.core.windows.net/fs/file", &ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: f},
	})
	require.NoError(t, err)
	return client
}

func TestDownloadBufferDefaultInitialChunkSize(t *testing.T) {
	size := 2*blob.DefaultDownloadBlockSize + 1
	f := &fakeChunkedFile{size: size}
	client := newFakeChunkedFileClient(t, f)

	n, err := client.DownloadBuffer(context.Background(), make([]byte, size), nil)
	require.NoError(t, err)
	require.Equal(t, size, n)

	// The initial read is the default chunk size, the same as the subsequent reads.
	require.Len(t, f.readRanges, 3)
	require.Equal(t, fmt.Sprintf("bytes=0-%d", blob.DefaultDownloadBlockSize-1), f.readRanges[0])
	remaining := f.readRanges[1:]
	slices.Sort(remaining)
	require.Equal(t, []string{
		fmt.Sprintf("bytes=%d-%d", blob.DefaultDownloadBlockSize, 2*blob.DefaultDownloadBlockSize-1),
		fmt.Sprintf("bytes=%d-%d", 2*blob.DefaultDownloadBlockSize, 2*blob.DefaultDownloadBlockSize),
	}, remaining)
}

func TestDownloadBufferExplicitChunkSizeInitialChunk(t *testing.T) {
	const chunkSize = 1024 * 1024
	f := &fakeChunkedFile{size: chunkSize + 1}
	client := newFakeChunkedFileClient(t, f)

	n, err := client.DownloadBuffer(context.Background(), make([]byte, f.size), &DownloadBufferOptions{ChunkSize: chunkSize})
	require.NoError(t, err)
	require.Equal(t, f.size, n)
	require.Equal(t, []string{fmt.Sprintf("bytes=0-%d", chunkSize-1), fmt.Sprintf("bytes=%d-%d", chunkSize, chunkSize)}, f.readRanges)
}

func TestUploadBufferDefaultInitialChunkSize(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		// Content no larger than the default chunk size is uploaded with a single append.
		f := &fakeChunkedFile{}
		client := newFakeChunkedFileClient(t, f)

		require.NoError(t, client.UploadBuffer(context.Background(), make([]byte, MaxAppendBytes), nil))
		require.Equal(t, map[int64]int64{0: MaxAppendBytes}, f.appendSizes)
		require.Equal(t, 1, f.flushCount)
	})

	t.Run("partitioned", func(t *testing.T) {
		// Larger content is appended in default sized chunks, starting with the first one.
		f := &fakeChunkedFile{}
		client := newFakeChunkedFileClient(t, f)

		require.NoError(t, client.UploadBuffer(context.Background(), make([]byte, MaxAppendBytes+1), nil))
		require.Equal(t, map[int64]int64{0: MaxAppendBytes, MaxAppendBytes: 1}, f.appendSizes)
		require.Equal(t, 1, f.flushCount)
	})
}

func TestUploadBufferExplicitChunkSizeInitialChunk(t *testing.T) {
	const chunkSize = 1024 * 1024
	f := &fakeChunkedFile{}
	client := newFakeChunkedFileClient(t, f)

	require.NoError(t, client.UploadBuffer(context.Background(), make([]byte, 2*chunkSize+1), &UploadBufferOptions{ChunkSize: chunkSize}))
	require.Equal(t, map[int64]int64{0: chunkSize, chunkSize: chunkSize, 2 * chunkSize: 1}, f.appendSizes)
	require.Equal(t, 1, f.flushCount)
}
