// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

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

	sharedKeyCred, ok := cred.(*SharedKeyCredential)
	require.True(t, ok)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewTableClientFromConnectionString("tableName", connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(&client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.u, "https://"))
	require.True(t, strings.Contains(client.client.con.u, "core.windows.net"))
}

func TestConnectionStringParserHTTP(t *testing.T) {
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "http://dummyaccount.table.core.windows.net")
	require.NotNil(t, cred)

	sharedKeyCred, ok := cred.(*SharedKeyCredential)
	require.True(t, ok)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewTableClientFromConnectionString("tableName", connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(&client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.u, "http://"))
	require.True(t, strings.Contains(client.client.con.u, "core.windows.net"))
}

func TestConnectionStringParserBasic(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "https://dummyaccount.table.core.windows.net")
	require.NotNil(t, cred)

	sharedKeyCred, ok := cred.(*SharedKeyCredential)
	require.True(t, ok)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewTableClientFromConnectionString("tableName", connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(&client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.u, "https://"))
	require.True(t, strings.Contains(client.client.con.u, "core.windows.net"))
}

func TestConnectionStringParserCustomDomain(t *testing.T) {
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;TableEndpoint=www.mydomain.com;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, serviceURL, "www.mydomain.com")
	require.NotNil(t, cred)

	sharedKeyCred, ok := cred.(*SharedKeyCredential)
	require.True(t, ok)
	require.Equal(t, sharedKeyCred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewTableClientFromConnectionString("tableName", connStr, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, client.cred.accountName, "dummyaccount")
	require.Equal(t, getAccountKey(&client.cred), "secretkeykey")
	require.True(t, strings.HasPrefix(client.client.con.u, "www."))
	require.True(t, strings.Contains(client.client.con.u, "mydomain.com"))
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
		require.Contains(t, err.Error(), ErrConnectionString.Error())
	}
}
