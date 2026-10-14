// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func newTestResponse(statusCode int, body string) *http.Response {
	header := http.Header{"Content-Type": []string{"application/xml"}}
	if statusCode != http.StatusOK {
		header.Set("x-ms-error-code", "ResourceNotFound")
	}
	req, _ := http.NewRequest(http.MethodGet, "https://account.file.core.windows.net/share", nil)
	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

const errorBody = `<?xml version="1.0" encoding="utf-8"?><Error><Code>ResourceNotFound</Code><Message>The specified resource does not exist.</Message></Error>`

// The exported HandleResponse wrappers are used by the hand-written pagers. They must accept 200 OK and reject other
// status codes; the generated handlers they wrap only accept the status codes passed to them.
func TestListHandleResponseWrappers(t *testing.T) {
	testCases := []struct {
		name   string
		body   string
		handle func(resp *http.Response) error
	}{
		{
			name: "ListFilesAndDirectoriesSegment",
			body: `<?xml version="1.0" encoding="utf-8"?><EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/" ShareName="share" DirectoryPath="dir"><Entries><File><Name>file.txt</Name><Properties><Content-Length>5</Content-Length></Properties></File></Entries><NextMarker /></EnumerationResults>`,
			handle: func(resp *http.Response) error {
				result, err := (&DirectoryClient{}).ListFilesAndDirectoriesSegmentHandleResponse(resp)
				if err == nil {
					require.Len(t, result.Segment.Files, 1)
					require.Equal(t, "file.txt", *result.Segment.Files[0].Name)
				}
				return err
			},
		},
		{
			name: "ListSharesSegment",
			body: `<?xml version="1.0" encoding="utf-8"?><EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/"><Shares><Share><Name>share1</Name><Properties><Last-Modified>Mon, 11 May 2026 04:19:33 GMT</Last-Modified><Etag>"0x8DEAF1480FD31EF"</Etag><Quota>100</Quota></Properties></Share></Shares><NextMarker /></EnumerationResults>`,
			handle: func(resp *http.Response) error {
				result, err := (&ServiceClient{}).ListSharesSegmentHandleResponse(resp)
				if err == nil {
					require.Len(t, result.Shares, 1)
					require.Equal(t, "share1", *result.Shares[0].Name)
				}
				return err
			},
		},
		{
			name: "ListAllRanges",
			body: `<?xml version="1.0" encoding="utf-8"?><Ranges><Range><Start>0</Start><End>511</End></Range><NextMarker /></Ranges>`,
			handle: func(resp *http.Response) error {
				result, err := (&FileClient{}).ListAllRangesHandleResponse(resp)
				if err == nil {
					require.Len(t, result.Ranges, 1)
					require.Equal(t, int64(511), *result.Ranges[0].End)
				}
				return err
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"/OK", func(t *testing.T) {
			require.NoError(t, tc.handle(newTestResponse(http.StatusOK, tc.body)))
		})
		for _, statusCode := range []int{http.StatusCreated, http.StatusNotFound, http.StatusInternalServerError} {
			t.Run(tc.name+"/"+http.StatusText(statusCode), func(t *testing.T) {
				err := tc.handle(newTestResponse(statusCode, errorBody))
				var respErr *azcore.ResponseError
				require.True(t, errors.As(err, &respErr), "expected *azcore.ResponseError, got %v", err)
				require.Equal(t, statusCode, respErr.StatusCode)
			})
		}
	}
}
