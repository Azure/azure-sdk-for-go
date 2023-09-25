//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhealthcareapis

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

// IotConnectorFhirDestinationClient contains the methods for the IotConnectorFhirDestination group.
// Don't use this type directly, use NewIotConnectorFhirDestinationClient() instead.
type IotConnectorFhirDestinationClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewIotConnectorFhirDestinationClient creates a new instance of IotConnectorFhirDestinationClient with the specified values.
//   - subscriptionID - The subscription identifier.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewIotConnectorFhirDestinationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*IotConnectorFhirDestinationClient, error) {
	cl, err := arm.NewClient(moduleName+".IotConnectorFhirDestinationClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &IotConnectorFhirDestinationClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates an IoT Connector FHIR destination resource with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01
//   - resourceGroupName - The name of the resource group that contains the service instance.
//   - workspaceName - The name of workspace resource.
//   - iotConnectorName - The name of IoT Connector resource.
//   - fhirDestinationName - The name of IoT Connector FHIR destination resource.
//   - iotFhirDestination - The parameters for creating or updating an IoT Connector FHIR destination resource.
//   - options - IotConnectorFhirDestinationClientBeginCreateOrUpdateOptions contains the optional parameters for the IotConnectorFhirDestinationClient.BeginCreateOrUpdate
//     method.
func (client *IotConnectorFhirDestinationClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, iotFhirDestination IotFhirDestination, options *IotConnectorFhirDestinationClientBeginCreateOrUpdateOptions) (*runtime.Poller[IotConnectorFhirDestinationClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, iotFhirDestination, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[IotConnectorFhirDestinationClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[IotConnectorFhirDestinationClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates an IoT Connector FHIR destination resource with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01
func (client *IotConnectorFhirDestinationClient) createOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, iotFhirDestination IotFhirDestination, options *IotConnectorFhirDestinationClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, iotFhirDestination, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
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
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, iotFhirDestination); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes an IoT Connector FHIR destination.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01
//   - resourceGroupName - The name of the resource group that contains the service instance.
//   - workspaceName - The name of workspace resource.
//   - iotConnectorName - The name of IoT Connector resource.
//   - fhirDestinationName - The name of IoT Connector FHIR destination resource.
//   - options - IotConnectorFhirDestinationClientBeginDeleteOptions contains the optional parameters for the IotConnectorFhirDestinationClient.BeginDelete
//     method.
func (client *IotConnectorFhirDestinationClient) BeginDelete(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientBeginDeleteOptions) (*runtime.Poller[IotConnectorFhirDestinationClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[IotConnectorFhirDestinationClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[IotConnectorFhirDestinationClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes an IoT Connector FHIR destination.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01
func (client *IotConnectorFhirDestinationClient) deleteOperation(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, options)
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
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the properties of the specified Iot Connector FHIR destination.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01
//   - resourceGroupName - The name of the resource group that contains the service instance.
//   - workspaceName - The name of workspace resource.
//   - iotConnectorName - The name of IoT Connector resource.
//   - fhirDestinationName - The name of IoT Connector FHIR destination resource.
//   - options - IotConnectorFhirDestinationClientGetOptions contains the optional parameters for the IotConnectorFhirDestinationClient.Get
//     method.
func (client *IotConnectorFhirDestinationClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, iotConnectorName string, fhirDestinationName string, options *IotConnectorFhirDestinationClientGetOptions) (IotConnectorFhirDestinationClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, iotConnectorName, fhirDestinationName, options)
	if err != nil {
		return IotConnectorFhirDestinationClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return IotConnectorFhirDestinationClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return IotConnectorFhirDestinationClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
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
func (client *IotConnectorFhirDestinationClient) getHandleResponse(resp *http.Response) (IotConnectorFhirDestinationClientGetResponse, error) {
	result := IotConnectorFhirDestinationClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.IotFhirDestination); err != nil {
		return IotConnectorFhirDestinationClientGetResponse{}, err
	}
	return result, nil
}

