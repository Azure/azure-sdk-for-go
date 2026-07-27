// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const successfulBlockListXML = `<?xml version="1.0" encoding="utf-8"?>
<BlockList>
  <CommittedBlocks>
    <Block>
      <Name>YmxvY2swMDAwMDE=</Name>
      <Size>20</Size>
    </Block>
  </CommittedBlocks>
</BlockList>`

func TestGetBlockListCPUTimeHeaders(t *testing.T) {
	const (
		crc64Header  = "x-ms-test-dedupe-crc64-cpu-time-us"
		sha256Header = "x-ms-test-dedupe-sha256-cpu-time-us"
	)

	int64Ptr := func(value int64) *int64 {
		return &value
	}

	tests := []struct {
		name       string
		headers    map[string]string
		wantCRC64  *int64
		wantSHA256 *int64
	}{
		{
			name:       "both valid positive",
			headers:    map[string]string{crc64Header: "123", sha256Header: "456"},
			wantCRC64:  int64Ptr(123),
			wantSHA256: int64Ptr(456),
		},
		{
			name:       "zero",
			headers:    map[string]string{crc64Header: "0", sha256Header: "0"},
			wantCRC64:  int64Ptr(0),
			wantSHA256: int64Ptr(0),
		},
		{
			name:    "both missing",
			headers: map[string]string{},
		},
		{
			name:      "SHA256 missing",
			headers:   map[string]string{crc64Header: "123"},
			wantCRC64: int64Ptr(123),
		},
		{
			name:       "CRC64 missing",
			headers:    map[string]string{sha256Header: "456"},
			wantSHA256: int64Ptr(456),
		},
		{
			name:    "malformed",
			headers: map[string]string{crc64Header: "not-a-number", sha256Header: "1.5"},
		},
		{
			name:    "negative",
			headers: map[string]string{crc64Header: "-1", sha256Header: "-2"},
		},
		{
			name:    "int64 overflow",
			headers: map[string]string{crc64Header: "9223372036854775808", sha256Header: "9223372036854775808"},
		},
		{
			name:      "valid CRC64 and invalid SHA256",
			headers:   map[string]string{crc64Header: "123", sha256Header: "invalid"},
			wantCRC64: int64Ptr(123),
		},
		{
			name:       "invalid CRC64 and valid SHA256",
			headers:    map[string]string{crc64Header: "-1", sha256Header: "456"},
			wantSHA256: int64Ptr(456),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			headers := http.Header{}
			for name, value := range test.headers {
				headers.Set(name, value)
			}

			httpResponse := &http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader(successfulBlockListXML)),
			}

			response, err := (&BlockBlobClient{}).getBlockListHandleResponse(httpResponse)
			require.NoError(t, err)
			require.Equal(t, test.wantCRC64, response.CRC64CPUTimeUS)
			require.Equal(t, test.wantSHA256, response.SHA256CPUTimeUS)

			require.Len(t, response.BlockList.CommittedBlocks, 1)
			require.Equal(t, "YmxvY2swMDAwMDE=", *response.BlockList.CommittedBlocks[0].Name)
			require.Equal(t, int64(20), *response.BlockList.CommittedBlocks[0].Size)
		})
	}
}
