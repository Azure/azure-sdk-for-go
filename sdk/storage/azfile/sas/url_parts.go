// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"net/url"
	"strings"
)

const (
	shareSnapshot = "sharesnapshot"
)

// IPEndpointStyleInfo is used for IP endpoint style URL when working with Azure storage emulator.
// Ex: "https://10.132.141.33/accountname/sharename"
type IPEndpointStyleInfo struct {
	AccountName string // "" if not using IP endpoint style
}

// URLParts object represents the components that make up an Azure Storage Share/Directory/File URL. You parse an
// existing URL into its parts by calling NewFileURLParts(). You construct a URL from parts by calling URL().
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type URLParts struct {
	Scheme              string              // Ex: "https://"
	Host                string              // Ex: "account.share.core.windows.net", "10.132.141.33", "10.132.141.33:80"
	IPEndpointStyleInfo IPEndpointStyleInfo // Useful Parts for IP endpoint style URL.
	ShareName           string              // Share name, Ex: "myshare"
	DirectoryOrFilePath string              // Path of directory or file, Ex: "mydirectory/myfile"
	ShareSnapshot       string              // IsZero is true if not a snapshot
	SAS                 QueryParameters
	UnparsedParams      string
}

// ParseURL parses a URL initializing URLParts' fields including any SAS-related & sharesnapshot query parameters.
// Any other query parameters remain in the UnparsedParams field.
func ParseURL(u string) (URLParts, error) {
	uri, err := url.Parse(u)
	if err != nil {
		return URLParts{}, err
	}

	up := URLParts{
		Scheme: uri.Scheme,
		Host:   uri.Host,
	}

	if uri.Path != "" {
		path := uri.Path
		if path[0] == '/' {
			path = path[1:]
		}
		if shared.IsIPEndpointStyle(up.Host) {
			if accountEndIndex := strings.Index(path, "/"); accountEndIndex == -1 { // Slash not found; path has account name & no share, path of directory or file
				up.IPEndpointStyleInfo.AccountName = path
				path = "" // no ShareName present in the URL so path should be empty
			} else {
				up.IPEndpointStyleInfo.AccountName = path[:accountEndIndex] // The account name is the part between the slashes
				path = path[accountEndIndex+1:]
			}
		}

		shareEndIndex := strings.Index(path, "/") // Find the next slash (if it exists)
		if shareEndIndex == -1 {                  // Slash not found; path has share name & no path of directory or file
			up.ShareName = path
		} else { // Slash found; path has share name & path of directory or file
			up.ShareName = path[:shareEndIndex]
			up.DirectoryOrFilePath = path[shareEndIndex+1:]
		}
	}

	// Convert the query parameters to a case-sensitive map & trim whitespace
	paramsMap := uri.Query()

	up.ShareSnapshot = "" // Assume no snapshot
	if snapshotStr, ok := caseInsensitiveValues(paramsMap).Get(shareSnapshot); ok {
		up.ShareSnapshot = snapshotStr[0]
		// If we recognized the query parameter, remove it from the map
		delete(paramsMap, shareSnapshot)
	}

	up.SAS = NewQueryParameters(paramsMap, true)
	up.UnparsedParams = paramsMap.Encode()
	return up, nil
}

// String returns a URL object whose fields are initialized from the URLParts fields. The URL's RawQuery
// field contains the SAS, snapshot, and unparsed query parameters.
func (up URLParts) String() string {
	path := ""
	// Concatenate account name for IP endpoint style URL
	if shared.IsIPEndpointStyle(up.Host) && up.IPEndpointStyleInfo.AccountName != "" {
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

	// If no snapshot is initially provided, fill it in from the SAS query properties to help the user
	if up.ShareSnapshot == "" && !up.SAS.ShareSnapshotTime().IsZero() {
		up.ShareSnapshot = up.SAS.ShareSnapshotTime().Format(SnapshotTimeFormat)
	}

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
