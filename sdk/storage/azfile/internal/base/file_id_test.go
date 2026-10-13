// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/stretchr/testify/require"
)

func TestFileID(t *testing.T) {
	require.Equal(t, "123", FileID("https://account.file.core.windows.net/share?fileid=123"))
	require.Equal(t, "123", FileID("https://account.file.core.windows.net/share?sharesnapshot=2011-03-09T01:42:34.9360000Z&FileId=123"))
	require.Empty(t, FileID("https://account.file.core.windows.net/share/dir/file"))
	require.Empty(t, FileID("https://account.file.core.windows.net/share"))
	require.Empty(t, FileID("://not a url"))
}

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
		{"path/FileIDOnly", pathURL, FileIDOnly, fileerror.ErrFileIDRequired},
		{"fileID/PathOnly", fileIDURL, PathOnly, fileerror.ErrFileIDNotSupported},
		{"fileID/PathOrFileID", fileIDURL, PathOrFileID, nil},
		{"fileID/FileIDOnly", fileIDURL, FileIDOnly, nil},
		{"pathAndFileID/PathOnly", pathAndFileID, PathOnly, fileerror.ErrFileIDNotSupported},
		{"pathAndFileID/PathOrFileID", pathAndFileID, PathOrFileID, fileerror.ErrFileIDNotSupported},
		{"pathAndFileID/FileIDOnly", pathAndFileID, FileIDOnly, fileerror.ErrFileIDNotSupported},
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
