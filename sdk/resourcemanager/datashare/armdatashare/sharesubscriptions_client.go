//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatashare

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

// ShareSubscriptionsClient contains the methods for the ShareSubscriptions group.
// Don't use this type directly, use NewShareSubscriptionsClient() instead.
type ShareSubscriptionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewShareSubscriptionsClient creates a new instance of ShareSubscriptionsClient with the specified values.
//   - subscriptionID - The subscription identifier
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewShareSubscriptionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ShareSubscriptionsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ShareSubscriptionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCancelSynchronization - Request to cancel a synchronization.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of the shareSubscription.
//   - shareSubscriptionSynchronization - Share Subscription Synchronization payload.
//   - options - ShareSubscriptionsClientBeginCancelSynchronizationOptions contains the optional parameters for the ShareSubscriptionsClient.BeginCancelSynchronization
//     method.
func (client *ShareSubscriptionsClient) BeginCancelSynchronization(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, shareSubscriptionSynchronization ShareSubscriptionSynchronization, options *ShareSubscriptionsClientBeginCancelSynchronizationOptions) (*runtime.Poller[ShareSubscriptionsClientCancelSynchronizationResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.cancelSynchronization(ctx, resourceGroupName, accountName, shareSubscriptionName, shareSubscriptionSynchronization, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ShareSubscriptionsClientCancelSynchronizationResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ShareSubscriptionsClientCancelSynchronizationResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CancelSynchronization - Request to cancel a synchronization.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
func (client *ShareSubscriptionsClient) cancelSynchronization(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, shareSubscriptionSynchronization ShareSubscriptionSynchronization, options *ShareSubscriptionsClientBeginCancelSynchronizationOptions) (*http.Response, error) {
	var err error
	const operationName = "ShareSubscriptionsClient.BeginCancelSynchronization"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.cancelSynchronizationCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, shareSubscriptionSynchronization, options)
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

// cancelSynchronizationCreateRequest creates the CancelSynchronization request.
func (client *ShareSubscriptionsClient) cancelSynchronizationCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, shareSubscriptionSynchronization ShareSubscriptionSynchronization, options *ShareSubscriptionsClientBeginCancelSynchronizationOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}/cancelSynchronization"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, shareSubscriptionSynchronization); err != nil {
		return nil, err
	}
	return req, nil
}

// Create - Create a shareSubscription in an account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of the shareSubscription.
//   - shareSubscription - create parameters for shareSubscription
//   - options - ShareSubscriptionsClientCreateOptions contains the optional parameters for the ShareSubscriptionsClient.Create
//     method.
func (client *ShareSubscriptionsClient) Create(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, shareSubscription ShareSubscription, options *ShareSubscriptionsClientCreateOptions) (ShareSubscriptionsClientCreateResponse, error) {
	var err error
	const operationName = "ShareSubscriptionsClient.Create"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, shareSubscription, options)
	if err != nil {
		return ShareSubscriptionsClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ShareSubscriptionsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return ShareSubscriptionsClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *ShareSubscriptionsClient) createCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, shareSubscription ShareSubscription, options *ShareSubscriptionsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, shareSubscription); err != nil {
		return nil, err
	}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *ShareSubscriptionsClient) createHandleResponse(resp *http.Response) (ShareSubscriptionsClientCreateResponse, error) {
	result := ShareSubscriptionsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ShareSubscription); err != nil {
		return ShareSubscriptionsClientCreateResponse{}, err
	}
	return result, nil
}

// BeginDelete - Delete a shareSubscription in an account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of the shareSubscription.
//   - options - ShareSubscriptionsClientBeginDeleteOptions contains the optional parameters for the ShareSubscriptionsClient.BeginDelete
//     method.
func (client *ShareSubscriptionsClient) BeginDelete(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientBeginDeleteOptions) (*runtime.Poller[ShareSubscriptionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, accountName, shareSubscriptionName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ShareSubscriptionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ShareSubscriptionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a shareSubscription in an account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
func (client *ShareSubscriptionsClient) deleteOperation(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ShareSubscriptionsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, options)
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
func (client *ShareSubscriptionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a shareSubscription in an account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of the shareSubscription.
//   - options - ShareSubscriptionsClientGetOptions contains the optional parameters for the ShareSubscriptionsClient.Get method.
func (client *ShareSubscriptionsClient) Get(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientGetOptions) (ShareSubscriptionsClientGetResponse, error) {
	var err error
	const operationName = "ShareSubscriptionsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, options)
	if err != nil {
		return ShareSubscriptionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ShareSubscriptionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ShareSubscriptionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ShareSubscriptionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ShareSubscriptionsClient) getHandleResponse(resp *http.Response) (ShareSubscriptionsClientGetResponse, error) {
	result := ShareSubscriptionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ShareSubscription); err != nil {
		return ShareSubscriptionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByAccountPager - List share subscriptions in an account
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - options - ShareSubscriptionsClientListByAccountOptions contains the optional parameters for the ShareSubscriptionsClient.NewListByAccountPager
//     method.
func (client *ShareSubscriptionsClient) NewListByAccountPager(resourceGroupName string, accountName string, options *ShareSubscriptionsClientListByAccountOptions) *runtime.Pager[ShareSubscriptionsClientListByAccountResponse] {
	return runtime.NewPager(runtime.PagingHandler[ShareSubscriptionsClientListByAccountResponse]{
		More: func(page ShareSubscriptionsClientListByAccountResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ShareSubscriptionsClientListByAccountResponse) (ShareSubscriptionsClientListByAccountResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ShareSubscriptionsClient.NewListByAccountPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByAccountCreateRequest(ctx, resourceGroupName, accountName, options)
			}, nil)
			if err != nil {
				return ShareSubscriptionsClientListByAccountResponse{}, err
			}
			return client.listByAccountHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByAccountCreateRequest creates the ListByAccount request.
func (client *ShareSubscriptionsClient) listByAccountCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *ShareSubscriptionsClientListByAccountOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions"
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
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByAccountHandleResponse handles the ListByAccount response.
func (client *ShareSubscriptionsClient) listByAccountHandleResponse(resp *http.Response) (ShareSubscriptionsClientListByAccountResponse, error) {
	result := ShareSubscriptionsClientListByAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ShareSubscriptionList); err != nil {
		return ShareSubscriptionsClientListByAccountResponse{}, err
	}
	return result, nil
}

// NewListSourceShareSynchronizationSettingsPager - Get synchronization settings set on a share
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of the shareSubscription.
//   - options - ShareSubscriptionsClientListSourceShareSynchronizationSettingsOptions contains the optional parameters for the
//     ShareSubscriptionsClient.NewListSourceShareSynchronizationSettingsPager method.
func (client *ShareSubscriptionsClient) NewListSourceShareSynchronizationSettingsPager(resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientListSourceShareSynchronizationSettingsOptions) *runtime.Pager[ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse]{
		More: func(page ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse) (ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ShareSubscriptionsClient.NewListSourceShareSynchronizationSettingsPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listSourceShareSynchronizationSettingsCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, options)
			}, nil)
			if err != nil {
				return ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse{}, err
			}
			return client.listSourceShareSynchronizationSettingsHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listSourceShareSynchronizationSettingsCreateRequest creates the ListSourceShareSynchronizationSettings request.
func (client *ShareSubscriptionsClient) listSourceShareSynchronizationSettingsCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientListSourceShareSynchronizationSettingsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}/listSourceShareSynchronizationSettings"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSourceShareSynchronizationSettingsHandleResponse handles the ListSourceShareSynchronizationSettings response.
func (client *ShareSubscriptionsClient) listSourceShareSynchronizationSettingsHandleResponse(resp *http.Response) (ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse, error) {
	result := ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SourceShareSynchronizationSettingList); err != nil {
		return ShareSubscriptionsClientListSourceShareSynchronizationSettingsResponse{}, err
	}
	return result, nil
}

// NewListSynchronizationDetailsPager - List synchronization details
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of the share subscription.
//   - shareSubscriptionSynchronization - Share Subscription Synchronization payload.
//   - options - ShareSubscriptionsClientListSynchronizationDetailsOptions contains the optional parameters for the ShareSubscriptionsClient.NewListSynchronizationDetailsPager
//     method.
func (client *ShareSubscriptionsClient) NewListSynchronizationDetailsPager(resourceGroupName string, accountName string, shareSubscriptionName string, shareSubscriptionSynchronization ShareSubscriptionSynchronization, options *ShareSubscriptionsClientListSynchronizationDetailsOptions) *runtime.Pager[ShareSubscriptionsClientListSynchronizationDetailsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ShareSubscriptionsClientListSynchronizationDetailsResponse]{
		More: func(page ShareSubscriptionsClientListSynchronizationDetailsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ShareSubscriptionsClientListSynchronizationDetailsResponse) (ShareSubscriptionsClientListSynchronizationDetailsResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ShareSubscriptionsClient.NewListSynchronizationDetailsPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listSynchronizationDetailsCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, shareSubscriptionSynchronization, options)
			}, nil)
			if err != nil {
				return ShareSubscriptionsClientListSynchronizationDetailsResponse{}, err
			}
			return client.listSynchronizationDetailsHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listSynchronizationDetailsCreateRequest creates the ListSynchronizationDetails request.
func (client *ShareSubscriptionsClient) listSynchronizationDetailsCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, shareSubscriptionSynchronization ShareSubscriptionSynchronization, options *ShareSubscriptionsClientListSynchronizationDetailsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}/listSynchronizationDetails"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, shareSubscriptionSynchronization); err != nil {
		return nil, err
	}
	return req, nil
}

// listSynchronizationDetailsHandleResponse handles the ListSynchronizationDetails response.
func (client *ShareSubscriptionsClient) listSynchronizationDetailsHandleResponse(resp *http.Response) (ShareSubscriptionsClientListSynchronizationDetailsResponse, error) {
	result := ShareSubscriptionsClientListSynchronizationDetailsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SynchronizationDetailsList); err != nil {
		return ShareSubscriptionsClientListSynchronizationDetailsResponse{}, err
	}
	return result, nil
}

// NewListSynchronizationsPager - List synchronizations of a share subscription
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of the share subscription.
//   - options - ShareSubscriptionsClientListSynchronizationsOptions contains the optional parameters for the ShareSubscriptionsClient.NewListSynchronizationsPager
//     method.
func (client *ShareSubscriptionsClient) NewListSynchronizationsPager(resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientListSynchronizationsOptions) *runtime.Pager[ShareSubscriptionsClientListSynchronizationsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ShareSubscriptionsClientListSynchronizationsResponse]{
		More: func(page ShareSubscriptionsClientListSynchronizationsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ShareSubscriptionsClientListSynchronizationsResponse) (ShareSubscriptionsClientListSynchronizationsResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ShareSubscriptionsClient.NewListSynchronizationsPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listSynchronizationsCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, options)
			}, nil)
			if err != nil {
				return ShareSubscriptionsClientListSynchronizationsResponse{}, err
			}
			return client.listSynchronizationsHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listSynchronizationsCreateRequest creates the ListSynchronizations request.
func (client *ShareSubscriptionsClient) listSynchronizationsCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, options *ShareSubscriptionsClientListSynchronizationsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}/listSynchronizations"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listSynchronizationsHandleResponse handles the ListSynchronizations response.
func (client *ShareSubscriptionsClient) listSynchronizationsHandleResponse(resp *http.Response) (ShareSubscriptionsClientListSynchronizationsResponse, error) {
	result := ShareSubscriptionsClientListSynchronizationsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ShareSubscriptionSynchronizationList); err != nil {
		return ShareSubscriptionsClientListSynchronizationsResponse{}, err
	}
	return result, nil
}

// BeginSynchronize - Initiate a copy
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareSubscriptionName - The name of share subscription
//   - synchronize - Synchronize payload
//   - options - ShareSubscriptionsClientBeginSynchronizeOptions contains the optional parameters for the ShareSubscriptionsClient.BeginSynchronize
//     method.
func (client *ShareSubscriptionsClient) BeginSynchronize(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, synchronize Synchronize, options *ShareSubscriptionsClientBeginSynchronizeOptions) (*runtime.Poller[ShareSubscriptionsClientSynchronizeResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.synchronize(ctx, resourceGroupName, accountName, shareSubscriptionName, synchronize, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ShareSubscriptionsClientSynchronizeResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ShareSubscriptionsClientSynchronizeResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Synchronize - Initiate a copy
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
func (client *ShareSubscriptionsClient) synchronize(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, synchronize Synchronize, options *ShareSubscriptionsClientBeginSynchronizeOptions) (*http.Response, error) {
	var err error
	const operationName = "ShareSubscriptionsClient.BeginSynchronize"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.synchronizeCreateRequest(ctx, resourceGroupName, accountName, shareSubscriptionName, synchronize, options)
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

// synchronizeCreateRequest creates the Synchronize request.
func (client *ShareSubscriptionsClient) synchronizeCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareSubscriptionName string, synchronize Synchronize, options *ShareSubscriptionsClientBeginSynchronizeOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shareSubscriptions/{shareSubscriptionName}/synchronize"
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
	if shareSubscriptionName == "" {
		return nil, errors.New("parameter shareSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareSubscriptionName}", url.PathEscape(shareSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, synchronize); err != nil {
		return nil, err
	}
	return req, nil
}
