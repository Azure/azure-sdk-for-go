//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armautomanage

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

// ConfigurationProfileHCRPAssignmentsClient contains the methods for the ConfigurationProfileHCRPAssignments group.
// Don't use this type directly, use NewConfigurationProfileHCRPAssignmentsClient() instead.
type ConfigurationProfileHCRPAssignmentsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewConfigurationProfileHCRPAssignmentsClient creates a new instance of ConfigurationProfileHCRPAssignmentsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewConfigurationProfileHCRPAssignmentsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConfigurationProfileHCRPAssignmentsClient, error) {
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
	client := &ConfigurationProfileHCRPAssignmentsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Creates an association between a ARC machine and Automanage configuration profile
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-04
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the Arc machine.
// configurationProfileAssignmentName - Name of the configuration profile assignment. Only default is supported.
// parameters - Parameters supplied to the create or update configuration profile assignment.
// options - ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateOptions contains the optional parameters for the ConfigurationProfileHCRPAssignmentsClient.CreateOrUpdate
// method.
func (client *ConfigurationProfileHCRPAssignmentsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, machineName string, configurationProfileAssignmentName string, parameters ConfigurationProfileAssignment, options *ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateOptions) (ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, machineName, configurationProfileAssignmentName, parameters, options)
	if err != nil {
		return ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ConfigurationProfileHCRPAssignmentsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, machineName string, configurationProfileAssignmentName string, parameters ConfigurationProfileAssignment, options *ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/providers/Microsoft.Automanage/configurationProfileAssignments/{configurationProfileAssignmentName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if configurationProfileAssignmentName == "" {
		return nil, errors.New("parameter configurationProfileAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationProfileAssignmentName}", url.PathEscape(configurationProfileAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ConfigurationProfileHCRPAssignmentsClient) createOrUpdateHandleResponse(resp *http.Response) (ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateResponse, error) {
	result := ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConfigurationProfileAssignment); err != nil {
		return ConfigurationProfileHCRPAssignmentsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a configuration profile assignment
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-04
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the Arc machine.
// configurationProfileAssignmentName - Name of the configuration profile assignment
// options - ConfigurationProfileHCRPAssignmentsClientDeleteOptions contains the optional parameters for the ConfigurationProfileHCRPAssignmentsClient.Delete
// method.
func (client *ConfigurationProfileHCRPAssignmentsClient) Delete(ctx context.Context, resourceGroupName string, machineName string, configurationProfileAssignmentName string, options *ConfigurationProfileHCRPAssignmentsClientDeleteOptions) (ConfigurationProfileHCRPAssignmentsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, machineName, configurationProfileAssignmentName, options)
	if err != nil {
		return ConfigurationProfileHCRPAssignmentsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ConfigurationProfileHCRPAssignmentsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return ConfigurationProfileHCRPAssignmentsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return ConfigurationProfileHCRPAssignmentsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ConfigurationProfileHCRPAssignmentsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, machineName string, configurationProfileAssignmentName string, options *ConfigurationProfileHCRPAssignmentsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/providers/Microsoft.Automanage/configurationProfileAssignments/{configurationProfileAssignmentName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if configurationProfileAssignmentName == "" {
		return nil, errors.New("parameter configurationProfileAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationProfileAssignmentName}", url.PathEscape(configurationProfileAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get information about a configuration profile assignment
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-05-04
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the Arc machine.
// configurationProfileAssignmentName - The configuration profile assignment name.
// options - ConfigurationProfileHCRPAssignmentsClientGetOptions contains the optional parameters for the ConfigurationProfileHCRPAssignmentsClient.Get
// method.
func (client *ConfigurationProfileHCRPAssignmentsClient) Get(ctx context.Context, resourceGroupName string, machineName string, configurationProfileAssignmentName string, options *ConfigurationProfileHCRPAssignmentsClientGetOptions) (ConfigurationProfileHCRPAssignmentsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, machineName, configurationProfileAssignmentName, options)
	if err != nil {
		return ConfigurationProfileHCRPAssignmentsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ConfigurationProfileHCRPAssignmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ConfigurationProfileHCRPAssignmentsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ConfigurationProfileHCRPAssignmentsClient) getCreateRequest(ctx context.Context, resourceGroupName string, machineName string, configurationProfileAssignmentName string, options *ConfigurationProfileHCRPAssignmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/providers/Microsoft.Automanage/configurationProfileAssignments/{configurationProfileAssignmentName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if configurationProfileAssignmentName == "" {
		return nil, errors.New("parameter configurationProfileAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationProfileAssignmentName}", url.PathEscape(configurationProfileAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ConfigurationProfileHCRPAssignmentsClient) getHandleResponse(resp *http.Response) (ConfigurationProfileHCRPAssignmentsClientGetResponse, error) {
	result := ConfigurationProfileHCRPAssignmentsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConfigurationProfileAssignment); err != nil {
		return ConfigurationProfileHCRPAssignmentsClientGetResponse{}, err
	}
	return result, nil
}
