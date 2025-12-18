//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sas

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseURLIPStyle(t *testing.T) {
	urlWithIP := "https://127.0.0.1:5000/fakestorageaccount"
	blobURLParts, err := ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, blobURLParts.Scheme, "https")
	require.Equal(t, blobURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, blobURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")

	urlWithIP = "https://127.0.0.1:5000/fakestorageaccount/fakecontainer"
	blobURLParts, err = ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, blobURLParts.Scheme, "https")
	require.Equal(t, blobURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, blobURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")
	require.Equal(t, blobURLParts.FileSystemName, "fakecontainer")

	urlWithIP = "https://127.0.0.1:5000/fakestorageaccount/fakecontainer/fakeblob"
	blobURLParts, err = ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, blobURLParts.Scheme, "https")
	require.Equal(t, blobURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, blobURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")
	require.Equal(t, blobURLParts.FileSystemName, "fakecontainer")
	require.Equal(t, blobURLParts.PathName, "fakeblob")
}

func TestParseURL(t *testing.T) {
	testStorageAccount := "fakestorageaccount"
	host := fmt.Sprintf("%s.blob.core.windows.net", testStorageAccount)
	testContainer := "fakecontainer"
	fileNames := []string{"/._.TESTT.txt", "/.gitignore/dummyfile1"}

	const sasStr = "?sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"

	for _, fileName := range fileNames {
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s?%s", testStorageAccount, testContainer, fileName, sasStr)
		blobURLParts, err := ParseURL(urlWithVersion)
		require.NoError(t, err)

		require.Equal(t, blobURLParts.Scheme, "https")
		require.Equal(t, blobURLParts.Host, host)
		require.Equal(t, blobURLParts.FileSystemName, testContainer)

		validateSAS(t, sasStr, blobURLParts.SAS)
	}

	for _, fileName := range fileNames {
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s?%s", testStorageAccount, testContainer, fileName, sasStr)
		blobURLParts, err := ParseURL(urlWithVersion)
		require.NoError(t, err)

		require.Equal(t, blobURLParts.Scheme, "https")
		require.Equal(t, blobURLParts.Host, host)
		require.Equal(t, blobURLParts.FileSystemName, testContainer)

		validateSAS(t, sasStr, blobURLParts.SAS)
	}
}
