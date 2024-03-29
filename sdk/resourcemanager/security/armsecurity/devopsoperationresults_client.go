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

// DevOpsOperationResultsClient contains the methods for the DevOpsOperationResults group.
// Don't use this type directly, use NewDevOpsOperationResultsClient() instead.
type DevOpsOperationResultsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewDevOpsOperationResultsClient creates a new instance of DevOpsOperationResultsClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewDevOpsOperationResultsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DevOpsOperationResultsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &DevOpsOperationResultsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Get devops long running operation result.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - securityConnectorName - The security connector name.
//   - operationResultID - The operation result Id.
//   - options - DevOpsOperationResultsClientGetOptions contains the optional parameters for the DevOpsOperationResultsClient.Get
//     method.
func (client *DevOpsOperationResultsClient) Get(ctx context.Context, resourceGroupName string, securityConnectorName string, operationResultID string, options *DevOpsOperationResultsClientGetOptions) (DevOpsOperationResultsClientGetResponse, error) {
	var err error
	const operationName = "DevOpsOperationResultsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, securityConnectorName, operationResultID, options)
	if err != nil {
		return DevOpsOperationResultsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DevOpsOperationResultsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DevOpsOperationResultsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *DevOpsOperationResultsClient) getCreateRequest(ctx context.Context, resourceGroupName string, securityConnectorName string, operationResultID string, options *DevOpsOperationResultsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/securityConnectors/{securityConnectorName}/devops/default/operationResults/{operationResultId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if securityConnectorName == "" {
		return nil, errors.New("parameter securityConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityConnectorName}", url.PathEscape(securityConnectorName))
	if operationResultID == "" {
		return nil, errors.New("parameter operationResultID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationResultId}", url.PathEscape(operationResultID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DevOpsOperationResultsClient) getHandleResponse(resp *http.Response) (DevOpsOperationResultsClientGetResponse, error) {
	result := DevOpsOperationResultsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.OperationStatusResult); err != nil {
		return DevOpsOperationResultsClientGetResponse{}, err
	}
	return result, nil
}
