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

	if nextLink := nextLinkFromHeader(resp.Header.Get("Link")); nextLink != "" {
		result.NextLink = to.Ptr(nextLink)
	}
	return result, err
}

// CheckKeyValuesPagerResponse is a custom response type for paginated HEAD requests.
// It wraps the generated CheckKeyValuesResponse and adds NextLink for pagination support.
type CheckKeyValuesPagerResponse struct {
	AzureAppConfigurationClientCheckKeyValuesResponse

	// NextLink contains the URL for the next page, parsed from the Link header response.
	NextLink *string
}

// NewCheckKeyValuesPagerWithMatchConditions uses HEAD requests (CheckKeyValues) for bandwidth optimization.
// Returns only ETag and SyncToken headers without response body, useful for detecting changes during watch/refresh.
func (client *AzureAppConfigurationClient) NewCheckKeyValuesPagerWithMatchConditions(matchConditions []azcore.MatchConditions, options *AzureAppConfigurationClientCheckKeyValuesOptions) *runtime.Pager[CheckKeyValuesPagerResponse] {
	return runtime.NewPager(runtime.PagingHandler[CheckKeyValuesPagerResponse]{
		More: func(page CheckKeyValuesPagerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CheckKeyValuesPagerResponse) (CheckKeyValuesPagerResponse, error) {
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
				return client.checkKeyValuesCreateRequest(ctx, options)
			}, &runtime.FetcherForNextLinkOptions{
				NextReq: func(ctx context.Context, encodedNextLink string) (*policy.Request, error) {
					return client.checkNextPageCreateRequestWithMatchConditions(ctx, encodedNextLink, curCondition)
				},
				StatusCodes: []int{http.StatusNotModified},
			})
			if err != nil {
				return CheckKeyValuesPagerResponse{}, err
			}
			return client.checkKeyValuesHandleResponseWithLinkHeader(resp)
		},
	})
}

// adds match conditions to the HEAD request for the next page
func (a *AzureAppConfigurationClient) checkNextPageCreateRequestWithMatchConditions(ctx context.Context, nextLink string, matchConditions azcore.MatchConditions) (*policy.Request, error) {
	urlPath := nextLink
	req, err := runtime.NewRequest(ctx, http.MethodHead, runtime.JoinPaths(a.endpoint, urlPath))
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

// parses the nextLink URL from the Link response header for HEAD requests
func (a *AzureAppConfigurationClient) checkKeyValuesHandleResponseWithLinkHeader(resp *http.Response) (CheckKeyValuesPagerResponse, error) {
	genResult, err := a.checkKeyValuesHandleResponse(resp)
	if err != nil {
		return CheckKeyValuesPagerResponse{}, err
	}

	result := CheckKeyValuesPagerResponse{
		AzureAppConfigurationClientCheckKeyValuesResponse: genResult,
	}

	link := resp.Header.Get("Link")
	if link == "" {
		return result, nil
	}

	// the link header format is <nextLinkURL>; rel="next"
	// extract the values between < and >
	if endIndex := strings.Index(link, ">"); endIndex > 0 {
		result.NextLink = to.Ptr(link[1:endIndex])
	}
	return result, nil
}

// CreateSnapshot exports internal createSnapshot
func (a *AzureAppConfigurationClient) CreateSnapshot(ctx context.Context, contentType CreateSnapshotRequestContentType, name string, entity Snapshot, options *AzureAppConfigurationClientBeginCreateSnapshotOptions) (*http.Response, error) {
	return a.createSnapshot(ctx, contentType, name, entity, options)
}

// getNextPageCreateRequest creates the getNextPageCreateRequest request.
func (client *AzureAppConfigurationClient) getNextPageCreateRequest(ctx context.Context, nextLink string) (*policy.Request, error) {
	urlPath := nextLink
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// ffAcceptHeader matches the Accept header used by the generated feature flag pager requests.
const ffAcceptHeader = "application/json;profile=\"https://azconfig.io/mime-profiles/ffset\";charset=utf-8, application/problem+json"

// NewGetFeatureFlagsPagerWithLinkHeader wraps the generated feature flag pager so that the nextLink
// can also be discovered from the "Link" response header. Azure App Configuration returns
// pagination information via the Link header rather than the JSON body for feature flag list
// responses.
func (client *AzureAppConfigurationFeatureFlagClient) NewGetFeatureFlagsPagerWithLinkHeader(options *AzureAppConfigurationFeatureFlagClientGetFeatureFlagsOptions) *runtime.Pager[AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse] {
	return runtime.NewPager(runtime.PagingHandler[AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse]{
		More: func(page AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse) (AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse, error) {
			nextLink := ""
			if page != nil && page.NextLink != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.getFeatureFlagsCreateRequest(ctx, options)
			}, &runtime.FetcherForNextLinkOptions{
				NextReq: func(ctx context.Context, encodedNextLink string) (*policy.Request, error) {
					return client.getFeatureFlagNextPageCreateRequest(ctx, encodedNextLink)
				},
			})
			if err != nil {
				return AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse{}, err
			}
			return client.getFeatureFlagsHandleResponseWithLinkHeader(resp)
		},
	})
}

// NewGetFeatureFlagRevisionsPagerWithLinkHeader wraps the generated feature flag revisions pager
// so that the nextLink can also be discovered from the "Link" response header.
func (client *AzureAppConfigurationFeatureFlagClient) NewGetFeatureFlagRevisionsPagerWithLinkHeader(options *AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsOptions) *runtime.Pager[AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse] {
	return runtime.NewPager(runtime.PagingHandler[AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse]{
		More: func(page AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse) (AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse, error) {
			nextLink := ""
			if page != nil && page.NextLink != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.getFeatureFlagRevisionsCreateRequest(ctx, options)
			}, &runtime.FetcherForNextLinkOptions{
				NextReq: func(ctx context.Context, encodedNextLink string) (*policy.Request, error) {
					return client.getFeatureFlagNextPageCreateRequest(ctx, encodedNextLink)
				},
			})
			if err != nil {
				return AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse{}, err
			}
			return client.getFeatureFlagRevisionsHandleResponseWithLinkHeader(resp)
		},
	})
}

// getFeatureFlagNextPageCreateRequest builds a request for a subsequent page of feature flags or
// revisions from the encoded nextLink returned by the service.
func (client *AzureAppConfigurationFeatureFlagClient) getFeatureFlagNextPageCreateRequest(ctx context.Context, nextLink string) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, nextLink))
	if err != nil {
		return nil, err
	}
	req.Raw().Header["Accept"] = []string{ffAcceptHeader}
	return req, nil
}

// getFeatureFlagsHandleResponseWithLinkHeader parses the nextLink from the "Link" response header
// when it is not populated in the JSON body.
func (client *AzureAppConfigurationFeatureFlagClient) getFeatureFlagsHandleResponseWithLinkHeader(resp *http.Response) (AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse, error) {
	result, err := client.getFeatureFlagsHandleResponse(resp)
	if err != nil {
		return AzureAppConfigurationFeatureFlagClientGetFeatureFlagsResponse{}, err
	}
	if result.NextLink != nil {
		return result, nil
	}
	if nextLink := nextLinkFromHeader(resp.Header.Get("Link")); nextLink != "" {
		result.NextLink = to.Ptr(nextLink)
	}
	return result, nil
}

// getFeatureFlagRevisionsHandleResponseWithLinkHeader parses the nextLink from the "Link" response
// header when it is not populated in the JSON body.
func (client *AzureAppConfigurationFeatureFlagClient) getFeatureFlagRevisionsHandleResponseWithLinkHeader(resp *http.Response) (AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse, error) {
	result, err := client.getFeatureFlagRevisionsHandleResponse(resp)
	if err != nil {
		return AzureAppConfigurationFeatureFlagClientGetFeatureFlagRevisionsResponse{}, err
	}
	if result.NextLink != nil {
		return result, nil
	}
	if nextLink := nextLinkFromHeader(resp.Header.Get("Link")); nextLink != "" {
		result.NextLink = to.Ptr(nextLink)
	}
	return result, nil
}

// nextLinkFromHeader extracts the URL from a Link header value of the form
// `<nextLinkURL>; rel="next"`. Returns "" if the header does not contain a URL.
func nextLinkFromHeader(link string) string {
	if link == "" {
		return ""
	}
	if endIndex := strings.Index(link, ">"); endIndex > 0 {
		return link[1:endIndex]
	}
	return ""
}
