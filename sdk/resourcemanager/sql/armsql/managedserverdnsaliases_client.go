//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsql

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

// ManagedServerDNSAliasesClient contains the methods for the ManagedServerDNSAliases group.
// Don't use this type directly, use NewManagedServerDNSAliasesClient() instead.
type ManagedServerDNSAliasesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewManagedServerDNSAliasesClient creates a new instance of ManagedServerDNSAliasesClient with the specified values.
//   - subscriptionID - The subscription ID that identifies an Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewManagedServerDNSAliasesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ManagedServerDNSAliasesClient, error) {
	cl, err := arm.NewClient(moduleName+".ManagedServerDNSAliasesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ManagedServerDNSAliasesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginAcquire - Acquires managed server DNS alias from another managed server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - managedInstanceName - The name of the managed instance.
//   - options - ManagedServerDNSAliasesClientBeginAcquireOptions contains the optional parameters for the ManagedServerDNSAliasesClient.BeginAcquire
//     method.
func (client *ManagedServerDNSAliasesClient) BeginAcquire(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, parameters ManagedServerDNSAliasAcquisition, options *ManagedServerDNSAliasesClientBeginAcquireOptions) (*runtime.Poller[ManagedServerDNSAliasesClientAcquireResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.acquire(ctx, resourceGroupName, managedInstanceName, dnsAliasName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[ManagedServerDNSAliasesClientAcquireResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ManagedServerDNSAliasesClientAcquireResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Acquire - Acquires managed server DNS alias from another managed server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
func (client *ManagedServerDNSAliasesClient) acquire(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, parameters ManagedServerDNSAliasAcquisition, options *ManagedServerDNSAliasesClientBeginAcquireOptions) (*http.Response, error) {
	var err error
	req, err := client.acquireCreateRequest(ctx, resourceGroupName, managedInstanceName, dnsAliasName, parameters, options)
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

// acquireCreateRequest creates the Acquire request.
func (client *ManagedServerDNSAliasesClient) acquireCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, parameters ManagedServerDNSAliasAcquisition, options *ManagedServerDNSAliasesClientBeginAcquireOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/dnsAliases/{dnsAliasName}/acquire"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if dnsAliasName == "" {
		return nil, errors.New("parameter dnsAliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dnsAliasName}", url.PathEscape(dnsAliasName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginCreateOrUpdate - Creates a managed server DNS alias.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - managedInstanceName - The name of the managed instance.
//   - options - ManagedServerDNSAliasesClientBeginCreateOrUpdateOptions contains the optional parameters for the ManagedServerDNSAliasesClient.BeginCreateOrUpdate
//     method.
func (client *ManagedServerDNSAliasesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, parameters ManagedServerDNSAliasCreation, options *ManagedServerDNSAliasesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ManagedServerDNSAliasesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, managedInstanceName, dnsAliasName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[ManagedServerDNSAliasesClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ManagedServerDNSAliasesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates a managed server DNS alias.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
func (client *ManagedServerDNSAliasesClient) createOrUpdate(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, parameters ManagedServerDNSAliasCreation, options *ManagedServerDNSAliasesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, managedInstanceName, dnsAliasName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ManagedServerDNSAliasesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, parameters ManagedServerDNSAliasCreation, options *ManagedServerDNSAliasesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/dnsAliases/{dnsAliasName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if dnsAliasName == "" {
		return nil, errors.New("parameter dnsAliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dnsAliasName}", url.PathEscape(dnsAliasName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes the managed server DNS alias with the given name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - managedInstanceName - The name of the managed instance.
//   - options - ManagedServerDNSAliasesClientBeginDeleteOptions contains the optional parameters for the ManagedServerDNSAliasesClient.BeginDelete
//     method.
func (client *ManagedServerDNSAliasesClient) BeginDelete(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, options *ManagedServerDNSAliasesClientBeginDeleteOptions) (*runtime.Poller[ManagedServerDNSAliasesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, managedInstanceName, dnsAliasName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[ManagedServerDNSAliasesClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ManagedServerDNSAliasesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the managed server DNS alias with the given name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
func (client *ManagedServerDNSAliasesClient) deleteOperation(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, options *ManagedServerDNSAliasesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, managedInstanceName, dnsAliasName, options)
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
func (client *ManagedServerDNSAliasesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, options *ManagedServerDNSAliasesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/dnsAliases/{dnsAliasName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if dnsAliasName == "" {
		return nil, errors.New("parameter dnsAliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dnsAliasName}", url.PathEscape(dnsAliasName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Gets a server DNS alias.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - managedInstanceName - The name of the managed instance.
//   - options - ManagedServerDNSAliasesClientGetOptions contains the optional parameters for the ManagedServerDNSAliasesClient.Get
//     method.
func (client *ManagedServerDNSAliasesClient) Get(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, options *ManagedServerDNSAliasesClientGetOptions) (ManagedServerDNSAliasesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, managedInstanceName, dnsAliasName, options)
	if err != nil {
		return ManagedServerDNSAliasesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ManagedServerDNSAliasesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ManagedServerDNSAliasesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ManagedServerDNSAliasesClient) getCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, dnsAliasName string, options *ManagedServerDNSAliasesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/dnsAliases/{dnsAliasName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if dnsAliasName == "" {
		return nil, errors.New("parameter dnsAliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dnsAliasName}", url.PathEscape(dnsAliasName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ManagedServerDNSAliasesClient) getHandleResponse(resp *http.Response) (ManagedServerDNSAliasesClientGetResponse, error) {
	result := ManagedServerDNSAliasesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedServerDNSAlias); err != nil {
		return ManagedServerDNSAliasesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByManagedInstancePager - Gets a list of managed server DNS aliases for a managed server.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - managedInstanceName - The name of the managed instance.
//   - options - ManagedServerDNSAliasesClientListByManagedInstanceOptions contains the optional parameters for the ManagedServerDNSAliasesClient.NewListByManagedInstancePager
//     method.
func (client *ManagedServerDNSAliasesClient) NewListByManagedInstancePager(resourceGroupName string, managedInstanceName string, options *ManagedServerDNSAliasesClientListByManagedInstanceOptions) *runtime.Pager[ManagedServerDNSAliasesClientListByManagedInstanceResponse] {
	return runtime.NewPager(runtime.PagingHandler[ManagedServerDNSAliasesClientListByManagedInstanceResponse]{
		More: func(page ManagedServerDNSAliasesClientListByManagedInstanceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ManagedServerDNSAliasesClientListByManagedInstanceResponse) (ManagedServerDNSAliasesClientListByManagedInstanceResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByManagedInstanceCreateRequest(ctx, resourceGroupName, managedInstanceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ManagedServerDNSAliasesClientListByManagedInstanceResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ManagedServerDNSAliasesClientListByManagedInstanceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ManagedServerDNSAliasesClientListByManagedInstanceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByManagedInstanceHandleResponse(resp)
		},
	})
}

// listByManagedInstanceCreateRequest creates the ListByManagedInstance request.
func (client *ManagedServerDNSAliasesClient) listByManagedInstanceCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, options *ManagedServerDNSAliasesClientListByManagedInstanceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/dnsAliases"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByManagedInstanceHandleResponse handles the ListByManagedInstance response.
func (client *ManagedServerDNSAliasesClient) listByManagedInstanceHandleResponse(resp *http.Response) (ManagedServerDNSAliasesClientListByManagedInstanceResponse, error) {
	result := ManagedServerDNSAliasesClientListByManagedInstanceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedServerDNSAliasListResult); err != nil {
		return ManagedServerDNSAliasesClientListByManagedInstanceResponse{}, err
	}
	return result, nil
}
