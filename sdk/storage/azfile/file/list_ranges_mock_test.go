// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/stretchr/testify/require"
)

func newMockFileClient(t *testing.T, responses ...testcommon.MockResponse) (*file.Client, *testcommon.MockTransport) {
	transport := testcommon.NewMockTransport(responses...)
	client, err := file.NewClientWithNoCredential("https://account.file.core.windows.net/share/dir/file.txt", &file.ClientOptions{
		ClientOptions: testcommon.MockClientOptions(transport),
	})
	require.NoError(t, err)
	return client, transport
}

func TestListRangesPagerMocked(t *testing.T) {
	_require := require.New(t)
	client, transport := newMockFileClient(t,
		testcommon.MockResponse{StatusCode: http.StatusOK, Body: `<?xml version="1.0" encoding="utf-8"?>
<Ranges>
	<Range><Start>0</Start><End>511</End></Range>
	<ClearRange><Start>512</Start><End>1023</End></ClearRange>
	<NextMarker>marker1</NextMarker>
</Ranges>`},
		testcommon.MockResponse{StatusCode: http.StatusOK, Body: `<?xml version="1.0" encoding="utf-8"?>
<Ranges>
	<Range><Start>1024</Start><End>2047</End></Range>
	<NextMarker />
</Ranges>`},
	)

	pager := client.NewListRangesPager(&file.ListRangesOptions{MaxResults: to.Ptr(int32(1))})

	_require.True(pager.More())
	page, err := pager.NextPage(context.Background())
	_require.NoError(err)
	_require.Len(page.Ranges, 1)
	_require.Equal(int64(0), *page.Ranges[0].Start)
	_require.Equal(int64(511), *page.Ranges[0].End)
	_require.Len(page.ClearRanges, 1)
	_require.Equal(int64(512), *page.ClearRanges[0].Start)
	_require.Equal(int64(1023), *page.ClearRanges[0].End)
	_require.Equal("marker1", *page.NextMarker)

	_require.True(pager.More())
	page, err = pager.NextPage(context.Background())
	_require.NoError(err)
	_require.Len(page.Ranges, 1)
	_require.Equal(int64(1024), *page.Ranges[0].Start)
	_require.Equal(int64(2047), *page.Ranges[0].End)
	_require.Empty(page.ClearRanges)
	_require.False(pager.More())

	_require.Len(transport.Requests, 2)
	_require.Equal("1", transport.Requests[0].URL.Query().Get("maxresults"))
	_require.Empty(transport.Requests[0].URL.Query().Get("marker"))
	_require.Equal("marker1", transport.Requests[1].URL.Query().Get("marker"))
}

func TestListRangesPagerErrorMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockFileClient(t, testcommon.MockResponse{
		StatusCode: http.StatusNotFound,
		ErrorCode:  string(fileerror.ResourceNotFound),
		Body:       `<?xml version="1.0" encoding="utf-8"?><Error><Code>ResourceNotFound</Code><Message>The specified resource does not exist.</Message></Error>`,
	})

	_, err := client.NewListRangesPager(nil).NextPage(context.Background())
	var respErr *azcore.ResponseError
	_require.True(errors.As(err, &respErr))
	_require.Equal(http.StatusNotFound, respErr.StatusCode)
	_require.True(fileerror.HasCode(err, fileerror.ResourceNotFound))
}
