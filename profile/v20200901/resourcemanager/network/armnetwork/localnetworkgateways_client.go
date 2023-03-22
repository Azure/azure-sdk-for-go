//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armnetwork

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profile/v20200901/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// LocalNetworkGatewaysClient contains the methods for the LocalNetworkGateways group.
// Don't use this type directly, use NewLocalNetworkGatewaysClient() instead.
type LocalNetworkGatewaysClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewLocalNetworkGatewaysClient creates a new instance of LocalNetworkGatewaysClient with the specified values.
//   - subscriptionID - The subscription credentials which uniquely identify the Microsoft Azure subscription. The subscription
//     ID forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewLocalNetworkGatewaysClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LocalNetworkGatewaysClient, error) {
	cl, err := arm.NewClient(internal.ModuleName+"/armnetwork.LocalNetworkGatewaysClient", internal.ModuleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &LocalNetworkGatewaysClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a local network gateway in the specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-11-01
//   - resourceGroupName - The name of the resource group.
//   - localNetworkGatewayName - The name of the local network gateway.
//   - parameters - Parameters supplied to the create or update local network gateway operation.
//   - options - LocalNetworkGatewaysClientBeginCreateOrUpdateOptions contains the optional parameters for the LocalNetworkGatewaysClient.BeginCreateOrUpdate
//     method.
func (client *LocalNetworkGatewaysClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, parameters LocalNetworkGateway, options *LocalNetworkGatewaysClientBeginCreateOrUpdateOptions) (*runtime.Poller[LocalNetworkGatewaysClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, localNetworkGatewayName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[LocalNetworkGatewaysClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[LocalNetworkGatewaysClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a local network gateway in the specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-11-01
func (client *LocalNetworkGatewaysClient) createOrUpdate(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, parameters LocalNetworkGateway, options *LocalNetworkGatewaysClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, localNetworkGatewayName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *LocalNetworkGatewaysClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, parameters LocalNetworkGateway, options *LocalNetworkGatewaysClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/localNetworkGateways/{localNetworkGatewayName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localNetworkGatewayName == "" {
		return nil, errors.New("parameter localNetworkGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localNetworkGatewayName}", url.PathEscape(localNetworkGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Deletes the specified local network gateway.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-11-01
//   - resourceGroupName - The name of the resource group.
//   - localNetworkGatewayName - The name of the local network gateway.
//   - options - LocalNetworkGatewaysClientBeginDeleteOptions contains the optional parameters for the LocalNetworkGatewaysClient.BeginDelete
//     method.
func (client *LocalNetworkGatewaysClient) BeginDelete(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, options *LocalNetworkGatewaysClientBeginDeleteOptions) (*runtime.Poller[LocalNetworkGatewaysClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, localNetworkGatewayName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[LocalNetworkGatewaysClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[LocalNetworkGatewaysClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the specified local network gateway.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-11-01
func (client *LocalNetworkGatewaysClient) deleteOperation(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, options *LocalNetworkGatewaysClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, localNetworkGatewayName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *LocalNetworkGatewaysClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, options *LocalNetworkGatewaysClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/localNetworkGateways/{localNetworkGatewayName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localNetworkGatewayName == "" {
		return nil, errors.New("parameter localNetworkGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localNetworkGatewayName}", url.PathEscape(localNetworkGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Gets the specified local network gateway in a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-11-01
//   - resourceGroupName - The name of the resource group.
//   - localNetworkGatewayName - The name of the local network gateway.
//   - options - LocalNetworkGatewaysClientGetOptions contains the optional parameters for the LocalNetworkGatewaysClient.Get
//     method.
func (client *LocalNetworkGatewaysClient) Get(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, options *LocalNetworkGatewaysClientGetOptions) (LocalNetworkGatewaysClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, localNetworkGatewayName, options)
	if err != nil {
		return LocalNetworkGatewaysClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return LocalNetworkGatewaysClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LocalNetworkGatewaysClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *LocalNetworkGatewaysClient) getCreateRequest(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, options *LocalNetworkGatewaysClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/localNetworkGateways/{localNetworkGatewayName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localNetworkGatewayName == "" {
		return nil, errors.New("parameter localNetworkGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localNetworkGatewayName}", url.PathEscape(localNetworkGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *LocalNetworkGatewaysClient) getHandleResponse(resp *http.Response) (LocalNetworkGatewaysClientGetResponse, error) {
	result := LocalNetworkGatewaysClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocalNetworkGateway); err != nil {
		return LocalNetworkGatewaysClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets all the local network gateways in a resource group.
//
// Generated from API version 2018-11-01
//   - resourceGroupName - The name of the resource group.
//   - options - LocalNetworkGatewaysClientListOptions contains the optional parameters for the LocalNetworkGatewaysClient.NewListPager
//     method.
func (client *LocalNetworkGatewaysClient) NewListPager(resourceGroupName string, options *LocalNetworkGatewaysClientListOptions) *runtime.Pager[LocalNetworkGatewaysClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[LocalNetworkGatewaysClientListResponse]{
		More: func(page LocalNetworkGatewaysClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *LocalNetworkGatewaysClientListResponse) (LocalNetworkGatewaysClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return LocalNetworkGatewaysClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return LocalNetworkGatewaysClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LocalNetworkGatewaysClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *LocalNetworkGatewaysClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *LocalNetworkGatewaysClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/localNetworkGateways"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *LocalNetworkGatewaysClient) listHandleResponse(resp *http.Response) (LocalNetworkGatewaysClientListResponse, error) {
	result := LocalNetworkGatewaysClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocalNetworkGatewayListResult); err != nil {
		return LocalNetworkGatewaysClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdateTags - Updates a local network gateway tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-11-01
//   - resourceGroupName - The name of the resource group.
//   - localNetworkGatewayName - The name of the local network gateway.
//   - parameters - Parameters supplied to update local network gateway tags.
//   - options - LocalNetworkGatewaysClientBeginUpdateTagsOptions contains the optional parameters for the LocalNetworkGatewaysClient.BeginUpdateTags
//     method.
func (client *LocalNetworkGatewaysClient) BeginUpdateTags(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, parameters TagsObject, options *LocalNetworkGatewaysClientBeginUpdateTagsOptions) (*runtime.Poller[LocalNetworkGatewaysClientUpdateTagsResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.updateTags(ctx, resourceGroupName, localNetworkGatewayName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[LocalNetworkGatewaysClientUpdateTagsResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[LocalNetworkGatewaysClientUpdateTagsResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// UpdateTags - Updates a local network gateway tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-11-01
func (client *LocalNetworkGatewaysClient) updateTags(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, parameters TagsObject, options *LocalNetworkGatewaysClientBeginUpdateTagsOptions) (*http.Response, error) {
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, localNetworkGatewayName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *LocalNetworkGatewaysClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, localNetworkGatewayName string, parameters TagsObject, options *LocalNetworkGatewaysClientBeginUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/localNetworkGateways/{localNetworkGatewayName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if localNetworkGatewayName == "" {
		return nil, errors.New("parameter localNetworkGatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{localNetworkGatewayName}", url.PathEscape(localNetworkGatewayName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}
