package storage

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// BlockListType is used to filter out types of blocks in a Get Blocks List call
// for a block blob.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179400.aspx for all
// block types.
type BlockListType string

// Filters for listing blocks in block blobs
const (
	BlockListTypeAll         BlockListType = "all"
	BlockListTypeCommitted   BlockListType = "committed"
	BlockListTypeUncommitted BlockListType = "uncommitted"
)

// Maximum sizes (per REST API) for various concepts
const (
	MaxBlobBlockSize = 100 * 1024 * 1024
	MaxBlobPageSize  = 4 * 1024 * 1024
)

// BlockStatus defines states a block for a block blob can
// be in.
type BlockStatus string

// List of statuses that can be used to refer to a block in a block list
const (
	BlockStatusUncommitted BlockStatus = "Uncommitted"
	BlockStatusCommitted   BlockStatus = "Committed"
	BlockStatusLatest      BlockStatus = "Latest"
)

// Block is used to create Block entities for Put Block List
// call.
type Block struct {
	ID     string
	Status BlockStatus
}

// BlockListResponse contains the response fields from Get Block List call.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179400.aspx
type BlockListResponse struct {
	XMLName           xml.Name        `xml:"BlockList"`
	CommittedBlocks   []BlockResponse `xml:"CommittedBlocks>Block"`
	UncommittedBlocks []BlockResponse `xml:"UncommittedBlocks>Block"`
}

// BlockResponse contains the block information returned
// in the GetBlockListCall.
type BlockResponse struct {
	Name string `xml:"Name"`
	Size int64  `xml:"Size"`
}

// CreateBlockBlob initializes an empty block blob with no blocks.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179451.aspx
func (b *Blob) CreateBlockBlob() error {
	return b.CreateBlockBlobFromReader(0, nil, nil)
}

// CreateBlockBlobFromReader initializes a block blob using data from
// reader. Size must be the number of bytes read from reader. To
// create an empty blob, use size==0 and reader==nil.
//
// The API rejects requests with size > 256 MiB (but this limit is not
// checked by the SDK). To write a larger blob, use CreateBlockBlob,
// PutBlock, and PutBlockList.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179451.aspx
func (b *Blob) CreateBlockBlobFromReader(size uint64, blob io.Reader, extraHeaders map[string]string) error {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), nil)
	extraHeaders = b.Container.bsc.client.protectUserAgent(extraHeaders)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeBlock)
	headers["Content-Length"] = strconv.FormatUint(size, 10)

	for k, v := range extraHeaders {
		headers[k] = v
	}

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, blob, b.Container.bsc.auth)
	if err != nil {
		return err
	}
	readAndCloseBody(resp.body)
	return checkRespCode(resp.statusCode, []int{http.StatusCreated})
}

// PutBlock saves the given data chunk to the specified block blob with
// given ID.
//
// The API rejects chunks larger than 100 MB (but this limit is not
// checked by the SDK).
//
// See https://msdn.microsoft.com/en-us/library/azure/dd135726.aspx
func (b *Blob) PutBlock(blockID string, chunk []byte) error {
	return b.PutBlockWithLength(blockID, uint64(len(chunk)), bytes.NewReader(chunk), nil)
}

// PutBlockWithLength saves the given data stream of exactly specified size to
// the block blob with given ID. It is an alternative to PutBlocks where data
// comes as stream but the length is known in advance.
//
// The API rejects requests with size > 100 MB (but this limit is not
// checked by the SDK).
//
// See https://msdn.microsoft.com/en-us/library/azure/dd135726.aspx
func (b *Blob) PutBlockWithLength(blockID string, size uint64, blob io.Reader, extraHeaders map[string]string) error {
	query := url.Values{
		"comp":    {"block"},
		"blockid": {blockID},
	}
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), query)
	extraHeaders = b.Container.bsc.client.protectUserAgent(extraHeaders)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeBlock)
	headers["Content-Length"] = strconv.FormatUint(size, 10)

	for k, v := range extraHeaders {
		headers[k] = v
	}

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, blob, b.Container.bsc.auth)
	if err != nil {
		return err
	}

	readAndCloseBody(resp.body)
	return checkRespCode(resp.statusCode, []int{http.StatusCreated})
}

// PutBlockList saves list of blocks to the specified block blob.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179467.aspx
func (b *Blob) PutBlockList(blocks []Block) error {
	blockListXML := prepareBlockListRequest(blocks)

	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), url.Values{"comp": {"blocklist"}})
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["Content-Length"] = strconv.Itoa(len(blockListXML))

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, strings.NewReader(blockListXML), b.Container.bsc.auth)
	if err != nil {
		return err
	}
	readAndCloseBody(resp.body)
	return checkRespCode(resp.statusCode, []int{http.StatusCreated})
}

// GetBlockList retrieves list of blocks in the specified block blob.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179400.aspx
func (b *Blob) GetBlockList(blockType BlockListType) (BlockListResponse, error) {
	params := url.Values{
		"comp":          {"blocklist"},
		"blocklisttype": {string(blockType)},
	}
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), params)
	headers := b.Container.bsc.client.getStandardHeaders()

	var out BlockListResponse
	resp, err := b.Container.bsc.client.exec(http.MethodGet, uri, headers, nil, b.Container.bsc.auth)
	if err != nil {
		return out, err
	}
	defer resp.body.Close()

	err = xmlUnmarshal(resp.body, &out)
	return out, err
}
