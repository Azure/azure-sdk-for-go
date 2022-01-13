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

// IotConnectorFhirDestinationClient contains the methods for the IotConnectorFhirDestination group.
// Don't use this type directly, use NewIotConnectorFhirDestinationClient() instead.
type IotConnectorFhirDestinationClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewIotConnectorFhirDestinationClient creates a new instance of IotConnectorFhirDestinationClient with the specified values.
// subscriptionID - The subscription identifier.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewIotConnectorFhirDestinationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *IotConnectorFhirDestinationClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &IotConnectorFhirDestinationClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// BeginCreateOrUpdate - Creates or updates an IoT Connector FHIR destination resource with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the service instance.
// workspaceName - The name of workspace resource.
// iotConnectorName - The name of IoT Connector resource.
// fhirDestinationName - The name of IoT Connector FHIR destination resource.
// iotFhirDestination - The parameters for creating or updating an IoT Connector FHIR destination resource.
// options - IotConnectorFhirDestinationClientBeginCreateOrUpdateOptions contains the optional parameters for the IotConnectorFhirDestinationClient.BeginCreateOrUpdate
// method.
func (client *IotConnectorFhirDestinationClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, iotFhirDestination IotFhirDestination, options *IotConnectorFhirDestinationClientBeginCreateOrUpdateOptions) (IotConnectorFhirDestinationClientCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, iotFhirDestination, options)
	if err != nil {
		return IotConnectorFhirDestinationClientCreateOrUpdatePollerResponse{}, err
	}
	result := IotConnectorFhirDestinationClientCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("IotConnectorFhirDestinationClient.CreateOrUpdate", "", resp, client.pl)
	if err != nil {
		return IotConnectorFhirDestinationClientCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &IotConnectorFhirDestinationClientCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Creates or updates an IoT Connector FHIR destination resource with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *IotConnectorFhirDestinationClient) createOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, iotFhirDestination IotFhirDestination, options *IotConnectorFhirDestinationClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, iotFhirDestination, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *IotConnectorFhirDestinationClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, iotFhirDestination IotFhirDestination, options *IotConnectorFhirDestinationClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/iotconnectors/{iotConnectorName}/fhirdestinations/{fhirDestinationName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if iotConnectorName == "" {
		return nil, errors.New("parameter iotConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{iotConnectorName}", url.PathEscape(iotConnectorName))
	if fhirDestinationName == "" {
		return nil, errors.New("parameter fhirDestinationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{fhirDestinationName}", url.PathEscape(fhirDestinationName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, iotFhirDestination)
}

// BeginDelete - Deletes an IoT Connector FHIR destination.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the service instance.
// workspaceName - The name of workspace resource.
// iotConnectorName - The name of IoT Connector resource.
// fhirDestinationName - The name of IoT Connector FHIR destination resource.
// options - IotConnectorFhirDestinationClientBeginDeleteOptions contains the optional parameters for the IotConnectorFhirDestinationClient.BeginDelete
// method.
func (client *IotConnectorFhirDestinationClient) BeginDelete(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientBeginDeleteOptions) (IotConnectorFhirDestinationClientDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, options)
	if err != nil {
		return IotConnectorFhirDestinationClientDeletePollerResponse{}, err
	}
	result := IotConnectorFhirDestinationClientDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("IotConnectorFhirDestinationClient.Delete", "", resp, client.pl)
	if err != nil {
		return IotConnectorFhirDestinationClientDeletePollerResponse{}, err
	}
	result.Poller = &IotConnectorFhirDestinationClientDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Deletes an IoT Connector FHIR destination.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *IotConnectorFhirDestinationClient) deleteOperation(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *IotConnectorFhirDestinationClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/iotconnectors/{iotConnectorName}/fhirdestinations/{fhirDestinationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if iotConnectorName == "" {
		return nil, errors.New("parameter iotConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{iotConnectorName}", url.PathEscape(iotConnectorName))
	if fhirDestinationName == "" {
		return nil, errors.New("parameter fhirDestinationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{fhirDestinationName}", url.PathEscape(fhirDestinationName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets the properties of the specified Iot Connector FHIR destination.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the service instance.
// workspaceName - The name of workspace resource.
// iotConnectorName - The name of IoT Connector resource.
// fhirDestinationName - The name of IoT Connector FHIR destination resource.
// options - IotConnectorFhirDestinationClientGetOptions contains the optional parameters for the IotConnectorFhirDestinationClient.Get
// method.
func (client *IotConnectorFhirDestinationClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientGetOptions) (IotConnectorFhirDestinationClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, options)
	if err != nil {
		return IotConnectorFhirDestinationClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return IotConnectorFhirDestinationClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return IotConnectorFhirDestinationClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *IotConnectorFhirDestinationClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/iotconnectors/{iotConnectorName}/fhirdestinations/{fhirDestinationName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if iotConnectorName == "" {
		return nil, errors.New("parameter iotConnectorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{iotConnectorName}", url.PathEscape(iotConnectorName))
	if fhirDestinationName == "" {
		return nil, errors.New("parameter fhirDestinationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{fhirDestinationName}", url.PathEscape(fhirDestinationName))
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
func (client *IotConnectorFhirDestinationClient) getHandleResponse(resp *http.Response) (IotConnectorFhirDestinationClientGetResponse, error) {
	result := IotConnectorFhirDestinationClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.IotFhirDestination); err != nil {
		return IotConnectorFhirDestinationClientGetResponse{}, err
	}
	return result, nil
}
