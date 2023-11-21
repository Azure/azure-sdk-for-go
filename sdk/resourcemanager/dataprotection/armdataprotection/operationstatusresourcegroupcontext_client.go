//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdataprotection

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

// OperationStatusResourceGroupContextClient contains the methods for the OperationStatusResourceGroupContext group.
// Don't use this type directly, use NewOperationStatusResourceGroupContextClient() instead.
type OperationStatusResourceGroupContextClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewOperationStatusResourceGroupContextClient creates a new instance of OperationStatusResourceGroupContextClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewOperationStatusResourceGroupContextClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationStatusResourceGroupContextClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &OperationStatusResourceGroupContextClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Gets the operation status for an operation over a ResourceGroup's context.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - OperationStatusResourceGroupContextClientGetOptions contains the optional parameters for the OperationStatusResourceGroupContextClient.Get
//     method.
func (client *OperationStatusResourceGroupContextClient) Get(ctx context.Context, resourceGroupName string, operationID string, options *OperationStatusResourceGroupContextClientGetOptions) (OperationStatusResourceGroupContextClientGetResponse, error) {
	var err error
	const operationName = "OperationStatusResourceGroupContextClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, operationID, options)
	if err != nil {
		return OperationStatusResourceGroupContextClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return OperationStatusResourceGroupContextClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return OperationStatusResourceGroupContextClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *OperationStatusResourceGroupContextClient) getCreateRequest(ctx context.Context, resourceGroupName string, operationID string, options *OperationStatusResourceGroupContextClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/operationStatus/{operationId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if operationID == "" {
		return nil, errors.New("parameter operationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationId}", url.PathEscape(operationID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *OperationStatusResourceGroupContextClient) getHandleResponse(resp *http.Response) (OperationStatusResourceGroupContextClientGetResponse, error) {
	result := OperationStatusResourceGroupContextClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.OperationResource); err != nil {
		return OperationStatusResourceGroupContextClientGetResponse{}, err
	}
	return result, nil
}
