//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhanaonazure

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

// SapMonitorsClient contains the methods for the SapMonitors group.
// Don't use this type directly, use NewSapMonitorsClient() instead.
type SapMonitorsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewSapMonitorsClient creates a new instance of SapMonitorsClient with the specified values.
// subscriptionID - Subscription ID which uniquely identify Microsoft Azure subscription. The subscription ID forms part of
// the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewSapMonitorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SapMonitorsClient, error) {
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
	client := &SapMonitorsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreate - Creates a SAP monitor for the specified subscription, resource group, and resource name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-02-07-preview
// resourceGroupName - Name of the resource group.
// sapMonitorName - Name of the SAP monitor resource.
// sapMonitorParameter - Request body representing a SAP Monitor
// options - SapMonitorsClientBeginCreateOptions contains the optional parameters for the SapMonitorsClient.BeginCreate method.
func (client *SapMonitorsClient) BeginCreate(ctx context.Context, resourceGroupName string, sapMonitorName string, sapMonitorParameter SapMonitor, options *SapMonitorsClientBeginCreateOptions) (*runtime.Poller[SapMonitorsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, sapMonitorName, sapMonitorParameter, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[SapMonitorsClientCreateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[SapMonitorsClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Creates a SAP monitor for the specified subscription, resource group, and resource name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-02-07-preview
func (client *SapMonitorsClient) create(ctx context.Context, resourceGroupName string, sapMonitorName string, sapMonitorParameter SapMonitor, options *SapMonitorsClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, sapMonitorName, sapMonitorParameter, options)
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

// createCreateRequest creates the Create request.
func (client *SapMonitorsClient) createCreateRequest(ctx context.Context, resourceGroupName string, sapMonitorName string, sapMonitorParameter SapMonitor, options *SapMonitorsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HanaOnAzure/sapMonitors/{sapMonitorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapMonitorName == "" {
		return nil, errors.New("parameter sapMonitorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapMonitorName}", url.PathEscape(sapMonitorName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-02-07-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, sapMonitorParameter)
}

// BeginDelete - Deletes a SAP monitor with the specified subscription, resource group, and monitor name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-02-07-preview
// resourceGroupName - Name of the resource group.
// sapMonitorName - Name of the SAP monitor resource.
// options - SapMonitorsClientBeginDeleteOptions contains the optional parameters for the SapMonitorsClient.BeginDelete method.
func (client *SapMonitorsClient) BeginDelete(ctx context.Context, resourceGroupName string, sapMonitorName string, options *SapMonitorsClientBeginDeleteOptions) (*runtime.Poller[SapMonitorsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, sapMonitorName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[SapMonitorsClientDeleteResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[SapMonitorsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Deletes a SAP monitor with the specified subscription, resource group, and monitor name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-02-07-preview
func (client *SapMonitorsClient) deleteOperation(ctx context.Context, resourceGroupName string, sapMonitorName string, options *SapMonitorsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, sapMonitorName, options)
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
func (client *SapMonitorsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, sapMonitorName string, options *SapMonitorsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HanaOnAzure/sapMonitors/{sapMonitorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapMonitorName == "" {
		return nil, errors.New("parameter sapMonitorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapMonitorName}", url.PathEscape(sapMonitorName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-02-07-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets properties of a SAP monitor for the specified subscription, resource group, and resource name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-02-07-preview
// resourceGroupName - Name of the resource group.
// sapMonitorName - Name of the SAP monitor resource.
// options - SapMonitorsClientGetOptions contains the optional parameters for the SapMonitorsClient.Get method.
func (client *SapMonitorsClient) Get(ctx context.Context, resourceGroupName string, sapMonitorName string, options *SapMonitorsClientGetOptions) (SapMonitorsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, sapMonitorName, options)
	if err != nil {
		return SapMonitorsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SapMonitorsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SapMonitorsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SapMonitorsClient) getCreateRequest(ctx context.Context, resourceGroupName string, sapMonitorName string, options *SapMonitorsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HanaOnAzure/sapMonitors/{sapMonitorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapMonitorName == "" {
		return nil, errors.New("parameter sapMonitorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapMonitorName}", url.PathEscape(sapMonitorName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-02-07-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SapMonitorsClient) getHandleResponse(resp *http.Response) (SapMonitorsClientGetResponse, error) {
	result := SapMonitorsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SapMonitor); err != nil {
		return SapMonitorsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of SAP monitors in the specified subscription. The operations returns various properties of
// each SAP monitor.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-02-07-preview
// options - SapMonitorsClientListOptions contains the optional parameters for the SapMonitorsClient.List method.
func (client *SapMonitorsClient) NewListPager(options *SapMonitorsClientListOptions) *runtime.Pager[SapMonitorsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[SapMonitorsClientListResponse]{
		More: func(page SapMonitorsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SapMonitorsClientListResponse) (SapMonitorsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SapMonitorsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return SapMonitorsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SapMonitorsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *SapMonitorsClient) listCreateRequest(ctx context.Context, options *SapMonitorsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HanaOnAzure/sapMonitors"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-02-07-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SapMonitorsClient) listHandleResponse(resp *http.Response) (SapMonitorsClientListResponse, error) {
	result := SapMonitorsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SapMonitorListResult); err != nil {
		return SapMonitorsClientListResponse{}, err
	}
	return result, nil
}

// Update - Patches the Tags field of a SAP monitor for the specified subscription, resource group, and monitor name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-02-07-preview
// resourceGroupName - Name of the resource group.
// sapMonitorName - Name of the SAP monitor resource.
// tagsParameter - Request body that only contains the new Tags field
// options - SapMonitorsClientUpdateOptions contains the optional parameters for the SapMonitorsClient.Update method.
func (client *SapMonitorsClient) Update(ctx context.Context, resourceGroupName string, sapMonitorName string, tagsParameter Tags, options *SapMonitorsClientUpdateOptions) (SapMonitorsClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, sapMonitorName, tagsParameter, options)
	if err != nil {
		return SapMonitorsClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SapMonitorsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SapMonitorsClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *SapMonitorsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, sapMonitorName string, tagsParameter Tags, options *SapMonitorsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HanaOnAzure/sapMonitors/{sapMonitorName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapMonitorName == "" {
		return nil, errors.New("parameter sapMonitorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapMonitorName}", url.PathEscape(sapMonitorName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-02-07-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, tagsParameter)
}

// updateHandleResponse handles the Update response.
func (client *SapMonitorsClient) updateHandleResponse(resp *http.Response) (SapMonitorsClientUpdateResponse, error) {
	result := SapMonitorsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SapMonitor); err != nil {
		return SapMonitorsClientUpdateResponse{}, err
	}
	return result, nil
}
