// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"strings"
)

func getAccountKey(cred *SharedKeyCredential) string {
	return base64.StdEncoding.EncodeToString(cred.accountKey.Load().([]byte))
}

func (s *azblobTestSuite) TestConnectionStringParser() {
	_assert := assert.New(s.T())

	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.Nil(err)
	_assert.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.Nil(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func (s *azblobTestSuite) TestConnectionStringParserHTTP() {
	_assert := assert.New(s.T())

	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.Nil(err)
	_assert.Equal(serviceURL, "http://dummyaccount.blob.core.windows.net")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.Nil(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func (s *azblobTestSuite) TestConnectionStringParserBasic() {
	_assert := assert.New(s.T())
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.Nil(err)
	_assert.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.Nil(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func (s *azblobTestSuite) TestConnectionStringParserCustomDomain() {
	_assert := assert.New(s.T())
	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;BlobEndpoint=www.mydomain.com;"
	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
	_assert.Nil(err)
	_assert.Equal(serviceURL, "www.mydomain.com")
	_assert.NotNil(sharedKeyCred)

	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.Nil(err)
	_assert.NotNil(client)
	_assert.Equal(sharedKeyCred.accountName, "dummyaccount")
	_assert.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "www."))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "mydomain.com"))
}

func (s *azblobTestSuite) TestConnectionStringParserInvalid() {
	_assert := assert.New(s.T())
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
		_assert.NotNil(err)
		_assert.Contains(err.Error(), errConnectionString.Error())
	}
}

func (s *azblobTestSuite) TestConnectionStringSAS() {
	_assert := assert.New(s.T())
	connStr := "AccountName=dummyaccount;SharedAccessSignature=fakesharedaccesssignature;"
	serviceURL, cred, err := parseConnectionString(connStr)
	_assert.Nil(err)
	_assert.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net/?fakesharedaccesssignature")
	_assert.Nil(cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.Nil(err)
	_assert.NotNil(client)
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "https://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.windows.net"))
}

func (s *azblobTestSuite) TestConnectionStringChinaCloud() {
	_assert := assert.New(s.T())
	connStr := "AccountName=dummyaccountname;AccountKey=secretkeykey;DefaultEndpointsProtocol=http;EndpointSuffix=core.chinacloudapi.cn;"
	serviceURL, cred, err := parseConnectionString(connStr)
	_assert.Nil(err)
	_assert.Equal(serviceURL, "http://dummyaccountname.blob.core.chinacloudapi.cn")
	_assert.NotNil(cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.Nil(err)
	_assert.NotNil(client)
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "core.chinacloudapi.cn"))

	_assert.NotNil(client.sharedKey)
	_assert.Equal(client.sharedKey.accountName, "dummyaccountname")
	_assert.Equal(getAccountKey(client.sharedKey), "secretkeykey")
}

func (s *azblobTestSuite) TestConnectionStringAzurite() {
	_assert := assert.New(s.T())
	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;BlobEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
	serviceURL, cred, err := parseConnectionString(connStr)
	_assert.Nil(err)
	_assert.Equal(serviceURL, "http://local-machine:11002/custom/account/path/faketokensignature")
	_assert.NotNil(cred)

	client, err := NewServiceClientFromConnectionString(connStr, nil)
	_assert.Nil(err)
	_assert.NotNil(client)
	_assert.True(strings.HasPrefix(client.client.con.Endpoint(), "http://"))
	_assert.True(strings.Contains(client.client.con.Endpoint(), "http://local-machine:11002/custom/account/path/faketokensignature"))

	_assert.NotNil(client.sharedKey)
	_assert.Equal(client.sharedKey.accountName, "dummyaccountname")
	_assert.Equal(getAccountKey(client.sharedKey), "secretkeykey")
}
