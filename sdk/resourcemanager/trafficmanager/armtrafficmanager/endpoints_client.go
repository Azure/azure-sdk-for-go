//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armtrafficmanager

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

// EndpointsClient contains the methods for the Endpoints group.
// Don't use this type directly, use NewEndpointsClient() instead.
type EndpointsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewEndpointsClient creates a new instance of EndpointsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewEndpointsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*EndpointsClient, error) {
	cl, err := arm.NewClient(moduleName+".EndpointsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &EndpointsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a Traffic Manager endpoint.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - profileName - The name of the Traffic Manager profile.
//   - endpointType - The type of the Traffic Manager endpoint to be created or updated.
//   - endpointName - The name of the Traffic Manager endpoint to be created or updated.
//   - parameters - The Traffic Manager endpoint parameters supplied to the CreateOrUpdate operation.
//   - options - EndpointsClientCreateOrUpdateOptions contains the optional parameters for the EndpointsClient.CreateOrUpdate
//     method.
func (client *EndpointsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, parameters Endpoint, options *EndpointsClientCreateOrUpdateOptions) (EndpointsClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, profileName, endpointType, endpointName, parameters, options)
	if err != nil {
		return EndpointsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return EndpointsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *EndpointsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, parameters Endpoint, options *EndpointsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficmanagerprofiles/{profileName}/{endpointType}/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointType == "" {
		return nil, errors.New("parameter endpointType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointType}", url.PathEscape(string(endpointType)))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *EndpointsClient) createOrUpdateHandleResponse(resp *http.Response) (EndpointsClientCreateOrUpdateResponse, error) {
	result := EndpointsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Endpoint); err != nil {
		return EndpointsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a Traffic Manager endpoint.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - profileName - The name of the Traffic Manager profile.
//   - endpointType - The type of the Traffic Manager endpoint to be deleted.
//   - endpointName - The name of the Traffic Manager endpoint to be deleted.
//   - options - EndpointsClientDeleteOptions contains the optional parameters for the EndpointsClient.Delete method.
func (client *EndpointsClient) Delete(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, options *EndpointsClientDeleteOptions) (EndpointsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, profileName, endpointType, endpointName, options)
	if err != nil {
		return EndpointsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return EndpointsClientDeleteResponse{}, err
	}
	resp, err := client.deleteHandleResponse(httpResp)
	return resp, err
}

// deleteCreateRequest creates the Delete request.
func (client *EndpointsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, options *EndpointsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficmanagerprofiles/{profileName}/{endpointType}/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointType == "" {
		return nil, errors.New("parameter endpointType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointType}", url.PathEscape(string(endpointType)))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *EndpointsClient) deleteHandleResponse(resp *http.Response) (EndpointsClientDeleteResponse, error) {
	result := EndpointsClientDeleteResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeleteOperationResult); err != nil {
		return EndpointsClientDeleteResponse{}, err
	}
	return result, nil
}

// Get - Gets a Traffic Manager endpoint.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - profileName - The name of the Traffic Manager profile.
//   - endpointType - The type of the Traffic Manager endpoint.
//   - endpointName - The name of the Traffic Manager endpoint.
//   - options - EndpointsClientGetOptions contains the optional parameters for the EndpointsClient.Get method.
func (client *EndpointsClient) Get(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, options *EndpointsClientGetOptions) (EndpointsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, profileName, endpointType, endpointName, options)
	if err != nil {
		return EndpointsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EndpointsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *EndpointsClient) getCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, options *EndpointsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficmanagerprofiles/{profileName}/{endpointType}/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointType == "" {
		return nil, errors.New("parameter endpointType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointType}", url.PathEscape(string(endpointType)))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *EndpointsClient) getHandleResponse(resp *http.Response) (EndpointsClientGetResponse, error) {
	result := EndpointsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Endpoint); err != nil {
		return EndpointsClientGetResponse{}, err
	}
	return result, nil
}

// Update - Update a Traffic Manager endpoint.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - profileName - The name of the Traffic Manager profile.
//   - endpointType - The type of the Traffic Manager endpoint to be updated.
//   - endpointName - The name of the Traffic Manager endpoint to be updated.
//   - parameters - The Traffic Manager endpoint parameters supplied to the Update operation.
//   - options - EndpointsClientUpdateOptions contains the optional parameters for the EndpointsClient.Update method.
func (client *EndpointsClient) Update(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, parameters Endpoint, options *EndpointsClientUpdateOptions) (EndpointsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, profileName, endpointType, endpointName, parameters, options)
	if err != nil {
		return EndpointsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EndpointsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *EndpointsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, profileName string, endpointType EndpointType, endpointName string, parameters Endpoint, options *EndpointsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficmanagerprofiles/{profileName}/{endpointType}/{endpointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if profileName == "" {
		return nil, errors.New("parameter profileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{profileName}", url.PathEscape(profileName))
	if endpointType == "" {
		return nil, errors.New("parameter endpointType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointType}", url.PathEscape(string(endpointType)))
	if endpointName == "" {
		return nil, errors.New("parameter endpointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", url.PathEscape(endpointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *EndpointsClient) updateHandleResponse(resp *http.Response) (EndpointsClientUpdateResponse, error) {
	result := EndpointsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Endpoint); err != nil {
		return EndpointsClientUpdateResponse{}, err
	}
	return result, nil
}

