//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armavs

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

// ScriptCmdletsClient contains the methods for the ScriptCmdlets group.
// Don't use this type directly, use NewScriptCmdletsClient() instead.
type ScriptCmdletsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewScriptCmdletsClient creates a new instance of ScriptCmdletsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewScriptCmdletsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ScriptCmdletsClient, error) {
	cl, err := arm.NewClient(moduleName+".ScriptCmdletsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ScriptCmdletsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Return information about a script cmdlet resource in a specific package on a private cloud
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - privateCloudName - Name of the private cloud
//   - scriptPackageName - Name of the script package in the private cloud
//   - scriptCmdletName - Name of the script cmdlet resource in the script package in the private cloud
//   - options - ScriptCmdletsClientGetOptions contains the optional parameters for the ScriptCmdletsClient.Get method.
func (client *ScriptCmdletsClient) Get(ctx context.Context, resourceGroupName string, privateCloudName string, scriptPackageName string, scriptCmdletName string, options *ScriptCmdletsClientGetOptions) (ScriptCmdletsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, privateCloudName, scriptPackageName, scriptCmdletName, options)
	if err != nil {
		return ScriptCmdletsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ScriptCmdletsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ScriptCmdletsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ScriptCmdletsClient) getCreateRequest(ctx context.Context, resourceGroupName string, privateCloudName string, scriptPackageName string, scriptCmdletName string, options *ScriptCmdletsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}/scriptPackages/{scriptPackageName}/scriptCmdlets/{scriptCmdletName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if privateCloudName == "" {
		return nil, errors.New("parameter privateCloudName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateCloudName}", url.PathEscape(privateCloudName))
	if scriptPackageName == "" {
		return nil, errors.New("parameter scriptPackageName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scriptPackageName}", url.PathEscape(scriptPackageName))
	if scriptCmdletName == "" {
		return nil, errors.New("parameter scriptCmdletName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scriptCmdletName}", url.PathEscape(scriptCmdletName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ScriptCmdletsClient) getHandleResponse(resp *http.Response) (ScriptCmdletsClientGetResponse, error) {
	result := ScriptCmdletsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScriptCmdlet); err != nil {
		return ScriptCmdletsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List script cmdlet resources available for a private cloud to create a script execution resource on a private
// cloud
//
// Generated from API version 2022-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - privateCloudName - Name of the private cloud
//   - scriptPackageName - Name of the script package in the private cloud
//   - options - ScriptCmdletsClientListOptions contains the optional parameters for the ScriptCmdletsClient.NewListPager method.
func (client *ScriptCmdletsClient) NewListPager(resourceGroupName string, privateCloudName string, scriptPackageName string, options *ScriptCmdletsClientListOptions) *runtime.Pager[ScriptCmdletsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ScriptCmdletsClientListResponse]{
		More: func(page ScriptCmdletsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ScriptCmdletsClientListResponse) (ScriptCmdletsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, privateCloudName, scriptPackageName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ScriptCmdletsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ScriptCmdletsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ScriptCmdletsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ScriptCmdletsClient) listCreateRequest(ctx context.Context, resourceGroupName string, privateCloudName string, scriptPackageName string, options *ScriptCmdletsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AVS/privateClouds/{privateCloudName}/scriptPackages/{scriptPackageName}/scriptCmdlets"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if privateCloudName == "" {
		return nil, errors.New("parameter privateCloudName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateCloudName}", url.PathEscape(privateCloudName))
	if scriptPackageName == "" {
		return nil, errors.New("parameter scriptPackageName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scriptPackageName}", url.PathEscape(scriptPackageName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ScriptCmdletsClient) listHandleResponse(resp *http.Response) (ScriptCmdletsClientListResponse, error) {
	result := ScriptCmdletsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScriptCmdletsList); err != nil {
		return ScriptCmdletsClientListResponse{}, err
	}
	return result, nil
}
