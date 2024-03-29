//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armquantum

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

// WorkspaceClient contains the methods for the Workspace group.
// Don't use this type directly, use NewWorkspaceClient() instead.
type WorkspaceClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewWorkspaceClient creates a new instance of WorkspaceClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkspaceClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkspaceClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkspaceClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CheckNameAvailability - Check the availability of the resource name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-11-13-preview
//   - locationName - Location.
//   - checkNameAvailabilityParameters - The name and type of the resource.
//   - options - WorkspaceClientCheckNameAvailabilityOptions contains the optional parameters for the WorkspaceClient.CheckNameAvailability
//     method.
func (client *WorkspaceClient) CheckNameAvailability(ctx context.Context, locationName string, checkNameAvailabilityParameters CheckNameAvailabilityParameters, options *WorkspaceClientCheckNameAvailabilityOptions) (WorkspaceClientCheckNameAvailabilityResponse, error) {
	var err error
	const operationName = "WorkspaceClient.CheckNameAvailability"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.checkNameAvailabilityCreateRequest(ctx, locationName, checkNameAvailabilityParameters, options)
	if err != nil {
		return WorkspaceClientCheckNameAvailabilityResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkspaceClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WorkspaceClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.checkNameAvailabilityHandleResponse(httpResp)
	return resp, err
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *WorkspaceClient) checkNameAvailabilityCreateRequest(ctx context.Context, locationName string, checkNameAvailabilityParameters CheckNameAvailabilityParameters, options *WorkspaceClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Quantum/locations/{locationName}/checkNameAvailability"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-11-13-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, checkNameAvailabilityParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *WorkspaceClient) checkNameAvailabilityHandleResponse(resp *http.Response) (WorkspaceClientCheckNameAvailabilityResponse, error) {
	result := WorkspaceClientCheckNameAvailabilityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameAvailabilityResult); err != nil {
		return WorkspaceClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// ListKeys - Get the keys to use with the Quantum APIs. A key is used to authenticate and authorize access to the Quantum
// REST APIs. Only one key is needed at a time; two are given to provide seamless key
// regeneration.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-11-13-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the quantum workspace resource.
//   - options - WorkspaceClientListKeysOptions contains the optional parameters for the WorkspaceClient.ListKeys method.
func (client *WorkspaceClient) ListKeys(ctx context.Context, resourceGroupName string, workspaceName string, options *WorkspaceClientListKeysOptions) (WorkspaceClientListKeysResponse, error) {
	var err error
	const operationName = "WorkspaceClient.ListKeys"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listKeysCreateRequest(ctx, resourceGroupName, workspaceName, options)
	if err != nil {
		return WorkspaceClientListKeysResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkspaceClientListKeysResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WorkspaceClientListKeysResponse{}, err
	}
	resp, err := client.listKeysHandleResponse(httpResp)
	return resp, err
}

// listKeysCreateRequest creates the ListKeys request.
func (client *WorkspaceClient) listKeysCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, options *WorkspaceClientListKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Quantum/workspaces/{workspaceName}/listKeys"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-11-13-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listKeysHandleResponse handles the ListKeys response.
func (client *WorkspaceClient) listKeysHandleResponse(resp *http.Response) (WorkspaceClientListKeysResponse, error) {
	result := WorkspaceClientListKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListKeysResult); err != nil {
		return WorkspaceClientListKeysResponse{}, err
	}
	return result, nil
}

// RegenerateKeys - Regenerate either the primary or secondary key for use with the Quantum APIs. The old key will stop working
// immediately.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-11-13-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the quantum workspace resource.
//   - keySpecification - Which key to regenerate: primary or secondary.
//   - options - WorkspaceClientRegenerateKeysOptions contains the optional parameters for the WorkspaceClient.RegenerateKeys
//     method.
func (client *WorkspaceClient) RegenerateKeys(ctx context.Context, resourceGroupName string, workspaceName string, keySpecification APIKeys, options *WorkspaceClientRegenerateKeysOptions) (WorkspaceClientRegenerateKeysResponse, error) {
	var err error
	const operationName = "WorkspaceClient.RegenerateKeys"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.regenerateKeysCreateRequest(ctx, resourceGroupName, workspaceName, keySpecification, options)
	if err != nil {
		return WorkspaceClientRegenerateKeysResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkspaceClientRegenerateKeysResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return WorkspaceClientRegenerateKeysResponse{}, err
	}
	return WorkspaceClientRegenerateKeysResponse{}, nil
}

// regenerateKeysCreateRequest creates the RegenerateKeys request.
func (client *WorkspaceClient) regenerateKeysCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, keySpecification APIKeys, options *WorkspaceClientRegenerateKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Quantum/workspaces/{workspaceName}/regenerateKey"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-11-13-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, keySpecification); err != nil {
		return nil, err
	}
	return req, nil
}
