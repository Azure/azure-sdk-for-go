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
	queueURLParts, err := ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, queueURLParts.Scheme, "https")
	require.Equal(t, queueURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, queueURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")

	urlWithIP = "https://127.0.0.1:5000/fakestorageaccount/fakequeue"
	queueURLParts, err = ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, queueURLParts.Scheme, "https")
	require.Equal(t, queueURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, queueURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")
	require.Equal(t, queueURLParts.QueueName, "fakequeue")

	urlWithIP = "https://127.0.0.1:5000/fakestorageaccount/fakequeue"
	queueURLParts, err = ParseURL(urlWithIP)
	require.NoError(t, err)
	require.Equal(t, queueURLParts.Scheme, "https")
	require.Equal(t, queueURLParts.Host, "127.0.0.1:5000")
	require.Equal(t, queueURLParts.IPEndpointStyleInfo.AccountName, "fakestorageaccount")
	require.Equal(t, queueURLParts.QueueName, "fakequeue")
}

func TestParseURL(t *testing.T) {
	testStorageAccount := "fakestorageaccount"
	host := fmt.Sprintf("%s.queue.core.windows.net", testStorageAccount)
	testQueue := "fakequeue"

	const sasStr = "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"

	urlWithVersion := fmt.Sprintf("https://%s.queue.core.windows.net/%s?%s", testStorageAccount, testQueue, sasStr)
	queueURLParts, err := ParseURL(urlWithVersion)
	require.NoError(t, err)

	require.Equal(t, queueURLParts.Scheme, "https")
	require.Equal(t, queueURLParts.Host, host)
	require.Equal(t, queueURLParts.QueueName, testQueue)
	validateSAS(t, sasStr, queueURLParts.SAS)
}
