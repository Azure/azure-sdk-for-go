//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armservicefabric

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

// ApplicationsClient contains the methods for the Applications group.
// Don't use this type directly, use NewApplicationsClient() instead.
type ApplicationsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewApplicationsClient creates a new instance of ApplicationsClient with the specified values.
// subscriptionID - The customer subscription identifier.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewApplicationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ApplicationsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &ApplicationsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update a Service Fabric application resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// clusterName - The name of the cluster resource.
// applicationName - The name of the application resource.
// parameters - The application resource.
// options - ApplicationsClientBeginCreateOrUpdateOptions contains the optional parameters for the ApplicationsClient.BeginCreateOrUpdate
// method.
func (client *ApplicationsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, parameters ApplicationResource, options *ApplicationsClientBeginCreateOrUpdateOptions) (*armruntime.Poller[ApplicationsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, applicationName, parameters, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[ApplicationsClientCreateOrUpdateResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[ApplicationsClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Create or update a Service Fabric application resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ApplicationsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, parameters ApplicationResource, options *ApplicationsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, applicationName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ApplicationsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, parameters ApplicationResource, options *ApplicationsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/clusters/{clusterName}/applications/{applicationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Delete a Service Fabric application resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// clusterName - The name of the cluster resource.
// applicationName - The name of the application resource.
// options - ApplicationsClientBeginDeleteOptions contains the optional parameters for the ApplicationsClient.BeginDelete
// method.
func (client *ApplicationsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, options *ApplicationsClientBeginDeleteOptions) (*armruntime.Poller[ApplicationsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, applicationName, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[ApplicationsClientDeleteResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[ApplicationsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Delete a Service Fabric application resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ApplicationsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, options *ApplicationsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, applicationName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ApplicationsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, options *ApplicationsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/clusters/{clusterName}/applications/{applicationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Get a Service Fabric application resource created or in the process of being created in the Service Fabric cluster
// resource.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// clusterName - The name of the cluster resource.
// applicationName - The name of the application resource.
// options - ApplicationsClientGetOptions contains the optional parameters for the ApplicationsClient.Get method.
func (client *ApplicationsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, options *ApplicationsClientGetOptions) (ApplicationsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, applicationName, options)
	if err != nil {
		return ApplicationsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ApplicationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ApplicationsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ApplicationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, options *ApplicationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/clusters/{clusterName}/applications/{applicationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ApplicationsClient) getHandleResponse(resp *http.Response) (ApplicationsClientGetResponse, error) {
	result := ApplicationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationResource); err != nil {
		return ApplicationsClientGetResponse{}, err
	}
	return result, nil
}

// List - Gets all application resources created or in the process of being created in the Service Fabric cluster resource.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// clusterName - The name of the cluster resource.
// options - ApplicationsClientListOptions contains the optional parameters for the ApplicationsClient.List method.
func (client *ApplicationsClient) List(ctx context.Context, resourceGroupName string, clusterName string, options *ApplicationsClientListOptions) (ApplicationsClientListResponse, error) {
	req, err := client.listCreateRequest(ctx, resourceGroupName, clusterName, options)
	if err != nil {
		return ApplicationsClientListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ApplicationsClientListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ApplicationsClientListResponse{}, runtime.NewResponseError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *ApplicationsClient) listCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, options *ApplicationsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/clusters/{clusterName}/applications"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ApplicationsClient) listHandleResponse(resp *http.Response) (ApplicationsClientListResponse, error) {
	result := ApplicationsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplicationResourceList); err != nil {
		return ApplicationsClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update a Service Fabric application resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// clusterName - The name of the cluster resource.
// applicationName - The name of the application resource.
// parameters - The application resource for patch operations.
// options - ApplicationsClientBeginUpdateOptions contains the optional parameters for the ApplicationsClient.BeginUpdate
// method.
func (client *ApplicationsClient) BeginUpdate(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, parameters ApplicationResourceUpdate, options *ApplicationsClientBeginUpdateOptions) (*armruntime.Poller[ApplicationsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, clusterName, applicationName, parameters, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[ApplicationsClientUpdateResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[ApplicationsClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Update a Service Fabric application resource with the specified name.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ApplicationsClient) update(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, parameters ApplicationResourceUpdate, options *ApplicationsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, clusterName, applicationName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *ApplicationsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, applicationName string, parameters ApplicationResourceUpdate, options *ApplicationsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/clusters/{clusterName}/applications/{applicationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if applicationName == "" {
		return nil, errors.New("parameter applicationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{applicationName}", url.PathEscape(applicationName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}
