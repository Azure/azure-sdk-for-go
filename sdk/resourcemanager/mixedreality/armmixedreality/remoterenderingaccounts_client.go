//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmixedreality

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

// RemoteRenderingAccountsClient contains the methods for the RemoteRenderingAccounts group.
// Don't use this type directly, use NewRemoteRenderingAccountsClient() instead.
type RemoteRenderingAccountsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewRemoteRenderingAccountsClient creates a new instance of RemoteRenderingAccountsClient with the specified values.
//   - subscriptionID - The Azure subscription ID. This is a GUID-formatted string (e.g. 00000000-0000-0000-0000-000000000000)
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewRemoteRenderingAccountsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RemoteRenderingAccountsClient, error) {
	cl, err := arm.NewClient(moduleName+".RemoteRenderingAccountsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &RemoteRenderingAccountsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Create - Creating or Updating a Remote Rendering Account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-01-preview
//   - resourceGroupName - Name of an Azure resource group.
//   - accountName - Name of an Mixed Reality Account.
//   - remoteRenderingAccount - Remote Rendering Account parameter.
//   - options - RemoteRenderingAccountsClientCreateOptions contains the optional parameters for the RemoteRenderingAccountsClient.Create
//     method.
func (client *RemoteRenderingAccountsClient) Create(ctx context.Context, resourceGroupName string, accountName string, remoteRenderingAccount RemoteRenderingAccount, options *RemoteRenderingAccountsClientCreateOptions) (RemoteRenderingAccountsClientCreateResponse, error) {
	var err error
	req, err := client.createCreateRequest(ctx, resourceGroupName, accountName, remoteRenderingAccount, options)
	if err != nil {
		return RemoteRenderingAccountsClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RemoteRenderingAccountsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return RemoteRenderingAccountsClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *RemoteRenderingAccountsClient) createCreateRequest(ctx context.Context, resourceGroupName string, accountName string, remoteRenderingAccount RemoteRenderingAccount, options *RemoteRenderingAccountsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MixedReality/remoteRenderingAccounts/{accountName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, remoteRenderingAccount); err != nil {
	return nil, err
}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *RemoteRenderingAccountsClient) createHandleResponse(resp *http.Response) (RemoteRenderingAccountsClientCreateResponse, error) {
	result := RemoteRenderingAccountsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RemoteRenderingAccount); err != nil {
		return RemoteRenderingAccountsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a Remote Rendering Account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-01-preview
//   - resourceGroupName - Name of an Azure resource group.
//   - accountName - Name of an Mixed Reality Account.
//   - options - RemoteRenderingAccountsClientDeleteOptions contains the optional parameters for the RemoteRenderingAccountsClient.Delete
//     method.
func (client *RemoteRenderingAccountsClient) Delete(ctx context.Context, resourceGroupName string, accountName string, options *RemoteRenderingAccountsClientDeleteOptions) (RemoteRenderingAccountsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, options)
	if err != nil {
		return RemoteRenderingAccountsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RemoteRenderingAccountsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return RemoteRenderingAccountsClientDeleteResponse{}, err
	}
	return RemoteRenderingAccountsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *RemoteRenderingAccountsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *RemoteRenderingAccountsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MixedReality/remoteRenderingAccounts/{accountName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Retrieve a Remote Rendering Account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-01-preview
//   - resourceGroupName - Name of an Azure resource group.
//   - accountName - Name of an Mixed Reality Account.
//   - options - RemoteRenderingAccountsClientGetOptions contains the optional parameters for the RemoteRenderingAccountsClient.Get
//     method.
func (client *RemoteRenderingAccountsClient) Get(ctx context.Context, resourceGroupName string, accountName string, options *RemoteRenderingAccountsClientGetOptions) (RemoteRenderingAccountsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, options)
	if err != nil {
		return RemoteRenderingAccountsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RemoteRenderingAccountsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RemoteRenderingAccountsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *RemoteRenderingAccountsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *RemoteRenderingAccountsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MixedReality/remoteRenderingAccounts/{accountName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RemoteRenderingAccountsClient) getHandleResponse(resp *http.Response) (RemoteRenderingAccountsClientGetResponse, error) {
	result := RemoteRenderingAccountsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RemoteRenderingAccount); err != nil {
		return RemoteRenderingAccountsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - List Resources by Resource Group
//
// Generated from API version 2021-03-01-preview
//   - resourceGroupName - Name of an Azure resource group.
//   - options - RemoteRenderingAccountsClientListByResourceGroupOptions contains the optional parameters for the RemoteRenderingAccountsClient.NewListByResourceGroupPager
//     method.
func (client *RemoteRenderingAccountsClient) NewListByResourceGroupPager(resourceGroupName string, options *RemoteRenderingAccountsClientListByResourceGroupOptions) (*runtime.Pager[RemoteRenderingAccountsClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[RemoteRenderingAccountsClientListByResourceGroupResponse]{
		More: func(page RemoteRenderingAccountsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RemoteRenderingAccountsClientListByResourceGroupResponse) (RemoteRenderingAccountsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return RemoteRenderingAccountsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return RemoteRenderingAccountsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RemoteRenderingAccountsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *RemoteRenderingAccountsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *RemoteRenderingAccountsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MixedReality/remoteRenderingAccounts"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *RemoteRenderingAccountsClient) listByResourceGroupHandleResponse(resp *http.Response) (RemoteRenderingAccountsClientListByResourceGroupResponse, error) {
	result := RemoteRenderingAccountsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RemoteRenderingAccountPage); err != nil {
		return RemoteRenderingAccountsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - List Remote Rendering Accounts by Subscription
//
// Generated from API version 2021-03-01-preview
//   - options - RemoteRenderingAccountsClientListBySubscriptionOptions contains the optional parameters for the RemoteRenderingAccountsClient.NewListBySubscriptionPager
//     method.
func (client *RemoteRenderingAccountsClient) NewListBySubscriptionPager(options *RemoteRenderingAccountsClientListBySubscriptionOptions) (*runtime.Pager[RemoteRenderingAccountsClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[RemoteRenderingAccountsClientListBySubscriptionResponse]{
		More: func(page RemoteRenderingAccountsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RemoteRenderingAccountsClientListBySubscriptionResponse) (RemoteRenderingAccountsClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return RemoteRenderingAccountsClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return RemoteRenderingAccountsClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RemoteRenderingAccountsClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *RemoteRenderingAccountsClient) listBySubscriptionCreateRequest(ctx context.Context, options *RemoteRenderingAccountsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.MixedReality/remoteRenderingAccounts"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *RemoteRenderingAccountsClient) listBySubscriptionHandleResponse(resp *http.Response) (RemoteRenderingAccountsClientListBySubscriptionResponse, error) {
	result := RemoteRenderingAccountsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RemoteRenderingAccountPage); err != nil {
		return RemoteRenderingAccountsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// ListKeys - List Both of the 2 Keys of a Remote Rendering Account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-01-preview
//   - resourceGroupName - Name of an Azure resource group.
//   - accountName - Name of an Mixed Reality Account.
//   - options - RemoteRenderingAccountsClientListKeysOptions contains the optional parameters for the RemoteRenderingAccountsClient.ListKeys
//     method.
func (client *RemoteRenderingAccountsClient) ListKeys(ctx context.Context, resourceGroupName string, accountName string, options *RemoteRenderingAccountsClientListKeysOptions) (RemoteRenderingAccountsClientListKeysResponse, error) {
	var err error
	req, err := client.listKeysCreateRequest(ctx, resourceGroupName, accountName, options)
	if err != nil {
		return RemoteRenderingAccountsClientListKeysResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RemoteRenderingAccountsClientListKeysResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RemoteRenderingAccountsClientListKeysResponse{}, err
	}
	resp, err := client.listKeysHandleResponse(httpResp)
	return resp, err
}

// listKeysCreateRequest creates the ListKeys request.
func (client *RemoteRenderingAccountsClient) listKeysCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *RemoteRenderingAccountsClientListKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MixedReality/remoteRenderingAccounts/{accountName}/listKeys"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listKeysHandleResponse handles the ListKeys response.
func (client *RemoteRenderingAccountsClient) listKeysHandleResponse(resp *http.Response) (RemoteRenderingAccountsClientListKeysResponse, error) {
	result := RemoteRenderingAccountsClientListKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccountKeys); err != nil {
		return RemoteRenderingAccountsClientListKeysResponse{}, err
	}
	return result, nil
}

// RegenerateKeys - Regenerate specified Key of a Remote Rendering Account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-01-preview
//   - resourceGroupName - Name of an Azure resource group.
//   - accountName - Name of an Mixed Reality Account.
//   - regenerate - Required information for key regeneration.
//   - options - RemoteRenderingAccountsClientRegenerateKeysOptions contains the optional parameters for the RemoteRenderingAccountsClient.RegenerateKeys
//     method.
func (client *RemoteRenderingAccountsClient) RegenerateKeys(ctx context.Context, resourceGroupName string, accountName string, regenerate AccountKeyRegenerateRequest, options *RemoteRenderingAccountsClientRegenerateKeysOptions) (RemoteRenderingAccountsClientRegenerateKeysResponse, error) {
	var err error
	req, err := client.regenerateKeysCreateRequest(ctx, resourceGroupName, accountName, regenerate, options)
	if err != nil {
		return RemoteRenderingAccountsClientRegenerateKeysResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RemoteRenderingAccountsClientRegenerateKeysResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RemoteRenderingAccountsClientRegenerateKeysResponse{}, err
	}
	resp, err := client.regenerateKeysHandleResponse(httpResp)
	return resp, err
}

// regenerateKeysCreateRequest creates the RegenerateKeys request.
func (client *RemoteRenderingAccountsClient) regenerateKeysCreateRequest(ctx context.Context, resourceGroupName string, accountName string, regenerate AccountKeyRegenerateRequest, options *RemoteRenderingAccountsClientRegenerateKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MixedReality/remoteRenderingAccounts/{accountName}/regenerateKeys"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, regenerate); err != nil {
	return nil, err
}
	return req, nil
}

// regenerateKeysHandleResponse handles the RegenerateKeys response.
func (client *RemoteRenderingAccountsClient) regenerateKeysHandleResponse(resp *http.Response) (RemoteRenderingAccountsClientRegenerateKeysResponse, error) {
	result := RemoteRenderingAccountsClientRegenerateKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccountKeys); err != nil {
		return RemoteRenderingAccountsClientRegenerateKeysResponse{}, err
	}
	return result, nil
}

// Update - Updating a Remote Rendering Account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-03-01-preview
//   - resourceGroupName - Name of an Azure resource group.
//   - accountName - Name of an Mixed Reality Account.
//   - remoteRenderingAccount - Remote Rendering Account parameter.
//   - options - RemoteRenderingAccountsClientUpdateOptions contains the optional parameters for the RemoteRenderingAccountsClient.Update
//     method.
func (client *RemoteRenderingAccountsClient) Update(ctx context.Context, resourceGroupName string, accountName string, remoteRenderingAccount RemoteRenderingAccount, options *RemoteRenderingAccountsClientUpdateOptions) (RemoteRenderingAccountsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, accountName, remoteRenderingAccount, options)
	if err != nil {
		return RemoteRenderingAccountsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RemoteRenderingAccountsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RemoteRenderingAccountsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *RemoteRenderingAccountsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, remoteRenderingAccount RemoteRenderingAccount, options *RemoteRenderingAccountsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MixedReality/remoteRenderingAccounts/{accountName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, remoteRenderingAccount); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *RemoteRenderingAccountsClient) updateHandleResponse(resp *http.Response) (RemoteRenderingAccountsClientUpdateResponse, error) {
	result := RemoteRenderingAccountsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RemoteRenderingAccount); err != nil {
		return RemoteRenderingAccountsClientUpdateResponse{}, err
	}
	return result, nil
}

