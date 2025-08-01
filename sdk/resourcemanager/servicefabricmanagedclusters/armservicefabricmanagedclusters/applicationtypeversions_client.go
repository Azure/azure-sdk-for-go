// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armservicefabricmanagedclusters

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

// ApplicationTypeVersionsClient contains the methods for the ApplicationTypeVersions group.
// Don't use this type directly, use NewApplicationTypeVersionsClient() instead.
type ApplicationTypeVersionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewApplicationTypeVersionsClient creates a new instance of ApplicationTypeVersionsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewApplicationTypeVersionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ApplicationTypeVersionsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ApplicationTypeVersionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update a Service Fabric managed application type version resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-03-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster resource.
//   - applicationTypeName - The name of the application type name resource.
//   - version - The application type version.
//   - parameters - The application type version resource.
//   - options - ApplicationTypeVersionsClientBeginCreateOrUpdateOptions contains the optional parameters for the ApplicationTypeVersionsClient.BeginCreateOrUpdate
//     method.
func (client *ApplicationTypeVersionsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, parameters ApplicationTypeVersionResource, options *ApplicationTypeVersionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ApplicationTypeVersionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, applicationTypeName, version, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ApplicationTypeVersionsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ApplicationTypeVersionsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create or update a Service Fabric managed application type version resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-03-01-preview
func (client *ApplicationTypeVersionsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, parameters ApplicationTypeVersionResource, options *ApplicationTypeVersionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ApplicationTypeVersionsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, applicationTypeName, version, parameters, options)
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

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ApplicationTypeVersionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, parameters ApplicationTypeVersionResource, _ *ApplicationTypeVersionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{clusterName}/applicationTypes/{applicationTypeName}/versions/{version}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationTypeName == "" {
		return nil, errors.New("parameter applicationTypeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationTypeName}", url.PathEscape(applicationTypeName))
	if version == "" {
		return nil, errors.New("parameter version cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{version}", url.PathEscape(version))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete a Service Fabric managed application type version resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-03-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster resource.
//   - applicationTypeName - The name of the application type name resource.
//   - version - The application type version.
//   - options - ApplicationTypeVersionsClientBeginDeleteOptions contains the optional parameters for the ApplicationTypeVersionsClient.BeginDelete
//     method.
func (client *ApplicationTypeVersionsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, options *ApplicationTypeVersionsClientBeginDeleteOptions) (*runtime.Poller[ApplicationTypeVersionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, applicationTypeName, version, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ApplicationTypeVersionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ApplicationTypeVersionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a Service Fabric managed application type version resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-03-01-preview
func (client *ApplicationTypeVersionsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, options *ApplicationTypeVersionsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ApplicationTypeVersionsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, applicationTypeName, version, options)
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
func (client *ApplicationTypeVersionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, _ *ApplicationTypeVersionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{clusterName}/applicationTypes/{applicationTypeName}/versions/{version}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationTypeName == "" {
		return nil, errors.New("parameter applicationTypeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationTypeName}", url.PathEscape(applicationTypeName))
	if version == "" {
		return nil, errors.New("parameter version cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{version}", url.PathEscape(version))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a Service Fabric managed application type version resource created or in the process of being created in the
// Service Fabric managed application type name resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-03-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster resource.
//   - applicationTypeName - The name of the application type name resource.
//   - version - The application type version.
//   - options - ApplicationTypeVersionsClientGetOptions contains the optional parameters for the ApplicationTypeVersionsClient.Get
//     method.
func (client *ApplicationTypeVersionsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, options *ApplicationTypeVersionsClientGetOptions) (ApplicationTypeVersionsClientGetResponse, error) {
	var err error
	const operationName = "ApplicationTypeVersionsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, applicationTypeName, version, options)
	if err != nil {
		return ApplicationTypeVersionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationTypeVersionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationTypeVersionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ApplicationTypeVersionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, _ *ApplicationTypeVersionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{clusterName}/applicationTypes/{applicationTypeName}/versions/{version}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationTypeName == "" {
		return nil, errors.New("parameter applicationTypeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationTypeName}", url.PathEscape(applicationTypeName))
	if version == "" {
		return nil, errors.New("parameter version cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{version}", url.PathEscape(version))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ApplicationTypeVersionsClient) getHandleResponse(resp *http.Response) (ApplicationTypeVersionsClientGetResponse, error) {
	result := ApplicationTypeVersionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationTypeVersionResource); err != nil {
		return ApplicationTypeVersionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByApplicationTypesPager - Gets all application type version resources created or in the process of being created
// in the Service Fabric managed application type name resource.
//
// Generated from API version 2025-03-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster resource.
//   - applicationTypeName - The name of the application type name resource.
//   - options - ApplicationTypeVersionsClientListByApplicationTypesOptions contains the optional parameters for the ApplicationTypeVersionsClient.NewListByApplicationTypesPager
//     method.
func (client *ApplicationTypeVersionsClient) NewListByApplicationTypesPager(resourceGroupName string, clusterName string, applicationTypeName string, options *ApplicationTypeVersionsClientListByApplicationTypesOptions) *runtime.Pager[ApplicationTypeVersionsClientListByApplicationTypesResponse] {
	return runtime.NewPager(runtime.PagingHandler[ApplicationTypeVersionsClientListByApplicationTypesResponse]{
		More: func(page ApplicationTypeVersionsClientListByApplicationTypesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ApplicationTypeVersionsClientListByApplicationTypesResponse) (ApplicationTypeVersionsClientListByApplicationTypesResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ApplicationTypeVersionsClient.NewListByApplicationTypesPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByApplicationTypesCreateRequest(ctx, resourceGroupName, clusterName, applicationTypeName, options)
			}, nil)
			if err != nil {
				return ApplicationTypeVersionsClientListByApplicationTypesResponse{}, err
			}
			return client.listByApplicationTypesHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByApplicationTypesCreateRequest creates the ListByApplicationTypes request.
func (client *ApplicationTypeVersionsClient) listByApplicationTypesCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, _ *ApplicationTypeVersionsClientListByApplicationTypesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{clusterName}/applicationTypes/{applicationTypeName}/versions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationTypeName == "" {
		return nil, errors.New("parameter applicationTypeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationTypeName}", url.PathEscape(applicationTypeName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByApplicationTypesHandleResponse handles the ListByApplicationTypes response.
func (client *ApplicationTypeVersionsClient) listByApplicationTypesHandleResponse(resp *http.Response) (ApplicationTypeVersionsClientListByApplicationTypesResponse, error) {
	result := ApplicationTypeVersionsClientListByApplicationTypesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationTypeVersionResourceList); err != nil {
		return ApplicationTypeVersionsClientListByApplicationTypesResponse{}, err
	}
	return result, nil
}

// Update - Updates the tags of an application type version resource of a given managed cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-03-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster resource.
//   - applicationTypeName - The name of the application type name resource.
//   - version - The application type version.
//   - parameters - The application type version resource updated tags.
//   - options - ApplicationTypeVersionsClientUpdateOptions contains the optional parameters for the ApplicationTypeVersionsClient.Update
//     method.
func (client *ApplicationTypeVersionsClient) Update(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, parameters ApplicationTypeVersionUpdateParameters, options *ApplicationTypeVersionsClientUpdateOptions) (ApplicationTypeVersionsClientUpdateResponse, error) {
	var err error
	const operationName = "ApplicationTypeVersionsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, clusterName, applicationTypeName, version, parameters, options)
	if err != nil {
		return ApplicationTypeVersionsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ApplicationTypeVersionsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ApplicationTypeVersionsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *ApplicationTypeVersionsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationTypeName string, version string, parameters ApplicationTypeVersionUpdateParameters, _ *ApplicationTypeVersionsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{clusterName}/applicationTypes/{applicationTypeName}/versions/{version}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationTypeName == "" {
		return nil, errors.New("parameter applicationTypeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationTypeName}", url.PathEscape(applicationTypeName))
	if version == "" {
		return nil, errors.New("parameter version cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{version}", url.PathEscape(version))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *ApplicationTypeVersionsClient) updateHandleResponse(resp *http.Response) (ApplicationTypeVersionsClientUpdateResponse, error) {
	result := ApplicationTypeVersionsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationTypeVersionResource); err != nil {
		return ApplicationTypeVersionsClientUpdateResponse{}, err
	}
	return result, nil
}
