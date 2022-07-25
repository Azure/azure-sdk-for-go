//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
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

// ParsedConnectionString is parsed connection string
type ParsedConnectionString = shared.ParsedConnectionString

// ParseConnectionString returns ParsedConnectionString
func ParseConnectionString(connectionString string) (ParsedConnectionString, error) {
	return shared.ParseConnectionString(connectionString)
}

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
	return exported.ParseBlobURL(u)
}

// isIPEndpointStyle checkes if URL's host is IP, in this case the storage account endpoint will be composed as:
// http(s)://IP(:port)/storageaccount/container/...
// As url's Host property, host could be both host or host:port
// nolint
func isIPEndpointStyle(host string) bool {
	return exported.IsIPEndpointStyle(host)
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
