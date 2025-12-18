//go:build go1.18

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
	fileURLParts, err := ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, fileURLParts.Scheme, "https")
	require.Equal(t, fileURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, fileURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")

	urlWithIP = "https://127.0.0.1:5000/fakestorageaccount/fakeshare"
	fileURLParts, err = ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, fileURLParts.Scheme, "https")
	require.Equal(t, fileURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, fileURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")
	require.Equal(t, fileURLParts.ShareName, "fakeshare")

	urlWithIP = "https://127.0.0.1:5000/fakestorageaccount/fakeshare/fakefile"
	fileURLParts, err = ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, fileURLParts.Scheme, "https")
	require.Equal(t, fileURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, fileURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")
	require.Equal(t, fileURLParts.ShareName, "fakeshare")
	require.Equal(t, fileURLParts.DirectoryOrFilePath, "fakefile")
}

func TestParseURL(t *testing.T) {
	testStorageAccount := "fakestorageaccount"
	host := fmt.Sprintf("%s.file.core.windows.net", testStorageAccount)
	testShare := "fakeshare"
	fileNames := []string{"/._.TESTT.txt", "/.gitignore/dummyfile1"}

	const sasStr = "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"

	for _, fileName := range fileNames {
		sasURL := fmt.Sprintf("https://%s.file.core.windows.net/%s%s?%s", testStorageAccount, testShare, fileName, sasStr)
		fileURLParts, err := ParseURL(sasURL)
		require.NoError(t, err)

		require.Equal(t, fileURLParts.Scheme, "https")
		require.Equal(t, fileURLParts.Host, host)
		require.Equal(t, fileURLParts.ShareName, testShare)

		validateSAS(t, sasStr, fileURLParts.SAS)
	}

	for _, fileName := range fileNames {
		shareSnapshotID := "2011-03-09T01:42:34Z"
		sasWithShareSnapshotID := "?sharesnapshot=" + shareSnapshotID + "&" + sasStr
		urlWithShareSnapshot := fmt.Sprintf("https://%s.file.core.windows.net/%s%s%s", testStorageAccount, testShare, fileName, sasWithShareSnapshotID)
		fileURLParts, err := ParseURL(urlWithShareSnapshot)
		require.NoError(t, err)

		require.Equal(t, fileURLParts.Scheme, "https")
		require.Equal(t, fileURLParts.Host, host)
		require.Equal(t, fileURLParts.ShareName, testShare)

		validateSAS(t, sasStr, fileURLParts.SAS)
	}
}
