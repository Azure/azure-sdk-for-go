package azfile

import (
	"net"
	"net/url"
	"strings"
)

const (
	shareSnapshot      = "sharesnapshot"
	SnapshotTimeFormat = "2006-01-02T15:04:05.0000000Z07:00"
)

// A FileURLParts object represents the components that make up an Azure Storage Share/Directory/File URL. You parse an
// existing URL into its parts by calling NewFileURLParts(). You construct a URL from parts by calling URL().
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type FileURLParts struct {
	Scheme              string // Ex: "https://"
	Host                string // Ex: "account.share.core.windows.net", "10.132.141.33", "10.132.141.33:80"
	ShareName           string // Share name, Ex: "myshare"
	DirectoryOrFilePath string // Path of directory or file, Ex: "mydirectory/myfile"
	ShareSnapshot       string // IsZero is true if not a snapshot
	SAS                 SASQueryParameters
	UnparsedParams      string
	IPEndpointStyleInfo IPEndpointStyleInfo // Useful Parts for IP endpoint style URL.
}

// IPEndpointStyleInfo is used for IP endpoint style URL.
// It's commonly used when working with Azure storage emulator or testing environments.
// Ex: "https://10.132.141.33/accountname/sharename"
type IPEndpointStyleInfo struct {
	AccountName string // "" if not using IP endpoint style
}

// isIPEndpointStyle checkes if URL's host is IP, in this case the storage account endpoint will be composed as:
// http(s)://IP(:port)/storageaccount/share(||container||etc)/...
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

// NewFileURLParts parses a URL initializing FileURLParts' fields including any SAS-related & sharesnapshot query parameters. Any other
// query parameters remain in the UnparsedParams field. This method overwrites all fields in the FileURLParts object.
func NewFileURLParts(u string) FileURLParts {
	uri, _ := url.Parse(u)

	up := FileURLParts{
		Scheme:              uri.Scheme,
		Host:                uri.Host,
		IPEndpointStyleInfo: IPEndpointStyleInfo{},
	}

	if uri.Path != "" {
		path := uri.Path

		if path[0] == '/' {
			path = path[1:]
		}

		if isIPEndpointStyle(up.Host) {
			if accountEndIndex := strings.Index(path, "/"); accountEndIndex == -1 { // Slash not found; path has account name & no share, path of directory or file
				up.IPEndpointStyleInfo.AccountName = path
			} else {
				up.IPEndpointStyleInfo.AccountName = path[:accountEndIndex] // The account name is the part between the slashes

				path = path[accountEndIndex+1:]
				// Find the next slash (if it exists)
				if shareEndIndex := strings.Index(path, "/"); shareEndIndex == -1 { // Slash not found; path has share name & no path of directory or file
					up.ShareName = path
				} else { // Slash found; path has share name & path of directory or file
					up.ShareName = path[:shareEndIndex]
					up.DirectoryOrFilePath = path[shareEndIndex+1:]
				}
			}
		} else {
			// Find the next slash (if it exists)
			if shareEndIndex := strings.Index(path, "/"); shareEndIndex == -1 { // Slash not found; path has share name & no path of directory or file
				up.ShareName = path
			} else { // Slash found; path has share name & path of directory or file
				up.ShareName = path[:shareEndIndex]
				up.DirectoryOrFilePath = path[shareEndIndex+1:]
			}
		}
	}

	// Convert the query parameters to a case-sensitive map & trim whitespace
	paramsMap := uri.Query()

	if snapshotStr, ok := caseInsensitiveValues(paramsMap).Get(shareSnapshot); ok {
		up.ShareSnapshot = snapshotStr[0]
		// If we recognized the query parameter, remove it from the map
		delete(paramsMap, shareSnapshot)
	}
	up.SAS = newSASQueryParameters(paramsMap, true)
	up.UnparsedParams = paramsMap.Encode()
	return up
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

// URL returns a URL object whose fields are initialized from the FileURLParts fields. The URL's RawQuery
// field contains the SAS, snapshot, and unparsed query parameters.
func (up FileURLParts) URL() string {
	path := ""
	// Concatenate account name for IP endpoint style URL
	if isIPEndpointStyle(up.Host) && up.IPEndpointStyleInfo.AccountName != "" {
		path += "/" + up.IPEndpointStyleInfo.AccountName
	}
	// Concatenate share & path of directory or file (if they exist)
	if up.ShareName != "" {
		path += "/" + up.ShareName
		if up.DirectoryOrFilePath != "" {
			path += "/" + up.DirectoryOrFilePath
		}
	}

	rawQuery := up.UnparsedParams

	// Concatenate share snapshot query parameter (if it exists)
	if up.ShareSnapshot != "" {
		if len(rawQuery) > 0 {
			rawQuery += "&"
		}
		rawQuery += shareSnapshot + "=" + up.ShareSnapshot
	}
	sas := up.SAS.Encode()
	if sas != "" {
		if len(rawQuery) > 0 {
			rawQuery += "&"
		}
		rawQuery += sas
	}
	u := url.URL{
		Scheme:   up.Scheme,
		Host:     up.Host,
		Path:     path,
		RawQuery: rawQuery,
	}
	return u.String()
}
