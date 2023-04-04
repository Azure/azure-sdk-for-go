//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhealthcareapis

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

// OperationResultsClient contains the methods for the OperationResults group.
// Don't use this type directly, use NewOperationResultsClient() instead.
type OperationResultsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewOperationResultsClient creates a new instance of OperationResultsClient with the specified values.
//   - subscriptionID - The subscription identifier.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewOperationResultsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationResultsClient, error) {
	cl, err := arm.NewClient(moduleName+".OperationResultsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &OperationResultsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Get the operation result for a long running operation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-31-preview
//   - locationName - The location of the operation.
//   - operationResultID - The ID of the operation result to get.
//   - options - OperationResultsClientGetOptions contains the optional parameters for the OperationResultsClient.Get method.
func (client *OperationResultsClient) Get(ctx context.Context, locationName string, operationResultID string, options *OperationResultsClientGetOptions) (OperationResultsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, locationName, operationResultID, options)
	if err != nil {
		return OperationResultsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return OperationResultsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return OperationResultsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *OperationResultsClient) getCreateRequest(ctx context.Context, locationName string, operationResultID string, options *OperationResultsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HealthcareApis/locations/{locationName}/operationresults/{operationResultId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if operationResultID == "" {
		return nil, errors.New("parameter operationResultID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationResultId}", url.PathEscape(operationResultID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *OperationResultsClient) getHandleResponse(resp *http.Response) (OperationResultsClientGetResponse, error) {
	result := OperationResultsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.OperationResultsDescription); err != nil {
		return OperationResultsClientGetResponse{}, err
	}
	return result, nil
}
