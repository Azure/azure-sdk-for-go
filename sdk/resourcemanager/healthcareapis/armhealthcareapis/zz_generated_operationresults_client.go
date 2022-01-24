//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhealthcareapis

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// OperationResultsClient contains the methods for the OperationResults group.
// Don't use this type directly, use NewOperationResultsClient() instead.
type OperationResultsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewOperationResultsClient creates a new instance of OperationResultsClient with the specified values.
// subscriptionID - The subscription identifier.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewOperationResultsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *OperationResultsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &OperationResultsClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// Get - Get the operation result for a long running operation.
// If the operation fails it returns an *azcore.ResponseError type.
// locationName - The location of the operation.
// operationResultID - The ID of the operation result to get.
// options - OperationResultsClientGetOptions contains the optional parameters for the OperationResultsClient.Get method.
func (client *OperationResultsClient) Get(ctx context.Context, locationName string, operationResultID string, options *OperationResultsClientGetOptions) (OperationResultsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, locationName, operationResultID, options)
	if err != nil {
		return OperationResultsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *OperationResultsClient) getHandleResponse(resp *http.Response) (OperationResultsClientGetResponse, error) {
	result := OperationResultsClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.OperationResultsDescription); err != nil {
		return OperationResultsClientGetResponse{}, err
	}
	return result, nil
}
