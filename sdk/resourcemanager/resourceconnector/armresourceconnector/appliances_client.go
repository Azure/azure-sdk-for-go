//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armresourceconnector

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

// AppliancesClient contains the methods for the Appliances group.
// Don't use this type directly, use NewAppliancesClient() instead.
type AppliancesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewAppliancesClient creates a new instance of AppliancesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAppliancesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AppliancesClient, error) {
	cl, err := arm.NewClient(moduleName+".AppliancesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AppliancesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates an Appliance in the specified Subscription and Resource Group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Appliances name.
//   - parameters - Parameters supplied to create or update an Appliance.
//   - options - AppliancesClientBeginCreateOrUpdateOptions contains the optional parameters for the AppliancesClient.BeginCreateOrUpdate
//     method.
func (client *AppliancesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, parameters Appliance, options *AppliancesClientBeginCreateOrUpdateOptions) (*runtime.Poller[AppliancesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, resourceName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AppliancesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AppliancesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates an Appliance in the specified Subscription and Resource Group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
func (client *AppliancesClient) createOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, parameters Appliance, options *AppliancesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, parameters, options)
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
func (client *AppliancesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, parameters Appliance, options *AppliancesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances/{resourceName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes an Appliance with the specified Resource Name, Resource Group, and Subscription Id.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Appliances name.
//   - options - AppliancesClientBeginDeleteOptions contains the optional parameters for the AppliancesClient.BeginDelete method.
func (client *AppliancesClient) BeginDelete(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientBeginDeleteOptions) (*runtime.Poller[AppliancesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, resourceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AppliancesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AppliancesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes an Appliance with the specified Resource Name, Resource Group, and Subscription Id.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
func (client *AppliancesClient) deleteOperation(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AppliancesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances/{resourceName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the details of an Appliance with a specified resource group and name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Appliances name.
//   - options - AppliancesClientGetOptions contains the optional parameters for the AppliancesClient.Get method.
func (client *AppliancesClient) Get(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientGetOptions) (AppliancesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return AppliancesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AppliancesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AppliancesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AppliancesClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances/{resourceName}"
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
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AppliancesClient) getHandleResponse(resp *http.Response) (AppliancesClientGetResponse, error) {
	result := AppliancesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Appliance); err != nil {
		return AppliancesClientGetResponse{}, err
	}
	return result, nil
}

// GetTelemetryConfig - Gets the telemetry config.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - options - AppliancesClientGetTelemetryConfigOptions contains the optional parameters for the AppliancesClient.GetTelemetryConfig
//     method.
func (client *AppliancesClient) GetTelemetryConfig(ctx context.Context, options *AppliancesClientGetTelemetryConfigOptions) (AppliancesClientGetTelemetryConfigResponse, error) {
	var err error
	req, err := client.getTelemetryConfigCreateRequest(ctx, options)
	if err != nil {
		return AppliancesClientGetTelemetryConfigResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AppliancesClientGetTelemetryConfigResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AppliancesClientGetTelemetryConfigResponse{}, err
	}
	resp, err := client.getTelemetryConfigHandleResponse(httpResp)
	return resp, err
}

// getTelemetryConfigCreateRequest creates the GetTelemetryConfig request.
func (client *AppliancesClient) getTelemetryConfigCreateRequest(ctx context.Context, options *AppliancesClientGetTelemetryConfigOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ResourceConnector/telemetryconfig"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getTelemetryConfigHandleResponse handles the GetTelemetryConfig response.
func (client *AppliancesClient) getTelemetryConfigHandleResponse(resp *http.Response) (AppliancesClientGetTelemetryConfigResponse, error) {
	result := AppliancesClientGetTelemetryConfigResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplianceGetTelemetryConfigResult); err != nil {
		return AppliancesClientGetTelemetryConfigResponse{}, err
	}
	return result, nil
}

// GetUpgradeGraph - Gets the upgrade graph of an Appliance with a specified resource group and name and specific release
// train.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Appliances name.
//   - upgradeGraph - Upgrade graph version, ex - stable
//   - options - AppliancesClientGetUpgradeGraphOptions contains the optional parameters for the AppliancesClient.GetUpgradeGraph
//     method.
func (client *AppliancesClient) GetUpgradeGraph(ctx context.Context, resourceGroupName string, resourceName string, upgradeGraph string, options *AppliancesClientGetUpgradeGraphOptions) (AppliancesClientGetUpgradeGraphResponse, error) {
	var err error
	req, err := client.getUpgradeGraphCreateRequest(ctx, resourceGroupName, resourceName, upgradeGraph, options)
	if err != nil {
		return AppliancesClientGetUpgradeGraphResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AppliancesClientGetUpgradeGraphResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AppliancesClientGetUpgradeGraphResponse{}, err
	}
	resp, err := client.getUpgradeGraphHandleResponse(httpResp)
	return resp, err
}

// getUpgradeGraphCreateRequest creates the GetUpgradeGraph request.
func (client *AppliancesClient) getUpgradeGraphCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, upgradeGraph string, options *AppliancesClientGetUpgradeGraphOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances/{resourceName}/upgradeGraphs/{upgradeGraph}"
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
	if upgradeGraph == "" {
		return nil, errors.New("parameter upgradeGraph cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{upgradeGraph}", url.PathEscape(upgradeGraph))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getUpgradeGraphHandleResponse handles the GetUpgradeGraph response.
func (client *AppliancesClient) getUpgradeGraphHandleResponse(resp *http.Response) (AppliancesClientGetUpgradeGraphResponse, error) {
	result := AppliancesClientGetUpgradeGraphResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UpgradeGraph); err != nil {
		return AppliancesClientGetUpgradeGraphResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Gets a list of Appliances in the specified subscription and resource group. The operation
// returns properties of each Appliance.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - AppliancesClientListByResourceGroupOptions contains the optional parameters for the AppliancesClient.NewListByResourceGroupPager
//     method.
func (client *AppliancesClient) NewListByResourceGroupPager(resourceGroupName string, options *AppliancesClientListByResourceGroupOptions) (*runtime.Pager[AppliancesClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AppliancesClientListByResourceGroupResponse]{
		More: func(page AppliancesClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AppliancesClientListByResourceGroupResponse) (AppliancesClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AppliancesClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AppliancesClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AppliancesClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *AppliancesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *AppliancesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances"
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
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *AppliancesClient) listByResourceGroupHandleResponse(resp *http.Response) (AppliancesClientListByResourceGroupResponse, error) {
	result := AppliancesClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplianceListResult); err != nil {
		return AppliancesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Gets a list of Appliances in the specified subscription. The operation returns properties
// of each Appliance
//
// Generated from API version 2022-10-27
//   - options - AppliancesClientListBySubscriptionOptions contains the optional parameters for the AppliancesClient.NewListBySubscriptionPager
//     method.
func (client *AppliancesClient) NewListBySubscriptionPager(options *AppliancesClientListBySubscriptionOptions) (*runtime.Pager[AppliancesClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AppliancesClientListBySubscriptionResponse]{
		More: func(page AppliancesClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AppliancesClientListBySubscriptionResponse) (AppliancesClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AppliancesClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AppliancesClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AppliancesClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *AppliancesClient) listBySubscriptionCreateRequest(ctx context.Context, options *AppliancesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ResourceConnector/appliances"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *AppliancesClient) listBySubscriptionHandleResponse(resp *http.Response) (AppliancesClientListBySubscriptionResponse, error) {
	result := AppliancesClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplianceListResult); err != nil {
		return AppliancesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// ListClusterUserCredential - Returns the cluster user credentials for the dedicated appliance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Appliances name.
//   - options - AppliancesClientListClusterUserCredentialOptions contains the optional parameters for the AppliancesClient.ListClusterUserCredential
//     method.
func (client *AppliancesClient) ListClusterUserCredential(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientListClusterUserCredentialOptions) (AppliancesClientListClusterUserCredentialResponse, error) {
	var err error
	req, err := client.listClusterUserCredentialCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return AppliancesClientListClusterUserCredentialResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AppliancesClientListClusterUserCredentialResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AppliancesClientListClusterUserCredentialResponse{}, err
	}
	resp, err := client.listClusterUserCredentialHandleResponse(httpResp)
	return resp, err
}

// listClusterUserCredentialCreateRequest creates the ListClusterUserCredential request.
func (client *AppliancesClient) listClusterUserCredentialCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientListClusterUserCredentialOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances/{resourceName}/listClusterUserCredential"
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
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listClusterUserCredentialHandleResponse handles the ListClusterUserCredential response.
func (client *AppliancesClient) listClusterUserCredentialHandleResponse(resp *http.Response) (AppliancesClientListClusterUserCredentialResponse, error) {
	result := AppliancesClientListClusterUserCredentialResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplianceListCredentialResults); err != nil {
		return AppliancesClientListClusterUserCredentialResponse{}, err
	}
	return result, nil
}

// ListKeys - Returns the cluster customer credentials for the dedicated appliance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Appliances name.
//   - options - AppliancesClientListKeysOptions contains the optional parameters for the AppliancesClient.ListKeys method.
func (client *AppliancesClient) ListKeys(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientListKeysOptions) (AppliancesClientListKeysResponse, error) {
	var err error
	req, err := client.listKeysCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return AppliancesClientListKeysResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AppliancesClientListKeysResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AppliancesClientListKeysResponse{}, err
	}
	resp, err := client.listKeysHandleResponse(httpResp)
	return resp, err
}

// listKeysCreateRequest creates the ListKeys request.
func (client *AppliancesClient) listKeysCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *AppliancesClientListKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances/{resourceName}/listkeys"
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
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	if options != nil && options.ArtifactType != nil {
		reqQP.Set("artifactType", *options.ArtifactType)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listKeysHandleResponse handles the ListKeys response.
func (client *AppliancesClient) listKeysHandleResponse(resp *http.Response) (AppliancesClientListKeysResponse, error) {
	result := AppliancesClientListKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplianceListKeysResults); err != nil {
		return AppliancesClientListKeysResponse{}, err
	}
	return result, nil
}

// NewListOperationsPager - Lists all available Appliances operations.
//
// Generated from API version 2022-10-27
//   - options - AppliancesClientListOperationsOptions contains the optional parameters for the AppliancesClient.NewListOperationsPager
//     method.
func (client *AppliancesClient) NewListOperationsPager(options *AppliancesClientListOperationsOptions) (*runtime.Pager[AppliancesClientListOperationsResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AppliancesClientListOperationsResponse]{
		More: func(page AppliancesClientListOperationsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AppliancesClientListOperationsResponse) (AppliancesClientListOperationsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listOperationsCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AppliancesClientListOperationsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AppliancesClientListOperationsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AppliancesClientListOperationsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listOperationsHandleResponse(resp)
		},
	})
}

// listOperationsCreateRequest creates the ListOperations request.
func (client *AppliancesClient) listOperationsCreateRequest(ctx context.Context, options *AppliancesClientListOperationsOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.ResourceConnector/operations"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listOperationsHandleResponse handles the ListOperations response.
func (client *AppliancesClient) listOperationsHandleResponse(resp *http.Response) (AppliancesClientListOperationsResponse, error) {
	result := AppliancesClientListOperationsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ApplianceOperationsList); err != nil {
		return AppliancesClientListOperationsResponse{}, err
	}
	return result, nil
}

// Update - Updates an Appliance with the specified Resource Name in the specified Resource Group and Subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-27
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - Appliances name.
//   - parameters - The updatable fields of an existing Appliance.
//   - options - AppliancesClientUpdateOptions contains the optional parameters for the AppliancesClient.Update method.
func (client *AppliancesClient) Update(ctx context.Context, resourceGroupName string, resourceName string, parameters PatchableAppliance, options *AppliancesClientUpdateOptions) (AppliancesClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, parameters, options)
	if err != nil {
		return AppliancesClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AppliancesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AppliancesClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *AppliancesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, parameters PatchableAppliance, options *AppliancesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ResourceConnector/appliances/{resourceName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-27")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *AppliancesClient) updateHandleResponse(resp *http.Response) (AppliancesClientUpdateResponse, error) {
	result := AppliancesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Appliance); err != nil {
		return AppliancesClientUpdateResponse{}, err
	}
	return result, nil
}

