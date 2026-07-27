//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stretchr/testify/require"
)

type getBlockListCPUTimeHeadersPolicy struct{}

func (getBlockListCPUTimeHeadersPolicy) Do(req *policy.Request) (*http.Response, error) {
	headers := http.Header{}
	headers.Set("x-ms-test-dedupe-crc64-cpu-time-us", "123")
	headers.Set("x-ms-test-dedupe-sha256-cpu-time-us", "456")

	const body = `<?xml version="1.0" encoding="utf-8"?>
<BlockList>
  <CommittedBlocks>
    <Block>
      <Name>YmxvY2swMDAwMDE=</Name>
      <Size>20</Size>
    </Block>
  </CommittedBlocks>
</BlockList>`

	return &http.Response{
		Request:    req.Raw(),
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func TestGetBlockListCPUTimeHeaders(t *testing.T) {
	client, err := blockblob.NewClientWithNoCredential("https://fake.blob.core.windows.net/container/blob", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			PerCallPolicies: []policy.Policy{getBlockListCPUTimeHeadersPolicy{}},
		},
	})
	require.NoError(t, err)

	response, err := client.GetBlockList(context.Background(), blockblob.BlockListTypeCommitted, nil)
	require.NoError(t, err)

	var crc64CPUTimeUS *int64 = response.CRC64CPUTimeUS
	var sha256CPUTimeUS *int64 = response.SHA256CPUTimeUS
	require.Equal(t, int64(123), *crc64CPUTimeUS)
	require.Equal(t, int64(456), *sha256CPUTimeUS)

	require.Len(t, response.BlockList.CommittedBlocks, 1)
	require.Equal(t, "YmxvY2swMDAwMDE=", *response.BlockList.CommittedBlocks[0].Name)
	require.Equal(t, int64(20), *response.BlockList.CommittedBlocks[0].Size)
}
