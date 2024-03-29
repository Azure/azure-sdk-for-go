//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

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

// MdeOnboardingsClient contains the methods for the MdeOnboardings group.
// Don't use this type directly, use NewMdeOnboardingsClient() instead.
type MdeOnboardingsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewMdeOnboardingsClient creates a new instance of MdeOnboardingsClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewMdeOnboardingsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MdeOnboardingsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &MdeOnboardingsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - The default configuration or data needed to onboard the machine to MDE
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-10-01-preview
//   - options - MdeOnboardingsClientGetOptions contains the optional parameters for the MdeOnboardingsClient.Get method.
func (client *MdeOnboardingsClient) Get(ctx context.Context, options *MdeOnboardingsClientGetOptions) (MdeOnboardingsClientGetResponse, error) {
	var err error
	const operationName = "MdeOnboardingsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, options)
	if err != nil {
		return MdeOnboardingsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return MdeOnboardingsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return MdeOnboardingsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *MdeOnboardingsClient) getCreateRequest(ctx context.Context, options *MdeOnboardingsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/mdeOnboardings/default"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *MdeOnboardingsClient) getHandleResponse(resp *http.Response) (MdeOnboardingsClientGetResponse, error) {
	result := MdeOnboardingsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MdeOnboardingData); err != nil {
		return MdeOnboardingsClientGetResponse{}, err
	}
	return result, nil
}

// List - The configuration or data needed to onboard the machine to MDE
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-10-01-preview
//   - options - MdeOnboardingsClientListOptions contains the optional parameters for the MdeOnboardingsClient.List method.
func (client *MdeOnboardingsClient) List(ctx context.Context, options *MdeOnboardingsClientListOptions) (MdeOnboardingsClientListResponse, error) {
	var err error
	const operationName = "MdeOnboardingsClient.List"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listCreateRequest(ctx, options)
	if err != nil {
		return MdeOnboardingsClientListResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return MdeOnboardingsClientListResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return MdeOnboardingsClientListResponse{}, err
	}
	resp, err := client.listHandleResponse(httpResp)
	return resp, err
}

// listCreateRequest creates the List request.
func (client *MdeOnboardingsClient) listCreateRequest(ctx context.Context, options *MdeOnboardingsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/mdeOnboardings"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *MdeOnboardingsClient) listHandleResponse(resp *http.Response) (MdeOnboardingsClientListResponse, error) {
	result := MdeOnboardingsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MdeOnboardingDataList); err != nil {
		return MdeOnboardingsClientListResponse{}, err
	}
	return result, nil
}
