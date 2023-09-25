//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridcontainerservice

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

// AgentPoolClient contains the methods for the AgentPool group.
// Don't use this type directly, use NewAgentPoolClient() instead.
type AgentPoolClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewAgentPoolClient creates a new instance of AgentPoolClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAgentPoolClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AgentPoolClient, error) {
	cl, err := arm.NewClient(moduleName+".AgentPoolClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AgentPoolClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates the agent pool in the Hybrid AKS provisioned cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Parameter for the name of the provisioned cluster
//   - agentPoolName - Parameter for the name of the agent pool in the provisioned cluster
//   - options - AgentPoolClientBeginCreateOrUpdateOptions contains the optional parameters for the AgentPoolClient.BeginCreateOrUpdate
//     method.
func (client *AgentPoolClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, agentPool AgentPool, options *AgentPoolClientBeginCreateOrUpdateOptions) (*runtime.Poller[AgentPoolClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, resourceName, agentPoolName, agentPool, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AgentPoolClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AgentPoolClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates the agent pool in the Hybrid AKS provisioned cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
func (client *AgentPoolClient) createOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, agentPool AgentPool, options *AgentPoolClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, agentPoolName, agentPool, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *AgentPoolClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, agentPool AgentPool, options *AgentPoolClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridContainerService/provisionedClusters/{resourceName}/agentPools/{agentPoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if agentPoolName == "" {
		return nil, errors.New("parameter agentPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{agentPoolName}", url.PathEscape(agentPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, agentPool); err != nil {
	return nil, err
}
	return req, nil
}

// Delete - Deletes the agent pool in the Hybrid AKS provisioned cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Parameter for the name of the provisioned cluster
//   - agentPoolName - Parameter for the name of the agent pool in the provisioned cluster
//   - options - AgentPoolClientDeleteOptions contains the optional parameters for the AgentPoolClient.Delete method.
func (client *AgentPoolClient) Delete(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, options *AgentPoolClientDeleteOptions) (AgentPoolClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, agentPoolName, options)
	if err != nil {
		return AgentPoolClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AgentPoolClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return AgentPoolClientDeleteResponse{}, err
	}
	return AgentPoolClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AgentPoolClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, options *AgentPoolClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridContainerService/provisionedClusters/{resourceName}/agentPools/{agentPoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if agentPoolName == "" {
		return nil, errors.New("parameter agentPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{agentPoolName}", url.PathEscape(agentPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the agent pool in the Hybrid AKS provisioned cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Parameter for the name of the provisioned cluster
//   - agentPoolName - Parameter for the name of the agent pool in the provisioned cluster
//   - options - AgentPoolClientGetOptions contains the optional parameters for the AgentPoolClient.Get method.
func (client *AgentPoolClient) Get(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, options *AgentPoolClientGetOptions) (AgentPoolClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, agentPoolName, options)
	if err != nil {
		return AgentPoolClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AgentPoolClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AgentPoolClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AgentPoolClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, options *AgentPoolClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridContainerService/provisionedClusters/{resourceName}/agentPools/{agentPoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if agentPoolName == "" {
		return nil, errors.New("parameter agentPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{agentPoolName}", url.PathEscape(agentPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AgentPoolClient) getHandleResponse(resp *http.Response) (AgentPoolClientGetResponse, error) {
	result := AgentPoolClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AgentPool); err != nil {
		return AgentPoolClientGetResponse{}, err
	}
	return result, nil
}

// ListByProvisionedCluster - Gets the agent pools in the Hybrid AKS provisioned cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Parameter for the name of the provisioned cluster
//   - options - AgentPoolClientListByProvisionedClusterOptions contains the optional parameters for the AgentPoolClient.ListByProvisionedCluster
//     method.
func (client *AgentPoolClient) ListByProvisionedCluster(ctx context.Context, resourceGroupName string, resourceName string, options *AgentPoolClientListByProvisionedClusterOptions) (AgentPoolClientListByProvisionedClusterResponse, error) {
	var err error
	req, err := client.listByProvisionedClusterCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return AgentPoolClientListByProvisionedClusterResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AgentPoolClientListByProvisionedClusterResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AgentPoolClientListByProvisionedClusterResponse{}, err
	}
	resp, err := client.listByProvisionedClusterHandleResponse(httpResp)
	return resp, err
}

// listByProvisionedClusterCreateRequest creates the ListByProvisionedCluster request.
func (client *AgentPoolClient) listByProvisionedClusterCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *AgentPoolClientListByProvisionedClusterOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridContainerService/provisionedClusters/{resourceName}/agentPools"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByProvisionedClusterHandleResponse handles the ListByProvisionedCluster response.
func (client *AgentPoolClient) listByProvisionedClusterHandleResponse(resp *http.Response) (AgentPoolClientListByProvisionedClusterResponse, error) {
	result := AgentPoolClientListByProvisionedClusterResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AgentPoolListResult); err != nil {
		return AgentPoolClientListByProvisionedClusterResponse{}, err
	}
	return result, nil
}

// Update - Updates the agent pool in the Hybrid AKS provisioned cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Parameter for the name of the provisioned cluster
//   - agentPoolName - Parameter for the name of the agent pool in the provisioned cluster
//   - options - AgentPoolClientUpdateOptions contains the optional parameters for the AgentPoolClient.Update method.
func (client *AgentPoolClient) Update(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, agentPool AgentPool, options *AgentPoolClientUpdateOptions) (AgentPoolClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, agentPoolName, agentPool, options)
	if err != nil {
		return AgentPoolClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AgentPoolClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return AgentPoolClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *AgentPoolClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, agentPoolName string, agentPool AgentPool, options *AgentPoolClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridContainerService/provisionedClusters/{resourceName}/agentPools/{agentPoolName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if agentPoolName == "" {
		return nil, errors.New("parameter agentPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{agentPoolName}", url.PathEscape(agentPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, agentPool); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *AgentPoolClient) updateHandleResponse(resp *http.Response) (AgentPoolClientUpdateResponse, error) {
	result := AgentPoolClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AgentPool); err != nil {
		return AgentPoolClientUpdateResponse{}, err
	}
	return result, nil
}

