//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdatabricks

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

// VNetPeeringClient contains the methods for the VNetPeering group.
// Don't use this type directly, use NewVNetPeeringClient() instead.
type VNetPeeringClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewVNetPeeringClient creates a new instance of VNetPeeringClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewVNetPeeringClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VNetPeeringClient, error) {
	cl, err := arm.NewClient(moduleName+".VNetPeeringClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &VNetPeeringClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates vNet Peering for workspace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-04-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - peeringName - The name of the workspace vNet peering.
//   - virtualNetworkPeeringParameters - Parameters supplied to the create workspace vNet Peering.
//   - options - VNetPeeringClientBeginCreateOrUpdateOptions contains the optional parameters for the VNetPeeringClient.BeginCreateOrUpdate
//     method.
func (client *VNetPeeringClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, virtualNetworkPeeringParameters VirtualNetworkPeering, options *VNetPeeringClientBeginCreateOrUpdateOptions) (*runtime.Poller[VNetPeeringClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, workspaceName, peeringName, virtualNetworkPeeringParameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[VNetPeeringClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[VNetPeeringClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates vNet Peering for workspace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-04-01-preview
func (client *VNetPeeringClient) createOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, virtualNetworkPeeringParameters VirtualNetworkPeering, options *VNetPeeringClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, peeringName, virtualNetworkPeeringParameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *VNetPeeringClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, virtualNetworkPeeringParameters VirtualNetworkPeering, options *VNetPeeringClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Databricks/workspaces/{workspaceName}/virtualNetworkPeerings/{peeringName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if peeringName == "" {
		return nil, errors.New("parameter peeringName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{peeringName}", url.PathEscape(peeringName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, virtualNetworkPeeringParameters)
}

// BeginDelete - Deletes the workspace vNetPeering.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-04-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - peeringName - The name of the workspace vNet peering.
//   - options - VNetPeeringClientBeginDeleteOptions contains the optional parameters for the VNetPeeringClient.BeginDelete method.
func (client *VNetPeeringClient) BeginDelete(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, options *VNetPeeringClientBeginDeleteOptions) (*runtime.Poller[VNetPeeringClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, workspaceName, peeringName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[VNetPeeringClientDeleteResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[VNetPeeringClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the workspace vNetPeering.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-04-01-preview
func (client *VNetPeeringClient) deleteOperation(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, options *VNetPeeringClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, peeringName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *VNetPeeringClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, options *VNetPeeringClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Databricks/workspaces/{workspaceName}/virtualNetworkPeerings/{peeringName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if peeringName == "" {
		return nil, errors.New("parameter peeringName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{peeringName}", url.PathEscape(peeringName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the workspace vNet Peering.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-04-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - peeringName - The name of the workspace vNet peering.
//   - options - VNetPeeringClientGetOptions contains the optional parameters for the VNetPeeringClient.Get method.
func (client *VNetPeeringClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, options *VNetPeeringClientGetOptions) (VNetPeeringClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, peeringName, options)
	if err != nil {
		return VNetPeeringClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return VNetPeeringClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return VNetPeeringClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VNetPeeringClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, peeringName string, options *VNetPeeringClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Databricks/workspaces/{workspaceName}/virtualNetworkPeerings/{peeringName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if peeringName == "" {
		return nil, errors.New("parameter peeringName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{peeringName}", url.PathEscape(peeringName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VNetPeeringClient) getHandleResponse(resp *http.Response) (VNetPeeringClientGetResponse, error) {
	result := VNetPeeringClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualNetworkPeering); err != nil {
		return VNetPeeringClientGetResponse{}, err
	}
	return result, nil
}

// NewListByWorkspacePager - Lists the workspace vNet Peerings.
//
// Generated from API version 2021-04-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - options - VNetPeeringClientListByWorkspaceOptions contains the optional parameters for the VNetPeeringClient.NewListByWorkspacePager
//     method.
func (client *VNetPeeringClient) NewListByWorkspacePager(resourceGroupName string, workspaceName string, options *VNetPeeringClientListByWorkspaceOptions) *runtime.Pager[VNetPeeringClientListByWorkspaceResponse] {
	return runtime.NewPager(runtime.PagingHandler[VNetPeeringClientListByWorkspaceResponse]{
		More: func(page VNetPeeringClientListByWorkspaceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VNetPeeringClientListByWorkspaceResponse) (VNetPeeringClientListByWorkspaceResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByWorkspaceCreateRequest(ctx, resourceGroupName, workspaceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VNetPeeringClientListByWorkspaceResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return VNetPeeringClientListByWorkspaceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VNetPeeringClientListByWorkspaceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByWorkspaceHandleResponse(resp)
		},
	})
}

// listByWorkspaceCreateRequest creates the ListByWorkspace request.
func (client *VNetPeeringClient) listByWorkspaceCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, options *VNetPeeringClientListByWorkspaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Databricks/workspaces/{workspaceName}/virtualNetworkPeerings"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByWorkspaceHandleResponse handles the ListByWorkspace response.
func (client *VNetPeeringClient) listByWorkspaceHandleResponse(resp *http.Response) (VNetPeeringClientListByWorkspaceResponse, error) {
	result := VNetPeeringClientListByWorkspaceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualNetworkPeeringList); err != nil {
		return VNetPeeringClientListByWorkspaceResponse{}, err
	}
	return result, nil
}
