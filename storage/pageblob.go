package storage

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// GetPageRangesResponse contains the response fields from
// Get Page Ranges call.
//
// See https://msdn.microsoft.com/en-us/library/azure/ee691973.aspx
type GetPageRangesResponse struct {
	XMLName  xml.Name    `xml:"PageList"`
	PageList []PageRange `xml:"PageRange"`
}

// PageRange contains information about a page of a page blob from
// Get Pages Range call.
//
// See https://msdn.microsoft.com/en-us/library/azure/ee691973.aspx
type PageRange struct {
	Start int64 `xml:"Start"`
	End   int64 `xml:"End"`
}

var (
	errBlobCopyAborted    = errors.New("storage: blob copy is aborted")
	errBlobCopyIDMismatch = errors.New("storage: blob copy id is a mismatch")
)

// PageWriteType defines the type updates that are going to be
// done on the page blob.
type PageWriteType string

// Types of operations on page blobs
const (
	PageWriteTypeUpdate PageWriteType = "update"
	PageWriteTypeClear  PageWriteType = "clear"
)

// PutPage writes a range of pages to a page blob or clears the given range.
// In case of 'clear' writes, given chunk is discarded. Ranges must be aligned
// with 512-byte boundaries and chunk must be of size multiplies by 512.
//
// See https://msdn.microsoft.com/en-us/library/ee691975.aspx
func (b *Blob) PutPage(startByte, endByte int64, writeType PageWriteType, chunk []byte, extraHeaders map[string]string) error {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), url.Values{"comp": {"page"}})
	extraHeaders = b.Container.bsc.client.protectUserAgent(extraHeaders)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypePage)
	headers["x-ms-page-write"] = string(writeType)
	headers["x-ms-range"] = fmt.Sprintf("bytes=%v-%v", startByte, endByte)
	for k, v := range extraHeaders {
		headers[k] = v
	}
	var contentLength int64
	var data io.Reader
	if writeType == PageWriteTypeClear {
		contentLength = 0
		data = bytes.NewReader([]byte{})
	} else {
		contentLength = int64(len(chunk))
		data = bytes.NewReader(chunk)
	}
	headers["Content-Length"] = strconv.FormatInt(contentLength, 10)

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, data, b.Container.bsc.auth)
	if err != nil {
		return err
	}
	readAndCloseBody(resp.body)
	return checkRespCode(resp.statusCode, []int{http.StatusCreated})
}

// GetPageRanges returns the list of valid page ranges for a page blob.
//
// See https://msdn.microsoft.com/en-us/library/azure/ee691973.aspx
func (b *Blob) GetPageRanges() (GetPageRangesResponse, error) {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), url.Values{"comp": {"pagelist"}})
	headers := b.Container.bsc.client.getStandardHeaders()

	var out GetPageRangesResponse
	resp, err := b.Container.bsc.client.exec(http.MethodGet, uri, headers, nil, b.Container.bsc.auth)
	if err != nil {
		return out, err
	}
	defer resp.body.Close()

	if err = checkRespCode(resp.statusCode, []int{http.StatusOK}); err != nil {
		return out, err
	}
	err = xmlUnmarshal(resp.body, &out)
	return out, err
}

// PutPageBlob initializes an empty page blob with specified name and maximum
// size in bytes (size must be aligned to a 512-byte boundary). A page blob must
// be created using this method before writing pages.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179451.aspx
func (b *Blob) PutPageBlob(size int64, extraHeaders map[string]string) error {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), nil)
	extraHeaders = b.Container.bsc.client.protectUserAgent(extraHeaders)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypePage)
	headers["x-ms-blob-content-length"] = strconv.FormatInt(size, 10)

	for k, v := range extraHeaders {
		headers[k] = v
	}

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, nil, b.Container.bsc.auth)
	if err != nil {
		return err
	}
	readAndCloseBody(resp.body)
	return checkRespCode(resp.statusCode, []int{http.StatusCreated})
}
