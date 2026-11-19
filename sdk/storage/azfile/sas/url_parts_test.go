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

func TestParseURLFileID(t *testing.T) {
	const fileID = "12384898975283830"
	for _, u := range []string{
		"https://account.file.core.windows.net/share?fileid=" + fileID,
		"https://account.file.core.windows.net/share?FileId=" + fileID,
		"https://account.file.core.windows.net/share?FILEID=" + fileID,
	} {
		fileURLParts, err := ParseURL(u)
		require.NoError(t, err)
		require.Equal(t, "https", fileURLParts.Scheme)
		require.Equal(t, "account.file.core.windows.net", fileURLParts.Host)
		require.Equal(t, "share", fileURLParts.ShareName)
		require.Empty(t, fileURLParts.DirectoryOrFilePath)
		require.Equal(t, fileID, fileURLParts.FileID)
		require.Empty(t, fileURLParts.UnparsedParams, "the file ID must not also be kept as an unparsed parameter")
		require.Empty(t, fileURLParts.SAS.Encode())

		// the file ID is emitted exactly once, using the canonical parameter name
		require.Equal(t, "https://account.file.core.windows.net/share?fileid="+fileID, fileURLParts.String())
	}
}

func TestParseURLFileIDEmptyForPathBasedURL(t *testing.T) {
	fileURLParts, err := ParseURL("https://account.file.core.windows.net/share/directory/file")
	require.NoError(t, err)
	require.Empty(t, fileURLParts.FileID)
	require.Equal(t, "directory/file", fileURLParts.DirectoryOrFilePath)
	require.Equal(t, "https://account.file.core.windows.net/share/directory/file", fileURLParts.String())
}

func TestParseURLFileIDAndSnapshot(t *testing.T) {
	const u = "https://account.file.core.windows.net/share?fileid=12384898975283830&sharesnapshot=2011-03-09T01:42:34.9360000Z"
	fileURLParts, err := ParseURL(u)
	require.NoError(t, err)
	require.Equal(t, "share", fileURLParts.ShareName)
	require.Empty(t, fileURLParts.DirectoryOrFilePath)
	require.Equal(t, "12384898975283830", fileURLParts.FileID)
	require.Equal(t, "2011-03-09T01:42:34.9360000Z", fileURLParts.ShareSnapshot)
	require.Empty(t, fileURLParts.UnparsedParams)
	require.Equal(t, u, fileURLParts.String())
}

func TestParseURLFileIDAndSAS(t *testing.T) {
	const sasQuery = "se=2015-04-30T02%3A23%3A26Z&sig=Z%2FRHIX5Xcg0Mq2rqI3OlWTjEg2tYkboXr1P9ZUXDtkk%3D&sip=168.1.5.60-168.1.5.70&sp=rw&spr=https&sr=b&st=2015-04-29T22%3A18%3A26Z&sv=2015-04-05"
	fileURLParts, err := ParseURL("https://account.file.core.windows.net/share?fileid=12384898975283830&" + sasQuery)
	require.NoError(t, err)
	require.Equal(t, "share", fileURLParts.ShareName)
	require.Empty(t, fileURLParts.DirectoryOrFilePath)
	require.Equal(t, "12384898975283830", fileURLParts.FileID)
	require.Empty(t, fileURLParts.UnparsedParams)
	require.Equal(t, "rw", fileURLParts.SAS.Permissions())
	require.Equal(t, "Z/RHIX5Xcg0Mq2rqI3OlWTjEg2tYkboXr1P9ZUXDtkk=", fileURLParts.SAS.Signature())
	require.Equal(t, "2015-04-05", fileURLParts.SAS.Version())

	// the file ID is emitted before the SAS
	require.Equal(t, "https://account.file.core.windows.net/share?fileid=12384898975283830&"+sasQuery, fileURLParts.String())
}

func TestURLPartsSetFileID(t *testing.T) {
	fileURLParts, err := ParseURL("https://account.file.core.windows.net/share")
	require.NoError(t, err)

	fileURLParts.FileID = "12384898975283830"
	require.Equal(t, "https://account.file.core.windows.net/share?fileid=12384898975283830", fileURLParts.String())

	fileURLParts.FileID = "99999999999999999"
	require.Equal(t, "https://account.file.core.windows.net/share?fileid=99999999999999999", fileURLParts.String())

	// clearing the file ID removes it from the URL
	fileURLParts.FileID = ""
	require.Equal(t, "https://account.file.core.windows.net/share", fileURLParts.String())
}

func TestURLPartsFileIDIsEscaped(t *testing.T) {
	fileURLParts, err := ParseURL("https://account.file.core.windows.net/share")
	require.NoError(t, err)

	fileURLParts.FileID = "a&b=c"
	u := fileURLParts.String()
	require.Equal(t, "https://account.file.core.windows.net/share?fileid=a%26b%3Dc", u)

	roundTrip, err := ParseURL(u)
	require.NoError(t, err)
	require.Equal(t, "a&b=c", roundTrip.FileID)
	require.Empty(t, roundTrip.UnparsedParams)
}
