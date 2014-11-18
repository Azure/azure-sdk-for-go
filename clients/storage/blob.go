package storage

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type BlobStorageClient struct {
	client StorageClient
}

// TODO (ahmetalpbalkan) use
type Container struct {
	Name       string
	Properties ContainerProperties
	Metadata   map[string]string
}

// TODO (ahmetalpbalkan) use
type ContainerProperties struct {
	LastModified  time.Time
	Etag          string
	LeaseStatus   string
	LeaseState    string
	LeaseDuration string
}

// TODO (ahmetalpbalkan) use
type ContainerListResponse struct {
	Prefix     string
	Marker     string
	NextMarker string
	MaxResults int64
	Containers []Container
}

const (
	BlobTypeBlock = "BlockBlob"
	BlobTypePage  = "PageBlob"
)

func (b BlobStorageClient) ListContainers() (*http.Response, error) {
	// TODO (ahmetb) pagination

	verb := "GET"
	uri := b.client.getEndpoint(blobServiceName, "", url.Values{"comp": {"list"}})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) GetContainer(name string) (*http.Response, error) {
	verb := "GET"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) CreateContainer(name string) (*http.Response, error) {
	verb := "PUT"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	headers["Content-Length"] = "0"
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) DeleteContainer(name string) (*http.Response, error) {
	verb := "DELETE"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) PutBlob(name, container, blobType string, blob []byte) (*http.Response, error) {
	verb := "PUT"
	path := fmt.Sprintf("%s/%s", name, container)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	headers["x-ms-blob-type"] = blobType
	headers["Content-Length"] = fmt.Sprintf("%v", len(blob))
	return b.client.exec(verb, uri, headers, bytes.NewReader(blob))
}

func (b BlobStorageClient) DeleteBlob(name, container string) (*http.Response, error) {
	verb := "DELETE"
	path := fmt.Sprintf("%s/%s", name, container)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}
