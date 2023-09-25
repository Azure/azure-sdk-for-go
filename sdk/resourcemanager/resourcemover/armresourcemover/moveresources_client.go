//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armresourcemover

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

// MoveResourcesClient contains the methods for the MoveResources group.
// Don't use this type directly, use NewMoveResourcesClient() instead.
type MoveResourcesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewMoveResourcesClient creates a new instance of MoveResourcesClient with the specified values.
//   - subscriptionID - The Subscription ID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewMoveResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MoveResourcesClient, error) {
	cl, err := arm.NewClient(moduleName+".MoveResourcesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &MoveResourcesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreate - Creates or updates a Move Resource in the move collection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01
//   - resourceGroupName - The Resource Group Name.
//   - moveCollectionName - The Move Collection Name.
//   - moveResourceName - The Move Resource Name.
//   - options - MoveResourcesClientBeginCreateOptions contains the optional parameters for the MoveResourcesClient.BeginCreate
//     method.
func (client *MoveResourcesClient) BeginCreate(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientBeginCreateOptions) (*runtime.Poller[MoveResourcesClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, moveCollectionName, moveResourceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[MoveResourcesClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[MoveResourcesClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Creates or updates a Move Resource in the move collection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01
func (client *MoveResourcesClient) create(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientBeginCreateOptions) (*http.Response, error) {
	var err error
	req, err := client.createCreateRequest(ctx, resourceGroupName, moveCollectionName, moveResourceName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createCreateRequest creates the Create request.
func (client *MoveResourcesClient) createCreateRequest(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/moveCollections/{moveCollectionName}/moveResources/{moveResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if moveCollectionName == "" {
		return nil, errors.New("parameter moveCollectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{moveCollectionName}", url.PathEscape(moveCollectionName))
	if moveResourceName == "" {
		return nil, errors.New("parameter moveResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{moveResourceName}", url.PathEscape(moveResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Body != nil {
		if err := runtime.MarshalAsJSON(req, *options.Body); err != nil {
	return nil, err
}
		return req, nil
	}
	return req, nil
}

// BeginDelete - Deletes a Move Resource from the move collection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01
//   - resourceGroupName - The Resource Group Name.
//   - moveCollectionName - The Move Collection Name.
//   - moveResourceName - The Move Resource Name.
//   - options - MoveResourcesClientBeginDeleteOptions contains the optional parameters for the MoveResourcesClient.BeginDelete
//     method.
func (client *MoveResourcesClient) BeginDelete(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientBeginDeleteOptions) (*runtime.Poller[MoveResourcesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, moveCollectionName, moveResourceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[MoveResourcesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[MoveResourcesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a Move Resource from the move collection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01
func (client *MoveResourcesClient) deleteOperation(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, moveCollectionName, moveResourceName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *MoveResourcesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/moveCollections/{moveCollectionName}/moveResources/{moveResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if moveCollectionName == "" {
		return nil, errors.New("parameter moveCollectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{moveCollectionName}", url.PathEscape(moveCollectionName))
	if moveResourceName == "" {
		return nil, errors.New("parameter moveResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{moveResourceName}", url.PathEscape(moveResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the Move Resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01
//   - resourceGroupName - The Resource Group Name.
//   - moveCollectionName - The Move Collection Name.
//   - moveResourceName - The Move Resource Name.
//   - options - MoveResourcesClientGetOptions contains the optional parameters for the MoveResourcesClient.Get method.
func (client *MoveResourcesClient) Get(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientGetOptions) (MoveResourcesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, moveCollectionName, moveResourceName, options)
	if err != nil {
		return MoveResourcesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return MoveResourcesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return MoveResourcesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *MoveResourcesClient) getCreateRequest(ctx context.Context, resourceGroupName string, moveCollectionName string, moveResourceName string, options *MoveResourcesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/moveCollections/{moveCollectionName}/moveResources/{moveResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if moveCollectionName == "" {
		return nil, errors.New("parameter moveCollectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{moveCollectionName}", url.PathEscape(moveCollectionName))
	if moveResourceName == "" {
		return nil, errors.New("parameter moveResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{moveResourceName}", url.PathEscape(moveResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *MoveResourcesClient) getHandleResponse(resp *http.Response) (MoveResourcesClientGetResponse, error) {
	result := MoveResourcesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MoveResource); err != nil {
		return MoveResourcesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists the Move Resources in the move collection.
//
// Generated from API version 2021-08-01
//   - resourceGroupName - The Resource Group Name.
//   - moveCollectionName - The Move Collection Name.
//   - options - MoveResourcesClientListOptions contains the optional parameters for the MoveResourcesClient.NewListPager method.
func (client *MoveResourcesClient) NewListPager(resourceGroupName string, moveCollectionName string, options *MoveResourcesClientListOptions) (*runtime.Pager[MoveResourcesClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[MoveResourcesClientListResponse]{
		More: func(page MoveResourcesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *MoveResourcesClientListResponse) (MoveResourcesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, moveCollectionName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return MoveResourcesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return MoveResourcesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return MoveResourcesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *MoveResourcesClient) listCreateRequest(ctx context.Context, resourceGroupName string, moveCollectionName string, options *MoveResourcesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Migrate/moveCollections/{moveCollectionName}/moveResources"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if moveCollectionName == "" {
		return nil, errors.New("parameter moveCollectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{moveCollectionName}", url.PathEscape(moveCollectionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *MoveResourcesClient) listHandleResponse(resp *http.Response) (MoveResourcesClientListResponse, error) {
	result := MoveResourcesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MoveResourceCollection); err != nil {
		return MoveResourcesClientListResponse{}, err
	}
	return result, nil
}

