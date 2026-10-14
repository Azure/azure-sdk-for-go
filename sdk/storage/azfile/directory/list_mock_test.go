// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/stretchr/testify/require"
)

func newMockDirectoryClient(t *testing.T, responses ...testcommon.MockResponse) (*directory.Client, *testcommon.MockTransport) {
	transport := testcommon.NewMockTransport(responses...)
	client, err := directory.NewClientWithNoCredential("https://account.file.core.windows.net/share/dir", &directory.ClientOptions{
		ClientOptions: testcommon.MockClientOptions(transport),
	})
	require.NoError(t, err)
	return client, transport
}

func TestListFilesAndDirectoriesPagerMocked(t *testing.T) {
	_require := require.New(t)
	client, transport := newMockDirectoryClient(t,
		testcommon.MockResponse{StatusCode: http.StatusOK, Body: `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="share" DirectoryPath="dir">
	<Entries>
		<File><Name>file.txt</Name><Properties><Content-Length>5</Content-Length></Properties></File>
	</Entries>
	<NextMarker>marker1</NextMarker>
</EnumerationResults>`},
		testcommon.MockResponse{StatusCode: http.StatusOK, Body: `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="share" DirectoryPath="dir">
	<Entries>
		<Directory><Name>subdir</Name><Properties /></Directory>
	</Entries>
	<NextMarker />
</EnumerationResults>`},
	)

	var files, dirs []string
	pager := client.NewListFilesAndDirectoriesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, f := range page.Segment.Files {
			files = append(files, *f.Name)
		}
		for _, d := range page.Segment.Directories {
			dirs = append(dirs, *d.Name)
		}
	}
	_require.Equal([]string{"file.txt"}, files)
	_require.Equal([]string{"subdir"}, dirs)
	_require.Len(transport.Requests, 2)
	_require.Equal("marker1", transport.Requests[1].URL.Query().Get("marker"))
}

func TestListFilesAndDirectoriesPagerErrorMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockDirectoryClient(t, testcommon.MockResponse{
		StatusCode: http.StatusNotFound,
		ErrorCode:  string(fileerror.ParentNotFound),
		Body:       `<?xml version="1.0" encoding="utf-8"?><Error><Code>ParentNotFound</Code><Message>The specified parent path does not exist.</Message></Error>`,
	})

	_, err := client.NewListFilesAndDirectoriesPager(nil).NextPage(context.Background())
	var respErr *azcore.ResponseError
	_require.True(errors.As(err, &respErr))
	_require.Equal(http.StatusNotFound, respErr.StatusCode)
	_require.True(fileerror.HasCode(err, fileerror.ParentNotFound))
}
