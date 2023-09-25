//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armlogic

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

// IntegrationServiceEnvironmentNetworkHealthClient contains the methods for the IntegrationServiceEnvironmentNetworkHealth group.
// Don't use this type directly, use NewIntegrationServiceEnvironmentNetworkHealthClient() instead.
type IntegrationServiceEnvironmentNetworkHealthClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewIntegrationServiceEnvironmentNetworkHealthClient creates a new instance of IntegrationServiceEnvironmentNetworkHealthClient with the specified values.
//   - subscriptionID - The subscription id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewIntegrationServiceEnvironmentNetworkHealthClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*IntegrationServiceEnvironmentNetworkHealthClient, error) {
	cl, err := arm.NewClient(moduleName+".IntegrationServiceEnvironmentNetworkHealthClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &IntegrationServiceEnvironmentNetworkHealthClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - Gets the integration service environment network health.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroup - The resource group.
//   - integrationServiceEnvironmentName - The integration service environment name.
//   - options - IntegrationServiceEnvironmentNetworkHealthClientGetOptions contains the optional parameters for the IntegrationServiceEnvironmentNetworkHealthClient.Get
//     method.
func (client *IntegrationServiceEnvironmentNetworkHealthClient) Get(ctx context.Context, resourceGroup string, integrationServiceEnvironmentName string, options *IntegrationServiceEnvironmentNetworkHealthClientGetOptions) (IntegrationServiceEnvironmentNetworkHealthClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroup, integrationServiceEnvironmentName, options)
	if err != nil {
		return IntegrationServiceEnvironmentNetworkHealthClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return IntegrationServiceEnvironmentNetworkHealthClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return IntegrationServiceEnvironmentNetworkHealthClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *IntegrationServiceEnvironmentNetworkHealthClient) getCreateRequest(ctx context.Context, resourceGroup string, integrationServiceEnvironmentName string, options *IntegrationServiceEnvironmentNetworkHealthClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Logic/integrationServiceEnvironments/{integrationServiceEnvironmentName}/health/network"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroup == "" {
		return nil, errors.New("parameter resourceGroup cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroup}", url.PathEscape(resourceGroup))
	if integrationServiceEnvironmentName == "" {
		return nil, errors.New("parameter integrationServiceEnvironmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{integrationServiceEnvironmentName}", url.PathEscape(integrationServiceEnvironmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *IntegrationServiceEnvironmentNetworkHealthClient) getHandleResponse(resp *http.Response) (IntegrationServiceEnvironmentNetworkHealthClientGetResponse, error) {
	result := IntegrationServiceEnvironmentNetworkHealthClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Value); err != nil {
		return IntegrationServiceEnvironmentNetworkHealthClientGetResponse{}, err
	}
	return result, nil
}

