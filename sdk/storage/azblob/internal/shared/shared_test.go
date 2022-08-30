//go:build go1.18
// +build go1.18

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
		//require.Contains(t, err.Error(), errConnectionString.Error())
	}
}

func TestParseConnectionString(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.blob.core.windows.net", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringHTTP(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://dummyaccount.blob.core.windows.net", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringBasic(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "https://dummyaccount.blob.core.windows.net", parsed.ServiceURL)
	require.Equal(t, "dummyaccount", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestParseConnectionStringCustomDomain(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;BlobEndpoint=www.mydomain.com;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "www.mydomain.com", parsed.ServiceURL)
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

func TestParseConnectionStringChinaCloud(t *testing.T) {
	connStr := "AccountName=dummyaccountname;AccountKey=secretkeykey;DefaultEndpointsProtocol=http;EndpointSuffix=core.chinacloudapi.cn;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://dummyaccountname.blob.core.chinacloudapi.cn", parsed.ServiceURL)
	require.Equal(t, "dummyaccountname", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}

func TestCParseConnectionStringAzurite(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;BlobEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://local-machine:11002/custom/account/path/faketokensignature", parsed.ServiceURL)
	require.Equal(t, "dummyaccountname", parsed.AccountName)
	require.Equal(t, "secretkeykey", parsed.AccountKey)
}
