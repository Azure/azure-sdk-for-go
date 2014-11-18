package storage

import (
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

func (b BlobStorageClient) ListContainers() (*http.Response, error) {
	verb := "GET"
	uri := b.client.getEndpoint(blobServiceName, "", url.Values{"comp": {"list"}})

	headers := b.client.getStandardHeaders()
	canonicalizedHeaders := b.client.buildCanonicalizedHeader(headers)
	canonicalizedResource, err := b.client.buildCanonicalizedResource(uri)

	if err != nil {
		return nil, err
	}

	contentEncoding := ""
	contentLanguage := ""
	contentLength := ""
	contentMD5 := ""
	contentType := ""
	date := ""
	ifModifiedSince := ""
	ifMatch := ""
	ifNoneMatch := ""
	ifUnmodifiedSince := ""
	Range := ""

	canonicalizedString := b.client.buildCanonicalizedString(verb, contentEncoding, contentLanguage, contentLength, contentMD5, contentType,
		date, ifModifiedSince, ifMatch, ifNoneMatch, ifUnmodifiedSince, Range, canonicalizedHeaders, canonicalizedResource)
	authHeader, err := b.client.createAuthorizationHeader(canonicalizedString)
	if err != nil {
		return nil, err
	}

	headers["Authorization"] = authHeader
	resp, err := b.client.exec(verb, uri, headers, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
