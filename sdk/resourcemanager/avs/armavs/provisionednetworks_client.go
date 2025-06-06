// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armavs

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

// ProvisionedNetworksClient contains the methods for the ProvisionedNetworks group.
// Don't use this type directly, use NewProvisionedNetworksClient() instead.
type ProvisionedNetworksClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewProvisionedNetworksClient creates a new instance of ProvisionedNetworksClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewProvisionedNetworksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ProvisionedNetworksClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ProvisionedNetworksClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Get a ProvisionedNetwork
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - privateCloudName - Name of the private cloud
//   - provisionedNetworkName - Name of the cloud link.
//   - options - ProvisionedNetworksClientGetOptions contains the optional parameters for the ProvisionedNetworksClient.Get method.
func (client *ProvisionedNetworksClient) Get(ctx context.Context, resourceGroupName string, privateCloudName string, provisionedNetworkName string, options *ProvisionedNetworksClientGetOptions) (ProvisionedNetworksClientGetResponse, error) {
	var err error
	const operationName = "ProvisionedNetworksClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, privateCloudName, provisionedNetworkName, options)
	if err != nil {
		return ProvisionedNetworksClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ProvisionedNetworksClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ProvisionedNetworksClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ProvisionedNetworksClient) getCreateRequest(ctx context.Context, resourceGroupName string, privateCloudName string, provisionedNetworkName string, _ *ProvisionedNetworksClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}/provisionedNetworks/{provisionedNetworkName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if privateCloudName == "" {
		return nil, errors.New("parameter privateCloudName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateCloudName}", url.PathEscape(privateCloudName))
	if provisionedNetworkName == "" {
		return nil, errors.New("parameter provisionedNetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{provisionedNetworkName}", url.PathEscape(provisionedNetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ProvisionedNetworksClient) getHandleResponse(resp *http.Response) (ProvisionedNetworksClientGetResponse, error) {
	result := ProvisionedNetworksClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ProvisionedNetwork); err != nil {
		return ProvisionedNetworksClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List ProvisionedNetwork resources by PrivateCloud
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - privateCloudName - Name of the private cloud
//   - options - ProvisionedNetworksClientListOptions contains the optional parameters for the ProvisionedNetworksClient.NewListPager
//     method.
func (client *ProvisionedNetworksClient) NewListPager(resourceGroupName string, privateCloudName string, options *ProvisionedNetworksClientListOptions) *runtime.Pager[ProvisionedNetworksClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ProvisionedNetworksClientListResponse]{
		More: func(page ProvisionedNetworksClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ProvisionedNetworksClientListResponse) (ProvisionedNetworksClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ProvisionedNetworksClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, privateCloudName, options)
			}, nil)
			if err != nil {
				return ProvisionedNetworksClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *ProvisionedNetworksClient) listCreateRequest(ctx context.Context, resourceGroupName string, privateCloudName string, _ *ProvisionedNetworksClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}/provisionedNetworks"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if privateCloudName == "" {
		return nil, errors.New("parameter privateCloudName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateCloudName}", url.PathEscape(privateCloudName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ProvisionedNetworksClient) listHandleResponse(resp *http.Response) (ProvisionedNetworksClientListResponse, error) {
	result := ProvisionedNetworksClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ProvisionedNetworkListResult); err != nil {
		return ProvisionedNetworksClientListResponse{}, err
	}
	return result, nil
}
