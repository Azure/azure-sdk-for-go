// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/stretchr/testify/require"
)

func TestNewResponseError(t *testing.T) {
	fakeURL, err := url.Parse("https://contoso.com")
	require.NoError(t, err)
	err = NewResponseError(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{ "code": "ErrorItsBroken", "message": "it's not working" }`)),
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	})
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, http.StatusInternalServerError, respErr.StatusCode)
	require.EqualValues(t, "ErrorItsBroken", respErr.ErrorCode)
	require.NotNil(t, respErr.RawResponse)
}

func TestNewResponseErrorWithErrorCode(t *testing.T) {
	fakeURL, err := url.Parse("https://contoso.com")
	require.NoError(t, err)
	err = NewResponseErrorWithErrorCode(&http.Response{
		Status:     "the system is down",
		StatusCode: http.StatusInternalServerError,
		Request: &http.Request{
			Method: http.MethodGet,
			URL:    fakeURL,
		},
	}, "ErrorItsBroken")
	var respErr *exported.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.EqualValues(t, http.StatusInternalServerError, respErr.StatusCode)
	require.EqualValues(t, "ErrorItsBroken", respErr.ErrorCode)
	require.NotNil(t, respErr.RawResponse)
}
