// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

type featureFlagPagerTransport struct {
	headers []http.Header
}

func (t *featureFlagPagerTransport) Do(req *http.Request) (*http.Response, error) {
	t.headers = append(t.headers, req.Header.Clone())

	header := http.Header{
		"Content-Type": []string{"application/json"},
		"ETag":         []string{`"response-page"`},
	}
	if len(t.headers) == 1 {
		header.Set("Link", `</ff?After=next>; rel="next"`)
	}

	return &http.Response{
		Request:    req,
		StatusCode: http.StatusNotModified,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader("{}")),
	}, nil
}

func TestNewGetFeatureFlagsPagerWithLinkHeaderMatchConditions(t *testing.T) {
	transport := &featureFlagPagerTransport{}
	client, err := azcore.NewClient("azappconfig", "v0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Transport: transport,
	})
	require.NoError(t, err)

	featureFlagClient := NewAzureAppConfigurationClient("https://example.azconfig.io", client).NewAzureAppConfigurationFeatureFlagClient()
	pager := featureFlagClient.NewGetFeatureFlagsPagerWithLinkHeader([]azcore.MatchConditions{
		{IfNoneMatch: to.Ptr(azcore.ETag(`"page-1"`))},
		{IfNoneMatch: to.Ptr(azcore.ETag(`"page-2"`))},
	}, &AzureAppConfigurationFeatureFlagClientGetFeatureFlagsOptions{})

	page, err := pager.NextPage(context.Background())
	require.NoError(t, err)
	require.NotNil(t, page.NextLink)
	require.NotNil(t, page.ETag)
	require.Equal(t, `"response-page"`, *page.ETag)

	_, err = pager.NextPage(context.Background())
	require.NoError(t, err)
	require.Len(t, transport.headers, 2)
	require.Equal(t, `"page-1"`, transport.headers[0].Get("If-None-Match"))
	require.Equal(t, `"page-2"`, transport.headers[1].Get("If-None-Match"))
}