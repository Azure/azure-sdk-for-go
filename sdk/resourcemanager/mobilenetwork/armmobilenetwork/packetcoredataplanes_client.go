//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmobilenetwork

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

// PacketCoreDataPlanesClient contains the methods for the PacketCoreDataPlanes group.
// Don't use this type directly, use NewPacketCoreDataPlanesClient() instead.
type PacketCoreDataPlanesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewPacketCoreDataPlanesClient creates a new instance of PacketCoreDataPlanesClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPacketCoreDataPlanesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PacketCoreDataPlanesClient, error) {
	cl, err := arm.NewClient(moduleName+".PacketCoreDataPlanesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PacketCoreDataPlanesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a packet core data plane. Must be created in the same location as its parent packet
// core control plane.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - packetCoreControlPlaneName - The name of the packet core control plane.
//   - packetCoreDataPlaneName - The name of the packet core data plane.
//   - parameters - Parameters supplied to the create or update packet core data plane operation.
//   - options - PacketCoreDataPlanesClientBeginCreateOrUpdateOptions contains the optional parameters for the PacketCoreDataPlanesClient.BeginCreateOrUpdate
//     method.
func (client *PacketCoreDataPlanesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, parameters PacketCoreDataPlane, options *PacketCoreDataPlanesClientBeginCreateOrUpdateOptions) (*runtime.Poller[PacketCoreDataPlanesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, packetCoreControlPlaneName, packetCoreDataPlaneName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PacketCoreDataPlanesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[PacketCoreDataPlanesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a packet core data plane. Must be created in the same location as its parent packet
// core control plane.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01
func (client *PacketCoreDataPlanesClient) createOrUpdate(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, parameters PacketCoreDataPlane, options *PacketCoreDataPlanesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, packetCoreControlPlaneName, packetCoreDataPlaneName, parameters, options)
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
func (client *PacketCoreDataPlanesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, parameters PacketCoreDataPlane, options *PacketCoreDataPlanesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{packetCoreControlPlaneName}/packetCoreDataPlanes/{packetCoreDataPlaneName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if packetCoreControlPlaneName == "" {
		return nil, errors.New("parameter packetCoreControlPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreControlPlaneName}", url.PathEscape(packetCoreControlPlaneName))
	if packetCoreDataPlaneName == "" {
		return nil, errors.New("parameter packetCoreDataPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreDataPlaneName}", url.PathEscape(packetCoreDataPlaneName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes the specified packet core data plane.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - packetCoreControlPlaneName - The name of the packet core control plane.
//   - packetCoreDataPlaneName - The name of the packet core data plane.
//   - options - PacketCoreDataPlanesClientBeginDeleteOptions contains the optional parameters for the PacketCoreDataPlanesClient.BeginDelete
//     method.
func (client *PacketCoreDataPlanesClient) BeginDelete(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, options *PacketCoreDataPlanesClientBeginDeleteOptions) (*runtime.Poller[PacketCoreDataPlanesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, packetCoreControlPlaneName, packetCoreDataPlaneName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PacketCoreDataPlanesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[PacketCoreDataPlanesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the specified packet core data plane.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01
func (client *PacketCoreDataPlanesClient) deleteOperation(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, options *PacketCoreDataPlanesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, packetCoreControlPlaneName, packetCoreDataPlaneName, options)
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
func (client *PacketCoreDataPlanesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, options *PacketCoreDataPlanesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{packetCoreControlPlaneName}/packetCoreDataPlanes/{packetCoreDataPlaneName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if packetCoreControlPlaneName == "" {
		return nil, errors.New("parameter packetCoreControlPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreControlPlaneName}", url.PathEscape(packetCoreControlPlaneName))
	if packetCoreDataPlaneName == "" {
		return nil, errors.New("parameter packetCoreDataPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreDataPlaneName}", url.PathEscape(packetCoreDataPlaneName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets information about the specified packet core data plane.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - packetCoreControlPlaneName - The name of the packet core control plane.
//   - packetCoreDataPlaneName - The name of the packet core data plane.
//   - options - PacketCoreDataPlanesClientGetOptions contains the optional parameters for the PacketCoreDataPlanesClient.Get
//     method.
func (client *PacketCoreDataPlanesClient) Get(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, options *PacketCoreDataPlanesClientGetOptions) (PacketCoreDataPlanesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, packetCoreControlPlaneName, packetCoreDataPlaneName, options)
	if err != nil {
		return PacketCoreDataPlanesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PacketCoreDataPlanesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PacketCoreDataPlanesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PacketCoreDataPlanesClient) getCreateRequest(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, options *PacketCoreDataPlanesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{packetCoreControlPlaneName}/packetCoreDataPlanes/{packetCoreDataPlaneName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if packetCoreControlPlaneName == "" {
		return nil, errors.New("parameter packetCoreControlPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreControlPlaneName}", url.PathEscape(packetCoreControlPlaneName))
	if packetCoreDataPlaneName == "" {
		return nil, errors.New("parameter packetCoreDataPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreDataPlaneName}", url.PathEscape(packetCoreDataPlaneName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PacketCoreDataPlanesClient) getHandleResponse(resp *http.Response) (PacketCoreDataPlanesClientGetResponse, error) {
	result := PacketCoreDataPlanesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PacketCoreDataPlane); err != nil {
		return PacketCoreDataPlanesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByPacketCoreControlPlanePager - Lists all the packet core data planes associated with a packet core control plane.
//
// Generated from API version 2023-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - packetCoreControlPlaneName - The name of the packet core control plane.
//   - options - PacketCoreDataPlanesClientListByPacketCoreControlPlaneOptions contains the optional parameters for the PacketCoreDataPlanesClient.NewListByPacketCoreControlPlanePager
//     method.
func (client *PacketCoreDataPlanesClient) NewListByPacketCoreControlPlanePager(resourceGroupName string, packetCoreControlPlaneName string, options *PacketCoreDataPlanesClientListByPacketCoreControlPlaneOptions) (*runtime.Pager[PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse]{
		More: func(page PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse) (PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByPacketCoreControlPlaneCreateRequest(ctx, resourceGroupName, packetCoreControlPlaneName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByPacketCoreControlPlaneHandleResponse(resp)
		},
	})
}

// listByPacketCoreControlPlaneCreateRequest creates the ListByPacketCoreControlPlane request.
func (client *PacketCoreDataPlanesClient) listByPacketCoreControlPlaneCreateRequest(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, options *PacketCoreDataPlanesClientListByPacketCoreControlPlaneOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{packetCoreControlPlaneName}/packetCoreDataPlanes"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if packetCoreControlPlaneName == "" {
		return nil, errors.New("parameter packetCoreControlPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreControlPlaneName}", url.PathEscape(packetCoreControlPlaneName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByPacketCoreControlPlaneHandleResponse handles the ListByPacketCoreControlPlane response.
func (client *PacketCoreDataPlanesClient) listByPacketCoreControlPlaneHandleResponse(resp *http.Response) (PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse, error) {
	result := PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PacketCoreDataPlaneListResult); err != nil {
		return PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse{}, err
	}
	return result, nil
}

// UpdateTags - Updates packet core data planes tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - packetCoreControlPlaneName - The name of the packet core control plane.
//   - packetCoreDataPlaneName - The name of the packet core data plane.
//   - parameters - Parameters supplied to update packet core data plane tags.
//   - options - PacketCoreDataPlanesClientUpdateTagsOptions contains the optional parameters for the PacketCoreDataPlanesClient.UpdateTags
//     method.
func (client *PacketCoreDataPlanesClient) UpdateTags(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, parameters TagsObject, options *PacketCoreDataPlanesClientUpdateTagsOptions) (PacketCoreDataPlanesClientUpdateTagsResponse, error) {
	var err error
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, packetCoreControlPlaneName, packetCoreDataPlaneName, parameters, options)
	if err != nil {
		return PacketCoreDataPlanesClientUpdateTagsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PacketCoreDataPlanesClientUpdateTagsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PacketCoreDataPlanesClientUpdateTagsResponse{}, err
	}
	resp, err := client.updateTagsHandleResponse(httpResp)
	return resp, err
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *PacketCoreDataPlanesClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, packetCoreDataPlaneName string, parameters TagsObject, options *PacketCoreDataPlanesClientUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/{packetCoreControlPlaneName}/packetCoreDataPlanes/{packetCoreDataPlaneName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if packetCoreControlPlaneName == "" {
		return nil, errors.New("parameter packetCoreControlPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreControlPlaneName}", url.PathEscape(packetCoreControlPlaneName))
	if packetCoreDataPlaneName == "" {
		return nil, errors.New("parameter packetCoreDataPlaneName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{packetCoreDataPlaneName}", url.PathEscape(packetCoreDataPlaneName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client *PacketCoreDataPlanesClient) updateTagsHandleResponse(resp *http.Response) (PacketCoreDataPlanesClientUpdateTagsResponse, error) {
	result := PacketCoreDataPlanesClientUpdateTagsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PacketCoreDataPlane); err != nil {
		return PacketCoreDataPlanesClientUpdateTagsResponse{}, err
	}
	return result, nil
}

