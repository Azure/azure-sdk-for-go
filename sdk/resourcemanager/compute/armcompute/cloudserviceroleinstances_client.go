//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

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

// CloudServiceRoleInstancesClient contains the methods for the CloudServiceRoleInstances group.
// Don't use this type directly, use NewCloudServiceRoleInstancesClient() instead.
type CloudServiceRoleInstancesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewCloudServiceRoleInstancesClient creates a new instance of CloudServiceRoleInstancesClient with the specified values.
//   - subscriptionID - Subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCloudServiceRoleInstancesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CloudServiceRoleInstancesClient, error) {
	cl, err := arm.NewClient(moduleName+".CloudServiceRoleInstancesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CloudServiceRoleInstancesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginDelete - Deletes a role instance from a cloud service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - roleInstanceName - Name of the role instance.
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientBeginDeleteOptions contains the optional parameters for the CloudServiceRoleInstancesClient.BeginDelete
//     method.
func (client *CloudServiceRoleInstancesClient) BeginDelete(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginDeleteOptions) (*runtime.Poller[CloudServiceRoleInstancesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[CloudServiceRoleInstancesClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CloudServiceRoleInstancesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a role instance from a cloud service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
func (client *CloudServiceRoleInstancesClient) deleteOperation(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "CloudServiceRoleInstancesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
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
func (client *CloudServiceRoleInstancesClient) deleteCreateRequest(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances/{roleInstanceName}"
	if roleInstanceName == "" {
		return nil, errors.New("parameter roleInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleInstanceName}", url.PathEscape(roleInstanceName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a role instance from a cloud service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - roleInstanceName - Name of the role instance.
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientGetOptions contains the optional parameters for the CloudServiceRoleInstancesClient.Get
//     method.
func (client *CloudServiceRoleInstancesClient) Get(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientGetOptions) (CloudServiceRoleInstancesClientGetResponse, error) {
	var err error
	const operationName = "CloudServiceRoleInstancesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
	if err != nil {
		return CloudServiceRoleInstancesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CloudServiceRoleInstancesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CloudServiceRoleInstancesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *CloudServiceRoleInstancesClient) getCreateRequest(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances/{roleInstanceName}"
	if roleInstanceName == "" {
		return nil, errors.New("parameter roleInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleInstanceName}", url.PathEscape(roleInstanceName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", string(*options.Expand))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CloudServiceRoleInstancesClient) getHandleResponse(resp *http.Response) (CloudServiceRoleInstancesClientGetResponse, error) {
	result := CloudServiceRoleInstancesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RoleInstance); err != nil {
		return CloudServiceRoleInstancesClientGetResponse{}, err
	}
	return result, nil
}

// GetInstanceView - Retrieves information about the run-time state of a role instance in a cloud service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - roleInstanceName - Name of the role instance.
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientGetInstanceViewOptions contains the optional parameters for the CloudServiceRoleInstancesClient.GetInstanceView
//     method.
func (client *CloudServiceRoleInstancesClient) GetInstanceView(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientGetInstanceViewOptions) (CloudServiceRoleInstancesClientGetInstanceViewResponse, error) {
	var err error
	const operationName = "CloudServiceRoleInstancesClient.GetInstanceView"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getInstanceViewCreateRequest(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
	if err != nil {
		return CloudServiceRoleInstancesClientGetInstanceViewResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CloudServiceRoleInstancesClientGetInstanceViewResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CloudServiceRoleInstancesClientGetInstanceViewResponse{}, err
	}
	resp, err := client.getInstanceViewHandleResponse(httpResp)
	return resp, err
}

// getInstanceViewCreateRequest creates the GetInstanceView request.
func (client *CloudServiceRoleInstancesClient) getInstanceViewCreateRequest(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientGetInstanceViewOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances/{roleInstanceName}/instanceView"
	if roleInstanceName == "" {
		return nil, errors.New("parameter roleInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleInstanceName}", url.PathEscape(roleInstanceName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getInstanceViewHandleResponse handles the GetInstanceView response.
func (client *CloudServiceRoleInstancesClient) getInstanceViewHandleResponse(resp *http.Response) (CloudServiceRoleInstancesClientGetInstanceViewResponse, error) {
	result := CloudServiceRoleInstancesClientGetInstanceViewResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RoleInstanceView); err != nil {
		return CloudServiceRoleInstancesClientGetInstanceViewResponse{}, err
	}
	return result, nil
}

// GetRemoteDesktopFile - Gets a remote desktop file for a role instance in a cloud service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - roleInstanceName - Name of the role instance.
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientGetRemoteDesktopFileOptions contains the optional parameters for the CloudServiceRoleInstancesClient.GetRemoteDesktopFile
//     method.
func (client *CloudServiceRoleInstancesClient) GetRemoteDesktopFile(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientGetRemoteDesktopFileOptions) (CloudServiceRoleInstancesClientGetRemoteDesktopFileResponse, error) {
	var err error
	const operationName = "CloudServiceRoleInstancesClient.GetRemoteDesktopFile"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getRemoteDesktopFileCreateRequest(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
	if err != nil {
		return CloudServiceRoleInstancesClientGetRemoteDesktopFileResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CloudServiceRoleInstancesClientGetRemoteDesktopFileResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CloudServiceRoleInstancesClientGetRemoteDesktopFileResponse{}, err
	}
	return CloudServiceRoleInstancesClientGetRemoteDesktopFileResponse{Body: httpResp.Body}, nil
}

// getRemoteDesktopFileCreateRequest creates the GetRemoteDesktopFile request.
func (client *CloudServiceRoleInstancesClient) getRemoteDesktopFileCreateRequest(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientGetRemoteDesktopFileOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances/{roleInstanceName}/remoteDesktopFile"
	if roleInstanceName == "" {
		return nil, errors.New("parameter roleInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleInstanceName}", url.PathEscape(roleInstanceName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	runtime.SkipBodyDownload(req)
	req.Raw().Header["Accept"] = []string{"application/x-rdp"}
	return req, nil
}

// NewListPager - Gets the list of all role instances in a cloud service. Use nextLink property in the response to get the
// next page of role instances. Do this till nextLink is null to fetch all the role instances.
//
// Generated from API version 2022-09-04
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientListOptions contains the optional parameters for the CloudServiceRoleInstancesClient.NewListPager
//     method.
func (client *CloudServiceRoleInstancesClient) NewListPager(resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientListOptions) *runtime.Pager[CloudServiceRoleInstancesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[CloudServiceRoleInstancesClientListResponse]{
		More: func(page CloudServiceRoleInstancesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CloudServiceRoleInstancesClientListResponse) (CloudServiceRoleInstancesClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "CloudServiceRoleInstancesClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, cloudServiceName, options)
			}, nil)
			if err != nil {
				return CloudServiceRoleInstancesClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *CloudServiceRoleInstancesClient) listCreateRequest(ctx context.Context, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", string(*options.Expand))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *CloudServiceRoleInstancesClient) listHandleResponse(resp *http.Response) (CloudServiceRoleInstancesClientListResponse, error) {
	result := CloudServiceRoleInstancesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RoleInstanceListResult); err != nil {
		return CloudServiceRoleInstancesClientListResponse{}, err
	}
	return result, nil
}

// BeginRebuild - The Rebuild Role Instance asynchronous operation reinstalls the operating system on instances of web roles
// or worker roles and initializes the storage resources that are used by them. If you do not
// want to initialize storage resources, you can use Reimage Role Instance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - roleInstanceName - Name of the role instance.
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientBeginRebuildOptions contains the optional parameters for the CloudServiceRoleInstancesClient.BeginRebuild
//     method.
func (client *CloudServiceRoleInstancesClient) BeginRebuild(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginRebuildOptions) (*runtime.Poller[CloudServiceRoleInstancesClientRebuildResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.rebuild(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[CloudServiceRoleInstancesClientRebuildResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CloudServiceRoleInstancesClientRebuildResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Rebuild - The Rebuild Role Instance asynchronous operation reinstalls the operating system on instances of web roles or
// worker roles and initializes the storage resources that are used by them. If you do not
// want to initialize storage resources, you can use Reimage Role Instance.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
func (client *CloudServiceRoleInstancesClient) rebuild(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginRebuildOptions) (*http.Response, error) {
	var err error
	const operationName = "CloudServiceRoleInstancesClient.BeginRebuild"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.rebuildCreateRequest(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
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

// rebuildCreateRequest creates the Rebuild request.
func (client *CloudServiceRoleInstancesClient) rebuildCreateRequest(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginRebuildOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances/{roleInstanceName}/rebuild"
	if roleInstanceName == "" {
		return nil, errors.New("parameter roleInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleInstanceName}", url.PathEscape(roleInstanceName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginReimage - The Reimage Role Instance asynchronous operation reinstalls the operating system on instances of web roles
// or worker roles.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - roleInstanceName - Name of the role instance.
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientBeginReimageOptions contains the optional parameters for the CloudServiceRoleInstancesClient.BeginReimage
//     method.
func (client *CloudServiceRoleInstancesClient) BeginReimage(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginReimageOptions) (*runtime.Poller[CloudServiceRoleInstancesClientReimageResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.reimage(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[CloudServiceRoleInstancesClientReimageResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CloudServiceRoleInstancesClientReimageResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Reimage - The Reimage Role Instance asynchronous operation reinstalls the operating system on instances of web roles or
// worker roles.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
func (client *CloudServiceRoleInstancesClient) reimage(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginReimageOptions) (*http.Response, error) {
	var err error
	const operationName = "CloudServiceRoleInstancesClient.BeginReimage"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.reimageCreateRequest(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
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

// reimageCreateRequest creates the Reimage request.
func (client *CloudServiceRoleInstancesClient) reimageCreateRequest(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginReimageOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances/{roleInstanceName}/reimage"
	if roleInstanceName == "" {
		return nil, errors.New("parameter roleInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleInstanceName}", url.PathEscape(roleInstanceName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginRestart - The Reboot Role Instance asynchronous operation requests a reboot of a role instance in the cloud service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - roleInstanceName - Name of the role instance.
//   - resourceGroupName - Name of the resource group.
//   - cloudServiceName - Name of the cloud service.
//   - options - CloudServiceRoleInstancesClientBeginRestartOptions contains the optional parameters for the CloudServiceRoleInstancesClient.BeginRestart
//     method.
func (client *CloudServiceRoleInstancesClient) BeginRestart(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginRestartOptions) (*runtime.Poller[CloudServiceRoleInstancesClientRestartResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.restart(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[CloudServiceRoleInstancesClientRestartResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CloudServiceRoleInstancesClientRestartResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Restart - The Reboot Role Instance asynchronous operation requests a reboot of a role instance in the cloud service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
func (client *CloudServiceRoleInstancesClient) restart(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginRestartOptions) (*http.Response, error) {
	var err error
	const operationName = "CloudServiceRoleInstancesClient.BeginRestart"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.restartCreateRequest(ctx, roleInstanceName, resourceGroupName, cloudServiceName, options)
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

// restartCreateRequest creates the Restart request.
func (client *CloudServiceRoleInstancesClient) restartCreateRequest(ctx context.Context, roleInstanceName string, resourceGroupName string, cloudServiceName string, options *CloudServiceRoleInstancesClientBeginRestartOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/cloudServices/{cloudServiceName}/roleInstances/{roleInstanceName}/restart"
	if roleInstanceName == "" {
		return nil, errors.New("parameter roleInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleInstanceName}", url.PathEscape(roleInstanceName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if cloudServiceName == "" {
		return nil, errors.New("parameter cloudServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cloudServiceName}", url.PathEscape(cloudServiceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}
