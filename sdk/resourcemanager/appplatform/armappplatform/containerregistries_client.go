//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappplatform

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

// ContainerRegistriesClient contains the methods for the ContainerRegistries group.
// Don't use this type directly, use NewContainerRegistriesClient() instead.
type ContainerRegistriesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewContainerRegistriesClient creates a new instance of ContainerRegistriesClient with the specified values.
//   - subscriptionID - Gets subscription ID which uniquely identify the Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewContainerRegistriesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContainerRegistriesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ContainerRegistriesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update container registry resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - containerRegistryName - The name of the container registry.
//   - containerRegistryResource - Parameters for the create or update operation
//   - options - ContainerRegistriesClientBeginCreateOrUpdateOptions contains the optional parameters for the ContainerRegistriesClient.BeginCreateOrUpdate
//     method.
func (client *ContainerRegistriesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, containerRegistryResource ContainerRegistryResource, options *ContainerRegistriesClientBeginCreateOrUpdateOptions) (*runtime.Poller[ContainerRegistriesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, serviceName, containerRegistryName, containerRegistryResource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContainerRegistriesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ContainerRegistriesClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create or update container registry resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
func (client *ContainerRegistriesClient) createOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, containerRegistryResource ContainerRegistryResource, options *ContainerRegistriesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ContainerRegistriesClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, containerRegistryName, containerRegistryResource, options)
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
func (client *ContainerRegistriesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, containerRegistryResource ContainerRegistryResource, options *ContainerRegistriesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/containerRegistries/{containerRegistryName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if containerRegistryName == "" {
		return nil, errors.New("parameter containerRegistryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerRegistryName}", url.PathEscape(containerRegistryName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, containerRegistryResource); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete a container registry resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - containerRegistryName - The name of the container registry.
//   - options - ContainerRegistriesClientBeginDeleteOptions contains the optional parameters for the ContainerRegistriesClient.BeginDelete
//     method.
func (client *ContainerRegistriesClient) BeginDelete(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, options *ContainerRegistriesClientBeginDeleteOptions) (*runtime.Poller[ContainerRegistriesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, serviceName, containerRegistryName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContainerRegistriesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ContainerRegistriesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a container registry resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
func (client *ContainerRegistriesClient) deleteOperation(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, options *ContainerRegistriesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ContainerRegistriesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, containerRegistryName, options)
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
func (client *ContainerRegistriesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, options *ContainerRegistriesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/containerRegistries/{containerRegistryName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if containerRegistryName == "" {
		return nil, errors.New("parameter containerRegistryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerRegistryName}", url.PathEscape(containerRegistryName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the container registries resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - containerRegistryName - The name of the container registry.
//   - options - ContainerRegistriesClientGetOptions contains the optional parameters for the ContainerRegistriesClient.Get method.
func (client *ContainerRegistriesClient) Get(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, options *ContainerRegistriesClientGetOptions) (ContainerRegistriesClientGetResponse, error) {
	var err error
	const operationName = "ContainerRegistriesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, containerRegistryName, options)
	if err != nil {
		return ContainerRegistriesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ContainerRegistriesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ContainerRegistriesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ContainerRegistriesClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, options *ContainerRegistriesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/containerRegistries/{containerRegistryName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if containerRegistryName == "" {
		return nil, errors.New("parameter containerRegistryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerRegistryName}", url.PathEscape(containerRegistryName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ContainerRegistriesClient) getHandleResponse(resp *http.Response) (ContainerRegistriesClientGetResponse, error) {
	result := ContainerRegistriesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerRegistryResource); err != nil {
		return ContainerRegistriesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List container registries resource.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - options - ContainerRegistriesClientListOptions contains the optional parameters for the ContainerRegistriesClient.NewListPager
//     method.
func (client *ContainerRegistriesClient) NewListPager(resourceGroupName string, serviceName string, options *ContainerRegistriesClientListOptions) *runtime.Pager[ContainerRegistriesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ContainerRegistriesClientListResponse]{
		More: func(page ContainerRegistriesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ContainerRegistriesClientListResponse) (ContainerRegistriesClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ContainerRegistriesClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, serviceName, options)
			}, nil)
			if err != nil {
				return ContainerRegistriesClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *ContainerRegistriesClient) listCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, options *ContainerRegistriesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/containerRegistries"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ContainerRegistriesClient) listHandleResponse(resp *http.Response) (ContainerRegistriesClientListResponse, error) {
	result := ContainerRegistriesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContainerRegistryResourceCollection); err != nil {
		return ContainerRegistriesClientListResponse{}, err
	}
	return result, nil
}

// BeginValidate - Check if the container registry properties are valid.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - containerRegistryName - The name of the container registry.
//   - containerRegistryProperties - Parameters for the validate operation
//   - options - ContainerRegistriesClientBeginValidateOptions contains the optional parameters for the ContainerRegistriesClient.BeginValidate
//     method.
func (client *ContainerRegistriesClient) BeginValidate(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, containerRegistryProperties ContainerRegistryProperties, options *ContainerRegistriesClientBeginValidateOptions) (*runtime.Poller[ContainerRegistriesClientValidateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.validate(ctx, resourceGroupName, serviceName, containerRegistryName, containerRegistryProperties, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ContainerRegistriesClientValidateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ContainerRegistriesClientValidateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Validate - Check if the container registry properties are valid.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
func (client *ContainerRegistriesClient) validate(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, containerRegistryProperties ContainerRegistryProperties, options *ContainerRegistriesClientBeginValidateOptions) (*http.Response, error) {
	var err error
	const operationName = "ContainerRegistriesClient.BeginValidate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.validateCreateRequest(ctx, resourceGroupName, serviceName, containerRegistryName, containerRegistryProperties, options)
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

// validateCreateRequest creates the Validate request.
func (client *ContainerRegistriesClient) validateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, containerRegistryName string, containerRegistryProperties ContainerRegistryProperties, options *ContainerRegistriesClientBeginValidateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/containerRegistries/{containerRegistryName}/validate"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if containerRegistryName == "" {
		return nil, errors.New("parameter containerRegistryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{containerRegistryName}", url.PathEscape(containerRegistryName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, containerRegistryProperties); err != nil {
		return nil, err
	}
	return req, nil
}
