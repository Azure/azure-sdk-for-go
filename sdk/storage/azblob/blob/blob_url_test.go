//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	testStorageAccount := "fakestorageaccount"
	host := fmt.Sprintf("%s.blob.core.windows.net", testStorageAccount)
	testContainer := "fakecontainer"
	fileNames := []string{"/._.TESTT.txt", "/.gitignore/dummyfile1"}

	const sasStr = "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"

	for _, fileName := range fileNames {
		snapshotID, versionID := "", "2021-10-25T05:41:32.5526810Z"
		sasWithVersionID := "?versionId=" + versionID + "&" + sasStr
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s%s", testStorageAccount, testContainer, fileName, sasWithVersionID)
		blobURLParts, err := ParseURL(urlWithVersion)
		require.NoError(t, err)

		require.Equal(t, blobURLParts.Scheme, "https")
		require.Equal(t, blobURLParts.Host, host)
		require.Equal(t, blobURLParts.ContainerName, testContainer)
		require.Equal(t, blobURLParts.VersionID, versionID)
		require.Equal(t, blobURLParts.Snapshot, snapshotID)

		validateSAS(t, sasStr, blobURLParts.SAS)
	}

	for _, fileName := range fileNames {
		snapshotID, versionID := "2011-03-09T01:42:34Z", ""
		sasWithSnapshotID := "?snapshot=" + snapshotID + "&" + sasStr
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s%s", testStorageAccount, testContainer, fileName, sasWithSnapshotID)
		blobURLParts, err := ParseURL(urlWithVersion)
		require.NoError(t, err)

		require.Equal(t, blobURLParts.Scheme, "https")
		require.Equal(t, blobURLParts.Host, host)
		require.Equal(t, blobURLParts.ContainerName, testContainer)
		require.Equal(t, blobURLParts.VersionID, versionID)
		require.Equal(t, blobURLParts.Snapshot, snapshotID)

		validateSAS(t, sasStr, blobURLParts.SAS)
	}

	//urlWithIP := "https://127.0.0.1:5000/"
}

func validateSAS(t *testing.T, sasStr string, parameters sas.QueryParameters) {
	sasCompMap := make(map[string]string)
	for _, sasComp := range strings.Split(sasStr, "&") {
		comp := strings.Split(sasComp, "=")
		sasCompMap[comp[0]] = comp[1]
	}

	require.Equal(t, parameters.Version(), sasCompMap["sv"])
	require.Equal(t, parameters.Services(), sasCompMap["ss"])
	require.Equal(t, parameters.ResourceTypes(), sasCompMap["srt"])
	require.Equal(t, string(parameters.Protocol()), sasCompMap["spr"])
	if _, ok := sasCompMap["st"]; ok {
		startTime, _, err := sas.ParseTime(sasCompMap["st"])
		require.NoError(t, err)
		require.Equal(t, parameters.StartTime(), startTime)
	}
	if _, ok := sasCompMap["se"]; ok {
		endTime, _, err := sas.ParseTime(sasCompMap["se"])
		require.NoError(t, err)
		require.Equal(t, parameters.ExpiryTime(), endTime)
	}

	if _, ok := sasCompMap["snapshot"]; ok {
		snapshotTime, _, err := sas.ParseTime(sasCompMap["snapshot"])
		require.NoError(t, err)
		require.Equal(t, parameters.SnapshotTime(), snapshotTime)
	}
	ipRange := parameters.IPRange()
	require.Equal(t, ipRange.String(), sasCompMap["sip"])
	require.Equal(t, parameters.Identifier(), sasCompMap["si"])
	require.Equal(t, parameters.Resource(), sasCompMap["sr"])
	require.Equal(t, parameters.Permissions(), sasCompMap["sp"])

	sign, err := url.QueryUnescape(sasCompMap["sig"])
	require.NoError(t, err)

	require.Equal(t, parameters.Signature(), sign)
	require.Equal(t, parameters.CacheControl(), sasCompMap["rscc"])
	require.Equal(t, parameters.ContentDisposition(), sasCompMap["rscd"])
	require.Equal(t, parameters.ContentEncoding(), sasCompMap["rsce"])
	require.Equal(t, parameters.ContentLanguage(), sasCompMap["rscl"])
	require.Equal(t, parameters.ContentType(), sasCompMap["rsct"])
	require.Equal(t, parameters.SignedOID(), sasCompMap["skoid"])
	require.Equal(t, parameters.SignedTID(), sasCompMap["sktid"])

	if _, ok := sasCompMap["skt"]; ok {
		signedStart, _, err := sas.ParseTime(sasCompMap["skt"])
		require.NoError(t, err)
		require.Equal(t, parameters.SignedStart(), signedStart)
	}
	require.Equal(t, parameters.SignedService(), sasCompMap["sks"])

	if _, ok := sasCompMap["ske"]; ok {
		signedExpiry, _, err := sas.ParseTime(sasCompMap["ske"])
		require.NoError(t, err)
		require.Equal(t, parameters.SignedExpiry(), signedExpiry)
	}

	require.Equal(t, parameters.SignedVersion(), sasCompMap["skv"])
	require.Equal(t, parameters.SignedDirectoryDepth(), sasCompMap["sdd"])
	require.Equal(t, parameters.PreauthorizedAgentObjectID(), sasCompMap["saoid"])
	require.Equal(t, parameters.AgentObjectID(), sasCompMap["suoid"])
	require.Equal(t, parameters.SignedCorrelationID(), sasCompMap["scid"])
}
