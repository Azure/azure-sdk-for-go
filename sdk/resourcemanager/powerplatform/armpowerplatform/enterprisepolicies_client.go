//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpowerplatform

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

// EnterprisePoliciesClient contains the methods for the EnterprisePolicies group.
// Don't use this type directly, use NewEnterprisePoliciesClient() instead.
type EnterprisePoliciesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewEnterprisePoliciesClient creates a new instance of EnterprisePoliciesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewEnterprisePoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*EnterprisePoliciesClient, error) {
	cl, err := arm.NewClient(moduleName+".EnterprisePoliciesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &EnterprisePoliciesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates an EnterprisePolicy
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-10-30-preview
//   - enterprisePolicyName - Name of the EnterprisePolicy.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - parameters - Parameters supplied to create or update EnterprisePolicy.
//   - options - EnterprisePoliciesClientCreateOrUpdateOptions contains the optional parameters for the EnterprisePoliciesClient.CreateOrUpdate
//     method.
func (client *EnterprisePoliciesClient) CreateOrUpdate(ctx context.Context, enterprisePolicyName string, resourceGroupName string, parameters EnterprisePolicy, options *EnterprisePoliciesClientCreateOrUpdateOptions) (EnterprisePoliciesClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, enterprisePolicyName, resourceGroupName, parameters, options)
	if err != nil {
		return EnterprisePoliciesClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EnterprisePoliciesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return EnterprisePoliciesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *EnterprisePoliciesClient) createOrUpdateCreateRequest(ctx context.Context, enterprisePolicyName string, resourceGroupName string, parameters EnterprisePolicy, options *EnterprisePoliciesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerPlatform/enterprisePolicies/{enterprisePolicyName}"
	if enterprisePolicyName == "" {
		return nil, errors.New("parameter enterprisePolicyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{enterprisePolicyName}", url.PathEscape(enterprisePolicyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-30-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *EnterprisePoliciesClient) createOrUpdateHandleResponse(resp *http.Response) (EnterprisePoliciesClientCreateOrUpdateResponse, error) {
	result := EnterprisePoliciesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnterprisePolicy); err != nil {
		return EnterprisePoliciesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete an EnterprisePolicy
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-10-30-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - enterprisePolicyName - Name of the EnterprisePolicy
//   - options - EnterprisePoliciesClientDeleteOptions contains the optional parameters for the EnterprisePoliciesClient.Delete
//     method.
func (client *EnterprisePoliciesClient) Delete(ctx context.Context, resourceGroupName string, enterprisePolicyName string, options *EnterprisePoliciesClientDeleteOptions) (EnterprisePoliciesClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, enterprisePolicyName, options)
	if err != nil {
		return EnterprisePoliciesClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EnterprisePoliciesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return EnterprisePoliciesClientDeleteResponse{}, err
	}
	return EnterprisePoliciesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *EnterprisePoliciesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, enterprisePolicyName string, options *EnterprisePoliciesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerPlatform/enterprisePolicies/{enterprisePolicyName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if enterprisePolicyName == "" {
		return nil, errors.New("parameter enterprisePolicyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{enterprisePolicyName}", url.PathEscape(enterprisePolicyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-30-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get information about an EnterprisePolicy
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-10-30-preview
//   - enterprisePolicyName - The EnterprisePolicy name.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - EnterprisePoliciesClientGetOptions contains the optional parameters for the EnterprisePoliciesClient.Get method.
func (client *EnterprisePoliciesClient) Get(ctx context.Context, enterprisePolicyName string, resourceGroupName string, options *EnterprisePoliciesClientGetOptions) (EnterprisePoliciesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, enterprisePolicyName, resourceGroupName, options)
	if err != nil {
		return EnterprisePoliciesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EnterprisePoliciesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EnterprisePoliciesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *EnterprisePoliciesClient) getCreateRequest(ctx context.Context, enterprisePolicyName string, resourceGroupName string, options *EnterprisePoliciesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerPlatform/enterprisePolicies/{enterprisePolicyName}"
	if enterprisePolicyName == "" {
		return nil, errors.New("parameter enterprisePolicyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{enterprisePolicyName}", url.PathEscape(enterprisePolicyName))
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
	reqQP.Set("api-version", "2020-10-30-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *EnterprisePoliciesClient) getHandleResponse(resp *http.Response) (EnterprisePoliciesClientGetResponse, error) {
	result := EnterprisePoliciesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnterprisePolicy); err != nil {
		return EnterprisePoliciesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Retrieve a list of EnterprisePolicies within a given resource group
//
// Generated from API version 2020-10-30-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - EnterprisePoliciesClientListByResourceGroupOptions contains the optional parameters for the EnterprisePoliciesClient.NewListByResourceGroupPager
//     method.
func (client *EnterprisePoliciesClient) NewListByResourceGroupPager(resourceGroupName string, options *EnterprisePoliciesClientListByResourceGroupOptions) (*runtime.Pager[EnterprisePoliciesClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[EnterprisePoliciesClientListByResourceGroupResponse]{
		More: func(page EnterprisePoliciesClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EnterprisePoliciesClientListByResourceGroupResponse) (EnterprisePoliciesClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return EnterprisePoliciesClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return EnterprisePoliciesClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return EnterprisePoliciesClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *EnterprisePoliciesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *EnterprisePoliciesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerPlatform/enterprisePolicies"
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
	reqQP.Set("api-version", "2020-10-30-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *EnterprisePoliciesClient) listByResourceGroupHandleResponse(resp *http.Response) (EnterprisePoliciesClientListByResourceGroupResponse, error) {
	result := EnterprisePoliciesClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnterprisePolicyList); err != nil {
		return EnterprisePoliciesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Retrieve a list of EnterprisePolicies within a subscription
//
// Generated from API version 2020-10-30-preview
//   - options - EnterprisePoliciesClientListBySubscriptionOptions contains the optional parameters for the EnterprisePoliciesClient.NewListBySubscriptionPager
//     method.
func (client *EnterprisePoliciesClient) NewListBySubscriptionPager(options *EnterprisePoliciesClientListBySubscriptionOptions) (*runtime.Pager[EnterprisePoliciesClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[EnterprisePoliciesClientListBySubscriptionResponse]{
		More: func(page EnterprisePoliciesClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EnterprisePoliciesClientListBySubscriptionResponse) (EnterprisePoliciesClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return EnterprisePoliciesClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return EnterprisePoliciesClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return EnterprisePoliciesClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *EnterprisePoliciesClient) listBySubscriptionCreateRequest(ctx context.Context, options *EnterprisePoliciesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.PowerPlatform/enterprisePolicies"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-30-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *EnterprisePoliciesClient) listBySubscriptionHandleResponse(resp *http.Response) (EnterprisePoliciesClientListBySubscriptionResponse, error) {
	result := EnterprisePoliciesClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnterprisePolicyList); err != nil {
		return EnterprisePoliciesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Updates an EnterprisePolicy
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-10-30-preview
//   - enterprisePolicyName - Name of the EnterprisePolicy.
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - parameters - Parameters supplied to update EnterprisePolicy.
//   - options - EnterprisePoliciesClientUpdateOptions contains the optional parameters for the EnterprisePoliciesClient.Update
//     method.
func (client *EnterprisePoliciesClient) Update(ctx context.Context, enterprisePolicyName string, resourceGroupName string, parameters PatchEnterprisePolicy, options *EnterprisePoliciesClientUpdateOptions) (EnterprisePoliciesClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, enterprisePolicyName, resourceGroupName, parameters, options)
	if err != nil {
		return EnterprisePoliciesClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EnterprisePoliciesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EnterprisePoliciesClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *EnterprisePoliciesClient) updateCreateRequest(ctx context.Context, enterprisePolicyName string, resourceGroupName string, parameters PatchEnterprisePolicy, options *EnterprisePoliciesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.PowerPlatform/enterprisePolicies/{enterprisePolicyName}"
	if enterprisePolicyName == "" {
		return nil, errors.New("parameter enterprisePolicyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{enterprisePolicyName}", url.PathEscape(enterprisePolicyName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-30-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *EnterprisePoliciesClient) updateHandleResponse(resp *http.Response) (EnterprisePoliciesClientUpdateResponse, error) {
	result := EnterprisePoliciesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnterprisePolicy); err != nil {
		return EnterprisePoliciesClientUpdateResponse{}, err
	}
	return result, nil
}

