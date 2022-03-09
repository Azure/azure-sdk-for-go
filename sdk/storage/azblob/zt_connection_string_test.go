// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getAccountKey(cred *SharedKeyCredential) string {
	return base64.StdEncoding.EncodeToString(cred.accountKey.Load().([]byte))
}

func TestConnectionStringParser(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "https://dummyaccount.blob.core.windows.net")
	require.NotNil(t, sharedKeyCred)

	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	require.True(t, strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserHTTP(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "http://dummyaccount.blob.core.windows.net")
	require.NotNil(t, sharedKeyCred)

	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	require.True(t, strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserBasic(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "https://dummyaccount.blob.core.windows.net")
	require.NotNil(t, sharedKeyCred)

	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	require.True(t, strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserCustomDomain(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;BlobEndpoint=www.mydomain.com;"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "www.mydomain.com")
	require.NotNil(t, sharedKeyCred)

	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.Endpoint(), "www."))
	require.True(t, strings.Contains(client.client.con.Endpoint(), "mydomain.com"))
}

func TestConnectionStringParserInvalid(t *testing.T) {
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
		_, _, err := parseConnectionString(badConnStr)
		require.Error(t, err)
		require.Contains(t, err.Error(), errConnectionString.Error())
	}
}

func TestConnectionStringSAS(t *testing.T) {
	connStr := "AccountName=dummyaccount;SharedAccessSignature=fakesharedaccesssignature;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "https://dummyaccount.blob.core.windows.net/?fakesharedaccesssignature")
	require.Nil(t, cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	require.True(t, strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringChinaCloud(t *testing.T) {
	connStr := "AccountName=dummyaccountname;AccountKey=secretkeykey;DefaultEndpointsProtocol=http;EndpointSuffix=core.chinacloudapi.cn;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "http://dummyaccountname.blob.core.chinacloudapi.cn")
	require.NotNil(t, cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	require.True(t, strings.Contains(client.client.con.Endpoint(), "core.chinacloudapi.cn"))

	require.NotNil(t, client.sharedKey)
	require.Equal(t, client.sharedKey.accountName, "dummyaccountname")
	require.Equal(t, getAccountKey(client.sharedKey), "secretkeykey")
}

func TestConnectionStringAzurite(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;BlobEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "http://local-machine:11002/custom/account/path/faketokensignature")
	require.NotNil(t, cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	require.True(t, strings.Contains(client.client.con.Endpoint(), "http://local-machine:11002/custom/account/path/faketokensignature"))

	require.NotNil(t, client.sharedKey)
	require.Equal(t, client.sharedKey.accountName, "dummyaccountname")
	require.Equal(t, getAccountKey(client.sharedKey), "secretkeykey")
}
