//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

const (
	TokenScope = "https://storage.azure.com/.default"
)

const (
	HeaderAuthorization     = "Authorization"
	HeaderXmsDate           = "x-ms-date"
	HeaderContentLength     = "Content-Length"
	HeaderContentEncoding   = "Content-Encoding"
	HeaderContentLanguage   = "Content-Language"
	HeaderContentType       = "Content-Type"
	HeaderContentMD5        = "Content-MD5"
	HeaderIfModifiedSince   = "If-Modified-Since"
	HeaderIfMatch           = "If-Match"
	HeaderIfNoneMatch       = "If-None-Match"
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	HeaderRange             = "Range"
	HeaderXmsVersion        = "x-ms-version"
	HeaderXmsRequestID      = "x-ms-request-id"
)

const crc64Polynomial uint64 = 0x9A6C9329AC4BC9B5

var CRC64Table = crc64.MakeTable(crc64Polynomial)

// CopyOptions returns a zero-value T if opts is nil.
// If opts is not nil, a copy is made and its address returned.
func CopyOptions[T any](opts *T) *T {
	if opts == nil {
		return new(T)
	}
	cp := *opts
	return &cp
}

var errConnectionString = errors.New("connection string is either blank or malformed. The expected connection string " +
	"should contain key value pairs separated by semicolons. For example 'DefaultEndpointsProtocol=https;AccountName=<accountName>;" +
	"AccountKey=<accountKey>;EndpointSuffix=core.windows.net'")

type ParsedConnectionString struct {
	ServiceURL  string
	AccountName string
	AccountKey  string
}

func GetURLs(url string) (string, string) {
	blobURL := strings.Replace(url, ".dfs.", ".blob.", 1)
	dfsURL := strings.Replace(url, ".blob.", ".dfs.", 1)
	return blobURL, dfsURL
}

func ParseConnectionString(connectionString string) (ParsedConnectionString, error) {
	const (
		defaultScheme = "https"
		defaultSuffix = "core.windows.net"
	)

	connStrMap := make(map[string]string)
	connectionString = strings.TrimRight(connectionString, ";")

	splitString := strings.Split(connectionString, ";")
	if len(splitString) == 0 {
		return ParsedConnectionString{}, errConnectionString
	}
	for _, stringPart := range splitString {
		parts := strings.SplitN(stringPart, "=", 2)
		if len(parts) != 2 {
			return ParsedConnectionString{}, errConnectionString
		}
		connStrMap[parts[0]] = parts[1]
	}

	protocol, ok := connStrMap["DefaultEndpointsProtocol"]
	if !ok {
		protocol = defaultScheme
	}

	suffix, ok := connStrMap["EndpointSuffix"]
	if !ok {
		suffix = defaultSuffix
	}

	blobEndpoint, has_blobEndpoint := connStrMap["BlobEndpoint"]
	accountName, has_accountName := connStrMap["AccountName"]

	var serviceURL string
	if has_blobEndpoint {
		serviceURL = blobEndpoint
	} else if has_accountName {
		serviceURL = fmt.Sprintf("%v://%v.blob.%v", protocol, accountName, suffix)
	} else {
		return ParsedConnectionString{}, errors.New("connection string needs either AccountName or BlobEndpoint")
	}

	if !strings.HasSuffix(serviceURL, "/") {
		// add a trailing slash to be consistent with the portal
		serviceURL += "/"
	}

	accountKey, has_accountKey := connStrMap["AccountKey"]
	sharedAccessSignature, has_sharedAccessSignature := connStrMap["SharedAccessSignature"]

	if has_accountName && has_accountKey {
		return ParsedConnectionString{
			ServiceURL:  serviceURL,
			AccountName: accountName,
			AccountKey:  accountKey,
		}, nil
	} else if has_sharedAccessSignature {
		return ParsedConnectionString{
			ServiceURL: fmt.Sprintf("%v?%v", serviceURL, sharedAccessSignature),
		}, nil
	} else {
		return ParsedConnectionString{}, errors.New("connection string needs either AccountKey or SharedAccessSignature")
	}

}

func SerializeBlobTagsToStrPtr(tagsMap map[string]string) *string {
	if len(tagsMap) == 0 {
		return nil
	}
	tags := make([]string, 0)
	for key, val := range tagsMap {
		tags = append(tags, url.QueryEscape(key)+"="+url.QueryEscape(val))
	}
	blobTagsString := strings.Join(tags, "&")
	return &blobTagsString
}

func ValidateSeekableStreamAt0AndGetCount(body io.ReadSeeker) (int64, error) {
	if body == nil { // nil body's are "logically" seekable to 0 and are 0 bytes long
		return 0, nil
	}

	err := validateSeekableStreamAt0(body)
	if err != nil {
		return 0, err
	}

	count, err := body.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, errors.New("body stream must be seekable")
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// return an error if body is not a valid seekable stream at 0
func validateSeekableStreamAt0(body io.ReadSeeker) error {
	if body == nil { // nil body's are "logically" seekable to 0
		return nil
	}
	if pos, err := body.Seek(0, io.SeekCurrent); pos != 0 || err != nil {
		// Help detect programmer error
		if err != nil {
			return errors.New("body stream must be seekable")
		}
		return errors.New("body stream must be set to position 0")
	}
	return nil
}

func RangeToString(offset, count int64) string {
	return "bytes=" + strconv.FormatInt(offset, 10) + "-" + strconv.FormatInt(offset+count-1, 10)
}

type nopCloser struct {
	io.ReadSeeker
}

func (n nopCloser) Close() error {
	return nil
}

// NopCloser returns a ReadSeekCloser with a no-op close method wrapping the provided io.ReadSeeker.
func NopCloser(rs io.ReadSeeker) io.ReadSeekCloser {
	return nopCloser{rs}
}

func GenerateLeaseID(leaseID *string) (*string, error) {
	if leaseID == nil {
		generatedUuid, err := uuid.New()
		if err != nil {
			return nil, err
		}
		leaseID = to.Ptr(generatedUuid.String())
	}
	return leaseID, nil
}

func GetClientOptions[T any](o *T) *T {
	if o == nil {
		return new(T)
	}
	return o
}

// IsIPEndpointStyle checkes if URL's host is IP, in this case the storage account endpoint will be composed as:
// http(s)://IP(:port)/storageaccount/container/...
// As url's Host property, host could be both host or host:port
func IsIPEndpointStyle(host string) bool {
	if host == "" {
		return false
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	// For IPv6, there could be case where SplitHostPort fails for cannot finding port.
	// In this case, eliminate the '[' and ']' in the URL.
	// For details about IPv6 URL, please refer to https://tools.ietf.org/html/rfc2732
	if host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1 : len(host)-1]
	}
	return net.ParseIP(host) != nil
}

// EscapeSplitPaths is utility function to escape the individual strings by eliminating "/" in the path
func EscapeSplitPaths(filePath string) string {
	names := strings.Split(filePath, "/")
	path := make([]string, len(names))
	for i, name := range names {
		path[i] = url.PathEscape(name)
	}
	escapedPathUrl := strings.Join(path, "/")
	return escapedPathUrl
}
