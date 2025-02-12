//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstorage

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

// DeletedAccountsClient contains the methods for the DeletedAccounts group.
// Don't use this type directly, use NewDeletedAccountsClient() instead.
type DeletedAccountsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewDeletedAccountsClient creates a new instance of DeletedAccountsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewDeletedAccountsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DeletedAccountsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &DeletedAccountsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Get properties of specified deleted account resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-01-01
//   - deletedAccountName - Name of the deleted storage account.
//   - location - The location of the deleted storage account.
//   - options - DeletedAccountsClientGetOptions contains the optional parameters for the DeletedAccountsClient.Get method.
func (client *DeletedAccountsClient) Get(ctx context.Context, deletedAccountName string, location string, options *DeletedAccountsClientGetOptions) (DeletedAccountsClientGetResponse, error) {
	var err error
	const operationName = "DeletedAccountsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, deletedAccountName, location, options)
	if err != nil {
		return DeletedAccountsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeletedAccountsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DeletedAccountsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *DeletedAccountsClient) getCreateRequest(ctx context.Context, deletedAccountName string, location string, options *DeletedAccountsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Storage/locations/{location}/deletedAccounts/{deletedAccountName}"
	if deletedAccountName == "" {
		return nil, errors.New("parameter deletedAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{deletedAccountName}", url.PathEscape(deletedAccountName))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DeletedAccountsClient) getHandleResponse(resp *http.Response) (DeletedAccountsClientGetResponse, error) {
	result := DeletedAccountsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedAccount); err != nil {
		return DeletedAccountsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists deleted accounts under the subscription.
//
// Generated from API version 2024-01-01
//   - options - DeletedAccountsClientListOptions contains the optional parameters for the DeletedAccountsClient.NewListPager
//     method.
func (client *DeletedAccountsClient) NewListPager(options *DeletedAccountsClientListOptions) *runtime.Pager[DeletedAccountsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[DeletedAccountsClientListResponse]{
		More: func(page DeletedAccountsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DeletedAccountsClientListResponse) (DeletedAccountsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "DeletedAccountsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return DeletedAccountsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *DeletedAccountsClient) listCreateRequest(ctx context.Context, options *DeletedAccountsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Storage/deletedAccounts"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *DeletedAccountsClient) listHandleResponse(resp *http.Response) (DeletedAccountsClientListResponse, error) {
	result := DeletedAccountsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedAccountListResult); err != nil {
		return DeletedAccountsClientListResponse{}, err
	}
	return result, nil
}
