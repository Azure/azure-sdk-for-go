//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridnetwork

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
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PublishersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a publisher.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - publisherName - The name of the publisher.
//   - options - PublishersClientBeginCreateOrUpdateOptions contains the optional parameters for the PublishersClient.BeginCreateOrUpdate
//     method.
func (client *PublishersClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientBeginCreateOrUpdateOptions) (*runtime.Poller[PublishersClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, publisherName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PublishersClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[PublishersClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates or updates a publisher.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
func (client *PublishersClient) createOrUpdate(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "PublishersClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, publisherName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PublishersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/publishers/{publisherName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Parameters != nil {
		if err := runtime.MarshalAsJSON(req, *options.Parameters); err != nil {
			return nil, err
		}
		return req, nil
	}
	return req, nil
}

// BeginDelete - Deletes the specified publisher.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - publisherName - The name of the publisher.
//   - options - PublishersClientBeginDeleteOptions contains the optional parameters for the PublishersClient.BeginDelete method.
func (client *PublishersClient) BeginDelete(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientBeginDeleteOptions) (*runtime.Poller[PublishersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, publisherName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PublishersClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[PublishersClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes the specified publisher.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
func (client *PublishersClient) deleteOperation(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "PublishersClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, publisherName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *PublishersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/publishers/{publisherName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets information about the specified publisher.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - publisherName - The name of the publisher.
//   - options - PublishersClientGetOptions contains the optional parameters for the PublishersClient.Get method.
func (client *PublishersClient) Get(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientGetOptions) (PublishersClientGetResponse, error) {
	var err error
	const operationName = "PublishersClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, publisherName, options)
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
func (client *PublishersClient) getCreateRequest(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/publishers/{publisherName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01")
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

// NewListByResourceGroupPager - Lists all the publishers in a resource group.
//
// Generated from API version 2023-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - PublishersClientListByResourceGroupOptions contains the optional parameters for the PublishersClient.NewListByResourceGroupPager
//     method.
func (client *PublishersClient) NewListByResourceGroupPager(resourceGroupName string, options *PublishersClientListByResourceGroupOptions) *runtime.Pager[PublishersClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[PublishersClientListByResourceGroupResponse]{
		More: func(page PublishersClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PublishersClientListByResourceGroupResponse) (PublishersClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PublishersClient.NewListByResourceGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			}, nil)
			if err != nil {
				return PublishersClientListByResourceGroupResponse{}, err
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *PublishersClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *PublishersClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/publishers"
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
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *PublishersClient) listByResourceGroupHandleResponse(resp *http.Response) (PublishersClientListByResourceGroupResponse, error) {
	result := PublishersClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublisherListResult); err != nil {
		return PublishersClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Lists all the publishers in a subscription.
//
// Generated from API version 2023-09-01
//   - options - PublishersClientListBySubscriptionOptions contains the optional parameters for the PublishersClient.NewListBySubscriptionPager
//     method.
func (client *PublishersClient) NewListBySubscriptionPager(options *PublishersClientListBySubscriptionOptions) *runtime.Pager[PublishersClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[PublishersClientListBySubscriptionResponse]{
		More: func(page PublishersClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PublishersClientListBySubscriptionResponse) (PublishersClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PublishersClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return PublishersClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *PublishersClient) listBySubscriptionCreateRequest(ctx context.Context, options *PublishersClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/publishers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01")
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

// Update - Update a publisher resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - publisherName - The name of the publisher.
//   - options - PublishersClientUpdateOptions contains the optional parameters for the PublishersClient.Update method.
func (client *PublishersClient) Update(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientUpdateOptions) (PublishersClientUpdateResponse, error) {
	var err error
	const operationName = "PublishersClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, publisherName, options)
	if err != nil {
		return PublishersClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PublishersClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PublishersClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *PublishersClient) updateCreateRequest(ctx context.Context, resourceGroupName string, publisherName string, options *PublishersClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/publishers/{publisherName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publisherName == "" {
		return nil, errors.New("parameter publisherName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publisherName}", url.PathEscape(publisherName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Parameters != nil {
		if err := runtime.MarshalAsJSON(req, *options.Parameters); err != nil {
			return nil, err
		}
		return req, nil
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *PublishersClient) updateHandleResponse(resp *http.Response) (PublishersClientUpdateResponse, error) {
	result := PublishersClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Publisher); err != nil {
		return PublishersClientUpdateResponse{}, err
	}
	return result, nil
}
