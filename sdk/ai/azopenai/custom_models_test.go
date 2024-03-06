//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestParseResponseError(t *testing.T) {
	bodyBytes, err := os.ReadFile("testdata/content_filter_response_error.json")
	require.NoError(t, err)

	buff := bytes.NewBuffer(bodyBytes)

	fakeURL, err := url.Parse("https://openai-something.microsoft.com")
	require.NoError(t, err)

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(buff),
		Request: &http.Request{
			Method: "POST",
			URL:    fakeURL,
		},
	}

	err = newContentFilterResponseError(resp)

	// this is the outer error, which is the standard Azure response error.
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusBadRequest, respErr.StatusCode)
	require.Equal(t, "content_filter", respErr.ErrorCode)

	// Azure also returns error information when content filtering happens.
	var contentFilterErr *ContentFilterResponseError
	require.ErrorAs(t, err, &contentFilterErr)

	// we're still a response error
	require.Equal(t, http.StatusBadRequest, respErr.StatusCode)
	require.Equal(t, "content_filter", respErr.ErrorCode)

	contentFilterResults := contentFilterErr.ContentFilterResults

	// this comment was considered violent, so it was filtered.
	require.Equal(t, &ContentFilterResult{
		Filtered: to.Ptr(true),
		Severity: to.Ptr(ContentFilterSeverityMedium)}, contentFilterResults.Violence)

	require.Equal(t, &ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(ContentFilterSeveritySafe)}, contentFilterResults.Hate)
	require.Equal(t, &ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(ContentFilterSeveritySafe)}, contentFilterResults.SelfHarm)
	require.Equal(t, &ContentFilterResult{Filtered: to.Ptr(false), Severity: to.Ptr(ContentFilterSeveritySafe)}, contentFilterResults.Sexual)

	require.NotNil(t, contentFilterResults)
}
