package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type BlobStorageClient struct {
	client StorageClient
}

type Container struct {
	Name       string              `xml:"Name"`
	Properties ContainerProperties `xml:"Properties"`
	// TODO (ahmetalpbalkan) Metadata
}

type Blob struct {
	Name       string         `xml:"Name"`
	Properties BlobProperties `xml:"Properties"`
	// TODO (ahmetalpbalkan) Metadata
}

type ContainerProperties struct {
	LastModified  string `xml:"Last-Modified"`
	Etag          string `xml:"Etag"`
	LeaseStatus   string `xml:"LeaseStatus"`
	LeaseState    string `xml:"LeaseState"`
	LeaseDuration string `xml:"LeaseDuration"`
	// TODO (ahmetalpbalkan) remaining fields
}

type BlobProperties struct {
	LastModified          string `xml:"Last-Modified"`
	Etag                  string `xml:"Etag"`
	ContentMD5            string `xml:"Content-MD5"`
	ContentLength         int64  `xml:"Content-Length"`
	ContentType           string `xml:"Content-Type"`
	ContentEncoding       string `xml:"Content-Encoding"`
	SequenceNumber        int64  `xml:"x-ms-blob-sequence-number"`
	CopyId                string `xml:"CopyId"`
	CopyStatus            string `xml:"CopyStatus"`
	CopySource            string `xml:"CopySource"`
	CopyProgress          string `xml:"CopyProgress"`
	CopyCompletionTime    string `xml:"CopyCompletionTime"`
	CopyStatusDescription string `xml:"CopyStatusDescription"`
}

type BlobListResponse struct {
	XMLName    xml.Name `xml:"EnumerationResults"`
	Xmlns      string   `xml:"xmlns,attr"`
	Prefix     string   `xml:"Prefix"`
	Marker     string   `xml:"Marker"`
	NextMarker string   `xml:"NextMarker"`
	MaxResults int64    `xml:"MaxResults"`
	Blobs      []Blob   `xml:"Blobs>Blob"`
}

type ContainerListResponse struct {
	XMLName    xml.Name    `xml:"EnumerationResults"`
	Xmlns      string      `xml:"xmlns,attr"`
	Prefix     string      `xml:"Prefix"`
	Marker     string      `xml:"Marker"`
	NextMarker string      `xml:"NextMarker"`
	MaxResults int64       `xml:"MaxResults"`
	Containers []Container `xml:"Containers>Container"`
}

type ListContainersParameters struct {
	Prefix     string
	Marker     string
	Include    string
	MaxResults uint
	Timeout    uint
}

func (p ListContainersParameters) GetParameters() url.Values {
	out := url.Values{}

	if p.Prefix != "" {
		out.Set("prefix", p.Prefix)
	}
	if p.Marker != "" {
		out.Set("marker", p.Marker)
	}
	if p.Include != "" {
		out.Set("include", p.Include)
	}
	if p.MaxResults != 0 {
		out.Set("maxresults", fmt.Sprintf("%v", p.MaxResults))
	}
	if p.Timeout != 0 {
		out.Set("timeout", fmt.Sprintf("%v", p.Timeout))
	}

	return out
}

type ListBlobsParameters struct {
	Prefix     string
	Delimiter  string
	Marker     string
	Include    string
	MaxResults uint
	Timeout    uint
}

func (p ListBlobsParameters) GetParameters() url.Values {
	out := url.Values{}

	if p.Prefix != "" {
		out.Set("prefix", p.Prefix)
	}
	if p.Delimiter != "" {
		out.Set("delimiter", p.Delimiter)
	}
	if p.Marker != "" {
		out.Set("marker", p.Marker)
	}
	if p.Include != "" {
		out.Set("include", p.Include)
	}
	if p.MaxResults != 0 {
		out.Set("maxresults", fmt.Sprintf("%v", p.MaxResults))
	}
	if p.Timeout != 0 {
		out.Set("timeout", fmt.Sprintf("%v", p.Timeout))
	}

	return out
}

type BlobType string

const (
	BlobTypeBlock BlobType = "BlockBlob"
	BlobTypePage  BlobType = "PageBlob"
)

const (
	blobCopyStatusPending = "pending"
	blobCopyStatusSuccess = "success"
	blobCopyStatusAborted = "aborted"
	blobCopyStatusFailed  = "failed"
)

type BlockListType string

const (
	BlockListTypeAll         BlockListType = "all"
	BlockListTypeCommitted   BlockListType = "committed"
	BlockListTypeUncommitted BlockListType = "uncommitted"
)

type ContainerAccessType string

const (
	ContainerAccessTypePrivate   ContainerAccessType = ""
	ContainerAccessTypeBlob      ContainerAccessType = "blob"
	ContainerAccessTypeContainer ContainerAccessType = "container"
)

const MaxBlobBlockSize = 64 * 1024 * 1024

type BlockStatus string

const (
	blockStatusUncommitted BlockStatus = "Uncommitted"
	blockStatusCommitted   BlockStatus = "Committed"
	blockStatusLatest      BlockStatus = "Latest"
)

type Block struct {
	Id     string
	Status BlockStatus
}

type BlockListResponse struct {
	XMLName           xml.Name        `xml:"BlockList"`
	CommittedBlocks   []BlockResponse `xml:"CommittedBlocks>Block"`
	UncommittedBlocks []BlockResponse `xml:"UncommittedBlocks>Block"`
}

type BlockResponse struct {
	Name string `xml:"Name"`
	Size int64  `xml:"Size"`
}

var (
	ErrNotCreated  = errors.New("storage: operation has returned a successful error code other than 201 Created.")
	ErrNotAccepted = errors.New("storage: operation has returned a successful error code other than 202 Accepted.")

	errBlobCopyAborted    = errors.New("storage: blob copy is aborted")
	errBlobCopyIdMismatch = errors.New("storage: blob copy id is a mismatch")
)

const errUnexpectedStatus = "storage: was expecting status code: %d, got: %d"

func (b BlobStorageClient) ListContainers(params ListContainersParameters) (ContainerListResponse, error) {
	q := mergeParams(params.GetParameters(), url.Values{"comp": {"list"}})
	uri := b.client.getEndpoint(blobServiceName, "", q)
	headers := b.client.getStandardHeaders()

	var out ContainerListResponse
	resp, err := b.client.exec("GET", uri, headers, nil)
	if err != nil {
		return out, err
	}

	err = xmlUnmarshal(resp.body, &out)
	return out, err
}

func (b BlobStorageClient) CreateContainer(name string, access ContainerAccessType) error {
	resp, err := b.createContainer(name, access)
	if err != nil {
		return err
	}
	if resp.statusCode != http.StatusCreated {
		return ErrNotCreated
	}
	return nil
}

func (b BlobStorageClient) CreateContainerIfNotExists(name string, access ContainerAccessType) (bool, error) {
	resp, err := b.createContainer(name, access)
	if resp != nil && (resp.statusCode == http.StatusCreated || resp.statusCode == http.StatusConflict) {
		return resp.statusCode == http.StatusCreated, nil
	}
	return false, err
}

func (b BlobStorageClient) createContainer(name string, access ContainerAccessType) (*storageResponse, error) {
	verb := "PUT"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	headers["Content-Length"] = "0"
	if access != "" {
		headers["x-ms-blob-public-access"] = string(access)
	}
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) ContainerExists(container string) (bool, error) {
	verb := "HEAD"
	path := fmt.Sprintf("%s", container)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{"restype": {"container"}})
	headers := b.client.getStandardHeaders()

	resp, err := b.client.exec(verb, uri, headers, nil)
	if resp != nil && (resp.statusCode == http.StatusOK || resp.statusCode == http.StatusNotFound) {
		return resp.statusCode == http.StatusOK, nil
	}
	return false, err
}

func (b BlobStorageClient) DeleteContainer(name string) error {
	resp, err := b.deleteContainer(name)
	if err != nil {
		return err
	}
	if resp.statusCode != http.StatusAccepted {
		return ErrNotAccepted
	}
	return nil
}

func (b BlobStorageClient) DeleteContainerIfExists(container string) (bool, error) {
	resp, err := b.deleteContainer(container)
	if resp != nil && (resp.statusCode == http.StatusAccepted || resp.statusCode == http.StatusNotFound) {
		return resp.statusCode == http.StatusAccepted, nil
	}
	return false, err
}

func (b BlobStorageClient) deleteContainer(name string) (*storageResponse, error) {
	verb := "DELETE"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) ListBlobs(container string, params ListBlobsParameters) (BlobListResponse, error) {
	q := mergeParams(params.GetParameters(), url.Values{
		"restype": {"container"},
		"comp":    {"list"}})
	uri := b.client.getEndpoint(blobServiceName, container, q)
	headers := b.client.getStandardHeaders()

	var out BlobListResponse
	resp, err := b.client.exec("GET", uri, headers, nil)
	if err != nil {
		return out, err
	}

	err = xmlUnmarshal(resp.body, &out)
	return out, err
}

func (b BlobStorageClient) BlobExists(container, name string) (bool, error) {
	verb := "HEAD"
	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	resp, err := b.client.exec(verb, uri, headers, nil)
	if resp != nil && (resp.statusCode == http.StatusOK || resp.statusCode == http.StatusNotFound) {
		return resp.statusCode == http.StatusOK, nil
	}
	return false, err
}

func (b BlobStorageClient) GetBlobUrl(container, name string) string {
	if container == "" {
		container = "$root"
	}
	path := fmt.Sprintf("%s/%s", container, name)
	return b.client.getEndpoint(blobServiceName, path, url.Values{})
}

func (b BlobStorageClient) GetBlob(container, name string) (io.ReadCloser, error) {
	verb := "GET"
	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	resp, err := b.client.exec(verb, uri, headers, nil)
	if err != nil {
		return nil, err
	}

	if resp.statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatus, http.StatusOK, resp.statusCode)
	}
	return resp.body, nil
}

func (b BlobStorageClient) GetBlobProperties(container, name string) (*BlobProperties, error) {
	verb := "HEAD"
	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	resp, err := b.client.exec(verb, uri, headers, nil)
	if err != nil {
		return nil, err
	}

	if resp.statusCode != http.StatusOK {
		return nil, fmt.Errorf(errUnexpectedStatus, http.StatusOK, resp.statusCode)
	}

	var contentLength, sequenceNum int64
	contentLengthStr := resp.headers.Get("Content-Length")
	if contentLengthStr != "" {
		contentLength, err = strconv.ParseInt(contentLengthStr, 0, 64)
		if err != nil {
			return nil, err
		}
	}

	sequenceNumStr := resp.headers.Get("Content-Length")
	if sequenceNumStr != "" {
		sequenceNum, err = strconv.ParseInt(sequenceNumStr, 0, 64)
		if err != nil {
			return nil, err
		}
	}

	return &BlobProperties{
		LastModified:          resp.headers.Get("Last-Modified"),
		Etag:                  resp.headers.Get("Etag"),
		ContentMD5:            resp.headers.Get("Content-MD5"),
		ContentLength:         contentLength,
		ContentEncoding:       resp.headers.Get("Content-Encodng"),
		SequenceNumber:        sequenceNum,
		CopyCompletionTime:    resp.headers.Get("x-ms-copy-completion-time"),
		CopyStatusDescription: resp.headers.Get("x-ms-copy-status-description"),
		CopyId:                resp.headers.Get("x-ms-copy-id"),
		CopyProgress:          resp.headers.Get("x-ms-copy-progress"),
		CopySource:            resp.headers.Get("x-ms-copy-source"),
		CopyStatus:            resp.headers.Get("x-ms-copy-status"),
	}, nil
}

func (b BlobStorageClient) PutBlockBlob(container, name string, blob io.Reader) error { // TODO (ahmetalpbalkan) consider ReadCloser and closing
	return b.putBlockBlob(container, name, blob, MaxBlobBlockSize)
}

func (b BlobStorageClient) putBlockBlob(container, name string, blob io.Reader, chunkSize int) error {
	if chunkSize <= 0 || chunkSize > MaxBlobBlockSize {
		chunkSize = MaxBlobBlockSize
	}

	chunk := make([]byte, chunkSize)
	n, err := blob.Read(chunk)
	if err != nil && err != io.EOF {
		return err
	}

	if err == io.EOF {
		// Fits into one block
		return b.putSingleBlockBlob(container, name, chunk[:n])
	} else {
		// Does not fit into one block. Upload block by block then commit the block list
		blockList := []Block{}
		blockNum := 0

		// Put blocks
		for {
			id := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", blockNum)))
			data := chunk[:n]
			err = b.PutBlock(container, name, id, data)
			if err != nil {
				return err
			}
			blockList = append(blockList, Block{id, blockStatusLatest})

			// Read next block
			n, err = blob.Read(chunk)
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				break
			}

			blockNum++
		}

		// Commit block list
		return b.PutBlockList(container, name, blockList)
	}
}

func (b BlobStorageClient) putSingleBlockBlob(container, name string, chunk []byte) error {
	if len(chunk) > MaxBlobBlockSize {
		return fmt.Errorf("storage: provided chunk (%d bytes) cannot fit into single-block blob (max %d bytes)", len(chunk), MaxBlobBlockSize)
	}

	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})
	headers := b.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeBlock)
	headers["Content-Length"] = fmt.Sprintf("%v", len(chunk))

	resp, err := b.client.exec("PUT", uri, headers, bytes.NewReader(chunk))
	if err != nil {
		return err
	}
	if resp.statusCode != http.StatusCreated {
		return ErrNotCreated
	}

	return nil
}

func (b BlobStorageClient) PutBlock(container, name, blockId string, chunk []byte) error {
	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{"comp": {"block"}, "blockid": {blockId}})
	headers := b.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeBlock)
	headers["Content-Length"] = fmt.Sprintf("%v", len(chunk))

	resp, err := b.client.exec("PUT", uri, headers, bytes.NewReader(chunk))
	if err != nil {
		return err
	}
	if resp.statusCode != http.StatusCreated {
		return ErrNotCreated
	}

	return nil
}

func (b BlobStorageClient) PutBlockList(container, name string, blocks []Block) error {
	blockListXml := prepareBlockListRequest(blocks)

	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{"comp": {"blocklist"}})
	headers := b.client.getStandardHeaders()
	headers["Content-Length"] = fmt.Sprintf("%v", len(blockListXml))

	resp, err := b.client.exec("PUT", uri, headers, strings.NewReader(blockListXml))
	if err != nil {
		return err
	}
	if resp.statusCode != http.StatusCreated {
		return ErrNotCreated
	}
	return nil
}

func (b BlobStorageClient) GetBlockList(container, name string, blockType BlockListType) (BlockListResponse, error) {
	params := url.Values{"comp": {"blocklist"}, "blocklisttype": {string(blockType)}}
	uri := b.client.getEndpoint(blobServiceName, fmt.Sprintf("%s/%s", container, name), params)
	headers := b.client.getStandardHeaders()

	var out BlockListResponse
	resp, err := b.client.exec("GET", uri, headers, nil)
	if err != nil {
		return out, err
	}

	err = xmlUnmarshal(resp.body, &out)
	return out, err
}

func (b BlobStorageClient) CopyBlob(container, name, sourceBlob string) error {
	copyId, err := b.startBlobCopy(container, name, sourceBlob)
	if err != nil {
		return err
	}

	return b.waitForBlobCopy(container, name, copyId)
}

func (b BlobStorageClient) startBlobCopy(container, name, sourceBlob string) (string, error) {
	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	headers["Content-Length"] = "0"
	headers["x-ms-copy-source"] = sourceBlob

	resp, err := b.client.exec("PUT", uri, headers, nil)
	if err != nil {
		return "", err
	}
	if resp.statusCode != http.StatusAccepted && resp.statusCode != http.StatusCreated {
		return "", fmt.Errorf(errUnexpectedStatus, []int{http.StatusAccepted, http.StatusCreated}, resp.statusCode)
	}

	copyId := resp.headers.Get("x-ms-copy-id")
	if copyId == "" {
		return "", errors.New("Got empty copy id header")
	}
	return copyId, nil
}

func (b BlobStorageClient) waitForBlobCopy(container, name, copyId string) error {
	for {
		props, err := b.GetBlobProperties(container, name)
		if err != nil {
			return err
		}

		if props.CopyId != copyId {
			return errBlobCopyIdMismatch
		}

		switch props.CopyStatus {
		case blobCopyStatusSuccess:
			return nil
		case blobCopyStatusPending:
			continue
		case blobCopyStatusAborted:
			return errBlobCopyAborted
		case blobCopyStatusFailed:
			return fmt.Errorf("storage: blob copy failed. Id=%s Description=%s", props.CopyId, props.CopyStatusDescription)
		default:
			return fmt.Errorf("storage: unhandled blob copy status: '%s'", props.CopyStatus)
		}
	}
}

func (b BlobStorageClient) DeleteBlob(container, name string) error {
	resp, err := b.deleteBlob(container, name)
	if err != nil {
		return err
	}
	if resp.statusCode != http.StatusAccepted {
		return ErrNotAccepted
	}
	return nil
}

func (b BlobStorageClient) GetBlobSASURI(container, name string, expiry time.Time, permissions string) (string, error) {
	var (
		signedPermissions = permissions
		blobUrl           = b.GetBlobUrl(container, name)
	)
	canonicalizedResource, err := b.client.buildCanonicalizedResource(blobUrl)
	if err != nil {
		return "", err
	}
	signedExpiry := expiry.Format(time.RFC3339)
	signedResource := "b"

	stringToSign, err := blobSASStringToSign(b.client.apiVersion, canonicalizedResource, signedExpiry, signedPermissions)
	if err != nil {
		return "", err
	}

	sig := b.client.computeHmac256(stringToSign)
	sasParams := url.Values{
		"sv":  {b.client.apiVersion},
		"se":  {signedExpiry},
		"sr":  {signedResource},
		"sp":  {signedPermissions},
		"sig": {sig},
	}

	sasUrl, err := url.Parse(blobUrl)
	if err != nil {
		return "", err
	}
	sasUrl.RawQuery = sasParams.Encode()
	return sasUrl.String(), nil
}

func blobSASStringToSign(signedVersion, canonicalizedResource, signedExpiry, signedPermissions string) (string, error) {
	var signedStart, signedIdentifier, rscc, rscd, rsce, rscl, rsct string

	// reference: http://msdn.microsoft.com/en-us/library/azure/dn140255.aspx
	if signedVersion >= "2013-08-15" {
		return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", signedPermissions, signedStart, signedExpiry, canonicalizedResource, signedIdentifier, signedVersion, rscc, rscd, rsce, rscl, rsct), nil
	} else {
		return "", errors.New("storage: not implemented SAS for versions earlier than 2013-08-15")
	}
}

func (b BlobStorageClient) DeleteBlobIfExists(container, name string) (bool, error) {
	resp, err := b.deleteBlob(container, name)
	if resp != nil && (resp.statusCode == http.StatusAccepted || resp.statusCode == http.StatusNotFound) {
		return resp.statusCode == http.StatusAccepted, nil
	}
	return false, err
}

func (b BlobStorageClient) deleteBlob(container, name string) (*storageResponse, error) {
	verb := "DELETE"
	path := fmt.Sprintf("%s/%s", container, name)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}
