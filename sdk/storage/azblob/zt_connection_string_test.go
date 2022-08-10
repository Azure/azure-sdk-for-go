//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//
//
//import (
//	"encoding/base64"
//	"github.com/stretchr/testify/require"
//	"strings"
//)
//
//func getAccountKey(cred *SharedKeyCredential) string {
//	return base64.StdEncoding.EncodeToString(cred.accountKey.Load().([]byte))
//}
//
//func (s *azblobTestSuite) TestConnectionStringParser() {
//	_require := require.New(s.T())
//
//	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
//	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
//	_require.Nil(err)
//	_require.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net")
//	_require.NotNil(sharedKeyCred)
//
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//
//	client := NewServiceClientFromConnectionString(connStr, nil)
//	_require.Nil(err)
//	_require.NotNil(client)
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//	_require.True(strings.HasPrefix(client.client.endpoint, "https://"))
//	_require.True(strings.Contains(client.client.endpoint, "core.windows.net"))
//}
//
//func (s *azblobTestSuite) TestConnectionStringParserHTTP() {
//	_require := require.New(s.T())
//
//	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
//	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
//	_require.Nil(err)
//	_require.Equal(serviceURL, "http://dummyaccount.blob.core.windows.net")
//	_require.NotNil(sharedKeyCred)
//
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//
//	client := NewServiceClientFromConnectionString(connStr, nil)
//	_require.Nil(err)
//	_require.NotNil(client)
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//	_require.True(strings.HasPrefix(client.client.endpoint, "http://"))
//	_require.True(strings.Contains(client.client.endpoint, "core.windows.net"))
//}
//
//func (s *azblobTestSuite) TestConnectionStringParserBasic() {
//	_require := require.New(s.T())
//	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
//	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
//	_require.Nil(err)
//	_require.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net")
//	_require.NotNil(sharedKeyCred)
//
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//
//	client := NewServiceClientFromConnectionString(connStr, nil)
//	_require.Nil(err)
//	_require.NotNil(client)
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//	_require.True(strings.HasPrefix(client.client.endpoint, "https://"))
//	_require.True(strings.Contains(client.client.endpoint, "core.windows.net"))
//}
//
//func (s *azblobTestSuite) TestConnectionStringParserCustomDomain() {
//	_require := require.New(s.T())
//	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;BlobEndpoint=www.mydomain.com;"
//	serviceURL, sharedKeyCred, err := parseConnectionString(connStr)
//	_require.Nil(err)
//	_require.Equal(serviceURL, "www.mydomain.com")
//	_require.NotNil(sharedKeyCred)
//
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//
//	client := NewServiceClientFromConnectionString(connStr, nil)
//	_require.Nil(err)
//	_require.NotNil(client)
//	_require.Equal(sharedKeyCred.accountName, "dummyaccount")
//	_require.Equal(getAccountKey(sharedKeyCred), "secretkeykey")
//	_require.True(strings.HasPrefix(client.client.endpoint, "www."))
//	_require.True(strings.Contains(client.client.endpoint, "mydomain.com"))
//}
//
//func (s *azblobTestSuite) TestConnectionStringParserInvalid() {
//	_require := require.New(s.T())
//	badConnectionStrings := []string{
//		"",
//		"foobar",
//		"foo;bar;baz",
//		"foo=;bar=;",
//		"=",
//		";",
//		"=;==",
//		"foobar=baz=foo",
//	}
//
//	for _, badConnStr := range badConnectionStrings {
//		_, _, err := parseConnectionString(badConnStr)
//		_require.NotNil(err)
//		_require.Contains(err.Error(), errConnectionString.Error())
//	}
//}
//
//func (s *azblobTestSuite) TestConnectionStringSAS() {
//	_require := require.New(s.T())
//	connStr := "AccountName=dummyaccount;SharedAccessSignature=fakesharedaccesssignature;"
//	serviceURL, cred, err := parseConnectionString(connStr)
//	_require.Nil(err)
//	_require.Equal(serviceURL, "https://dummyaccount.blob.core.windows.net/?fakesharedaccesssignature")
//	_require.Nil(cred)
//
//	client := NewServiceClientFromConnectionString(connStr, nil)
//	_require.Nil(err)
//	_require.NotNil(client)
//	_require.True(strings.HasPrefix(client.client.endpoint, "https://"))
//	_require.True(strings.Contains(client.client.endpoint, "core.windows.net"))
//}
//
//func (s *azblobTestSuite) TestConnectionStringChinaCloud() {
//	_require := require.New(s.T())
//	connStr := "AccountName=dummyaccountname;AccountKey=secretkeykey;DefaultEndpointsProtocol=http;EndpointSuffix=core.chinacloudapi.cn;"
//	serviceURL, cred, err := parseConnectionString(connStr)
//	_require.Nil(err)
//	_require.Equal(serviceURL, "http://dummyaccountname.blob.core.chinacloudapi.cn")
//	_require.NotNil(cred)
//
//	client := NewServiceClientFromConnectionString(connStr, nil)
//	_require.Nil(err)
//	_require.NotNil(client)
//	_require.True(strings.HasPrefix(client.client.endpoint, "http://"))
//	_require.True(strings.Contains(client.client.endpoint, "core.chinacloudapi.cn"))
//
//	_require.NotNil(client.sharedKey)
//	_require.Equal(client.sharedKey.accountName, "dummyaccountname")
//	_require.Equal(getAccountKey(client.sharedKey), "secretkeykey")
//}
//
//func (s *azblobTestSuite) TestConnectionStringAzurite() {
//	_require := require.New(s.T())
//	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccountname;AccountKey=secretkeykey;BlobEndpoint=http://local-machine:11002/custom/account/path/faketokensignature;"
//	serviceURL, cred, err := parseConnectionString(connStr)
//	_require.Nil(err)
//	_require.Equal(serviceURL, "http://local-machine:11002/custom/account/path/faketokensignature")
//	_require.NotNil(cred)
//
//	client := NewServiceClientFromConnectionString(connStr, nil)
//	_require.Nil(err)
//	_require.NotNil(client)
//	_require.True(strings.HasPrefix(client.client.endpoint, "http://"))
//	_require.True(strings.Contains(client.client.endpoint, "http://local-machine:11002/custom/account/path/faketokensignature"))
//
//	_require.NotNil(client.sharedKey)
//	_require.Equal(client.sharedKey.accountName, "dummyaccountname")
//	_require.Equal(getAccountKey(client.sharedKey), "secretkeykey")
//}
