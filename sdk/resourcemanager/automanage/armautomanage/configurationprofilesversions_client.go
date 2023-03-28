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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ConfigurationProfilesVersionsClient contains the methods for the ConfigurationProfilesVersions group.
// Don't use this type directly, use NewConfigurationProfilesVersionsClient() instead.
type ConfigurationProfilesVersionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewConfigurationProfilesVersionsClient creates a new instance of ConfigurationProfilesVersionsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewConfigurationProfilesVersionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ConfigurationProfilesVersionsClient, error) {
	cl, err := arm.NewClient(moduleName+".ConfigurationProfilesVersionsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ConfigurationProfilesVersionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates a configuration profile version
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-04
//   - configurationProfileName - Name of the configuration profile.
//   - versionName - The configuration profile version name.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - parameters - Parameters supplied to create or update configuration profile.
//   - options - ConfigurationProfilesVersionsClientCreateOrUpdateOptions contains the optional parameters for the ConfigurationProfilesVersionsClient.CreateOrUpdate
//     method.
func (client *ConfigurationProfilesVersionsClient) CreateOrUpdate(ctx context.Context, configurationProfileName string, versionName string, resourceGroupName string, parameters ConfigurationProfile, options *ConfigurationProfilesVersionsClientCreateOrUpdateOptions) (ConfigurationProfilesVersionsClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, configurationProfileName, versionName, resourceGroupName, parameters, options)
	if err != nil {
		return ConfigurationProfilesVersionsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ConfigurationProfilesVersionsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return ConfigurationProfilesVersionsClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ConfigurationProfilesVersionsClient) createOrUpdateCreateRequest(ctx context.Context, configurationProfileName string, versionName string, resourceGroupName string, parameters ConfigurationProfile, options *ConfigurationProfilesVersionsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automanage/configurationProfiles/{configurationProfileName}/versions/{versionName}"
	if configurationProfileName == "" {
		return nil, errors.New("parameter configurationProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationProfileName}", url.PathEscape(configurationProfileName))
	if versionName == "" {
		return nil, errors.New("parameter versionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{versionName}", url.PathEscape(versionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
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
func (client *ConfigurationProfilesVersionsClient) createOrUpdateHandleResponse(resp *http.Response) (ConfigurationProfilesVersionsClientCreateOrUpdateResponse, error) {
	result := ConfigurationProfilesVersionsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConfigurationProfile); err != nil {
		return ConfigurationProfilesVersionsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a configuration profile version
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-04
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - configurationProfileName - Name of the configuration profile
//   - versionName - The configuration profile version name.
//   - options - ConfigurationProfilesVersionsClientDeleteOptions contains the optional parameters for the ConfigurationProfilesVersionsClient.Delete
//     method.
func (client *ConfigurationProfilesVersionsClient) Delete(ctx context.Context, resourceGroupName string, configurationProfileName string, versionName string, options *ConfigurationProfilesVersionsClientDeleteOptions) (ConfigurationProfilesVersionsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, configurationProfileName, versionName, options)
	if err != nil {
		return ConfigurationProfilesVersionsClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ConfigurationProfilesVersionsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return ConfigurationProfilesVersionsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return ConfigurationProfilesVersionsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ConfigurationProfilesVersionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, configurationProfileName string, versionName string, options *ConfigurationProfilesVersionsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automanage/configurationProfiles/{configurationProfileName}/versions/{versionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if configurationProfileName == "" {
		return nil, errors.New("parameter configurationProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationProfileName}", url.PathEscape(configurationProfileName))
	if versionName == "" {
		return nil, errors.New("parameter versionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{versionName}", url.PathEscape(versionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get information about a configuration profile version
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-04
//   - configurationProfileName - The configuration profile name.
//   - versionName - The configuration profile version name.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - ConfigurationProfilesVersionsClientGetOptions contains the optional parameters for the ConfigurationProfilesVersionsClient.Get
//     method.
func (client *ConfigurationProfilesVersionsClient) Get(ctx context.Context, configurationProfileName string, versionName string, resourceGroupName string, options *ConfigurationProfilesVersionsClientGetOptions) (ConfigurationProfilesVersionsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, configurationProfileName, versionName, resourceGroupName, options)
	if err != nil {
		return ConfigurationProfilesVersionsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ConfigurationProfilesVersionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ConfigurationProfilesVersionsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ConfigurationProfilesVersionsClient) getCreateRequest(ctx context.Context, configurationProfileName string, versionName string, resourceGroupName string, options *ConfigurationProfilesVersionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automanage/configurationProfiles/{configurationProfileName}/versions/{versionName}"
	if configurationProfileName == "" {
		return nil, errors.New("parameter configurationProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationProfileName}", url.PathEscape(configurationProfileName))
	if versionName == "" {
		return nil, errors.New("parameter versionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{versionName}", url.PathEscape(versionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
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
func (client *ConfigurationProfilesVersionsClient) getHandleResponse(resp *http.Response) (ConfigurationProfilesVersionsClientGetResponse, error) {
	result := ConfigurationProfilesVersionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConfigurationProfile); err != nil {
		return ConfigurationProfilesVersionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListChildResourcesPager - Retrieve a list of configuration profile version for a configuration profile
//
// Generated from API version 2022-05-04
//   - configurationProfileName - Name of the configuration profile.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - ConfigurationProfilesVersionsClientListChildResourcesOptions contains the optional parameters for the ConfigurationProfilesVersionsClient.NewListChildResourcesPager
//     method.
func (client *ConfigurationProfilesVersionsClient) NewListChildResourcesPager(configurationProfileName string, resourceGroupName string, options *ConfigurationProfilesVersionsClientListChildResourcesOptions) *runtime.Pager[ConfigurationProfilesVersionsClientListChildResourcesResponse] {
	return runtime.NewPager(runtime.PagingHandler[ConfigurationProfilesVersionsClientListChildResourcesResponse]{
		More: func(page ConfigurationProfilesVersionsClientListChildResourcesResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ConfigurationProfilesVersionsClientListChildResourcesResponse) (ConfigurationProfilesVersionsClientListChildResourcesResponse, error) {
			req, err := client.listChildResourcesCreateRequest(ctx, configurationProfileName, resourceGroupName, options)
			if err != nil {
				return ConfigurationProfilesVersionsClientListChildResourcesResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ConfigurationProfilesVersionsClientListChildResourcesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ConfigurationProfilesVersionsClientListChildResourcesResponse{}, runtime.NewResponseError(resp)
			}
			return client.listChildResourcesHandleResponse(resp)
		},
	})
}

// listChildResourcesCreateRequest creates the ListChildResources request.
func (client *ConfigurationProfilesVersionsClient) listChildResourcesCreateRequest(ctx context.Context, configurationProfileName string, resourceGroupName string, options *ConfigurationProfilesVersionsClientListChildResourcesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automanage/configurationProfiles/{configurationProfileName}/versions"
	if configurationProfileName == "" {
		return nil, errors.New("parameter configurationProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{configurationProfileName}", url.PathEscape(configurationProfileName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listChildResourcesHandleResponse handles the ListChildResources response.
func (client *ConfigurationProfilesVersionsClient) listChildResourcesHandleResponse(resp *http.Response) (ConfigurationProfilesVersionsClientListChildResourcesResponse, error) {
	result := ConfigurationProfilesVersionsClientListChildResourcesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ConfigurationProfileList); err != nil {
		return ConfigurationProfilesVersionsClientListChildResourcesResponse{}, err
	}
	return result, nil
}
