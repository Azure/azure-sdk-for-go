// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectionStringParser(t *testing.T) {
	require := require.New(t)

	connStr := "DefaultEndpointsProtocol=https;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(err)
	require.Equal(serviceURL, "https://dummyaccount.table.core.windows.net")
	require.NotNil(cred)
}

func TestConnectionStringParserHTTP(t *testing.T) {
	require := require.New(t)

	connStr := "DefaultEndpointsProtocol=http;AccountName=dummyaccount;AccountKey=secretkeykey;EndpointSuffix=core.windows.net"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(err)
	require.Equal(serviceURL, "http://dummyaccount.table.core.windows.net")
	require.NotNil(cred)
}

func TestConnectionStringParserBasic(t *testing.T) {
	require := require.New(t)

	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(err)
	require.Equal(serviceURL, "https://dummyaccount.table.core.windows.net")
	require.NotNil(cred)
}

func TestConnectionStringParserCustomDomain(t *testing.T) {
	require := require.New(t)

	connStr := "AccountName=dummyaccount;AccountKey=secretkeykey;TableEndpoint=www.mydomain.com;"
	serviceURL, cred, err := parseConnectionString(connStr)
	require.NoError(err)
	require.Equal(serviceURL, "www.mydomain.com")
	require.NotNil(cred)
}
