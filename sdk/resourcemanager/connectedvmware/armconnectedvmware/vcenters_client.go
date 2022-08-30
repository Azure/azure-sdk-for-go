//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by Microsoft (R) AutoRest Code Generator.Changes may cause incorrect behavior and will be lost if the code
// is regenerated.
// DO NOT EDIT.

package armconnectedvmware

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
	"strconv"
	"strings"
)

// VCentersClient contains the methods for the VCenters group.
// Don't use this type directly, use NewVCentersClient() instead.
type VCentersClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewVCentersClient creates a new instance of VCentersClient with the specified values.
// subscriptionID - The Subscription ID.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewVCentersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VCentersClient, error) {
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
	client := &VCentersClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreate - Create Or Update vCenter.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
// resourceGroupName - The Resource Group Name.
// vcenterName - Name of the vCenter.
// body - Request payload.
// options - VCentersClientBeginCreateOptions contains the optional parameters for the VCentersClient.BeginCreate method.
func (client *VCentersClient) BeginCreate(ctx context.Context, resourceGroupName string, vcenterName string, body VCenter, options *VCentersClientBeginCreateOptions) (*runtime.Poller[VCentersClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, vcenterName, body, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[VCentersClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[VCentersClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Create Or Update vCenter.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
func (client *VCentersClient) create(ctx context.Context, resourceGroupName string, vcenterName string, body VCenter, options *VCentersClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, vcenterName, body, options)
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
func (client *VCentersClient) createCreateRequest(ctx context.Context, resourceGroupName string, vcenterName string, body VCenter, options *VCentersClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ConnectedVMwarevSphere/vcenters/{vcenterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vcenterName == "" {
		return nil, errors.New("parameter vcenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vcenterName}", url.PathEscape(vcenterName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// BeginDelete - Implements vCenter DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
// resourceGroupName - The Resource Group Name.
// vcenterName - Name of the vCenter.
// options - VCentersClientBeginDeleteOptions contains the optional parameters for the VCentersClient.BeginDelete method.
func (client *VCentersClient) BeginDelete(ctx context.Context, resourceGroupName string, vcenterName string, options *VCentersClientBeginDeleteOptions) (*runtime.Poller[VCentersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, vcenterName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[VCentersClientDeleteResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[VCentersClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Implements vCenter DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
func (client *VCentersClient) deleteOperation(ctx context.Context, resourceGroupName string, vcenterName string, options *VCentersClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, vcenterName, options)
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
func (client *VCentersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, vcenterName string, options *VCentersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ConnectedVMwarevSphere/vcenters/{vcenterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vcenterName == "" {
		return nil, errors.New("parameter vcenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vcenterName}", url.PathEscape(vcenterName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-10-preview")
	if options != nil && options.Force != nil {
		reqQP.Set("force", strconv.FormatBool(*options.Force))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Implements vCenter GET method.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
// resourceGroupName - The Resource Group Name.
// vcenterName - Name of the vCenter.
// options - VCentersClientGetOptions contains the optional parameters for the VCentersClient.Get method.
func (client *VCentersClient) Get(ctx context.Context, resourceGroupName string, vcenterName string, options *VCentersClientGetOptions) (VCentersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, vcenterName, options)
	if err != nil {
		return VCentersClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VCentersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VCentersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VCentersClient) getCreateRequest(ctx context.Context, resourceGroupName string, vcenterName string, options *VCentersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ConnectedVMwarevSphere/vcenters/{vcenterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vcenterName == "" {
		return nil, errors.New("parameter vcenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vcenterName}", url.PathEscape(vcenterName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VCentersClient) getHandleResponse(resp *http.Response) (VCentersClientGetResponse, error) {
	result := VCentersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VCenter); err != nil {
		return VCentersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List of vCenters in a subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
// options - VCentersClientListOptions contains the optional parameters for the VCentersClient.List method.
func (client *VCentersClient) NewListPager(options *VCentersClientListOptions) *runtime.Pager[VCentersClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[VCentersClientListResponse]{
		More: func(page VCentersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VCentersClientListResponse) (VCentersClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VCentersClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VCentersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VCentersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *VCentersClient) listCreateRequest(ctx context.Context, options *VCentersClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ConnectedVMwarevSphere/vcenters"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *VCentersClient) listHandleResponse(resp *http.Response) (VCentersClientListResponse, error) {
	result := VCentersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VCentersList); err != nil {
		return VCentersClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - List of vCenters in a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
// resourceGroupName - The Resource Group Name.
// options - VCentersClientListByResourceGroupOptions contains the optional parameters for the VCentersClient.ListByResourceGroup
// method.
func (client *VCentersClient) NewListByResourceGroupPager(resourceGroupName string, options *VCentersClientListByResourceGroupOptions) *runtime.Pager[VCentersClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[VCentersClientListByResourceGroupResponse]{
		More: func(page VCentersClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VCentersClientListByResourceGroupResponse) (VCentersClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VCentersClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VCentersClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VCentersClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *VCentersClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *VCentersClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ConnectedVMwarevSphere/vcenters"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *VCentersClient) listByResourceGroupHandleResponse(resp *http.Response) (VCentersClientListByResourceGroupResponse, error) {
	result := VCentersClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VCentersList); err != nil {
		return VCentersClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// Update - API to update certain properties of the vCenter resource.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-01-10-preview
// resourceGroupName - The Resource Group Name.
// vcenterName - Name of the vCenter.
// body - Resource properties to update.
// options - VCentersClientUpdateOptions contains the optional parameters for the VCentersClient.Update method.
func (client *VCentersClient) Update(ctx context.Context, resourceGroupName string, vcenterName string, body ResourcePatch, options *VCentersClientUpdateOptions) (VCentersClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, vcenterName, body, options)
	if err != nil {
		return VCentersClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VCentersClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VCentersClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *VCentersClient) updateCreateRequest(ctx context.Context, resourceGroupName string, vcenterName string, body ResourcePatch, options *VCentersClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ConnectedVMwarevSphere/vcenters/{vcenterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vcenterName == "" {
		return nil, errors.New("parameter vcenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vcenterName}", url.PathEscape(vcenterName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// updateHandleResponse handles the Update response.
func (client *VCentersClient) updateHandleResponse(resp *http.Response) (VCentersClientUpdateResponse, error) {
	result := VCentersClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VCenter); err != nil {
		return VCentersClientUpdateResponse{}, err
	}
	return result, nil
}
