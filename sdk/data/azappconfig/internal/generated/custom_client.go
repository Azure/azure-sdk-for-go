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

func (client *AzureAppConfigurationClient) NewGetKeyValuesPagerWithNextLinks(options *AzureAppConfigurationClientGetKeyValuesOptions) *runtime.Pager[AzureAppConfigurationClientGetKeyValuesResponse] {
	if options == nil {
		options = &AzureAppConfigurationClientGetKeyValuesOptions{}
	}
	return newPagerWithNextLinks(client, options.AcceptDatetime,
		func(ctx context.Context) (*policy.Request, error) {
			return client.getKeyValuesCreateRequest(ctx, options)
		},
		client.getKeyValuesHandleResponseWithLinkHeader,
		func(page AzureAppConfigurationClientGetKeyValuesResponse) *string { return page.NextLink },
	)
}

func (client *AzureAppConfigurationClient) NewGetRevisionsPagerWithNextLinks(options *AzureAppConfigurationClientGetRevisionsOptions) *runtime.Pager[AzureAppConfigurationClientGetRevisionsResponse] {
	if options == nil {
		options = &AzureAppConfigurationClientGetRevisionsOptions{}
	}
	return newPagerWithNextLinks(client, options.AcceptDatetime,
		func(ctx context.Context) (*policy.Request, error) {
			return client.getRevisionsCreateRequest(ctx, options)
		},
		func(resp *http.Response) (AzureAppConfigurationClientGetRevisionsResponse, error) {
			result, err := client.getRevisionsHandleResponse(resp)
			if err == nil && result.NextLink == nil {
				result.NextLink = nextLinkFromHeader(resp.Header)
			}
			return result, err
		},
		func(page AzureAppConfigurationClientGetRevisionsResponse) *string { return page.NextLink },
	)
}

func (client *AzureAppConfigurationClient) NewGetSnapshotsPagerWithNextLinks(options *AzureAppConfigurationClientGetSnapshotsOptions) *runtime.Pager[AzureAppConfigurationClientGetSnapshotsResponse] {
	return newPagerWithNextLinks(client, nil,
		func(ctx context.Context) (*policy.Request, error) {
			return client.getSnapshotsCreateRequest(ctx, options)
		},
		client.getSnapshotsHandleResponse,
		func(page AzureAppConfigurationClientGetSnapshotsResponse) *string { return page.NextLink },
	)
}

// The service can return endpoint-relative next links, which the generated pagers don't resolve.
func newPagerWithNextLinks[T any](client *AzureAppConfigurationClient, acceptDatetime *string, createRequest func(context.Context) (*policy.Request, error), handleResponse func(*http.Response) (T, error), getNextLink func(T) *string) *runtime.Pager[T] {
	return runtime.NewPager(runtime.PagingHandler[T]{
		More: func(page T) bool {
			nextLink := getNextLink(page)
			return nextLink != nil && *nextLink != ""
		},
		Fetcher: func(ctx context.Context, page *T) (T, error) {
			nextLink := ""
			if page != nil {
				nextLink = *getNextLink(*page)
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, createRequest, &runtime.FetcherForNextLinkOptions{
				NextReq: client.nextPageRequest(http.MethodGet, acceptDatetime, azcore.MatchConditions{}),
			})
			if err != nil {
				var zero T
				return zero, err
			}
			return handleResponse(resp)
		},
	})
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
				NextReq:     client.nextPageRequest(http.MethodGet, options.AcceptDatetime, curCondition),
				StatusCodes: []int{http.StatusNotModified},
			})
			if err != nil {
				return AzureAppConfigurationClientGetKeyValuesResponse{}, err
			}
			return client.getKeyValuesHandleResponseWithLinkHeader(resp)
		},
	})
}

func (a *AzureAppConfigurationClient) getKeyValuesHandleResponseWithLinkHeader(resp *http.Response) (AzureAppConfigurationClientGetKeyValuesResponse, error) {
	result, err := a.getKeyValuesHandleResponse(resp)
	if err == nil && result.NextLink == nil {
		result.NextLink = nextLinkFromHeader(resp.Header)
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
				NextReq:     client.nextPageRequest(http.MethodHead, options.AcceptDatetime, curCondition),
				StatusCodes: []int{http.StatusNotModified},
			})
			if err != nil {
				return CheckKeyValuesPagerResponse{}, err
			}
			return client.checkKeyValuesHandleResponseWithLinkHeader(resp)
		},
	})
}

func (a *AzureAppConfigurationClient) checkKeyValuesHandleResponseWithLinkHeader(resp *http.Response) (CheckKeyValuesPagerResponse, error) {
	genResult, err := a.checkKeyValuesHandleResponse(resp)
	if err != nil {
		return CheckKeyValuesPagerResponse{}, err
	}

	return CheckKeyValuesPagerResponse{
		AzureAppConfigurationClientCheckKeyValuesResponse: genResult,
		NextLink: nextLinkFromHeader(resp.Header),
	}, nil
}

// CreateSnapshot exports internal createSnapshot
func (a *AzureAppConfigurationClient) CreateSnapshot(ctx context.Context, contentType CreateSnapshotRequestContentType, name string, entity Snapshot, options *AzureAppConfigurationClientBeginCreateSnapshotOptions) (*http.Response, error) {
	return a.createSnapshot(ctx, contentType, name, entity, options)
}

func (client *AzureAppConfigurationClient) nextPageRequest(method string, acceptDatetime *string, conditions azcore.MatchConditions) func(context.Context, string) (*policy.Request, error) {
	return func(ctx context.Context, nextLink string) (*policy.Request, error) {
		req, err := runtime.NewRequestForNextLink(ctx, method, client.endpoint, nextLink)
		if err != nil {
			return nil, err
		}
		h := req.Raw().Header
		if method == http.MethodGet {
			h.Set("Accept", "application/json")
		}
		if acceptDatetime != nil {
			h.Set("Accept-Datetime", *acceptDatetime)
		}
		if conditions.IfMatch != nil {
			h.Set("If-Match", string(*conditions.IfMatch))
		}
		if conditions.IfNoneMatch != nil {
			h.Set("If-None-Match", string(*conditions.IfNoneMatch))
		}
		return req, nil
	}
}

// Historical responses can contain both rel="original" and rel="next".
func nextLinkFromHeader(h http.Header) *string {
	for _, links := range h.Values("Link") {
		for links != "" {
			target, params, ok := strings.Cut(strings.TrimSpace(links), ">")
			if !ok || !strings.HasPrefix(target, "<") {
				break
			}
			links = strings.TrimSpace(params)
			for strings.HasPrefix(links, ";") {
				var param string
				param, links = cutLinkParam(links[1:])
				name, value, ok := strings.Cut(param, "=")
				if !ok || !strings.EqualFold(strings.TrimSpace(name), "rel") {
					continue
				}
				value = strings.TrimSpace(value)
				if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
					value = value[1 : len(value)-1]
				}
				for rel := range strings.FieldsSeq(value) {
					if strings.EqualFold(rel, "next") {
						return to.Ptr(target[1:])
					}
				}
			}
			links, ok = strings.CutPrefix(links, ",")
			if !ok {
				break
			}
		}
	}
	return nil
}

// Commas and semicolons inside a quoted link parameter aren't separators.
func cutLinkParam(s string) (string, string) {
	quoted := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			quoted = !quoted
		case '\\':
			if quoted {
				i++
			}
		case ',', ';':
			if !quoted {
				return s[:i], s[i:]
			}
		}
	}
	return s, ""
}
