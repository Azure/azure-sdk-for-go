package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/runtime"
	"net"
	"net/url"
	"strings"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// ParsedConnectionString
type ParsedConnectionString = shared.ParsedConnectionString

func ParseConnectionString(connectionString string) (ParsedConnectionString, error) {
	return shared.ParseConnectionString(connectionString)
}

type nopClosingTransferManager struct {
	runtime.TransferManager
}

func (n *nopClosingTransferManager) Close() {
	// do nothing
}

const (
	snapshot  = "snapshot"
	versionId = "versionid"
)

// IPEndpointStyleInfo is used for IP endpoint style URL when working with Azure storage emulator.
// Ex: "https://10.132.141.33/accountname/containername"
type IPEndpointStyleInfo = exported.IPEndpointStyleInfo

// BlobURLParts object represents the components that make up an Azure Storage Container/Blob URL. You parse an
// existing URL into its parts by calling NewBlobURLParts(). You construct a URL from parts by calling URL().
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type BlobURLParts = exported.BlobURLParts

// ParseBlobURL parses a URL initializing BlobURLParts' fields including any SAS-related & snapshot query parameters. Any other
// query parameters remain in the UnparsedParams field. This method overwrites all fields in the BlobURLParts object.
func ParseBlobURL(u string) (BlobURLParts, error) {
	uri, err := url.Parse(u)
	if err != nil {
		return BlobURLParts{}, err
	}

	up := BlobURLParts{
		Scheme: uri.Scheme,
		Host:   uri.Host,
	}

	// Find the container & blob names (if any)
	if uri.Path != "" {
		path := uri.Path
		if path[0] == '/' {
			path = path[1:] // If path starts with a slash, remove it
		}
		if isIPEndpointStyle(up.Host) {
			if accountEndIndex := strings.Index(path, "/"); accountEndIndex == -1 { // Slash not found; path has account name & no container name or blob
				up.IPEndpointStyleInfo.AccountName = path
				path = "" // No ContainerName present in the URL so path should be empty
			} else {
				up.IPEndpointStyleInfo.AccountName = path[:accountEndIndex] // The account name is the part between the slashes
				path = path[accountEndIndex+1:]                             // path refers to portion after the account name now (container & blob names)
			}
		}

		containerEndIndex := strings.Index(path, "/") // Find the next slash (if it exists)
		if containerEndIndex == -1 {                  // Slash not found; path has container name & no blob name
			up.ContainerName = path
		} else {
			up.ContainerName = path[:containerEndIndex] // The container name is the part between the slashes
			up.BlobName = path[containerEndIndex+1:]    // The blob name is after the container slash
		}
	}

	// Convert the query parameters to a case-sensitive map & trim whitespace
	paramsMap := uri.Query()

	up.Snapshot = "" // Assume no snapshot
	if snapshotStr, ok := caseInsensitiveValues(paramsMap).Get(snapshot); ok {
		up.Snapshot = snapshotStr[0]
		// If we recognized the query parameter, remove it from the map
		delete(paramsMap, snapshot)
	}

	up.VersionID = "" // Assume no versionID
	if versionIDs, ok := caseInsensitiveValues(paramsMap).Get(versionId); ok {
		up.VersionID = versionIDs[0]
		// If we recognized the query parameter, remove it from the map
		delete(paramsMap, versionId)   // delete "versionid" from paramsMap
		delete(paramsMap, "versionId") // delete "versionId" from paramsMap
	}

	up.SAS = exported.NewSASQueryParameters(paramsMap, true)
	up.UnparsedParams = paramsMap.Encode()
	return up, nil
}

// BlobURLToString returns a URL object whose fields are initialized from the BlobURLParts fields. The URL's RawQuery
// field contains the SAS, snapshot, and unparsed query parameters.
func BlobURLToString(up BlobURLParts) string {
	return up.URL()
}

// isIPEndpointStyle checkes if URL's host is IP, in this case the storage account endpoint will be composed as:
// http(s)://IP(:port)/storageaccount/container/...
// As url's Host property, host could be both host or host:port
func isIPEndpointStyle(host string) bool {
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

type caseInsensitiveValues url.Values // map[string][]string

func (values caseInsensitiveValues) Get(key string) ([]string, bool) {
	key = strings.ToLower(key)
	for k, v := range values {
		if strings.ToLower(k) == key {
			return v, true
		}
	}
	return []string{}, false
}
