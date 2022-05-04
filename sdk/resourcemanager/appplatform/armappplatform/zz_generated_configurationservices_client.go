//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappplatform

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

// ConfigurationServicesClient contains the methods for the ConfigurationServices group.
// Don't use this type directly, use NewConfigurationServicesClient() instead.
type ConfigurationServicesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewConfigurationServicesClient creates a new instance of ConfigurationServicesClient with the specified values.
// subscriptionID - Gets subscription ID which uniquely identify the Microsoft Azure subscription. The subscription ID forms
// part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewConfigurationServicesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConfigurationServicesClient, error) {
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
	client := &ConfigurationServicesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create the default Application Configuration Service or update the existing Application Configuration
// Service.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// configurationServiceName - The name of Application Configuration Service.
// configurationServiceResource - Parameters for the update operation
// options - ConfigurationServicesClientBeginCreateOrUpdateOptions contains the optional parameters for the ConfigurationServicesClient.BeginCreateOrUpdate
// method.
func (client *ConfigurationServicesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, configurationServiceResource ConfigurationServiceResource, options *ConfigurationServicesClientBeginCreateOrUpdateOptions) (*armruntime.Poller[ConfigurationServicesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, serviceName, configurationServiceName, configurationServiceResource, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[ConfigurationServicesClientCreateOrUpdateResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[ConfigurationServicesClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Create the default Application Configuration Service or update the existing Application Configuration
// Service.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ConfigurationServicesClient) createOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, configurationServiceResource ConfigurationServiceResource, options *ConfigurationServicesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, configurationServiceName, configurationServiceResource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ConfigurationServicesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, configurationServiceResource ConfigurationServiceResource, options *ConfigurationServicesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/configurationServices/{configurationServiceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if configurationServiceName == "" {
		return nil, errors.New("parameter configurationServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationServiceName}", url.PathEscape(configurationServiceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, configurationServiceResource)
}

// BeginDelete - Disable the default Application Configuration Service.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// configurationServiceName - The name of Application Configuration Service.
// options - ConfigurationServicesClientBeginDeleteOptions contains the optional parameters for the ConfigurationServicesClient.BeginDelete
// method.
func (client *ConfigurationServicesClient) BeginDelete(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, options *ConfigurationServicesClientBeginDeleteOptions) (*armruntime.Poller[ConfigurationServicesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, serviceName, configurationServiceName, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[ConfigurationServicesClientDeleteResponse]{
			FinalStateVia: armruntime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[ConfigurationServicesClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Disable the default Application Configuration Service.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ConfigurationServicesClient) deleteOperation(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, options *ConfigurationServicesClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, configurationServiceName, options)
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
func (client *ConfigurationServicesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, options *ConfigurationServicesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/configurationServices/{configurationServiceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if configurationServiceName == "" {
		return nil, errors.New("parameter configurationServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationServiceName}", url.PathEscape(configurationServiceName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Get the Application Configuration Service and its properties.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// configurationServiceName - The name of Application Configuration Service.
// options - ConfigurationServicesClientGetOptions contains the optional parameters for the ConfigurationServicesClient.Get
// method.
func (client *ConfigurationServicesClient) Get(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, options *ConfigurationServicesClientGetOptions) (ConfigurationServicesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, configurationServiceName, options)
	if err != nil {
		return ConfigurationServicesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ConfigurationServicesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ConfigurationServicesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ConfigurationServicesClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, options *ConfigurationServicesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/configurationServices/{configurationServiceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if configurationServiceName == "" {
		return nil, errors.New("parameter configurationServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationServiceName}", url.PathEscape(configurationServiceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ConfigurationServicesClient) getHandleResponse(resp *http.Response) (ConfigurationServicesClientGetResponse, error) {
	result := ConfigurationServicesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConfigurationServiceResource); err != nil {
		return ConfigurationServicesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Handles requests to list all resources in a Service.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// options - ConfigurationServicesClientListOptions contains the optional parameters for the ConfigurationServicesClient.List
// method.
func (client *ConfigurationServicesClient) NewListPager(resourceGroupName string, serviceName string, options *ConfigurationServicesClientListOptions) *runtime.Pager[ConfigurationServicesClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[ConfigurationServicesClientListResponse]{
		More: func(page ConfigurationServicesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ConfigurationServicesClientListResponse) (ConfigurationServicesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, serviceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ConfigurationServicesClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ConfigurationServicesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ConfigurationServicesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ConfigurationServicesClient) listCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, options *ConfigurationServicesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/configurationServices"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ConfigurationServicesClient) listHandleResponse(resp *http.Response) (ConfigurationServicesClientListResponse, error) {
	result := ConfigurationServicesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConfigurationServiceResourceCollection); err != nil {
		return ConfigurationServicesClientListResponse{}, err
	}
	return result, nil
}

// BeginValidate - Check if the Application Configuration Service settings are valid.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// configurationServiceName - The name of Application Configuration Service.
// settings - Application Configuration Service settings to be validated
// options - ConfigurationServicesClientBeginValidateOptions contains the optional parameters for the ConfigurationServicesClient.BeginValidate
// method.
func (client *ConfigurationServicesClient) BeginValidate(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, settings ConfigurationServiceSettings, options *ConfigurationServicesClientBeginValidateOptions) (*armruntime.Poller[ConfigurationServicesClientValidateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.validate(ctx, resourceGroupName, serviceName, configurationServiceName, settings, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller(resp, client.pl, &armruntime.NewPollerOptions[ConfigurationServicesClientValidateResponse]{
			FinalStateVia: armruntime.FinalStateViaLocation,
		})
	} else {
		return armruntime.NewPollerFromResumeToken[ConfigurationServicesClientValidateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Validate - Check if the Application Configuration Service settings are valid.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *ConfigurationServicesClient) validate(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, settings ConfigurationServiceSettings, options *ConfigurationServicesClientBeginValidateOptions) (*http.Response, error) {
	req, err := client.validateCreateRequest(ctx, resourceGroupName, serviceName, configurationServiceName, settings, options)
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

// validateCreateRequest creates the Validate request.
func (client *ConfigurationServicesClient) validateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, configurationServiceName string, settings ConfigurationServiceSettings, options *ConfigurationServicesClientBeginValidateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/configurationServices/{configurationServiceName}/validate"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if configurationServiceName == "" {
		return nil, errors.New("parameter configurationServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationServiceName}", url.PathEscape(configurationServiceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, settings)
}
