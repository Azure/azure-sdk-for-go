// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armstandbypool

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

// StandbyVirtualMachinePoolsClient contains the methods for the StandbyVirtualMachinePools group.
// Don't use this type directly, use NewStandbyVirtualMachinePoolsClient() instead.
type StandbyVirtualMachinePoolsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewStandbyVirtualMachinePoolsClient creates a new instance of StandbyVirtualMachinePoolsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewStandbyVirtualMachinePoolsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*StandbyVirtualMachinePoolsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &StandbyVirtualMachinePoolsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a StandbyVirtualMachinePoolResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - standbyVirtualMachinePoolName - Name of the standby virtual machine pool
//   - resource - Resource create parameters.
//   - options - StandbyVirtualMachinePoolsClientBeginCreateOrUpdateOptions contains the optional parameters for the StandbyVirtualMachinePoolsClient.BeginCreateOrUpdate
//     method.
func (client *StandbyVirtualMachinePoolsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, resource StandbyVirtualMachinePoolResource, options *StandbyVirtualMachinePoolsClientBeginCreateOrUpdateOptions) (*runtime.Poller[StandbyVirtualMachinePoolsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, standbyVirtualMachinePoolName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[StandbyVirtualMachinePoolsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[StandbyVirtualMachinePoolsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create a StandbyVirtualMachinePoolResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
func (client *StandbyVirtualMachinePoolsClient) createOrUpdate(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, resource StandbyVirtualMachinePoolResource, options *StandbyVirtualMachinePoolsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "StandbyVirtualMachinePoolsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, standbyVirtualMachinePoolName, resource, options)
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
func (client *StandbyVirtualMachinePoolsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, resource StandbyVirtualMachinePoolResource, _ *StandbyVirtualMachinePoolsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StandbyPool/standbyVirtualMachinePools/{standbyVirtualMachinePoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if standbyVirtualMachinePoolName == "" {
		return nil, errors.New("parameter standbyVirtualMachinePoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{standbyVirtualMachinePoolName}", url.PathEscape(standbyVirtualMachinePoolName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete a StandbyVirtualMachinePoolResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - standbyVirtualMachinePoolName - Name of the standby virtual machine pool
//   - options - StandbyVirtualMachinePoolsClientBeginDeleteOptions contains the optional parameters for the StandbyVirtualMachinePoolsClient.BeginDelete
//     method.
func (client *StandbyVirtualMachinePoolsClient) BeginDelete(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, options *StandbyVirtualMachinePoolsClientBeginDeleteOptions) (*runtime.Poller[StandbyVirtualMachinePoolsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, standbyVirtualMachinePoolName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[StandbyVirtualMachinePoolsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[StandbyVirtualMachinePoolsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a StandbyVirtualMachinePoolResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
func (client *StandbyVirtualMachinePoolsClient) deleteOperation(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, options *StandbyVirtualMachinePoolsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "StandbyVirtualMachinePoolsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, standbyVirtualMachinePoolName, options)
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
func (client *StandbyVirtualMachinePoolsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, _ *StandbyVirtualMachinePoolsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StandbyPool/standbyVirtualMachinePools/{standbyVirtualMachinePoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if standbyVirtualMachinePoolName == "" {
		return nil, errors.New("parameter standbyVirtualMachinePoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{standbyVirtualMachinePoolName}", url.PathEscape(standbyVirtualMachinePoolName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a StandbyVirtualMachinePoolResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - standbyVirtualMachinePoolName - Name of the standby virtual machine pool
//   - options - StandbyVirtualMachinePoolsClientGetOptions contains the optional parameters for the StandbyVirtualMachinePoolsClient.Get
//     method.
func (client *StandbyVirtualMachinePoolsClient) Get(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, options *StandbyVirtualMachinePoolsClientGetOptions) (StandbyVirtualMachinePoolsClientGetResponse, error) {
	var err error
	const operationName = "StandbyVirtualMachinePoolsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, standbyVirtualMachinePoolName, options)
	if err != nil {
		return StandbyVirtualMachinePoolsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return StandbyVirtualMachinePoolsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return StandbyVirtualMachinePoolsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *StandbyVirtualMachinePoolsClient) getCreateRequest(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, _ *StandbyVirtualMachinePoolsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StandbyPool/standbyVirtualMachinePools/{standbyVirtualMachinePoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if standbyVirtualMachinePoolName == "" {
		return nil, errors.New("parameter standbyVirtualMachinePoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{standbyVirtualMachinePoolName}", url.PathEscape(standbyVirtualMachinePoolName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *StandbyVirtualMachinePoolsClient) getHandleResponse(resp *http.Response) (StandbyVirtualMachinePoolsClientGetResponse, error) {
	result := StandbyVirtualMachinePoolsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StandbyVirtualMachinePoolResource); err != nil {
		return StandbyVirtualMachinePoolsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - List StandbyVirtualMachinePoolResource resources by resource group
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - StandbyVirtualMachinePoolsClientListByResourceGroupOptions contains the optional parameters for the StandbyVirtualMachinePoolsClient.NewListByResourceGroupPager
//     method.
func (client *StandbyVirtualMachinePoolsClient) NewListByResourceGroupPager(resourceGroupName string, options *StandbyVirtualMachinePoolsClientListByResourceGroupOptions) *runtime.Pager[StandbyVirtualMachinePoolsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[StandbyVirtualMachinePoolsClientListByResourceGroupResponse]{
		More: func(page StandbyVirtualMachinePoolsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *StandbyVirtualMachinePoolsClientListByResourceGroupResponse) (StandbyVirtualMachinePoolsClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "StandbyVirtualMachinePoolsClient.NewListByResourceGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			}, nil)
			if err != nil {
				return StandbyVirtualMachinePoolsClientListByResourceGroupResponse{}, err
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *StandbyVirtualMachinePoolsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, _ *StandbyVirtualMachinePoolsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StandbyPool/standbyVirtualMachinePools"
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
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *StandbyVirtualMachinePoolsClient) listByResourceGroupHandleResponse(resp *http.Response) (StandbyVirtualMachinePoolsClientListByResourceGroupResponse, error) {
	result := StandbyVirtualMachinePoolsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StandbyVirtualMachinePoolResourceListResult); err != nil {
		return StandbyVirtualMachinePoolsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - List StandbyVirtualMachinePoolResource resources by subscription ID
//
// Generated from API version 2024-03-01
//   - options - StandbyVirtualMachinePoolsClientListBySubscriptionOptions contains the optional parameters for the StandbyVirtualMachinePoolsClient.NewListBySubscriptionPager
//     method.
func (client *StandbyVirtualMachinePoolsClient) NewListBySubscriptionPager(options *StandbyVirtualMachinePoolsClientListBySubscriptionOptions) *runtime.Pager[StandbyVirtualMachinePoolsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[StandbyVirtualMachinePoolsClientListBySubscriptionResponse]{
		More: func(page StandbyVirtualMachinePoolsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *StandbyVirtualMachinePoolsClientListBySubscriptionResponse) (StandbyVirtualMachinePoolsClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "StandbyVirtualMachinePoolsClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return StandbyVirtualMachinePoolsClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *StandbyVirtualMachinePoolsClient) listBySubscriptionCreateRequest(ctx context.Context, _ *StandbyVirtualMachinePoolsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.StandbyPool/standbyVirtualMachinePools"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *StandbyVirtualMachinePoolsClient) listBySubscriptionHandleResponse(resp *http.Response) (StandbyVirtualMachinePoolsClientListBySubscriptionResponse, error) {
	result := StandbyVirtualMachinePoolsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StandbyVirtualMachinePoolResourceListResult); err != nil {
		return StandbyVirtualMachinePoolsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Update a StandbyVirtualMachinePoolResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - standbyVirtualMachinePoolName - Name of the standby virtual machine pool
//   - properties - The resource properties to be updated.
//   - options - StandbyVirtualMachinePoolsClientUpdateOptions contains the optional parameters for the StandbyVirtualMachinePoolsClient.Update
//     method.
func (client *StandbyVirtualMachinePoolsClient) Update(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, properties StandbyVirtualMachinePoolResourceUpdate, options *StandbyVirtualMachinePoolsClientUpdateOptions) (StandbyVirtualMachinePoolsClientUpdateResponse, error) {
	var err error
	const operationName = "StandbyVirtualMachinePoolsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, standbyVirtualMachinePoolName, properties, options)
	if err != nil {
		return StandbyVirtualMachinePoolsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return StandbyVirtualMachinePoolsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return StandbyVirtualMachinePoolsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *StandbyVirtualMachinePoolsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, standbyVirtualMachinePoolName string, properties StandbyVirtualMachinePoolResourceUpdate, _ *StandbyVirtualMachinePoolsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StandbyPool/standbyVirtualMachinePools/{standbyVirtualMachinePoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if standbyVirtualMachinePoolName == "" {
		return nil, errors.New("parameter standbyVirtualMachinePoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{standbyVirtualMachinePoolName}", url.PathEscape(standbyVirtualMachinePoolName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, properties); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *StandbyVirtualMachinePoolsClient) updateHandleResponse(resp *http.Response) (StandbyVirtualMachinePoolsClientUpdateResponse, error) {
	result := StandbyVirtualMachinePoolsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StandbyVirtualMachinePoolResource); err != nil {
		return StandbyVirtualMachinePoolsClientUpdateResponse{}, err
	}
	return result, nil
}
