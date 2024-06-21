//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmobilenetwork

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

// PacketCoreControlPlaneVersionsClient contains the methods for the PacketCoreControlPlaneVersions group.
// Don't use this type directly, use NewPacketCoreControlPlaneVersionsClient() instead.
type PacketCoreControlPlaneVersionsClient struct {
	internal *arm.Client
}

// NewPacketCoreControlPlaneVersionsClient creates a new instance of PacketCoreControlPlaneVersionsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPacketCoreControlPlaneVersionsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*PacketCoreControlPlaneVersionsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PacketCoreControlPlaneVersionsClient{
		internal: cl,
	}
	return client, nil
}

// Get - Gets information about the specified packet core control plane version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-04-01
//   - versionName - The name of the packet core control plane version.
//   - options - PacketCoreControlPlaneVersionsClientGetOptions contains the optional parameters for the PacketCoreControlPlaneVersionsClient.Get
//     method.
func (client *PacketCoreControlPlaneVersionsClient) Get(ctx context.Context, versionName string, options *PacketCoreControlPlaneVersionsClientGetOptions) (PacketCoreControlPlaneVersionsClientGetResponse, error) {
	var err error
	const operationName = "PacketCoreControlPlaneVersionsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, versionName, options)
	if err != nil {
		return PacketCoreControlPlaneVersionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PacketCoreControlPlaneVersionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PacketCoreControlPlaneVersionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PacketCoreControlPlaneVersionsClient) getCreateRequest(ctx context.Context, versionName string, options *PacketCoreControlPlaneVersionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.MobileNetwork/packetCoreControlPlaneVersions/{versionName}"
	if versionName == "" {
		return nil, errors.New("parameter versionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{versionName}", url.PathEscape(versionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PacketCoreControlPlaneVersionsClient) getHandleResponse(resp *http.Response) (PacketCoreControlPlaneVersionsClientGetResponse, error) {
	result := PacketCoreControlPlaneVersionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PacketCoreControlPlaneVersion); err != nil {
		return PacketCoreControlPlaneVersionsClientGetResponse{}, err
	}
	return result, nil
}

// GetBySubscription - Gets information about the specified packet core control plane version.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-04-01
//   - versionName - The name of the packet core control plane version.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - options - PacketCoreControlPlaneVersionsClientGetBySubscriptionOptions contains the optional parameters for the PacketCoreControlPlaneVersionsClient.GetBySubscription
//     method.
func (client *PacketCoreControlPlaneVersionsClient) GetBySubscription(ctx context.Context, versionName string, subscriptionID string, options *PacketCoreControlPlaneVersionsClientGetBySubscriptionOptions) (PacketCoreControlPlaneVersionsClientGetBySubscriptionResponse, error) {
	var err error
	const operationName = "PacketCoreControlPlaneVersionsClient.GetBySubscription"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getBySubscriptionCreateRequest(ctx, versionName, subscriptionID, options)
	if err != nil {
		return PacketCoreControlPlaneVersionsClientGetBySubscriptionResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PacketCoreControlPlaneVersionsClientGetBySubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PacketCoreControlPlaneVersionsClientGetBySubscriptionResponse{}, err
	}
	resp, err := client.getBySubscriptionHandleResponse(httpResp)
	return resp, err
}

// getBySubscriptionCreateRequest creates the GetBySubscription request.
func (client *PacketCoreControlPlaneVersionsClient) getBySubscriptionCreateRequest(ctx context.Context, versionName string, subscriptionID string, options *PacketCoreControlPlaneVersionsClientGetBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.MobileNetwork/packetCoreControlPlaneVersions/{versionName}"
	if versionName == "" {
		return nil, errors.New("parameter versionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{versionName}", url.PathEscape(versionName))
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getBySubscriptionHandleResponse handles the GetBySubscription response.
func (client *PacketCoreControlPlaneVersionsClient) getBySubscriptionHandleResponse(resp *http.Response) (PacketCoreControlPlaneVersionsClientGetBySubscriptionResponse, error) {
	result := PacketCoreControlPlaneVersionsClientGetBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PacketCoreControlPlaneVersion); err != nil {
		return PacketCoreControlPlaneVersionsClientGetBySubscriptionResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all supported packet core control planes versions.
//
// Generated from API version 2024-04-01
//   - options - PacketCoreControlPlaneVersionsClientListOptions contains the optional parameters for the PacketCoreControlPlaneVersionsClient.NewListPager
//     method.
func (client *PacketCoreControlPlaneVersionsClient) NewListPager(options *PacketCoreControlPlaneVersionsClientListOptions) *runtime.Pager[PacketCoreControlPlaneVersionsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[PacketCoreControlPlaneVersionsClientListResponse]{
		More: func(page PacketCoreControlPlaneVersionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PacketCoreControlPlaneVersionsClientListResponse) (PacketCoreControlPlaneVersionsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PacketCoreControlPlaneVersionsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return PacketCoreControlPlaneVersionsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *PacketCoreControlPlaneVersionsClient) listCreateRequest(ctx context.Context, options *PacketCoreControlPlaneVersionsClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.MobileNetwork/packetCoreControlPlaneVersions"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PacketCoreControlPlaneVersionsClient) listHandleResponse(resp *http.Response) (PacketCoreControlPlaneVersionsClientListResponse, error) {
	result := PacketCoreControlPlaneVersionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PacketCoreControlPlaneVersionListResult); err != nil {
		return PacketCoreControlPlaneVersionsClientListResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Lists all supported packet core control planes versions.
//
// Generated from API version 2024-04-01
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - options - PacketCoreControlPlaneVersionsClientListBySubscriptionOptions contains the optional parameters for the PacketCoreControlPlaneVersionsClient.NewListBySubscriptionPager
//     method.
func (client *PacketCoreControlPlaneVersionsClient) NewListBySubscriptionPager(subscriptionID string, options *PacketCoreControlPlaneVersionsClientListBySubscriptionOptions) *runtime.Pager[PacketCoreControlPlaneVersionsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[PacketCoreControlPlaneVersionsClientListBySubscriptionResponse]{
		More: func(page PacketCoreControlPlaneVersionsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PacketCoreControlPlaneVersionsClientListBySubscriptionResponse) (PacketCoreControlPlaneVersionsClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PacketCoreControlPlaneVersionsClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, subscriptionID, options)
			}, nil)
			if err != nil {
				return PacketCoreControlPlaneVersionsClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *PacketCoreControlPlaneVersionsClient) listBySubscriptionCreateRequest(ctx context.Context, subscriptionID string, options *PacketCoreControlPlaneVersionsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.MobileNetwork/packetCoreControlPlaneVersions"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *PacketCoreControlPlaneVersionsClient) listBySubscriptionHandleResponse(resp *http.Response) (PacketCoreControlPlaneVersionsClientListBySubscriptionResponse, error) {
	result := PacketCoreControlPlaneVersionsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PacketCoreControlPlaneVersionListResult); err != nil {
		return PacketCoreControlPlaneVersionsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}
