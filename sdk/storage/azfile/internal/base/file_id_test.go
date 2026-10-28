// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/stretchr/testify/require"
)

func TestCheckAddressing(t *testing.T) {
	const (
		pathURL       = "https://account.file.core.windows.net/share/dir/file"
		fileIDURL     = "https://account.file.core.windows.net/share?fileid=123"
		pathAndFileID = "https://account.file.core.windows.net/share/dir/file?fileid=123"
	)

	for _, tc := range []struct {
		name       string
		url        string
		addressing Addressing
		wantErr    error
	}{
		{"path/PathOnly", pathURL, PathOnly, nil},
		{"path/PathOrFileID", pathURL, PathOrFileID, nil},
		{"path/FileIDOnly", pathURL, FileIDOnly, fileerror.FileIDRequired},
		{"fileID/PathOnly", fileIDURL, PathOnly, fileerror.FileIDNotSupported},
		{"fileID/PathOrFileID", fileIDURL, PathOrFileID, nil},
		{"fileID/FileIDOnly", fileIDURL, FileIDOnly, nil},
		{"pathAndFileID/PathOnly", pathAndFileID, PathOnly, fileerror.FileIDNotSupported},
		{"pathAndFileID/PathOrFileID", pathAndFileID, PathOrFileID, fileerror.FileIDNotSupported},
		{"pathAndFileID/FileIDOnly", pathAndFileID, FileIDOnly, fileerror.FileIDNotSupported},
		// an unparsable URL is left for the request to report, as before file ID support
		{"invalid/PathOnly", "://not a url", PathOnly, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckAddressing(tc.url, "file.Client.Op", tc.addressing)
			if tc.wantErr == nil {
				require.NoError(t, err)
				return
			}
			require.ErrorIs(t, err, tc.wantErr)
			require.Contains(t, err.Error(), "file.Client.Op")
		})
	}
}
