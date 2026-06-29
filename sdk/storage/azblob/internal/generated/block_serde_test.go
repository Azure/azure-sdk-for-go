// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBlockListUnmarshalXMLWithHashes verifies that GetBlockList responses
// containing the per-block Offset, Crc64, and Sha256 elements are decoded into
// the Block model, with the base64-encoded hashes decoded to raw bytes.
func TestBlockListUnmarshalXMLWithHashes(t *testing.T) {
	const body = `<?xml version="1.0" encoding="utf-8"?>
<BlockList>
  <CommittedBlocks>
    <Block>
      <Name>YmxvY2stMDAwMDAx</Name>
      <Size>4194304</Size>
      <Offset>0</Offset>
      <Crc64>q83vEjRWeJA=</Crc64>
      <Sha256>7j3H4Y7l0QFvZ2f9Pq6c8oP2Y9W3xV8QmQkF0x9abcA=</Sha256>
    </Block>
    <Block>
      <Name>YmxvY2stMDAwMDAy</Name>
      <Size>8388608</Size>
      <Offset>4194304</Offset>
      <Crc64>AbCdEfGhIjk=</Crc64>
      <Sha256>rM2L7p8s9q0tU1vW2xY3zA4bC5dE6fG7hI8jK9lMnO0=</Sha256>
    </Block>
  </CommittedBlocks>
</BlockList>`

	var bl BlockList
	require.NoError(t, xml.Unmarshal([]byte(body), &bl))
	require.Len(t, bl.CommittedBlocks, 2)
	require.Empty(t, bl.UncommittedBlocks)

	first := bl.CommittedBlocks[0]
	require.Equal(t, "YmxvY2stMDAwMDAx", *first.Name)
	require.Equal(t, int64(4194304), *first.Size)
	require.Equal(t, int64(0), *first.Offset)
	require.Len(t, first.Crc64, 8)
	require.Len(t, first.Sha256, 32)
	require.Equal(t, "q83vEjRWeJA=", base64.StdEncoding.EncodeToString(first.Crc64))
	require.Equal(t, "7j3H4Y7l0QFvZ2f9Pq6c8oP2Y9W3xV8QmQkF0x9abcA=", base64.StdEncoding.EncodeToString(first.Sha256))

	second := bl.CommittedBlocks[1]
	require.Equal(t, "YmxvY2stMDAwMDAy", *second.Name)
	require.Equal(t, int64(8388608), *second.Size)
	require.Equal(t, int64(4194304), *second.Offset)
	require.Equal(t, "AbCdEfGhIjk=", base64.StdEncoding.EncodeToString(second.Crc64))
	require.Equal(t, "rM2L7p8s9q0tU1vW2xY3zA4bC5dE6fG7hI8jK9lMnO0=", base64.StdEncoding.EncodeToString(second.Sha256))
}

// TestBlockListUnmarshalXMLWithoutHashes verifies that the legacy GetBlockList
// response (without Offset/Crc64/Sha256) still decodes, leaving the new fields
// nil for backward compatibility.
func TestBlockListUnmarshalXMLWithoutHashes(t *testing.T) {
	const body = `<?xml version="1.0" encoding="utf-8"?>
<BlockList>
  <CommittedBlocks>
    <Block>
      <Name>YmxvY2stMDAwMDAx</Name>
      <Size>4194304</Size>
    </Block>
  </CommittedBlocks>
</BlockList>`

	var bl BlockList
	require.NoError(t, xml.Unmarshal([]byte(body), &bl))
	require.Len(t, bl.CommittedBlocks, 1)

	block := bl.CommittedBlocks[0]
	require.Equal(t, "YmxvY2stMDAwMDAx", *block.Name)
	require.Equal(t, int64(4194304), *block.Size)
	require.Nil(t, block.Offset)
	require.Nil(t, block.Crc64)
	require.Nil(t, block.Sha256)
}

// TestGetBlockListCreateRequestInclude verifies that GetBlockList emits the
// include=crc64,sha256 query parameter (which XStore gates the hashes on) only
// when the Include option is set, and omits it otherwise for back-compat.
func TestGetBlockListCreateRequestInclude(t *testing.T) {
	client := &BlockBlobClient{endpoint: "https://acct.blob.core.windows.net/container/blob"}

	req, err := client.getBlockListCreateRequest(
		context.Background(),
		BlockListTypeCommitted,
		&BlockBlobClientGetBlockListOptions{
			Include: []BlockListIncludeItem{BlockListIncludeItemCrc64, BlockListIncludeItemSha256},
		},
		nil, nil)
	require.NoError(t, err)
	q := req.Raw().URL.Query()
	require.Equal(t, "committed", q.Get("blocklisttype"))
	require.Equal(t, "blocklist", q.Get("comp"))
	require.Equal(t, "crc64,sha256", q.Get("include"))

	req, err = client.getBlockListCreateRequest(
		context.Background(),
		BlockListTypeCommitted,
		nil, nil, nil)
	require.NoError(t, err)
	require.Empty(t, req.Raw().URL.Query().Get("include"))
}
