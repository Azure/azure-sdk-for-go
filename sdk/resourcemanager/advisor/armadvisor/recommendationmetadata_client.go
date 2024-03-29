//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armadvisor

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// RecommendationMetadataClient contains the methods for the RecommendationMetadata group.
// Don't use this type directly, use NewRecommendationMetadataClient() instead.
type RecommendationMetadataClient struct {
	internal *arm.Client
}

// NewRecommendationMetadataClient creates a new instance of RecommendationMetadataClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewRecommendationMetadataClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*RecommendationMetadataClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &RecommendationMetadataClient{
		internal: cl,
	}
	return client, nil
}

// Get - Gets the metadata entity.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-01-01
//   - name - Name of metadata entity.
//   - options - RecommendationMetadataClientGetOptions contains the optional parameters for the RecommendationMetadataClient.Get
//     method.
func (client *RecommendationMetadataClient) Get(ctx context.Context, name string, options *RecommendationMetadataClientGetOptions) (RecommendationMetadataClientGetResponse, error) {
	var err error
	const operationName = "RecommendationMetadataClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, name, options)
	if err != nil {
		return RecommendationMetadataClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RecommendationMetadataClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RecommendationMetadataClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *RecommendationMetadataClient) getCreateRequest(ctx context.Context, name string, options *RecommendationMetadataClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Advisor/metadata/{name}"
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RecommendationMetadataClient) getHandleResponse(resp *http.Response) (RecommendationMetadataClientGetResponse, error) {
	result := RecommendationMetadataClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MetadataEntity); err != nil {
		return RecommendationMetadataClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets the list of metadata entities.
//
// Generated from API version 2020-01-01
//   - options - RecommendationMetadataClientListOptions contains the optional parameters for the RecommendationMetadataClient.NewListPager
//     method.
func (client *RecommendationMetadataClient) NewListPager(options *RecommendationMetadataClientListOptions) *runtime.Pager[RecommendationMetadataClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[RecommendationMetadataClientListResponse]{
		More: func(page RecommendationMetadataClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RecommendationMetadataClientListResponse) (RecommendationMetadataClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "RecommendationMetadataClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return RecommendationMetadataClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *RecommendationMetadataClient) listCreateRequest(ctx context.Context, options *RecommendationMetadataClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Advisor/metadata"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *RecommendationMetadataClient) listHandleResponse(resp *http.Response) (RecommendationMetadataClientListResponse, error) {
	result := RecommendationMetadataClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MetadataEntityListResult); err != nil {
		return RecommendationMetadataClientListResponse{}, err
	}
	return result, nil
}
