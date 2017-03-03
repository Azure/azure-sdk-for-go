package storage

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
)

// PutAppendBlob initializes an empty append blob with specified name. An
// append blob must be created using this method before appending blocks.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd179451.aspx
func (b *Blob) PutAppendBlob(extraHeaders map[string]string) error {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), nil)
	extraHeaders = b.Container.bsc.client.protectUserAgent(extraHeaders)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeAppend)

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

// AppendBlock appends a block to an append blob.
//
// See https://msdn.microsoft.com/en-us/library/azure/mt427365.aspx
func (b *Blob) AppendBlock(chunk []byte, extraHeaders map[string]string) error {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), url.Values{"comp": {"appendblock"}})
	extraHeaders = b.Container.bsc.client.protectUserAgent(extraHeaders)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeAppend)
	headers["Content-Length"] = strconv.Itoa(len(chunk))

	for k, v := range extraHeaders {
		headers[k] = v
	}

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, bytes.NewReader(chunk), b.Container.bsc.auth)
	if err != nil {
		return err
	}
	readAndCloseBody(resp.body)
	return checkRespCode(resp.statusCode, []int{http.StatusCreated})
}
