//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"github.com/stretchr/testify/require"
)

func (s *azblobTestSuite) TestURLParsing() {
	_require := require.New(s.T())
	testStorageAccount := "fakestorageaccount"
	host := fmt.Sprintf("%s.blob.core.windows.net", testStorageAccount)
	testContainer := "fakecontainer"
	fileNames := []string{"/._.TESTT.txt", "/.gitignore/dummyfile1"}

	sas := "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"

	for _, fileName := range fileNames {
		snapshotID, versionID := "", "2021-10-25T05:41:32.5526810Z"
		sasWithVersionID := "?versionId=" + versionID + "&" + sas
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s%s", testStorageAccount, testContainer, fileName, sasWithVersionID)
		blobURLParts, err := azblob.ParseBlobURL(urlWithVersion)
		_require.Nil(err)

		_require.Equal(blobURLParts.Scheme, "https")
		_require.Equal(blobURLParts.Host, host)
		_require.Equal(blobURLParts.ContainerName, testContainer)
		_require.Equal(blobURLParts.VersionID, versionID)
		_require.Equal(blobURLParts.Snapshot, snapshotID)

		validateSAS(_require, sas, blobURLParts.SAS)
	}

	for _, fileName := range fileNames {
		snapshotID, versionID := "2011-03-09T01:42:34Z", ""
		sasWithSnapshotID := "?snapshot=" + snapshotID + "&" + sas
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s%s", testStorageAccount, testContainer, fileName, sasWithSnapshotID)
		blobURLParts, err := azblob.ParseBlobURL(urlWithVersion)
		_require.Nil(err)

		_require.Equal(blobURLParts.Scheme, "https")
		_require.Equal(blobURLParts.Host, host)
		_require.Equal(blobURLParts.ContainerName, testContainer)
		_require.Equal(blobURLParts.VersionID, versionID)
		_require.Equal(blobURLParts.Snapshot, snapshotID)

		validateSAS(_require, sas, blobURLParts.SAS)
	}

	//urlWithIP := "https://127.0.0.1:5000/"

}
