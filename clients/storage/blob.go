package storage

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/url"
)

type BlobStorageClient struct {
	client StorageClient
}

// TODO (ahmetalpbalkan) use
type Container struct {
	Name       string              `xml:"Name"`
	Properties ContainerProperties `xml:"Properties"`
	Metadata   map[string]string   `xml:"Metadata"`
}

// TODO (ahmetalpbalkan) use
type ContainerProperties struct {
	LastModified  string `xml:"Last-Modified"`
	Etag          string `xml:"Etag"`
	LeaseStatus   string `xml:"LeaseStatus"`
	LeaseState    string `xml:"LeaseState"`
	LeaseDuration string `xml:"LeaseDuration"`
}

type ContainerListResponse struct {
	XMLName    xml.Name      `xml:"EnumerationResults"`
	Xmlns      string        `xml:"xmlns,attr"`
	Prefix     string        `xml:"Prefix"`
	Marker     string        `xml:"Marker"`
	NextMarker string        `xml:"NextMarker"`
	MaxResults int64         `xml:"MaxResults"`
	Containers ContainerList `xml:"Containers"`
}

type ContainerList struct {
	XMLName    xml.Name    `xml:"Containers"`
	Containers []Container `xml:"Container"`
}

type BlobType string

const (
	BlobTypeBlock BlobType = "BlockBlob"
	BlobTypePage  BlobType = "PageBlob"
)

type ContainerAccessType string

const (
	ContainerAccessTypePrivate   ContainerAccessType = ""
	ContainerAccessTypeBlob      ContainerAccessType = "blob"
	ContainerAccessTypeContainer ContainerAccessType = "container"
)

func (b BlobStorageClient) ListContainers() (ContainerListResponse, error) {
	// TODO (ahmetb) pagination
	verb := "GET"
	uri := b.client.getEndpoint(blobServiceName, "", url.Values{"comp": {"list"}})
	headers := b.client.getStandardHeaders()

	var out ContainerListResponse
	resp, err := b.client.exec(verb, uri, headers, nil)
	if err != nil {
		return out, err
	}

	fmt.Println(string(resp))
	err = xml.Unmarshal(resp, &out)
	return out, err
}

func (b BlobStorageClient) ContainerExists(name string) ([]byte, error) {
	verb := "GET"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) CreateContainer(name string, access ContainerAccessType) ([]byte, error) {
	verb := "PUT"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	headers["Content-Length"] = "0"
	if access != "" {
		headers["x-ms-blob-public-access"] = string(access)
	}
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) DeleteContainer(name string) ([]byte, error) {
	verb := "DELETE"
	uri := b.client.getEndpoint(blobServiceName, name, url.Values{"restype": {"container"}})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}

func (b BlobStorageClient) PutBlob(name, container string, blobType BlobType, blob []byte) ([]byte, error) {
	verb := "PUT"
	path := fmt.Sprintf("%s/%s", name, container)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(blobType)
	headers["Content-Length"] = fmt.Sprintf("%v", len(blob))
	return b.client.exec(verb, uri, headers, bytes.NewReader(blob))
}

func (b BlobStorageClient) DeleteBlob(name, container string) ([]byte, error) {
	verb := "DELETE"
	path := fmt.Sprintf("%s/%s", name, container)
	uri := b.client.getEndpoint(blobServiceName, path, url.Values{})

	headers := b.client.getStandardHeaders()
	return b.client.exec(verb, uri, headers, nil)
}
