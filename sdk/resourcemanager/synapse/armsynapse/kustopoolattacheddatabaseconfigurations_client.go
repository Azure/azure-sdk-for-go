//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsynapse

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

// KustoPoolAttachedDatabaseConfigurationsClient contains the methods for the KustoPoolAttachedDatabaseConfigurations group.
// Don't use this type directly, use NewKustoPoolAttachedDatabaseConfigurationsClient() instead.
type KustoPoolAttachedDatabaseConfigurationsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewKustoPoolAttachedDatabaseConfigurationsClient creates a new instance of KustoPoolAttachedDatabaseConfigurationsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewKustoPoolAttachedDatabaseConfigurationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*KustoPoolAttachedDatabaseConfigurationsClient, error) {
	cl, err := arm.NewClient(moduleName+".KustoPoolAttachedDatabaseConfigurationsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &KustoPoolAttachedDatabaseConfigurationsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates an attached database configuration.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - attachedDatabaseConfigurationName - The name of the attached database configuration.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - parameters - The database parameters supplied to the CreateOrUpdate operation.
//   - options - KustoPoolAttachedDatabaseConfigurationsClientBeginCreateOrUpdateOptions contains the optional parameters for
//     the KustoPoolAttachedDatabaseConfigurationsClient.BeginCreateOrUpdate method.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) BeginCreateOrUpdate(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, parameters AttachedDatabaseConfiguration, options *KustoPoolAttachedDatabaseConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[KustoPoolAttachedDatabaseConfigurationsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, workspaceName, kustoPoolName, attachedDatabaseConfigurationName, resourceGroupName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[KustoPoolAttachedDatabaseConfigurationsClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[KustoPoolAttachedDatabaseConfigurationsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates an attached database configuration.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
func (client *KustoPoolAttachedDatabaseConfigurationsClient) createOrUpdate(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, parameters AttachedDatabaseConfiguration, options *KustoPoolAttachedDatabaseConfigurationsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, workspaceName, kustoPoolName, attachedDatabaseConfigurationName, resourceGroupName, parameters, options)
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
func (client *KustoPoolAttachedDatabaseConfigurationsClient) createOrUpdateCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, parameters AttachedDatabaseConfiguration, options *KustoPoolAttachedDatabaseConfigurationsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/attachedDatabaseConfigurations/{attachedDatabaseConfigurationName}"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if attachedDatabaseConfigurationName == "" {
		return nil, errors.New("parameter attachedDatabaseConfigurationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{attachedDatabaseConfigurationName}", url.PathEscape(attachedDatabaseConfigurationName))
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
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes the attached database configuration with the given name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - attachedDatabaseConfigurationName - The name of the attached database configuration.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - KustoPoolAttachedDatabaseConfigurationsClientBeginDeleteOptions contains the optional parameters for the KustoPoolAttachedDatabaseConfigurationsClient.BeginDelete
//     method.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) BeginDelete(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, options *KustoPoolAttachedDatabaseConfigurationsClientBeginDeleteOptions) (*runtime.Poller[KustoPoolAttachedDatabaseConfigurationsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, workspaceName, kustoPoolName, attachedDatabaseConfigurationName, resourceGroupName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[KustoPoolAttachedDatabaseConfigurationsClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[KustoPoolAttachedDatabaseConfigurationsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the attached database configuration with the given name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
func (client *KustoPoolAttachedDatabaseConfigurationsClient) deleteOperation(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, options *KustoPoolAttachedDatabaseConfigurationsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, workspaceName, kustoPoolName, attachedDatabaseConfigurationName, resourceGroupName, options)
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
func (client *KustoPoolAttachedDatabaseConfigurationsClient) deleteCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, options *KustoPoolAttachedDatabaseConfigurationsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/attachedDatabaseConfigurations/{attachedDatabaseConfigurationName}"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if attachedDatabaseConfigurationName == "" {
		return nil, errors.New("parameter attachedDatabaseConfigurationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{attachedDatabaseConfigurationName}", url.PathEscape(attachedDatabaseConfigurationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Returns an attached database configuration.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - attachedDatabaseConfigurationName - The name of the attached database configuration.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - KustoPoolAttachedDatabaseConfigurationsClientGetOptions contains the optional parameters for the KustoPoolAttachedDatabaseConfigurationsClient.Get
//     method.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) Get(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, options *KustoPoolAttachedDatabaseConfigurationsClientGetOptions) (KustoPoolAttachedDatabaseConfigurationsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, workspaceName, kustoPoolName, attachedDatabaseConfigurationName, resourceGroupName, options)
	if err != nil {
		return KustoPoolAttachedDatabaseConfigurationsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return KustoPoolAttachedDatabaseConfigurationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return KustoPoolAttachedDatabaseConfigurationsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) getCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, attachedDatabaseConfigurationName string, resourceGroupName string, options *KustoPoolAttachedDatabaseConfigurationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/attachedDatabaseConfigurations/{attachedDatabaseConfigurationName}"
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if kustoPoolName == "" {
		return nil, errors.New("parameter kustoPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{kustoPoolName}", url.PathEscape(kustoPoolName))
	if attachedDatabaseConfigurationName == "" {
		return nil, errors.New("parameter attachedDatabaseConfigurationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{attachedDatabaseConfigurationName}", url.PathEscape(attachedDatabaseConfigurationName))
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
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) getHandleResponse(resp *http.Response) (KustoPoolAttachedDatabaseConfigurationsClientGetResponse, error) {
	result := KustoPoolAttachedDatabaseConfigurationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AttachedDatabaseConfiguration); err != nil {
		return KustoPoolAttachedDatabaseConfigurationsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByKustoPoolPager - Returns the list of attached database configurations of the given Kusto Pool.
//
// Generated from API version 2021-06-01-preview
//   - workspaceName - The name of the workspace.
//   - kustoPoolName - The name of the Kusto pool.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolOptions contains the optional parameters for the
//     KustoPoolAttachedDatabaseConfigurationsClient.NewListByKustoPoolPager method.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) NewListByKustoPoolPager(workspaceName string, kustoPoolName string, resourceGroupName string, options *KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolOptions) (*runtime.Pager[KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse]) {
	return runtime.NewPager(runtime.PagingHandler[KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse]{
		More: func(page KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse) (KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse, error) {
			req, err := client.listByKustoPoolCreateRequest(ctx, workspaceName, kustoPoolName, resourceGroupName, options)
			if err != nil {
				return KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByKustoPoolHandleResponse(resp)
		},
	})
}

// listByKustoPoolCreateRequest creates the ListByKustoPool request.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) listByKustoPoolCreateRequest(ctx context.Context, workspaceName string, kustoPoolName string, resourceGroupName string, options *KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/kustoPools/{kustoPoolName}/attachedDatabaseConfigurations"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByKustoPoolHandleResponse handles the ListByKustoPool response.
func (client *KustoPoolAttachedDatabaseConfigurationsClient) listByKustoPoolHandleResponse(resp *http.Response) (KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse, error) {
	result := KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AttachedDatabaseConfigurationListResult); err != nil {
		return KustoPoolAttachedDatabaseConfigurationsClientListByKustoPoolResponse{}, err
	}
	return result, nil
}

