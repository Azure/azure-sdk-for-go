//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsynapse

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

// KustoPoolChildResourceClient contains the methods for the KustoPoolChildResource group.
// Don't use this type directly, use NewKustoPoolChildResourceClient() instead.
type KustoPoolChildResourceClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewKustoPoolChildResourceClient creates a new instance of KustoPoolChildResourceClient with the specified values.
func NewKustoPoolChildResourceClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *KustoPoolChildResourceClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &KustoPoolChildResourceClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// CheckNameAvailability - Checks that the Kusto Pool child resource name is valid and is not already in use.
// If the operation fails it returns the *ErrorResponse error type.
func (client *KustoPoolChildResourceClient) CheckNameAvailability(ctx context.Context, workspaceName string, kustoPoolName string, resourceGroupName string, resourceName DatabaseCheckNameRequest, options *KustoPoolChildResourceCheckNameAvailabilityOptions) (KustoPoolChildResourceCheckNameAvailabilityResponse, error) {
	req, err := client.checkNameAvailabilityCreateRequest(ctx, workspaceName, kustoPoolName, resourceGroupName, resourceName, options)
	if err != nil {
		return KustoPoolChildResourceCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return KustoPoolChildResourceCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return KustoPoolChildResourceCheckNameAvailabilityResponse{}, client.checkNameAvailabilityHandleError(resp)
	}
	return client.checkNameAvailabilityHandleResponse(resp)
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *KustoPoolChildResourceClient) checkNameAvailabilityCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, resourceGroupName string, resourceName DatabaseCheckNameRequest, options *KustoPoolChildResourceCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/checkNameAvailability"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, resourceName)
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *KustoPoolChildResourceClient) checkNameAvailabilityHandleResponse(resp *http.Response) (KustoPoolChildResourceCheckNameAvailabilityResponse, error) {
	result := KustoPoolChildResourceCheckNameAvailabilityResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.CheckNameResult); err != nil {
		return KustoPoolChildResourceCheckNameAvailabilityResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// checkNameAvailabilityHandleError handles the CheckNameAvailability error response.
func (client *KustoPoolChildResourceClient) checkNameAvailabilityHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
