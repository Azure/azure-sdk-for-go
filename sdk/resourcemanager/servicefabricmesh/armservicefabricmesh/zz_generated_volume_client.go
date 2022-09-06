//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armservicefabricmesh

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

// VolumeClient contains the methods for the Volume group.
// Don't use this type directly, use NewVolumeClient() instead.
type VolumeClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewVolumeClient creates a new instance of VolumeClient with the specified values.
// subscriptionID - The customer subscription identifier
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewVolumeClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*VolumeClient, error) {
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
	client := &VolumeClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Creates a volume resource with the specified name, description and properties. If a volume resource with the same
// name exists, then it is updated with the specified description and properties.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-01-preview
// resourceGroupName - Azure resource group name
// volumeResourceName - The identity of the volume.
// volumeResourceDescription - Description for creating a Volume resource.
// options - VolumeClientCreateOptions contains the optional parameters for the VolumeClient.Create method.
func (client *VolumeClient) Create(ctx context.Context, resourceGroupName string, volumeResourceName string, volumeResourceDescription VolumeResourceDescription, options *VolumeClientCreateOptions) (VolumeClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, volumeResourceName, volumeResourceDescription, options)
	if err != nil {
		return VolumeClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VolumeClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return VolumeClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *VolumeClient) createCreateRequest(ctx context.Context, resourceGroupName string, volumeResourceName string, volumeResourceDescription VolumeResourceDescription, options *VolumeClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabricMesh/volumes/{volumeResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{volumeResourceName}", volumeResourceName)
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, volumeResourceDescription)
}

// createHandleResponse handles the Create response.
func (client *VolumeClient) createHandleResponse(resp *http.Response) (VolumeClientCreateResponse, error) {
	result := VolumeClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResourceDescription); err != nil {
		return VolumeClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the volume resource identified by the name.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-01-preview
// resourceGroupName - Azure resource group name
// volumeResourceName - The identity of the volume.
// options - VolumeClientDeleteOptions contains the optional parameters for the VolumeClient.Delete method.
func (client *VolumeClient) Delete(ctx context.Context, resourceGroupName string, volumeResourceName string, options *VolumeClientDeleteOptions) (VolumeClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, volumeResourceName, options)
	if err != nil {
		return VolumeClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VolumeClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return VolumeClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return VolumeClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *VolumeClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, volumeResourceName string, options *VolumeClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabricMesh/volumes/{volumeResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{volumeResourceName}", volumeResourceName)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the information about the volume resource with the given name. The information include the description and other
// properties of the volume.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-01-preview
// resourceGroupName - Azure resource group name
// volumeResourceName - The identity of the volume.
// options - VolumeClientGetOptions contains the optional parameters for the VolumeClient.Get method.
func (client *VolumeClient) Get(ctx context.Context, resourceGroupName string, volumeResourceName string, options *VolumeClientGetOptions) (VolumeClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, volumeResourceName, options)
	if err != nil {
		return VolumeClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return VolumeClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VolumeClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VolumeClient) getCreateRequest(ctx context.Context, resourceGroupName string, volumeResourceName string, options *VolumeClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabricMesh/volumes/{volumeResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{volumeResourceName}", volumeResourceName)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VolumeClient) getHandleResponse(resp *http.Response) (VolumeClientGetResponse, error) {
	result := VolumeClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResourceDescription); err != nil {
		return VolumeClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Gets the information about all volume resources in a given resource group. The information
// include the description and other properties of the Volume.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-01-preview
// resourceGroupName - Azure resource group name
// options - VolumeClientListByResourceGroupOptions contains the optional parameters for the VolumeClient.ListByResourceGroup
// method.
func (client *VolumeClient) NewListByResourceGroupPager(resourceGroupName string, options *VolumeClientListByResourceGroupOptions) *runtime.Pager[VolumeClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[VolumeClientListByResourceGroupResponse]{
		More: func(page VolumeClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VolumeClientListByResourceGroupResponse) (VolumeClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VolumeClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VolumeClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VolumeClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *VolumeClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *VolumeClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabricMesh/volumes"
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
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *VolumeClient) listByResourceGroupHandleResponse(resp *http.Response) (VolumeClientListByResourceGroupResponse, error) {
	result := VolumeClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResourceDescriptionList); err != nil {
		return VolumeClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Gets the information about all volume resources in a given resource group. The information
// include the description and other properties of the volume.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-09-01-preview
// options - VolumeClientListBySubscriptionOptions contains the optional parameters for the VolumeClient.ListBySubscription
// method.
func (client *VolumeClient) NewListBySubscriptionPager(options *VolumeClientListBySubscriptionOptions) *runtime.Pager[VolumeClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[VolumeClientListBySubscriptionResponse]{
		More: func(page VolumeClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VolumeClientListBySubscriptionResponse) (VolumeClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VolumeClientListBySubscriptionResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return VolumeClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VolumeClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *VolumeClient) listBySubscriptionCreateRequest(ctx context.Context, options *VolumeClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ServiceFabricMesh/volumes"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *VolumeClient) listBySubscriptionHandleResponse(resp *http.Response) (VolumeClientListBySubscriptionResponse, error) {
	result := VolumeClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VolumeResourceDescriptionList); err != nil {
		return VolumeClientListBySubscriptionResponse{}, err
	}
	return result, nil
}
