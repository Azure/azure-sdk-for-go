//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// APIManagementClient contains the methods for the APIManagementClient group.
// Don't use this type directly, use NewAPIManagementClient() instead.
type APIManagementClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewAPIManagementClient creates a new instance of APIManagementClient with the specified values.
func NewAPIManagementClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *APIManagementClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &APIManagementClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// BeginPerformConnectivityCheckAsync - Performs a connectivity check between the API Management service and a given destination, and returns metrics for
// the connection, as well as errors encountered while trying to establish it.
// If the operation fails it returns the *ErrorResponse error type.
func (client *APIManagementClient) BeginPerformConnectivityCheckAsync(ctx context.Context, resourceGroupName string, serviceName string, connectivityCheckRequestParams ConnectivityCheckRequest, options *APIManagementClientBeginPerformConnectivityCheckAsyncOptions) (APIManagementClientPerformConnectivityCheckAsyncPollerResponse, error) {
	resp, err := client.performConnectivityCheckAsync(ctx, resourceGroupName, serviceName, connectivityCheckRequestParams, options)
	if err != nil {
		return APIManagementClientPerformConnectivityCheckAsyncPollerResponse{}, err
	}
	result := APIManagementClientPerformConnectivityCheckAsyncPollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("APIManagementClient.PerformConnectivityCheckAsync", "location", resp, client.pl, client.performConnectivityCheckAsyncHandleError)
	if err != nil {
		return APIManagementClientPerformConnectivityCheckAsyncPollerResponse{}, err
	}
	result.Poller = &APIManagementClientPerformConnectivityCheckAsyncPoller{
		pt: pt,
	}
	return result, nil
}

// PerformConnectivityCheckAsync - Performs a connectivity check between the API Management service and a given destination, and returns metrics for the
// connection, as well as errors encountered while trying to establish it.
// If the operation fails it returns the *ErrorResponse error type.
func (client *APIManagementClient) performConnectivityCheckAsync(ctx context.Context, resourceGroupName string, serviceName string, connectivityCheckRequestParams ConnectivityCheckRequest, options *APIManagementClientBeginPerformConnectivityCheckAsyncOptions) (*http.Response, error) {
	req, err := client.performConnectivityCheckAsyncCreateRequest(ctx, resourceGroupName, serviceName, connectivityCheckRequestParams, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.performConnectivityCheckAsyncHandleError(resp)
	}
	return resp, nil
}

// performConnectivityCheckAsyncCreateRequest creates the PerformConnectivityCheckAsync request.
func (client *APIManagementClient) performConnectivityCheckAsyncCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, connectivityCheckRequestParams ConnectivityCheckRequest, options *APIManagementClientBeginPerformConnectivityCheckAsyncOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/connectivityCheck"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, connectivityCheckRequestParams)
}

// performConnectivityCheckAsyncHandleError handles the PerformConnectivityCheckAsync error response.
func (client *APIManagementClient) performConnectivityCheckAsyncHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
