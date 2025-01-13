//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

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
		require.ErrorIs(t, err, errConnectionString)
		require.Zero(t, parsed)
	}
}

func TestParseConnectionString(t *testing.T) {
	connStr := "Endpoint=http://abc.com;AccessKey=ABC;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://abc.com/", parsed.Endpoint)
	require.Equal(t, "ABC", parsed.AccessKey)
}

func TestParseConnectionStringLowercase(t *testing.T) {
	connStr := "Endpoint=http://abc.com;accessKey=ABC;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://abc.com/", parsed.Endpoint)
	require.Equal(t, "ABC", parsed.AccessKey)
}

func TestParseConnectionStringWithPort(t *testing.T) {
	connStr := "Endpoint=http://abc.com:8080;accessKey=ABC;Port=8088;"
	parsed, err := ParseConnectionString(connStr)
	require.NoError(t, err)
	require.Equal(t, "http://abc.com:8088/", parsed.Endpoint)
	require.Equal(t, "ABC", parsed.AccessKey)
}
