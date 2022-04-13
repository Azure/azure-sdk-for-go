//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybridcompute

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

// MachineExtensionsClient contains the methods for the MachineExtensions group.
// Don't use this type directly, use NewMachineExtensionsClient() instead.
type MachineExtensionsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewMachineExtensionsClient creates a new instance of MachineExtensionsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewMachineExtensionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*MachineExtensionsClient, error) {
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
	client := &MachineExtensionsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreateOrUpdate - The operation to create or update the extension.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the machine where the extension should be created or updated.
// extensionName - The name of the machine extension.
// extensionParameters - Parameters supplied to the Create Machine Extension operation.
// options - MachineExtensionsClientBeginCreateOrUpdateOptions contains the optional parameters for the MachineExtensionsClient.BeginCreateOrUpdate
// method.
func (client *MachineExtensionsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, machineName string, extensionName string, extensionParameters MachineExtension, options *MachineExtensionsClientBeginCreateOrUpdateOptions) (*armruntime.Poller[MachineExtensionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, machineName, extensionName, extensionParameters, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[MachineExtensionsClientCreateOrUpdateResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[MachineExtensionsClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - The operation to create or update the extension.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *MachineExtensionsClient) createOrUpdate(ctx context.Context, resourceGroupName string, machineName string, extensionName string, extensionParameters MachineExtension, options *MachineExtensionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, machineName, extensionName, extensionParameters, options)
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

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *MachineExtensionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, machineName string, extensionName string, extensionParameters MachineExtension, options *MachineExtensionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/extensions/{extensionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if extensionName == "" {
		return nil, errors.New("parameter extensionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{extensionName}", url.PathEscape(extensionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, extensionParameters)
}

// BeginDelete - The operation to delete the extension.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the machine where the extension should be deleted.
// extensionName - The name of the machine extension.
// options - MachineExtensionsClientBeginDeleteOptions contains the optional parameters for the MachineExtensionsClient.BeginDelete
// method.
func (client *MachineExtensionsClient) BeginDelete(ctx context.Context, resourceGroupName string, machineName string, extensionName string, options *MachineExtensionsClientBeginDeleteOptions) (*armruntime.Poller[MachineExtensionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, machineName, extensionName, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[MachineExtensionsClientDeleteResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[MachineExtensionsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - The operation to delete the extension.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *MachineExtensionsClient) deleteOperation(ctx context.Context, resourceGroupName string, machineName string, extensionName string, options *MachineExtensionsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, machineName, extensionName, options)
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
func (client *MachineExtensionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, machineName string, extensionName string, options *MachineExtensionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/extensions/{extensionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if extensionName == "" {
		return nil, errors.New("parameter extensionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{extensionName}", url.PathEscape(extensionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - The operation to get the extension.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the machine containing the extension.
// extensionName - The name of the machine extension.
// options - MachineExtensionsClientGetOptions contains the optional parameters for the MachineExtensionsClient.Get method.
func (client *MachineExtensionsClient) Get(ctx context.Context, resourceGroupName string, machineName string, extensionName string, options *MachineExtensionsClientGetOptions) (MachineExtensionsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, machineName, extensionName, options)
	if err != nil {
		return MachineExtensionsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return MachineExtensionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return MachineExtensionsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *MachineExtensionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, machineName string, extensionName string, options *MachineExtensionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/extensions/{extensionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if extensionName == "" {
		return nil, errors.New("parameter extensionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{extensionName}", url.PathEscape(extensionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *MachineExtensionsClient) getHandleResponse(resp *http.Response) (MachineExtensionsClientGetResponse, error) {
	result := MachineExtensionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MachineExtension); err != nil {
		return MachineExtensionsClientGetResponse{}, err
	}
	return result, nil
}

// List - The operation to get all extensions of a non-Azure machine
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the machine containing the extension.
// options - MachineExtensionsClientListOptions contains the optional parameters for the MachineExtensionsClient.List method.
func (client *MachineExtensionsClient) List(resourceGroupName string, machineName string, options *MachineExtensionsClientListOptions) *runtime.Pager[MachineExtensionsClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[MachineExtensionsClientListResponse]{
		More: func(page MachineExtensionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *MachineExtensionsClientListResponse) (MachineExtensionsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, machineName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return MachineExtensionsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return MachineExtensionsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return MachineExtensionsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *MachineExtensionsClient) listCreateRequest(ctx context.Context, resourceGroupName string, machineName string, options *MachineExtensionsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/extensions"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	reqQP.Set("api-version", "2021-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *MachineExtensionsClient) listHandleResponse(resp *http.Response) (MachineExtensionsClientListResponse, error) {
	result := MachineExtensionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MachineExtensionsListResult); err != nil {
		return MachineExtensionsClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - The operation to create or update the extension.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// machineName - The name of the machine where the extension should be created or updated.
// extensionName - The name of the machine extension.
// extensionParameters - Parameters supplied to the Create Machine Extension operation.
// options - MachineExtensionsClientBeginUpdateOptions contains the optional parameters for the MachineExtensionsClient.BeginUpdate
// method.
func (client *MachineExtensionsClient) BeginUpdate(ctx context.Context, resourceGroupName string, machineName string, extensionName string, extensionParameters MachineExtensionUpdate, options *MachineExtensionsClientBeginUpdateOptions) (*armruntime.Poller[MachineExtensionsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, machineName, extensionName, extensionParameters, options)
		if err != nil {
			return nil, err
		}
		return armruntime.NewPoller[MachineExtensionsClientUpdateResponse](resp, client.pl, nil)
	} else {
		return armruntime.NewPollerFromResumeToken[MachineExtensionsClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - The operation to create or update the extension.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *MachineExtensionsClient) update(ctx context.Context, resourceGroupName string, machineName string, extensionName string, extensionParameters MachineExtensionUpdate, options *MachineExtensionsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, machineName, extensionName, extensionParameters, options)
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

// updateCreateRequest creates the Update request.
func (client *MachineExtensionsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, machineName string, extensionName string, extensionParameters MachineExtensionUpdate, options *MachineExtensionsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridCompute/machines/{machineName}/extensions/{extensionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if machineName == "" {
		return nil, errors.New("parameter machineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{machineName}", url.PathEscape(machineName))
	if extensionName == "" {
		return nil, errors.New("parameter extensionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{extensionName}", url.PathEscape(extensionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, extensionParameters)
}
