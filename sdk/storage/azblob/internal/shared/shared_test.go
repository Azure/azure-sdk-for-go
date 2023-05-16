//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
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
		//require.Contains(t, err.Error(), errConnectionString.Error())
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
	require.Nil(t, blobTags)

	// Case 2
	tags = map[string]string{}
	blobTags = SerializeBlobTags(tags)
	require.Nil(t, blobTags)

	// Case 3
	tags = map[string]string{
		"foo": "bar",
		"az":  "sdk",
		"sdk": "storage",
	}
	blobTags = SerializeBlobTags(tags)
	require.NotNil(t, blobTags)
	for _, tagPtr := range (*blobTags).BlobTagSet {
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
