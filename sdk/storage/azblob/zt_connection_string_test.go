// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getAccountKey(cred *SharedKeyCredential) string {
	return base64.StdEncoding.EncodeToString(cred.accountKey.Load().([]byte))
}

func TestConnectionStringParser(t *testing.T) {
	_assert := assert.New(t)

	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.NoError(err)
	_assert.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.NoError(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserHTTP(t *testing.T) {
	_assert := assert.New(t)

	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.NoError(err)
	_assert.Equal(serviceURL, "http://dummyaccount.blob.core.windows.net")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.NoError(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserBasic(t *testing.T) {
	_assert := assert.New(t)
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.NoError(err)
	_assert.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.NoError(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringParserCustomDomain(t *testing.T) {
	_assert := assert.New(t)
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;BlobEndpoint=www.mydomain.com;"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.NoError(err)
	_assert.Equal(serviceURL, "www.mydomain.com")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.NoError(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "www."))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "mydomain.com"))
}

func TestConnectionStringParserInvalid(t *testing.T) {
	_assert := assert.New(t)
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
		_assert.Error(err)
		_assert.Contains(err.Error(), errConnectionString.Error())
	}
}

func TestConnectionStringSAS(t *testing.T) {
	_assert := assert.New(t)
	connStr := "AccountName=dummyaccount;SharedAccessSignature=fakesharedaccesssignature;"
	serviceURL, cred, err := parseConnectionString(connStr)
	_assert.NoError(err)
	_assert.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net/?fakesharedaccesssignature")
	_assert.Nil(cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.NoError(err)
	_assert.NotNil(client)
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func TestConnectionStringChinaCloud(t *testing.T) {
	_assert := assert.New(t)
	connStr := "AccountName=dummyaccountname;AccountKey=secretkeykey;DefaultEndpointsProtocol=http;EndpointSuffix=core.chinacloudapi.cn;"
	serviceURL, cred, err := parseConnectionString(connStr)
	_assert.NoError(err)
	_assert.Equal(serviceURL, "http://dummyaccountname.blob.core.chinacloudapi.cn")
	_assert.NotNil(cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.NoError(err)
	_assert.NotNil(client)
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.chinacloudapi.cn"))

	_assert.NotNil(client.sharedKey)
	_assert.Equal(client.sharedKey.accountName, "dummyaccountname")
	_assert.Equal(getAccountKey(client.sharedKey), "secretkeykey")
}

func TestConnectionStringAzurite(t *testing.T) {
	_assert := assert.New(t)
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;BlobEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
	serviceURL, cred, err := parseConnectionString(connStr)
	_assert.NoError(err)
	_assert.Equal(serviceURL, "http://local-machine:11002/custom/account/path/faketokensignature")
	_assert.NotNil(cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.NoError(err)
	_assert.NotNil(client)
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "http://local-machine:11002/custom/account/path/faketokensignature"))

	_assert.NotNil(client.sharedKey)
	_assert.Equal(client.sharedKey.accountName, "dummyaccountname")
	_assert.Equal(getAccountKey(client.sharedKey), "secretkeykey")
}
