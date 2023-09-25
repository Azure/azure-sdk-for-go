//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmediaservices

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

// OperationStatusesClient contains the methods for the MediaServicesOperationStatuses group.
// Don't use this type directly, use NewOperationStatusesClient() instead.
type OperationStatusesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewOperationStatusesClient creates a new instance of OperationStatusesClient with the specified values.
//   - subscriptionID - The unique identifier for a Microsoft Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewOperationStatusesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationStatusesClient, error) {
	cl, err := arm.NewClient(moduleName+".OperationStatusesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &OperationStatusesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - Get media service operation status.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01
//   - locationName - Location name.
//   - operationID - Operation ID.
//   - options - OperationStatusesClientGetOptions contains the optional parameters for the OperationStatusesClient.Get method.
func (client *OperationStatusesClient) Get(ctx context.Context, locationName string, operationID string, options *OperationStatusesClientGetOptions) (OperationStatusesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, locationName, operationID, options)
	if err != nil {
		return OperationStatusesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return OperationStatusesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return OperationStatusesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *OperationStatusesClient) getCreateRequest(ctx context.Context, locationName string, operationID string, options *OperationStatusesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Media/locations/{locationName}/mediaServicesOperationStatuses/{operationId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if operationID == "" {
		return nil, errors.New("parameter operationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationId}", url.PathEscape(operationID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *OperationStatusesClient) getHandleResponse(resp *http.Response) (OperationStatusesClientGetResponse, error) {
	result := OperationStatusesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MediaServiceOperationStatus); err != nil {
		return OperationStatusesClientGetResponse{}, err
	}
	return result, nil
}

