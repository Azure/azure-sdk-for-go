// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"net/url"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConnectionStringInvalid(t *testing.T) {
	badConnectionStrings := []string{
		"",
		"foobar",
		"foo;bar;baz",
		"foo=;bar=;",
		"=",
		";",
		"=;==",
		"foobar=baz=foo",
	}

	for _, badConnStr := range badConnectionStrings {
		parsed, err := ParseConnectionString(badConnStr)
		require.Error(t, err)
		require.Zero(t, parsed)
	}
}

func TestParseConnectionString(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.blob.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringHTTP(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://dummyaccount.blob.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringSuffixTrailingSlash(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net/"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.blob.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringBasic(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.blob.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringCustomDomain(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;BlobEndpoint=www.mydomain.com;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "www.mydomain.com/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringSAS(t *testing.T) {
	connStr := "AccountName=dummyaccount;SharedAccessSignature=fakesharedaccesssignature;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.blob.core.windows.net/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringSASAndEndpointAndAccountName(t *testing.T) {
	connStr := "AccountName=devstoreaccount1;SharedAccessSignature=fakesharedaccesssignature;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://127.0.0.1:10000/devstoreaccount1/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringSASSuffixTrailingSlash(t *testing.T) {
	connStr := "AccountName=dummyaccount;SharedAccessSignature=fakesharedaccesssignature;EndpointSuffix=core.windows.net/"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.blob.core.windows.net/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringSASAndEndpoint(t *testing.T) {
	connStr := "SharedAccessSignature=fakesharedaccesssignature;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://127.0.0.1:10000/devstoreaccount1/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringSASAndEndpointTrailingSlash(t *testing.T) {
	connStr := "SharedAccessSignature=fakesharedaccesssignature;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1/;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://127.0.0.1:10000/devstoreaccount1/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringChinaCloud(t *testing.T) {
	connStr := "AccountName=dummyaccountname;AccountKey=secretkeykey;DefaultEndpointsProtocol=http;EndpointSuffix=core.chinacloudapi.cn;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://dummyaccountname.blob.core.chinacloudapi.cn/", parsed.ServiceURL)
	require.Equal(t, "dummyaccountname", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringAzurite(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;BlobEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://local-machine:11002/custom/account/path/faketokensignature/", parsed.ServiceURL)
	require.Equal(t, "dummyaccountname", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestSerializeBlobTags(t *testing.T) {
	var tags map[string]string

	// Case 1
	tags = nil
	blobTags := SerializeBlobTags(tags)
	require.NotNil(t, blobTags)
	require.Len(t, blobTags.BlobTagSet, 0)

	// Case 2
	tags = map[string]string{}
	blobTags = SerializeBlobTags(tags)
	require.NotNil(t, blobTags)
	require.Len(t, blobTags.BlobTagSet, 0)

	// Case 3
	tags = map[string]string{
		"foo": "bar",
		"az":  "sdk",
		"sdk": "storage",
	}
	blobTags = SerializeBlobTags(tags)
	require.NotNil(t, blobTags)
	for _, tagPtr := range blobTags.BlobTagSet {
		require.Contains(t, tags, *tagPtr.Key)
		require.Equal(t, tags[*tagPtr.Key], *tagPtr.Value)
		delete(tags, *tagPtr.Key)
	}
	require.Len(t, tags, 0)
}

func TestSerializeBlobTagsToStrPtr(t *testing.T) {
	var tags map[string]string

	// Case 1
	tags = nil
	tagsStr := SerializeBlobTagsToStrPtr(tags)
	require.Nil(t, tagsStr)

	// Case 2
	tags = map[string]string{}
	tagsStr = SerializeBlobTagsToStrPtr(tags)
	require.Nil(t, tagsStr)

	// Case 3
	tags = map[string]string{
		"foo": "bar",
		"az":  "sdk",
		"sdk": "storage",
	}
	tagsStr = SerializeBlobTagsToStrPtr(tags)
	require.NotNil(t, tagsStr)
	// split string on &
	kvPairs := strings.Split(*tagsStr, "&")
	for _, kv := range kvPairs {
		pair := strings.Split(kv, "=")
		require.Contains(t, tags, pair[0])
		require.Equal(t, tags[pair[0]], pair[1])
		delete(tags, pair[0])
	}
	require.Len(t, tags, 0)
}

func TestIsIPEndpointStyle(t *testing.T) {
	require.False(t, IsIPEndpointStyle(""))
	require.False(t, IsIPEndpointStyle(":0"))

	require.True(t, IsIPEndpointStyle("127.0.0.1"))
	require.True(t, IsIPEndpointStyle("127.0.0.1:80"))
}

func TestDefaultConcurrencyValue_InBounds(t *testing.T) {
	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "")
	val := DefaultConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))
}

func TestDefaultConcurrencyValue_Deterministic(t *testing.T) {
	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "")
	val1 := DefaultConcurrencyValue()
	val2 := DefaultConcurrencyValue()
	require.Equal(t, val1, val2)
}

func TestDefaultConcurrencyValue_MatchesCPU(t *testing.T) {
	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "")
	cpus := runtime.NumCPU()
	val := DefaultConcurrencyValue()
	if cpus < 8 {
		require.Equal(t, uint16(8), val)
	} else if cpus > 96 {
		require.Equal(t, uint16(96), val)
	} else {
		require.Equal(t, uint16(cpus), val)
	}
}

func TestDefaultConcurrencyValue_LegacyEnvVar(t *testing.T) {
	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "true")
	require.Equal(t, uint16(DefaultConcurrency), DefaultConcurrencyValue())

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "TRUE")
	require.Equal(t, uint16(DefaultConcurrency), DefaultConcurrencyValue())

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "false")
	val := DefaultConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "")
	val = DefaultConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))
}

func TestDefaultStreamConcurrencyValue_InBounds(t *testing.T) {
	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "")
	val := DefaultStreamConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))
}

func TestDefaultStreamConcurrencyValue_LegacyEnvVar(t *testing.T) {
	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "true")
	require.Equal(t, uint16(1), DefaultStreamConcurrencyValue())

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "TRUE")
	require.Equal(t, uint16(1), DefaultStreamConcurrencyValue())

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "false")
	val := DefaultStreamConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))
}

func TestGetServiceURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expectedURL string
		expectError bool
		errorSubstr string
	}{
		// Standard-style URLs
		{"StandardBlobURL", "https://account.blob.core.windows.net/container/blob", "https://account.blob.core.windows.net/", false, ""},
		{"StandardContainerURL", "https://account.blob.core.windows.net/container", "https://account.blob.core.windows.net/", false, ""},
		{"StandardServiceURL", "https://account.blob.core.windows.net/", "https://account.blob.core.windows.net/", false, ""},
		{"StandardServiceURLNoTrailingSlash", "https://account.blob.core.windows.net", "https://account.blob.core.windows.net/", false, ""},

		// IP-style URLs (Azurite/emulator)
		{"IPStyleBlobURL", "https://127.0.0.1:10000/devstoreaccount1/container/blob", "https://127.0.0.1:10000/devstoreaccount1/", false, ""},
		{"IPStyleContainerURL", "https://127.0.0.1:10000/devstoreaccount1/container", "https://127.0.0.1:10000/devstoreaccount1/", false, ""},
		{"IPStyleServiceURL", "https://127.0.0.1:10000/devstoreaccount1/", "https://127.0.0.1:10000/devstoreaccount1/", false, ""},
		{"IPStyleServiceURLNoTrailingSlash", "https://127.0.0.1:10000/devstoreaccount1", "https://127.0.0.1:10000/devstoreaccount1/", false, ""},

		// SAS token preservation
		{"StandardURLWithSAS", "https://account.blob.core.windows.net/container/blob?sv=2021-06-08&ss=b&srt=sco&sig=test", "https://account.blob.core.windows.net/?sv=2021-06-08&ss=b&srt=sco&sig=test", false, ""},
		{"IPStyleURLWithSAS", "https://127.0.0.1:10000/devstoreaccount1/container?sv=2021-06-08&sig=test", "https://127.0.0.1:10000/devstoreaccount1/?sv=2021-06-08&sig=test", false, ""},

		// HTTP scheme
		{"HTTPScheme", "http://account.blob.core.windows.net/container", "http://account.blob.core.windows.net/", false, ""},

		// Custom domain
		{"CustomDomain", "https://mydomain.com/container/blob", "https://mydomain.com/", false, ""},

		// China cloud
		{"ChinaCloud", "https://account.blob.core.chinacloudapi.cn/container", "https://account.blob.core.chinacloudapi.cn/", false, ""},
		// Error cases
		{"IPStyleMissingAccountName", "https://127.0.0.1:10000/", "", true, "missing the account name"},
		{"IPStyleEmptyPath", "https://127.0.0.1:10000", "", true, "missing the account name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceURL, err := GetServiceURL(tt.inputURL)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorSubstr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedURL, serviceURL)
			}
		})
	}
}

func TestGetContainerAndBlobName(t *testing.T) {
	tests := []struct {
		name              string
		inputURL          string
		expectedContainer string
		expectedBlob      string
		expectError       bool
		errorSubstr       string
	}{
		// Standard-style URLs: container is the first path segment
		{"StandardBlobURL", "https://account.blob.core.windows.net/container/blob", "container", "blob", false, ""},
		{"StandardNestedBlobURL", "https://account.blob.core.windows.net/container/dir/sub/blob.txt", "container", "dir/sub/blob.txt", false, ""},
		{"StandardContainerURL", "https://account.blob.core.windows.net/container", "container", "", false, ""},
		{"StandardContainerURLTrailingSlash", "https://account.blob.core.windows.net/container/", "container", "", false, ""},

		// IP-style URLs: first segment is the account, container is the second
		{"IPStyleBlobURL", "https://127.0.0.1:10000/devstoreaccount1/container/blob", "container", "blob", false, ""},
		{"IPStyleNestedBlobURL", "https://127.0.0.1:10000/devstoreaccount1/container/dir/blob.txt", "container", "dir/blob.txt", false, ""},
		{"IPStyleContainerURL", "https://127.0.0.1:10000/devstoreaccount1/container", "container", "", false, ""},

		// Query strings must not leak into the names
		{"StandardURLWithSAS", "https://account.blob.core.windows.net/container/blob?sv=2021-06-08&sig=test", "container", "blob", false, ""},
		{"IPStyleURLWithSAS", "https://127.0.0.1:10000/devstoreaccount1/container?sv=2021-06-08&sig=test", "container", "", false, ""},

		// Encoded blob names are returned decoded
		{"BlobNameWithSpace", "https://account.blob.core.windows.net/container/my%20blob.txt", "container", "my blob.txt", false, ""},

		// Custom domain and sovereign clouds behave as standard-style
		{"CustomDomain", "https://mydomain.com/container/blob", "container", "blob", false, ""},
		{"ChinaCloud", "https://account.blob.core.chinacloudapi.cn/container/blob", "container", "blob", false, ""},

		// Error cases
		{"StandardServiceURL", "https://account.blob.core.windows.net/", "", "", true, "missing the container name"},
		{"StandardServiceURLNoTrailingSlash", "https://account.blob.core.windows.net", "", "", true, "missing the container name"},
		{"IPStyleAccountOnly", "https://127.0.0.1:10000/devstoreaccount1", "", "", true, "missing the container name"},
		{"IPStyleAccountOnlyTrailingSlash", "https://127.0.0.1:10000/devstoreaccount1/", "", "", true, "missing the container name"},
		{"IPStyleEmptyPath", "https://127.0.0.1:10000", "", "", true, "missing the container name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.inputURL)
			require.NoError(t, err)

			container, blob, err := GetContainerAndBlobName(u)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorSubstr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedContainer, container)
				require.Equal(t, tt.expectedBlob, blob)
			}
		})
	}
}

func TestGetContainerAndBlobNameNilURL(t *testing.T) {
	_, _, err := GetContainerAndBlobName(nil)
	require.Error(t, err)
}

func TestGetAccountName(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expected    string
		expectError bool
	}{
		// Standard-style URLs
		{"StandardBlobURL", "https://account.blob.core.windows.net/container/blob", "account", false},
		{"StandardContainerURL", "https://account.blob.core.windows.net/container", "account", false},
		{"StandardServiceURL", "https://account.blob.core.windows.net/", "account", false},
		{"StandardServiceURLNoTrailingSlash", "https://account.blob.core.windows.net", "account", false},
		{"StandardURLWithPort", "https://account.blob.core.windows.net:443/container", "account", false},
		{"StandardURLWithSAS", "https://account.blob.core.windows.net/container/blob?sig=test", "account", false},

		// IP-style URLs (Azurite/emulator): the account is the first path segment
		{"IPStyleBlobURL", "https://127.0.0.1:10000/devstoreaccount1/container/blob", "devstoreaccount1", false},
		{"IPStyleServiceURL", "https://127.0.0.1:10000/devstoreaccount1/", "devstoreaccount1", false},
		{"IPStyleServiceURLNoTrailingSlash", "http://127.0.0.1:10000/devstoreaccount1", "devstoreaccount1", false},

		// Error cases
		{"IPStyleMissingAccount", "https://127.0.0.1:10000/", "", true},
		{"IPStyleNoPath", "https://127.0.0.1:10000", "", true},
		{"HostWithoutSubdomain", "https://localhost/container/blob", "", true},
		{"EmptyHost", "/container/blob", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountName, err := GetAccountName(tt.inputURL)
			if tt.expectError {
				require.Error(t, err)
				require.Empty(t, accountName)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, accountName)
		})
	}
}
