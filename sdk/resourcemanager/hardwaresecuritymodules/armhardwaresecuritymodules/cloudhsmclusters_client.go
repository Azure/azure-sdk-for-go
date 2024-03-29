//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhardwaresecuritymodules

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

// CloudHsmClustersClient contains the methods for the CloudHsmClusters group.
// Don't use this type directly, use NewCloudHsmClustersClient() instead.
type CloudHsmClustersClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewCloudHsmClustersClient creates a new instance of CloudHsmClustersClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCloudHsmClustersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CloudHsmClustersClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CloudHsmClustersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or Update a Cloud HSM Cluster in the specified subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudHsmClusterName - The name of the Cloud HSM Cluster within the specified resource group. Cloud HSM Cluster names must
//     be between 3 and 24 characters in length.
//   - body - Parameters to create Cloud HSM Cluster
//   - options - CloudHsmClustersClientBeginCreateOrUpdateOptions contains the optional parameters for the CloudHsmClustersClient.BeginCreateOrUpdate
//     method.
func (client *CloudHsmClustersClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body CloudHsmCluster, options *CloudHsmClustersClientBeginCreateOrUpdateOptions) (*runtime.Poller[CloudHsmClustersClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, cloudHsmClusterName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CloudHsmClustersClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaOriginalURI,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[CloudHsmClustersClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create or Update a Cloud HSM Cluster in the specified subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
func (client *CloudHsmClustersClient) createOrUpdate(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body CloudHsmCluster, options *CloudHsmClustersClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "CloudHsmClustersClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, cloudHsmClusterName, body, options)
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
func (client *CloudHsmClustersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body CloudHsmCluster, options *CloudHsmClustersClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/{cloudHsmClusterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudHsmClusterName == "" {
		return nil, errors.New("parameter cloudHsmClusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudHsmClusterName}", url.PathEscape(cloudHsmClusterName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes the specified Cloud HSM Cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudHsmClusterName - The name of the Cloud HSM Cluster within the specified resource group. Cloud HSM Cluster names must
//     be between 3 and 24 characters in length.
//   - options - CloudHsmClustersClientBeginDeleteOptions contains the optional parameters for the CloudHsmClustersClient.BeginDelete
//     method.
func (client *CloudHsmClustersClient) BeginDelete(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *CloudHsmClustersClientBeginDeleteOptions) (*runtime.Poller[CloudHsmClustersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, cloudHsmClusterName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CloudHsmClustersClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[CloudHsmClustersClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes the specified Cloud HSM Cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
func (client *CloudHsmClustersClient) deleteOperation(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *CloudHsmClustersClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "CloudHsmClustersClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, cloudHsmClusterName, options)
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
func (client *CloudHsmClustersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *CloudHsmClustersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/{cloudHsmClusterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudHsmClusterName == "" {
		return nil, errors.New("parameter cloudHsmClusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudHsmClusterName}", url.PathEscape(cloudHsmClusterName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified Cloud HSM Cluster
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudHsmClusterName - The name of the Cloud HSM Cluster within the specified resource group. Cloud HSM Cluster names must
//     be between 3 and 24 characters in length.
//   - options - CloudHsmClustersClientGetOptions contains the optional parameters for the CloudHsmClustersClient.Get method.
func (client *CloudHsmClustersClient) Get(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *CloudHsmClustersClientGetOptions) (CloudHsmClustersClientGetResponse, error) {
	var err error
	const operationName = "CloudHsmClustersClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, cloudHsmClusterName, options)
	if err != nil {
		return CloudHsmClustersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CloudHsmClustersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CloudHsmClustersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *CloudHsmClustersClient) getCreateRequest(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *CloudHsmClustersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/{cloudHsmClusterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudHsmClusterName == "" {
		return nil, errors.New("parameter cloudHsmClusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudHsmClusterName}", url.PathEscape(cloudHsmClusterName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CloudHsmClustersClient) getHandleResponse(resp *http.Response) (CloudHsmClustersClientGetResponse, error) {
	result := CloudHsmClustersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CloudHsmCluster); err != nil {
		return CloudHsmClustersClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - The List operation gets information about the Cloud HSM Clusters associated with the subscription
// and within the specified resource group.
//
// Generated from API version 2023-12-10-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - CloudHsmClustersClientListByResourceGroupOptions contains the optional parameters for the CloudHsmClustersClient.NewListByResourceGroupPager
//     method.
func (client *CloudHsmClustersClient) NewListByResourceGroupPager(resourceGroupName string, options *CloudHsmClustersClientListByResourceGroupOptions) *runtime.Pager[CloudHsmClustersClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[CloudHsmClustersClientListByResourceGroupResponse]{
		More: func(page CloudHsmClustersClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CloudHsmClustersClientListByResourceGroupResponse) (CloudHsmClustersClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "CloudHsmClustersClient.NewListByResourceGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			}, nil)
			if err != nil {
				return CloudHsmClustersClientListByResourceGroupResponse{}, err
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *CloudHsmClustersClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *CloudHsmClustersClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters"
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
	if options != nil && options.Skiptoken != nil {
		reqQP.Set("$skiptoken", *options.Skiptoken)
	}
	reqQP.Set("api-version", "2023-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *CloudHsmClustersClient) listByResourceGroupHandleResponse(resp *http.Response) (CloudHsmClustersClientListByResourceGroupResponse, error) {
	result := CloudHsmClustersClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CloudHsmClusterListResult); err != nil {
		return CloudHsmClustersClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - The List operation gets information about the Cloud HSM Clusters associated with the subscription.
//
// Generated from API version 2023-12-10-preview
//   - options - CloudHsmClustersClientListBySubscriptionOptions contains the optional parameters for the CloudHsmClustersClient.NewListBySubscriptionPager
//     method.
func (client *CloudHsmClustersClient) NewListBySubscriptionPager(options *CloudHsmClustersClientListBySubscriptionOptions) *runtime.Pager[CloudHsmClustersClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[CloudHsmClustersClientListBySubscriptionResponse]{
		More: func(page CloudHsmClustersClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CloudHsmClustersClientListBySubscriptionResponse) (CloudHsmClustersClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "CloudHsmClustersClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return CloudHsmClustersClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *CloudHsmClustersClient) listBySubscriptionCreateRequest(ctx context.Context, options *CloudHsmClustersClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Skiptoken != nil {
		reqQP.Set("$skiptoken", *options.Skiptoken)
	}
	reqQP.Set("api-version", "2023-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *CloudHsmClustersClient) listBySubscriptionHandleResponse(resp *http.Response) (CloudHsmClustersClientListBySubscriptionResponse, error) {
	result := CloudHsmClustersClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CloudHsmClusterListResult); err != nil {
		return CloudHsmClustersClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update a Cloud HSM Cluster in the specified subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudHsmClusterName - The name of the Cloud HSM Cluster within the specified resource group. Cloud HSM Cluster names must
//     be between 3 and 24 characters in length.
//   - body - Parameters to create Cloud HSM Cluster
//   - options - CloudHsmClustersClientBeginUpdateOptions contains the optional parameters for the CloudHsmClustersClient.BeginUpdate
//     method.
func (client *CloudHsmClustersClient) BeginUpdate(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body CloudHsmClusterPatchParameters, options *CloudHsmClustersClientBeginUpdateOptions) (*runtime.Poller[CloudHsmClustersClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, cloudHsmClusterName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CloudHsmClustersClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[CloudHsmClustersClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Update a Cloud HSM Cluster in the specified subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
func (client *CloudHsmClustersClient) update(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body CloudHsmClusterPatchParameters, options *CloudHsmClustersClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "CloudHsmClustersClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, cloudHsmClusterName, body, options)
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
func (client *CloudHsmClustersClient) updateCreateRequest(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body CloudHsmClusterPatchParameters, options *CloudHsmClustersClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/{cloudHsmClusterName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudHsmClusterName == "" {
		return nil, errors.New("parameter cloudHsmClusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudHsmClusterName}", url.PathEscape(cloudHsmClusterName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
		return nil, err
	}
	return req, nil
}
