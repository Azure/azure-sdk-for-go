//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

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

// VPNSitesClient contains the methods for the VPNSites group.
// Don't use this type directly, use NewVPNSitesClient() instead.
type VPNSitesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewVPNSitesClient creates a new instance of VPNSitesClient with the specified values.
//   - subscriptionID - The subscription credentials which uniquely identify the Microsoft Azure subscription. The subscription
//     ID forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewVPNSitesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VPNSitesClient, error) {
	cl, err := arm.NewClient(moduleName+".VPNSitesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &VPNSitesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates a VpnSite resource if it doesn't exist else updates the existing VpnSite.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The resource group name of the VpnSite.
//   - vpnSiteName - The name of the VpnSite being created or updated.
//   - vpnSiteParameters - Parameters supplied to create or update VpnSite.
//   - options - VPNSitesClientBeginCreateOrUpdateOptions contains the optional parameters for the VPNSitesClient.BeginCreateOrUpdate
//     method.
func (client *VPNSitesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, vpnSiteName string, vpnSiteParameters VPNSite, options *VPNSitesClientBeginCreateOrUpdateOptions) (*runtime.Poller[VPNSitesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, vpnSiteName, vpnSiteParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[VPNSitesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[VPNSitesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates a VpnSite resource if it doesn't exist else updates the existing VpnSite.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
func (client *VPNSitesClient) createOrUpdate(ctx context.Context, resourceGroupName string, vpnSiteName string, vpnSiteParameters VPNSite, options *VPNSitesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "VPNSitesClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, vpnSiteName, vpnSiteParameters, options)
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
func (client *VPNSitesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, vpnSiteName string, vpnSiteParameters VPNSite, options *VPNSitesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnSites/{vpnSiteName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vpnSiteName == "" {
		return nil, errors.New("parameter vpnSiteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vpnSiteName}", url.PathEscape(vpnSiteName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, vpnSiteParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes a VpnSite.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The resource group name of the VpnSite.
//   - vpnSiteName - The name of the VpnSite being deleted.
//   - options - VPNSitesClientBeginDeleteOptions contains the optional parameters for the VPNSitesClient.BeginDelete method.
func (client *VPNSitesClient) BeginDelete(ctx context.Context, resourceGroupName string, vpnSiteName string, options *VPNSitesClientBeginDeleteOptions) (*runtime.Poller[VPNSitesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, vpnSiteName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[VPNSitesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[VPNSitesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a VpnSite.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
func (client *VPNSitesClient) deleteOperation(ctx context.Context, resourceGroupName string, vpnSiteName string, options *VPNSitesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "VPNSitesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, vpnSiteName, options)
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
func (client *VPNSitesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, vpnSiteName string, options *VPNSitesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnSites/{vpnSiteName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vpnSiteName == "" {
		return nil, errors.New("parameter vpnSiteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vpnSiteName}", url.PathEscape(vpnSiteName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Retrieves the details of a VPN site.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The resource group name of the VpnSite.
//   - vpnSiteName - The name of the VpnSite being retrieved.
//   - options - VPNSitesClientGetOptions contains the optional parameters for the VPNSitesClient.Get method.
func (client *VPNSitesClient) Get(ctx context.Context, resourceGroupName string, vpnSiteName string, options *VPNSitesClientGetOptions) (VPNSitesClientGetResponse, error) {
	var err error
	const operationName = "VPNSitesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, vpnSiteName, options)
	if err != nil {
		return VPNSitesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return VPNSitesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return VPNSitesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *VPNSitesClient) getCreateRequest(ctx context.Context, resourceGroupName string, vpnSiteName string, options *VPNSitesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnSites/{vpnSiteName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vpnSiteName == "" {
		return nil, errors.New("parameter vpnSiteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vpnSiteName}", url.PathEscape(vpnSiteName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VPNSitesClient) getHandleResponse(resp *http.Response) (VPNSitesClientGetResponse, error) {
	result := VPNSitesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VPNSite); err != nil {
		return VPNSitesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all the VpnSites in a subscription.
//
// Generated from API version 2023-02-01
//   - options - VPNSitesClientListOptions contains the optional parameters for the VPNSitesClient.NewListPager method.
func (client *VPNSitesClient) NewListPager(options *VPNSitesClientListOptions) *runtime.Pager[VPNSitesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[VPNSitesClientListResponse]{
		More: func(page VPNSitesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VPNSitesClientListResponse) (VPNSitesClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "VPNSitesClient.NewListPager")
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VPNSitesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return VPNSitesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VPNSitesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *VPNSitesClient) listCreateRequest(ctx context.Context, options *VPNSitesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/vpnSites"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *VPNSitesClient) listHandleResponse(resp *http.Response) (VPNSitesClientListResponse, error) {
	result := VPNSitesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListVPNSitesResult); err != nil {
		return VPNSitesClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Lists all the vpnSites in a resource group.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The resource group name of the VpnSite.
//   - options - VPNSitesClientListByResourceGroupOptions contains the optional parameters for the VPNSitesClient.NewListByResourceGroupPager
//     method.
func (client *VPNSitesClient) NewListByResourceGroupPager(resourceGroupName string, options *VPNSitesClientListByResourceGroupOptions) *runtime.Pager[VPNSitesClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[VPNSitesClientListByResourceGroupResponse]{
		More: func(page VPNSitesClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VPNSitesClientListByResourceGroupResponse) (VPNSitesClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "VPNSitesClient.NewListByResourceGroupPager")
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VPNSitesClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return VPNSitesClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VPNSitesClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *VPNSitesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *VPNSitesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnSites"
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
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *VPNSitesClient) listByResourceGroupHandleResponse(resp *http.Response) (VPNSitesClientListByResourceGroupResponse, error) {
	result := VPNSitesClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListVPNSitesResult); err != nil {
		return VPNSitesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// UpdateTags - Updates VpnSite tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The resource group name of the VpnSite.
//   - vpnSiteName - The name of the VpnSite being updated.
//   - vpnSiteParameters - Parameters supplied to update VpnSite tags.
//   - options - VPNSitesClientUpdateTagsOptions contains the optional parameters for the VPNSitesClient.UpdateTags method.
func (client *VPNSitesClient) UpdateTags(ctx context.Context, resourceGroupName string, vpnSiteName string, vpnSiteParameters TagsObject, options *VPNSitesClientUpdateTagsOptions) (VPNSitesClientUpdateTagsResponse, error) {
	var err error
	const operationName = "VPNSitesClient.UpdateTags"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, vpnSiteName, vpnSiteParameters, options)
	if err != nil {
		return VPNSitesClientUpdateTagsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return VPNSitesClientUpdateTagsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return VPNSitesClientUpdateTagsResponse{}, err
	}
	resp, err := client.updateTagsHandleResponse(httpResp)
	return resp, err
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *VPNSitesClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, vpnSiteName string, vpnSiteParameters TagsObject, options *VPNSitesClientUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnSites/{vpnSiteName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vpnSiteName == "" {
		return nil, errors.New("parameter vpnSiteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vpnSiteName}", url.PathEscape(vpnSiteName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, vpnSiteParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client *VPNSitesClient) updateTagsHandleResponse(resp *http.Response) (VPNSitesClientUpdateTagsResponse, error) {
	result := VPNSitesClientUpdateTagsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VPNSite); err != nil {
		return VPNSitesClientUpdateTagsResponse{}, err
	}
	return result, nil
}
