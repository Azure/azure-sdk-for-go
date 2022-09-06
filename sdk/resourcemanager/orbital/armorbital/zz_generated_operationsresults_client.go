//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armorbital

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// OperationsResultsClient contains the methods for the OperationsResults group.
// Don't use this type directly, use NewOperationsResultsClient() instead.
type OperationsResultsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewOperationsResultsClient creates a new instance of OperationsResultsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewOperationsResultsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OperationsResultsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &OperationsResultsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginGet - Returns operation results.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// location - The name of Azure region.
// operationID - The ID of an ongoing async operation.
// options - OperationsResultsClientBeginGetOptions contains the optional parameters for the OperationsResultsClient.BeginGet
// method.
func (client *OperationsResultsClient) BeginGet(ctx context.Context, location string, operationID string, options *OperationsResultsClientBeginGetOptions) (*runtime.Poller[OperationsResultsClientGetResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.get(ctx, location, operationID, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[OperationsResultsClientGetResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[OperationsResultsClientGetResponse](options.ResumeToken, client.pl, nil)
	}
}

// Get - Returns operation results.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
func (client *OperationsResultsClient) get(ctx context.Context, location string, operationID string, options *OperationsResultsClientBeginGetOptions) (*http.Response, error) {
	req, err := client.getCreateRequest(ctx, location, operationID, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// getCreateRequest creates the Get request.
func (client *OperationsResultsClient) getCreateRequest(ctx context.Context, location string, operationID string, options *OperationsResultsClientBeginGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Orbital/locations/{location}/operationResults/{operationId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if operationID == "" {
		return nil, errors.New("parameter operationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{operationId}", url.PathEscape(operationID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}
