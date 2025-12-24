// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
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
	require.Equal(t, "https://dummyaccount.queue.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringHTTP(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://dummyaccount.queue.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringSuffixTrailingSlash(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net/"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.queue.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringBasic(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.queue.core.windows.net/", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringCustomDomain(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;QueueEndpoint=www.mydomain.com;"
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
	require.Equal(t, "https://dummyaccount.queue.core.windows.net/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringSASAndEndpointAndCustomDomain(t *testing.T) {
	connStr := "AccountName=devstoreaccount1;SharedAccessSignature=fakesharedaccesssignature;QueueEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
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
	require.Equal(t, "https://dummyaccount.queue.core.windows.net/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringSASAndEndpoint(t *testing.T) {
	connStr := "SharedAccessSignature=fakesharedaccesssignature;QueueEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://127.0.0.1:10000/devstoreaccount1/?fakesharedaccesssignature", parsed.ServiceURL)
	require.Empty(t, parsed.AccountName)
	require.Empty(t, parsed.AccountKey)
}

func TestParseConnectionStringSASAndEndpointTrailingSlash(t *testing.T) {
	connStr := "SharedAccessSignature=fakesharedaccesssignature;QueueEndpoint=http://127.0.0.1:10000/devstoreaccount1/;"
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
	require.Equal(t, "http://dummyaccountname.queue.core.chinacloudapi.cn/", parsed.ServiceURL)
	require.Equal(t, "dummyaccountname", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringAzurite(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;QueueEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://local-machine:11002/custom/account/path/faketokensignature/", parsed.ServiceURL)
	require.Equal(t, "dummyaccountname", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}
