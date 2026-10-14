// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"
	"github.com/stretchr/testify/require"
)

func newMockServiceClient(t *testing.T, responses ...testcommon.MockResponse) (*service.Client, *testcommon.MockTransport) {
	transport := testcommon.NewMockTransport(responses...)
	client, err := service.NewClientWithNoCredential("https://account.file.core.windows.net/", &service.ClientOptions{
		ClientOptions: testcommon.MockClientOptions(transport),
	})
	require.NoError(t, err)
	return client, transport
}

func TestListSharesPagerMocked(t *testing.T) {
	_require := require.New(t)
	client, transport := newMockServiceClient(t,
		testcommon.MockResponse{StatusCode: http.StatusOK, Body: `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/">
	<Shares>
		<Share><Name>share1</Name><Properties><Last-Modified>Mon, 11 May 2026 04:19:33 GMT</Last-Modified><Etag>"0x8DEAF1480FD31EF"</Etag><Quota>100</Quota></Properties></Share>
	</Shares>
	<NextMarker>marker1</NextMarker>
</EnumerationResults>`},
		testcommon.MockResponse{StatusCode: http.StatusOK, Body: `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/">
	<Shares>
		<Share><Name>share2</Name><Properties><Last-Modified>Mon, 11 May 2026 04:19:33 GMT</Last-Modified><Etag>"0x8DEAF1480FD31F0"</Etag><Quota>200</Quota></Properties></Share>
	</Shares>
	<NextMarker />
</EnumerationResults>`},
	)

	var names []string
	pager := client.NewListSharesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, s := range page.Shares {
			names = append(names, *s.Name)
		}
	}
	_require.Equal([]string{"share1", "share2"}, names)
	_require.Len(transport.Requests, 2)
	_require.Equal("marker1", transport.Requests[1].URL.Query().Get("marker"))
}

func TestListSharesPagerErrorMocked(t *testing.T) {
	_require := require.New(t)
	client, _ := newMockServiceClient(t, testcommon.MockResponse{
		StatusCode: http.StatusForbidden,
		ErrorCode:  string(fileerror.AuthorizationFailure),
		Body:       `<?xml version="1.0" encoding="utf-8"?><Error><Code>AuthorizationFailure</Code><Message>This request is not authorized to perform this operation.</Message></Error>`,
	})

	_, err := client.NewListSharesPager(nil).NextPage(context.Background())
	var respErr *azcore.ResponseError
	_require.True(errors.As(err, &respErr))
	_require.Equal(http.StatusForbidden, respErr.StatusCode)
	_require.True(fileerror.HasCode(err, fileerror.AuthorizationFailure))
}
