//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcustomerlockbox

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

// GetClient contains the methods for the Get group.
// Don't use this type directly, use NewGetClient() instead.
type GetClient struct {
	internal *arm.Client
}

// NewGetClient creates a new instance of GetClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewGetClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*GetClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &GetClient{
		internal: cl,
	}
	return client, nil
}

// TenantOptedIn - Get Customer Lockbox request
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-02-28-preview
//   - tenantID - The Azure tenant ID. This is a GUID-formatted string (e.g. 00000000-0000-0000-0000-000000000000)
//   - options - GetClientTenantOptedInOptions contains the optional parameters for the GetClient.TenantOptedIn method.
func (client *GetClient) TenantOptedIn(ctx context.Context, tenantID string, options *GetClientTenantOptedInOptions) (GetClientTenantOptedInResponse, error) {
	var err error
	const operationName = "GetClient.TenantOptedIn"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.tenantOptedInCreateRequest(ctx, tenantID, options)
	if err != nil {
		return GetClientTenantOptedInResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GetClientTenantOptedInResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return GetClientTenantOptedInResponse{}, err
	}
	resp, err := client.tenantOptedInHandleResponse(httpResp)
	return resp, err
}

// tenantOptedInCreateRequest creates the TenantOptedIn request.
func (client *GetClient) tenantOptedInCreateRequest(ctx context.Context, tenantID string, options *GetClientTenantOptedInOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.CustomerLockbox/tenantOptedIn/{tenantId}"
	if tenantID == "" {
		return nil, errors.New("parameter tenantID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{tenantId}", url.PathEscape(tenantID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-02-28-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// tenantOptedInHandleResponse handles the TenantOptedIn response.
func (client *GetClient) tenantOptedInHandleResponse(resp *http.Response) (GetClientTenantOptedInResponse, error) {
	result := GetClientTenantOptedInResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TenantOptInResponse); err != nil {
		return GetClientTenantOptedInResponse{}, err
	}
	return result, nil
}
