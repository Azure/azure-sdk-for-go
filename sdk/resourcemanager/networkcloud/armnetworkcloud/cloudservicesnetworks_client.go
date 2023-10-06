//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetworkcloud

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

// CloudServicesNetworksClient contains the methods for the CloudServicesNetworks group.
// Don't use this type directly, use NewCloudServicesNetworksClient() instead.
type CloudServicesNetworksClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewCloudServicesNetworksClient creates a new instance of CloudServicesNetworksClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCloudServicesNetworksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CloudServicesNetworksClient, error) {
	cl, err := arm.NewClient(moduleName+".CloudServicesNetworksClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CloudServicesNetworksClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a new cloud services network or update the properties of the existing cloud services network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-07-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudServicesNetworkName - The name of the cloud services network.
//   - cloudServicesNetworkParameters - The request body.
//   - options - CloudServicesNetworksClientBeginCreateOrUpdateOptions contains the optional parameters for the CloudServicesNetworksClient.BeginCreateOrUpdate
//     method.
func (client *CloudServicesNetworksClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, cloudServicesNetworkParameters CloudServicesNetwork, options *CloudServicesNetworksClientBeginCreateOrUpdateOptions) (*runtime.Poller[CloudServicesNetworksClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, cloudServicesNetworkName, cloudServicesNetworkParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CloudServicesNetworksClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CloudServicesNetworksClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a new cloud services network or update the properties of the existing cloud services network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-07-01
func (client *CloudServicesNetworksClient) createOrUpdate(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, cloudServicesNetworkParameters CloudServicesNetwork, options *CloudServicesNetworksClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, cloudServicesNetworkName, cloudServicesNetworkParameters, options)
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
func (client *CloudServicesNetworksClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, cloudServicesNetworkParameters CloudServicesNetwork, options *CloudServicesNetworksClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/cloudServicesNetworks/{cloudServicesNetworkName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServicesNetworkName == "" {
		return nil, errors.New("parameter cloudServicesNetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServicesNetworkName}", url.PathEscape(cloudServicesNetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, cloudServicesNetworkParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete the provided cloud services network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-07-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudServicesNetworkName - The name of the cloud services network.
//   - options - CloudServicesNetworksClientBeginDeleteOptions contains the optional parameters for the CloudServicesNetworksClient.BeginDelete
//     method.
func (client *CloudServicesNetworksClient) BeginDelete(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, options *CloudServicesNetworksClientBeginDeleteOptions) (*runtime.Poller[CloudServicesNetworksClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, cloudServicesNetworkName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CloudServicesNetworksClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CloudServicesNetworksClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete the provided cloud services network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-07-01
func (client *CloudServicesNetworksClient) deleteOperation(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, options *CloudServicesNetworksClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, cloudServicesNetworkName, options)
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
func (client *CloudServicesNetworksClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, options *CloudServicesNetworksClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/cloudServicesNetworks/{cloudServicesNetworkName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServicesNetworkName == "" {
		return nil, errors.New("parameter cloudServicesNetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServicesNetworkName}", url.PathEscape(cloudServicesNetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get properties of the provided cloud services network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-07-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudServicesNetworkName - The name of the cloud services network.
//   - options - CloudServicesNetworksClientGetOptions contains the optional parameters for the CloudServicesNetworksClient.Get
//     method.
func (client *CloudServicesNetworksClient) Get(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, options *CloudServicesNetworksClientGetOptions) (CloudServicesNetworksClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, cloudServicesNetworkName, options)
	if err != nil {
		return CloudServicesNetworksClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CloudServicesNetworksClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CloudServicesNetworksClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *CloudServicesNetworksClient) getCreateRequest(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, options *CloudServicesNetworksClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/cloudServicesNetworks/{cloudServicesNetworkName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServicesNetworkName == "" {
		return nil, errors.New("parameter cloudServicesNetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServicesNetworkName}", url.PathEscape(cloudServicesNetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CloudServicesNetworksClient) getHandleResponse(resp *http.Response) (CloudServicesNetworksClientGetResponse, error) {
	result := CloudServicesNetworksClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CloudServicesNetwork); err != nil {
		return CloudServicesNetworksClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Get a list of cloud services networks in the provided resource group.
//
// Generated from API version 2023-07-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - CloudServicesNetworksClientListByResourceGroupOptions contains the optional parameters for the CloudServicesNetworksClient.NewListByResourceGroupPager
//     method.
func (client *CloudServicesNetworksClient) NewListByResourceGroupPager(resourceGroupName string, options *CloudServicesNetworksClientListByResourceGroupOptions) *runtime.Pager[CloudServicesNetworksClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[CloudServicesNetworksClientListByResourceGroupResponse]{
		More: func(page CloudServicesNetworksClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CloudServicesNetworksClientListByResourceGroupResponse) (CloudServicesNetworksClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CloudServicesNetworksClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CloudServicesNetworksClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CloudServicesNetworksClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *CloudServicesNetworksClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *CloudServicesNetworksClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/cloudServicesNetworks"
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
	reqQP.Set("api-version", "2023-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *CloudServicesNetworksClient) listByResourceGroupHandleResponse(resp *http.Response) (CloudServicesNetworksClientListByResourceGroupResponse, error) {
	result := CloudServicesNetworksClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CloudServicesNetworkList); err != nil {
		return CloudServicesNetworksClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Get a list of cloud services networks in the provided subscription.
//
// Generated from API version 2023-07-01
//   - options - CloudServicesNetworksClientListBySubscriptionOptions contains the optional parameters for the CloudServicesNetworksClient.NewListBySubscriptionPager
//     method.
func (client *CloudServicesNetworksClient) NewListBySubscriptionPager(options *CloudServicesNetworksClientListBySubscriptionOptions) *runtime.Pager[CloudServicesNetworksClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[CloudServicesNetworksClientListBySubscriptionResponse]{
		More: func(page CloudServicesNetworksClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CloudServicesNetworksClientListBySubscriptionResponse) (CloudServicesNetworksClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CloudServicesNetworksClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CloudServicesNetworksClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CloudServicesNetworksClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *CloudServicesNetworksClient) listBySubscriptionCreateRequest(ctx context.Context, options *CloudServicesNetworksClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.NetworkCloud/cloudServicesNetworks"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *CloudServicesNetworksClient) listBySubscriptionHandleResponse(resp *http.Response) (CloudServicesNetworksClientListBySubscriptionResponse, error) {
	result := CloudServicesNetworksClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CloudServicesNetworkList); err != nil {
		return CloudServicesNetworksClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update properties of the provided cloud services network, or update the tags associated with it. Properties
// and tag updates can be done independently.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-07-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudServicesNetworkName - The name of the cloud services network.
//   - cloudServicesNetworkUpdateParameters - The request body.
//   - options - CloudServicesNetworksClientBeginUpdateOptions contains the optional parameters for the CloudServicesNetworksClient.BeginUpdate
//     method.
func (client *CloudServicesNetworksClient) BeginUpdate(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, cloudServicesNetworkUpdateParameters CloudServicesNetworkPatchParameters, options *CloudServicesNetworksClientBeginUpdateOptions) (*runtime.Poller[CloudServicesNetworksClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, cloudServicesNetworkName, cloudServicesNetworkUpdateParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CloudServicesNetworksClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CloudServicesNetworksClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Update properties of the provided cloud services network, or update the tags associated with it. Properties and
// tag updates can be done independently.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-07-01
func (client *CloudServicesNetworksClient) update(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, cloudServicesNetworkUpdateParameters CloudServicesNetworkPatchParameters, options *CloudServicesNetworksClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, cloudServicesNetworkName, cloudServicesNetworkUpdateParameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *CloudServicesNetworksClient) updateCreateRequest(ctx context.Context, resourceGroupName string, cloudServicesNetworkName string, cloudServicesNetworkUpdateParameters CloudServicesNetworkPatchParameters, options *CloudServicesNetworksClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/cloudServicesNetworks/{cloudServicesNetworkName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServicesNetworkName == "" {
		return nil, errors.New("parameter cloudServicesNetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServicesNetworkName}", url.PathEscape(cloudServicesNetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, cloudServicesNetworkUpdateParameters); err != nil {
		return nil, err
	}
	return req, nil
}