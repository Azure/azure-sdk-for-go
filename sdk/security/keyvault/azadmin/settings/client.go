//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package settings

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// Client contains the methods for the Client group.
// Don't use this type directly, use NewClient() instead.
type Client struct {
	endpoint string
	pl       runtime.Pipeline
}

// GetSetting - Retrieves the setting object of a specified setting name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 7.4-preview.1
// settingName - The name of the account setting. Must be a valid settings option.
// options - ClientGetSettingOptions contains the optional parameters for the Client.GetSetting method.
func (client *Client) GetSetting(ctx context.Context, settingName string, options *ClientGetSettingOptions) (ClientGetSettingResponse, error) {
	req, err := client.getSettingCreateRequest(ctx, settingName, options)
	if err != nil {
		return ClientGetSettingResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ClientGetSettingResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ClientGetSettingResponse{}, runtime.NewResponseError(resp)
	}
	return client.getSettingHandleResponse(resp)
}

// getSettingCreateRequest creates the GetSetting request.
func (client *Client) getSettingCreateRequest(ctx context.Context, settingName string, options *ClientGetSettingOptions) (*policy.Request, error) {
	urlPath := "/settings/{setting-name}"
	if settingName == "" {
		return nil, errors.New("parameter settingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{setting-name}", url.PathEscape(settingName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.4-preview.1")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getSettingHandleResponse handles the GetSetting response.
func (client *Client) getSettingHandleResponse(resp *http.Response) (ClientGetSettingResponse, error) {
	result := ClientGetSettingResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Setting); err != nil {
		return ClientGetSettingResponse{}, err
	}
	return result, nil
}

// GetSettings - Retrieves a list of all the available account settings that can be configured.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 7.4-preview.1
// options - ClientGetSettingsOptions contains the optional parameters for the Client.GetSettings method.
func (client *Client) GetSettings(ctx context.Context, options *ClientGetSettingsOptions) (ClientGetSettingsResponse, error) {
	req, err := client.getSettingsCreateRequest(ctx, options)
	if err != nil {
		return ClientGetSettingsResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ClientGetSettingsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ClientGetSettingsResponse{}, runtime.NewResponseError(resp)
	}
	return client.getSettingsHandleResponse(resp)
}

// getSettingsCreateRequest creates the GetSettings request.
func (client *Client) getSettingsCreateRequest(ctx context.Context, options *ClientGetSettingsOptions) (*policy.Request, error) {
	urlPath := "/settings"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.4-preview.1")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getSettingsHandleResponse handles the GetSettings response.
func (client *Client) getSettingsHandleResponse(resp *http.Response) (ClientGetSettingsResponse, error) {
	result := ClientGetSettingsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListResult); err != nil {
		return ClientGetSettingsResponse{}, err
	}
	return result, nil
}

// UpdateSetting - Description of the pool setting to be updated
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 7.4-preview.1
// settingName - The name of the account setting. Must be a valid settings option.
// parameters - The parameters to update an account setting.
// options - ClientUpdateSettingOptions contains the optional parameters for the Client.UpdateSetting method.
func (client *Client) UpdateSetting(ctx context.Context, settingName string, parameters UpdateSettingRequest, options *ClientUpdateSettingOptions) (ClientUpdateSettingResponse, error) {
	req, err := client.updateSettingCreateRequest(ctx, settingName, parameters, options)
	if err != nil {
		return ClientUpdateSettingResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ClientUpdateSettingResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ClientUpdateSettingResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateSettingHandleResponse(resp)
}

// updateSettingCreateRequest creates the UpdateSetting request.
func (client *Client) updateSettingCreateRequest(ctx context.Context, settingName string, parameters UpdateSettingRequest, options *ClientUpdateSettingOptions) (*policy.Request, error) {
	urlPath := "/settings/{setting-name}"
	if settingName == "" {
		return nil, errors.New("parameter settingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{setting-name}", url.PathEscape(settingName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.endpoint, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "7.4-preview.1")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateSettingHandleResponse handles the UpdateSetting response.
func (client *Client) updateSettingHandleResponse(resp *http.Response) (ClientUpdateSettingResponse, error) {
	result := ClientUpdateSettingResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Setting); err != nil {
		return ClientUpdateSettingResponse{}, err
	}
	return result, nil
}
