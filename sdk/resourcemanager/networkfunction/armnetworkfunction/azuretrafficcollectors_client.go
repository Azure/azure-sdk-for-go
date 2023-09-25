//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetworkfunction

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

// AzureTrafficCollectorsClient contains the methods for the AzureTrafficCollectors group.
// Don't use this type directly, use NewAzureTrafficCollectorsClient() instead.
type AzureTrafficCollectorsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewAzureTrafficCollectorsClient creates a new instance of AzureTrafficCollectorsClient with the specified values.
//   - subscriptionID - Azure Subscription ID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAzureTrafficCollectorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AzureTrafficCollectorsClient, error) {
	cl, err := arm.NewClient(moduleName+".AzureTrafficCollectorsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AzureTrafficCollectorsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a Azure Traffic Collector resource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group.
//   - azureTrafficCollectorName - Azure Traffic Collector name
//   - parameters - The parameters to provide for the created Azure Traffic Collector.
//   - options - AzureTrafficCollectorsClientBeginCreateOrUpdateOptions contains the optional parameters for the AzureTrafficCollectorsClient.BeginCreateOrUpdate
//     method.
func (client *AzureTrafficCollectorsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, parameters AzureTrafficCollector, options *AzureTrafficCollectorsClientBeginCreateOrUpdateOptions) (*runtime.Poller[AzureTrafficCollectorsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, azureTrafficCollectorName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AzureTrafficCollectorsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AzureTrafficCollectorsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a Azure Traffic Collector resource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
func (client *AzureTrafficCollectorsClient) createOrUpdate(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, parameters AzureTrafficCollector, options *AzureTrafficCollectorsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, azureTrafficCollectorName, parameters, options)
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
func (client *AzureTrafficCollectorsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, parameters AzureTrafficCollector, options *AzureTrafficCollectorsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkFunction/azureTrafficCollectors/{azureTrafficCollectorName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if azureTrafficCollectorName == "" {
		return nil, errors.New("parameter azureTrafficCollectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureTrafficCollectorName}", url.PathEscape(azureTrafficCollectorName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes a specified Azure Traffic Collector resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group.
//   - azureTrafficCollectorName - Azure Traffic Collector name
//   - options - AzureTrafficCollectorsClientBeginDeleteOptions contains the optional parameters for the AzureTrafficCollectorsClient.BeginDelete
//     method.
func (client *AzureTrafficCollectorsClient) BeginDelete(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, options *AzureTrafficCollectorsClientBeginDeleteOptions) (*runtime.Poller[AzureTrafficCollectorsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, azureTrafficCollectorName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AzureTrafficCollectorsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AzureTrafficCollectorsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a specified Azure Traffic Collector resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
func (client *AzureTrafficCollectorsClient) deleteOperation(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, options *AzureTrafficCollectorsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, azureTrafficCollectorName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AzureTrafficCollectorsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, options *AzureTrafficCollectorsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkFunction/azureTrafficCollectors/{azureTrafficCollectorName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if azureTrafficCollectorName == "" {
		return nil, errors.New("parameter azureTrafficCollectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureTrafficCollectorName}", url.PathEscape(azureTrafficCollectorName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified Azure Traffic Collector in a specified resource group
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group.
//   - azureTrafficCollectorName - Azure Traffic Collector name
//   - options - AzureTrafficCollectorsClientGetOptions contains the optional parameters for the AzureTrafficCollectorsClient.Get
//     method.
func (client *AzureTrafficCollectorsClient) Get(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, options *AzureTrafficCollectorsClientGetOptions) (AzureTrafficCollectorsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, azureTrafficCollectorName, options)
	if err != nil {
		return AzureTrafficCollectorsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AzureTrafficCollectorsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AzureTrafficCollectorsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AzureTrafficCollectorsClient) getCreateRequest(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, options *AzureTrafficCollectorsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkFunction/azureTrafficCollectors/{azureTrafficCollectorName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if azureTrafficCollectorName == "" {
		return nil, errors.New("parameter azureTrafficCollectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureTrafficCollectorName}", url.PathEscape(azureTrafficCollectorName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AzureTrafficCollectorsClient) getHandleResponse(resp *http.Response) (AzureTrafficCollectorsClientGetResponse, error) {
	result := AzureTrafficCollectorsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureTrafficCollector); err != nil {
		return AzureTrafficCollectorsClientGetResponse{}, err
	}
	return result, nil
}

// UpdateTags - Updates the specified Azure Traffic Collector tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-11-01
//   - resourceGroupName - The name of the resource group.
//   - azureTrafficCollectorName - Azure Traffic Collector name
//   - parameters - Parameters supplied to update Azure Traffic Collector tags.
//   - options - AzureTrafficCollectorsClientUpdateTagsOptions contains the optional parameters for the AzureTrafficCollectorsClient.UpdateTags
//     method.
func (client *AzureTrafficCollectorsClient) UpdateTags(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, parameters TagsObject, options *AzureTrafficCollectorsClientUpdateTagsOptions) (AzureTrafficCollectorsClientUpdateTagsResponse, error) {
	var err error
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, azureTrafficCollectorName, parameters, options)
	if err != nil {
		return AzureTrafficCollectorsClientUpdateTagsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AzureTrafficCollectorsClientUpdateTagsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AzureTrafficCollectorsClientUpdateTagsResponse{}, err
	}
	resp, err := client.updateTagsHandleResponse(httpResp)
	return resp, err
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *AzureTrafficCollectorsClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, azureTrafficCollectorName string, parameters TagsObject, options *AzureTrafficCollectorsClientUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkFunction/azureTrafficCollectors/{azureTrafficCollectorName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if azureTrafficCollectorName == "" {
		return nil, errors.New("parameter azureTrafficCollectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{azureTrafficCollectorName}", url.PathEscape(azureTrafficCollectorName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client *AzureTrafficCollectorsClient) updateTagsHandleResponse(resp *http.Response) (AzureTrafficCollectorsClientUpdateTagsResponse, error) {
	result := AzureTrafficCollectorsClientUpdateTagsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AzureTrafficCollector); err != nil {
		return AzureTrafficCollectorsClientUpdateTagsResponse{}, err
	}
	return result, nil
}

