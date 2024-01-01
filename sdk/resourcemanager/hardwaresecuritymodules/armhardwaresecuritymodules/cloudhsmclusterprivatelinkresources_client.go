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

// CloudHsmClusterPrivateLinkResourcesClient contains the methods for the CloudHsmClusterPrivateLinkResources group.
// Don't use this type directly, use NewCloudHsmClusterPrivateLinkResourcesClient() instead.
type CloudHsmClusterPrivateLinkResourcesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewCloudHsmClusterPrivateLinkResourcesClient creates a new instance of CloudHsmClusterPrivateLinkResourcesClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCloudHsmClusterPrivateLinkResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CloudHsmClusterPrivateLinkResourcesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CloudHsmClusterPrivateLinkResourcesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// ListByCloudHsmCluster - Gets the private link resources supported for the Cloud Hsm Cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-10-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cloudHsmClusterName - The name of the Cloud HSM Cluster within the specified resource group. Cloud HSM Cluster names must
//     be between 3 and 24 characters in length.
//   - options - CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterOptions contains the optional parameters for the
//     CloudHsmClusterPrivateLinkResourcesClient.ListByCloudHsmCluster method.
func (client *CloudHsmClusterPrivateLinkResourcesClient) ListByCloudHsmCluster(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterOptions) (CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse, error) {
	var err error
	const operationName = "CloudHsmClusterPrivateLinkResourcesClient.ListByCloudHsmCluster"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listByCloudHsmClusterCreateRequest(ctx, resourceGroupName, cloudHsmClusterName, options)
	if err != nil {
		return CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse{}, err
	}
	resp, err := client.listByCloudHsmClusterHandleResponse(httpResp)
	return resp, err
}

// listByCloudHsmClusterCreateRequest creates the ListByCloudHsmCluster request.
func (client *CloudHsmClusterPrivateLinkResourcesClient) listByCloudHsmClusterCreateRequest(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/{cloudHsmClusterName}/privateLinkResources"
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

// listByCloudHsmClusterHandleResponse handles the ListByCloudHsmCluster response.
func (client *CloudHsmClusterPrivateLinkResourcesClient) listByCloudHsmClusterHandleResponse(resp *http.Response) (CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse, error) {
	result := CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateLinkResourceListResult); err != nil {
		return CloudHsmClusterPrivateLinkResourcesClientListByCloudHsmClusterResponse{}, err
	}
	return result, nil
}
