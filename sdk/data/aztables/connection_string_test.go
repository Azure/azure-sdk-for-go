// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

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
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "https://dummyaccount.table.core.windows.net")
	require.NotNil(t, cred)

	require.Equal(t, cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(cred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.cred)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "https://"))
	require.True(t, strings.Contains(client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserHTTP(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "http://dummyaccount.table.core.windows.net")
	require.NotNil(t, cred)

	require.Equal(t, cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(cred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.cred)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "http://"))
	require.True(t, strings.Contains(client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserBasic(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "https://dummyaccount.table.core.windows.net")
	require.NotNil(t, cred)

	require.Equal(t, cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(cred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.cred)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "https://"))
	require.True(t, strings.Contains(client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserCustomDomain(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;TableEndpoint=www.mydomain.com;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "www.mydomain.com")
	require.NotNil(t, cred)

	require.Equal(t, cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(cred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.cred)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "www."))
	require.True(t, strings.Contains(client.con.Endpoint(), "mydomain.com"))
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
	require.Equal(t, serviceURL, "https://dummyaccount.table.core.windows.net/?fakesharedaccesssignature")
	require.Nil(t, cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "https://"))
	require.True(t, strings.Contains(client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringCosmos(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccountname;AccountKey=secretkeykey;TableEndpoint=https://dummyaccountname.table.cosmos.azure.com:443/;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "https://dummyaccountname.table.cosmos.azure.com:443/")
	require.NotNil(t, cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "https://"))
	require.True(t, strings.Contains(client.con.Endpoint(), "cosmos.azure.com:443"))

	require.NotNil(t, client.cred)
	require.Equal(t, client.cred.accountName, "dummyaccountname")
	require.Equal(t, getAccountKey(client.cred), "secretkeykey")
}

func TestConnectionStringChinaCloud(t *testing.T) {
	connStr := "AccountName=dummyaccountname;AccountKey=secretkeykey;DefaultEndpointsProtocol=http;EndpointSuffix=core.chinacloudapi.cn;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "http://dummyaccountname.table.core.chinacloudapi.cn")
	require.NotNil(t, cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "http://"))
	require.True(t, strings.Contains(client.con.Endpoint(), "core.chinacloudapi.cn"))

	require.Equal(t, client.cred.accountName, "dummyaccountname")
	require.Equal(t, getAccountKey(client.cred), "secretkeykey")
}

func TestConnectionStringAzurite(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;TableEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "http://local-machine:11002/custom/account/path/faketokensignature")
	require.NotNil(t, cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, strings.HasPrefix(client.con.Endpoint(), "http://"))
	require.True(t, strings.Contains(client.con.Endpoint(), "http://local-machine:11002/custom/account/path/faketokensignature"))
	require.NotNil(t, client.cred)
	require.Equal(t, client.cred.accountName, "dummyaccountname")
	require.Equal(t, getAccountKey(client.cred), "secretkeykey")
}
