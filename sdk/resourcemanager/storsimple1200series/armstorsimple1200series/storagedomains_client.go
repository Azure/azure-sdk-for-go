//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstorsimple1200series

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

// StorageDomainsClient contains the methods for the StorageDomains group.
// Don't use this type directly, use NewStorageDomainsClient() instead.
type StorageDomainsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewStorageDomainsClient creates a new instance of StorageDomainsClient with the specified values.
//   - subscriptionID - The subscription id
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewStorageDomainsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*StorageDomainsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &StorageDomainsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates the storage domain.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2016-10-01
//   - storageDomainName - The storage domain name.
//   - resourceGroupName - The resource group name
//   - managerName - The manager name
//   - storageDomain - The storageDomain.
//   - options - StorageDomainsClientBeginCreateOrUpdateOptions contains the optional parameters for the StorageDomainsClient.BeginCreateOrUpdate
//     method.
func (client *StorageDomainsClient) BeginCreateOrUpdate(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, storageDomain StorageDomain, options *StorageDomainsClientBeginCreateOrUpdateOptions) (*runtime.Poller[StorageDomainsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, storageDomainName, resourceGroupName, managerName, storageDomain, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[StorageDomainsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[StorageDomainsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates or updates the storage domain.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2016-10-01
func (client *StorageDomainsClient) createOrUpdate(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, storageDomain StorageDomain, options *StorageDomainsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "StorageDomainsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, storageDomainName, resourceGroupName, managerName, storageDomain, options)
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
func (client *StorageDomainsClient) createOrUpdateCreateRequest(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, storageDomain StorageDomain, options *StorageDomainsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/storageDomains/{storageDomainName}"
	if storageDomainName == "" {
		return nil, errors.New("parameter storageDomainName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageDomainName}", url.PathEscape(storageDomainName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managerName == "" {
		return nil, errors.New("parameter managerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managerName}", url.PathEscape(managerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, storageDomain); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes the storage domain.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2016-10-01
//   - storageDomainName - The storage domain name.
//   - resourceGroupName - The resource group name
//   - managerName - The manager name
//   - options - StorageDomainsClientBeginDeleteOptions contains the optional parameters for the StorageDomainsClient.BeginDelete
//     method.
func (client *StorageDomainsClient) BeginDelete(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, options *StorageDomainsClientBeginDeleteOptions) (*runtime.Poller[StorageDomainsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, storageDomainName, resourceGroupName, managerName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[StorageDomainsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[StorageDomainsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes the storage domain.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2016-10-01
func (client *StorageDomainsClient) deleteOperation(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, options *StorageDomainsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "StorageDomainsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, storageDomainName, resourceGroupName, managerName, options)
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
func (client *StorageDomainsClient) deleteCreateRequest(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, options *StorageDomainsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/storageDomains/{storageDomainName}"
	if storageDomainName == "" {
		return nil, errors.New("parameter storageDomainName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageDomainName}", url.PathEscape(storageDomainName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managerName == "" {
		return nil, errors.New("parameter managerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managerName}", url.PathEscape(managerName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Returns the properties of the specified storage domain name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2016-10-01
//   - storageDomainName - The storage domain name.
//   - resourceGroupName - The resource group name
//   - managerName - The manager name
//   - options - StorageDomainsClientGetOptions contains the optional parameters for the StorageDomainsClient.Get method.
func (client *StorageDomainsClient) Get(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, options *StorageDomainsClientGetOptions) (StorageDomainsClientGetResponse, error) {
	var err error
	const operationName = "StorageDomainsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, storageDomainName, resourceGroupName, managerName, options)
	if err != nil {
		return StorageDomainsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return StorageDomainsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return StorageDomainsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *StorageDomainsClient) getCreateRequest(ctx context.Context, storageDomainName string, resourceGroupName string, managerName string, options *StorageDomainsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/storageDomains/{storageDomainName}"
	if storageDomainName == "" {
		return nil, errors.New("parameter storageDomainName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageDomainName}", url.PathEscape(storageDomainName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managerName == "" {
		return nil, errors.New("parameter managerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managerName}", url.PathEscape(managerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *StorageDomainsClient) getHandleResponse(resp *http.Response) (StorageDomainsClientGetResponse, error) {
	result := StorageDomainsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StorageDomain); err != nil {
		return StorageDomainsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByManagerPager - Retrieves all the storage domains in a manager.
//
// Generated from API version 2016-10-01
//   - resourceGroupName - The resource group name
//   - managerName - The manager name
//   - options - StorageDomainsClientListByManagerOptions contains the optional parameters for the StorageDomainsClient.NewListByManagerPager
//     method.
func (client *StorageDomainsClient) NewListByManagerPager(resourceGroupName string, managerName string, options *StorageDomainsClientListByManagerOptions) *runtime.Pager[StorageDomainsClientListByManagerResponse] {
	return runtime.NewPager(runtime.PagingHandler[StorageDomainsClientListByManagerResponse]{
		More: func(page StorageDomainsClientListByManagerResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *StorageDomainsClientListByManagerResponse) (StorageDomainsClientListByManagerResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "StorageDomainsClient.NewListByManagerPager")
			req, err := client.listByManagerCreateRequest(ctx, resourceGroupName, managerName, options)
			if err != nil {
				return StorageDomainsClientListByManagerResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return StorageDomainsClientListByManagerResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return StorageDomainsClientListByManagerResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByManagerHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByManagerCreateRequest creates the ListByManager request.
func (client *StorageDomainsClient) listByManagerCreateRequest(ctx context.Context, resourceGroupName string, managerName string, options *StorageDomainsClientListByManagerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorSimple/managers/{managerName}/storageDomains"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managerName == "" {
		return nil, errors.New("parameter managerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managerName}", url.PathEscape(managerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2016-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByManagerHandleResponse handles the ListByManager response.
func (client *StorageDomainsClient) listByManagerHandleResponse(resp *http.Response) (StorageDomainsClientListByManagerResponse, error) {
	result := StorageDomainsClientListByManagerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StorageDomainList); err != nil {
		return StorageDomainsClientListByManagerResponse{}, err
	}
	return result, nil
}
