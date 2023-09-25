//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybriddatamanager

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

// PublicKeysClient contains the methods for the PublicKeys group.
// Don't use this type directly, use NewPublicKeysClient() instead.
type PublicKeysClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewPublicKeysClient creates a new instance of PublicKeysClient with the specified values.
//   - subscriptionID - The Subscription Id
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPublicKeysClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PublicKeysClient, error) {
	cl, err := arm.NewClient(moduleName+".PublicKeysClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PublicKeysClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - This method gets the public keys.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-06-01
//   - publicKeyName - Name of the public key.
//   - resourceGroupName - The Resource Group Name
//   - dataManagerName - The name of the DataManager Resource within the specified resource group. DataManager names must be between
//     3 and 24 characters in length and use any alphanumeric and underscore only
//   - options - PublicKeysClientGetOptions contains the optional parameters for the PublicKeysClient.Get method.
func (client *PublicKeysClient) Get(ctx context.Context, publicKeyName string, resourceGroupName string, dataManagerName string, options *PublicKeysClientGetOptions) (PublicKeysClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, publicKeyName, resourceGroupName, dataManagerName, options)
	if err != nil {
		return PublicKeysClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PublicKeysClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PublicKeysClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PublicKeysClient) getCreateRequest(ctx context.Context, publicKeyName string, resourceGroupName string, dataManagerName string, options *PublicKeysClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridData/dataManagers/{dataManagerName}/publicKeys/{publicKeyName}"
	if publicKeyName == "" {
		return nil, errors.New("parameter publicKeyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publicKeyName}", url.PathEscape(publicKeyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataManagerName == "" {
		return nil, errors.New("parameter dataManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataManagerName}", url.PathEscape(dataManagerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PublicKeysClient) getHandleResponse(resp *http.Response) (PublicKeysClientGetResponse, error) {
	result := PublicKeysClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublicKey); err != nil {
		return PublicKeysClientGetResponse{}, err
	}
	return result, nil
}

// NewListByDataManagerPager - This method gets the list view of public keys, however it will only have one element.
//
// Generated from API version 2019-06-01
//   - resourceGroupName - The Resource Group Name
//   - dataManagerName - The name of the DataManager Resource within the specified resource group. DataManager names must be between
//     3 and 24 characters in length and use any alphanumeric and underscore only
//   - options - PublicKeysClientListByDataManagerOptions contains the optional parameters for the PublicKeysClient.NewListByDataManagerPager
//     method.
func (client *PublicKeysClient) NewListByDataManagerPager(resourceGroupName string, dataManagerName string, options *PublicKeysClientListByDataManagerOptions) (*runtime.Pager[PublicKeysClientListByDataManagerResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PublicKeysClientListByDataManagerResponse]{
		More: func(page PublicKeysClientListByDataManagerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PublicKeysClientListByDataManagerResponse) (PublicKeysClientListByDataManagerResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByDataManagerCreateRequest(ctx, resourceGroupName, dataManagerName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PublicKeysClientListByDataManagerResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PublicKeysClientListByDataManagerResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PublicKeysClientListByDataManagerResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByDataManagerHandleResponse(resp)
		},
	})
}

// listByDataManagerCreateRequest creates the ListByDataManager request.
func (client *PublicKeysClient) listByDataManagerCreateRequest(ctx context.Context, resourceGroupName string, dataManagerName string, options *PublicKeysClientListByDataManagerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridData/dataManagers/{dataManagerName}/publicKeys"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataManagerName == "" {
		return nil, errors.New("parameter dataManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataManagerName}", url.PathEscape(dataManagerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByDataManagerHandleResponse handles the ListByDataManager response.
func (client *PublicKeysClient) listByDataManagerHandleResponse(resp *http.Response) (PublicKeysClientListByDataManagerResponse, error) {
	result := PublicKeysClientListByDataManagerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublicKeyList); err != nil {
		return PublicKeysClientListByDataManagerResponse{}, err
	}
	return result, nil
}

