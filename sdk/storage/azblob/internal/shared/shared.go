// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

const (
	TokenScope = "https://storage.azure.com/.default"
)

const (
	HeaderAuthorization      = "Authorization"
	HeaderXmsDate            = "x-ms-date"
	HeaderContentLength      = "Content-Length"
	HeaderContentEncoding    = "Content-Encoding"
	HeaderContentLanguage    = "Content-Language"
	HeaderContentType        = "Content-Type"
	HeaderContentMD5         = "Content-MD5"
	HeaderIfModifiedSince    = "If-Modified-Since"
	HeaderIfMatch            = "If-Match"
	HeaderIfNoneMatch        = "If-None-Match"
	HeaderIfUnmodifiedSince  = "If-Unmodified-Since"
	HeaderRange              = "Range"
	HeaderXmsVersion         = "x-ms-version"
	HeaderXmsRequestID       = "x-ms-request-id"
	HeaderXmsClientRequestID = "x-ms-client-request-id"
	HeaderDate               = "Date"
	HeaderXmsStructuredBody  = "x-ms-structured-body"
	HeaderXmsBlobType        = "x-ms-blob-type"
	HeaderXmsCopySource      = "x-ms-copy-source"
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

// SerializeBlobTags converts tags to generated.BlobTags
func SerializeBlobTags(tagsMap map[string]string) *generated.BlobTags {
	blobTagSet := make([]*generated.BlobTag, 0)
	for key, val := range tagsMap {
		newKey, newVal := key, val
		blobTagSet = append(blobTagSet, &generated.BlobTag{Key: &newKey, Value: &newVal})
	}
	return &generated.BlobTags{BlobTagSet: blobTagSet}
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
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	if host == "" {
		return false
	}
	// For IPv6, there could be case where SplitHostPort fails for cannot finding port.
	// In this case, eliminate the '[' and ']' in the URL.
	// For details about IPv6 URL, please refer to https://tools.ietf.org/html/rfc2732
	if host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1 : len(host)-1]
	}
	return net.ParseIP(host) != nil
}

// pathStylePorts are the ports at which an endpoint addresses its account in the first path
// segment even when its host is a name, as the storage emulators do (for example
// http://localhost:10000/devstoreaccount1). The list matches the other Azure Storage SDKs.
var pathStylePorts = map[string]bool{
	"10000": true, "10001": true, "10002": true, "10003": true, "10004": true,
	"10100": true, "10101": true, "10102": true, "10103": true, "10104": true,
	"11000": true, "11001": true, "11002": true, "11003": true, "11004": true,
	"11100": true, "11101": true, "11102": true, "11103": true, "11104": true,
}

// IsPathStyleURL reports whether u carries its account name in the first path segment rather
// than in its host: an IP host such as https://127.0.0.1:10000/account/..., or an emulator port
// such as http://localhost:10000/account/....
func IsPathStyleURL(u *url.URL) bool {
	return IsIPEndpointStyle(u.Host) || pathStylePorts[u.Port()]
}

// GetServiceURL reduces a service, container or blob URL to the blob service endpoint, discarding
// the container and blob path segments, the query (including any SAS) and the fragment.
// For example, "https://account.blob.core.windows.net/container/blob?sv=..." returns
// "https://account.blob.core.windows.net/", and the path-style
// "https://127.0.0.1:10000/account/container" returns "https://127.0.0.1:10000/account/".
func GetServiceURL(storageURL string) (string, error) {
	u, err := url.Parse(storageURL)
	if err != nil {
		return "", errors.New("the storage URL could not be parsed")
	}
	u.Path, u.RawPath, u.RawQuery, u.Fragment, u.RawFragment = "/", "", "", "", ""
	if IsPathStyleURL(u) {
		account, ok := firstPathSegment(storageURL)
		if !ok {
			return "", errors.New("path-style endpoint URL is missing the account name path segment")
		}
		u.Path = "/" + account + "/"
	}
	return u.String(), nil
}

// GetAccountName extracts the storage account name from a service, container or blob URL, the
// way the other Azure Storage SDKs do:
//   - path-style URLs (see IsPathStyleURL) carry it in the first path segment;
//   - otherwise it is the first label of a host whose remainder contains "blob", for example
//     "account" for account.blob.core.windows.net, with any "-ipv6" or "-dualstack" suffix and then
//     any "-secondary" suffix removed.
//
// Any other host, such as a custom domain, has no derivable account name and returns an error.
// The error never includes the URL's query.
func GetAccountName(storageURL string) (string, error) {
	u, err := url.Parse(storageURL)
	if err != nil {
		return "", errors.New("the storage URL could not be parsed")
	}
	if IsPathStyleURL(u) {
		account, ok := firstPathSegment(storageURL)
		if !ok {
			return "", errors.New("path-style endpoint URL is missing the account name path segment")
		}
		return account, nil
	}
	if account := accountNameFromHost(u.Hostname(), "blob"); account != "" {
		return account, nil
	}
	return "", fmt.Errorf("could not determine the account name from host %q", u.Hostname())
}

// accountNameFromHost returns the account name in host for the given service sub-domain, or ""
// when host isn't a storage host for that service.
func accountNameFromHost(host, serviceSubDomain string) string {
	dot := strings.Index(host, ".")
	if dot <= 0 || !strings.Contains(host[dot:], serviceSubDomain) {
		return ""
	}
	account := host[:dot]
	// trim in this order to handle names such as "account-secondary-ipv6"
	if trimmed, ok := strings.CutSuffix(account, "-ipv6"); ok {
		account = trimmed
	} else if trimmed, ok := strings.CutSuffix(account, "-dualstack"); ok {
		account = trimmed
	}
	account, _ = strings.CutSuffix(account, "-secondary")
	return account
}

// firstPathSegment returns the first segment of the URL's path.
func firstPathSegment(storageURL string) (string, bool) {
	u, err := url.Parse(storageURL)
	if err != nil {
		return "", false
	}
	segment, _, _ := strings.Cut(strings.TrimPrefix(u.Path, "/"), "/")
	return segment, segment != ""
}

// GetContainerAndBlobName splits a parsed container or blob URL into its container and blob names.
// For standard-style endpoints (e.g. "https://account.blob.core.windows.net/container/blob") the
// container is the first path segment. For path-style endpoints (see IsPathStyleURL, e.g.
// "https://127.0.0.1:10000/account/container/blob") the first path segment is the account name, so
// the container is the second.
// blob is empty when the URL addresses a container rather than a blob.
// It accepts the parsed URL rather than a string so callers on the request path do not have to
// re-serialize and re-parse a URL they already hold.
// Returns an error if the container name cannot be determined.
func GetContainerAndBlobName(u *url.URL) (container string, blob string, err error) {
	if u == nil {
		return "", "", errors.New("a URL is required to determine the container name")
	}

	path := strings.TrimPrefix(u.Path, "/")
	if IsPathStyleURL(u) {
		// path-style: scheme://host:port/accountName/container/blob...
		// drop the account name segment
		_, remainder, found := strings.Cut(path, "/")
		if !found {
			return "", "", errors.New("path-style endpoint URL is missing the container name path segment")
		}
		path = remainder
	}

	container, blob, _ = strings.Cut(path, "/")
	if container == "" {
		return "", "", errors.New("URL is missing the container name path segment")
	}
	return container, blob, nil
}

// ReadAtLeast reads from r into buf until it has read at least min bytes.
// It returns the number of bytes copied and an error.
// The EOF error is returned if no bytes were read or
// EOF happened after reading fewer than min bytes.
// If min is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
// On return, n >= min if and only if err == nil.
// If r returns an error having read at least min bytes, the error is dropped.
// This method is same as io.ReadAtLeast except that it does not
// return io.ErrUnexpectedEOF when fewer than min bytes are read.
func ReadAtLeast(r io.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	}
	return
}

// HeaderValue returns the value of the named header from h.
//
// It looks the name up both the way http.Header.Get does, which canonicalizes it
// ("x-ms-blob-type" becomes "X-Ms-Blob-Type"), and as a literal map key. The generated clients
// assign directly into the header map with the wire casing, e.g.
//
//	req.Raw().Header["x-ms-blob-type"] = []string{"BlockBlob"}
//
// so those headers are invisible to Get and would otherwise read as absent.
func HeaderValue(h http.Header, name string) string {
	if v := h.Get(name); v != "" {
		return v
	}
	if v := h[name]; len(v) > 0 {
		return v[0]
	}
	return ""
}

// ParseURLWithoutQuery returns rawURL without its user information, query and fragment, so that
// it can be logged or returned in an error without disclosing a SAS or other credentials.
func ParseURLWithoutQuery(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("the URL could not be parsed")
	}
	u.User, u.RawQuery, u.ForceQuery, u.Fragment, u.RawFragment = nil, "", false, "", ""
	return u.String(), nil
}
