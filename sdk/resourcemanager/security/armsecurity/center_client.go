//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// CenterClient contains the methods for the SecurityCenter group.
// Don't use this type directly, use NewCenterClient() instead.
type CenterClient struct {
	internal *arm.Client
}

// NewCenterClient creates a new instance of CenterClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCenterClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*CenterClient, error) {
	cl, err := arm.NewClient(moduleName+".CenterClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CenterClient{
	internal: cl,
	}
	return client, nil
}

// GetSensitivitySettings - Gets data sensitivity settings for sensitive data discovery
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-15-preview
//   - options - CenterClientGetSensitivitySettingsOptions contains the optional parameters for the CenterClient.GetSensitivitySettings
//     method.
func (client *CenterClient) GetSensitivitySettings(ctx context.Context, options *CenterClientGetSensitivitySettingsOptions) (CenterClientGetSensitivitySettingsResponse, error) {
	var err error
	req, err := client.getSensitivitySettingsCreateRequest(ctx, options)
	if err != nil {
		return CenterClientGetSensitivitySettingsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CenterClientGetSensitivitySettingsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CenterClientGetSensitivitySettingsResponse{}, err
	}
	resp, err := client.getSensitivitySettingsHandleResponse(httpResp)
	return resp, err
}

// getSensitivitySettingsCreateRequest creates the GetSensitivitySettings request.
func (client *CenterClient) getSensitivitySettingsCreateRequest(ctx context.Context, options *CenterClientGetSensitivitySettingsOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Security/sensitivitySettings/current"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-15-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getSensitivitySettingsHandleResponse handles the GetSensitivitySettings response.
func (client *CenterClient) getSensitivitySettingsHandleResponse(resp *http.Response) (CenterClientGetSensitivitySettingsResponse, error) {
	result := CenterClientGetSensitivitySettingsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GetSensitivitySettingsResponse); err != nil {
		return CenterClientGetSensitivitySettingsResponse{}, err
	}
	return result, nil
}

// UpdateSensitivitySettings - Updates data sensitivity settings for sensitive data discovery
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-15-preview
//   - sensitivitySettings - The data sensitivity settings to update
//   - options - CenterClientUpdateSensitivitySettingsOptions contains the optional parameters for the CenterClient.UpdateSensitivitySettings
//     method.
func (client *CenterClient) UpdateSensitivitySettings(ctx context.Context, sensitivitySettings UpdateSensitivitySettingsRequest, options *CenterClientUpdateSensitivitySettingsOptions) (CenterClientUpdateSensitivitySettingsResponse, error) {
	var err error
	req, err := client.updateSensitivitySettingsCreateRequest(ctx, sensitivitySettings, options)
	if err != nil {
		return CenterClientUpdateSensitivitySettingsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CenterClientUpdateSensitivitySettingsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CenterClientUpdateSensitivitySettingsResponse{}, err
	}
	resp, err := client.updateSensitivitySettingsHandleResponse(httpResp)
	return resp, err
}

// updateSensitivitySettingsCreateRequest creates the UpdateSensitivitySettings request.
func (client *CenterClient) updateSensitivitySettingsCreateRequest(ctx context.Context, sensitivitySettings UpdateSensitivitySettingsRequest, options *CenterClientUpdateSensitivitySettingsOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Security/sensitivitySettings/current"
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-15-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, sensitivitySettings); err != nil {
	return nil, err
}
	return req, nil
}

// updateSensitivitySettingsHandleResponse handles the UpdateSensitivitySettings response.
func (client *CenterClient) updateSensitivitySettingsHandleResponse(resp *http.Response) (CenterClientUpdateSensitivitySettingsResponse, error) {
	result := CenterClientUpdateSensitivitySettingsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GetSensitivitySettingsResponse); err != nil {
		return CenterClientUpdateSensitivitySettingsResponse{}, err
	}
	return result, nil
}

