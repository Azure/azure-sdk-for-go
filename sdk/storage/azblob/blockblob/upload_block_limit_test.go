// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stretchr/testify/require"
)

type countingTransport struct{ calls atomic.Int32 }

func (c *countingTransport) Do(req *http.Request) (*http.Response, error) {
	c.calls.Add(1)
	return &http.Response{StatusCode: http.StatusCreated, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
}

// A block count above 65,535 once wrapped around in a uint16 conversion, which bypassed the
// MaxBlocks check and panicked while staging the first block.
func TestUploadBufferBlockLimitExceeded(t *testing.T) {
	for _, size := range []int{blockblob.MaxBlocks + 1, 1<<16 + 1, 64 * 1024} {
		tr := &countingTransport{}
		client, err := blockblob.NewClientWithNoCredential("https://account.blob.core.windows.net/container/blob", &blockblob.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
		})
		require.NoError(t, err)

		_, err = client.UploadBuffer(context.Background(), make([]byte, size), &blockblob.UploadBufferOptions{BlockSize: 1})
		require.ErrorContains(t, err, "block limit exceeded", "size %d", size)
		require.Zero(t, tr.calls.Load(), "nothing is uploaded for size %d", size)
	}
}

// A negative block size once sent small content down the staged path, where the negative block
// count panicked.
func TestUploadBufferNegativeBlockSize(t *testing.T) {
	for _, size := range []int{10, 64 * 1024} {
		tr := &countingTransport{}
		client, err := blockblob.NewClientWithNoCredential("https://account.blob.core.windows.net/container/blob", &blockblob.ClientOptions{
			ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
		})
		require.NoError(t, err)

		_, err = client.UploadBuffer(context.Background(), make([]byte, size), &blockblob.UploadBufferOptions{BlockSize: -1})
		require.ErrorContains(t, err, "BlockSize must not be negative", "size %d", size)
		require.Zero(t, tr.calls.Load(), "nothing is uploaded for size %d", size)
	}
}

// headerTransport records the headers of each request by its comp query value.
type headerTransport struct {
	mu      sync.Mutex
	headers map[string][]http.Header
}

func (h *headerTransport) Do(req *http.Request) (*http.Response, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.headers == nil {
		h.headers = map[string][]http.Header{}
	}
	comp := req.URL.Query().Get("comp")
	// the generated code may set headers under non-canonical keys, so canonicalize them
	canonical := http.Header{}
	for k, v := range req.Header {
		canonical[http.CanonicalHeaderKey(k)] = append(canonical[http.CanonicalHeaderKey(k)], v...)
	}
	h.headers[comp] = append(h.headers[comp], canonical)
	return &http.Response{StatusCode: http.StatusCreated, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
}

// Content larger than the block size is staged; the commit must keep the upload's access
// conditions, as the single-request path does.
func TestUploadBufferStagedCommitKeepsAccessConditions(t *testing.T) {
	tr := &headerTransport{}
	client, err := blockblob.NewClientWithNoCredential("https://account.blob.core.windows.net/container/blob", &blockblob.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: tr, Retry: policy.RetryOptions{MaxRetries: -1}},
	})
	require.NoError(t, err)

	etag := azcore.ETag("\"0x8D\"")
	_, err = client.UploadBuffer(context.Background(), make([]byte, 8*1024*1024), &blockblob.UploadBufferOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
			LeaseAccessConditions:    &blob.LeaseAccessConditions{LeaseID: to.Ptr("lease-id")},
		},
	})
	require.NoError(t, err)

	require.Len(t, tr.headers["block"], 2, "8 MiB is staged in two default-size blocks")
	for _, h := range tr.headers["block"] {
		require.Equal(t, "lease-id", h.Get("x-ms-lease-id"), "Put Block")
	}
	require.Len(t, tr.headers["blocklist"], 1)
	commit := tr.headers["blocklist"][0]
	require.Equal(t, string(etag), commit.Get("If-Match"), "Put Block List keeps If-Match")
	require.Equal(t, "lease-id", commit.Get("x-ms-lease-id"), "Put Block List keeps the lease")
}
