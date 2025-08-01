// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcontainerservice

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

// ManagedNamespacesClient contains the methods for the ManagedNamespaces group.
// Don't use this type directly, use NewManagedNamespacesClient() instead.
type ManagedNamespacesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewManagedNamespacesClient creates a new instance of ManagedNamespacesClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewManagedNamespacesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ManagedNamespacesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ManagedNamespacesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a managed namespace in the specified managed cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-05-02-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the managed cluster resource.
//   - managedNamespaceName - The name of the managed namespace.
//   - parameters - The namespace to create or update.
//   - options - ManagedNamespacesClientBeginCreateOrUpdateOptions contains the optional parameters for the ManagedNamespacesClient.BeginCreateOrUpdate
//     method.
func (client *ManagedNamespacesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, parameters ManagedNamespace, options *ManagedNamespacesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedNamespacesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, resourceName, managedNamespaceName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ManagedNamespacesClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ManagedNamespacesClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates or updates a managed namespace in the specified managed cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-05-02-preview
func (client *ManagedNamespacesClient) createOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, parameters ManagedNamespace, options *ManagedNamespacesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ManagedNamespacesClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, managedNamespaceName, parameters, options)
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
func (client *ManagedNamespacesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, parameters ManagedNamespace, _ *ManagedNamespacesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/managedNamespaces/{managedNamespaceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if managedNamespaceName == "" {
		return nil, errors.New("parameter managedNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedNamespaceName}", url.PathEscape(managedNamespaceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-05-02-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes a namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-05-02-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the managed cluster resource.
//   - managedNamespaceName - The name of the managed namespace.
//   - options - ManagedNamespacesClientBeginDeleteOptions contains the optional parameters for the ManagedNamespacesClient.BeginDelete
//     method.
func (client *ManagedNamespacesClient) BeginDelete(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, options *ManagedNamespacesClientBeginDeleteOptions) (*runtime.Poller[ManagedNamespacesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, resourceName, managedNamespaceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ManagedNamespacesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ManagedNamespacesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes a namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-05-02-preview
func (client *ManagedNamespacesClient) deleteOperation(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, options *ManagedNamespacesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ManagedNamespacesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, managedNamespaceName, options)
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
func (client *ManagedNamespacesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, _ *ManagedNamespacesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/managedNamespaces/{managedNamespaceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if managedNamespaceName == "" {
		return nil, errors.New("parameter managedNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedNamespaceName}", url.PathEscape(managedNamespaceName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-05-02-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified namespace of a managed cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-05-02-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the managed cluster resource.
//   - managedNamespaceName - The name of the managed namespace.
//   - options - ManagedNamespacesClientGetOptions contains the optional parameters for the ManagedNamespacesClient.Get method.
func (client *ManagedNamespacesClient) Get(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, options *ManagedNamespacesClientGetOptions) (ManagedNamespacesClientGetResponse, error) {
	var err error
	const operationName = "ManagedNamespacesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, managedNamespaceName, options)
	if err != nil {
		return ManagedNamespacesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedNamespacesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ManagedNamespacesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ManagedNamespacesClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, _ *ManagedNamespacesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/managedNamespaces/{managedNamespaceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if managedNamespaceName == "" {
		return nil, errors.New("parameter managedNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedNamespaceName}", url.PathEscape(managedNamespaceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-05-02-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ManagedNamespacesClient) getHandleResponse(resp *http.Response) (ManagedNamespacesClientGetResponse, error) {
	result := ManagedNamespacesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedNamespace); err != nil {
		return ManagedNamespacesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByManagedClusterPager - Gets a list of managed namespaces in the specified managed cluster.
//
// Generated from API version 2025-05-02-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the managed cluster resource.
//   - options - ManagedNamespacesClientListByManagedClusterOptions contains the optional parameters for the ManagedNamespacesClient.NewListByManagedClusterPager
//     method.
func (client *ManagedNamespacesClient) NewListByManagedClusterPager(resourceGroupName string, resourceName string, options *ManagedNamespacesClientListByManagedClusterOptions) *runtime.Pager[ManagedNamespacesClientListByManagedClusterResponse] {
	return runtime.NewPager(runtime.PagingHandler[ManagedNamespacesClientListByManagedClusterResponse]{
		More: func(page ManagedNamespacesClientListByManagedClusterResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ManagedNamespacesClientListByManagedClusterResponse) (ManagedNamespacesClientListByManagedClusterResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ManagedNamespacesClient.NewListByManagedClusterPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByManagedClusterCreateRequest(ctx, resourceGroupName, resourceName, options)
			}, nil)
			if err != nil {
				return ManagedNamespacesClientListByManagedClusterResponse{}, err
			}
			return client.listByManagedClusterHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByManagedClusterCreateRequest creates the ListByManagedCluster request.
func (client *ManagedNamespacesClient) listByManagedClusterCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, _ *ManagedNamespacesClientListByManagedClusterOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/managedNamespaces"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-05-02-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByManagedClusterHandleResponse handles the ListByManagedCluster response.
func (client *ManagedNamespacesClient) listByManagedClusterHandleResponse(resp *http.Response) (ManagedNamespacesClientListByManagedClusterResponse, error) {
	result := ManagedNamespacesClientListByManagedClusterResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedNamespaceListResult); err != nil {
		return ManagedNamespacesClientListByManagedClusterResponse{}, err
	}
	return result, nil
}

// ListCredential - Lists the credentials of a namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-05-02-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the managed cluster resource.
//   - managedNamespaceName - The name of the managed namespace.
//   - options - ManagedNamespacesClientListCredentialOptions contains the optional parameters for the ManagedNamespacesClient.ListCredential
//     method.
func (client *ManagedNamespacesClient) ListCredential(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, options *ManagedNamespacesClientListCredentialOptions) (ManagedNamespacesClientListCredentialResponse, error) {
	var err error
	const operationName = "ManagedNamespacesClient.ListCredential"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listCredentialCreateRequest(ctx, resourceGroupName, resourceName, managedNamespaceName, options)
	if err != nil {
		return ManagedNamespacesClientListCredentialResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedNamespacesClientListCredentialResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ManagedNamespacesClientListCredentialResponse{}, err
	}
	resp, err := client.listCredentialHandleResponse(httpResp)
	return resp, err
}

// listCredentialCreateRequest creates the ListCredential request.
func (client *ManagedNamespacesClient) listCredentialCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, _ *ManagedNamespacesClientListCredentialOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/managedNamespaces/{managedNamespaceName}/listCredential"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if managedNamespaceName == "" {
		return nil, errors.New("parameter managedNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedNamespaceName}", url.PathEscape(managedNamespaceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-05-02-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listCredentialHandleResponse handles the ListCredential response.
func (client *ManagedNamespacesClient) listCredentialHandleResponse(resp *http.Response) (ManagedNamespacesClientListCredentialResponse, error) {
	result := ManagedNamespacesClientListCredentialResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CredentialResults); err != nil {
		return ManagedNamespacesClientListCredentialResponse{}, err
	}
	return result, nil
}

// Update - Updates tags on a managed namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-05-02-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the managed cluster resource.
//   - managedNamespaceName - The name of the managed namespace.
//   - parameters - Parameters supplied to the patch namespace operation, we only support patch tags for now.
//   - options - ManagedNamespacesClientUpdateOptions contains the optional parameters for the ManagedNamespacesClient.Update
//     method.
func (client *ManagedNamespacesClient) Update(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, parameters TagsObject, options *ManagedNamespacesClientUpdateOptions) (ManagedNamespacesClientUpdateResponse, error) {
	var err error
	const operationName = "ManagedNamespacesClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, managedNamespaceName, parameters, options)
	if err != nil {
		return ManagedNamespacesClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedNamespacesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ManagedNamespacesClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *ManagedNamespacesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, managedNamespaceName string, parameters TagsObject, _ *ManagedNamespacesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/managedNamespaces/{managedNamespaceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if managedNamespaceName == "" {
		return nil, errors.New("parameter managedNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedNamespaceName}", url.PathEscape(managedNamespaceName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-05-02-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *ManagedNamespacesClient) updateHandleResponse(resp *http.Response) (ManagedNamespacesClientUpdateResponse, error) {
	result := ManagedNamespacesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedNamespace); err != nil {
		return ManagedNamespacesClientUpdateResponse{}, err
	}
	return result, nil
}
