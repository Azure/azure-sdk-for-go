//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armedgemarketplace

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PublishersClient contains the methods for the Publishers group.
// Don't use this type directly, use NewPublishersClient() instead.
type PublishersClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewPublishersClient creates a new instance of PublishersClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPublishersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PublishersClient, error) {
	cl, err := arm.NewClient(moduleName+".PublishersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PublishersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Get a Publisher
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-08-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource.
//   - publisherName - Name of the publisher
//   - options - PublishersClientGetOptions contains the optional parameters for the PublishersClient.Get method.
func (client *PublishersClient) Get(ctx context.Context, resourceURI string, publisherName string, options *PublishersClientGetOptions) (PublishersClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceURI, publisherName, options)
	if err != nil {
		return PublishersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PublishersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PublishersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PublishersClient) getCreateRequest(ctx context.Context, resourceURI string, publisherName string, options *PublishersClientGetOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.EdgeMarketplace/publishers/{publisherName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PublishersClient) getHandleResponse(resp *http.Response) (PublishersClientGetResponse, error) {
	result := PublishersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Publisher); err != nil {
		return PublishersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List Publisher resources by parent
//
// Generated from API version 2023-08-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource.
//   - options - PublishersClientListOptions contains the optional parameters for the PublishersClient.NewListPager method.
func (client *PublishersClient) NewListPager(resourceURI string, options *PublishersClientListOptions) *runtime.Pager[PublishersClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[PublishersClientListResponse]{
		More: func(page PublishersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PublishersClientListResponse) (PublishersClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceURI, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PublishersClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PublishersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PublishersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *PublishersClient) listCreateRequest(ctx context.Context, resourceURI string, options *PublishersClientListOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.EdgeMarketplace/publishers"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-08-01-preview")
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Maxpagesize != nil {
		reqQP.Set("maxpagesize", strconv.FormatInt(int64(*options.Maxpagesize), 10))
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PublishersClient) listHandleResponse(resp *http.Response) (PublishersClientListResponse, error) {
	result := PublishersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublisherListResult); err != nil {
		return PublishersClientListResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - List Publisher resources in subscription
//
// Generated from API version 2023-08-01-preview
//   - options - PublishersClientListBySubscriptionOptions contains the optional parameters for the PublishersClient.NewListBySubscriptionPager
//     method.
func (client *PublishersClient) NewListBySubscriptionPager(options *PublishersClientListBySubscriptionOptions) *runtime.Pager[PublishersClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[PublishersClientListBySubscriptionResponse]{
		More: func(page PublishersClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PublishersClientListBySubscriptionResponse) (PublishersClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PublishersClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PublishersClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PublishersClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *PublishersClient) listBySubscriptionCreateRequest(ctx context.Context, options *PublishersClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.EdgeMarketplace/publishers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-08-01-preview")
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Maxpagesize != nil {
		reqQP.Set("maxpagesize", strconv.FormatInt(int64(*options.Maxpagesize), 10))
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *PublishersClient) listBySubscriptionHandleResponse(resp *http.Response) (PublishersClientListBySubscriptionResponse, error) {
	result := PublishersClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublisherListResult); err != nil {
		return PublishersClientListBySubscriptionResponse{}, err
	}
	return result, nil
}
