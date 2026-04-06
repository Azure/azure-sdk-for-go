// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

func NewAzureAppConfigurationClient(endpoint string, client *azcore.Client) *AzureAppConfigurationClient {
	return &AzureAppConfigurationClient{
		internal: client,
		endpoint: endpoint,
	}
}

func (a *AzureAppConfigurationClient) Pipeline() runtime.Pipeline {
	return a.internal.Pipeline()
}

func (a *AzureAppConfigurationClient) Tracer() tracing.Tracer {
	return a.internal.Tracer()
}

// copy of NewGetKeyValuesPager with slice of match conditions and other tweaks
func (client *AzureAppConfigurationClient) NewGetKeyValuesPagerWithMatchConditions(matchConditions []azcore.MatchConditions, options *AzureAppConfigurationClientGetKeyValuesOptions) *runtime.Pager[AzureAppConfigurationClientGetKeyValuesResponse] {
	return runtime.NewPager(runtime.PagingHandler[AzureAppConfigurationClientGetKeyValuesResponse]{
		More: func(page AzureAppConfigurationClientGetKeyValuesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AzureAppConfigurationClientGetKeyValuesResponse) (AzureAppConfigurationClientGetKeyValuesResponse, error) {
			curCondition := azcore.MatchConditions{}
			if len(matchConditions) > 0 {
				curCondition = matchConditions[0]
				matchConditions = matchConditions[1:]
			}
			options.IfMatch = (*string)(curCondition.IfMatch)
			options.IfNoneMatch = (*string)(curCondition.IfNoneMatch)
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.getKeyValuesCreateRequest(ctx, options)
			}, &runtime.FetcherForNextLinkOptions{
				NextReq: func(ctx context.Context, encodedNextLink string) (*policy.Request, error) {
					return client.getNextPageCreateRequestWithMatchConditions(ctx, encodedNextLink, curCondition)
				},
				StatusCodes: []int{http.StatusNotModified},
			})
			if err != nil {
				return AzureAppConfigurationClientGetKeyValuesResponse{}, err
			}
			return client.getKeyValuesHandleResponseWithLinkHeader(resp)
		},
	})
}

// adds match conditions to the request created in getNextPageCreateRequest
func (a *AzureAppConfigurationClient) getNextPageCreateRequestWithMatchConditions(ctx context.Context, nextLink string, matchConditions azcore.MatchConditions) (*policy.Request, error) {
	req, err := a.getNextPageCreateRequest(ctx, nextLink)
	if err != nil {
		return nil, err
	}
	if matchConditions.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{*(*string)(matchConditions.IfMatch)}
	}
	if matchConditions.IfNoneMatch != nil {
		req.Raw().Header["If-None-Match"] = []string{*(*string)(matchConditions.IfNoneMatch)}
	}
	return req, nil
}

// parses the nextLink URL from the Link response header
func (a *AzureAppConfigurationClient) getKeyValuesHandleResponseWithLinkHeader(resp *http.Response) (AzureAppConfigurationClientGetKeyValuesResponse, error) {
	result, err := a.getKeyValuesHandleResponse(resp)
	if err != nil {
		return AzureAppConfigurationClientGetKeyValuesResponse{}, err
	}
	if result.NextLink != nil {
		return result, err
	}

	link := resp.Header.Get("Link")
	if link == "" {
		return result, err
	}

	// the link header format is <nextLinkURL>; rel="next"
	// extract the values between < and >
	if endIndex := strings.Index(link, ">"); endIndex > 0 {
		result.NextLink = to.Ptr(link[1:endIndex])
	}
	return result, err
}
