// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

// rangeServer is a transport that serves ranged reads of content, the way the service does.
type rangeServer struct {
	content []byte

	mu     sync.Mutex
	ranges []string
}

func (r *rangeServer) Do(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet {
		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	}
	rng := rawHeaderValue(req.Header, "x-ms-range")
	r.mu.Lock()
	r.ranges = append(r.ranges, rng)
	r.mu.Unlock()

	size := int64(len(r.content))
	start, end := requestedRange(rng)
	if end >= size {
		end = size - 1
	}
	header := http.Header{}
	header.Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
	header.Set("ETag", `"etag"`)
	return &http.Response{
		StatusCode: http.StatusPartialContent,
		Body:       io.NopCloser(bytes.NewReader(r.content[start : end+1])),
		Header:     header,
	}, nil
}

// lockedWriterAt is an io.WriterAt that is neither a []byte nor an *os.File, like a caller's
// wrapper around a file.
type lockedWriterAt struct {
	mu  sync.Mutex
	buf []byte
}

func (w *lockedWriterAt) WriteAt(p []byte, off int64) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if end := int(off) + len(p); end > len(w.buf) {
		w.buf = append(w.buf, make([]byte, end-len(w.buf))...)
	}
	return copy(w.buf[off:], p), nil
}

func TestDownloadWriterAt(t *testing.T) {
	content := make([]byte, 1050)
	for i := range content {
		content[i] = byte(i % 251)
	}
	srv := &rangeServer{content: content}
	client, err := NewClientWithNoCredential("https://fake/container/blob", &ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: srv},
	})
	require.NoError(t, err)

	w := &lockedWriterAt{}
	n, err := client.DownloadWriterAt(context.Background(), w, &DownloadWriterAtOptions{
		BlockSize:   100,
		Concurrency: 4,
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(content)), n)
	require.Equal(t, content, w.buf)
	require.Len(t, srv.ranges, 11, "the initial read plus one read per remaining block")
}

func TestDownloadWriterAtNilOptions(t *testing.T) {
	content := []byte("small blob")
	srv := &rangeServer{content: content}
	client, err := NewClientWithNoCredential("https://fake/container/blob", &ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: srv},
	})
	require.NoError(t, err)

	w := &lockedWriterAt{}
	n, err := client.DownloadWriterAt(context.Background(), w, nil)
	require.NoError(t, err)
	require.Equal(t, int64(len(content)), n)
	require.Equal(t, content, w.buf)
	require.Len(t, srv.ranges, 1, "a blob smaller than one block takes a single read")
}
