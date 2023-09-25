//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armeventhub

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

// ApplicationGroupClient contains the methods for the ApplicationGroup group.
// Don't use this type directly, use NewApplicationGroupClient() instead.
type ApplicationGroupClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewApplicationGroupClient creates a new instance of ApplicationGroupClient with the specified values.
//   - subscriptionID - Subscription credentials that uniquely identify a Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewApplicationGroupClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ApplicationGroupClient, error) {
	cl, err := arm.NewClient(moduleName+".ApplicationGroupClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ApplicationGroupClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdateApplicationGroup - Creates or updates an ApplicationGroup for a Namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - applicationGroupName - The Application Group name
//   - parameters - The ApplicationGroup.
//   - options - ApplicationGroupClientCreateOrUpdateApplicationGroupOptions contains the optional parameters for the ApplicationGroupClient.CreateOrUpdateApplicationGroup
//     method.
func (client *ApplicationGroupClient) CreateOrUpdateApplicationGroup(ctx context.Context, resourceGroupName string, namespaceName string, applicationGroupName string, parameters ApplicationGroup, options *ApplicationGroupClientCreateOrUpdateApplicationGroupOptions) (ApplicationGroupClientCreateOrUpdateApplicationGroupResponse, error) {
	var err error
	req, err := client.createOrUpdateApplicationGroupCreateRequest(ctx, resourceGroupName, namespaceName, applicationGroupName, parameters, options)
	if err != nil {
		return ApplicationGroupClientCreateOrUpdateApplicationGroupResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationGroupClientCreateOrUpdateApplicationGroupResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationGroupClientCreateOrUpdateApplicationGroupResponse{}, err
	}
	resp, err := client.createOrUpdateApplicationGroupHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateApplicationGroupCreateRequest creates the CreateOrUpdateApplicationGroup request.
func (client *ApplicationGroupClient) createOrUpdateApplicationGroupCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, applicationGroupName string, parameters ApplicationGroup, options *ApplicationGroupClientCreateOrUpdateApplicationGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/applicationGroups/{applicationGroupName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if applicationGroupName == "" {
		return nil, errors.New("parameter applicationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationGroupName}", url.PathEscape(applicationGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateApplicationGroupHandleResponse handles the CreateOrUpdateApplicationGroup response.
func (client *ApplicationGroupClient) createOrUpdateApplicationGroupHandleResponse(resp *http.Response) (ApplicationGroupClientCreateOrUpdateApplicationGroupResponse, error) {
	result := ApplicationGroupClientCreateOrUpdateApplicationGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationGroup); err != nil {
		return ApplicationGroupClientCreateOrUpdateApplicationGroupResponse{}, err
	}
	return result, nil
}

// Delete - Deletes an ApplicationGroup for a Namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - applicationGroupName - The Application Group name
//   - options - ApplicationGroupClientDeleteOptions contains the optional parameters for the ApplicationGroupClient.Delete method.
func (client *ApplicationGroupClient) Delete(ctx context.Context, resourceGroupName string, namespaceName string, applicationGroupName string, options *ApplicationGroupClientDeleteOptions) (ApplicationGroupClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, namespaceName, applicationGroupName, options)
	if err != nil {
		return ApplicationGroupClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationGroupClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationGroupClientDeleteResponse{}, err
	}
	return ApplicationGroupClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ApplicationGroupClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, applicationGroupName string, options *ApplicationGroupClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/applicationGroups/{applicationGroupName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if applicationGroupName == "" {
		return nil, errors.New("parameter applicationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationGroupName}", url.PathEscape(applicationGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets an ApplicationGroup for a Namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - applicationGroupName - The Application Group name
//   - options - ApplicationGroupClientGetOptions contains the optional parameters for the ApplicationGroupClient.Get method.
func (client *ApplicationGroupClient) Get(ctx context.Context, resourceGroupName string, namespaceName string, applicationGroupName string, options *ApplicationGroupClientGetOptions) (ApplicationGroupClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, namespaceName, applicationGroupName, options)
	if err != nil {
		return ApplicationGroupClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationGroupClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationGroupClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ApplicationGroupClient) getCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, applicationGroupName string, options *ApplicationGroupClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/applicationGroups/{applicationGroupName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if applicationGroupName == "" {
		return nil, errors.New("parameter applicationGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationGroupName}", url.PathEscape(applicationGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ApplicationGroupClient) getHandleResponse(resp *http.Response) (ApplicationGroupClientGetResponse, error) {
	result := ApplicationGroupClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationGroup); err != nil {
		return ApplicationGroupClientGetResponse{}, err
	}
	return result, nil
}

// NewListByNamespacePager - Gets a list of application groups for a Namespace.
//
// Generated from API version 2022-10-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - options - ApplicationGroupClientListByNamespaceOptions contains the optional parameters for the ApplicationGroupClient.NewListByNamespacePager
//     method.
func (client *ApplicationGroupClient) NewListByNamespacePager(resourceGroupName string, namespaceName string, options *ApplicationGroupClientListByNamespaceOptions) (*runtime.Pager[ApplicationGroupClientListByNamespaceResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ApplicationGroupClientListByNamespaceResponse]{
		More: func(page ApplicationGroupClientListByNamespaceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ApplicationGroupClientListByNamespaceResponse) (ApplicationGroupClientListByNamespaceResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByNamespaceCreateRequest(ctx, resourceGroupName, namespaceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ApplicationGroupClientListByNamespaceResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ApplicationGroupClientListByNamespaceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ApplicationGroupClientListByNamespaceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByNamespaceHandleResponse(resp)
		},
	})
}

// listByNamespaceCreateRequest creates the ListByNamespace request.
func (client *ApplicationGroupClient) listByNamespaceCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, options *ApplicationGroupClientListByNamespaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/applicationGroups"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByNamespaceHandleResponse handles the ListByNamespace response.
func (client *ApplicationGroupClient) listByNamespaceHandleResponse(resp *http.Response) (ApplicationGroupClientListByNamespaceResponse, error) {
	result := ApplicationGroupClientListByNamespaceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationGroupListResult); err != nil {
		return ApplicationGroupClientListByNamespaceResponse{}, err
	}
	return result, nil
}

