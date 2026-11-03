// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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
