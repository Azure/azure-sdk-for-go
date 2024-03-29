//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagedapplications

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

// ApplicationsClient contains the methods for the Applications group.
// Don't use this type directly, use NewApplicationsClient() instead.
type ApplicationsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewApplicationsClient creates a new instance of ApplicationsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewApplicationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ApplicationsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ApplicationsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates a new managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - applicationName - The name of the managed application.
//   - parameters - Parameters supplied to the create or update a managed application.
//   - options - ApplicationsClientBeginCreateOrUpdateOptions contains the optional parameters for the ApplicationsClient.BeginCreateOrUpdate
//     method.
func (client *ApplicationsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, applicationName string, parameters Application, options *ApplicationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ApplicationsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, applicationName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ApplicationsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ApplicationsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates a new managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
func (client *ApplicationsClient) createOrUpdate(ctx context.Context, resourceGroupName string, applicationName string, parameters Application, options *ApplicationsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ApplicationsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, applicationName, parameters, options)
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
func (client *ApplicationsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, parameters Application, options *ApplicationsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applications/{applicationName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginCreateOrUpdateByID - Creates a new managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - applicationID - The fully qualified ID of the managed application, including the managed application name and the managed
//     application resource type. Use the format,
//     /subscriptions/{guid}/resourceGroups/{resource-group-name}/Microsoft.Solutions/applications/{application-name}
//   - parameters - Parameters supplied to the create or update a managed application.
//   - options - ApplicationsClientBeginCreateOrUpdateByIDOptions contains the optional parameters for the ApplicationsClient.BeginCreateOrUpdateByID
//     method.
func (client *ApplicationsClient) BeginCreateOrUpdateByID(ctx context.Context, applicationID string, parameters Application, options *ApplicationsClientBeginCreateOrUpdateByIDOptions) (*runtime.Poller[ApplicationsClientCreateOrUpdateByIDResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdateByID(ctx, applicationID, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ApplicationsClientCreateOrUpdateByIDResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ApplicationsClientCreateOrUpdateByIDResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdateByID - Creates a new managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
func (client *ApplicationsClient) createOrUpdateByID(ctx context.Context, applicationID string, parameters Application, options *ApplicationsClientBeginCreateOrUpdateByIDOptions) (*http.Response, error) {
	var err error
	const operationName = "ApplicationsClient.BeginCreateOrUpdateByID"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateByIDCreateRequest(ctx, applicationID, parameters, options)
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

// createOrUpdateByIDCreateRequest creates the CreateOrUpdateByID request.
func (client *ApplicationsClient) createOrUpdateByIDCreateRequest(ctx context.Context, applicationID string, parameters Application, options *ApplicationsClientBeginCreateOrUpdateByIDOptions) (*policy.Request, error) {
	urlPath := "/{applicationId}"
	urlPath = strings.ReplaceAll(urlPath, "{applicationId}", applicationID)
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes the managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - applicationName - The name of the managed application.
//   - options - ApplicationsClientBeginDeleteOptions contains the optional parameters for the ApplicationsClient.BeginDelete
//     method.
func (client *ApplicationsClient) BeginDelete(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientBeginDeleteOptions) (*runtime.Poller[ApplicationsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, applicationName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ApplicationsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ApplicationsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes the managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
func (client *ApplicationsClient) deleteOperation(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ApplicationsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, applicationName, options)
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
func (client *ApplicationsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applications/{applicationName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginDeleteByID - Deletes the managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - applicationID - The fully qualified ID of the managed application, including the managed application name and the managed
//     application resource type. Use the format,
//     /subscriptions/{guid}/resourceGroups/{resource-group-name}/Microsoft.Solutions/applications/{application-name}
//   - options - ApplicationsClientBeginDeleteByIDOptions contains the optional parameters for the ApplicationsClient.BeginDeleteByID
//     method.
func (client *ApplicationsClient) BeginDeleteByID(ctx context.Context, applicationID string, options *ApplicationsClientBeginDeleteByIDOptions) (*runtime.Poller[ApplicationsClientDeleteByIDResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteByID(ctx, applicationID, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ApplicationsClientDeleteByIDResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ApplicationsClientDeleteByIDResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// DeleteByID - Deletes the managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
func (client *ApplicationsClient) deleteByID(ctx context.Context, applicationID string, options *ApplicationsClientBeginDeleteByIDOptions) (*http.Response, error) {
	var err error
	const operationName = "ApplicationsClient.BeginDeleteByID"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteByIDCreateRequest(ctx, applicationID, options)
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

// deleteByIDCreateRequest creates the DeleteByID request.
func (client *ApplicationsClient) deleteByIDCreateRequest(ctx context.Context, applicationID string, options *ApplicationsClientBeginDeleteByIDOptions) (*policy.Request, error) {
	urlPath := "/{applicationId}"
	urlPath = strings.ReplaceAll(urlPath, "{applicationId}", applicationID)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - applicationName - The name of the managed application.
//   - options - ApplicationsClientGetOptions contains the optional parameters for the ApplicationsClient.Get method.
func (client *ApplicationsClient) Get(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientGetOptions) (ApplicationsClientGetResponse, error) {
	var err error
	const operationName = "ApplicationsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, applicationName, options)
	if err != nil {
		return ApplicationsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNotFound) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ApplicationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applications/{applicationName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ApplicationsClient) getHandleResponse(resp *http.Response) (ApplicationsClientGetResponse, error) {
	result := ApplicationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Application); err != nil {
		return ApplicationsClientGetResponse{}, err
	}
	return result, nil
}

// GetByID - Gets the managed application.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - applicationID - The fully qualified ID of the managed application, including the managed application name and the managed
//     application resource type. Use the format,
//     /subscriptions/{guid}/resourceGroups/{resource-group-name}/Microsoft.Solutions/applications/{application-name}
//   - options - ApplicationsClientGetByIDOptions contains the optional parameters for the ApplicationsClient.GetByID method.
func (client *ApplicationsClient) GetByID(ctx context.Context, applicationID string, options *ApplicationsClientGetByIDOptions) (ApplicationsClientGetByIDResponse, error) {
	var err error
	const operationName = "ApplicationsClient.GetByID"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getByIDCreateRequest(ctx, applicationID, options)
	if err != nil {
		return ApplicationsClientGetByIDResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationsClientGetByIDResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNotFound) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationsClientGetByIDResponse{}, err
	}
	resp, err := client.getByIDHandleResponse(httpResp)
	return resp, err
}

// getByIDCreateRequest creates the GetByID request.
func (client *ApplicationsClient) getByIDCreateRequest(ctx context.Context, applicationID string, options *ApplicationsClientGetByIDOptions) (*policy.Request, error) {
	urlPath := "/{applicationId}"
	urlPath = strings.ReplaceAll(urlPath, "{applicationId}", applicationID)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getByIDHandleResponse handles the GetByID response.
func (client *ApplicationsClient) getByIDHandleResponse(resp *http.Response) (ApplicationsClientGetByIDResponse, error) {
	result := ApplicationsClientGetByIDResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Application); err != nil {
		return ApplicationsClientGetByIDResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Gets all the applications within a resource group.
//
// Generated from API version 2018-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - ApplicationsClientListByResourceGroupOptions contains the optional parameters for the ApplicationsClient.NewListByResourceGroupPager
//     method.
func (client *ApplicationsClient) NewListByResourceGroupPager(resourceGroupName string, options *ApplicationsClientListByResourceGroupOptions) *runtime.Pager[ApplicationsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[ApplicationsClientListByResourceGroupResponse]{
		More: func(page ApplicationsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ApplicationsClientListByResourceGroupResponse) (ApplicationsClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ApplicationsClient.NewListByResourceGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			}, nil)
			if err != nil {
				return ApplicationsClientListByResourceGroupResponse{}, err
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ApplicationsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ApplicationsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applications"
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
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ApplicationsClient) listByResourceGroupHandleResponse(resp *http.Response) (ApplicationsClientListByResourceGroupResponse, error) {
	result := ApplicationsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationListResult); err != nil {
		return ApplicationsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Gets all the applications within a subscription.
//
// Generated from API version 2018-06-01
//   - options - ApplicationsClientListBySubscriptionOptions contains the optional parameters for the ApplicationsClient.NewListBySubscriptionPager
//     method.
func (client *ApplicationsClient) NewListBySubscriptionPager(options *ApplicationsClientListBySubscriptionOptions) *runtime.Pager[ApplicationsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[ApplicationsClientListBySubscriptionResponse]{
		More: func(page ApplicationsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ApplicationsClientListBySubscriptionResponse) (ApplicationsClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ApplicationsClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return ApplicationsClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *ApplicationsClient) listBySubscriptionCreateRequest(ctx context.Context, options *ApplicationsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Solutions/applications"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *ApplicationsClient) listBySubscriptionHandleResponse(resp *http.Response) (ApplicationsClientListBySubscriptionResponse, error) {
	result := ApplicationsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationListResult); err != nil {
		return ApplicationsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Updates an existing managed application. The only value that can be updated via PATCH currently is the tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - applicationName - The name of the managed application.
//   - options - ApplicationsClientUpdateOptions contains the optional parameters for the ApplicationsClient.Update method.
func (client *ApplicationsClient) Update(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientUpdateOptions) (ApplicationsClientUpdateResponse, error) {
	var err error
	const operationName = "ApplicationsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, applicationName, options)
	if err != nil {
		return ApplicationsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *ApplicationsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, applicationName string, options *ApplicationsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Solutions/applications/{applicationName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
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
func (client *ApplicationsClient) updateHandleResponse(resp *http.Response) (ApplicationsClientUpdateResponse, error) {
	result := ApplicationsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Application); err != nil {
		return ApplicationsClientUpdateResponse{}, err
	}
	return result, nil
}

// UpdateByID - Updates an existing managed application. The only value that can be updated via PATCH currently is the tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-06-01
//   - applicationID - The fully qualified ID of the managed application, including the managed application name and the managed
//     application resource type. Use the format,
//     /subscriptions/{guid}/resourceGroups/{resource-group-name}/Microsoft.Solutions/applications/{application-name}
//   - options - ApplicationsClientUpdateByIDOptions contains the optional parameters for the ApplicationsClient.UpdateByID method.
func (client *ApplicationsClient) UpdateByID(ctx context.Context, applicationID string, options *ApplicationsClientUpdateByIDOptions) (ApplicationsClientUpdateByIDResponse, error) {
	var err error
	const operationName = "ApplicationsClient.UpdateByID"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateByIDCreateRequest(ctx, applicationID, options)
	if err != nil {
		return ApplicationsClientUpdateByIDResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationsClientUpdateByIDResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationsClientUpdateByIDResponse{}, err
	}
	resp, err := client.updateByIDHandleResponse(httpResp)
	return resp, err
}

// updateByIDCreateRequest creates the UpdateByID request.
func (client *ApplicationsClient) updateByIDCreateRequest(ctx context.Context, applicationID string, options *ApplicationsClientUpdateByIDOptions) (*policy.Request, error) {
	urlPath := "/{applicationId}"
	urlPath = strings.ReplaceAll(urlPath, "{applicationId}", applicationID)
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
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

// updateByIDHandleResponse handles the UpdateByID response.
func (client *ApplicationsClient) updateByIDHandleResponse(resp *http.Response) (ApplicationsClientUpdateByIDResponse, error) {
	result := ApplicationsClientUpdateByIDResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Application); err != nil {
		return ApplicationsClientUpdateByIDResponse{}, err
	}
	return result, nil
}
